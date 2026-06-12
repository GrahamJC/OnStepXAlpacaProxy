package connection

import (
	"fmt"
	"strconv"
)

func (osx *OnstepXDevice) GetSiteLatitude() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get latitude
	return osx.sendReceiveDegrees("GtH")
}

func (osx *OnstepXDevice) SetSiteLatitude(lat float64) error {

	// Check value
	if (lat < -90) || (lat > 90) {
		return ErrInvalidValue
	}

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Set latitude
	return osx.sendBool(fmt.Sprintf("St%s", formatDegrees(lat, true)))
}

func (osx *OnstepXDevice) GetSiteLongitude() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get longitude
	return osx.sendReceiveDegrees("GgH")
}

func (osx *OnstepXDevice) SetSiteLongitude(lng float64) error {

	// Check value
	if (lng < -180) || (lng > 360) {
		return ErrInvalidValue
	}

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Set longitude
	return osx.sendBool(fmt.Sprintf("Sg%s", formatDegrees(lng, false)))
}

func (osx *OnstepXDevice) GetSiteElevation() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get elevation
	return osx.sendReceiveFloat("Gv")
}

func (osx *OnstepXDevice) SetSiteElevation(elevation float64) error {

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Set elevation
	return osx.sendBool(fmt.Sprintf("Sv%3.1f", elevation))
}

func (osx *OnstepXDevice) GetSiteUTCOffset() (int, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get UTC offset
	return osx.sendReceiveInt("GG")
}

func (osx *OnstepXDevice) SetSiteUTCOffset(offset int) error {

	// Check value
	if offset < -12 || offset > 12 {
		return ErrInvalidValue
	}

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Set UTC offset
	return osx.sendBool("SG" + strconv.Itoa(offset))
}
