package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/image"
	"github.com/spf13/cobra"
)

func init() {
	computeCmd.AddCommand(computeDescribeImagesCmd)
}

var computeDescribeImagesCmd = &cobra.Command{
	Use:     "describe-images",
	Aliases: []string{"images"},
	Short:   "List available compute images",
	Long:    "List available images for compute instances. This command uses the Glance Image API.",
	Run: func(cmd *cobra.Command, args []string) {
		creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
		client := image.NewClient(getRegion(), creds, nil, debug)
		ctx := context.Background()

		result, err := client.ListImages(ctx, nil)
		if err != nil {
			exitWithError("Failed to list images", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tVISIBILITY\tSIZE (MB)\tOS\tCREATED")
		for _, img := range result.Images {
			sizeInMB := img.Size / (1024 * 1024)
			createdDate := img.CreatedAt.Format("2006-01-02")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
				img.ID, img.Name, img.Status, img.Visibility, sizeInMB, img.OSDistro, createdDate)
		}
		w.Flush()
	},
}
