package postgresql

import "time"

// ResponseHeader represents common API response header
type ResponseHeader struct {
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
	IsSuccessful  bool   `json:"isSuccessful"`
}

// Project & Region Types
type Region struct {
	RegionCode string `json:"regionCode"`
	IsEnabled  bool   `json:"isEnabled"`
}

type RegionsResponse struct {
	Header  *ResponseHeader `json:"header"`
	Regions []Region        `json:"regions"`
}

type Member struct {
	MemberID     string `json:"memberId"`
	MemberName   string `json:"memberName"`
	EmailAddress string `json:"emailAddress"`
	PhoneNumber  string `json:"phoneNumber"`
}

type MembersResponse struct {
	Header  *ResponseHeader `json:"header"`
	Members []Member        `json:"members"`
}

// Network Types
type Subnet struct {
	SubnetID         string `json:"subnetId"`
	SubnetName       string `json:"subnetName"`
	SubnetCIDR       string `json:"subnetCidr"`
	UsingGateway     bool   `json:"usingGateway"`
	AvailableIPCount int    `json:"availableIpCount"`
}

type SubnetsResponse struct {
	Header  *ResponseHeader `json:"header"`
	Subnets []Subnet        `json:"subnets"`
}

type StorageTypesResponse struct {
	Header       *ResponseHeader `json:"header"`
	StorageTypes []string        `json:"storageTypes"`
}

// DB Flavors & Versions
type DBFlavor struct {
	DBFlavorID   string `json:"dbFlavorId"`
	DBFlavorName string `json:"dbFlavorName"`
	RAM          int    `json:"ram"`
	VCPUs        int    `json:"vcpus"`
}

type DBFlavorsResponse struct {
	Header    *ResponseHeader `json:"header"`
	DBFlavors []DBFlavor      `json:"dbFlavors"`
}

type DBVersion struct {
	DBVersionCode      string `json:"dbVersionCode"`
	DBMajorVersionCode string `json:"dbMajorVersionCode"`
	Name               string `json:"name"`
	CanCreate          bool   `json:"canCreate"`
}

type DBVersionsResponse struct {
	Header     *ResponseHeader `json:"header"`
	DBVersions []DBVersion     `json:"dbVersions"`
}

type BackupSchedule struct {
	BackupWndBgnTime      string `json:"backupWndBgnTime"`
	BackupWndDuration     string `json:"backupWndDuration"`
	BackupRetryExpireTime string `json:"backupRetryExpireTime,omitempty"`
}

// DB Instance Group
type DBInstanceGroup struct {
	DBInstanceGroupID   string    `json:"dbInstanceGroupId"`
	DBInstanceGroupName string    `json:"dbInstanceGroupName"`
	ReplicationType     string    `json:"replicationType"`
	CreatedYmdt         time.Time `json:"createdYmdt"`
	UpdatedYmdt         time.Time `json:"updatedYmdt"`
}

type DBInstanceGroupsResponse struct {
	Header           *ResponseHeader   `json:"header"`
	DBInstanceGroups []DBInstanceGroup `json:"dbInstanceGroups"`
}

type DBInstanceGroupResponse struct {
	Header *ResponseHeader `json:"header"`
	DBInstanceGroup
}

type DBInstance struct {
	DBInstanceID              string    `json:"dbInstanceId"`
	DBInstanceGroupID         string    `json:"dbInstanceGroupId"`
	DBInstanceName            string    `json:"dbInstanceName"`
	Description               string    `json:"description,omitempty"`
	DBVersion                 string    `json:"dbVersion"`
	DBPort                    int       `json:"dbPort"`
	DBInstanceType            string    `json:"dbInstanceType"`
	DBInstanceStatus          string    `json:"dbInstanceStatus"`
	ProgressStatus            string    `json:"progressStatus"`
	DBFlavorID                string    `json:"dbFlavorId"`
	ParameterGroupID          string    `json:"parameterGroupId"`
	DBSecurityGroupIDs        []string  `json:"dbSecurityGroupIds"`
	NotificationGroupIDs      []string  `json:"notificationGroupIds"`
	UseDeletionProtection     bool      `json:"useDeletionProtection"`
	NeedToApplyParameterGroup bool      `json:"needToApplyParameterGroup"`
	NeedMigration             bool      `json:"needMigration"`
	OSVersion                 string    `json:"osVersion,omitempty"`
	StorageType               string    `json:"storageType,omitempty"`
	StorageSize               int       `json:"storageSize,omitempty"`
	CreatedYmdt               time.Time `json:"createdYmdt"`
	UpdatedYmdt               time.Time `json:"updatedYmdt"`
}

type DBInstancesResponse struct {
	Header      *ResponseHeader `json:"header"`
	DBInstances []DBInstance    `json:"dbInstances"`
}

type DBInstanceResponse struct {
	Header *ResponseHeader `json:"header"`
	DBInstance
}

// Create DB Instance Request
type CreateDBInstanceRequest struct {
	DBInstanceName          string   `json:"dbInstanceName"`
	DBInstanceCandidateName string   `json:"dbInstanceCandidateName,omitempty"`
	Description             string   `json:"description,omitempty"`
	DBFlavorID              string   `json:"dbFlavorId"`
	DBVersion               string   `json:"dbVersion"`
	DBPort                  int      `json:"dbPort,omitempty"`
	DatabaseName            string   `json:"databaseName"`
	DBUserName              string   `json:"dbUserName"`
	DBPassword              string   `json:"dbPassword"`
	ParameterGroupID        string   `json:"parameterGroupId"`
	DBSecurityGroupIDs      []string `json:"dbSecurityGroupIds,omitempty"`
	UserGroupIDs            []string `json:"userGroupIds,omitempty"`
	UseHighAvailability     bool     `json:"useHighAvailability,omitempty"`
	UseDeletionProtection   bool     `json:"useDeletionProtection,omitempty"`
	Network                 struct {
		SubnetID         string `json:"subnetId"`
		AvailabilityZone string `json:"availabilityZone,omitempty"`
		UsePublicAccess  bool   `json:"usePublicAccess,omitempty"`
	} `json:"network"`
	Storage struct {
		StorageType string `json:"storageType"`
		StorageSize int    `json:"storageSize"`
	} `json:"storage"`
	Backup struct {
		BackupPeriod    int              `json:"backupPeriod"`
		BackupSchedules []BackupSchedule `json:"backupSchedules,omitempty"`
	} `json:"backup"`
}

// Modify DB Instance Request
// API: PUT /v1.0/db-instances/{dbInstanceId}
type ModifyDBInstanceRequest struct {
	DBInstanceName     string   `json:"dbInstanceName,omitempty"`
	Description        string   `json:"description,omitempty"`
	DBPort             int      `json:"dbPort,omitempty"`
	DBFlavorID         string   `json:"dbFlavorId,omitempty"`
	ParameterGroupID   string   `json:"parameterGroupId,omitempty"`
	DBSecurityGroupIDs []string `json:"dbSecurityGroupIds,omitempty"`
	ExecuteBackup      bool     `json:"executeBackup,omitempty"`
}

// ModifyStorageInfoRequest for PUT /v1.0/db-instances/{dbInstanceId}/storage-info
type ModifyStorageInfoRequest struct {
	StorageSize int `json:"storageSize"`
}

// RestartInstanceRequest represents a request to restart a database instance
type RestartInstanceRequest struct {
	// UseOnlineFailover enables restart using failover (HA instances only)
	// When true, minimizes downtime by failing over to standby before restart
	UseOnlineFailover bool `json:"useOnlineFailover,omitempty"`
	// ExecuteBackup triggers a backup before restart
	ExecuteBackup bool `json:"executeBackup,omitempty"`
}

// High Availability
type HighAvailability struct {
	UseHighAvailability bool      `json:"useHighAvailability"`
	PingInterval        int       `json:"pingInterval"`
	ReplicateStatus     string    `json:"replicateStatus"`
	ReplicateDelay      int       `json:"replicateDelay"`
	StandbyDBInstanceID string    `json:"standbyDbInstanceId"`
	CreatedYmdt         time.Time `json:"createdYmdt"`
	UpdatedYmdt         time.Time `json:"updatedYmdt"`
}

type HighAvailabilityResponse struct {
	Header *ResponseHeader `json:"header"`
	HighAvailability
}

// Extensions (PostgreSQL specific) - managed at db-instance-group level
type ExtensionDatabase struct {
	DBInstanceGroupExtensionID     string `json:"dbInstanceGroupExtensionId"`
	DBInstanceGroupExtensionStatus string `json:"dbInstanceGroupExtensionStatus"`
	DatabaseID                     string `json:"databaseId"`
	DatabaseName                   string `json:"databaseName"`
	ReservedAction                 string `json:"reservedAction"`
	ErrorReason                    string `json:"errorReason,omitempty"`
}

type Extension struct {
	ExtensionID     string              `json:"extensionId"`
	ExtensionName   string              `json:"extensionName"`
	ExtensionStatus string              `json:"extensionStatus"`
	Databases       []ExtensionDatabase `json:"databases"`
}

type ExtensionsResponse struct {
	Header        *ResponseHeader `json:"header"`
	Extensions    []Extension     `json:"extensions"`
	IsNeedToApply bool            `json:"isNeedToApply"`
}

type InstallExtensionRequest struct {
	DatabaseID  string `json:"databaseId"`
	SchemaName  string `json:"schemaName"`
	WithCascade bool   `json:"withCascade,omitempty"`
}

// Database Management (PostgreSQL specific)
type Database struct {
	DatabaseID   string    `json:"databaseId"`
	DatabaseName string    `json:"databaseName"`
	Owner        string    `json:"owner"`
	Encoding     string    `json:"encoding"`
	Collate      string    `json:"collate"`
	CType        string    `json:"ctype"`
	Size         int64     `json:"size"`
	TableCount   int       `json:"tableCount"`
	CreatedYmdt  time.Time `json:"createdYmdt"`
}

type DatabasesResponse struct {
	Header    *ResponseHeader `json:"header"`
	Databases []Database      `json:"databases"`
}

type CreateDatabaseRequest struct {
	DatabaseName string `json:"databaseName"`
	Owner        string `json:"owner"`
	Encoding     string `json:"encoding,omitempty"`
	Collate      string `json:"collate,omitempty"`
	CType        string `json:"ctype,omitempty"`
}

// User Management
type DBUser struct {
	DBUserID      string    `json:"dbUserId"`
	DBUserName    string    `json:"dbUserName"`
	AuthorityType string    `json:"authorityType"`
	DBUserStatus  string    `json:"dbUserStatus"`
	CreatedYmdt   time.Time `json:"createdYmdt"`
	UpdatedYmdt   time.Time `json:"updatedYmdt"`
}

type DBUsersResponse struct {
	Header  *ResponseHeader `json:"header"`
	DBUsers []DBUser        `json:"dbUsers"`
}

type CreateDBUserRequest struct {
	DBUserName      string `json:"dbUserName"`
	DBPassword      string `json:"dbPassword"`
	AuthorityType   string `json:"authorityType,omitempty"`
	IsSuperuser     bool   `json:"isSuperuser,omitempty"`
	CanCreateDB     bool   `json:"canCreateDb,omitempty"`
	CanCreateRole   bool   `json:"canCreateRole,omitempty"`
	CanLogin        bool   `json:"canLogin,omitempty"`
	IsReplication   bool   `json:"isReplication,omitempty"`
	ConnectionLimit int    `json:"connectionLimit,omitempty"`
}

// HBA Rules (PostgreSQL specific)
type HBARule struct {
	HBARuleID           string     `json:"hbaRuleId"`
	HBARuleStatus       string     `json:"hbaRuleStatus"`
	Order               int        `json:"order"`
	DatabaseApplyType   string     `json:"databaseApplyType"`
	DBUserApplyTypeCode string     `json:"dbUserApplyTypeCode"`
	Databases           []Database `json:"databases"`
	DBUsers             []DBUser   `json:"dbUsers"`
	Address             string     `json:"address"`
	AuthMethod          string     `json:"authMethod"`
	ReservedAction      string     `json:"reservedAction"`
	Applicable          bool       `json:"applicable"`
}

type HBARulesResponse struct {
	Header   *ResponseHeader `json:"header"`
	HBARules []HBARule       `json:"hbaRules"`
}

type CreateHBARuleRequest struct {
	// Connection type: HOST, HOSTSSL, HOSTNOSSL
	ConnectionType string `json:"connectionType,omitempty"`
	// Database apply type: ENTIRE (all databases) or SELECTED (specific databases)
	DatabaseApplyType string `json:"databaseApplyType"`
	// DB user apply type: ENTIRE (all users) or USER_CUSTOM (specific users)
	DBUserApplyType string `json:"dbUserApplyType"`
	// Address in CIDR format (e.g., 0.0.0.0/0, 192.168.1.0/24)
	Address string `json:"address"`
	// Authentication method: SCRAM_SHA_256, MD5, TRUST
	AuthMethod string `json:"authMethod"`
	// Specific database IDs (required when databaseApplyType is SELECTED)
	DatabaseIds []string `json:"databaseIds,omitempty"`
	// Specific user IDs (required when dbUserApplyType is USER_CUSTOM)
	DBUserIds []string `json:"dbUserIds,omitempty"`
}

// Backup
type Backup struct {
	BackupID       string    `json:"backupId"`
	BackupName     string    `json:"backupName"`
	DBInstanceID   string    `json:"dbInstanceId"`
	DBInstanceName string    `json:"dbInstanceName"`
	BackupType     string    `json:"backupType"`
	BackupStatus   string    `json:"backupStatus"`
	BackupSize     int64     `json:"backupSize"`
	DBVersion      string    `json:"dbVersion"`
	DBDataSize     int64     `json:"dbDataSize"`
	BinLogPosition string    `json:"binLogPosition,omitempty"`
	CreatedYmdt    time.Time `json:"createdYmdt"`
	UpdatedYmdt    time.Time `json:"updatedYmdt"`
}

type BackupsResponse struct {
	Header  *ResponseHeader `json:"header"`
	Backups []Backup        `json:"backups"`
}

type CreateBackupRequest struct {
	BackupName string `json:"backupName"`
}

type ExportBackupRequest struct {
	TenantID        string `json:"tenantId"`
	UserName        string `json:"userName"`
	Password        string `json:"password"`
	TargetContainer string `json:"targetContainer"`
	ObjectPath      string `json:"objectPath"`
}

type RestoreBackupRequest struct {
	BackupID              string                `json:"-"` // Used in URL path, not body
	DBInstanceName        string                `json:"dbInstanceName"`
	Description           string                `json:"description,omitempty"`
	DBFlavorID            string                `json:"dbFlavorId"`
	DBPort                int                   `json:"dbPort,omitempty"`
	ParameterGroupID      string                `json:"parameterGroupId"`
	DBSecurityGroupIDs    []string              `json:"dbSecurityGroupIds,omitempty"`
	UseHighAvailability   bool                  `json:"useHighAvailability,omitempty"`
	UseDeletionProtection bool                  `json:"useDeletionProtection,omitempty"`
	Network               *RestoreNetworkConfig `json:"network"`
	Storage               *RestoreStorageConfig `json:"storage"`
	Backup                *RestoreBackupConfig  `json:"backup"`
}

// RestoreNetworkConfig for restore operation
type RestoreNetworkConfig struct {
	SubnetID         string `json:"subnetId"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

// RestoreStorageConfig for restore operation
type RestoreStorageConfig struct {
	StorageType string `json:"storageType"`
	StorageSize int    `json:"storageSize"`
}

// RestoreBackupConfig for restore operation
type RestoreBackupConfig struct {
	BackupPeriod    int                     `json:"backupPeriod"`
	BackupSchedules []RestoreBackupSchedule `json:"backupSchedules"`
}

// RestoreBackupSchedule for restore operation
type RestoreBackupSchedule struct {
	BackupWndBgnTime  string `json:"backupWndBgnTime"`
	BackupWndDuration string `json:"backupWndDuration"`
}

// DB Security Groups
type DBSecurityGroup struct {
	DBSecurityGroupID   string                `json:"dbSecurityGroupId"`
	DBSecurityGroupName string                `json:"dbSecurityGroupName"`
	Description         string                `json:"description,omitempty"`
	IsDefault           bool                  `json:"isDefault"`
	ProgressStatus      string                `json:"progressStatus"`
	Rules               []DBSecurityGroupRule `json:"rules,omitempty"`
	CreatedYmdt         time.Time             `json:"createdYmdt"`
	UpdatedYmdt         time.Time             `json:"updatedYmdt"`
}

type DBSecurityGroupsResponse struct {
	Header           *ResponseHeader   `json:"header"`
	DBSecurityGroups []DBSecurityGroup `json:"dbSecurityGroups"`
}

type DBSecurityGroupResponse struct {
	Header          *ResponseHeader `json:"header"`
	DBSecurityGroup DBSecurityGroup `json:"dbSecurityGroup"`
}

type DBSecurityGroupRule struct {
	RuleID      string    `json:"ruleId"`
	Description string    `json:"description,omitempty"`
	Direction   string    `json:"direction"`
	EtherType   string    `json:"etherType"`
	Port        *PortSpec `json:"port,omitempty"`
	CIDR        string    `json:"cidr"`
	CreatedYmdt time.Time `json:"createdYmdt"`
	UpdatedYmdt time.Time `json:"updatedYmdt"`
}

type CreateDBSecurityGroupRequest struct {
	DBSecurityGroupName string `json:"dbSecurityGroupName"`
	Description         string `json:"description,omitempty"`
}

// PortSpec defines the port specification for security group rules
type PortSpec struct {
	PortType string `json:"portType"`
	MinPort  int    `json:"minPort,omitempty"`
	MaxPort  int    `json:"maxPort,omitempty"`
}

type CreateDBSecurityGroupRuleRequest struct {
	Description string    `json:"description,omitempty"`
	Direction   string    `json:"direction"`
	EtherType   string    `json:"etherType"`
	Port        *PortSpec `json:"port,omitempty"`
	Protocol    string    `json:"protocol,omitempty"`
	CIDR        string    `json:"cidr"`
}

// Parameter Groups
type ParameterGroup struct {
	ParameterGroupID   string    `json:"parameterGroupId"`
	ParameterGroupName string    `json:"parameterGroupName"`
	Description        string    `json:"description,omitempty"`
	DBVersion          string    `json:"dbVersion"`
	IsDefault          bool      `json:"isDefault"`
	CreatedYmdt        time.Time `json:"createdYmdt"`
	UpdatedYmdt        time.Time `json:"updatedYmdt"`
}

type ParameterGroupsResponse struct {
	Header          *ResponseHeader  `json:"header"`
	ParameterGroups []ParameterGroup `json:"parameterGroups"`
}

type ParameterGroupResponse struct {
	Header *ResponseHeader `json:"header"`
	ParameterGroup
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	ParameterID   string `json:"parameterId"`
	ParameterName string `json:"parameterName"`
	Value         string `json:"value"`
	DefaultValue  string `json:"defaultValue"`
	AllowedValues string `json:"allowedValues,omitempty"`
	UpdateType    string `json:"updateType"`
	ApplyType     string `json:"applyType"`
	DataType      string `json:"dataType"`
	Description   string `json:"description,omitempty"`
	MinValue      string `json:"minValue,omitempty"`
	MaxValue      string `json:"maxValue,omitempty"`
}

type CreateParameterGroupRequest struct {
	ParameterGroupName string `json:"parameterGroupName"`
	Description        string `json:"description,omitempty"`
	DBVersion          string `json:"dbVersion"`
}

type ModifyParametersRequest struct {
	ModifiedParameters []struct {
		ParameterID string `json:"parameterId"`
		Value       string `json:"value"`
	} `json:"modifiedParameters"`
}

// User Groups
type UserGroup struct {
	UserGroupID   string    `json:"userGroupId"`
	UserGroupName string    `json:"userGroupName"`
	Description   string    `json:"description,omitempty"`
	UserIDs       []string  `json:"userIds"`
	CreatedYmdt   time.Time `json:"createdYmdt"`
	UpdatedYmdt   time.Time `json:"updatedYmdt"`
}

type UserGroupsResponse struct {
	Header     *ResponseHeader `json:"header"`
	UserGroups []UserGroup     `json:"userGroups"`
}

type UserGroupResponse struct {
	Header *ResponseHeader `json:"header"`
	UserGroup
}

type CreateUserGroupRequest struct {
	UserGroupName string   `json:"userGroupName"`
	Description   string   `json:"description,omitempty"`
	UserIDs       []string `json:"userIds"`
}

// Notification Groups
type NotificationGroup struct {
	NotificationGroupID   string `json:"notificationGroupId"`
	NotificationGroupName string `json:"notificationGroupName"`
	Description           string `json:"description,omitempty"`
	NotificationType      string `json:"notificationType"`
	IsEnabled             bool   `json:"isEnabled"`
	Recipients            []struct {
		RecipientType string `json:"recipientType"`
		Recipient     string `json:"recipient"`
	} `json:"recipients"`
	CreatedYmdt time.Time `json:"createdYmdt"`
	UpdatedYmdt time.Time `json:"updatedYmdt"`
}

type NotificationGroupsResponse struct {
	Header             *ResponseHeader     `json:"header"`
	NotificationGroups []NotificationGroup `json:"notificationGroups"`
}

type NotificationGroupResponse struct {
	Header *ResponseHeader `json:"header"`
	NotificationGroup
}

type CreateNotificationGroupRequest struct {
	NotificationGroupName string `json:"notificationGroupName"`
	Description           string `json:"description,omitempty"`
	NotificationType      string `json:"notificationType"`
	IsEnabled             bool   `json:"isEnabled"`
	Recipients            []struct {
		RecipientType string `json:"recipientType"`
		Recipient     string `json:"recipient"`
	} `json:"recipients"`
}

// Watchdog (PostgreSQL specific)
type Watchdog struct {
	WatchdogID   string    `json:"watchdogId"`
	WatchdogName string    `json:"watchdogName"`
	DBInstanceID string    `json:"dbInstanceId"`
	IsEnabled    bool      `json:"isEnabled"`
	QueryTimeout int       `json:"queryTimeout"`
	CreatedYmdt  time.Time `json:"createdYmdt"`
	UpdatedYmdt  time.Time `json:"updatedYmdt"`
}

type WatchdogsResponse struct {
	Header    *ResponseHeader `json:"header"`
	Watchdogs []Watchdog      `json:"watchdogs"`
}

type CreateWatchdogRequest struct {
	WatchdogName string `json:"watchdogName"`
	IsEnabled    bool   `json:"isEnabled"`
	QueryTimeout int    `json:"queryTimeout"`
}

// Monitoring & Events
type Metric struct {
	MetricID    string `json:"metricId"`
	MetricName  string `json:"metricName"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

type MetricsResponse struct {
	Header  *ResponseHeader `json:"header"`
	Metrics []Metric        `json:"metrics"`
}

type MetricStatistic struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

type MetricStatisticsResponse struct {
	Header     *ResponseHeader   `json:"header"`
	Statistics []MetricStatistic `json:"statistics"`
}

type Event struct {
	EventID      string    `json:"eventId"`
	EventCode    string    `json:"eventCode"`
	EventName    string    `json:"eventName"`
	Category     string    `json:"category"`
	DBInstanceID string    `json:"dbInstanceId,omitempty"`
	SourceType   string    `json:"sourceType"`
	SourceID     string    `json:"sourceId,omitempty"`
	Message      string    `json:"message"`
	EventYmdt    time.Time `json:"eventYmdt"`
}

type EventsResponse struct {
	Header *ResponseHeader `json:"header"`
	Events []Event         `json:"events"`
}

type EventParams struct {
	InstanceID string
	StartTime  string
	EndTime    string
	EventCode  string
	SourceType string
	Page       int
	Size       int
}

type MonitoringItem struct {
	MonitoringItemID   string `json:"monitoringItemId"`
	MonitoringItemName string `json:"monitoringItemName"`
	Description        string `json:"description"`
	Threshold          int    `json:"threshold"`
	IsEnabled          bool   `json:"isEnabled"`
}

type MonitoringItemsResponse struct {
	Header          *ResponseHeader  `json:"header"`
	MonitoringItems []MonitoringItem `json:"monitoringItems"`
}

// Job Response
type JobIDResponse struct {
	Header *ResponseHeader `json:"header"`
	JobID  string          `json:"jobId"`
}

// Job Details
type JobDetail struct {
	JobID        string    `json:"jobId"`
	JobType      string    `json:"jobType"`
	JobStatus    string    `json:"jobStatus"`
	ResourceID   string    `json:"resourceId"`
	ResourceName string    `json:"resourceName"`
	CreatedYmdt  time.Time `json:"createdYmdt"`
	UpdatedYmdt  time.Time `json:"updatedYmdt"`
}

type JobDetailResponse struct {
	Header *ResponseHeader `json:"header"`
	JobDetail
}

type NetworkEndpoint struct {
	Domain         string `json:"domain"`
	IPAddress      string `json:"ipAddress"`
	IsPublicAccess bool   `json:"isPublicAccess"`
}

type NetworkInfoResponse struct {
	Header    *ResponseHeader   `json:"header"`
	EndPoints []NetworkEndpoint `json:"endPoints"`
}

type ModifyDeletionProtectionRequest struct {
	UseDeletionProtection bool `json:"useDeletionProtection"`
}

type ModifyHighAvailabilityRequest struct {
	UseHighAvailability bool `json:"useHighAvailability"`
}

type ReplicaNetwork struct {
	AvailabilityZone string `json:"availabilityZone"`
}

type CreateReplicaRequest struct {
	DBInstanceName        string          `json:"dbInstanceName"`
	Description           string          `json:"description,omitempty"`
	DBFlavorID            string          `json:"dbFlavorId,omitempty"`
	DBPort                int             `json:"dbPort,omitempty"`
	ParameterGroupID      string          `json:"parameterGroupId,omitempty"`
	DBSecurityGroupIDs    []string        `json:"dbSecurityGroupIds,omitempty"`
	UserGroupIDs          []string        `json:"userGroupIds,omitempty"`
	UseDeletionProtection bool            `json:"useDeletionProtection,omitempty"`
	Network               *ReplicaNetwork `json:"network"`
}

type LogFile struct {
	LogFileName string    `json:"logFileName"`
	LogFileSize int64     `json:"logFileSize"`
	CreatedYmdt time.Time `json:"createdYmdt"`
	UpdatedYmdt time.Time `json:"updatedYmdt"`
}

type LogFilesResponse struct {
	Header   *ResponseHeader `json:"header"`
	LogFiles []LogFile       `json:"logFiles"`
}

type CreateInstanceGroupRequest struct {
	DBInstanceGroupName string `json:"dbInstanceGroupName"`
	ReplicationType     string `json:"replicationType,omitempty"`
}

type UpdateInstanceGroupRequest struct {
	DBInstanceGroupName string `json:"dbInstanceGroupName"`
}

type LogFileParams struct {
	InstanceID string
	Page       int
	Size       int
}

type Schema struct {
	SchemaName string   `json:"schemaName"`
	Owner      string   `json:"owner"`
	Tables     []string `json:"tables,omitempty"`
	Size       int64    `json:"size"`
}

type DatabaseSchemaResponse struct {
	Header  *ResponseHeader `json:"header"`
	Schemas []Schema        `json:"schemas"`
}

type Permission struct {
	DatabaseName string   `json:"databaseName"`
	Privileges   []string `json:"privileges"`
}

type GrantPermissionsRequest struct {
	Permissions []Permission `json:"permissions"`
}

type ListInstancesOutput = DBInstancesResponse
type GetInstanceOutput = DBInstanceResponse
type CreateInstanceInput = CreateDBInstanceRequest
type ModifyInstanceInput = ModifyDBInstanceRequest
type JobOutput = JobIDResponse
type ListInstanceGroupsOutput = DBInstanceGroupsResponse
type InstanceGroupOutput = DBInstanceGroupResponse
type ListFlavorsOutput = DBFlavorsResponse
type ListVersionsOutput = DBVersionsResponse
type ListStorageTypesOutput = StorageTypesResponse
type ListSecurityGroupsOutput = DBSecurityGroupsResponse
type SecurityGroupOutput = DBSecurityGroupResponse
type CreateSecurityGroupInput = CreateDBSecurityGroupRequest
type SecurityGroupIDOutput = JobIDResponse
type CreateSecurityGroupRuleInput = CreateDBSecurityGroupRuleRequest
type SecurityGroupRuleOutput = JobIDResponse
type ListParameterGroupsOutput = ParameterGroupsResponse
type ParameterGroupOutput = ParameterGroupResponse
type CreateParameterGroupInput = CreateParameterGroupRequest
type ModifyParametersInput = ModifyParametersRequest
type ParameterGroupIDOutput = JobIDResponse
type ListBackupsOutput = BackupsResponse
type CreateBackupInput = CreateBackupRequest
type RestoreBackupInput = RestoreBackupRequest
type ListDBUsersOutput = DBUsersResponse
type CreateDBUserInput = CreateDBUserRequest
type UpdateDBUserInput = CreateDBUserRequest
type ListDatabasesOutput = DatabasesResponse
type CreateDatabaseInput = CreateDatabaseRequest
type DatabaseIDOutput = JobIDResponse
type ListSubnetsOutput = SubnetsResponse
type NetworkInfoOutput = NetworkInfoResponse
type ModifyStorageInfoInput = ModifyStorageInfoRequest
type ModifyDeletionProtectionInput = ModifyDeletionProtectionRequest
type CreateReplicaInput = CreateReplicaRequest
type EnableHAInput = ModifyHighAvailabilityRequest
type ListNotificationGroupsOutput = NotificationGroupsResponse
type NotificationGroupOutput = NotificationGroupResponse
type CreateNotificationGroupInput = CreateNotificationGroupRequest
type NotificationGroupIDOutput = JobIDResponse
type ListLogFilesOutput = LogFilesResponse
