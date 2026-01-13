package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-cli/pkg/interactive"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/mariadb"
	"github.com/spf13/cobra"
)

var rdsMariaDBCmd = &cobra.Command{
	Use:   "rds-mariadb",
	Short: "Manage RDS for MariaDB instances",
	Long:  `Manage RDS for MariaDB instances, backups, parameter groups, and more.`,
}

// Instance Commands
var mariadbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all MariaDB instances",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListInstances(context.Background())
		if err != nil {
			exitWithError("failed to list instances", err)
		}
		printMariaDBInstances(result)
	},
}

var mariadbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new MariaDB instance",
	Long: `Create a new MariaDB instance.

If required flags are not provided and running in a terminal,
interactive mode will be activated to guide you through the setup.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		client := newMariaDBClient()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		version, _ := cmd.Flags().GetString("version")
		userName, _ := cmd.Flags().GetString("user-name")
		password, _ := cmd.Flags().GetString("password")
		port, _ := cmd.Flags().GetInt("port")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		availabilityZone, _ := cmd.Flags().GetString("availability-zone")
		storageType, _ := cmd.Flags().GetString("storage-type")
		storageSize, _ := cmd.Flags().GetInt("storage-size")
		paramGroupID, _ := cmd.Flags().GetString("parameter-group-id")
		securityGroupIDs, _ := cmd.Flags().GetStringSlice("security-group-ids")
		useHA, _ := cmd.Flags().GetBool("use-ha")
		deletionProtection, _ := cmd.Flags().GetBool("deletion-protection")
		backupPeriod, _ := cmd.Flags().GetInt("backup-period")
		backupStartTime, _ := cmd.Flags().GetString("backup-start-time")

		missingRequired := name == "" || flavorID == "" || version == "" || userName == "" || password == "" || subnetID == "" || availabilityZone == "" || paramGroupID == "" || storageType == "" || storageSize == 0

		if missingRequired && interactive.CanRunInteractive() {
			azOptions := fetchAvailabilityZoneOptions(ctx)
			interactiveHandler := interactive.NewMariaDBInteractive(ctx, client, getRegion(), azOptions)
			interactiveHandler.SetDefinitions()
			pm := interactiveHandler.GetPromptManager()

			pm.SetProvidedValues(map[string]interface{}{
				"name":                name,
				"version":             version,
				"flavor-id":           flavorID,
				"user-name":           userName,
				"password":            password,
				"subnet-id":           subnetID,
				"availability-zone":   availabilityZone,
				"storage-type":        storageType,
				"storage-size":        storageSize,
				"port":                port,
				"parameter-group-id":  paramGroupID,
				"ha":                  useHA,
				"deletion-protection": deletionProtection,
				"backup-period":       backupPeriod,
				"backup-start-time":   backupStartTime,
			})

			values, err := pm.CollectValues()
			if err != nil {
				exitWithError("interactive mode failed", err)
			}

			pm.ShowSummary("MariaDB Instance Configuration")
			confirmed, err := pm.ConfirmExecution("Create this MariaDB instance?")
			if err != nil || !confirmed {
				fmt.Println("Operation cancelled.")
				return
			}

			name = values["name"].(string)
			version = values["version"].(string)
			flavorID = values["flavor-id"].(string)
			userName = values["user-name"].(string)
			password = values["password"].(string)
			subnetID = values["subnet-id"].(string)
			availabilityZone = values["availability-zone"].(string)
			if v, ok := values["storage-type"].(string); ok && v != "" {
				storageType = v
			}
			if v, ok := values["storage-size"].(int); ok && v > 0 {
				storageSize = v
			}
			if v, ok := values["port"].(int); ok && v > 0 {
				port = v
			}
			if v, ok := values["parameter-group-id"].(string); ok {
				paramGroupID = v
			}
			if v, ok := values["ha"].(bool); ok {
				useHA = v
			}
			if v, ok := values["deletion-protection"].(bool); ok {
				deletionProtection = v
			}
			if v, ok := values["backup-period"].(int); ok {
				backupPeriod = v
			}
			if v, ok := values["backup-start-time"].(string); ok && v != "" {
				backupStartTime = v
			}
		} else if missingRequired {
			exitWithError("required flags: --name, --flavor-id, --version, --user-name, --password, --subnet-id, --availability-zone, --parameter-group-id, --storage-type, --storage-size", nil)
		}

		if backupStartTime == "" {
			backupStartTime = "00:00"
		}

		input := &mariadb.CreateInstanceInput{
			DBInstanceName:        name,
			Description:           description,
			DBFlavorID:            flavorID,
			DBVersion:             version,
			DBUserName:            userName,
			DBPassword:            password,
			DBPort:                port,
			ParameterGroupID:      paramGroupID,
			DBSecurityGroupIDs:    securityGroupIDs,
			UseHighAvailability:   useHA,
			UseDeletionProtection: deletionProtection,
			Network: &mariadb.Network{
				SubnetID:         subnetID,
				AvailabilityZone: availabilityZone,
			},
			Storage: &mariadb.Storage{
				StorageType: storageType,
				StorageSize: storageSize,
			},
			Backup: &mariadb.BackupConfig{
				BackupPeriod: backupPeriod,
				BackupSchedules: []mariadb.BackupSchedule{
					{BackupWndBgnTime: backupStartTime, BackupWndDuration: "TWO_HOURS"},
				},
			},
		}

		result, err := client.CreateInstance(ctx, input)
		if err != nil {
			exitWithError("failed to create instance", err)
		}
		fmt.Printf("Instance creation initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbGetCmd = &cobra.Command{
	Use:   "get [instance-id]",
	Short: "Get details of a MariaDB instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.GetInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get instance", err)
		}
		printMariaDBInstance(result)
	},
}

var mariadbDeleteCmd = &cobra.Command{
	Use:   "delete [instance-id]",
	Short: "Delete a MariaDB instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.DeleteInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete instance", err)
		}
		fmt.Printf("Instance deletion initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbModifyCmd = &cobra.Command{
	Use:   "modify [instance-id]",
	Short: "Modify a MariaDB instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		port, _ := cmd.Flags().GetInt("port")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		paramGroupID, _ := cmd.Flags().GetString("parameter-group-id")
		securityGroupIDs, _ := cmd.Flags().GetStringSlice("security-group-ids")

		input := &mariadb.ModifyInstanceInput{}
		hasChanges := false

		if name != "" {
			input.DBInstanceName = name
			hasChanges = true
		}
		if description != "" {
			input.Description = description
			hasChanges = true
		}
		if port > 0 {
			if port < 3306 || port > 43306 {
				exitWithError("port must be between 3306 and 43306", nil)
			}
			input.DBPort = port
			hasChanges = true
		}
		if flavorID != "" {
			input.DBFlavorID = flavorID
			hasChanges = true
		}
		if paramGroupID != "" {
			input.ParameterGroupID = paramGroupID
			hasChanges = true
		}
		if len(securityGroupIDs) > 0 {
			input.DBSecurityGroupIDs = securityGroupIDs
			hasChanges = true
		}

		if !hasChanges {
			exitWithError("at least one modification flag is required", nil)
		}

		client := newMariaDBClient()
		result, err := client.ModifyInstance(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to modify instance", err)
		}
		printMariaDBInstance(result)
	},
}

var mariadbForceRestartCmd = &cobra.Command{
	Use:   "force-restart [instance-id]",
	Short: "Force restart a MariaDB instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ForceRestartInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to force restart instance", err)
		}
		fmt.Printf("Force restart initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbStartCmd = &cobra.Command{
	Use:   "start [instance-id]",
	Short: "Start a stopped MariaDB instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.StartInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to start instance", err)
		}
		fmt.Printf("Instance start initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbStopCmd = &cobra.Command{
	Use:   "stop [instance-id]",
	Short: "Stop a running MariaDB instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.StopInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to stop instance", err)
		}
		fmt.Printf("Instance stop initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbRestartCmd = &cobra.Command{
	Use:   "restart [instance-id]",
	Short: "Restart a MariaDB instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		useFailover, _ := cmd.Flags().GetBool("use-failover")
		executeBackup, _ := cmd.Flags().GetBool("execute-backup")
		client := newMariaDBClient()
		req := &mariadb.RestartInstanceRequest{
			UseOnlineFailover: useFailover,
			ExecuteBackup:     executeBackup,
		}
		result, err := client.RestartInstance(context.Background(), args[0], req)
		if err != nil {
			exitWithError("failed to restart instance", err)
		}
		fmt.Printf("Instance restart initiated. Job ID: %s\n", result.JobID)
	},
}

// HA Commands
var mariadbHACmd = &cobra.Command{
	Use:   "ha",
	Short: "Manage High Availability for MariaDB instances",
}

var mariadbHAEnableCmd = &cobra.Command{
	Use:   "enable [instance-id]",
	Short: "Enable High Availability",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pingInterval, _ := cmd.Flags().GetInt("ping-interval")
		input := &mariadb.EnableHAInput{
			UseHighAvailability: true,
			PingInterval:        pingInterval,
		}
		client := newMariaDBClient()
		result, err := client.EnableHighAvailability(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to enable HA", err)
		}
		fmt.Printf("HA enable initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbHADisableCmd = &cobra.Command{
	Use:   "disable [instance-id]",
	Short: "Disable High Availability",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.DisableHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to disable HA", err)
		}
		fmt.Printf("HA disable initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbHAPauseCmd = &cobra.Command{
	Use:   "pause [instance-id]",
	Short: "Pause High Availability",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.PauseHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to pause HA", err)
		}
		fmt.Printf("HA pause initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbHAResumeCmd = &cobra.Command{
	Use:   "resume [instance-id]",
	Short: "Resume High Availability",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ResumeHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to resume HA", err)
		}
		fmt.Printf("HA resume initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbHARepairCmd = &cobra.Command{
	Use:   "repair [instance-id]",
	Short: "Repair High Availability (recreate standby instance)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.RepairHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to repair HA", err)
		}
		fmt.Printf("HA repair initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbHASplitCmd = &cobra.Command{
	Use:   "split [instance-id]",
	Short: "Split High Availability (separate standby into independent instance)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.SplitHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to split HA", err)
		}
		fmt.Printf("HA split initiated. Job ID: %s\n", result.JobID)
	},
}

// Replica Commands
var mariadbReplicaCmd = &cobra.Command{
	Use:   "replica",
	Short: "Manage read replicas",
}

var mariadbCreateReplicaCmd = &cobra.Command{
	Use:   "create [source-instance-id]",
	Short: "Create a read replica from a master instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		az, _ := cmd.Flags().GetString("availability-zone")
		port, _ := cmd.Flags().GetInt("port")

		if name == "" {
			exitWithError("--name is required", nil)
		}
		if az == "" {
			exitWithError("--availability-zone is required", nil)
		}

		input := &mariadb.CreateReplicaInput{
			DBInstanceName: name,
			Description:    description,
			DBFlavorID:     flavorID,
			DBPort:         port,
			Network: &mariadb.ReplicaNetwork{
				AvailabilityZone: az,
			},
		}

		client := newMariaDBClient()
		result, err := client.CreateReplica(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create replica", err)
		}
		fmt.Printf("Replica creation initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbPromoteReplicaCmd = &cobra.Command{
	Use:   "promote [replica-instance-id]",
	Short: "Promote a read replica to standalone master",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.PromoteReplica(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to promote replica", err)
		}
		fmt.Printf("Replica promotion initiated. Job ID: %s\n", result.JobID)
	},
}

// Resource Commands
var mariadbFlavorsCmd = &cobra.Command{
	Use:   "flavors",
	Short: "List available MariaDB flavors",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListFlavors(context.Background())
		if err != nil {
			exitWithError("failed to list flavors", err)
		}
		printMariaDBFlavors(result)
	},
}

var mariadbVersionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "List available MariaDB versions",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListVersions(context.Background())
		if err != nil {
			exitWithError("failed to list versions", err)
		}
		printMariaDBVersions(result)
	},
}

var mariadbBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage backups",
}

var mariadbBackupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List backups",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("instance-id")
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")
		client := newMariaDBClient()
		result, err := client.ListBackups(context.Background(), instanceID, "", page, size)
		if err != nil {
			exitWithError("failed to list backups", err)
		}
		printMariaDBBackups(result)
	},
}

var mariadbBackupCreateCmd = &cobra.Command{
	Use:   "create [instance-id]",
	Short: "Create a backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			exitWithError("--name is required", nil)
		}
		input := &mariadb.CreateBackupInput{BackupName: name}
		client := newMariaDBClient()
		result, err := client.CreateBackup(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create backup", err)
		}
		fmt.Printf("Backup creation initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbBackupDeleteCmd = &cobra.Command{
	Use:   "delete [backup-id]",
	Short: "Delete a backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.DeleteBackup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete backup", err)
		}
		fmt.Printf("Backup deletion initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbBackupExportCmd = &cobra.Command{
	Use:   "export [backup-id]",
	Short: "Export a backup to object storage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tenantID, _ := cmd.Flags().GetString("tenant-id")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		targetContainer, _ := cmd.Flags().GetString("target-container")
		objectPath, _ := cmd.Flags().GetString("object-path")

		if tenantID == "" || username == "" || password == "" || targetContainer == "" || objectPath == "" {
			exitWithError("required flags: --tenant-id, --username, --password, --target-container, --object-path", nil)
		}

		input := &mariadb.ExportBackupInput{
			TenantID:        tenantID,
			Username:        username,
			Password:        password,
			TargetContainer: targetContainer,
			ObjectPath:      objectPath,
		}

		client := newMariaDBClient()
		result, err := client.ExportBackup(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to export backup", err)
		}
		fmt.Printf("Backup export initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbBackupRestoreCmd = &cobra.Command{
	Use:   "restore [backup-id]",
	Short: "Restore a backup to a new instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		az, _ := cmd.Flags().GetString("availability-zone")
		paramGroupID, _ := cmd.Flags().GetString("parameter-group-id")

		if name == "" || flavorID == "" || az == "" || paramGroupID == "" {
			exitWithError("required flags: --name, --flavor-id, --availability-zone, --parameter-group-id", nil)
		}

		input := &mariadb.RestoreBackupInput{
			DBInstanceName:   name,
			DBFlavorID:       flavorID,
			ParameterGroupID: paramGroupID,
			AvailabilityZone: az,
		}

		client := newMariaDBClient()
		result, err := client.RestoreBackup(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to restore backup", err)
		}
		fmt.Printf("Backup restore initiated. Job ID: %s\n", result.JobID)
	},
}

var mariadbBackupToObjectStorageCmd = &cobra.Command{
	Use:   "backup-to-object-storage [instance-id]",
	Short: "Backup directly to object storage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tenantID, _ := cmd.Flags().GetString("tenant-id")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		targetContainer, _ := cmd.Flags().GetString("target-container")
		objectPath, _ := cmd.Flags().GetString("object-path")

		if tenantID == "" || username == "" || password == "" || targetContainer == "" || objectPath == "" {
			exitWithError("required flags: --tenant-id, --username, --password, --target-container, --object-path", nil)
		}

		input := &mariadb.BackupToObjectStorageInput{
			TenantID:        tenantID,
			Username:        username,
			Password:        password,
			TargetContainer: targetContainer,
			ObjectPath:      objectPath,
		}

		client := newMariaDBClient()
		result, err := client.BackupToObjectStorage(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to backup to object storage", err)
		}
		fmt.Printf("Backup to object storage initiated. Job ID: %s\n", result.JobID)
	},
}

func init() {
	rootCmd.AddCommand(rdsMariaDBCmd)

	// Instance commands
	rdsMariaDBCmd.AddCommand(mariadbListCmd)
	rdsMariaDBCmd.AddCommand(mariadbCreateCmd)
	rdsMariaDBCmd.AddCommand(mariadbGetCmd)
	rdsMariaDBCmd.AddCommand(mariadbDeleteCmd)
	rdsMariaDBCmd.AddCommand(mariadbModifyCmd)
	rdsMariaDBCmd.AddCommand(mariadbStartCmd)
	rdsMariaDBCmd.AddCommand(mariadbStopCmd)
	rdsMariaDBCmd.AddCommand(mariadbRestartCmd)
	rdsMariaDBCmd.AddCommand(mariadbForceRestartCmd)

	mariadbRestartCmd.Flags().Bool("use-failover", false, "Use online failover during restart")
	mariadbRestartCmd.Flags().Bool("execute-backup", false, "Execute backup before restart")

	// Modify command flags
	mariadbModifyCmd.Flags().String("name", "", "New instance name")
	mariadbModifyCmd.Flags().String("description", "", "New description")
	mariadbModifyCmd.Flags().Int("port", 0, "New MariaDB port")
	mariadbModifyCmd.Flags().String("flavor-id", "", "New flavor ID")
	mariadbModifyCmd.Flags().String("parameter-group-id", "", "New parameter group ID")
	mariadbModifyCmd.Flags().StringSlice("security-group-ids", nil, "New security group IDs")

	mariadbCreateCmd.Flags().String("name", "", "Instance name (required)")
	mariadbCreateCmd.Flags().String("description", "", "Instance description")
	mariadbCreateCmd.Flags().String("flavor-id", "", "Flavor ID (required)")
	mariadbCreateCmd.Flags().String("version", "", "MariaDB version (required)")
	mariadbCreateCmd.Flags().String("user-name", "", "Admin user name (required)")
	mariadbCreateCmd.Flags().String("password", "", "Admin user password (required)")
	mariadbCreateCmd.Flags().Int("port", 3306, "MariaDB port")
	mariadbCreateCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	mariadbCreateCmd.Flags().String("availability-zone", "", "Availability zone (required, e.g. kr-pub-a)")
	mariadbCreateCmd.Flags().String("storage-type", "", "Storage type (from API)")
	mariadbCreateCmd.Flags().Int("storage-size", 20, "Storage size in GB")
	mariadbCreateCmd.Flags().String("parameter-group-id", "", "Parameter group ID")
	mariadbCreateCmd.Flags().StringSlice("security-group-ids", nil, "Security group IDs")
	mariadbCreateCmd.Flags().Bool("use-ha", false, "Enable High Availability")
	mariadbCreateCmd.Flags().Bool("deletion-protection", false, "Enable deletion protection")
	mariadbCreateCmd.Flags().Int("backup-period", 0, "Backup retention period (days)")
	mariadbCreateCmd.Flags().String("backup-start-time", "", "Backup start time (HH:MM)")

	// HA commands
	rdsMariaDBCmd.AddCommand(mariadbHACmd)
	mariadbHACmd.AddCommand(mariadbHAEnableCmd)
	mariadbHACmd.AddCommand(mariadbHADisableCmd)
	mariadbHACmd.AddCommand(mariadbHAPauseCmd)
	mariadbHACmd.AddCommand(mariadbHAResumeCmd)
	mariadbHACmd.AddCommand(mariadbHARepairCmd)
	mariadbHACmd.AddCommand(mariadbHASplitCmd)
	mariadbHAEnableCmd.Flags().Int("ping-interval", 3, "Ping interval in seconds")

	// Replica commands
	rdsMariaDBCmd.AddCommand(mariadbReplicaCmd)
	mariadbReplicaCmd.AddCommand(mariadbCreateReplicaCmd)
	mariadbReplicaCmd.AddCommand(mariadbPromoteReplicaCmd)
	mariadbCreateReplicaCmd.Flags().String("name", "", "Replica instance name (required)")
	mariadbCreateReplicaCmd.Flags().String("description", "", "Description")
	mariadbCreateReplicaCmd.Flags().String("flavor-id", "", "Flavor ID (optional, defaults to source)")
	mariadbCreateReplicaCmd.Flags().String("availability-zone", "", "Availability zone (e.g. kr-pub-a)")

	// Resource commands
	rdsMariaDBCmd.AddCommand(mariadbFlavorsCmd)
	rdsMariaDBCmd.AddCommand(mariadbVersionsCmd)

	// Backup commands
	rdsMariaDBCmd.AddCommand(mariadbBackupCmd)
	mariadbBackupCmd.AddCommand(mariadbBackupListCmd)
	mariadbBackupCmd.AddCommand(mariadbBackupCreateCmd)
	mariadbBackupCmd.AddCommand(mariadbBackupDeleteCmd)
	mariadbBackupCmd.AddCommand(mariadbBackupExportCmd)
	mariadbBackupCmd.AddCommand(mariadbBackupRestoreCmd)
	rdsMariaDBCmd.AddCommand(mariadbBackupToObjectStorageCmd)

	mariadbBackupListCmd.Flags().String("instance-id", "", "Filter by instance ID")
	mariadbBackupListCmd.Flags().Int("page", 0, "Page number")
	mariadbBackupListCmd.Flags().Int("size", 20, "Page size")
	mariadbBackupCreateCmd.Flags().String("name", "", "Backup name (required)")

	mariadbBackupExportCmd.Flags().String("tenant-id", "", "Object storage tenant ID (required)")
	mariadbBackupExportCmd.Flags().String("username", "", "Object storage username (required)")
	mariadbBackupExportCmd.Flags().String("password", "", "Object storage password (required)")
	mariadbBackupExportCmd.Flags().String("target-container", "", "Target container name (required)")
	mariadbBackupExportCmd.Flags().String("object-path", "", "Object path in container (required)")

	mariadbBackupRestoreCmd.Flags().String("name", "", "New instance name (required)")
	mariadbBackupRestoreCmd.Flags().String("flavor-id", "", "Flavor ID (required)")
	mariadbBackupRestoreCmd.Flags().String("availability-zone", "", "Availability zone (required)")
	mariadbBackupRestoreCmd.Flags().String("parameter-group-id", "", "Parameter group ID (required)")

	mariadbBackupToObjectStorageCmd.Flags().String("tenant-id", "", "Object storage tenant ID (required)")
	mariadbBackupToObjectStorageCmd.Flags().String("username", "", "Object storage username (required)")
	mariadbBackupToObjectStorageCmd.Flags().String("password", "", "Object storage password (required)")
	mariadbBackupToObjectStorageCmd.Flags().String("target-container", "", "Target container name (required)")
	mariadbBackupToObjectStorageCmd.Flags().String("object-path", "", "Object path in container (required)")
}

func newMariaDBClient() *mariadb.Client {
	ak := getAccessKey()
	sk := getSecretKey()
	appKey := getMariaDBAppKey()
	if appKey == "" {
		exitWithError("appkey is required (set via --appkey, NHN_CLOUD_APPKEY, or ~/.nhncloud/credentials)", nil)
	}
	var creds credentials.Credentials
	if ak != "" && sk != "" {
		creds = credentials.NewStatic(ak, sk)
	}
	return mariadb.NewClient(getRegion(), appKey, creds, debug)
}

func printMariaDBInstances(result *mariadb.ListInstancesOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tVERSION")
	for _, inst := range result.DBInstances {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", inst.DBInstanceID, inst.DBInstanceName, inst.DBInstanceStatus, inst.DBVersion)
	}
	w.Flush()
}

func printMariaDBInstance(result *mariadb.GetInstanceOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result)
		return
	}
	fmt.Printf("ID:       %s\n", result.DBInstanceID)
	fmt.Printf("Name:     %s\n", result.DBInstanceName)
	fmt.Printf("Status:   %s\n", result.DBInstanceStatus)
	fmt.Printf("Version:  %s\n", result.DBVersion)
	fmt.Printf("Storage:  %d GB (%s)\n", result.StorageSize, result.StorageType)
	fmt.Printf("Port:     %d\n", result.DBPort)
	fmt.Printf("Created:  %s\n", result.CreatedYmdt)
}

func printMariaDBFlavors(result *mariadb.ListFlavorsOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tVCPU\tRAM(MB)")
	for _, f := range result.DBFlavors {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", f.FlavorID, f.FlavorName, f.Vcpus, f.Ram)
	}
	w.Flush()
}

func printMariaDBVersions(result *mariadb.ListVersionsOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VERSION\tDISPLAY NAME")
	for _, v := range result.DBVersions {
		fmt.Fprintf(w, "%s\t%s\n", v.DBVersion, v.DBVersionName)
	}
	w.Flush()
}

func printMariaDBBackups(result *mariadb.ListBackupsOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tSIZE(MB)\tCREATED")
	for _, b := range result.Backups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", b.BackupID, b.BackupName, b.BackupStatus, b.BackupSize, b.CreatedYmdt)
	}
	w.Flush()
}
