package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// Complete Notification Group Commands
// ============================================================================

var getNotificationGroupCmd = &cobra.Command{
	Use:   "get-notification-group",
	Short: "Get details of a notification group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("notification-group-id")
		if groupID == "" {
			exitWithError("--notification-group-id is required", nil)
		}

		client := newMySQLClient()
		result, err := client.GetNotificationGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to get notification group", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			fmt.Printf("Notification Group ID: %s\n", result.NotificationGroup.NotificationGroupID)
			fmt.Printf("Name: %s\n", result.NotificationGroup.NotificationGroupName)
			fmt.Printf("Enabled: %v\n", result.NotificationGroup.IsEnabled)
		}
	},
}

var updateNotificationGroupCmd = &cobra.Command{
	Use:   "update-notification-group",
	Short: "Update a notification group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("notification-group-id")
		name, _ := cmd.Flags().GetString("notification-group-name")
		enabled, _ := cmd.Flags().GetBool("enabled")

		if groupID == "" {
			exitWithError("--notification-group-id is required", nil)
		}

		client := newMySQLClient()
		req := &mysql.UpdateNotificationGroupRequest{
			NotificationGroupName: &name,
			IsEnabled:             &enabled,
		}

		_, err := client.UpdateNotificationGroup(context.Background(), groupID, req)
		if err != nil {
			exitWithError("failed to update notification group", err)
		}

		fmt.Printf("Notification group updated successfully\n")
	},
}

// ============================================================================
// Parameter Group Utility
// ============================================================================

var copyDBParameterGroupCmd = &cobra.Command{
	Use:   "copy-db-parameter-group",
	Short: "Copy a DB parameter group",
	Run: func(cmd *cobra.Command, args []string) {
		sourceGroupID, _ := cmd.Flags().GetString("source-parameter-group-id")
		targetGroupName, _ := cmd.Flags().GetString("target-parameter-group-name")

		if sourceGroupID == "" || targetGroupName == "" {
			exitWithError("--source-parameter-group-id and --target-parameter-group-name are required", nil)
		}

		client := newMySQLClient()
		req := &mysql.CopyParameterGroupRequest{
			ParameterGroupName: targetGroupName,
		}

		result, err := client.CopyParameterGroup(context.Background(), sourceGroupID, req)
		if err != nil {
			exitWithError("failed to copy parameter group", err)
		}

		fmt.Printf("Parameter group copied: %s\n", result.ParameterGroupID)
	},
}

var resetDBParameterGroupCmd = &cobra.Command{
	Use:   "reset-db-parameter-group",
	Short: "Reset a DB parameter group to default values",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("parameter-group-id")

		if groupID == "" {
			exitWithError("--parameter-group-id is required", nil)
		}

		client := newMySQLClient()
		_, err := client.ResetParameterGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to reset parameter group", err)
		}

		fmt.Printf("Parameter group reset successfully\n")
	},
}

// ============================================================================
// Backup Export
// ============================================================================

var exportBackupCmd = &cobra.Command{
	Use:   "export-backup-to-object-storage",
	Short: "Export a backup to object storage",
	Run: func(cmd *cobra.Command, args []string) {
		backupID, _ := cmd.Flags().GetString("backup-id")
		tenantID, _ := cmd.Flags().GetString("tenant-id")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		if backupID == "" || tenantID == "" || username == "" || password == "" {
			exitWithError("--backup-id, --tenant-id, --username, and --password are required", nil)
		}

		client := newMySQLClient()
		req := &mysql.ExportBackupRequest{
			TenantID: tenantID,
			Username: username,
			Password: password,
		}

		result, err := client.ExportBackup(context.Background(), backupID, req)
		if err != nil {
			exitWithError("failed to export backup", err)
		}

		fmt.Printf("Backup export initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Security Group Rules (Informational)
// ============================================================================

var getDBSecurityGroupCmd = &cobra.Command{
	Use:   "get-db-security-group",
	Short: "Get details of a DB security group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("db-security-group-id")
		if groupID == "" {
			exitWithError("--db-security-group-id is required", nil)
		}

		client := newMySQLClient()
		result, err := client.GetSecurityGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to get security group", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			fmt.Printf("Security Group ID: %s\n", result.DBSecurityGroup.DBSecurityGroupID)
			fmt.Printf("Name: %s\n", result.DBSecurityGroup.DBSecurityGroupName)
			fmt.Printf("Description: %s\n", result.DBSecurityGroup.Description)
			fmt.Printf("Rules: %d\n", len(result.DBSecurityGroup.Rules))
		}
	},
}

var getDBParameterGroupCmd = &cobra.Command{
	Use:   "get-db-parameter-group",
	Short: "Get details of a DB parameter group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("parameter-group-id")
		if groupID == "" {
			exitWithError("--parameter-group-id is required", nil)
		}

		client := newMySQLClient()
		result, err := client.GetParameterGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to get parameter group", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			fmt.Printf("Parameter Group ID: %s\n", result.ParameterGroup.ParameterGroupID)
			fmt.Printf("Name: %s\n", result.ParameterGroup.ParameterGroupName)
			fmt.Printf("Description: %s\n", result.ParameterGroup.Description)
		}
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	// Notification groups
	rdsMySQLCmd.AddCommand(getNotificationGroupCmd)
	rdsMySQLCmd.AddCommand(updateNotificationGroupCmd)

	getNotificationGroupCmd.Flags().String("notification-group-id", "", "Notification group ID (required)")

	updateNotificationGroupCmd.Flags().String("notification-group-id", "", "Notification group ID (required)")
	updateNotificationGroupCmd.Flags().String("notification-group-name", "", "New notification group name")
	updateNotificationGroupCmd.Flags().Bool("enabled", true, "Enable the notification group")

	// Parameter group utilities
	rdsMySQLCmd.AddCommand(copyDBParameterGroupCmd)
	rdsMySQLCmd.AddCommand(resetDBParameterGroupCmd)

	copyDBParameterGroupCmd.Flags().String("source-parameter-group-id", "", "Source parameter group ID (required)")
	copyDBParameterGroupCmd.Flags().String("target-parameter-group-name", "", "Target parameter group name (required)")

	resetDBParameterGroupCmd.Flags().String("parameter-group-id", "", "Parameter group ID (required)")

	// Backup export
	rdsMySQLCmd.AddCommand(exportBackupCmd)

	exportBackupCmd.Flags().String("backup-id", "", "Backup ID (required)")
	exportBackupCmd.Flags().String("tenant-id", "", "Tenant ID for object storage (required)")
	exportBackupCmd.Flags().String("username", "", "Object storage username (required)")
	exportBackupCmd.Flags().String("password", "", "Object storage password (required)")

	// Get commands
	rdsMySQLCmd.AddCommand(getDBSecurityGroupCmd)
	rdsMySQLCmd.AddCommand(getDBParameterGroupCmd)

	getDBSecurityGroupCmd.Flags().String("db-security-group-id", "", "Security group ID (required)")
	getDBParameterGroupCmd.Flags().String("parameter-group-id", "", "Parameter group ID (required)")
}
