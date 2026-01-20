package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	ErrorCode  string `json:"error_code,omitempty"`
}

func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("API Error %d [%s]: %s", e.StatusCode, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("API Error %d: %s", e.StatusCode, e.Message)
}

type Client struct {
	BaseURL       string
	HTTPClient    *http.Client
	TokenProvider TokenProvider
	Debug         bool
	UserAgent     string
}

type ClientOption func(*Client)

func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = hc
	}
}

func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.Debug = debug
	}
}

func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.UserAgent = ua
	}
}

func NewClient(baseURL string, tokenProvider TokenProvider, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL:       strings.TrimSuffix(baseURL, "/"),
		TokenProvider: tokenProvider,
		UserAgent:     "nhn-cloud-sdk-go/0.1.0",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Request(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}

	endpointPath := endpoint
	endpointQuery := ""
	if idx := strings.Index(endpoint, "?"); idx != -1 {
		endpointPath = endpoint[:idx]
		endpointQuery = endpoint[idx+1:]
	}
	u.Path = path.Join(u.Path, endpointPath)
	if endpointQuery != "" {
		u.RawQuery = endpointQuery
	}

	var reqBody io.Reader
	var contentType string

	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = strings.NewReader(v)
			contentType = "text/plain"
		case url.Values:
			reqBody = strings.NewReader(v.Encode())
			contentType = "application/x-www-form-urlencoded"
		case io.Reader:
			reqBody = v
			contentType = "application/octet-stream"
		default:
			jsonData, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}
			reqBody = bytes.NewReader(jsonData)
			contentType = "application/json"
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if c.TokenProvider != nil {
		token, err := c.TokenProvider.GetToken(ctx)
		if err != nil {
			return fmt.Errorf("failed to get access token: %w", err)
		}
		c.TokenProvider.SetAuthHeader(req, token)
	}

	if c.Debug {
		c.debugRequest(req, body)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if c.Debug {
		c.debugResponse(resp, respBody)
	}

	if resp.StatusCode >= 400 {
		apiError := &APIError{
			StatusCode: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
		}

		if len(respBody) > 0 {
			var errResp map[string]interface{}
			if json.Unmarshal(respBody, &errResp) == nil {
				if msg, ok := errResp["message"].(string); ok {
					apiError.Message = msg
				}
				if code, ok := errResp["error_code"].(string); ok {
					apiError.ErrorCode = code
				}
				if header, ok := errResp["header"].(map[string]interface{}); ok {
					if msg, ok := header["resultMessage"].(string); ok {
						apiError.Message = msg
					}
					if code, ok := header["resultCode"].(float64); ok {
						apiError.ErrorCode = fmt.Sprintf("%d", int(code))
					}
				}
			}
		}

		return apiError
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func (c *Client) GET(ctx context.Context, endpoint string, result interface{}) error {
	return c.Request(ctx, http.MethodGet, endpoint, nil, result)
}

func (c *Client) POST(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, endpoint, body, result)
}

func (c *Client) PUT(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPut, endpoint, body, result)
}

func (c *Client) PATCH(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPatch, endpoint, body, result)
}

func (c *Client) DELETE(ctx context.Context, endpoint string, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, endpoint, nil, result)
}

func (c *Client) DeleteWithBody(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, endpoint, body, result)
}

func (c *Client) debugRequest(req *http.Request, body interface{}) {
	fmt.Printf("=== REQUEST ===\n")
	fmt.Printf("%s %s\n", req.Method, req.URL.String())

	fmt.Printf("Headers:\n")
	for name, values := range req.Header {
		for _, value := range values {
			lowerName := strings.ToLower(name)
			if lowerName == "authorization" || lowerName == "x-nhn-authorization" || lowerName == "x-auth-token" {
				fmt.Printf("  %s: ***\n", name)
			} else {
				fmt.Printf("  %s: %s\n", name, value)
			}
		}
	}

	if body != nil {
		switch v := body.(type) {
		case string:
			fmt.Printf("Body: %s\n", v)
		case url.Values:
			fmt.Printf("Body: %s\n", v.Encode())
		default:
			if jsonData, err := json.MarshalIndent(body, "", "  "); err == nil {
				fmt.Printf("Body: %s\n", string(jsonData))
			} else {
				fmt.Printf("Body: %+v\n", body)
			}
		}
	}
	fmt.Printf("===============\n")
}

func (c *Client) debugResponse(resp *http.Response, body []byte) {
	fmt.Printf("=== RESPONSE ===\n")
	fmt.Printf("Status: %s\n", resp.Status)

	fmt.Printf("Headers:\n")
	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", name, value)
		}
	}

	if len(body) > 0 {
		var jsonData interface{}
		if json.Unmarshal(body, &jsonData) == nil {
			if prettyJSON, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
				fmt.Printf("Body: %s\n", string(prettyJSON))
			} else {
				fmt.Printf("Body: %s\n", string(body))
			}
		} else {
			fmt.Printf("Body: %s\n", string(body))
		}
	}
	fmt.Printf("================\n")
}
