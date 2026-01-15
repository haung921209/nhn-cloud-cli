package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-cli/pkg/auth"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mariadb"
	"github.com/spf13/cobra"
)

var rdsMariaDBCmd = &cobra.Command{
	Use:   "rds-mariadb",
	Short: "Manage RDS for MariaDB instances",
	Long:  `Manage RDS for MariaDB instances, backups, parameter groups, and more.`,
}

// ============================================================================
// Instance Commands
// ============================================================================

var describeMariaDBInstancesCmd = &cobra.Command{
	Use:   "describe-db-instances",
	Short: "Describe MariaDB DB instances",
	Long: `Describes one or more MariaDB DB instances.
If --db-instance-identifier is specified, describes a specific instance.
Otherwise, describes all instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")

		if instanceID != "" {
			// Describe specific instance
			resolvedID, err := resolveMariaDBInstanceIdentifier(client, instanceID)
			if err != nil {
				exitWithError("failed to resolve instance identifier", err)
			}

			result, err := client.GetInstance(context.Background(), resolvedID)
			if err != nil {
				exitWithError("failed to describe instance", err)
			}
			mariadbPrintInstanceDetail(result)
		} else {
			// List all instances
			result, err := client.ListInstances(context.Background())
			if err != nil {
				exitWithError("failed to list instances", err)
			}
			mariadbPrintInstanceList(result)
		}
	},
}

// ============================================================================
// Helper Functions
// ============================================================================

func newMariaDBClient() *mariadb.Client {
	cfg, err := auth.GetMariaDBConfig()
	if err != nil {
		exitWithError("failed to load MariaDB credentials", err)
	}

	client, err := mariadb.NewClient(cfg)
	if err != nil {
		exitWithError("failed to create MariaDB client", err)
	}

	return client
}

// resolveMariaDBInstanceIdentifier resolves an instance identifier (name or ID) to an ID
func resolveMariaDBInstanceIdentifier(client *mariadb.Client, identifier string) (string, error) {
	// If identifier looks like a UUID, return it as-is
	if len(identifier) == 36 && identifier[8] == '-' && identifier[13] == '-' {
		return identifier, nil
	}

	// Otherwise, treat as name and look it up
	result, err := client.ListInstances(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to list instances: %w", err)
	}

	for _, inst := range result.DBInstances {
		if inst.DBInstanceName == identifier {
			return inst.DBInstanceID, nil
		}
	}

	return "", fmt.Errorf("instance not found: %s", identifier)
}

// ============================================================================
// Print Functions
// ============================================================================

func mariadbPrintInstanceList(result *mariadb.ListInstancesResponse) {
	if output == "json" {
		mariadbPrintJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	// Similar to MySQL, ListInstances API might not return Flavor info, so omitting it from table
	fmt.Fprintln(w, "DB_INSTANCE_ID\tNAME\tSTATUS\tVERSION")
	for _, inst := range result.DBInstances {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			inst.DBInstanceID,
			inst.DBInstanceName,
			inst.DBInstanceStatus,
			inst.DBVersion,
		)
	}
	w.Flush()
}

func mariadbPrintInstanceDetail(result *mariadb.GetInstanceResponse) {
	if output == "json" {
		mariadbPrintJSON(result)
		return
	}

	inst := result.DatabaseInstance
	fmt.Printf("DB Instance: %s\n", inst.DBInstanceName)
	fmt.Printf("  ID: %s\n", inst.DBInstanceID)
	// DBInstanceGroupID and DBInstanceType are not in MariaDB SDK

	fmt.Printf("  Status: %s\n", inst.DBInstanceStatus)
	if inst.ProgressStatus != "" {
		fmt.Printf("  Progress: %s\n", inst.ProgressStatus)
	}
	fmt.Printf("  Flavor: %s\n", inst.DBFlavorID)
	// DBFlavorName is optional
	if inst.DBFlavorName != "" {
		fmt.Printf("  Flavor Name: %s\n", inst.DBFlavorName)
	}
	fmt.Printf("  Version: %s\n", inst.DBVersion)
	fmt.Printf("  Port: %d\n", inst.DBPort)
	if inst.ParameterGroupID != "" {
		fmt.Printf("  Parameter Group: %s\n", inst.ParameterGroupID)
	}
	if len(inst.DBSecurityGroupIDs) > 0 {
		fmt.Printf("  Security Groups: %v\n", inst.DBSecurityGroupIDs)
	}
	if len(inst.NotificationGroupIDs) > 0 {
		fmt.Printf("  Notification Groups: %v\n", inst.NotificationGroupIDs)
	}

	// MariaDB SDK only has CreatedAt/UpdatedAt
	if inst.CreatedAt != "" {
		fmt.Printf("  Created: %s\n", inst.CreatedAt)
	}
	if inst.UpdatedAt != "" {
		fmt.Printf("  Updated: %s\n", inst.UpdatedAt)
	}
}

func mariadbPrintJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func init() {
	rootCmd.AddCommand(rdsMariaDBCmd)

	rdsMariaDBCmd.AddCommand(describeMariaDBInstancesCmd)
	describeMariaDBInstancesCmd.Flags().String("db-instance-identifier", "", "DB instance identifier")

	// create-db-instance
	rdsMariaDBCmd.AddCommand(createMariaDBInstanceCmd)
	createMariaDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createMariaDBInstanceCmd.Flags().String("db-flavor-id", "", "DB flavor ID (required)")
	createMariaDBInstanceCmd.Flags().String("engine-version", "", "Engine version (required)")
	createMariaDBInstanceCmd.Flags().String("master-username", "", "Master username (required)")
	createMariaDBInstanceCmd.Flags().String("master-user-password", "", "Master user password (required)")
	createMariaDBInstanceCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	createMariaDBInstanceCmd.Flags().String("availability-zone", "", "Availability zone (required, e.g. kr-pub-a)")
	createMariaDBInstanceCmd.Flags().String("db-parameter-group-id", "", "DB parameter group ID (required)")

	// Optional flags
	createMariaDBInstanceCmd.Flags().String("description", "", "Instance description")
	createMariaDBInstanceCmd.Flags().Int("allocated-storage", 20, "Allocated storage in GB")
	createMariaDBInstanceCmd.Flags().Int("port", 3306, "Database port")
	createMariaDBInstanceCmd.Flags().String("storage-type", "General SSD", "Storage type")
	createMariaDBInstanceCmd.Flags().StringSlice("db-security-group-ids", nil, "DB security group IDs")
	createMariaDBInstanceCmd.Flags().Bool("multi-az", false, "Enable multi-AZ deployment")
	createMariaDBInstanceCmd.Flags().Int("backup-retention-period", 0, "Backup retention period in days")
	createMariaDBInstanceCmd.Flags().String("backup-window", "00:00", "Backup window time (HH:MM)")

	// modify-db-instance
	rdsMariaDBCmd.AddCommand(modifyMariaDBInstanceCmd)
	modifyMariaDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	modifyMariaDBInstanceCmd.Flags().String("new-db-instance-identifier", "", "New DB instance identifier")
	modifyMariaDBInstanceCmd.Flags().String("db-flavor-id", "", "New DB flavor ID")
	modifyMariaDBInstanceCmd.Flags().StringSlice("db-security-group-ids", nil, "New DB security group IDs (comma-separated)")
	modifyMariaDBInstanceCmd.Flags().Int("port", 0, "New database port")

	// delete-db-instance
	rdsMariaDBCmd.AddCommand(deleteMariaDBInstanceCmd)
	deleteMariaDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
}

var modifyMariaDBInstanceCmd = &cobra.Command{
	Use:   "modify-db-instance",
	Short: "Modify a MariaDB DB instance",
	Long:  `Modifies settings for a MariaDB DB instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		dbInstanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		req := &mariadb.ModifyInstanceRequest{}
		hasChanges := false

		// Collect modifications
		if cmd.Flags().Changed("new-db-instance-identifier") {
			newID, _ := cmd.Flags().GetString("new-db-instance-identifier")
			req.DBInstanceName = &newID
			hasChanges = true
		}
		if cmd.Flags().Changed("db-flavor-id") {
			flavorID, _ := cmd.Flags().GetString("db-flavor-id")
			req.DBFlavorID = &flavorID
			hasChanges = true
		}
		if cmd.Flags().Changed("port") {
			port, _ := cmd.Flags().GetInt("port")
			req.DBPort = &port
			hasChanges = true
		}
		if cmd.Flags().Changed("db-security-group-ids") {
			sgs, _ := cmd.Flags().GetStringSlice("db-security-group-ids")
			req.DBSecurityGroupIDs = sgs
			hasChanges = true
		}

		if !hasChanges {
			exitWithError("at least one modification parameter required", nil)
		}

		result, err := client.ModifyInstance(context.Background(), dbInstanceID, req)
		if err != nil {
			exitWithError("failed to modify instance", err)
		}

		fmt.Printf("DB instance modification initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var deleteMariaDBInstanceCmd = &cobra.Command{
	Use:   "delete-db-instance",
	Short: "Delete a MariaDB DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		dbInstanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.DeleteInstance(context.Background(), dbInstanceID)
		if err != nil {
			exitWithError("failed to delete instance", err)
		}

		fmt.Printf("DB instance deletion initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

// getResolvedMariaDBInstanceID is a helper that gets and resolves instance ID from command flags
func getResolvedMariaDBInstanceID(cmd *cobra.Command, client *mariadb.Client) (string, error) {
	identifier, _ := cmd.Flags().GetString("db-instance-identifier")
	if identifier == "" {
		return "", fmt.Errorf("--db-instance-identifier is required")
	}
	return resolveMariaDBInstanceIdentifier(client, identifier)
}

var createMariaDBInstanceCmd = &cobra.Command{
	Use:   "create-db-instance",
	Short: "Create a new MariaDB DB instance",
	Long: `Creates a new MariaDB DB instance.

Example:
  nhncloud rds-mariadb create-db-instance \
    --db-instance-identifier mydb \
    --db-flavor-id <flavor-uuid> \
    --engine-version <version-uuid> \
    --master-username admin \
    --master-user-password SecurePass123 \
    --allocated-storage 20 \
    --subnet-id <subnet-uuid> \
    --availability-zone kr-pub-a`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		client := newMariaDBClient()

		// Required parameters
		dbInstanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		dbFlavorID, _ := cmd.Flags().GetString("db-flavor-id")
		engineVersion, _ := cmd.Flags().GetString("engine-version")
		masterUsername, _ := cmd.Flags().GetString("master-username")
		masterPassword, _ := cmd.Flags().GetString("master-user-password")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		availabilityZone, _ := cmd.Flags().GetString("availability-zone")
		parameterGroupID, _ := cmd.Flags().GetString("db-parameter-group-id")

		// Optional parameters
		description, _ := cmd.Flags().GetString("description")
		allocatedStorage, _ := cmd.Flags().GetInt("allocated-storage")
		port, _ := cmd.Flags().GetInt("port")
		storageType, _ := cmd.Flags().GetString("storage-type")
		securityGroupIDs, _ := cmd.Flags().GetStringSlice("db-security-group-ids")
		multiAZ, _ := cmd.Flags().GetBool("multi-az")
		backupRetentionPeriod, _ := cmd.Flags().GetInt("backup-retention-period")
		backupWindow, _ := cmd.Flags().GetString("backup-window")

		// Validation (basic checks, SDK handles most)
		if dbInstanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}
		if dbFlavorID == "" {
			exitWithError("--db-flavor-id is required", nil)
		}
		if engineVersion == "" {
			exitWithError("--engine-version is required", nil)
		}
		if masterUsername == "" {
			exitWithError("--master-username is required", nil)
		}
		if masterPassword == "" {
			exitWithError("--master-user-password is required", nil)
		}
		if subnetID == "" {
			exitWithError("--subnet-id is required", nil)
		}
		if availabilityZone == "" {
			exitWithError("--availability-zone is required", nil)
		}
		if parameterGroupID == "" {
			exitWithError("--db-parameter-group-id is required", nil)
		}

		// Defaults
		if allocatedStorage == 0 {
			allocatedStorage = 20
		}
		if port == 0 {
			port = 3306
		}
		if storageType == "" {
			storageType = "General SSD"
		}
		if backupWindow == "" {
			backupWindow = "00:00"
		}

		// Build request
		req := &mariadb.CreateInstanceRequest{
			DBInstanceName:     dbInstanceID,
			Description:        description,
			DBFlavorID:         dbFlavorID,
			DBVersion:          engineVersion,
			DBUserName:         masterUsername,
			DBPassword:         masterPassword,
			DBPort:             &port,
			ParameterGroupID:   parameterGroupID,
			DBSecurityGroupIDs: securityGroupIDs,
			Network: mariadb.CreateInstanceNetworkConfig{
				SubnetID:         subnetID,
				AvailabilityZone: availabilityZone,
			},
			Storage: mariadb.CreateInstanceStorageConfig{
				StorageType: storageType,
				StorageSize: allocatedStorage,
			},
			Backup: mariadb.CreateInstanceBackupConfig{
				BackupPeriod: backupRetentionPeriod,
				BackupSchedules: []mariadb.CreateInstanceBackupSchedule{
					{
						BackupWndBgnTime:  backupWindow,
						BackupWndDuration: "TWO_HOURS",
					},
				},
			},
		}

		if multiAZ {
			req.UseHighAvailability = &multiAZ
		}

		result, err := client.CreateInstance(ctx, req)
		if err != nil {
			exitWithError("failed to create instance", err)
		}

		fmt.Printf("DB instance creation initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
		fmt.Printf("\nTo wait for completion, run:\n")
		// Note: Waiter not implemented yet
		fmt.Printf("  nhncloud rds-mariadb describe-db-instances --db-instance-identifier %s\n", dbInstanceID)
	},
}
