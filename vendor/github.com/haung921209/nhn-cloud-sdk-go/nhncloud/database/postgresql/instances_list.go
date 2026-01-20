package postgresql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// ListInstancesResponse is the response for ListInstances
type ListInstancesResponse struct {
	PostgreSQLResponse
	DBInstances []DatabaseInstance `json:"dbInstances"`
}

// ListInstances retrieves a list of PostgreSQL instances.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#_1
func (c *Client) ListInstances(ctx context.Context) (*ListInstancesResponse, error) {
	path := "/v1.0/db-instances"
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListInstancesResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetInstanceResponse is the response for GetInstance
type GetInstanceResponse struct {
	PostgreSQLResponse
	DatabaseInstance
}

// GetInstance retrieves details of a specific PostgreSQL instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#_2
func (c *Client) GetInstance(ctx context.Context, instanceID string) (*GetInstanceResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s", instanceID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result GetInstanceResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
