package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// DB User Commands
// ============================================================================

var dbUserCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "Manage database users",
}

var listDBUsersCmd = &cobra.Command{
	Use:   "list [instance-id]",
	Short: "List all database users for an instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListDBUsers(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list database users", err)
		}
		printDBUsers(result)
	},
}

var createDBUserCmd = &cobra.Command{
	Use:   "create [instance-id]",
	Short: "Create a database user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userName, _ := cmd.Flags().GetString("name")
		password, _ := cmd.Flags().GetString("password")
		hostIP, _ := cmd.Flags().GetString("host")
		authPlugin, _ := cmd.Flags().GetString("auth-plugin")
		authorities, _ := cmd.Flags().GetStringSlice("authorities")

		if userName == "" || password == "" {
			exitWithError("--name and --password are required", nil)
		}

		authorityType := "READ"
		if len(authorities) > 0 {
			authorityType = authorities[0]
		}
		input := &mysql.CreateDBUserInput{
			DBUserName:           userName,
			DBPassword:           password,
			Host:                 hostIP,
			AuthorityType:        authorityType,
			AuthenticationPlugin: authPlugin,
		}

		client := newMySQLClient()
		result, err := client.CreateDBUser(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create database user", err)
		}
		fmt.Printf("Database user creation initiated. Job ID: %s\n", result.JobID)
	},
}

var updateDBUserCmd = &cobra.Command{
	Use:   "update [instance-id] [user-id]",
	Short: "Update a database user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")
		authorities, _ := cmd.Flags().GetStringSlice("authorities")

		if password == "" && len(authorities) == 0 {
			exitWithError("at least --password or --authorities is required", nil)
		}

		authorityType := ""
		if len(authorities) > 0 {
			authorityType = authorities[0]
		}
		input := &mysql.UpdateDBUserInput{
			AuthorityType: authorityType,
		}

		client := newMySQLClient()
		result, err := client.UpdateDBUser(context.Background(), args[0], args[1], input)
		if err != nil {
			exitWithError("failed to update database user", err)
		}
		fmt.Printf("Database user update initiated. Job ID: %s\n", result.JobID)
	},
}

var deleteDBUserCmd = &cobra.Command{
	Use:   "delete [instance-id] [user-id]",
	Short: "Delete a database user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.DeleteDBUser(context.Background(), args[0], args[1])
		if err != nil {
			exitWithError("failed to delete database user", err)
		}
		fmt.Printf("Database user deletion initiated. Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Schema Commands
// ============================================================================

var schemaCmd = &cobra.Command{
	Use:     "schema",
	Aliases: []string{"schemas", "database", "db"},
	Short:   "Manage database schemas",
}

var listSchemasCmd = &cobra.Command{
	Use:   "list [instance-id]",
	Short: "List all schemas for an instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListSchemas(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list schemas", err)
		}
		printSchemas(result)
	},
}

var createSchemaCmd = &cobra.Command{
	Use:   "create [instance-id]",
	Short: "Create a database schema",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		input := &mysql.CreateSchemaInput{
			DBSchemaName: name,
		}

		client := newMySQLClient()
		result, err := client.CreateSchema(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create schema", err)
		}
		fmt.Printf("Schema created. ID: %s, Job ID: %s\n", result.DBSchemaID, result.JobID)
	},
}

var deleteSchemaCmd = &cobra.Command{
	Use:   "delete [instance-id] [schema-id]",
	Short: "Delete a database schema",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.DeleteSchema(context.Background(), args[0], args[1])
		if err != nil {
			exitWithError("failed to delete schema", err)
		}
		fmt.Printf("Schema deletion initiated. Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Notification Group Commands
// ============================================================================

var notificationGroupCmd = &cobra.Command{
	Use:     "notification-group",
	Aliases: []string{"ng"},
	Short:   "Manage notification groups",
}

var listNotificationGroupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notification groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListNotificationGroups(context.Background())
		if err != nil {
			exitWithError("failed to list notification groups", err)
		}
		printNotificationGroups(result)
	},
}

var getNotificationGroupCmd = &cobra.Command{
	Use:   "get [notification-group-id]",
	Short: "Get details of a notification group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.GetNotificationGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get notification group", err)
		}
		printNotificationGroupDetail(result)
	},
}

var createNotificationGroupCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a notification group",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		notifyEmail, _ := cmd.Flags().GetBool("email")
		notifySMS, _ := cmd.Flags().GetBool("sms")
		enabled, _ := cmd.Flags().GetBool("enabled")
		instanceIDs, _ := cmd.Flags().GetStringSlice("instance-ids")
		userGroupIDs, _ := cmd.Flags().GetStringSlice("user-group-ids")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		notificationType := "EMAIL"
		if notifySMS {
			notificationType = "SMS"
		}
		_ = notifyEmail
		_ = instanceIDs
		_ = userGroupIDs
		input := &mysql.CreateNotificationGroupInput{
			NotificationGroupName: name,
			NotificationType:      notificationType,
			IsEnabled:             enabled,
		}

		client := newMySQLClient()
		result, err := client.CreateNotificationGroup(context.Background(), input)
		if err != nil {
			exitWithError("failed to create notification group", err)
		}
		fmt.Printf("Notification group created. ID: %s\n", result.NotificationGroupID)
	},
}

var deleteNotificationGroupCmd = &cobra.Command{
	Use:   "delete [notification-group-id]",
	Short: "Delete a notification group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		_, err := client.DeleteNotificationGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete notification group", err)
		}
		fmt.Println("Notification group deleted successfully.")
	},
}

// ============================================================================
// Log and Metrics Commands
// ============================================================================

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Manage logs",
}

var listLogsCmd = &cobra.Command{
	Use:   "list [instance-id]",
	Short: "List log files for an instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListLogFiles(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list log files", err)
		}
		printLogFiles(result)
	},
}

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "View metrics",
}

var listMetricsCmd = &cobra.Command{
	Use:   "list",
	Short: "List available metrics",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListMetrics(context.Background())
		if err != nil {
			exitWithError("failed to list metrics", err)
		}
		printMetrics(result)
	},
}

var getMetricStatsCmd = &cobra.Command{
	Use:   "stats [instance-id]",
	Short: "Get metric statistics for an instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		interval, _ := cmd.Flags().GetInt("interval")

		if from == "" || to == "" {
			exitWithError("--from and --to are required (format: YYYY-MM-DD HH:MM)", nil)
		}

		var intervalPtr *int
		if interval > 0 {
			intervalPtr = &interval
		}

		client := newMySQLClient()
		result, err := client.GetMetricStatistics(context.Background(), args[0], from, to, intervalPtr)
		if err != nil {
			exitWithError("failed to get metric statistics", err)
		}
		printMetricStats(result)
	},
}

// ============================================================================
// Network Info Commands
// ============================================================================

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage network settings",
}

var getNetworkInfoCmd = &cobra.Command{
	Use:   "get [instance-id]",
	Short: "Get network information for an instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.GetNetworkInfo(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get network info", err)
		}
		printNetworkInfo(result)
	},
}

var setPublicAccessCmd = &cobra.Command{
	Use:   "public-access [instance-id]",
	Short: "Enable or disable public access",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		enable, _ := cmd.Flags().GetBool("enable")
		disable, _ := cmd.Flags().GetBool("disable")

		if enable == disable {
			exitWithError("specify either --enable or --disable", nil)
		}

		input := &mysql.ModifyNetworkInfoInput{
			UsePublicAccess: enable,
		}

		client := newMySQLClient()
		result, err := client.ModifyNetworkInfo(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to modify network info", err)
		}
		fmt.Printf("Network modification initiated. Job ID: %s\n", result.JobID)
	},
}

var resizeStorageCmd = &cobra.Command{
	Use:   "resize-storage [instance-id]",
	Short: "Resize storage for an instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		size, _ := cmd.Flags().GetInt("size")

		if size <= 0 {
			exitWithError("--size is required and must be positive", nil)
		}

		input := &mysql.ModifyStorageInfoInput{
			StorageSize: size,
		}

		client := newMySQLClient()
		result, err := client.ModifyStorageInfo(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to resize storage", err)
		}
		fmt.Printf("Storage resize initiated. Job ID: %s\n", result.JobID)
	},
}

var setDeletionProtectionCmd = &cobra.Command{
	Use:   "deletion-protection [instance-id]",
	Short: "Enable or disable deletion protection",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		enable, _ := cmd.Flags().GetBool("enable")
		disable, _ := cmd.Flags().GetBool("disable")

		if enable == disable {
			exitWithError("specify either --enable or --disable", nil)
		}

		input := &mysql.ModifyDeletionProtectionInput{
			UseDeletionProtection: enable,
		}

		client := newMySQLClient()
		result, err := client.ModifyDeletionProtection(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to modify deletion protection", err)
		}
		fmt.Printf("Deletion protection modified. Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	// DB User commands
	rdsMySQLCmd.AddCommand(dbUserCmd)
	dbUserCmd.AddCommand(listDBUsersCmd)
	dbUserCmd.AddCommand(createDBUserCmd)
	dbUserCmd.AddCommand(updateDBUserCmd)
	dbUserCmd.AddCommand(deleteDBUserCmd)

	createDBUserCmd.Flags().String("name", "", "User name (required)")
	createDBUserCmd.Flags().String("password", "", "Password (required)")
	createDBUserCmd.Flags().String("host", "%", "Host IP (default: % for any)")
	createDBUserCmd.Flags().String("auth-plugin", "", "Auth plugin")
	createDBUserCmd.Flags().StringSlice("authorities", nil, "Authorities")

	updateDBUserCmd.Flags().String("password", "", "New password")
	updateDBUserCmd.Flags().StringSlice("authorities", nil, "New authorities")

	// Schema commands
	rdsMySQLCmd.AddCommand(schemaCmd)
	schemaCmd.AddCommand(listSchemasCmd)
	schemaCmd.AddCommand(createSchemaCmd)
	schemaCmd.AddCommand(deleteSchemaCmd)

	createSchemaCmd.Flags().String("name", "", "Schema name (required)")

	// Notification Group commands
	rdsMySQLCmd.AddCommand(notificationGroupCmd)
	notificationGroupCmd.AddCommand(listNotificationGroupsCmd)
	notificationGroupCmd.AddCommand(getNotificationGroupCmd)
	notificationGroupCmd.AddCommand(createNotificationGroupCmd)
	notificationGroupCmd.AddCommand(deleteNotificationGroupCmd)

	createNotificationGroupCmd.Flags().String("name", "", "Notification group name (required)")
	createNotificationGroupCmd.Flags().Bool("email", true, "Enable email notifications")
	createNotificationGroupCmd.Flags().Bool("sms", false, "Enable SMS notifications")
	createNotificationGroupCmd.Flags().Bool("enabled", true, "Enable the notification group")
	createNotificationGroupCmd.Flags().StringSlice("instance-ids", nil, "Instance IDs to monitor")
	createNotificationGroupCmd.Flags().StringSlice("user-group-ids", nil, "User group IDs to notify")

	// Log commands
	rdsMySQLCmd.AddCommand(logCmd)
	logCmd.AddCommand(listLogsCmd)

	// Metrics commands
	rdsMySQLCmd.AddCommand(metricsCmd)
	metricsCmd.AddCommand(listMetricsCmd)
	metricsCmd.AddCommand(getMetricStatsCmd)

	getMetricStatsCmd.Flags().String("from", "", "Start time (YYYY-MM-DD HH:MM)")
	getMetricStatsCmd.Flags().String("to", "", "End time (YYYY-MM-DD HH:MM)")
	getMetricStatsCmd.Flags().Int("interval", 0, "Aggregation interval in minutes")

	// Network commands
	rdsMySQLCmd.AddCommand(networkCmd)
	networkCmd.AddCommand(getNetworkInfoCmd)
	networkCmd.AddCommand(setPublicAccessCmd)

	setPublicAccessCmd.Flags().Bool("enable", false, "Enable public access")
	setPublicAccessCmd.Flags().Bool("disable", false, "Disable public access")

	// Storage resize
	rdsMySQLCmd.AddCommand(resizeStorageCmd)
	resizeStorageCmd.Flags().Int("size", 0, "New storage size in GB (required)")

	// Deletion protection
	rdsMySQLCmd.AddCommand(setDeletionProtectionCmd)
	setDeletionProtectionCmd.Flags().Bool("enable", false, "Enable deletion protection")
	setDeletionProtectionCmd.Flags().Bool("disable", false, "Disable deletion protection")
}

// ============================================================================
// Print Functions
// ============================================================================

func printDBUsers(result *mysql.ListDBUsersOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tHOST\tAUTH PLUGIN\tCREATED")
	for _, u := range result.DBUsers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			u.DBUserID, u.DBUserName, u.Host, u.AuthenticationPlugin, u.CreatedYmdt)
	}
	w.Flush()
}

func printSchemas(result *mysql.ListSchemasOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCREATED")
	for _, s := range result.DBSchemas {
		fmt.Fprintf(w, "%s\t%s\t%s\n", s.DBSchemaId, s.DBSchemaName, s.CreatedYmdt)
	}
	w.Flush()
}

func printNotificationGroups(result *mysql.ListNotificationGroupsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tTYPE\tENABLED")
	for _, ng := range result.NotificationGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\n",
			ng.NotificationGroupID, ng.NotificationGroupName, ng.NotificationType, ng.IsEnabled)
	}
	w.Flush()
}

func printNotificationGroupDetail(result *mysql.NotificationGroupOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	fmt.Printf("ID:           %s\n", result.NotificationGroup.NotificationGroupID)
	fmt.Printf("Name:         %s\n", result.NotificationGroup.NotificationGroupName)
	fmt.Printf("Type:         %s\n", result.NotificationGroup.NotificationType)
	fmt.Printf("Enabled:      %v\n", result.NotificationGroup.IsEnabled)
	fmt.Printf("Created:      %s\n", result.NotificationGroup.CreatedYmdt)
	fmt.Printf("Updated:      %s\n", result.NotificationGroup.UpdatedYmdt)
}

func printLogFiles(result *mysql.ListLogFilesOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "FILE NAME\tSIZE\tCREATED")
	for _, l := range result.LogFiles {
		fmt.Fprintf(w, "%s\t%d bytes\t%s\n",
			l.LogFileName, l.LogFileSize, l.CreatedYmdt)
	}
	w.Flush()
}

func printMetrics(result *mysql.ListMetricsOutput) {
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

func printMetricStats(result *mysql.MetricStatisticsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	for _, ms := range result.MetricStatistics {
		fmt.Printf("Metric: %s (%s)\n", ms.MeasureName, ms.Unit)
		fmt.Printf("  Data points: %d\n", len(ms.Values))
	}
}

func printNetworkInfo(result *mysql.NetworkInfoOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	fmt.Printf("Availability Zone:   %s\n", result.AvailabilityZone)
	fmt.Printf("Subnet ID:           %s\n", result.Subnet.SubnetID)
	fmt.Printf("Subnet Name:         %s\n", result.Subnet.SubnetName)
	fmt.Println("Endpoints:")
	for _, ep := range result.EndPoints {
		fmt.Printf("  - %s: %s (%s)\n", ep.EndPointType, ep.Domain, ep.IPAddress)
	}
}
