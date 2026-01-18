package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	computeCmd.AddCommand(computeDescribeFlavorsCmd)
}

var computeDescribeFlavorsCmd = &cobra.Command{
	Use:     "describe-flavors",
	Aliases: []string{"flavors"},
	Short:   "List available compute flavors",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		result, err := client.ListFlavors(ctx)
		if err != nil {
			exitWithError("Failed to list flavors", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tVCPUs\tRAM (MB)\tDISK (GB)")
		for _, f := range result.Flavors {
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\n",
				f.ID, f.Name, f.VCPUs, f.RAM, f.Disk)
		}
		w.Flush()
	},
}
