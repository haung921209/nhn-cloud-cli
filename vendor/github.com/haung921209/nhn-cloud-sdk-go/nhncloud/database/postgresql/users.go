package postgresql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// DBUser represents a PostgreSQL database user
type DBUser struct {
	DBUserID      string `json:"dbUserId"`
	DBUserName    string `json:"dbUserName"`
	AuthorityType string `json:"authorityType,omitempty"`
	// PostgreSQL-specific fields
	IsSuperuser     bool `json:"isSuperuser,omitempty"`
	CanCreateDB     bool `json:"canCreateDB,omitempty"`
	CanCreateRole   bool `json:"canCreateRole,omitempty"`
	CanLogin        bool `json:"canLogin,omitempty"`
	IsReplication   bool `json:"isReplication,omitempty"`
	ConnectionLimit int  `json:"connectionLimit,omitempty"`
}

// ListDBUsersResponse is the response for ListDBUsers
type ListDBUsersResponse struct {
	PostgreSQLResponse
	DBUsers []DBUser `json:"dbUsers"`
}

// ListDBUsers retrieves all database users for a PostgreSQL instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#list-users
func (c *Client) ListDBUsers(ctx context.Context, instanceID string) (*ListDBUsersResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/db-users", instanceID)
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

// CreateDBUserRequest is the request for creating a database user with PostgreSQL-specific options
type CreateDBUserRequest struct {
	DBUserName    string `json:"dbUserName"`
	DBPassword    string `json:"dbPassword"` // 4-16 chars
	AuthorityType string `json:"authorityType,omitempty"`
	// PostgreSQL-specific extended options
	IsSuperuser     *bool `json:"isSuperuser,omitempty"`
	CanCreateDB     *bool `json:"canCreateDB,omitempty"`
	CanCreateRole   *bool `json:"canCreateRole,omitempty"`
	CanLogin        *bool `json:"canLogin,omitempty"`
	IsReplication   *bool `json:"isReplication,omitempty"`
	ConnectionLimit *int  `json:"connectionLimit,omitempty"`
}

// CreateDBUserResponse is the response for CreateDBUser
type CreateDBUserResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// CreateDBUser creates a new database user with PostgreSQL-specific privileges.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#create-user
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
	if len(req.DBPassword) < 4 || len(req.DBPassword) > 16 {
		return nil, &core.ValidationError{
			Field:   "DBPassword",
			Message: "password must be 4-16 characters",
		}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/db-users", instanceID)
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
	DBPassword    *string `json:"dbPassword,omitempty"`
	AuthorityType *string `json:"authorityType,omitempty"`
	// PostgreSQL-specific options
	IsSuperuser     *bool `json:"isSuperuser,omitempty"`
	CanCreateDB     *bool `json:"canCreateDB,omitempty"`
	CanCreateRole   *bool `json:"canCreateRole,omitempty"`
	CanLogin        *bool `json:"canLogin,omitempty"`
	IsReplication   *bool `json:"isReplication,omitempty"`
	ConnectionLimit *int  `json:"connectionLimit,omitempty"`
}

// UpdateDBUserResponse is the response for UpdateDBUser
type UpdateDBUserResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// UpdateDBUser updates a database user's properties.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#update-user
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

	path := fmt.Sprintf("/v1.0/db-instances/%s/db-users/%s", instanceID, userID)
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
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// DeleteDBUser deletes a database user.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#delete-user
func (c *Client) DeleteDBUser(ctx context.Context, instanceID, userID string) (*DeleteDBUserResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if userID == "" {
		return nil, &core.ValidationError{Field: "userID", Message: "user ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/db-users/%s", instanceID, userID)
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
