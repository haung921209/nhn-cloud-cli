package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/mysql"
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

var listInstancesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all MySQL instances",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListInstances(context.Background())
		if err != nil {
			exitWithError("failed to list instances", err)
		}
		printInstanceList(result)
	},
}

var getInstanceCmd = &cobra.Command{
	Use:   "get [instance-id]",
	Short: "Get details of a MySQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.GetInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get instance", err)
		}
		printInstanceDetail(result)
	},
}

var createInstanceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new MySQL instance",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		version, _ := cmd.Flags().GetString("version")
		userName, _ := cmd.Flags().GetString("user-name")
		password, _ := cmd.Flags().GetString("password")
		port, _ := cmd.Flags().GetInt("port")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		storageType, _ := cmd.Flags().GetString("storage-type")
		storageSize, _ := cmd.Flags().GetInt("storage-size")
		paramGroupID, _ := cmd.Flags().GetString("parameter-group-id")
		securityGroupIDs, _ := cmd.Flags().GetStringSlice("security-group-ids")
		useHA, _ := cmd.Flags().GetBool("use-ha")
		deletionProtection, _ := cmd.Flags().GetBool("deletion-protection")
		backupPeriod, _ := cmd.Flags().GetInt("backup-period")
		backupStartTime, _ := cmd.Flags().GetString("backup-start-time")

		// Validation
		if name == "" || flavorID == "" || version == "" || userName == "" || password == "" || subnetID == "" {
			exitWithError("required flags: --name, --flavor-id, --version, --user-name, --password, --subnet-id", nil)
		}

		input := &mysql.CreateInstanceInput{
			Name:                  name,
			FlavorID:              flavorID,
			Version:               version,
			UserName:              userName,
			Password:              password,
			Port:                  port,
			ParameterGroupID:      paramGroupID,
			SecurityGroupIDs:      securityGroupIDs,
			UseHighAvailability:   useHA,
			UseDeletionProtection: deletionProtection,
			Network: &mysql.NetworkConfig{
				SubnetID: subnetID,
			},
			Storage: &mysql.StorageConfig{
				StorageType: storageType,
				StorageSize: storageSize,
			},
		}

		if backupPeriod > 0 {
			input.Backup = &mysql.BackupConfig{
				BackupPeriod: backupPeriod,
			}
			if backupStartTime != "" {
				input.Backup.BackupSchedules = []mysql.BackupSchedule{
					{BackupWndBgnTime: backupStartTime, BackupWndDuration: "02:00"},
				}
			}
		}

		client := newMySQLClient()
		result, err := client.CreateInstance(context.Background(), input)
		if err != nil {
			exitWithError("failed to create instance", err)
		}
		fmt.Printf("Instance creation initiated. Job ID: %s\n", result.JobID)
	},
}

var modifyInstanceCmd = &cobra.Command{
	Use:   "modify [instance-id]",
	Short: "Modify a MySQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		port, _ := cmd.Flags().GetInt("port")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		paramGroupID, _ := cmd.Flags().GetString("parameter-group-id")
		securityGroupIDs, _ := cmd.Flags().GetStringSlice("security-group-ids")

		input := &mysql.ModifyInstanceInput{}
		hasChanges := false

		if name != "" {
			input.Name = name
			hasChanges = true
		}
		if description != "" {
			input.Description = description
			hasChanges = true
		}
		if port > 0 {
			input.Port = port
			hasChanges = true
		}
		if flavorID != "" {
			input.FlavorID = flavorID
			hasChanges = true
		}
		if paramGroupID != "" {
			input.ParameterGroupID = paramGroupID
			hasChanges = true
		}
		if len(securityGroupIDs) > 0 {
			input.SecurityGroupIDs = securityGroupIDs
			hasChanges = true
		}

		if !hasChanges {
			exitWithError("at least one modification flag is required", nil)
		}

		client := newMySQLClient()
		result, err := client.ModifyInstance(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to modify instance", err)
		}
		printInstanceDetail(result)
	},
}

var deleteInstanceCmd = &cobra.Command{
	Use:   "delete [instance-id]",
	Short: "Delete a MySQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.DeleteInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete instance", err)
		}
		fmt.Printf("Instance deletion initiated. Job ID: %s\n", result.JobID)
	},
}

var startInstanceCmd = &cobra.Command{
	Use:   "start [instance-id]",
	Short: "Start a stopped MySQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.StartInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to start instance", err)
		}
		fmt.Printf("Instance start initiated. Job ID: %s\n", result.JobID)
	},
}

var stopInstanceCmd = &cobra.Command{
	Use:   "stop [instance-id]",
	Short: "Stop a running MySQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.StopInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to stop instance", err)
		}
		fmt.Printf("Instance stop initiated. Job ID: %s\n", result.JobID)
	},
}

var restartInstanceCmd = &cobra.Command{
	Use:   "restart [instance-id]",
	Short: "Restart a MySQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		useFailover, _ := cmd.Flags().GetBool("use-failover")
		client := newMySQLClient()
		result, err := client.RestartInstance(context.Background(), args[0], useFailover)
		if err != nil {
			exitWithError("failed to restart instance", err)
		}
		fmt.Printf("Instance restart initiated. Job ID: %s\n", result.JobID)
	},
}

var forceRestartInstanceCmd = &cobra.Command{
	Use:   "force-restart [instance-id]",
	Short: "Force restart a MySQL instance (kills all connections)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ForceRestartInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to force restart instance", err)
		}
		fmt.Printf("Instance force restart initiated. Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// High Availability Commands
// ============================================================================

var haCmd = &cobra.Command{
	Use:   "ha",
	Short: "Manage High Availability for MySQL instances",
}

var haEnableCmd = &cobra.Command{
	Use:   "enable [instance-id]",
	Short: "Enable High Availability for an instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pingInterval, _ := cmd.Flags().GetInt("ping-interval")
		failoverWait, _ := cmd.Flags().GetInt("failover-wait")

		input := &mysql.EnableHAInput{
			UseHighAvailability:     true,
			PingInterval:            pingInterval,
			FailoverReplWaitingTime: failoverWait,
		}

		client := newMySQLClient()
		result, err := client.EnableHighAvailability(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to enable HA", err)
		}
		fmt.Printf("HA enable initiated. Job ID: %s\n", result.JobID)
	},
}

var haDisableCmd = &cobra.Command{
	Use:   "disable [instance-id]",
	Short: "Disable High Availability for an instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.DisableHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to disable HA", err)
		}
		fmt.Printf("HA disable initiated. Job ID: %s\n", result.JobID)
	},
}

var haPauseCmd = &cobra.Command{
	Use:   "pause [instance-id]",
	Short: "Pause High Availability monitoring",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.PauseHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to pause HA", err)
		}
		fmt.Printf("HA pause initiated. Job ID: %s\n", result.JobID)
	},
}

var haResumeCmd = &cobra.Command{
	Use:   "resume [instance-id]",
	Short: "Resume High Availability monitoring",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ResumeHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to resume HA", err)
		}
		fmt.Printf("HA resume initiated. Job ID: %s\n", result.JobID)
	},
}

var haRepairCmd = &cobra.Command{
	Use:   "repair [instance-id]",
	Short: "Repair High Availability replication",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.RepairHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to repair HA", err)
		}
		fmt.Printf("HA repair initiated. Job ID: %s\n", result.JobID)
	},
}

var haSplitCmd = &cobra.Command{
	Use:   "split [instance-id]",
	Short: "Split HA standby as independent instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.SplitHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to split HA", err)
		}
		fmt.Printf("HA split initiated. Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Replica Commands
// ============================================================================

var replicaCmd = &cobra.Command{
	Use:   "replica",
	Short: "Manage read replicas",
}

var createReplicaCmd = &cobra.Command{
	Use:   "create [source-instance-id]",
	Short: "Create a read replica from a master instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		az, _ := cmd.Flags().GetString("availability-zone")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		input := &mysql.CreateReplicaInput{
			DBInstanceName:   name,
			Description:      description,
			DBFlavorID:       flavorID,
			AvailabilityZone: az,
		}

		client := newMySQLClient()
		result, err := client.CreateReplica(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create replica", err)
		}
		fmt.Printf("Replica creation initiated. Job ID: %s\n", result.JobID)
	},
}

var promoteReplicaCmd = &cobra.Command{
	Use:   "promote [replica-instance-id]",
	Short: "Promote a read replica to standalone master",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.PromoteReplica(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to promote replica", err)
		}
		fmt.Printf("Replica promotion initiated. Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Resource Listing Commands
// ============================================================================

var listFlavorsCmd = &cobra.Command{
	Use:   "flavors",
	Short: "List available MySQL flavors",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListFlavors(context.Background())
		if err != nil {
			exitWithError("failed to list flavors", err)
		}
		printFlavors(result)
	},
}

var listVersionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "List available MySQL versions",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListVersions(context.Background())
		if err != nil {
			exitWithError("failed to list versions", err)
		}
		printVersions(result)
	},
}

var listStorageTypesCmd = &cobra.Command{
	Use:   "storage-types",
	Short: "List available storage types",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListStorageTypes(context.Background())
		if err != nil {
			exitWithError("failed to list storage types", err)
		}
		printStorageTypes(result)
	},
}

var listSubnetsCmd = &cobra.Command{
	Use:   "subnets",
	Short: "List available subnets for RDS",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListSubnets(context.Background())
		if err != nil {
			exitWithError("failed to list subnets", err)
		}
		printSubnets(result)
	},
}

// ============================================================================
// Backup Commands
// ============================================================================

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage backups",
}

var listBackupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List backups",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("instance-id")
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")

		client := newMySQLClient()
		result, err := client.ListBackups(context.Background(), instanceID, "", page, size)
		if err != nil {
			exitWithError("failed to list backups", err)
		}
		printBackups(result)
	},
}

var createBackupCmd = &cobra.Command{
	Use:   "create [instance-id]",
	Short: "Create a backup for an instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			exitWithError("--name is required", nil)
		}

		input := &mysql.CreateBackupInput{
			BackupName: name,
		}

		client := newMySQLClient()
		result, err := client.CreateBackup(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create backup", err)
		}
		fmt.Printf("Backup creation initiated. Job ID: %s\n", result.JobID)
	},
}

var restoreBackupCmd = &cobra.Command{
	Use:   "restore [backup-id]",
	Short: "Restore a backup to a new instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("instance-name")
		if name == "" {
			exitWithError("--instance-name is required", nil)
		}

		input := &mysql.RestoreBackupInput{
			DBInstanceName: name,
		}

		client := newMySQLClient()
		result, err := client.RestoreBackup(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to restore backup", err)
		}
		fmt.Printf("Backup restore initiated. Job ID: %s\n", result.JobID)
	},
}

var deleteBackupCmd = &cobra.Command{
	Use:   "delete [backup-id]",
	Short: "Delete a backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.DeleteBackup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete backup", err)
		}
		fmt.Printf("Backup deletion initiated. Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Instance Group Commands
// ============================================================================

var instanceGroupCmd = &cobra.Command{
	Use:   "instance-groups",
	Short: "Manage instance groups",
}

var listInstanceGroupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all instance groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListInstanceGroups(context.Background())
		if err != nil {
			exitWithError("failed to list instance groups", err)
		}
		printInstanceGroups(result)
	},
}

var getInstanceGroupCmd = &cobra.Command{
	Use:   "get [group-id]",
	Short: "Get details of an instance group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.GetInstanceGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get instance group", err)
		}
		printInstanceGroupDetail(result)
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	rootCmd.AddCommand(rdsMySQLCmd)

	// Instance commands
	rdsMySQLCmd.AddCommand(listInstancesCmd)
	rdsMySQLCmd.AddCommand(getInstanceCmd)
	rdsMySQLCmd.AddCommand(createInstanceCmd)
	rdsMySQLCmd.AddCommand(modifyInstanceCmd)
	rdsMySQLCmd.AddCommand(deleteInstanceCmd)
	rdsMySQLCmd.AddCommand(startInstanceCmd)
	rdsMySQLCmd.AddCommand(stopInstanceCmd)
	rdsMySQLCmd.AddCommand(restartInstanceCmd)
	rdsMySQLCmd.AddCommand(forceRestartInstanceCmd)

	// Create instance flags
	createInstanceCmd.Flags().String("name", "", "Instance name (required)")
	createInstanceCmd.Flags().String("flavor-id", "", "Flavor ID (required)")
	createInstanceCmd.Flags().String("version", "", "MySQL version (required)")
	createInstanceCmd.Flags().String("user-name", "", "Admin user name (required)")
	createInstanceCmd.Flags().String("password", "", "Admin user password (required)")
	createInstanceCmd.Flags().Int("port", 3306, "MySQL port")
	createInstanceCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	createInstanceCmd.Flags().String("storage-type", "SSD", "Storage type (SSD, HDD)")
	createInstanceCmd.Flags().Int("storage-size", 20, "Storage size in GB")
	createInstanceCmd.Flags().String("parameter-group-id", "", "Parameter group ID")
	createInstanceCmd.Flags().StringSlice("security-group-ids", nil, "Security group IDs")
	createInstanceCmd.Flags().Bool("use-ha", false, "Enable High Availability")
	createInstanceCmd.Flags().Bool("deletion-protection", false, "Enable deletion protection")
	createInstanceCmd.Flags().Int("backup-period", 0, "Backup retention period (days)")
	createInstanceCmd.Flags().String("backup-start-time", "", "Backup start time (HH:MM)")

	// Modify instance flags
	modifyInstanceCmd.Flags().String("name", "", "New instance name")
	modifyInstanceCmd.Flags().String("description", "", "New description")
	modifyInstanceCmd.Flags().Int("port", 0, "New MySQL port")
	modifyInstanceCmd.Flags().String("flavor-id", "", "New flavor ID")
	modifyInstanceCmd.Flags().String("parameter-group-id", "", "New parameter group ID")
	modifyInstanceCmd.Flags().StringSlice("security-group-ids", nil, "New security group IDs")

	// Restart flags
	restartInstanceCmd.Flags().Bool("use-failover", false, "Use online failover during restart (HA only)")

	// HA commands
	rdsMySQLCmd.AddCommand(haCmd)
	haCmd.AddCommand(haEnableCmd)
	haCmd.AddCommand(haDisableCmd)
	haCmd.AddCommand(haPauseCmd)
	haCmd.AddCommand(haResumeCmd)
	haCmd.AddCommand(haRepairCmd)
	haCmd.AddCommand(haSplitCmd)

	haEnableCmd.Flags().Int("ping-interval", 3, "Ping interval in seconds")
	haEnableCmd.Flags().Int("failover-wait", 30, "Failover replication waiting time in seconds")

	// Replica commands
	rdsMySQLCmd.AddCommand(replicaCmd)
	replicaCmd.AddCommand(createReplicaCmd)
	replicaCmd.AddCommand(promoteReplicaCmd)

	createReplicaCmd.Flags().String("name", "", "Replica instance name (required)")
	createReplicaCmd.Flags().String("description", "", "Description")
	createReplicaCmd.Flags().String("flavor-id", "", "Flavor ID (defaults to source)")
	createReplicaCmd.Flags().String("availability-zone", "", "Availability zone")

	// Resource listing commands
	rdsMySQLCmd.AddCommand(listFlavorsCmd)
	rdsMySQLCmd.AddCommand(listVersionsCmd)
	rdsMySQLCmd.AddCommand(listStorageTypesCmd)
	rdsMySQLCmd.AddCommand(listSubnetsCmd)

	// Backup commands
	rdsMySQLCmd.AddCommand(backupCmd)
	backupCmd.AddCommand(listBackupsCmd)
	backupCmd.AddCommand(createBackupCmd)
	backupCmd.AddCommand(restoreBackupCmd)
	backupCmd.AddCommand(deleteBackupCmd)

	listBackupsCmd.Flags().String("instance-id", "", "Filter by instance ID")
	listBackupsCmd.Flags().Int("page", 0, "Page number")
	listBackupsCmd.Flags().Int("size", 20, "Page size")

	createBackupCmd.Flags().String("name", "", "Backup name (required)")
	restoreBackupCmd.Flags().String("instance-name", "", "New instance name (required)")

	// Instance group commands
	rdsMySQLCmd.AddCommand(instanceGroupCmd)
	instanceGroupCmd.AddCommand(listInstanceGroupsCmd)
	instanceGroupCmd.AddCommand(getInstanceGroupCmd)
}

// ============================================================================
// Helper Functions
// ============================================================================

func newMySQLClient() *mysql.Client {
	ak := os.Getenv("NHN_CLOUD_ACCESS_KEY")
	sk := os.Getenv("NHN_CLOUD_SECRET_KEY")

	appKey := getAppKey()
	if appKey == "" {
		exitWithError("appkey is required", nil)
	}

	var creds credentials.Credentials
	if ak != "" && sk != "" {
		creds = credentials.NewStatic(ak, sk)
	}

	return mysql.NewClient(getRegion(), appKey, creds, debug)
}

// ============================================================================
// Print Functions
// ============================================================================

func printInstanceList(result *mysql.ListInstancesOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tFLAVOR\tVERSION\tHA")
	for _, inst := range result.Instances {
		ha := "No"
		if inst.UseHighAvailability {
			ha = "Yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			inst.ID, inst.Name, inst.Status, inst.FlavorID, inst.Version, ha)
	}
	w.Flush()
}

func printInstanceDetail(result *mysql.GetInstanceOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	fmt.Printf("ID:                   %s\n", result.ID)
	fmt.Printf("Name:                 %s\n", result.Name)
	fmt.Printf("Status:               %s\n", result.Status)
	fmt.Printf("Version:              %s\n", result.Version)
	fmt.Printf("Flavor:               %s\n", result.FlavorID)
	fmt.Printf("Storage:              %d GB (%s)\n", result.StorageSize, result.StorageType)
	fmt.Printf("Port:                 %d\n", result.Port)
	fmt.Printf("High Availability:    %v\n", result.UseHighAvailability)
	fmt.Printf("Deletion Protection:  %v\n", result.UseDeletionProtection)
	fmt.Printf("Created:              %s\n", result.CreatedAt)
	fmt.Printf("Updated:              %s\n", result.UpdatedAt)
}

func printFlavors(result *mysql.ListFlavorsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tVCPU\tRAM(MB)")
	for _, f := range result.Flavors {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", f.ID, f.Name, f.VCPUs, f.RAM)
	}
	w.Flush()
}

func printVersions(result *mysql.ListVersionsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VERSION\tDISPLAY NAME")
	for _, v := range result.Versions {
		fmt.Fprintf(w, "%s\t%s\n", v.DBVersion, v.DisplayName)
	}
	w.Flush()
}

func printStorageTypes(result *mysql.ListStorageTypesOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tMIN SIZE(GB)\tMAX SIZE(GB)")
	for _, s := range result.StorageTypes {
		fmt.Fprintf(w, "%s\t%d\t%d\n", s.StorageType, s.MinSize, s.MaxSize)
	}
	w.Flush()
}

func printSubnets(result *mysql.ListSubnetsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCIDR\tAZ\tVPC")
	for _, s := range result.Subnets {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			s.ID, s.SubnetName, s.SubnetCIDR, s.AvailabilityZone, s.VPCName)
	}
	w.Flush()
}

func printBackups(result *mysql.ListBackupsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tSIZE(MB)\tCREATED")
	for _, b := range result.Backups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
			b.ID, b.Name, b.Status, b.Size, b.CreatedAt)
	}
	w.Flush()
}

func printInstanceGroups(result *mysql.ListInstanceGroupsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tREPLICATION TYPE\tINSTANCES\tCREATED")
	for _, g := range result.InstanceGroups {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
			g.ID, g.ReplicationType, len(g.Instances), g.CreatedAt)
	}
	w.Flush()
}

func printInstanceGroupDetail(result *mysql.InstanceGroupOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	fmt.Printf("ID:               %s\n", result.ID)
	fmt.Printf("Replication Type: %s\n", result.ReplicationType)
	fmt.Printf("Created:          %s\n", result.CreatedAt)
	fmt.Printf("Updated:          %s\n", result.UpdatedAt)
	fmt.Println("\nInstances:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  ID\tNAME\tTYPE\tSTATUS")
	for _, inst := range result.Instances {
		fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n",
			inst.ID, inst.Name, inst.Type, inst.Status)
	}
	w.Flush()
}

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}
