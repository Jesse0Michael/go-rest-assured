package assured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/pborman/uuid"
)

// Client
type Client struct {
	Errc chan error
	Options
}

// NewClient creates a new go-rest-assured client
func NewClient(opts ...Option) *Client {
	c := Client{
		Errc:    make(chan error),
		Options: DefaultOptions,
	}
	c.Options.applyOptions(opts...)

	if c.ctx == nil {
		c.ctx = context.Background()
	}
	c.ctx, c.cancel = context.WithCancel(c.Options.ctx)

	if c.Options.Port == 0 {
		if listen, err := net.Listen("tcp", ":0"); err == nil {
			c.Options.Port = listen.Addr().(*net.TCPAddr).Port
			listen.Close()
		}
	}
	c.startApplicationHTTPListener()
	return &c
}

// url returns the url to used by the client internally
func (c *Client) url() string {
	schema := "http"
	if c.tlsCertFile != "" && c.tlsKeyFile != "" {
		schema = "https"
	}
	return fmt.Sprintf("%s://%s:%d", schema, c.host, c.Port)
}

// URL returns the url to use to test you stubbed endpoints
func (c *Client) URL() string {
	return fmt.Sprintf("%s/when", c.url())
}

// Close is used to close the running service
func (c *Client) Close() {
	c.cancel()
}

// Given stubs assured Call(s)
func (c *Client) Given(calls ...Call) error {
	for _, call := range calls {
		// Default method to GET
		if call.Method == "" {
			call.Method = http.MethodGet
		}

		// Sanitize Path
		call.Path = strings.Trim(call.Path, "/")

		req, err := http.NewRequest(call.Method, fmt.Sprintf("%s/given/%s", c.url(), call.Path), bytes.NewReader(call.Response))
		if err != nil {
			return err
		}
		if call.StatusCode != 0 {
			req.Header.Set(AssuredStatus, strconv.Itoa(call.StatusCode))
		}
		if call.Delay > 0 {
			req.Header.Set(AssuredDelay, strconv.Itoa(call.Delay))
		}
		for key, value := range call.Headers {
			req.Header.Set(key, value)
		}

		// Create callbacks
		callbacks := make([]*http.Request, len(call.Callbacks))
		callbackKey := uuid.New()
		for i, callback := range call.Callbacks {
			if callback.Target == "" {
				return fmt.Errorf("cannot stub callback without target")
			}
			callbackReq, err := http.NewRequest(callback.Method, fmt.Sprintf("%s/callback", c.url()), bytes.NewReader(callback.Response))
			if err != nil {
				return err
			}
			callbackReq.Header.Set(AssuredCallbackTarget, callback.Target)
			callbackReq.Header.Set(AssuredCallbackKey, callbackKey)
			if callback.Delay > 0 {
				callbackReq.Header.Set(AssuredCallbackDelay, strconv.Itoa(callback.Delay))
			}
			for key, value := range callback.Headers {
				callbackReq.Header.Set(key, value)
			}
			callbacks[i] = callbackReq
		}
		if len(callbacks) > 0 {
			req.Header.Set(AssuredCallbackKey, callbackKey)
		}

		if _, err = c.httpClient.Do(req); err != nil {
			return err
		}
		for _, cReq := range callbacks {
			if _, err = c.httpClient.Do(cReq); err != nil {
				return err
			}
		}
	}
	return nil
}

// Verify returns all of the calls made against a stubbed method and path
func (c *Client) Verify(method, path string) ([]Call, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/verify/%s", c.url(), path), nil)
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

	var calls []Call
	if err = json.NewDecoder(resp.Body).Decode(&calls); err != nil {
		return nil, err
	}
	return calls, nil
}

// Clear assured calls for a Method and Path
func (c *Client) Clear(method, path string) error {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/clear/%s", c.url(), path), nil)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Do(req)
	return err
}

// ClearAll clears all assured calls
func (c *Client) ClearAll() error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/clear", c.url()), nil)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Do(req)
	return err
}
