package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// Notification Group Commands
// ============================================================================

var describeNotificationGroupsCmd = &cobra.Command{
	Use:   "describe-notification-groups",
	Short: "Describe notification groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListNotificationGroups(context.Background())
		if err != nil {
			exitWithError("failed to list notification groups", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			for _, ng := range result.NotificationGroups {
				fmt.Printf("%s: %s\n", ng.NotificationGroupID, ng.NotificationGroupName)
			}
		}
	},
}

var createNotificationGroupCmd = &cobra.Command{
	Use:   "create-notification-group",
	Short: "Create a notification group",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("notification-group-name")
		if name == "" {
			exitWithError("--notification-group-name is required", nil)
		}

		client := newMySQLClient()
		req := &mysql.CreateNotificationGroupRequest{
			NotificationGroupName: name,
		}

		result, err := client.CreateNotificationGroup(context.Background(), req)
		if err != nil {
			exitWithError("failed to create notification group", err)
		}

		fmt.Printf("Notification group created: %s\n", result.NotificationGroupID)
	},
}

var deleteNotificationGroupCmd = &cobra.Command{
	Use:   "delete-notification-group",
	Short: "Delete a notification group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("notification-group-id")
		if groupID == "" {
			exitWithError("--notification-group-id is required", nil)
		}

		client := newMySQLClient()
		_, err := client.DeleteNotificationGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to delete notification group", err)
		}

		fmt.Printf("Notification group deleted successfully\n")
	},
}

// ============================================================================
// Log Commands
// ============================================================================

var describeLogFilesCmd = &cobra.Command{
	Use:   "describe-log-files",
	Short: "Describe log files for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.ListLogFiles(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to list log files", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			for _, log := range result.LogFiles {
				fmt.Printf("%s - %s (%d bytes)\n", log.LogFileName, log.ModifiedAt, log.LogFileSize)
			}
		}
	},
}

// ============================================================================
// Metrics Commands
// ============================================================================

var describeMetricsCmd = &cobra.Command{
	Use:   "describe-metrics",
	Short: "Describe available metrics",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListMetrics(context.Background())
		if err != nil {
			exitWithError("failed to list metrics", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			for _, metric := range result.Metrics {
				fmt.Printf("%s: %s\n", metric.MetricName, metric.Unit)
			}
		}
	},
}

// ============================================================================
// Network Commands
// ============================================================================

var describeNetworkInfoCmd = &cobra.Command{
	Use:   "describe-network-info",
	Short: "Describe network information for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.GetNetworkInfo(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to get network info", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			fmt.Printf("Subnet ID: %s\n", result.NetworkInfo.Subnet.SubnetID)
			fmt.Printf("Availability Zone: %s\n", result.NetworkInfo.AvailabilityZone)
		}
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	// Notification groups
	rdsMySQLCmd.AddCommand(describeNotificationGroupsCmd)
	rdsMySQLCmd.AddCommand(createNotificationGroupCmd)
	rdsMySQLCmd.AddCommand(deleteNotificationGroupCmd)

	// Logs
	rdsMySQLCmd.AddCommand(describeLogFilesCmd)

	// Metrics
	rdsMySQLCmd.AddCommand(describeMetricsCmd)

	// Network
	rdsMySQLCmd.AddCommand(describeNetworkInfoCmd)

	// Notification group flags
	createNotificationGroupCmd.Flags().String("notification-group-name", "", "Notification group name (required)")
	deleteNotificationGroupCmd.Flags().String("notification-group-id", "", "Notification group ID (required)")

	// Log flags
	describeLogFilesCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	// Network flags
	describeNetworkInfoCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
}
