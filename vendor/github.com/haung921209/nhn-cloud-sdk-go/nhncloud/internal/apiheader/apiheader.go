// Package apiheader reads the result header that NHN Cloud application
// services put in every response body, success or not:
// {"header": {"isSuccessful": false, "resultCode": N, "resultMessage": "..."}}.
// Most of these APIs answer HTTP 200 even on failure, so the header is the
// only signal.
package apiheader

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Error is a response whose header says isSuccessful=false.
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("API error %s: %s", e.Code, e.Message)
}

// Check returns an *Error when body carries a header with isSuccessful=false,
// and nil otherwise (no body, not JSON, no header, or a successful one).
func Check(body []byte) error {
	var r struct {
		Header *struct {
			IsSuccessful  *bool           `json:"isSuccessful"`
			ResultCode    json.RawMessage `json:"resultCode"`
			ResultMessage string          `json:"resultMessage"`
		} `json:"header"`
	}
	if len(body) == 0 || json.Unmarshal(body, &r) != nil || r.Header == nil || r.Header.IsSuccessful == nil || *r.Header.IsSuccessful {
		return nil
	}
	// resultCode is a number in most services and a string in some.
	return &Error{Code: strings.Trim(string(r.Header.ResultCode), `"`), Message: r.Header.ResultMessage}
}
