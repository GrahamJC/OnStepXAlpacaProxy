package connection

import "strings"

func (osx *OnstepXDevice) GetAtPark() (bool, error) {

	// Must be connected
	if !osx.IsConnected() {
		return false, ErrNotConnected
	}

	// Get status from OnStepX
	sts, err := osx.sendReceive("GU")
	if err != nil {
		return false, err
	}
	return strings.Contains(sts, "P"), nil
}

func (osx *OnstepXDevice) GetIsParking() (bool, error) {

	// Must be connected
	if !osx.IsConnected() {
		return false, ErrNotConnected
	}

	// Get status from OnStepX
	sts, err := osx.sendReceive("GU")
	if err != nil {
		return false, err
	}
	return strings.Contains(sts, "I"), nil
}

func (osx *OnstepXDevice) StartPark() error {

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Start move to home position
	return osx.send("hP")
}

func (osx *OnstepXDevice) Unpark() error {

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Restore OnStepX to operation
	return osx.sendBool("hR")
}

func (osx *OnstepXDevice) SetPark() error {

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Set OnStepX park to current position
	return osx.sendBool("hQ")
}
