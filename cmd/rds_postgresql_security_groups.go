package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
	"github.com/spf13/cobra"
)

func init() {
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLSecurityGroupsCmd)

	rdsPostgreSQLCmd.AddCommand(getPostgreSQLSecurityGroupCmd)
	getPostgreSQLSecurityGroupCmd.Flags().String("db-security-group-identifier", "", "DB security group identifier (Required)")
	getPostgreSQLSecurityGroupCmd.MarkFlagRequired("db-security-group-identifier")

	rdsPostgreSQLCmd.AddCommand(createPostgreSQLSecurityGroupCmd)
	createPostgreSQLSecurityGroupCmd.Flags().String("db-security-group-name", "", "Security group name (required)")
	createPostgreSQLSecurityGroupCmd.Flags().String("description", "", "Description")
	createPostgreSQLSecurityGroupCmd.Flags().String("cidr", "", "Initial CIDR block (required for PostgreSQL)")
	createPostgreSQLSecurityGroupCmd.Flags().Int("port", 0, "Specific port (e.g. 5432, 15432)")
	createPostgreSQLSecurityGroupCmd.Flags().Int("min-port", 0, "Minimum port for range")
	createPostgreSQLSecurityGroupCmd.Flags().Int("max-port", 0, "Maximum port for range")

	rdsPostgreSQLCmd.AddCommand(authorizePostgreSQLSecurityGroupIngressCmd)
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().String("db-security-group-identifier", "", "DB security group identifier (Required)")
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().String("cidr", "", "CIDR block to allow (Required)")
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().String("description", "", "Rule description")
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().Int("port", 0, "Specific port (e.g. 5432, 15432)")
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().Int("min-port", 0, "Minimum port for range")
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().Int("max-port", 0, "Maximum port for range")
	authorizePostgreSQLSecurityGroupIngressCmd.MarkFlagRequired("db-security-group-identifier")
	authorizePostgreSQLSecurityGroupIngressCmd.MarkFlagRequired("cidr")

	rdsPostgreSQLCmd.AddCommand(deletePostgreSQLSecurityGroupCmd)
	deletePostgreSQLSecurityGroupCmd.Flags().String("db-security-group-identifier", "", "Security group identifier (required)")
	deletePostgreSQLSecurityGroupCmd.MarkFlagRequired("db-security-group-identifier")
}

var describePostgreSQLSecurityGroupsCmd = &cobra.Command{
	Use:   "describe-db-security-groups",
	Short: "Describe PostgreSQL DB security groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		result, err := client.ListSecurityGroups(context.Background())
		if err != nil {
			exitWithError("failed to list security groups", err)
		}

		if output == "json" {
			postgresqlPrintJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME")
		for _, sg := range result.DBSecurityGroups {
			fmt.Fprintf(w, "%s\t%s\n",
				sg.DBSecurityGroupID,
				sg.DBSecurityGroupName,
			)
		}
		w.Flush()
	},
}

var getPostgreSQLSecurityGroupCmd = &cobra.Command{
	Use:   "get-db-security-group",
	Short: "Get details of a PostgreSQL DB security group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		identifier, _ := cmd.Flags().GetString("db-security-group-identifier")

		// Resolve identifier if it's a name (Naive implementation: Assume ID for now)
		groupID := identifier

		result, err := client.GetSecurityGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to get security group", err)
		}

		if output == "json" {
			postgresqlPrintJSON(result)
			return
		}

		fmt.Printf("ID: %s\n", result.DBSecurityGroup.DBSecurityGroupID)
		fmt.Printf("Name: %s\n", result.DBSecurityGroup.DBSecurityGroupName)
		fmt.Println("Rules:")
		for _, rule := range result.DBSecurityGroup.Rules {
			// Handle potential nil pointers safely
			minPort := 0
			maxPort := 0
			if rule.Port.MinPort != nil {
				minPort = *rule.Port.MinPort
			}
			if rule.Port.MaxPort != nil {
				maxPort = *rule.Port.MaxPort
			}

			fmt.Printf("  - RuleID: %s, Protocol: %s, Port: %d-%d, CIDR: %s\n",
				rule.RuleID,
				rule.EtherType,
				minPort,
				maxPort,
				rule.CIDR)
		}
	},
}

var createPostgreSQLSecurityGroupCmd = &cobra.Command{
	Use:   "create-db-security-group",
	Short: "Create a PostgreSQL DB security group",
	Long: `Creates a PostgreSQL DB security group.
IMPORTANT: PostgreSQL requires at least one security rule to be defined at creation.
You must provide --cidr for the initial ingress rule.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("db-security-group-name")
		description, _ := cmd.Flags().GetString("description")
		cidr, _ := cmd.Flags().GetString("cidr")
		port, _ := cmd.Flags().GetInt("port")
		minPortFlag, _ := cmd.Flags().GetInt("min-port")
		maxPortFlag, _ := cmd.Flags().GetInt("max-port")

		if name == "" {
			exitWithError("--db-security-group-name is required", nil)
		}
		if cidr == "" {
			exitWithError("--cidr is required (PostgreSQL requires initial rule)", nil)
		}

		// Resolve port range. Priority: --port (single) > --min-port/--max-port > default 5432.
		var minPort, maxPort int
		switch {
		case port > 0:
			minPort = port
			maxPort = port
		case minPortFlag > 0 && maxPortFlag > 0:
			minPort = minPortFlag
			maxPort = maxPortFlag
		default:
			minPort = 5432
			maxPort = 5432
		}

		rule := postgresql.SecurityRule{
			Direction: "INGRESS",
			EtherType: "IPV4",
			CIDR:      cidr,
			Port: postgresql.RulePort{
				PortType: "PORT_RANGE",
				MinPort:  &minPort,
				MaxPort:  &maxPort,
			},
		}

		client := newPostgreSQLClient()
		req := &postgresql.CreateSecurityGroupRequest{
			DBSecurityGroupName: name,
			Description:         description,
			Rules:               []postgresql.SecurityRule{rule},
		}

		result, err := client.CreateSecurityGroup(context.Background(), req)
		if err != nil {
			exitWithError("failed to create security group", err)
		}

		fmt.Printf("Security group created: %s\n", result.DBSecurityGroupID)
	},
}

var authorizePostgreSQLSecurityGroupIngressCmd = &cobra.Command{
	Use:   "authorize-db-security-group-ingress",
	Short: "Authorize ingress rule for PostgreSQL DB security group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		identifier, _ := cmd.Flags().GetString("db-security-group-identifier")
		cidr, _ := cmd.Flags().GetString("cidr")
		description, _ := cmd.Flags().GetString("description")
		port, _ := cmd.Flags().GetInt("port")
		minPortFlag, _ := cmd.Flags().GetInt("min-port")
		maxPortFlag, _ := cmd.Flags().GetInt("max-port")

		// Resolve port range. Priority: --port (single) > --min-port/--max-port > default 5432.
		var minPort, maxPort int
		switch {
		case port > 0:
			minPort = port
			maxPort = port
		case minPortFlag > 0 && maxPortFlag > 0:
			minPort = minPortFlag
			maxPort = maxPortFlag
		default:
			minPort = 5432
			maxPort = 5432
		}

		req := &postgresql.CreateSecurityRuleRequest{
			Description: description,
			CIDR:        cidr,
			Port: postgresql.RulePort{
				PortType: "PORT_RANGE",
				MinPort:  &minPort,
				MaxPort:  &maxPort,
			},
			Direction: "INGRESS",
			EtherType: "IPV4",
		}

		// Assume identifier is ID
		resp, err := client.CreateSecurityRule(context.Background(), identifier, req)
		if err != nil {
			exitWithError("failed to authorize security group ingress", err)
		}

		fmt.Printf("Security rule created: %s\n", resp.RuleID)
	},
}

var deletePostgreSQLSecurityGroupCmd = &cobra.Command{
	Use:   "delete-db-security-group",
	Short: "Delete a PostgreSQL DB security group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("db-security-group-identifier")
		if groupID == "" {
			exitWithError("--db-security-group-identifier is required", nil)
		}

		client := newPostgreSQLClient()
		_, err := client.DeleteSecurityGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to delete security group", err)
		}

		fmt.Printf("Security group deleted successfully\n")
	},
}
