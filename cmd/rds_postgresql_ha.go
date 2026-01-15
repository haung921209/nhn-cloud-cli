package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
	"github.com/spf13/cobra"
)

// ============================================================================
// HA Commands
// ============================================================================

var enablePostgreSQLMultiAZCmd = &cobra.Command{
	Use:   "enable-multi-az",
	Short: "Enable High Availability (Multi-AZ)",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		pingInterval, _ := cmd.Flags().GetInt("ping-interval")

		// Always send pingInterval with default value
		req := &postgresql.EnableHARequest{
			UseHighAvailability: true,
			PingInterval:        &pingInterval,
		}

		result, err := client.EnableHA(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to enable Multi-AZ", err)
		}

		fmt.Printf("Multi-AZ enablement initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var disablePostgreSQLMultiAZCmd = &cobra.Command{
	Use:   "disable-multi-az",
	Short: "Disable High Availability (Multi-AZ)",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.DisableHA(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to disable Multi-AZ", err)
		}

		fmt.Printf("Multi-AZ disablement initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var pausePostgreSQLMultiAZCmd = &cobra.Command{
	Use:   "pause-multi-az",
	Short: "Pause Multi-AZ monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.PauseHA(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to pause Multi-AZ", err)
		}

		fmt.Printf("Multi-AZ paused.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var resumePostgreSQLMultiAZCmd = &cobra.Command{
	Use:   "resume-multi-az",
	Short: "Resume Multi-AZ monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ResumeHA(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to resume Multi-AZ", err)
		}

		fmt.Printf("Multi-AZ resumed.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Read Replica Commands
// ============================================================================

var createPostgreSQLReadReplicaCmd = &cobra.Command{
	Use:   "create-read-replica",
	Short: "Create a read replica",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		sourceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve source instance ID", err)
		}

		replicaName, _ := cmd.Flags().GetString("replica-identifier")
		if replicaName == "" {
			exitWithError("--replica-identifier is required", nil)
		}

		req := &postgresql.CreateReplicaRequest{
			DBInstanceName: replicaName,
		}

		result, err := client.CreateReplica(context.Background(), sourceID, req)
		if err != nil {
			exitWithError("failed to create read replica", err)
		}

		fmt.Printf("Read replica creation initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var promotePostgreSQLReadReplicaCmd = &cobra.Command{
	Use:   "promote-read-replica",
	Short: "Promote a read replica to standalone instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		replicaID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve replica instance ID", err)
		}

		result, err := client.PromoteReplica(context.Background(), replicaID)
		if err != nil {
			exitWithError("failed to promote read replica", err)
		}

		fmt.Printf("Read replica promotion initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

func init() {
	// HA Commands
	rdsPostgreSQLCmd.AddCommand(enablePostgreSQLMultiAZCmd)
	rdsPostgreSQLCmd.AddCommand(disablePostgreSQLMultiAZCmd)
	rdsPostgreSQLCmd.AddCommand(pausePostgreSQLMultiAZCmd)
	rdsPostgreSQLCmd.AddCommand(resumePostgreSQLMultiAZCmd)

	enablePostgreSQLMultiAZCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	enablePostgreSQLMultiAZCmd.Flags().Int("ping-interval", 10, "Ping interval in seconds (default: 10)")

	disablePostgreSQLMultiAZCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	pausePostgreSQLMultiAZCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	resumePostgreSQLMultiAZCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	// Replica Commands
	rdsPostgreSQLCmd.AddCommand(createPostgreSQLReadReplicaCmd)
	rdsPostgreSQLCmd.AddCommand(promotePostgreSQLReadReplicaCmd)

	createPostgreSQLReadReplicaCmd.Flags().String("db-instance-identifier", "", "Source DB instance identifier (required)")
	createPostgreSQLReadReplicaCmd.Flags().String("replica-identifier", "", "Read replica identifier/name (required)")

	promotePostgreSQLReadReplicaCmd.Flags().String("db-instance-identifier", "", "Read replica instance identifier (required)")
}
