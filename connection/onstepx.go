package connection

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrNotConnected    = errors.New("Not connected")
	ErrInvalidValue    = errors.New("Invalid value")
	ErrFailed          = errors.New("Request failed")
	ErrInvalidResponse = errors.New("Invalid response")
)

type PierSide int

const (
	PierSideEast = 0
	PierSideWest = 1
	PierSideNone = 2
)

type OnstepXDevice struct {
	transport Transport
}

// Construction
func NewOnstepX(transport Transport) *OnstepXDevice {
	return &OnstepXDevice{
		transport: transport,
	}
}

// Status
func (osx *OnstepXDevice) IsConnected() bool {
	return osx.transport.IsConnected()
}

func (osx *OnstepXDevice) Connect() error {
	return osx.transport.Open()
}

func (osx *OnstepXDevice) Disconnect() error {
	return osx.transport.Close()
}

// Version information
func (osx *OnstepXDevice) GetVersionProduct() (string, error) {
	if !osx.IsConnected() {
		return "", ErrNotConnected
	}
	return osx.sendReceive("GVP")
}

func (osx *OnstepXDevice) GetVersionNumber() (string, error) {
	if !osx.IsConnected() {
		return "", ErrNotConnected
	}
	return osx.sendReceive("GVN")
}

func (osx *OnstepXDevice) GetVersionFull() (string, error) {
	if !osx.IsConnected() {
		return "", ErrNotConnected
	}
	return osx.sendReceive("GVM")
}

// Helpers
func (osx *OnstepXDevice) send(cmd string) error {
	return osx.transport.Send(cmd)
}

func (osx *OnstepXDevice) sendReceive(cmd string) (string, error) {
	return osx.transport.SendReceive(cmd)
}

func (osx *OnstepXDevice) sendBool(cmd string) error {
	rsp, err := osx.transport.SendReceive(cmd)
	if err != nil {
		return fmt.Errorf("%s: %w", cmd, err)
	} else if strings.TrimSpace(rsp) == "0" {
		return ErrFailed
	} else if strings.TrimSpace(rsp) != "1" {
		return ErrInvalidResponse
	}
	return nil
}

func (osx *OnstepXDevice) sendReceiveDegrees(cmd string) (float64, error) {
	rsp, err := osx.transport.SendReceive(cmd)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", cmd, err)
	}
	return parseDegrees(rsp)
}

func (osx *OnstepXDevice) sendReceiveHours(cmd string) (float64, error) {
	rsp, err := osx.transport.SendReceive(cmd)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", cmd, err)
	}
	return parseHours(rsp)
}

func (osx *OnstepXDevice) sendReceiveInt(cmd string) (int, error) {
	rsp, err := osx.transport.SendReceive(cmd)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", cmd, err)
	}
	return strconv.Atoi(rsp)
}

func (osx *OnstepXDevice) sendReceiveFloat(cmd string) (float64, error) {
	rsp, err := osx.transport.SendReceive(cmd)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", cmd, err)
	}
	return strconv.ParseFloat(rsp, 64)
}

// Coordinate conversions
func parseDegrees(s string) (float64, error) {
	re := regexp.MustCompile(`^(?P<sign>[\+-]?)(?P<deg>[0-9]*)(?:[:\*](?P<min>[0-9]*)(?:[:'](?P<sec>[0-9\.]*))?)?$`)
	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return 0, fmt.Errorf("error parsing degrees: %s", s)
	}
	sign := 1.0
	if match[1] == "-" {
		sign = -1.0
	}
	deg, _ := strconv.ParseFloat(match[2], 64)
	min, _ := strconv.ParseFloat(match[3], 64)
	sec, _ := strconv.ParseFloat(match[4], 64)
	return sign * (deg + min/60.0 + sec/3600.0), nil
}

func formatDegrees(deg float64) string {
	d := int(math.Floor(deg))
	m := int(math.Floor((deg - float64(d)) * 60.0))
	s := (deg - float64(d) - float64(m)/60.0) * 3600.0
	return fmt.Sprintf("%02d*%02d:%04.1f", d, m, s)
}

func formatDegreesSigned(deg float64) string {
	sign := ""
	if deg < 0 {
		sign = "-"
		deg = -deg
	} else {
		sign = "+"
	}
	d := int(math.Floor(deg))
	m := int(math.Floor((deg - float64(d)) * 60.0))
	s := (deg - float64(d) - float64(m)/60.0) * 3600.0
	return fmt.Sprintf("%s%02d*%02d:%04.1f", sign, d, m, s)
}

func parseHours(s string) (float64, error) {
	re := regexp.MustCompile(`^(?P<sign>[\+-]?)(?P<hrs>[0-9]*)(?:[:](?P<min>[0-9]*)(?:[:'](?P<sec>[0-9\.]*))?)?$`)
	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return 0, fmt.Errorf("error parsing hours: %s", s)
	}
	sign := 1.0
	if match[1] == "-" {
		sign = -1.0
	}
	hrs, _ := strconv.ParseFloat(match[2], 64)
	min, _ := strconv.ParseFloat(match[3], 64)
	sec, _ := strconv.ParseFloat(match[4], 64)
	return sign * (hrs + min/60.0 + sec/3600.0), nil
}

func formatHours(hrs float64) string {
	h := int(math.Floor(hrs))
	m := int(math.Floor((hrs - float64(h)) * 60.0))
	s := (hrs - float64(h) - float64(m)/60.0) * 3600.0
	return fmt.Sprintf("%02d*%02d:%04.1f", h, m, s)
}

/*



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
*/
