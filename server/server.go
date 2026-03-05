package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"

	"onstepx-alpaca-proxy/alpaca"
	"onstepx-alpaca-proxy/config"
	"onstepx-alpaca-proxy/onstepx"
)

var (
	onstepxDevice onstepx.OnStepXDevice
	telescopeApi  *alpaca.TelescopeAPI
)

// Start initializes and starts the HTTP server serving the UI from the provided filesystem.
func Start(uiFS fs.FS, appVersion string, device onstepx.OnStepXDevice) error {

	// Save OnStepX device
	onstepxDevice = device
	telescopeApi = alpaca.NewTelescopeAPI(appVersion, 0, device)

	// Setup router
	alpaca.SetupManagementHandlers(appVersion)
	telescopeApi.SetupRoutes()
	setupUIRoutes(uiFS, appVersion)

	// Bind listener
	conf := config.Get()
	addr := fmt.Sprintf("%s:%d", conf.ListenAddress, conf.NetworkPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("FATAL: Could not bind to address UI port. Please check your configuration.", "addr", addr, "error", err)
		return err
	}

	// Run server
	slog.Info("Starting Alpaca API server", "addr", addr)
	return http.Serve(listener, nil)
}

func setupUIRoutes(uiFS fs.FS, appVersion string) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/setup" {
			// Serve the SPA entry point
			http.ServeFileFS(w, r, uiFS, "index.html")
		} else {
			// Serve static assets
			http.FileServer(http.FS(uiFS)).ServeHTTP(w, r)
		}
	})
}

// UI support handlers
func handleGetProxyVersion(appVersion string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Version string `json:"version"`
		}{
			Version: appVersion,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func handleGetFirmwareVersion(w http.ResponseWriter, r *http.Request) {

	// Must be connected
	if !onstepxDevice.IsConnected() {
		http.Error(w, "Device not connected", http.StatusServiceUnavailable)
	}

	// Get version from OnStepX device
	version, err := onstepxDevice.GetFirmwareVersion()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Return version
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Version string `json:"version"`
	}{
		Version: version,
	}
	json.NewEncoder(w).Encode(response)
}

func handleGetPosition(w http.ResponseWriter, r *http.Request) {

	var err error
	response := struct {
		RA  float64 `json:"ra"`
		Dec float64 `json:"dec"`
		Alt float64 `json:"alt"`
		Az  float64 `json:"az"`
	}{}

	// Must be connected
	if !onstepxDevice.IsConnected() {
		http.Error(w, "Device not connected", http.StatusServiceUnavailable)
	}

	// Get current position from OnStepX device
	response.RA, err = onstepxDevice.GetRightAscension()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response.Dec, err = onstepxDevice.GetDeclination()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response.Alt, err = onstepxDevice.GetAltitude()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response.Az, err = onstepxDevice.GetAzimuth()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Return position
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/*
func handleGetStatus(w http.ResponseWriter, r *http.Request) {

	// Get version from OnStepX device
	version := ""
	if onstepxDevice.IsConnected() {
		version = onstepxDevice.GetFirmwareVersion()
	}

	// Return OnStepX device status
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		IsHome     bool `json:"isHome"`
		IsParked   bool `json:"isParked"`
		IsSlewing  bool `json:"isSlewing"`
		IsTracking bool `json:"isTracking"`
		Position   struct {
			RA  float64 `json:"ra"`
			Dec float64 `json:"dec"`
			Alt float64 `json:"alt"`
			Az  float64 `json:"az"`
		} `json:"position"`
		TrackingRate int8 `json:"trackingRate"`
		GuideRate    float64 `json:"guideRate"`
		MaxSlewSpeed float64 `json:"maxSlewSpeed"`
		GotoRate     float64 `json:"gotoRate"`
	}{
		IsHome:     true,
		IsParked:   true,
		IsSlewing:  false,
		IsTracking: false,
	}
	json.NewEncoder(w).Encode(response)
}
*/

/*
func setupAlpacaDeviceRoutes(api *alpaca.API) {
	// Redirects for ASCOM client setup requests
	http.HandleFunc("/setup/v1/switch/0/setup", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/setup", http.StatusFound) })
	http.HandleFunc("/setup/v1/observingconditions/0/setup", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/setup", http.StatusFound) })

	// Common handlers
	commonHandlers := map[string]http.HandlerFunc{
		"description":      api.HandleDeviceDescription,
		"driverinfo":       api.HandleDriverInfo,
		"driverversion":    api.HandleDriverVersion,
		"connected":        api.HandleConnected,
		"interfaceversion": api.HandleInterfaceVersion,
	}

	// Switch device
	switchHandlers := map[string]http.HandlerFunc{
		"maxswitch":            api.HandleSwitchMaxSwitch,
		"getswitchname":        api.HandleSwitchGetSwitchName,
		"setswitchname":        api.HandleSwitchSetSwitchName,
		"canwrite":             api.HandleSwitchCanWrite,
		"getswitch":            api.HandleSwitchGetSwitch,
		"getswitchvalue":       api.HandleSwitchGetSwitchValue,
		"setswitchvalue":       api.HandleSwitchSetSwitchValue,
		"setswitch":            api.HandleSwitchSetSwitchValue, // Alias
		"getswitchdescription": api.HandleSwitchGetSwitchDescription,
		"maxswitchvalue":       api.HandleSwitchMaxSwitchValue,
		"minswitchvalue":       api.HandleSwitchMinSwitchValue,
		"switchstep":           api.HandleSwitchSwitchStep,
		"name":                 api.HandleDeviceName("SV241 Power Switch"),
		"supportedactions":     api.HandleSwitchSupportedActions,
		"action":               api.HandleSwitchAction,
	}
	for k, v := range commonHandlers {
		switchHandlers[k] = v
	}
	http.HandleFunc("/api/v1/switch/0/", alpaca.Handler(deviceMux(switchHandlers, api)))

	// ObservingConditions device
	obsCondHandlers := map[string]http.HandlerFunc{
		"temperature":         api.HandleObsCondTemperature,
		"humidity":            api.HandleObsCondHumidity,
		"dewpoint":            api.HandleObsCondDewPoint,
		"name":                api.HandleDeviceName("SV241 Environment"),
		"supportedactions":    api.HandleSupportedActions,
		"action":              api.HandleObsCondAction,
		"averageperiod":       api.HandleObsCondAveragePeriod,
		"sensordescription":   api.HandleObsCondSensorDescription,
		"timesincelastupdate": api.HandleObsCondTimeSinceLastUpdate,
		"refresh":             api.HandleObsCondRefresh,
		"cloudcover":          api.HandleObsCondNotImplemented,
		"pressure":            api.HandleObsCondNotImplemented,
		"rainrate":            api.HandleObsCondNotImplemented,
		"skybrightness":       api.HandleObsCondNotImplemented,
		"skyquality":          api.HandleObsCondNotImplemented,
		"skytemperature":      api.HandleObsCondNotImplemented,
		"starfwhm":            api.HandleObsCondNotImplemented,
		"winddirection":       api.HandleObsCondNotImplemented,
		"windgust":            api.HandleObsCondNotImplemented,
		"windspeed":           api.HandleObsCondNotImplemented,
	}
	for k, v := range commonHandlers {
		obsCondHandlers[k] = v
	}
	http.HandleFunc("/api/v1/observingconditions/0/", alpaca.Handler(deviceMux(obsCondHandlers, api)))
}

// deviceMux creates a handler that routes to sub-handlers based on the final URL path segment.
func deviceMux(handlers map[string]http.HandlerFunc, api *alpaca.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSuffix(r.URL.Path, "/")
		lastSlash := strings.LastIndex(path, "/")
		if lastSlash == -1 {
			alpaca.ErrorResponse(w, r, http.StatusNotFound, 0x404, "Invalid URL path.")
			return
		}
		method := strings.ToLower(path[lastSlash+1:])

		if handler, ok := handlers[method]; ok {
			handler(w, r)
		} else {
			alpaca.ErrorResponse(w, r, http.StatusNotFound, 0x404, fmt.Sprintf("Method '%s' not found on this device.", method))
		}
	}
}

// --- API Handlers ---

func handleGetFirmwareConfig(w http.ResponseWriter, r *http.Request) {
	resp, err := serial.SendCommand(`{"get":"config"}`, false, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, resp)
}

func handleSetFirmwareConfig(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var js json.RawMessage
	if json.Unmarshal(body, &js) != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	command := fmt.Sprintf(`{"sc":%s}`, string(body))
	logger.Debug("Sending to device: %s", command)
	resp, err := serial.SendCommand(command, true, 10*time.Second)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send command to device: %v", err), http.StatusServiceUnavailable)
		return
	}

	// Trigger a switch map sync in case standard switches were enabled/disabled
	go serial.SyncFirmwareConfig()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, resp)
}

func handleGetPowerStatus(w http.ResponseWriter, r *http.Request) {
	serial.Status.RLock()
	defer serial.Status.RUnlock()
	if serial.Status.Data == nil {
		http.Error(w, "Status cache is not yet populated", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serial.Status.Data)
}

func handleSetAllPower(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var payload struct {
		State bool `json:"state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	stateInt := 0
	if payload.State {
		stateInt = 1
	}
	command := fmt.Sprintf(`{"set":{"all":%d}}`, stateInt)
	responseJSON, err := serial.SendCommand(command, true, 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send command to device: %v", err), http.StatusServiceUnavailable)
		return
	}
	var statusData map[string]map[string]interface{}
	if json.Unmarshal([]byte(responseJSON), &statusData) == nil {
		serial.Status.Lock()
		serial.Status.Data = statusData["status"]
		serial.Status.Unlock()
	}
	w.WriteHeader(http.StatusOK)
}

func handleGetLiveStatus(w http.ResponseWriter, r *http.Request) {
	serial.Conditions.RLock()
	defer serial.Conditions.RUnlock()
	if serial.Conditions.Data == nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{}")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serial.Conditions.Data)
}

func handleDeviceCommand(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// We need to check the command type to see if we should wait for a response.
	var commandPayload struct {
		Command string `json:"command"`
	}
	// We ignore the error here because the body might be a different type of command JSON
	// and we want to handle those generically.
	json.Unmarshal(body, &commandPayload)

	commandJSON := string(body)

	// Fire-and-forget commands
	if commandPayload.Command == "reboot" || commandPayload.Command == "factory_reset" {
		logger.Info("Received command '%s' from web UI. Sending to device.", commandJSON)
		serial.SendCommand(commandJSON, true, 0) // Don't wait for response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"Command sent successfully"}`) // Return valid JSON
		return
	}

	// For all other commands, send and wait for a response.
	if commandPayload.Command == "dry_sensor" {
		logger.Info("Received command '%s' from web UI. Sending to device.", commandJSON)
	} else {
		logger.Debug("Received generic command from web UI: %s", commandJSON)
	}

	// Use a timeout that's appropriate for commands that might take a moment.
	resp, err := serial.SendCommand(commandJSON, true, 5*time.Second)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send command to device: %v", err), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// The device response is expected to be JSON, so we can just pass it through.
	fmt.Fprint(w, resp)
}

func handleDownloadLog(w http.ResponseWriter, r *http.Request) {
	logPath := logger.GetLogFilePath()
	if logPath == "" {
		http.Error(w, "Log file path not available", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=\"proxy.log\"")
	http.ServeFile(w, r, logPath)
}

func handleGetFirmwareVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Version string `json:"version"`
	}{
		Version: serial.GetFirmwareVersion(),
	}
	json.NewEncoder(w).Encode(response)
}

func handleGetProxyVersion(appVersion string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Version string `json:"version"`
		}{
			Version: appVersion,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func handleCreateBackup(w http.ResponseWriter, r *http.Request) {
	logger.Info("Creating combined configuration backup...")
	firmwareConfigJSON, err := serial.SendCommand(`{"get":"config"}`, true, 0)
	if err != nil {
		http.Error(w, "Failed to get firmware configuration", http.StatusInternalServerError)
		return
	}
	backup := config.CombinedConfig{
		ProxyConfig:    config.Get(),
		FirmwareConfig: json.RawMessage(firmwareConfigJSON),
	}
	backupJSON, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		http.Error(w, "Failed to create backup file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename="sv241_backup.json"`)
	w.Write(backupJSON)
	logger.Info("Successfully created and sent configuration backup.")
}

func handleRestoreBackup(w http.ResponseWriter, r *http.Request) {
	logger.Info("Restoring combined configuration from backup...")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var backup config.CombinedConfig
	if err := json.Unmarshal(body, &backup); err != nil {
		http.Error(w, "Invalid backup file format", http.StatusBadRequest)
		return
	}
	if backup.ProxyConfig == nil || backup.FirmwareConfig == nil {
		http.Error(w, "Incomplete backup file", http.StatusBadRequest)
		return
	}

	// Restore Firmware Config
	compactFirmwareConfig, _ := json.Marshal(backup.FirmwareConfig)
	firmwareCommand := fmt.Sprintf(`{"sc":%s}`, string(compactFirmwareConfig))
	if _, err := serial.SendCommand(firmwareCommand, true, 10*time.Second); err != nil {
		http.Error(w, fmt.Sprintf("Failed to send firmware config to device: %v", err), http.StatusServiceUnavailable)
		return
	}
	logger.Info("Firmware configuration restored successfully.")

	// Restore Proxy Config
	conf := config.Get()
	conf.NetworkPort = backup.ProxyConfig.NetworkPort
	conf.ListenAddress = backup.ProxyConfig.ListenAddress
	conf.LogLevel = backup.ProxyConfig.LogLevel
	conf.SwitchNames = backup.ProxyConfig.SwitchNames
	conf.HeaterAutoEnableLeader = backup.ProxyConfig.HeaterAutoEnableLeader
	conf.HistoryRetentionNights = backup.ProxyConfig.HistoryRetentionNights
	conf.TelemetryInterval = backup.ProxyConfig.TelemetryInterval
	conf.EnableAlpacaVoltageControl = backup.ProxyConfig.EnableAlpacaVoltageControl
	conf.EnableMasterPower = backup.ProxyConfig.EnableMasterPower
	conf.AutoDetectPort = backup.ProxyConfig.AutoDetectPort
	conf.SerialPortName = "" // Clear port to trigger auto-detection
	logger.Info("Serial port name cleared to trigger auto-detection.")
	logger.SetLevelFromString(conf.LogLevel)

	if err := config.Save(); err != nil {
		http.Error(w, "Failed to save proxy configuration", http.StatusInternalServerError)
		return
	}
	logger.Info("Proxy configuration restored successfully.")

	// Synchronously attempt to reconnect so the user comes back to a connected system
	logger.Info("Restore: Disconnecting current session...")
	serial.Reconnect("") // Ensure we are disconnected first to free the port

	// Give the OS a moment to release the serial port handle
	logger.Info("Restore: Waiting for port to release...")
	time.Sleep(1 * time.Second)

	logger.Info("Restore: attempting immediate auto-detection...")
	foundPort, err := serial.FindPort()
	if err == nil {
		logger.Info("Restore: Immediate auto-detection found port '%s'. Reconnecting...", foundPort)
		serial.Reconnect(foundPort)
		fmt.Fprintf(w, "Configuration restored successfully. Connected to %s.", foundPort)
	} else {
		logger.Warn("Restore: Immediate auto-detection failed: %v. Background task will retry.", err)
		// Leave it to the background task
		go serial.Reconnect("")
		fmt.Fprint(w, "Configuration restored successfully. Logic will retry connection in background.")
	}
}

// handleSerialRelease closes the serial port to allow external tools (e.g., web flasher) to access it.
func handleSerialRelease(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logger.Info("API request to release serial port received.")
	err := serial.ReleasePort()

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		logger.Error("Failed to release serial port: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Serial port released. Auto-reconnect is paused.",
	})
}

// handleSerialResume resumes auto-reconnect after flashing is complete.
func handleSerialResume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logger.Info("API request to resume serial reconnect received.")
	serial.ResumeReconnect()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Auto-reconnect resumed.",
	})
}
*/
