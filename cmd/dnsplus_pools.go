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
	dnsplusCmd.AddCommand(dnsDescribePoolsCmd)
	dnsplusCmd.AddCommand(dnsCreatePoolCmd)
	dnsplusCmd.AddCommand(dnsDeletePoolCmd)

	dnsDescribePoolsCmd.Flags().String("gslb-id", "", "GSLB ID (required)")
	dnsDescribePoolsCmd.MarkFlagRequired("gslb-id")

	dnsCreatePoolCmd.Flags().String("gslb-id", "", "GSLB ID (required)")
	dnsCreatePoolCmd.Flags().String("name", "", "Pool name")
	dnsCreatePoolCmd.Flags().String("description", "", "Pool description")
	dnsCreatePoolCmd.Flags().Int("priority", 1, "Pool priority")
	dnsCreatePoolCmd.Flags().Int("weight", 1, "Pool weight")
	dnsCreatePoolCmd.Flags().String("pool-region", "", "Pool region")
	dnsCreatePoolCmd.MarkFlagRequired("gslb-id")
	dnsCreatePoolCmd.MarkFlagRequired("name")

	dnsDeletePoolCmd.Flags().String("gslb-id", "", "GSLB ID (required)")
	dnsDeletePoolCmd.Flags().StringSlice("pool-ids", nil, "Pool IDs (required)")
	dnsDeletePoolCmd.MarkFlagRequired("gslb-id")
	dnsDeletePoolCmd.MarkFlagRequired("pool-ids")
}

var dnsDescribePoolsCmd = &cobra.Command{
	Use:   "describe-pools",
	Short: "List pools in a GSLB",
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID, _ := cmd.Flags().GetString("gslb-id")
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListPools(ctx, gslbID)
		if err != nil {
			return fmt.Errorf("failed to list pools: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.PoolList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tPRIORITY\tWEIGHT\tENDPOINTS")
		for _, pool := range result.PoolList {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%d\n",
				pool.PoolID,
				pool.PoolName,
				pool.PoolStatus,
				pool.Priority,
				pool.Weight,
				pool.EndpointCount,
			)
		}
		return w.Flush()
	},
}

var dnsCreatePoolCmd = &cobra.Command{
	Use:   "create-pool",
	Short: "Create a pool",
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID, _ := cmd.Flags().GetString("gslb-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		priority, _ := cmd.Flags().GetInt("priority")
		weight, _ := cmd.Flags().GetInt("weight")
		region, _ := cmd.Flags().GetString("pool-region")

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.CreatePoolInput{
			PoolName:    name,
			Description: description,
			Priority:    priority,
			Weight:      weight,
			Region:      region,
		}

		result, err := client.CreatePool(ctx, gslbID, input)
		if err != nil {
			return fmt.Errorf("failed to create pool: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Pool)
		}

		fmt.Printf("Pool created successfully: %s (%s)\n", result.Pool.PoolName, result.Pool.PoolID)
		return nil
	},
}

var dnsDeletePoolCmd = &cobra.Command{
	Use:   "delete-pools",
	Short: "Delete pools",
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID, _ := cmd.Flags().GetString("gslb-id")
		poolIDs, _ := cmd.Flags().GetStringSlice("pool-ids")

		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeletePools(ctx, gslbID, poolIDs)
		if err != nil {
			return fmt.Errorf("failed to delete pools: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Pool(s) deleted successfully\n")
		return nil
	},
}
