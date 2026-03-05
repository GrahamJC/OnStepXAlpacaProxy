package alpaca

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	serverTransactionID uint32
)

type AlpacaIDs struct {
	ClientID            uint32
	ClientTransactionID uint32
	ServerTransactionID uint32
}

// Handler is a middleware that wraps HTTP handlers to provide Alpaca-specific functionality.
// It parses ClientTransactionID and ClientID from the request form.
func Handler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Parse query string and form data
		if err := r.ParseForm(); err != nil {
			slog.Warn("Error parsing form data for request", "method", r.Method, "path", r.URL.Path, "error", err)
		}

		// Log request
		params := ""
		for k, v := range r.Form {
			if s := strings.ToUpper(k); s != "CLIENTID" && s != "CLIENTTRANSACTIONID" {
				sep := ""
				if params != "" {
					sep = "&"
				}
				params = params + sep + k + "=" + v[0]
			}
		}
		slog.Debug("HTTP Request", "method", r.Method, "path", r.URL.Path, "params", params)

		// Add ClientID, ClientTransactionID and ServerTransactionID to context
		alpacaIDs := AlpacaIDs{
			ClientID:            0,
			ClientTransactionID: 0,
			ServerTransactionID: 0,
		}
		if idStr, ok := getFormValue(r, "ClientID"); ok {
			ui, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				slog.Debug(fmt.Sprintf("Invalid ClientID - %s", idStr))
			} else {
				alpacaIDs.ClientID = uint32(ui)
			}
		}
		if idStr, ok := getFormValue(r, "ClientTransactionID"); ok {
			ui, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				slog.Debug(fmt.Sprintf("Invalid ClientTransactionID - %s", idStr))
			} else {
				alpacaIDs.ClientTransactionID = uint32(ui)
			}
		}
		alpacaIDs.ServerTransactionID = atomic.AddUint32(&serverTransactionID, 1)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "AlpacaIDs", alpacaIDs)
		r = r.WithContext(ctx)

		// Call next handler
		fn(w, r)
	}
}

// GetFormValueIgnoreCase retrieves the first value for a given key from the request form, case-insensitively.
// The ASCOM conformance checker requires case-sensitivity for PUT parameters, so we handle that.
func getFormValue(r *http.Request, key string) (string, bool) {
	if r.Method == "PUT" {
		if values, ok := r.Form[key]; ok {
			if len(values) > 0 {
				return values[0], true
			}
			return "", true // Key exists, but has no value.
		}
		return "", false // Key not found with correct case.
	}

	// For GET and other methods, be case-insensitive.
	for k, values := range r.Form {
		if strings.EqualFold(k, key) {
			if len(values) > 0 {
				return values[0], true
			}
			return "", true // Key exists but has no value.
		}
	}
	return "", false
}

func checkBoolParam(w http.ResponseWriter, r *http.Request, name string) (bool, bool) {
	valStr, ok := getFormValue(r, name)
	if !ok {
		BadRequestResponse(w, r, fmt.Sprintf("Parameter '%s' not found", name))
		return false, false
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		BadRequestResponse(w, r, fmt.Sprintf("Parameter '%s' invalid format '%s'", name, valStr))
		return false, false
	}
	return val, true
}

func checkFloatParam(w http.ResponseWriter, r *http.Request, name string, minVal float64, maxVal float64) (float64, bool) {
	valStr, ok := getFormValue(r, name)
	if !ok {
		BadRequestResponse(w, r, fmt.Sprintf("Parameter '%s' not found", name))
		return 0, false
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		BadRequestResponse(w, r, fmt.Sprintf("Parameter '%s' invalid format '%s'", name, valStr))
		return 0, false
	}
	if !math.IsNaN(minVal) && val < minVal {
		ErrorResponse(w, r, 0x0401, fmt.Sprintf("Parameter '%s' invalid value %.1f (min %.1f)", name, val, minVal))
		return 0, false
	}
	if !math.IsNaN(maxVal) && val > maxVal {
		ErrorResponse(w, r, 0x0401, fmt.Sprintf("Parameter '%s' invalid value %.1f (max %.1f)", name, val, maxVal))
		return 0, false
	}
	return val, true
}

func checkIntParam(w http.ResponseWriter, r *http.Request, name string, minVal int64, maxVal int64) (int64, bool) {
	valStr, ok := getFormValue(r, name)
	if !ok {
		BadRequestResponse(w, r, fmt.Sprintf("Parameter '%s' not found", name))
		return 0, false
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		BadRequestResponse(w, r, fmt.Sprintf("Parameter '%s' invalid format '%s'", name, valStr))
		return 0, false
	}
	if minVal != math.MinInt64 && val < minVal {
		ErrorResponse(w, r, 0x0401, fmt.Sprintf("Parameter '%s' invalid value %d (min %d)", name, val, minVal))
		return 0, false
	}
	if maxVal != math.MaxInt64 && val > maxVal {
		ErrorResponse(w, r, 0x0401, fmt.Sprintf("Parameter '%s' invalid value %d (max %d)", name, val, maxVal))
		return 0, false
	}
	return val, true
}
