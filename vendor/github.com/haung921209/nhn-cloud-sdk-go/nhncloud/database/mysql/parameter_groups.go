package mysql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// ParameterGroup represents a database parameter group
type ParameterGroup struct {
	ParameterGroupID   string      `json:"parameterGroupId"`
	ParameterGroupName string      `json:"parameterGroupName"`
	Description        string      `json:"description,omitempty"`
	DBVersion          string      `json:"dbVersion"`
	Parameters         []Parameter `json:"parameters,omitempty"`
	CreatedAt          string      `json:"createdAt,omitempty"`
	UpdatedAt          string      `json:"updatedAt,omitempty"`
}

// Parameter represents a database parameter
type Parameter struct {
	ParameterID   string `json:"parameterId"`
	ParameterName string `json:"parameterName"`
	Value         string `json:"value"`
	DefaultValue  string `json:"defaultValue,omitempty"`
	AllowedValues string `json:"allowedValues,omitempty"`
	DataType      string `json:"dataType,omitempty"`
	IsModifiable  bool   `json:"isModifiable,omitempty"`
	ApplyType     string `json:"applyType,omitempty"`
}

// ListParameterGroupsResponse is the response for ListParameterGroups
type ListParameterGroupsResponse struct {
	MySQLResponse
	ParameterGroups []ParameterGroup `json:"parameterGroups"`
}

// ListParameterGroups retrieves all parameter groups.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#parameter-group_1
func (c *Client) ListParameterGroups(ctx context.Context) (*ListParameterGroupsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "/v3.0/parameter-groups", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListParameterGroupsResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetParameterGroupResponse is the response for GetParameterGroup
type GetParameterGroupResponse struct {
	MySQLResponse
	ParameterGroup ParameterGroup `json:"parameterGroup"`
}

// GetParameterGroup retrieves a specific parameter group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#parameter-group_2
func (c *Client) GetParameterGroup(ctx context.Context, groupID string) (*GetParameterGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "parameter group ID is required"}
	}

	path := fmt.Sprintf("/v3.0/parameter-groups/%s", groupID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result GetParameterGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateParameterGroupRequest is the request for creating a parameter group
type CreateParameterGroupRequest struct {
	ParameterGroupName string `json:"parameterGroupName"`
	Description        string `json:"description,omitempty"`
	DBVersion          string `json:"dbVersion"`
}

// CreateParameterGroupResponse is the response for CreateParameterGroup
type CreateParameterGroupResponse struct {
	MySQLResponse
	ParameterGroupID string `json:"parameterGroupId"`
}

// CreateParameterGroup creates a new parameter group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#parameter-group_3
func (c *Client) CreateParameterGroup(ctx context.Context, req *CreateParameterGroupRequest) (*CreateParameterGroupResponse, error) {
	if req.ParameterGroupName == "" {
		return nil, &core.ValidationError{Field: "ParameterGroupName", Message: "parameter group name is required"}
	}
	if req.DBVersion == "" {
		return nil, &core.ValidationError{Field: "DBVersion", Message: "DB version is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "/v3.0/parameter-groups", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateParameterGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CopyParameterGroupRequest is the request for copying a parameter group
type CopyParameterGroupRequest struct {
	ParameterGroupName string `json:"parameterGroupName"`
	Description        string `json:"description,omitempty"`
}

// CopyParameterGroupResponse is the response for CopyParameterGroup
type CopyParameterGroupResponse struct {
	MySQLResponse
	ParameterGroupID string `json:"parameterGroupId"`
}

// CopyParameterGroup copies an existing parameter group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#parameter-group_4
func (c *Client) CopyParameterGroup(ctx context.Context, groupID string, req *CopyParameterGroupRequest) (*CopyParameterGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "parameter group ID is required"}
	}
	if req.ParameterGroupName == "" {
		return nil, &core.ValidationError{Field: "ParameterGroupName", Message: "parameter group name is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v3.0/parameter-groups/%s/copy", groupID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CopyParameterGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateParameterGroupRequest is the request for updating a parameter group
type UpdateParameterGroupRequest struct {
	ParameterGroupName *string `json:"parameterGroupName,omitempty"`
	Description        *string `json:"description,omitempty"`
}

// UpdateParameterGroupResponse is the response for UpdateParameterGroup
type UpdateParameterGroupResponse struct {
	MySQLResponse
}

// UpdateParameterGroup updates a parameter group's metadata.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#parameter-group_5
func (c *Client) UpdateParameterGroup(ctx context.Context, groupID string, req *UpdateParameterGroupRequest) (*UpdateParameterGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "parameter group ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v3.0/parameter-groups/%s", groupID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result UpdateParameterGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ModifiedParameter represents a parameter modification
type ModifiedParameter struct {
	ParameterID string `json:"parameterId"`
	Value       string `json:"value"`
}

// ModifyParametersRequest is the request for modifying parameters
type ModifyParametersRequest struct {
	ModifiedParameters []ModifiedParameter `json:"modifiedParameters"`
}

// ModifyParametersResponse is the response for ModifyParameters
type ModifyParametersResponse struct {
	MySQLResponse
}

// ModifyParameters modifies parameters within a parameter group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#parameter-group_6
func (c *Client) ModifyParameters(ctx context.Context, groupID string, req *ModifyParametersRequest) (*ModifyParametersResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "parameter group ID is required"}
	}
	if len(req.ModifiedParameters) == 0 {
		return nil, &core.ValidationError{Field: "ModifiedParameters", Message: "at least one parameter modification is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v3.0/parameter-groups/%s/parameters", groupID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result ModifyParametersResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ResetParameterGroupResponse is the response for ResetParameterGroup
type ResetParameterGroupResponse struct {
	MySQLResponse
}

// ResetParameterGroup resets a parameter group to default values.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#parameter-group_7
func (c *Client) ResetParameterGroup(ctx context.Context, groupID string) (*ResetParameterGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "parameter group ID is required"}
	}

	path := fmt.Sprintf("/v3.0/parameter-groups/%s/reset", groupID)
	req, err := http.NewRequestWithContext(ctx, "PUT", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ResetParameterGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteParameterGroupResponse is the response for DeleteParameterGroup
type DeleteParameterGroupResponse struct {
	MySQLResponse
}

// DeleteParameterGroup deletes a parameter group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#parameter-group_8
func (c *Client) DeleteParameterGroup(ctx context.Context, groupID string) (*DeleteParameterGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "parameter group ID is required"}
	}

	path := fmt.Sprintf("/v3.0/parameter-groups/%s", groupID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteParameterGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
