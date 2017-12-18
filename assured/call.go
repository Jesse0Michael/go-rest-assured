package assured

import (
	"fmt"
)

// Call is a structure containing a request that is stubbed or made
type Call struct {
	Path       string            `json:"path"`
	Method     string            `json:"method"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Response   []byte            `json:"response,omitempty"`
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

// Callback is a structure containing a callback that is stubbed
type Callback struct {
	Target   string            `json:"target"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	Response []byte            `json:"response,omitempty"`
}
