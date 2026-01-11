package interactive

import (
	"context"
	"fmt"
	"strings"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/mariadb"
)

type MariaDBInteractive struct {
	client    *mariadb.Client
	pm        *PromptManager
	region    string
	azOptions []Option
}

func NewMariaDBInteractive(ctx context.Context, client *mariadb.Client, region string, azOptions []Option) *MariaDBInteractive {
	return &MariaDBInteractive{
		client:    client,
		pm:        NewPromptManager(ctx),
		region:    region,
		azOptions: azOptions,
	}
}

func (m *MariaDBInteractive) GetCreateDefinitions() []ParameterDef {
	return []ParameterDef{
		{
			Name:        "name",
			DisplayName: "Instance Name",
			Required:    true,
			Type:        TypeString,
			Validator:   ValidateInstanceName,
			Description: "Unique name for the database instance (4-50 chars, alphanumeric and hyphens)",
		},
		{
			Name:        "version",
			DisplayName: "MariaDB Version",
			Required:    true,
			Type:        TypeSelect,
			Fetcher:     m.fetchDBVersions,
			Default:     "latest",
			Description: "MariaDB engine version",
		},
		{
			Name:        "flavor-id",
			DisplayName: "Instance Type",
			Required:    true,
			Type:        TypeSelect,
			Fetcher:     m.fetchDBFlavors,
			Description: "Instance specification (CPU, memory)",
		},
		{
			Name:        "user-name",
			DisplayName: "Admin Username",
			Required:    true,
			Type:        TypeString,
			Validator:   ValidateUsername,
			Default:     "admin",
			Description: "Database administrator username",
		},
		{
			Name:        "password",
			DisplayName: "Admin Password",
			Required:    true,
			Type:        TypePassword,
			Validator:   ValidatePassword,
			Description: "Database administrator password (min 8 chars, mixed case + numbers)",
		},
		{
			Name:        "subnet-id",
			DisplayName: "Subnet",
			Required:    true,
			Type:        TypeSelect,
			Fetcher:     m.fetchSubnets,
			Description: "Network subnet for the database instance",
		},
		{
			Name:        "availability-zone",
			DisplayName: "Availability Zone",
			Required:    true,
			Type:        TypeSelect,
			Fetcher:     m.fetchAvailabilityZones,
			Description: "Availability zone for the instance",
		},
		{
			Name:        "storage-type",
			DisplayName: "Storage Type",
			Required:    true,
			Type:        TypeSelect,
			Fetcher:     m.fetchStorageTypes,
			Description: "Storage type for the database",
		},
		{
			Name:        "storage-size",
			DisplayName: "Storage Size (GB)",
			Required:    false,
			Type:        TypeInt,
			Default:     100,
			Validator:   ValidateStorageSize,
			Description: "Storage size in GB (20-6000)",
		},
		{
			Name:        "port",
			DisplayName: "Database Port",
			Required:    false,
			Type:        TypeInt,
			Default:     3306,
			Validator:   ValidatePort,
			Description: "Port number for MariaDB connections",
		},
		{
			Name:        "parameter-group-id",
			DisplayName: "Parameter Group",
			Required:    true,
			Type:        TypeSelect,
			Fetcher:     m.fetchParameterGroups,
			Description: "Database configuration parameter group (required)",
		},
		{
			Name:        "security-group-ids",
			DisplayName: "Security Groups",
			Required:    false,
			Type:        TypeMultiSelect,
			Fetcher:     m.fetchSecurityGroups,
			Description: "Security groups for network access control",
		},
		{
			Name:        "ha",
			DisplayName: "Enable High Availability",
			Required:    false,
			Type:        TypeConfirm,
			Default:     false,
			Description: "Enable Multi-AZ deployment for high availability",
		},
		{
			Name:        "deletion-protection",
			DisplayName: "Enable Deletion Protection",
			Required:    false,
			Type:        TypeConfirm,
			Default:     false,
			Description: "Prevent accidental deletion of the instance",
		},
		{
			Name:        "backup-period",
			DisplayName: "Backup Retention Period (days)",
			Required:    false,
			Type:        TypeInt,
			Default:     7,
			Validator:   ValidateBackupPeriod,
			Description: "Number of days to retain backups (0-35)",
		},
		{
			Name:        "backup-start-time",
			DisplayName: "Backup Window Start Time",
			Required:    false,
			Type:        TypeString,
			Default:     "02:00",
			Validator:   ValidateTimeFormat,
			Description: "Daily backup window start time (HH:MM format)",
		},
	}
}

func (m *MariaDBInteractive) fetchDBVersions(ctx context.Context) ([]Option, error) {
	response, err := m.client.ListVersions(ctx)
	if err != nil {
		return []Option{
			{Value: "10.6.8", Display: "MariaDB 10.6.8", Description: "MariaDB 10.6.8 (API unavailable)"},
			{Value: "10.5.17", Display: "MariaDB 10.5.17", Description: "MariaDB 10.5.17 (API unavailable)"},
		}, nil
	}

	if response == nil || len(response.DBVersions) == 0 {
		return []Option{
			{Value: "10.6.8", Display: "MariaDB 10.6.8", Description: "MariaDB 10.6.8 (default)"},
		}, nil
	}

	var options []Option
	for _, version := range response.DBVersions {
		options = append(options, Option{
			Value:       version.DBVersion,
			Display:     fmt.Sprintf("MariaDB %s", version.DBVersion),
			Description: version.DBVersionName,
		})
	}
	return options, nil
}

func (m *MariaDBInteractive) fetchDBFlavors(ctx context.Context) ([]Option, error) {
	response, err := m.client.ListFlavors(ctx)
	if err != nil {
		return []Option{
			{Value: "m2.c1m2", Display: "m2.c1m2", Description: "1 vCPU, 2GB RAM"},
			{Value: "m2.c2m4", Display: "m2.c2m4", Description: "2 vCPU, 4GB RAM"},
			{Value: "m2.c4m8", Display: "m2.c4m8", Description: "4 vCPU, 8GB RAM"},
		}, nil
	}

	if response == nil || len(response.DBFlavors) == 0 {
		return []Option{
			{Value: "m2.c2m4", Display: "m2.c2m4", Description: "2 vCPU, 4GB RAM (default)"},
		}, nil
	}

	var options []Option
	for _, flavor := range response.DBFlavors {
		description := fmt.Sprintf("%d vCPU, %dGB RAM", flavor.Vcpus, flavor.Ram/1024)
		options = append(options, Option{
			Value:       flavor.FlavorID,
			Display:     flavor.FlavorID,
			Description: description,
		})
	}
	return options, nil
}

func (m *MariaDBInteractive) fetchParameterGroups(ctx context.Context) ([]Option, error) {
	response, err := m.client.ListParameterGroups(ctx)
	if err != nil {
		return []Option{
			{Value: "", Display: "Default (skip)", Description: "Use default parameter group"},
		}, nil
	}

	if response == nil || len(response.ParameterGroups) == 0 {
		return []Option{
			{Value: "", Display: "Default (none available)", Description: "No parameter groups available"},
		}, nil
	}

	selectedVersion := ""
	if v, ok := m.pm.values["version"].(string); ok {
		selectedVersion = v
	}

	var options []Option
	for _, group := range response.ParameterGroups {
		if selectedVersion != "" && group.DBVersion != selectedVersion {
			continue
		}
		options = append(options, Option{
			Value:       group.ParameterGroupID,
			Display:     group.ParameterGroupName,
			Description: fmt.Sprintf("%s (%s)", group.DBVersion, group.Description),
		})
	}

	if len(options) == 0 {
		return []Option{
			{Value: "", Display: "No matching parameter groups", Description: fmt.Sprintf("No parameter groups for version %s", selectedVersion)},
		}, nil
	}

	return options, nil
}

func (m *MariaDBInteractive) fetchSecurityGroups(ctx context.Context) ([]Option, error) {
	response, err := m.client.ListSecurityGroups(ctx)
	if err != nil {
		return []Option{
			{Value: "", Display: "Default (skip)", Description: "Use default security group"},
		}, nil
	}

	if response == nil || len(response.DBSecurityGroups) == 0 {
		return []Option{
			{Value: "", Display: "Default (none available)", Description: "No security groups available"},
		}, nil
	}

	var options []Option
	for _, group := range response.DBSecurityGroups {
		options = append(options, Option{
			Value:       group.DBSecurityGroupID,
			Display:     group.DBSecurityGroupName,
			Description: group.Description,
		})
	}
	return options, nil
}

func (m *MariaDBInteractive) fetchSubnets(ctx context.Context) ([]Option, error) {
	response, err := m.client.ListSubnets(ctx)
	if err != nil {
		return []Option{
			{Value: "", Display: "Manual subnet required", Description: "Please provide --subnet-id"},
		}, nil
	}

	if response == nil || len(response.Subnets) == 0 {
		return []Option{
			{Value: "", Display: "No subnets available", Description: "Please create a subnet first"},
		}, nil
	}

	var options []Option
	for _, subnet := range response.Subnets {
		description := fmt.Sprintf("%s (%s)", subnet.SubnetCidr, subnet.SubnetName)
		if subnet.AvailableIpCount > 0 {
			description += fmt.Sprintf(" - %d IPs available", subnet.AvailableIpCount)
		}
		options = append(options, Option{
			Value:       subnet.SubnetID,
			Display:     subnet.SubnetName,
			Description: description,
		})
	}
	return options, nil
}

func (m *MariaDBInteractive) fetchAvailabilityZones(ctx context.Context) ([]Option, error) {
	if len(m.azOptions) > 0 {
		return m.azOptions, nil
	}

	allZones := map[string][]Option{
		"kr1": {
			{Value: "kr-pub-a", Display: "kr-pub-a", Description: "Korea Public zone A"},
			{Value: "kr-pub-b", Display: "kr-pub-b", Description: "Korea Public zone B"},
		},
		"kr2": {
			{Value: "kr2-pub-a", Display: "kr2-pub-a", Description: "Korea 2 Public zone A"},
			{Value: "kr2-pub-b", Display: "kr2-pub-b", Description: "Korea 2 Public zone B"},
		},
		"jp1": {
			{Value: "jp-pub-a", Display: "jp-pub-a", Description: "Japan Public zone A"},
			{Value: "jp-pub-b", Display: "jp-pub-b", Description: "Japan Public zone B"},
		},
	}

	if zones, ok := allZones[m.region]; ok {
		return zones, nil
	}
	return allZones["kr1"], nil
}

func (m *MariaDBInteractive) fetchStorageTypes(ctx context.Context) ([]Option, error) {
	response, err := m.client.ListStorageTypes(ctx)
	if err != nil {
		return []Option{
			{Value: "General SSD", Display: "General SSD", Description: "Default storage type"},
		}, nil
	}

	if response == nil || len(response.StorageTypes) == 0 {
		return []Option{
			{Value: "General SSD", Display: "General SSD", Description: "Default storage type"},
		}, nil
	}

	var options []Option
	for _, storageType := range response.StorageTypes {
		options = append(options, Option{
			Value:   storageType,
			Display: storageType,
		})
	}
	return options, nil
}

func (m *MariaDBInteractive) GetPromptManager() *PromptManager {
	return m.pm
}

func (m *MariaDBInteractive) SetDefinitions() {
	m.pm.SetDefinitions(m.GetCreateDefinitions())
}

func (m *MariaDBInteractive) GetModifyDefinitions() []ParameterDef {
	return []ParameterDef{
		{
			Name:        "name",
			DisplayName: "Instance Name",
			Required:    false,
			Type:        TypeString,
			Validator:   ValidateInstanceName,
			Description: "New name for the database instance",
		},
		{
			Name:        "description",
			DisplayName: "Description",
			Required:    false,
			Type:        TypeString,
			Description: "New description for the database instance",
		},
		{
			Name:        "flavor-id",
			DisplayName: "Instance Type",
			Required:    false,
			Type:        TypeSelect,
			Fetcher:     m.fetchDBFlavors,
			Description: "New instance specification (CPU, memory)",
		},
		{
			Name:        "storage-size",
			DisplayName: "Storage Size (GB)",
			Required:    false,
			Type:        TypeInt,
			Validator:   ValidateStorageSize,
			Description: "New storage size in GB (can only be increased)",
		},
		{
			Name:        "parameter-group-id",
			DisplayName: "Parameter Group",
			Required:    false,
			Type:        TypeSelect,
			Fetcher:     m.fetchParameterGroups,
			Description: "New parameter group",
		},
		{
			Name:        "security-group-ids",
			DisplayName: "Security Groups",
			Required:    false,
			Type:        TypeMultiSelect,
			Fetcher:     m.fetchSecurityGroups,
			Description: "New security groups",
		},
		{
			Name:        "deletion-protection",
			DisplayName: "Deletion Protection",
			Required:    false,
			Type:        TypeConfirm,
			Description: "Enable/disable deletion protection",
		},
	}
}

func (m *MariaDBInteractive) SetModifyDefinitions() {
	m.pm.SetDefinitions(m.GetModifyDefinitions())
}

func (m *MariaDBInteractive) ShowCurrentConfiguration(instance *mariadb.DatabaseInstance) {
	fmt.Println("\n" + strings.Repeat("═", 50))
	fmt.Printf("📋 Current MariaDB Instance Configuration\n")
	fmt.Println(strings.Repeat("═", 50))
	fmt.Printf("%-20s: %s\n", "Instance ID", instance.DBInstanceID)
	fmt.Printf("%-20s: %s\n", "Instance Name", instance.DBInstanceName)
	fmt.Printf("%-20s: %s\n", "Status", instance.DBInstanceStatus)
	fmt.Printf("%-20s: %s\n", "Version", instance.DBVersion)
	fmt.Printf("%-20s: %s\n", "Flavor ID", instance.DBFlavorID)
	fmt.Println(strings.Repeat("─", 50))
}
