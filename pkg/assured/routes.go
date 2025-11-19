package assured

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func routes(
	logger *slog.Logger,
	calls *Store[Call],
	records *Store[Record],
	httpClient *http.Client,
	trackRecords bool,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/assured/health", handleHealth)
	mux.HandleFunc("/assured/given", handleGiven(logger, calls))
	mux.HandleFunc("/assured/verify", handleVerify(records, trackRecords))
	mux.HandleFunc("/assured/clear", handleClear(logger, calls, records))
	mux.HandleFunc("/assured/clearall", handleClearAll(logger, calls, records))
	mux.HandleFunc("/", handleWhen(logger, httpClient, calls, records, trackRecords))

	return mux
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

// decodeAssuredRecord converts an http request into an assured Record object
func decodeAssuredRecord(req *http.Request) Record {
	record := Record{
		Path:    strings.Trim(req.URL.Path, "/"),
		Method:  req.Method,
		Cookies: req.Cookies(),
	}

	// Set headers
	headers := map[string]string{}
	for key, value := range req.Header {
		headers[key] = value[0]
	}
	record.Headers = headers

	// Set query
	query := map[string]string{}
	for key, value := range req.URL.Query() {
		query[key] = value[0]
	}
	record.Query = query

	// Set response body
	if req.Body != nil {
		defer func() { _ = req.Body.Close() }()
		if bytes, err := io.ReadAll(req.Body); err == nil {
			record.Body = bytes
		}
	}

	return record
}

// encodeAssuredCall writes the assured Call to the http response as it is intended to be stubbed
func encodeAssuredCall(w http.ResponseWriter, i interface{}) error {
	switch resp := i.(type) {
	case Call:
		for key, value := range resp.Headers {
			w.Header().Set(key, value)
		}
		if resp.StatusCode > 0 {
			w.WriteHeader(resp.StatusCode)
		}
		_, _ = w.Write([]byte(resp.String()))
	case []Record:
		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(resp)
	}
	return nil
}
