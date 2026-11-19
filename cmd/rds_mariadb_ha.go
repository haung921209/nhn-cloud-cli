package cmd

import (
	"context"

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

		// Always set pingInterval - API requires it (default 10 if not specified)
		candidate, _ := cmd.Flags().GetString("db-instance-candidate-name")
		if candidate == "" {
			exitWithError("--db-instance-candidate-name is required", nil)
		}
		req := &mariadb.EnableHARequest{
			UseHighAvailability:     true,
			PingInterval:            &pingInterval,
			DBInstanceCandidateName: candidate,
		}

		if cmd.Flags().Changed("replication-mode") {
			req.ReplicationMode = replicationMode
		}

		result, err := client.EnableHA(context.Background(), dbInstanceID, req)
		if err != nil {
			exitWithError("failed to enable Multi-AZ", err)
		}

		if err := printWriteResult(result, "Multi-AZ enablement initiated.\nJob ID: %s\n", result.JobID); err != nil {
			exitWithError("failed to print result", err)
		}
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

		if err := printWriteResult(result, "Multi-AZ disablement initiated.\nJob ID: %s\n", result.JobID); err != nil {
			exitWithError("failed to print result", err)
		}
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

		if err := printWriteResult(result, "Multi-AZ paused.\nJob ID: %s\n", result.JobID); err != nil {
			exitWithError("failed to print result", err)
		}
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

		if err := printWriteResult(result, "Multi-AZ resumed.\nJob ID: %s\n", result.JobID); err != nil {
			exitWithError("failed to print result", err)
		}
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

		zone, _ := cmd.Flags().GetString("availability-zone")
		if zone == "" {
			exitWithError("--availability-zone is required", nil)
		}
		req := &mariadb.CreateReplicaRequest{
			DBInstanceName: replicaName,
			Network:        mariadb.ReplicaNetworkConfig{AvailabilityZone: zone},
		}

		result, err := client.CreateReplica(context.Background(), sourceID, req)
		if err != nil {
			exitWithError("failed to create read replica", err)
		}

		if err := printWriteResult(result, "Read replica creation initiated.\nJob ID: %s\n", result.JobID); err != nil {
			exitWithError("failed to print result", err)
		}
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

		if err := printWriteResult(result, "Read replica promotion initiated.\nJob ID: %s\n", result.JobID); err != nil {
			exitWithError("failed to print result", err)
		}
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
	enableMultiAZMariaDBCmd.Flags().String("db-instance-candidate-name", "", "Name of the standby (candidate) DB instance (required)")

	disableMultiAZMariaDBCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	pauseMultiAZMariaDBCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	resumeMultiAZMariaDBCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	// Replica Commands
	rdsMariaDBCmd.AddCommand(createMariaDBReadReplicaCmd)
	rdsMariaDBCmd.AddCommand(promoteMariaDBReadReplicaCmd)

	createMariaDBReadReplicaCmd.Flags().String("db-instance-identifier", "", "Source DB instance identifier (required)")
	createMariaDBReadReplicaCmd.Flags().String("replica-identifier", "", "Read replica identifier/name (required)")
	createMariaDBReadReplicaCmd.Flags().String("availability-zone", "", "Availability zone of the replica (required, e.g. kr-pub-a)")

	promoteMariaDBReadReplicaCmd.Flags().String("db-instance-identifier", "", "Read replica instance identifier (required)")
}
