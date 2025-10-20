package jules

import (
	"fmt"
	"time"
)

const (
	// DefaultBaseURL is the default API base URL
	DefaultBaseURL = "https://julius.googleapis.com/v1alpha"

	// DefaultTimeout is the default request timeout
	DefaultTimeout = 30 * time.Second

	// DefaultMaxRetries is the default maximum number of retries
	DefaultMaxRetries = 3

	// DefaultRetryBackoffFactor is the default exponential backoff factor
	DefaultRetryBackoffFactor = 1.0

	// DefaultMaxBackoff is the default maximum backoff duration
	DefaultMaxBackoff = 10 * time.Second

	// DefaultPollInterval is the default polling interval for wait operations
	DefaultPollInterval = 5 * time.Second

	// DefaultSessionTimeout is the default timeout for session completion
	DefaultSessionTimeout = 600 * time.Second
)

// Config holds the configuration for the Jules API client
type Config struct {
	// APIKey is the API key for authentication (required)
	APIKey string

	// BaseURL is the base URL for the API
	BaseURL string

	// Timeout is the HTTP request timeout
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts
	MaxRetries int

	// RetryBackoffFactor is the exponential backoff factor for retries
	RetryBackoffFactor float64

	// MaxBackoff is the maximum backoff duration between retries
	MaxBackoff time.Duration

	// VerifySSL controls SSL certificate verification
	VerifySSL bool
}

// NewConfig creates a new Config with default values
func NewConfig(apiKey string) (*Config, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	return &Config{
		APIKey:             apiKey,
		BaseURL:            DefaultBaseURL,
		Timeout:            DefaultTimeout,
		MaxRetries:         DefaultMaxRetries,
		RetryBackoffFactor: DefaultRetryBackoffFactor,
		MaxBackoff:         DefaultMaxBackoff,
		VerifySSL:          true,
	}, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	if c.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if c.MaxRetries < 0 {
		return fmt.Errorf("max retries must be non-negative")
	}
	if c.RetryBackoffFactor < 0 {
		return fmt.Errorf("retry backoff factor must be non-negative")
	}
	if c.MaxBackoff <= 0 {
		return fmt.Errorf("max backoff must be positive")
	}
	return nil
}
