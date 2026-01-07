package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/natgateway"
	"github.com/spf13/cobra"
)

var natGatewayCmd = &cobra.Command{
	Use:     "nat-gateway",
	Aliases: []string{"nat", "natgw"},
	Short:   "Manage NAT Gateways",
}

func init() {
	rootCmd.AddCommand(natGatewayCmd)

	natGatewayCmd.AddCommand(natGatewayListCmd)
	natGatewayCmd.AddCommand(natGatewayGetCmd)
	natGatewayCmd.AddCommand(natGatewayCreateCmd)
	natGatewayCmd.AddCommand(natGatewayUpdateCmd)
	natGatewayCmd.AddCommand(natGatewayDeleteCmd)

	natGatewayCreateCmd.Flags().String("name", "", "NAT gateway name (required)")
	natGatewayCreateCmd.Flags().String("description", "", "Description")
	natGatewayCreateCmd.Flags().String("vpc-id", "", "VPC ID (required)")
	natGatewayCreateCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	natGatewayCreateCmd.Flags().String("floating-ip-id", "", "Floating IP ID (optional)")
	natGatewayCreateCmd.MarkFlagRequired("name")
	natGatewayCreateCmd.MarkFlagRequired("vpc-id")
	natGatewayCreateCmd.MarkFlagRequired("subnet-id")

	natGatewayUpdateCmd.Flags().String("name", "", "NAT gateway name")
	natGatewayUpdateCmd.Flags().String("description", "", "Description")
}

func newNATGatewayClient() *natgateway.Client {
	return natgateway.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

var natGatewayListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all NAT gateways",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNATGatewayClient()
		result, err := client.ListNATGateways(context.Background())
		if err != nil {
			exitWithError("Failed to list NAT gateways", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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
	Use:   "get [nat-gateway-id]",
	Short: "Get NAT gateway details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNATGatewayClient()
		result, err := client.GetNATGateway(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get NAT gateway", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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
	Use:   "create",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("NAT gateway created: %s\n", result.NATGateway.ID)
		fmt.Printf("Name: %s\n", result.NATGateway.Name)
		fmt.Printf("Status: %s\n", result.NATGateway.Status)
	},
}

var natGatewayUpdateCmd = &cobra.Command{
	Use:   "update [nat-gateway-id]",
	Short: "Update a NAT gateway",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNATGatewayClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &natgateway.UpdateNATGatewayInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateNATGateway(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update NAT gateway", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("NAT gateway updated: %s\n", result.NATGateway.ID)
		fmt.Printf("Name: %s\n", result.NATGateway.Name)
	},
}

var natGatewayDeleteCmd = &cobra.Command{
	Use:   "delete [nat-gateway-id]",
	Short: "Delete a NAT gateway",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNATGatewayClient()
		if err := client.DeleteNATGateway(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete NAT gateway", err)
		}
		fmt.Printf("NAT gateway %s deleted\n", args[0])
	},
}
