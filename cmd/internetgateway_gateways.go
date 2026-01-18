package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/internetgateway"
	"github.com/spf13/cobra"
)

func init() {
	internetGatewayCmd.AddCommand(igwDescribeCmd)
	internetGatewayCmd.AddCommand(igwGetCmd)
	internetGatewayCmd.AddCommand(igwCreateCmd)
	internetGatewayCmd.AddCommand(igwDeleteCmd)

	igwGetCmd.Flags().String("gateway-id", "", "Internet Gateway ID (required)")
	igwGetCmd.MarkFlagRequired("gateway-id")

	igwCreateCmd.Flags().String("name", "", "Internet gateway name (required)")
	igwCreateCmd.Flags().String("routing-table-id", "", "Routing table ID (required)")
	igwCreateCmd.Flags().String("external-network-id", "", "External network ID (optional)")
	igwCreateCmd.MarkFlagRequired("name")
	igwCreateCmd.MarkFlagRequired("routing-table-id")

	igwDeleteCmd.Flags().String("gateway-id", "", "Internet Gateway ID (required)")
	igwDeleteCmd.MarkFlagRequired("gateway-id")
}

var igwDescribeCmd = &cobra.Command{
	Use:     "describe-internet-gateways",
	Aliases: []string{"list-internet-gateways", "list"},
	Short:   "List all internet gateways",
	Run: func(cmd *cobra.Command, args []string) {
		client := newInternetGatewayClient()
		result, err := client.ListInternetGateways(context.Background())
		if err != nil {
			exitWithError("Failed to list internet gateways", err)
		}

		if output == "json" {
			printJSON(result)
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
	Use:     "describe-internet-gateway",
	Aliases: []string{"get-internet-gateway", "get"},
	Short:   "Get internet gateway details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newInternetGatewayClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("gateway-id")

		result, err := client.GetInternetGateway(ctx, id)
		if err != nil {
			exitWithError("Failed to get internet gateway", err)
		}

		if output == "json" {
			printJSON(result)
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
	Use:   "create-internet-gateway",
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
			printJSON(result)
			return
		}

		fmt.Printf("Internet gateway created: %s\n", result.InternetGateway.ID)
		fmt.Printf("Name: %s\n", result.InternetGateway.Name)
		fmt.Printf("State: %s\n", result.InternetGateway.State)
	},
}

var igwDeleteCmd = &cobra.Command{
	Use:   "delete-internet-gateway",
	Short: "Delete an internet gateway",
	Run: func(cmd *cobra.Command, args []string) {
		client := newInternetGatewayClient()
		id, _ := cmd.Flags().GetString("gateway-id")
		if err := client.DeleteInternetGateway(context.Background(), id); err != nil {
			exitWithError("Failed to delete internet gateway", err)
		}
		fmt.Printf("Internet gateway %s deleted\n", id)
	},
}
