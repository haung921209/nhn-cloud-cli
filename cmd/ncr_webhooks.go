package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
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
	ncrCreateWebhookCmd.Flags().StringSlice("events", []string{"PUSH_ARTIFACT"}, "Event types: PUSH_ARTIFACT, PULL_ARTIFACT, DELETE_ARTIFACT")
	ncrCreateWebhookCmd.Flags().String("target-type", "http", "Notification type: http or slack")
	ncrCreateWebhookCmd.Flags().Bool("enabled", true, "Enable the webhook")
	ncrCreateWebhookCmd.MarkFlagRequired("registry-id")
	ncrCreateWebhookCmd.MarkFlagRequired("name")
	ncrCreateWebhookCmd.MarkFlagRequired("target-url")

	ncrDeleteWebhookCmd.Flags().String("registry-id", "", "Registry ID (required)")
	ncrDeleteWebhookCmd.Flags().String("webhook-id", "", "Webhook policy ID from describe-webhooks (required)")
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

		if structuredOutput() {
			printOutput(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTARGETS\tEVENTS\tENABLED")
		for _, wh := range result.Webhooks {
			var targets []string
			for _, t := range wh.Targets {
				targets = append(targets, t.Address)
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%v\n",
				wh.ID, wh.Name, strings.Join(targets, ","), strings.Join(wh.EventTypes, ","), wh.Enabled)
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
		targetType, _ := cmd.Flags().GetString("target-type")
		enabled, _ := cmd.Flags().GetBool("enabled")

		input := &ncr.CreateWebhookInput{
			Name:       name,
			Enabled:    enabled,
			EventTypes: events,
			Targets:    []ncr.WebhookTarget{{Address: targetURL, Type: targetType}},
		}

		if _, err := client.CreateWebhook(ctx, registryID, input); err != nil {
			exitWithError("Failed to create webhook", err)
		}
		// The create response is the common header only; find the policy by name.
		list, err := client.ListWebhooks(ctx, registryID)
		if err != nil {
			exitWithError("Webhook created, but reading it back failed", err)
		}
		for _, wh := range list.Webhooks {
			if wh.Name == name {
				printWriteResult(wh, "Webhook created\nID:   %d\nName: %s\n", wh.ID, wh.Name)
				return
			}
		}
		exitWithError("Webhook created, but no policy named "+name+" was listed", nil)
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
