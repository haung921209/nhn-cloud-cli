package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mariadb"
	"github.com/spf13/cobra"
)

// ============================================================================
// HA Commands
// ============================================================================

var enableMultiAZMariaDBCmd = &cobra.Command{
	Use:   "enable-multi-az",
	Short: "Enable High Availability (Multi-AZ)",
	Long: `Enables Multi-AZ High Availability for a MariaDB instance.
This creates a standby replica in a different availability zone.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		dbInstanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		pingInterval, _ := cmd.Flags().GetInt("ping-interval")
		replicationMode, _ := cmd.Flags().GetString("replication-mode")

		req := &mariadb.EnableHARequest{
			UseHighAvailability: true,
		}

		if cmd.Flags().Changed("ping-interval") {
			req.PingInterval = &pingInterval
		}
		if cmd.Flags().Changed("replication-mode") {
			req.ReplicationMode = replicationMode
		}

		result, err := client.EnableHA(context.Background(), dbInstanceID, req)
		if err != nil {
			exitWithError("failed to enable Multi-AZ", err)
		}

		fmt.Printf("Multi-AZ enablement initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var disableMultiAZMariaDBCmd = &cobra.Command{
	Use:   "disable-multi-az",
	Short: "Disable High Availability (Multi-AZ)",
	Long:  `Disables Multi-AZ High Availability for a MariaDB instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		dbInstanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.DisableHA(context.Background(), dbInstanceID)
		if err != nil {
			exitWithError("failed to disable Multi-AZ", err)
		}

		fmt.Printf("Multi-AZ disablement initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var pauseMultiAZMariaDBCmd = &cobra.Command{
	Use:   "pause-multi-az",
	Short: "Pause Multi-AZ monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		dbInstanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.PauseHA(context.Background(), dbInstanceID)
		if err != nil {
			exitWithError("failed to pause Multi-AZ", err)
		}

		fmt.Printf("Multi-AZ paused.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var resumeMultiAZMariaDBCmd = &cobra.Command{
	Use:   "resume-multi-az",
	Short: "Resume Multi-AZ monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		dbInstanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ResumeHA(context.Background(), dbInstanceID)
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

var createMariaDBReadReplicaCmd = &cobra.Command{
	Use:   "create-read-replica",
	Short: "Create a read replica",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		sourceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve source instance ID", err)
		}

		replicaName, _ := cmd.Flags().GetString("replica-identifier")
		if replicaName == "" {
			exitWithError("--replica-identifier is required", nil)
		}

		req := &mariadb.CreateReplicaRequest{
			DBInstanceName: replicaName,
		}

		// Add other replica config options as needed based on SDK support
		// Currently SDK CreateReplicaRequest mainly takes DBInstanceName

		result, err := client.CreateReplica(context.Background(), sourceID, req)
		if err != nil {
			exitWithError("failed to create read replica", err)
		}

		fmt.Printf("Read replica creation initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var promoteMariaDBReadReplicaCmd = &cobra.Command{
	Use:   "promote-read-replica",
	Short: "Promote a read replica to standalone instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		// Here db-instance-identifier is the replica ID
		replicaID, err := getResolvedMariaDBInstanceID(cmd, client)
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
	rdsMariaDBCmd.AddCommand(enableMultiAZMariaDBCmd)
	rdsMariaDBCmd.AddCommand(disableMultiAZMariaDBCmd)
	rdsMariaDBCmd.AddCommand(pauseMultiAZMariaDBCmd)
	rdsMariaDBCmd.AddCommand(resumeMultiAZMariaDBCmd)

	enableMultiAZMariaDBCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	enableMultiAZMariaDBCmd.Flags().Int("ping-interval", 10, "Ping interval in seconds")
	enableMultiAZMariaDBCmd.Flags().String("replication-mode", "async", "Replication mode (async, semi-sync)")

	disableMultiAZMariaDBCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	pauseMultiAZMariaDBCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	resumeMultiAZMariaDBCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	// Replica Commands
	rdsMariaDBCmd.AddCommand(createMariaDBReadReplicaCmd)
	rdsMariaDBCmd.AddCommand(promoteMariaDBReadReplicaCmd)

	createMariaDBReadReplicaCmd.Flags().String("db-instance-identifier", "", "Source DB instance identifier (required)")
	createMariaDBReadReplicaCmd.Flags().String("replica-identifier", "", "Read replica identifier/name (required)")

	promoteMariaDBReadReplicaCmd.Flags().String("db-instance-identifier", "", "Read replica instance identifier (required)")
}
