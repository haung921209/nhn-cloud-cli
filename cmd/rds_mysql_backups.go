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
		dbInstanceID, err := getResolvedInstanceID(cmd, newMySQLClient())
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
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
		dbInstanceID, err := getResolvedInstanceID(cmd, newMySQLClient())
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}
		snapshotName, _ := cmd.Flags().GetString("db-snapshot-identifier")

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
// Backup Configuration Commands
// ============================================================================

var describeDBBackupInfoCmd = &cobra.Command{
	Use:   "describe-db-backup-info",
	Short: "Describe backup configuration for an instance",
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, err := getResolvedInstanceID(cmd, newMySQLClient())
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		client := newMySQLClient()
		result, err := client.GetBackupInfo(context.Background(), dbInstanceID)
		if err != nil {
			exitWithError("failed to get backup info", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Backup Configuration for %s:\n", dbInstanceID)
		fmt.Printf("  Backup Period: %d days\n", result.BackupPeriod)
		fmt.Printf("  Use Backup Lock: %v\n", result.UseBackupLock)
		fmt.Printf("  Wait Timeout: %d seconds\n", result.FtwrlWaitTimeout)
		fmt.Printf("  Retry Count: %d\n", result.BackupRetryCount)

		if len(result.BackupSchedules) > 0 {
			fmt.Println("  Schedules:")
			for _, schedule := range result.BackupSchedules {
				fmt.Printf("    - Start: %s, Duration: %s\n", schedule.BackupWndBgnTime, schedule.BackupWndDuration)
			}
		}
	},
}

var modifyDBBackupInfoCmd = &cobra.Command{
	Use:   "modify-db-backup-info",
	Short: "Modify backup configuration (enable automatic backups)",
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, err := getResolvedInstanceID(cmd, newMySQLClient())
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		backupPeriod, _ := cmd.Flags().GetInt("backup-period")
		useBackupLock, _ := cmd.Flags().GetBool("use-backup-lock")

		// Optional schedule
		startTime, _ := cmd.Flags().GetString("backup-window-start")
		duration, _ := cmd.Flags().GetString("backup-window-duration")

		client := newMySQLClient()
		req := &mysql.ModifyBackupInfoRequest{
			BackupPeriod:  backupPeriod,
			UseBackupLock: &useBackupLock,
		}

		if startTime != "" && duration != "" {
			req.BackupSchedules = []mysql.BackupSchedule{
				{
					BackupWndBgnTime:  startTime,
					BackupWndDuration: duration,
				},
			}
		}

		result, err := client.ModifyBackupInfo(context.Background(), dbInstanceID, req)
		if err != nil {
			exitWithError("failed to modify backup info", err)
		}

		fmt.Printf("Backup info modification initiated.\n")
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
	rdsMySQLCmd.AddCommand(describeDBBackupInfoCmd)
	rdsMySQLCmd.AddCommand(modifyDBBackupInfoCmd)

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

	// describe-db-backup-info
	describeDBBackupInfoCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	// modify-db-backup-info
	modifyDBBackupInfoCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	modifyDBBackupInfoCmd.Flags().Int("backup-period", 1, "Backup retention period in days (0-730). Set to 0 to disable.")
	modifyDBBackupInfoCmd.Flags().Bool("use-backup-lock", true, "Use backup lock (required for HA, default: true)")
	modifyDBBackupInfoCmd.Flags().String("backup-window-start", "", "Backup window start time (HH:MM:SS)")
	modifyDBBackupInfoCmd.Flags().String("backup-window-duration", "", "Backup window duration (ONE_HOUR, TWO_HOURS, etc.)")
}
