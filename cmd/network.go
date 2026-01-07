package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/floatingip"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/securitygroup"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/vpc"
	"github.com/spf13/cobra"
)

var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "Manage VPCs and subnets",
}

var networkSecurityGroupCmd = &cobra.Command{
	Use:     "security-group",
	Aliases: []string{"sg"},
	Short:   "Manage network security groups",
}

var floatingIPCmd = &cobra.Command{
	Use:     "floating-ip",
	Aliases: []string{"fip"},
	Short:   "Manage floating IPs",
}

func init() {
	rootCmd.AddCommand(vpcCmd)
	rootCmd.AddCommand(networkSecurityGroupCmd)
	rootCmd.AddCommand(floatingIPCmd)

	vpcCmd.AddCommand(vpcListCmd)
	vpcCmd.AddCommand(vpcGetCmd)
	vpcCmd.AddCommand(vpcSubnetsCmd)

	networkSecurityGroupCmd.AddCommand(sgListCmd)
	networkSecurityGroupCmd.AddCommand(sgGetCmd)
	networkSecurityGroupCmd.AddCommand(sgCreateCmd)
	networkSecurityGroupCmd.AddCommand(sgDeleteCmd)
	networkSecurityGroupCmd.AddCommand(sgRuleCreateCmd)

	sgCreateCmd.Flags().String("name", "", "Security group name (required)")
	sgCreateCmd.Flags().String("description", "", "Security group description")
	sgCreateCmd.MarkFlagRequired("name")

	sgRuleCreateCmd.Flags().String("security-group-id", "", "Security group ID (required)")
	sgRuleCreateCmd.Flags().String("direction", "ingress", "Direction (ingress/egress)")
	sgRuleCreateCmd.Flags().String("ethertype", "IPv4", "Ethertype (IPv4/IPv6)")
	sgRuleCreateCmd.Flags().String("protocol", "", "Protocol (tcp/udp/icmp)")
	sgRuleCreateCmd.Flags().Int("port-min", 0, "Minimum port")
	sgRuleCreateCmd.Flags().Int("port-max", 0, "Maximum port")
	sgRuleCreateCmd.Flags().String("remote-ip", "0.0.0.0/0", "Remote IP prefix")
	sgRuleCreateCmd.MarkFlagRequired("security-group-id")

	floatingIPCmd.AddCommand(fipListCmd)
	floatingIPCmd.AddCommand(fipGetCmd)
	floatingIPCmd.AddCommand(fipCreateCmd)
	floatingIPCmd.AddCommand(fipDeleteCmd)
	floatingIPCmd.AddCommand(fipAssociateCmd)
	floatingIPCmd.AddCommand(fipDisassociateCmd)

	fipCreateCmd.Flags().String("network-id", "", "External network ID (required)")
	fipCreateCmd.MarkFlagRequired("network-id")

	fipAssociateCmd.Flags().String("port-id", "", "Port ID to associate (required)")
	fipAssociateCmd.MarkFlagRequired("port-id")
}

func getIdentityCreds() credentials.IdentityCredentials {
	return credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
}

var vpcListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPCs",
	Run: func(cmd *cobra.Command, args []string) {
		client := vpc.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		result, err := client.ListVPCs(ctx)
		if err != nil {
			exitWithError("Failed to list VPCs", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tCIDR\tSTATUS")
		for _, v := range result.VPCs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", v.ID, v.Name, v.CIDRv4, v.State)
		}
		w.Flush()
	},
}

var vpcGetCmd = &cobra.Command{
	Use:   "get [vpc-id]",
	Short: "Get VPC details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := vpc.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		result, err := client.GetVPC(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get VPC", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		v := result.VPC
		fmt.Printf("ID:     %s\n", v.ID)
		fmt.Printf("Name:   %s\n", v.Name)
		fmt.Printf("CIDR:   %s\n", v.CIDRv4)
		fmt.Printf("Status: %s\n", v.State)
	},
}

var vpcSubnetsCmd = &cobra.Command{
	Use:   "subnets",
	Short: "List subnets",
	Run: func(cmd *cobra.Command, args []string) {
		client := vpc.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		result, err := client.ListSubnets(ctx)
		if err != nil {
			exitWithError("Failed to list subnets", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tCIDR\tNETWORK_ID")
		for _, s := range result.Subnets {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.ID, s.Name, s.CIDR, s.NetworkID)
		}
		w.Flush()
	},
}

var sgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all security groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := securitygroup.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		result, err := client.ListSecurityGroups(ctx)
		if err != nil {
			exitWithError("Failed to list security groups", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION")
		for _, sg := range result.SecurityGroups {
			fmt.Fprintf(w, "%s\t%s\t%s\n", sg.ID, sg.Name, sg.Description)
		}
		w.Flush()
	},
}

var sgGetCmd = &cobra.Command{
	Use:   "get [security-group-id]",
	Short: "Get security group details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := securitygroup.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		result, err := client.GetSecurityGroup(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get security group", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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
			if rule.PortRangeMin != nil && *rule.PortRangeMin != 0 {
				port = fmt.Sprintf("%d-%d", *rule.PortRangeMin, *rule.PortRangeMax)
			}
			remote := rule.RemoteIPPrefix
			if remote == "" {
				remote = rule.RemoteGroupID
			}
			fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n", rule.Direction, rule.Protocol, port, remote)
		}
		w.Flush()
	},
}

var sgCreateCmd = &cobra.Command{
	Use:   "create",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Security group created: %s\n", result.SecurityGroup.ID)
		fmt.Printf("Name: %s\n", result.SecurityGroup.Name)
	},
}

var sgDeleteCmd = &cobra.Command{
	Use:   "delete [security-group-id]",
	Short: "Delete a security group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := securitygroup.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		if err := client.DeleteSecurityGroup(ctx, args[0]); err != nil {
			exitWithError("Failed to delete security group", err)
		}

		fmt.Printf("Security group %s deleted\n", args[0])
	},
}

var sgRuleCreateCmd = &cobra.Command{
	Use:   "rule-create",
	Short: "Create a security group rule",
	Run: func(cmd *cobra.Command, args []string) {
		client := securitygroup.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		sgID, _ := cmd.Flags().GetString("security-group-id")
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
			exitWithError("Failed to create security group rule", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Security group rule created: %s\n", result.SecurityGroupRule.ID)
	},
}

var fipListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all floating IPs",
	Run: func(cmd *cobra.Command, args []string) {
		client := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		result, err := client.ListFloatingIPs(ctx)
		if err != nil {
			exitWithError("Failed to list floating IPs", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tFLOATING_IP\tFIXED_IP\tSTATUS\tPORT_ID")
		for _, fip := range result.FloatingIPs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				fip.ID, fip.FloatingIPAddress, fip.FixedIPAddress, fip.Status, fip.PortID)
		}
		w.Flush()
	},
}

var fipGetCmd = &cobra.Command{
	Use:   "get [floating-ip-id]",
	Short: "Get floating IP details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		result, err := client.GetFloatingIP(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get floating IP", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fip := result.FloatingIP
		fmt.Printf("ID:          %s\n", fip.ID)
		fmt.Printf("Floating IP: %s\n", fip.FloatingIPAddress)
		fmt.Printf("Fixed IP:    %s\n", fip.FixedIPAddress)
		fmt.Printf("Status:      %s\n", fip.Status)
		fmt.Printf("Port ID:     %s\n", fip.PortID)
		fmt.Printf("Tenant ID:   %s\n", fip.TenantID)
	},
}

var fipCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new floating IP",
	Run: func(cmd *cobra.Command, args []string) {
		client := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		networkID, _ := cmd.Flags().GetString("network-id")

		input := &floatingip.CreateFloatingIPInput{
			FloatingNetworkID: networkID,
		}

		result, err := client.CreateFloatingIP(ctx, input)
		if err != nil {
			exitWithError("Failed to create floating IP", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Floating IP created: %s\n", result.FloatingIP.FloatingIPAddress)
		fmt.Printf("ID: %s\n", result.FloatingIP.ID)
	},
}

var fipDeleteCmd = &cobra.Command{
	Use:   "delete [floating-ip-id]",
	Short: "Delete a floating IP",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		if err := client.DeleteFloatingIP(ctx, args[0]); err != nil {
			exitWithError("Failed to delete floating IP", err)
		}

		fmt.Printf("Floating IP %s deleted\n", args[0])
	},
}

var fipAssociateCmd = &cobra.Command{
	Use:   "associate [floating-ip-id]",
	Short: "Associate floating IP with a port",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		portID, _ := cmd.Flags().GetString("port-id")

		input := &floatingip.UpdateFloatingIPInput{
			PortID: &portID,
		}

		result, err := client.UpdateFloatingIP(ctx, args[0], input)
		if err != nil {
			exitWithError("Failed to associate floating IP", err)
		}

		fmt.Printf("Floating IP %s associated with port %s\n",
			result.FloatingIP.FloatingIPAddress, portID)
	},
}

var fipDisassociateCmd = &cobra.Command{
	Use:   "disassociate [floating-ip-id]",
	Short: "Disassociate floating IP from port",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		emptyPort := ""
		input := &floatingip.UpdateFloatingIPInput{
			PortID: &emptyPort,
		}

		_, err := client.UpdateFloatingIP(ctx, args[0], input)
		if err != nil {
			exitWithError("Failed to disassociate floating IP", err)
		}

		fmt.Printf("Floating IP %s disassociated\n", args[0])
	},
}
