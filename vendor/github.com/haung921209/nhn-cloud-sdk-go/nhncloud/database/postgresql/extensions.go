package postgresql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// Extension represents a PostgreSQL extension
type Extension struct {
	ExtensionID     string   `json:"extensionId"`
	ExtensionName   string   `json:"extensionName"` // e.g., postgis, hstore, uuid-ossp
	ExtensionStatus string   `json:"extensionStatus"`
	Databases       []string `json:"databases"` // Databases where extension is installed
}

// ListExtensionsResponse is the response for ListExtensions
type ListExtensionsResponse struct {
	PostgreSQLResponse
	Extensions    []Extension `json:"extensions"`
	IsNeedToApply bool        `json:"isNeedToApply"` // Pending changes need apply
}

// ListExtensions retrieves all extensions for a PostgreSQL instance group.
//
// IMPORTANT: Extensions operate at INSTANCE GROUP level, not instance level.
//
// Known Issue (CSP-012): All Extensions APIs currently return 404 error.
// This API is implemented for completeness but is not functional.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#list-extensions
func (c *Client) ListExtensions(ctx context.Context, instanceGroupID string) (*ListExtensionsResponse, error) {
	if instanceGroupID == "" {
		return nil, &core.ValidationError{Field: "instanceGroupID", Message: "instance group ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instance-groups/%s/extensions", instanceGroupID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListExtensionsResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// InstallExtensionRequest is the request for installing an extension
type InstallExtensionRequest struct {
	DatabaseID  string `json:"databaseId"`
	SchemaName  string `json:"schemaName"`
	WithCascade *bool  `json:"withCascade,omitempty"` // Install dependencies
}

// Install ExtensionResponse is the response for InstallExtension
type InstallExtensionResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// InstallExtension installs a PostgreSQL extension in a database.
//
// IMPORTANT: Extensions operate at INSTANCE GROUP level.
//
// Known Issue (CSP-012): All Extensions APIs currently return 404 error.
// This API is implemented for completeness but is not functional.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#install-extension
func (c *Client) InstallExtension(ctx context.Context, instanceGroupID, extensionID string, req *InstallExtensionRequest) (*InstallExtensionResponse, error) {
	if instanceGroupID == "" {
		return nil, &core.ValidationError{Field: "instanceGroupID", Message: "instance group ID is required"}
	}
	if extensionID == "" {
		return nil, &core.ValidationError{Field: "extensionID", Message: "extension ID is required"}
	}
	if req.DatabaseID == "" {
		return nil, &core.ValidationError{Field: "DatabaseID", Message: "database ID is required"}
	}
	if req.SchemaName == "" {
		return nil, &core.ValidationError{Field: "SchemaName", Message: "schema name is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instance-groups/%s/extensions/%s", instanceGroupID, extensionID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result InstallExtensionResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteExtensionResponse is the response for DeleteExtension
type DeleteExtensionResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// DeleteExtension deletes an installed extension from a database.
//
// Known Issue (CSP-012): All Extensions APIs currently return 404 error.
// This API is implemented for completeness but is not functional.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#delete-extension
func (c *Client) DeleteExtension(ctx context.Context, instanceGroupID, extensionInstanceID string, withCascade bool) (*DeleteExtensionResponse, error) {
	if instanceGroupID == "" {
		return nil, &core.ValidationError{Field: "instanceGroupID", Message: "instance group ID is required"}
	}
	if extensionInstanceID == "" {
		return nil, &core.ValidationError{Field: "extensionInstanceID", Message: "extension instance ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instance-groups/%s/extensions/%s?withCascade=%t",
		instanceGroupID, extensionInstanceID, withCascade)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteExtensionResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ApplyExtensionsResponse is the response for ApplyExtensions
type ApplyExtensionsResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// ApplyExtensions applies pending extension changes.
//
// Known Issue (CSP-012): All Extensions APIs currently return 404 error.
// This API is implemented for completeness but is not functional.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#apply-extensions
func (c *Client) ApplyExtensions(ctx context.Context, instanceGroupID string) (*ApplyExtensionsResponse, error) {
	if instanceGroupID == "" {
		return nil, &core.ValidationError{Field: "instanceGroupID", Message: "instance group ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instance-groups/%s/extensions/apply", instanceGroupID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ApplyExtensionsResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SyncExtensionsResponse is the response for SyncExtensions
type SyncExtensionsResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// SyncExtensions synchronizes extension status.
//
// Known Issue (CSP-012): All Extensions APIs currently return 404 error.
// This API is implemented for completeness but is not functional.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#sync-extensions
func (c *Client) SyncExtensions(ctx context.Context, instanceGroupID string) (*SyncExtensionsResponse, error) {
	if instanceGroupID == "" {
		return nil, &core.ValidationError{Field: "instanceGroupID", Message: "instance group ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instance-groups/%s/extensions/sync", instanceGroupID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result SyncExtensionsResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
