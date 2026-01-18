package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/dnsplus"
	"github.com/spf13/cobra"
)

func init() {
	dnsplusCmd.AddCommand(dnsDescribeRecordSetsCmd)
	dnsplusCmd.AddCommand(dnsCreateRecordSetCmd)
	dnsplusCmd.AddCommand(dnsDeleteRecordSetCmd)

	dnsDescribeRecordSetsCmd.Flags().String("zone-id", "", "Zone ID (required)")
	dnsDescribeRecordSetsCmd.MarkFlagRequired("zone-id")

	dnsCreateRecordSetCmd.Flags().String("zone-id", "", "Zone ID (required)")
	dnsCreateRecordSetCmd.Flags().String("name", "", "Record set name")
	dnsCreateRecordSetCmd.Flags().String("type", "", "Record type (A, AAAA, CNAME, MX, TXT, etc.)")
	dnsCreateRecordSetCmd.Flags().Int("ttl", 300, "TTL in seconds")
	dnsCreateRecordSetCmd.Flags().StringSlice("record", nil, "Record content (can specify multiple)")
	dnsCreateRecordSetCmd.MarkFlagRequired("zone-id")
	dnsCreateRecordSetCmd.MarkFlagRequired("name")
	dnsCreateRecordSetCmd.MarkFlagRequired("type")
	dnsCreateRecordSetCmd.MarkFlagRequired("record")

	dnsDeleteRecordSetCmd.Flags().String("zone-id", "", "Zone ID (required)")
	dnsDeleteRecordSetCmd.Flags().StringSlice("recordset-ids", nil, "Record set IDs (required)")
	dnsDeleteRecordSetCmd.MarkFlagRequired("zone-id")
	dnsDeleteRecordSetCmd.MarkFlagRequired("recordset-ids")
}

var dnsDescribeRecordSetsCmd = &cobra.Command{
	Use:     "describe-record-sets",
	Aliases: []string{"list-record-sets"},
	Short:   "List record sets in a zone",
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneID, _ := cmd.Flags().GetString("zone-id")
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListRecordSets(ctx, zoneID)
		if err != nil {
			return fmt.Errorf("failed to list record sets: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.RecordSetList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tTTL\tRECORDS")
		for _, rs := range result.RecordSetList {
			records := make([]string, len(rs.RecordList))
			for i, r := range rs.RecordList {
				records[i] = r.RecordContent
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				rs.RecordSetID,
				rs.RecordSetName,
				rs.RecordSetType,
				rs.TTL,
				strings.Join(records, ", "),
			)
		}
		return w.Flush()
	},
}

var dnsCreateRecordSetCmd = &cobra.Command{
	Use:   "create-record-set",
	Short: "Create a record set",
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneID, _ := cmd.Flags().GetString("zone-id")
		name, _ := cmd.Flags().GetString("name")
		recordType, _ := cmd.Flags().GetString("type")
		ttl, _ := cmd.Flags().GetInt("ttl")
		records, _ := cmd.Flags().GetStringSlice("record")

		client := newDNSPlusClient()
		ctx := context.Background()

		recordList := make([]dnsplus.Record, len(records))
		for i, r := range records {
			recordList[i] = dnsplus.Record{RecordContent: r}
		}

		input := &dnsplus.CreateRecordSetInput{
			RecordSetName: name,
			RecordSetType: recordType,
			TTL:           ttl,
			RecordList:    recordList,
		}

		result, err := client.CreateRecordSet(ctx, zoneID, input)
		if err != nil {
			return fmt.Errorf("failed to create record set: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.RecordSet)
		}

		fmt.Printf("Record set created successfully: %s (%s)\n", result.RecordSet.RecordSetName, result.RecordSet.RecordSetID)
		return nil
	},
}

var dnsDeleteRecordSetCmd = &cobra.Command{
	Use:   "delete-record-sets",
	Short: "Delete record sets",
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneID, _ := cmd.Flags().GetString("zone-id")
		recordsetIDs, _ := cmd.Flags().GetStringSlice("recordset-ids")

		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeleteRecordSets(ctx, zoneID, recordsetIDs)
		if err != nil {
			return fmt.Errorf("failed to delete record sets: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Record set(s) deleted successfully\n")
		return nil
	},
}
