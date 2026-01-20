package mariadb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// CreateInstanceRequest is the request for creating a MariaDB instance
// All fields from official API spec: https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/
type CreateInstanceRequest struct {
	DBInstanceName          string                      `json:"dbInstanceName"`
	DBInstanceCandidateName string                      `json:"dbInstanceCandidateName,omitempty"`
	Description             string                      `json:"description,omitempty"`
	DBFlavorID              string                      `json:"dbFlavorId"`
	DBVersion               string                      `json:"dbVersion"`
	DBUserName              string                      `json:"dbUserName"`
	DBPassword              string                      `json:"dbPassword"`
	DBPort                  *int                        `json:"dbPort,omitempty"`
	ParameterGroupID        string                      `json:"parameterGroupId"`
	DBSecurityGroupIDs      []string                    `json:"dbSecurityGroupIds,omitempty"`
	UserGroupIDs            []string                    `json:"userGroupIds,omitempty"`
	NotificationGroupIDs    []string                    `json:"notificationGroupIds,omitempty"`
	Network                 CreateInstanceNetworkConfig `json:"network"`
	Storage                 CreateInstanceStorageConfig `json:"storage"`
	Backup                  CreateInstanceBackupConfig  `json:"backup"`
	UseHighAvailability     *bool                       `json:"useHighAvailability,omitempty"`
	ReplicationMode         string                      `json:"replicationMode,omitempty"`
	UseDeletionProtection   *bool                       `json:"useDeletionProtection,omitempty"`
	AuthenticationPlugin    string                      `json:"authenticationPlugin,omitempty"` // MariaDB: NATIVE, SHA256
	// Note: TLSOption not supported in MariaDB
}

// CreateInstanceNetworkConfig specifies network configuration for instance creation
type CreateInstanceNetworkConfig struct {
	SubnetID         string `json:"subnetId"`
	AvailabilityZone string `json:"availabilityZone"`
	UsePublicAccess  *bool  `json:"usePublicAccess,omitempty"`
}

// CreateInstanceStorageConfig specifies storage configuration
type CreateInstanceStorageConfig struct {
	StorageType string `json:"storageType"`
	StorageSize int    `json:"storageSize"`
}

// CreateInstanceBackupConfig specifies backup configuration
type CreateInstanceBackupConfig struct {
	BackupPeriod     int                            `json:"backupPeriod"`
	BackupSchedules  []CreateInstanceBackupSchedule `json:"backupSchedules"`
	BackupRetryCount *int                           `json:"backupRetryCount,omitempty"`
}

// CreateInstanceBackupSchedule specifies a backup schedule
type CreateInstanceBackupSchedule struct {
	BackupWndBgnTime  string `json:"backupWndBgnTime"`
	BackupWndDuration string `json:"backupWndDuration"`
}

// Validate validates the create instance request per official API constraints
func (r *CreateInstanceRequest) Validate() error {
	if r.DBInstanceName == "" {
		return &core.ValidationError{Field: "DBInstanceName", Message: "instance name is required"}
	}
	if len(r.DBInstanceName) > 100 {
		return &core.ValidationError{Field: "DBInstanceName", Message: "instance name must be <= 100 characters"}
	}
	if r.DBFlavorID == "" {
		return &core.ValidationError{Field: "DBFlavorID", Message: "flavor ID is required"}
	}
	if r.DBVersion == "" {
		return &core.ValidationError{Field: "DBVersion", Message: "DB version is required"}
	}
	if r.DBUserName == "" {
		return &core.ValidationError{Field: "DBUserName", Message: "username is required"}
	}
	if r.DBPassword == "" {
		return &core.ValidationError{Field: "DBPassword", Message: "password is required"}
	}
	// API constraint: password 4-16 characters (official spec)
	if len(r.DBPassword) < 4 || len(r.DBPassword) > 16 {
		return &core.ValidationError{
			Field:   "DBPassword",
			Message: "password must be 4-16 characters (per API spec)",
		}
	}
	// MariaDB port range: 3306-43306 (official spec)
	if r.DBPort != nil && (*r.DBPort < 3306 || *r.DBPort > 43306) {
		return &core.ValidationError{
			Field:   "DBPort",
			Message: "port must be between 3306 and 43306 (MariaDB spec)",
		}
	}
	if r.ParameterGroupID == "" {
		return &core.ValidationError{Field: "ParameterGroupID", Message: "parameter group ID is required"}
	}
	if r.Network.SubnetID == "" {
		return &core.ValidationError{Field: "Network.SubnetID", Message: "subnet ID is required"}
	}
	if r.Network.AvailabilityZone == "" {
		return &core.ValidationError{Field: "Network.AvailabilityZone", Message: "availability zone is required"}
	}
	if r.Storage.StorageType == "" {
		return &core.ValidationError{Field: "Storage.StorageType", Message: "storage type is required"}
	}
	if r.Storage.StorageSize < 20 || r.Storage.StorageSize > 2048 {
		return &core.ValidationError{Field: "Storage.StorageSize", Message: "storage size must be 20-2048 GB"}
	}
	if r.Backup.BackupPeriod < 0 || r.Backup.BackupPeriod > 730 {
		return &core.ValidationError{Field: "Backup.BackupPeriod", Message: "backup period must be 0-730 days"}
	}
	if len(r.Backup.BackupSchedules) == 0 {
		return &core.ValidationError{Field: "Backup.BackupSchedules", Message: "at least one backup schedule is required"}
	}
	return nil
}

// CreateInstanceResponse is the response for CreateInstance
type CreateInstanceResponse struct {
	MariaDBResponse
	JobID string `json:"jobId"`
}

// CreateInstance creates a new MariaDB database instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db_3
func (c *Client) CreateInstance(ctx context.Context, req *CreateInstanceRequest) (*CreateInstanceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "/v3.0/db-instances", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateInstanceResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ModifyInstanceRequest is the request for modifying an instance
type ModifyInstanceRequest struct {
	DBInstanceName          *string  `json:"dbInstanceName,omitempty"`
	DBInstanceCandidateName *string  `json:"dbInstanceCandidateName,omitempty"`
	Description             *string  `json:"description,omitempty"`
	DBPort                  *int     `json:"dbPort,omitempty"`
	DBVersion               *string  `json:"dbVersion,omitempty"`
	UseDummy                *bool    `json:"useDummy,omitempty"`
	DBFlavorID              *string  `json:"dbFlavorId,omitempty"`
	ParameterGroupID        *string  `json:"parameterGroupId,omitempty"`
	DBSecurityGroupIDs      []string `json:"dbSecurityGroupIds,omitempty"`
	ExecuteBackup           *bool    `json:"executeBackup,omitempty"`
	UseOnlineFailover       *bool    `json:"useOnlineFailover,omitempty"`
}

// ModifyInstanceResponse is the response for ModifyInstance
type ModifyInstanceResponse struct {
	MariaDBResponse
	JobID string `json:"jobId"`
}

// ModifyInstance modifies an existing MariaDB instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db_4
func (c *Client) ModifyInstance(ctx context.Context, instanceID string, req *ModifyInstanceRequest) (*ModifyInstanceResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result ModifyInstanceResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteInstanceResponse is the response for DeleteInstance
type DeleteInstanceResponse struct {
	MariaDBResponse
	JobID string `json:"jobId"`
}

// DeleteInstance deletes a MariaDB database instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#db_5
func (c *Client) DeleteInstance(ctx context.Context, instanceID string) (*DeleteInstanceResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s", instanceID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteInstanceResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
