package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	internetGatewayCmd.AddCommand(igwDescribeExternalNetworksCmd)
}

var igwDescribeExternalNetworksCmd = &cobra.Command{
	Use:     "describe-external-networks",
	Aliases: []string{"list-external-networks"},
	Short:   "List available external networks",
	Run: func(cmd *cobra.Command, args []string) {
		client := newInternetGatewayClient()
		result, err := client.ListExternalNetworks(context.Background())
		if err != nil {
			exitWithError("Failed to list external networks", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tROUTER_EXTERNAL")
		for _, net := range result.Networks {
			fmt.Fprintf(w, "%s\t%s\t%v\n",
				net.ID, net.Name, net.RouterExternal)
		}
		w.Flush()
	},
}
