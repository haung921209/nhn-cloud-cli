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
// Parameter Group Commands
// ============================================================================

var describePostgreSQLParameterGroupsCmd = &cobra.Command{
	Use:   "describe-db-parameter-groups",
	Short: "Describe PostgreSQL parameter groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		groupID, _ := cmd.Flags().GetString("db-parameter-group-id")

		if groupID != "" {
			result, err := client.GetParameterGroup(context.Background(), groupID)
			if err != nil {
				exitWithError("failed to get parameter group", err)
			}
			if output == "json" {
				postgresqlPrintJSON(result)
			} else {
				postgresqlPrintParameterGroupDetail(result)
			}
		} else {
			result, err := client.ListParameterGroups(context.Background())
			if err != nil {
				exitWithError("failed to list parameter groups", err)
			}
			if output == "json" {
				postgresqlPrintJSON(result)
			} else {
				postgresqlPrintParameterGroupList(result)
			}
		}
	},
}

var createPostgreSQLParameterGroupCmd = &cobra.Command{
	Use:   "create-db-parameter-group",
	Short: "Create a PostgreSQL parameter group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		name, _ := cmd.Flags().GetString("db-parameter-group-name")
		description, _ := cmd.Flags().GetString("description")
		dbVersion, _ := cmd.Flags().GetString("db-parameter-group-family")

		if name == "" {
			exitWithError("--db-parameter-group-name is required", nil)
		}
		if dbVersion == "" {
			exitWithError("--db-parameter-group-family is required (e.g., POSTGRESQL_V14_6)", nil)
		}

		req := &postgresql.CreateParameterGroupRequest{
			ParameterGroupName: name,
			Description:        description,
			DBVersion:          dbVersion,
		}

		result, err := client.CreateParameterGroup(context.Background(), req)
		if err != nil {
			exitWithError("failed to create parameter group", err)
		}

		fmt.Printf("Parameter group created.\n")
		fmt.Printf("ID: %s\n", result.ParameterGroupID)
	},
}

var deletePostgreSQLParameterGroupCmd = &cobra.Command{
	Use:   "delete-db-parameter-group",
	Short: "Delete a PostgreSQL parameter group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		groupID, _ := cmd.Flags().GetString("db-parameter-group-id")
		if groupID == "" {
			exitWithError("--db-parameter-group-id is required", nil)
		}

		_, err := client.DeleteParameterGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to delete parameter group", err)
		}

		fmt.Printf("Parameter group deleted successfully.\n")
	},
}

var resetPostgreSQLParameterGroupCmd = &cobra.Command{
	Use:   "reset-db-parameter-group",
	Short: "Reset a PostgreSQL parameter group to defaults",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		groupID, _ := cmd.Flags().GetString("db-parameter-group-id")
		if groupID == "" {
			exitWithError("--db-parameter-group-id is required", nil)
		}

		_, err := client.ResetParameterGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to reset parameter group", err)
		}

		fmt.Printf("Parameter group reset successfully.\n")
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func postgresqlPrintParameterGroupList(result *postgresql.ListParameterGroupsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PG_ID\tNAME\tVERSION\tDESCRIPTION")
	for _, pg := range result.ParameterGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			pg.ParameterGroupID,
			pg.ParameterGroupName,
			pg.DBVersion,
			pg.Description,
		)
	}
	w.Flush()
}

func postgresqlPrintParameterGroupDetail(result *postgresql.GetParameterGroupResponse) {
	pg := result.ParameterGroup
	fmt.Printf("ID: %s\n", pg.ParameterGroupID)
	fmt.Printf("Name: %s\n", pg.ParameterGroupName)
	fmt.Printf("Version: %s\n", pg.DBVersion)
	fmt.Printf("Description: %s\n", pg.Description)

	if len(pg.Parameters) > 0 {
		fmt.Println("\nParameters:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVALUE\tMODIFIABLE")
		for _, p := range pg.Parameters {
			fmt.Fprintf(w, "%s\t%s\t%v\n", p.ParameterName, p.Value, p.IsModifiable)
		}
		w.Flush()
	}
}

func init() {
	// Parameter Group commands
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLParameterGroupsCmd)
	rdsPostgreSQLCmd.AddCommand(createPostgreSQLParameterGroupCmd)
	rdsPostgreSQLCmd.AddCommand(deletePostgreSQLParameterGroupCmd)
	rdsPostgreSQLCmd.AddCommand(resetPostgreSQLParameterGroupCmd)

	describePostgreSQLParameterGroupsCmd.Flags().String("db-parameter-group-id", "", "Parameter group ID")

	createPostgreSQLParameterGroupCmd.Flags().String("db-parameter-group-name", "", "Parameter group name (required)")
	createPostgreSQLParameterGroupCmd.Flags().String("db-parameter-group-family", "", "DB version (required, e.g., POSTGRESQL_V14_6)")
	createPostgreSQLParameterGroupCmd.Flags().String("description", "", "Description")

	deletePostgreSQLParameterGroupCmd.Flags().String("db-parameter-group-id", "", "Parameter group ID (required)")
	resetPostgreSQLParameterGroupCmd.Flags().String("db-parameter-group-id", "", "Parameter group ID (required)")
}
