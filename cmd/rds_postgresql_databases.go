package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
	"github.com/spf13/cobra"
)

// ============================================================================
// Database Commands (PostgreSQL-specific)
// ============================================================================

var describePostgreSQLDatabasesCmd = &cobra.Command{
	Use:   "describe-databases",
	Short: "Describe PostgreSQL databases",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ListDatabases(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to list databases", err)
		}

		if output == "json" {
			postgresqlPrintJSON(result)
		} else {
			postgresqlPrintDatabaseList(result)
		}
	},
}

var createPostgreSQLDatabaseCmd = &cobra.Command{
	Use:   "create-database",
	Short: "Create a PostgreSQL database",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		dbName, _ := cmd.Flags().GetString("database-name")
		owner, _ := cmd.Flags().GetString("owner")
		encoding, _ := cmd.Flags().GetString("encoding")

		if dbName == "" {
			exitWithError("--database-name is required", nil)
		}

		req := &postgresql.CreateDatabaseRequest{
			DatabaseName: dbName,
		}
		if owner != "" {
			req.Owner = owner
		}
		if encoding != "" {
			req.Encoding = encoding
		}

		result, err := client.CreateDatabase(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to create database", err)
		}

		fmt.Printf("Database creation initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

var deletePostgreSQLDatabaseCmd = &cobra.Command{
	Use:   "delete-database",
	Short: "Delete a PostgreSQL database",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		dbID, _ := cmd.Flags().GetString("database-id")
		if dbID == "" {
			exitWithError("--database-id is required", nil)
		}

		_, err = client.DeleteDatabase(context.Background(), instanceID, dbID)
		if err != nil {
			exitWithError("failed to delete database", err)
		}

		fmt.Printf("Database deleted successfully.\n")
	},
}

// ============================================================================
// User Commands
// ============================================================================

var describePostgreSQLUsersCmd = &cobra.Command{
	Use:   "describe-db-users",
	Short: "Describe PostgreSQL DB users",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ListDBUsers(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to list DB users", err)
		}

		if output == "json" {
			postgresqlPrintJSON(result)
		} else {
			postgresqlPrintUserList(result)
		}
	},
}

var createPostgreSQLUserCmd = &cobra.Command{
	Use:   "create-db-user",
	Short: "Create a PostgreSQL database user",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		username, _ := cmd.Flags().GetString("db-user-name")
		password, _ := cmd.Flags().GetString("db-password")

		if username == "" {
			exitWithError("--db-user-name is required", nil)
		}
		if password == "" {
			exitWithError("--db-password is required", nil)
		}

		req := &postgresql.CreateDBUserRequest{
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

var deletePostgreSQLUserCmd = &cobra.Command{
	Use:   "delete-db-user",
	Short: "Delete a PostgreSQL database user",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
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

		fmt.Printf("DB user deleted successfully.\n")
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func postgresqlPrintDatabaseList(result *postgresql.ListDatabasesResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATABASE_ID\tNAME\tOWNER\tENCODING\tSIZE")
	for _, db := range result.Databases {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
			db.DatabaseID,
			db.DatabaseName,
			db.Owner,
			db.Encoding,
			db.Size,
		)
	}
	w.Flush()
}

func postgresqlPrintUserList(result *postgresql.ListDBUsersResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "USER_ID\tUSERNAME\tSUPERUSER\tCREATEDB")
	for _, user := range result.DBUsers {
		fmt.Fprintf(w, "%s\t%s\t%v\t%v\n",
			user.DBUserID,
			user.DBUserName,
			user.IsSuperuser,
			user.CanCreateDB,
		)
	}
	w.Flush()
}

func init() {
	// Database commands
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLDatabasesCmd)
	rdsPostgreSQLCmd.AddCommand(createPostgreSQLDatabaseCmd)
	rdsPostgreSQLCmd.AddCommand(deletePostgreSQLDatabaseCmd)

	describePostgreSQLDatabasesCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	createPostgreSQLDatabaseCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createPostgreSQLDatabaseCmd.Flags().String("database-name", "", "Database name (required)")
	createPostgreSQLDatabaseCmd.Flags().String("owner", "", "Database owner username")
	createPostgreSQLDatabaseCmd.Flags().String("encoding", "", "Character encoding (default: UTF8)")

	deletePostgreSQLDatabaseCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	deletePostgreSQLDatabaseCmd.Flags().String("database-id", "", "Database ID (required)")

	// User commands
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLUsersCmd)
	rdsPostgreSQLCmd.AddCommand(createPostgreSQLUserCmd)
	rdsPostgreSQLCmd.AddCommand(deletePostgreSQLUserCmd)

	describePostgreSQLUsersCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	createPostgreSQLUserCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createPostgreSQLUserCmd.Flags().String("db-user-name", "", "Database username (required)")
	createPostgreSQLUserCmd.Flags().String("db-password", "", "Database password (required)")

	deletePostgreSQLUserCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	deletePostgreSQLUserCmd.Flags().String("db-user-id", "", "DB user ID (required)")
}
