package cmd

import (
	"context"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// High Availability Commands
// ============================================================================

var enableMultiAZCmd = &cobra.Command{
	Use:   "enable-multi-az",
	Short: "Enable multi-AZ (HA) for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, err := getResolvedInstanceID(cmd, newMySQLClient())
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		pingInterval, _ := cmd.Flags().GetInt("ping-interval")
		candidate, _ := cmd.Flags().GetString("db-instance-candidate-name")
		if candidate == "" {
			exitWithError("--db-instance-candidate-name is required", nil)
		}

		client := newMySQLClient()
		req := &mysql.EnableHARequest{
			UseHighAvailability:     true,
			PingInterval:            &pingInterval,
			DBInstanceCandidateName: candidate,
		}
		result, err := client.EnableHA(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to enable multi-AZ", err)
		}

		printWriteResult(result, "Multi-AZ enabled.\nJob ID: %s\n", result.JobID)
	},
}

var disableMultiAZCmd = &cobra.Command{
	Use:   "disable-multi-az",
	Short: "Disable multi-AZ (HA) for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, err := getResolvedInstanceID(cmd, newMySQLClient())
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		client := newMySQLClient()
		result, err := client.DisableHA(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to disable multi-AZ", err)
		}

		printWriteResult(result, "Multi-AZ disabled.\nJob ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	rdsMySQLCmd.AddCommand(enableMultiAZCmd)
	rdsMySQLCmd.AddCommand(disableMultiAZCmd)

	// enable-multi-az flags
	enableMultiAZCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	enableMultiAZCmd.Flags().Int("ping-interval", 300, "Ping interval in seconds for HA monitoring (default: 300)")
	enableMultiAZCmd.Flags().String("db-instance-candidate-name", "", "Name of the standby (candidate) DB instance (required)")

	// disable-multi-az flags
	disableMultiAZCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
}
