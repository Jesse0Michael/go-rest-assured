package assured

import (
	// "net/http"

	kitlog "github.com/go-kit/kit/log"
)

// Client
type Client struct {
	Logger kitlog.Logger
}

// NewClient creates a new go-rest-assured client with the given parameters
func NewClient(l kitlog.Logger) *Client {
	return &Client{
		Logger: l,
	}
}

// Run starts the go-rest-assured service through the client
// func (c *Client) Run() error {

// }

// // Close is used to close the running service
// func (c *Client) Close() {

// }

// // URL returns the URL needed to use the stubbed assured endpoints
// func (c *Client) URL() (string, error) {

// }

// // Given stubs an assured Call
// func (c *Client) Given(call *Call) error {
// 	//Status.HttpOK if not set on call
// 	//strip leading `/`? use url library to strip query and stuff
// }

// // Verify returns all of the calls made against a stubbed method and path
// func (c *Client) Verify(method, path string) ([]*Call, error) {

// }

// // Clear assured calls for a Method and Path
// func (c *Client) Clear(method, path string) error {

// }

// // Clear assured calls for a Method and Path
// func (c *Client) Clear(method, path string) error {

// }
