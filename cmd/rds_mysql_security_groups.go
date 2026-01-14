package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// Security Group Commands
// ============================================================================

var describeDBSecurityGroupsCmd = &cobra.Command{
	Use:   "describe-db-security-groups",
	Short: "Describe MySQL DB security groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListSecurityGroups(context.Background())
		if err != nil {
			exitWithError("failed to list security groups", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			for _, sg := range result.DBSecurityGroups {
				fmt.Printf("%s: %s (%d rules)\n", sg.DBSecurityGroupID, sg.DBSecurityGroupName, len(sg.Rules))
			}
		}
	},
}

var createDBSecurityGroupCmd = &cobra.Command{
	Use:   "create-db-security-group",
	Short: "Create a DB security group",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("db-security-group-name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			exitWithError("--db-security-group-name is required", nil)
		}

		client := newMySQLClient()
		req := &mysql.CreateSecurityGroupRequest{
			DBSecurityGroupName: name,
			Description:         description,
		}

		result, err := client.CreateSecurityGroup(context.Background(), req)
		if err != nil {
			exitWithError("failed to create security group", err)
		}

		fmt.Printf("Security group created: %s\n", result.DBSecurityGroupID)
	},
}

var authorizeDBSecurityGroupIngressCmd = &cobra.Command{
	Use:   "authorize-db-security-group-ingress",
	Short: "Authorize ingress rule for DB security group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("db-security-group-id")
		cidr, _ := cmd.Flags().GetString("cidr")

		if groupID != "" && cidr != "" {
			fmt.Println("NOTE: NHN Cloud requires security group rules to be specified during group creation.")
		} else {
			fmt.Println("NOTE: NHN Cloud requires security group rules to be specified during group creation.")
		}
		fmt.Println("Rules cannot be added separately. Please recreate the security group with the desired rules.")
	},
}

var deleteDBSecurityGroupCmd = &cobra.Command{
	Use:   "delete-db-security-group",
	Short: "Delete a DB security group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("db-security-group-id")
		if groupID == "" {
			exitWithError("--db-security-group-id is required", nil)
		}

		client := newMySQLClient()
		_, err := client.DeleteSecurityGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to delete security group", err)
		}

		fmt.Printf("Security group deleted successfully\n")
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	rdsMySQLCmd.AddCommand(describeDBSecurityGroupsCmd)
	rdsMySQLCmd.AddCommand(createDBSecurityGroupCmd)
	rdsMySQLCmd.AddCommand(authorizeDBSecurityGroupIngressCmd)
	rdsMySQLCmd.AddCommand(deleteDBSecurityGroupCmd)

	// create flags
	createDBSecurityGroupCmd.Flags().String("db-security-group-name", "", "Security group name (required)")
	createDBSecurityGroupCmd.Flags().String("description", "", "Description")

	// authorize ingress flags
	authorizeDBSecurityGroupIngressCmd.Flags().String("db-security-group-id", "", "Security group ID (required)")
	authorizeDBSecurityGroupIngressCmd.Flags().String("cidr", "", "CIDR block (required, e.g., 0.0.0.0/0)")
	authorizeDBSecurityGroupIngressCmd.Flags().String("description", "", "Rule description")

	// delete flags
	deleteDBSecurityGroupCmd.Flags().String("db-security-group-id", "", "Security group ID (required)")
}
