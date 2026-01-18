package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncr"
	"github.com/spf13/cobra"
)

func init() {
	ncrCmd.AddCommand(ncrDescribeWebhooksCmd)
	ncrCmd.AddCommand(ncrCreateWebhookCmd)
	ncrCmd.AddCommand(ncrDeleteWebhookCmd)

	ncrDescribeWebhooksCmd.Flags().String("registry-id", "", "Registry ID (required)")
	ncrDescribeWebhooksCmd.MarkFlagRequired("registry-id")

	ncrCreateWebhookCmd.Flags().String("registry-id", "", "Registry ID (required)")
	ncrCreateWebhookCmd.Flags().String("name", "", "Webhook name (required)")
	ncrCreateWebhookCmd.Flags().String("target-url", "", "Webhook target URL (required)")
	ncrCreateWebhookCmd.Flags().StringSlice("events", []string{"push"}, "Events to trigger (push, delete)")
	ncrCreateWebhookCmd.MarkFlagRequired("registry-id")
	ncrCreateWebhookCmd.MarkFlagRequired("name")
	ncrCreateWebhookCmd.MarkFlagRequired("target-url")

	ncrDeleteWebhookCmd.Flags().String("registry-id", "", "Registry ID (required)")
	ncrDeleteWebhookCmd.Flags().String("webhook-id", "", "Webhook ID (required)")
	ncrDeleteWebhookCmd.MarkFlagRequired("registry-id")
	ncrDeleteWebhookCmd.MarkFlagRequired("webhook-id")
}

var ncrDescribeWebhooksCmd = &cobra.Command{
	Use:   "describe-webhooks",
	Short: "List webhooks for a registry",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()
		registryID, _ := cmd.Flags().GetString("registry-id")

		result, err := client.ListWebhooks(ctx, registryID)
		if err != nil {
			exitWithError("Failed to list webhooks", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTARGET_URL\tENABLED")
		for _, wh := range result.Webhooks {
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n",
				wh.ID, wh.Name, wh.TargetURL, wh.Enabled)
		}
		w.Flush()
	},
}

var ncrCreateWebhookCmd = &cobra.Command{
	Use:   "create-webhook",
	Short: "Create a webhook for a registry",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		registryID, _ := cmd.Flags().GetString("registry-id")
		name, _ := cmd.Flags().GetString("name")
		targetURL, _ := cmd.Flags().GetString("target-url")
		events, _ := cmd.Flags().GetStringSlice("events")

		input := &ncr.CreateWebhookInput{
			Name:      name,
			TargetURL: targetURL,
			Events:    events,
		}

		result, err := client.CreateWebhook(ctx, registryID, input)
		if err != nil {
			exitWithError("Failed to create webhook", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Webhook created successfully!\n")
		fmt.Printf("ID:   %s\n", result.ID)
		fmt.Printf("Name: %s\n", result.Name)
	},
}

var ncrDeleteWebhookCmd = &cobra.Command{
	Use:   "delete-webhook",
	Short: "Delete a webhook",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()
		registryID, _ := cmd.Flags().GetString("registry-id")
		webhookID, _ := cmd.Flags().GetString("webhook-id")

		if err := client.DeleteWebhook(ctx, registryID, webhookID); err != nil {
			exitWithError("Failed to delete webhook", err)
		}

		fmt.Printf("Webhook %s deleted successfully\n", webhookID)
	},
}
