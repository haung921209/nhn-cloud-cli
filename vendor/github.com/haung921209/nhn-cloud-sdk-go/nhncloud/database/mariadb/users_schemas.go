package mariadb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// DBUser represents a database user
type DBUser struct {
	DBUserID             string `json:"dbUserId,omitempty"`
	DBUserName           string `json:"dbUserName"`
	Host                 string `json:"host"`
	AuthorityType        string `json:"authorityType"`
	AuthenticationPlugin string `json:"authenticationPlugin,omitempty"`
	TLSOption            string `json:"tlsOption,omitempty"`
	CreatedAt            string `json:"createdAt,omitempty"`
}

// ListDBUsersResponse is the response for ListDBUsers
type ListDBUsersResponse struct {
	MariaDBResponse
	DBUsers []DBUser `json:"dbUsers"`
}

// ListDBUsers retrieves all database users for an instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db-user_1
func (c *Client) ListDBUsers(ctx context.Context, instanceID string) (*ListDBUsersResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/db-users", instanceID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListDBUsersResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateDBUserRequest is the request for creating a database user
type CreateDBUserRequest struct {
	DBUserName           string `json:"dbUserName"`
	DBPassword           string `json:"dbPassword"`
	Host                 string `json:"host"`
	AuthorityType        string `json:"authorityType"`
	AuthenticationPlugin string `json:"authenticationPlugin,omitempty"` // MariaDB: NATIVE, SHA256 (no CACHING_SHA2)
	// Note: TLSOption not supported in MariaDB
}

// CreateDBUserResponse is the response for CreateDBUser
type CreateDBUserResponse struct {
	MariaDBResponse
	JobID string `json:"jobId"`
}

// CreateDBUser creates a new database user.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db-user_2
func (c *Client) CreateDBUser(ctx context.Context, instanceID string, req *CreateDBUserRequest) (*CreateDBUserResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if req.DBUserName == "" {
		return nil, &core.ValidationError{Field: "DBUserName", Message: "username is required"}
	}
	if req.DBPassword == "" {
		return nil, &core.ValidationError{Field: "DBPassword", Message: "password is required"}
	}
	// API constraint: password 4-16 characters
	if len(req.DBPassword) < 4 || len(req.DBPassword) > 16 {
		return nil, &core.ValidationError{Field: "DBPassword", Message: "password must be 4-16 characters"}
	}
	if req.Host == "" {
		return nil, &core.ValidationError{Field: "Host", Message: "host is required"}
	}
	if req.AuthorityType == "" {
		return nil, &core.ValidationError{Field: "AuthorityType", Message: "authority type is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/db-users", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateDBUserResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateDBUserRequest is the request for updating a database user
type UpdateDBUserRequest struct {
	DBPassword           *string `json:"dbPassword,omitempty"`
	AuthorityType        *string `json:"authorityType,omitempty"`
	AuthenticationPlugin *string `json:"authenticationPlugin,omitempty"` // MariaDB: NATIVE, SHA256
	// Note: TLSOption not supported in MariaDB
}

// UpdateDBUserResponse is the response for UpdateDBUser
type UpdateDBUserResponse struct {
	MariaDBResponse
	JobID string `json:"jobId"`
}

// UpdateDBUser updates an existing database user.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db-user_3
func (c *Client) UpdateDBUser(ctx context.Context, instanceID, userID string, req *UpdateDBUserRequest) (*UpdateDBUserResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if userID == "" {
		return nil, &core.ValidationError{Field: "userID", Message: "user ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/db-users/%s", instanceID, userID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result UpdateDBUserResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteDBUserResponse is the response for DeleteDBUser
type DeleteDBUserResponse struct {
	MariaDBResponse
	JobID string `json:"jobId"`
}

// DeleteDBUser deletes a database user.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db-user_4
func (c *Client) DeleteDBUser(ctx context.Context, instanceID, userID string) (*DeleteDBUserResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if userID == "" {
		return nil, &core.ValidationError{Field: "userID", Message: "user ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/db-users/%s", instanceID, userID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteDBUserResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DBSchema represents a database schema
type DBSchema struct {
	DBSchemaID   string `json:"dbSchemaId,omitempty"`
	DBSchemaName string `json:"dbSchemaName"`
	CreatedAt    string `json:"createdAt,omitempty"`
}

// ListSchemasResponse is the response for ListSchemas
type ListSchemasResponse struct {
	MariaDBResponse
	DBSchemas []DBSchema `json:"dbSchemas"`
}

// ListSchemas retrieves all database schemas for an instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db-schema_1
func (c *Client) ListSchemas(ctx context.Context, instanceID string) (*ListSchemasResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/db-schemas", instanceID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListSchemasResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateSchemaRequest is the request for creating a schema
type CreateSchemaRequest struct {
	DBSchemaName string `json:"dbSchemaName"`
}

// CreateSchemaResponse is the response for CreateSchema
type CreateSchemaResponse struct {
	MariaDBResponse
	JobID string `json:"jobId"`
}

// CreateSchema creates a new database schema.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db-schema_2
func (c *Client) CreateSchema(ctx context.Context, instanceID string, req *CreateSchemaRequest) (*CreateSchemaResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if req.DBSchemaName == "" {
		return nil, &core.ValidationError{Field: "DBSchemaName", Message: "schema name is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/db-schemas", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateSchemaResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteSchemaResponse is the response for DeleteSchema
type DeleteSchemaResponse struct {
	MariaDBResponse
	JobID string `json:"jobId"`
}

// DeleteSchema deletes a database schema.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db-schema_3
func (c *Client) DeleteSchema(ctx context.Context, instanceID, schemaID string) (*DeleteSchemaResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if schemaID == "" {
		return nil, &core.ValidationError{Field: "schemaID", Message: "schema ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/db-schemas/%s", instanceID, schemaID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteSchemaResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
