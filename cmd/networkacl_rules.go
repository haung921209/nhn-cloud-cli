package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/networkacl"
	"github.com/spf13/cobra"
)

func init() {
	networkACLCmd.AddCommand(aclDescribeRulesCmd)
	networkACLCmd.AddCommand(aclGetRuleCmd)
	networkACLCmd.AddCommand(aclCreateRuleCmd)
	networkACLCmd.AddCommand(aclUpdateRuleCmd)
	networkACLCmd.AddCommand(aclDeleteRuleCmd)

	aclDescribeRulesCmd.Flags().String("acl-id", "", "Filter by ACL ID")

	aclGetRuleCmd.Flags().String("rule-id", "", "Rule ID (required)")
	aclGetRuleCmd.MarkFlagRequired("rule-id")

	aclCreateRuleCmd.Flags().String("acl-id", "", "ACL ID (required)")
	aclCreateRuleCmd.Flags().String("description", "", "Description")
	aclCreateRuleCmd.Flags().String("protocol", "", "Protocol (tcp/udp/icmp/any)")
	aclCreateRuleCmd.Flags().String("ethertype", "IPv4", "Ethertype (IPv4/IPv6)")
	aclCreateRuleCmd.Flags().String("src-ip", "", "Source IP CIDR")
	aclCreateRuleCmd.Flags().String("dst-ip", "", "Destination IP CIDR")
	aclCreateRuleCmd.Flags().Int("src-port-min", 0, "Source port min")
	aclCreateRuleCmd.Flags().Int("src-port-max", 0, "Source port max")
	aclCreateRuleCmd.Flags().Int("dst-port-min", 0, "Destination port min")
	aclCreateRuleCmd.Flags().Int("dst-port-max", 0, "Destination port max")
	aclCreateRuleCmd.Flags().String("policy", "allow", "Policy (allow/deny)")
	aclCreateRuleCmd.Flags().Int("order", 1, "Rule order (priority)")
	aclCreateRuleCmd.MarkFlagRequired("acl-id")
	aclCreateRuleCmd.MarkFlagRequired("policy")

	aclUpdateRuleCmd.Flags().String("rule-id", "", "Rule ID (required)")
	aclUpdateRuleCmd.Flags().String("description", "", "Description")
	aclUpdateRuleCmd.Flags().String("protocol", "", "Protocol")
	aclUpdateRuleCmd.Flags().String("src-ip", "", "Source IP CIDR")
	aclUpdateRuleCmd.Flags().String("dst-ip", "", "Destination IP CIDR")
	aclUpdateRuleCmd.Flags().Int("src-port-min", 0, "Source port min")
	aclUpdateRuleCmd.Flags().Int("src-port-max", 0, "Source port max")
	aclUpdateRuleCmd.Flags().Int("dst-port-min", 0, "Destination port min")
	aclUpdateRuleCmd.Flags().Int("dst-port-max", 0, "Destination port max")
	aclUpdateRuleCmd.Flags().String("policy", "", "Policy")
	aclUpdateRuleCmd.Flags().Int("order", 0, "Rule order")
	aclUpdateRuleCmd.MarkFlagRequired("rule-id")

	aclDeleteRuleCmd.Flags().String("rule-id", "", "Rule ID (required)")
	aclDeleteRuleCmd.MarkFlagRequired("rule-id")
}

var aclDescribeRulesCmd = &cobra.Command{
	Use:     "describe-rules",
	Aliases: []string{"list-rules"},
	Short:   "List ACL rules",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		aclID, _ := cmd.Flags().GetString("acl-id")

		var result *networkacl.ListACLRulesOutput
		var err error

		if aclID != "" {
			result, err = client.ListACLRulesByACL(ctx, aclID)
		} else {
			result, err = client.ListACLRules(ctx)
		}

		if err != nil {
			exitWithError("Failed to list ACL rules", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tACL_ID\tPROTOCOL\tPOLICY\tORDER\tSRC_IP\tDST_IP")
		for _, rule := range result.ACLRules {
			protocol := rule.Protocol
			if protocol == "" {
				protocol = "any"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
				rule.ID, rule.ACLID, protocol, rule.Policy, rule.OrderNum, rule.SrcIPPrefix, rule.DstIPPrefix)
		}
		w.Flush()
	},
}

var aclGetRuleCmd = &cobra.Command{
	Use:     "describe-rule",
	Aliases: []string{"get-rule"},
	Short:   "Get ACL rule details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("rule-id")

		result, err := client.GetACLRule(ctx, id)
		if err != nil {
			exitWithError("Failed to get ACL rule", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		rule := result.ACLRule
		fmt.Printf("ID:          %s\n", rule.ID)
		fmt.Printf("ACL ID:      %s\n", rule.ACLID)
		fmt.Printf("Description: %s\n", rule.Description)
		fmt.Printf("Protocol:    %s\n", rule.Protocol)
		fmt.Printf("Ethertype:   %s\n", rule.EtherType)
		fmt.Printf("Policy:      %s\n", rule.Policy)
		fmt.Printf("Order:       %d\n", rule.OrderNum)
		fmt.Printf("Src IP:      %s\n", rule.SrcIPPrefix)
		fmt.Printf("Dst IP:      %s\n", rule.DstIPPrefix)
		if rule.SrcPortMin != nil {
			fmt.Printf("Src Port:    %d-%d\n", *rule.SrcPortMin, *rule.SrcPortMax)
		}
		if rule.DstPortMin != nil {
			fmt.Printf("Dst Port:    %d-%d\n", *rule.DstPortMin, *rule.DstPortMax)
		}
	},
}

var aclCreateRuleCmd = &cobra.Command{
	Use:   "create-rule",
	Short: "Create a new ACL rule",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		aclID, _ := cmd.Flags().GetString("acl-id")
		description, _ := cmd.Flags().GetString("description")
		protocol, _ := cmd.Flags().GetString("protocol")
		ethertype, _ := cmd.Flags().GetString("ethertype")
		srcIP, _ := cmd.Flags().GetString("src-ip")
		dstIP, _ := cmd.Flags().GetString("dst-ip")
		srcPortMin, _ := cmd.Flags().GetInt("src-port-min")
		srcPortMax, _ := cmd.Flags().GetInt("src-port-max")
		dstPortMin, _ := cmd.Flags().GetInt("dst-port-min")
		dstPortMax, _ := cmd.Flags().GetInt("dst-port-max")
		policy, _ := cmd.Flags().GetString("policy")
		order, _ := cmd.Flags().GetInt("order")

		input := &networkacl.CreateACLRuleInput{
			ACLID:       aclID,
			Description: description,
			Protocol:    protocol,
			EtherType:   ethertype,
			SrcIPPrefix: srcIP,
			DstIPPrefix: dstIP,
			Policy:      policy,
			OrderNum:    order,
		}

		if srcPortMin > 0 {
			input.SrcPortMin = &srcPortMin
			input.SrcPortMax = &srcPortMax
		}
		if dstPortMin > 0 {
			input.DstPortMin = &dstPortMin
			input.DstPortMax = &dstPortMax
		}

		result, err := client.CreateACLRule(ctx, input)
		if err != nil {
			exitWithError("Failed to create ACL rule", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("ACL rule created: %s\n", result.ACLRule.ID)
	},
}

var aclUpdateRuleCmd = &cobra.Command{
	Use:   "update-rule",
	Short: "Update an ACL rule",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("rule-id")
		description, _ := cmd.Flags().GetString("description")
		protocol, _ := cmd.Flags().GetString("protocol")
		srcIP, _ := cmd.Flags().GetString("src-ip")
		dstIP, _ := cmd.Flags().GetString("dst-ip")
		srcPortMin, _ := cmd.Flags().GetInt("src-port-min")
		srcPortMax, _ := cmd.Flags().GetInt("src-port-max")
		dstPortMin, _ := cmd.Flags().GetInt("dst-port-min")
		dstPortMax, _ := cmd.Flags().GetInt("dst-port-max")
		policy, _ := cmd.Flags().GetString("policy")
		order, _ := cmd.Flags().GetInt("order")

		input := &networkacl.UpdateACLRuleInput{
			Description: description,
			Protocol:    protocol,
			SrcIPPrefix: srcIP,
			DstIPPrefix: dstIP,
			Policy:      policy,
		}

		if srcPortMin > 0 {
			input.SrcPortMin = &srcPortMin
			input.SrcPortMax = &srcPortMax
		}
		if dstPortMin > 0 {
			input.DstPortMin = &dstPortMin
			input.DstPortMax = &dstPortMax
		}
		if order > 0 {
			input.OrderNum = &order
		}

		result, err := client.UpdateACLRule(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update ACL rule", err)
		}

		fmt.Printf("ACL rule updated: %s\n", result.ACLRule.ID)
	},
}

var aclDeleteRuleCmd = &cobra.Command{
	Use:   "delete-rule",
	Short: "Delete an ACL rule",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("rule-id")

		if err := client.DeleteACLRule(ctx, id); err != nil {
			exitWithError("Failed to delete ACL rule", err)
		}
		fmt.Printf("ACL rule %s deleted\n", id)
	},
}
