package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// Security Group Commands
// ============================================================================

var securityGroupCmd = &cobra.Command{
	Use:     "security-group",
	Aliases: []string{"sg"},
	Short:   "Manage DB security groups",
}

var listSecurityGroupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DB security groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListSecurityGroups(context.Background())
		if err != nil {
			exitWithError("failed to list security groups", err)
		}
		printSecurityGroups(result)
	},
}

var getSecurityGroupCmd = &cobra.Command{
	Use:   "get [security-group-id]",
	Short: "Get details of a security group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.GetSecurityGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get security group", err)
		}
		printSecurityGroupDetail(result)
	},
}

var createSecurityGroupCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new security group",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		input := &mysql.CreateSecurityGroupInput{
			DBSecurityGroupName: name,
			Description:         description,
		}

		client := newMySQLClient()
		result, err := client.CreateSecurityGroup(context.Background(), input)
		if err != nil {
			exitWithError("failed to create security group", err)
		}
		fmt.Printf("Security group created. ID: %s\n", result.DBSecurityGroupID)
	},
}

var updateSecurityGroupCmd = &cobra.Command{
	Use:   "update [security-group-id]",
	Short: "Update a security group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" && description == "" {
			exitWithError("at least --name or --description is required", nil)
		}

		input := &mysql.UpdateSecurityGroupInput{
			DBSecurityGroupName: name,
			Description:         description,
		}

		client := newMySQLClient()
		_, err := client.UpdateSecurityGroup(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to update security group", err)
		}
		fmt.Println("Security group updated successfully.")
	},
}

var deleteSecurityGroupCmd = &cobra.Command{
	Use:   "delete [security-group-id]",
	Short: "Delete a security group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		_, err := client.DeleteSecurityGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete security group", err)
		}
		fmt.Println("Security group deleted successfully.")
	},
}

// Security Group Rules
var sgRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Manage security group rules",
}

var createSGRuleCmd = &cobra.Command{
	Use:   "create [security-group-id]",
	Short: "Create a security group rule",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		direction, _ := cmd.Flags().GetString("direction")
		etherType, _ := cmd.Flags().GetString("ether-type")
		port, _ := cmd.Flags().GetInt("port")
		cidr, _ := cmd.Flags().GetString("cidr")

		if direction == "" || cidr == "" {
			exitWithError("--direction and --cidr are required", nil)
		}

		input := &mysql.CreateSecurityGroupRuleInput{
			Direction: direction,
			EtherType: etherType,
			Port:      mysql.Port{PortType: "TCP", MinPort: &port, MaxPort: &port},
			CIDR:      cidr,
		}

		client := newMySQLClient()
		result, err := client.CreateSecurityGroupRule(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create security group rule", err)
		}
		fmt.Printf("Security group rule created. ID: %s\n", result.RuleID)
	},
}

var updateSGRuleCmd = &cobra.Command{
	Use:   "update [security-group-id] [rule-id]",
	Short: "Update a security group rule",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		direction, _ := cmd.Flags().GetString("direction")
		etherType, _ := cmd.Flags().GetString("ether-type")
		port, _ := cmd.Flags().GetInt("port")
		cidr, _ := cmd.Flags().GetString("cidr")

		input := &mysql.UpdateSecurityGroupRuleInput{
			Direction: direction,
			EtherType: etherType,
			Port:      &mysql.Port{PortType: "TCP", MinPort: &port, MaxPort: &port},
			CIDR:      cidr,
		}

		client := newMySQLClient()
		_, err := client.UpdateSecurityGroupRule(context.Background(), args[0], args[1], input)
		if err != nil {
			exitWithError("failed to update security group rule", err)
		}
		fmt.Println("Security group rule updated successfully.")
	},
}

var deleteSGRuleCmd = &cobra.Command{
	Use:   "delete [security-group-id] [rule-id]",
	Short: "Delete a security group rule",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		_, err := client.DeleteSecurityGroupRule(context.Background(), args[0], args[1])
		if err != nil {
			exitWithError("failed to delete security group rule", err)
		}
		fmt.Println("Security group rule deleted successfully.")
	},
}

// ============================================================================
// Parameter Group Commands
// ============================================================================

var parameterGroupCmd = &cobra.Command{
	Use:     "parameter-group",
	Aliases: []string{"pg"},
	Short:   "Manage parameter groups",
}

var listParameterGroupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all parameter groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListParameterGroups(context.Background())
		if err != nil {
			exitWithError("failed to list parameter groups", err)
		}
		printParameterGroups(result)
	},
}

var getParameterGroupCmd = &cobra.Command{
	Use:   "get [parameter-group-id]",
	Short: "Get details of a parameter group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.GetParameterGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get parameter group", err)
		}
		printParameterGroupDetail(result)
	},
}

var createParameterGroupCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new parameter group",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		dbVersion, _ := cmd.Flags().GetString("version")

		if name == "" || dbVersion == "" {
			exitWithError("--name and --version are required", nil)
		}

		input := &mysql.CreateParameterGroupInput{
			ParameterGroupName: name,
			Description:        description,
			DBVersion:          dbVersion,
		}

		client := newMySQLClient()
		result, err := client.CreateParameterGroup(context.Background(), input)
		if err != nil {
			exitWithError("failed to create parameter group", err)
		}
		fmt.Printf("Parameter group created. ID: %s\n", result.ParameterGroupID)
	},
}

var copyParameterGroupCmd = &cobra.Command{
	Use:   "copy [source-parameter-group-id]",
	Short: "Copy a parameter group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		input := &mysql.CopyParameterGroupInput{
			ParameterGroupName: name,
			Description:        description,
		}

		client := newMySQLClient()
		result, err := client.CopyParameterGroup(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to copy parameter group", err)
		}
		fmt.Printf("Parameter group copied. ID: %s\n", result.ParameterGroupID)
	},
}

var updateParameterGroupCmd = &cobra.Command{
	Use:   "update [parameter-group-id]",
	Short: "Update a parameter group name/description",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" && description == "" {
			exitWithError("at least --name or --description is required", nil)
		}

		input := &mysql.UpdateParameterGroupInput{
			ParameterGroupName: name,
			Description:        description,
		}

		client := newMySQLClient()
		_, err := client.UpdateParameterGroup(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to update parameter group", err)
		}
		fmt.Println("Parameter group updated successfully.")
	},
}

var modifyParametersCmd = &cobra.Command{
	Use:   "modify [parameter-group-id]",
	Short: "Modify parameters in a parameter group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		params, _ := cmd.Flags().GetStringToString("set")

		if len(params) == 0 {
			exitWithError("--set is required (e.g., --set max_connections=200)", nil)
		}

		var paramValues []struct {
			ParameterID string `json:"parameterId"`
			Value       string `json:"value"`
		}
		for name, value := range params {
			paramValues = append(paramValues, struct {
				ParameterID string `json:"parameterId"`
				Value       string `json:"value"`
			}{
				ParameterID: name,
				Value:       value,
			})
		}

		input := &mysql.ModifyParametersInput{
			ModifiedParameters: paramValues,
		}

		client := newMySQLClient()
		_, err := client.ModifyParameters(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to modify parameters", err)
		}
		fmt.Println("Parameters modified successfully.")
	},
}

var resetParameterGroupCmd = &cobra.Command{
	Use:   "reset [parameter-group-id]",
	Short: "Reset all parameters to default values",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		_, err := client.ResetParameterGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to reset parameter group", err)
		}
		fmt.Println("Parameter group reset successfully.")
	},
}

var deleteParameterGroupCmd = &cobra.Command{
	Use:   "delete [parameter-group-id]",
	Short: "Delete a parameter group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		_, err := client.DeleteParameterGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete parameter group", err)
		}
		fmt.Println("Parameter group deleted successfully.")
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	// Security Group commands
	rdsMySQLCmd.AddCommand(securityGroupCmd)
	securityGroupCmd.AddCommand(listSecurityGroupsCmd)
	securityGroupCmd.AddCommand(getSecurityGroupCmd)
	securityGroupCmd.AddCommand(createSecurityGroupCmd)
	securityGroupCmd.AddCommand(updateSecurityGroupCmd)
	securityGroupCmd.AddCommand(deleteSecurityGroupCmd)
	securityGroupCmd.AddCommand(sgRuleCmd)

	createSecurityGroupCmd.Flags().String("name", "", "Security group name (required)")
	createSecurityGroupCmd.Flags().String("description", "", "Description")

	updateSecurityGroupCmd.Flags().String("name", "", "New name")
	updateSecurityGroupCmd.Flags().String("description", "", "New description")

	// Security Group Rule commands
	sgRuleCmd.AddCommand(createSGRuleCmd)
	sgRuleCmd.AddCommand(updateSGRuleCmd)
	sgRuleCmd.AddCommand(deleteSGRuleCmd)

	createSGRuleCmd.Flags().String("direction", "ingress", "Direction (ingress/egress)")
	createSGRuleCmd.Flags().String("ether-type", "IPV4", "Ether type (IPV4/IPV6)")
	createSGRuleCmd.Flags().Int("port", 3306, "Port number")
	createSGRuleCmd.Flags().String("cidr", "", "CIDR block (required)")

	updateSGRuleCmd.Flags().String("direction", "", "Direction (ingress/egress)")
	updateSGRuleCmd.Flags().String("ether-type", "", "Ether type (IPV4/IPV6)")
	updateSGRuleCmd.Flags().Int("port", 0, "Port number")
	updateSGRuleCmd.Flags().String("cidr", "", "CIDR block")

	// Parameter Group commands
	rdsMySQLCmd.AddCommand(parameterGroupCmd)
	parameterGroupCmd.AddCommand(listParameterGroupsCmd)
	parameterGroupCmd.AddCommand(getParameterGroupCmd)
	parameterGroupCmd.AddCommand(createParameterGroupCmd)
	parameterGroupCmd.AddCommand(copyParameterGroupCmd)
	parameterGroupCmd.AddCommand(updateParameterGroupCmd)
	parameterGroupCmd.AddCommand(modifyParametersCmd)
	parameterGroupCmd.AddCommand(resetParameterGroupCmd)
	parameterGroupCmd.AddCommand(deleteParameterGroupCmd)

	listParameterGroupsCmd.Flags().String("version", "", "Filter by DB version")

	createParameterGroupCmd.Flags().String("name", "", "Parameter group name (required)")
	createParameterGroupCmd.Flags().String("description", "", "Description")
	createParameterGroupCmd.Flags().String("version", "", "DB version (required)")

	copyParameterGroupCmd.Flags().String("name", "", "New parameter group name (required)")
	copyParameterGroupCmd.Flags().String("description", "", "Description")

	updateParameterGroupCmd.Flags().String("name", "", "New name")
	updateParameterGroupCmd.Flags().String("description", "", "New description")

	modifyParametersCmd.Flags().StringToString("set", nil, "Parameters to set (e.g., --set max_connections=200)")
}

// ============================================================================
// Print Functions
// ============================================================================

func printSecurityGroups(result *mysql.ListSecurityGroupsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tRULES")
	for _, sg := range result.DBSecurityGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
			sg.DBSecurityGroupID, sg.DBSecurityGroupName, sg.Description, len(sg.Rules))
	}
	w.Flush()
}

func printSecurityGroupDetail(result *mysql.SecurityGroupOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	fmt.Printf("ID:          %s\n", result.DBSecurityGroup.DBSecurityGroupID)
	fmt.Printf("Name:        %s\n", result.DBSecurityGroup.DBSecurityGroupName)
	fmt.Printf("Description: %s\n", result.DBSecurityGroup.Description)
	fmt.Printf("Created:     %s\n", result.DBSecurityGroup.CreatedYmdt)
	fmt.Printf("Updated:     %s\n", result.DBSecurityGroup.UpdatedYmdt)
	fmt.Println("\nRules:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  ID\tDIRECTION\tTYPE\tCIDR")
	for _, rule := range result.DBSecurityGroup.Rules {
		fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n",
			rule.RuleID, rule.Direction, rule.EtherType, rule.CIDR)
	}
	w.Flush()
}

func printParameterGroups(result *mysql.ListParameterGroupsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tVERSION\tSTATUS")
	for _, pg := range result.ParameterGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			pg.ParameterGroupID, pg.ParameterGroupName, pg.DBVersion, pg.ParameterGroupStatus)
	}
	w.Flush()
}

func printParameterGroupDetail(result *mysql.ParameterGroupOutput) {
	if output == "json" {
		printJSON(result)
		return
	}

	fmt.Printf("ID:          %s\n", result.ParameterGroupID)
	fmt.Printf("Name:        %s\n", result.ParameterGroupName)
	fmt.Printf("Description: %s\n", result.Description)
	fmt.Printf("Version:     %s\n", result.DBVersion)
	fmt.Printf("Status:      %s\n", result.ParameterGroupStatus)
	fmt.Printf("Created:     %s\n", result.CreatedYmdt)
	fmt.Printf("Updated:     %s\n", result.UpdatedYmdt)

	if len(result.Parameters) > 0 {
		fmt.Printf("\nParameters (%d):\n", len(result.Parameters))
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  NAME\tVALUE\tDEFAULT\tUPDATE TYPE")
		for _, p := range result.Parameters {
			fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n",
				p.ParameterName, p.Value, p.DefaultValue, p.UpdateType)
		}
		w.Flush()
	}
}
