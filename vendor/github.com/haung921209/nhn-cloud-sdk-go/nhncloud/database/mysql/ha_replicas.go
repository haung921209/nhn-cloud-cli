package mysql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// EnableHARequest is the request for enabling high availability
type EnableHARequest struct {
	UseHighAvailability bool   `json:"useHighAvailability"`
	PingInterval        *int   `json:"pingInterval,omitempty"`
	ReplicationMode     string `json:"replicationMode,omitempty"`
}

// EnableHAResponse is the response for EnableHA
type EnableHAResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// EnableHA enables high availability for an instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_58
func (c *Client) EnableHA(ctx context.Context, instanceID string, req *EnableHARequest) (*EnableHAResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/high-availability", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result EnableHAResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DisableHARequest is the request for disabling high availability
type DisableHARequest struct {
	UseHighAvailability bool `json:"useHighAvailability"` // Should be false
}

// DisableHAResponse is the response for DisableHA
type DisableHAResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// DisableHA disables high availability for an instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_58
func (c *Client) DisableHA(ctx context.Context, instanceID string) (*DisableHAResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	req := DisableHARequest{UseHighAvailability: false}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/high-availability", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result DisableHAResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// PauseHAResponse is the response for PauseHA
type PauseHAResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// PauseHA pauses high availability monitoring.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_59
func (c *Client) PauseHA(ctx context.Context, instanceID string) (*PauseHAResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/high-availability/pause", instanceID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result PauseHAResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ResumeHAResponse is the response for ResumeHA
type ResumeHAResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// ResumeHA resumes high availability monitoring.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_60
func (c *Client) ResumeHA(ctx context.Context, instanceID string) (*ResumeHAResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/high-availability/resume", instanceID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ResumeHAResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RepairHAResponse is the response for RepairHA
type RepairHAResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// RepairHA repairs high availability configuration.
// Note: May return 500 error on healthy instances (known issue CSP-011).
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_61
func (c *Client) RepairHA(ctx context.Context, instanceID string) (*RepairHAResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/high-availability/repair", instanceID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result RepairHAResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SplitHAResponse is the response for SplitHA
type SplitHAResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// SplitHA splits a high availability setup into separate instances.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_62
func (c *Client) SplitHA(ctx context.Context, instanceID string) (*SplitHAResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/high-availability/split", instanceID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result SplitHAResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateReplicaRequest is the request for creating a read replica
type CreateReplicaRequest struct {
	DBInstanceName string `json:"dbInstanceName"`
	// Additional replica configuration fields
}

// CreateReplicaResponse is the response for CreateReplica
type CreateReplicaResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// CreateReplica creates a read replica from an instance.
// Note: May return 500 error (known issue CSP-009).
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_63
func (c *Client) CreateReplica(ctx context.Context, instanceID string, req *CreateReplicaRequest) (*CreateReplicaResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/replicate", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateReplicaResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// PromoteReplicaResponse is the response for PromoteReplica
type PromoteReplicaResponse struct {
	MySQLResponse
	JobID string `json:"jobId"`
}

// PromoteReplica promotes a read replica to a standalone instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_64
func (c *Client) PromoteReplica(ctx context.Context, instanceID string) (*PromoteReplicaResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/promote", instanceID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result PromoteReplicaResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
