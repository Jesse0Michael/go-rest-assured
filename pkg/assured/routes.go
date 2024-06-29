package assured

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func routes(
	logger *slog.Logger,
	assuredCalls *CallStore,
	madeCalls *CallStore,
	callbackCalls *CallStore,
	httpClient *http.Client,
	trackMadeCalls bool,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/given/{path...}", handleGiven(logger, assuredCalls))
	mux.HandleFunc("/callback", handleGivenCallback(logger, callbackCalls))
	mux.HandleFunc("/when/{path...}", handleWhen(logger, httpClient, assuredCalls, madeCalls, callbackCalls, trackMadeCalls))
	mux.HandleFunc("/verify/{path...}", handleVerify(logger, madeCalls, trackMadeCalls))
	mux.HandleFunc("/clear/{path...}", handleClear(logger, assuredCalls, madeCalls, callbackCalls, trackMadeCalls))
	mux.HandleFunc("/clear", handleClearAll(logger, assuredCalls, madeCalls, callbackCalls, trackMadeCalls))

	return mux
}

// func decode[T any](r *http.Request) (T, error) {
// 	var v T
// 	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
// 		return v, fmt.Errorf("decode json: %w", err)
// 	}
// 	return v, nil
// }

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

// decodeAssuredCall converts an http request into an assured Call object
func decodeAssuredCall(req *http.Request) Call {
	method := req.Method
	if m := req.Header.Get(AssuredMethod); m != "" {
		method = m
	}

	call := Call{
		Path:       req.PathValue("path"),
		Method:     method,
		StatusCode: http.StatusOK,
	}

	// Set status code override
	if statusCode, err := strconv.ParseInt(req.Header.Get(AssuredStatus), 10, 64); err == nil {
		call.StatusCode = int(statusCode)
	}

	// Set headers
	headers := map[string]string{}
	for key, value := range req.Header {
		headers[key] = value[0]
	}
	call.Headers = headers

	// Set query
	query := map[string]string{}
	for key, value := range req.URL.Query() {
		query[key] = value[0]
	}
	call.Query = query

	// Set response body
	if req.Body != nil {
		defer req.Body.Close()
		if bytes, err := io.ReadAll(req.Body); err == nil {
			call.Response = bytes
		}
	}

	return call
}

// decodeAssuredCallback converts an http request into an assured Callback object
func decodeAssuredCallback(req *http.Request) (Call, error) {
	call := Call{
		Method:     req.Method,
		StatusCode: http.StatusCreated,
	}

	// Require headers
	if len(req.Header[AssuredCallbackKey]) == 0 {
		return call, fmt.Errorf("'%s' header required for callback", AssuredCallbackKey)
	}
	if len(req.Header[AssuredCallbackTarget]) == 0 {
		return call, fmt.Errorf("'%s' header required for callback", AssuredCallbackTarget)
	}

	// Set headers
	headers := map[string]string{}
	for key, value := range req.Header {
		headers[key] = value[0]
	}
	call.Headers = headers

	// Set response body
	if req.Body != nil {
		defer req.Body.Close()
		if bytes, err := io.ReadAll(req.Body); err == nil {
			call.Response = bytes
		}
	}

	return call, nil
}

// encodeAssuredCall writes the assured Call to the http response as it is intended to be stubbed
func encodeAssuredCall(w http.ResponseWriter, i interface{}) error {
	switch resp := i.(type) {
	case *Call:
		for key, value := range resp.Headers {
			if !strings.HasPrefix(key, "Assured-") {
				w.Header().Set(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write([]byte(resp.String()))
	case []*Call:
		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(resp)
	}
	return nil
}
