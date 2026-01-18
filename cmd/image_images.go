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
	imageCmd.AddCommand(imgDescribeCmd)
	imageCmd.AddCommand(imgGetCmd)
	imageCmd.AddCommand(imgCreateCmd)
	imageCmd.AddCommand(imgDeleteCmd)

	// List flags
	imgDescribeCmd.Flags().String("name", "", "Filter by image name")
	imgDescribeCmd.Flags().String("status", "", "Filter by status (active, queued, saving, etc.)")
	imgDescribeCmd.Flags().String("visibility", "", "Filter by visibility (public, private, shared, community)")
	imgDescribeCmd.Flags().String("os-type", "", "Filter by OS type (linux, windows)")
	imgDescribeCmd.Flags().String("os-distro", "", "Filter by OS distribution (ubuntu, centos, etc.)")
	imgDescribeCmd.Flags().Int("limit", 0, "Limit number of results")

	// Get flags
	imgGetCmd.Flags().String("image-id", "", "Image ID (required)")
	imgGetCmd.MarkFlagRequired("image-id")

	// Create flags
	imgCreateCmd.Flags().String("name", "", "Image name (required)")
	imgCreateCmd.Flags().String("disk-format", "raw", "Disk format (raw, qcow2, vhd, vmdk, iso)")
	imgCreateCmd.Flags().String("container-format", "bare", "Container format (bare, ovf, ova, docker)")
	imgCreateCmd.Flags().Int("min-disk", 0, "Minimum disk size in GB")
	imgCreateCmd.Flags().Int("min-ram", 0, "Minimum RAM in MB")
	imgCreateCmd.Flags().Bool("protected", false, "Prevent image from being deleted")
	imgCreateCmd.Flags().String("os-distro", "", "OS distribution")
	imgCreateCmd.Flags().String("os-version", "", "OS version")
	imgCreateCmd.Flags().String("os-type", "", "OS type (linux, windows)")
	imgCreateCmd.MarkFlagRequired("name")

	// Delete flags
	imgDeleteCmd.Flags().String("image-id", "", "Image ID (required)")
	imgDeleteCmd.MarkFlagRequired("image-id")
}

var imgDescribeCmd = &cobra.Command{
	Use:     "describe-images",
	Aliases: []string{"list-images", "list"},
	Short:   "List all images",
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		status, _ := cmd.Flags().GetString("status")
		visibility, _ := cmd.Flags().GetString("visibility")
		osType, _ := cmd.Flags().GetString("os-type")
		osDistro, _ := cmd.Flags().GetString("os-distro")
		limit, _ := cmd.Flags().GetInt("limit")

		input := &image.ListImagesInput{
			Name:       name,
			Status:     status,
			Visibility: visibility,
			OSType:     osType,
			OSDistro:   osDistro,
			Limit:      limit,
		}

		result, err := client.ListImages(ctx, input)
		if err != nil {
			exitWithError("Failed to list images", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tVISIBILITY\tSIZE (MB)\tOS\tCREATED")
		for _, img := range result.Images {
			sizeMB := img.Size / (1024 * 1024)
			osInfo := img.OSDistro
			if osInfo == "" {
				osInfo = img.OSType
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
				img.ID, img.Name, img.Status, img.Visibility, sizeMB, osInfo, img.CreatedAt.Format("2006-01-02"))
		}
		w.Flush()
	},
}

var imgGetCmd = &cobra.Command{
	Use:     "describe-image",
	Aliases: []string{"get-image", "get"},
	Short:   "Get image details",
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("image-id")

		result, err := client.GetImage(ctx, id)
		if err != nil {
			exitWithError("Failed to get image", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ID:               %s\n", result.ID)
		fmt.Printf("Name:             %s\n", result.Name)
		fmt.Printf("Status:           %s\n", result.Status)
		fmt.Printf("Visibility:       %s\n", result.Visibility)
		fmt.Printf("Protected:        %v\n", result.Protected)
		fmt.Printf("Owner:            %s\n", result.Owner)
		fmt.Printf("Size:             %d bytes (%.2f MB)\n", result.Size, float64(result.Size)/(1024*1024))
		fmt.Printf("Min Disk:         %d GB\n", result.MinDisk)
		fmt.Printf("Min RAM:          %d MB\n", result.MinRAM)
		fmt.Printf("Disk Format:      %s\n", result.DiskFormat)
		fmt.Printf("Container Format: %s\n", result.ContainerFormat)
		if result.OSType != "" {
			fmt.Printf("OS Type:          %s\n", result.OSType)
		}
		if result.OSDistro != "" {
			fmt.Printf("OS Distro:        %s\n", result.OSDistro)
		}
		if result.OSVersion != "" {
			fmt.Printf("OS Version:       %s\n", result.OSVersion)
		}
		if result.Checksum != "" {
			fmt.Printf("Checksum:         %s\n", result.Checksum)
		}
		if len(result.Tags) > 0 {
			fmt.Printf("Tags:             %v\n", result.Tags)
		}
		fmt.Printf("Created:          %s\n", result.CreatedAt)
		fmt.Printf("Updated:          %s\n", result.UpdatedAt)
	},
}

var imgCreateCmd = &cobra.Command{
	Use:   "create-image",
	Short: "Create a new image (metadata only)",
	Long:  `Create image metadata. To upload image data, use the Glance API directly.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		diskFormat, _ := cmd.Flags().GetString("disk-format")
		containerFormat, _ := cmd.Flags().GetString("container-format")
		minDisk, _ := cmd.Flags().GetInt("min-disk")
		minRAM, _ := cmd.Flags().GetInt("min-ram")
		protected, _ := cmd.Flags().GetBool("protected")
		osDistro, _ := cmd.Flags().GetString("os-distro")
		osVersion, _ := cmd.Flags().GetString("os-version")
		osType, _ := cmd.Flags().GetString("os-type")

		input := &image.CreateImageInput{
			Name:            name,
			DiskFormat:      diskFormat,
			ContainerFormat: containerFormat,
			MinDisk:         minDisk,
			MinRAM:          minRAM,
			Protected:       protected,
			OSDistro:        osDistro,
			OSVersion:       osVersion,
			OSType:          osType,
		}

		result, err := client.CreateImage(ctx, input)
		if err != nil {
			exitWithError("Failed to create image", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Image created successfully!\n")
		fmt.Printf("ID:     %s\n", result.ID)
		fmt.Printf("Name:   %s\n", result.Name)
		fmt.Printf("Status: %s\n", result.Status)
	},
}

var imgDeleteCmd = &cobra.Command{
	Use:   "delete-image",
	Short: "Delete an image",
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("image-id")

		if err := client.DeleteImage(ctx, id); err != nil {
			exitWithError("Failed to delete image", err)
		}

		fmt.Printf("Image %s deleted successfully\n", id)
	},
}
