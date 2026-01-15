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
// Backup Commands
// ============================================================================

var describePostgreSQLBackupsCmd = &cobra.Command{
	Use:   "describe-db-snapshots",
	Short: "Describe PostgreSQL backups/snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ListBackups(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to list backups", err)
		}

		if output == "json" {
			postgresqlPrintJSON(result)
		} else {
			postgresqlPrintBackupList(result)
		}
	},
}

var createPostgreSQLBackupCmd = &cobra.Command{
	Use:   "create-db-snapshot",
	Short: "Create a PostgreSQL backup/snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		snapshotName, _ := cmd.Flags().GetString("db-snapshot-identifier")
		if snapshotName == "" {
			exitWithError("--db-snapshot-identifier is required", nil)
		}

		req := &postgresql.CreateBackupRequest{
			BackupName: snapshotName,
		}

		result, err := client.CreateBackup(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to create backup", err)
		}

		fmt.Printf("Backup creation initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var deletePostgreSQLBackupCmd = &cobra.Command{
	Use:   "delete-db-snapshot",
	Short: "Delete a PostgreSQL backup/snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		backupID, _ := cmd.Flags().GetString("db-snapshot-identifier")
		if backupID == "" {
			exitWithError("--db-snapshot-identifier is required (backup UUID)", nil)
		}

		_, err := client.DeleteBackup(context.Background(), backupID)
		if err != nil {
			exitWithError("failed to delete backup", err)
		}

		fmt.Printf("Backup deleted successfully.\n")
	},
}

// ============================================================================
// Security Group Commands
// ============================================================================

var describePostgreSQLSecurityGroupsCmd = &cobra.Command{
	Use:   "describe-db-security-groups",
	Short: "Describe PostgreSQL security groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		result, err := client.ListSecurityGroups(context.Background())
		if err != nil {
			exitWithError("failed to list security groups", err)
		}

		if output == "json" {
			postgresqlPrintJSON(result)
		} else {
			postgresqlPrintSecurityGroupList(result)
		}
	},
}

var createPostgreSQLSecurityGroupCmd = &cobra.Command{
	Use:   "create-db-security-group",
	Short: "Create a PostgreSQL security group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		name, _ := cmd.Flags().GetString("db-security-group-name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			exitWithError("--db-security-group-name is required", nil)
		}

		req := &postgresql.CreateSecurityGroupRequest{
			DBSecurityGroupName: name,
			Description:         description,
		}

		result, err := client.CreateSecurityGroup(context.Background(), req)
		if err != nil {
			exitWithError("failed to create security group", err)
		}

		fmt.Printf("Security group created.\n")
		fmt.Printf("ID: %s\n", result.DBSecurityGroupID)
	},
}

var deletePostgreSQLSecurityGroupCmd = &cobra.Command{
	Use:   "delete-db-security-group",
	Short: "Delete a PostgreSQL security group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		sgID, _ := cmd.Flags().GetString("db-security-group-identifier")
		if sgID == "" {
			exitWithError("--db-security-group-identifier is required", nil)
		}

		_, err := client.DeleteSecurityGroup(context.Background(), sgID)
		if err != nil {
			exitWithError("failed to delete security group", err)
		}

		fmt.Printf("Security group deleted successfully.\n")
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func postgresqlPrintBackupList(result *postgresql.ListBackupsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "BACKUP_ID\tNAME\tSTATUS\tSIZE\tCREATED")
	for _, b := range result.Backups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
			b.BackupID,
			b.BackupName,
			b.BackupStatus,
			b.BackupSize,
			b.CreatedAt,
		)
	}
	w.Flush()
}

func postgresqlPrintSecurityGroupList(result *postgresql.ListSecurityGroupsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SG_ID\tNAME\tDESCRIPTION")
	for _, sg := range result.DBSecurityGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			sg.DBSecurityGroupID,
			sg.DBSecurityGroupName,
			sg.Description,
		)
	}
	w.Flush()
}

func init() {
	// Backup commands
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLBackupsCmd)
	rdsPostgreSQLCmd.AddCommand(createPostgreSQLBackupCmd)
	rdsPostgreSQLCmd.AddCommand(deletePostgreSQLBackupCmd)

	describePostgreSQLBackupsCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	createPostgreSQLBackupCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createPostgreSQLBackupCmd.Flags().String("db-snapshot-identifier", "", "Backup/snapshot name (required)")

	deletePostgreSQLBackupCmd.Flags().String("db-snapshot-identifier", "", "Backup UUID (required)")

	// Security Group commands
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLSecurityGroupsCmd)
	rdsPostgreSQLCmd.AddCommand(createPostgreSQLSecurityGroupCmd)
	rdsPostgreSQLCmd.AddCommand(deletePostgreSQLSecurityGroupCmd)

	createPostgreSQLSecurityGroupCmd.Flags().String("db-security-group-name", "", "Security group name (required)")
	createPostgreSQLSecurityGroupCmd.Flags().String("description", "", "Description")

	deletePostgreSQLSecurityGroupCmd.Flags().String("db-security-group-identifier", "", "Security group ID (required)")
}
