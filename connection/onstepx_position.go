package connection

func (osx *OnstepXDevice) GetRightAscension() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get current RA
	return osx.sendReceiveHours("GRH")
}

func (osx *OnstepXDevice) GetDeclination() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get current Dec
	return osx.sendReceiveDegrees("GDH")
}

func (osx *OnstepXDevice) GetAltitude() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get current altitude
	return osx.sendReceiveFloat("GAH")
}

func (osx *OnstepXDevice) GetAzimuth() (float64, error) {

	// Must be connected
	if !osx.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get current azimuth
	return osx.sendReceiveDegrees("GZH")
}
