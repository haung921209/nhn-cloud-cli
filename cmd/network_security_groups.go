package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/securitygroup"
	"github.com/spf13/cobra"
)

func init() {
	networkCmd.AddCommand(networkDescribeSecurityGroupsCmd)
	networkCmd.AddCommand(networkCreateSecurityGroupCmd)
	networkCmd.AddCommand(networkDeleteSecurityGroupCmd)
	networkCmd.AddCommand(networkAuthorizeSecurityGroupIngressCmd)
	networkCmd.AddCommand(networkDeleteSecurityGroupRuleCmd)

	networkDescribeSecurityGroupsCmd.Flags().String("group-id", "", "Security Group ID")

	networkCreateSecurityGroupCmd.Flags().String("name", "", "Security group name (required)")
	networkCreateSecurityGroupCmd.Flags().String("description", "", "Security group description")
	networkCreateSecurityGroupCmd.MarkFlagRequired("name")

	networkDeleteSecurityGroupCmd.Flags().String("group-id", "", "Security Group ID (required)")
	networkDeleteSecurityGroupCmd.MarkFlagRequired("group-id")

	networkAuthorizeSecurityGroupIngressCmd.Flags().String("group-id", "", "Security group ID (required)")
	networkAuthorizeSecurityGroupIngressCmd.Flags().String("direction", "ingress", "Direction (ingress/egress)")
	networkAuthorizeSecurityGroupIngressCmd.Flags().String("ethertype", "IPv4", "Ethertype (IPv4/IPv6)")
	networkAuthorizeSecurityGroupIngressCmd.Flags().String("protocol", "", "Protocol (tcp/udp/icmp)")
	networkAuthorizeSecurityGroupIngressCmd.Flags().Int("port-min", 0, "Minimum port")
	networkAuthorizeSecurityGroupIngressCmd.Flags().Int("port-max", 0, "Maximum port")
	networkAuthorizeSecurityGroupIngressCmd.Flags().String("remote-ip", "0.0.0.0/0", "Remote IP prefix")
	networkAuthorizeSecurityGroupIngressCmd.MarkFlagRequired("group-id")

	networkDeleteSecurityGroupRuleCmd.Flags().String("rule-id", "", "Security Group Rule ID (required)")
	networkDeleteSecurityGroupRuleCmd.MarkFlagRequired("rule-id")
}

var networkDescribeSecurityGroupsCmd = &cobra.Command{
	Use:   "describe-security-groups",
	Short: "Describe security groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := securitygroup.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()
		groupID, _ := cmd.Flags().GetString("group-id")

		if groupID != "" {
			result, err := client.GetSecurityGroup(ctx, groupID)
			if err != nil {
				exitWithError("Failed to get security group", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			sg := result.SecurityGroup
			fmt.Printf("ID:          %s\n", sg.ID)
			fmt.Printf("Name:        %s\n", sg.Name)
			fmt.Printf("Description: %s\n", sg.Description)
			fmt.Printf("\nRules:\n")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "  DIRECTION\tPROTOCOL\tPORT\tREMOTE")
			for _, rule := range sg.Rules {
				port := ""
				if rule.PortRangeMin != nil {
					port = fmt.Sprintf("%d-%d", *rule.PortRangeMin, *rule.PortRangeMax)
				}
				remote := rule.RemoteIPPrefix
				if remote == "" {
					remote = rule.RemoteGroupID
				}
				protocol := ""
				if rule.Protocol != nil {
					protocol = *rule.Protocol
				}
				fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n", rule.Direction, protocol, port, remote)
			}
			w.Flush()
		} else {
			result, err := client.ListSecurityGroups(ctx)
			if err != nil {
				exitWithError("Failed to list security groups", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION")
			for _, sg := range result.SecurityGroups {
				fmt.Fprintf(w, "%s\t%s\t%s\n", sg.ID, sg.Name, sg.Description)
			}
			w.Flush()
		}
	},
}

var networkCreateSecurityGroupCmd = &cobra.Command{
	Use:   "create-security-group",
	Short: "Create a new security group",
	Run: func(cmd *cobra.Command, args []string) {
		client := securitygroup.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &securitygroup.CreateSecurityGroupInput{
			Name:        name,
			Description: description,
		}

		result, err := client.CreateSecurityGroup(ctx, input)
		if err != nil {
			exitWithError("Failed to create security group", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Security group created: %s\n", result.SecurityGroup.ID)
	},
}

var networkDeleteSecurityGroupCmd = &cobra.Command{
	Use:   "delete-security-group",
	Short: "Delete a security group",
	Run: func(cmd *cobra.Command, args []string) {
		client := securitygroup.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()
		groupID, _ := cmd.Flags().GetString("group-id")

		if err := client.DeleteSecurityGroup(ctx, groupID); err != nil {
			exitWithError("Failed to delete security group", err)
		}

		fmt.Printf("Security group %s deleted\n", groupID)
	},
}

var networkAuthorizeSecurityGroupIngressCmd = &cobra.Command{
	Use:   "authorize-security-group-ingress",
	Short: "Create a security group rule (ingress/egress)",
	Long:  "Authorize ingress (or egress using --direction) traffic. Maps to creating a rule.",
	Run: func(cmd *cobra.Command, args []string) {
		client := securitygroup.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		sgID, _ := cmd.Flags().GetString("group-id")
		direction, _ := cmd.Flags().GetString("direction")
		ethertype, _ := cmd.Flags().GetString("ethertype")
		protocol, _ := cmd.Flags().GetString("protocol")
		portMin, _ := cmd.Flags().GetInt("port-min")
		portMax, _ := cmd.Flags().GetInt("port-max")
		remoteIP, _ := cmd.Flags().GetString("remote-ip")

		input := &securitygroup.CreateRuleInput{
			SecurityGroupID: sgID,
			Direction:       direction,
			EtherType:       ethertype,
			Protocol:        protocol,
			PortRangeMin:    &portMin,
			PortRangeMax:    &portMax,
			RemoteIPPrefix:  remoteIP,
		}

		result, err := client.CreateRule(ctx, input)
		if err != nil {
			exitWithError("Failed to authorize security group ingress", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Security group rule created: %s\n", result.SecurityGroupRule.ID)
	},
}

var networkDeleteSecurityGroupRuleCmd = &cobra.Command{
	Use:   "delete-security-group-rule",
	Short: "Delete a security group rule",
	Run: func(cmd *cobra.Command, args []string) {
		client := securitygroup.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()
		ruleID, _ := cmd.Flags().GetString("rule-id")

		if err := client.DeleteRule(ctx, ruleID); err != nil {
			exitWithError("Failed to delete security group rule", err)
		}

		fmt.Printf("Security group rule %s deleted\n", ruleID)
	},
}
