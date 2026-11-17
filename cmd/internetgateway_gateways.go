package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/internetgateway"
	"github.com/spf13/cobra"
)

func init() {
	internetGatewayCmd.AddCommand(igwDescribeCmd)
	internetGatewayCmd.AddCommand(igwGetCmd)
	internetGatewayCmd.AddCommand(igwCreateCmd)
	internetGatewayCmd.AddCommand(igwDeleteCmd)

	igwGetCmd.Flags().String("gateway-id", "", requiredHint("Internet Gateway ID", hintIGWID))
	igwGetCmd.MarkFlagRequired("gateway-id")

	igwCreateCmd.Flags().String("name", "", "Internet gateway name (required)")
	igwCreateCmd.Flags().String("routing-table-id", "", "Routing table to attach the new gateway to ("+hintRoutingTableID+")")
	igwCreateCmd.Flags().String("external-network-id", "", "External network ID ("+hintExtNetID+")")
	igwCreateCmd.MarkFlagRequired("name")

	igwDeleteCmd.Flags().String("gateway-id", "", requiredHint("Internet Gateway ID", hintIGWID))
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
		name, _ := cmd.Flags().GetString("name")
		routingTableID, _ := cmd.Flags().GetString("routing-table-id")
		externalNetworkID, _ := cmd.Flags().GetString("external-network-id")

		// The create API takes only a name and an external network. A
		// routing_table_id in the create body is ignored (the gateway stays
		// unattached), so attaching is a separate call.
		var res internetgateway.GetInternetGatewayOutput
		gw := map[string]string{"name": name}
		if externalNetworkID != "" {
			gw["external_network_id"] = externalNetworkID
		}
		raw := apiCall(http.MethodPost, "/v2.0/internetgateways", map[string]interface{}{"internetgateway": gw}, &res)
		gwID := res.InternetGateway.ID
		attachHint := fmt.Sprintf("nhncloud network attach-routing-table-gateway --routing-table-id <id> --gateway-id %s", gwID)

		if routingTableID != "" {
			body := map[string]string{"gateway_id": gwID}
			if err := networkRequest(context.Background(), http.MethodPut,
				"/v2.0/routingtables/"+url.PathEscape(routingTableID)+"/attach_gateway", body, nil); err != nil {
				exitWithError(fmt.Sprintf("Internet gateway %s was created but attaching it to routing table %s failed; retry with: %s",
					gwID, routingTableID, strings.Replace(attachHint, "<id>", routingTableID, 1)), err)
			}
		}

		emit(raw, func() {
			fmt.Printf("Internet gateway created: %s\n", gwID)
			fmt.Printf("Name: %s\n", res.InternetGateway.Name)
			fmt.Printf("State: %s\n", res.InternetGateway.State)
			if routingTableID != "" {
				fmt.Printf("Attached to routing table: %s\n", routingTableID)
			} else {
				fmt.Printf("Attach it with: %s\n", attachHint)
			}
		})
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
