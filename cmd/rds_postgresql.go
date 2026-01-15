package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-cli/pkg/auth"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
	"github.com/spf13/cobra"
)

// ============================================================================
// Base Command
// ============================================================================

var rdsPostgreSQLCmd = &cobra.Command{
	Use:     "rds-postgresql",
	Aliases: []string{"rds-pg", "rds-postgres"},
	Short:   "Manage RDS for PostgreSQL resources",
	Long: `Manage NHN Cloud RDS for PostgreSQL instances, databases, users, and more.

Note: PostgreSQL uses Bearer token authentication, different from MySQL/MariaDB.`,
}

func init() {
	rootCmd.AddCommand(rdsPostgreSQLCmd)
}

// ============================================================================
// Client Helper
// ============================================================================

func newPostgreSQLClient() *postgresql.Client {
	cfg, err := auth.GetPostgreSQLConfig()
	if err != nil {
		exitWithError("failed to get PostgreSQL config", err)
	}

	client, err := postgresql.NewClient(cfg)
	if err != nil {
		exitWithError("failed to create PostgreSQL client", err)
	}

	return client
}

// ============================================================================
// Instance Identifier Resolution
// ============================================================================

func resolvePostgreSQLInstanceIdentifier(client *postgresql.Client, identifier string) (string, error) {
	// If it looks like a UUID, return as-is
	if len(identifier) == 36 && identifier[8] == '-' && identifier[13] == '-' {
		return identifier, nil
	}

	// Otherwise, search by name
	instances, err := client.ListInstances(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to list instances: %w", err)
	}

	for _, inst := range instances.DBInstances {
		if inst.DBInstanceName == identifier {
			return inst.DBInstanceID, nil
		}
	}

	return "", fmt.Errorf("instance not found: %s", identifier)
}

func getResolvedPostgreSQLInstanceID(cmd *cobra.Command, client *postgresql.Client) (string, error) {
	identifier, _ := cmd.Flags().GetString("db-instance-identifier")
	if identifier == "" {
		return "", fmt.Errorf("--db-instance-identifier is required")
	}
	return resolvePostgreSQLInstanceIdentifier(client, identifier)
}

// ============================================================================
// Output Functions
// ============================================================================

func postgresqlPrintJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func postgresqlPrintInstanceList(result *postgresql.ListInstancesResponse) {
	if output == "json" {
		postgresqlPrintJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DB_INSTANCE_ID\tNAME\tSTATUS\tTYPE\tVERSION")
	for _, inst := range result.DBInstances {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			inst.DBInstanceID,
			inst.DBInstanceName,
			inst.DBInstanceStatus,
			inst.DBInstanceType,
			inst.DBVersion,
		)
	}
	w.Flush()
}

func postgresqlPrintInstanceDetail(result *postgresql.GetInstanceResponse) {
	if output == "json" {
		postgresqlPrintJSON(result)
		return
	}

	inst := result.DatabaseInstance
	fmt.Printf("DB Instance ID: %s\n", inst.DBInstanceID)
	fmt.Printf("Name: %s\n", inst.DBInstanceName)
	fmt.Printf("Status: %s\n", inst.DBInstanceStatus)
	fmt.Printf("Type: %s\n", inst.DBInstanceType)
	fmt.Printf("Version: %s\n", inst.DBVersion)
	fmt.Printf("Port: %d\n", inst.DBPort)
	fmt.Printf("Flavor: %s\n", inst.DBFlavorName)
	fmt.Printf("Created: %s\n", inst.CreatedAt)
	fmt.Printf("Updated: %s\n", inst.UpdatedAt)

	if inst.Network.SubnetID != "" {
		fmt.Printf("\nNetwork:\n")
		fmt.Printf("  Subnet: %s\n", inst.Network.SubnetName)
		fmt.Printf("  Availability Zone: %s\n", inst.Network.AvailabilityZone)
		fmt.Printf("  Public Access: %v\n", inst.Network.UsePublicAccess)
		if inst.Network.IPAddress != "" {
			fmt.Printf("  IP Address: %s\n", inst.Network.IPAddress)
		}
	}

	if inst.HighAvailability != nil && inst.HighAvailability.Use {
		fmt.Printf("\nHigh Availability:\n")
		fmt.Printf("  Enabled: %v\n", inst.HighAvailability.Use)
		fmt.Printf("  Availability Zone: %s\n", inst.HighAvailability.AvailabilityZone)
	}
}

// ============================================================================
// Instance Commands
// ============================================================================

var describePostgreSQLInstancesCmd = &cobra.Command{
	Use:   "describe-db-instances",
	Short: "Describe PostgreSQL DB instances",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		identifier, _ := cmd.Flags().GetString("db-instance-identifier")

		if identifier != "" {
			instanceID, err := resolvePostgreSQLInstanceIdentifier(client, identifier)
			if err != nil {
				exitWithError("failed to resolve instance identifier", err)
			}

			result, err := client.GetInstance(context.Background(), instanceID)
			if err != nil {
				exitWithError("failed to get instance", err)
			}
			postgresqlPrintInstanceDetail(result)
		} else {
			result, err := client.ListInstances(context.Background())
			if err != nil {
				exitWithError("failed to list instances", err)
			}
			postgresqlPrintInstanceList(result)
		}
	},
}

var deletePostgreSQLInstanceCmd = &cobra.Command{
	Use:   "delete-db-instance",
	Short: "Delete a PostgreSQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.DeleteInstance(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to delete instance", err)
		}

		fmt.Printf("Instance deletion initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var startPostgreSQLInstanceCmd = &cobra.Command{
	Use:   "start-db-instance",
	Short: "Start a stopped PostgreSQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.StartInstance(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to start instance", err)
		}

		fmt.Printf("Instance start initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var stopPostgreSQLInstanceCmd = &cobra.Command{
	Use:   "stop-db-instance",
	Short: "Stop a running PostgreSQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.StopInstance(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to stop instance", err)
		}

		fmt.Printf("Instance stop initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var rebootPostgreSQLInstanceCmd = &cobra.Command{
	Use:   "reboot-db-instance",
	Short: "Reboot a PostgreSQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		useOnlineFailover, _ := cmd.Flags().GetBool("use-online-failover")
		executeBackup, _ := cmd.Flags().GetBool("execute-backup")

		req := &postgresql.RestartInstanceRequest{}
		if cmd.Flags().Changed("use-online-failover") {
			req.UseOnlineFailover = &useOnlineFailover
		}
		if cmd.Flags().Changed("execute-backup") {
			req.ExecuteBackup = &executeBackup
		}

		result, err := client.RestartInstance(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to reboot instance", err)
		}

		fmt.Printf("Instance reboot initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

func init() {
	// Instance commands
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLInstancesCmd)
	rdsPostgreSQLCmd.AddCommand(deletePostgreSQLInstanceCmd)
	rdsPostgreSQLCmd.AddCommand(startPostgreSQLInstanceCmd)
	rdsPostgreSQLCmd.AddCommand(stopPostgreSQLInstanceCmd)
	rdsPostgreSQLCmd.AddCommand(rebootPostgreSQLInstanceCmd)

	describePostgreSQLInstancesCmd.Flags().String("db-instance-identifier", "", "DB instance identifier")

	deletePostgreSQLInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	startPostgreSQLInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	stopPostgreSQLInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	rebootPostgreSQLInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	rebootPostgreSQLInstanceCmd.Flags().Bool("use-online-failover", false, "Use online failover during reboot (HA only)")
	rebootPostgreSQLInstanceCmd.Flags().Bool("execute-backup", false, "Execute backup before reboot")
}
