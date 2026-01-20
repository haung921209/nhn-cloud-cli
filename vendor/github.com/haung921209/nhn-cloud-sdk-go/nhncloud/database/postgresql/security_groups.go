package postgresql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// SecurityGroup represents a database security group
type SecurityGroup struct {
	DBSecurityGroupID   string         `json:"dbSecurityGroupId"`
	DBSecurityGroupName string         `json:"dbSecurityGroupName"`
	Description         string         `json:"description,omitempty"`
	Rules               []SecurityRule `json:"rules,omitempty"`
	CreatedAt           string         `json:"createdAt,omitempty"`
	UpdatedAt           string         `json:"updatedAt,omitempty"`
}

// SecurityRule represents a security group rule
type SecurityRule struct {
	RuleID      string   `json:"ruleId,omitempty"`
	Description string   `json:"description,omitempty"`
	Direction   string   `json:"direction"`
	EtherType   string   `json:"etherType"`
	Port        RulePort `json:"port"`
	CIDR        string   `json:"cidr"`
}

// RulePort represents port configuration for a security rule
type RulePort struct {
	PortType string `json:"portType"`
	MinPort  *int   `json:"minPort,omitempty"`
	MaxPort  *int   `json:"maxPort,omitempty"`
}

// ListSecurityGroupsResponse is the response for ListSecurityGroups
type ListSecurityGroupsResponse struct {
	PostgreSQLResponse
	DBSecurityGroups []SecurityGroup `json:"dbSecurityGroups"`
}

// ListSecurityGroups retrieves all database security groups.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#db-security-group_1
func (c *Client) ListSecurityGroups(ctx context.Context) (*ListSecurityGroupsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "/v1.0/db-security-groups", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListSecurityGroupsResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSecurityGroupResponse is the response for GetSecurityGroup
type GetSecurityGroupResponse struct {
	PostgreSQLResponse
	DBSecurityGroup SecurityGroup `json:"dbSecurityGroup"`
}

// GetSecurityGroup retrieves a specific database security group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#db-security-group_2
func (c *Client) GetSecurityGroup(ctx context.Context, groupID string) (*GetSecurityGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "security group ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-security-groups/%s", groupID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result GetSecurityGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateSecurityGroupRequest is the request for creating a security group
type CreateSecurityGroupRequest struct {
	DBSecurityGroupName string         `json:"dbSecurityGroupName"`
	Description         string         `json:"description,omitempty"`
	Rules               []SecurityRule `json:"rules,omitempty"`
}

// CreateSecurityGroupResponse is the response for CreateSecurityGroup
type CreateSecurityGroupResponse struct {
	PostgreSQLResponse
	DBSecurityGroupID string `json:"dbSecurityGroupId"`
}

// CreateSecurityGroup creates a new database security group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#db-security-group_3
func (c *Client) CreateSecurityGroup(ctx context.Context, req *CreateSecurityGroupRequest) (*CreateSecurityGroupResponse, error) {
	if req.DBSecurityGroupName == "" {
		return nil, &core.ValidationError{Field: "DBSecurityGroupName", Message: "security group name is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "/v1.0/db-security-groups", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateSecurityGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateSecurityGroupRequest is the request for updating a security group
type UpdateSecurityGroupRequest struct {
	DBSecurityGroupName *string `json:"dbSecurityGroupName,omitempty"`
	Description         *string `json:"description,omitempty"`
}

// UpdateSecurityGroupResponse is the response for UpdateSecurityGroup
type UpdateSecurityGroupResponse struct {
	PostgreSQLResponse
}

// UpdateSecurityGroup updates a database security group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#db-security-group_4
func (c *Client) UpdateSecurityGroup(ctx context.Context, groupID string, req *UpdateSecurityGroupRequest) (*UpdateSecurityGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "security group ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-security-groups/%s", groupID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result UpdateSecurityGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteSecurityGroupResponse is the response for DeleteSecurityGroup
type DeleteSecurityGroupResponse struct {
	PostgreSQLResponse
}

// DeleteSecurityGroup deletes a database security group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#db-security-group_5
func (c *Client) DeleteSecurityGroup(ctx context.Context, groupID string) (*DeleteSecurityGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "security group ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-security-groups/%s", groupID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteSecurityGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateSecurityRuleRequest is the request for creating a security rule
type CreateSecurityRuleRequest struct {
	Description string   `json:"description,omitempty"`
	Direction   string   `json:"direction"`
	EtherType   string   `json:"etherType"`
	Port        RulePort `json:"port"`
	CIDR        string   `json:"cidr"`
}

// CreateSecurityRuleResponse is the response for CreateSecurityRule
type CreateSecurityRuleResponse struct {
	PostgreSQLResponse
	RuleID string `json:"ruleId"`
}

// CreateSecurityRule creates a new security rule in a security group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#db-security-group_6
func (c *Client) CreateSecurityRule(ctx context.Context, groupID string, req *CreateSecurityRuleRequest) (*CreateSecurityRuleResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "security group ID is required"}
	}
	if req.Direction == "" {
		return nil, &core.ValidationError{Field: "Direction", Message: "direction is required"}
	}
	if req.EtherType == "" {
		return nil, &core.ValidationError{Field: "EtherType", Message: "ether type is required"}
	}
	if req.CIDR == "" {
		return nil, &core.ValidationError{Field: "CIDR", Message: "CIDR is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-security-groups/%s/rules", groupID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateSecurityRuleResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateSecurityRuleRequest is the request for updating a security rule
type UpdateSecurityRuleRequest struct {
	Description *string   `json:"description,omitempty"`
	Direction   *string   `json:"direction,omitempty"`
	EtherType   *string   `json:"etherType,omitempty"`
	Port        *RulePort `json:"port,omitempty"`
	CIDR        *string   `json:"cidr,omitempty"`
}

// UpdateSecurityRuleResponse is the response for UpdateSecurityRule
type UpdateSecurityRuleResponse struct {
	PostgreSQLResponse
}

// UpdateSecurityRule updates a security rule in a security group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#db-security-group_7
func (c *Client) UpdateSecurityRule(ctx context.Context, groupID, ruleID string, req *UpdateSecurityRuleRequest) (*UpdateSecurityRuleResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "security group ID is required"}
	}
	if ruleID == "" {
		return nil, &core.ValidationError{Field: "ruleID", Message: "rule ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-security-groups/%s/rules/%s", groupID, ruleID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result UpdateSecurityRuleResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteSecurityRuleResponse is the response for DeleteSecurityRule
type DeleteSecurityRuleResponse struct {
	PostgreSQLResponse
}

// DeleteSecurityRule deletes a security rule from a security group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#db-security-group_8
func (c *Client) DeleteSecurityRule(ctx context.Context, groupID, ruleID string) (*DeleteSecurityRuleResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "security group ID is required"}
	}
	if ruleID == "" {
		return nil, &core.ValidationError{Field: "ruleID", Message: "rule ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-security-groups/%s/rules/%s", groupID, ruleID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteSecurityRuleResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
