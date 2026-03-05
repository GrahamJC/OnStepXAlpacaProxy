package alpaca

import (
	"fmt"
	"net/http"
	"time"

	"onstepx-alpaca-proxy/onstepx"
)

// API holds all dependencies for the Alpaca API handlers.
type TelescopeAPI struct {
	appVersion  string
	deviceNo    int
	device      onstepx.OnStepXDevice
	name        string
	description string
	info        string
	apiVersion  int
}

// Alpaca enumerations
const (
	PierSideUnknown = -1
	PierSideEast    = 0
	PierSideWest    = 1
)

const (
	EquOther       = 0
	EquTopocentric = 1
	EquJ2000       = 2
	EquJ2050       = 3
	EquB1950       = 4
)

const (
	TrackingSidereal = 0
	TrackingLunar    = 1
	TrackingSolar    = 2
	TrackingKing     = 3
)

// NewAPI creates a new telescope API instance.
func NewTelescopeAPI(appVersion string, deviceNo int, device onstepx.OnStepXDevice) *TelescopeAPI {
	return &TelescopeAPI{
		appVersion:  appVersion,
		deviceNo:    deviceNo,
		device:      device,
		name:        "OnStepX Telescope",
		description: "ASCOM Alpaca OnStepX proxy",
		info:        "A Go-based ASCOM Alpaca proxy driver for OnStepX mounts.",
		apiVersion:  1,
	}
}

type apiHandlers struct {
	GetHandler http.HandlerFunc
	PutHandler http.HandlerFunc
}

func (api *TelescopeAPI) SetupRoutes() {

	// Redirects for ASCOM client setup requests
	http.HandleFunc(fmt.Sprintf("/setup/v1/telescope/%d/setup", api.deviceNo), func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/setup", http.StatusFound) })

	// Telecope V4
	handlers := map[string]apiHandlers{
		// Common methods
		"action":           apiHandlers{nil, api.PutAction},
		"commandblind":     apiHandlers{nil, api.PutCommandBlind},
		"commandbool":      apiHandlers{nil, api.PutCommandBool},
		"commandstring":    apiHandlers{nil, api.PutCommandString},
		"connect":          apiHandlers{nil, api.PutConnect},
		"connected":        apiHandlers{api.GetConnected, api.PutConnected},
		"connecting":       apiHandlers{api.GetConnecting, nil},
		"description":      apiHandlers{api.GetDescription, nil},
		"devicestate":      apiHandlers{api.GetDeviceState, nil},
		"disconnect":       apiHandlers{nil, api.PutDisconnect},
		"driverinfo":       apiHandlers{api.GetDriverInfo, nil},
		"driverversion":    apiHandlers{api.GetDriverVersion, nil},
		"interfaceversion": apiHandlers{api.GetInterfaceVersion, nil},
		"name":             apiHandlers{api.GetName, nil},
		"supportedactions": apiHandlers{api.GetSupportedActions, nil},
		// Telescope properties
		"alignmentmode":            apiHandlers{api.GetAlignmentMode, nil},
		"altitude":                 apiHandlers{api.GetAltitude, nil},
		"aperturearea":             apiHandlers{api.GetApertureArea, nil},
		"aperturediameter":         apiHandlers{api.GetApertureDiameter, nil},
		"athome":                   apiHandlers{api.GetAtHome, nil},
		"atpark":                   apiHandlers{api.GetAtPark, nil},
		"azimuth":                  apiHandlers{api.GetAzimuth, nil},
		"canfindhome":              apiHandlers{api.GetCanFindHome, nil},
		"canpark":                  apiHandlers{api.GetCanPark, nil},
		"canpulseguide":            apiHandlers{api.GetCanPulseGuide, nil},
		"cansetdeclinationrate":    apiHandlers{api.GetCanSetDeclinationRate, nil},
		"cansetrightascensionrate": apiHandlers{api.GetCanSetRightAscensionRate, nil},
		"cansetguiderates":         apiHandlers{api.GetCanSetGuideRates, nil},
		"cansetpark":               apiHandlers{api.GetCanSetPark, nil},
		"cansetpierside":           apiHandlers{api.GetCanSetPierSide, nil},
		"cansettracking":           apiHandlers{api.GetCanSetTracking, nil},
		"canslew":                  apiHandlers{api.GetCanSlew, nil},
		"canslewaltaz":             apiHandlers{api.GetCanSlewAltAz, nil},
		"canslewaltazasync":        apiHandlers{api.GetCanSlewAltAzAsync, nil},
		"canslewasync":             apiHandlers{api.GetCanSlewAsync, nil},
		"cansync":                  apiHandlers{api.GetCanSync, nil},
		"cansyncaltaz":             apiHandlers{api.GetCanSyncAltAz, nil},
		"canunpark":                apiHandlers{api.GetCanUnpark, nil},
		"declination":              apiHandlers{api.GetDeclination, nil},
		"declinationrate":          apiHandlers{api.GetDeclinationRate, api.PutDeclinationRate},
		"doesrefraction":           apiHandlers{api.GetDoesRefraction, api.PutDoesRefraction},
		"equatorialsystem":         apiHandlers{api.GetEquatorialSystem, nil},
		"focallength":              apiHandlers{api.GetFocalLength, nil},
		"guideratedeclination":     apiHandlers{api.GetGuideRateDeclination, api.PutGuideRateDeclination},
		"guideraterightascension":  apiHandlers{api.GetGuideRateRightAscension, api.PutGuideRateRightAscension},
		"ispulseguiding":           apiHandlers{api.GetIsPulseGuiding, nil},
		"rightascension":           apiHandlers{api.GetRightAscension, nil},
		"rightascensionrate":       apiHandlers{api.GetRightAscensionRate, api.PutRightAscensionRate},
		"sideofpier":               apiHandlers{api.GetSideOfPier, api.PutSideOfPier},
		"siderealtime":             apiHandlers{api.GetSiderealTime, nil},
		"siteelevation":            apiHandlers{api.GetSiteElevation, api.PutSiteElevation},
		"sitelatitude":             apiHandlers{api.GetSiteLatitude, api.PutSiteLatitude},
		"sitelongitude":            apiHandlers{api.GetSiteLongitude, api.PutSiteLongitude},
		"slewing":                  apiHandlers{api.GetSlewing, nil},
		"slewsettletime":           apiHandlers{api.GetSlewSettleTime, api.PutSlewSettleTime},
		"targetdeclination":        apiHandlers{api.GetTargetDeclination, api.PutTargetDeclination},
		"targetrightascension":     apiHandlers{api.GetTargetRightAscension, api.PutTargetRightAscension},
		"tracking":                 apiHandlers{api.GetTracking, api.PutTracking},
		"trackingrate":             apiHandlers{api.GetTrackingRate, api.PutTrackingRate},
		"trackingrates":            apiHandlers{api.GetTrackingRates, nil},
		"utcdate":                  apiHandlers{api.GetUTCDate, api.PutUTCDate},
		// Telescope methods
		"abortslew":              apiHandlers{nil, api.PutAbortSlew},
		"axisrates":              apiHandlers{api.GetAxisRates, nil},
		"canmoveaxis":            apiHandlers{api.GetCanMoveAxis, nil},
		"destinationsideofpier":  apiHandlers{api.GetDestinationSideOfPier, nil},
		"findhome":               apiHandlers{nil, api.PutFindHome},
		"moveaxis":               apiHandlers{nil, api.PutMoveAxis},
		"park":                   apiHandlers{nil, api.PutPark},
		"pulseguide":             apiHandlers{nil, api.PutPulseGuide},
		"setpark":                apiHandlers{nil, api.PutSetPark},
		"setupdialog":            apiHandlers{nil, api.PutSetupDialog},
		"slewtoaltaz":            apiHandlers{nil, api.PutSlewToAltAz},
		"slewtoaltazasync":       apiHandlers{nil, api.PutSlewToAltAzAsync},
		"slewtocoordinates":      apiHandlers{nil, api.PutSlewToCoordinates},
		"slewtocoordinatesasync": apiHandlers{nil, api.PutSlewToCoordinatesAsync},
		"slewtotarget":           apiHandlers{nil, api.PutSlewToTarget},
		"slewtotargetasync":      apiHandlers{nil, api.PutSlewToTargetAsync},
		"synctoaltaz":            apiHandlers{nil, api.PutSyncToAltAz},
		"synctocoordinates":      apiHandlers{nil, api.PutSyncToCoordinates},
		"synctotarget":           apiHandlers{nil, api.PutSyncToTarget},
		"unpark":                 apiHandlers{nil, api.PutUnpark},
	}
	for k, v := range handlers {
		if v.GetHandler != nil {
			http.HandleFunc(fmt.Sprintf("GET /api/v1/telescope/%d/%s", api.deviceNo, k), Handler(v.GetHandler))
		}
		if v.PutHandler != nil {
			http.HandleFunc(fmt.Sprintf("PUT /api/v1/telescope/%d/%s", api.deviceNo, k), Handler(v.PutHandler))
		}
	}
}

// --- Common Device Handlers ---
func (a *TelescopeAPI) PutAction(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r, "Action")
}

func (a *TelescopeAPI) PutCommandBlind(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r, "CommandBlind")
}

func (a *TelescopeAPI) PutCommandBool(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r, "CommandBool")
}

func (a *TelescopeAPI) PutCommandString(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r, "CommandString")
}

func (a *TelescopeAPI) PutConnect(w http.ResponseWriter, r *http.Request) {
	if !a.device.IsConnected() {
		a.device.Connect()
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetConnected(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, a.device.IsConnected())
}

func (a *TelescopeAPI) PutConnected(w http.ResponseWriter, r *http.Request) {
	connect, ok := checkBoolParam(w, r, "Connected")
	if !ok {
		return
	}

	// Connect/disconnect as requested
	if connect && !a.device.IsConnected() {
		if !a.device.Connect() {
			InternalErrorResponse(w, r, "Connected failed. Please check the USB connection.")
			return
		}
	} else if !connect && a.device.IsConnected() {
		a.device.Disconnect()
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetConnecting(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, false)
}

func (a *TelescopeAPI) GetDescription(w http.ResponseWriter, r *http.Request) {
	StringResponse(w, r, a.description)
}

func (a *TelescopeAPI) GetDeviceState(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get device state
	var state = map[string]any{
		"AtHome":         false,
		"AtPark":         false,
		"IsPulseGuiding": false,
		"Slewing":        false,
		"Tracking":       false,
		"RightAscension": 0,
		"Declination":    0,
		"Altitude":       0,
		"Azimuth":        0,
		"SideOfPier":     PierSideUnknown,
		"SiderealTime":   0,
		"UTCDate":        "1900-01-01T00:00:00",
	}

	// Get status from OnStepX
	var err error
	if state["AtHome"], err = a.device.GetAtHome(); err != nil {
		InternalErrorResponse(w, r, "Failed to get home state from OnStepX")
		return
	}
	if state["AtPark"], err = a.device.GetAtPark(); err != nil {
		InternalErrorResponse(w, r, "Failed to get park state from OnStepX")
		return
	}
	if state["Tracking"], err = a.device.GetIsTracking(); err != nil {
		InternalErrorResponse(w, r, "Failed to get tracking state from OnStepX")
		return
	}
	if state["Slewing"], err = a.device.GetIsSlewing(); err != nil {
		InternalErrorResponse(w, r, "Failed to get slewing state from OnStepX")
		return
	}
	if state["IsPulseGuiding"], err = a.device.GetIsPulseGuiding(); err != nil {
		InternalErrorResponse(w, r, "Failed to get pulse guiding state from OnStepX")
		return
	}
	if state["RightAscension"], err = a.device.GetRightAscension(); err != nil {
		InternalErrorResponse(w, r, "Failed to get right ascension from OnStepX")
		return
	}
	if state["Declination"], err = a.device.GetDeclination(); err != nil {
		InternalErrorResponse(w, r, "Failed to get declination from OnStepX")
		return
	}
	if state["Altitude"], err = a.device.GetAltitude(); err != nil {
		InternalErrorResponse(w, r, "Failed to get altitude from OnStepX")
		return
	}
	if state["Azimuth"], err = a.device.GetAzimuth(); err != nil {
		InternalErrorResponse(w, r, "Failed to get azimuth from OnStepX")
		return
	}
	ps, err := a.device.GetPierSide()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get pier side from OnStepX")
		return
	}
	switch ps {
	case onstepx.PierSideEast:
		state["SideOfPier"] = PierSideEast
	case onstepx.PierSideWest:
		state["SideOfPier"] = PierSideWest
	default:
		state["SideOfPier"] = PierSideUnknown
	}
	if state["SiderealTime"], err = a.device.GetSiderealTime(); err != nil {
		InternalErrorResponse(w, r, "Failed to get sidereal time from OnStepX")
		return
	}
	utc, err := a.device.GetUTCDateTime()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get UTC date/time from OnStepX")
		return
	}
	state["UTCDate"] = utc.Format(time.RFC3339)

	// Return as an array of name/value pairs
	type StateItem struct {
		Name  string
		Value any
	}
	response := []StateItem{}
	for k, v := range state {
		response = append(response, StateItem{Name: k, Value: v})
	}
	AnyResponse(w, r, response)
}

func (a *TelescopeAPI) PutDisconnect(w http.ResponseWriter, r *http.Request) {
	if a.device.IsConnected() {
		a.device.Disconnect()
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetDriverInfo(w http.ResponseWriter, r *http.Request) {
	StringResponse(w, r, a.info)
}

func (a *TelescopeAPI) GetDriverVersion(w http.ResponseWriter, r *http.Request) {
	StringResponse(w, r, a.appVersion)
}

func (a *TelescopeAPI) GetInterfaceVersion(w http.ResponseWriter, r *http.Request) {
	IntResponse(w, r, 4)
}

func (a *TelescopeAPI) GetName(w http.ResponseWriter, r *http.Request) {
	StringResponse(w, r, a.name)
}

func (a *TelescopeAPI) GetSupportedActions(w http.ResponseWriter, r *http.Request) {
	StringListResponse(w, r, []string{})
}

// --- Telescope V4 Property Handlers ---
func (a *TelescopeAPI) GetAlignmentMode(w http.ResponseWriter, r *http.Request) {
	PropertyNotImplementedResponse(w, r, "AlignmentMode")
}

func (a *TelescopeAPI) GetAltitude(w http.ResponseWriter, r *http.Request) {
	// Check device is connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get altitude and return it
	alt, err := a.device.GetAltitude()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get altitude from OnStepX")
		return
	}
	FloatResponse(w, r, alt)
}

func (a *TelescopeAPI) GetApertureArea(w http.ResponseWriter, r *http.Request) {
	PropertyNotImplementedResponse(w, r, "ApertureArea")
}

func (a *TelescopeAPI) GetApertureDiameter(w http.ResponseWriter, r *http.Request) {
	PropertyNotImplementedResponse(w, r, "ApertureDiameter")
}

func (a *TelescopeAPI) GetAtHome(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX home status
	ah, err := a.device.GetAtHome()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get home status from OnStepX")
		return
	}
	BoolResponse(w, r, ah)
}

func (a *TelescopeAPI) GetAtPark(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX park status
	ap, err := a.device.GetAtPark()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get park status from OnStepX")
		return
	}
	BoolResponse(w, r, ap)
}

func (a *TelescopeAPI) GetAzimuth(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX azimuth
	az, err := a.device.GetAzimuth()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get azimuth from OnStepX")
		return
	}
	FloatResponse(w, r, az)
}

func (a *TelescopeAPI) GetCanFindHome(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) GetCanPark(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) GetCanPulseGuide(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) GetCanSetDeclinationRate(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) GetCanSetGuideRates(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, false)
}

func (a *TelescopeAPI) GetCanSetPark(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) GetCanSetPierSide(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, false)
}

func (a *TelescopeAPI) GetCanSetRightAscensionRate(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) GetCanSetTracking(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) GetCanSlew(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, false)
}

func (a *TelescopeAPI) GetCanSlewAltAz(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, false)
}

func (a *TelescopeAPI) GetCanSlewAltAzAsync(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) GetCanSlewAsync(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) GetCanSync(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) GetCanSyncAltAz(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, false)
}

func (a *TelescopeAPI) GetCanUnpark(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) GetDeclination(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX declination
	dec, err := a.device.GetDeclination()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get declination from OnStepX")
		return
	}
	FloatResponse(w, r, dec)
}

func (a *TelescopeAPI) GetDeclinationRate(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX declination rate offset
	offset, err := a.device.GetDeclinationRate()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get declination rate offset from OnStepX")
		return
	}

	// Adjust from OnStepX units (arcsec per sidereal second) to Alpaca units (arcsec per UTC second)
	offset *= 1.00273791
	FloatResponse(w, r, offset)
}

func (a *TelescopeAPI) PutDeclinationRate(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	offset, ok := checkFloatParam(w, r, "DeclinationRate", -1800, 1800)
	if !ok {
		return
	}

	// Check tracking rate (must be sidereal)
	if !a.checkTrackingSidereal(w, r) {
		return
	}

	// Adjust from Alpaca units (arcsec per UTC second) to OnStepX units (arcsec per sidereal second)
	offset /= 1.00273791

	// Set OnStepX declination rate offset
	if err := a.device.SetDeclinationRate(offset); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX declination rate offset")
		return
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetDoesRefraction(w http.ResponseWriter, r *http.Request) {
	BoolResponse(w, r, true)
}

func (a *TelescopeAPI) PutDoesRefraction(w http.ResponseWriter, r *http.Request) {
	PropertyReadOnlyResponse(w, r, "DoesRefraction")
}

func (a *TelescopeAPI) GetEquatorialSystem(w http.ResponseWriter, r *http.Request) {
	IntResponse(w, r, EquTopocentric)
}

func (a *TelescopeAPI) GetFocalLength(w http.ResponseWriter, r *http.Request) {
	PropertyNotImplementedResponse(w, r, "FocalLength")
}

func (a *TelescopeAPI) GetGuideRateDeclination(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX pulse guide rate
	rate, err := a.device.GetPulseGuideRate()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get guide rate from OnStepX")
		return
	}
	FloatResponse(w, r, rate)
}

func (a *TelescopeAPI) PutGuideRateDeclination(w http.ResponseWriter, r *http.Request) {
	PropertyReadOnlyResponse(w, r, "GuideRateDeclination")
}

func (a *TelescopeAPI) GetGuideRateRightAscension(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX pulse guide rate
	rate, err := a.device.GetPulseGuideRate()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get guide rate from OnStepX")
		return
	}
	FloatResponse(w, r, rate)
}

func (a *TelescopeAPI) PutGuideRateRightAscension(w http.ResponseWriter, r *http.Request) {
	PropertyReadOnlyResponse(w, r, "GuideRateRightAscension")
}

func (a *TelescopeAPI) GetIsPulseGuiding(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX status
	pg, err := a.device.GetIsPulseGuiding()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get pulse guiding status from OnStepX")
		return
	}
	BoolResponse(w, r, pg)
}

func (a *TelescopeAPI) GetRightAscension(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX right ascension
	ra, err := a.device.GetRightAscension()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get right ascension from OnStepX")
		return
	}
	FloatResponse(w, r, ra)
}

func (a *TelescopeAPI) GetRightAscensionRate(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX right ascension rate offset
	offset, err := a.device.GetRightAscensionRate()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get right ascension rate offset from OnStepX")
		return
	}

	// Adjust from OnStepX units (arcsec per sidereal second) to Alpaca units (RA second per siderial second)
	offset /= 15.0
	FloatResponse(w, r, offset)
}

func (a *TelescopeAPI) PutRightAscensionRate(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	offset, ok := checkFloatParam(w, r, "RightAscensionRate", -1800, 1800)
	if !ok {
		return
	}

	// Check tracking rate (must be sidereal)
	if !a.checkTrackingSidereal(w, r) {
		return
	}

	// Adjust from Alpaca units (RA second per siderial second) to OnStepX units (arcsec per sidereal second)
	offset *= 15.0

	// Set OnStepX right ascension rate offset
	if err := a.device.SetRightAscensionRate(offset); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX right ascension rate offset")
		return
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetSideOfPier(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX status
	ps, err := a.device.GetPierSide()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get pier side from OnStepX")
		return
	}
	switch ps {
	case onstepx.PierSideEast:
		IntResponse(w, r, PierSideEast)
	case onstepx.PierSideWest:
		IntResponse(w, r, PierSideWest)
	default:
		IntResponse(w, r, PierSideUnknown)
	}
}

func (a *TelescopeAPI) PutSideOfPier(w http.ResponseWriter, r *http.Request) {
	PropertyReadOnlyResponse(w, r, "SideOfPier")
}

func (a *TelescopeAPI) GetSiderealTime(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX sidereal time
	st, err := a.device.GetSiderealTime()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get sidereal time from OnStepX")
		return
	}
	FloatResponse(w, r, st)
}

func (a *TelescopeAPI) GetSiteElevation(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get site elevation
	elv, err := a.device.GetSiteElevation()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get SiteElevation from device")
		return
	}

	// Return result
	FloatResponse(w, r, elv)
}

func (a *TelescopeAPI) PutSiteElevation(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	elv, ok := checkFloatParam(w, r, "SiteElevation", -300, 10000)
	if !ok {
		return
	}

	// Set OnStepX site elevation
	if err := a.device.SetSiteElevation(elv); err != nil {
		InternalErrorResponse(w, r, "Failed to set SiteElevation on device")
		return
	}

	// Success
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetSiteLatitude(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get site elevation
	lat, err := a.device.GetSiteLatitude()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get SiteLatitude from device")
		return
	}

	// Return result
	FloatResponse(w, r, lat)
}

func (a *TelescopeAPI) PutSiteLatitude(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	lat, ok := checkFloatParam(w, r, "SiteLatitude", -90, 90)
	if !ok {
		return
	}

	// Set OnStepX site latitude
	if err := a.device.SetSiteLatitude(lat); err != nil {
		InternalErrorResponse(w, r, "Failed to set SiteLatitude on device")
		return
	}

	// Success
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetSiteLongitude(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get site elevation
	long, err := a.device.GetSiteLongitude()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get SiteLongitude from device")
		return
	}

	// Return result
	FloatResponse(w, r, long)
}

func (a *TelescopeAPI) PutSiteLongitude(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	long, ok := checkFloatParam(w, r, "SiteLongitude", -180, 180)
	if !ok {
		return
	}

	// Set OnStepX site longitude
	if err := a.device.SetSiteLongitude(long); err != nil {
		InternalErrorResponse(w, r, "Failed to set SiteLongitude on device")
		return
	}

	// Success
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetSlewing(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Check if OnStepX is slewing, guiding (not pulse guiding), homing or parking
	isSlewing, err := a.device.GetIsSlewing()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get slewing status from OnStepX")
		return
	}
	if !isSlewing {
		isSlewing, err = a.device.GetIsGuiding()
		if err != nil {
			InternalErrorResponse(w, r, "Failed to get guiding status from OnStepX")
			return
		}
	}
	if !isSlewing {
		isSlewing, err = a.device.GetIsHoming()
		if err != nil {
			InternalErrorResponse(w, r, "Failed to get homing status from OnStepX")
			return
		}
	}
	if !isSlewing {
		isSlewing, err = a.device.GetIsParking()
		if err != nil {
			InternalErrorResponse(w, r, "Failed to get parking status from OnStepX")
			return
		}
	}
	BoolResponse(w, r, isSlewing)
}

func (a *TelescopeAPI) GetSlewSettleTime(w http.ResponseWriter, r *http.Request) {
	PropertyNotImplementedResponse(w, r, "SlewSettleTime")
}

func (a *TelescopeAPI) PutSlewSettleTime(w http.ResponseWriter, r *http.Request) {
	PropertyNotImplementedResponse(w, r, "SlewSettleTime")
}

func (a *TelescopeAPI) GetTargetDeclination(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX target declination
	dec, err := a.device.GetTargetDeclination()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get target declination from OnStepX")
		return
	}
	FloatResponse(w, r, dec)
}

func (a *TelescopeAPI) PutTargetDeclination(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	dec, ok := checkFloatParam(w, r, "TargetDeclination", -90, 90)
	if !ok {
		return
	}

	// Set OnStepX target declination
	if err := a.device.SetTargetDeclination(dec); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX target declination")
		return
	}

	// Success
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetTargetRightAscension(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX target right ascension
	ra, err := a.device.GetTargetRightAscension()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get target right ascension from OnStepX")
		return
	}
	FloatResponse(w, r, ra)
}

func (a *TelescopeAPI) PutTargetRightAscension(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	ra, ok := checkFloatParam(w, r, "TargetRightAscension", 0, 24)
	if !ok {
		return
	}

	// Set OnStepX target right ascension
	if err := a.device.SetTargetRightAscension(ra); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX target right ascension")
		return
	}

	// Success
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetTracking(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get tracking status from OnStepX
	trk, err := a.device.GetIsTracking()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get tracking status from OnStepX")
		return
	}
	BoolResponse(w, r, trk)
}

func (a *TelescopeAPI) PutTracking(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	trk, ok := checkBoolParam(w, r, "Tracking")
	if !ok {
		return
	}

	// Check if mount is parked
	if trk && !a.checkNotParked(w, r) {
		return
	}

	// Set OnStepX tracking state
	err := a.device.SetTracking(trk)
	if err != nil {
		InternalErrorResponse(w, r, "Failed to start/stop tracking")
		return
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetTrackingRate(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX tracking rate
	rate, err := a.device.GetTrackingRate()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get tracking rate from OnStepX")
		return
	}
	switch rate {
	case "Sidereal":
		IntResponse(w, r, TrackingSidereal)
	case "Lunar":
		IntResponse(w, r, TrackingLunar)
	case "Solar":
		IntResponse(w, r, TrackingSolar)
	case "King":
		IntResponse(w, r, TrackingKing)
	default:
		InternalErrorResponse(w, r, fmt.Sprintf("Invalid tracking rate '%s' returned by OnStepX", rate))
	}
}

func (a *TelescopeAPI) PutTrackingRate(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	rate, ok := checkIntParam(w, r, "TrackingRate", TrackingSidereal, TrackingKing)
	if !ok {
		return
	}

	// Set OnStepX tracking rate
	rateStr := ""
	switch rate {
	case TrackingSidereal:
		rateStr = "Sidereal"
	case TrackingLunar:
		rateStr = "Lunar"
	case TrackingSolar:
		rateStr = "Solar"
	case TrackingKing:
		rateStr = "King"
	}
	if err := a.device.SetTrackingRate(rateStr); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX tracking rate")
		return
	}

	// Success
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) GetTrackingRates(w http.ResponseWriter, r *http.Request) {
	AnyResponse(w, r, [...]int{0, 1, 2, 3})
}

func (a *TelescopeAPI) GetUTCDate(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Get OnStepX UTC date/time
	utc, err := a.device.GetUTCDateTime()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get UTC date/time from OnStepX")
		return
	}
	StringResponse(w, r, utc.Format(time.RFC3339))
}

func (a *TelescopeAPI) PutUTCDate(w http.ResponseWriter, r *http.Request) {
	PropertyReadOnlyResponse(w, r, "UTCDate")
}

// --- Telescope V4 Method Handlers ---
func (a *TelescopeAPI) PutAbortSlew(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Check if mount is parked
	if !a.checkNotParked(w, r) {
		return
	}

	// Abort OnStepX slew
	if err := a.device.AbortSlew(); err != nil {
		InternalErrorResponse(w, r, "Failed to abort OnStepX slew")
		return
	}
	SuccessResponse(w, r)
}

type Rates struct {
	Minimum float64 `json:"Minimum"`
	Maximum float64 `json:"Maximum"`
}

func (a *TelescopeAPI) GetAxisRates(w http.ResponseWriter, r *http.Request) {
	// Validate parameters
	axis, ok := checkIntParam(w, r, "Axis", 0, 2)
	if !ok {
		return
	}

	// Return 0-4 deg/sec
	if axis == 0 || axis == 1 {
		AnyResponse(w, r, [...]Rates{{Minimum: 0.0, Maximum: 4.0}})
	} else {
		AnyResponse(w, r, [...]Rates{})
	}
}

func (a *TelescopeAPI) GetCanMoveAxis(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	axis, ok := checkIntParam(w, r, "Axis", 0, 2)
	if !ok {
		return
	}

	// Can move right ascension and declination axes
	if axis == 0 || axis == 1 {
		BoolResponse(w, r, true)
	} else {
		BoolResponse(w, r, false)
	}
}

func (a *TelescopeAPI) GetDestinationSideOfPier(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	ra, ok := checkFloatParam(w, r, "RightAscension", 0, 24)
	if !ok {
		return
	}
	dec, ok := checkFloatParam(w, r, "Declination", -90, 90)
	if !ok {
		return
	}

	// Save current OnStepX target and set it to request position
	origRa, errRa := a.device.GetTargetRightAscension()
	origDec, errDec := a.device.GetTargetDeclination()
	if errRa != nil || errDec != nil {
		InternalErrorResponse(w, r, "Failed to save OnStepX target position")
		return
	}
	errRa = a.device.SetTargetRightAscension(ra)
	errDec = a.device.SetTargetDeclination(dec)
	if errRa != nil || errDec != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX target position")
		return
	}

	// Get destination side of pier and restore original target position
	ps, err := a.device.GetTargetPierSide()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get OnStepX destination pier side")
		return
	}
	errRa = a.device.SetTargetRightAscension(origRa)
	errDec = a.device.SetTargetDeclination(origDec)
	if errRa != nil || errDec != nil {
		InternalErrorResponse(w, r, "Failed to restore OnStepX target position")
		return
	}
	switch ps {
	case onstepx.PierSideEast:
		IntResponse(w, r, PierSideEast)
	case onstepx.PierSideWest:
		IntResponse(w, r, PierSideWest)
	default:
		IntResponse(w, r, PierSideUnknown)
	}
}

func (a *TelescopeAPI) PutFindHome(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Check if mount is parked or slewing (abort slew if necessary)
	if !a.checkNotParked(w, r) {
		return
	}
	if !a.checkNotSlewing(w, r, true) {
		return
	}

	// Start move to home position
	if err := a.device.StartHome(); err != nil {
		InternalErrorResponse(w, r, "Failed to start OnStepX move to home position")
		return
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) PutMoveAxis(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	axis, ok := checkIntParam(w, r, "Axis", 0, 1)
	if !ok {
		return
	}
	rate, ok := checkFloatParam(w, r, "Rate", -4, 4)
	if !ok {
		return
	}

	// Check if mount is parked
	if !a.checkNotParked(w, r) {
		return
	}

	// Start moving axis at requested rate
	switch axis {
	case 0:
		if rate == 0 {
			err := a.device.StopRightAscensionSlew()
			if err != nil {
				InternalErrorResponse(w, r, "Failed to abort right ascension slew")
				return
			}
		} else if rate > 0 {
			err := a.device.SetRightAscensionSlewRate(rate)
			if err != nil {
				InternalErrorResponse(w, r, "Failed to set right ascension slew rate")
				return
			}
			err = a.device.StartRightAscensionSlewEast()
			if err != nil {
				InternalErrorResponse(w, r, "Failed to start right ascension slew")
				return
			}
		} else {
			err := a.device.SetRightAscensionSlewRate(-rate)
			if err != nil {
				InternalErrorResponse(w, r, "Failed to set right ascension slew rate")
				return
			}
			err = a.device.StartRightAscensionSlewWest()
			if err != nil {
				InternalErrorResponse(w, r, "Failed to start right ascension slew")
				return
			}
		}
	case 1:
		if rate == 0 {
			err := a.device.StopDeclinationSlew()
			if err != nil {
				InternalErrorResponse(w, r, "Failed to abort declination slew")
				return
			}
		} else if rate > 0 {
			err := a.device.SetDeclinationSlewRate(rate)
			if err != nil {
				InternalErrorResponse(w, r, "Failed to set declination slew rate")
				return
			}
			err = a.device.StartDeclinationSlewNorth()
			if err != nil {
				InternalErrorResponse(w, r, "Failed to start rdeclination slew")
				return
			}
		} else {
			err := a.device.SetDeclinationSlewRate(-rate)
			if err != nil {
				InternalErrorResponse(w, r, "Failed to set declination slew rate")
				return
			}
			err = a.device.StartDeclinationSlewSouth()
			if err != nil {
				InternalErrorResponse(w, r, "Failed to start rdeclination slew")
				return
			}
		}
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) PutPark(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Check if already parked (or parking)
	atPark, err := a.device.GetAtPark()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get park state")
		return
	}
	isParking, err := a.device.GetIsParking()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get parking state")
		return
	}

	// Start move to park position (if not already parked/parking)
	if !atPark && !isParking {
		if err := a.device.StartPark(); err != nil {
			InternalErrorResponse(w, r, "Failed to start OnStepX park operation")
			return
		}
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) PutPulseGuide(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Check if mount is parked or slewing
	if !a.checkNotParked(w, r) {
		return
	}
	if !a.checkNotSlewing(w, r, false) {
		return
	}

	// Validate parameters
	dirn, ok := checkIntParam(w, r, "Direction", 0, 3)
	if !ok {
		return
	}
	dirnStr := ""
	switch dirn {
	case 0:
		dirnStr = "n"
	case 1:
		dirnStr = "s"
	case 2:
		dirnStr = "e"
	case 3:
		dirnStr = "w"
	default:
		BadRequestResponse(w, r, "Invalid direction")
	}
	durn, ok := checkIntParam(w, r, "Duration", 0, 5000)
	if !ok {
		return
	}

	// Start pulse guiding (if duration not zero)
	if durn > 0 {
		err := a.device.StartPulseGuide(dirnStr, int(durn))
		if err != nil {
			InternalErrorResponse(w, r, "Failed to start pulse guide")
			return
		}
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) PutSetPark(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Check if parked (or parking), slewing or tracking
	if !a.checkNotParked(w, r) || !a.checkNotTracking(w, r) || !a.checkNotTracking(w, r) {
		return
	}

	// Set OnStepX park position
	if err := a.device.SetPark(); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX park position")
		return
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) PutSetupDialog(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r, "SetupDialog")
}

func (a *TelescopeAPI) PutSlewToAltAz(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r, "SlewToAltAz")
}

func (a *TelescopeAPI) PutSlewToAltAzAsync(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	alt, ok := checkFloatParam(w, r, "Altitude", -90, 90)
	if !ok {
		return
	}
	az, ok := checkFloatParam(w, r, "Azimuth", 0, 360)
	if !ok {
		return
	}

	// Check if mount is parked or slewing (abort slew if necessary)
	if !a.checkNotParked(w, r) {
		return
	}
	if !a.checkNotSlewing(w, r, true) {
		return
	}

	// Set OnStepX target position
	if err := a.device.SetTargetAltitude(alt); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX target position")
		return
	}
	if err := a.device.SetTargetAzimuth(az); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX target position")
		return
	}

	// Slew to target position
	if err := a.device.StartSlewToTargetAltAz(); err != nil {
		InternalErrorResponse(w, r, "Failed to start slew to target")
		return
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) PutSlewToCoordinates(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r, "SlewToCoordinates")
}

func (a *TelescopeAPI) PutSlewToCoordinatesAsync(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	ra, ok := checkFloatParam(w, r, "RightAscension", 0, 24)
	if !ok {
		return
	}
	dec, ok := checkFloatParam(w, r, "Declination", -90, 90)
	if !ok {
		return
	}

	// Check if mount is parked or slewing (abort slew if necessary)
	if !a.checkNotParked(w, r) {
		return
	}
	if !a.checkNotSlewing(w, r, true) {
		return
	}

	// Set OnStepX target position
	if err := a.device.SetTargetRightAscension(ra); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX target position")
		return
	}
	if err := a.device.SetTargetDeclination(dec); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX target position")
		return
	}

	// Start slew to target position
	if err := a.device.StartSlewToTarget(); err != nil {
		InternalErrorResponse(w, r, "Failed to start slew to target")
		return
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) PutSlewToTarget(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r, "SlewToTarget")
}

func (a *TelescopeAPI) PutSlewToTargetAsync(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Check if mount is parked or slewing (abort slew if necessary)
	if !a.checkNotParked(w, r) {
		return
	}
	if !a.checkNotSlewing(w, r, true) {
		return
	}

	// Start slew to target position
	if err := a.device.StartSlewToTarget(); err != nil {
		InternalErrorResponse(w, r, "Failed to start slew to target")
		return
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) PutSyncToAltAz(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r, "SyncToAltAz")
}

func (a *TelescopeAPI) PutSyncToCoordinates(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Validate parameters
	ra, ok := checkFloatParam(w, r, "RightAscension", 0, 24)
	if !ok {
		return
	}
	dec, ok := checkFloatParam(w, r, "Declination", -90, 90)
	if !ok {
		return
	}

	// Check if mount is parked or slewing
	if !a.checkNotParked(w, r) {
		return
	}
	if !a.checkNotSlewing(w, r, false) {
		return
	}

	// Set OnStepX target position
	if err := a.device.SetTargetRightAscension(ra); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX target position")
		return
	}
	if err := a.device.SetTargetDeclination(dec); err != nil {
		InternalErrorResponse(w, r, "Failed to set OnStepX target position")
		return
	}

	// Sync to target position
	if err := a.device.SyncToTarget(); err != nil {
		InternalErrorResponse(w, r, "Failed to sync to target position")
		return
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) PutSyncToTarget(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Check if mount is parked or slewing
	if !a.checkNotParked(w, r) {
		return
	}
	if !a.checkNotSlewing(w, r, false) {
		return
	}

	// Sync to target position
	if err := a.device.SyncToTarget(); err != nil {
		InternalErrorResponse(w, r, "Failed to sync to target position")
		return
	}
	SuccessResponse(w, r)
}

func (a *TelescopeAPI) PutUnpark(w http.ResponseWriter, r *http.Request) {
	// Check device connected
	if !a.checkConnected(w, r) {
		return
	}

	// Restore mount to operation if parked or at home position
	atPark, err := a.device.GetAtPark()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get park state")
		return
	}
	atHome, err := a.device.GetAtHome()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get home state")
		return
	}
	if atPark || atHome {
		if err := a.device.Unpark(); err != nil {
			InternalErrorResponse(w, r, "Failed to unpark OnStepX")
			return
		}
	}
	SuccessResponse(w, r)
}

// Helpers
func (a *TelescopeAPI) checkConnected(w http.ResponseWriter, r *http.Request) bool {
	if !a.device.IsConnected() {
		ErrorResponse(w, r, 0x407, "Device not connected")
		return false
	}
	return true
}

func (a *TelescopeAPI) checkNotParked(w http.ResponseWriter, r *http.Request) bool {
	atPark, err := a.device.GetAtPark()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get parked status")
		return false
	}
	isParking, err := a.device.GetIsParking()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get parking status")
		return false
	}
	if atPark || isParking {
		ErrorResponse(w, r, 0x0408, "Invalid while parked/parking")
		return false
	}
	return true
}

func (a *TelescopeAPI) checkNotTracking(w http.ResponseWriter, r *http.Request) bool {
	isTracking, err := a.device.GetIsTracking()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get tracking status")
		return false
	}
	if isTracking {
		ErrorResponse(w, r, 0x040B, "Invalid while tracking")
		return false
	}
	return true
}

func (a *TelescopeAPI) checkNotSlewing(w http.ResponseWriter, r *http.Request, abort bool) bool {

	// Check if slewing
	isSlewing, err := a.device.GetIsSlewing()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get slewing status")
		return false
	}
	if isSlewing {
		if !abort {
			ErrorResponse(w, r, 0x040B, "Invalid while slewing")
			return false
		}

		// Abort slew
		if err := a.device.AbortSlew(); err != nil {
			InternalErrorResponse(w, r, "Failed to abort slew")
			return false
		}
		for range 200 {
			isSlewing, err = a.device.GetIsSlewing()
			if err != nil {
				InternalErrorResponse(w, r, "Failed to get OnStepX slew status")
				return false
			}
			if !isSlewing {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
	return !isSlewing
}

func (a *TelescopeAPI) checkTrackingSidereal(w http.ResponseWriter, r *http.Request) bool {

	// Check tgracking rate (must be sidereal)
	trkRate, err := a.device.GetTrackingRate()
	if err != nil {
		InternalErrorResponse(w, r, "Failed to get OnStepX tracking rate")
		return false
	}
	if trkRate != "Sidereal" {
		ErrorResponse(w, r, 0x040B, "Only allowed when tracking is sidereal")
		return false
	}
	return true
}
