// Package transport provides HTTP transport with retry and middleware support.
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/errors"
)

type TokenProvider interface {
	GetBearerToken() (string, error)
}

// Client is an HTTP client with retry support and middleware.
type Client struct {
	httpClient    *http.Client
	baseURL       string
	headers       map[string]string
	tokenProvider TokenProvider

	maxAttempts       int
	initialBackoff    time.Duration
	maxBackoff        time.Duration
	backoffMultiplier float64
	retryableCodes    map[int]bool

	debug bool
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// NewClient creates a new HTTP client with the given options.
func NewClient(baseURL string, opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:           strings.TrimSuffix(baseURL, "/"),
		headers:           make(map[string]string),
		maxAttempts:       3,
		initialBackoff:    100 * time.Millisecond,
		maxBackoff:        5 * time.Second,
		backoffMultiplier: 2.0,
		retryableCodes: map[int]bool{
			408: true, 429: true, 500: true, 502: true, 503: true, 504: true,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// WithRetry configures retry behavior.
func WithRetry(maxAttempts int, initialBackoff, maxBackoff time.Duration) ClientOption {
	return func(c *Client) {
		c.maxAttempts = maxAttempts
		c.initialBackoff = initialBackoff
		c.maxBackoff = maxBackoff
	}
}

// WithoutRetry disables retries.
func WithoutRetry() ClientOption {
	return func(c *Client) {
		c.maxAttempts = 1
	}
}

// WithHeader adds a default header to all requests.
func WithHeader(key, value string) ClientOption {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// WithDebug enables debug logging.
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithAppKeyAuth sets App Key based authentication headers (for RDS APIs).
func WithAppKeyAuth(appKey, accessKeyID, secretAccessKey string) ClientOption {
	return func(c *Client) {
		c.headers["X-TC-APP-KEY"] = appKey
		c.headers["X-TC-AUTHENTICATION-ID"] = accessKeyID
		c.headers["X-TC-AUTHENTICATION-SECRET"] = secretAccessKey
	}
}

// WithBearerAuth sets Bearer token authentication.
func WithBearerAuth(token string) ClientOption {
	return func(c *Client) {
		c.headers["Authorization"] = "Bearer " + token
	}
}

// WithDynamicBearerAuth sets dynamic Bearer token authentication using a token provider.
func WithDynamicBearerAuth(appKey string, provider TokenProvider) ClientOption {
	return func(c *Client) {
		c.headers["X-TC-APP-KEY"] = appKey
		c.tokenProvider = provider
	}
}

// Request represents an HTTP request to be executed.
type Request struct {
	Method  string
	Path    string
	Query   url.Values
	Body    interface{}
	Headers map[string]string
}

// Response represents an HTTP response.
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// Do executes an HTTP request with retry logic.
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	var lastErr error
	backoff := c.initialBackoff

	for attempt := 1; attempt <= c.maxAttempts; attempt++ {
		resp, err := c.doOnce(ctx, req, attempt)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Check if retryable
		if !c.shouldRetry(err, attempt) {
			return nil, err
		}

		// Wait before retry
		if attempt < c.maxAttempts {
			select {
			case <-ctx.Done():
				return nil, &errors.TimeoutError{Cause: ctx.Err()}
			case <-time.After(c.jitteredBackoff(backoff)):
			}
			backoff = time.Duration(float64(backoff) * c.backoffMultiplier)
			if backoff > c.maxBackoff {
				backoff = c.maxBackoff
			}
		}
	}

	return nil, lastErr
}

func (c *Client) doOnce(ctx context.Context, req *Request, attempt int) (*Response, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	reqPath := req.Path
	var pathQuery url.Values
	if idx := strings.Index(reqPath, "?"); idx != -1 {
		queryStr := reqPath[idx+1:]
		reqPath = reqPath[:idx]
		pathQuery, _ = url.ParseQuery(queryStr)
	}

	u.Path = path.Join(u.Path, reqPath)

	query := make(url.Values)
	for k, v := range pathQuery {
		query[k] = v
	}
	for k, v := range req.Query {
		query[k] = v
	}
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	// Prepare body
	var bodyReader io.Reader
	if req.Body != nil {
		jsonData, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "nhn-cloud-sdk-go/1.0.0")

	for k, v := range c.headers {
		httpReq.Header.Set(k, v)
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Set dynamic bearer token if provider exists (for OAuth2-based services like PostgreSQL)
	if c.tokenProvider != nil {
		token, err := c.tokenProvider.GetBearerToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get bearer token: %w", err)
		}
		httpReq.Header.Set("X-NHN-AUTHORIZATION", "Bearer "+token)
	}

	// Debug logging
	if c.debug {
		c.logRequest(httpReq, req.Body, attempt)
	}

	// Execute request
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, &errors.NetworkError{Cause: err}
	}
	defer httpResp.Body.Close()

	// Read response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Debug logging
	if c.debug {
		c.logResponse(httpResp, body)
	}

	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       body,
	}

	// Handle HTTP error responses
	if httpResp.StatusCode >= 400 {
		return resp, c.parseErrorResponse(resp)
	}

	// Handle API-level errors (HTTP 200 but isSuccessful=false)
	// NHN Cloud APIs return HTTP 200 with error details in body
	if len(body) > 0 {
		if apiErr := c.checkAPIError(resp); apiErr != nil {
			return resp, apiErr
		}
	}

	return resp, nil
}

func (c *Client) shouldRetry(err error, attempt int) bool {
	if attempt >= c.maxAttempts {
		return false
	}
	return errors.IsRetryable(err)
}

func (c *Client) jitteredBackoff(base time.Duration) time.Duration {
	// Add up to 20% jitter
	jitter := time.Duration(rand.Float64() * 0.2 * float64(base))
	return base + jitter
}

func (c *Client) parseErrorResponse(resp *Response) error {
	var apiResp struct {
		Header struct {
			ResultCode    int    `json:"resultCode"`
			ResultMessage string `json:"resultMessage"`
			IsSuccessful  bool   `json:"isSuccessful"`
		} `json:"header"`
		Message   string `json:"message"`
		ErrorCode string `json:"error_code"`
	}

	message := http.StatusText(resp.StatusCode)
	code := ""

	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, &apiResp); err == nil {
			if apiResp.Header.ResultMessage != "" {
				message = apiResp.Header.ResultMessage
			} else if apiResp.Message != "" {
				message = apiResp.Message
			}
			code = apiResp.ErrorCode
		}
	}

	requestID := resp.Headers.Get("X-Request-Id")
	return errors.FromHTTPResponse(resp.StatusCode, code, message, requestID)
}

func (c *Client) checkAPIError(resp *Response) error {
	var apiResp struct {
		Header struct {
			ResultCode    int    `json:"resultCode"`
			ResultMessage string `json:"resultMessage"`
			IsSuccessful  bool   `json:"isSuccessful"`
		} `json:"header"`
	}

	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil
	}

	if apiResp.Header.ResultCode != 0 && !apiResp.Header.IsSuccessful {
		requestID := resp.Headers.Get("X-Request-Id")
		return errors.FromHTTPResponse(resp.StatusCode, fmt.Sprintf("%d", apiResp.Header.ResultCode), apiResp.Header.ResultMessage, requestID)
	}

	return nil
}

func (c *Client) logRequest(req *http.Request, body interface{}, attempt int) {
	fmt.Printf("=== REQUEST (attempt %d) ===\n", attempt)
	fmt.Printf("%s %s\n", req.Method, req.URL.String())
	fmt.Printf("Headers:\n")
	for name, values := range req.Header {
		for _, value := range values {
			// Mask sensitive headers
			lower := strings.ToLower(name)
			if strings.Contains(lower, "secret") || strings.Contains(lower, "key") || strings.Contains(lower, "auth") {
				fmt.Printf("  %s: ***\n", name)
			} else {
				fmt.Printf("  %s: %s\n", name, value)
			}
		}
	}
	if body != nil {
		if jsonData, err := json.MarshalIndent(body, "", "  "); err == nil {
			fmt.Printf("Body: %s\n", string(jsonData))
		}
	}
	fmt.Printf("===========================\n")
}

func (c *Client) logResponse(resp *http.Response, body []byte) {
	fmt.Printf("=== RESPONSE ===\n")
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Headers:\n")
	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", name, value)
		}
	}
	if len(body) > 0 {
		// Truncate very long responses
		if len(body) > 2000 {
			fmt.Printf("Body (truncated): %s...\n", string(body[:2000]))
		} else {
			fmt.Printf("Body: %s\n", string(body))
		}
	}
	fmt.Printf("================\n")
}

// --- Convenience methods ---

// GET performs a GET request.
func (c *Client) GET(ctx context.Context, path string, result interface{}) error {
	resp, err := c.Do(ctx, &Request{Method: "GET", Path: path})
	if err != nil {
		return err
	}
	if result != nil && len(resp.Body) > 0 {
		return json.Unmarshal(resp.Body, result)
	}
	return nil
}

// POST performs a POST request.
func (c *Client) POST(ctx context.Context, path string, body, result interface{}) error {
	resp, err := c.Do(ctx, &Request{Method: "POST", Path: path, Body: body})
	if err != nil {
		return err
	}
	if result != nil && len(resp.Body) > 0 {
		return json.Unmarshal(resp.Body, result)
	}
	return nil
}

// PUT performs a PUT request.
func (c *Client) PUT(ctx context.Context, path string, body, result interface{}) error {
	resp, err := c.Do(ctx, &Request{Method: "PUT", Path: path, Body: body})
	if err != nil {
		return err
	}
	if result != nil && len(resp.Body) > 0 {
		return json.Unmarshal(resp.Body, result)
	}
	return nil
}

// DELETE performs a DELETE request.
func (c *Client) DELETE(ctx context.Context, path string, result interface{}) error {
	resp, err := c.Do(ctx, &Request{Method: "DELETE", Path: path})
	if err != nil {
		return err
	}
	if result != nil && len(resp.Body) > 0 {
		return json.Unmarshal(resp.Body, result)
	}
	return nil
}

// --- Pagination support ---

// Paginator provides iteration over paginated API responses.
type Paginator[T any] struct {
	client  *Client
	path    string
	query   url.Values
	page    int
	size    int
	hasMore bool
	extract func([]byte) ([]T, int, error) // extracts items and total count
}

// NewPaginator creates a new paginator.
func NewPaginator[T any](client *Client, path string, pageSize int, extractor func([]byte) ([]T, int, error)) *Paginator[T] {
	return &Paginator[T]{
		client:  client,
		path:    path,
		query:   make(url.Values),
		page:    0,
		size:    pageSize,
		hasMore: true,
		extract: extractor,
	}
}

// HasMorePages returns true if there are more pages to fetch.
func (p *Paginator[T]) HasMorePages() bool {
	return p.hasMore
}

// NextPage fetches the next page of results.
func (p *Paginator[T]) NextPage(ctx context.Context) ([]T, error) {
	if !p.hasMore {
		return nil, nil
	}

	p.query.Set("page", fmt.Sprintf("%d", p.page))
	p.query.Set("size", fmt.Sprintf("%d", p.size))

	resp, err := p.client.Do(ctx, &Request{
		Method: "GET",
		Path:   p.path,
		Query:  p.query,
	})
	if err != nil {
		return nil, err
	}

	items, total, err := p.extract(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to extract page items: %w", err)
	}

	p.page++
	// Check if we've fetched all items
	fetchedSoFar := p.page * p.size
	p.hasMore = fetchedSoFar < total && len(items) > 0

	return items, nil
}

// All fetches all pages and returns combined results.
func (p *Paginator[T]) All(ctx context.Context) ([]T, error) {
	var all []T
	for p.HasMorePages() {
		items, err := p.NextPage(ctx)
		if err != nil {
			return all, err
		}
		all = append(all, items...)
	}
	return all, nil
}

// unused but may be useful
var _ = math.MaxInt
