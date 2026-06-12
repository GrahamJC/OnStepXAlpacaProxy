package connection

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"go.bug.st/serial"
)

const (
	SerialClosed     = 1
	SerialConnecting = 2
	SerialOpen       = 3
	SerialError      = 4
)

type SerialStatus int

type SerialTransport struct {
	portName string
	baudRate int
	timeout  time.Duration
	status   SerialStatus
	lock     sync.Mutex
	port     serial.Port
}

func NewSerialConnection(portName string, baudRate int, timeoutMs int) *SerialTransport {
	os := SerialTransport{}
	os.portName = portName
	os.baudRate = baudRate
	os.timeout = time.Duration(timeoutMs) * time.Millisecond
	os.status = SerialClosed
	return &os
}

// Transport interface
func (os *SerialTransport) IsConnected() bool {
	return os.status == SerialOpen
}

func (os *SerialTransport) Open() error {

	// Ignore if already connected
	if os.status == SerialOpen {
		return nil
	}

	// Reject if not closed
	if os.status != SerialClosed {
		return errors.New("Open called when not closed")
	}

	// Check that we have a COM Port and baud rate
	if os.portName == "" {
		return errors.New("Open with COM port not set")
	}
	if os.baudRate == 0 {
		return errors.New("Connect with baud rate not set")
	}

	// Lock connection
	os.lock.Lock()
	defer os.lock.Unlock()

	// Open the port
	var err error
	mode := &serial.Mode{BaudRate: os.baudRate}
	os.port, err = serial.Open(os.portName, mode)
	if err != nil {
		return fmt.Errorf("Error opening serial port: %w", err)
	}
	os.status = SerialOpen

	// Done
	return nil
}

func (os *SerialTransport) Close() error {

	// Ignore if already closed
	if os.status == SerialClosed {
		return nil
	}

	// Reject if not open
	if os.status != SerialOpen {
		return errors.New("Close called when not open")
	}

	// Lock connection
	os.lock.Lock()
	defer os.lock.Unlock()

	// Close the port
	if err := os.port.Close(); err != nil {
		return fmt.Errorf("Error closing serial port: %w", err)
	}
	os.status = SerialClosed

	// Done
	return nil
}

func (os *SerialTransport) Send(cmd string) error {

	// Check that we're connected
	if os.status != SerialOpen {
		return errors.New("Send called when not connected")
	}

	// Lock connection
	os.lock.Lock()
	defer os.lock.Unlock()

	// Send the command
	fmt.Printf(":%s#\n", cmd)
	_, err := os.port.Write([]byte(":" + cmd + "#"))
	if err != nil {
		return fmt.Errorf("Send: %w", err)
	}

	// Done
	return nil
}

func (os *SerialTransport) SendReceive(cmd string) (string, error) {

	// Check that we're connected
	if os.status != SerialOpen {
		return "", errors.New("SendReceive called when not connected")
	}

	// Lock connection
	os.lock.Lock()
	defer os.lock.Unlock()

	// Send the command
	_, err := os.port.Write([]byte(":" + cmd + "#"))
	if err != nil {
		return "", fmt.Errorf("SendReceive: %w", err)
	}

	// Read the response until a '#' character is received or timeout occurs (there is also
	// some special handling for commands that return 0 or 1 with no terminating '#)
	var response []byte
	os.port.SetReadTimeout(10 * time.Millisecond)
	readBuf := make([]byte, 1)
	deadline := time.Now().Add(os.timeout)
	for {
		if time.Now().After(deadline) {
			return "", errors.New("read timeout")
		}
		n, err := os.port.Read(readBuf)
		if err != nil {
			return "", fmt.Errorf("SendReceive: %w", err)
		}
		if n > 0 {
			b := readBuf[0]
			if b == '#' {
				break
			}
			response = append(response, b)
		} else if len(response) > 0 && (string(response) == "0" || string(response) == "1") {
			// Special handling for 0/1 responses with no terminating '#'
			break
		}
	}

	// Return result
	fmt.Printf(":%s# => %s\n", cmd, string(response))
	return string(response), nil
}
