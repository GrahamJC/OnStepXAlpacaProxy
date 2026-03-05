package onstepx

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.bug.st/serial"
)

var (
	ErrNotConnected = errors.New("Device not connected")
	ErrInvalidValue = errors.New("Invalid value")
)

type PierSide int

const (
	PierSideEast = 0
	PierSideWest = 1
	PierSideNone = 2
)

type DeviceStatus struct {
	IsHome     bool
	IsParked   bool
	IsTracking bool
	IsSlewing  bool
	IsGuiding  bool
	PierSide   int
}

type OnStepXDevice interface {
	// Properties
	Connect() bool
	Disconnect()
	IsConnected() bool
	ComPort() string
	BaudRate() int
	// Basic OnStepX commands
	GetFirmwareVersion() (string, error)
	GetLocalDateTime() (time.Time, error)
	SetLocalDateTime(time.Time) error
	GetSiderealTime() (float64, error)
	GetUTCDateTime() (time.Time, error)
	GetSiteLatitude() (float64, error)
	SetSiteLatitude(float64) error
	GetSiteLongitude() (float64, error)
	SetSiteLongitude(float64) error
	GetSiteElevation() (float64, error)
	SetSiteElevation(float64) error
	GetSiteUTCOffset() (int, error)
	SetSiteUTCOffset(int) error
	GetRightAscension() (float64, error)
	GetDeclination() (float64, error)
	GetAltitude() (float64, error)
	GetAzimuth() (float64, error)
	GetAtHome() (bool, error)
	GetIsHoming() (bool, error)
	GetAtPark() (bool, error)
	GetIsParking() (bool, error)
	GetIsTracking() (bool, error)
	GetIsSlewing() (bool, error)
	GetIsGuiding() (bool, error)
	GetIsPulseGuiding() (bool, error)
	GetPierSide() (PierSide, error)
	GetRightAscensionRate() (float64, error)
	SetRightAscensionRate(float64) error
	GetDeclinationRate() (float64, error)
	SetDeclinationRate(float64) error
	GetPulseGuideRate() (float64, error)
	GetTargetRightAscension() (float64, error)
	SetTargetRightAscension(float64) error
	GetTargetDeclination() (float64, error)
	SetTargetDeclination(float64) error
	SetTargetAltitude(float64) error
	SetTargetAzimuth(float64) error
	GetTargetPierSide() (PierSide, error)
	GetTrackingRate() (string, error)
	SetTrackingRate(string) error
	SetTracking(bool) error
	AbortSlew() error
	StartHome() error
	SetRightAscensionSlewRate(float64) error
	StartRightAscensionSlewEast() error
	StartRightAscensionSlewWest() error
	StopRightAscensionSlew() error
	SetDeclinationSlewRate(float64) error
	StartDeclinationSlewNorth() error
	StartDeclinationSlewSouth() error
	StopDeclinationSlew() error
	StartPark() error
	Unpark() error
	SetPark() error
	StartSlewToTarget() error
	StartSlewToTargetAltAz() error
	SyncToTarget() error
	StartPulseGuide(string, int) error
	// Helpers
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

	// Close command channel (terminates goroutine to process commands)
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

	// Get OnStepX version
	version, err := od.sendCommand("GVN", RspHash, 0)
	if err != nil {
		return "", err
	}
	return version, nil
}

func (od *onstepxDevice) GetLocalDateTime() (time.Time, error) {
	// Result
	var result time.Time

	// Must be connected
	if err := od.checkConnected("GetLocalDateTime"); err != nil {
		return result, err
	}

	// Get local date and time
	rspDate, err := od.sendCommand("GC", RspHash, 0)
	if err != nil {
		return result, err
	}
	rspTime, err := od.sendCommand("GL", RspHash, 0)
	if err != nil {
		return result, err
	}

	// Parse result
	result, err = time.Parse("01/02/06 03:04:05", rspDate+" "+rspTime)
	if err != nil {
		slog.Debug("OnStepX request GC/GL response bad format", "response", rspDate+" "+rspTime, "error", err)
		return result, err
	}
	return result, nil
}

func (od *onstepxDevice) SetLocalDateTime(t time.Time) error {
	// Must be connected
	if err := od.checkConnected("SetLocalDateTime"); err != nil {
		return err
	}

	// Set local date/time
	rspDate, err := od.sendCommand("SC"+t.Format("01/02/06"), RspOne, 0)
	if err != nil {
		return err
	} else if rspDate != "1" {
		slog.Warn("failed to set OnStepX date")
		return fmt.Errorf("failed to set OnStepX date")
	}
	rspTime, err := od.sendCommand("SL"+t.Format("03:04:05"), RspOne, 0)
	if err != nil {
		return err
	} else if rspTime != "1" {
		slog.Warn("failed to set OnStepX time")
		return fmt.Errorf("failed to set OnStepX time")
	}
	return nil
}

func (od *onstepxDevice) GetSiderealTime() (float64, error) {
	// Result
	var result float64

	// Must be connected
	if err := od.checkConnected("GetSiderealTime"); err != nil {
		return result, err
	}

	// Get sidereal time
	rsp, err := od.sendCommand("GS", RspHash, 0)
	if err != nil {
		return result, err
	}

	// Parse result
	result, err = parseHHMMSS(rsp)
	if err != nil {
		slog.Debug("OnStepX GS response bad format", "response", rsp, "error", err)
		return result, err
	}
	return result, nil
}

func (od *onstepxDevice) GetUTCDateTime() (time.Time, error) {
	// Result
	var result time.Time

	// Must be connected
	if err := od.checkConnected("GetUTCDateTime"); err != nil {
		return result, err
	}

	// Get UTC date and time
	rspDate, err := od.sendCommand("GX81", RspHash, 0)
	if err != nil {
		return result, err
	}
	rspTime, err := od.sendCommand("GX80", RspHash, 0)
	if err != nil {
		return result, err
	}

	// Parse result
	result, err = time.Parse("01/02/06 15:04:05", rspDate+" "+rspTime)
	if err != nil {
		slog.Debug("OnStepX GX81/GX80 response bad format", "response", rspDate+" "+rspTime, "error", err)
		return result, err
	}
	return result, nil
}

func (od *onstepxDevice) GetSiteLatitude() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetSiteLatitude"); err != nil {
		return 0, err
	}

	// Get latitude
	rsp, err := od.sendCommand("GtH", RspHash, 0)
	if err != nil {
		return 0, err
	}
	lat, err := parseDDMMSS(rsp)
	if err != nil {
		slog.Debug("OnStepX Gt response bad format", "response", rsp, "error", err)
		return 0, err
	}
	return lat, nil
}

func (od *onstepxDevice) SetSiteLatitude(lat float64) error {
	// Must be connected
	if err := od.checkConnected("SetSiteLatitude"); err != nil {
		return err
	}

	// Set latitude
	rsp, err := od.sendCommand("St"+formatDDMMSS(lat, true, false), RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Warn("failed to set OnStepX site latitude")
		return fmt.Errorf("failed to set OnStepX site latitude")
	}
	return nil
}

func (od *onstepxDevice) GetSiteLongitude() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetSiteLongitude"); err != nil {
		return 0, err
	}

	// Get longitude
	rsp, err := od.sendCommand("GgH", RspHash, 0)
	if err != nil {
		return 0.0, err
	}
	long, err := parseDDMMSS(rsp)
	if err != nil {
		slog.Debug("OnStepX Gg response bad format", "response", rsp, "error", err)
		return 0.0, err
	}
	return long, nil
}

func (od *onstepxDevice) SetSiteLongitude(long float64) error {
	// Must be connected
	if err := od.checkConnected("SetSiteLongitude"); err != nil {
		return err
	}

	// Set Longitude
	rsp, err := od.sendCommand("Sg"+formatDDMMSS(long, false, true), RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Warn("failed to set OnStepX site longitude")
		return fmt.Errorf("failed to set OnStepX site longitude")
	}
	return nil
}

func (od *onstepxDevice) GetSiteElevation() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetSiteElevation"); err != nil {
		return 0, err
	}

	// Get elevation
	rsp, err := od.sendCommand("Gv", RspHash, 0)
	if err != nil {
		return 0.0, err
	}
	elevation, err := strconv.ParseFloat(rsp, 64)
	if err != nil {
		slog.Debug("OnStepX Gv response bad format", "response", rsp, "error", err)
		return 0.0, err
	}
	return elevation, nil
}

func (od *onstepxDevice) SetSiteElevation(elevation float64) error {
	// Must be connected
	if err := od.checkConnected("SetSiteElevation"); err != nil {
		return err
	}

	// Set elevation
	rsp, err := od.sendCommand(fmt.Sprintf("Sv%3.1f", elevation), RspOne, 0)
	if err != nil {
		return nil
	} else if rsp != "1" {
		slog.Warn("failed to set OnStepX site elevation")
		return fmt.Errorf("failed to set OnStepX site elevation")
	}
	return nil
}

func (od *onstepxDevice) GetSiteUTCOffset() (int, error) {
	// Must be connected
	if err := od.checkConnected("GetSiteUTCOffset"); err != nil {
		return 0, err
	}

	// Get UTC offset
	rsp, err := od.sendCommand("GG", RspHash, 0)
	if err != nil {
		return 0, err
	}
	offset, err := parseHHMM(rsp)
	if err != nil {
		slog.Debug("OnStepX GG response bad format", "response", rsp, "error", err)
		return 0, err
	}
	return int(math.Trunc(float64(offset))), nil
}

func (od *onstepxDevice) SetSiteUTCOffset(offset int) error {
	// Must be connected
	if err := od.checkConnected("SetSiteUTCOffset"); err != nil {
		return err
	}

	// Check min/max values
	if offset < -13 || offset > 12 {
		slog.Debug("UTC offset out of range (-13 to 12)", "value", offset)
		return ErrInvalidValue
	}

	// Set UTC offset
	rsp, err := od.sendCommand("SG"+formatHH(offset), RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Debug("failed to set OnStepX site UTC offset")
		return fmt.Errorf("failed to set OnStepX site UTC offset")
	}
	return nil
}

func (od *onstepxDevice) GetRightAscension() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetRightAscension"); err != nil {
		return 0, err
	}

	// Get current RA
	rsp, err := od.sendCommand("GRH", RspHash, 0)
	if err != nil {
		return 0.0, err
	}
	ra, err := parseHHMMSS(rsp)
	if err != nil {
		slog.Debug("OnStepX response GR bad format", "response", rsp, "error", err)
		return 0.0, err
	}
	return ra, nil
}

func (od *onstepxDevice) GetDeclination() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetDeclination"); err != nil {
		return 0, err
	}

	// Get current Dec
	rsp, err := od.sendCommand("GDH", RspHash, 0)
	if err != nil {
		return 0.0, err
	}
	dec, err := parseDDMMSS(rsp)
	if err != nil {
		slog.Debug("OnStepX response GD bad format", "response", rsp, "error", err)
		return 0.0, err
	}
	return dec, nil
}

func (od *onstepxDevice) GetAltitude() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetAltitude"); err != nil {
		return 0, err
	}

	// Get current altitude
	rsp, err := od.sendCommand("GAH", RspHash, 0)
	if err != nil {
		return 0.0, err
	}
	alt, err := parseDDMMSS(rsp)
	if err != nil {
		slog.Debug("OnStepX response GA bad format", "response", rsp, "error", err)
		return 0.0, err
	}
	return alt, nil
}

func (od *onstepxDevice) GetAzimuth() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetAzimuth"); err != nil {
		return 0, err
	}

	// Get current azimuth
	rsp, err := od.sendCommand("GZH", RspHash, 0)
	if err != nil {
		return 0.0, err
	}
	az, err := parseDDMMSS(rsp)
	if err != nil {
		slog.Debug("OnStepX response GZ bad format", "response", rsp, "error", err)
		return 0.0, err
	}
	return az, nil
}

func (od *onstepxDevice) GetAtHome() (bool, error) {
	// Must be connected
	if err := od.checkConnected("GetAtHome"); err != nil {
		return false, err
	}

	// Get status from OnStepX
	sts, err := od.sendCommand("GU", RspHash, 0)
	if err != nil {
		return false, err
	}
	return strings.Contains(sts, "H"), nil
}

func (od *onstepxDevice) GetIsHoming() (bool, error) {
	// Must be connected
	if err := od.checkConnected("GetIsHoming"); err != nil {
		return false, err
	}

	// Get status from OnStepX
	sts, err := od.sendCommand("GU", RspHash, 0)
	if err != nil {
		return false, err
	}
	return strings.Contains(sts, "h"), nil
}

func (od *onstepxDevice) GetAtPark() (bool, error) {
	// Must be connected
	if err := od.checkConnected("GetAtPark"); err != nil {
		return false, err
	}

	// Get status from OnStepX
	sts, err := od.sendCommand("GU", RspHash, 0)
	if err != nil {
		return false, err
	}
	return strings.Contains(sts, "P"), nil
}

func (od *onstepxDevice) GetIsParking() (bool, error) {
	// Must be connected
	if err := od.checkConnected("GetIsParking"); err != nil {
		return false, err
	}

	// Get status from OnStepX
	sts, err := od.sendCommand("GU", RspHash, 0)
	if err != nil {
		return false, err
	}
	return strings.Contains(sts, "I"), nil
}

func (od *onstepxDevice) GetIsTracking() (bool, error) {
	// Must be connected
	if err := od.checkConnected("GetIsTracking"); err != nil {
		return false, err
	}

	// Get status from OnStepX
	sts, err := od.sendCommand("GU", RspHash, 0)
	if err != nil {
		return false, err
	}
	return !strings.Contains(sts, "n"), nil
}

func (od *onstepxDevice) GetIsSlewing() (bool, error) {
	// Must be connected
	if err := od.checkConnected("GetIsSlewing"); err != nil {
		return false, err
	}

	// Get status from OnStepX
	sts, err := od.sendCommand("GU", RspHash, 0)
	if err != nil {
		return false, err
	}
	return !strings.Contains(sts, "N"), nil
}

func (od *onstepxDevice) GetIsGuiding() (bool, error) {
	// Must be connected
	if err := od.checkConnected("GetIsGuiding"); err != nil {
		return false, err
	}

	// Get status from OnStepX
	sts, err := od.sendCommand("GU", RspHash, 0)
	if err != nil {
		return false, err
	}
	return strings.Contains(sts, "g"), nil
}

func (od *onstepxDevice) GetIsPulseGuiding() (bool, error) {
	// Must be connected
	if err := od.checkConnected("GetIsPulseGuiding"); err != nil {
		return false, err
	}

	// Get status from OnStepX
	sts, err := od.sendCommand("GU", RspHash, 0)
	if err != nil {
		return false, err
	}
	return strings.Contains(sts, "G"), nil
}

func (od *onstepxDevice) GetPierSide() (PierSide, error) {
	// Must be connected
	if err := od.checkConnected("GetPierSide"); err != nil {
		return PierSideNone, err
	}

	// Get pier side from OnStepX
	ps, err := od.sendCommand("Gm", RspHash, 0)
	if err != nil {
		return PierSideNone, err
	}
	switch ps {
	case "E":
		return PierSideEast, nil
	case "W":
		return PierSideWest, nil
	case "N":
		return PierSideNone, nil
	default:
		return PierSideNone, fmt.Errorf("Invalid pier side '%s' returned by OnStepX", ps)
	}
}

func (od *onstepxDevice) GetRightAscensionRate() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetRightAscensionRate"); err != nil {
		return 0, err
	}

	// Get right ascension rate offset from OnStepX
	rsp, err := od.sendCommand("GXTR", RspHash, 0)
	if err != nil {
		return 0, err
	}
	offset, err := strconv.ParseFloat(rsp, 64)
	if err != nil {
		slog.Debug("OnStepX GXTR response bad format", "response", rsp, "error", err)
		return 0, err
	}
	return offset, nil
}

func (od *onstepxDevice) SetRightAscensionRate(offset float64) error {
	// Must be connected
	if err := od.checkConnected("GetRightAscensionRate"); err != nil {
		return err
	}

	// Check min/max values
	if offset < -1800 || offset > 1800 {
		slog.Debug("Right ascension rate offset out of range (-1800 to 1800)", "value", offset)
		return ErrInvalidValue
	}

	// Set declination rate offset
	rsp, err := od.sendCommand(fmt.Sprintf("SXTR,%f", offset), RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Debug("failed to set OnStepX right ascension rate offset")
		return fmt.Errorf("failed to set OnStepX right ascension rate offset")
	}
	return nil
}

func (od *onstepxDevice) GetDeclinationRate() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetDeclinationRate"); err != nil {
		return 0, err
	}

	// Get declination rate offset from OnStepX
	rsp, err := od.sendCommand("GXTD", RspHash, 0)
	if err != nil {
		return 0, err
	}
	offset, err := strconv.ParseFloat(rsp, 64)
	if err != nil {
		slog.Debug("OnStepX GXTD response bad format", "response", rsp, "error", err)
		return 0, err
	}
	return offset, nil
}

func (od *onstepxDevice) SetDeclinationRate(offset float64) error {
	// Must be connected
	if err := od.checkConnected("GetDeclinationRate"); err != nil {
		return err
	}

	// Check min/max values
	if offset < -1800 || offset > 1800 {
		slog.Debug("Declination rate offset out of range (-1800 to 1800)", "value", offset)
		return ErrInvalidValue
	}

	// Set declination rate offset
	rsp, err := od.sendCommand(fmt.Sprintf("SXTD,%f", offset), RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Debug("failed to set OnStepX declination rate offset")
		return fmt.Errorf("failed to set OnStepX declination rate offset")
	}
	return nil
}

func (od *onstepxDevice) GetPulseGuideRate() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetGuideRateRightAscension"); err != nil {
		return 0, err
	}

	// Get guide rate from OnStepX
	rsp, err := od.sendCommand("GU", RspHash, 0)
	if err != nil {
		return 0, err
	}
	rate := 0.0
	switch idx := len(rsp) - 3; rsp[idx : idx+1] {
	case "0":
		rate = 0.25
	case "1":
		rate = 0.5
	case "2":
		rate = 1.0
	default:
		slog.Debug("OnStepX GU response bad format (pulse guide rate)", "response", rsp, "error", err)
		return 0, fmt.Errorf("failed to get OnStepX pulse guide rate")
	}
	return rate * 15 / 3600, nil
}

func (od *onstepxDevice) GetTargetRightAscension() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetTargetRightAscension"); err != nil {
		return 0, err
	}

	// Get trrget rigth ascension
	rsp, err := od.sendCommand("GrH", RspHash, 0)
	if err != nil {
		return 0.0, err
	}
	ra, err := parseHHMMSS(rsp)
	if err != nil {
		slog.Debug("OnStepX response Gr bad format", "response", rsp, "error", err)
		return 0.0, err
	}
	return ra, nil
}

func (od *onstepxDevice) SetTargetRightAscension(ra float64) error {
	// Must be connected
	if err := od.checkConnected("SetTargetRightAscension"); err != nil {
		return err
	}

	// Check min/max values
	if ra < 0 || ra >= 24 {
		slog.Debug("Target right ascension out of range (0 to 24)", "value", ra)
		return ErrInvalidValue
	}

	// Set target right ascension
	rsp, err := od.sendCommand("Sr"+formatHHMMSS(ra), RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Debug("failed to get OnStepX target right ascension")
		return fmt.Errorf("failed to get OnStepX target right ascension")
	}
	return nil
}

func (od *onstepxDevice) GetTargetDeclination() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetTargetDeclination"); err != nil {
		return 0, err
	}

	// Get target declination
	rsp, err := od.sendCommand("GdH", RspHash, 0)
	if err != nil {
		return 0.0, err
	}
	dec, err := parseDDMMSS(rsp)
	if err != nil {
		slog.Debug("OnStepX response Gd bad format", "response", rsp, "error", err)
		return 0.0, err
	}
	return dec, nil
}

func (od *onstepxDevice) SetTargetDeclination(dec float64) error {
	// Must be connected
	if err := od.checkConnected("SetTargetDeclination"); err != nil {
		return err
	}

	// Check min/max values
	if dec < -90 || dec > 90 {
		return ErrInvalidValue
	}

	// Set target declination
	rsp, err := od.sendCommand("Sd"+formatDDMMSS(dec, true, false), RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Debug("failed to get OnStepX target declination")
		return fmt.Errorf("failed to get OnStepX target declination")
	}
	return nil
}

func (od *onstepxDevice) GetTargetAltitude() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetTargetAltitude"); err != nil {
		return 0, err
	}

	// Get target altitude
	rsp, err := od.sendCommand("GaH", RspHash, 0)
	if err != nil {
		return 0.0, err
	}
	dec, err := parseDDMMSS(rsp)
	if err != nil {
		slog.Debug("OnStepX response GaH bad format", "response", rsp, "error", err)
		return 0.0, err
	}
	return dec, nil
}

func (od *onstepxDevice) SetTargetAltitude(alt float64) error {
	// Must be connected
	if err := od.checkConnected("SetTargetAltitude"); err != nil {
		return err
	}

	// Check min/max values
	if alt < -90 || alt > 90 {
		slog.Debug("Target altitude out of range (-90 to 90)", "value", alt)
		return ErrInvalidValue
	}

	// Set target altitude
	rsp, err := od.sendCommand("Sa"+formatDDMMSS(alt, true, false), RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Debug("failed to set OnStepX target altitude")
		return fmt.Errorf("failed to set OnStepX target altitude")
	}
	return nil
}

func (od *onstepxDevice) GetTargetAzimuth() (float64, error) {
	// Must be connected
	if err := od.checkConnected("GetTargetAzimuth"); err != nil {
		return 0, err
	}

	// Get target azimuth
	rsp, err := od.sendCommand("GzH", RspHash, 0)
	if err != nil {
		return 0.0, err
	}
	dec, err := parseDDMMSS(rsp)
	if err != nil {
		slog.Debug("OnStepX response GzH bad format", "response", rsp, "error", err)
		return 0.0, err
	}
	return dec, nil
}

func (od *onstepxDevice) SetTargetAzimuth(az float64) error {
	// Must be connected
	if err := od.checkConnected("SetTargetAzimuth"); err != nil {
		return err
	}

	// Check min/max values
	if az < 0 || az > 360 {
		slog.Debug("Target azimuth out of range (0 to 360)", "value", az)
		return ErrInvalidValue
	}

	// Set target azimuth
	rsp, err := od.sendCommand("Sz"+formatDDMMSS(az, false, true), RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Debug("failed to set OnStepX target azimuth")
		return fmt.Errorf("failed to set OnStepX target azimuth")
	}
	return nil
}

func (od *onstepxDevice) GetTargetPierSide() (PierSide, error) {
	// Must be connected
	if err := od.checkConnected("GetTargetPierSide"); err != nil {
		return PierSideNone, err
	}

	// Get OnStepX target pier side
	rsp, err := od.sendCommand("MD", RspOne, 0)
	if err != nil {
		return PierSideNone, err
	}
	switch rsp {
	case "0":
		return PierSideEast, nil
	case "1":
		return PierSideWest, nil
	case "2":
		return PierSideNone, nil
	default:
		slog.Debug(fmt.Sprintf("Invalid response '%s' to OnStepX command 'MD'", rsp))
		return PierSideNone, fmt.Errorf("failed to get pier side from OnStepX")
	}
}

func (od *onstepxDevice) GetTrackingRate() (string, error) {
	// Must be connected
	if err := od.checkConnected("GetTrackingRate"); err != nil {
		return "", err
	}

	// Get tracking rate
	rsp, err := od.sendCommand("GU", RspHash, 0)
	if err != nil {
		return "", err
	}
	if strings.Contains(rsp, "(") {
		return "Lunar", nil
	} else if strings.Contains(rsp, "O") {
		return "Solar", nil
	} else if strings.Contains(rsp, "k") {
		return "King", nil
	} else {
		return "Sidereal", nil
	}
}

func (od *onstepxDevice) SetTrackingRate(rate string) error {
	// Must be connected
	if err := od.checkConnected("SetTrackingRate"); err != nil {
		return err
	}

	// Set tracking rate
	cmd := ""
	switch rate {
	case "Sidereal":
		cmd = "TQ"
	case "Lunar":
		cmd = "TL"
	case "Solar":
		cmd = "TS"
	case "King":
		cmd = "TK"
	default:
		return errors.New("Invalid tracking rate")
	}
	_, err := od.sendCommand(cmd, RspNone, 0)
	if err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) SetTracking(state bool) error {
	// Must be connected
	if err := od.checkConnected("SetTracking"); err != nil {
		return err
	}

	// Start/stop tracking
	if state {
		rsp, err := od.sendCommand("Te", RspOne, 0)
		if err != nil {
			return err
		} else if rsp != "1" {
			slog.Debug("failed to start tracking")
			return fmt.Errorf("failed to start tracking")
		}
	} else {
		rsp, err := od.sendCommand("Td", RspOne, 0)
		if err != nil {
			return err
		} else if rsp != "1" {
			slog.Debug("failed to stop tracking")
			return fmt.Errorf("failed to stop tracking")
		}
	}
	return nil
}

func (od *onstepxDevice) AbortSlew() error {
	// Must be connected
	if err := od.checkConnected("AbortSlew"); err != nil {
		return err
	}

	// Abort OnStepX slew
	_, err := od.sendCommand("Q", RspNone, 0)
	if err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) StartHome() error {
	// Must be connected
	if err := od.checkConnected("FindHome"); err != nil {
		return err
	}

	// Start move to home position
	_, err := od.sendCommand("hC", RspNone, 0)
	if err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) SetRightAscensionSlewRate(degPerSec float64) error {
	// Must be connected
	if err := od.checkConnected("SetRightAscensionSlewRate"); err != nil {
		return err
	}

	// Set right ascension slew rate
	if _, err := od.sendCommand(fmt.Sprintf("RA%08.5f", degPerSec), RspNone, 0); err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) StartRightAscensionSlewEast() error {
	// Must be connected
	if err := od.checkConnected("StartRightAscensionSlewEast"); err != nil {
		return err
	}

	// Start right ascension move east
	_, err := od.sendCommand("Me", RspNone, 0)
	if err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) StartRightAscensionSlewWest() error {
	// Must be connected
	if err := od.checkConnected("StartRightAscensionSlewWest"); err != nil {
		return err
	}

	// Start right ascension move west
	_, err := od.sendCommand("Mw", RspNone, 0)
	if err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) StopRightAscensionSlew() error {
	// Must be connected
	if err := od.checkConnected("StopRightAscensionSlew"); err != nil {
		return err
	}

	// Stop right ascension movement
	_, err := od.sendCommand("Qe", RspNone, 0)
	if err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) SetDeclinationSlewRate(degPerSec float64) error {
	// Must be connected
	if err := od.checkConnected("SetDeclinationSlewRate"); err != nil {
		return err
	}

	// Set declination slew rate
	if _, err := od.sendCommand(fmt.Sprintf("RE%08.5f", degPerSec), RspNone, 0); err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) StartDeclinationSlewNorth() error {
	// Must be connected
	if err := od.checkConnected("StartDeclinationSlewNorth"); err != nil {
		return err
	}

	// Start declination move north
	_, err := od.sendCommand("Mn", RspNone, 0)
	if err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) StartDeclinationSlewSouth() error {
	// Must be connected
	if err := od.checkConnected("StartDeclinationSlewSouth"); err != nil {
		return err
	}

	// Start declination move south
	_, err := od.sendCommand("Ms", RspNone, 0)
	if err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) StopDeclinationSlew() error {
	// Must be connected
	if err := od.checkConnected("StopDeclinationSlew"); err != nil {
		return err
	}

	// Stop declination movement
	_, err := od.sendCommand("Qn", RspNone, 0)
	if err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) StartPark() error {
	// Must be connected
	if err := od.checkConnected("StartPark"); err != nil {
		return err
	}

	// Start move to park position
	rsp, err := od.sendCommand("hP", RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Debug("failed to start OnStepX park operation")
		return fmt.Errorf("failed to start OnStepX park operation")
	}
	return nil
}

func (od *onstepxDevice) Unpark() error {
	// Must be connected
	if err := od.checkConnected("Unpark"); err != nil {
		return err
	}

	// Restore OnStepX to operation
	rsp, err := od.sendCommand("hR", RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Debug("failed to unpark OnStepX")
		return fmt.Errorf("failed to unpark OnStepX")
	}
	return nil
}

func (od *onstepxDevice) SetPark() error {
	// Must be connected
	if err := od.checkConnected("SetPark"); err != nil {
		return err
	}

	// Set OnStepX park to current position
	rsp, err := od.sendCommand("hQ", RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Debug("failed to set OnStepX park position")
		return fmt.Errorf("failed to set OnStepX park position")
	}
	return nil
}

func (od *onstepxDevice) StartSlewToTarget() error {
	// Must be connected
	if err := od.checkConnected("StartSlewToTarget"); err != nil {
		return err
	}

	// Start slew to target
	rsp, err := od.sendCommand("MS", RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "0" {
		slog.Debug("failed to start slew to target", "code", rsp)
		return fmt.Errorf("failed to start slew to target")
	}
	return nil
}

func (od *onstepxDevice) StartSlewToTargetAltAz() error {
	// Must be connected
	if err := od.checkConnected("StartSlewToTargetAltAz"); err != nil {
		return err
	}

	// Start slew to target
	rsp, err := od.sendCommand("MA", RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "0" {
		slog.Debug("failed to start slew to target", "code", rsp)
		return fmt.Errorf("failed to start slew to target")
	}
	return nil
}

func (od *onstepxDevice) SyncToTarget() error {
	// Must be connected
	if err := od.checkConnected("SyncToTarget"); err != nil {
		return err
	}

	// Sync to target
	_, err := od.sendCommand("CS", RspNone, 0)
	if err != nil {
		return err
	}
	return nil
}

func (od *onstepxDevice) StartPulseGuide(dirn string, durn int) error {
	// Must be connected
	if err := od.checkConnected("SyncToTarget"); err != nil {
		return err
	}

	// Start pulse guiding
	rsp, err := od.sendCommand(fmt.Sprintf("MG%s%d", dirn, durn), RspOne, 0)
	if err != nil {
		return err
	} else if rsp != "1" {
		slog.Debug("pulse guide failed")
		return fmt.Errorf("pulse guide failed")
	}
	return nil
}

// Helpers
func (od *onstepxDevice) checkConnected(action string) error {
	// Lock device
	od.mutex.Lock()
	defer od.mutex.Unlock()

	// Check connected state
	if !od.isConnected {
		slog.Debug("device not connected", "action", action)
		return ErrNotConnected
	}
	return nil
}
