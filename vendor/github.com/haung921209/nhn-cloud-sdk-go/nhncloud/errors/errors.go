// Package errors provides typed error handling for NHN Cloud SDK.
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError represents an error returned by the NHN Cloud API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	RequestID  string
	Retryable  bool
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("nhncloud: %s (code=%s, status=%d)", e.Message, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("nhncloud: %s (status=%d)", e.Message, e.StatusCode)
}

// NotFoundError indicates the requested resource was not found.
type NotFoundError struct {
	Resource   string
	ResourceID string
	APIError
}

func (e *NotFoundError) Error() string {
	if e.ResourceID != "" {
		return fmt.Sprintf("nhncloud: %s '%s' not found", e.Resource, e.ResourceID)
	}
	return fmt.Sprintf("nhncloud: %s not found", e.Resource)
}

// AuthenticationError indicates an authentication failure.
type AuthenticationError struct {
	APIError
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("nhncloud: authentication failed: %s", e.Message)
}

// RateLimitError indicates the request was rate limited.
type RateLimitError struct {
	RetryAfter int // seconds
	APIError
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("nhncloud: rate limited, retry after %d seconds", e.RetryAfter)
}

// ValidationError indicates invalid request parameters.
type ValidationError struct {
	Field  string
	Reason string
	APIError
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("nhncloud: validation failed for '%s': %s", e.Field, e.Reason)
	}
	return fmt.Sprintf("nhncloud: validation failed: %s", e.Message)
}

// NetworkError indicates a network-level failure.
type NetworkError struct {
	Cause error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("nhncloud: network error: %v", e.Cause)
}

func (e *NetworkError) Unwrap() error {
	return e.Cause
}

// TimeoutError indicates the request timed out.
type TimeoutError struct {
	Cause error
}

func (e *TimeoutError) Error() string {
	return "nhncloud: request timed out"
}

func (e *TimeoutError) Unwrap() error {
	return e.Cause
}

// --- Helper functions for error checking ---

// IsNotFound returns true if the error indicates a resource was not found.
func IsNotFound(err error) bool {
	var notFound *NotFoundError
	if errors.As(err, &notFound) {
		return true
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsAuthentication returns true if the error indicates an authentication failure.
func IsAuthentication(err error) bool {
	var authErr *AuthenticationError
	if errors.As(err, &authErr) {
		return true
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized || apiErr.StatusCode == http.StatusForbidden
	}
	return false
}

// IsRateLimited returns true if the error indicates rate limiting.
func IsRateLimited(err error) bool {
	var rateErr *RateLimitError
	if errors.As(err, &rateErr) {
		return true
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// IsRetryable returns true if the error is potentially retryable.
func IsRetryable(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Retryable
	}
	var netErr *NetworkError
	if errors.As(err, &netErr) {
		return true
	}
	var timeoutErr *TimeoutError
	return errors.As(err, &timeoutErr)
}

// IsValidation returns true if the error indicates a validation failure.
func IsValidation(err error) bool {
	var valErr *ValidationError
	return errors.As(err, &valErr)
}

// IsTimeout returns true if the error indicates a timeout.
func IsTimeout(err error) bool {
	var timeoutErr *TimeoutError
	return errors.As(err, &timeoutErr)
}

// --- Error construction from HTTP response ---

// FromHTTPResponse creates an appropriate error from an HTTP response.
func FromHTTPResponse(statusCode int, code, message, requestID string) error {
	baseErr := APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		RequestID:  requestID,
		Retryable:  statusCode >= 500 || statusCode == http.StatusTooManyRequests,
	}

	switch statusCode {
	case http.StatusNotFound:
		return &NotFoundError{APIError: baseErr}
	case http.StatusUnauthorized, http.StatusForbidden:
		return &AuthenticationError{APIError: baseErr}
	case http.StatusTooManyRequests:
		return &RateLimitError{APIError: baseErr}
	case http.StatusBadRequest:
		return &ValidationError{APIError: baseErr}
	default:
		return &baseErr
	}
}
