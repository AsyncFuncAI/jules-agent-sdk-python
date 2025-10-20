package jules

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

// BaseClient is the base HTTP client for making API requests
type BaseClient struct {
	config       *Config
	httpClient   *http.Client
	requestCount int
	errorCount   int
}

// NewBaseClient creates a new BaseClient
func NewBaseClient(config *Config) (*BaseClient, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Create HTTP client with connection pooling
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	if !config.VerifySSL {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	return &BaseClient{
		config:     config,
		httpClient: httpClient,
	}, nil
}

// Get performs a GET request
func (c *BaseClient) Get(ctx context.Context, path string) (map[string]interface{}, error) {
	return c.request(ctx, http.MethodGet, path, nil)
}

// Post performs a POST request
func (c *BaseClient) Post(ctx context.Context, path string, body interface{}) (map[string]interface{}, error) {
	return c.request(ctx, http.MethodPost, path, body)
}

// request performs an HTTP request with retry logic
func (c *BaseClient) request(ctx context.Context, method, path string, body interface{}) (map[string]interface{}, error) {
	url := c.config.BaseURL + path
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		// Apply backoff delay for retries
		if attempt > 0 {
			backoff := c.calculateBackoff(attempt)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		// Prepare request body
		var bodyReader io.Reader
		if body != nil {
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewReader(bodyBytes)
		}

		// Create request
		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Goog-Api-Key", c.config.APIKey)

		// Execute request
		c.requestCount++
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			// Retry on network errors
			continue
		}

		// Read response body
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		// Handle HTTP errors
		if resp.StatusCode >= 400 {
			var responseData map[string]interface{}
			if len(respBody) > 0 {
				json.Unmarshal(respBody, &responseData)
			}

			lastErr = c.handleError(resp.StatusCode, respBody, responseData)

			// Retry on 5xx errors
			if resp.StatusCode >= 500 {
				c.errorCount++
				continue
			}

			// Don't retry on client errors (4xx)
			return nil, lastErr
		}

		// Parse successful response
		var result map[string]interface{}
		if len(respBody) > 0 {
			if err := json.Unmarshal(respBody, &result); err != nil {
				return nil, fmt.Errorf("failed to unmarshal response: %w", err)
			}
		}

		return result, nil
	}

	if lastErr != nil {
		c.errorCount++
		return nil, lastErr
	}

	return nil, fmt.Errorf("request failed after %d attempts", c.config.MaxRetries+1)
}

// calculateBackoff calculates the exponential backoff duration
func (c *BaseClient) calculateBackoff(attempt int) time.Duration {
	backoff := c.config.RetryBackoffFactor * math.Pow(2, float64(attempt-1))
	backoffDuration := time.Duration(backoff * float64(time.Second))
	if backoffDuration > c.config.MaxBackoff {
		backoffDuration = c.config.MaxBackoff
	}
	return backoffDuration
}

// handleError converts HTTP errors to appropriate error types
func (c *BaseClient) handleError(statusCode int, body []byte, response map[string]interface{}) error {
	message := string(body)
	if errMsg, ok := response["error"].(map[string]interface{}); ok {
		if msg, ok := errMsg["message"].(string); ok {
			message = msg
		}
	}

	switch statusCode {
	case 400:
		return NewValidationError(message, response)
	case 401:
		return NewAuthenticationError(message, response)
	case 404:
		return NewNotFoundError(message, response)
	case 429:
		// Try to extract Retry-After header
		retryAfter := 0
		if ra, ok := response["retryAfter"].(float64); ok {
			retryAfter = int(ra)
		}
		return NewRateLimitError(message, retryAfter, response)
	default:
		if statusCode >= 500 {
			return NewServerError(message, statusCode, response)
		}
		return &APIError{
			Message:    message,
			StatusCode: statusCode,
			Response:   response,
		}
	}
}

// Close closes the HTTP client and releases resources
func (c *BaseClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// Stats returns statistics about the client
func (c *BaseClient) Stats() map[string]int {
	return map[string]int{
		"request_count": c.requestCount,
		"error_count":   c.errorCount,
	}
}
