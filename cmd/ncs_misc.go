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
	ncsCmd.AddCommand(ncsDescribeTemplatesCmd)
	ncsCmd.AddCommand(ncsGetTemplateCmd)
	ncsCmd.AddCommand(ncsDescribeEventsCmd)

	ncsGetTemplateCmd.Flags().String("template-id", "", "Template ID (required)")
	ncsGetTemplateCmd.MarkFlagRequired("template-id")

	ncsDescribeEventsCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsDescribeEventsCmd.MarkFlagRequired("workload-id")
}

var ncsDescribeTemplatesCmd = &cobra.Command{
	Use:     "describe-templates",
	Aliases: []string{"list-templates", "templates"},
	Short:   "List workload templates",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		result, err := client.ListTemplates(ctx)
		if err != nil {
			exitWithError("Failed to list templates", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tVERSION\tTYPE\tCREATED")
		for _, t := range result.Templates {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				t.ID, t.Name, t.Version, t.Type, t.CreatedAt)
		}
		w.Flush()
	},
}

var ncsGetTemplateCmd = &cobra.Command{
	Use:     "describe-template",
	Aliases: []string{"get-template"},
	Short:   "Get template details",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		templateID, _ := cmd.Flags().GetString("template-id")

		result, err := client.GetTemplate(ctx, templateID)
		if err != nil {
			exitWithError("Failed to get template", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ID:          %s\n", result.ID)
		fmt.Printf("Name:        %s\n", result.Name)
		fmt.Printf("Version:     %s\n", result.Version)
		fmt.Printf("Type:        %s\n", result.Type)
		fmt.Printf("Description: %s\n", result.Description)
		fmt.Printf("Is Public:   %v\n", result.IsPublic)
		fmt.Printf("Created:     %s\n", result.CreatedAt)
	},
}

var ncsDescribeEventsCmd = &cobra.Command{
	Use:     "describe-events",
	Aliases: []string{"list-events", "events"},
	Short:   "List workload events",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")

		result, err := client.GetWorkloadEvents(ctx, workloadID)
		if err != nil {
			exitWithError("Failed to list events", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIME\tTYPE\tREASON\tMESSAGE")
		for _, e := range result.Events {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				e.LastTime, e.EventType, e.Reason, e.Message)
		}
		w.Flush()
	},
}
