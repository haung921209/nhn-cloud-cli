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
	dnsplusCmd.AddCommand(dnsDescribeZonesCmd)
	dnsplusCmd.AddCommand(dnsCreateZoneCmd)
	dnsplusCmd.AddCommand(dnsUpdateZoneCmd)
	dnsplusCmd.AddCommand(dnsDeleteZoneCmd)

	dnsCreateZoneCmd.Flags().String("name", "", "Zone name (e.g., example.com)")
	dnsCreateZoneCmd.Flags().String("description", "", "Zone description")
	dnsCreateZoneCmd.MarkFlagRequired("name")

	dnsUpdateZoneCmd.Flags().String("zone-id", "", "Zone ID (required)")
	dnsUpdateZoneCmd.Flags().String("description", "", "Zone description")
	dnsUpdateZoneCmd.Flags().String("status", "", "Zone status (USE, STOP)")
	dnsUpdateZoneCmd.MarkFlagRequired("zone-id")

	dnsDeleteZoneCmd.Flags().StringSlice("zone-ids", nil, "Zone IDs (required)")
	dnsDeleteZoneCmd.MarkFlagRequired("zone-ids")
}

var dnsDescribeZonesCmd = &cobra.Command{
	Use:     "describe-zones",
	Aliases: []string{"list-zones"},
	Short:   "List DNS zones",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListZones(ctx)
		if err != nil {
			return fmt.Errorf("failed to list zones: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.ZoneList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tRECORDS\tCREATED")
		for _, zone := range result.ZoneList {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				zone.ZoneID,
				zone.ZoneName,
				zone.ZoneStatus,
				zone.RecordSetCount,
				zone.CreatedAt.Format("2006-01-02"),
			)
		}
		return w.Flush()
	},
}

var dnsCreateZoneCmd = &cobra.Command{
	Use:   "create-zone",
	Short: "Create a DNS zone",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.CreateZoneInput{
			ZoneName:    name,
			Description: description,
		}

		result, err := client.CreateZone(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to create zone: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Zone)
		}

		if result.Zone == nil {
			fmt.Println("Zone created successfully (No details returned)")
			return nil
		}

		fmt.Printf("Zone created successfully: %s (%s)\n", result.Zone.ZoneName, result.Zone.ZoneID)
		return nil
	},
}

var dnsUpdateZoneCmd = &cobra.Command{
	Use:   "update-zone",
	Short: "Update a DNS zone",
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneID, _ := cmd.Flags().GetString("zone-id")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.UpdateZoneInput{
			Description: description,
			ZoneStatus:  status,
		}

		result, err := client.UpdateZone(ctx, zoneID, input)
		if err != nil {
			return fmt.Errorf("failed to update zone: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Zone)
		}

		fmt.Printf("Zone updated successfully: %s\n", result.Zone.ZoneID)
		return nil
	},
}

var dnsDeleteZoneCmd = &cobra.Command{
	Use:     "delete-zones",
	Aliases: []string{"delete-zone"},
	Short:   "Delete DNS zones",
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneIDs, _ := cmd.Flags().GetStringSlice("zone-ids")
		// Support singular flag if plural is empty? (Requires adding flag first)
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeleteZones(ctx, zoneIDs)
		if err != nil {
			return fmt.Errorf("failed to delete zones: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Zone deletion initiated for %d zone(s)\n", len(zoneIDs))
		return nil
	},
}
