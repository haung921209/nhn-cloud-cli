// Package certmanager provides Certificate Manager service client
package certmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultBaseURL = "https://certmanager.api.nhncloudservice.com"

// Client represents a Certificate Manager API client
type Client struct {
	baseURL     string
	appKey      string
	accessKeyID string
	secretKey   string
	httpClient  *http.Client
	debug       bool
}

// NewClient creates a new Certificate Manager client
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
	}
}

// buildPath constructs the full API path
func (c *Client) buildPath(path string) string {
	return fmt.Sprintf("/certmanager/v1.0/appkeys/%s%s", c.appKey, path)
}

// doRequest performs an HTTP request
func (c *Client) doRequest(ctx context.Context, method, path string) ([]byte, error) {
	fullURL := c.baseURL + c.buildPath(path)
	if c.debug {
		fmt.Printf("[DEBUG] %s %s\n", method, fullURL)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add user authentication headers
	if c.accessKeyID != "" && c.secretKey != "" {
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

// ListCertificates lists all certificates
func (c *Client) ListCertificates(ctx context.Context) (*ListCertificatesOutput, error) {
	data, err := c.doRequest(ctx, "GET", "/certificates")
	if err != nil {
		return nil, err
	}

	var result ListCertificatesOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DownloadCertificateFiles downloads certificate files
func (c *Client) DownloadCertificateFiles(ctx context.Context, certificateName string) (*DownloadCertificateFilesOutput, error) {
	path := fmt.Sprintf("/certificates/%s/files", certificateName)
	data, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var result DownloadCertificateFilesOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
