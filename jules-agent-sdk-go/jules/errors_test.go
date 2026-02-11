package jules

import (
	"strings"
	"testing"
)

func TestAPIError(t *testing.T) {
	tests := []struct {
		name       string
		err        *APIError
		wantString string
	}{
		{
			name: "error with status code",
			err: &APIError{
				Message:    "test error",
				StatusCode: 400,
			},
			wantString: "Jules API error (status 400): test error",
		},
		{
			name: "error without status code",
			err: &APIError{
				Message:    "test error",
				StatusCode: 0,
			},
			wantString: "Jules API error: test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantString {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.wantString)
			}
		})
	}
}

func TestNewAuthenticationError(t *testing.T) {
	message := "authentication failed"
	response := map[string]interface{}{"error": "invalid token"}

	err := NewAuthenticationError(message, response)

	if err.Message != message {
		t.Errorf("NewAuthenticationError() Message = %v, want %v", err.Message, message)
	}
	if err.StatusCode != 401 {
		t.Errorf("NewAuthenticationError() StatusCode = %v, want 401", err.StatusCode)
	}
	if err.Response == nil {
		t.Error("NewAuthenticationError() Response is nil")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("NewAuthenticationError() Error() should contain '401', got %v", err.Error())
	}
}

func TestNewNotFoundError(t *testing.T) {
	message := "resource not found"
	response := map[string]interface{}{"error": "not found"}

	err := NewNotFoundError(message, response)

	if err.Message != message {
		t.Errorf("NewNotFoundError() Message = %v, want %v", err.Message, message)
	}
	if err.StatusCode != 404 {
		t.Errorf("NewNotFoundError() StatusCode = %v, want 404", err.StatusCode)
	}
	if err.Response == nil {
		t.Error("NewNotFoundError() Response is nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("NewNotFoundError() Error() should contain '404', got %v", err.Error())
	}
}

func TestNewValidationError(t *testing.T) {
	message := "validation failed"
	response := map[string]interface{}{"error": "invalid input"}

	err := NewValidationError(message, response)

	if err.Message != message {
		t.Errorf("NewValidationError() Message = %v, want %v", err.Message, message)
	}
	if err.StatusCode != 400 {
		t.Errorf("NewValidationError() StatusCode = %v, want 400", err.StatusCode)
	}
	if err.Response == nil {
		t.Error("NewValidationError() Response is nil")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("NewValidationError() Error() should contain '400', got %v", err.Error())
	}
}

func TestNewRateLimitError(t *testing.T) {
	message := "rate limit exceeded"
	retryAfter := 60
	response := map[string]interface{}{"error": "too many requests"}

	err := NewRateLimitError(message, retryAfter, response)

	if err.Message != message {
		t.Errorf("NewRateLimitError() Message = %v, want %v", err.Message, message)
	}
	if err.StatusCode != 429 {
		t.Errorf("NewRateLimitError() StatusCode = %v, want 429", err.StatusCode)
	}
	if err.RetryAfter != retryAfter {
		t.Errorf("NewRateLimitError() RetryAfter = %v, want %v", err.RetryAfter, retryAfter)
	}
	if err.Response == nil {
		t.Error("NewRateLimitError() Response is nil")
	}
	if !strings.Contains(err.Error(), "429") {
		t.Errorf("NewRateLimitError() Error() should contain '429', got %v", err.Error())
	}
}

func TestNewServerError(t *testing.T) {
	message := "internal server error"
	statusCode := 500
	response := map[string]interface{}{"error": "server error"}

	err := NewServerError(message, statusCode, response)

	if err.Message != message {
		t.Errorf("NewServerError() Message = %v, want %v", err.Message, message)
	}
	if err.StatusCode != statusCode {
		t.Errorf("NewServerError() StatusCode = %v, want %v", err.StatusCode, statusCode)
	}
	if err.Response == nil {
		t.Error("NewServerError() Response is nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("NewServerError() Error() should contain '500', got %v", err.Error())
	}
}

func TestNewTimeoutError(t *testing.T) {
	message := "operation timed out"

	err := NewTimeoutError(message)

	if err.Message != message {
		t.Errorf("NewTimeoutError() Message = %v, want %v", err.Message, message)
	}
	if !strings.Contains(err.Error(), "Timeout") {
		t.Errorf("NewTimeoutError() Error() should contain 'Timeout', got %v", err.Error())
	}
	if !strings.Contains(err.Error(), message) {
		t.Errorf("NewTimeoutError() Error() should contain message, got %v", err.Error())
	}
}

func TestErrorTypes(t *testing.T) {
	// Verify that all error types implement the error interface
	var _ error = &APIError{}
	var _ error = &AuthenticationError{}
	var _ error = &NotFoundError{}
	var _ error = &ValidationError{}
	var _ error = &RateLimitError{}
	var _ error = &ServerError{}
	var _ error = &TimeoutError{}
}
