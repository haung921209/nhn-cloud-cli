// Package nas provides NAS (Network Attached Storage) service client for NHN Cloud.
package nas

import (
	"time"
)

// Header represents API response header
type Header struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// Paging represents pagination information
type Paging struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalCount int `json:"totalCount"`
}

// EncryptionKey represents encryption key information
type EncryptionKey struct {
	KeyID   string `json:"keyId"`
	Version int    `json:"version"`
}

// Encryption represents volume encryption settings
type Encryption struct {
	Enabled bool            `json:"enabled"`
	Keys    []EncryptionKey `json:"keys"`
}

// Interface represents volume interface configuration
type Interface struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Status   string `json:"status"`
	SubnetID string `json:"subnetId"`
	TenantID string `json:"tenantId"`
}

// VolumeMirror represents volume replication configuration
type VolumeMirror struct {
	ID                 string    `json:"id"`
	DirectionChangedAt *string   `json:"directionChangedAt"`
	DstProjectID       string    `json:"dstProjectId"`
	DstRegion          string    `json:"dstRegion"`
	DstTenantID        string    `json:"dstTenantId"`
	DstVolumeID        string    `json:"dstVolumeId"`
	DstVolumeName      string    `json:"dstVolumeName"`
	SrcProjectID       string    `json:"srcProjectId"`
	SrcRegion          string    `json:"srcRegion"`
	SrcTenantID        string    `json:"srcTenantId"`
	SrcVolumeID        string    `json:"srcVolumeId"`
	SrcVolumeName      string    `json:"srcVolumeName"`
	CreatedAt          time.Time `json:"createdAt"`
}

// SnapshotSchedule represents snapshot scheduling configuration
type SnapshotSchedule struct {
	Time       string `json:"time"`
	TimeOffset string `json:"timeOffset"`
}

// SnapshotPolicy represents volume snapshot policy
type SnapshotPolicy struct {
	MaxScheduledCount int              `json:"maxScheduledCount"`
	ReservePercent    int              `json:"reservePercent"`
	Schedule          SnapshotSchedule `json:"schedule"`
}

// MountProtocol represents volume mount protocol settings
type MountProtocol struct {
	CIFSAuthIDs []string `json:"cifsAuthIds"`
	Protocol    string   `json:"protocol"`
}

// Volume represents a NAS volume
type Volume struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Status         string         `json:"status"`
	Description    *string        `json:"description"`
	SizeGB         int            `json:"sizeGb"`
	ProjectID      string         `json:"projectId"`
	TenantID       string         `json:"tenantId"`
	ACL            []string       `json:"acl"`
	Encryption     Encryption     `json:"encryption"`
	Interfaces     []Interface    `json:"interfaces"`
	Mirrors        []VolumeMirror `json:"mirrors"`
	MountProtocol  MountProtocol  `json:"mountProtocol"`
	SnapshotPolicy SnapshotPolicy `json:"snapshotPolicy"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

// Usage represents volume usage information
type Usage struct {
	SnapshotReserveGB int `json:"snapshotReserveGb"`
	UsedGB            int `json:"usedGb"`
}

// Snapshot represents a volume snapshot
type Snapshot struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Preserved        bool      `json:"preserved"`
	Size             int       `json:"size"`
	ReclaimableSpace *int      `json:"reclaimableSpace,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
}

// RestoreHistory represents snapshot restore history
type RestoreHistory struct {
	RequestedAt   time.Time  `json:"requestedAt"`
	RequestedIP   string     `json:"requestedIp"`
	RequestedUser string     `json:"requestedUser"`
	RestoredAt    *time.Time `json:"restoredAt"`
	Result        string     `json:"result"`
	SnapshotID    string     `json:"snapshotId"`
	SnapshotName  string     `json:"snapshotName"`
	VolumeID      string     `json:"volumeId"`
}

// VolumeMirrorStat represents volume mirror statistics
type VolumeMirrorStat struct {
	LastSuccessTransferBytes   int       `json:"lastSuccessTransferBytes"`
	LastSuccessTransferEndTime time.Time `json:"lastSuccessTransferEndTime"`
	LastTransferBytes          int       `json:"lastTransferBytes"`
	LastTransferEndTime        time.Time `json:"lastTransferEndTime"`
	LastTransferStatus         string    `json:"lastTransferStatus"`
}

// --- Input Types ---

// ListVolumesInput represents options for listing volumes
type ListVolumesInput struct {
	SizeGB       *int   `json:"sizeGb,omitempty"`
	MaxSizeGB    *int   `json:"maxSizeGb,omitempty"`
	MinSizeGB    *int   `json:"minSizeGb,omitempty"`
	Name         string `json:"name,omitempty"`
	NameContains string `json:"nameContains,omitempty"`
	SubnetID     string `json:"subnetId,omitempty"`
	Limit        *int   `json:"limit,omitempty"`
	Page         *int   `json:"page,omitempty"`
}

// CreateVolumeInterfaceInput represents interface configuration for volume creation
type CreateVolumeInterfaceInput struct {
	SubnetID string `json:"subnetId"`
}

// CreateVolumeInput represents the request to create a volume
type CreateVolumeInput struct {
	Name           string                       `json:"name"`
	Description    *string                      `json:"description,omitempty"`
	SizeGB         int                          `json:"sizeGb"`
	Encryption     *Encryption                  `json:"encryption,omitempty"`
	Interfaces     []CreateVolumeInterfaceInput `json:"interfaces,omitempty"`
	MountProtocol  *MountProtocol               `json:"mountProtocol,omitempty"`
	SnapshotPolicy *SnapshotPolicy              `json:"snapshotPolicy,omitempty"`
}

// UpdateVolumeInput represents the request to update a volume
type UpdateVolumeInput struct {
	Description    *string         `json:"description,omitempty"`
	MountProtocol  *MountProtocol  `json:"mountProtocol,omitempty"`
	SnapshotPolicy *SnapshotPolicy `json:"snapshotPolicy,omitempty"`
}

// CreateInterfaceInput represents the request to create an interface
type CreateInterfaceInput struct {
	SubnetID string `json:"subnetId"`
}

// CreateSnapshotInput represents the request to create a snapshot
type CreateSnapshotInput struct {
	Name string `json:"name"`
}

// ListRestoreHistoriesInput represents options for listing restore histories
type ListRestoreHistoriesInput struct {
	Limit *int `json:"limit,omitempty"`
	Page  *int `json:"page,omitempty"`
}

// CreateVolumeMirrorDstInput represents destination volume for mirror
type CreateVolumeMirrorDstInput struct {
	Name           string                       `json:"name"`
	Description    *string                      `json:"description,omitempty"`
	Encryption     *Encryption                  `json:"encryption,omitempty"`
	Interfaces     []CreateVolumeInterfaceInput `json:"interfaces,omitempty"`
	MountProtocol  *MountProtocol               `json:"mountProtocol,omitempty"`
	SnapshotPolicy *SnapshotPolicy              `json:"snapshotPolicy,omitempty"`
}

// CreateVolumeMirrorInput represents the request to create volume mirror
type CreateVolumeMirrorInput struct {
	DstRegion   string                     `json:"dstRegion"`
	DstTenantID string                     `json:"dstTenantId"`
	DstVolume   CreateVolumeMirrorDstInput `json:"dstVolume"`
}

// --- Output Types ---

// ListVolumesOutput represents the response for volumes list
type ListVolumesOutput struct {
	Header  Header   `json:"header"`
	Paging  Paging   `json:"paging"`
	Volumes []Volume `json:"volumes"`
}

// GetVolumeOutput represents the response for single volume
type GetVolumeOutput struct {
	Header Header `json:"header"`
	Volume Volume `json:"volume"`
}

// CreateVolumeOutput represents the response for volume creation
type CreateVolumeOutput struct {
	Header Header `json:"header"`
	Volume Volume `json:"volume"`
}

// UpdateVolumeOutput represents the response for volume update
type UpdateVolumeOutput struct {
	Header Header `json:"header"`
	Volume Volume `json:"volume"`
}

// GetUsageOutput represents the response for volume usage
type GetUsageOutput struct {
	Header Header `json:"header"`
	Usage  Usage  `json:"usage"`
}

// CreateInterfaceOutput represents the response for interface creation
type CreateInterfaceOutput struct {
	Header    Header    `json:"header"`
	Interface Interface `json:"interface"`
}

// ListSnapshotsOutput represents the response for snapshots list
type ListSnapshotsOutput struct {
	Header    Header     `json:"header"`
	Snapshots []Snapshot `json:"snapshots"`
}

// GetSnapshotOutput represents the response for single snapshot
type GetSnapshotOutput struct {
	Header   Header   `json:"header"`
	Snapshot Snapshot `json:"snapshot"`
}

// CreateSnapshotOutput represents the response for snapshot creation
type CreateSnapshotOutput struct {
	Header   Header   `json:"header"`
	Snapshot Snapshot `json:"snapshot"`
}

// ListRestoreHistoriesOutput represents the response for restore histories
type ListRestoreHistoriesOutput struct {
	Header           Header           `json:"header"`
	Paging           Paging           `json:"paging"`
	RestoreHistories []RestoreHistory `json:"restoreHistories"`
}

// CreateVolumeMirrorOutput represents the response for volume mirror creation
type CreateVolumeMirrorOutput struct {
	Header       Header       `json:"header"`
	VolumeMirror VolumeMirror `json:"volumeMirror"`
}

// GetVolumeMirrorStatOutput represents the response for volume mirror statistics
type GetVolumeMirrorStatOutput struct {
	Header           Header           `json:"header"`
	VolumeMirrorStat VolumeMirrorStat `json:"volumeMirrorStat"`
}

// --- Constants ---

// Volume status constants
const (
	VolumeStatusCreating  = "creating"
	VolumeStatusAvailable = "available"
	VolumeStatusDeleting  = "deleting"
	VolumeStatusError     = "error"
	VolumeStatusModifying = "modifying"
)

// Interface status constants
const (
	InterfaceStatusCreating  = "creating"
	InterfaceStatusAvailable = "available"
	InterfaceStatusDeleting  = "deleting"
	InterfaceStatusError     = "error"
)

// Mount protocol constants
const (
	ProtocolNFS  = "NFS"
	ProtocolCIFS = "CIFS"
)

// Mirror status constants
const (
	MirrorStatusSuccess    = "success"
	MirrorStatusInProgress = "in-progress"
	MirrorStatusFailed     = "failed"
)
