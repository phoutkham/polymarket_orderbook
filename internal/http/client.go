package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Request defines the structure for an HTTP request
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Params  map[string]string
	Body    any
}

// Client is a wrapper around http.Client for reusability
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new Client with a specified timeout
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Do executes an HTTP request and returns the response body
func (c *Client) Do(reqConfig Request) ([]byte, error) {
	// Parse URL and add query parameters
	u, err := url.Parse(reqConfig.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := u.Query()
	for k, v := range reqConfig.Params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	// Handle body
	var bodyReader io.Reader
	if reqConfig.Body != nil {
		switch b := reqConfig.Body.(type) {
		case io.Reader:
			bodyReader = b
		default:
			jsonData, err := json.Marshal(reqConfig.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal body: %w", err)
			}
			bodyReader = bytes.NewReader(jsonData)
		}
	}

	// Create request
	httpReq, err := http.NewRequest(reqConfig.Method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for k, v := range reqConfig.Headers {
		httpReq.Header.Set(k, v)
	}

	// Set default JSON header if body is provided and not already set
	if reqConfig.Body != nil && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respBody, fmt.Errorf("request returned non-2xx status: %d", resp.StatusCode)
	}

	return respBody, nil
}

func (c *Client) GET(path string) ([]byte, error) {
	req := Request{
		Method:  "GET",
		URL:     path,
		Headers: map[string]string{"Accept": "application/json"},
	}

	// Execute the request
	body, err := c.Do(req)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		if len(body) > 0 {
			fmt.Printf("Response Body: %s\n", string(body))
		}
		return nil, err
	}

	return body, nil
}

func (c *Client) POST(path string, _body interface{}) ([]byte, error) {
	req := Request{
		Method:  "POST",
		URL:     path,
		Headers: map[string]string{"Accept": "application/json"},
		Body:    _body,
	}

	// Execute the request
	body, err := c.Do(req)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		if len(body) > 0 {
			fmt.Printf("Response Body: %s\n", string(body))
		}
		return nil, err
	}

	return body, nil
}
