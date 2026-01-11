package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncr"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/spf13/cobra"
)

var ncrCmd = &cobra.Command{
	Use:     "ncr",
	Aliases: []string{"registry"},
	Short:   "Manage NHN Container Registry (NCR)",
	Long:    `Manage container registries, images, and tags.`,
}

func init() {
	rootCmd.AddCommand(ncrCmd)

	ncrCmd.AddCommand(ncrListCmd)
	ncrCmd.AddCommand(ncrGetCmd)
	ncrCmd.AddCommand(ncrCreateCmd)
	ncrCmd.AddCommand(ncrDeleteCmd)

	ncrCmd.AddCommand(ncrImagesCmd)
	ncrCmd.AddCommand(ncrImageGetCmd)
	ncrCmd.AddCommand(ncrImageDeleteCmd)
	ncrCmd.AddCommand(ncrTagsCmd)
	ncrCmd.AddCommand(ncrTagDeleteCmd)
	ncrCmd.AddCommand(ncrScanCmd)
	ncrCmd.AddCommand(ncrScanResultCmd)

	ncrCmd.AddCommand(ncrWebhooksCmd)
	ncrCmd.AddCommand(ncrWebhookCreateCmd)
	ncrCmd.AddCommand(ncrWebhookDeleteCmd)

	ncrCreateCmd.Flags().String("name", "", "Registry name (required)")
	ncrCreateCmd.Flags().String("description", "", "Registry description")
	ncrCreateCmd.Flags().Bool("public", false, "Make registry public")
	ncrCreateCmd.MarkFlagRequired("name")

	ncrWebhookCreateCmd.Flags().String("registry-id", "", "Registry ID (required)")
	ncrWebhookCreateCmd.Flags().String("name", "", "Webhook name (required)")
	ncrWebhookCreateCmd.Flags().String("target-url", "", "Webhook target URL (required)")
	ncrWebhookCreateCmd.Flags().StringSlice("events", []string{"push"}, "Events to trigger (push, delete)")
	ncrWebhookCreateCmd.MarkFlagRequired("registry-id")
	ncrWebhookCreateCmd.MarkFlagRequired("name")
	ncrWebhookCreateCmd.MarkFlagRequired("target-url")
}

func getNCRClient() *ncr.Client {
	creds := credentials.NewStatic(getAccessKey(), getSecretKey())
	return ncr.NewClient(getRegion(), getNCRAppKey(), creds, nil, debug)
}

var ncrListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all container registries",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		result, err := client.ListRegistries(ctx)
		if err != nil {
			exitWithError("Failed to list registries", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tURI\tPUBLIC\tSTATUS\tCREATED")
		for _, r := range result.Registries {
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\t%s\n",
				r.ID, r.Name, r.URI, r.IsPublic, r.Status, r.CreatedAt)
		}
		w.Flush()
	},
}

var ncrGetCmd = &cobra.Command{
	Use:   "get [registry-id]",
	Short: "Get registry details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		result, err := client.GetRegistry(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get registry", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ID:      %s\n", result.ID)
		fmt.Printf("Name:    %s\n", result.Name)
		fmt.Printf("URI:     %s\n", result.URI)
		fmt.Printf("Public:  %v\n", result.IsPublic)
		fmt.Printf("Status:  %s\n", result.Status)
		fmt.Printf("Created: %s\n", result.CreatedAt)
	},
}

var ncrCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new container registry",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		isPublic, _ := cmd.Flags().GetBool("public")

		input := &ncr.CreateRegistryInput{
			Name:        name,
			Description: description,
			IsPublic:    isPublic,
		}

		result, err := client.CreateRegistry(ctx, input)
		if err != nil {
			exitWithError("Failed to create registry", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Registry created successfully!\n")
		fmt.Printf("ID:   %s\n", result.ID)
		fmt.Printf("Name: %s\n", result.Name)
		fmt.Printf("URI:  %s\n", result.URI)
	},
}

var ncrDeleteCmd = &cobra.Command{
	Use:   "delete [registry-id]",
	Short: "Delete a container registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		if err := client.DeleteRegistry(ctx, args[0]); err != nil {
			exitWithError("Failed to delete registry", err)
		}

		fmt.Printf("Registry %s deleted successfully\n", args[0])
	},
}

var ncrImagesCmd = &cobra.Command{
	Use:   "images [registry-id]",
	Short: "List images in a registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		result, err := client.ListImages(ctx, args[0])
		if err != nil {
			exitWithError("Failed to list images", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tSIZE\tPULL_COUNT\tCREATED")
		for _, img := range result.Images {
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
				img.Name, formatSize(img.Size), img.PullCount, img.CreatedAt)
		}
		w.Flush()
	},
}

var ncrImageGetCmd = &cobra.Command{
	Use:   "image-get [registry-id] [image-name]",
	Short: "Get image details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		result, err := client.GetImage(ctx, args[0], args[1])
		if err != nil {
			exitWithError("Failed to get image", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Name:       %s\n", result.Name)
		fmt.Printf("Size:       %s\n", formatSize(result.Size))
		fmt.Printf("Pull Count: %d\n", result.PullCount)
		fmt.Printf("Created:    %s\n", result.CreatedAt)
		if len(result.Tags) > 0 {
			fmt.Printf("Tags:       %v\n", result.Tags)
		}
	},
}

var ncrImageDeleteCmd = &cobra.Command{
	Use:   "image-delete [registry-id] [image-name]",
	Short: "Delete an image",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		if err := client.DeleteImage(ctx, args[0], args[1]); err != nil {
			exitWithError("Failed to delete image", err)
		}

		fmt.Printf("Image %s deleted successfully\n", args[1])
	},
}

var ncrTagsCmd = &cobra.Command{
	Use:   "tags [registry-id] [image-name]",
	Short: "List tags for an image",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		result, err := client.ListTags(ctx, args[0], args[1])
		if err != nil {
			exitWithError("Failed to list tags", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TAG\tDIGEST\tSIZE\tCREATED")
		for _, tag := range result.Tags {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				tag.Name, truncateDigest(tag.Digest), formatSize(tag.Size), tag.CreatedAt)
		}
		w.Flush()
	},
}

var ncrTagDeleteCmd = &cobra.Command{
	Use:   "tag-delete [registry-id] [image-name] [tag]",
	Short: "Delete a tag",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		if err := client.DeleteTag(ctx, args[0], args[1], args[2]); err != nil {
			exitWithError("Failed to delete tag", err)
		}

		fmt.Printf("Tag %s deleted successfully\n", args[2])
	},
}

var ncrScanCmd = &cobra.Command{
	Use:   "scan [registry-id] [image-name] [tag]",
	Short: "Scan an image for vulnerabilities",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		if err := client.ScanImage(ctx, args[0], args[1], args[2]); err != nil {
			exitWithError("Failed to initiate scan", err)
		}

		fmt.Printf("Scan initiated for %s:%s\n", args[1], args[2])
	},
}

var ncrScanResultCmd = &cobra.Command{
	Use:   "scan-result [registry-id] [image-name] [tag]",
	Short: "Get vulnerability scan results",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		result, err := client.GetImageScanResult(ctx, args[0], args[1], args[2])
		if err != nil {
			exitWithError("Failed to get scan result", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Scan Status:    %s\n", result.Status)
		fmt.Printf("Scan Completed: %s\n", result.ScanCompletedAt)
		if result.Summary != nil {
			fmt.Printf("Vulnerabilities:\n")
			fmt.Printf("  Critical: %d\n", result.Summary.Critical)
			fmt.Printf("  High:     %d\n", result.Summary.High)
			fmt.Printf("  Medium:   %d\n", result.Summary.Medium)
			fmt.Printf("  Low:      %d\n", result.Summary.Low)
			fmt.Printf("  Total:    %d\n", result.Summary.Total)
		}
	},
}

var ncrWebhooksCmd = &cobra.Command{
	Use:   "webhooks [registry-id]",
	Short: "List webhooks for a registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		result, err := client.ListWebhooks(ctx, args[0])
		if err != nil {
			exitWithError("Failed to list webhooks", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var ncrWebhookCreateCmd = &cobra.Command{
	Use:   "webhook-create",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Webhook created successfully!\n")
		fmt.Printf("ID:   %s\n", result.ID)
		fmt.Printf("Name: %s\n", result.Name)
	},
}

var ncrWebhookDeleteCmd = &cobra.Command{
	Use:   "webhook-delete [registry-id] [webhook-id]",
	Short: "Delete a webhook",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		if err := client.DeleteWebhook(ctx, args[0], args[1]); err != nil {
			exitWithError("Failed to delete webhook", err)
		}

		fmt.Printf("Webhook %s deleted successfully\n", args[1])
	},
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func truncateDigest(digest string) string {
	if len(digest) > 19 {
		return digest[:19] + "..."
	}
	return digest
}
