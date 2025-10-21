package assured

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
		b, err := json.Marshal(call)
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, call.Method, fmt.Sprintf("%s/assured/given", c.url()), bytes.NewReader(b))
		if err != nil {
			return err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			var apiError APIError
			if err = json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
				return err
			}
			if apiError.Error != "" {
				return errors.New(apiError.Error)
			}
			return fmt.Errorf("failure to stub assured call")
		}
	}
	return nil
}

// Verify returns all of the calls made against a stubbed method and path
func (c *Client) Verify(ctx context.Context, method, path string) ([]Call, error) {
	body := Call{
		Method: method,
		Path:   path,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/assured/verify", c.url()), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		var apiError APIError
		if err = json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return nil, err
		}
		if apiError.Error != "" {
			return nil, errors.New(apiError.Error)
		}
		return nil, fmt.Errorf("failure to stub assured call")
	}

	var calls []Call
	if err = json.NewDecoder(resp.Body).Decode(&calls); err != nil {
		return nil, err
	}
	return calls, nil
}

// Clear assured calls for a Method and Path
func (c *Client) Clear(ctx context.Context, method, path string) error {
	body := Call{
		Method: method,
		Path:   path,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/assured/clear", c.url()), bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var apiError APIError
		if err = json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return err
		}
		if apiError.Error != "" {
			return errors.New(apiError.Error)
		}
		return fmt.Errorf("failure to clear assured call")
	}
	return err
}

// ClearAll clears all assured calls
func (c *Client) ClearAll(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/assured/clearall", c.url()), nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var apiError APIError
		if err = json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return err
		}
		if apiError.Error != "" {
			return errors.New(apiError.Error)
		}
		return fmt.Errorf("failure to clear all assured calls")
	}
	return err
}
