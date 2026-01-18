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
	dnsplusCmd.AddCommand(dnsDescribeGSLBCmd)
	dnsplusCmd.AddCommand(dnsCreateGSLBCmd)
	dnsplusCmd.AddCommand(dnsDeleteGSLBCmd)

	dnsCreateGSLBCmd.Flags().String("name", "", "GSLB name")
	dnsCreateGSLBCmd.Flags().String("description", "", "GSLB description")
	dnsCreateGSLBCmd.Flags().String("routing-type", "", "Routing type (FAILOVER, RANDOM, GEOLOCATION)")
	dnsCreateGSLBCmd.Flags().Int("ttl", 300, "TTL in seconds")
	dnsCreateGSLBCmd.Flags().String("health-check-id", "", "Health check ID")
	dnsCreateGSLBCmd.MarkFlagRequired("name")
	dnsCreateGSLBCmd.MarkFlagRequired("routing-type")

	dnsDeleteGSLBCmd.Flags().StringSlice("gslb-ids", nil, "GSLB IDs (required)")
	dnsDeleteGSLBCmd.MarkFlagRequired("gslb-ids")
}

var dnsDescribeGSLBCmd = &cobra.Command{
	Use:     "describe-gslbs",
	Aliases: []string{"list-gslbs"},
	Short:   "List GSLBs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListGSLBs(ctx)
		if err != nil {
			return fmt.Errorf("failed to list GSLBs: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.GslbList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tDOMAIN\tSTATUS\tROUTING\tPOOLS")
		for _, gslb := range result.GslbList {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\n",
				gslb.GslbID,
				gslb.GslbName,
				gslb.GslbDomain,
				gslb.GslbStatus,
				gslb.RoutingType,
				gslb.PoolCount,
			)
		}
		return w.Flush()
	},
}

var dnsCreateGSLBCmd = &cobra.Command{
	Use:   "create-gslb",
	Short: "Create a GSLB",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		routingType, _ := cmd.Flags().GetString("routing-type")
		ttl, _ := cmd.Flags().GetInt("ttl")
		healthCheckID, _ := cmd.Flags().GetString("health-check-id")

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.CreateGSLBInput{
			GslbName:      name,
			Description:   description,
			RoutingType:   routingType,
			TTL:           ttl,
			HealthCheckID: healthCheckID,
		}

		result, err := client.CreateGSLB(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to create GSLB: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Gslb)
		}

		fmt.Printf("GSLB created successfully: %s (%s)\n", result.Gslb.GslbName, result.Gslb.GslbID)
		return nil
	},
}

var dnsDeleteGSLBCmd = &cobra.Command{
	Use:   "delete-gslbs",
	Short: "Delete GSLBs",
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbIDs, _ := cmd.Flags().GetStringSlice("gslb-ids")
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeleteGSLBs(ctx, gslbIDs)
		if err != nil {
			return fmt.Errorf("failed to delete GSLBs: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("GSLB(s) deleted successfully\n")
		return nil
	},
}
