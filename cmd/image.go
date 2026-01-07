package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/image"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"images", "img"},
	Short:   "Manage Glance images",
	Long:    `Manage images including list, get, create, delete, and tag management.`,
}

var imageTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage image tags",
}

var imageMemberCmd = &cobra.Command{
	Use:   "member",
	Short: "Manage image members (sharing)",
}

func init() {
	rootCmd.AddCommand(imageCmd)

	// Main image commands
	imageCmd.AddCommand(imageListCmd)
	imageCmd.AddCommand(imageGetCmd)
	imageCmd.AddCommand(imageCreateCmd)
	imageCmd.AddCommand(imageDeleteCmd)

	// Tag subcommands
	imageCmd.AddCommand(imageTagCmd)
	imageTagCmd.AddCommand(imageTagAddCmd)
	imageTagCmd.AddCommand(imageTagRemoveCmd)

	// Member subcommands
	imageCmd.AddCommand(imageMemberCmd)
	imageMemberCmd.AddCommand(imageMemberListCmd)
	imageMemberCmd.AddCommand(imageMemberAddCmd)
	imageMemberCmd.AddCommand(imageMemberUpdateCmd)
	imageMemberCmd.AddCommand(imageMemberRemoveCmd)

	// List flags
	imageListCmd.Flags().String("name", "", "Filter by image name")
	imageListCmd.Flags().String("status", "", "Filter by status (active, queued, saving, etc.)")
	imageListCmd.Flags().String("visibility", "", "Filter by visibility (public, private, shared, community)")
	imageListCmd.Flags().String("os-type", "", "Filter by OS type (linux, windows)")
	imageListCmd.Flags().String("os-distro", "", "Filter by OS distribution (ubuntu, centos, etc.)")
	imageListCmd.Flags().Int("limit", 0, "Limit number of results")

	// Create flags
	imageCreateCmd.Flags().String("name", "", "Image name (required)")
	imageCreateCmd.Flags().String("disk-format", "raw", "Disk format (raw, qcow2, vhd, vmdk, iso)")
	imageCreateCmd.Flags().String("container-format", "bare", "Container format (bare, ovf, ova, docker)")
	imageCreateCmd.Flags().Int("min-disk", 0, "Minimum disk size in GB")
	imageCreateCmd.Flags().Int("min-ram", 0, "Minimum RAM in MB")
	imageCreateCmd.Flags().Bool("protected", false, "Prevent image from being deleted")
	imageCreateCmd.Flags().String("os-distro", "", "OS distribution")
	imageCreateCmd.Flags().String("os-version", "", "OS version")
	imageCreateCmd.Flags().String("os-type", "", "OS type (linux, windows)")
	imageCreateCmd.MarkFlagRequired("name")

	// Member add flags
	imageMemberAddCmd.Flags().String("member", "", "Member tenant ID (required)")
	imageMemberAddCmd.MarkFlagRequired("member")

	// Member update flags
	imageMemberUpdateCmd.Flags().String("status", "", "Member status: accepted, pending, rejected (required)")
	imageMemberUpdateCmd.MarkFlagRequired("status")
}

func getImageClient() *image.Client {
	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	return image.NewClient(getRegion(), creds, nil, debug)
}

// ============ Image Commands ============

var imageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all images",
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

var imageGetCmd = &cobra.Command{
	Use:   "get [image-id]",
	Short: "Get image details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()

		result, err := client.GetImage(ctx, args[0])
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

var imageCreateCmd = &cobra.Command{
	Use:   "create",
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

var imageDeleteCmd = &cobra.Command{
	Use:   "delete [image-id]",
	Short: "Delete an image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()

		if err := client.DeleteImage(ctx, args[0]); err != nil {
			exitWithError("Failed to delete image", err)
		}

		fmt.Printf("Image %s deleted successfully\n", args[0])
	},
}

// ============ Tag Commands ============

var imageTagAddCmd = &cobra.Command{
	Use:   "add [image-id] [tag]",
	Short: "Add a tag to an image",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()

		if err := client.AddTag(ctx, args[0], args[1]); err != nil {
			exitWithError("Failed to add tag", err)
		}

		fmt.Printf("Tag '%s' added to image %s\n", args[1], args[0])
	},
}

var imageTagRemoveCmd = &cobra.Command{
	Use:   "remove [image-id] [tag]",
	Short: "Remove a tag from an image",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()

		if err := client.RemoveTag(ctx, args[0], args[1]); err != nil {
			exitWithError("Failed to remove tag", err)
		}

		fmt.Printf("Tag '%s' removed from image %s\n", args[1], args[0])
	},
}

// ============ Member Commands ============

var imageMemberListCmd = &cobra.Command{
	Use:   "list [image-id]",
	Short: "List members of an image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()

		result, err := client.ListImageMembers(ctx, args[0])
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

var imageMemberAddCmd = &cobra.Command{
	Use:   "add [image-id]",
	Short: "Add a member to share an image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()

		memberID, _ := cmd.Flags().GetString("member")

		input := &image.CreateImageMemberInput{
			Member: memberID,
		}

		result, err := client.AddImageMember(ctx, args[0], input)
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

var imageMemberUpdateCmd = &cobra.Command{
	Use:   "update [image-id] [member-id]",
	Short: "Update image member status (accept/reject shared image)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()

		status, _ := cmd.Flags().GetString("status")

		input := &image.UpdateImageMemberInput{
			Status: status,
		}

		result, err := client.UpdateImageMember(ctx, args[0], args[1], input)
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

var imageMemberRemoveCmd = &cobra.Command{
	Use:   "remove [image-id] [member-id]",
	Short: "Remove a member from an image",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getImageClient()
		ctx := context.Background()

		if err := client.RemoveImageMember(ctx, args[0], args[1]); err != nil {
			exitWithError("Failed to remove image member", err)
		}

		fmt.Printf("Member %s removed from image %s\n", args[1], args[0])
	},
}
