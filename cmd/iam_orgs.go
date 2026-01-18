package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	iamCmd.AddCommand(iamDescribeOrganizationsCmd)
	iamCmd.AddCommand(iamDescribeOrganizationCmd) // Optional alias or just merge
}

var iamDescribeOrganizationsCmd = &cobra.Command{
	Use:     "describe-organizations",
	Aliases: []string{"list-organizations", "orgs"},
	Short:   "List all organizations",
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()

		result, err := client.ListOrganizations(ctx)
		if err != nil {
			exitWithError("Failed to list organizations", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tCREATED")
		for _, o := range result.Organizations() {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				o.ID, o.Name, o.Status, o.CreatedAt)
		}
		w.Flush()
	},
}

var iamDescribeOrganizationCmd = &cobra.Command{
	Use:   "describe-organization [org-id]",
	Short: "Get organization details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()

		result, err := client.GetOrganization(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get organization", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		o := result.Organization
		fmt.Printf("ID:          %s\n", o.ID)
		fmt.Printf("Name:        %s\n", o.Name)
		fmt.Printf("Status:      %s\n", o.Status)
		fmt.Printf("Description: %s\n", o.Description)
		fmt.Printf("Created:     %s\n", o.CreatedAt)
	},
}
