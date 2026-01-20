// Package resourcewatcher provides Resource Watcher (Governance) service client
package resourcewatcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultBaseURL = "https://resource-watcher.api.nhncloudservice.com"

// Client represents a Resource Watcher API client
type Client struct {
	baseURL     string
	appKey      string
	accessKeyID string
	secretKey   string
	httpClient  *http.Client
	debug       bool
}

// NewClient creates a new Resource Watcher client
// appKey: The app key for the Resource Watcher service
// accessKeyID: User Access Key ID for authentication
// secretKey: Secret Access Key for authentication
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
func (c *Client) buildPath(version, resource string) string {
	return fmt.Sprintf("/resource-watcher/%s/appkeys/%s/%s", version, c.appKey, resource)
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

	fullURL := c.baseURL + path
	if c.debug {
		fmt.Printf("[DEBUG] %s %s\n", method, fullURL)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TC-AUTHENTICATION-ID", c.accessKeyID)
	req.Header.Set("X-TC-AUTHENTICATION-SECRET", c.secretKey)

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

// ============== Event Alarm APIs (v2.0) ==============

// CreateEventAlarm creates a new event alarm
func (c *Client) CreateEventAlarm(ctx context.Context, input *CreateEventAlarmInput) (*CreateEventAlarmOutput, error) {
	path := c.buildPath("v2.0", "event-alarms")

	data, err := c.doRequest(ctx, "POST", path, input)
	if err != nil {
		return nil, fmt.Errorf("create event alarm: %w", err)
	}

	var result CreateEventAlarmOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// GetEventAlarm retrieves details of a specific event alarm
func (c *Client) GetEventAlarm(ctx context.Context, alarmID string) (*GetEventAlarmOutput, error) {
	path := c.buildPath("v2.0", fmt.Sprintf("event-alarms/%s", alarmID))

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get event alarm: %w", err)
	}

	var result GetEventAlarmOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// SearchEventAlarms searches for event alarms
func (c *Client) SearchEventAlarms(ctx context.Context, input *SearchEventAlarmsInput) (*SearchEventAlarmsOutput, error) {
	path := c.buildPath("v2.0", "event-alarms/search")

	data, err := c.doRequest(ctx, "POST", path, input)
	if err != nil {
		return nil, fmt.Errorf("search event alarms: %w", err)
	}

	var result SearchEventAlarmsOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateEventAlarm updates an existing event alarm
func (c *Client) UpdateEventAlarm(ctx context.Context, alarmID string, input *UpdateEventAlarmInput) (*SimpleOutput, error) {
	path := c.buildPath("v2.0", fmt.Sprintf("event-alarms/%s", alarmID))

	data, err := c.doRequest(ctx, "PUT", path, input)
	if err != nil {
		return nil, fmt.Errorf("update event alarm: %w", err)
	}

	var result SimpleOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteEventAlarm deletes an event alarm
func (c *Client) DeleteEventAlarm(ctx context.Context, alarmID string) (*SimpleOutput, error) {
	path := c.buildPath("v2.0", fmt.Sprintf("event-alarms/%s", alarmID))

	data, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, fmt.Errorf("delete event alarm: %w", err)
	}

	var result SimpleOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteEventAlarms deletes multiple event alarms
func (c *Client) DeleteEventAlarms(ctx context.Context, input *DeleteEventAlarmsInput) (*SimpleOutput, error) {
	path := c.buildPath("v2.0", "event-alarms")

	data, err := c.doRequest(ctx, "DELETE", path, input)
	if err != nil {
		return nil, fmt.Errorf("delete event alarms: %w", err)
	}

	var result SimpleOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// ============== Alarm History APIs (v1.0) ==============

// GetAlarmHistory retrieves a specific alarm history record
func (c *Client) GetAlarmHistory(ctx context.Context, alarmID, historyID string) (*GetAlarmHistoryOutput, error) {
	path := c.buildPath("v1.0", fmt.Sprintf("alarms/%s/alarm-history/%s", alarmID, historyID))

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get alarm history: %w", err)
	}

	var result GetAlarmHistoryOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// SearchAlarmHistory searches for alarm history
func (c *Client) SearchAlarmHistory(ctx context.Context, alarmID string, input *SearchAlarmHistoryInput) (*SearchAlarmHistoryOutput, error) {
	path := c.buildPath("v1.0", fmt.Sprintf("alarms/%s/alarm-history/search", alarmID))

	data, err := c.doRequest(ctx, "POST", path, input)
	if err != nil {
		return nil, fmt.Errorf("search alarm history: %w", err)
	}

	var result SearchAlarmHistoryOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// ============== Event APIs (v1.0) ==============

// ListEvents retrieves a list of events
func (c *Client) ListEvents(ctx context.Context) (*ListEventsOutput, error) {
	path := c.buildPath("v1.0", "events")

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	var result ListEventsOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// GetEvent retrieves details of a specific event
func (c *Client) GetEvent(ctx context.Context, productID, eventID string) (*GetEventOutput, error) {
	path := c.buildPath("v1.0", fmt.Sprintf("events/%s/%s", productID, eventID))

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}

	var result GetEventOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// ============== Resource Group APIs (v1.0) ==============

// ListResourceGroups retrieves a list of resource groups
func (c *Client) ListResourceGroups(ctx context.Context) (*ListResourceGroupsOutput, error) {
	path := c.buildPath("v1.0", "resource-groups")

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list resource groups: %w", err)
	}

	var result ListResourceGroupsOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// ============== Resource Tag APIs (v1.0) ==============

// ListResourceTags retrieves a list of resource tags
func (c *Client) ListResourceTags(ctx context.Context) (*ListResourceTagsOutput, error) {
	path := c.buildPath("v1.0", "resource-tags")

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list resource tags: %w", err)
	}

	var result ListResourceTagsOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}
