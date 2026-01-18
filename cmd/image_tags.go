package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	imageCmd.AddCommand(imgAddTagCmd)
	imageCmd.AddCommand(imgRemoveTagCmd)

	imgAddTagCmd.Flags().String("image-id", "", "Image ID (required)")
	imgAddTagCmd.Flags().String("tag", "", "Tag to add (required)")
	imgAddTagCmd.MarkFlagRequired("image-id")
	imgAddTagCmd.MarkFlagRequired("tag")

	imgRemoveTagCmd.Flags().String("image-id", "", "Image ID (required)")
	imgRemoveTagCmd.Flags().String("tag", "", "Tag to remove (required)")
	imgRemoveTagCmd.MarkFlagRequired("image-id")
	imgRemoveTagCmd.MarkFlagRequired("tag")
}

var imgAddTagCmd = &cobra.Command{
	Use:   "add-tag",
	Short: "Add a tag to an image",
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()
		imageID, _ := cmd.Flags().GetString("image-id")
		tag, _ := cmd.Flags().GetString("tag")

		if err := client.AddTag(ctx, imageID, tag); err != nil {
			exitWithError("Failed to add tag", err)
		}

		fmt.Printf("Tag '%s' added to image %s\n", tag, imageID)
	},
}

var imgRemoveTagCmd = &cobra.Command{
	Use:   "remove-tag",
	Short: "Remove a tag from an image",
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()
		imageID, _ := cmd.Flags().GetString("image-id")
		tag, _ := cmd.Flags().GetString("tag")

		if err := client.RemoveTag(ctx, imageID, tag); err != nil {
			exitWithError("Failed to remove tag", err)
		}

		fmt.Printf("Tag '%s' removed from image %s\n", tag, imageID)
	},
}
