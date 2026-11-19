package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	ncrCmd.AddCommand(ncrDescribeRepositoriesCmd)
	ncrCmd.AddCommand(ncrDescribeImagesCmd)
	ncrCmd.AddCommand(ncrDeleteRepositoryCmd)
	ncrCmd.AddCommand(ncrDeleteImageCmd)

	ncrDescribeRepositoriesCmd.Flags().String("registry-id", "", "Registry ID (required)")
	ncrDescribeRepositoriesCmd.Flags().String("repository-name", "", "Image name without the registry prefix, for details")
	ncrDescribeRepositoriesCmd.MarkFlagRequired("registry-id")

	ncrDescribeImagesCmd.Flags().String("registry-id", "", "Registry ID (required)")
	ncrDescribeImagesCmd.Flags().String("repository-name", "", "Image name without the registry prefix (required)")
	ncrDescribeImagesCmd.MarkFlagRequired("registry-id")
	ncrDescribeImagesCmd.MarkFlagRequired("repository-name")

	ncrDeleteRepositoryCmd.Flags().String("registry-id", "", "Registry ID (required)")
	ncrDeleteRepositoryCmd.Flags().String("repository-name", "", "Repository name (required)")
	ncrDeleteRepositoryCmd.MarkFlagRequired("registry-id")
	ncrDeleteRepositoryCmd.MarkFlagRequired("repository-name")

	ncrDeleteImageCmd.Flags().String("registry-id", "", "Registry ID (required)")
	ncrDeleteImageCmd.Flags().String("repository-name", "", "Repository name (required)")
	ncrDeleteImageCmd.Flags().String("image-tag", "", "Image tag (required)")
	ncrDeleteImageCmd.MarkFlagRequired("registry-id")
	ncrDeleteImageCmd.MarkFlagRequired("repository-name")
	ncrDeleteImageCmd.MarkFlagRequired("image-tag")
}

var ncrDescribeRepositoriesCmd = &cobra.Command{
	Use:   "describe-repositories",
	Short: "List repositories (images) in a registry",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()
		registryID, _ := cmd.Flags().GetString("registry-id")
		repoName, _ := cmd.Flags().GetString("repository-name")

		if repoName != "" {
			// Get Image Details (Repository details)
			result, err := client.GetImage(ctx, registryID, repoName)
			if err != nil {
				exitWithError("Failed to get repository", err)
			}
			if structuredOutput() {
				printOutput(result)
				return
			}
			img := result.Image
			fmt.Printf("Name:       %s\n", img.Name)
			fmt.Printf("Artifacts:  %d\n", img.ArtifactCount)
			fmt.Printf("Pull Count: %d\n", img.PullCount)
			fmt.Printf("Created:    %s\n", img.CreatedAt)
		} else {
			// List Repositories (List Images)
			result, err := client.ListImages(ctx, registryID)
			if err != nil {
				exitWithError("Failed to list repositories", err)
			}
			if structuredOutput() {
				printOutput(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tARTIFACTS\tPULL_COUNT\tCREATED")
			for _, img := range result.Images {
				fmt.Fprintf(w, "%s\t%d\t%d\t%s\n",
					img.Name, img.ArtifactCount, img.PullCount, img.CreatedAt)
			}
			w.Flush()
		}
	},
}

var ncrDescribeImagesCmd = &cobra.Command{
	Use:   "describe-images",
	Short: "List images (tags) in a repository",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()
		registryID, _ := cmd.Flags().GetString("registry-id")
		repoName, _ := cmd.Flags().GetString("repository-name")

		result, err := client.ListArtifacts(ctx, registryID, repoName)
		if err != nil {
			exitWithError("Failed to list images", err)
		}

		if structuredOutput() {
			printOutput(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TAG\tDIGEST\tSIZE\tPUSHED")
		for _, a := range result.Artifacts {
			for _, tag := range a.Tags {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					tag.Name, truncateDigest(a.Digest), formatSize(a.Size), tag.PushTime)
			}
		}
		w.Flush()
	},
}

var ncrDeleteRepositoryCmd = &cobra.Command{
	Use:   "delete-repository",
	Short: "Delete a repository (image)",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()
		registryID, _ := cmd.Flags().GetString("registry-id")
		repoName, _ := cmd.Flags().GetString("repository-name")

		if err := client.DeleteImage(ctx, registryID, repoName); err != nil {
			exitWithError("Failed to delete repository", err)
		}

		fmt.Printf("Repository %s deleted successfully\n", repoName)
	},
}

var ncrDeleteImageCmd = &cobra.Command{
	Use:   "delete-image",
	Short: "Delete a specific image tag",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()
		registryID, _ := cmd.Flags().GetString("registry-id")
		repoName, _ := cmd.Flags().GetString("repository-name")
		tag, _ := cmd.Flags().GetString("image-tag")

		if err := client.DeleteTag(ctx, registryID, repoName, tag); err != nil {
			exitWithError("Failed to delete image tag", err)
		}

		fmt.Printf("Image tag %s deleted successfully\n", tag)
	},
}

func truncateDigest(digest string) string {
	if len(digest) > 19 {
		return digest[:19] + "..."
	}
	return digest
}
