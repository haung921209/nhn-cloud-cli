package mariadb

// ResponseHeader represents common API response header
type ResponseHeader struct {
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
	IsSuccessful  bool   `json:"isSuccessful"`
}

// DatabaseInstanceGroup represents a MariaDB database instance group
type DatabaseInstanceGroup struct {
	DBInstanceGroupID string `json:"dbInstanceGroupId"`
	ReplicationType   string `json:"replicationType"`
	CreatedYmdt       string `json:"createdYmdt"`
	UpdatedYmdt       string `json:"updatedYmdt"`
}

// DatabaseInstanceGroupsResponse represents the response for listing instance groups
type DatabaseInstanceGroupsResponse struct {
	Header           *ResponseHeader         `json:"header"`
	DBInstanceGroups []DatabaseInstanceGroup `json:"dbInstanceGroups"`
}

// DatabaseInstance represents a MariaDB database instance
type DatabaseInstance struct {
	DBInstanceID     string `json:"dbInstanceId"`
	DBInstanceName   string `json:"dbInstanceName"`
	DBInstanceStatus string `json:"dbInstanceStatus"`
	Description      string `json:"description,omitempty"`
	DBVersion        string `json:"dbVersion"`
	DBPort           int    `json:"dbPort"`

	// Storage
	StorageType string `json:"storageType"`
	StorageSize int    `json:"storageSize"`

	// Network
	SubnetID           string   `json:"subnetId"`
	DBSecurityGroupIDs []string `json:"dbSecurityGroupIds,omitempty"`

	// Configuration
	DBFlavorID       string `json:"dbFlavorId"`
	ParameterGroupID string `json:"parameterGroupId,omitempty"`

	// MariaDB specific
	AuthenticationPlugin string `json:"authenticationPlugin,omitempty"`

	// Protection
	UseDeletionProtection bool `json:"useDeletionProtection,omitempty"`

	// Timestamps
	CreatedYmdt string `json:"createdYmdt"`
	UpdatedYmdt string `json:"updatedYmdt"`
}

// DatabaseInstanceResponse represents the response for database instance operations
type DatabaseInstanceResponse struct {
	Header *ResponseHeader `json:"header"`
	DatabaseInstance
}

// DatabaseInstancesResponse represents the response for listing database instances
type DatabaseInstancesResponse struct {
	Header      *ResponseHeader    `json:"header"`
	DBInstances []DatabaseInstance `json:"dbInstances"`
}

// Network represents network configuration for database instance
type Network struct {
	SubnetID         string `json:"subnetId"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	UsePublicAccess  bool   `json:"usePublicAccess,omitempty"`
}

// Storage represents storage configuration for database instance
type Storage struct {
	StorageType string `json:"storageType"`
	StorageSize int    `json:"storageSize"`
}

// BackupSchedule represents a backup schedule configuration
type BackupSchedule struct {
	BackupWndBgnTime  string `json:"backupWndBgnTime"`
	BackupWndDuration string `json:"backupWndDuration"`
}

// BackupConfig represents backup configuration for database instance
type BackupConfig struct {
	BackupPeriod    int              `json:"backupPeriod"`
	BackupSchedules []BackupSchedule `json:"backupSchedules"`
}

// RestartInstanceRequest represents a request to restart a database instance
type RestartInstanceRequest struct {
	// UseOnlineFailover enables restart using failover (HA instances only)
	// When true, minimizes downtime by failing over to standby before restart
	UseOnlineFailover bool `json:"useOnlineFailover,omitempty"`
	// ExecuteBackup triggers a backup before restart
	ExecuteBackup bool `json:"executeBackup,omitempty"`
}

// CreateDatabaseInstanceRequest represents MariaDB instance creation request
type CreateDatabaseInstanceRequest struct {
	DBInstanceName          string `json:"dbInstanceName"`
	DBInstanceCandidateName string `json:"dbInstanceCandidateName,omitempty"`
	Description             string `json:"description,omitempty"`
	DBFlavorID              string `json:"dbFlavorId"`
	DBVersion               string `json:"dbVersion"`
	DBUserName              string `json:"dbUserName"`
	DBPassword              string `json:"dbPassword"`
	DBPort                  int    `json:"dbPort,omitempty"`

	// Configuration
	ParameterGroupID     string   `json:"parameterGroupId,omitempty"`
	DBSecurityGroupIDs   []string `json:"dbSecurityGroupIds,omitempty"`
	UserGroupIDs         []string `json:"userGroupIds,omitempty"`
	NotificationGroupIDs []string `json:"notificationGroupIds,omitempty"`

	// Nested objects (required by API)
	Network *Network      `json:"network"`
	Storage *Storage      `json:"storage"`
	Backup  *BackupConfig `json:"backup"`

	// HA and other options
	UseHighAvailability   bool   `json:"useHighAvailability,omitempty"`
	ReplicationMode       string `json:"replicationMode,omitempty"`
	UseDeletionProtection bool   `json:"useDeletionProtection,omitempty"`

	// MariaDB specific
	AuthenticationPlugin string `json:"authenticationPlugin,omitempty"`
}

// ModifyDatabaseInstanceRequest represents a request to modify a database instance
// API: PUT /v3.0/db-instances/{dbInstanceId}
type ModifyDatabaseInstanceRequest struct {
	DBInstanceName          string   `json:"dbInstanceName,omitempty"`
	DBInstanceCandidateName string   `json:"dbInstanceCandidateName,omitempty"`
	Description             string   `json:"description,omitempty"`
	DBPort                  int      `json:"dbPort,omitempty"`
	DBFlavorID              string   `json:"dbFlavorId,omitempty"`
	ParameterGroupID        string   `json:"parameterGroupId,omitempty"`
	DBSecurityGroupIDs      []string `json:"dbSecurityGroupIds,omitempty"`
	ExecuteBackup           bool     `json:"executeBackup,omitempty"`
	UseOnlineFailover       bool     `json:"useOnlineFailover,omitempty"`
}

// ModifyStorageInfoRequest for PUT /v3.0/db-instances/{dbInstanceId}/storage-info
type ModifyStorageInfoRequest struct {
	StorageSize       int  `json:"storageSize"`
	UseOnlineFailover bool `json:"useOnlineFailover,omitempty"`
}

// ModifyDeletionProtectionRequest for PUT /v3.0/db-instances/{dbInstanceId}/deletion-protection
type ModifyDeletionProtectionRequest struct {
	UseDeletionProtection bool `json:"useDeletionProtection"`
}

// Parameter Groups
type Parameter struct {
	ParameterID       string `json:"parameterId"`
	ParameterName     string `json:"parameterName"`
	FileParameterName string `json:"fileParameterName"`
	Value             string `json:"value"`
	DefaultValue      string `json:"defaultValue"`
	AllowedValue      string `json:"allowedValue"`
	UpdateType        string `json:"updateType"`
	ApplyType         string `json:"applyType"`
}

type ParameterGroup struct {
	ParameterGroupID     string      `json:"parameterGroupId"`
	ParameterGroupName   string      `json:"parameterGroupName"`
	Description          string      `json:"description,omitempty"`
	DBVersion            string      `json:"dbVersion"`
	ParameterGroupStatus string      `json:"parameterGroupStatus"`
	Parameters           []Parameter `json:"parameters,omitempty"`
	CreatedYmdt          string      `json:"createdYmdt"`
	UpdatedYmdt          string      `json:"updatedYmdt"`
}

type ParameterGroupsResponse struct {
	Header          *ResponseHeader  `json:"header"`
	ParameterGroups []ParameterGroup `json:"parameterGroups"`
}

type ParameterGroupResponse struct {
	Header         *ResponseHeader `json:"header"`
	ParameterGroup                 // Embedded fields
}

type CreateParameterGroupRequest struct {
	ParameterGroupName string `json:"parameterGroupName"`
	Description        string `json:"description,omitempty"`
	DBVersion          string `json:"dbVersion"`
}

type CopyParameterGroupRequest struct {
	ParameterGroupName string `json:"parameterGroupName"`
	Description        string `json:"description,omitempty"`
}

type UpdateParameterGroupRequest struct {
	ParameterGroupName string `json:"parameterGroupName,omitempty"`
	Description        string `json:"description,omitempty"`
}

type ModifyParametersRequest struct {
	ModifiedParameters []struct {
		ParameterID string `json:"parameterId"`
		Value       string `json:"value"`
	} `json:"modifiedParameters"`
}

type ParameterGroupIDResponse struct {
	Header           *ResponseHeader `json:"header"`
	ParameterGroupID string          `json:"parameterGroupId"`
}

// DB Security Groups
type Port struct {
	PortType string `json:"portType"`
	MinPort  *int   `json:"minPort,omitempty"`
	MaxPort  *int   `json:"maxPort,omitempty"`
}

type SecurityRule struct {
	RuleID      string `json:"ruleId"`
	Description string `json:"description,omitempty"`
	Direction   string `json:"direction"`
	EtherType   string `json:"etherType"`
	Port        Port   `json:"port"`
	CIDR        string `json:"cidr"`
	CreatedYmdt string `json:"createdYmdt"`
	UpdatedYmdt string `json:"updatedYmdt"`
}

type DBSecurityGroup struct {
	DBSecurityGroupID   string         `json:"dbSecurityGroupId"`
	DBSecurityGroupName string         `json:"dbSecurityGroupName"`
	Description         string         `json:"description,omitempty"`
	ProgressStatus      string         `json:"progressStatus"`
	Rules               []SecurityRule `json:"rules,omitempty"`
	CreatedYmdt         string         `json:"createdYmdt"`
	UpdatedYmdt         string         `json:"updatedYmdt"`
}

type DBSecurityGroupsResponse struct {
	Header           *ResponseHeader   `json:"header"`
	DBSecurityGroups []DBSecurityGroup `json:"dbSecurityGroups"`
}

type DBSecurityGroupResponse struct {
	Header          *ResponseHeader `json:"header"`
	DBSecurityGroup DBSecurityGroup `json:"dbSecurityGroup"`
}

type CreateDBSecurityGroupRequest struct {
	DBSecurityGroupName string `json:"dbSecurityGroupName"`
	Description         string `json:"description,omitempty"`
	Rules               []struct {
		Description string `json:"description,omitempty"`
		Direction   string `json:"direction"`
		EtherType   string `json:"etherType"`
		Port        Port   `json:"port"`
		CIDR        string `json:"cidr"`
	} `json:"rules"`
}

type UpdateDBSecurityGroupRequest struct {
	DBSecurityGroupName string `json:"dbSecurityGroupName,omitempty"`
	Description         string `json:"description,omitempty"`
}

type DBSecurityGroupIDResponse struct {
	Header            *ResponseHeader `json:"header"`
	DBSecurityGroupID string          `json:"dbSecurityGroupId"`
}

type CreateDBSecurityGroupRuleRequest struct {
	Description string `json:"description,omitempty"`
	Direction   string `json:"direction"`
	EtherType   string `json:"etherType"`
	Port        Port   `json:"port"`
	CIDR        string `json:"cidr"`
}

type CreateDBSecurityGroupRuleResponse struct {
	Header *ResponseHeader `json:"header"`
	RuleID string          `json:"ruleId"`
}

type UpdateDBSecurityGroupRuleRequest struct {
	Description string `json:"description,omitempty"`
	Direction   string `json:"direction,omitempty"`
	EtherType   string `json:"etherType,omitempty"`
	Port        *Port  `json:"port,omitempty"`
	CIDR        string `json:"cidr,omitempty"`
}

// DB Users
type DBUser struct {
	DBUserID             string `json:"dbUserId"`
	DBUserName           string `json:"dbUserName"`
	Host                 string `json:"host"`
	AuthorityType        string `json:"authorityType"`
	DBUserStatus         string `json:"dbUserStatus"`
	AuthenticationPlugin string `json:"authenticationPlugin,omitempty"`
	CreatedYmdt          string `json:"createdYmdt"`
	UpdatedYmdt          string `json:"updatedYmdt"`
}

type DBUsersResponse struct {
	Header  *ResponseHeader `json:"header"`
	DBUsers []DBUser        `json:"dbUsers"`
}

type CreateDBUserRequest struct {
	DBUserName           string `json:"dbUserName"`
	DBPassword           string `json:"dbPassword"`
	Host                 string `json:"host"`
	AuthorityType        string `json:"authorityType"`
	AuthenticationPlugin string `json:"authenticationPlugin,omitempty"`
}

type UpdateDBUserRequest struct {
	AuthorityType string `json:"authorityType,omitempty"`
}

type JobIDResponse struct {
	Header *ResponseHeader `json:"header"`
	JobID  string          `json:"jobId"`
}

// NetworkSubnet represents a RDS-specific network subnet
type NetworkSubnet struct {
	SubnetID         string `json:"subnetId"`
	SubnetName       string `json:"subnetName"`
	SubnetCidr       string `json:"subnetCidr"`
	UsingGateway     bool   `json:"usingGateway"`
	AvailableIpCount int    `json:"availableIpCount"`
}

// NetworkSubnetsResponse represents the response for RDS network subnets
type NetworkSubnetsResponse struct {
	Header  *ResponseHeader `json:"header"`
	Subnets []NetworkSubnet `json:"subnets"`
}

// Backups
type Backup struct {
	BackupID     string `json:"backupId"`
	BackupName   string `json:"backupName"`
	BackupStatus string `json:"backupStatus"`
	DBInstanceID string `json:"dbInstanceId"`
	DBVersion    string `json:"dbVersion"`
	BackupType   string `json:"backupType"`
	BackupSize   int64  `json:"backupSize"`
	CreatedYmdt  string `json:"createdYmdt"`
	UpdatedYmdt  string `json:"updatedYmdt"`
}

type BackupsResponse struct {
	Header      *ResponseHeader `json:"header"`
	TotalCounts int             `json:"totalCounts"`
	Backups     []Backup        `json:"backups"`
}

type CreateBackupRequest struct {
	BackupName string `json:"backupName"`
}

type BackupToObjectStorageRequest struct {
	TenantID        string `json:"tenantId"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	TargetContainer string `json:"targetContainer"`
	ObjectPath      string `json:"objectPath"`
}

type RestoreBackupRequest struct {
	DBInstanceName        string   `json:"dbInstanceName"`
	DBFlavorID            string   `json:"dbFlavorId"`
	DBPort                int      `json:"dbPort,omitempty"`
	ParameterGroupID      string   `json:"parameterGroupId,omitempty"`
	DBSecurityGroupIDs    []string `json:"dbSecurityGroupIds,omitempty"`
	UserGroupIDs          []string `json:"userGroupIds,omitempty"`
	UsePublicAccess       bool     `json:"usePublicAccess,omitempty"`
	AvailabilityZone      string   `json:"availabilityZone,omitempty"`
	RestoreBackupID       string   `json:"restoreBackupId"`
	RestoreBinLogFileName string   `json:"restoreBinLogFileName,omitempty"`
	RestoreBinLogPosition int64    `json:"restoreBinLogPosition,omitempty"`
}

type ExportBackupRequest struct {
	TenantID        string `json:"tenantId"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	TargetContainer string `json:"targetContainer"`
	ObjectPath      string `json:"objectPath"`
}

// Flavors & Versions
type DBFlavor struct {
	FlavorID   string `json:"dbFlavorId"`
	FlavorName string `json:"dbFlavorName"`
	Ram        int    `json:"ram"`
	Vcpus      int    `json:"vcpus"`
	Disk       int    `json:"disk"`
}

type DBFlavorsResponse struct {
	Header    *ResponseHeader `json:"header"`
	DBFlavors []DBFlavor      `json:"dbFlavors"`
}

type DBVersion struct {
	DBVersion     string `json:"dbVersion"`
	DBVersionName string `json:"dbVersionName"`
}

type DBVersionsResponse struct {
	Header     *ResponseHeader `json:"header"`
	DBVersions []DBVersion     `json:"dbVersions"`
}

// Storage Types
type StorageTypesResponse struct {
	Header       *ResponseHeader `json:"header"`
	StorageTypes []string        `json:"storageTypes"`
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
	} `json:"recipients,omitempty"`
	CreatedYmdt string `json:"createdYmdt"`
	UpdatedYmdt string `json:"updatedYmdt"`
}

type NotificationGroupsResponse struct {
	Header             *ResponseHeader     `json:"header"`
	NotificationGroups []NotificationGroup `json:"notificationGroups"`
}

type NotificationGroupResponse struct {
	Header            *ResponseHeader   `json:"header"`
	NotificationGroup NotificationGroup `json:"notificationGroup"`
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

type UpdateNotificationGroupRequest struct {
	NotificationGroupName string `json:"notificationGroupName,omitempty"`
	Description           string `json:"description,omitempty"`
	NotificationType      string `json:"notificationType,omitempty"`
	IsEnabled             bool   `json:"isEnabled"`
	Recipients            []struct {
		RecipientType string `json:"recipientType"`
		Recipient     string `json:"recipient"`
	} `json:"recipients,omitempty"`
}

type NotificationGroupIDResponse struct {
	Header              *ResponseHeader `json:"header"`
	NotificationGroupID string          `json:"notificationGroupId"`
}

// Log Files
type LogFile struct {
	LogFileName string `json:"logFileName"`
	LogFileSize int64  `json:"logFileSize"`
	CreatedYmdt string `json:"createdYmdt"`
	UpdatedYmdt string `json:"updatedYmdt"`
}

type LogFilesResponse struct {
	Header   *ResponseHeader `json:"header"`
	LogFiles []LogFile       `json:"logFiles"`
}

// Network/Subnets
type Subnet struct {
	SubnetID         string `json:"subnetId"`
	SubnetName       string `json:"subnetName"`
	SubnetCidr       string `json:"subnetCidr"`
	UsingGateway     bool   `json:"usingGateway"`
	AvailableIpCount int    `json:"availableIpCount"`
}

type SubnetsResponse struct {
	Header  *ResponseHeader `json:"header"`
	Subnets []Subnet        `json:"subnets"`
}

// DatabaseInstanceInGroup represents an individual instance within a group
type DatabaseInstanceInGroup struct {
	DBInstanceID     string `json:"dbInstanceId"`
	DBInstanceType   string `json:"dbInstanceType"`
	DBInstanceStatus string `json:"dbInstanceStatus"`
}

// DatabaseInstanceGroupDetail represents detailed information about an instance group
type DatabaseInstanceGroupDetail struct {
	Header            *ResponseHeader           `json:"header"`
	DBInstanceGroupID string                    `json:"dbInstanceGroupId"`
	ReplicationType   string                    `json:"replicationType"`
	DBInstances       []DatabaseInstanceInGroup `json:"dbInstances"`
	CreatedYmdt       string                    `json:"createdYmdt"`
	UpdatedYmdt       string                    `json:"updatedYmdt"`
}

// DatabaseInstanceGroupResponse represents the response for getting instance group details
type DatabaseInstanceGroupResponse = DatabaseInstanceGroupDetail

// NetworkEndpoint represents a network endpoint for database connections
type NetworkEndpoint struct {
	Domain       string `json:"domain"`
	IPAddress    string `json:"ipAddress"`
	EndPointType string `json:"endPointType"`
}

// SubnetInfo represents subnet information
type SubnetInfo struct {
	SubnetID   string `json:"subnetId"`
	SubnetName string `json:"subnetName"`
	SubnetCidr string `json:"subnetCidr"`
}

// NetworkInfoResponse represents the response for instance network information
type NetworkInfoResponse struct {
	Header           *ResponseHeader   `json:"header"`
	AvailabilityZone string            `json:"availabilityZone"`
	Subnet           SubnetInfo        `json:"subnet"`
	EndPoints        []NetworkEndpoint `json:"endPoints"`
}

// DB Schemas (Database Management)
type DBSchema struct {
	DBSchemaID   string `json:"dbSchemaId"`
	DBSchemaName string `json:"dbSchemaName"`
	CreatedYmdt  string `json:"createdYmdt"`
	UpdatedYmdt  string `json:"updatedYmdt"`
}

type DBSchemasResponse struct {
	Header    *ResponseHeader `json:"header"`
	DBSchemas []DBSchema      `json:"dbSchemas"`
}

type DBSchemaResponse struct {
	Header   *ResponseHeader `json:"header"`
	DBSchema DBSchema        `json:"dbSchema"`
}

type CreateDBSchemaRequest struct {
	DBSchemaName string `json:"dbSchemaName"`
}

type DBSchemaIDResponse struct {
	Header     *ResponseHeader `json:"header"`
	DBSchemaID string          `json:"dbSchemaId"`
}

// Replica Types

// CreateReplicaRequest represents a request to create a read replica
type CreateReplicaRequest struct {
	DBInstanceName          string          `json:"dbInstanceName"`
	DBInstanceCandidateName string          `json:"dbInstanceCandidateName,omitempty"`
	Description             string          `json:"description,omitempty"`
	DBFlavorID              string          `json:"dbFlavorId,omitempty"`
	DBPort                  int             `json:"dbPort,omitempty"`
	ParameterGroupID        string          `json:"parameterGroupId,omitempty"`
	DBSecurityGroupIDs      []string        `json:"dbSecurityGroupIds,omitempty"`
	UserGroupIDs            []string        `json:"userGroupIds,omitempty"`
	UseDefaultNotification  bool            `json:"useDefaultNotification,omitempty"`
	UseDeletionProtection   bool            `json:"useDeletionProtection,omitempty"`
	Network                 *ReplicaNetwork `json:"network"`
	Storage                 *ReplicaStorage `json:"storage,omitempty"`
	Backup                  *ReplicaBackup  `json:"backup,omitempty"`
}

// ReplicaNetwork represents network configuration for replica creation
type ReplicaNetwork struct {
	AvailabilityZone string `json:"availabilityZone"`
	UsePublicAccess  bool   `json:"usePublicAccess,omitempty"`
}

// ReplicaStorage represents storage configuration for replica creation
type ReplicaStorage struct {
	StorageType string `json:"storageType,omitempty"`
	StorageSize int    `json:"storageSize,omitempty"`
}

// ReplicaBackup represents backup configuration for replica creation
type ReplicaBackup struct {
	BackupPeriod    int              `json:"backupPeriod,omitempty"`
	BackupSchedules []BackupSchedule `json:"backupSchedules,omitempty"`
}

type Metric struct {
	MeasureName string `json:"measureName"`
	Unit        string `json:"unit"`
}

type MetricsResponse struct {
	Header  *ResponseHeader `json:"header"`
	Metrics []Metric        `json:"metrics"`
}

type MetricStatistic struct {
	MeasureName string          `json:"measureName"`
	Unit        string          `json:"unit"`
	Values      [][]interface{} `json:"values"`
}

type MetricStatisticsResponse struct {
	Header           *ResponseHeader   `json:"header"`
	MetricStatistics []MetricStatistic `json:"metricStatistics"`
}

type ModifyHighAvailabilityRequest struct {
	UseHighAvailability     bool `json:"useHighAvailability"`
	PingInterval            int  `json:"pingInterval,omitempty"`
	FailoverReplWaitingTime int  `json:"failoverReplWaitingTime,omitempty"`
}

type ModifyNetworkInfoRequest struct {
	UsePublicAccess bool `json:"usePublicAccess"`
}

type ListInstancesOutput = DatabaseInstancesResponse
type GetInstanceOutput = DatabaseInstanceResponse
type CreateInstanceInput = CreateDatabaseInstanceRequest
type CreateInstanceOutput = JobIDResponse
type ModifyInstanceInput = ModifyDatabaseInstanceRequest
type JobOutput = JobIDResponse
type ListInstanceGroupsOutput = DatabaseInstanceGroupsResponse
type InstanceGroupOutput = DatabaseInstanceGroupResponse
type ListFlavorsOutput = DBFlavorsResponse
type ListVersionsOutput = DBVersionsResponse
type ListStorageTypesOutput = StorageTypesResponse
type ListSecurityGroupsOutput = DBSecurityGroupsResponse
type SecurityGroupOutput = DBSecurityGroupResponse
type CreateSecurityGroupInput = CreateDBSecurityGroupRequest
type UpdateSecurityGroupInput = UpdateDBSecurityGroupRequest
type SecurityGroupIDOutput = DBSecurityGroupIDResponse
type CreateSecurityGroupRuleInput = CreateDBSecurityGroupRuleRequest
type UpdateSecurityGroupRuleInput = UpdateDBSecurityGroupRuleRequest
type SecurityGroupRuleOutput = CreateDBSecurityGroupRuleResponse
type ListParameterGroupsOutput = ParameterGroupsResponse
type ParameterGroupOutput = ParameterGroupResponse
type CreateParameterGroupInput = CreateParameterGroupRequest
type CopyParameterGroupInput = CopyParameterGroupRequest
type UpdateParameterGroupInput = UpdateParameterGroupRequest
type ModifyParametersInput = ModifyParametersRequest
type ParameterGroupIDOutput = ParameterGroupIDResponse
type ListBackupsOutput = BackupsResponse
type CreateBackupInput = CreateBackupRequest
type BackupToObjectStorageInput = BackupToObjectStorageRequest
type RestoreBackupInput = RestoreBackupRequest
type ExportBackupInput = ExportBackupRequest
type ListDBUsersOutput = DBUsersResponse
type CreateDBUserInput = CreateDBUserRequest
type UpdateDBUserInput = UpdateDBUserRequest
type ListNotificationGroupsOutput = NotificationGroupsResponse
type NotificationGroupOutput = NotificationGroupResponse
type CreateNotificationGroupInput = CreateNotificationGroupRequest
type UpdateNotificationGroupInput = UpdateNotificationGroupRequest
type NotificationGroupIDOutput = NotificationGroupIDResponse
type ListLogFilesOutput = LogFilesResponse
type ListSchemasOutput = DBSchemasResponse
type SchemaOutput = DBSchemaResponse
type CreateSchemaInput = CreateDBSchemaRequest
type SchemaIDOutput = DBSchemaIDResponse
type ListMetricsOutput = MetricsResponse
type MetricStatisticsOutput = MetricStatisticsResponse
type ListSubnetsOutput = SubnetsResponse
type NetworkInfoOutput = NetworkInfoResponse
type EnableHAInput = ModifyHighAvailabilityRequest
type ModifyStorageInfoInput = ModifyStorageInfoRequest
type ModifyDeletionProtectionInput = ModifyDeletionProtectionRequest
type ModifyNetworkInfoInput = ModifyNetworkInfoRequest
type CreateReplicaInput = CreateReplicaRequest
