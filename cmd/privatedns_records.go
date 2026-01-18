package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/privatedns"
	"github.com/spf13/cobra"
)

func init() {
	privateDNSCmd.AddCommand(pdnsDescribeRecordsCmd)
	privateDNSCmd.AddCommand(pdnsGetRecordCmd)
	privateDNSCmd.AddCommand(pdnsCreateRecordCmd)
	privateDNSCmd.AddCommand(pdnsUpdateRecordCmd)
	privateDNSCmd.AddCommand(pdnsDeleteRecordCmd)

	pdnsDescribeRecordsCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsDescribeRecordsCmd.MarkFlagRequired("zone-id")

	pdnsGetRecordCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsGetRecordCmd.Flags().String("record-id", "", "Record ID (required)")
	pdnsGetRecordCmd.MarkFlagRequired("zone-id")
	pdnsGetRecordCmd.MarkFlagRequired("record-id")

	pdnsCreateRecordCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsCreateRecordCmd.Flags().String("name", "", "Record name (required)")
	pdnsCreateRecordCmd.Flags().String("type", "", "Record type: A, AAAA, CNAME, MX, TXT (required)")
	pdnsCreateRecordCmd.Flags().Int("ttl", 300, "TTL in seconds")
	pdnsCreateRecordCmd.Flags().StringSlice("records", []string{}, "Record values")
	pdnsCreateRecordCmd.MarkFlagRequired("zone-id")
	pdnsCreateRecordCmd.MarkFlagRequired("name")
	pdnsCreateRecordCmd.MarkFlagRequired("type")
	pdnsCreateRecordCmd.MarkFlagRequired("records")

	pdnsUpdateRecordCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsUpdateRecordCmd.Flags().String("record-id", "", "Record ID (required)")
	pdnsUpdateRecordCmd.Flags().Int("ttl", 0, "TTL in seconds")
	pdnsUpdateRecordCmd.Flags().StringSlice("records", []string{}, "Record values")
	pdnsUpdateRecordCmd.MarkFlagRequired("zone-id")
	pdnsUpdateRecordCmd.MarkFlagRequired("record-id")

	pdnsDeleteRecordCmd.Flags().String("zone-id", "", "Zone ID (required)")
	pdnsDeleteRecordCmd.Flags().String("record-id", "", "Record ID (required)")
	pdnsDeleteRecordCmd.MarkFlagRequired("zone-id")
	pdnsDeleteRecordCmd.MarkFlagRequired("record-id")
}

var pdnsDescribeRecordsCmd = &cobra.Command{
	Use:     "describe-records",
	Aliases: []string{"list-records"},
	Short:   "List all record sets in a zone",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		ctx := context.Background()
		zoneID, _ := cmd.Flags().GetString("zone-id")

		result, err := client.ListRRSets(ctx, zoneID)
		if err != nil {
			exitWithError("Failed to list records", err)
		}

		if output == "json" {
			printJSON(result)
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

var pdnsGetRecordCmd = &cobra.Command{
	Use:     "describe-record",
	Aliases: []string{"get-record"},
	Short:   "Get record set details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		ctx := context.Background()
		zoneID, _ := cmd.Flags().GetString("zone-id")
		recordID, _ := cmd.Flags().GetString("record-id")

		result, err := client.GetRRSet(ctx, zoneID, recordID)
		if err != nil {
			exitWithError("Failed to get record", err)
		}

		if output == "json" {
			printJSON(result)
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

var pdnsCreateRecordCmd = &cobra.Command{
	Use:   "create-record",
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
			printJSON(result)
			return
		}

		fmt.Printf("Record created: %s\n", result.RRSet.ID)
		fmt.Printf("Name: %s\n", result.RRSet.Name)
		fmt.Printf("Type: %s\n", result.RRSet.Type)
		fmt.Printf("TTL: %d\n", result.RRSet.TTL)
	},
}

var pdnsUpdateRecordCmd = &cobra.Command{
	Use:   "update-record",
	Short: "Update a record set",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		zoneID, _ := cmd.Flags().GetString("zone-id")
		recordID, _ := cmd.Flags().GetString("record-id")
		ttl, _ := cmd.Flags().GetInt("ttl")
		records, _ := cmd.Flags().GetStringSlice("records")

		input := &privatedns.UpdateRRSetInput{
			TTL:     ttl,
			Records: records,
		}

		result, err := client.UpdateRRSet(context.Background(), zoneID, recordID, input)
		if err != nil {
			exitWithError("Failed to update record", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Record updated: %s\n", result.RRSet.ID)
	},
}

var pdnsDeleteRecordCmd = &cobra.Command{
	Use:   "delete-record",
	Short: "Delete a record set",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPrivateDNSClient()
		zoneID, _ := cmd.Flags().GetString("zone-id")
		recordID, _ := cmd.Flags().GetString("record-id")

		if err := client.DeleteRRSet(context.Background(), zoneID, recordID); err != nil {
			exitWithError("Failed to delete record", err)
		}
		fmt.Printf("Record %s deleted\n", recordID)
	},
}
