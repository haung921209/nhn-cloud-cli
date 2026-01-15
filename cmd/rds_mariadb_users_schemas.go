package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mariadb"
	"github.com/spf13/cobra"
)

// ============================================================================
// DB User Commands
// ============================================================================

var describeMariaDBUsersCmd = &cobra.Command{
	Use:   "describe-db-users",
	Short: "Describe MariaDB DB users",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		instanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ListDBUsers(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to list DB users", err)
		}

		if output == "json" {
			mariadbPrintJSON(result)
		} else {
			mariadbPrintUserList(result)
		}
	},
}

var createMariaDBUserCmd = &cobra.Command{
	Use:   "create-db-user",
	Short: "Create a MariaDB database user",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		instanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		username, _ := cmd.Flags().GetString("db-user-name")
		password, _ := cmd.Flags().GetString("db-password")
		host, _ := cmd.Flags().GetString("host")
		authorityType, _ := cmd.Flags().GetString("authority-type")

		if username == "" {
			exitWithError("--db-user-name is required", nil)
		}
		if password == "" {
			exitWithError("--db-password is required (4-16 characters)", nil)
		}
		if host == "" {
			exitWithError("--host is required (e.g., '%' for all hosts)", nil)
		}
		if authorityType == "" {
			exitWithError("--authority-type is required (READ, WRITE, DDL, etc.)", nil)
		}

		req := &mariadb.CreateDBUserRequest{
			DBUserName:    username,
			DBPassword:    password,
			Host:          host,
			AuthorityType: authorityType,
			// AuthenticationPlugin: optional, MariaDB uses NATIVE or SHA256
		}

		result, err := client.CreateDBUser(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to create DB user", err)
		}

		fmt.Printf("DB user creation initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var deleteMariaDBUserCmd = &cobra.Command{
	Use:   "delete-db-user",
	Short: "Delete a MariaDB database user",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		instanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		userID, _ := cmd.Flags().GetString("db-user-id")
		if userID == "" {
			exitWithError("--db-user-id is required", nil)
		}

		_, err = client.DeleteDBUser(context.Background(), instanceID, userID)
		if err != nil {
			exitWithError("failed to delete DB user", err)
		}

		fmt.Printf("DB user deleted successfully\n")
	},
}

// ============================================================================
// DB Schema Commands
// ============================================================================

var describeMariaDBSchemasCmd = &cobra.Command{
	Use:   "describe-db-schemas",
	Short: "Describe MariaDB DB schemas",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		instanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ListSchemas(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to list schemas", err)
		}

		if output == "json" {
			mariadbPrintJSON(result)
		} else {
			mariadbPrintSchemaList(result)
		}
	},
}

var createMariaDBSchemaCmd = &cobra.Command{
	Use:   "create-db-schema",
	Short: "Create a MariaDB database schema",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		instanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		schemaName, _ := cmd.Flags().GetString("db-schema-name")
		if schemaName == "" {
			exitWithError("--db-schema-name is required", nil)
		}

		req := &mariadb.CreateSchemaRequest{
			DBSchemaName: schemaName,
		}

		result, err := client.CreateSchema(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to create schema", err)
		}

		fmt.Printf("Schema creation initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var deleteMariaDBSchemaCmd = &cobra.Command{
	Use:   "delete-db-schema",
	Short: "Delete a MariaDB database schema",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		instanceID, err := getResolvedMariaDBInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		schemaID, _ := cmd.Flags().GetString("db-schema-id")
		if schemaID == "" {
			exitWithError("--db-schema-id is required", nil)
		}

		_, err = client.DeleteSchema(context.Background(), instanceID, schemaID)
		if err != nil {
			exitWithError("failed to delete schema", err)
		}

		fmt.Printf("Schema deleted successfully\n")
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func mariadbPrintUserList(result *mariadb.ListDBUsersResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tUSERNAME\tHOST\tAUTHORITY")
	for _, user := range result.DBUsers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			user.DBUserID,
			user.DBUserName,
			user.Host,
			user.AuthorityType,
		)
	}
	w.Flush()
}

func mariadbPrintSchemaList(result *mariadb.ListSchemasResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME")
	for _, schema := range result.DBSchemas {
		fmt.Fprintf(w, "%s\t%s\n",
			schema.DBSchemaID,
			schema.DBSchemaName,
		)
	}
	w.Flush()
}

func init() {
	// DB Users
	rdsMariaDBCmd.AddCommand(describeMariaDBUsersCmd)
	rdsMariaDBCmd.AddCommand(createMariaDBUserCmd)
	rdsMariaDBCmd.AddCommand(deleteMariaDBUserCmd)

	describeMariaDBUsersCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	createMariaDBUserCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createMariaDBUserCmd.Flags().String("db-user-name", "", "Database username (required)")
	createMariaDBUserCmd.Flags().String("db-password", "", "Database password, 4-16 chars (required)")
	createMariaDBUserCmd.Flags().String("host", "", "Host pattern (required, e.g., '%' for all hosts)")
	createMariaDBUserCmd.Flags().String("authority-type", "", "Authority type (required: READ, WRITE, DDL, etc.)")

	deleteMariaDBUserCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	deleteMariaDBUserCmd.Flags().String("db-user-id", "", "DB user ID (required)")

	// DB Schemas
	rdsMariaDBCmd.AddCommand(describeMariaDBSchemasCmd)
	rdsMariaDBCmd.AddCommand(createMariaDBSchemaCmd)
	rdsMariaDBCmd.AddCommand(deleteMariaDBSchemaCmd)

	describeMariaDBSchemasCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	createMariaDBSchemaCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createMariaDBSchemaCmd.Flags().String("db-schema-name", "", "Schema name (required)")

	deleteMariaDBSchemaCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	deleteMariaDBSchemaCmd.Flags().String("db-schema-id", "", "Schema ID (required)")
}
