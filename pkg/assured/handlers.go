package assured

import (
	"bytes"
	"context"
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
func handleGiven(logger *slog.Logger, calls *Store[Call]) http.HandlerFunc {
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

		calls.Add(call)
		logger.InfoContext(r.Context(), "assured call set", "key", call.Key())

		_ = encode(w, http.StatusOK, call)
	}
}

// handleWhen is used to respond to a given assured call
func handleWhen(logger *slog.Logger, httpClient *http.Client, calls *Store[Call], records *Store[Record], trackRecords bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		record := decodeAssuredRecord(r)
		matched := calls.Get(record.Key())
		if len(matched) == 0 {
			logger.InfoContext(r.Context(), "assured call not found", "key", record.Key())
			_ = encode(w, http.StatusNotFound, APIError{"no assured calls"})
			return
		}

		if trackRecords {
			records.Add(record)
		}
		assured := matched[0]
		calls.Rotate(assured)

		// Trigger callbacks, if applicable
		for _, callback := range assured.Callbacks {
			go sendCallback(r.Context(), logger, httpClient, callback)
		}

		// Delay response
		time.Sleep(time.Duration(assured.Delay) * time.Second)

		logger.InfoContext(r.Context(), "assured call responded", "key", record.Key())
		_ = encodeAssuredCall(w, assured)
	}
}

// handleVerify returns all matching assured calls, used to verify a particular call
func handleVerify(records *Store[Record], trackRecords bool) http.HandlerFunc {
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

		if !trackRecords {
			_ = encode(w, http.StatusNotFound, APIError{"tracking records is disabled"})
			return
		}

		calls := records.Get(req.Key())
		_ = encodeAssuredCall(w, calls)
	}
}

// handleClear is used to clear a specific assured call
func handleClear(logger *slog.Logger, calls *Store[Call], records *Store[Record]) http.HandlerFunc {
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

		calls.Clear(req.Key())
		records.Clear(req.Key())
		logger.InfoContext(r.Context(), "cleared calls for path", "key", req.Key())
	}
}

// handleClearAll is used to clear all assured calls
func handleClearAll(logger *slog.Logger, calls *Store[Call], records *Store[Record]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		calls.ClearAll()
		records.ClearAll()
		logger.InfoContext(r.Context(), "cleared all calls")
	}
}

// sendCallback sends a given callback to its target
func sendCallback(ctx context.Context, logger *slog.Logger, httpClient *http.Client, callback Callback) {
	req, err := http.NewRequest(callback.Method, callback.Target, bytes.NewBuffer(callback.Response))
	if err != nil {
		logger.InfoContext(ctx, "failed to build callback request", "target", callback.Target, "error", err)
		return
	}
	for key, value := range callback.Headers {
		req.Header.Set(key, value)
	}
	// Delay callback, if applicable
	time.Sleep(time.Duration(callback.Delay) * time.Second)
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.InfoContext(ctx, "failed to reach callback target", "target", callback.Target, "error", err)
		return
	}
	logger.InfoContext(ctx, "sent callback to target", "target", callback.Target, "status_code", resp.StatusCode)
}
