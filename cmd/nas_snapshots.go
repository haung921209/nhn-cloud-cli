package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/nas"
	"github.com/spf13/cobra"
)

func init() {
	nasCmd.AddCommand(nasDescribeSnapshotsCmd)
	nasCmd.AddCommand(nasCreateSnapshotCmd)
	nasCmd.AddCommand(nasDeleteSnapshotCmd)
	nasCmd.AddCommand(nasRestoreSnapshotCmd)

	nasDescribeSnapshotsCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasDescribeSnapshotsCmd.MarkFlagRequired("volume-id")

	nasCreateSnapshotCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasCreateSnapshotCmd.Flags().String("name", "", "Snapshot name (required)")
	nasCreateSnapshotCmd.MarkFlagRequired("volume-id")
	nasCreateSnapshotCmd.MarkFlagRequired("name")

	nasDeleteSnapshotCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasDeleteSnapshotCmd.Flags().String("snapshot-id", "", "Snapshot ID (required)")
	nasDeleteSnapshotCmd.MarkFlagRequired("volume-id")
	nasDeleteSnapshotCmd.MarkFlagRequired("snapshot-id")

	nasRestoreSnapshotCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasRestoreSnapshotCmd.Flags().String("snapshot-id", "", "Snapshot ID (required)")
	nasRestoreSnapshotCmd.MarkFlagRequired("volume-id")
	nasRestoreSnapshotCmd.MarkFlagRequired("snapshot-id")
}

var nasDescribeSnapshotsCmd = &cobra.Command{
	Use:     "describe-snapshots",
	Aliases: []string{"list-snapshots"},
	Short:   "List snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		ctx := context.Background()
		volID, _ := cmd.Flags().GetString("volume-id")

		result, err := client.ListSnapshots(ctx, volID)
		if err != nil {
			exitWithError("Failed to list snapshots", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSIZE\tPRESERVED\tCREATED")
		for _, s := range result.Snapshots {
			fmt.Fprintf(w, "%s\t%s\t%d\t%v\t%s\n",
				s.ID, s.Name, s.Size, s.Preserved, s.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		w.Flush()
	},
}

var nasCreateSnapshotCmd = &cobra.Command{
	Use:   "create-snapshot",
	Short: "Create a new snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		ctx := context.Background()
		volID, _ := cmd.Flags().GetString("volume-id")
		name, _ := cmd.Flags().GetString("name")

		input := &nas.CreateSnapshotInput{
			Name: name,
		}

		result, err := client.CreateSnapshot(ctx, volID, input)
		if err != nil {
			exitWithError("Failed to create snapshot", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Snapshot created: %s (%s)\n", result.Snapshot.Name, result.Snapshot.ID)
	},
}

var nasDeleteSnapshotCmd = &cobra.Command{
	Use:   "delete-snapshot",
	Short: "Delete a snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		ctx := context.Background()
		volID, _ := cmd.Flags().GetString("volume-id")
		snapID, _ := cmd.Flags().GetString("snapshot-id")

		if err := client.DeleteSnapshot(ctx, volID, snapID); err != nil {
			exitWithError("Failed to delete snapshot", err)
		}

		fmt.Printf("Snapshot %s deleted\n", snapID)
	},
}

var nasRestoreSnapshotCmd = &cobra.Command{
	Use:   "restore-snapshot",
	Short: "Restore volume from snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		ctx := context.Background()
		volID, _ := cmd.Flags().GetString("volume-id")
		snapID, _ := cmd.Flags().GetString("snapshot-id")

		if err := client.RestoreSnapshot(ctx, volID, snapID); err != nil {
			exitWithError("Failed to restore snapshot", err)
		}

		fmt.Printf("Volume restored from snapshot %s\n", snapID)
	},
}
