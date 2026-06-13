package connection

import (
	"strings"
)

func (osx *OnstepXDevice) GetAtHome() (bool, error) {

	// Must be connected
	if !osx.IsConnected() {
		return false, ErrNotConnected
	}

	// Get status from OnStepX
	sts, err := osx.sendReceive("GU")
	if err != nil {
		return false, err
	}
	return strings.Contains(sts, "H"), nil
}

func (osx *OnstepXDevice) GetIsHoming() (bool, error) {

	// Must be connected
	if !osx.IsConnected() {
		return false, ErrNotConnected
	}

	// Get status from OnStepX
	sts, err := osx.sendReceive("GU")
	if err != nil {
		return false, err
	}
	return strings.Contains(sts, "H"), nil
}

func (osx *OnstepXDevice) StartHome() error {

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Start move to home position
	return osx.send("hC")
}
