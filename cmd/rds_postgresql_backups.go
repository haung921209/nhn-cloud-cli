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

		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")
		result, err := client.ListBackups(context.Background(), instanceID, page, size)
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

		if err := printWriteResult(result, "Backup creation initiated.\nJob ID: %s\n", result.JobID); err != nil {
			exitWithError("failed to print result", err)
		}
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

		result, err := client.DeleteBackup(context.Background(), backupID)
		if err != nil {
			exitWithError("failed to delete backup", err)
		}

		if err := printWriteResult(result, "Backup deleted successfully.\n"); err != nil {
			exitWithError("failed to print result", err)
		}
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

func init() {
	// Backup commands
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLBackupsCmd)
	rdsPostgreSQLCmd.AddCommand(createPostgreSQLBackupCmd)
	rdsPostgreSQLCmd.AddCommand(deletePostgreSQLBackupCmd)

	describePostgreSQLBackupsCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	describePostgreSQLBackupsCmd.Flags().Int("page", 1, "Page number (from 1)")
	describePostgreSQLBackupsCmd.Flags().Int("size", 20, "Page size (1-100)")

	createPostgreSQLBackupCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createPostgreSQLBackupCmd.Flags().String("db-snapshot-identifier", "", "Backup/snapshot name (required)")

	deletePostgreSQLBackupCmd.Flags().String("db-snapshot-identifier", "", "Backup UUID (required)")
}
