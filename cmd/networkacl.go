package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/networkacl"
	"github.com/spf13/cobra"
)

var networkACLCmd = &cobra.Command{
	Use:     "network-acl",
	Aliases: []string{"acl", "nacl"},
	Short:   "Manage Network ACLs",
}

var aclRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Manage ACL rules",
}

var aclBindingCmd = &cobra.Command{
	Use:   "binding",
	Short: "Manage ACL bindings",
}

func init() {
	rootCmd.AddCommand(networkACLCmd)

	// ACL commands
	networkACLCmd.AddCommand(aclListCmd)
	networkACLCmd.AddCommand(aclGetCmd)
	networkACLCmd.AddCommand(aclCreateCmd)
	networkACLCmd.AddCommand(aclUpdateCmd)
	networkACLCmd.AddCommand(aclDeleteCmd)
	networkACLCmd.AddCommand(aclRuleCmd)
	networkACLCmd.AddCommand(aclBindingCmd)

	aclCreateCmd.Flags().String("name", "", "ACL name (required)")
	aclCreateCmd.Flags().String("description", "", "ACL description")
	aclCreateCmd.MarkFlagRequired("name")

	aclUpdateCmd.Flags().String("name", "", "ACL name")
	aclUpdateCmd.Flags().String("description", "", "ACL description")

	// Rule commands
	aclRuleCmd.AddCommand(aclRuleListCmd)
	aclRuleCmd.AddCommand(aclRuleGetCmd)
	aclRuleCmd.AddCommand(aclRuleCreateCmd)
	aclRuleCmd.AddCommand(aclRuleUpdateCmd)
	aclRuleCmd.AddCommand(aclRuleDeleteCmd)

	aclRuleListCmd.Flags().String("acl-id", "", "Filter by ACL ID")

	aclRuleCreateCmd.Flags().String("acl-id", "", "ACL ID (required)")
	aclRuleCreateCmd.Flags().String("description", "", "Rule description")
	aclRuleCreateCmd.Flags().String("protocol", "", "Protocol (tcp/udp/icmp or empty for any)")
	aclRuleCreateCmd.Flags().String("ethertype", "IPv4", "Ethertype (IPv4/IPv6)")
	aclRuleCreateCmd.Flags().String("src-ip", "", "Source IP CIDR")
	aclRuleCreateCmd.Flags().String("dst-ip", "", "Destination IP CIDR")
	aclRuleCreateCmd.Flags().Int("src-port-min", 0, "Source port min")
	aclRuleCreateCmd.Flags().Int("src-port-max", 0, "Source port max")
	aclRuleCreateCmd.Flags().Int("dst-port-min", 0, "Destination port min")
	aclRuleCreateCmd.Flags().Int("dst-port-max", 0, "Destination port max")
	aclRuleCreateCmd.Flags().String("policy", "allow", "Policy (allow/deny)")
	aclRuleCreateCmd.Flags().Int("order", 1, "Rule order/priority")
	aclRuleCreateCmd.MarkFlagRequired("acl-id")
	aclRuleCreateCmd.MarkFlagRequired("policy")

	aclRuleUpdateCmd.Flags().String("description", "", "Rule description")
	aclRuleUpdateCmd.Flags().String("protocol", "", "Protocol")
	aclRuleUpdateCmd.Flags().String("src-ip", "", "Source IP CIDR")
	aclRuleUpdateCmd.Flags().String("dst-ip", "", "Destination IP CIDR")
	aclRuleUpdateCmd.Flags().Int("src-port-min", 0, "Source port min")
	aclRuleUpdateCmd.Flags().Int("src-port-max", 0, "Source port max")
	aclRuleUpdateCmd.Flags().Int("dst-port-min", 0, "Destination port min")
	aclRuleUpdateCmd.Flags().Int("dst-port-max", 0, "Destination port max")
	aclRuleUpdateCmd.Flags().String("policy", "", "Policy (allow/deny)")
	aclRuleUpdateCmd.Flags().Int("order", 0, "Rule order/priority")

	// Binding commands
	aclBindingCmd.AddCommand(aclBindingListCmd)
	aclBindingCmd.AddCommand(aclBindingGetCmd)
	aclBindingCmd.AddCommand(aclBindingCreateCmd)
	aclBindingCmd.AddCommand(aclBindingDeleteCmd)

	aclBindingListCmd.Flags().String("acl-id", "", "Filter by ACL ID")

	aclBindingCreateCmd.Flags().String("acl-id", "", "ACL ID (required)")
	aclBindingCreateCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	aclBindingCreateCmd.MarkFlagRequired("acl-id")
	aclBindingCreateCmd.MarkFlagRequired("subnet-id")
}

func newNetworkACLClient() *networkacl.Client {
	return networkacl.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

// ============== ACL Commands ==============

var aclListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ACLs",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		result, err := client.ListACLs(context.Background())
		if err != nil {
			exitWithError("Failed to list ACLs", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tSHARED")
		for _, acl := range result.ACLs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n",
				acl.ID, acl.Name, acl.Description, acl.Shared)
		}
		w.Flush()
	},
}

var aclGetCmd = &cobra.Command{
	Use:   "get [acl-id]",
	Short: "Get ACL details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		result, err := client.GetACL(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get ACL", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		acl := result.ACL
		fmt.Printf("ID:          %s\n", acl.ID)
		fmt.Printf("Name:        %s\n", acl.Name)
		fmt.Printf("Description: %s\n", acl.Description)
		fmt.Printf("Shared:      %v\n", acl.Shared)
		fmt.Printf("Tenant ID:   %s\n", acl.TenantID)
		fmt.Printf("Created:     %s\n", acl.CreateTime)
		fmt.Printf("Updated:     %s\n", acl.UpdateTime)
	},
}

var aclCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new ACL",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &networkacl.CreateACLInput{
			Name:        name,
			Description: description,
		}

		result, err := client.CreateACL(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create ACL", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ACL created: %s\n", result.ACL.ID)
		fmt.Printf("Name: %s\n", result.ACL.Name)
	},
}

var aclUpdateCmd = &cobra.Command{
	Use:   "update [acl-id]",
	Short: "Update an ACL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &networkacl.UpdateACLInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateACL(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update ACL", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ACL updated: %s\n", result.ACL.ID)
		fmt.Printf("Name: %s\n", result.ACL.Name)
	},
}

var aclDeleteCmd = &cobra.Command{
	Use:   "delete [acl-id]",
	Short: "Delete an ACL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		if err := client.DeleteACL(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete ACL", err)
		}
		fmt.Printf("ACL %s deleted\n", args[0])
	},
}

// ============== ACL Rule Commands ==============

var aclRuleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ACL rules",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		aclID, _ := cmd.Flags().GetString("acl-id")

		var result *networkacl.ListACLRulesOutput
		var err error

		if aclID != "" {
			result, err = client.ListACLRulesByACL(context.Background(), aclID)
		} else {
			result, err = client.ListACLRules(context.Background())
		}

		if err != nil {
			exitWithError("Failed to list ACL rules", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var aclRuleGetCmd = &cobra.Command{
	Use:   "get [rule-id]",
	Short: "Get ACL rule details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		result, err := client.GetACLRule(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get ACL rule", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var aclRuleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new ACL rule",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
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

		result, err := client.CreateACLRule(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create ACL rule", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ACL rule created: %s\n", result.ACLRule.ID)
		fmt.Printf("Policy: %s\n", result.ACLRule.Policy)
	},
}

var aclRuleUpdateCmd = &cobra.Command{
	Use:   "update [rule-id]",
	Short: "Update an ACL rule",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
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

		result, err := client.UpdateACLRule(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update ACL rule", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ACL rule updated: %s\n", result.ACLRule.ID)
	},
}

var aclRuleDeleteCmd = &cobra.Command{
	Use:   "delete [rule-id]",
	Short: "Delete an ACL rule",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		if err := client.DeleteACLRule(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete ACL rule", err)
		}
		fmt.Printf("ACL rule %s deleted\n", args[0])
	},
}

// ============== ACL Binding Commands ==============

var aclBindingListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ACL bindings",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		aclID, _ := cmd.Flags().GetString("acl-id")

		var result *networkacl.ListACLBindingsOutput
		var err error

		if aclID != "" {
			result, err = client.ListACLBindingsByACL(context.Background(), aclID)
		} else {
			result, err = client.ListACLBindings(context.Background())
		}

		if err != nil {
			exitWithError("Failed to list ACL bindings", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tACL_ID\tSUBNET_ID")
		for _, binding := range result.ACLBindings {
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				binding.ID, binding.ACLID, binding.SubnetID)
		}
		w.Flush()
	},
}

var aclBindingGetCmd = &cobra.Command{
	Use:   "get [binding-id]",
	Short: "Get ACL binding details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		result, err := client.GetACLBinding(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get ACL binding", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		binding := result.ACLBinding
		fmt.Printf("ID:        %s\n", binding.ID)
		fmt.Printf("ACL ID:    %s\n", binding.ACLID)
		fmt.Printf("Subnet ID: %s\n", binding.SubnetID)
		fmt.Printf("Tenant ID: %s\n", binding.TenantID)
		fmt.Printf("Created:   %s\n", binding.CreateTime)
	},
}

var aclBindingCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new ACL binding",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		aclID, _ := cmd.Flags().GetString("acl-id")
		subnetID, _ := cmd.Flags().GetString("subnet-id")

		input := &networkacl.CreateACLBindingInput{
			ACLID:    aclID,
			SubnetID: subnetID,
		}

		result, err := client.CreateACLBinding(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create ACL binding", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ACL binding created: %s\n", result.ACLBinding.ID)
		fmt.Printf("ACL ID:    %s\n", result.ACLBinding.ACLID)
		fmt.Printf("Subnet ID: %s\n", result.ACLBinding.SubnetID)
	},
}

var aclBindingDeleteCmd = &cobra.Command{
	Use:   "delete [binding-id]",
	Short: "Delete an ACL binding",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		if err := client.DeleteACLBinding(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete ACL binding", err)
		}
		fmt.Printf("ACL binding %s deleted\n", args[0])
	},
}
