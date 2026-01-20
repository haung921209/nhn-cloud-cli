// Package core provides error types for the SDK.
package core

import "fmt"

// HTTPError represents an HTTP-level error
type HTTPError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s, body: %s", e.StatusCode, e.Status, e.Body)
}

// APIError represents an API-level error (successful HTTP but failed API call)
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

// ParseError represents a JSON parsing error
type ParseError struct {
	StatusCode int
	Body       string
	Err        error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("failed to parse response (HTTP %d): %v, body: %s",
		e.StatusCode, e.Err, e.Body)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// ValidationError represents request validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}
