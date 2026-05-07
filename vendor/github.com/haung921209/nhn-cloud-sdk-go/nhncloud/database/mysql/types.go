package mysql

import "github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"

// MySQLResponse is the common response wrapper for MySQL APIs
type MySQLResponse struct {
	Header core.ResponseHeader `json:"header"`
}

// GetHeader implements core.WithHeader
func (r *MySQLResponse) GetHeader() *core.ResponseHeader {
	return &r.Header
}

// InstanceStatus represents the status of a database instance
type InstanceStatus string

const (
	InstanceStatusAvailable       InstanceStatus = "AVAILABLE"
	InstanceStatusBeforeCreate    InstanceStatus = "BEFORE_CREATE"
	InstanceStatusStorageFull     InstanceStatus = "STORAGE_FULL"
	InstanceStatusFailToCreate    InstanceStatus = "FAIL_TO_CREATE"
	InstanceStatusFailToConnect   InstanceStatus = "FAIL_TO_CONNECT"
	InstanceStatusReplicationStop InstanceStatus = "REPLICATION_STOP"
	InstanceStatusFailover        InstanceStatus = "FAILOVER"
	InstanceStatusShutdown        InstanceStatus = "SHUTDOWN"
	InstanceStatusDeleted         InstanceStatus = "DELETED"
)

// DatabaseInstance represents a MySQL database instance
// All fields from official API v4.0 specification
type DatabaseInstance struct {
	DBInstanceID          string                   `json:"dbInstanceId"`
	DBInstanceGroupID     string                   `json:"dbInstanceGroupId,omitempty"`
	DBInstanceName        string                   `json:"dbInstanceName"`
	Description           string                   `json:"description"`
	DBInstanceType        string                   `json:"dbInstanceType,omitempty"`
	DBInstanceStatus      InstanceStatus           `json:"dbInstanceStatus"`
	DBVersion             string                   `json:"dbVersion"`
	DBPort                int                      `json:"dbPort"`
	DBFlavorID            string                   `json:"dbFlavorId"`
	DBFlavorName          string                   `json:"dbFlavorName,omitempty"`
	ParameterGroupID      string                   `json:"parameterGroupId"`
	ParameterGroupName    string                   `json:"parameterGroupName,omitempty"`
	DBSecurityGroupIDs    []string                 `json:"dbSecurityGroupIds,omitempty"`
	DBSecurityGroupNames  []string                 `json:"dbSecurityGroupNames,omitempty"`
	UserGroupIDs          []string                 `json:"userGroupIds,omitempty"`
	NotificationGroupIDs  []string                 `json:"notificationGroupIds,omitempty"`
	Network               *DatabaseInstanceNetwork `json:"network,omitempty"`
	Storage               *DatabaseInstanceStorage `json:"storage,omitempty"`
	Backup                *DatabaseInstanceBackup  `json:"backup,omitempty"`
	HighAvailability      *DatabaseInstanceHA      `json:"highAvailability,omitempty"`
	ReadReplicaCount      int                      `json:"readReplicaCount,omitempty"`
	ProgressStatus        string                   `json:"progressStatus,omitempty"`
	UseDeletionProtection bool                     `json:"useDeletionProtection,omitempty"`
	UseSlowQueryAnalysis  bool                     `json:"useSlowQueryAnalysis,omitempty"`
	AuthenticationPlugin  string                   `json:"authenticationPlugin,omitempty"`
	NeedToApplyParamGroup bool                     `json:"needToApplyParameterGroup,omitempty"`
	NeedMigration         bool                     `json:"needMigration,omitempty"`
	SupportUpgrade        bool                     `json:"supportUpgrade,omitempty"`
	CreatedYmdt           string                   `json:"createdYmdt"`
	UpdatedYmdt           string                   `json:"updatedYmdt"`
}

// DatabaseInstanceNetwork represents network configuration
type DatabaseInstanceNetwork struct {
	SubnetID         string `json:"subnetId"`
	SubnetName       string `json:"subnetName,omitempty"`
	AvailabilityZone string `json:"availabilityZone"`
	UsePublicAccess  bool   `json:"usePublicAccess"`
	DomainName       string `json:"domainName,omitempty"`
	IPAddress        string `json:"ipAddress,omitempty"`
	FloatingIP       string `json:"floatingIp,omitempty"`
	PublicIP         string `json:"publicIp,omitempty"`
}

// DatabaseInstanceStorage represents storage configuration
type DatabaseInstanceStorage struct {
	StorageType string `json:"storageType"`
	StorageSize int    `json:"storageSize"`
	// Ref: docs/api-specs/database/rds-mysql-v4.0.md#db-인스턴스-생성하기 (storage block)
	StorageAutoscale *StorageAutoscale `json:"storageAutoscale,omitempty"`
}

// StorageAutoscale represents the data-storage auto-scale configuration block.
// Optional sub-block of `storage` on instance create / modify / restore /
// replicate endpoints, and a top-level field on the storage-info endpoints.
//
// Ref: docs/api-specs/database/rds-mysql-v4.0.md#db-인스턴스-생성하기
// Ref: docs/api-specs/database/rds-mysql-v4.0.md#데이터-스토리지-정보-보기
// Ref: docs/api-specs/database/rds-mysql-v4.0.md#데이터-스토리지-정보-수정하기
type StorageAutoscale struct {
	// UseStorageAutoscale: 스토리지 자동 확장 여부
	UseStorageAutoscale *bool `json:"useStorageAutoscale,omitempty"`
	// Threshold: 자동 확장 조건(%) (50–95)
	Threshold *int `json:"threshold,omitempty"`
	// MaxStorageSize: 자동 확장 최대 크기(GB) (max 4096)
	MaxStorageSize *int `json:"maxStorageSize,omitempty"`
	// CooldownTime: 자동 확장 쿨다운 시간(분) (10–1440). Spec spelling is "cooldownTime".
	CooldownTime *int `json:"cooldownTime,omitempty"`
}

// DatabaseInstanceBackup represents backup configuration
type DatabaseInstanceBackup struct {
	BackupPeriod     int              `json:"backupPeriod"`
	BackupSchedules  []BackupSchedule `json:"backupSchedules,omitempty"`
	BackupRetryCount int              `json:"backupRetryCount,omitempty"`
}

// BackupSchedule represents a backup schedule
type BackupSchedule struct {
	BackupWndBgnTime  string `json:"backupWndBgnTime"`
	BackupWndDuration string `json:"backupWndDuration"`
}

// DatabaseInstanceHA represents high availability configuration
type DatabaseInstanceHA struct {
	Use               bool   `json:"use"`
	AvailabilityZone  string `json:"availabilityZone,omitempty"`
	CandidateMasterID string `json:"candidateMasterId,omitempty"`
	PingInterval      int    `json:"pingInterval,omitempty"`
	ReplicationMode   string `json:"replicationMode,omitempty"`
}
