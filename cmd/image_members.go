package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/image"
	"github.com/spf13/cobra"
)

func init() {
	imageCmd.AddCommand(imgDescribeMembersCmd)
	imageCmd.AddCommand(imgAddMemberCmd)
	imageCmd.AddCommand(imgUpdateMemberCmd)
	imageCmd.AddCommand(imgRemoveMemberCmd)

	imgDescribeMembersCmd.Flags().String("image-id", "", "Image ID (required)")
	imgDescribeMembersCmd.MarkFlagRequired("image-id")

	imgAddMemberCmd.Flags().String("image-id", "", "Image ID (required)")
	imgAddMemberCmd.Flags().String("member", "", "Member tenant ID to share with (required)")
	imgAddMemberCmd.MarkFlagRequired("image-id")
	imgAddMemberCmd.MarkFlagRequired("member")

	imgUpdateMemberCmd.Flags().String("image-id", "", "Image ID (required)")
	imgUpdateMemberCmd.Flags().String("member-id", "", "Member ID (same as tenant ID) (required)")
	imgUpdateMemberCmd.Flags().String("status", "", "Member status: accepted, pending, rejected (required)")
	imgUpdateMemberCmd.MarkFlagRequired("image-id")
	imgUpdateMemberCmd.MarkFlagRequired("member-id")
	imgUpdateMemberCmd.MarkFlagRequired("status")

	imgRemoveMemberCmd.Flags().String("image-id", "", "Image ID (required)")
	imgRemoveMemberCmd.Flags().String("member-id", "", "Member ID (required)")
	imgRemoveMemberCmd.MarkFlagRequired("image-id")
	imgRemoveMemberCmd.MarkFlagRequired("member-id")
}

var imgDescribeMembersCmd = &cobra.Command{
	Use:     "describe-members",
	Aliases: []string{"list-members"},
	Short:   "List members of an image",
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()
		imageID, _ := cmd.Flags().GetString("image-id")

		result, err := client.ListImageMembers(ctx, imageID)
		if err != nil {
			exitWithError("Failed to list image members", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		if len(result.Members) == 0 {
			fmt.Println("No members found for this image")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "MEMBER_ID\tSTATUS\tCREATED\tUPDATED")
		for _, m := range result.Members {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				m.MemberID, m.Status, m.CreatedAt.Format("2006-01-02"), m.UpdatedAt.Format("2006-01-02"))
		}
		w.Flush()
	},
}

var imgAddMemberCmd = &cobra.Command{
	Use:   "add-member",
	Short: "Add a member to share an image",
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()
		imageID, _ := cmd.Flags().GetString("image-id")
		memberID, _ := cmd.Flags().GetString("member")

		input := &image.CreateImageMemberInput{
			Member: memberID,
		}

		result, err := client.AddImageMember(ctx, imageID, input)
		if err != nil {
			exitWithError("Failed to add image member", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Member added successfully!\n")
		fmt.Printf("Image ID:  %s\n", result.ImageID)
		fmt.Printf("Member ID: %s\n", result.MemberID)
		fmt.Printf("Status:    %s\n", result.Status)
	},
}

var imgUpdateMemberCmd = &cobra.Command{
	Use:   "update-member",
	Short: "Update image member status (accept/reject shared image)",
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()
		imageID, _ := cmd.Flags().GetString("image-id")
		memberID, _ := cmd.Flags().GetString("member-id")
		status, _ := cmd.Flags().GetString("status")

		input := &image.UpdateImageMemberInput{
			Status: status,
		}

		result, err := client.UpdateImageMember(ctx, imageID, memberID, input)
		if err != nil {
			exitWithError("Failed to update image member", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Member status updated to '%s'\n", result.Status)
	},
}

var imgRemoveMemberCmd = &cobra.Command{
	Use:   "remove-member",
	Short: "Remove a member from an image",
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()
		imageID, _ := cmd.Flags().GetString("image-id")
		memberID, _ := cmd.Flags().GetString("member-id")

		if err := client.RemoveImageMember(ctx, imageID, memberID); err != nil {
			exitWithError("Failed to remove image member", err)
		}

		fmt.Printf("Member %s removed from image %s\n", memberID, imageID)
	},
}
