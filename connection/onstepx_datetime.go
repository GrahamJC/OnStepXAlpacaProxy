package connection

import (
	"fmt"
	"time"
)

func (osx *OnstepXDevice) GetLocalDateTime() (time.Time, error) {

	// Must be connected
	if !osx.IsConnected() {
		return time.Time{}, ErrNotConnected
	}

	// Get local date and time
	rspDate, err := osx.sendReceive("GC")
	if err != nil {
		return time.Time{}, err
	}
	rspTime, err := osx.sendReceive("GL")
	if err != nil {
		return time.Time{}, err
	}

	// Parse result
	result, err := time.Parse("01/02/06 15:04:05", rspDate+" "+rspTime)
	if err != nil {
		return time.Time{}, err
	}
	return result, nil
}

func (osx *OnstepXDevice) SetLocalDateTime(t time.Time) error {

	// Must be connected
	if !osx.IsConnected() {
		return ErrNotConnected
	}

	// Set local date/time
	err := osx.sendBool("SC" + t.Format("01/02/06"))
	if err != nil {
		return err
	}
	err = osx.sendBool("SL" + t.Format("15:04:05"))
	if err != nil {
		return err
	}
	return nil
}

func (osx *OnstepXDevice) GetSiderealTime() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get sidereal time
	result, err := osx.sendReceiveHours("GS")
	if err != nil {
		return 0, fmt.Errorf("GS: %w", err)
	}
	return result, nil
}

func (osx *OnstepXDevice) GetUTCDateTime() (time.Time, error) {

	// Must be connected
	if !osx.IsConnected() {
		return time.Time{}, ErrNotConnected
	}

	// Get UTC date and time
	rspDate, err := osx.sendReceive("GX81")
	if err != nil {
		return time.Time{}, err
	}
	rspTime, err := osx.sendReceive("GX80")
	if err != nil {
		return time.Time{}, err
	}

	// Parse result
	result, err := time.Parse("01/02/06 15:04:05", rspDate+" "+rspTime)
	if err != nil {
		return time.Time{}, err
	}
	return result, nil
}
