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
	blockStorageCmd.AddCommand(bsDescribeSnapshotsCmd)
	blockStorageCmd.AddCommand(bsCreateSnapshotCmd)
	blockStorageCmd.AddCommand(bsDeleteSnapshotCmd)

	bsDescribeSnapshotsCmd.Flags().String("snapshot-id", "", "Snapshot ID")

	bsCreateSnapshotCmd.Flags().String("volume-id", "", "Volume ID to snapshot (required)")
	bsCreateSnapshotCmd.Flags().String("name", "", "Snapshot name")
	bsCreateSnapshotCmd.Flags().String("description", "", "Snapshot description")
	bsCreateSnapshotCmd.Flags().Bool("force", false, "Force snapshot of in-use volume")
	bsCreateSnapshotCmd.MarkFlagRequired("volume-id")

	bsDeleteSnapshotCmd.Flags().String("snapshot-id", "", "Snapshot ID (required)")
	bsDeleteSnapshotCmd.MarkFlagRequired("snapshot-id")
}

var bsDescribeSnapshotsCmd = &cobra.Command{
	Use:   "describe-snapshots",
	Short: "Describe snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("snapshot-id")

		if id != "" {
			result, err := client.GetSnapshot(ctx, id)
			if err != nil {
				exitWithError("Failed to get snapshot", err)
			}
			if output == "json" {
				printJSON(result)
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
		} else {
			result, err := client.ListSnapshots(ctx)
			if err != nil {
				exitWithError("Failed to list snapshots", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tSTATUS\tSIZE\tVOLUME_ID\tCREATED")
			for _, s := range result.Snapshots {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d GB\t%s\t%s\n",
					s.ID, s.Name, s.Status, s.Size, s.VolumeID, s.CreatedAt)
			}
			w.Flush()
		}
	},
}

var bsCreateSnapshotCmd = &cobra.Command{
	Use:   "create-snapshot",
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
			printJSON(result)
			return
		}

		fmt.Printf("Snapshot created successfully!\n")
		fmt.Printf("ID:     %s\n", result.Snapshot.ID)
		fmt.Printf("Name:   %s\n", result.Snapshot.Name)
		fmt.Printf("Status: %s\n", result.Snapshot.Status)
	},
}

var bsDeleteSnapshotCmd = &cobra.Command{
	Use:   "delete-snapshot",
	Short: "Delete a snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		client := getBlockStorageClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("snapshot-id")

		if err := client.DeleteSnapshot(ctx, id); err != nil {
			exitWithError("Failed to delete snapshot", err)
		}

		fmt.Printf("Snapshot %s deleted successfully\n", id)
	},
}
