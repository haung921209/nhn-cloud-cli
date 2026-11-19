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
	// --notify-email/--notify-sms are booleans now; reject a stray value
	// such as "--notify-email ops@example.com" instead of ignoring it.
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		name, _ := cmd.Flags().GetString("name")
		enabled, _ := cmd.Flags().GetBool("enabled")
		email, _ := cmd.Flags().GetBool("notify-email")
		sms, _ := cmd.Flags().GetBool("notify-sms")
		instanceIDs, _ := cmd.Flags().GetStringSlice("db-instance-ids")
		userGroupIDs, _ := cmd.Flags().GetStringSlice("user-group-ids")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		req := &postgresql.CreateNotificationGroupRequest{
			NotificationGroupName: name,
			IsEnabled:             enabled,
			NotifyEmail:           email,
			NotifySms:             sms,
			DBInstanceIDs:         instanceIDs,
			UserGroupIDs:          userGroupIDs,
		}

		result, err := client.CreateNotificationGroup(context.Background(), req)
		if err != nil {
			exitWithError("failed to create notification group", err)
		}

		if err := printWriteResult(result, "Notification group created.\nID: %s\n", result.NotificationGroupID); err != nil {
			exitWithError("failed to print result", err)
		}
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

		result, err := client.DeleteNotificationGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to delete notification group", err)
		}

		if err := printWriteResult(result, "Notification group deleted successfully.\n"); err != nil {
			exitWithError("failed to print result", err)
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
	fmt.Printf("Email: %v\n", ng.NotifyEmail)
	fmt.Printf("SMS: %v\n", ng.NotifySms)
}

func init() {
	// Notification Group commands
	rdsPostgreSQLCmd.AddCommand(postgresqlDescribeNotificationGroupsCmd)
	rdsPostgreSQLCmd.AddCommand(postgresqlCreateNotificationGroupCmd)
	rdsPostgreSQLCmd.AddCommand(postgresqlDeleteNotificationGroupCmd)

	postgresqlDescribeNotificationGroupsCmd.Flags().String("notification-group-id", "", "Specific notification group ID")

	postgresqlCreateNotificationGroupCmd.Flags().String("name", "", "Notification group name (required)")
	postgresqlCreateNotificationGroupCmd.Flags().Bool("enabled", true, "Enable notifications")
	postgresqlCreateNotificationGroupCmd.Flags().Bool("notify-email", true, "Notify by email")
	postgresqlCreateNotificationGroupCmd.Flags().Bool("notify-sms", true, "Notify by SMS")
	postgresqlCreateNotificationGroupCmd.Flags().StringSlice("db-instance-ids", nil, "DB instance IDs to watch")
	postgresqlCreateNotificationGroupCmd.Flags().StringSlice("user-group-ids", nil, "User group IDs to notify")

	postgresqlDeleteNotificationGroupCmd.Flags().String("notification-group-id", "", "Notification group ID (required)")
}
