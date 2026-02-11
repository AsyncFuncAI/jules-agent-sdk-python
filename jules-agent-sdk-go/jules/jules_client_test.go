package jules

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid API key",
			apiKey:  "test-api-key-123",
			wantErr: false,
		},
		{
			name:    "empty API key",
			apiKey:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if client == nil {
					t.Error("NewClient() returned nil client")
					return
				}
				defer client.Close()

				// Verify API clients are initialized
				if client.Sessions == nil {
					t.Error("Sessions API is nil")
				}
				if client.Activities == nil {
					t.Error("Activities API is nil")
				}
				if client.Sources == nil {
					t.Error("Sources API is nil")
				}
			}
		})
	}
}

func TestNewClientWithConfig(t *testing.T) {
	config := &Config{
		APIKey:             "custom-api-key",
		BaseURL:            "https://custom.example.com",
		Timeout:            60 * time.Second,
		MaxRetries:         5,
		RetryBackoffFactor: 2.0,
		MaxBackoff:         30 * time.Second,
		VerifySSL:          false,
	}

	client, err := NewClientWithConfig(config)
	if err != nil {
		t.Fatalf("NewClientWithConfig() error = %v", err)
	}
	defer client.Close()

	if client == nil {
		t.Fatal("NewClientWithConfig() returned nil client")
	}

	// Verify API clients are initialized
	if client.Sessions == nil {
		t.Error("Sessions API is nil")
	}
	if client.Activities == nil {
		t.Error("Activities API is nil")
	}
	if client.Sources == nil {
		t.Error("Sources API is nil")
	}

	// Verify base client was configured correctly
	if client.baseClient == nil {
		t.Fatal("Base client is nil")
	}
}

func TestNewClientWithInvalidConfig(t *testing.T) {
	invalidConfig := &Config{
		APIKey:  "",
		BaseURL: "https://example.com",
		Timeout: 30 * time.Second,
	}

	client, err := NewClientWithConfig(invalidConfig)
	if err == nil {
		t.Error("NewClientWithConfig() expected error with invalid config, got nil")
		if client != nil {
			client.Close()
		}
	}
}

func TestClientClose(t *testing.T) {
	client, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	// Close should not return an error
	if err := client.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Calling Close multiple times should be safe
	if err := client.Close(); err != nil {
		t.Errorf("Second Close() error = %v", err)
	}
}

func TestClientStats(t *testing.T) {
	client, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer client.Close()

	stats := client.Stats()
	if stats == nil {
		t.Fatal("Stats() returned nil")
	}

	// Verify stats structure
	if _, ok := stats["request_count"]; !ok {
		t.Error("Stats missing request_count")
	}
	if _, ok := stats["error_count"]; !ok {
		t.Error("Stats missing error_count")
	}

	// Initial stats should be zero
	if stats["request_count"] != 0 {
		t.Errorf("Initial request_count = %v, want 0", stats["request_count"])
	}
	if stats["error_count"] != 0 {
		t.Errorf("Initial error_count = %v, want 0", stats["error_count"])
	}
}

func TestClientAPIsShareBaseClient(t *testing.T) {
	client, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer client.Close()

	// Verify all APIs share the same base client
	if client.Sessions.client != client.baseClient {
		t.Error("Sessions API does not share base client")
	}
	if client.Activities.client != client.baseClient {
		t.Error("Activities API does not share base client")
	}
	if client.Sources.client != client.baseClient {
		t.Error("Sources API does not share base client")
	}
}

func TestClientDefaultConfiguration(t *testing.T) {
	client, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer client.Close()

	// Verify default configuration values
	config := client.baseClient.config
	if config.BaseURL != DefaultBaseURL {
		t.Errorf("BaseURL = %v, want %v", config.BaseURL, DefaultBaseURL)
	}
	if config.Timeout != DefaultTimeout {
		t.Errorf("Timeout = %v, want %v", config.Timeout, DefaultTimeout)
	}
	if config.MaxRetries != DefaultMaxRetries {
		t.Errorf("MaxRetries = %v, want %v", config.MaxRetries, DefaultMaxRetries)
	}
	if config.RetryBackoffFactor != DefaultRetryBackoffFactor {
		t.Errorf("RetryBackoffFactor = %v, want %v", config.RetryBackoffFactor, DefaultRetryBackoffFactor)
	}
	if config.MaxBackoff != DefaultMaxBackoff {
		t.Errorf("MaxBackoff = %v, want %v", config.MaxBackoff, DefaultMaxBackoff)
	}
	if !config.VerifySSL {
		t.Error("VerifySSL should be true by default")
	}
}

func TestClientIntegration(t *testing.T) {
	// This is an integration test that verifies all APIs work together
	// Using a mock server
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"sources": []interface{}{
					map[string]interface{}{
						"id":   "source-1",
						"name": "sources/test-repo",
					},
				},
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-api-key")
	client, err := NewClientWithConfig(config)
	if err != nil {
		t.Fatalf("NewClientWithConfig() error = %v", err)
	}
	defer client.Close()

	// Make a request through the client
	ctx := context.Background()
	sources, err := client.Sources.List(ctx, nil)
	if err != nil {
		t.Fatalf("Sources.List() error = %v", err)
	}

	if len(sources.Sources) != 1 {
		t.Errorf("Sources count = %v, want 1", len(sources.Sources))
	}

	// Verify stats were updated
	stats := client.Stats()
	if stats["request_count"] != 1 {
		t.Errorf("request_count = %v, want 1", stats["request_count"])
	}
}

func TestClientNilConfig(t *testing.T) {
	// NewClientWithConfig should handle validation properly
	client, err := NewClientWithConfig(nil)
	if err == nil {
		t.Error("NewClientWithConfig(nil) expected error, got nil")
		if client != nil {
			client.Close()
		}
	}
}
