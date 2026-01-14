package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// Backup Commands (AWS-style: snapshot terminology)
// ============================================================================

var describeDBSnapshotsCmd = &cobra.Command{
	Use:   "describe-db-snapshots",
	Short: "Describe MySQL DB snapshots (backups)",
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if dbInstanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.ListBackups(context.Background(), dbInstanceID)
		if err != nil {
			exitWithError("failed to list backups", err)
		}

		printBackupList(result)
	},
}

var createDBSnapshotCmd = &cobra.Command{
	Use:   "create-db-snapshot",
	Short: "Create a DB snapshot (manual backup)",
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		snapshotName, _ := cmd.Flags().GetString("db-snapshot-identifier")

		if dbInstanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}
		if snapshotName == "" {
			exitWithError("--db-snapshot-identifier is required", nil)
		}

		client := newMySQLClient()
		req := &mysql.CreateBackupRequest{
			BackupName: snapshotName,
		}

		result, err := client.CreateBackup(context.Background(), dbInstanceID, req)
		if err != nil {
			exitWithError("failed to create snapshot", err)
		}

		fmt.Printf("DB snapshot creation initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var deleteDBSnapshotCmd = &cobra.Command{
	Use:   "delete-db-snapshot",
	Short: "Delete a DB snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		snapshotID, _ := cmd.Flags().GetString("db-snapshot-identifier")
		if snapshotID == "" {
			exitWithError("--db-snapshot-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.DeleteBackup(context.Background(), snapshotID)
		if err != nil {
			exitWithError("failed to delete snapshot", err)
		}

		fmt.Printf("DB snapshot deletion initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var restoreDBInstanceFromSnapshotCmd = &cobra.Command{
	Use:   "restore-db-instance-from-snapshot",
	Short: "Restore a DB instance from a snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		snapshotID, _ := cmd.Flags().GetString("db-snapshot-identifier")
		newInstanceID, _ := cmd.Flags().GetString("db-instance-identifier")

		if snapshotID == "" {
			exitWithError("--db-snapshot-identifier is required", nil)
		}

		client := newMySQLClient()
		req := &mysql.RestoreBackupRequest{}
		if newInstanceID != "" {
			req.DBInstanceName = newInstanceID
		}

		result, err := client.RestoreBackup(context.Background(), snapshotID, req)
		if err != nil {
			exitWithError("failed to restore from snapshot", err)
		}

		fmt.Printf("DB instance restore initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func printBackupList(result *mysql.ListBackupsResponse) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "BACKUP_ID\tNAME\tSTATUS\tTYPE\tCREATED")
	for _, backup := range result.Backups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			backup.BackupID,
			backup.BackupName,
			backup.BackupStatus,
			backup.BackupType,
			backup.CreatedAt,
		)
	}
	w.Flush()
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	rdsMySQLCmd.AddCommand(describeDBSnapshotsCmd)
	rdsMySQLCmd.AddCommand(createDBSnapshotCmd)
	rdsMySQLCmd.AddCommand(deleteDBSnapshotCmd)
	rdsMySQLCmd.AddCommand(restoreDBInstanceFromSnapshotCmd)

	// describe-db-snapshots
	describeDBSnapshotsCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	// create-db-snapshot
	createDBSnapshotCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createDBSnapshotCmd.Flags().String("db-snapshot-identifier", "", "Snapshot identifier/name (required)")

	// delete-db-snapshot
	deleteDBSnapshotCmd.Flags().String("db-snapshot-identifier", "", "Snapshot ID to delete (required)")

	// restore-db-instance-from-snapshot
	restoreDBInstanceFromSnapshotCmd.Flags().String("db-snapshot-identifier", "", "Snapshot ID to restore from (required)")
	restoreDBInstanceFromSnapshotCmd.Flags().String("db-instance-identifier", "", "New instance identifier (optional)")
}
