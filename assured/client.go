package assured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	"github.com/phayes/freeport"
)

// Client
type Client struct {
	logger     kitlog.Logger
	ctx        context.Context
	errc       chan error
	port       *int
	httpClient *http.Client
}

// NewClient creates a new go-rest-assured client with the given parameters
func NewClient(l kitlog.Logger, port *int) *Client {
	return &Client{
		logger:     l,
		ctx:        context.Background(),
		errc:       make(chan error),
		port:       port,
		httpClient: &http.Client{},
	}
}

// Port gets the port that go-rest-assured will run on, or an open port if not set
func (c *Client) Port() int {
	if c.port != nil {
		return *c.port
	}
	port := freeport.GetPort()
	c.port = &port
	return port
}

// Run starts the go-rest-assured service through the client
func (c *Client) Run() {
	StartApplicationHTTPListener(c.Port(), c.logger, c.ctx, c.errc)
}

// Close is used to close the running service
func (c *Client) Close() {

}

// Given stubs an assured Call
func (c *Client) Given(call *Call) error {
	var req *http.Request
	var err error

	if call.Response == nil {
		req, err = http.NewRequest(call.Method, fmt.Sprintf("given/%s", call.Path), nil)
	} else {
		req, err = http.NewRequest(call.Method, fmt.Sprintf("given/%s", call.Path), bytes.NewReader(call.Response))
	}
	if err != nil {
		return err
	}
	if call.StatusCode != 0 {
		req.Header.Set("Assert-Status", string(call.StatusCode))
	}

	_, err = c.httpClient.Do(req)
	return err
}

// Verify returns all of the calls made against a stubbed method and path
func (c *Client) Verify(method, path string) ([]*Call, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("verify/%s", path), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failure to verify calls")
	}
	defer resp.Body.Close()

	var calls []*Call
	if err = json.NewDecoder(resp.Body).Decode(calls); err != nil {
		return nil, err
	}
	return calls, nil
}

// Clear assured calls for a Method and Path
func (c *Client) Clear(method, path string) error {
	req, err := http.NewRequest(method, fmt.Sprintf("clear/%s", path), nil)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Do(req)
	return err
}

// ClearAll clears all assured calls
func (c *Client) ClearAll() error {
	req, err := http.NewRequest(http.MethodDelete, "clear", nil)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Do(req)
	return err
}
