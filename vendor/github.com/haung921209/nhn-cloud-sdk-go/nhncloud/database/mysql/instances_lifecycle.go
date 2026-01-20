package mysql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// StartInstanceResponse is the response for StartInstance
type StartInstanceResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// StartInstance starts a stopped MySQL instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#db_6
func (c *Client) StartInstance(ctx context.Context, instanceID string) (*StartInstanceResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/start", instanceID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result StartInstanceResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// StopInstanceResponse is the response for StopInstance
type StopInstanceResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// StopInstance stops a running MySQL instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#db_7
func (c *Client) StopInstance(ctx context.Context, instanceID string) (*StopInstanceResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/stop", instanceID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result StopInstanceResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RestartInstanceRequest is the request for restarting an instance
type RestartInstanceRequest struct {
	UseOnlineFailover *bool `json:"useOnlineFailover,omitempty"`
	ExecuteBackup     *bool `json:"executeBackup,omitempty"`
}

// RestartInstanceResponse is the response for RestartInstance
type RestartInstanceResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// RestartInstance restarts a MySQL instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#db_8
func (c *Client) RestartInstance(ctx context.Context, instanceID string, req *RestartInstanceRequest) (*RestartInstanceResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	var body []byte
	var err error
	if req != nil {
		body, err = json.Marshal(req)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/restart", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result RestartInstanceResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ForceRestartInstanceResponse is the response for ForceRestartInstance
type ForceRestartInstanceResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// ForceRestartInstance force restarts a MySQL instance.
// Use this when normal restart fails.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#db_9
func (c *Client) ForceRestartInstance(ctx context.Context, instanceID string) (*ForceRestartInstanceResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/force-restart", instanceID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ForceRestartInstanceResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
