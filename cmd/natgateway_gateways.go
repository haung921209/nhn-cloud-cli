package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/natgateway"
	"github.com/spf13/cobra"
)

func init() {
	natGatewayCmd.AddCommand(natGatewayDescribeCmd)
	natGatewayCmd.AddCommand(natGatewayGetCmd)
	natGatewayCmd.AddCommand(natGatewayCreateCmd)
	natGatewayCmd.AddCommand(natGatewayUpdateCmd)
	natGatewayCmd.AddCommand(natGatewayDeleteCmd)

	natGatewayGetCmd.Flags().String("gateway-id", "", "NAT Gateway ID (required)")
	natGatewayGetCmd.MarkFlagRequired("gateway-id")

	natGatewayCreateCmd.Flags().String("name", "", "NAT gateway name (required)")
	natGatewayCreateCmd.Flags().String("description", "", "Description")
	natGatewayCreateCmd.Flags().String("vpc-id", "", "VPC ID (required)")
	natGatewayCreateCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	natGatewayCreateCmd.Flags().String("floating-ip-id", "", "Floating IP ID (optional)")
	natGatewayCreateCmd.MarkFlagRequired("name")
	natGatewayCreateCmd.MarkFlagRequired("vpc-id")
	natGatewayCreateCmd.MarkFlagRequired("subnet-id")

	natGatewayUpdateCmd.Flags().String("gateway-id", "", "NAT Gateway ID (required)")
	natGatewayUpdateCmd.Flags().String("name", "", "NAT gateway name")
	natGatewayUpdateCmd.Flags().String("description", "", "Description")
	natGatewayUpdateCmd.MarkFlagRequired("gateway-id")

	natGatewayDeleteCmd.Flags().String("gateway-id", "", "NAT Gateway ID (required)")
	natGatewayDeleteCmd.MarkFlagRequired("gateway-id")
}

var natGatewayDescribeCmd = &cobra.Command{
	Use:     "describe-nat-gateways",
	Aliases: []string{"list-nat-gateways"},
	Short:   "List all NAT gateways",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNATGatewayClient()
		result, err := client.ListNATGateways(context.Background())
		if err != nil {
			exitWithError("Failed to list NAT gateways", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tFLOATING_IP\tSTATUS\tSTATE")
		for _, gw := range result.NATGateways {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				gw.ID, gw.Name, gw.FloatingIPAddress, gw.Status, gw.State)
		}
		w.Flush()
	},
}

var natGatewayGetCmd = &cobra.Command{
	Use:     "describe-nat-gateway",
	Aliases: []string{"get-nat-gateway"},
	Short:   "Get NAT gateway details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNATGatewayClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("gateway-id")

		result, err := client.GetNATGateway(ctx, id)
		if err != nil {
			exitWithError("Failed to get NAT gateway", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		gw := result.NATGateway
		fmt.Printf("ID:              %s\n", gw.ID)
		fmt.Printf("Name:            %s\n", gw.Name)
		fmt.Printf("Description:     %s\n", gw.Description)
		fmt.Printf("VPC ID:          %s\n", gw.VPCID)
		fmt.Printf("Subnet ID:       %s\n", gw.SubnetID)
		fmt.Printf("Floating IP ID:  %s\n", gw.FloatingIPID)
		fmt.Printf("Floating IP:     %s\n", gw.FloatingIPAddress)
		fmt.Printf("Status:          %s\n", gw.Status)
		fmt.Printf("State:           %s\n", gw.State)
		fmt.Printf("Tenant ID:       %s\n", gw.TenantID)
		fmt.Printf("Created At:      %s\n", gw.CreatedAt)
		fmt.Printf("Updated At:      %s\n", gw.UpdatedAt)
	},
}

var natGatewayCreateCmd = &cobra.Command{
	Use:   "create-nat-gateway",
	Short: "Create a new NAT gateway",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNATGatewayClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		vpcID, _ := cmd.Flags().GetString("vpc-id")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		floatingIPID, _ := cmd.Flags().GetString("floating-ip-id")

		input := &natgateway.CreateNATGatewayInput{
			Name:         name,
			Description:  description,
			VPCID:        vpcID,
			SubnetID:     subnetID,
			FloatingIPID: floatingIPID,
		}

		result, err := client.CreateNATGateway(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create NAT gateway", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("NAT gateway created: %s\n", result.NATGateway.ID)
		fmt.Printf("Name: %s\n", result.NATGateway.Name)
		fmt.Printf("Status: %s\n", result.NATGateway.Status)
	},
}

var natGatewayUpdateCmd = &cobra.Command{
	Use:   "update-nat-gateway",
	Short: "Update a NAT gateway",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNATGatewayClient()
		id, _ := cmd.Flags().GetString("gateway-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &natgateway.UpdateNATGatewayInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateNATGateway(context.Background(), id, input)
		if err != nil {
			exitWithError("Failed to update NAT gateway", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("NAT gateway updated: %s\n", result.NATGateway.ID)
		fmt.Printf("Name: %s\n", result.NATGateway.Name)
	},
}

var natGatewayDeleteCmd = &cobra.Command{
	Use:   "delete-nat-gateway",
	Short: "Delete a NAT gateway",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNATGatewayClient()
		id, _ := cmd.Flags().GetString("gateway-id")
		if err := client.DeleteNATGateway(context.Background(), id); err != nil {
			exitWithError("Failed to delete NAT gateway", err)
		}
		fmt.Printf("NAT gateway %s deleted\n", id)
	},
}
