package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

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
		client := newMariaDBClient()
		result, err := client.RestartInstance(context.Background(), args[0], useFailover)
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

func init() {
	rootCmd.AddCommand(rdsMariaDBCmd)

	// Instance commands
	rdsMariaDBCmd.AddCommand(mariadbListCmd)
	rdsMariaDBCmd.AddCommand(mariadbGetCmd)
	rdsMariaDBCmd.AddCommand(mariadbDeleteCmd)
	rdsMariaDBCmd.AddCommand(mariadbStartCmd)
	rdsMariaDBCmd.AddCommand(mariadbStopCmd)
	rdsMariaDBCmd.AddCommand(mariadbRestartCmd)
	mariadbRestartCmd.Flags().Bool("use-failover", false, "Use online failover during restart")

	// HA commands
	rdsMariaDBCmd.AddCommand(mariadbHACmd)
	mariadbHACmd.AddCommand(mariadbHAEnableCmd)
	mariadbHACmd.AddCommand(mariadbHADisableCmd)
	mariadbHACmd.AddCommand(mariadbHAPauseCmd)
	mariadbHACmd.AddCommand(mariadbHAResumeCmd)
	mariadbHAEnableCmd.Flags().Int("ping-interval", 3, "Ping interval in seconds")

	// Resource commands
	rdsMariaDBCmd.AddCommand(mariadbFlavorsCmd)
	rdsMariaDBCmd.AddCommand(mariadbVersionsCmd)

	// Backup commands
	rdsMariaDBCmd.AddCommand(mariadbBackupCmd)
	mariadbBackupCmd.AddCommand(mariadbBackupListCmd)
	mariadbBackupCmd.AddCommand(mariadbBackupCreateCmd)
	mariadbBackupCmd.AddCommand(mariadbBackupDeleteCmd)
	mariadbBackupListCmd.Flags().String("instance-id", "", "Filter by instance ID")
	mariadbBackupListCmd.Flags().Int("page", 0, "Page number")
	mariadbBackupListCmd.Flags().Int("size", 20, "Page size")
	mariadbBackupCreateCmd.Flags().String("name", "", "Backup name (required)")
}

func newMariaDBClient() *mariadb.Client {
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
	return mariadb.NewClient(getRegion(), appKey, creds, debug)
}

func printMariaDBInstances(result *mariadb.ListInstancesOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tVERSION\tHA")
	for _, inst := range result.Instances {
		ha := "No"
		if inst.UseHighAvailability {
			ha = "Yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", inst.ID, inst.Name, inst.Status, inst.Version, ha)
	}
	w.Flush()
}

func printMariaDBInstance(result *mariadb.GetInstanceOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	fmt.Printf("ID:       %s\n", result.ID)
	fmt.Printf("Name:     %s\n", result.Name)
	fmt.Printf("Status:   %s\n", result.Status)
	fmt.Printf("Version:  %s\n", result.Version)
	fmt.Printf("Storage:  %d GB (%s)\n", result.StorageSize, result.StorageType)
	fmt.Printf("Port:     %d\n", result.Port)
	fmt.Printf("HA:       %v\n", result.UseHighAvailability)
	fmt.Printf("Created:  %s\n", result.CreatedAt)
}

func printMariaDBFlavors(result *mariadb.ListFlavorsOutput) {
	if output == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tVCPU\tRAM(MB)")
	for _, f := range result.Flavors {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", f.ID, f.Name, f.VCPUs, f.RAM)
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
	for _, v := range result.Versions {
		fmt.Fprintf(w, "%s\t%s\n", v.DBVersion, v.DisplayName)
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
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", b.ID, b.Name, b.Status, b.Size, b.CreatedAt)
	}
	w.Flush()
}
