package assured

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client wraps client specific configuration and behavior.
type Client struct {
	ClientOptions
}

// NewClient creates a new assured client
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		ClientOptions: DefaultClientOptions,
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

		req, err := http.NewRequestWithContext(ctx, call.Method, c.assuredURL("assured/given"), bytes.NewReader(b))
		if err != nil {
			return err
		}

		if err := c.process(req, nil); err != nil {
			return err
		}
	}
	return nil
}

// Verify returns all of the records made against a stubbed method and path
func (c *Client) Verify(ctx context.Context, method, path string) ([]Record, error) {
	body := Call{
		Method: method,
		Path:   path,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.assuredURL("assured/verify"), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	var records []Record
	if err = c.process(req, &records); err != nil {
		return nil, err
	}
	return records, nil
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.assuredURL("assured/clear"), bytes.NewReader(b))
	if err != nil {
		return err
	}
	return c.process(req, nil)
}

// ClearAll clears all assured calls
func (c *Client) ClearAll(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.assuredURL("assured/clearall"), nil)
	if err != nil {
		return err
	}
	return c.process(req, nil)
}

func (c *Client) assuredURL(path string) string {
	base := strings.TrimRight(c.baseURL, "/")
	return fmt.Sprintf("%s/%s", base, strings.TrimPrefix(path, "/"))
}

// process executes an HTTP request, applies shared error handling, and optionally unmarshals JSON into out.
func (c *Client) process(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		message := "unexpected response"
		if len(bodyBytes) > 0 {
			var apiError APIError
			if err = json.Unmarshal(bodyBytes, &apiError); err == nil && apiError.Error != "" {
				message = apiError.Error
			} else if trimmed := strings.TrimSpace(string(bodyBytes)); trimmed != "" {
				message = trimmed
			}
		}
		return fmt.Errorf("%d:%s", resp.StatusCode, message)
	}
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return err
		}
	}
	return nil
}
