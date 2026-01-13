package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/block"
	"github.com/spf13/cobra"
)

var blockStorageCmd = &cobra.Command{
	Use:     "block-storage",
	Aliases: []string{"volume", "bs"},
	Short:   "Manage Block Storage volumes and snapshots",
	Long:    `Manage block storage volumes, snapshots, and volume types.`,
}

func init() {
	rootCmd.AddCommand(blockStorageCmd)

	blockStorageCmd.AddCommand(bsListCmd)
	blockStorageCmd.AddCommand(bsGetCmd)
	blockStorageCmd.AddCommand(bsCreateCmd)
	blockStorageCmd.AddCommand(bsDeleteCmd)
	blockStorageCmd.AddCommand(bsUpdateCmd)
	blockStorageCmd.AddCommand(bsExtendCmd)
	blockStorageCmd.AddCommand(bsAttachCmd)
	blockStorageCmd.AddCommand(bsDetachCmd)
	blockStorageCmd.AddCommand(bsTypesCmd)

	blockStorageCmd.AddCommand(bsSnapshotListCmd)
	blockStorageCmd.AddCommand(bsSnapshotGetCmd)
	blockStorageCmd.AddCommand(bsSnapshotCreateCmd)
	blockStorageCmd.AddCommand(bsSnapshotDeleteCmd)

	bsCreateCmd.Flags().String("name", "", "Volume name")
	bsCreateCmd.Flags().Int("size", 10, "Volume size in GB (required)")
	bsCreateCmd.Flags().String("type", "", "Volume type")
	bsCreateCmd.Flags().String("availability-zone", "", "Availability zone")
	bsCreateCmd.Flags().String("snapshot-id", "", "Create from snapshot")
	bsCreateCmd.Flags().String("source-volume-id", "", "Create from existing volume")
	bsCreateCmd.Flags().String("description", "", "Volume description")
	bsCreateCmd.MarkFlagRequired("size")

	bsUpdateCmd.Flags().String("name", "", "New volume name")
	bsUpdateCmd.Flags().String("description", "", "New description")

	bsExtendCmd.Flags().Int("size", 0, "New size in GB (required)")
	bsExtendCmd.MarkFlagRequired("size")

	bsAttachCmd.Flags().String("server-id", "", "Server ID to attach to (required)")
	bsAttachCmd.Flags().String("device", "", "Device path (e.g., /dev/vdb)")
	bsAttachCmd.MarkFlagRequired("server-id")

	bsSnapshotCreateCmd.Flags().String("volume-id", "", "Volume ID to snapshot (required)")
	bsSnapshotCreateCmd.Flags().String("name", "", "Snapshot name")
	bsSnapshotCreateCmd.Flags().String("description", "", "Snapshot description")
	bsSnapshotCreateCmd.Flags().Bool("force", false, "Force snapshot of in-use volume")
	bsSnapshotCreateCmd.MarkFlagRequired("volume-id")
}

func getBlockStorageClient() *block.Client {
	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	return block.NewClient(getRegion(), creds, nil, debug)
}

var bsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all volumes",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		result, err := client.ListVolumes(ctx)
		if err != nil {
			exitWithError("Failed to list volumes", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tSIZE\tTYPE\tBOOTABLE\tCREATED")
		for _, v := range result.Volumes {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d GB\t%s\t%s\t%s\n",
				v.ID, v.Name, v.Status, v.Size, v.VolumeType, v.Bootable, v.CreatedAt)
		}
		w.Flush()
	},
}

var bsGetCmd = &cobra.Command{
	Use:   "get [volume-id]",
	Short: "Get volume details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		result, err := client.GetVolume(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get volume", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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
	},
}

var bsCreateCmd = &cobra.Command{
	Use:   "create",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Volume created successfully!\n")
		fmt.Printf("ID:     %s\n", result.Volume.ID)
		fmt.Printf("Name:   %s\n", result.Volume.Name)
		fmt.Printf("Size:   %d GB\n", result.Volume.Size)
		fmt.Printf("Status: %s\n", result.Volume.Status)
	},
}

var bsDeleteCmd = &cobra.Command{
	Use:   "delete [volume-id]",
	Short: "Delete a volume",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		if err := client.DeleteVolume(ctx, args[0]); err != nil {
			exitWithError("Failed to delete volume", err)
		}

		fmt.Printf("Volume %s deleted successfully\n", args[0])
	},
}

var bsUpdateCmd = &cobra.Command{
	Use:   "update [volume-id]",
	Short: "Update volume name or description",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &block.UpdateVolumeInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateVolume(ctx, args[0], input)
		if err != nil {
			exitWithError("Failed to update volume", err)
		}

		fmt.Printf("Volume %s updated successfully\n", result.Volume.ID)
	},
}

var bsExtendCmd = &cobra.Command{
	Use:   "extend [volume-id]",
	Short: "Extend volume size",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		newSize, _ := cmd.Flags().GetInt("size")

		if newSize <= 0 || newSize > 1000 {
			exitWithError("new size must be between 1 and 1000 GB", nil)
		}

		if err := client.ExtendVolume(ctx, args[0], newSize); err != nil {
			exitWithError("Failed to extend volume", err)
		}

		fmt.Printf("Volume %s extended to %d GB\n", args[0], newSize)
	},
}

var bsAttachCmd = &cobra.Command{
	Use:   "attach [volume-id]",
	Short: "Attach volume to a server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		serverID, _ := cmd.Flags().GetString("server-id")
		device, _ := cmd.Flags().GetString("device")

		if err := client.AttachVolume(ctx, args[0], serverID, device); err != nil {
			exitWithError("Failed to attach volume", err)
		}

		fmt.Printf("Volume %s attached to server %s\n", args[0], serverID)
	},
}

var bsDetachCmd = &cobra.Command{
	Use:   "detach [volume-id]",
	Short: "Detach volume from server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		if err := client.DetachVolume(ctx, args[0]); err != nil {
			exitWithError("Failed to detach volume", err)
		}

		fmt.Printf("Volume %s detached\n", args[0])
	},
}

var bsTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List available volume types",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		result, err := client.ListVolumeTypes(ctx)
		if err != nil {
			exitWithError("Failed to list volume types", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var bsSnapshotListCmd = &cobra.Command{
	Use:   "snapshot-list",
	Short: "List all snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		result, err := client.ListSnapshots(ctx)
		if err != nil {
			exitWithError("Failed to list snapshots", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tSIZE\tVOLUME_ID\tCREATED")
		for _, s := range result.Snapshots {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d GB\t%s\t%s\n",
				s.ID, s.Name, s.Status, s.Size, s.VolumeID, s.CreatedAt)
		}
		w.Flush()
	},
}

var bsSnapshotGetCmd = &cobra.Command{
	Use:   "snapshot-get [snapshot-id]",
	Short: "Get snapshot details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		result, err := client.GetSnapshot(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get snapshot", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		s := result.Snapshot
		fmt.Printf("ID:          %s\n", s.ID)
		fmt.Printf("Name:        %s\n", s.Name)
		fmt.Printf("Status:      %s\n", s.Status)
		fmt.Printf("Size:        %d GB\n", s.Size)
		fmt.Printf("Volume ID:   %s\n", s.VolumeID)
		fmt.Printf("Description: %s\n", s.Description)
		fmt.Printf("Created:     %s\n", s.CreatedAt)
	},
}

var bsSnapshotCreateCmd = &cobra.Command{
	Use:   "snapshot-create",
	Short: "Create a snapshot from a volume",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		volumeID, _ := cmd.Flags().GetString("volume-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		force, _ := cmd.Flags().GetBool("force")

		input := &block.CreateSnapshotInput{
			VolumeID:    volumeID,
			Name:        name,
			Description: description,
			Force:       force,
		}

		result, err := client.CreateSnapshot(ctx, input)
		if err != nil {
			exitWithError("Failed to create snapshot", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Snapshot created successfully!\n")
		fmt.Printf("ID:     %s\n", result.Snapshot.ID)
		fmt.Printf("Name:   %s\n", result.Snapshot.Name)
		fmt.Printf("Status: %s\n", result.Snapshot.Status)
	},
}

var bsSnapshotDeleteCmd = &cobra.Command{
	Use:   "snapshot-delete [snapshot-id]",
	Short: "Delete a snapshot",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()

		if err := client.DeleteSnapshot(ctx, args[0]); err != nil {
			exitWithError("Failed to delete snapshot", err)
		}

		fmt.Printf("Snapshot %s deleted successfully\n", args[0])
	},
}
