package jules

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestNewBaseClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				APIKey:             "test-key",
				BaseURL:            "https://example.com",
				Timeout:            30 * time.Second,
				MaxRetries:         3,
				RetryBackoffFactor: 1.0,
				MaxBackoff:         10 * time.Second,
				VerifySSL:          true,
			},
			wantErr: false,
		},
		{
			name: "invalid config",
			config: &Config{
				APIKey:  "",
				BaseURL: "https://example.com",
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewBaseClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBaseClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewBaseClient() returned nil client")
			}
			if client != nil {
				client.Close()
			}
		})
	}
}

func TestBaseClientGet(t *testing.T) {
	// Create mock server
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: http.StatusOK,
			Body: map[string]interface{}{
				"id":   "test-123",
				"name": "test",
			},
		},
	})
	defer mockServer.Close()

	// Create client
	config := NewTestConfig(mockServer.URL, "test-key")
	client, err := NewBaseClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Make GET request
	ctx := context.Background()
	result, err := client.Get(ctx, "/test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	// Verify result
	if result["id"] != "test-123" {
		t.Errorf("Get() id = %v, want test-123", result["id"])
	}

	// Verify request
	req := mockServer.GetLastRequest()
	if req.Method != http.MethodGet {
		t.Errorf("Request method = %v, want GET", req.Method)
	}
	if req.Header.Get("X-Goog-Api-Key") != "test-key" {
		t.Errorf("API key header = %v, want test-key", req.Header.Get("X-Goog-Api-Key"))
	}
}

func TestBaseClientPost(t *testing.T) {
	// Create mock server
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: http.StatusOK,
			Body: map[string]interface{}{
				"id":      "session-123",
				"status":  "created",
			},
		},
	})
	defer mockServer.Close()

	// Create client
	config := NewTestConfig(mockServer.URL, "test-key")
	client, err := NewBaseClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Make POST request
	ctx := context.Background()
	body := map[string]interface{}{
		"prompt": "test prompt",
		"source": "sources/test",
	}
	result, err := client.Post(ctx, "/sessions", body)
	if err != nil {
		t.Fatalf("Post() error = %v", err)
	}

	// Verify result
	if result["id"] != "session-123" {
		t.Errorf("Post() id = %v, want session-123", result["id"])
	}

	// Verify request
	req := mockServer.GetLastRequest()
	if req.Method != http.MethodPost {
		t.Errorf("Request method = %v, want POST", req.Method)
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %v, want application/json", req.Header.Get("Content-Type"))
	}
}

func TestBaseClientErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   map[string]interface{}
		expectedErrType string
	}{
		{
			name:       "400 validation error",
			statusCode: 400,
			responseBody: map[string]interface{}{
				"error": map[string]interface{}{
					"message": "validation failed",
				},
			},
			expectedErrType: "*jules.ValidationError",
		},
		{
			name:       "401 authentication error",
			statusCode: 401,
			responseBody: map[string]interface{}{
				"error": map[string]interface{}{
					"message": "unauthorized",
				},
			},
			expectedErrType: "*jules.AuthenticationError",
		},
		{
			name:       "404 not found error",
			statusCode: 404,
			responseBody: map[string]interface{}{
				"error": map[string]interface{}{
					"message": "not found",
				},
			},
			expectedErrType: "*jules.NotFoundError",
		},
		{
			name:       "429 rate limit error",
			statusCode: 429,
			responseBody: map[string]interface{}{
				"error": map[string]interface{}{
					"message": "rate limit exceeded",
				},
			},
			expectedErrType: "*jules.RateLimitError",
		},
		{
			name:       "500 server error",
			statusCode: 500,
			responseBody: map[string]interface{}{
				"error": map[string]interface{}{
					"message": "internal server error",
				},
			},
			expectedErrType: "*jules.ServerError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			mockServer := NewMockServer(t, []MockResponse{
				{
					StatusCode: tt.statusCode,
					Body:       tt.responseBody,
				},
			})
			defer mockServer.Close()

			// Create client with no retries for error tests
			config := NewTestConfig(mockServer.URL, "test-key")
			config.MaxRetries = 0
			client, err := NewBaseClient(config)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}
			defer client.Close()

			// Make request
			ctx := context.Background()
			_, err = client.Get(ctx, "/test")

			// Verify error type
			if err == nil {
				t.Errorf("Expected error, got nil")
				return
			}

			switch tt.expectedErrType {
			case "*jules.ValidationError":
				if _, ok := err.(*ValidationError); !ok {
					t.Errorf("Expected ValidationError, got %T", err)
				}
			case "*jules.AuthenticationError":
				if _, ok := err.(*AuthenticationError); !ok {
					t.Errorf("Expected AuthenticationError, got %T", err)
				}
			case "*jules.NotFoundError":
				if _, ok := err.(*NotFoundError); !ok {
					t.Errorf("Expected NotFoundError, got %T", err)
				}
			case "*jules.RateLimitError":
				if _, ok := err.(*RateLimitError); !ok {
					t.Errorf("Expected RateLimitError, got %T", err)
				}
			case "*jules.ServerError":
				if _, ok := err.(*ServerError); !ok {
					t.Errorf("Expected ServerError, got %T", err)
				}
			}
		})
	}
}

func TestBaseClientRetry(t *testing.T) {
	// Create mock server that fails twice then succeeds
	mockServer := NewMockServer(t, []MockResponse{
		{StatusCode: 500, Body: map[string]interface{}{"error": "server error"}},
		{StatusCode: 500, Body: map[string]interface{}{"error": "server error"}},
		{StatusCode: 200, Body: map[string]interface{}{"id": "success"}},
	})
	defer mockServer.Close()

	// Create client with retries
	config := NewTestConfig(mockServer.URL, "test-key")
	config.MaxRetries = 3
	config.RetryBackoffFactor = 0.1 // Fast retries for testing
	client, err := NewBaseClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Make request
	ctx := context.Background()
	result, err := client.Get(ctx, "/test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	// Verify success
	if result["id"] != "success" {
		t.Errorf("Get() id = %v, want success", result["id"])
	}

	// Verify retry count
	if mockServer.GetRequestCount() != 3 {
		t.Errorf("Request count = %v, want 3", mockServer.GetRequestCount())
	}
}

func TestBaseClientStats(t *testing.T) {
	// Create mock server
	mockServer := NewMockServer(t, []MockResponse{
		{StatusCode: 200, Body: map[string]interface{}{"id": "1"}},
		{StatusCode: 200, Body: map[string]interface{}{"id": "2"}},
		{StatusCode: 500, Body: map[string]interface{}{"error": "error"}},
	})
	defer mockServer.Close()

	// Create client
	config := NewTestConfig(mockServer.URL, "test-key")
	config.MaxRetries = 0
	client, err := NewBaseClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Make successful requests
	client.Get(ctx, "/test1")
	client.Get(ctx, "/test2")

	// Make failed request
	client.Get(ctx, "/test3")

	// Verify stats
	stats := client.Stats()
	if stats["request_count"] != 3 {
		t.Errorf("request_count = %v, want 3", stats["request_count"])
	}
	if stats["error_count"] != 1 {
		t.Errorf("error_count = %v, want 1", stats["error_count"])
	}
}

func TestCalculateBackoff(t *testing.T) {
	config := &Config{
		APIKey:             "test",
		BaseURL:            "https://example.com",
		Timeout:            30 * time.Second,
		MaxRetries:         3,
		RetryBackoffFactor: 1.0,
		MaxBackoff:         10 * time.Second,
	}

	client, _ := NewBaseClient(config)
	defer client.Close()

	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{1, 1 * time.Second},
		{2, 2 * time.Second},
		{3, 4 * time.Second},
		{4, 8 * time.Second},
		{5, 10 * time.Second}, // Capped at MaxBackoff
		{10, 10 * time.Second}, // Capped at MaxBackoff
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := client.calculateBackoff(tt.attempt)
			if got != tt.want {
				t.Errorf("calculateBackoff(%d) = %v, want %v", tt.attempt, got, tt.want)
			}
		})
	}
}

func TestBaseClientContextCancellation(t *testing.T) {
	// Create mock server with delay
	mockServer := NewMockServer(t, []MockResponse{
		{StatusCode: 200, Body: map[string]interface{}{"id": "test"}},
	})
	defer mockServer.Close()

	// Create client
	config := NewTestConfig(mockServer.URL, "test-key")
	client, err := NewBaseClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Make request with cancelled context
	_, err = client.Get(ctx, "/test")
	if err == nil {
		t.Error("Expected error with cancelled context, got nil")
	}
}

func TestBaseClientJSONParsing(t *testing.T) {
	// Create mock server with complex JSON
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"nested": map[string]interface{}{
					"value": "test",
					"count": float64(123),
				},
				"array": []interface{}{"a", "b", "c"},
			},
		},
	})
	defer mockServer.Close()

	// Create client
	config := NewTestConfig(mockServer.URL, "test-key")
	client, err := NewBaseClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Make request
	ctx := context.Background()
	result, err := client.Get(ctx, "/test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	// Verify nested object
	nested, ok := result["nested"].(map[string]interface{})
	if !ok {
		t.Fatal("nested is not a map")
	}
	if nested["value"] != "test" {
		t.Errorf("nested.value = %v, want test", nested["value"])
	}

	// Verify array
	arr, ok := result["array"].([]interface{})
	if !ok {
		t.Fatal("array is not a slice")
	}
	if len(arr) != 3 {
		t.Errorf("array length = %v, want 3", len(arr))
	}
}
