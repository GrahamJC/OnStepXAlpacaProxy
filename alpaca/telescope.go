package alpaca

import (
	"fmt"
	"net/http"
	"strconv"

	"onstepx-alpaca-proxy/onstepx"
)

// API holds all dependencies for the Alpaca API handlers.
type TelescopeAPI struct {
	appVersion string
	deviceNo   int
	device     onstepx.OnStepXDevice
}

// NewAPI creates a new telescope API instance.
func NewTelescopeAPI(appVersion string, deviceNo int, device onstepx.OnStepXDevice) *TelescopeAPI {
	return &TelescopeAPI{
		appVersion: appVersion,
		deviceNo:   deviceNo,
		device:     device,
	}
}

func (api *TelescopeAPI) SetupRoutes() {

	// Redirects for ASCOM client setup requests
	http.HandleFunc(fmt.Sprintf("/setup/v1/telescope/%d/setup", api.deviceNo), func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/setup", http.StatusFound) })

	// Telecope V4
	handlers := map[string]http.HandlerFunc{
		"name":                     api.HandleDeviceName("OnStepX Telescope"),
		"description":              api.HandleDeviceDescription,
		"driverinfo":               api.HandleDriverInfo,
		"driverversion":            api.HandleDriverVersion,
		"connected":                api.HandleConnected,
		"connecting":               api.HandleConnecting,
		"interfaceversion":         api.HandleInterfaceVersion,
		"supportedactions":         api.HandleSupportedActions,
		"action":                   api.HandleAction,
		"abortslew":                api.HandleAbortSlew,
		"axisRates":                api.HandleAxisRates,
		"canmoveaxis":              api.HandleCanMoveAxis,
		"commandblind":             api.HandleCommandBlind,
		"commandbool":              api.HandleCommandBool,
		"commandstring":            api.HandleCommandString,
		"connect":                  api.HandleConnect,
		"destinationsideofpier":    api.HandleDestinationSideOfPier,
		"disconnect":               api.HandleDisconnect,
		"park":                     api.HandlePark,
		"pulseguide":               api.HandlePulseGuide,
		"setpark":                  api.HandleSetPark,
		"setupdialog":              api.HandleSetupDialog,
		"slewtoaltaz":              api.HandleSlewToAltAz,
		"slewtoaltazasync":         api.HandleSlewToAltAzAsynch,
		"slewtocoordinates":        api.HandleSlewToCoordinates,
		"slewtocoordinatesasync":   api.HandleSlewToCoordinatesAsync,
		"slewtotarget":             api.HandleSlewToTarget,
		"slewtotargetasync":        api.HandleSlewToTargetAsync,
		"synctoaltaz":              api.HandleSyncToAltAz,
		"synctocoordinates":        api.HandleSyncToCoordinates,
		"synctotarget":             api.HandleSyncToTarget,
		"unpark":                   api.HandleUnpark,
		"alignmentmode":            api.HandleAlignmentMode,
		"altitude":                 api.HandleAltitude,
		"aperturearea":             api.HandleApertureArea,
		"aperturediameter":         api.HandleApertureDiameter,
		"athome":                   api.HandleAtHome,
		"atpark":                   api.HandleAtPark,
		"azimuth":                  api.HandleAzimuth,
		"canfindhome":              api.HandleCanFindHome,
		"canpark":                  api.HandleCanPark,
		"canpulseguid":             api.HandleCanPulseGuide,
		"cnsetpark":                api.HandleCanSetPark,
		"cansetpierside":           api.HandleCanSetPierSide,
		"cansetrightascensionrate": api.HandleCanSetRightAscensionRate,
		"cansettracking":           api.HandleCanSetTracking,
		"canslew":                  api.HandleCanSlew,
		"canslewaltaz":             api.HandleCanSlewAltAz,
		"canslewaltazasync":        api.HandleCanSlewAltAzAsync,
		"canslewasync":             api.HandleCanSlewAsync,
		"canssync":                 api.HandleCanSync,
		"cansyncaltaz":             api.HandleCanSyncAltAz,
		"canunpark":                api.HandleCanUnpark,
		"decelination":             api.HandleDeclination,
		"declinationrate":          api.HandleDeclinationRate,
		"devicestate":              api.HandleDeviceState,
		"doesrefraction":           api.HandleDoesRefraction,
		"equatorialsystem":         api.HandleEquatorialSystem,
		"focallength":              api.HandleFocalLength,
		"guideratedeclination":     api.HandleGuideRateDeclination,
		"guideraterightascension":  api.HandleGuideRateRightAscension,
		"sideofpier":               api.HandleSideOfPier,
		"siderealtime":             api.HandleSiderealTime,
		"siteelevation":            api.HandleSiteElevation,
		"sitelatitude":             api.HandleSiteLatitude,
		"sitelongitude":            api.HandleSiteLongitude,
		"slewsettletime":           api.HandleSlewSettleTime,
		"slewing":                  api.HandleSlewing,
		"targetdeclination":        api.HandleTargetDeclination,
		"targetrightascension":     api.HandleTargetRightAscension,
		"tracking":                 api.HandleTracking,
		"trackingrate":             api.HandleTrackingRate,
		"utcdate":                  api.HandleUTCDate,
	}
	for k, v := range handlers {
		http.HandleFunc(fmt.Sprintf("/api/v1/telescope/%d/%s", api.deviceNo, k), Handler(v))
	}
}

// --- Common Device Handlers ---

func (a *TelescopeAPI) HandleDeviceName(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		StringResponse(w, r, name)
	}
}

func (a *TelescopeAPI) HandleDeviceDescription(w http.ResponseWriter, r *http.Request) {
	StringResponse(w, r, "OnStepX Proxy Driver")
}

func (a *TelescopeAPI) HandleDriverInfo(w http.ResponseWriter, r *http.Request) {
	StringResponse(w, r, "A Go-based ASCOM Alpaca proxy driver for OnStepX mounts.")
}

func (a *TelescopeAPI) HandleDriverVersion(w http.ResponseWriter, r *http.Request) {
	StringResponse(w, r, a.appVersion)
}

func (a *TelescopeAPI) HandleInterfaceVersion(w http.ResponseWriter, r *http.Request) {
	IntResponse(w, r, 1) // Switch and ObsCond are both Interface Version 1
}

func (a *TelescopeAPI) HandleConnected(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		connectedStr, ok := GetFormValueIgnoreCase(r, "Connected")
		if !ok {
			ErrorResponse(w, r, http.StatusOK, 0x400, "Missing Connected parameter for PUT request")
			return
		}
		connected, err := strconv.ParseBool(connectedStr)
		if err != nil {
			ErrorResponse(w, r, http.StatusOK, 0x400, fmt.Sprintf("Invalid value for Connected: '%s'", connectedStr))
			return
		}
		// When client tries to connect, verify hardware is available
		if connected && !a.device.IsConnected() {
			ErrorResponse(w, r, http.StatusOK, 0x400, "OnStepX device not connected. Please check the USB connection.")
			return
		}
		// The connection is managed automatically, so we just acknowledge.
		EmptyResponse(w, r)
		return
	}
	// For GET, report the actual connection status.
	BoolResponse(w, r, a.device.IsConnected())
}

func (a *TelescopeAPI) HandleConnecting(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSupportedActions(w http.ResponseWriter, r *http.Request) {
	StringListResponse(w, r, []string{})
}

func (a *TelescopeAPI) HandleAction(w http.ResponseWriter, r *http.Request) {
	action, ok := GetFormValueIgnoreCase(r, "Action")
	if !ok {
		ErrorResponse(w, r, http.StatusOK, 0x400, "Missing Action parameter")
		return
	}

	ErrorResponse(w, r, http.StatusOK, 0x400, fmt.Sprintf("Action '%s' is not supported.", action))
}

// --- Telescope V4 Handlers ---
func (a *TelescopeAPI) HandleAbortSlew(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleAxisRates(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanMoveAxis(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCommandBlind(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCommandBool(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCommandString(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleConnect(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleDestinationSideOfPier(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleDisconnect(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandlePark(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandlePulseGuide(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSetPark(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSetupDialog(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSlewToAltAz(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSlewToAltAzAsynch(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSlewToCoordinates(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSlewToCoordinatesAsync(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSlewToTarget(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSlewToTargetAsync(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSyncToAltAz(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSyncToCoordinates(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSyncToTarget(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleUnpark(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

// Properties
func (a *TelescopeAPI) HandleAlignmentMode(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleAltitude(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleApertureArea(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleApertureDiameter(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleAtHome(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleAtPark(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleAzimuth(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanFindHome(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanPark(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanPulseGuide(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSetDeclinationRates(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSetGuideRate(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSetPark(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSetPierSide(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSetRightAscensionRate(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSetTracking(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSlew(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSlewAltAz(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSlewAltAzAsync(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSlewAsync(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSync(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanSyncAltAz(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleCanUnpark(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleDeclination(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleDeclinationRate(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleDeviceState(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleDoesRefraction(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleEquatorialSystem(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleFocalLength(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleGuideRateDeclination(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleGuideRateRightAscension(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleIsPulseGuiding(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleRightAscension(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleRightAscensionRate(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSideOfPier(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSiderealTime(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSiteElevation(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSiteLatitude(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSiteLongitude(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSlewSettleTime(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleSlewing(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleTargetDeclination(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleTargetRightAscension(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleTracking(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}
func (a *TelescopeAPI) HandleTrackingRate(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleTrackingRates(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}

func (a *TelescopeAPI) HandleUTCDate(w http.ResponseWriter, r *http.Request) {
	NotImplementedResponse(w, r)
}
