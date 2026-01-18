package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	flowlogCmd.AddCommand(flDescribePortsCmd)
	flowlogCmd.AddCommand(flGetPortCmd)

	flGetPortCmd.Flags().String("port-id", "", "Logging Port ID (required)")
	flGetPortCmd.MarkFlagRequired("port-id")
}

var flDescribePortsCmd = &cobra.Command{
	Use:     "describe-logging-ports",
	Aliases: []string{"list-logging-ports", "list-ports"},
	Short:   "List all flow log logging ports",
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		result, err := client.ListLoggingPorts(context.Background())
		if err != nil {
			exitWithError("Failed to list logging ports", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tLOGGER_ID\tPORT_ID\tSTATE")
		for _, port := range result.LoggingPorts {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				port.ID, port.LoggerID, port.PortID, port.State)
		}
		w.Flush()
	},
}

var flGetPortCmd = &cobra.Command{
	Use:     "describe-logging-port",
	Aliases: []string{"get-logging-port", "get-port"},
	Short:   "Get flow log logging port details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("port-id")

		result, err := client.GetLoggingPort(ctx, id)
		if err != nil {
			exitWithError("Failed to get logging port", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		port := result.LoggingPort
		fmt.Printf("ID:         %s\n", port.ID)
		fmt.Printf("Logger ID:  %s\n", port.LoggerID)
		fmt.Printf("Port ID:    %s\n", port.PortID)
		fmt.Printf("State:      %s\n", port.State)
		fmt.Printf("Created At: %s\n", port.CreatedAt)
	},
}
