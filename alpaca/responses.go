package alpaca

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// --- Response Structs ---

type Response struct {
	ClientTransactionID uint32 `json:"ClientTransactionID"`
	ServerTransactionID uint32 `json:"ServerTransactionID"`
	ErrorNumber         int    `json:"ErrorNumber"`
	ErrorMessage        string `json:"ErrorMessage"`
}

type ValueResponse struct {
	Response
	Value interface{} `json:"Value"`
}

// --- Management API Response ---

// ManagementValueResponse is for management endpoints that don't use the standard handler.
func ManagementValueResponse(w http.ResponseWriter, r *http.Request, value interface{}) {
	response := struct {
		Value               interface{} `json:"Value"`
		ClientTransactionID uint32      `json:"ClientTransactionID"`
		ServerTransactionID uint32      `json:"ServerTransactionID"`
		ErrorNumber         int         `json:"ErrorNumber"`
		ErrorMessage        string      `json:"ErrorMessage"`
	}{
		Value:               value,
		ClientTransactionID: 0, // Not available in this context
		ServerTransactionID: 0, // Not stateful for this endpoint
		ErrorNumber:         0,
		ErrorMessage:        "",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// --- Standard Alpaca Responses ---
func BadRequestResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Set("Content-Type", "text")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func InternalErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Set("Content-Type", "text")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

func writeResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func SuccessResponse(w http.ResponseWriter, r *http.Request) {
	alpacaIDs, _ := r.Context().Value("AlpacaIDs").(AlpacaIDs)
	resp := Response{
		ClientTransactionID: alpacaIDs.ClientTransactionID,
		ServerTransactionID: alpacaIDs.ServerTransactionID,
	}
	writeResponse(w, r, resp)
}

func StringListResponse(w http.ResponseWriter, r *http.Request, value []string) {
	alpacaIDs, _ := r.Context().Value("AlpacaIDs").(AlpacaIDs)
	resp := ValueResponse{
		Response: Response{
			ClientTransactionID: alpacaIDs.ClientTransactionID,
			ServerTransactionID: alpacaIDs.ServerTransactionID,
		},
		Value: value,
	}
	writeResponse(w, r, resp)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, errNum int, errMsg string) {
	alpacaIDs, _ := r.Context().Value("AlpacaIDs").(AlpacaIDs)
	resp := Response{
		ClientTransactionID: alpacaIDs.ClientTransactionID,
		ServerTransactionID: alpacaIDs.ServerTransactionID,
		ErrorNumber:         errNum,
		ErrorMessage:        errMsg,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	slog.Error(fmt.Sprintf("Alpaca request failed with error %d: %s", errNum, errMsg))
	json.NewEncoder(w).Encode(resp)
}

func NotImplementedResponse(w http.ResponseWriter, r *http.Request, name string) {
	ErrorResponse(w, r, 0x400, fmt.Sprintf("Method '%s' not implemented", name))
}

func PropertyNotImplementedResponse(w http.ResponseWriter, r *http.Request, name string) {
	ErrorResponse(w, r, 0x400, fmt.Sprintf("Property '%s' not implemented", name))
}

func PropertyReadOnlyResponse(w http.ResponseWriter, r *http.Request, name string) {
	ErrorResponse(w, r, 0x400, fmt.Sprintf("Property '%s' cannot be written", name))
}

func AnyResponse(w http.ResponseWriter, r *http.Request, value interface{}) {
	alpacaIDs, _ := r.Context().Value("AlpacaIDs").(AlpacaIDs)
	resp := ValueResponse{
		Response: Response{
			ClientTransactionID: alpacaIDs.ClientTransactionID,
			ServerTransactionID: alpacaIDs.ServerTransactionID,
		},
		Value: value,
	}
	writeResponse(w, r, resp)
}

func StringResponse(w http.ResponseWriter, r *http.Request, value string) {
	alpacaIDs, _ := r.Context().Value("AlpacaIDs").(AlpacaIDs)
	resp := ValueResponse{
		Response: Response{
			ClientTransactionID: alpacaIDs.ClientTransactionID,
			ServerTransactionID: alpacaIDs.ServerTransactionID,
		},
		Value: value,
	}
	writeResponse(w, r, resp)
}

func IntResponse(w http.ResponseWriter, r *http.Request, value int) {
	alpacaIDs, _ := r.Context().Value("AlpacaIDs").(AlpacaIDs)
	resp := ValueResponse{
		Response: Response{
			ClientTransactionID: alpacaIDs.ClientTransactionID,
			ServerTransactionID: alpacaIDs.ServerTransactionID,
		},
		Value: value,
	}
	writeResponse(w, r, resp)
}

func FloatResponse(w http.ResponseWriter, r *http.Request, value float64) {
	// Round to 8 decimal places to avoid IEEE 754 floating-point precision issues
	// (e.g., 3.4 showing as 3.4000000000000004)
	//value = math.Round(value*100000000) / 100000000
	alpacaIDs, _ := r.Context().Value("AlpacaIDs").(AlpacaIDs)

	resp := ValueResponse{
		Response: Response{
			ClientTransactionID: alpacaIDs.ClientTransactionID,
			ServerTransactionID: alpacaIDs.ServerTransactionID,
		},
		Value: value,
	}
	writeResponse(w, r, resp)
}

func BoolResponse(w http.ResponseWriter, r *http.Request, value bool) {
	alpacaIDs, _ := r.Context().Value("AlpacaIDs").(AlpacaIDs)
	resp := ValueResponse{
		Response: Response{
			ClientTransactionID: alpacaIDs.ClientTransactionID,
			ServerTransactionID: alpacaIDs.ServerTransactionID,
		},
		Value: value,
	}
	writeResponse(w, r, resp)
}
