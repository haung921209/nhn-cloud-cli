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
	nasCmd.AddCommand(nasDescribeVolumesCmd)
	nasCmd.AddCommand(nasCreateVolumeCmd)
	nasCmd.AddCommand(nasUpdateVolumeCmd)
	nasCmd.AddCommand(nasDeleteVolumeCmd)
	nasCmd.AddCommand(nasVolumeUsageCmd) // Extra command specific to NAS

	nasDescribeVolumesCmd.Flags().String("name", "", "Filter by exact name")
	nasDescribeVolumesCmd.Flags().String("name-contains", "", "Filter by name containing string")
	nasDescribeVolumesCmd.Flags().String("subnet-id", "", "Filter by subnet ID")

	nasCreateVolumeCmd.Flags().String("name", "", "Volume name (required)")
	nasCreateVolumeCmd.Flags().Int("size", 0, "Volume size in GB (required)")
	nasCreateVolumeCmd.Flags().String("description", "", "Description")
	nasCreateVolumeCmd.Flags().String("subnet-id", "", "Subnet ID for interface")
	nasCreateVolumeCmd.Flags().String("protocol", "NFS", "Mount protocol: NFS or CIFS")
	nasCreateVolumeCmd.Flags().Bool("encryption", false, "Enable encryption")
	nasCreateVolumeCmd.MarkFlagRequired("name")
	nasCreateVolumeCmd.MarkFlagRequired("size")

	nasUpdateVolumeCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasUpdateVolumeCmd.Flags().String("description", "", "Description")
	nasUpdateVolumeCmd.Flags().String("protocol", "", "Mount protocol")
	nasUpdateVolumeCmd.MarkFlagRequired("volume-id")

	nasDeleteVolumeCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasDeleteVolumeCmd.MarkFlagRequired("volume-id")

	nasVolumeUsageCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasVolumeUsageCmd.MarkFlagRequired("volume-id")
}

var nasDescribeVolumesCmd = &cobra.Command{
	Use:     "describe-volumes",
	Aliases: []string{"list-volumes"},
	Short:   "List or describe NAS volumes",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		ctx := context.Background()

		input := &nas.ListVolumesInput{}
		name, _ := cmd.Flags().GetString("name")
		if name != "" {
			input.Name = name
		}
		nameContains, _ := cmd.Flags().GetString("name-contains")
		if nameContains != "" {
			input.NameContains = nameContains
		}
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		if subnetID != "" {
			input.SubnetID = subnetID
		}

		result, err := client.ListVolumes(ctx, input)
		if err != nil {
			exitWithError("Failed to list volumes", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSIZE_GB\tSTATUS\tPROTOCOL\tINTERFACES")
		for _, v := range result.Volumes {
			ifCount := len(v.Interfaces)
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%d\n",
				v.ID, v.Name, v.SizeGB, v.Status, v.MountProtocol.Protocol, ifCount)
		}
		w.Flush()
	},
}

var nasCreateVolumeCmd = &cobra.Command{
	Use:   "create-volume",
	Short: "Create a new NAS volume",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("name")
		size, _ := cmd.Flags().GetInt("size")
		description, _ := cmd.Flags().GetString("description")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		protocol, _ := cmd.Flags().GetString("protocol")
		encryption, _ := cmd.Flags().GetBool("encryption")

		input := &nas.CreateVolumeInput{
			Name:   name,
			SizeGB: size,
			MountProtocol: &nas.MountProtocol{
				Protocol: protocol,
			},
		}

		if description != "" {
			input.Description = &description
		}
		if subnetID != "" {
			input.Interfaces = []nas.CreateVolumeInterfaceInput{
				{SubnetID: subnetID},
			}
		}
		if encryption {
			input.Encryption = &nas.Encryption{Enabled: true}
		}

		result, err := client.CreateVolume(ctx, input)
		if err != nil {
			exitWithError("Failed to create volume", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Volume created: %s (%s)\n", result.Volume.Name, result.Volume.ID)
	},
}

var nasUpdateVolumeCmd = &cobra.Command{
	Use:   "update-volume",
	Short: "Update a NAS volume",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("volume-id")
		description, _ := cmd.Flags().GetString("description")
		protocol, _ := cmd.Flags().GetString("protocol")

		input := &nas.UpdateVolumeInput{}
		if description != "" {
			input.Description = &description
		}
		if protocol != "" {
			input.MountProtocol = &nas.MountProtocol{Protocol: protocol}
		}

		result, err := client.UpdateVolume(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update volume", err)
		}

		fmt.Printf("Volume updated: %s\n", result.Volume.ID)
	},
}

var nasDeleteVolumeCmd = &cobra.Command{
	Use:   "delete-volume",
	Short: "Delete a NAS volume",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("volume-id")

		if err := client.DeleteVolume(ctx, id); err != nil {
			exitWithError("Failed to delete volume", err)
		}

		fmt.Printf("Volume %s deleted\n", id)
	},
}

var nasVolumeUsageCmd = &cobra.Command{
	Use:   "get-volume-usage",
	Short: "Get volume usage information",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("volume-id")

		result, err := client.GetVolumeUsage(ctx, id)
		if err != nil {
			exitWithError("Failed to get volume usage", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Volume ID:              %s\n", id)
		fmt.Printf("Used (GB):              %d\n", result.Usage.UsedGB)
		fmt.Printf("Snapshot Reserve (GB):  %d\n", result.Usage.SnapshotReserveGB)
	},
}
