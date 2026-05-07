package postgresql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// HBARuleConnectionType represents the connection type for HBA rules
type HBARuleConnectionType string

const (
	HBARuleConnectionTypeHost      HBARuleConnectionType = "HOST"
	HBARuleConnectionTypeHostSSL   HBARuleConnectionType = "HOSTSSL"
	HBARuleConnectionTypeHostNoSSL HBARuleConnectionType = "HOSTNOSSL"
)

// HBARuleApplyType represents how rules are applied
type HBARuleApplyType string

const (
	HBARuleApplyTypeEntire   HBARuleApplyType = "ENTIRE"
	HBARuleApplyTypeSelected HBARuleApplyType = "SELECTED"
)

// HBARuleUserApplyType represents how user rules are applied
type HBARuleUserApplyType string

const (
	HBARuleUserApplyTypeEntire     HBARuleUserApplyType = "ENTIRE"
	HBARuleUserApplyTypeUserCustom HBARuleUserApplyType = "USER_CUSTOM"
)

// HBARuleAuthMethod represents authentication methods
type HBARuleAuthMethod string

const (
	HBARuleAuthMethodScramSHA256 HBARuleAuthMethod = "SCRAM_SHA_256"
	HBARuleAuthMethodMD5         HBARuleAuthMethod = "MD5"
	HBARuleAuthMethodTrust       HBARuleAuthMethod = "TRUST"
)

// HBARule represents a pg_hba.conf access control rule
type HBARule struct {
	HBARuleID         string                `json:"hbaRuleId"`
	Order             int                   `json:"order"`
	ConnectionType    HBARuleConnectionType `json:"connectionType"`
	DatabaseApplyType HBARuleApplyType      `json:"databaseApplyType"`
	DBUserApplyType   HBARuleUserApplyType  `json:"dbUserApplyType"`
	Address           string                `json:"address"` // CIDR notation
	AuthMethod        HBARuleAuthMethod     `json:"authMethod"`
	DatabaseIDs       []string              `json:"databaseIds,omitempty"`
	DBUserIDs         []string              `json:"dbUserIds,omitempty"`
	Applicable        bool                  `json:"applicable"`
}

// ListHBARulesResponse is the response for ListHBARules
type ListHBARulesResponse struct {
	PostgreSQLResponse
	HBARules []HBARule `json:"hbaRules"`
}

// ListHBARules retrieves all HBA rules for a PostgreSQL instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#list-hba-rules
func (c *Client) ListHBARules(ctx context.Context, instanceID string) (*ListHBARulesResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/hba-rules", instanceID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListHBARulesResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateHBARuleRequest is the request for creating an HBA rule
type CreateHBARuleRequest struct {
	ConnectionType    string   `json:"connectionType,omitempty"` // HOST, HOSTSSL, HOSTNOSSL
	DatabaseApplyType string   `json:"databaseApplyType"`        // ENTIRE or SELECTED
	DBUserApplyType   string   `json:"dbUserApplyType"`          // ENTIRE or USER_CUSTOM
	Address           string   `json:"address"`                  // CIDR (e.g., 0.0.0.0/0)
	AuthMethod        string   `json:"authMethod"`               // SCRAM_SHA_256, MD5, TRUST
	DatabaseIDs       []string `json:"databaseIds,omitempty"`    // When SELECTED
	DBUserIDs         []string `json:"dbUserIds,omitempty"`      // When USER_CUSTOM
}

// CreateHBARuleResponse is the response for CreateHBARule
type CreateHBARuleResponse struct {
	PostgreSQLResponse
	HBARuleID string `json:"hbaRuleId"`
}

// CreateHBARule creates a new HBA rule for a PostgreSQL instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#create-hba-rule
func (c *Client) CreateHBARule(ctx context.Context, instanceID string, req *CreateHBARuleRequest) (*CreateHBARuleResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if req.DatabaseApplyType == "" {
		return nil, &core.ValidationError{Field: "DatabaseApplyType", Message: "database apply type is required"}
	}
	if req.DBUserApplyType == "" {
		return nil, &core.ValidationError{Field: "DBUserApplyType", Message: "DB user apply type is required"}
	}
	if req.Address == "" {
		return nil, &core.ValidationError{Field: "Address", Message: "address (CIDR) is required"}
	}
	if req.AuthMethod == "" {
		return nil, &core.ValidationError{Field: "AuthMethod", Message: "auth method is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/hba-rules", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateHBARuleResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ModifyHBARuleRequest is the request for modifying an HBA rule
type ModifyHBARuleRequest struct {
	ConnectionType    *string  `json:"connectionType,omitempty"`
	DatabaseApplyType *string  `json:"databaseApplyType,omitempty"`
	DBUserApplyType   *string  `json:"dbUserApplyType,omitempty"`
	Address           *string  `json:"address,omitempty"`
	AuthMethod        *string  `json:"authMethod,omitempty"`
	DatabaseIDs       []string `json:"databaseIds,omitempty"`
	DBUserIDs         []string `json:"dbUserIds,omitempty"`
}

// ModifyHBARuleResponse is the response for ModifyHBARule
type ModifyHBARuleResponse struct {
	PostgreSQLResponse
}

// ModifyHBARule modifies an existing HBA rule.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#modify-hba-rule
func (c *Client) ModifyHBARule(ctx context.Context, instanceID, ruleID string, req *ModifyHBARuleRequest) (*ModifyHBARuleResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if ruleID == "" {
		return nil, &core.ValidationError{Field: "ruleID", Message: "rule ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/hba-rules/%s", instanceID, ruleID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result ModifyHBARuleResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteHBARuleResponse is the response for DeleteHBARule
type DeleteHBARuleResponse struct {
	PostgreSQLResponse
}

// DeleteHBARule deletes an HBA rule.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#delete-hba-rule
func (c *Client) DeleteHBARule(ctx context.Context, instanceID, ruleID string) (*DeleteHBARuleResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if ruleID == "" {
		return nil, &core.ValidationError{Field: "ruleID", Message: "rule ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/hba-rules/%s", instanceID, ruleID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteHBARuleResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ReorderHBARulesRequest is the request for reordering HBA rules
type ReorderHBARulesRequest struct {
	HBARuleIDs []string `json:"hbaRuleIds"` // New order of rule IDs
}

// ReorderHBARulesResponse is the response for ReorderHBARules
type ReorderHBARulesResponse struct {
	PostgreSQLResponse
}

// ReorderHBARules changes the order of HBA rules.
// Rule order is important as PostgreSQL processes rules from top to bottom.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#reorder-hba-rules
func (c *Client) ReorderHBARules(ctx context.Context, instanceID string, req *ReorderHBARulesRequest) (*ReorderHBARulesResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if len(req.HBARuleIDs) == 0 {
		return nil, &core.ValidationError{Field: "HBARuleIDs", Message: "at least one rule ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/hba-rules/orders", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result ReorderHBARulesResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ApplyHBARulesResponse is the response for ApplyHBARules
type ApplyHBARulesResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// ApplyHBARules applies pending HBA rule changes to the PostgreSQL instance.
// Changes are not active until this method is called.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/#apply-hba-rules
func (c *Client) ApplyHBARules(ctx context.Context, instanceID string) (*ApplyHBARulesResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/hba-rules/apply", instanceID)
	req, err := http.NewRequestWithContext(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ApplyHBARulesResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
