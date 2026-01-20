package mysql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// UserGroup represents a user group
type UserGroup struct {
	UserGroupID       string   `json:"userGroupId"`
	UserGroupName     string   `json:"userGroupName"`
	UserGroupTypeCode string   `json:"userGroupTypeCode,omitempty"`
	Members           []Member `json:"members,omitempty"`
	CreatedYmdt       string   `json:"createdYmdt,omitempty"`
	UpdatedYmdt       string   `json:"updatedYmdt,omitempty"`
}

// Member represents a user group member
type Member struct {
	MemberID string `json:"memberId"`
}

// ListUserGroupsResponse is the response for ListUserGroups
type ListUserGroupsResponse struct {
	MySQLResponse
	UserGroups []UserGroup `json:"userGroups"`
}

// ListUserGroups retrieves all user groups.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_69
func (c *Client) ListUserGroups(ctx context.Context) (*ListUserGroupsResponse, error) {
	path := "/v3.0/user-groups"
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListUserGroupsResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetUserGroupResponse is the response for GetUserGroup
type GetUserGroupResponse struct {
	MySQLResponse
	UserGroup
}

// GetUserGroup retrieves a specific user group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_70
func (c *Client) GetUserGroup(ctx context.Context, groupID string) (*GetUserGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "group ID is required"}
	}

	path := fmt.Sprintf("/v3.0/user-groups/%s", groupID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result GetUserGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateUserGroupRequest is the request for creating a user group
type CreateUserGroupRequest struct {
	UserGroupName string   `json:"userGroupName"`
	MemberIDs     []string `json:"memberIds,omitempty"`
	SelectAllYN   bool     `json:"selectAllYN,omitempty"`
}

// CreateUserGroupResponse is the response for CreateUserGroup
type CreateUserGroupResponse struct {
	MySQLResponse
	UserGroupID string `json:"userGroupId"`
}

// CreateUserGroup creates a new user group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_71
func (c *Client) CreateUserGroup(ctx context.Context, req *CreateUserGroupRequest) (*CreateUserGroupResponse, error) {
	if req.UserGroupName == "" {
		return nil, &core.ValidationError{Field: "UserGroupName", Message: "user group name is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	path := "/v3.0/user-groups"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateUserGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateUserGroupRequest is the request for updating a user group
type UpdateUserGroupRequest struct {
	UserGroupName string   `json:"userGroupName,omitempty"`
	MemberIDs     []string `json:"memberIds,omitempty"`
}

// UpdateUserGroupResponse is the response for UpdateUserGroup
type UpdateUserGroupResponse struct {
	MySQLResponse
}

// UpdateUserGroup updates a user group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_72
func (c *Client) UpdateUserGroup(ctx context.Context, groupID string, req *UpdateUserGroupRequest) (*UpdateUserGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "group ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v3.0/user-groups/%s", groupID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result UpdateUserGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteUserGroupResponse is the response for DeleteUserGroup
type DeleteUserGroupResponse struct {
	MySQLResponse
}

// DeleteUserGroup deletes a user group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/#_73
func (c *Client) DeleteUserGroup(ctx context.Context, groupID string) (*DeleteUserGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "group ID is required"}
	}

	path := fmt.Sprintf("/v3.0/user-groups/%s", groupID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteUserGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
