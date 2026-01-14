package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// Read Replica Commands
// ============================================================================

var createReadReplicaCmd = &cobra.Command{
	Use:   "create-read-replica",
	Short: "Create a read replica for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		replicaName, _ := cmd.Flags().GetString("replica-name")

		if instanceID == "" || replicaName == "" {
			exitWithError("--db-instance-identifier and --replica-name are required", nil)
		}

		client := newMySQLClient()
		req := &mysql.CreateReplicaRequest{
			DBInstanceName: replicaName,
		}

		result, err := client.CreateReplica(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to create read replica", err)
		}

		fmt.Printf("Read replica creation initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var promoteReadReplicaCmd = &cobra.Command{
	Use:   "promote-read-replica",
	Short: "Promote a read replica to a standalone DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		replicaID, _ := cmd.Flags().GetString("db-instance-identifier")
		if replicaID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.PromoteReplica(context.Background(), replicaID)
		if err != nil {
			exitWithError("failed to promote read replica", err)
		}

		fmt.Printf("Read replica promotion initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Advanced HA Commands
// ============================================================================

var pauseMultiAZCmd = &cobra.Command{
	Use:   "pause-multi-az",
	Short: "Pause multi-AZ (HA) for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.PauseHA(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to pause HA", err)
		}

		fmt.Printf("HA pause initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var resumeMultiAZCmd = &cobra.Command{
	Use:   "resume-multi-az",
	Short: "Resume multi-AZ (HA) for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.ResumeHA(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to resume HA", err)
		}

		fmt.Printf("HA resume initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var repairMultiAZCmd = &cobra.Command{
	Use:   "repair-multi-az",
	Short: "Repair multi-AZ (HA) for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.RepairHA(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to repair HA", err)
		}

		fmt.Printf("HA repair initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var splitMultiAZCmd = &cobra.Command{
	Use:   "split-multi-az",
	Short: "Split multi-AZ (HA) into separate instances",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		client := newMySQLClient()
		result, err := client.SplitHA(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to split HA", err)
		}

		fmt.Printf("HA split initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Network and Storage Management
// ============================================================================

var modifyDBInstanceNetworkCmd = &cobra.Command{
	Use:   "modify-db-instance-network",
	Short: "Modify network settings for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		enablePublicAccess, _ := cmd.Flags().GetBool("enable-public-access")
		disablePublicAccess, _ := cmd.Flags().GetBool("disable-public-access")

		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		if enablePublicAccess && disablePublicAccess {
			exitWithError("cannot specify both --enable-public-access and --disable-public-access", nil)
		}

		usePublicAccess := enablePublicAccess

		client := newMySQLClient()
		req := &mysql.ModifyNetworkInfoRequest{
			UsePublicAccess: usePublicAccess,
		}

		result, err := client.ModifyNetworkInfo(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to modify network", err)
		}

		fmt.Printf("Network modification initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var modifyDBInstanceStorageCmd = &cobra.Command{
	Use:   "modify-db-instance-storage",
	Short: "Modify storage size for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		storageSize, _ := cmd.Flags().GetInt("allocated-storage")

		if instanceID == "" || storageSize <= 0 {
			exitWithError("--db-instance-identifier and --allocated-storage are required", nil)
		}

		client := newMySQLClient()
		req := &mysql.ModifyStorageInfoRequest{
			StorageSize: storageSize,
		}

		result, err := client.ModifyStorageInfo(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to modify storage", err)
		}

		fmt.Printf("Storage modification initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var modifyDBInstanceDeletionProtectionCmd = &cobra.Command{
	Use:   "modify-deletion-protection",
	Short: "Enable or disable deletion protection for a DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		enable, _ := cmd.Flags().GetBool("enable")
		disable, _ := cmd.Flags().GetBool("disable")

		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}

		if enable && disable {
			exitWithError("cannot specify both --enable and --disable", nil)
		}

		if !enable && !disable {
			exitWithError("must specify either --enable or --disable", nil)
		}

		useDeletionProtection := enable

		client := newMySQLClient()
		req := &mysql.ModifyDeletionProtectionRequest{
			UseDeletionProtection: useDeletionProtection,
		}

		_, err := client.ModifyDeletionProtection(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to modify deletion protection", err)
		}

		fmt.Printf("Deletion protection modified successfully\n")
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	// Read Replica commands
	rdsMySQLCmd.AddCommand(createReadReplicaCmd)
	rdsMySQLCmd.AddCommand(promoteReadReplicaCmd)

	createReadReplicaCmd.Flags().String("db-instance-identifier", "", "Source DB instance identifier (required)")
	createReadReplicaCmd.Flags().String("replica-name", "", "Name for the read replica (required)")

	promoteReadReplicaCmd.Flags().String("db-instance-identifier", "", "Read replica identifier (required)")

	// Advanced HA commands
	rdsMySQLCmd.AddCommand(pauseMultiAZCmd)
	rdsMySQLCmd.AddCommand(resumeMultiAZCmd)
	rdsMySQLCmd.AddCommand(repairMultiAZCmd)
	rdsMySQLCmd.AddCommand(splitMultiAZCmd)

	pauseMultiAZCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	resumeMultiAZCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	repairMultiAZCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	splitMultiAZCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	// Network and Storage commands
	rdsMySQLCmd.AddCommand(modifyDBInstanceNetworkCmd)
	rdsMySQLCmd.AddCommand(modifyDBInstanceStorageCmd)
	rdsMySQLCmd.AddCommand(modifyDBInstanceDeletionProtectionCmd)

	modifyDBInstanceNetworkCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	modifyDBInstanceNetworkCmd.Flags().Bool("enable-public-access", false, "Enable public access")
	modifyDBInstanceNetworkCmd.Flags().Bool("disable-public-access", false, "Disable public access")

	modifyDBInstanceStorageCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	modifyDBInstanceStorageCmd.Flags().Int("allocated-storage", 0, "New storage size in GB (required)")

	modifyDBInstanceDeletionProtectionCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	modifyDBInstanceDeletionProtectionCmd.Flags().Bool("enable", false, "Enable deletion protection")
	modifyDBInstanceDeletionProtectionCmd.Flags().Bool("disable", false, "Disable deletion protection")
}
