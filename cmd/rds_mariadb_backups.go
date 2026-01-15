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
// Backup (Snapshot) Commands
// ============================================================================

var describeMariaDBSnapshotsCmd = &cobra.Command{
	Use:   "describe-db-snapshots",
	Short: "Describe MariaDB DB snapshots (backups)",
	Long: `Describes MariaDB DB snapshots (backups).
Current API limitation: --db-instance-identifier is REQUIRED. Global listing is not supported.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		instanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ListBackups(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to list backups", err)
		}

		mariadbPrintBackupList(result)
	},
}

var createMariaDBSnapshotCmd = &cobra.Command{
	Use:   "create-db-snapshot",
	Short: "Create a MariaDB DB snapshot",
	Long:  `Creates a manual backup (snapshot) for a MariaDB instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		instanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		snapshotID, _ := cmd.Flags().GetString("db-snapshot-identifier")
		if snapshotID == "" {
			exitWithError("--db-snapshot-identifier is required", nil)
		}

		req := &mariadb.CreateBackupRequest{
			BackupName: snapshotID,
		}

		result, err := client.CreateBackup(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to create snapshot", err)
		}

		fmt.Printf("Snapshot creation initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var deleteMariaDBSnapshotCmd = &cobra.Command{
	Use:   "delete-db-snapshot",
	Short: "Delete a MariaDB DB snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		// Note: DeleteBackup in SDK takes backupID (UUID)
		// But CLI users might want to providing name or ID.
		// Since ListBackups requires instanceID, resolving backup name to ID globally is hard.
		// We will assume the user provides the exact Backup ID (UUID) for now, similar to MySQL CLI v2.0 constraint.
		snapshotID, _ := cmd.Flags().GetString("db-snapshot-identifier")
		if snapshotID == "" {
			exitWithError("--db-snapshot-identifier is required", nil)
		}

		result, err := client.DeleteBackup(context.Background(), snapshotID)
		if err != nil {
			exitWithError("failed to delete snapshot", err)
		}

		fmt.Printf("Snapshot deletion initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var restoreMariaDBInstanceFromSnapshotCmd = &cobra.Command{
	Use:   "restore-db-instance-from-db-snapshot",
	Short: "Restore a MariaDB DB instance from a snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		snapshotID, _ := cmd.Flags().GetString("db-snapshot-identifier")
		targetInstanceName, _ := cmd.Flags().GetString("db-instance-identifier")

		if snapshotID == "" {
			exitWithError("--db-snapshot-identifier is required", nil)
		}
		if targetInstanceName == "" {
			exitWithError("--db-instance-identifier is required (for new restored instance)", nil)
		}

		req := &mariadb.RestoreBackupRequest{
			DBInstanceName: targetInstanceName,
			// Additional restoration parameters can be added here if SDK supports them
		}

		result, err := client.RestoreBackup(context.Background(), snapshotID, req)
		if err != nil {
			exitWithError("failed to restore instance from snapshot", err)
		}

		fmt.Printf("Restoration initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func mariadbPrintBackupList(result *mariadb.ListBackupsResponse) {
	if output == "json" {
		mariadbPrintJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "BACKUP_ID\tNAME\tSTATUS\tCREATED_AT")
	for _, backup := range result.Backups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			backup.BackupID,
			backup.BackupName,
			backup.BackupStatus,
			backup.CreatedAt,
		)
	}
	w.Flush()
}

func init() {
	// describe-db-snapshots
	rdsMariaDBCmd.AddCommand(describeMariaDBSnapshotsCmd)
	describeMariaDBSnapshotsCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required for MariaDB)")
	describeMariaDBSnapshotsCmd.Flags().String("db-snapshot-identifier", "", "DB snapshot identifier (optional filter)")

	// create-db-snapshot
	rdsMariaDBCmd.AddCommand(createMariaDBSnapshotCmd)
	createMariaDBSnapshotCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createMariaDBSnapshotCmd.Flags().String("db-snapshot-identifier", "", "DB snapshot identifier/name (required)")

	// delete-db-snapshot
	rdsMariaDBCmd.AddCommand(deleteMariaDBSnapshotCmd)
	deleteMariaDBSnapshotCmd.Flags().String("db-snapshot-identifier", "", "DB snapshot identifier/ID (required)")

	// restore-db-instance-from-db-snapshot
	rdsMariaDBCmd.AddCommand(restoreMariaDBInstanceFromSnapshotCmd)
	restoreMariaDBInstanceFromSnapshotCmd.Flags().String("db-snapshot-identifier", "", "Source DB snapshot identifier/ID (required)")
	restoreMariaDBInstanceFromSnapshotCmd.Flags().String("db-instance-identifier", "", "Name for the new restored instance (required)")
}
