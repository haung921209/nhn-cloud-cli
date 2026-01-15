package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mariadb"
	"github.com/spf13/cobra"
)

// ============================================================================
// Notification Group Commands
// ============================================================================

var mariadbDescribeNotificationGroupsCmd = &cobra.Command{
	Use:   "describe-notification-groups",
	Short: "Describe notification groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		groupID, _ := cmd.Flags().GetString("notification-group-id")

		if groupID != "" {
			result, err := client.GetNotificationGroup(context.Background(), groupID)
			if err != nil {
				exitWithError("failed to get notification group", err)
			}
			if output == "json" {
				mariadbPrintJSON(result)
			} else {
				mariadbPrintNotificationGroupDetail(result)
			}
		} else {
			result, err := client.ListNotificationGroups(context.Background())
			if err != nil {
				exitWithError("failed to list notification groups", err)
			}
			if output == "json" {
				mariadbPrintJSON(result)
			} else {
				mariadbPrintNotificationGroupList(result)
			}
		}
	},
}

var mariadbCreateNotificationGroupCmd = &cobra.Command{
	Use:   "create-notification-group",
	Short: "Create a notification group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		name, _ := cmd.Flags().GetString("name")
		enabled, _ := cmd.Flags().GetBool("enabled")
		emails, _ := cmd.Flags().GetStringSlice("notify-email")
		sms, _ := cmd.Flags().GetStringSlice("notify-sms")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		req := &mariadb.CreateNotificationGroupRequest{
			NotificationGroupName: name,
			IsEnabled:             enabled,
			NotifyEmail:           emails,
			NotifySms:             sms,
		}

		result, err := client.CreateNotificationGroup(context.Background(), req)
		if err != nil {
			exitWithError("failed to create notification group", err)
		}

		fmt.Printf("Notification group created.\n")
		fmt.Printf("ID: %s\n", result.NotificationGroupID)
	},
}

var mariadbDeleteNotificationGroupCmd = &cobra.Command{
	Use:   "delete-notification-group",
	Short: "Delete a notification group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		groupID, _ := cmd.Flags().GetString("notification-group-id")
		if groupID == "" {
			exitWithError("--notification-group-id is required", nil)
		}

		_, err := client.DeleteNotificationGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to delete notification group", err)
		}

		fmt.Printf("Notification group deleted successfully.\n")
	},
}

// ============================================================================
// Log Commands
// ============================================================================

var mariadbDescribeLogsCmd = &cobra.Command{
	Use:   "describe-logs",
	Short: "Describe log files for an instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		instanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ListLogFiles(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to list log files", err)
		}

		if output == "json" {
			mariadbPrintJSON(result)
		} else {
			mariadbPrintLogFileList(result)
		}
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func mariadbPrintNotificationGroupList(result *mariadb.ListNotificationGroupsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "GROUP_ID\tNAME\tENABLED")
	for _, ng := range result.NotificationGroups {
		fmt.Fprintf(w, "%s\t%s\t%v\n",
			ng.NotificationGroupID,
			ng.NotificationGroupName,
			ng.IsEnabled,
		)
	}
	w.Flush()
}

func mariadbPrintNotificationGroupDetail(result *mariadb.GetNotificationGroupResponse) {
	ng := result.NotificationGroup
	fmt.Printf("ID: %s\n", ng.NotificationGroupID)
	fmt.Printf("Name: %s\n", ng.NotificationGroupName)
	fmt.Printf("Enabled: %v\n", ng.IsEnabled)
	if len(ng.NotifyEmail) > 0 {
		fmt.Printf("Emails: %v\n", ng.NotifyEmail)
	}
	if len(ng.NotifySms) > 0 {
		fmt.Printf("SMS: %v\n", ng.NotifySms)
	}
}

func mariadbPrintLogFileList(result *mariadb.ListLogFilesResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "FILE_NAME\tSIZE\tMODIFIED")
	for _, log := range result.LogFiles {
		fmt.Fprintf(w, "%s\t%d\t%s\n",
			log.LogFileName,
			log.LogFileSize,
			log.ModifiedAt,
		)
	}
	w.Flush()
}

func init() {
	// Notification Group commands
	rdsMariaDBCmd.AddCommand(mariadbDescribeNotificationGroupsCmd)
	rdsMariaDBCmd.AddCommand(mariadbCreateNotificationGroupCmd)
	rdsMariaDBCmd.AddCommand(mariadbDeleteNotificationGroupCmd)

	mariadbDescribeNotificationGroupsCmd.Flags().String("notification-group-id", "", "Specific notification group ID")

	mariadbCreateNotificationGroupCmd.Flags().String("name", "", "Notification group name (required)")
	mariadbCreateNotificationGroupCmd.Flags().Bool("enabled", true, "Enable notifications")
	mariadbCreateNotificationGroupCmd.Flags().StringSlice("notify-email", nil, "Email addresses for notifications")
	mariadbCreateNotificationGroupCmd.Flags().StringSlice("notify-sms", nil, "Phone numbers for SMS notifications")

	mariadbDeleteNotificationGroupCmd.Flags().String("notification-group-id", "", "Notification group ID (required)")

	// Log commands
	rdsMariaDBCmd.AddCommand(mariadbDescribeLogsCmd)
	mariadbDescribeLogsCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
}
