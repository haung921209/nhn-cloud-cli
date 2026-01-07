package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/internetgateway"
	"github.com/spf13/cobra"
)

var internetGatewayCmd = &cobra.Command{
	Use:     "internet-gateway",
	Aliases: []string{"igw"},
	Short:   "Manage Internet Gateways",
}

func init() {
	rootCmd.AddCommand(internetGatewayCmd)

	internetGatewayCmd.AddCommand(igwListCmd)
	internetGatewayCmd.AddCommand(igwGetCmd)
	internetGatewayCmd.AddCommand(igwCreateCmd)
	internetGatewayCmd.AddCommand(igwDeleteCmd)
	internetGatewayCmd.AddCommand(igwExternalNetworksCmd)

	igwCreateCmd.Flags().String("name", "", "Internet gateway name (required)")
	igwCreateCmd.Flags().String("routing-table-id", "", "Routing table ID (required)")
	igwCreateCmd.Flags().String("external-network-id", "", "External network ID (optional)")
	igwCreateCmd.MarkFlagRequired("name")
	igwCreateCmd.MarkFlagRequired("routing-table-id")
}

func newInternetGatewayClient() *internetgateway.Client {
	return internetgateway.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

var igwListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all internet gateways",
	Run: func(cmd *cobra.Command, args []string) {
		client := newInternetGatewayClient()
		result, err := client.ListInternetGateways(context.Background())
		if err != nil {
			exitWithError("Failed to list internet gateways", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tROUTING_TABLE_ID\tSTATE")
		for _, gw := range result.InternetGateways {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				gw.ID, gw.Name, gw.RoutingTableID, gw.State)
		}
		w.Flush()
	},
}

var igwGetCmd = &cobra.Command{
	Use:   "get [internet-gateway-id]",
	Short: "Get internet gateway details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newInternetGatewayClient()
		result, err := client.GetInternetGateway(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get internet gateway", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		gw := result.InternetGateway
		fmt.Printf("ID:                  %s\n", gw.ID)
		fmt.Printf("Name:                %s\n", gw.Name)
		fmt.Printf("Routing Table ID:    %s\n", gw.RoutingTableID)
		fmt.Printf("External Network ID: %s\n", gw.ExternalNetworkID)
		fmt.Printf("State:               %s\n", gw.State)
		fmt.Printf("Tenant ID:           %s\n", gw.TenantID)
		fmt.Printf("Created At:          %s\n", gw.CreateTime)
	},
}

var igwCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new internet gateway",
	Run: func(cmd *cobra.Command, args []string) {
		client := newInternetGatewayClient()
		name, _ := cmd.Flags().GetString("name")
		routingTableID, _ := cmd.Flags().GetString("routing-table-id")
		externalNetworkID, _ := cmd.Flags().GetString("external-network-id")

		input := &internetgateway.CreateInternetGatewayInput{
			Name:              name,
			RoutingTableID:    routingTableID,
			ExternalNetworkID: externalNetworkID,
		}

		result, err := client.CreateInternetGateway(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create internet gateway", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Internet gateway created: %s\n", result.InternetGateway.ID)
		fmt.Printf("Name: %s\n", result.InternetGateway.Name)
		fmt.Printf("State: %s\n", result.InternetGateway.State)
	},
}

var igwDeleteCmd = &cobra.Command{
	Use:   "delete [internet-gateway-id]",
	Short: "Delete an internet gateway",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newInternetGatewayClient()
		if err := client.DeleteInternetGateway(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete internet gateway", err)
		}
		fmt.Printf("Internet gateway %s deleted\n", args[0])
	},
}

var igwExternalNetworksCmd = &cobra.Command{
	Use:   "external-networks",
	Short: "List available external networks",
	Run: func(cmd *cobra.Command, args []string) {
		client := newInternetGatewayClient()
		result, err := client.ListExternalNetworks(context.Background())
		if err != nil {
			exitWithError("Failed to list external networks", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tROUTER_EXTERNAL")
		for _, net := range result.Networks {
			fmt.Fprintf(w, "%s\t%s\t%v\n",
				net.ID, net.Name, net.RouterExternal)
		}
		w.Flush()
	},
}
