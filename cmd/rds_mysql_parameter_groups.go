package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// Parameter Group Commands
// ============================================================================

var describeDBParameterGroupsCmd = &cobra.Command{
	Use:   "describe-db-parameter-groups",
	Short: "Describe MySQL DB parameter groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListParameterGroups(context.Background())
		if err != nil {
			exitWithError("failed to list parameter groups", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			for _, pg := range result.ParameterGroups {
				fmt.Printf("%s: %s\n", pg.ParameterGroupID, pg.ParameterGroupName)
			}
		}
	},
}

var createDBParameterGroupCmd = &cobra.Command{
	Use:   "create-db-parameter-group",
	Short: "Create a DB parameter group",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("db-parameter-group-name")
		description, _ := cmd.Flags().GetString("description")
		dbVersion, _ := cmd.Flags().GetString("engine-version")

		if name == "" {
			exitWithError("--db-parameter-group-name is required", nil)
		}
		if dbVersion == "" {
			exitWithError("--engine-version is required", nil)
		}

		client := newMySQLClient()
		req := &mysql.CreateParameterGroupRequest{
			ParameterGroupName: name,
			Description:        description,
			DBVersion:          dbVersion,
		}

		result, err := client.CreateParameterGroup(context.Background(), req)
		if err != nil {
			exitWithError("failed to create parameter group", err)
		}

		fmt.Printf("Parameter group created: %s\n", result.ParameterGroupID)
	},
}

var modifyDBParameterGroupCmd = &cobra.Command{
	Use:   "modify-db-parameter-group",
	Short: "Modify a DB parameter group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("db-parameter-group-id")
		name, _ := cmd.Flags().GetString("db-parameter-group-name")
		description, _ := cmd.Flags().GetString("description")

		if groupID == "" {
			exitWithError("--db-parameter-group-id is required", nil)
		}

		client := newMySQLClient()
		req := &mysql.UpdateParameterGroupRequest{}

		if name != "" {
			req.ParameterGroupName = &name
		}
		if description != "" {
			req.Description = &description
		}

		_, err := client.UpdateParameterGroup(context.Background(), groupID, req)
		if err != nil {
			exitWithError("failed to modify parameter group", err)
		}

		fmt.Printf("Parameter group modified successfully\n")
	},
}

var deleteDBParameterGroupCmd = &cobra.Command{
	Use:   "delete-db-parameter-group",
	Short: "Delete a DB parameter group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("db-parameter-group-id")
		if groupID == "" {
			exitWithError("--db-parameter-group-id is required", nil)
		}

		client := newMySQLClient()
		_, err := client.DeleteParameterGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to delete parameter group", err)
		}

		fmt.Printf("Parameter group deleted successfully\n")
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	rdsMySQLCmd.AddCommand(describeDBParameterGroupsCmd)
	rdsMySQLCmd.AddCommand(createDBParameterGroupCmd)
	rdsMySQLCmd.AddCommand(modifyDBParameterGroupCmd)
	rdsMySQLCmd.AddCommand(deleteDBParameterGroupCmd)

	// create flags
	createDBParameterGroupCmd.Flags().String("db-parameter-group-name", "", "Parameter group name (required)")
	createDBParameterGroupCmd.Flags().String("description", "", "Description")
	createDBParameterGroupCmd.Flags().String("engine-version", "", "Engine version (required)")

	// modify flags
	modifyDBParameterGroupCmd.Flags().String("db-parameter-group-id", "", "Parameter group ID (required)")
	modifyDBParameterGroupCmd.Flags().String("db-parameter-group-name", "", "New parameter group name")
	modifyDBParameterGroupCmd.Flags().String("description", "", "New description")

	// delete flags
	deleteDBParameterGroupCmd.Flags().String("db-parameter-group-id", "", "Parameter group ID (required)")
}
