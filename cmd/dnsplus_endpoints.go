package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/dnsplus"
	"github.com/spf13/cobra"
)

func init() {
	dnsplusCmd.AddCommand(dnsDescribeEndpointsCmd)
	dnsplusCmd.AddCommand(dnsCreateEndpointCmd)
	dnsplusCmd.AddCommand(dnsDeleteEndpointCmd)

	dnsDescribeEndpointsCmd.Flags().String("gslb-id", "", "GSLB ID (required)")
	dnsDescribeEndpointsCmd.Flags().String("pool-id", "", "Pool ID (required)")
	dnsDescribeEndpointsCmd.MarkFlagRequired("gslb-id")
	dnsDescribeEndpointsCmd.MarkFlagRequired("pool-id")

	dnsCreateEndpointCmd.Flags().String("gslb-id", "", "GSLB ID (required)")
	dnsCreateEndpointCmd.Flags().String("pool-id", "", "Pool ID (required)")
	dnsCreateEndpointCmd.Flags().String("address", "", "Endpoint address (IP or domain)")
	dnsCreateEndpointCmd.Flags().Int("weight", 1, "Endpoint weight")
	dnsCreateEndpointCmd.Flags().String("description", "", "Endpoint description")
	dnsCreateEndpointCmd.MarkFlagRequired("gslb-id")
	dnsCreateEndpointCmd.MarkFlagRequired("pool-id")
	dnsCreateEndpointCmd.MarkFlagRequired("address")

	dnsDeleteEndpointCmd.Flags().String("gslb-id", "", "GSLB ID (required)")
	dnsDeleteEndpointCmd.Flags().String("pool-id", "", "Pool ID (required)")
	dnsDeleteEndpointCmd.Flags().StringSlice("endpoint-ids", nil, "Endpoint IDs (required)")
	dnsDeleteEndpointCmd.MarkFlagRequired("gslb-id")
	dnsDeleteEndpointCmd.MarkFlagRequired("pool-id")
	dnsDeleteEndpointCmd.MarkFlagRequired("endpoint-ids")
}

var dnsDescribeEndpointsCmd = &cobra.Command{
	Use:   "describe-endpoints",
	Short: "List endpoints in a pool",
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID, _ := cmd.Flags().GetString("gslb-id")
		poolID, _ := cmd.Flags().GetString("pool-id")

		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListEndpoints(ctx, gslbID, poolID)
		if err != nil {
			return fmt.Errorf("failed to list endpoints: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.EndpointList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tADDRESS\tSTATUS\tWEIGHT\tHEALTH")
		for _, ep := range result.EndpointList {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				ep.EndpointID,
				ep.EndpointAddress,
				ep.EndpointStatus,
				ep.Weight,
				ep.HealthStatus,
			)
		}
		return w.Flush()
	},
}

var dnsCreateEndpointCmd = &cobra.Command{
	Use:   "create-endpoint",
	Short: "Create an endpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID, _ := cmd.Flags().GetString("gslb-id")
		poolID, _ := cmd.Flags().GetString("pool-id")
		address, _ := cmd.Flags().GetString("address")
		weight, _ := cmd.Flags().GetInt("weight")
		description, _ := cmd.Flags().GetString("description")

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.CreateEndpointInput{
			EndpointAddress: address,
			Weight:          weight,
			Description:     description,
		}

		result, err := client.CreateEndpoint(ctx, gslbID, poolID, input)
		if err != nil {
			return fmt.Errorf("failed to create endpoint: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Endpoint)
		}

		fmt.Printf("Endpoint created successfully: %s (%s)\n", result.Endpoint.EndpointAddress, result.Endpoint.EndpointID)
		return nil
	},
}

var dnsDeleteEndpointCmd = &cobra.Command{
	Use:   "delete-endpoints",
	Short: "Delete endpoints",
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID, _ := cmd.Flags().GetString("gslb-id")
		poolID, _ := cmd.Flags().GetString("pool-id")
		endpointIDs, _ := cmd.Flags().GetStringSlice("endpoint-ids")

		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeleteEndpoints(ctx, gslbID, poolID, endpointIDs)
		if err != nil {
			return fmt.Errorf("failed to delete endpoints: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Endpoint(s) deleted successfully\n")
		return nil
	},
}
