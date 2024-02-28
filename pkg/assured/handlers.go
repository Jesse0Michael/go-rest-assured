package assured

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type APIError struct {
	Error error `json:"error"`
}

// handleGiven is used to stub out a call for a given path
func handleGiven(logger *slog.Logger, assuredCalls *CallStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		call := decodeAssuredCall(r)

		assuredCalls.Add(&call)
		logger.With("path", call.ID()).Info("assured call set")

		_ = encode[Call](w, http.StatusOK, call)
	}
}

// handleGivenCallback is used to stub out callbacks for a callback key
func handleGivenCallback(logger *slog.Logger, callbackCalls *CallStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		call, err := decodeAssuredCallback(r)
		if err != nil {
			_ = encode[APIError](w, http.StatusBadRequest, APIError{err})
		}

		callbackCalls.AddAt(call.Headers[AssuredCallbackKey], &call)
		logger.With("key", call.Headers[AssuredCallbackKey], "target", call.Headers[AssuredCallbackTarget]).Info("assured callback set")

		_ = encode[Call](w, http.StatusOK, call)
	}
}

// handleWhen is used to respond to a given assured call
func handleWhen(logger *slog.Logger, httpClient *http.Client, assuredCalls, madeCalls, callbackCalls *CallStore, trackMadeCalls bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		call := decodeAssuredCall(r)
		calls := assuredCalls.Get(call.ID())
		if len(calls) == 0 {
			logger.With("path", call.ID()).Info("assured call not found")
			_ = encode[APIError](w, http.StatusNotFound, APIError{errors.New("No assured calls")})
			return
		}

		if trackMadeCalls {
			madeCalls.Add(&call)
		}
		assured := calls[0]
		assuredCalls.Rotate(assured)

		// Trigger callbacks, if applicable
		for _, callback := range callbackCalls.Get(assured.Headers[AssuredCallbackKey]) {
			go sendCallback(logger, httpClient, callback.Headers[AssuredCallbackTarget], callback)
		}

		// Delay response
		if delay, err := strconv.ParseInt(assured.Headers[AssuredDelay], 10, 64); err == nil {
			time.Sleep(time.Duration(delay) * time.Second)
		}

		logger.With("path", call.ID()).Info("assured call responded")
		_ = encodeAssuredCall(w, assured)
	}
}

// handleVerify returns all matching assured calls, used to verify a particular call
func handleVerify(logger *slog.Logger, madeCalls *CallStore, trackMadeCalls bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		call := decodeAssuredCall(r)

		if !trackMadeCalls {
			_ = encode[APIError](w, http.StatusNotFound, APIError{errors.New("Tracking made calls is disabled")})
			return
		}

		calls := madeCalls.Get(call.ID())
		_ = encodeAssuredCall(w, calls)
	}
}

// handleClear is used to clear a specific assured call
func handleClear(logger *slog.Logger, assuredCalls, madeCalls, callbackCalls *CallStore, trackMadeCalls bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		call := decodeAssuredCall(r)

		assuredCalls.Clear(call.ID())
		madeCalls.Clear(call.ID())
		logger.With("path", call.ID()).Info("cleared calls for path")
		if call.Headers[AssuredCallbackKey] != "" {
			callbackCalls.Clear(call.Headers[AssuredCallbackKey])
			logger.With("key", call.Headers[AssuredCallbackKey]).Info("cleared calls for key")
		}
	}
}

// handleClearAll is used to clear all assured calls
func handleClearAll(logger *slog.Logger, assuredCalls, madeCalls, callbackCalls *CallStore, trackMadeCalls bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assuredCalls.ClearAll()
		madeCalls.ClearAll()
		callbackCalls.ClearAll()
		logger.Info("cleared all calls")
	}
}

// sendCallback sends a given callback to its target
func sendCallback(logger *slog.Logger, httpClient *http.Client, target string, call *Call) {
	var delay int64
	if delayOverride, err := strconv.ParseInt(call.Headers[AssuredCallbackDelay], 10, 64); err == nil {
		delay = delayOverride
	}
	req, err := http.NewRequest(call.Method, target, bytes.NewBuffer(call.Response))
	if err != nil {
		logger.With("target", target, "error", err).Info("failed to build callback request")
		return
	}
	for key, value := range call.Headers {
		req.Header.Set(key, value)
	}
	// Delay callback, if applicable
	time.Sleep(time.Duration(delay) * time.Second)
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.With("target", target, "error", err).Info("failed to reach callback target")
		return
	}
	logger.With("target", target, "status_code", resp.StatusCode).Info("sent callback to target")
}
