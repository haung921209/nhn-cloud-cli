package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncs"
	"github.com/spf13/cobra"
)

func init() {
	ncsCmd.AddCommand(ncsDescribeVolumesCmd)
	ncsCmd.AddCommand(ncsAttachVolumeCmd)

	ncsAttachVolumeCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsAttachVolumeCmd.Flags().String("volume-id", "", "Volume ID (required)")
	ncsAttachVolumeCmd.Flags().String("mount-path", "", "Mount path in container (required)")
	ncsAttachVolumeCmd.Flags().Bool("read-only", false, "Mount as read-only")
	ncsAttachVolumeCmd.MarkFlagRequired("workload-id")
	ncsAttachVolumeCmd.MarkFlagRequired("volume-id")
	ncsAttachVolumeCmd.MarkFlagRequired("mount-path")
}

var ncsDescribeVolumesCmd = &cobra.Command{
	Use:     "describe-volumes",
	Aliases: []string{"list-volumes", "volumes"},
	Short:   "List volumes",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
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
		fmt.Fprintln(w, "ID\tNAME\tSIZE\tSTATUS\tCREATED")
		for _, v := range result.Volumes {
			fmt.Fprintf(w, "%s\t%s\t%dGB\t%s\t%s\n",
				v.VolumeID, v.Name, v.Size, v.Status, v.CreatedAt)
		}
		w.Flush()
	},
}

var ncsAttachVolumeCmd = &cobra.Command{
	Use:   "attach-volume",
	Short: "Attach a volume to a workload",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")
		volumeID, _ := cmd.Flags().GetString("volume-id")
		mountPath, _ := cmd.Flags().GetString("mount-path")
		readOnly, _ := cmd.Flags().GetBool("read-only")

		input := &ncs.VolumeAttachInput{
			VolumeID:  volumeID,
			MountPath: mountPath,
			ReadOnly:  readOnly,
		}

		_, err := client.AttachVolume(ctx, workloadID, input)
		if err != nil {
			exitWithError("Failed to attach volume", err)
		}

		fmt.Printf("Volume %s attached to workload %s at %s\n", volumeID, workloadID, mountPath)
	},
}
