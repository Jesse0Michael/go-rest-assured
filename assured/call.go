package assured

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// Call is a structure containing a request that is stubbed or made
type Call struct {
	Path       string            `json:"path"`
	Method     string            `json:"method"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Query      map[string]string `json:"query,omitempty"`
	Response   CallResponse      `json:"response,omitempty"`
	Callbacks  []Callback        `json:"callbacks,omitempty"`
}

// ID is used as a key when managing stubbed and made calls
func (c Call) ID() string {
	return fmt.Sprintf("%s:%s", c.Method, c.Path)
}

// String converts a Call's Response into a string
func (c Call) String() string {
	rawString := string(c.Response)

	// TODO: implement string replacements for special cases
	return rawString
}

// CallResponse allows control over the Call's Response encoding
type CallResponse []byte

// UnmarshalJSON is a custom implementation for JSON Unmarshalling for the CallResponse
// Unmarshalling will first check if the data is a local filepath that can be read
// Else it will check if the data is stringified JSON and un-stringify the data to use
// or Else it will just use the []byte
func (response *CallResponse) UnmarshalJSON(data []byte) error {
	unmarshaled := []byte{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		// The data is a []byte, so use it
		unmarshaled = data
	}

	if s, err := strconv.Unquote(string(unmarshaled)); err == nil {
		absPath, _ := filepath.Abs(s)
		if _, err := os.Stat(absPath); err == nil {
			// The data is a path that exists, therefore we will read the file
			if file, err := ioutil.ReadFile(absPath); err == nil {
				*response = file
				return nil
			}
		}
		// The data is stringified JSON, therefore we eill use the unquoted JSON
		*response = []byte(s)
		return nil
	}

	*response = unmarshaled
	return nil
}

// Callback is a structure containing a callback that is stubbed
type Callback struct {
	Target   string            `json:"target"`
	Method   string            `json:"method"`
	Delay    int               `json:"delay,omitempty"`
	Headers  map[string]string `json:"headers"`
	Response CallResponse      `json:"response,omitempty"`
}
