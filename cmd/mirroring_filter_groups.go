package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/mirroring"
	"github.com/spf13/cobra"
)

func init() {
	mirroringCmd.AddCommand(mirroringDescribeFilterGroupsCmd)
	mirroringCmd.AddCommand(mirroringGetFilterGroupCmd)
	mirroringCmd.AddCommand(mirroringCreateFilterGroupCmd)
	mirroringCmd.AddCommand(mirroringUpdateFilterGroupCmd)
	mirroringCmd.AddCommand(mirroringDeleteFilterGroupCmd)

	mirroringGetFilterGroupCmd.Flags().String("filter-group-id", "", "Filter group ID (required)")
	mirroringGetFilterGroupCmd.MarkFlagRequired("filter-group-id")

	mirroringCreateFilterGroupCmd.Flags().String("name", "", "Filter group name (required)")
	mirroringCreateFilterGroupCmd.Flags().String("description", "", "Description")
	mirroringCreateFilterGroupCmd.MarkFlagRequired("name")

	mirroringUpdateFilterGroupCmd.Flags().String("filter-group-id", "", "Filter group ID (required)")
	mirroringUpdateFilterGroupCmd.Flags().String("name", "", "Filter group name")
	mirroringUpdateFilterGroupCmd.Flags().String("description", "", "Description")
	mirroringUpdateFilterGroupCmd.MarkFlagRequired("filter-group-id")

	mirroringDeleteFilterGroupCmd.Flags().String("filter-group-id", "", "Filter group ID (required)")
	mirroringDeleteFilterGroupCmd.MarkFlagRequired("filter-group-id")
}

var mirroringDescribeFilterGroupsCmd = &cobra.Command{
	Use:     "describe-filter-groups",
	Aliases: []string{"list-filter-groups"},
	Short:   "List all mirroring filter groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		result, err := client.ListFilterGroups(context.Background())
		if err != nil {
			exitWithError("Failed to list filter groups", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATE\tFILTER_COUNT")
		for _, fg := range result.MirroringFilterGroups {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
				fg.ID, fg.Name, fg.State, len(fg.FilterIDs))
		}
		w.Flush()
	},
}

var mirroringGetFilterGroupCmd = &cobra.Command{
	Use:     "describe-filter-group",
	Aliases: []string{"get-filter-group"},
	Short:   "Get filter group details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("filter-group-id")

		result, err := client.GetFilterGroup(ctx, id)
		if err != nil {
			exitWithError("Failed to get filter group", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fg := result.MirroringFilterGroup
		fmt.Printf("ID:          %s\n", fg.ID)
		fmt.Printf("Name:        %s\n", fg.Name)
		fmt.Printf("Description: %s\n", fg.Description)
		fmt.Printf("State:       %s\n", fg.State)
		fmt.Printf("Filter IDs:  %v\n", fg.FilterIDs)
		fmt.Printf("Created At:  %s\n", fg.CreatedAt.Format("2006-01-02 15:04:05"))
	},
}

var mirroringCreateFilterGroupCmd = &cobra.Command{
	Use:   "create-filter-group",
	Short: "Create a new filter group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &mirroring.CreateFilterGroupInput{
			Name:        name,
			Description: description,
		}

		result, err := client.CreateFilterGroup(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create filter group", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Filter group created: %s\n", result.MirroringFilterGroup.ID)
		fmt.Printf("Name: %s\n", result.MirroringFilterGroup.Name)
	},
}

var mirroringUpdateFilterGroupCmd = &cobra.Command{
	Use:   "update-filter-group",
	Short: "Update a filter group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		id, _ := cmd.Flags().GetString("filter-group-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &mirroring.UpdateFilterGroupInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateFilterGroup(context.Background(), id, input)
		if err != nil {
			exitWithError("Failed to update filter group", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Filter group updated: %s\n", result.MirroringFilterGroup.ID)
	},
}

var mirroringDeleteFilterGroupCmd = &cobra.Command{
	Use:   "delete-filter-group",
	Short: "Delete a filter group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		id, _ := cmd.Flags().GetString("filter-group-id")
		if err := client.DeleteFilterGroup(context.Background(), id); err != nil {
			exitWithError("Failed to delete filter group", err)
		}
		fmt.Printf("Filter group %s deleted\n", id)
	},
}
