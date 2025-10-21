package assured

import (
	"bytes"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type APIError struct {
	Error string `json:"error"`
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// handleGiven is used to stub out a call for a given path
func handleGiven(logger *slog.Logger, assuredCalls *CallStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		call, err := decode[Call](r)
		if err != nil {
			_ = encode(w, http.StatusBadRequest, APIError{err.Error()})
			return
		}

		// Sanitize Path
		call.Path = strings.Trim(call.Path, "/")

		if call.Method == "" {
			call.Method = http.MethodGet
		}

		// validate http request
		_, err = http.NewRequest(call.Method, call.Path, nil)
		if err != nil {
			_ = encode(w, http.StatusBadRequest, APIError{err.Error()})
			return
		}

		for _, callback := range call.Callbacks {
			if callback.Target == "" {
				_ = encode(w, http.StatusBadRequest, APIError{"cannot stub callback without target"})
				return
			}
			_, err = http.NewRequest(callback.Method, callback.Target, nil)
			if err != nil {
				_ = encode(w, http.StatusBadRequest, APIError{err.Error()})
				return
			}
		}

		assuredCalls.Add(&call)
		logger.With("path", call.ID()).Info("assured call set")

		_ = encode(w, http.StatusOK, call)
	}
}

// handleWhen is used to respond to a given assured call
func handleWhen(logger *slog.Logger, httpClient *http.Client, assuredCalls, madeCalls *CallStore, trackMadeCalls bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		call := decodeAssuredCall(r)
		calls := assuredCalls.Get(call.ID())
		if len(calls) == 0 {
			logger.With("path", call.ID()).Info("assured call not found")
			_ = encode(w, http.StatusNotFound, APIError{"no assured calls"})
			return
		}

		if trackMadeCalls {
			madeCalls.Add(&call)
		}
		assured := calls[0]
		assuredCalls.Rotate(assured)

		// Trigger callbacks, if applicable
		for _, callback := range assured.Callbacks {
			go sendCallback(logger, httpClient, callback)
		}

		// Delay response
		time.Sleep(time.Duration(assured.Delay) * time.Second)

		logger.With("path", call.ID()).Info("assured call responded")
		_ = encodeAssuredCall(w, assured)
	}
}

// handleVerify returns all matching assured calls, used to verify a particular call
func handleVerify(madeCalls *CallStore, trackMadeCalls bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decode[Call](r)
		if err != nil {
			_ = encode(w, http.StatusBadRequest, APIError{err.Error()})
			return
		}

		_, err = http.NewRequest(req.Method, req.Path, nil)
		if err != nil {
			_ = encode(w, http.StatusBadRequest, APIError{err.Error()})
			return
		}

		if !trackMadeCalls {
			_ = encode(w, http.StatusNotFound, APIError{"tracking made calls is disabled"})
			return
		}

		calls := madeCalls.Get(req.ID())
		_ = encodeAssuredCall(w, calls)
	}
}

// handleClear is used to clear a specific assured call
func handleClear(logger *slog.Logger, assuredCalls, madeCalls *CallStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decode[Call](r)
		if err != nil {
			_ = encode(w, http.StatusBadRequest, APIError{err.Error()})
			return
		}

		_, err = http.NewRequest(req.Method, req.Path, nil)
		if err != nil {
			_ = encode(w, http.StatusBadRequest, APIError{err.Error()})
			return
		}

		assuredCalls.Clear(req.ID())
		madeCalls.Clear(req.ID())
		logger.With("path", req.ID()).Info("cleared calls for path")
	}
}

// handleClearAll is used to clear all assured calls
func handleClearAll(logger *slog.Logger, assuredCalls, madeCalls *CallStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assuredCalls.ClearAll()
		madeCalls.ClearAll()
		logger.Info("cleared all calls")
	}
}

// sendCallback sends a given callback to its target
func sendCallback(logger *slog.Logger, httpClient *http.Client, callback Callback) {
	req, err := http.NewRequest(callback.Method, callback.Target, bytes.NewBuffer(callback.Response))
	if err != nil {
		logger.With("target", callback.Target, "error", err).Info("failed to build callback request")
		return
	}
	for key, value := range callback.Headers {
		req.Header.Set(key, value)
	}
	// Delay callback, if applicable
	time.Sleep(time.Duration(callback.Delay) * time.Second)
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.With("target", callback.Target, "error", err).Info("failed to reach callback target")
		return
	}
	logger.With("target", callback.Target, "status_code", resp.StatusCode).Info("sent callback to target")
}
