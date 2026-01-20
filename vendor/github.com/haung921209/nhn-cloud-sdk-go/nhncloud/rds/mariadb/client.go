package mariadb

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
	baseURL := endpoint.ResolveWithAppKey(endpoint.ServiceRDSMariaDB, region, appKey)

	opts := []transport.ClientOption{
		transport.WithDebug(debug),
	}

	if creds != nil {
		opts = append(opts, transport.WithAppKeyAuth(
			appKey,
			creds.GetAccessKeyID(),
			creds.GetSecretAccessKey(),
		))
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

func (c *Client) CreateInstance(ctx context.Context, input *CreateInstanceInput) (*CreateInstanceOutput, error) {
	var out CreateInstanceOutput
	if err := c.transport.POST(ctx, "/db-instances", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ModifyInstance(ctx context.Context, instanceID string, input *ModifyInstanceInput) (*GetInstanceOutput, error) {
	var out GetInstanceOutput
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

func (c *Client) UpdateSecurityGroup(ctx context.Context, securityGroupID string, input *UpdateSecurityGroupInput) (*ResponseHeader, error) {
	var out ResponseHeader
	if err := c.transport.PUT(ctx, "/db-security-groups/"+securityGroupID, input, &out); err != nil {
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

func (c *Client) UpdateSecurityGroupRule(ctx context.Context, securityGroupID, ruleID string, input *UpdateSecurityGroupRuleInput) (*ResponseHeader, error) {
	var out ResponseHeader
	if err := c.transport.PUT(ctx, "/db-security-groups/"+securityGroupID+"/rules/"+ruleID, input, &out); err != nil {
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

func (c *Client) CopyParameterGroup(ctx context.Context, parameterGroupID string, input *CopyParameterGroupInput) (*ParameterGroupIDOutput, error) {
	var out ParameterGroupIDOutput
	if err := c.transport.POST(ctx, "/parameter-groups/"+parameterGroupID+"/copy", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateParameterGroup(ctx context.Context, parameterGroupID string, input *UpdateParameterGroupInput) (*ResponseHeader, error) {
	var out ResponseHeader
	if err := c.transport.PUT(ctx, "/parameter-groups/"+parameterGroupID, input, &out); err != nil {
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

func (c *Client) ResetParameterGroup(ctx context.Context, parameterGroupID string) (*ResponseHeader, error) {
	var out ResponseHeader
	if err := c.transport.PUT(ctx, "/parameter-groups/"+parameterGroupID+"/reset", nil, &out); err != nil {
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

func (c *Client) ListBackups(ctx context.Context, instanceID, dbVersion string, page, size int) (*ListBackupsOutput, error) {
	path := fmt.Sprintf("/backups?page=%d&size=%d", page, size)
	if instanceID != "" {
		path += "&dbInstanceId=" + url.QueryEscape(instanceID)
	}
	if dbVersion != "" {
		path += "&dbVersion=" + url.QueryEscape(dbVersion)
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

func (c *Client) BackupToObjectStorage(ctx context.Context, instanceID string, input *BackupToObjectStorageInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/backup-to-object-storage", input, &out); err != nil {
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

func (c *Client) ExportBackup(ctx context.Context, backupID string, input *ExportBackupInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/backups/"+backupID+"/export", input, &out); err != nil {
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

// Schemas

func (c *Client) ListSchemas(ctx context.Context, instanceID string) (*ListSchemasOutput, error) {
	var out ListSchemasOutput
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID+"/db-schemas", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateSchema(ctx context.Context, instanceID string, input *CreateSchemaInput) (*SchemaIDOutput, error) {
	var out SchemaIDOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/db-schemas", input, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteSchema(ctx context.Context, instanceID, schemaID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.DELETE(ctx, "/db-instances/"+instanceID+"/db-schemas/"+schemaID, &out); err != nil {
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

func (c *Client) ModifyNetworkInfo(ctx context.Context, instanceID string, input *ModifyNetworkInfoInput) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.PUT(ctx, "/db-instances/"+instanceID+"/network-info", input, &out); err != nil {
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

func (c *Client) SplitHighAvailability(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.transport.POST(ctx, "/db-instances/"+instanceID+"/high-availability/split", nil, &out); err != nil {
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

func (c *Client) UpdateNotificationGroup(ctx context.Context, notificationGroupID string, input *UpdateNotificationGroupInput) (*ResponseHeader, error) {
	var out ResponseHeader
	if err := c.transport.PUT(ctx, "/notification-groups/"+notificationGroupID, input, &out); err != nil {
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

// Logs and Metrics

func (c *Client) ListLogFiles(ctx context.Context, instanceID string) (*ListLogFilesOutput, error) {
	var out ListLogFilesOutput
	if err := c.transport.GET(ctx, "/db-instances/"+instanceID+"/log-files", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListMetrics(ctx context.Context) (*ListMetricsOutput, error) {
	var out ListMetricsOutput
	if err := c.transport.GET(ctx, "/metrics", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetMetricStatistics(ctx context.Context, instanceID, from, to string, interval *int) (*MetricStatisticsOutput, error) {
	path := fmt.Sprintf("/metric-statistics?dbInstanceId=%s&from=%s&to=%s",
		url.QueryEscape(instanceID), url.QueryEscape(from), url.QueryEscape(to))
	if interval != nil {
		path += fmt.Sprintf("&interval=%d", *interval)
	}
	var out MetricStatisticsOutput
	if err := c.transport.GET(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
