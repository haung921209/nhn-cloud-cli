package mariadb

import "github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"

// MariaDBResponse is the common response wrapper for MariaDB APIs
type MariaDBResponse struct {
	Header core.ResponseHeader `json:"header"`
}

// GetHeader implements core.WithHeader
func (r *MariaDBResponse) GetHeader() *core.ResponseHeader {
	return &r.Header
}

// InstanceStatus represents the status of a database instance
type InstanceStatus string

const (
	InstanceStatusAvailable    InstanceStatus = "AVAILABLE"
	InstanceStatusBeforeCreate InstanceStatus = "BEFORE_CREATE"
	InstanceStatusCreating     InstanceStatus = "CREATING"
	InstanceStatusModifying    InstanceStatus = "MODIFYING"
	InstanceStatusDeleting     InstanceStatus = "DELETING"
	InstanceStatusFailed       InstanceStatus = "FAILED"
	InstanceStatusFailToCreate InstanceStatus = "FAIL_TO_CREATE"
	InstanceStatusStopped      InstanceStatus = "STOPPED"
	InstanceStatusStopping     InstanceStatus = "STOPPING"
	InstanceStatusStarting     InstanceStatus = "STARTING"
	InstanceStatusRestarting   InstanceStatus = "RESTARTING"
	InstanceStatusBackingUp    InstanceStatus = "BACKING_UP"
	InstanceStatusRestoring    InstanceStatus = "RESTORING"
	InstanceStatusReplicating  InstanceStatus = "REPLICATING"
	InstanceStatusFailoverIng  InstanceStatus = "FAILOVER_ING"
)

// DatabaseInstance represents a MariaDB database instance
// All fields from official API specification
type DatabaseInstance struct {
	DBInstanceID          string                  `json:"dbInstanceId"`
	DBInstanceName        string                  `json:"dbInstanceName"`
	DBInstanceDescription string                  `json:"dbInstanceDescription"`
	DBInstanceStatus      InstanceStatus          `json:"dbInstanceStatus"`
	DBVersion             string                  `json:"dbVersion"`
	DBPort                int                     `json:"dbPort"`
	DBFlavorID            string                  `json:"dbFlavorId"`
	DBFlavorName          string                  `json:"dbFlavorName,omitempty"`
	ParameterGroupID      string                  `json:"parameterGroupId"`
	ParameterGroupName    string                  `json:"parameterGroupName,omitempty"`
	DBSecurityGroupIDs    []string                `json:"dbSecurityGroupIds,omitempty"`
	DBSecurityGroupNames  []string                `json:"dbSecurityGroupNames,omitempty"`
	UserGroupIDs          []string                `json:"userGroupIds,omitempty"`
	NotificationGroupIDs  []string                `json:"notificationGroupIds,omitempty"`
	Network               DatabaseInstanceNetwork `json:"network,omitempty"`
	Storage               DatabaseInstanceStorage `json:"storage,omitempty"`
	Backup                DatabaseInstanceBackup  `json:"backup,omitempty"`
	HighAvailability      *DatabaseInstanceHA     `json:"highAvailability,omitempty"`
	ReadReplicaCount      int                     `json:"readReplicaCount,omitempty"`
	ProgressStatus        string                  `json:"progressStatus,omitempty"`
	DeletionProtection    bool                    `json:"deletionProtection,omitempty"`
	CreatedAt             string                  `json:"createdAt"`
	UpdatedAt             string                  `json:"updatedAt"`
}

// DatabaseInstanceNetwork represents network configuration
type DatabaseInstanceNetwork struct {
	SubnetID         string `json:"subnetId"`
	SubnetName       string `json:"subnetName,omitempty"`
	AvailabilityZone string `json:"availabilityZone"`
	UsePublicAccess  bool   `json:"usePublicAccess"`
	DomainName       string `json:"domainName,omitempty"`
	IPAddress        string `json:"ipAddress,omitempty"`
}

// DatabaseInstanceStorage represents storage configuration
type DatabaseInstanceStorage struct {
	StorageType string `json:"storageType"`
	StorageSize int    `json:"storageSize"`
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
