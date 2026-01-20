package postgresql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// NetworkInfo represents network information for an instance
type NetworkInfo struct {
	AvailabilityZone string            `json:"availabilityZone"`
	Subnet           NetworkSubnet     `json:"subnet"`
	EndPoints        []NetworkEndPoint `json:"endPoints"`
}

// NetworkSubnet represents subnet information
type NetworkSubnet struct {
	SubnetID   string `json:"subnetId"`
	SubnetName string `json:"subnetName,omitempty"`
	SubnetCIDR string `json:"subnetCidr,omitempty"`
}

// NetworkEndPoint represents a connection endpoint
type NetworkEndPoint struct {
	Domain       string `json:"domain"`
	IPAddress    string `json:"ipAddress"`
	EndPointType string `json:"endPointType"` // MASTER, READ
}

// GetNetworkInfoResponse is the response for GetNetworkInfo
type GetNetworkInfoResponse struct {
	PostgreSQLResponse
	NetworkInfo
}

// GetNetworkInfo retrieves network information for an instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_65
func (c *Client) GetNetworkInfo(ctx context.Context, instanceID string) (*GetNetworkInfoResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/network-info", instanceID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result GetNetworkInfoResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ModifyNetworkInfoRequest is the request for modifying network configuration
type ModifyNetworkInfoRequest struct {
	UsePublicAccess bool `json:"usePublicAccess"`
}

// ModifyNetworkInfoResponse is the response for ModifyNetworkInfo
type ModifyNetworkInfoResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// ModifyNetworkInfo modifies the network configuration (public access).
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_66
func (c *Client) ModifyNetworkInfo(ctx context.Context, instanceID string, req *ModifyNetworkInfoRequest) (*ModifyNetworkInfoResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/network-info", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result ModifyNetworkInfoResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// StorageInfo represents storage information
type StorageInfo struct {
	StorageType string `json:"storageType"`
	StorageSize int    `json:"storageSize"`
}

// GetStorageInfoResponse is the response for GetStorageInfo
type GetStorageInfoResponse struct {
	PostgreSQLResponse
	StorageInfo
}

// GetStorageInfo retrieves storage information for an instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_66
func (c *Client) GetStorageInfo(ctx context.Context, instanceID string) (*GetStorageInfoResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/storage-info", instanceID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result GetStorageInfoResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ModifyStorageInfoRequest is the request for modifying storage size
type ModifyStorageInfoRequest struct {
	StorageSize       int   `json:"storageSize"`
	UseOnlineFailover *bool `json:"useOnlineFailover,omitempty"`
}

// ModifyStorageInfoResponse is the response for ModifyStorageInfo
type ModifyStorageInfoResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// ModifyStorageInfo modifies the storage size of an instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_67
func (c *Client) ModifyStorageInfo(ctx context.Context, instanceID string, req *ModifyStorageInfoRequest) (*ModifyStorageInfoResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if req.StorageSize < 20 || req.StorageSize > 2048 {
		return nil, &core.ValidationError{Field: "StorageSize", Message: "storage size must be 20-2048 GB"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/storage-info", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result ModifyStorageInfoResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ModifyDeletionProtectionRequest is the request for modifying deletion protection
type ModifyDeletionProtectionRequest struct {
	UseDeletionProtection bool `json:"useDeletionProtection"`
}

// ModifyDeletionProtectionResponse is the response for ModifyDeletionProtection
type ModifyDeletionProtectionResponse struct {
	PostgreSQLResponse
}

// ModifyDeletionProtection enables or disables deletion protection.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_68
func (c *Client) ModifyDeletionProtection(ctx context.Context, instanceID string, req *ModifyDeletionProtectionRequest) (*ModifyDeletionProtectionResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/deletion-protection", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result ModifyDeletionProtectionResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
