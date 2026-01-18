package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	resourceWatcherCmd.AddCommand(rwDescribeEventsCmd)
	resourceWatcherCmd.AddCommand(rwGetEventCmd)

	// Get flags
	rwGetEventCmd.Flags().String("product-id", "", "Product ID (required)")
	rwGetEventCmd.Flags().String("event-id", "", "Event ID (required)")
	rwGetEventCmd.MarkFlagRequired("product-id")
	rwGetEventCmd.MarkFlagRequired("event-id")
}

var rwDescribeEventsCmd = &cobra.Command{
	Use:     "describe-events",
	Aliases: []string{"list-events", "events"},
	Short:   "List available events",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		result, err := client.ListEvents(ctx)
		if err != nil {
			exitWithError("Failed to list events", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "EVENT_ID\tPRODUCT\tNAME\tTYPE")
		for _, e := range result.Events {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				e.EventID, e.ProductID, e.EventName, e.EventType)
		}
		w.Flush()
	},
}

var rwGetEventCmd = &cobra.Command{
	Use:     "describe-event",
	Aliases: []string{"get-event"},
	Short:   "Get event details",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()
		productID, _ := cmd.Flags().GetString("product-id")
		eventID, _ := cmd.Flags().GetString("event-id")

		result, err := client.GetEvent(ctx, productID, eventID)
		if err != nil {
			exitWithError("Failed to get event", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		e := result.Event
		fmt.Printf("Event ID:    %s\n", e.EventID)
		fmt.Printf("Product ID:  %s\n", e.ProductID)
		fmt.Printf("Event Name:  %s\n", e.EventName)
		fmt.Printf("Event Type:  %s\n", e.EventType)
		fmt.Printf("Description: %s\n", e.Description)
	},
}
