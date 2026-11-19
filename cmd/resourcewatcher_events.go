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
