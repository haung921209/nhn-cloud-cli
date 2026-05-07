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
// Security Group Commands
// ============================================================================

var describeMariaDBSecurityGroupsCmd = &cobra.Command{
	Use:   "describe-db-security-groups",
	Short: "Describe MariaDB DB security groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListSecurityGroups(context.Background())
		if err != nil {
			exitWithError("failed to list security groups", err)
		}

		if output == "json" {
			mariadbPrintJSON(result)
		} else {
			mariadbPrintSecurityGroupList(result)
		}
	},
}

var createMariaDBSecurityGroupCmd = &cobra.Command{
	Use:   "create-db-security-group",
	Short: "Create a MariaDB DB security group",
	Long: `Creates a MariaDB DB security group.
IMPORTANT: MariaDB requires at least one security rule to be defined at creation.
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
			exitWithError("--cidr is required (MariaDB requires initial rule)", nil)
		}

		// Resolve port range. Priority: --port (single) > --min-port/--max-port > default 3306.
		var minPort, maxPort int
		switch {
		case port > 0:
			minPort = port
			maxPort = port
		case minPortFlag > 0 && maxPortFlag > 0:
			minPort = minPortFlag
			maxPort = maxPortFlag
		default:
			minPort = 3306
			maxPort = 3306
		}

		rule := mariadb.SecurityRule{
			Direction: "INGRESS",
			EtherType: "IPV4",
			CIDR:      cidr,
			Port: mariadb.RulePort{
				PortType: "PORT_RANGE",
				MinPort:  &minPort,
				MaxPort:  &maxPort,
			},
		}

		client := newMariaDBClient()
		req := &mariadb.CreateSecurityGroupRequest{
			DBSecurityGroupName: name,
			Description:         description,
			Rules:               []mariadb.SecurityRule{rule},
		}

		result, err := client.CreateSecurityGroup(context.Background(), req)
		if err != nil {
			exitWithError("failed to create security group", err)
		}

		fmt.Printf("Security group created: %s\n", result.DBSecurityGroupID)
	},
}

var authorizeMariaDBSecurityGroupIngressCmd = &cobra.Command{
	Use:   "authorize-db-security-group-ingress",
	Short: "Authorize ingress rule for MariaDB security group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("db-security-group-identifier")
		cidr, _ := cmd.Flags().GetString("cidr")
		description, _ := cmd.Flags().GetString("description")
		port, _ := cmd.Flags().GetInt("port")
		minPortFlag, _ := cmd.Flags().GetInt("min-port")
		maxPortFlag, _ := cmd.Flags().GetInt("max-port")

		if groupID == "" {
			exitWithError("--db-security-group-identifier is required", nil)
		}
		if cidr == "" {
			exitWithError("--cidr is required", nil)
		}

		client := newMariaDBClient()

		// Resolve port range. Priority: --port (single) > --min-port/--max-port > default 3306.
		var minPort, maxPort int
		switch {
		case port > 0:
			minPort = port
			maxPort = port
		case minPortFlag > 0 && maxPortFlag > 0:
			minPort = minPortFlag
			maxPort = maxPortFlag
		default:
			minPort = 3306
			maxPort = 3306
		}

		req := &mariadb.CreateSecurityRuleRequest{
			Description: description,
			Direction:   "INGRESS",
			EtherType:   "IPV4",
			CIDR:        cidr,
			Port: mariadb.RulePort{
				PortType: "PORT_RANGE",
				MinPort:  &minPort,
				MaxPort:  &maxPort,
			},
		}

		result, err := client.CreateSecurityRule(context.Background(), groupID, req)
		if err != nil {
			exitWithError("failed to authorize security group ingress", err)
		}

		fmt.Printf("Security rule created: %s\n", result.RuleID)
	},
}

var deleteMariaDBSecurityGroupCmd = &cobra.Command{
	Use:   "delete-db-security-group",
	Short: "Delete a MariaDB DB security group",
	Run: func(cmd *cobra.Command, args []string) {
		groupID, _ := cmd.Flags().GetString("db-security-group-identifier")
		if groupID == "" {
			exitWithError("--db-security-group-identifier is required", nil)
		}

		client := newMariaDBClient()
		_, err := client.DeleteSecurityGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to delete security group", err)
		}

		fmt.Printf("Security group deleted successfully\n")
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func mariadbPrintSecurityGroupList(result *mariadb.ListSecurityGroupsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tRULES")
	for _, sg := range result.DBSecurityGroups {
		fmt.Fprintf(w, "%s\t%s\t%d\n",
			sg.DBSecurityGroupID,
			sg.DBSecurityGroupName,
			len(sg.Rules),
		)
	}
	w.Flush()
}

func init() {
	rdsMariaDBCmd.AddCommand(describeMariaDBSecurityGroupsCmd)
	rdsMariaDBCmd.AddCommand(createMariaDBSecurityGroupCmd)
	rdsMariaDBCmd.AddCommand(authorizeMariaDBSecurityGroupIngressCmd)
	rdsMariaDBCmd.AddCommand(deleteMariaDBSecurityGroupCmd)

	// create flags
	createMariaDBSecurityGroupCmd.Flags().String("db-security-group-name", "", "Security group name (required)")
	createMariaDBSecurityGroupCmd.Flags().String("description", "", "Description")
	createMariaDBSecurityGroupCmd.Flags().String("cidr", "", "Initial CIDR block (required for MariaDB)")
	createMariaDBSecurityGroupCmd.Flags().Int("port", 0, "Specific port (e.g. 3306, 13306)")
	createMariaDBSecurityGroupCmd.Flags().Int("min-port", 0, "Minimum port for range")
	createMariaDBSecurityGroupCmd.Flags().Int("max-port", 0, "Maximum port for range")

	// authorize ingress flags
	authorizeMariaDBSecurityGroupIngressCmd.Flags().String("db-security-group-identifier", "", "Security group identifier (required)")
	authorizeMariaDBSecurityGroupIngressCmd.Flags().String("cidr", "", "CIDR block (required, e.g., 0.0.0.0/0)")
	authorizeMariaDBSecurityGroupIngressCmd.Flags().String("description", "", "Rule description")
	authorizeMariaDBSecurityGroupIngressCmd.Flags().Int("port", 0, "Specific port (e.g. 3306, 13306)")
	authorizeMariaDBSecurityGroupIngressCmd.Flags().Int("min-port", 0, "Minimum port for range")
	authorizeMariaDBSecurityGroupIngressCmd.Flags().Int("max-port", 0, "Maximum port for range")

	// delete flags
	deleteMariaDBSecurityGroupCmd.Flags().String("db-security-group-identifier", "", "Security group identifier (required)")
}
