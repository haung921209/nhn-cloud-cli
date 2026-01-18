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
	transitHubCmd.AddCommand(thDescribeMulticastDomainsCmd)
	transitHubCmd.AddCommand(thCreateMulticastDomainCmd)
	transitHubCmd.AddCommand(thUpdateMulticastDomainCmd)
	transitHubCmd.AddCommand(thDeleteMulticastDomainCmd)

	thCreateMulticastDomainCmd.Flags().String("name", "", "Name (required)")
	thCreateMulticastDomainCmd.Flags().String("description", "", "Description")
	thCreateMulticastDomainCmd.Flags().String("transit-hub-id", "", "Transit Hub ID (required)")
	thCreateMulticastDomainCmd.MarkFlagRequired("name")
	thCreateMulticastDomainCmd.MarkFlagRequired("transit-hub-id")

	thUpdateMulticastDomainCmd.Flags().String("domain-id", "", "Multicast Domain ID (required)")
	thUpdateMulticastDomainCmd.Flags().String("name", "", "Name")
	thUpdateMulticastDomainCmd.Flags().String("description", "", "Description")
	thUpdateMulticastDomainCmd.MarkFlagRequired("domain-id")

	thDeleteMulticastDomainCmd.Flags().String("domain-id", "", "Multicast Domain ID (required)")
	thDeleteMulticastDomainCmd.MarkFlagRequired("domain-id")
}

var thDescribeMulticastDomainsCmd = &cobra.Command{
	Use:     "describe-multicast-domains",
	Aliases: []string{"list-multicast-domains"},
	Short:   "List multicast domains",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()

		result, err := client.ListMulticastDomains(ctx)
		if err != nil {
			exitWithError("Failed to list multicast domains", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTRANSIT_HUB_ID\tSTATUS")
		for _, m := range result.Domains {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", m.ID, m.Name, m.TransitHubID, m.Status)
		}
		w.Flush()
	},
}

var thCreateMulticastDomainCmd = &cobra.Command{
	Use:   "create-multicast-domain",
	Short: "Create a new multicast domain",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		thID, _ := cmd.Flags().GetString("transit-hub-id")

		input := &transithub.CreateMulticastDomainInput{
			Name:         name,
			Description:  desc,
			TransitHubID: thID,
		}

		result, err := client.CreateMulticastDomain(ctx, input)
		if err != nil {
			exitWithError("Failed to create multicast domain", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Multicast Domain created: %s (%s)\n", result.Domain.Name, result.Domain.ID)
	},
}

var thUpdateMulticastDomainCmd = &cobra.Command{
	Use:   "update-multicast-domain",
	Short: "Update a multicast domain",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("domain-id")
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")

		input := &transithub.UpdateMulticastDomainInput{
			Name:        name,
			Description: desc,
		}

		result, err := client.UpdateMulticastDomain(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update multicast domain", err)
		}

		fmt.Printf("Multicast Domain updated: %s\n", result.Domain.ID)
	},
}

var thDeleteMulticastDomainCmd = &cobra.Command{
	Use:   "delete-multicast-domain",
	Short: "Delete a multicast domain",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("domain-id")

		if err := client.DeleteMulticastDomain(ctx, id); err != nil {
			exitWithError("Failed to delete multicast domain", err)
		}

		fmt.Printf("Multicast Domain %s deleted\n", id)
	},
}
