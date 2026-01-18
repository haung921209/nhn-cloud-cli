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
	transitHubCmd.AddCommand(thDescribeAttachmentsCmd)
	transitHubCmd.AddCommand(thCreateAttachmentCmd)
	transitHubCmd.AddCommand(thUpdateAttachmentCmd)
	transitHubCmd.AddCommand(thDeleteAttachmentCmd)

	thCreateAttachmentCmd.Flags().String("name", "", "Attachment name (required)")
	thCreateAttachmentCmd.Flags().String("description", "", "Attachment description")
	thCreateAttachmentCmd.Flags().String("transit-hub-id", "", "Transit Hub ID (required)")
	thCreateAttachmentCmd.Flags().String("resource-type", "VPC", "Resource type (VPC, VPN)")
	thCreateAttachmentCmd.Flags().String("resource-id", "", "Resource ID (required)")
	thCreateAttachmentCmd.MarkFlagRequired("name")
	thCreateAttachmentCmd.MarkFlagRequired("transit-hub-id")
	thCreateAttachmentCmd.MarkFlagRequired("resource-id")

	thUpdateAttachmentCmd.Flags().String("attachment-id", "", "Attachment ID (required)")
	thUpdateAttachmentCmd.Flags().String("name", "", "Attachment name")
	thUpdateAttachmentCmd.Flags().String("description", "", "Attachment description")
	thUpdateAttachmentCmd.MarkFlagRequired("attachment-id")

	thDeleteAttachmentCmd.Flags().String("attachment-id", "", "Attachment ID (required)")
	thDeleteAttachmentCmd.MarkFlagRequired("attachment-id")
}

var thDescribeAttachmentsCmd = &cobra.Command{
	Use:     "describe-attachments",
	Aliases: []string{"list-attachments"},
	Short:   "List attachments",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()

		result, err := client.ListAttachments(ctx)
		if err != nil {
			exitWithError("Failed to list attachments", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTRANSIT_HUB\tRESOURCE_TYPE\tRESOURCE_ID\tSTATUS")
		for _, a := range result.Attachments {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				a.ID, a.Name, a.TransitHubID, a.ResourceType, a.ResourceID, a.Status)
		}
		w.Flush()
	},
}

var thCreateAttachmentCmd = &cobra.Command{
	Use:   "create-attachment",
	Short: "Create a new attachment",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		thID, _ := cmd.Flags().GetString("transit-hub-id")
		rType, _ := cmd.Flags().GetString("resource-type")
		rID, _ := cmd.Flags().GetString("resource-id")

		input := &transithub.CreateAttachmentInput{
			Name:         name,
			Description:  desc,
			TransitHubID: thID,
			ResourceType: rType,
			ResourceID:   rID,
		}

		result, err := client.CreateAttachment(ctx, input)
		if err != nil {
			exitWithError("Failed to create attachment", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Attachment created: %s (%s)\n", result.Attachment.Name, result.Attachment.ID)
	},
}

var thUpdateAttachmentCmd = &cobra.Command{
	Use:   "update-attachment",
	Short: "Update an attachment",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("attachment-id")
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")

		input := &transithub.UpdateAttachmentInput{
			Name:        name,
			Description: desc,
		}

		result, err := client.UpdateAttachment(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update attachment", err)
		}

		fmt.Printf("Attachment updated: %s\n", result.Attachment.ID)
	},
}

var thDeleteAttachmentCmd = &cobra.Command{
	Use:   "delete-attachment",
	Short: "Delete an attachment",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("attachment-id")

		if err := client.DeleteAttachment(ctx, id); err != nil {
			exitWithError("Failed to delete attachment", err)
		}

		fmt.Printf("Attachment %s deleted\n", id)
	},
}
