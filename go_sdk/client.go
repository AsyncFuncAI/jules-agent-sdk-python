package jules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL    = "https://jules.googleapis.com/v1alpha"
	defaultMaxRetries = 3
	DefaultTimeout    = 30 * time.Second
)

// Client is the base client for the Jules API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	maxRetries int
}

// NewClient creates a new Jules API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: DefaultTimeout},
		maxRetries: defaultMaxRetries,
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, path)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "jules-go-sdk/0.1.0")

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < c.maxRetries; i++ {
		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break // Success or non-retriable error
		}
		time.Sleep(time.Duration(i*i) * time.Second) // Exponential backoff
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("API error: %s", resp.Status)
	}

    if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

func (c *Client) get(path string, v interface{}) (*http.Response, error) {
    req, err := c.newRequest("GET", path, nil)
    if err != nil {
        return nil, err
    }
    return c.do(req, v)
}

func (c *Client) post(path string, body, v interface{}) (*http.Response, error) {
    req, err := c.newRequest("POST", path, body)
    if err != nil {
        return nil, err
    }
    return c.do(req, v)
}