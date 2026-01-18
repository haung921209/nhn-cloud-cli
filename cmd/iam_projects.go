package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	iamCmd.AddCommand(iamDescribeProjectsCmd)
	iamCmd.AddCommand(iamDescribeProjectCmd)

	iamDescribeProjectsCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamDescribeProjectsCmd.MarkFlagRequired("org-id")

	iamDescribeProjectCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamDescribeProjectCmd.MarkFlagRequired("org-id")
}

var iamDescribeProjectsCmd = &cobra.Command{
	Use:   "describe-projects",
	Short: "List projects in an organization",
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()
		orgID, _ := cmd.Flags().GetString("org-id")

		result, err := client.ListProjects(ctx, orgID)
		if err != nil {
			exitWithError("Failed to list projects", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tCREATED")
		for _, p := range result.Projects {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				p.ID, p.Name, p.Status, p.CreatedAt)
		}
		w.Flush()
	},
}

var iamDescribeProjectCmd = &cobra.Command{
	Use:   "describe-project [project-id]",
	Short: "Get project details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()
		orgID, _ := cmd.Flags().GetString("org-id")

		result, err := client.GetProject(ctx, orgID, args[0])
		if err != nil {
			exitWithError("Failed to get project", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		p := result.Project
		fmt.Printf("ID:          %s\n", p.ID)
		fmt.Printf("Name:        %s\n", p.Name)
		fmt.Printf("Status:      %s\n", p.Status)
		fmt.Printf("Description: %s\n", p.Description)
		fmt.Printf("Org ID:      %s\n", p.OrganizationID)
		fmt.Printf("Created:     %s\n", p.CreatedAt)
	},
}
