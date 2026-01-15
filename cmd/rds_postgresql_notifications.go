package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
	"github.com/spf13/cobra"
)

// ============================================================================
// Notification Group Commands
// ============================================================================

var postgresqlDescribeNotificationGroupsCmd = &cobra.Command{
	Use:   "describe-notification-groups",
	Short: "Describe notification groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		groupID, _ := cmd.Flags().GetString("notification-group-id")

		if groupID != "" {
			result, err := client.GetNotificationGroup(context.Background(), groupID)
			if err != nil {
				exitWithError("failed to get notification group", err)
			}
			if output == "json" {
				postgresqlPrintJSON(result)
			} else {
				postgresqlPrintNotificationGroupDetail(result)
			}
		} else {
			result, err := client.ListNotificationGroups(context.Background())
			if err != nil {
				exitWithError("failed to list notification groups", err)
			}
			if output == "json" {
				postgresqlPrintJSON(result)
			} else {
				postgresqlPrintNotificationGroupList(result)
			}
		}
	},
}

var postgresqlCreateNotificationGroupCmd = &cobra.Command{
	Use:   "create-notification-group",
	Short: "Create a notification group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		name, _ := cmd.Flags().GetString("name")
		enabled, _ := cmd.Flags().GetBool("enabled")
		emails, _ := cmd.Flags().GetStringSlice("notify-email")
		sms, _ := cmd.Flags().GetStringSlice("notify-sms")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		req := &postgresql.CreateNotificationGroupRequest{
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

var postgresqlDeleteNotificationGroupCmd = &cobra.Command{
	Use:   "delete-notification-group",
	Short: "Delete a notification group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

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

var postgresqlDescribeLogsCmd = &cobra.Command{
	Use:   "describe-logs",
	Short: "Describe log files for an instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ListLogFiles(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to list log files", err)
		}

		if output == "json" {
			postgresqlPrintJSON(result)
		} else {
			postgresqlPrintLogFileList(result)
		}
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func postgresqlPrintNotificationGroupList(result *postgresql.ListNotificationGroupsResponse) {
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

func postgresqlPrintNotificationGroupDetail(result *postgresql.GetNotificationGroupResponse) {
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

func postgresqlPrintLogFileList(result *postgresql.ListLogFilesResponse) {
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
	rdsPostgreSQLCmd.AddCommand(postgresqlDescribeNotificationGroupsCmd)
	rdsPostgreSQLCmd.AddCommand(postgresqlCreateNotificationGroupCmd)
	rdsPostgreSQLCmd.AddCommand(postgresqlDeleteNotificationGroupCmd)

	postgresqlDescribeNotificationGroupsCmd.Flags().String("notification-group-id", "", "Specific notification group ID")

	postgresqlCreateNotificationGroupCmd.Flags().String("name", "", "Notification group name (required)")
	postgresqlCreateNotificationGroupCmd.Flags().Bool("enabled", true, "Enable notifications")
	postgresqlCreateNotificationGroupCmd.Flags().StringSlice("notify-email", nil, "Email addresses for notifications")
	postgresqlCreateNotificationGroupCmd.Flags().StringSlice("notify-sms", nil, "Phone numbers for SMS notifications")

	postgresqlDeleteNotificationGroupCmd.Flags().String("notification-group-id", "", "Notification group ID (required)")

	// Log commands
	rdsPostgreSQLCmd.AddCommand(postgresqlDescribeLogsCmd)
	postgresqlDescribeLogsCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
}
