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
	mirroringCmd.AddCommand(mirroringDescribeFiltersCmd)
	mirroringCmd.AddCommand(mirroringGetFilterCmd)
	mirroringCmd.AddCommand(mirroringCreateFilterCmd)
	mirroringCmd.AddCommand(mirroringUpdateFilterCmd)
	mirroringCmd.AddCommand(mirroringDeleteFilterCmd)

	mirroringGetFilterCmd.Flags().String("filter-id", "", "Filter ID (required)")
	mirroringGetFilterCmd.MarkFlagRequired("filter-id")

	mirroringCreateFilterCmd.Flags().String("filter-group-id", "", "Filter group ID (required)")
	mirroringCreateFilterCmd.Flags().String("name", "", "Filter name (required)")
	mirroringCreateFilterCmd.Flags().String("protocol", "", "Protocol: tcp, udp, icmp")
	mirroringCreateFilterCmd.Flags().String("source-cidr", "", "Source CIDR")
	mirroringCreateFilterCmd.Flags().String("dest-cidr", "", "Destination CIDR")
	mirroringCreateFilterCmd.Flags().Int("source-port-min", 0, "Source port min")
	mirroringCreateFilterCmd.Flags().Int("source-port-max", 0, "Source port max")
	mirroringCreateFilterCmd.Flags().Int("dest-port-min", 0, "Destination port min")
	mirroringCreateFilterCmd.Flags().Int("dest-port-max", 0, "Destination port max")
	mirroringCreateFilterCmd.Flags().String("action", "accept", "Action: accept, drop")
	mirroringCreateFilterCmd.Flags().String("description", "", "Description")
	mirroringCreateFilterCmd.MarkFlagRequired("filter-group-id")
	mirroringCreateFilterCmd.MarkFlagRequired("name")

	mirroringUpdateFilterCmd.Flags().String("filter-id", "", "Filter ID (required)")
	mirroringUpdateFilterCmd.Flags().String("name", "", "Filter name")
	mirroringUpdateFilterCmd.Flags().String("protocol", "", "Protocol")
	mirroringUpdateFilterCmd.Flags().String("source-cidr", "", "Source CIDR")
	mirroringUpdateFilterCmd.Flags().String("dest-cidr", "", "Destination CIDR")
	mirroringUpdateFilterCmd.Flags().String("action", "", "Action: accept, drop")
	mirroringUpdateFilterCmd.Flags().String("description", "", "Description")
	mirroringUpdateFilterCmd.MarkFlagRequired("filter-id")

	mirroringDeleteFilterCmd.Flags().String("filter-id", "", "Filter ID (required)")
	mirroringDeleteFilterCmd.MarkFlagRequired("filter-id")
}

var mirroringDescribeFiltersCmd = &cobra.Command{
	Use:     "describe-filters",
	Aliases: []string{"list-filters"},
	Short:   "List all mirroring filters",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		result, err := client.ListFilters(context.Background())
		if err != nil {
			exitWithError("Failed to list filters", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tPROTOCOL\tACTION\tSTATE")
		for _, f := range result.MirroringFilters {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				f.ID, f.Name, f.Protocol, f.Action, f.State)
		}
		w.Flush()
	},
}

var mirroringGetFilterCmd = &cobra.Command{
	Use:     "describe-filter",
	Aliases: []string{"get-filter"},
	Short:   "Get filter details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("filter-id")

		result, err := client.GetFilter(ctx, id)
		if err != nil {
			exitWithError("Failed to get filter", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		f := result.MirroringFilter
		fmt.Printf("ID:              %s\n", f.ID)
		fmt.Printf("Name:            %s\n", f.Name)
		fmt.Printf("Description:     %s\n", f.Description)
		fmt.Printf("Filter Group ID: %s\n", f.FilterGroupID)
		fmt.Printf("Protocol:        %s\n", f.Protocol)
		fmt.Printf("Source CIDR:     %s\n", f.SourceCIDR)
		fmt.Printf("Dest CIDR:       %s\n", f.DestCIDR)
		fmt.Printf("Source Port:     %d-%d\n", f.SourcePortMin, f.SourcePortMax)
		fmt.Printf("Dest Port:       %d-%d\n", f.DestPortMin, f.DestPortMax)
		fmt.Printf("Action:          %s\n", f.Action)
		fmt.Printf("State:           %s\n", f.State)
		fmt.Printf("Created At:      %s\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
	},
}

var mirroringCreateFilterCmd = &cobra.Command{
	Use:   "create-filter",
	Short: "Create a new filter",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		filterGroupID, _ := cmd.Flags().GetString("filter-group-id")
		name, _ := cmd.Flags().GetString("name")
		protocol, _ := cmd.Flags().GetString("protocol")
		sourceCIDR, _ := cmd.Flags().GetString("source-cidr")
		destCIDR, _ := cmd.Flags().GetString("dest-cidr")
		sourcePortMin, _ := cmd.Flags().GetInt("source-port-min")
		sourcePortMax, _ := cmd.Flags().GetInt("source-port-max")
		destPortMin, _ := cmd.Flags().GetInt("dest-port-min")
		destPortMax, _ := cmd.Flags().GetInt("dest-port-max")
		action, _ := cmd.Flags().GetString("action")
		description, _ := cmd.Flags().GetString("description")

		input := &mirroring.CreateFilterInput{
			FilterGroupID: filterGroupID,
			Name:          name,
			Protocol:      protocol,
			SourceCIDR:    sourceCIDR,
			DestCIDR:      destCIDR,
			SourcePortMin: sourcePortMin,
			SourcePortMax: sourcePortMax,
			DestPortMin:   destPortMin,
			DestPortMax:   destPortMax,
			Action:        action,
			Description:   description,
		}

		result, err := client.CreateFilter(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create filter", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Filter created: %s\n", result.MirroringFilter.ID)
		fmt.Printf("Name: %s\n", result.MirroringFilter.Name)
	},
}

var mirroringUpdateFilterCmd = &cobra.Command{
	Use:   "update-filter",
	Short: "Update a filter",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		id, _ := cmd.Flags().GetString("filter-id")
		name, _ := cmd.Flags().GetString("name")
		protocol, _ := cmd.Flags().GetString("protocol")
		sourceCIDR, _ := cmd.Flags().GetString("source-cidr")
		destCIDR, _ := cmd.Flags().GetString("dest-cidr")
		action, _ := cmd.Flags().GetString("action")
		description, _ := cmd.Flags().GetString("description")

		input := &mirroring.UpdateFilterInput{
			Name:        name,
			Protocol:    protocol,
			SourceCIDR:  sourceCIDR,
			DestCIDR:    destCIDR,
			Action:      action,
			Description: description,
		}

		result, err := client.UpdateFilter(context.Background(), id, input)
		if err != nil {
			exitWithError("Failed to update filter", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Filter updated: %s\n", result.MirroringFilter.ID)
	},
}

var mirroringDeleteFilterCmd = &cobra.Command{
	Use:   "delete-filter",
	Short: "Delete a filter",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		id, _ := cmd.Flags().GetString("filter-id")
		if err := client.DeleteFilter(context.Background(), id); err != nil {
			exitWithError("Failed to delete filter", err)
		}
		fmt.Printf("Filter %s deleted\n", id)
	},
}
