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
// Security Group Commands
// ============================================================================

var describePostgreSQLSecurityGroupsCmd = &cobra.Command{
	Use:   "describe-db-security-groups",
	Short: "Describe PostgreSQL security groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		result, err := client.ListSecurityGroups(context.Background())
		if err != nil {
			exitWithError("failed to list security groups", err)
		}

		if output == "json" {
			postgresqlPrintJSON(result)
		} else {
			postgresqlPrintSecurityGroupList(result)
		}
	},
}

var createPostgreSQLSecurityGroupCmd = &cobra.Command{
	Use:   "create-db-security-group",
	Short: "Create a PostgreSQL security group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		name, _ := cmd.Flags().GetString("db-security-group-name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			exitWithError("--db-security-group-name is required", nil)
		}

		req := &postgresql.CreateSecurityGroupRequest{
			DBSecurityGroupName: name,
			Description:         description,
		}

		result, err := client.CreateSecurityGroup(context.Background(), req)
		if err != nil {
			exitWithError("failed to create security group", err)
		}

		fmt.Printf("Security group created.\n")
		fmt.Printf("ID: %s\n", result.DBSecurityGroupID)
	},
}

var deletePostgreSQLSecurityGroupCmd = &cobra.Command{
	Use:   "delete-db-security-group",
	Short: "Delete a PostgreSQL security group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		sgID, _ := cmd.Flags().GetString("db-security-group-identifier")
		if sgID == "" {
			exitWithError("--db-security-group-identifier is required", nil)
		}

		_, err := client.DeleteSecurityGroup(context.Background(), sgID)
		if err != nil {
			exitWithError("failed to delete security group", err)
		}

		fmt.Printf("Security group deleted successfully.\n")
	},
}

// ============================================================================
// Security Group Rule Commands
// ============================================================================

var authorizePostgreSQLSecurityGroupIngressCmd = &cobra.Command{
	Use:   "authorize-db-security-group-ingress",
	Short: "Authorize ingress to a DB security group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		groupID, _ := cmd.Flags().GetString("db-security-group-identifier")
		cidr, _ := cmd.Flags().GetString("cidr")
		port, _ := cmd.Flags().GetInt("port")
		desc, _ := cmd.Flags().GetString("description")

		if groupID == "" {
			exitWithError("--db-security-group-identifier is required", nil)
		}
		if cidr == "" {
			exitWithError("--cidr is required (e.g. 0.0.0.0/0)", nil)
		}
		if port == 0 {
			port = 5432
		}

		req := &postgresql.CreateSecurityRuleRequest{
			Direction:   "INGRESS",
			EtherType:   "IPv4",
			CIDR:        cidr,
			Description: desc,
			Port: postgresql.RulePort{
				PortType: "range",
				MinPort:  &port,
				MaxPort:  &port,
			},
		}

		result, err := client.CreateSecurityRule(context.Background(), groupID, req)
		if err != nil {
			exitWithError("failed to authorize ingress", err)
		}

		fmt.Printf("Ingress rule authorized. Rule ID: %s\n", result.RuleID)
	},
}

var revokePostgreSQLSecurityGroupIngressCmd = &cobra.Command{
	Use:   "revoke-db-security-group-ingress",
	Short: "Revoke ingress from a DB security group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		groupID, _ := cmd.Flags().GetString("db-security-group-identifier")
		ruleID, _ := cmd.Flags().GetString("security-group-rule-id")

		if groupID == "" {
			exitWithError("--db-security-group-identifier is required", nil)
		}
		if ruleID == "" {
			exitWithError("--security-group-rule-id is required", nil)
		}

		_, err := client.DeleteSecurityRule(context.Background(), groupID, ruleID)
		if err != nil {
			exitWithError("failed to revoke ingress", err)
		}

		fmt.Printf("Ingress rule revoked successfully.\n")
	},
}

// ============================================================================
// Print Helpers
// ============================================================================

func postgresqlPrintSecurityGroupList(result *postgresql.ListSecurityGroupsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SG_ID\tNAME\tDESCRIPTION\tRULES")
	for _, sg := range result.DBSecurityGroups {
		ruleCount := len(sg.Rules)
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
			sg.DBSecurityGroupID,
			sg.DBSecurityGroupName,
			sg.Description,
			ruleCount,
		)
	}
	w.Flush()
}

func init() {
	// Register commands
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLSecurityGroupsCmd)
	rdsPostgreSQLCmd.AddCommand(createPostgreSQLSecurityGroupCmd)
	rdsPostgreSQLCmd.AddCommand(deletePostgreSQLSecurityGroupCmd)

	rdsPostgreSQLCmd.AddCommand(authorizePostgreSQLSecurityGroupIngressCmd)
	rdsPostgreSQLCmd.AddCommand(revokePostgreSQLSecurityGroupIngressCmd)

	// Flags
	createPostgreSQLSecurityGroupCmd.Flags().String("db-security-group-name", "", "Security group name (required)")
	createPostgreSQLSecurityGroupCmd.Flags().String("description", "", "Description")

	deletePostgreSQLSecurityGroupCmd.Flags().String("db-security-group-identifier", "", "Security group ID (required)")

	authorizePostgreSQLSecurityGroupIngressCmd.Flags().String("db-security-group-identifier", "", "Security group ID (required)")
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().String("cidr", "", "CIDR range (required)")
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().Int("port", 5432, "Port (default 5432)")
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().String("description", "", "Rule description")

	revokePostgreSQLSecurityGroupIngressCmd.Flags().String("db-security-group-identifier", "", "Security group ID (required)")
	revokePostgreSQLSecurityGroupIngressCmd.Flags().String("security-group-rule-id", "", "Security group rule ID (required)")
}
