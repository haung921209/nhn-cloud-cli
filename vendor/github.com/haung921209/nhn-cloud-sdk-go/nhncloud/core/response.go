// Package core provides response parsing and error handling.
package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ResponseHeader is the standard NHN Cloud API response header
type ResponseHeader struct {
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
	IsSuccessful  bool   `json:"isSuccessful"`
}

// WithHeader interface for responses that have a header
type WithHeader interface {
	GetHeader() *ResponseHeader
}

// ParseResponse parses HTTP response into typed result
// This ensures complete response parsing - all fields must be in result struct
func ParseResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	// Read entire body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(body),
		}
	}

	// Unmarshal to result structure
	if err := json.Unmarshal(body, result); err != nil {
		return &ParseError{
			StatusCode: resp.StatusCode,
			Body:       string(body),
			Err:        err,
		}
	}

	// Check API-level error (if result has Header field)
	if withHeader, ok := result.(WithHeader); ok {
		header := withHeader.GetHeader()
		if header != nil {
			// Check for API error
			if header.ResultCode != 0 || !header.IsSuccessful {
				return &APIError{
					Code:    header.ResultCode,
					Message: header.ResultMessage,
				}
			}
		}
	}

	return nil
}
