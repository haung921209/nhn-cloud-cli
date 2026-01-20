package postgresql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// Backup represents a database backup
type Backup struct {
	BackupID          string `json:"backupId"`
	BackupName        string `json:"backupName"`
	BackupType        string `json:"backupType,omitempty"`
	BackupStatus      string `json:"backupStatus"`
	DBInstanceID      string `json:"dbInstanceId"`
	DBInstanceName    string `json:"dbInstanceName,omitempty"`
	BackupSize        int64  `json:"backupSize,omitempty"`
	BackupStartedAt   string `json:"backupStartedAt,omitempty"`
	BackupCompletedAt string `json:"backupCompletedAt,omitempty"`
	CreatedAt         string `json:"createdAt"`
}

// ListBackupsResponse is the response for ListBackups
type ListBackupsResponse struct {
	PostgreSQLResponse
	Backups []Backup `json:"backups"`
}

// ListBackups retrieves backups for an instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#backup_1
func (c *Client) ListBackups(ctx context.Context, instanceID string) (*ListBackupsResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	// Query parameters can be added: page, size, dbVersion
	path := fmt.Sprintf("/v1.0/backups?dbInstanceId=%s", instanceID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListBackupsResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateBackupRequest is the request for creating a backup
type CreateBackupRequest struct {
	BackupName string `json:"backupName"`
}

// CreateBackupResponse is the response for CreateBackup
type CreateBackupResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// CreateBackup creates a manual backup for an instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#backup_2
func (c *Client) CreateBackup(ctx context.Context, instanceID string, req *CreateBackupRequest) (*CreateBackupResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if req.BackupName == "" {
		return nil, &core.ValidationError{Field: "BackupName", Message: "backup name is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/backup", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateBackupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// BackupToObjectStorageRequest is the request for backing up to object storage
type BackupToObjectStorageRequest struct {
	TenantID        string `json:"tenantId"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	TargetContainer string `json:"targetContainer"`
	ObjectPath      string `json:"objectPath"`
}

// BackupToObjectStorageResponse is the response for BackupToObjectStorage
type BackupToObjectStorageResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// BackupToObjectStorage backs up an instance to object storage.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_44
func (c *Client) BackupToObjectStorage(ctx context.Context, instanceID string, req *BackupToObjectStorageRequest) (*BackupToObjectStorageResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/db-instances/%s/backup-to-object-storage", instanceID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result BackupToObjectStorageResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RestoreBackupRequest is the request for restoring a backup
type RestoreBackupRequest struct {
	DBInstanceName string `json:"dbInstanceName,omitempty"`
	// Additional restore options can be added here
}

// RestoreBackupResponse is the response for RestoreBackup
type RestoreBackupResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// RestoreBackup restores an instance from a backup.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#backup_4
func (c *Client) RestoreBackup(ctx context.Context, backupID string, req *RestoreBackupRequest) (*RestoreBackupResponse, error) {
	if backupID == "" {
		return nil, &core.ValidationError{Field: "backupID", Message: "backup ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/backups/%s/restore", backupID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result RestoreBackupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ExportBackupRequest is the request for exporting a backup
type ExportBackupRequest struct {
	TenantID        string `json:"tenantId"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	TargetContainer string `json:"targetContainer"`
	ObjectPath      string `json:"objectPath"`
}

// ExportBackupResponse is the response for ExportBackup
type ExportBackupResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// ExportBackup exports a backup to object storage.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#backup_5
func (c *Client) ExportBackup(ctx context.Context, backupID string, req *ExportBackupRequest) (*ExportBackupResponse, error) {
	if backupID == "" {
		return nil, &core.ValidationError{Field: "backupID", Message: "backup ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/backups/%s/export", backupID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result ExportBackupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteBackupResponse is the response for DeleteBackup
type DeleteBackupResponse struct {
	PostgreSQLResponse
	JobID string `json:"jobId"`
}

// DeleteBackup deletes a backup.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#backup_6
func (c *Client) DeleteBackup(ctx context.Context, backupID string) (*DeleteBackupResponse, error) {
	if backupID == "" {
		return nil, &core.ValidationError{Field: "backupID", Message: "backup ID is required"}
	}

	path := fmt.Sprintf("/v1.0/backups/%s", backupID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteBackupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
