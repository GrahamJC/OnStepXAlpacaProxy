package alpaca

import (
	"encoding/json"
	"net/http"
)

// --- Management Handlers ---

// AlpacaDescription defines the structure for the management/v1/description endpoint.
type AlpacaDescription struct {
	ServerName          string `json:"ServerName"`
	Manufacturer        string `json:"Manufacturer"`
	ManufacturerVersion string `json:"ManufacturerVersion"`
	Location            string `json:"Location"`
}

// AlpacaConfiguredDevice defines the structure for a single device in the management/v1/configureddevices endpoint.
type AlpacaConfiguredDevice struct {
	DeviceName   string `json:"DeviceName"`
	DeviceType   string `json:"DeviceType"`
	DeviceNumber int    `json:"DeviceNumber"`
	UniqueID     string `json:"UniqueID"`
}

func SetupManagementHandlers(appVersion string) {
	http.HandleFunc("/management/v1/description", HandleManagementDescription(appVersion))
	http.HandleFunc("/management/v1/configureddevices", HandleManagementConfiguredDevices)
	http.HandleFunc("/management/apiversions", HandleManagementApiVersions)
}

func HandleManagementDescription(appVersion string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		description := AlpacaDescription{
			ServerName:          "OnStepX Alpaca Proxy",
			Manufacturer:        "User-Made",
			ManufacturerVersion: appVersion,
			Location:            "My Observatory",
		}
		ManagementValueResponse(w, r, description)
	}
}

// HandleManagementConfiguredDevices is static and doesn't need the API struct receiver.
func HandleManagementConfiguredDevices(w http.ResponseWriter, r *http.Request) {
	devices := []AlpacaConfiguredDevice{
		{
			DeviceName:   "OnStepX Telescope",
			DeviceType:   "Telescope",
			DeviceNumber: 0,
			UniqueID:     "a354fbdd-e686-4b99-913a-3ebd78bddf4d", // Static GUID
		},
	}
	ManagementValueResponse(w, r, devices)
}

// HandleManagementApiVersions is static and doesn't need the API struct receiver.
func HandleManagementApiVersions(w http.ResponseWriter, r *http.Request) {
	// This endpoint doesn't use the standard alpaca handler.
	response := struct {
		Value               []int  `json:"Value"`
		ClientTransactionID uint32 `json:"ClientTransactionID"`
		ServerTransactionID uint32 `json:"ServerTransactionID"`
		ErrorNumber         int    `json:"ErrorNumber"`
		ErrorMessage        string `json:"ErrorMessage"`
	}{
		Value: []int{1},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
