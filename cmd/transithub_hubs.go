package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/transithub"
	"github.com/spf13/cobra"
)

func init() {
	transitHubCmd.AddCommand(thDescribeHubsCmd)
	transitHubCmd.AddCommand(thCreateHubCmd)
	transitHubCmd.AddCommand(thUpdateHubCmd)
	transitHubCmd.AddCommand(thDeleteHubCmd)

	thCreateHubCmd.Flags().String("name", "", "Hub name (required)")
	thCreateHubCmd.Flags().String("description", "", "Hub description")
	thCreateHubCmd.MarkFlagRequired("name")

	thUpdateHubCmd.Flags().String("hub-id", "", "Transit Hub ID (required)")
	thUpdateHubCmd.Flags().String("name", "", "Hub name")
	thUpdateHubCmd.Flags().String("description", "", "Hub description")
	thUpdateHubCmd.MarkFlagRequired("hub-id")

	thDeleteHubCmd.Flags().String("hub-id", "", "Transit Hub ID (required)")
	thDeleteHubCmd.MarkFlagRequired("hub-id")
}

var thDescribeHubsCmd = &cobra.Command{
	Use:     "describe-transit-hubs",
	Aliases: []string{"list-transit-hubs", "describe-hubs"},
	Short:   "List or describe Transit Hubs",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()

		// If args/flags supported specific ID get, we'd do it here.
		// For now, listing all.
		result, err := client.ListTransitHubs(ctx)
		if err != nil {
			exitWithError("Failed to list transit hubs", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tSTATE\tCREATED")
		for _, th := range result.TransitHubs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				th.ID, th.Name, th.Status, th.State, th.CreatedAt)
		}
		w.Flush()
	},
}

var thCreateHubCmd = &cobra.Command{
	Use:     "create-transit-hub",
	Aliases: []string{"create-hub"},
	Short:   "Create a new Transit Hub",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")

		input := &transithub.CreateTransitHubInput{
			Name:        name,
			Description: desc,
		}

		result, err := client.CreateTransitHub(ctx, input)
		if err != nil {
			exitWithError("Failed to create transit hub", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Transit Hub created: %s (%s)\n", result.TransitHub.Name, result.TransitHub.ID)
	},
}

var thUpdateHubCmd = &cobra.Command{
	Use:     "update-transit-hub",
	Aliases: []string{"update-hub"},
	Short:   "Update a Transit Hub",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("hub-id")
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")

		input := &transithub.UpdateTransitHubInput{
			Name:        name,
			Description: desc,
		}

		result, err := client.UpdateTransitHub(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update transit hub", err)
		}

		fmt.Printf("Transit Hub updated: %s\n", result.TransitHub.ID)
	},
}

var thDeleteHubCmd = &cobra.Command{
	Use:     "delete-transit-hub",
	Aliases: []string{"delete-hub"},
	Short:   "Delete a Transit Hub",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("hub-id")

		if err := client.DeleteTransitHub(ctx, id); err != nil {
			exitWithError("Failed to delete transit hub", err)
		}

		fmt.Printf("Transit Hub %s deleted\n", id)
	},
}
