package assured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	AssuredStatus         = "Assured-Status"
	AssuredMethod         = "Assured-Method"
	AssuredDelay          = "Assured-Delay"
	AssuredCallbackKey    = "Assured-Callback-Key"
	AssuredCallbackTarget = "Assured-Callback-Target"
	AssuredCallbackDelay  = "Assured-Callback-Delay"
)

// Client
type Client struct {
	Options
}

// NewClient creates a new go-rest-assured client
func NewClient(opts ...Option) *Client {
	c := &Client{
		Options: DefaultOptions,
	}
	c.applyOptions(opts...)
	return c
}

// Given stubs assured Call(s)
func (c *Client) Given(ctx context.Context, calls ...Call) error {
	for _, call := range calls {
		// Default method to GET
		if call.Method == "" {
			call.Method = http.MethodGet
		}

		// Sanitize Path
		call.Path = strings.Trim(call.Path, "/")

		req, err := http.NewRequestWithContext(ctx, call.Method, fmt.Sprintf("%s/given/%s", c.url(), call.Path), bytes.NewReader(call.Response))
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
		callbackKey := uuid.NewString()
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
func (c *Client) Verify(ctx context.Context, method, path string) ([]Call, error) {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/verify/%s", c.url(), path), nil)
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
	defer func() { _ = resp.Body.Close() }()

	var calls []Call
	if err = json.NewDecoder(resp.Body).Decode(&calls); err != nil {
		return nil, err
	}
	return calls, nil
}

// Clear assured calls for a Method and Path
func (c *Client) Clear(ctx context.Context, method, path string) error {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/clear/%s", c.url(), path), nil)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Do(req)
	return err
}

// ClearAll clears all assured calls
func (c *Client) ClearAll(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/clear", c.url()), nil)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Do(req)
	return err
}
