package connection

import (
	"fmt"
)

func (osx *OnstepXDevice) GetTargetRightAscension() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get target rigth ascension
	return osx.sendReceiveHours("GrH")
}

func (osx *OnstepXDevice) SetTargetRightAscension(ra float64) error {

	// Check value
	if ra < 0 || ra >= 24 {
		return ErrInvalidValue
	}

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Set target right ascension
	return osx.sendBool("Sr" + formatHours(ra))
}

func (osx *OnstepXDevice) GetTargetDeclination() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get target declination
	return osx.sendReceiveDegrees("GdH")
}

func (osx *OnstepXDevice) SetTargetDeclination(dec float64) error {

	// Check value
	if dec < -90 || dec > 90 {
		return ErrInvalidValue
	}

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Set target declination
	return osx.sendBool("Sd" + formatDegreesSigned(dec))
}

func (osx *OnstepXDevice) GetTargetAltitude() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get target altitude
	return osx.sendReceiveDegrees("GaH")
}

func (osx *OnstepXDevice) SetTargetAltitude(alt float64) error {

	// Check value
	if alt < -90 || alt > 90 {
		return ErrInvalidValue
	}

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Set target altitude
	return osx.sendBool("Sa" + formatDegreesSigned(alt))
}

func (osx *OnstepXDevice) GetTargetAzimuth() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get target azimuth
	return osx.sendReceiveDegrees("GzH")
}

func (osx *OnstepXDevice) SetTargetAzimuth(az float64) error {

	// Check value
	if az < 0 || az > 360 {
		return ErrInvalidValue
	}

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Set target azimuth
	return osx.sendBool("Sz" + formatDegrees(az))
}

func (osx *OnstepXDevice) GetTargetPierSide() (PierSide, error) {

	// Must be connected
	if !osx.IsConnected() {
		return PierSideNone, ErrNotConnected
	}

	// Get OnStepX target pier side
	rsp, err := osx.sendReceive("MD")
	if err != nil {
		return PierSideNone, fmt.Errorf("MD: %w", err)
	}
	switch rsp {
	case "0":
		return PierSideEast, nil
	case "1":
		return PierSideWest, nil
	case "2":
		return PierSideNone, nil
	default:
		return PierSideNone, ErrInvalidResponse
	}
}
