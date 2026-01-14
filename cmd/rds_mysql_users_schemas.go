package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// DB User Commands
// ============================================================================

var createDBUserCmd = &cobra.Command{
	Use:   "create-db-user",
	Short: "Create a database user",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		username, _ := cmd.Flags().GetString("db-user-name")
		password, _ := cmd.Flags().GetString("db-password")

		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}
		if username == "" {
			exitWithError("--db-user-name is required", nil)
		}
		if password == "" {
			exitWithError("--db-password is required (4-16 characters)", nil)
		}

		client := newMySQLClient()
		req := &mysql.CreateDBUserRequest{
			DBUserName: username,
			DBPassword: password,
		}

		result, err := client.CreateDBUser(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to create DB user", err)
		}

		fmt.Printf("DB user creation initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var deleteDBUserCmd = &cobra.Command{
	Use:   "delete-db-user",
	Short: "Delete a database user",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		userID, _ := cmd.Flags().GetString("db-user-id")

		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}
		if userID == "" {
			exitWithError("--db-user-id is required", nil)
		}

		client := newMySQLClient()
		_, err := client.DeleteDBUser(context.Background(), instanceID, userID)
		if err != nil {
			exitWithError("failed to delete DB user", err)
		}

		fmt.Printf("DB user deleted successfully\n")
	},
}

// ============================================================================
// DB Schema Commands
// ============================================================================

var createDBSchemaCmd = &cobra.Command{
	Use:   "create-db-schema",
	Short: "Create a database schema",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		schemaName, _ := cmd.Flags().GetString("db-schema-name")

		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}
		if schemaName == "" {
			exitWithError("--db-schema-name is required", nil)
		}

		client := newMySQLClient()
		req := &mysql.CreateSchemaRequest{
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

var deleteDBSchemaCmd = &cobra.Command{
	Use:   "delete-db-schema",
	Short: "Delete a database schema",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")
		schemaID, _ := cmd.Flags().GetString("db-schema-id")

		if instanceID == "" {
			exitWithError("--db-instance-identifier is required", nil)
		}
		if schemaID == "" {
			exitWithError("--db-schema-id is required", nil)
		}

		client := newMySQLClient()
		_, err := client.DeleteSchema(context.Background(), instanceID, schemaID)
		if err != nil {
			exitWithError("failed to delete schema", err)
		}

		fmt.Printf("Schema deleted successfully\n")
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	rdsMySQLCmd.AddCommand(createDBUserCmd)
	rdsMySQLCmd.AddCommand(deleteDBUserCmd)
	rdsMySQLCmd.AddCommand(createDBSchemaCmd)
	rdsMySQLCmd.AddCommand(deleteDBSchemaCmd)

	// create-db-user flags
	createDBUserCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createDBUserCmd.Flags().String("db-user-name", "", "Database username (required)")
	createDBUserCmd.Flags().String("db-password", "", "Database password, 4-16 chars (required)")

	// delete-db-user flags
	deleteDBUserCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	deleteDBUserCmd.Flags().String("db-user-id", "", "DB user ID (required)")

	// create-db-schema flags
	createDBSchemaCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createDBSchemaCmd.Flags().String("db-schema-name", "", "Schema name (required)")

	// delete-db-schema flags
	deleteDBSchemaCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	deleteDBSchemaCmd.Flags().String("db-schema-id", "", "Schema ID (required)")
}
