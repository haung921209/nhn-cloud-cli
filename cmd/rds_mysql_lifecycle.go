package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// Lifecycle Commands
// ============================================================================

var startDBInstanceCmd = &cobra.Command{
	Use:   "start-db-instance",
	Short: "Start a stopped MySQL DB instance",
	Long: `Starts a stopped MySQL DB instance.

Example:
  nhncloud rds-mysql start-db-instance --db-instance-identifier mydb`,
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if dbInstanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.StartInstance(context.Background(), dbInstanceID)
		if err != nil {
			exitWithError("failed to start instance", err)
		}

		fmt.Printf("DB instance start initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var stopDBInstanceCmd = &cobra.Command{
	Use:   "stop-db-instance",
	Short: "Stop a running MySQL DB instance",
	Long: `Stops a running MySQL DB instance.

Example:
  nhncloud rds-mysql stop-db-instance --db-instance-identifier mydb`,
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if dbInstanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.StopInstance(context.Background(), dbInstanceID)
		if err != nil {
			exitWithError("failed to stop instance", err)
		}

		fmt.Printf("DB instance stop initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var rebootDBInstanceCmd = &cobra.Command{
	Use:   "reboot-db-instance",
	Short: "Reboot a MySQL DB instance",
	Long: `Reboots a MySQL DB instance.

Example:
  nhncloud rds-mysql reboot-db-instance --db-instance-identifier mydb`,
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if dbInstanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		// Get optional parameters
		useOnlineFailover, _ := cmd.Flags().GetBool("use-online-failover")
		executeBackup, _ := cmd.Flags().GetBool("execute-backup")

		client := newMySQLClient()
		req := &mysql.RestartInstanceRequest{}

		if useOnlineFailover {
			req.UseOnlineFailover = &useOnlineFailover
		}
		if executeBackup {
			req.ExecuteBackup = &executeBackup
		}

		result, err := client.RestartInstance(context.Background(), dbInstanceID, req)
		if err != nil {
			exitWithError("failed to reboot instance", err)
		}

		fmt.Printf("DB instance reboot initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var forceRebootDBInstanceCmd = &cobra.Command{
	Use:   "force-reboot-db-instance",
	Short: "Force reboot a MySQL DB instance",
	Long: `Force reboots a MySQL DB instance by forcefully restarting it.

Example:
  nhncloud rds-mysql force-reboot-db-instance --db-instance-identifier mydb`,
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if dbInstanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.ForceRestartInstance(context.Background(), dbInstanceID)
		if err != nil {
			exitWithError("failed to force reboot instance", err)
		}

		fmt.Printf("DB instance force reboot initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	rdsMySQLCmd.AddCommand(startDBInstanceCmd)
	rdsMySQLCmd.AddCommand(stopDBInstanceCmd)
	rdsMySQLCmd.AddCommand(rebootDBInstanceCmd)
	rdsMySQLCmd.AddCommand(forceRebootDBInstanceCmd)

	// start-db-instance flags
	startDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	// stop-db-instance flags
	stopDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	// reboot-db-instance flags
	rebootDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	rebootDBInstanceCmd.Flags().Bool("use-online-failover", false, "Use online failover for HA instances")
	rebootDBInstanceCmd.Flags().Bool("execute-backup", false, "Execute backup before reboot")

	// force-reboot-db-instance flags
	forceRebootDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
}
