package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/postgresql"
	"github.com/spf13/cobra"
)

var rdsPostgreSQLCmd = &cobra.Command{
	Use:     "rds-postgresql",
	Aliases: []string{"rds-pg"},
	Short:   "Manage RDS for PostgreSQL instances",
	Long:    `Manage RDS for PostgreSQL instances, backups, parameter groups, and more.`,
}

// Instance Commands
var pgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all PostgreSQL instances",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListInstances(context.Background())
		if err != nil {
			exitWithError("failed to list instances", err)
		}
		printPGInstances(result)
	},
}

var pgCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new PostgreSQL instance",
	Run: func(cmd *cobra.Command, args []string) {
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

		if name == "" || flavorID == "" || version == "" || userName == "" || password == "" || subnetID == "" || availabilityZone == "" {
			exitWithError("required flags: --name, --flavor-id, --version, --user-name, --password, --subnet-id, --availability-zone", nil)
		}

		if backupStartTime == "" {
			backupStartTime = "00:00"
		}

		input := &postgresql.CreateInstanceInput{
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
		}
		input.Network.SubnetID = subnetID
		input.Network.AvailabilityZone = availabilityZone
		input.Storage.StorageType = storageType
		input.Storage.StorageSize = storageSize
		input.Backup.BackupPeriod = backupPeriod
		_ = backupStartTime

		client := newPostgreSQLClient()
		result, err := client.CreateInstance(context.Background(), input)
		if err != nil {
			exitWithError("failed to create instance", err)
		}
		fmt.Printf("Instance creation initiated. Job ID: %s\n", result.JobID)
	},
}

var pgGetCmd = &cobra.Command{
	Use:   "get [instance-id]",
	Short: "Get details of a PostgreSQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.GetInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get instance", err)
		}
		printPGInstance(result)
	},
}

var pgDeleteCmd = &cobra.Command{
	Use:   "delete [instance-id]",
	Short: "Delete a PostgreSQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.DeleteInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete instance", err)
		}
		fmt.Printf("Instance deletion initiated. Job ID: %s\n", result.JobID)
	},
}

var pgStartCmd = &cobra.Command{
	Use:   "start [instance-id]",
	Short: "Start a stopped PostgreSQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.StartInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to start instance", err)
		}
		fmt.Printf("Instance start initiated. Job ID: %s\n", result.JobID)
	},
}

var pgStopCmd = &cobra.Command{
	Use:   "stop [instance-id]",
	Short: "Stop a running PostgreSQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.StopInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to stop instance", err)
		}
		fmt.Printf("Instance stop initiated. Job ID: %s\n", result.JobID)
	},
}

var pgRestartCmd = &cobra.Command{
	Use:   "restart [instance-id]",
	Short: "Restart a PostgreSQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		useFailover, _ := cmd.Flags().GetBool("use-failover")
		client := newPostgreSQLClient()
		result, err := client.RestartInstance(context.Background(), args[0], useFailover)
		if err != nil {
			exitWithError("failed to restart instance", err)
		}
		fmt.Printf("Instance restart initiated. Job ID: %s\n", result.JobID)
	},
}

// HA Commands
var pgHACmd = &cobra.Command{
	Use:   "ha",
	Short: "Manage High Availability for PostgreSQL instances",
}

var pgHAEnableCmd = &cobra.Command{
	Use:   "enable [instance-id]",
	Short: "Enable High Availability",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pingInterval, _ := cmd.Flags().GetInt("ping-interval")
		_ = pingInterval // PostgreSQL HA doesn't support ping-interval
		input := &postgresql.EnableHAInput{
			UseHighAvailability: true,
		}
		client := newPostgreSQLClient()
		result, err := client.EnableHighAvailability(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to enable HA", err)
		}
		fmt.Printf("HA enable initiated. Job ID: %s\n", result.JobID)
	},
}

var pgHADisableCmd = &cobra.Command{
	Use:   "disable [instance-id]",
	Short: "Disable High Availability",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.DisableHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to disable HA", err)
		}
		fmt.Printf("HA disable initiated. Job ID: %s\n", result.JobID)
	},
}

var pgHAPauseCmd = &cobra.Command{
	Use:   "pause [instance-id]",
	Short: "Pause High Availability",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.PauseHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to pause HA", err)
		}
		fmt.Printf("HA pause initiated. Job ID: %s\n", result.JobID)
	},
}

var pgHAResumeCmd = &cobra.Command{
	Use:   "resume [instance-id]",
	Short: "Resume High Availability",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ResumeHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to resume HA", err)
		}
		fmt.Printf("HA resume initiated. Job ID: %s\n", result.JobID)
	},
}

// Replica Commands
var pgReplicaCmd = &cobra.Command{
	Use:   "replica",
	Short: "Manage read replicas",
}

var pgCreateReplicaCmd = &cobra.Command{
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

		input := &postgresql.CreateReplicaInput{
			DBInstanceName: name,
			Description:    description,
			DBFlavorID:     flavorID,
			Network: &postgresql.ReplicaNetwork{
				AvailabilityZone: az,
			},
		}

		client := newPostgreSQLClient()
		result, err := client.CreateReplica(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create replica", err)
		}
		fmt.Printf("Replica creation initiated. Job ID: %s\n", result.JobID)
	},
}

var pgPromoteReplicaCmd = &cobra.Command{
	Use:   "promote [replica-instance-id]",
	Short: "Promote a read replica to standalone master",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.PromoteReplica(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to promote replica", err)
		}
		fmt.Printf("Replica promotion initiated. Job ID: %s\n", result.JobID)
	},
}

// Resource Commands
var pgFlavorsCmd = &cobra.Command{
	Use:   "flavors",
	Short: "List available PostgreSQL flavors",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListFlavors(context.Background())
		if err != nil {
			exitWithError("failed to list flavors", err)
		}
		printPGFlavors(result)
	},
}

var pgVersionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "List available PostgreSQL versions",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListVersions(context.Background())
		if err != nil {
			exitWithError("failed to list versions", err)
		}
		printPGVersions(result)
	},
}

var pgBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage backups",
}

var pgBackupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List backups",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("instance-id")
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")
		client := newPostgreSQLClient()
		result, err := client.ListBackups(context.Background(), instanceID, page, size)
		if err != nil {
			exitWithError("failed to list backups", err)
		}
		printPGBackups(result)
	},
}

var pgBackupCreateCmd = &cobra.Command{
	Use:   "create [instance-id]",
	Short: "Create a backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			exitWithError("--name is required", nil)
		}
		input := &postgresql.CreateBackupInput{BackupName: name}
		client := newPostgreSQLClient()
		result, err := client.CreateBackup(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create backup", err)
		}
		fmt.Printf("Backup creation initiated. Job ID: %s\n", result.JobID)
	},
}

var pgBackupDeleteCmd = &cobra.Command{
	Use:   "delete [backup-id]",
	Short: "Delete a backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.DeleteBackup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete backup", err)
		}
		fmt.Printf("Backup deletion initiated. Job ID: %s\n", result.JobID)
	},
}

// Database Commands
var pgDatabaseCmd = &cobra.Command{
	Use:     "database",
	Aliases: []string{"db"},
	Short:   "Manage databases",
}

var pgDatabaseListCmd = &cobra.Command{
	Use:   "list [instance-id]",
	Short: "List databases in an instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListDatabases(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list databases", err)
		}
		printPGDatabases(result)
	},
}

var pgDatabaseCreateCmd = &cobra.Command{
	Use:   "create [instance-id]",
	Short: "Create a database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			exitWithError("--name is required", nil)
		}
		input := &postgresql.CreateDatabaseInput{DatabaseName: name}
		client := newPostgreSQLClient()
		result, err := client.CreateDatabase(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create database", err)
		}
		fmt.Printf("Database creation initiated. Job ID: %s\n", result.JobID)
	},
}

var pgDatabaseDeleteCmd = &cobra.Command{
	Use:   "delete [instance-id] [database-id]",
	Short: "Delete a database",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.DeleteDatabase(context.Background(), args[0], args[1])
		if err != nil {
			exitWithError("failed to delete database", err)
		}
		fmt.Printf("Database deletion initiated. Job ID: %s\n", result.JobID)
	},
}

func init() {
	rootCmd.AddCommand(rdsPostgreSQLCmd)

	// Instance commands
	rdsPostgreSQLCmd.AddCommand(pgListCmd)
	rdsPostgreSQLCmd.AddCommand(pgCreateCmd)
	rdsPostgreSQLCmd.AddCommand(pgGetCmd)
	rdsPostgreSQLCmd.AddCommand(pgDeleteCmd)
	rdsPostgreSQLCmd.AddCommand(pgStartCmd)
	rdsPostgreSQLCmd.AddCommand(pgStopCmd)
	rdsPostgreSQLCmd.AddCommand(pgRestartCmd)
	pgRestartCmd.Flags().Bool("use-failover", false, "Use online failover during restart")

	pgCreateCmd.Flags().String("name", "", "Instance name (required)")
	pgCreateCmd.Flags().String("description", "", "Instance description")
	pgCreateCmd.Flags().String("flavor-id", "", "Flavor ID (required)")
	pgCreateCmd.Flags().String("version", "", "PostgreSQL version (required)")
	pgCreateCmd.Flags().String("user-name", "", "Admin user name (required)")
	pgCreateCmd.Flags().String("password", "", "Admin user password (required)")
	pgCreateCmd.Flags().Int("port", 5432, "PostgreSQL port")
	pgCreateCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	pgCreateCmd.Flags().String("availability-zone", "", "Availability zone (required, e.g. kr-pub-a)")
	pgCreateCmd.Flags().String("storage-type", "General SSD", "Storage type")
	pgCreateCmd.Flags().Int("storage-size", 20, "Storage size in GB")
	pgCreateCmd.Flags().String("parameter-group-id", "", "Parameter group ID")
	pgCreateCmd.Flags().StringSlice("security-group-ids", nil, "Security group IDs")
	pgCreateCmd.Flags().Bool("use-ha", false, "Enable High Availability")
	pgCreateCmd.Flags().Bool("deletion-protection", false, "Enable deletion protection")
	pgCreateCmd.Flags().Int("backup-period", 0, "Backup retention period (days)")
	pgCreateCmd.Flags().String("backup-start-time", "", "Backup start time (HH:MM)")

	// HA commands
	rdsPostgreSQLCmd.AddCommand(pgHACmd)
	pgHACmd.AddCommand(pgHAEnableCmd)
	pgHACmd.AddCommand(pgHADisableCmd)
	pgHACmd.AddCommand(pgHAPauseCmd)
	pgHACmd.AddCommand(pgHAResumeCmd)
	pgHAEnableCmd.Flags().Int("ping-interval", 3, "Ping interval in seconds")

	// Replica commands
	rdsPostgreSQLCmd.AddCommand(pgReplicaCmd)
	pgReplicaCmd.AddCommand(pgCreateReplicaCmd)
	pgReplicaCmd.AddCommand(pgPromoteReplicaCmd)
	pgCreateReplicaCmd.Flags().String("name", "", "Replica instance name (required)")
	pgCreateReplicaCmd.Flags().String("description", "", "Description")
	pgCreateReplicaCmd.Flags().String("flavor-id", "", "Flavor ID (optional, defaults to source)")
	pgCreateReplicaCmd.Flags().String("availability-zone", "", "Availability zone (e.g. kr-pub-a)")

	// Resource commands
	rdsPostgreSQLCmd.AddCommand(pgFlavorsCmd)
	rdsPostgreSQLCmd.AddCommand(pgVersionsCmd)

	// Backup commands
	rdsPostgreSQLCmd.AddCommand(pgBackupCmd)
	pgBackupCmd.AddCommand(pgBackupListCmd)
	pgBackupCmd.AddCommand(pgBackupCreateCmd)
	pgBackupCmd.AddCommand(pgBackupDeleteCmd)
	pgBackupListCmd.Flags().String("instance-id", "", "Filter by instance ID")
	pgBackupListCmd.Flags().Int("page", 0, "Page number")
	pgBackupListCmd.Flags().Int("size", 20, "Page size")
	pgBackupCreateCmd.Flags().String("name", "", "Backup name (required)")

	// Database commands
	rdsPostgreSQLCmd.AddCommand(pgDatabaseCmd)
	pgDatabaseCmd.AddCommand(pgDatabaseListCmd)
	pgDatabaseCmd.AddCommand(pgDatabaseCreateCmd)
	pgDatabaseCmd.AddCommand(pgDatabaseDeleteCmd)
	pgDatabaseCreateCmd.Flags().String("name", "", "Database name (required)")
}

func newPostgreSQLClient() *postgresql.Client {
	ak := getAccessKey()
	sk := getSecretKey()
	appKey := getPostgreSQLAppKey()
	if appKey == "" {
		exitWithError("appkey is required (set via --appkey, NHN_CLOUD_APPKEY, or ~/.nhncloud/credentials)", nil)
	}
	var creds credentials.Credentials
	if ak != "" && sk != "" {
		creds = credentials.NewStatic(ak, sk)
	}
	return postgresql.NewClient(getRegion(), appKey, creds, debug)
}

func printPGInstances(result *postgresql.ListInstancesOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tVERSION\tTYPE")
	for _, inst := range result.DBInstances {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", inst.DBInstanceID, inst.DBInstanceName, inst.DBInstanceStatus, inst.DBVersion, inst.DBInstanceType)
	}
	w.Flush()
}

func printPGInstance(result *postgresql.GetInstanceOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	fmt.Printf("ID:       %s\n", result.DBInstanceID)
	fmt.Printf("Name:     %s\n", result.DBInstanceName)
	fmt.Printf("Status:   %s\n", result.DBInstanceStatus)
	fmt.Printf("Version:  %s\n", result.DBVersion)
	fmt.Printf("Port:     %d\n", result.DBPort)
	fmt.Printf("Type:     %s\n", result.DBInstanceType)
	fmt.Printf("Created:  %s\n", result.CreatedYmdt)
}

func printPGFlavors(result *postgresql.ListFlavorsOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tVCPU\tRAM(MB)")
	for _, f := range result.DBFlavors {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", f.DBFlavorID, f.DBFlavorName, f.VCPUs, f.RAM)
	}
	w.Flush()
}

func printPGVersions(result *postgresql.ListVersionsOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VERSION\tNAME")
	for _, v := range result.DBVersions {
		fmt.Fprintf(w, "%s\t%s\n", v.DBVersionCode, v.Name)
	}
	w.Flush()
}

func printPGBackups(result *postgresql.ListBackupsOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tSIZE(BYTES)\tCREATED")
	for _, b := range result.Backups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", b.BackupID, b.BackupName, b.BackupStatus, b.BackupSize, b.CreatedYmdt)
	}
	w.Flush()
}

func printPGDatabases(result *postgresql.ListDatabasesOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCREATED")
	for _, db := range result.Databases {
		fmt.Fprintf(w, "%s\t%s\t%s\n", db.DatabaseID, db.DatabaseName, db.CreatedYmdt)
	}
	w.Flush()
}
