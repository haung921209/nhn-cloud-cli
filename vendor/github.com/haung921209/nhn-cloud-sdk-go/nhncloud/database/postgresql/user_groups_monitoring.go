package postgresql

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
	PostgreSQLResponse
	UserGroups []UserGroup `json:"userGroups"`
}

// ListUserGroups retrieves all user groups.
func (c *Client) ListUserGroups(ctx context.Context) (*ListUserGroupsResponse, error) {
	path := "/v1.0/user-groups"
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
	PostgreSQLResponse
	UserGroup
}

// GetUserGroup retrieves a specific user group.
func (c *Client) GetUserGroup(ctx context.Context, groupID string) (*GetUserGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "group ID is required"}
	}

	path := fmt.Sprintf("/v1.0/user-groups/%s", groupID)
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
	PostgreSQLResponse
	UserGroupID string `json:"userGroupId"`
}

// CreateUserGroup creates a new user group.
func (c *Client) CreateUserGroup(ctx context.Context, req *CreateUserGroupRequest) (*CreateUserGroupResponse, error) {
	if req.UserGroupName == "" {
		return nil, &core.ValidationError{Field: "UserGroupName", Message: "user group name is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	path := "/v1.0/user-groups"
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

// DeleteUserGroupResponse is the response for DeleteUserGroup
type DeleteUserGroupResponse struct {
	PostgreSQLResponse
}

// DeleteUserGroup deletes a user group.
func (c *Client) DeleteUserGroup(ctx context.Context, groupID string) (*DeleteUserGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "group ID is required"}
	}

	path := fmt.Sprintf("/v1.0/user-groups/%s", groupID)
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
