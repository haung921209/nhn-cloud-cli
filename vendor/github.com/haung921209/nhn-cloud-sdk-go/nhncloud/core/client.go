// Package core provides common HTTP client functionality for NHN Cloud SDK.
package core

import (
	"context"
	"net/http"
	"time"
)

// Authenticator handles request authentication
type Authenticator interface {
	Authenticate(req *http.Request) error
}

// Client is the base HTTP client for all API calls
type Client struct {
	httpClient *http.Client
	baseURL    string
	auth       Authenticator
	options    ClientOptions
}

// ClientOptions configures the client behavior
type ClientOptions struct {
	Timeout      time.Duration
	MaxRetries   int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
	UserAgent    string
	Debug        bool
}

// DefaultClientOptions returns sensible default options
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
		UserAgent:    "nhn-cloud-sdk-go/2.0.0",
		Debug:        false,
	}
}

// NewClient creates a new HTTP client
func NewClient(baseURL string, auth Authenticator, opts *ClientOptions) *Client {
	if opts == nil {
		opts = DefaultClientOptions()
	}

	return &Client{
		httpClient: &http.Client{Timeout: opts.Timeout},
		baseURL:    baseURL,
		auth:       auth,
		options:    *opts,
	}
}

// Do executes an HTTP request with authentication
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Set full URL
	req.URL.Scheme = "https"
	req.URL.Host = c.baseURL

	// Add authentication
	if c.auth != nil {
		if err := c.auth.Authenticate(req); err != nil {
			return nil, err
		}
	}

	// Add common headers
	req.Header.Set("User-Agent", c.options.UserAgent)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	return c.httpClient.Do(req)
}

// BaseURL returns the base URL
func (c *Client) BaseURL() string {
	return c.baseURL
}
