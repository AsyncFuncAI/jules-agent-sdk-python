package jules

import (
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid API key",
			apiKey:  "test-api-key",
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
			config, err := NewConfig(tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if config.APIKey != tt.apiKey {
					t.Errorf("NewConfig() APIKey = %v, want %v", config.APIKey, tt.apiKey)
				}
				if config.BaseURL != DefaultBaseURL {
					t.Errorf("NewConfig() BaseURL = %v, want %v", config.BaseURL, DefaultBaseURL)
				}
				if config.Timeout != DefaultTimeout {
					t.Errorf("NewConfig() Timeout = %v, want %v", config.Timeout, DefaultTimeout)
				}
				if config.MaxRetries != DefaultMaxRetries {
					t.Errorf("NewConfig() MaxRetries = %v, want %v", config.MaxRetries, DefaultMaxRetries)
				}
				if config.RetryBackoffFactor != DefaultRetryBackoffFactor {
					t.Errorf("NewConfig() RetryBackoffFactor = %v, want %v", config.RetryBackoffFactor, DefaultRetryBackoffFactor)
				}
				if config.MaxBackoff != DefaultMaxBackoff {
					t.Errorf("NewConfig() MaxBackoff = %v, want %v", config.MaxBackoff, DefaultMaxBackoff)
				}
				if !config.VerifySSL {
					t.Errorf("NewConfig() VerifySSL = %v, want true", config.VerifySSL)
				}
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
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
			},
			wantErr: false,
		},
		{
			name: "empty API key",
			config: &Config{
				APIKey:             "",
				BaseURL:            "https://example.com",
				Timeout:            30 * time.Second,
				MaxRetries:         3,
				RetryBackoffFactor: 1.0,
				MaxBackoff:         10 * time.Second,
			},
			wantErr: true,
			errMsg:  "API key is required",
		},
		{
			name: "empty base URL",
			config: &Config{
				APIKey:             "test-key",
				BaseURL:            "",
				Timeout:            30 * time.Second,
				MaxRetries:         3,
				RetryBackoffFactor: 1.0,
				MaxBackoff:         10 * time.Second,
			},
			wantErr: true,
			errMsg:  "base URL is required",
		},
		{
			name: "zero timeout",
			config: &Config{
				APIKey:             "test-key",
				BaseURL:            "https://example.com",
				Timeout:            0,
				MaxRetries:         3,
				RetryBackoffFactor: 1.0,
				MaxBackoff:         10 * time.Second,
			},
			wantErr: true,
			errMsg:  "timeout must be positive",
		},
		{
			name: "negative max retries",
			config: &Config{
				APIKey:             "test-key",
				BaseURL:            "https://example.com",
				Timeout:            30 * time.Second,
				MaxRetries:         -1,
				RetryBackoffFactor: 1.0,
				MaxBackoff:         10 * time.Second,
			},
			wantErr: true,
			errMsg:  "max retries must be non-negative",
		},
		{
			name: "negative backoff factor",
			config: &Config{
				APIKey:             "test-key",
				BaseURL:            "https://example.com",
				Timeout:            30 * time.Second,
				MaxRetries:         3,
				RetryBackoffFactor: -1.0,
				MaxBackoff:         10 * time.Second,
			},
			wantErr: true,
			errMsg:  "retry backoff factor must be non-negative",
		},
		{
			name: "zero max backoff",
			config: &Config{
				APIKey:             "test-key",
				BaseURL:            "https://example.com",
				Timeout:            30 * time.Second,
				MaxRetries:         3,
				RetryBackoffFactor: 1.0,
				MaxBackoff:         0,
			},
			wantErr: true,
			errMsg:  "max backoff must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Config.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestDefaultConstants(t *testing.T) {
	if DefaultBaseURL != "https://julius.googleapis.com/v1alpha" {
		t.Errorf("DefaultBaseURL = %v, want https://julius.googleapis.com/v1alpha", DefaultBaseURL)
	}
	if DefaultTimeout != 30*time.Second {
		t.Errorf("DefaultTimeout = %v, want 30s", DefaultTimeout)
	}
	if DefaultMaxRetries != 3 {
		t.Errorf("DefaultMaxRetries = %v, want 3", DefaultMaxRetries)
	}
	if DefaultRetryBackoffFactor != 1.0 {
		t.Errorf("DefaultRetryBackoffFactor = %v, want 1.0", DefaultRetryBackoffFactor)
	}
	if DefaultMaxBackoff != 10*time.Second {
		t.Errorf("DefaultMaxBackoff = %v, want 10s", DefaultMaxBackoff)
	}
	if DefaultPollInterval != 5*time.Second {
		t.Errorf("DefaultPollInterval = %v, want 5s", DefaultPollInterval)
	}
	if DefaultSessionTimeout != 600*time.Second {
		t.Errorf("DefaultSessionTimeout = %v, want 600s", DefaultSessionTimeout)
	}
}
