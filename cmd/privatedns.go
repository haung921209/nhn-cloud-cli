package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/privatedns"
	"github.com/spf13/cobra"
)

var privateDNSCmd = &cobra.Command{
	Use:     "private-dns",
	Aliases: []string{"pdns", "privatedns"},
	Short:   "Manage Private DNS zones and records",
}

var pdnsZoneCmd = &cobra.Command{
	Use:     "zone",
	Aliases: []string{"zones"},
	Short:   "Manage Private DNS zones",
}

var pdnsRecordCmd = &cobra.Command{
	Use:     "record",
	Aliases: []string{"records", "rrset", "rrsets"},
	Short:   "Manage DNS record sets",
}

func init() {
	rootCmd.AddCommand(privateDNSCmd)

	// Zone commands
	privateDNSCmd.AddCommand(pdnsZoneCmd)
	pdnsZoneCmd.AddCommand(pdnsZoneListCmd)
	pdnsZoneCmd.AddCommand(pdnsZoneGetCmd)
	pdnsZoneCmd.AddCommand(pdnsZoneCreateCmd)
	pdnsZoneCmd.AddCommand(pdnsZoneUpdateCmd)
	pdnsZoneCmd.AddCommand(pdnsZoneDeleteCmd)

	// Record commands
	privateDNSCmd.AddCommand(pdnsRecordCmd)
	pdnsRecordCmd.AddCommand(pdnsRecordListCmd)
	pdnsRecordCmd.AddCommand(pdnsRecordGetCmd)
	pdnsRecordCmd.AddCommand(pdnsRecordCreateCmd)
	pdnsRecordCmd.AddCommand(pdnsRecordUpdateCmd)
	pdnsRecordCmd.AddCommand(pdnsRecordDeleteCmd)

	// Zone create flags
	pdnsZoneCreateCmd.Flags().String("name", "", "Zone name (e.g., example.local) (required)")
	pdnsZoneCreateCmd.Flags().String("description", "", "Description")
	pdnsZoneCreateCmd.Flags().String("vpc-id", "", "VPC ID (required)")
	pdnsZoneCreateCmd.MarkFlagRequired("name")
	pdnsZoneCreateCmd.MarkFlagRequired("vpc-id")

	// Zone update flags
	pdnsZoneUpdateCmd.Flags().String("description", "", "Description")

	// Record list flags
	pdnsRecordListCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsRecordListCmd.MarkFlagRequired("zone-id")

	// Record get flags
	pdnsRecordGetCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsRecordGetCmd.MarkFlagRequired("zone-id")

	// Record create flags
	pdnsRecordCreateCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsRecordCreateCmd.Flags().String("name", "", "Record name (required)")
	pdnsRecordCreateCmd.Flags().String("type", "", "Record type: A, AAAA, CNAME, MX, TXT, etc. (required)")
	pdnsRecordCreateCmd.Flags().Int("ttl", 300, "TTL in seconds")
	pdnsRecordCreateCmd.Flags().StringSlice("records", []string{}, "Record values (comma-separated or multiple --records)")
	pdnsRecordCreateCmd.MarkFlagRequired("zone-id")
	pdnsRecordCreateCmd.MarkFlagRequired("name")
	pdnsRecordCreateCmd.MarkFlagRequired("type")
	pdnsRecordCreateCmd.MarkFlagRequired("records")

	// Record update flags
	pdnsRecordUpdateCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsRecordUpdateCmd.Flags().Int("ttl", 0, "TTL in seconds")
	pdnsRecordUpdateCmd.Flags().StringSlice("records", []string{}, "Record values")
	pdnsRecordUpdateCmd.MarkFlagRequired("zone-id")

	// Record delete flags
	pdnsRecordDeleteCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsRecordDeleteCmd.MarkFlagRequired("zone-id")
}

func newPrivateDNSClient() *privatedns.Client {
	return privatedns.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

// ================================
// Zone Commands
// ================================

var pdnsZoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all private DNS zones",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		result, err := client.ListZones(context.Background())
		if err != nil {
			exitWithError("Failed to list zones", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var pdnsZoneGetCmd = &cobra.Command{
	Use:   "get [zone-id]",
	Short: "Get zone details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		result, err := client.GetZone(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get zone", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var pdnsZoneCreateCmd = &cobra.Command{
	Use:   "create",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Zone created: %s\n", result.Zone.ID)
		fmt.Printf("Name: %s\n", result.Zone.Name)
		fmt.Printf("State: %s\n", result.Zone.State)
	},
}

var pdnsZoneUpdateCmd = &cobra.Command{
	Use:   "update [zone-id]",
	Short: "Update a private DNS zone",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		description, _ := cmd.Flags().GetString("description")

		input := &privatedns.UpdateZoneInput{
			Description: description,
		}

		result, err := client.UpdateZone(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update zone", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Zone updated: %s\n", result.Zone.ID)
	},
}

var pdnsZoneDeleteCmd = &cobra.Command{
	Use:   "delete [zone-id]",
	Short: "Delete a private DNS zone",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		if err := client.DeleteZone(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete zone", err)
		}
		fmt.Printf("Zone %s deleted\n", args[0])
	},
}

// ================================
// Record Commands
// ================================

var pdnsRecordListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all record sets in a zone",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		zoneID, _ := cmd.Flags().GetString("zone-id")

		result, err := client.ListRRSets(context.Background(), zoneID)
		if err != nil {
			exitWithError("Failed to list records", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tTTL\tRECORDS\tSTATE")
		for _, r := range result.RRSets {
			records := strings.Join(r.Records, ", ")
			if len(records) > 40 {
				records = records[:37] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n", r.ID, r.Name, r.Type, r.TTL, records, r.State)
		}
		w.Flush()
	},
}

var pdnsRecordGetCmd = &cobra.Command{
	Use:   "get [record-id]",
	Short: "Get record set details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		zoneID, _ := cmd.Flags().GetString("zone-id")

		result, err := client.GetRRSet(context.Background(), zoneID, args[0])
		if err != nil {
			exitWithError("Failed to get record", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		r := result.RRSet
		fmt.Printf("ID:      %s\n", r.ID)
		fmt.Printf("Zone ID: %s\n", r.ZoneID)
		fmt.Printf("Name:    %s\n", r.Name)
		fmt.Printf("Type:    %s\n", r.Type)
		fmt.Printf("TTL:     %d\n", r.TTL)
		fmt.Printf("Records: %s\n", strings.Join(r.Records, ", "))
		fmt.Printf("State:   %s\n", r.State)
		fmt.Printf("Created: %s\n", r.CreatedAt)
		fmt.Printf("Updated: %s\n", r.UpdatedAt)
	},
}

var pdnsRecordCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new record set",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		zoneID, _ := cmd.Flags().GetString("zone-id")
		name, _ := cmd.Flags().GetString("name")
		recordType, _ := cmd.Flags().GetString("type")
		ttl, _ := cmd.Flags().GetInt("ttl")
		records, _ := cmd.Flags().GetStringSlice("records")

		input := &privatedns.CreateRRSetInput{
			Name:    name,
			Type:    recordType,
			TTL:     ttl,
			Records: records,
		}

		result, err := client.CreateRRSet(context.Background(), zoneID, input)
		if err != nil {
			exitWithError("Failed to create record", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Record created: %s\n", result.RRSet.ID)
		fmt.Printf("Name: %s\n", result.RRSet.Name)
		fmt.Printf("Type: %s\n", result.RRSet.Type)
		fmt.Printf("TTL: %d\n", result.RRSet.TTL)
	},
}

var pdnsRecordUpdateCmd = &cobra.Command{
	Use:   "update [record-id]",
	Short: "Update a record set",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		zoneID, _ := cmd.Flags().GetString("zone-id")
		ttl, _ := cmd.Flags().GetInt("ttl")
		records, _ := cmd.Flags().GetStringSlice("records")

		input := &privatedns.UpdateRRSetInput{
			TTL:     ttl,
			Records: records,
		}

		result, err := client.UpdateRRSet(context.Background(), zoneID, args[0], input)
		if err != nil {
			exitWithError("Failed to update record", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Record updated: %s\n", result.RRSet.ID)
	},
}

var pdnsRecordDeleteCmd = &cobra.Command{
	Use:   "delete [record-id]",
	Short: "Delete a record set",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		zoneID, _ := cmd.Flags().GetString("zone-id")

		if err := client.DeleteRRSet(context.Background(), zoneID, args[0]); err != nil {
			exitWithError("Failed to delete record", err)
		}
		fmt.Printf("Record %s deleted\n", args[0])
	},
}
