package jules

import "fmt"

// APIError is the base error type for all Jules API errors
type APIError struct {
	Message    string
	StatusCode int
	Response   map[string]interface{}
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("Jules API error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("Jules API error: %s", e.Message)
}

// AuthenticationError represents a 401 authentication error
type AuthenticationError struct {
	*APIError
}

// NewAuthenticationError creates a new AuthenticationError
func NewAuthenticationError(message string, response map[string]interface{}) *AuthenticationError {
	return &AuthenticationError{
		APIError: &APIError{
			Message:    message,
			StatusCode: 401,
			Response:   response,
		},
	}
}

// NotFoundError represents a 404 not found error
type NotFoundError struct {
	*APIError
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(message string, response map[string]interface{}) *NotFoundError {
	return &NotFoundError{
		APIError: &APIError{
			Message:    message,
			StatusCode: 404,
			Response:   response,
		},
	}
}

// ValidationError represents a 400 validation error
type ValidationError struct {
	*APIError
}

// NewValidationError creates a new ValidationError
func NewValidationError(message string, response map[string]interface{}) *ValidationError {
	return &ValidationError{
		APIError: &APIError{
			Message:    message,
			StatusCode: 400,
			Response:   response,
		},
	}
}

// RateLimitError represents a 429 rate limit error
type RateLimitError struct {
	*APIError
	RetryAfter int // Retry-After header value in seconds
}

// NewRateLimitError creates a new RateLimitError
func NewRateLimitError(message string, retryAfter int, response map[string]interface{}) *RateLimitError {
	return &RateLimitError{
		APIError: &APIError{
			Message:    message,
			StatusCode: 429,
			Response:   response,
		},
		RetryAfter: retryAfter,
	}
}

// ServerError represents a 5xx server error
type ServerError struct {
	*APIError
}

// NewServerError creates a new ServerError
func NewServerError(message string, statusCode int, response map[string]interface{}) *ServerError {
	return &ServerError{
		APIError: &APIError{
			Message:    message,
			StatusCode: statusCode,
			Response:   response,
		},
	}
}

// TimeoutError represents a timeout error
type TimeoutError struct {
	Message string
}

// Error implements the error interface
func (e *TimeoutError) Error() string {
	return fmt.Sprintf("Timeout: %s", e.Message)
}

// NewTimeoutError creates a new TimeoutError
func NewTimeoutError(message string) *TimeoutError {
	return &TimeoutError{Message: message}
}
