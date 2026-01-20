package postgresql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// Database represents a PostgreSQL database
// PostgreSQL uses /databases API (not /db-schemas like MySQL)
type Database struct {
	DatabaseID   string `json:"databaseId"`
	DatabaseName string `json:"databaseName"`
	Owner        string `json:"owner"`
	Encoding     string `json:"encoding"`   // e.g., UTF8
	Collate      string `json:"collate"`    // e.g., en_US.UTF-8
	Ctype        string `json:"ctype"`      // e.g., en_US.UTF-8
	Size         int64  `json:"size"`       // Size in bytes
	TableCount   int    `json:"tableCount"` // Number of tables
}

// ListDatabasesResponse is the response for ListDatabases
type ListDatabasesResponse struct {
	PostgreSQLResponse
	Databases []Database `json:"databases"`
}

// ListDatabases retrieves all databases for a PostgreSQL instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#list-databases
func (c *Client) ListDatabases(ctx context.Context, instanceID string) (*ListDatabasesResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/databases", instanceID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListDatabasesResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateDatabaseRequest is the request for creating a database
type CreateDatabaseRequest struct {
	DatabaseName string `json:"databaseName"`
	Owner        string `json:"owner,omitempty"`    // Owner username
	Encoding     string `json:"encoding,omitempty"` // Default: UTF8
	Collate      string `json:"collate,omitempty"`  // Collation
	Ctype        string `json:"ctype,omitempty"`    // Character type
}

// CreateDatabaseResponse is the response for CreateDatabase
type CreateDatabaseResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// CreateDatabase creates a new database in a PostgreSQL instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#create-database
func (c *Client) CreateDatabase(ctx context.Context, instanceID string, req *CreateDatabaseRequest) (*CreateDatabaseResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if req.DatabaseName == "" {
		return nil, &core.ValidationError{Field: "DatabaseName", Message: "database name is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/databases", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateDatabaseResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ModifyDatabaseRequest is the request for modifying a database
type ModifyDatabaseRequest struct {
	Owner string `json:"owner,omitempty"`
}

// ModifyDatabaseResponse is the response for ModifyDatabase
type ModifyDatabaseResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// ModifyDatabase modifies a database's properties.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#modify-database
func (c *Client) ModifyDatabase(ctx context.Context, instanceID, databaseID string, req *ModifyDatabaseRequest) (*ModifyDatabaseResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if databaseID == "" {
		return nil, &core.ValidationError{Field: "databaseID", Message: "database ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/databases/%s", instanceID, databaseID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result ModifyDatabaseResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteDatabaseResponse is the response for DeleteDatabase
type DeleteDatabaseResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// DeleteDatabase deletes a database from a PostgreSQL instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#delete-database
func (c *Client) DeleteDatabase(ctx context.Context, instanceID, databaseID string) (*DeleteDatabaseResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if databaseID == "" {
		return nil, &core.ValidationError{Field: "databaseID", Message: "database ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/databases/%s", instanceID, databaseID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteDatabaseResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
