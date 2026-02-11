package jules

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockServer represents a mock HTTP server for testing
type MockServer struct {
	*httptest.Server
	Requests  []*http.Request
	Responses []MockResponse
	nextIndex int
}

// MockResponse represents a mock HTTP response
type MockResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
}

// NewMockServer creates a new mock server for testing
func NewMockServer(t *testing.T, responses []MockResponse) *MockServer {
	mock := &MockServer{
		Requests:  make([]*http.Request, 0),
		Responses: responses,
	}

	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.Requests = append(mock.Requests, r)

		if mock.nextIndex >= len(mock.Responses) {
			http.Error(w, "No more mock responses", http.StatusInternalServerError)
			return
		}

		response := mock.Responses[mock.nextIndex]
		mock.nextIndex++

		// Set headers
		for k, v := range response.Headers {
			w.Header().Set(k, v)
		}

		// Set status code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)

		// Write body
		if response.Body != nil {
			json.NewEncoder(w).Encode(response.Body)
		}
	}))

	return mock
}

// GetLastRequest returns the last request received by the mock server
func (m *MockServer) GetLastRequest() *http.Request {
	if len(m.Requests) == 0 {
		return nil
	}
	return m.Requests[len(m.Requests)-1]
}

// GetRequestCount returns the number of requests received
func (m *MockServer) GetRequestCount() int {
	return len(m.Requests)
}

// NewTestConfig creates a test configuration
func NewTestConfig(baseURL, apiKey string) *Config {
	return &Config{
		APIKey:             apiKey,
		BaseURL:            baseURL,
		Timeout:            DefaultTimeout,
		MaxRetries:         3,
		RetryBackoffFactor: 1.0,
		MaxBackoff:         DefaultMaxBackoff,
		VerifySSL:          true,
	}
}
