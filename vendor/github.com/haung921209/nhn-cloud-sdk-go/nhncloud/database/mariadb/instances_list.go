package mariadb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// ListInstancesResponse is the response for ListInstances
type ListInstancesResponse struct {
	MariaDBResponse
	DBInstances []DatabaseInstance `json:"dbInstances"`
}

// ListInstances retrieves a list of MariaDB instances.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db_1
func (c *Client) ListInstances(ctx context.Context) (*ListInstancesResponse, error) {
	path := "/v3.0/db-instances"
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
	MariaDBResponse
	DatabaseInstance
}

// GetInstance retrieves details of a specific MariaDB instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db_2
func (c *Client) GetInstance(ctx context.Context, instanceID string) (*GetInstanceResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s", instanceID)
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
