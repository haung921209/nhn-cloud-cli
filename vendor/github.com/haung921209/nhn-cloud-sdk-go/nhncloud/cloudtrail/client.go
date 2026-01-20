// Package cloudtrail provides CloudTrail service client
package cloudtrail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultBaseURL = "https://cloud-trail.api.nhncloudservice.com"

// Client represents a CloudTrail API client
type Client struct {
	baseURL     string
	appKey      string
	accessKeyID string
	secretKey   string
	httpClient  *http.Client
	debug       bool
	useV2       bool // Use v2.0 API (requires user auth)
}

// NewClient creates a new CloudTrail client
func NewClient(appKey, accessKeyID, secretKey string, httpClient *http.Client, debug bool) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{
		baseURL:     DefaultBaseURL,
		appKey:      appKey,
		accessKeyID: accessKeyID,
		secretKey:   secretKey,
		httpClient:  httpClient,
		debug:       debug,
		useV2:       true, // Default to v2.0 for better security
	}
}

// SetUseV2 sets whether to use v2.0 API
func (c *Client) SetUseV2(useV2 bool) {
	c.useV2 = useV2
}

// buildPath constructs the full API path
func (c *Client) buildPath(path string) string {
	version := "v1.0"
	if c.useV2 {
		version = "v2.0"
	}
	return fmt.Sprintf("/cloud-trail/%s/appkeys/%s%s", version, c.appKey, path)
}

// doRequest performs an HTTP request
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
		if c.debug {
			fmt.Printf("[DEBUG] Request body: %s\n", string(jsonData))
		}
	}

	fullURL := c.baseURL + c.buildPath(path)
	if c.debug {
		fmt.Printf("[DEBUG] %s %s\n", method, fullURL)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// v2.0 requires user authentication headers
	if c.useV2 && c.accessKeyID != "" && c.secretKey != "" {
		req.Header.Set("X-TC-AUTHENTICATION-ID", c.accessKeyID)
		req.Header.Set("X-TC-AUTHENTICATION-SECRET", c.secretKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if c.debug {
		fmt.Printf("[DEBUG] Response status: %d\n", resp.StatusCode)
		fmt.Printf("[DEBUG] Response body: %s\n", string(respBody))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// SearchEvents searches CloudTrail events
func (c *Client) SearchEvents(ctx context.Context, input *SearchEventsInput) (*SearchEventsOutput, error) {
	data, err := c.doRequest(ctx, "POST", "/events/search", input)
	if err != nil {
		return nil, err
	}

	var result SearchEventsOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// SearchEventsSimple searches events with common defaults
func (c *Client) SearchEventsSimple(ctx context.Context, from, to time.Time, page, size int) (*SearchEventsOutput, error) {
	input := &SearchEventsInput{
		From: from,
		To:   to,
		Page: page,
		Size: size,
	}
	return c.SearchEvents(ctx, input)
}
