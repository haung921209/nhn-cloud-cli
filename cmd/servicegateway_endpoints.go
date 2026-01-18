package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	serviceGatewayCmd.AddCommand(sgDescribeEndpointsCmd)
	serviceGatewayCmd.AddCommand(sgGetEndpointCmd)

	sgGetEndpointCmd.Flags().String("endpoint-id", "", "Service Endpoint ID (required)")
	sgGetEndpointCmd.MarkFlagRequired("endpoint-id")
}

var sgDescribeEndpointsCmd = &cobra.Command{
	Use:     "describe-service-endpoints",
	Aliases: []string{"list-service-endpoints"},
	Short:   "List all available service endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		result, err := client.ListServiceEndpoints(context.Background())
		if err != nil {
			exitWithError("Failed to list service endpoints", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSERVICE_NAME\tREGION\tTYPE")
		for _, ep := range result.ServiceEndpoints {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				ep.ID, ep.Name, ep.ServiceName, ep.Region, ep.EndpointType)
		}
		w.Flush()
	},
}

var sgGetEndpointCmd = &cobra.Command{
	Use:     "describe-service-endpoint",
	Aliases: []string{"get-service-endpoint"},
	Short:   "Get service endpoint details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("endpoint-id")

		result, err := client.GetServiceEndpoint(ctx, id)
		if err != nil {
			exitWithError("Failed to get service endpoint", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		ep := result.ServiceEndpoint
		fmt.Printf("ID:            %s\n", ep.ID)
		fmt.Printf("Name:          %s\n", ep.Name)
		fmt.Printf("Service Name:  %s\n", ep.ServiceName)
		fmt.Printf("Description:   %s\n", ep.Description)
		fmt.Printf("Region:        %s\n", ep.Region)
		fmt.Printf("Endpoint Type: %s\n", ep.EndpointType)
		fmt.Printf("Created:       %s\n", ep.CreateTime)
	},
}
