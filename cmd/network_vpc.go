package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/vpc"
	"github.com/spf13/cobra"
)

func init() {
	// VPC Commands
	networkCmd.AddCommand(networkDescribeVPCsCmd)
	networkCmd.AddCommand(networkCreateVPCCmd)
	networkCmd.AddCommand(networkDeleteVPCCmd)

	// Subnet Commands
	networkCmd.AddCommand(networkDescribeSubnetsCmd)
	networkCmd.AddCommand(networkCreateSubnetCmd)
	networkCmd.AddCommand(networkDeleteSubnetCmd)

	networkCreateVPCCmd.Flags().String("name", "", "VPC name (required)")
	networkCreateVPCCmd.Flags().String("cidr", "", "VPC CIDR (required, e.g. 10.0.0.0/16)")
	networkCreateVPCCmd.MarkFlagRequired("name")
	networkCreateVPCCmd.MarkFlagRequired("cidr")

	networkDescribeVPCsCmd.Flags().String("vpc-id", "", "VPC ID")

	networkDeleteVPCCmd.Flags().String("vpc-id", "", "VPC ID (required)")
	networkDeleteVPCCmd.MarkFlagRequired("vpc-id")

	networkCreateSubnetCmd.Flags().String("name", "", "Subnet name (required)")
	networkCreateSubnetCmd.Flags().String("vpc-id", "", "VPC ID (required)")
	networkCreateSubnetCmd.Flags().String("cidr", "", "Subnet CIDR (required)")
	networkCreateSubnetCmd.Flags().String("gateway", "", "Gateway IP")
	networkCreateSubnetCmd.Flags().Bool("enable-dhcp", true, "Enable DHCP")
	networkCreateSubnetCmd.MarkFlagRequired("name")
	networkCreateSubnetCmd.MarkFlagRequired("vpc-id")
	networkCreateSubnetCmd.MarkFlagRequired("cidr")

	networkDescribeSubnetsCmd.Flags().String("subnet-id", "", "Subnet ID")

	networkDeleteSubnetCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	networkDeleteSubnetCmd.MarkFlagRequired("subnet-id")
}

var networkDescribeVPCsCmd = &cobra.Command{
	Use:   "describe-vpcs",
	Short: "Describe VPCs",
	Run: func(cmd *cobra.Command, args []string) {
		client := vpc.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()
		vpcID, _ := cmd.Flags().GetString("vpc-id")

		if vpcID != "" {
			result, err := client.GetVPC(ctx, vpcID)
			if err != nil {
				exitWithError("Failed to get VPC", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			v := result.VPC
			fmt.Printf("ID:     %s\n", v.ID)
			fmt.Printf("Name:   %s\n", v.Name)
			fmt.Printf("CIDR:   %s\n", v.CIDRv4)
			fmt.Printf("Status: %s\n", v.State)
		} else {
			result, err := client.ListVPCs(ctx)
			if err != nil {
				exitWithError("Failed to list VPCs", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tCIDR\tSTATUS")
			for _, v := range result.VPCs {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", v.ID, v.Name, v.CIDRv4, v.State)
			}
			w.Flush()
		}
	},
}

var networkCreateVPCCmd = &cobra.Command{
	Use:   "create-vpc",
	Short: "Create a new VPC",
	Run: func(cmd *cobra.Command, args []string) {
		client := vpc.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		cidr, _ := cmd.Flags().GetString("cidr")

		input := &vpc.CreateVPCInput{
			Name:   name,
			CIDRv4: cidr,
		}

		result, err := client.CreateVPC(ctx, input)
		if err != nil {
			exitWithError("Failed to create VPC", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("VPC created: %s\n", result.VPC.ID)
	},
}

var networkDeleteVPCCmd = &cobra.Command{
	Use:   "delete-vpc",
	Short: "Delete a VPC",
	Run: func(cmd *cobra.Command, args []string) {
		client := vpc.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()
		vpcID, _ := cmd.Flags().GetString("vpc-id")

		if err := client.DeleteVPC(ctx, vpcID); err != nil {
			exitWithError("Failed to delete VPC", err)
		}

		fmt.Printf("VPC %s deleted\n", vpcID)
	},
}

var networkDescribeSubnetsCmd = &cobra.Command{
	Use:   "describe-subnets",
	Short: "Describe subnets",
	Run: func(cmd *cobra.Command, args []string) {
		client := vpc.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()
		subnetID, _ := cmd.Flags().GetString("subnet-id")

		if subnetID != "" {
			result, err := client.GetSubnet(ctx, subnetID)
			if err != nil {
				exitWithError("Failed to get subnet", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			s := result.Subnet
			fmt.Printf("ID:         %s\n", s.ID)
			fmt.Printf("Name:       %s\n", s.Name)
			fmt.Printf("CIDR:       %s\n", s.CIDR)
			fmt.Printf("Network ID: %s\n", s.NetworkID)
			fmt.Printf("Gateway:    %s\n", s.GatewayIP)
			fmt.Printf("DHCP:       %v\n", s.EnableDHCP)
		} else {
			result, err := client.ListSubnets(ctx)
			if err != nil {
				exitWithError("Failed to list subnets", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tCIDR\tNETWORK_ID")
			for _, s := range result.Subnets {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.ID, s.Name, s.CIDR, s.NetworkID)
			}
			w.Flush()
		}
	},
}

var networkCreateSubnetCmd = &cobra.Command{
	Use:   "create-subnet",
	Short: "Create a new subnet",
	Run: func(cmd *cobra.Command, args []string) {
		client := vpc.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		vpcID, _ := cmd.Flags().GetString("vpc-id")
		cidr, _ := cmd.Flags().GetString("cidr")
		gateway, _ := cmd.Flags().GetString("gateway")

		input := &vpc.CreateSubnetInput{
			Name:      name,
			VPCID:     vpcID,
			CIDR:      cidr,
			GatewayIP: gateway,
		}

		result, err := client.CreateSubnet(ctx, input)
		if err != nil {
			exitWithError("Failed to create subnet", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Subnet created: %s\n", result.VPCSubnet.ID)
	},
}

var networkDeleteSubnetCmd = &cobra.Command{
	Use:   "delete-subnet",
	Short: "Delete a subnet",
	Run: func(cmd *cobra.Command, args []string) {
		client := vpc.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()
		subnetID, _ := cmd.Flags().GetString("subnet-id")

		if err := client.DeleteSubnet(ctx, subnetID); err != nil {
			exitWithError("Failed to delete subnet", err)
		}

		fmt.Printf("Subnet %s deleted\n", subnetID)
	},
}
