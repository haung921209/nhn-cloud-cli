// Package capture provides an opt-in http.RoundTripper that mirrors
// each response body to a directory specified by NHN_SDK_CAPTURE_DIR.
// When the env var is empty, RoundTrip is a transparent passthrough.
//
// File naming: <METHOD>_<path-slug>_<unix-nanos>.json
// where path-slug is the URL path with leading/trailing slashes trimmed
// and remaining slashes replaced with underscores (truncated to 80 chars).
//
// Capture is best-effort: any failure to write the capture file is
// silently ignored so that SDK callers never observe capture-related
// errors. The response body remains readable to the caller because we
// read it once and replace it with bytes.NewReader(body).
package capture

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EnvVar is the environment variable that, when non-empty, enables capture.
const EnvVar = "NHN_SDK_CAPTURE_DIR"

// Transport wraps an inner http.RoundTripper. When NHN_SDK_CAPTURE_DIR is
// set, it mirrors response bodies to that directory.
type Transport struct {
	Inner http.RoundTripper
}

// NewTransport returns a Transport wrapping inner. If inner is nil,
// http.DefaultTransport is used.
func NewTransport(inner http.RoundTripper) *Transport {
	if inner == nil {
		inner = http.DefaultTransport
	}
	return &Transport{Inner: inner}
}

// WrapTransport is a convenience wrapper used by code paths that want a
// plain http.RoundTripper return type without a concrete *Transport.
func WrapTransport(inner http.RoundTripper) http.RoundTripper {
	return NewTransport(inner)
}

// WrapClient ensures c.Transport is wrapped with a capture Transport.
// Safe to call on a nil client (returns a fresh client) or a client whose
// Transport is already a *Transport (returns c unmodified).
func WrapClient(c *http.Client) *http.Client {
	if c == nil {
		return &http.Client{Transport: NewTransport(http.DefaultTransport)}
	}
	if _, alreadyWrapped := c.Transport.(*Transport); alreadyWrapped {
		return c
	}
	c.Transport = NewTransport(c.Transport)
	return c
}

// RoundTrip implements http.RoundTripper.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Inner.RoundTrip(req)
	dir := os.Getenv(EnvVar)
	if err != nil || dir == "" || resp == nil || resp.Body == nil {
		return resp, err
	}

	body, readErr := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	// Always replace the body so the caller sees consistent semantics
	// regardless of whether the read succeeded.
	resp.Body = io.NopCloser(bytes.NewReader(body))
	if readErr != nil {
		return resp, nil
	}

	slug := strings.ReplaceAll(strings.Trim(req.URL.Path, "/"), "/", "_")
	if slug == "" {
		slug = "root"
	}
	if len(slug) > 80 {
		slug = slug[:80]
	}
	name := fmt.Sprintf("%s_%s_%d.json", req.Method, slug, time.Now().UnixNano())
	if mkErr := os.MkdirAll(dir, 0o755); mkErr == nil {
		_ = os.WriteFile(filepath.Join(dir, name), body, 0o644)
	}
	return resp, nil
}
