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

	rdsPostgreSQLCmd.AddCommand(authorizePostgreSQLSecurityGroupIngressCmd)
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().String("db-security-group-identifier", "", "DB security group identifier (Required)")
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().String("cidr", "", "CIDR block to allow (Required)")
	authorizePostgreSQLSecurityGroupIngressCmd.Flags().Int("port", 5432, "Port to allow")
	authorizePostgreSQLSecurityGroupIngressCmd.MarkFlagRequired("db-security-group-identifier")
	authorizePostgreSQLSecurityGroupIngressCmd.MarkFlagRequired("cidr")
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

var authorizePostgreSQLSecurityGroupIngressCmd = &cobra.Command{
	Use:   "authorize-db-security-group-ingress",
	Short: "Authorize ingress rule for PostgreSQL DB security group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		identifier, _ := cmd.Flags().GetString("db-security-group-identifier")
		cidr, _ := cmd.Flags().GetString("cidr")
		port, _ := cmd.Flags().GetInt("port")

		req := &postgresql.CreateSecurityRuleRequest{
			CIDR: cidr,
			Port: postgresql.RulePort{
				PortType: "PORT",
			},
			Direction: "INGRESS",
			EtherType: "IPV4",
		}

		// Handle port filtering
		if port > 0 {
			req.Port.MinPort = &port
			req.Port.MaxPort = &port
		}

		// Assume identifier is ID
		resp, err := client.CreateSecurityRule(context.Background(), identifier, req)
		if err != nil {
			exitWithError("failed to authorize security group ingress", err)
		}

		fmt.Printf("Security rule created: %s\n", resp.RuleID)
	},
}
