package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/block"
	"github.com/spf13/cobra"
)

func init() {
	blockStorageCmd.AddCommand(bsDescribeVolumesCmd)
	blockStorageCmd.AddCommand(bsCreateVolumeCmd)
	blockStorageCmd.AddCommand(bsDeleteVolumeCmd)
	blockStorageCmd.AddCommand(bsUpdateVolumeCmd)
	blockStorageCmd.AddCommand(bsExtendVolumeCmd)
	blockStorageCmd.AddCommand(bsAttachVolumeCmd)
	blockStorageCmd.AddCommand(bsDetachVolumeCmd)
	blockStorageCmd.AddCommand(bsDescribeVolumeTypesCmd)

	bsCreateVolumeCmd.Flags().String("name", "", "Volume name")
	bsCreateVolumeCmd.Flags().Int("size", 10, "Volume size in GB (required)")
	bsCreateVolumeCmd.Flags().String("type", "", "Volume type")
	bsCreateVolumeCmd.Flags().String("availability-zone", "", "Availability zone")
	bsCreateVolumeCmd.Flags().String("snapshot-id", "", "Create from snapshot")
	bsCreateVolumeCmd.Flags().String("source-volume-id", "", "Create from existing volume")
	bsCreateVolumeCmd.Flags().String("description", "", "Volume description")
	bsCreateVolumeCmd.MarkFlagRequired("size")

	bsDescribeVolumesCmd.Flags().String("volume-id", "", "Volume ID")

	bsDeleteVolumeCmd.Flags().String("volume-id", "", "Volume ID (required)")
	bsDeleteVolumeCmd.MarkFlagRequired("volume-id")

	bsUpdateVolumeCmd.Flags().String("volume-id", "", "Volume ID (required)")
	bsUpdateVolumeCmd.Flags().String("name", "", "New volume name")
	bsUpdateVolumeCmd.Flags().String("description", "", "New description")
	bsUpdateVolumeCmd.MarkFlagRequired("volume-id")

	bsExtendVolumeCmd.Flags().String("volume-id", "", "Volume ID (required)")
	bsExtendVolumeCmd.Flags().Int("size", 0, "New size in GB (required)")
	bsExtendVolumeCmd.MarkFlagRequired("volume-id")
	bsExtendVolumeCmd.MarkFlagRequired("size")

	bsAttachVolumeCmd.Flags().String("volume-id", "", "Volume ID (required)")
	bsAttachVolumeCmd.Flags().String("server-id", "", "Server ID to attach to (required)")
	bsAttachVolumeCmd.Flags().String("device", "", "Device path (e.g., /dev/vdb)")
	bsAttachVolumeCmd.MarkFlagRequired("volume-id")
	bsAttachVolumeCmd.MarkFlagRequired("server-id")

	bsDetachVolumeCmd.Flags().String("volume-id", "", "Volume ID (required)")
	bsDetachVolumeCmd.MarkFlagRequired("volume-id")
}

var bsDescribeVolumesCmd = &cobra.Command{
	Use:     "describe-volumes",
	Aliases: []string{"list-volumes", "list", "ls"},
	Short:   "Describe volumes",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("volume-id")

		if id != "" {
			result, err := client.GetVolume(ctx, id)
			if err != nil {
				exitWithError("Failed to get volume", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			v := result.Volume
			fmt.Printf("ID:                %s\n", v.ID)
			fmt.Printf("Name:              %s\n", v.Name)
			fmt.Printf("Status:            %s\n", v.Status)
			fmt.Printf("Size:              %d GB\n", v.Size)
			fmt.Printf("Type:              %s\n", v.VolumeType)
			fmt.Printf("Bootable:          %s\n", v.Bootable)
			fmt.Printf("Encrypted:         %v\n", v.Encrypted)
			fmt.Printf("Availability Zone: %s\n", v.AvailabilityZone)
			fmt.Printf("Description:       %s\n", v.Description)
			fmt.Printf("Created:           %s\n", v.CreatedAt)

			if len(v.Attachments) > 0 {
				fmt.Printf("\nAttachments:\n")
				for _, a := range v.Attachments {
					fmt.Printf("  - Server: %s, Device: %s\n", a.ServerID, a.Device)
				}
			}
		} else {
			result, err := client.ListVolumes(ctx)
			if err != nil {
				exitWithError("Failed to list volumes", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tSTATUS\tSIZE\tTYPE\tBOOTABLE\tCREATED")
			for _, v := range result.Volumes {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d GB\t%s\t%s\t%s\n",
					v.ID, v.Name, v.Status, v.Size, v.VolumeType, v.Bootable, v.CreatedAt)
			}
			w.Flush()
		}
	},
}

var bsCreateVolumeCmd = &cobra.Command{
	Use:   "create-volume",
	Short: "Create a new volume",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		size, _ := cmd.Flags().GetInt("size")
		volumeType, _ := cmd.Flags().GetString("type")
		az, _ := cmd.Flags().GetString("availability-zone")
		snapshotID, _ := cmd.Flags().GetString("snapshot-id")
		sourceVolID, _ := cmd.Flags().GetString("source-volume-id")
		description, _ := cmd.Flags().GetString("description")

		if size < 10 || size > 1000 {
			exitWithError("size must be between 10 and 1000 GB", nil)
		}

		input := &block.CreateVolumeInput{
			Name:             name,
			Size:             size,
			VolumeType:       volumeType,
			AvailabilityZone: az,
			SnapshotID:       snapshotID,
			SourceVolID:      sourceVolID,
			Description:      description,
		}

		result, err := client.CreateVolume(ctx, input)
		if err != nil {
			exitWithError("Failed to create volume", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Volume created successfully!\n")
		fmt.Printf("ID:     %s\n", result.Volume.ID)
		fmt.Printf("Name:   %s\n", result.Volume.Name)
		fmt.Printf("Size:   %d GB\n", result.Volume.Size)
		fmt.Printf("Status: %s\n", result.Volume.Status)
	},
}

var bsDeleteVolumeCmd = &cobra.Command{
	Use:   "delete-volume",
	Short: "Delete a volume",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("volume-id")

		if err := client.DeleteVolume(ctx, id); err != nil {
			exitWithError("Failed to delete volume", err)
		}

		fmt.Printf("Volume %s deleted successfully\n", id)
	},
}

var bsUpdateVolumeCmd = &cobra.Command{
	Use:   "update-volume",
	Short: "Update volume name or description",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("volume-id")

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &block.UpdateVolumeInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateVolume(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update volume", err)
		}

		fmt.Printf("Volume %s updated successfully\n", result.Volume.ID)
	},
}

var bsExtendVolumeCmd = &cobra.Command{
	Use:   "extend-volume",
	Short: "Extend volume size",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("volume-id")

		newSize, _ := cmd.Flags().GetInt("size")

		if newSize <= 0 || newSize > 1000 {
			exitWithError("new size must be between 1 and 1000 GB", nil)
		}

		if err := client.ExtendVolume(ctx, id, newSize); err != nil {
			exitWithError("Failed to extend volume", err)
		}

		fmt.Printf("Volume %s extended to %d GB\n", id, newSize)
	},
}

var bsAttachVolumeCmd = &cobra.Command{
	Use:   "attach-volume",
	Short: "Attach volume to a server",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("volume-id")

		serverID, _ := cmd.Flags().GetString("server-id")
		device, _ := cmd.Flags().GetString("device")

		if err := client.AttachVolume(ctx, id, serverID, device); err != nil {
			exitWithError("Failed to attach volume", err)
		}

		fmt.Printf("Volume %s attached to server %s\n", id, serverID)
	},
}

var bsDetachVolumeCmd = &cobra.Command{
	Use:   "detach-volume",
	Short: "Detach volume from server",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("volume-id")

		if err := client.DetachVolume(ctx, id); err != nil {
			exitWithError("Failed to detach volume", err)
		}

		fmt.Printf("Volume %s detached\n", id)
	},
}

var bsDescribeVolumeTypesCmd = &cobra.Command{
	Use:   "describe-volume-types",
	Short: "List available volume types",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		result, err := client.ListVolumeTypes(ctx)
		if err != nil {
			exitWithError("Failed to list volume types", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION")
		for _, t := range result.VolumeTypes {
			fmt.Fprintf(w, "%s\t%s\t%s\n", t.ID, t.Name, t.Description)
		}
		w.Flush()
	},
}
