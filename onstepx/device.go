package onstepx

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"sync"
	"time"

	"go.bug.st/serial"
)

var (
	ErrNotConnected = errors.New("Device not connected")
)

type OnStepXDevice interface {
	Connect() bool
	Disconnect()
	IsConnected() bool
	ComPort() string
	BaudRate() int
	GetFirmwareVersion() (string, error)
	GetSiteTime() (time.Time, error)
	SetSiteTime(time.Time) error
	GetSiteLatitude() (float32, error)
	SetSiteLatitude(float32) error
	GetSiteLongitude() (float32, error)
	SetSiteLongitude(float32) error
	GetSiteElevation() (float32, error)
	SetSiteElevation(float32) error
	GetSiteUTCOffset() (int, error)
	SetSiteUTCOffset(int) error
	GetPositionRA() (float32, error)
	GetPositionDec() (float32, error)
	GetPositionAlt() (float32, error)
	GetPositionAz() (float32, error)
}

type onstepxDevice struct {
	comPort     string
	baudRate    int
	isConnected bool
	port        serial.Port
	cmdChan     chan command
	mutex       sync.Mutex
}

func NewDevice(comPort string, baudRate int) OnStepXDevice {

	// If baudRate not specified use 9600
	if baudRate == 0 {
		baudRate = 9600
	}

	// Create and initialize a new device
	onstepx := &onstepxDevice{}
	onstepx.comPort = comPort
	onstepx.baudRate = baudRate

	// Return device
	return onstepx
}

func (od *onstepxDevice) Connect() bool {

	var err error

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Ignore if already connected
	if od.isConnected {
		slog.Debug("Connect called when already connected")
		return false
	}

	// Check that we have a COM Port and open it
	if od.comPort == "" {
		slog.Warn("Connect called with no COM port set")
		return false
	}
	mode := &serial.Mode{BaudRate: od.baudRate}
	od.port, err = serial.Open(od.comPort, mode)
	if err != nil {
		slog.Error("Connect failed to open COM port", "comPort", od.comPort, "error", err)
		return false
	}
	od.isConnected = true

	// Create request/response channels and start goroutine to process commands
	od.cmdChan = make(chan command)
	go processCommands(od.port, od.cmdChan)

	// Done
	return true
}

func (od *onstepxDevice) Disconnect() {

	var err error

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Ignore if not connected
	if !od.isConnected {
		slog.Debug("Disconnect() called when not connected")
		return
	}

	// Close command channel
	close(od.cmdChan)

	// Close port
	err = od.port.Close()
	if err != nil {
		slog.Error("Disconnect failed to close port", "comPort", od.comPort)
	}
	od.port = nil

	// Done
	od.isConnected = false
}

func (od *onstepxDevice) IsConnected() bool {

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Return connection status
	return od.isConnected
}

func (od *onstepxDevice) ComPort() string {
	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Return COM port
	return od.comPort
}

func (od *onstepxDevice) BaudRate() int {
	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Return baud rate
	return od.baudRate
}

func (od *onstepxDevice) GetFirmwareVersion() (string, error) {

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Get OnStepX version
	version, err := od.sendCommand("GVN", RspHash, 0)
	if err != nil {
		return "", fmt.Errorf("GetFirmwareVersion: %w", err)
	}
	return version, nil
}

func (od *onstepxDevice) GetSiteTime() (time.Time, error) {

	var err error
	var result time.Time

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("GetSiteTime called when not connected")
		return result, ErrNotConnected
	}

	// Get date and time
	rspDate, err := od.sendCommand("GC", RspHash, 0)
	if err != nil {
		slog.Debug("GetSiteTime request GC failed", "error", err)
		return result, err
	}
	rspTime, err := od.sendCommand("GL", RspHash, 0)
	if err != nil {
		slog.Debug("GetSiteTime request GL failed", "error", err)
		return result, err
	}

	// Parse result
	result, err = time.Parse("01/02/06 03:04:05", rspDate+" "+rspTime)
	if err != nil {
		slog.Debug("GetSiteTime failed to parse response", "error", err)
		return result, err
	}
	return result, nil
}

func (od *onstepxDevice) SetSiteTime(t time.Time) error {

	var err error

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("SetSiteTime called when not connected")
		return ErrNotConnected
	}

	rspDate, err := od.sendCommand("SC"+t.Format("01/02/06"), RspBool, 0)
	if err != nil || rspDate != "1" {
		slog.Warn("Connect failed to set date", "error", err)
		return fmt.Errorf("failed to set date: %w", err)
	} else {
		slog.Info("Connect set date successfully")
	}
	rspTime, err := od.sendCommand("SL"+t.Format("03:04:05"), RspBool, 0)
	if err != nil || rspTime != "1" {
		slog.Warn("Connect failed to set time", "error", err)
		return fmt.Errorf("failed to set time: %w", err)
	} else {
		slog.Info("Connect set time successfully")
	}
	return nil
}

func (od *onstepxDevice) GetSiteLatitude() (float32, error) {

	var err error

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("GetSiteLatitude called when not connected")
		return 0.0, ErrNotConnected
	}

	// Get latitude
	rsp, err := od.sendCommand("Gt", RspHash, 0)
	if err != nil {
		slog.Debug("failed to get latitude", "error", err)
		return 0.0, err
	}
	lat, err := parseDDMM(rsp)
	if err != nil {
		slog.Debug("parse DDMM failed", "error", err)
		return 0.0, err
	}
	return lat, nil
}

func (od *onstepxDevice) SetSiteLatitude(lat float32) error {

	var err error

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("SetSiteLatitude called when not connected")
		return ErrNotConnected
	}

	// Set latitude
	rsp, err := od.sendCommand("St"+formatDDMM(lat, true), RspBool, 0)
	if err != nil || rsp != "1" {
		slog.Debug("failed to set latitude", "error", err)
		return err
	}
	return nil
}

func (od *onstepxDevice) GetSiteLongitude() (float32, error) {

	var err error

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("GetSiteLongitude called when not connected")
		return 0.0, ErrNotConnected
	}

	// Get longitude
	rsp, err := od.sendCommand("Gg", RspHash, 0)
	if err != nil {
		slog.Debug("failed to get longitude", "error", err)
		return 0.0, err
	}
	long, err := parseDDMM(rsp)
	if err != nil {
		slog.Debug("parse DDMM failed", "error", err)
		return 0.0, err
	}
	return long, nil
}

func (od *onstepxDevice) SetSiteLongitude(long float32) error {

	var err error

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("SetSiteLongitude called when not connected")
		return ErrNotConnected
	}

	// Set Longitude
	rsp, err := od.sendCommand("Sg"+formatDDDMM(long), RspBool, 0)
	if err != nil || rsp != "1" {
		slog.Debug("failed to set longitude", "error", err)
		return err
	}
	return nil
}

func (od *onstepxDevice) GetSiteElevation() (float32, error) {

	var err error

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("GetSiteElevation called when not connected")
		return 0.0, ErrNotConnected
	}

	// Get elevation
	rsp, err := od.sendCommand("Gv", RspHash, 0)
	if err != nil {
		slog.Debug("failed to get elevation", "error", err)
		return 0.0, err
	}
	elevation, err := strconv.ParseFloat(rsp, 32)
	if err != nil {
		slog.Debug("parse float32 failed", "error", err)
		return 0.0, err
	}
	return float32(elevation), nil
}

func (od *onstepxDevice) SetSiteElevation(elevation float32) error {

	var err error

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("SetSiteElevation called when not connected")
		return ErrNotConnected
	}

	// Set elevation
	rsp, err := od.sendCommand(fmt.Sprintf("Sv%3.1f", elevation), RspBool, 0)
	if err != nil || rsp != "1" {
		slog.Debug("failed to set elevation", "error", err)
		return err
	}
	return nil
}

func (od *onstepxDevice) GetSiteUTCOffset() (int, error) {

	var err error

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("GetSiteUTCOffset called when not connected")
		return 0.0, ErrNotConnected
	}

	// Get UTC offset
	rsp, err := od.sendCommand("GG", RspHash, 0)
	if err != nil {
		slog.Debug("failed to get UTC offset", "error", err)
		return 0, err
	}
	offset, err := parseHHMM(rsp)
	if err != nil {
		slog.Debug("parse UTC offset failed", "error", err)
		return 0, err
	}
	return int(math.Trunc(float64(offset))), nil
}

func (od *onstepxDevice) SetSiteUTCOffset(offset int) error {

	var err error

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("SetSiteUTCOffset called when not connected")
		return ErrNotConnected
	}

	// Set UTC offset
	rsp, err := od.sendCommand("SG"+formatHH(offset, true), RspBool, 0)
	if err != nil || rsp != "1" {
		slog.Debug("failed to set UTC offset", "error", err)
		return err
	}
	return nil
}

func (od *onstepxDevice) GetPositionRA() (float32, error) {

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("GetPositionRA called when not connected")
		return 0.0, ErrNotConnected
	}

	// Get current RA
	rsp, err := od.sendCommand("GR", RspHash, 0)
	if err != nil {
		slog.Debug("GetPositionRA request GR failed", "error", err)
		return 0.0, err
	}
	ra, err := parseHHMMSS(rsp)
	if err != nil {
		slog.Debug("GetPositionRA parse response failed", "error", err)
		return 0.0, err
	}
	return ra, nil
}

func (od *onstepxDevice) GetPositionDec() (float32, error) {

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("GetPositionDec called when not connected")
		return 0.0, ErrNotConnected
	}

	// Get current Dec
	rsp, err := od.sendCommand("GD", RspHash, 0)
	if err != nil {
		slog.Debug("GetPositionDec request GD failed", "error", err)
		return 0.0, err
	}
	dec, err := parseDDMMSS(rsp)
	if err != nil {
		slog.Debug("GetPositionDec parse response failed", "error", err)
		return 0.0, err
	}
	return dec, nil
}

func (od *onstepxDevice) GetPositionAlt() (float32, error) {

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("GetPositionAlt called when not connected")
		return 0.0, ErrNotConnected
	}

	// Get current altitude
	rsp, err := od.sendCommand("GA", RspHash, 0)
	if err != nil {
		slog.Debug("GetPositionAlt request GA failed", "error", err)
		return 0.0, err
	}
	alt, err := parseDDMMSS(rsp)
	if err != nil {
		slog.Debug("GetPositionAlt parse response failed", "error", err)
		return 0.0, err
	}
	return alt, nil
}

func (od *onstepxDevice) GetPositionAz() (float32, error) {

	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Must be connected
	if !od.isConnected {
		slog.Debug("GetPositionAz called when not connected")
		return 0.0, ErrNotConnected
	}

	// Get current azimuth
	rsp, err := od.sendCommand("GZ", RspHash, 0)
	if err != nil {
		slog.Debug("GetPositionAz request GZ failed", "error", err)
		return 0.0, err
	}
	az, err := parseDDMMSS(rsp)
	if err != nil {
		slog.Debug("GetPositionAz parse response failed", "error", err)
		return 0.0, err
	}
	return az, nil
}
