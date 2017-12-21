package assured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	kitlog "github.com/go-kit/kit/log"
	"github.com/pborman/uuid"
)

// Client
type Client struct {
	Errc       chan error
	Port       int
	ctx        context.Context
	cancel     context.CancelFunc
	httpClient http.Client
}

// NewDefaultClient creates a new go-rest-assured client with default parameters
func NewDefaultClient() *Client {
	settings := Settings{
		Logger:         kitlog.NewLogfmtLogger(ioutil.Discard),
		HTTPClient:     *http.DefaultClient,
		TrackMadeCalls: true,
	}
	return NewClient(nil, settings)
}

// NewClient creates a new go-rest-assured client
func NewClient(root context.Context, settings Settings) *Client {
	if root == nil {
		root = context.Background()
	}
	if settings.Port == 0 {
		if listen, err := net.Listen("tcp", ":0"); err == nil {
			settings.Port = listen.Addr().(*net.TCPAddr).Port
			listen.Close()
		}
	}
	ctx, cancel := context.WithCancel(root)
	c := Client{
		Errc:       make(chan error),
		Port:       settings.Port,
		ctx:        ctx,
		cancel:     cancel,
		httpClient: settings.HTTPClient,
	}
	StartApplicationHTTPListener(c.ctx, c.Errc, settings)
	return &c
}

// URL returns the url to use to test you stubbed endpoints
func (c *Client) URL() string {
	return fmt.Sprintf("http://localhost:%d/when", c.Port)
}

// Close is used to close the running service
func (c *Client) Close() {
	c.cancel()
}

// Given stubs assured Call(s)
func (c *Client) Given(calls ...Call) error {
	for _, call := range calls {
		if call.Method == "" {
			return fmt.Errorf("cannot stub call without Method")
		}

		// Sanitize Path
		call.Path = strings.Trim(call.Path, "/")

		req, err := http.NewRequest(call.Method, fmt.Sprintf("http://localhost:%d/given/%s", c.Port, call.Path), bytes.NewReader(call.Response))
		if err != nil {
			return err
		}
		if call.StatusCode != 0 {
			req.Header.Set(AssuredStatus, fmt.Sprintf("%d", call.StatusCode))
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
			callbackReq, err := http.NewRequest(callback.Method, fmt.Sprintf("http://localhost:%d/callback", c.Port), bytes.NewReader(callback.Response))
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
	req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:%d/verify/%s", c.Port, path), nil)
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
	req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:%d/clear/%s", c.Port, path), nil)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Do(req)
	return err
}

// ClearAll clears all assured calls
func (c *Client) ClearAll() error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:%d/clear", c.Port), nil)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Do(req)
	return err
}
