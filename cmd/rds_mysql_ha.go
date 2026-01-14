package cmd

import (
	"context"
	"fmt"

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
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		req := &mysql.EnableHARequest{
			UseHighAvailability: true,
		}
		result, err := client.EnableHA(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to enable multi-AZ", err)
		}

		fmt.Printf("Multi-AZ enabled.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var disableMultiAZCmd = &cobra.Command{
	Use:   "disable-multi-az",
	Short: "Disable multi-AZ (HA) for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.DisableHA(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to disable multi-AZ", err)
		}

		fmt.Printf("Multi-AZ disabled.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
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

	// disable-multi-az flags
	disableMultiAZCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
}
