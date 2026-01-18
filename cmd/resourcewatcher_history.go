package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/resourcewatcher"
	"github.com/spf13/cobra"
)

func init() {
	resourceWatcherCmd.AddCommand(rwDescribeHistoryCmd)
	resourceWatcherCmd.AddCommand(rwGetHistoryCmd)

	// List flags
	rwDescribeHistoryCmd.Flags().String("alarm-id", "", "Alarm ID (required)")
	rwDescribeHistoryCmd.Flags().String("start", "", "Start datetime (ISO8601)")
	rwDescribeHistoryCmd.Flags().String("end", "", "End datetime (ISO8601)")
	rwDescribeHistoryCmd.Flags().Int("page", 0, "Page number")
	rwDescribeHistoryCmd.Flags().Int("size", 20, "Page size")
	rwDescribeHistoryCmd.MarkFlagRequired("alarm-id")

	// Get flags
	rwGetHistoryCmd.Flags().String("alarm-id", "", "Alarm ID (required)")
	rwGetHistoryCmd.Flags().String("history-id", "", "History ID (required)")
	rwGetHistoryCmd.MarkFlagRequired("alarm-id")
	rwGetHistoryCmd.MarkFlagRequired("history-id")
}

var rwDescribeHistoryCmd = &cobra.Command{
	Use:     "describe-alarm-history",
	Aliases: []string{"list-history", "history"},
	Short:   "List alarm history",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		alarmID, _ := cmd.Flags().GetString("alarm-id")
		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")

		input := &resourcewatcher.SearchAlarmHistoryInput{
			StartDateTime: start,
			EndDateTime:   end,
			Page:          page,
			Size:          size,
		}

		result, err := client.SearchAlarmHistory(ctx, alarmID, input)
		if err != nil {
			exitWithError("Failed to list alarm history", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "HISTORY_ID\tEVENT_NAME\tRESOURCE\tCREATED")
		for _, h := range result.Histories {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				h.AlarmHistoryID, h.EventName, h.ResourceName, h.CreatedDateTime)
		}
		w.Flush()
		fmt.Printf("\nTotal: %d\n", result.TotalCount)
	},
}

var rwGetHistoryCmd = &cobra.Command{
	Use:     "describe-alarm-history-detail",
	Aliases: []string{"get-history"},
	Short:   "Get alarm history details",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()
		alarmID, _ := cmd.Flags().GetString("alarm-id")
		historyID, _ := cmd.Flags().GetString("history-id")

		result, err := client.GetAlarmHistory(ctx, alarmID, historyID)
		if err != nil {
			exitWithError("Failed to get alarm history", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		h := result.History
		fmt.Printf("History ID:    %s\n", h.AlarmHistoryID)
		fmt.Printf("Alarm ID:      %s\n", h.AlarmID)
		fmt.Printf("Event ID:      %s\n", h.EventID)
		fmt.Printf("Event Name:    %s\n", h.EventName)
		fmt.Printf("Product ID:    %s\n", h.ProductID)
		fmt.Printf("Resource ID:   %s\n", h.ResourceID)
		fmt.Printf("Resource Name: %s\n", h.ResourceName)
		fmt.Printf("Created:       %s\n", h.CreatedDateTime)

		if len(h.AlarmSendResults) > 0 {
			fmt.Println("\nSend Results:")
			for _, r := range h.AlarmSendResults {
				fmt.Printf("  - %s (%s): %s at %s\n", r.TargetType, r.TargetID, r.SendStatus, r.SentDateTime)
			}
		}
	},
}
