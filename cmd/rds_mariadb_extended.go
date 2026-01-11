package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/mariadb"
	"github.com/spf13/cobra"
)

// ============================================================================
// MariaDB Extended Commands - Users, Schemas, Security Groups, etc.
// ============================================================================

// DB User Commands
var mariadbUserCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "Manage database users",
}

var mariadbListDBUsersCmd = &cobra.Command{
	Use:   "list [instance-id]",
	Short: "List all database users",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListDBUsers(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list users", err)
		}
		printMariaDBUsers(result)
	},
}

var mariadbCreateDBUserCmd = &cobra.Command{
	Use:   "create [instance-id]",
	Short: "Create a database user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userName, _ := cmd.Flags().GetString("name")
		password, _ := cmd.Flags().GetString("password")
		host, _ := cmd.Flags().GetString("host")

		if userName == "" || password == "" {
			exitWithError("--name and --password are required", nil)
		}

		input := &mariadb.CreateDBUserInput{
			DBUserName: userName,
			DBPassword: password,
			Host:       host,
		}

		client := newMariaDBClient()
		result, err := client.CreateDBUser(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create user", err)
		}
		fmt.Printf("User creation initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbDeleteDBUserCmd = &cobra.Command{
	Use:   "delete [instance-id] [user-id]",
	Short: "Delete a database user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.DeleteDBUser(context.Background(), args[0], args[1])
		if err != nil {
			exitWithError("failed to delete user", err)
		}
		fmt.Printf("User deletion initiated. Job ID: %s\n", result.JobID)
	},
}

// Schema Commands
var mariadbSchemaCmd = &cobra.Command{
	Use:     "schema",
	Aliases: []string{"schemas", "database", "db"},
	Short:   "Manage database schemas",
}

var mariadbListSchemasCmd = &cobra.Command{
	Use:   "list [instance-id]",
	Short: "List all schemas",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListSchemas(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list schemas", err)
		}
		printMariaDBSchemas(result)
	},
}

var mariadbCreateSchemaCmd = &cobra.Command{
	Use:   "create [instance-id]",
	Short: "Create a database schema",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			exitWithError("--name is required", nil)
		}

		input := &mariadb.CreateSchemaInput{DBSchemaName: name}
		client := newMariaDBClient()
		result, err := client.CreateSchema(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create schema", err)
		}
		fmt.Printf("Schema created. ID: %s\n", result.DBSchemaID)
	},
}

var mariadbDeleteSchemaCmd = &cobra.Command{
	Use:   "delete [instance-id] [schema-id]",
	Short: "Delete a database schema",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.DeleteSchema(context.Background(), args[0], args[1])
		if err != nil {
			exitWithError("failed to delete schema", err)
		}
		fmt.Printf("Schema deletion initiated. Job ID: %s\n", result.JobID)
	},
}

// Security Group Commands
var mariadbSGCmd = &cobra.Command{
	Use:     "security-group",
	Aliases: []string{"sg"},
	Short:   "Manage DB security groups",
}

var mariadbListSGCmd = &cobra.Command{
	Use:   "list",
	Short: "List all security groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListSecurityGroups(context.Background())
		if err != nil {
			exitWithError("failed to list security groups", err)
		}
		printMariaDBSecurityGroups(result)
	},
}

var mariadbGetSGCmd = &cobra.Command{
	Use:   "get [security-group-id]",
	Short: "Get security group details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.GetSecurityGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get security group", err)
		}
		printMariaDBSecurityGroupDetail(result)
	},
}

var mariadbCreateSGCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a security group",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		if name == "" {
			exitWithError("--name is required", nil)
		}

		input := &mariadb.CreateSecurityGroupInput{
			DBSecurityGroupName: name,
			Description:         desc,
		}
		client := newMariaDBClient()
		result, err := client.CreateSecurityGroup(context.Background(), input)
		if err != nil {
			exitWithError("failed to create security group", err)
		}
		fmt.Printf("Security group created. ID: %s\n", result.DBSecurityGroupID)
	},
}

var mariadbDeleteSGCmd = &cobra.Command{
	Use:   "delete [security-group-id]",
	Short: "Delete a security group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		_, err := client.DeleteSecurityGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete security group", err)
		}
		fmt.Println("Security group deleted.")
	},
}

// Parameter Group Commands
var mariadbPGCmd = &cobra.Command{
	Use:     "parameter-group",
	Aliases: []string{"pg"},
	Short:   "Manage parameter groups",
}

var mariadbListPGCmd = &cobra.Command{
	Use:   "list",
	Short: "List all parameter groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListParameterGroups(context.Background())
		if err != nil {
			exitWithError("failed to list parameter groups", err)
		}
		printMariaDBParameterGroups(result)
	},
}

var mariadbGetPGCmd = &cobra.Command{
	Use:   "get [parameter-group-id]",
	Short: "Get parameter group details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.GetParameterGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get parameter group", err)
		}
		printMariaDBParameterGroupDetail(result)
	},
}

var mariadbCreatePGCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a parameter group",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		version, _ := cmd.Flags().GetString("version")
		if name == "" || version == "" {
			exitWithError("--name and --version are required", nil)
		}

		input := &mariadb.CreateParameterGroupInput{
			ParameterGroupName: name,
			Description:        desc,
			DBVersion:          version,
		}
		client := newMariaDBClient()
		result, err := client.CreateParameterGroup(context.Background(), input)
		if err != nil {
			exitWithError("failed to create parameter group", err)
		}
		fmt.Printf("Parameter group created. ID: %s\n", result.ParameterGroupID)
	},
}

var mariadbDeletePGCmd = &cobra.Command{
	Use:   "delete [parameter-group-id]",
	Short: "Delete a parameter group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		_, err := client.DeleteParameterGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete parameter group", err)
		}
		fmt.Println("Parameter group deleted.")
	},
}

// Notification Group Commands
var mariadbNGCmd = &cobra.Command{
	Use:     "notification-group",
	Aliases: []string{"ng"},
	Short:   "Manage notification groups",
}

var mariadbListNGCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notification groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListNotificationGroups(context.Background())
		if err != nil {
			exitWithError("failed to list notification groups", err)
		}
		printMariaDBNotificationGroups(result)
	},
}

var mariadbGetNGCmd = &cobra.Command{
	Use:   "get [notification-group-id]",
	Short: "Get notification group details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.GetNotificationGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get notification group", err)
		}
		printMariaDBNotificationGroupDetail(result)
	},
}

var mariadbDeleteNGCmd = &cobra.Command{
	Use:   "delete [notification-group-id]",
	Short: "Delete a notification group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		_, err := client.DeleteNotificationGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete notification group", err)
		}
		fmt.Println("Notification group deleted.")
	},
}

// Log Commands
var mariadbLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Manage logs",
}

var mariadbListLogsCmd = &cobra.Command{
	Use:   "list [instance-id]",
	Short: "List log files",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListLogFiles(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list log files", err)
		}
		printMariaDBLogFiles(result)
	},
}

// Metrics Commands
var mariadbMetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "View metrics",
}

var mariadbListMetricsCmd = &cobra.Command{
	Use:   "list",
	Short: "List available metrics",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListMetrics(context.Background())
		if err != nil {
			exitWithError("failed to list metrics", err)
		}
		printMariaDBMetrics(result)
	},
}

// Resource Commands
var mariadbStorageTypesCmd = &cobra.Command{
	Use:   "storage-types",
	Short: "List available storage types",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListStorageTypes(context.Background())
		if err != nil {
			exitWithError("failed to list storage types", err)
		}
		for _, t := range result.StorageTypes {
			fmt.Println(t)
		}
	},
}

var mariadbSubnetsCmd = &cobra.Command{
	Use:   "subnets",
	Short: "List available subnets",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListSubnets(context.Background())
		if err != nil {
			exitWithError("failed to list subnets", err)
		}
		printMariaDBSubnets(result)
	},
}

// Network Commands
var mariadbNetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage network settings",
}

var mariadbGetNetworkCmd = &cobra.Command{
	Use:   "get [instance-id]",
	Short: "Get network information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.GetNetworkInfo(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get network info", err)
		}
		printMariaDBNetworkInfo(result)
	},
}

// Storage & Protection Commands
var mariadbResizeStorageCmd = &cobra.Command{
	Use:   "resize-storage [instance-id]",
	Short: "Resize storage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		size, _ := cmd.Flags().GetInt("size")
		if size <= 0 {
			exitWithError("--size is required and must be positive", nil)
		}

		input := &mariadb.ModifyStorageInfoInput{StorageSize: size}
		client := newMariaDBClient()
		result, err := client.ModifyStorageInfo(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to resize storage", err)
		}
		fmt.Printf("Storage resize initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbDeletionProtectionCmd = &cobra.Command{
	Use:   "deletion-protection [instance-id]",
	Short: "Enable or disable deletion protection",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		enable, _ := cmd.Flags().GetBool("enable")
		disable, _ := cmd.Flags().GetBool("disable")
		if enable == disable {
			exitWithError("specify either --enable or --disable", nil)
		}

		input := &mariadb.ModifyDeletionProtectionInput{UseDeletionProtection: enable}
		client := newMariaDBClient()
		result, err := client.ModifyDeletionProtection(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to modify deletion protection", err)
		}
		fmt.Printf("Deletion protection modified. Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Init
// ============================================================================

func init() {
	// User commands
	rdsMariaDBCmd.AddCommand(mariadbUserCmd)
	mariadbUserCmd.AddCommand(mariadbListDBUsersCmd)
	mariadbUserCmd.AddCommand(mariadbCreateDBUserCmd)
	mariadbUserCmd.AddCommand(mariadbDeleteDBUserCmd)
	mariadbCreateDBUserCmd.Flags().String("name", "", "User name (required)")
	mariadbCreateDBUserCmd.Flags().String("password", "", "Password (required)")
	mariadbCreateDBUserCmd.Flags().String("host", "%", "Host")

	// Schema commands
	rdsMariaDBCmd.AddCommand(mariadbSchemaCmd)
	mariadbSchemaCmd.AddCommand(mariadbListSchemasCmd)
	mariadbSchemaCmd.AddCommand(mariadbCreateSchemaCmd)
	mariadbSchemaCmd.AddCommand(mariadbDeleteSchemaCmd)
	mariadbCreateSchemaCmd.Flags().String("name", "", "Schema name (required)")

	// Security Group commands
	rdsMariaDBCmd.AddCommand(mariadbSGCmd)
	mariadbSGCmd.AddCommand(mariadbListSGCmd)
	mariadbSGCmd.AddCommand(mariadbGetSGCmd)
	mariadbSGCmd.AddCommand(mariadbCreateSGCmd)
	mariadbSGCmd.AddCommand(mariadbDeleteSGCmd)
	mariadbCreateSGCmd.Flags().String("name", "", "Name (required)")
	mariadbCreateSGCmd.Flags().String("description", "", "Description")

	// Parameter Group commands
	rdsMariaDBCmd.AddCommand(mariadbPGCmd)
	mariadbPGCmd.AddCommand(mariadbListPGCmd)
	mariadbPGCmd.AddCommand(mariadbGetPGCmd)
	mariadbPGCmd.AddCommand(mariadbCreatePGCmd)
	mariadbPGCmd.AddCommand(mariadbDeletePGCmd)
	mariadbCreatePGCmd.Flags().String("name", "", "Name (required)")
	mariadbCreatePGCmd.Flags().String("description", "", "Description")
	mariadbCreatePGCmd.Flags().String("version", "", "DB version (required)")

	// Notification Group commands
	rdsMariaDBCmd.AddCommand(mariadbNGCmd)
	mariadbNGCmd.AddCommand(mariadbListNGCmd)
	mariadbNGCmd.AddCommand(mariadbGetNGCmd)
	mariadbNGCmd.AddCommand(mariadbDeleteNGCmd)

	// Log commands
	rdsMariaDBCmd.AddCommand(mariadbLogCmd)
	mariadbLogCmd.AddCommand(mariadbListLogsCmd)

	// Metrics commands
	rdsMariaDBCmd.AddCommand(mariadbMetricsCmd)
	mariadbMetricsCmd.AddCommand(mariadbListMetricsCmd)

	// Resource commands
	rdsMariaDBCmd.AddCommand(mariadbStorageTypesCmd)
	rdsMariaDBCmd.AddCommand(mariadbSubnetsCmd)

	// Network commands
	rdsMariaDBCmd.AddCommand(mariadbNetworkCmd)
	mariadbNetworkCmd.AddCommand(mariadbGetNetworkCmd)

	// Storage & Protection commands
	rdsMariaDBCmd.AddCommand(mariadbResizeStorageCmd)
	mariadbResizeStorageCmd.Flags().Int("size", 0, "New storage size in GB (required)")
	rdsMariaDBCmd.AddCommand(mariadbDeletionProtectionCmd)
	mariadbDeletionProtectionCmd.Flags().Bool("enable", false, "Enable deletion protection")
	mariadbDeletionProtectionCmd.Flags().Bool("disable", false, "Disable deletion protection")
}

// ============================================================================
// Print Functions
// ============================================================================

func printMariaDBUsers(result *mariadb.ListDBUsersOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tHOST\tCREATED")
	for _, u := range result.DBUsers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", u.DBUserID, u.DBUserName, u.Host, u.CreatedYmdt)
	}
	w.Flush()
}

func printMariaDBSchemas(result *mariadb.ListSchemasOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCREATED")
	for _, s := range result.DBSchemas {
		fmt.Fprintf(w, "%s\t%s\t%s\n", s.DBSchemaID, s.DBSchemaName, s.CreatedYmdt)
	}
	w.Flush()
}

func printMariaDBSecurityGroups(result *mariadb.ListSecurityGroupsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tRULES")
	for _, sg := range result.DBSecurityGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", sg.DBSecurityGroupID, sg.DBSecurityGroupName, sg.Description, len(sg.Rules))
	}
	w.Flush()
}

func printMariaDBSecurityGroupDetail(result *mariadb.SecurityGroupOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	fmt.Printf("ID:          %s\n", result.DBSecurityGroup.DBSecurityGroupID)
	fmt.Printf("Name:        %s\n", result.DBSecurityGroup.DBSecurityGroupName)
	fmt.Printf("Description: %s\n", result.DBSecurityGroup.Description)
	fmt.Println("\nRules:")
	for _, r := range result.DBSecurityGroup.Rules {
		fmt.Printf("  - %s: %s %s\n", r.RuleID, r.Direction, r.CIDR)
	}
}

func printMariaDBParameterGroups(result *mariadb.ListParameterGroupsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tVERSION\tSTATUS")
	for _, pg := range result.ParameterGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", pg.ParameterGroupID, pg.ParameterGroupName, pg.DBVersion, pg.ParameterGroupStatus)
	}
	w.Flush()
}

func printMariaDBParameterGroupDetail(result *mariadb.ParameterGroupOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	fmt.Printf("ID:          %s\n", result.ParameterGroupID)
	fmt.Printf("Name:        %s\n", result.ParameterGroupName)
	fmt.Printf("Version:     %s\n", result.DBVersion)
	fmt.Printf("Status:      %s\n", result.ParameterGroupStatus)
	fmt.Printf("Parameters:  %d\n", len(result.Parameters))
}

func printMariaDBNotificationGroups(result *mariadb.ListNotificationGroupsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tTYPE\tENABLED")
	for _, ng := range result.NotificationGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", ng.NotificationGroupID, ng.NotificationGroupName, ng.NotificationType, ng.IsEnabled)
	}
	w.Flush()
}

func printMariaDBNotificationGroupDetail(result *mariadb.NotificationGroupOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	fmt.Printf("ID:      %s\n", result.NotificationGroup.NotificationGroupID)
	fmt.Printf("Name:    %s\n", result.NotificationGroup.NotificationGroupName)
	fmt.Printf("Type:    %s\n", result.NotificationGroup.NotificationType)
	fmt.Printf("Enabled: %v\n", result.NotificationGroup.IsEnabled)
}

func printMariaDBLogFiles(result *mariadb.ListLogFilesOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "FILE NAME\tSIZE\tCREATED")
	for _, l := range result.LogFiles {
		fmt.Fprintf(w, "%s\t%d\t%s\n", l.LogFileName, l.LogFileSize, l.CreatedYmdt)
	}
	w.Flush()
}

func printMariaDBMetrics(result *mariadb.ListMetricsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tUNIT")
	for _, m := range result.Metrics {
		fmt.Fprintf(w, "%s\t%s\n", m.MeasureName, m.Unit)
	}
	w.Flush()
}

func printMariaDBSubnets(result *mariadb.ListSubnetsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCIDR")
	for _, s := range result.Subnets {
		fmt.Fprintf(w, "%s\t%s\t%s\n", s.SubnetID, s.SubnetName, s.SubnetCidr)
	}
	w.Flush()
}

func printMariaDBNetworkInfo(result *mariadb.NetworkInfoOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	fmt.Printf("Availability Zone: %s\n", result.AvailabilityZone)
	fmt.Printf("Subnet ID:         %s\n", result.Subnet.SubnetID)
	fmt.Printf("Subnet Name:       %s\n", result.Subnet.SubnetName)
	fmt.Println("Endpoints:")
	for _, ep := range result.EndPoints {
		fmt.Printf("  - %s: %s (%s)\n", ep.EndPointType, ep.Domain, ep.IPAddress)
	}
}
