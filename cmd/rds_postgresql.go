package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-cli/pkg/interactive"
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
	Long: `Create a new PostgreSQL instance.

If required flags are not provided and running in a terminal,
interactive mode will be activated to guide you through the setup.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		client := newPostgreSQLClient()

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
		databaseName, _ := cmd.Flags().GetString("database-name")

		missingRequired := name == "" || flavorID == "" || version == "" || userName == "" || password == "" || subnetID == "" || availabilityZone == "" || paramGroupID == "" || storageType == "" || storageSize == 0

		if missingRequired && interactive.CanRunInteractive() {
			azOptions := fetchAvailabilityZoneOptions(ctx)
			interactiveHandler := interactive.NewPostgreSQLInteractive(ctx, client, getRegion(), azOptions)
			interactiveHandler.SetDefinitions()
			pm := interactiveHandler.GetPromptManager()

			pm.SetProvidedValues(map[string]interface{}{
				"name":                name,
				"version":             version,
				"flavor-id":           flavorID,
				"user-name":           userName,
				"password":            password,
				"database-name":       databaseName,
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

			pm.ShowSummary("PostgreSQL Instance Configuration")
			confirmed, err := pm.ConfirmExecution("Create this PostgreSQL instance?")
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
			if v, ok := values["database-name"].(string); ok && v != "" {
				databaseName = v
			}
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
		if databaseName == "" {
			databaseName = "mydb"
		}

		input := &postgresql.CreateInstanceInput{
			DBInstanceName:        name,
			Description:           description,
			DBFlavorID:            flavorID,
			DBVersion:             version,
			DBUserName:            userName,
			DBPassword:            password,
			DBPort:                port,
			DatabaseName:          databaseName,
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
		if backupStartTime != "" {
			if len(backupStartTime) == 5 {
				backupStartTime = backupStartTime + ":00"
			}
			input.Backup.BackupSchedules = []postgresql.BackupSchedule{
				{
					BackupWndBgnTime:  backupStartTime,
					BackupWndDuration: "ONE_HOUR",
				},
			}
		}

		result, err := client.CreateInstance(ctx, input)
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
		executeBackup, _ := cmd.Flags().GetBool("execute-backup")
		client := newPostgreSQLClient()
		req := &postgresql.RestartInstanceRequest{
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

var pgModifyCmd = &cobra.Command{
	Use:   "modify [instance-id]",
	Short: "Modify a PostgreSQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		port, _ := cmd.Flags().GetInt("port")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		paramGroupID, _ := cmd.Flags().GetString("parameter-group-id")
		securityGroupIDs, _ := cmd.Flags().GetStringSlice("security-group-ids")

		input := &postgresql.ModifyInstanceInput{}
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

		client := newPostgreSQLClient()
		result, err := client.ModifyInstance(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to modify instance", err)
		}
		fmt.Printf("Instance modification initiated. Job ID: %s\n", result.JobID)
	},
}

var pgForceRestartCmd = &cobra.Command{
	Use:   "force-restart [instance-id]",
	Short: "Force restart a PostgreSQL instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ForceRestartInstance(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to force restart instance", err)
		}
		fmt.Printf("Force restart initiated. Job ID: %s\n", result.JobID)
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

var pgHARepairCmd = &cobra.Command{
	Use:   "repair [instance-id]",
	Short: "Repair High Availability (recreate standby instance)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.RepairHighAvailability(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to repair HA", err)
		}
		fmt.Printf("HA repair initiated. Job ID: %s\n", result.JobID)
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

var pgHBARuleCmd = &cobra.Command{
	Use:     "hba",
	Aliases: []string{"hba-rule"},
	Short:   "Manage PostgreSQL HBA (Host-Based Authentication) rules",
}

var pgHBARuleListCmd = &cobra.Command{
	Use:   "list [instance-id]",
	Short: "List HBA rules for an instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListHBARules(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list HBA rules", err)
		}
		printPGHBARules(result)
	},
}

var pgHBARuleCreateCmd = &cobra.Command{
	Use:   "create [instance-id]",
	Short: "Create an HBA rule",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		connType, _ := cmd.Flags().GetString("connection-type")
		dbApplyType, _ := cmd.Flags().GetString("database-apply-type")
		userApplyType, _ := cmd.Flags().GetString("user-apply-type")
		address, _ := cmd.Flags().GetString("address")
		authMethod, _ := cmd.Flags().GetString("auth-method")

		if address == "" || authMethod == "" {
			exitWithError("--address and --auth-method are required", nil)
		}

		input := &postgresql.CreateHBARuleRequest{
			ConnectionType:    connType,
			DatabaseApplyType: dbApplyType,
			DBUserApplyType:   userApplyType,
			Address:           address,
			AuthMethod:        authMethod,
		}

		client := newPostgreSQLClient()
		result, err := client.CreateHBARule(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create HBA rule", err)
		}
		fmt.Printf("HBA rule created successfully:\n")
		fmt.Printf("  ID: %s\n", result.HBARuleID)
		fmt.Printf("  Status: %s\n", result.HBARuleStatus)
		fmt.Printf("  Address: %s\n", result.Address)
		fmt.Printf("  Auth Method: %s\n", result.AuthMethod)
	},
}

var pgHBARuleDeleteCmd = &cobra.Command{
	Use:   "delete [instance-id] [rule-id]",
	Short: "Delete an HBA rule",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		_, err := client.DeleteHBARule(context.Background(), args[0], args[1])
		if err != nil {
			exitWithError("failed to delete HBA rule", err)
		}
		fmt.Printf("HBA rule deleted successfully\n")
	},
}

var pgMetricsListCmd = &cobra.Command{
	Use:   "metrics-list [instance-id]",
	Short: "List available metrics",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListMetrics(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list metrics", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "METRIC_ID\tNAME\tUNIT\tDESCRIPTION")
		for _, m := range result.Metrics {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", m.MetricID, m.MetricName, m.Unit, m.Description)
		}
		w.Flush()
	},
}

var pgBackupExportCmd = &cobra.Command{
	Use:   "backup-export [backup-id]",
	Short: "Export backup to object storage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		tenantID, _ := cmd.Flags().GetString("tenant-id")
		containerName, _ := cmd.Flags().GetString("container-name")

		if tenantID == "" || containerName == "" {
			exitWithError("--tenant-id and --container-name are required", nil)
		}

		result, err := client.ExportBackup(context.Background(), args[0], tenantID, containerName)
		if err != nil {
			exitWithError("failed to export backup", err)
		}

		fmt.Printf("Backup export initiated. Job ID: %s\n", result.JobID)
	},
}

var pgBackupRestoreCmd = &cobra.Command{
	Use:   "backup-restore [backup-id]",
	Short: "Restore from backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.RestoreFromBackup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to restore backup", err)
		}

		fmt.Printf("Backup restore initiated. Job ID: %s\n", result.JobID)
	},
}

var pgExtensionCmd = &cobra.Command{
	Use:   "extension",
	Short: "Manage PostgreSQL extensions",
}

var pgExtensionListCmd = &cobra.Command{
	Use:   "list [instance-group-id]",
	Short: "List available extensions for an instance group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListExtensions(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list extensions", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			return
		}

		if len(result.Extensions) == 0 {
			fmt.Println("No extensions found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tDATABASES")
		for _, ext := range result.Extensions {
			dbCount := len(ext.Databases)
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
				ext.ExtensionID, ext.ExtensionName, ext.ExtensionStatus, dbCount)
		}
		w.Flush()

		if result.IsNeedToApply {
			fmt.Println("\nNote: Changes need to be applied to take effect")
		}
	},
}

var pgExtensionGetCmd = &cobra.Command{
	Use:   "get [instance-group-id] [extension-id]",
	Short: "Get extension details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.GetExtension(context.Background(), args[0], args[1])
		if err != nil {
			exitWithError("failed to get extension", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			return
		}

		fmt.Printf("ID:     %s\n", result.ExtensionID)
		fmt.Printf("Name:   %s\n", result.ExtensionName)
		fmt.Printf("Status: %s\n", result.ExtensionStatus)

		if len(result.Databases) > 0 {
			fmt.Println("\nInstalled Databases:")
			for _, db := range result.Databases {
				fmt.Printf("  - %s (%s)\n", db.DatabaseName, db.DBInstanceGroupExtensionStatus)
			}
		}
	},
}

var pgExtensionInstallCmd = &cobra.Command{
	Use:   "install [instance-group-id] [extension-id]",
	Short: "Install extension on a database",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		databaseID, _ := cmd.Flags().GetString("database-id")
		schemaName, _ := cmd.Flags().GetString("schema")
		withCascade, _ := cmd.Flags().GetBool("cascade")

		if databaseID == "" {
			exitWithError("--database-id is required", nil)
		}

		req := &postgresql.InstallExtensionRequest{
			DatabaseID:  databaseID,
			SchemaName:  schemaName,
			WithCascade: withCascade,
		}

		result, err := client.InstallExtension(context.Background(), args[0], args[1], req)
		if err != nil {
			exitWithError("failed to install extension", err)
		}

		fmt.Printf("Extension installation initiated. Job ID: %s\n", result.JobID)
	},
}

var pgExtensionUninstallCmd = &cobra.Command{
	Use:   "uninstall [instance-group-id] [extension-id]",
	Short: "Uninstall extension from a database",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		databaseID, _ := cmd.Flags().GetString("database-id")

		if databaseID == "" {
			exitWithError("--database-id is required", nil)
		}

		result, err := client.UninstallExtension(context.Background(), args[0], args[1], databaseID)
		if err != nil {
			exitWithError("failed to uninstall extension", err)
		}

		fmt.Printf("Extension uninstallation initiated. Job ID: %s\n", result.JobID)
	},
}

var pgStorageResizeCmd = &cobra.Command{
	Use:   "storage-resize [instance-id]",
	Short: "Resize database storage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		size, _ := cmd.Flags().GetInt("size")
		if size < 20 {
			exitWithError("Storage size must be at least 20 GB", nil)
		}

		result, err := client.ResizeStorage(context.Background(), args[0], size)
		if err != nil {
			exitWithError("failed to resize storage", err)
		}

		fmt.Printf("Storage resize initiated to %d GB. Job ID: %s\n", size, result.JobID)
	},
}

var pgEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List database events",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		instanceID, _ := cmd.Flags().GetString("instance-id")
		startTime, _ := cmd.Flags().GetString("start-time")
		endTime, _ := cmd.Flags().GetString("end-time")
		eventCode, _ := cmd.Flags().GetString("event-code")
		sourceType, _ := cmd.Flags().GetString("source-type")
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")

		params := &postgresql.EventParams{
			InstanceID: instanceID,
			StartTime:  startTime,
			EndTime:    endTime,
			EventCode:  eventCode,
			SourceType: sourceType,
			Page:       page,
			Size:       size,
		}

		result, err := client.ListEvents(context.Background(), params)
		if err != nil {
			exitWithError("failed to list events", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			return
		}

		if len(result.Events) == 0 {
			fmt.Println("No events found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "EVENT_TIME\tCATEGORY\tEVENT_NAME\tSOURCE_TYPE\tMESSAGE")
		for _, evt := range result.Events {
			msg := evt.Message
			if len(msg) > 50 {
				msg = msg[:47] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				evt.EventYmdt.Format("2006-01-02 15:04"), evt.Category, evt.EventName, evt.SourceType, msg)
		}
		w.Flush()
	},
}

var pgWatchdogCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "Manage database watchdog",
}

var pgWatchdogGetCmd = &cobra.Command{
	Use:   "get [instance-id]",
	Short: "Get watchdog status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.GetWatchdog(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get watchdog", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			return
		}

		fmt.Printf("Watchdog ID:    %s\n", result.WatchdogID)
		fmt.Printf("Name:           %s\n", result.WatchdogName)
		fmt.Printf("Enabled:        %v\n", result.IsEnabled)
		fmt.Printf("Query Timeout:  %d seconds\n", result.QueryTimeout)
		fmt.Printf("Created:        %s\n", result.CreatedYmdt)
	},
}

var pgWatchdogCreateCmd = &cobra.Command{
	Use:   "create [instance-id]",
	Short: "Create watchdog for instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		name, _ := cmd.Flags().GetString("name")
		enabled, _ := cmd.Flags().GetBool("enabled")
		timeout, _ := cmd.Flags().GetInt("timeout")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		req := &postgresql.CreateWatchdogRequest{
			WatchdogName: name,
			IsEnabled:    enabled,
			QueryTimeout: timeout,
		}

		result, err := client.CreateWatchdog(context.Background(), args[0], req)
		if err != nil {
			exitWithError("failed to create watchdog", err)
		}

		fmt.Printf("Watchdog created. Job ID: %s\n", result.JobID)
	},
}

var pgWatchdogUpdateCmd = &cobra.Command{
	Use:   "update [instance-id]",
	Short: "Update watchdog configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		name, _ := cmd.Flags().GetString("name")
		enabled, _ := cmd.Flags().GetBool("enabled")
		timeout, _ := cmd.Flags().GetInt("timeout")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		req := &postgresql.CreateWatchdogRequest{
			WatchdogName: name,
			IsEnabled:    enabled,
			QueryTimeout: timeout,
		}

		result, err := client.UpdateWatchdog(context.Background(), args[0], req)
		if err != nil {
			exitWithError("failed to update watchdog", err)
		}

		fmt.Printf("Watchdog updated. Job ID: %s\n", result.JobID)
	},
}

var pgWatchdogDeleteCmd = &cobra.Command{
	Use:   "delete [instance-id]",
	Short: "Delete watchdog",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.DeleteWatchdog(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete watchdog", err)
		}

		fmt.Printf("Watchdog deleted. Job ID: %s\n", result.JobID)
	},
}

var pgNotificationMonitoringCmd = &cobra.Command{
	Use:   "notification-monitoring [notification-group-id]",
	Short: "Get notification group monitoring items",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.GetNotificationGroupMonitoringItems(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get monitoring items", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			return
		}

		if len(result.MonitoringItems) == 0 {
			fmt.Println("No monitoring items configured")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTHRESHOLD\tENABLED\tDESCRIPTION")
		for _, item := range result.MonitoringItems {
			enabled := "No"
			if item.IsEnabled {
				enabled = "Yes"
			}
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n",
				item.MonitoringItemID, item.MonitoringItemName, item.Threshold, enabled, item.Description)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(rdsPostgreSQLCmd)

	// Instance commands
	rdsPostgreSQLCmd.AddCommand(pgListCmd)
	rdsPostgreSQLCmd.AddCommand(pgCreateCmd)
	rdsPostgreSQLCmd.AddCommand(pgGetCmd)
	rdsPostgreSQLCmd.AddCommand(pgDeleteCmd)
	rdsPostgreSQLCmd.AddCommand(pgModifyCmd)
	rdsPostgreSQLCmd.AddCommand(pgStartCmd)
	rdsPostgreSQLCmd.AddCommand(pgStopCmd)
	rdsPostgreSQLCmd.AddCommand(pgRestartCmd)
	rdsPostgreSQLCmd.AddCommand(pgForceRestartCmd)

	pgRestartCmd.Flags().Bool("use-failover", false, "Use online failover during restart")
	pgRestartCmd.Flags().Bool("execute-backup", false, "Execute backup before restart")

	pgModifyCmd.Flags().String("name", "", "New instance name")
	pgModifyCmd.Flags().String("description", "", "New description")
	pgModifyCmd.Flags().Int("port", 0, "New PostgreSQL port")
	pgModifyCmd.Flags().String("flavor-id", "", "New flavor ID")
	pgModifyCmd.Flags().String("parameter-group-id", "", "New parameter group ID")
	pgModifyCmd.Flags().StringSlice("security-group-ids", nil, "New security group IDs")

	pgCreateCmd.Flags().String("name", "", "Instance name (required)")
	pgCreateCmd.Flags().String("description", "", "Instance description")
	pgCreateCmd.Flags().String("flavor-id", "", "Flavor ID (required)")
	pgCreateCmd.Flags().String("version", "", "PostgreSQL version (required)")
	pgCreateCmd.Flags().String("user-name", "", "Admin user name (required)")
	pgCreateCmd.Flags().String("password", "", "Admin user password (required)")
	pgCreateCmd.Flags().String("database-name", "", "Initial database name (default: mydb)")
	pgCreateCmd.Flags().Int("port", 5432, "PostgreSQL port")
	pgCreateCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	pgCreateCmd.Flags().String("availability-zone", "", "Availability zone (required, e.g. kr-pub-a)")
	pgCreateCmd.Flags().String("storage-type", "", "Storage type (from API)")
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
	pgHACmd.AddCommand(pgHARepairCmd)
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

	// HBA Rule commands
	rdsPostgreSQLCmd.AddCommand(pgHBARuleCmd)
	pgHBARuleCmd.AddCommand(pgHBARuleListCmd)
	pgHBARuleCmd.AddCommand(pgHBARuleCreateCmd)
	pgHBARuleCmd.AddCommand(pgHBARuleDeleteCmd)
	pgHBARuleCreateCmd.Flags().String("connection-type", "HOST", "Connection type (HOST, HOSTSSL, HOSTNOSSL)")
	pgHBARuleCreateCmd.Flags().String("database-apply-type", "ENTIRE", "Database apply type (ENTIRE, SELECTED)")
	pgHBARuleCreateCmd.Flags().String("user-apply-type", "ENTIRE", "User apply type (ENTIRE, USER_CUSTOM)")
	pgHBARuleCreateCmd.Flags().String("address", "", "Address in CIDR format (e.g., 0.0.0.0/0) (required)")
	pgHBARuleCreateCmd.Flags().String("auth-method", "", "Auth method (SCRAM_SHA_256, TRUST, REJECT) (required)")

	// Metrics commands
	rdsPostgreSQLCmd.AddCommand(pgMetricsListCmd)

	// Backup export/restore commands
	rdsPostgreSQLCmd.AddCommand(pgBackupExportCmd)
	rdsPostgreSQLCmd.AddCommand(pgBackupRestoreCmd)
	pgBackupExportCmd.Flags().String("tenant-id", "", "Tenant ID (required)")
	pgBackupExportCmd.Flags().String("container-name", "", "Object storage container name (required)")

	// Extension commands
	rdsPostgreSQLCmd.AddCommand(pgExtensionCmd)
	pgExtensionCmd.AddCommand(pgExtensionListCmd)
	pgExtensionCmd.AddCommand(pgExtensionGetCmd)
	pgExtensionCmd.AddCommand(pgExtensionInstallCmd)
	pgExtensionCmd.AddCommand(pgExtensionUninstallCmd)
	pgExtensionInstallCmd.Flags().String("database-id", "", "Database ID to install extension on (required)")
	pgExtensionInstallCmd.Flags().String("schema", "public", "Schema name for the extension")
	pgExtensionInstallCmd.Flags().Bool("cascade", false, "Install with cascade (install dependencies)")
	pgExtensionUninstallCmd.Flags().String("database-id", "", "Database ID to uninstall extension from (required)")

	// Storage resize command
	rdsPostgreSQLCmd.AddCommand(pgStorageResizeCmd)
	pgStorageResizeCmd.Flags().Int("size", 0, "New storage size in GB (required, minimum 20)")
	pgStorageResizeCmd.MarkFlagRequired("size")

	// Events command
	rdsPostgreSQLCmd.AddCommand(pgEventsCmd)
	pgEventsCmd.Flags().String("instance-id", "", "Filter by instance ID")
	pgEventsCmd.Flags().String("start-time", "", "Start time (YYYY-MM-DD HH:MM:SS)")
	pgEventsCmd.Flags().String("end-time", "", "End time (YYYY-MM-DD HH:MM:SS)")
	pgEventsCmd.Flags().String("event-code", "", "Filter by event code")
	pgEventsCmd.Flags().String("source-type", "", "Filter by source type")
	pgEventsCmd.Flags().Int("page", 1, "Page number")
	pgEventsCmd.Flags().Int("size", 20, "Page size")

	// Watchdog commands
	rdsPostgreSQLCmd.AddCommand(pgWatchdogCmd)
	pgWatchdogCmd.AddCommand(pgWatchdogGetCmd)
	pgWatchdogCmd.AddCommand(pgWatchdogCreateCmd)
	pgWatchdogCmd.AddCommand(pgWatchdogUpdateCmd)
	pgWatchdogCmd.AddCommand(pgWatchdogDeleteCmd)
	pgWatchdogCreateCmd.Flags().String("name", "", "Watchdog name (required)")
	pgWatchdogCreateCmd.Flags().Bool("enabled", true, "Enable watchdog")
	pgWatchdogCreateCmd.Flags().Int("timeout", 60, "Query timeout in seconds")
	pgWatchdogUpdateCmd.Flags().String("name", "", "Watchdog name (required)")
	pgWatchdogUpdateCmd.Flags().Bool("enabled", true, "Enable watchdog")
	pgWatchdogUpdateCmd.Flags().Int("timeout", 60, "Query timeout in seconds")

	// Notification monitoring command
	rdsPostgreSQLCmd.AddCommand(pgNotificationMonitoringCmd)
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
	fmt.Printf("Storage:  %d GB (%s)\n", result.StorageSize, result.StorageType)
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

func printPGHBARules(result *postgresql.HBARulesResponse) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tORDER\tADDRESS\tAUTH_METHOD\tDB_APPLY\tUSER_APPLY")
	for _, rule := range result.HBARules {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
			rule.HBARuleID,
			rule.HBARuleStatus,
			rule.Order,
			rule.Address,
			rule.AuthMethod,
			rule.DatabaseApplyType,
			rule.DBUserApplyTypeCode)
	}
	w.Flush()
}
