package postgresql

import (
	"context"
	"fmt"
	"net/url"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/endpoint"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/transport"
)

type Client struct {
	transport *transport.Client
	region    string
	appKey    string
}

func NewClient(region, appKey string, creds credentials.Credentials, debug bool) *Client {
	baseURL := endpoint.ResolveWithAppKey(endpoint.ServiceRDSPostgreSQL, region, appKey)

	opts := []transport.ClientOption{
		transport.WithDebug(debug),
	}

	if creds != nil {
		// PostgreSQL v1.0 API uses OAuth2 Bearer token (X-NHN-AUTHORIZATION)
		// unlike MySQL/MariaDB v3.0 which uses X-TC-AUTHENTICATION-ID/SECRET
		tokenProvider := credentials.NewTokenProvider(
			creds.GetAccessKeyID(),
			creds.GetSecretAccessKey(),
		)
		opts = append(opts, transport.WithDynamicBearerAuth(appKey, tokenProvider))
	}

	return &Client{
		transport: transport.NewClient(baseURL, opts...),
		region:    region,
		appKey:    appKey,
	}
}

// Instance operations

func (c *Client) ListInstances(ctx context.Context) (*ListInstancesOutput, error) {
	var out ListInstancesOutput
	if err := c.transport.GET(ctx, "/db-instances", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetInstance(ctx context.Context, instanceID string) (*GetInstanceOutput, error) {
	var out GetInstanceOutput
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateInstance(ctx context.Context, input *CreateInstanceInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ModifyInstance(ctx context.Context, instanceID string, input *ModifyInstanceInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.PUT(ctx, "/db-instances/"+instanceID, input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteInstance(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.DELETE(ctx, "/db-instances/"+instanceID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) StartInstance(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/start", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) StopInstance(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/stop", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) RestartInstance(ctx context.Context, instanceID string, req *RestartInstanceRequest) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/restart", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ForceRestartInstance(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/force-restart", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Instance Groups

func (c *Client) ListInstanceGroups(ctx context.Context) (*ListInstanceGroupsOutput, error) {
	var out ListInstanceGroupsOutput
	if err := c.transport.GET(ctx, "/db-instance-groups", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetInstanceGroup(ctx context.Context, groupID string) (*InstanceGroupOutput, error) {
	var out InstanceGroupOutput
	if err := c.transport.GET(ctx, "/db-instance-groups/"+groupID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateInstanceGroup(ctx context.Context, req *CreateInstanceGroupRequest) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instance-groups", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateInstanceGroup(ctx context.Context, groupID string, req *UpdateInstanceGroupRequest) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.PUT(ctx, "/db-instance-groups/"+groupID, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteInstanceGroup(ctx context.Context, groupID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.DELETE(ctx, "/db-instance-groups/"+groupID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Flavors, Versions, Storage Types

func (c *Client) ListFlavors(ctx context.Context) (*ListFlavorsOutput, error) {
	var out ListFlavorsOutput
	if err := c.transport.GET(ctx, "/db-flavors", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListVersions(ctx context.Context) (*ListVersionsOutput, error) {
	var out ListVersionsOutput
	if err := c.transport.GET(ctx, "/db-versions", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListStorageTypes(ctx context.Context) (*ListStorageTypesOutput, error) {
	var out ListStorageTypesOutput
	if err := c.transport.GET(ctx, "/storage-types", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Security Groups

func (c *Client) ListSecurityGroups(ctx context.Context) (*ListSecurityGroupsOutput, error) {
	var out ListSecurityGroupsOutput
	if err := c.transport.GET(ctx, "/db-security-groups", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetSecurityGroup(ctx context.Context, securityGroupID string) (*SecurityGroupOutput, error) {
	var out SecurityGroupOutput
	if err := c.transport.GET(ctx, "/db-security-groups/"+securityGroupID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateSecurityGroup(ctx context.Context, input *CreateSecurityGroupInput) (*SecurityGroupIDOutput, error) {
	var out SecurityGroupIDOutput
	if err := c.transport.POST(ctx, "/db-security-groups", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteSecurityGroup(ctx context.Context, securityGroupID string) (*ResponseHeader, error) {
	var out ResponseHeader
	if err := c.transport.DELETE(ctx, "/db-security-groups/"+securityGroupID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateSecurityGroupRule(ctx context.Context, securityGroupID string, input *CreateSecurityGroupRuleInput) (*SecurityGroupRuleOutput, error) {
	var out SecurityGroupRuleOutput
	if err := c.transport.POST(ctx, "/db-security-groups/"+securityGroupID+"/rules", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteSecurityGroupRule(ctx context.Context, securityGroupID, ruleID string) (*ResponseHeader, error) {
	var out ResponseHeader
	if err := c.transport.DELETE(ctx, "/db-security-groups/"+securityGroupID+"/rules/"+ruleID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListParameterGroups(ctx context.Context) (*ListParameterGroupsOutput, error) {
	var out ListParameterGroupsOutput
	if err := c.transport.GET(ctx, "/parameter-groups", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetParameterGroup(ctx context.Context, parameterGroupID string) (*ParameterGroupOutput, error) {
	var out ParameterGroupOutput
	if err := c.transport.GET(ctx, "/parameter-groups/"+parameterGroupID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateParameterGroup(ctx context.Context, input *CreateParameterGroupInput) (*ParameterGroupIDOutput, error) {
	var out ParameterGroupIDOutput
	if err := c.transport.POST(ctx, "/parameter-groups", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ModifyParameters(ctx context.Context, parameterGroupID string, input *ModifyParametersInput) (*ResponseHeader, error) {
	var out ResponseHeader
	if err := c.transport.PUT(ctx, "/parameter-groups/"+parameterGroupID+"/parameters", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteParameterGroup(ctx context.Context, parameterGroupID string) (*ResponseHeader, error) {
	var out ResponseHeader
	if err := c.transport.DELETE(ctx, "/parameter-groups/"+parameterGroupID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Backups

func (c *Client) ListBackups(ctx context.Context, instanceID string, page, size int) (*ListBackupsOutput, error) {
	path := fmt.Sprintf("/backups?page=%d&size=%d", page, size)
	if instanceID != "" {
		path += "&dbInstanceId=" + url.QueryEscape(instanceID)
	}
	var out ListBackupsOutput
	if err := c.transport.GET(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateBackup(ctx context.Context, instanceID string, input *CreateBackupInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/backup", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) RestoreBackup(ctx context.Context, backupID string, input *RestoreBackupInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/backups/"+backupID+"/restore", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteBackup(ctx context.Context, backupID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.DELETE(ctx, "/backups/"+backupID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DB Users

func (c *Client) ListDBUsers(ctx context.Context, instanceID string) (*ListDBUsersOutput, error) {
	var out ListDBUsersOutput
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID+"/db-users", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateDBUser(ctx context.Context, instanceID string, input *CreateDBUserInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/db-users", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateDBUser(ctx context.Context, instanceID, userID string, input *UpdateDBUserInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.PUT(ctx, "/db-instances/"+instanceID+"/db-users/"+userID, input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteDBUser(ctx context.Context, instanceID, userID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.DELETE(ctx, "/db-instances/"+instanceID+"/db-users/"+userID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Databases

func (c *Client) ListDatabases(ctx context.Context, instanceID string) (*ListDatabasesOutput, error) {
	var out ListDatabasesOutput
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID+"/databases", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateDatabase(ctx context.Context, instanceID string, input *CreateDatabaseInput) (*DatabaseIDOutput, error) {
	var out DatabaseIDOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/databases", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteDatabase(ctx context.Context, instanceID, databaseID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.DELETE(ctx, "/db-instances/"+instanceID+"/databases/"+databaseID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateDatabase(ctx context.Context, instanceID, databaseID string, newOwner string) (*JobOutput, error) {
	req := map[string]string{"owner": newOwner}
	var out JobOutput
	if err := c.transport.PUT(ctx, "/db-instances/"+instanceID+"/databases/"+databaseID, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Network

func (c *Client) ListSubnets(ctx context.Context) (*ListSubnetsOutput, error) {
	var out ListSubnetsOutput
	if err := c.transport.GET(ctx, "/network/subnets", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetNetworkInfo(ctx context.Context, instanceID string) (*NetworkInfoOutput, error) {
	var out NetworkInfoOutput
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID+"/network-info", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ModifyStorageInfo(ctx context.Context, instanceID string, input *ModifyStorageInfoInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.PUT(ctx, "/db-instances/"+instanceID+"/storage-info", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ModifyDeletionProtection(ctx context.Context, instanceID string, input *ModifyDeletionProtectionInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.PUT(ctx, "/db-instances/"+instanceID+"/deletion-protection", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Replicas

func (c *Client) CreateReplica(ctx context.Context, instanceID string, input *CreateReplicaInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/replicate", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PromoteReplica(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/promote", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// High Availability

func (c *Client) EnableHighAvailability(ctx context.Context, instanceID string, input *EnableHAInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.PUT(ctx, "/db-instances/"+instanceID+"/high-availability", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DisableHighAvailability(ctx context.Context, instanceID string) (*JobOutput, error) {
	input := map[string]interface{}{"useHighAvailability": false}
	var out JobOutput
	if err := c.transport.PUT(ctx, "/db-instances/"+instanceID+"/high-availability", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PauseHighAvailability(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/high-availability/pause", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ResumeHighAvailability(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/high-availability/resume", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) RepairHighAvailability(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/high-availability/repair", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Notification Groups

func (c *Client) ListNotificationGroups(ctx context.Context) (*ListNotificationGroupsOutput, error) {
	var out ListNotificationGroupsOutput
	if err := c.transport.GET(ctx, "/notification-groups", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetNotificationGroup(ctx context.Context, notificationGroupID string) (*NotificationGroupOutput, error) {
	var out NotificationGroupOutput
	if err := c.transport.GET(ctx, "/notification-groups/"+notificationGroupID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateNotificationGroup(ctx context.Context, input *CreateNotificationGroupInput) (*NotificationGroupIDOutput, error) {
	var out NotificationGroupIDOutput
	if err := c.transport.POST(ctx, "/notification-groups", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteNotificationGroup(ctx context.Context, notificationGroupID string) (*ResponseHeader, error) {
	var out ResponseHeader
	if err := c.transport.DELETE(ctx, "/notification-groups/"+notificationGroupID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Logs

func (c *Client) ListLogFiles(ctx context.Context, instanceID string) (*ListLogFilesOutput, error) {
	var out ListLogFilesOutput
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID+"/log-files", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DownloadLogFile(ctx context.Context, instanceID, logFileID string) (string, error) {
	var out struct {
		Header      *ResponseHeader `json:"header"`
		DownloadURL string          `json:"downloadUrl"`
	}
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID+"/log-files/"+logFileID+"/download", &out); err != nil {
		return "", err
	}
	return out.DownloadURL, nil
}

func (c *Client) ListRegions(ctx context.Context) (*RegionsResponse, error) {
	var out RegionsResponse
	if err := c.transport.GET(ctx, "/project/regions", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListMembers(ctx context.Context) (*MembersResponse, error) {
	var out MembersResponse
	if err := c.transport.GET(ctx, "/project/members", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// HBA Rules

func (c *Client) ListHBARules(ctx context.Context, instanceID string) (*HBARulesResponse, error) {
	var out HBARulesResponse
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID+"/hba-rules", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetHBARule(ctx context.Context, instanceID, ruleID string) (*HBARule, error) {
	var out struct {
		Header  *ResponseHeader `json:"header"`
		HBARule HBARule         `json:"hbaRule"`
	}
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID+"/hba-rules/"+ruleID, &out); err != nil {
		return nil, err
	}
	return &out.HBARule, nil
}

func (c *Client) CreateHBARule(ctx context.Context, instanceID string, req *CreateHBARuleRequest) (*HBARule, error) {
	var out struct {
		Header  *ResponseHeader `json:"header"`
		HBARule HBARule         `json:"hbaRule"`
	}
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/hba-rules", req, &out); err != nil {
		return nil, err
	}
	return &out.HBARule, nil
}

func (c *Client) DeleteHBARule(ctx context.Context, instanceID, ruleID string) (*ResponseHeader, error) {
	var out ResponseHeader
	if err := c.transport.DELETE(ctx, "/db-instances/"+instanceID+"/hba-rules/"+ruleID, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListMetrics(ctx context.Context, instanceID string) (*MetricsResponse, error) {
	var out MetricsResponse
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID+"/metrics", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetMetricStatistics(ctx context.Context, instanceID, metricName string, startTime, endTime string) (*MetricStatisticsResponse, error) {
	path := fmt.Sprintf("/db-instances/%s/metrics/%s/statistics?startTime=%s&endTime=%s",
		instanceID, metricName, startTime, endTime)
	var out MetricStatisticsResponse
	if err := c.transport.GET(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ExportBackup(ctx context.Context, backupID string, tenantID string, containerName string) (*JobIDResponse, error) {
	req := map[string]string{
		"tenantId":      tenantID,
		"containerName": containerName,
	}
	var out JobIDResponse
	if err := c.transport.POST(ctx, "/backups/"+backupID+"/export", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) RestoreFromBackup(ctx context.Context, backupID string) (*JobIDResponse, error) {
	var out JobIDResponse
	if err := c.transport.POST(ctx, "/backups/"+backupID+"/restore", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListExtensions(ctx context.Context, instanceGroupID string) (*ExtensionsResponse, error) {
	var out ExtensionsResponse
	if err := c.transport.GET(ctx, "/db-instance-groups/"+instanceGroupID+"/extensions", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetExtension(ctx context.Context, instanceGroupID, extensionID string) (*Extension, error) {
	var out struct {
		Header    *ResponseHeader `json:"header"`
		Extension Extension       `json:"extension"`
	}
	if err := c.transport.GET(ctx, "/db-instance-groups/"+instanceGroupID+"/extensions/"+extensionID, &out); err != nil {
		return nil, err
	}
	return &out.Extension, nil
}

func (c *Client) InstallExtension(ctx context.Context, instanceGroupID, extensionID string, req *InstallExtensionRequest) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instance-groups/"+instanceGroupID+"/extensions/"+extensionID+"/install", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UninstallExtension(ctx context.Context, instanceGroupID, extensionID, databaseID string) (*JobOutput, error) {
	path := fmt.Sprintf("/db-instance-groups/%s/extensions/%s/uninstall?databaseId=%s",
		instanceGroupID, extensionID, databaseID)
	var out JobOutput
	if err := c.transport.POST(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ResizeStorage(ctx context.Context, instanceID string, newSizeGB int) (*JobOutput, error) {
	input := &ModifyStorageInfoInput{
		StorageSize: newSizeGB,
	}
	return c.ModifyStorageInfo(ctx, instanceID, input)
}

func (c *Client) ListEvents(ctx context.Context, params *EventParams) (*EventsResponse, error) {
	query := url.Values{}
	if params != nil {
		if params.InstanceID != "" {
			query.Set("dbInstanceId", params.InstanceID)
		}
		if params.StartTime != "" {
			query.Set("startTime", params.StartTime)
		}
		if params.EndTime != "" {
			query.Set("endTime", params.EndTime)
		}
		if params.EventCode != "" {
			query.Set("eventCode", params.EventCode)
		}
		if params.SourceType != "" {
			query.Set("sourceType", params.SourceType)
		}
		if params.Page > 0 {
			query.Set("page", fmt.Sprintf("%d", params.Page))
		}
		if params.Size > 0 {
			query.Set("size", fmt.Sprintf("%d", params.Size))
		}
	}

	path := "/events"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var out EventsResponse
	if err := c.transport.GET(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetWatchdog(ctx context.Context, instanceID string) (*Watchdog, error) {
	var out struct {
		Header   *ResponseHeader `json:"header"`
		Watchdog Watchdog        `json:"watchdog"`
	}
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID+"/watchdog", &out); err != nil {
		return nil, err
	}
	return &out.Watchdog, nil
}

func (c *Client) CreateWatchdog(ctx context.Context, instanceID string, req *CreateWatchdogRequest) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/watchdog", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateWatchdog(ctx context.Context, instanceID string, req *CreateWatchdogRequest) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.PUT(ctx, "/db-instances/"+instanceID+"/watchdog", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteWatchdog(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.DELETE(ctx, "/db-instances/"+instanceID+"/watchdog", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetNotificationGroupMonitoringItems(ctx context.Context, groupID string) (*MonitoringItemsResponse, error) {
	var out MonitoringItemsResponse
	if err := c.transport.GET(ctx, "/notification-groups/"+groupID+"/monitoring-items", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateNotificationGroupMonitoringItems(ctx context.Context, groupID string, items []MonitoringItem) (*JobOutput, error) {
	req := map[string]interface{}{
		"monitoringItems": items,
	}
	var out JobOutput
	if err := c.transport.PUT(ctx, "/notification-groups/"+groupID+"/monitoring-items", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
