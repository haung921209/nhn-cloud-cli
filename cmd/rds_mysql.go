package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-cli/pkg/auth"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

var rdsMySQLCmd = &cobra.Command{
	Use:   "rds-mysql",
	Short: "Manage RDS for MySQL instances",
	Long:  `Manage RDS for MySQL instances, backups, parameter groups, and more.`,
}

// ============================================================================
// Instance Commands
// ============================================================================

var describeDBInstancesCmd = &cobra.Command{
	Use:   "describe-db-instances",
	Short: "Describe MySQL DB instances",
	Long: `Describes one or more MySQL DB instances.
If --db-instance-identifier is specified, describes a specific instance.
Otherwise, describes all instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")

		if instanceID != "" {
			// Describe specific instance
			resolvedID, err := resolveInstanceIdentifier(client, instanceID)
			if err != nil {
				exitWithError("failed to resolve instance identifier", err)
			}

			result, err := client.GetInstance(context.Background(), resolvedID)
			if err != nil {
				exitWithError("failed to describe instance", err)
			}
			printInstanceDetail(result)
		} else {
			// List all instances
			result, err := client.ListInstances(context.Background())
			if err != nil {
				exitWithError("failed to list instances", err)
			}
			printInstanceList(result)
		}
	},
}

var createDBInstanceCmd = &cobra.Command{
	Use:   "create-db-instance",
	Short: "Create a new MySQL DB instance",
	Long: `Creates a new MySQL DB instance.

Example:
  nhncloud rds-mysql create-db-instance \
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
		client := newMySQLClient()

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

		// Validation
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

		// Default values
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
		req := &mysql.CreateInstanceRequest{
			DBInstanceName:     dbInstanceID,
			Description:        description,
			DBFlavorID:         dbFlavorID,
			DBVersion:          engineVersion,
			DBUserName:         masterUsername,
			DBPassword:         masterPassword,
			DBPort:             &port,
			ParameterGroupID:   parameterGroupID,
			DBSecurityGroupIDs: securityGroupIDs,
			Network: mysql.CreateInstanceNetworkConfig{
				SubnetID:         subnetID,
				AvailabilityZone: availabilityZone,
			},
			Storage: mysql.CreateInstanceStorageConfig{
				StorageType: storageType,
				StorageSize: allocatedStorage,
			},
			Backup: mysql.CreateInstanceBackupConfig{
				BackupPeriod: backupRetentionPeriod,
				BackupSchedules: []mysql.CreateInstanceBackupSchedule{
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
		fmt.Printf("  nhncloud rds-mysql wait db-instance-available --db-instance-identifier %s\n", dbInstanceID)
	},
}

var modifyDBInstanceCmd = &cobra.Command{
	Use:   "modify-db-instance",
	Short: "Modify a MySQL DB instance",
	Long:  `Modifies settings for a MySQL DB instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, err := getResolvedInstanceID(cmd, newMySQLClient())
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		client := newMySQLClient()
		req := &mysql.ModifyInstanceRequest{}
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

var deleteDBInstanceCmd = &cobra.Command{
	Use:   "delete-db-instance",
	Short: "Delete a MySQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, err := getResolvedInstanceID(cmd, newMySQLClient())
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		client := newMySQLClient()
		result, err := client.DeleteInstance(context.Background(), dbInstanceID)
		if err != nil {
			exitWithError("failed to delete instance", err)
		}

		fmt.Printf("DB instance deletion initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Helper Functions
// ============================================================================

func newMySQLClient() *mysql.Client {
	cfg, err := auth.GetMySQLConfig()
	if err != nil {
		exitWithError("failed to load MySQL credentials", err)
	}

	client, err := mysql.NewClient(cfg)
	if err != nil {
		exitWithError("failed to create MySQL client", err)
	}

	return client
}

// ============================================================================
// Print Functions
// ============================================================================

func printInstanceList(result *mysql.ListInstancesResponse) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
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

func printInstanceDetail(result *mysql.GetInstanceResponse) {
	if output == "json" {
		printJSON(result)
		return
	}

	inst := result.DatabaseInstance
	fmt.Printf("DB Instance: %s\n", inst.DBInstanceName)
	fmt.Printf("  ID: %s\n", inst.DBInstanceID)
	if inst.DBInstanceGroupID != "" {
		fmt.Printf("  Group ID: %s\n", inst.DBInstanceGroupID)
	}
	if inst.DBInstanceType != "" {
		fmt.Printf("  Type: %s\n", inst.DBInstanceType)
	}
	fmt.Printf("  Status: %s\n", inst.DBInstanceStatus)
	if inst.ProgressStatus != "" {
		fmt.Printf("  Progress: %s\n", inst.ProgressStatus)
	}
	fmt.Printf("  Flavor: %s\n", inst.DBFlavorID)
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
	if inst.CreatedYmdt != "" {
		fmt.Printf("  Created: %s\n", inst.CreatedYmdt)
	} else if inst.CreatedAt != "" {
		fmt.Printf("  Created: %s\n", inst.CreatedAt)
	}
	if inst.UpdatedYmdt != "" {
		fmt.Printf("  Updated: %s\n", inst.UpdatedYmdt)
	} else if inst.UpdatedAt != "" {
		fmt.Printf("  Updated: %s\n", inst.UpdatedAt)
	}
}

func printJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	rootCmd.AddCommand(rdsMySQLCmd)

	// Instance commands
	rdsMySQLCmd.AddCommand(describeDBInstancesCmd)
	rdsMySQLCmd.AddCommand(createDBInstanceCmd)
	rdsMySQLCmd.AddCommand(modifyDBInstanceCmd)
	rdsMySQLCmd.AddCommand(deleteDBInstanceCmd)

	// describe-db-instances flags
	describeDBInstancesCmd.Flags().String("db-instance-identifier", "", "DB instance identifier")

	// create-db-instance flags (required)
	createDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createDBInstanceCmd.Flags().String("db-flavor-id", "", "DB flavor ID (required)")
	createDBInstanceCmd.Flags().String("engine-version", "", "Engine version (required)")
	createDBInstanceCmd.Flags().String("master-username", "", "Master username (required)")
	createDBInstanceCmd.Flags().String("master-user-password", "", "Master user password (required)")
	createDBInstanceCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	createDBInstanceCmd.Flags().String("availability-zone", "", "Availability zone (required, e.g. kr-pub-a)")
	createDBInstanceCmd.Flags().String("db-parameter-group-id", "", "DB parameter group ID (required)")

	// create-db-instance flags (optional)
	createDBInstanceCmd.Flags().String("description", "", "Instance description")
	createDBInstanceCmd.Flags().Int("allocated-storage", 20, "Allocated storage in GB")
	createDBInstanceCmd.Flags().Int("port", 3306, "Database port")
	createDBInstanceCmd.Flags().String("storage-type", "General SSD", "Storage type")
	createDBInstanceCmd.Flags().StringSlice("db-security-group-ids", nil, "DB security group IDs")
	createDBInstanceCmd.Flags().Bool("multi-az", false, "Enable multi-AZ deployment")
	createDBInstanceCmd.Flags().Int("backup-retention-period", 0, "Backup retention period in days")
	createDBInstanceCmd.Flags().String("backup-window", "00:00", "Backup window time (HH:MM)")

	// modify-db-instance flags
	modifyDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	modifyDBInstanceCmd.Flags().String("new-db-instance-identifier", "", "New DB instance identifier")
	modifyDBInstanceCmd.Flags().String("db-flavor-id", "", "New DB flavor ID")
	modifyDBInstanceCmd.Flags().StringSlice("db-security-group-ids", nil, "New DB security group IDs (comma-separated)")
	modifyDBInstanceCmd.Flags().Int("port", 0, "New database port")

	// delete-db-instance flags
	deleteDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
}
