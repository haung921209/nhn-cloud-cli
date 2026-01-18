package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/privatedns"
	"github.com/spf13/cobra"
)

func init() {
	privateDNSCmd.AddCommand(pdnsDescribeZonesCmd)
	privateDNSCmd.AddCommand(pdnsGetZoneCmd)
	privateDNSCmd.AddCommand(pdnsCreateZoneCmd)
	privateDNSCmd.AddCommand(pdnsUpdateZoneCmd)
	privateDNSCmd.AddCommand(pdnsDeleteZoneCmd)

	pdnsGetZoneCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsGetZoneCmd.MarkFlagRequired("zone-id")

	pdnsCreateZoneCmd.Flags().String("name", "", "Zone name (required)")
	pdnsCreateZoneCmd.Flags().String("description", "", "Description")
	pdnsCreateZoneCmd.Flags().String("vpc-id", "", "VPC ID (required)")
	pdnsCreateZoneCmd.MarkFlagRequired("name")
	pdnsCreateZoneCmd.MarkFlagRequired("vpc-id")

	pdnsUpdateZoneCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsUpdateZoneCmd.Flags().String("description", "", "Description")
	pdnsUpdateZoneCmd.MarkFlagRequired("zone-id")

	pdnsDeleteZoneCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsDeleteZoneCmd.MarkFlagRequired("zone-id")
}

var pdnsDescribeZonesCmd = &cobra.Command{
	Use:     "describe-zones",
	Aliases: []string{"list-zones"},
	Short:   "List all private DNS zones",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		result, err := client.ListZones(context.Background())
		if err != nil {
			exitWithError("Failed to list zones", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tVPC_ID\tRECORDS\tSTATE")
		for _, z := range result.Zones {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", z.ID, z.Name, z.VPCID, z.RecordCount, z.State)
		}
		w.Flush()
	},
}

var pdnsGetZoneCmd = &cobra.Command{
	Use:     "describe-zone",
	Aliases: []string{"get-zone"},
	Short:   "Get zone details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("zone-id")

		result, err := client.GetZone(ctx, id)
		if err != nil {
			exitWithError("Failed to get zone", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		z := result.Zone
		fmt.Printf("ID:           %s\n", z.ID)
		fmt.Printf("Name:         %s\n", z.Name)
		fmt.Printf("Description:  %s\n", z.Description)
		fmt.Printf("VPC ID:       %s\n", z.VPCID)
		fmt.Printf("Record Count: %d\n", z.RecordCount)
		fmt.Printf("State:        %s\n", z.State)
		fmt.Printf("Created:      %s\n", z.CreatedAt)
		fmt.Printf("Updated:      %s\n", z.UpdatedAt)
	},
}

var pdnsCreateZoneCmd = &cobra.Command{
	Use:   "create-zone",
	Short: "Create a new private DNS zone",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		vpcID, _ := cmd.Flags().GetString("vpc-id")

		input := &privatedns.CreateZoneInput{
			Name:        name,
			Description: description,
			VPCID:       vpcID,
		}

		result, err := client.CreateZone(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create zone", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Zone created: %s\n", result.Zone.ID)
		fmt.Printf("Name: %s\n", result.Zone.Name)
		fmt.Printf("State: %s\n", result.Zone.State)
	},
}

var pdnsUpdateZoneCmd = &cobra.Command{
	Use:   "update-zone",
	Short: "Update a private DNS zone",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		id, _ := cmd.Flags().GetString("zone-id")
		description, _ := cmd.Flags().GetString("description")

		input := &privatedns.UpdateZoneInput{
			Description: description,
		}

		result, err := client.UpdateZone(context.Background(), id, input)
		if err != nil {
			exitWithError("Failed to update zone", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Zone updated: %s\n", result.Zone.ID)
	},
}

var pdnsDeleteZoneCmd = &cobra.Command{
	Use:   "delete-zone",
	Short: "Delete a private DNS zone",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		id, _ := cmd.Flags().GetString("zone-id")
		if err := client.DeleteZone(context.Background(), id); err != nil {
			exitWithError("Failed to delete zone", err)
		}
		fmt.Printf("Zone %s deleted\n", id)
	},
}
