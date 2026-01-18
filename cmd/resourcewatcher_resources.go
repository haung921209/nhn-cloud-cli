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
	resourceWatcherCmd.AddCommand(rwDescribeResourceGroupsCmd)
	resourceWatcherCmd.AddCommand(rwDescribeResourceTagsCmd)
}

var rwDescribeResourceGroupsCmd = &cobra.Command{
	Use:     "describe-resource-groups",
	Aliases: []string{"list-resource-groups", "resource-groups"},
	Short:   "List resource groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		result, err := client.ListResourceGroups(ctx)
		if err != nil {
			exitWithError("Failed to list resource groups", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "GROUP_ID\tNAME\tDESCRIPTION\tCREATED")
		for _, g := range result.ResourceGroups {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				g.ResourceGroupID, g.ResourceGroupName, g.Description, g.CreatedDateTime)
		}
		w.Flush()
	},
}

var rwDescribeResourceTagsCmd = &cobra.Command{
	Use:     "describe-resource-tags",
	Aliases: []string{"list-resource-tags", "resource-tags"},
	Short:   "List resource tags",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		result, err := client.ListResourceTags(ctx)
		if err != nil {
			exitWithError("Failed to list resource tags", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TAG_ID\tNAME\tKEY\tVALUE\tCREATED")
		for _, t := range result.ResourceTags {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				t.ResourceTagID, t.ResourceTagName, t.TagKey, t.TagValue, t.CreatedDateTime)
		}
		w.Flush()
	},
}
