package onstepx

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type ResponseType int

const (
	RspNone ResponseType = iota
	RspBool
	RspHash
)

type command struct {
	Request string
	RspChan chan string
	ErrChan chan error
	RspType ResponseType
	Timeout time.Duration
}

func (d *onstepxDevice) sendCommand(request string, rspType ResponseType, timeout time.Duration) (string, error) {

	// Must be connected
	if !d.isConnected {
		return "", ErrNotConnected
	}

	// If no timeout specified use default of 2 seconds
	if timeout == 0 {
		timeout = 2 * time.Second
	}

	// Create command and send it for processing
	rspChan := make(chan string)
	errChan := make(chan error)
	cmd := command{
		Request: request,
		RspChan: rspChan,
		ErrChan: errChan,
		RspType: rspType,
		Timeout: timeout,
	}
	d.cmdChan <- cmd

	// Wait for response, error or timeout
	select {
	case response := <-rspChan:
		return response, nil
	case err := <-errChan:
		return "", err
	case <-time.After(timeout):
		return "", errors.New("command timed out waiting for response")
	}
}

// Check available serial ports to find OnStepX device
func FindPort(baudRate int) (string, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		slog.Warn("enumerator.GetDetailedPortsList returned an error", "error", err)
	}
	if len(ports) == 0 {
		return "", errors.New("no serial ports found on the system")
	}

	slog.Info("probing for OnStepX device...", "Ports", len(ports))
	for _, port := range ports {
		slog.Debug("checking port", "name", port.Name, "isUSB", port.IsUSB, "VID", port.VID, "PID", port.PID)
		if port.IsUSB {
			slog.Info("probing port", "name", port.Name, "baudRate", baudRate)
			if probePortWithTimeout(port.Name, baudRate, 4*time.Second) {
				return port.Name, nil
			}
		} else {
			slog.Debug("skipping non-USB port", "name", port.Name)
		}
	}
	return "", errors.New("could not find OnStepX device on any USB serial port")
}

// probePortWithTimeout probes a port with a hard timeout that guarantees cleanup.
// Uses a goroutine for the actual probe, but closes the port if timeout occurs.
func probePortWithTimeout(portName string, baudRate int, timeout time.Duration) bool {
	resultChan := make(chan bool, 1)

	// Shared variable for port handle - allows cleanup on timeout
	var probePort serial.Port
	var probeMutex sync.Mutex

	go func() {
		mode := &serial.Mode{BaudRate: baudRate}
		p, err := serial.Open(portName, mode)
		if err != nil {
			slog.Warn("could not open port to probe", "name", portName, "error", err)
			resultChan <- false
			return
		}

		// Store port handle for potential cleanup
		probeMutex.Lock()
		probePort = p
		probeMutex.Unlock()

		// Set read timeout
		p.SetReadTimeout(2 * time.Second)

		_, err = p.Write([]byte(":GVP#"))
		if err != nil {
			slog.Debug("Port write failed", "name", portName, "error", err)
			p.Close()
			resultChan <- false
			return
		}

		reader := bufio.NewReader(p)
		line, err := reader.ReadString('#')
		p.Close() // Close immediately after read

		// Clear the shared handle since we closed it
		probeMutex.Lock()
		probePort = nil
		probeMutex.Unlock()

		if err != nil {
			slog.Debug("Port read failed or timed out", "name", portName, "error", err)
			resultChan <- false
			return
		}

		if line == "On-Step#" {
			slog.Info("Successfully probed port", "name", portName)
			resultChan <- true
			return
		}

		slog.Debug("Unexpected response", "name", portName, "response", line)
		resultChan <- false
	}()

	// Wait for result with hard timeout
	select {
	case success := <-resultChan:
		return success
	case <-time.After(timeout):
		slog.Warn("Probe timed out. Forcing cleanup.", "name", portName, "timeout", timeout)

		// Force close the port if goroutine is still holding it
		probeMutex.Lock()
		if probePort != nil {
			probePort.Close()
			probePort = nil
		}
		probeMutex.Unlock()

		return false
	}
}

func processCommands(port serial.Port, commands chan command) {

	var err error

	// Process commands until command channel is closed
	slog.Info("command processor started")
	for {

		// Get command
		cmd, ok := <-commands
		if !ok {
			// Command channel closed so terminate
			slog.Info("Command channel closed")
			break
		}

		// Discard any characters in read buffer
		port.ResetInputBuffer()

		// Send request to OnStepX device
		_, err = port.Write([]byte(":" + cmd.Request + "#"))
		if err != nil {
			slog.Error("Serial write failed", "error", err)
			cmd.ErrChan <- fmt.Errorf("failed to write to serial port: %w", err)
			continue
		}
		slog.Debug("request sent to OnStepX device: " + cmd.Request)

		// Get response (if any)
		var response []byte
		err = nil
		if cmd.RspType == RspBool {

			// Read a single character
			port.SetReadTimeout(cmd.Timeout)
			readBuf := make([]byte, 1)
			n, err := port.Read(readBuf)
			if err != nil || n != 1 {
				err = fmt.Errorf("serial read failed: %w", err)
				break
			}
			response = append(response, readBuf[0])

		} else if cmd.RspType == RspHash {

			// Read charcaters one at a time until '#' received
			port.SetReadTimeout(100 * time.Millisecond)
			readBuf := make([]byte, 1)
			start := time.Now()
			for {
				if time.Since(start) > cmd.Timeout {
					err = errors.New("read timeout")
					break
				}
				n, err := port.Read(readBuf)
				if err != nil {
					err = fmt.Errorf("serial read failed: %w", err)
					break
				}
				if n > 0 {
					b := readBuf[0]
					if b == '#' {
						break
					}
					response = append(response, b)
				}
			}
		}

		// Return error or response
		if err != nil {
			cmd.ErrChan <- err
		} else {
			rspString := string(response)
			if len(rspString) > 0 {
				slog.Debug("response received from OnStepX device: " + rspString)
			}
			cmd.RspChan <- rspString
		}
	}
}
