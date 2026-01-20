// Package dnsplus provides DNS Plus service client
package dnsplus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultBaseURL = "https://dnsplus.api.nhncloudservice.com"

// Client represents a DNS Plus API client
type Client struct {
	baseURL    string
	appKey     string
	httpClient *http.Client
	debug      bool
}

// NewClient creates a new DNS Plus client
func NewClient(appKey string, httpClient *http.Client, debug bool) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{
		baseURL:    DefaultBaseURL,
		appKey:     appKey,
		httpClient: httpClient,
		debug:      debug,
	}
}

// buildPath constructs the full API path
func (c *Client) buildPath(path string) string {
	return fmt.Sprintf("/dnsplus/v1.0/appkeys/%s%s", c.appKey, path)
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

// ================================
// DNS Zone Operations
// ================================

// ListZones lists all DNS zones
func (c *Client) ListZones(ctx context.Context) (*ListZonesOutput, error) {
	data, err := c.doRequest(ctx, "GET", "/zones", nil)
	if err != nil {
		return nil, err
	}

	var result ListZonesOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateZone creates a new DNS zone
func (c *Client) CreateZone(ctx context.Context, input *CreateZoneInput) (*ZoneOutput, error) {
	request := map[string]interface{}{
		"zone": input,
	}
	data, err := c.doRequest(ctx, "POST", "/zones", request)
	if err != nil {
		return nil, err
	}

	var result ZoneOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateZone updates a DNS zone
func (c *Client) UpdateZone(ctx context.Context, zoneID string, input *UpdateZoneInput) (*ZoneOutput, error) {
	path := fmt.Sprintf("/zones/%s", zoneID)
	request := map[string]interface{}{
		"zone": input,
	}
	data, err := c.doRequest(ctx, "PUT", path, request)
	if err != nil {
		return nil, err
	}

	var result ZoneOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteZones deletes DNS zones (async)
func (c *Client) DeleteZones(ctx context.Context, zoneIDs []string) (*APIResponse, error) {
	request := map[string]interface{}{
		"zoneIdList": zoneIDs,
	}
	data, err := c.doRequest(ctx, "DELETE", "/zones/async", request)
	if err != nil {
		return nil, err
	}

	var result APIResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ================================
// Record Set Operations
// ================================

// ListRecordSets lists record sets in a zone
func (c *Client) ListRecordSets(ctx context.Context, zoneID string) (*ListRecordSetsOutput, error) {
	path := fmt.Sprintf("/zones/%s/recordsets", zoneID)
	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result ListRecordSetsOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateRecordSet creates a new record set
func (c *Client) CreateRecordSet(ctx context.Context, zoneID string, input *CreateRecordSetInput) (*RecordSetOutput, error) {
	path := fmt.Sprintf("/zones/%s/recordsets", zoneID)
	request := map[string]interface{}{
		"recordset": input,
	}
	data, err := c.doRequest(ctx, "POST", path, request)
	if err != nil {
		return nil, err
	}

	var result RecordSetOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateRecordSet updates a record set
func (c *Client) UpdateRecordSet(ctx context.Context, zoneID, recordsetID string, input *UpdateRecordSetInput) (*RecordSetOutput, error) {
	path := fmt.Sprintf("/zones/%s/recordsets/%s", zoneID, recordsetID)
	request := map[string]interface{}{
		"recordset": input,
	}
	data, err := c.doRequest(ctx, "PUT", path, request)
	if err != nil {
		return nil, err
	}

	var result RecordSetOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteRecordSets deletes record sets
func (c *Client) DeleteRecordSets(ctx context.Context, zoneID string, recordsetIDs []string) (*APIResponse, error) {
	path := fmt.Sprintf("/zones/%s/recordsets", zoneID)
	request := map[string]interface{}{
		"recordsetIdList": recordsetIDs,
	}
	data, err := c.doRequest(ctx, "DELETE", path, request)
	if err != nil {
		return nil, err
	}

	var result APIResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ================================
// GSLB Operations
// ================================

// ListGSLBs lists all GSLBs
func (c *Client) ListGSLBs(ctx context.Context) (*ListGSLBsOutput, error) {
	data, err := c.doRequest(ctx, "GET", "/gslbs", nil)
	if err != nil {
		return nil, err
	}

	var result ListGSLBsOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateGSLB creates a new GSLB
func (c *Client) CreateGSLB(ctx context.Context, input *CreateGSLBInput) (*GSLBOutput, error) {
	request := map[string]interface{}{
		"gslb": input,
	}
	data, err := c.doRequest(ctx, "POST", "/gslbs", request)
	if err != nil {
		return nil, err
	}

	var result GSLBOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateGSLB updates a GSLB
func (c *Client) UpdateGSLB(ctx context.Context, gslbID string, input *UpdateGSLBInput) (*GSLBOutput, error) {
	path := fmt.Sprintf("/gslbs/%s", gslbID)
	request := map[string]interface{}{
		"gslb": input,
	}
	data, err := c.doRequest(ctx, "PUT", path, request)
	if err != nil {
		return nil, err
	}

	var result GSLBOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteGSLBs deletes GSLBs
func (c *Client) DeleteGSLBs(ctx context.Context, gslbIDs []string) (*APIResponse, error) {
	request := map[string]interface{}{
		"gslbIdList": gslbIDs,
	}
	data, err := c.doRequest(ctx, "DELETE", "/gslbs", request)
	if err != nil {
		return nil, err
	}

	var result APIResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ================================
// GSLB Pool Operations
// ================================

// ListPools lists pools in a GSLB
func (c *Client) ListPools(ctx context.Context, gslbID string) (*ListPoolsOutput, error) {
	path := fmt.Sprintf("/gslbs/%s/pools", gslbID)
	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result ListPoolsOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreatePool creates a new pool
func (c *Client) CreatePool(ctx context.Context, gslbID string, input *CreatePoolInput) (*PoolOutput, error) {
	path := fmt.Sprintf("/gslbs/%s/pools", gslbID)
	request := map[string]interface{}{
		"pool": input,
	}
	data, err := c.doRequest(ctx, "POST", path, request)
	if err != nil {
		return nil, err
	}

	var result PoolOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdatePool updates a pool
func (c *Client) UpdatePool(ctx context.Context, gslbID, poolID string, input *UpdatePoolInput) (*PoolOutput, error) {
	path := fmt.Sprintf("/gslbs/%s/pools/%s", gslbID, poolID)
	request := map[string]interface{}{
		"pool": input,
	}
	data, err := c.doRequest(ctx, "PUT", path, request)
	if err != nil {
		return nil, err
	}

	var result PoolOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeletePools deletes pools
func (c *Client) DeletePools(ctx context.Context, gslbID string, poolIDs []string) (*APIResponse, error) {
	path := fmt.Sprintf("/gslbs/%s/pools", gslbID)
	request := map[string]interface{}{
		"poolIdList": poolIDs,
	}
	data, err := c.doRequest(ctx, "DELETE", path, request)
	if err != nil {
		return nil, err
	}

	var result APIResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ================================
// GSLB Endpoint Operations
// ================================

// ListEndpoints lists endpoints in a pool
func (c *Client) ListEndpoints(ctx context.Context, gslbID, poolID string) (*ListEndpointsOutput, error) {
	path := fmt.Sprintf("/gslbs/%s/pools/%s/endpoints", gslbID, poolID)
	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result ListEndpointsOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateEndpoint creates a new endpoint
func (c *Client) CreateEndpoint(ctx context.Context, gslbID, poolID string, input *CreateEndpointInput) (*EndpointOutput, error) {
	path := fmt.Sprintf("/gslbs/%s/pools/%s/endpoints", gslbID, poolID)
	request := map[string]interface{}{
		"endpoint": input,
	}
	data, err := c.doRequest(ctx, "POST", path, request)
	if err != nil {
		return nil, err
	}

	var result EndpointOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateEndpoint updates an endpoint
func (c *Client) UpdateEndpoint(ctx context.Context, gslbID, poolID, endpointID string, input *UpdateEndpointInput) (*EndpointOutput, error) {
	path := fmt.Sprintf("/gslbs/%s/pools/%s/endpoints/%s", gslbID, poolID, endpointID)
	request := map[string]interface{}{
		"endpoint": input,
	}
	data, err := c.doRequest(ctx, "PUT", path, request)
	if err != nil {
		return nil, err
	}

	var result EndpointOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteEndpoints deletes endpoints
func (c *Client) DeleteEndpoints(ctx context.Context, gslbID, poolID string, endpointIDs []string) (*APIResponse, error) {
	path := fmt.Sprintf("/gslbs/%s/pools/%s/endpoints", gslbID, poolID)
	request := map[string]interface{}{
		"endpointIdList": endpointIDs,
	}
	data, err := c.doRequest(ctx, "DELETE", path, request)
	if err != nil {
		return nil, err
	}

	var result APIResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ================================
// Health Check Operations
// ================================

// ListHealthChecks lists all health checks
func (c *Client) ListHealthChecks(ctx context.Context) (*ListHealthChecksOutput, error) {
	data, err := c.doRequest(ctx, "GET", "/health-checks", nil)
	if err != nil {
		return nil, err
	}

	var result ListHealthChecksOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateHealthCheck creates a new health check
func (c *Client) CreateHealthCheck(ctx context.Context, input *CreateHealthCheckInput) (*HealthCheckOutput, error) {
	request := map[string]interface{}{
		"healthCheck": input,
	}
	data, err := c.doRequest(ctx, "POST", "/health-checks", request)
	if err != nil {
		return nil, err
	}

	var result HealthCheckOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateHealthCheck updates a health check
func (c *Client) UpdateHealthCheck(ctx context.Context, healthCheckID string, input *UpdateHealthCheckInput) (*HealthCheckOutput, error) {
	path := fmt.Sprintf("/health-checks/%s", healthCheckID)
	request := map[string]interface{}{
		"healthCheck": input,
	}
	data, err := c.doRequest(ctx, "PUT", path, request)
	if err != nil {
		return nil, err
	}

	var result HealthCheckOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteHealthChecks deletes health checks
func (c *Client) DeleteHealthChecks(ctx context.Context, healthCheckIDs []string) (*APIResponse, error) {
	request := map[string]interface{}{
		"healthCheckIdList": healthCheckIDs,
	}
	data, err := c.doRequest(ctx, "DELETE", "/health-checks", request)
	if err != nil {
		return nil, err
	}

	var result APIResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
