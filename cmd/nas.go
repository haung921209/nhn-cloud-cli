package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/nas"
	"github.com/spf13/cobra"
)

var nasCmd = &cobra.Command{
	Use:     "nas",
	Aliases: []string{"nas-storage"},
	Short:   "Manage NAS Storage volumes, snapshots, and interfaces",
}

var nasVolumeCmd = &cobra.Command{
	Use:     "volume",
	Aliases: []string{"volumes", "vol"},
	Short:   "Manage NAS volumes",
}

var nasSnapshotCmd = &cobra.Command{
	Use:     "snapshot",
	Aliases: []string{"snapshots", "snap"},
	Short:   "Manage NAS volume snapshots",
}

var nasInterfaceCmd = &cobra.Command{
	Use:     "interface",
	Aliases: []string{"interfaces", "if"},
	Short:   "Manage NAS volume interfaces",
}

func init() {
	rootCmd.AddCommand(nasCmd)

	// Volume commands
	nasCmd.AddCommand(nasVolumeCmd)
	nasVolumeCmd.AddCommand(nasVolumeListCmd)
	nasVolumeCmd.AddCommand(nasVolumeGetCmd)
	nasVolumeCmd.AddCommand(nasVolumeCreateCmd)
	nasVolumeCmd.AddCommand(nasVolumeUpdateCmd)
	nasVolumeCmd.AddCommand(nasVolumeDeleteCmd)
	nasVolumeCmd.AddCommand(nasVolumeUsageCmd)

	// Snapshot commands
	nasCmd.AddCommand(nasSnapshotCmd)
	nasSnapshotCmd.AddCommand(nasSnapshotListCmd)
	nasSnapshotCmd.AddCommand(nasSnapshotGetCmd)
	nasSnapshotCmd.AddCommand(nasSnapshotCreateCmd)
	nasSnapshotCmd.AddCommand(nasSnapshotDeleteCmd)
	nasSnapshotCmd.AddCommand(nasSnapshotRestoreCmd)

	// Interface commands
	nasCmd.AddCommand(nasInterfaceCmd)
	nasInterfaceCmd.AddCommand(nasInterfaceCreateCmd)
	nasInterfaceCmd.AddCommand(nasInterfaceDeleteCmd)

	// Volume list flags
	nasVolumeListCmd.Flags().String("name", "", "Filter by exact name")
	nasVolumeListCmd.Flags().String("name-contains", "", "Filter by name containing string")
	nasVolumeListCmd.Flags().String("subnet-id", "", "Filter by subnet ID")
	nasVolumeListCmd.Flags().Int("size", 0, "Filter by exact size (GB)")
	nasVolumeListCmd.Flags().Int("min-size", 0, "Filter by minimum size (GB)")
	nasVolumeListCmd.Flags().Int("max-size", 0, "Filter by maximum size (GB)")
	nasVolumeListCmd.Flags().Int("limit", 0, "Limit results")
	nasVolumeListCmd.Flags().Int("page", 0, "Page number")

	// Volume create flags
	nasVolumeCreateCmd.Flags().String("name", "", "Volume name (required)")
	nasVolumeCreateCmd.Flags().Int("size", 0, "Volume size in GB (required)")
	nasVolumeCreateCmd.Flags().String("description", "", "Description")
	nasVolumeCreateCmd.Flags().String("subnet-id", "", "Subnet ID for interface")
	nasVolumeCreateCmd.Flags().String("protocol", "NFS", "Mount protocol: NFS or CIFS")
	nasVolumeCreateCmd.Flags().Bool("encryption", false, "Enable encryption")
	nasVolumeCreateCmd.MarkFlagRequired("name")
	nasVolumeCreateCmd.MarkFlagRequired("size")

	// Volume update flags
	nasVolumeUpdateCmd.Flags().String("description", "", "Description")
	nasVolumeUpdateCmd.Flags().String("protocol", "", "Mount protocol: NFS or CIFS")

	// Snapshot list flags
	nasSnapshotListCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasSnapshotListCmd.MarkFlagRequired("volume-id")

	// Snapshot get flags
	nasSnapshotGetCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasSnapshotGetCmd.MarkFlagRequired("volume-id")

	// Snapshot create flags
	nasSnapshotCreateCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasSnapshotCreateCmd.Flags().String("name", "", "Snapshot name (required)")
	nasSnapshotCreateCmd.MarkFlagRequired("volume-id")
	nasSnapshotCreateCmd.MarkFlagRequired("name")

	// Snapshot delete flags
	nasSnapshotDeleteCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasSnapshotDeleteCmd.MarkFlagRequired("volume-id")

	// Snapshot restore flags
	nasSnapshotRestoreCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasSnapshotRestoreCmd.MarkFlagRequired("volume-id")

	// Interface create flags
	nasInterfaceCreateCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasInterfaceCreateCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	nasInterfaceCreateCmd.MarkFlagRequired("volume-id")
	nasInterfaceCreateCmd.MarkFlagRequired("subnet-id")

	// Interface delete flags
	nasInterfaceDeleteCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasInterfaceDeleteCmd.MarkFlagRequired("volume-id")
}

func newNASClient() *nas.Client {
	return nas.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

// ================================
// Volume Commands
// ================================

var nasVolumeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all NAS volumes",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()

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
		size, _ := cmd.Flags().GetInt("size")
		if size > 0 {
			input.SizeGB = &size
		}
		minSize, _ := cmd.Flags().GetInt("min-size")
		if minSize > 0 {
			input.MinSizeGB = &minSize
		}
		maxSize, _ := cmd.Flags().GetInt("max-size")
		if maxSize > 0 {
			input.MaxSizeGB = &maxSize
		}
		limit, _ := cmd.Flags().GetInt("limit")
		if limit > 0 {
			input.Limit = &limit
		}
		page, _ := cmd.Flags().GetInt("page")
		if page > 0 {
			input.Page = &page
		}

		result, err := client.ListVolumes(context.Background(), input)
		if err != nil {
			exitWithError("Failed to list volumes", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

		if result.Paging.TotalCount > 0 {
			fmt.Printf("\nTotal: %d volumes\n", result.Paging.TotalCount)
		}
	},
}

var nasVolumeGetCmd = &cobra.Command{
	Use:   "get [volume-id]",
	Short: "Get volume details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		result, err := client.GetVolume(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get volume", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		v := result.Volume
		fmt.Printf("ID:           %s\n", v.ID)
		fmt.Printf("Name:         %s\n", v.Name)
		if v.Description != nil {
			fmt.Printf("Description:  %s\n", *v.Description)
		}
		fmt.Printf("Size (GB):    %d\n", v.SizeGB)
		fmt.Printf("Status:       %s\n", v.Status)
		fmt.Printf("Protocol:     %s\n", v.MountProtocol.Protocol)
		fmt.Printf("Encryption:   %v\n", v.Encryption.Enabled)
		fmt.Printf("Project ID:   %s\n", v.ProjectID)
		fmt.Printf("Tenant ID:    %s\n", v.TenantID)
		fmt.Printf("Created:      %s\n", v.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:      %s\n", v.UpdatedAt.Format("2006-01-02 15:04:05"))

		if len(v.Interfaces) > 0 {
			fmt.Println("\nInterfaces:")
			for _, iface := range v.Interfaces {
				fmt.Printf("  - ID: %s, Path: %s, Status: %s, Subnet: %s\n",
					iface.ID, iface.Path, iface.Status, iface.SubnetID)
			}
		}

		if len(v.ACL) > 0 {
			fmt.Printf("\nACL:          %s\n", strings.Join(v.ACL, ", "))
		}

		fmt.Println("\nSnapshot Policy:")
		fmt.Printf("  Max Scheduled Count: %d\n", v.SnapshotPolicy.MaxScheduledCount)
		fmt.Printf("  Reserve Percent:     %d\n", v.SnapshotPolicy.ReservePercent)
		fmt.Printf("  Schedule Time:       %s (Offset: %s)\n",
			v.SnapshotPolicy.Schedule.Time, v.SnapshotPolicy.Schedule.TimeOffset)
	},
}

var nasVolumeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new NAS volume",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		name, _ := cmd.Flags().GetString("name")
		size, _ := cmd.Flags().GetInt("size")
		description, _ := cmd.Flags().GetString("description")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		protocol, _ := cmd.Flags().GetString("protocol")
		encryption, _ := cmd.Flags().GetBool("encryption")

		if size < 300 || size > 10240 {
			exitWithError("size must be between 300 and 10240 GB", nil)
		}

		if protocol != "" && protocol != "NFS" && protocol != "CIFS" {
			exitWithError("protocol must be NFS or CIFS", nil)
		}

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

		result, err := client.CreateVolume(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create volume", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Volume created: %s\n", result.Volume.ID)
		fmt.Printf("Name: %s\n", result.Volume.Name)
		fmt.Printf("Size (GB): %d\n", result.Volume.SizeGB)
		fmt.Printf("Status: %s\n", result.Volume.Status)
	},
}

var nasVolumeUpdateCmd = &cobra.Command{
	Use:   "update [volume-id]",
	Short: "Update a NAS volume",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		description, _ := cmd.Flags().GetString("description")
		protocol, _ := cmd.Flags().GetString("protocol")

		input := &nas.UpdateVolumeInput{}

		if description != "" {
			input.Description = &description
		}
		if protocol != "" {
			input.MountProtocol = &nas.MountProtocol{Protocol: protocol}
		}

		result, err := client.UpdateVolume(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update volume", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Volume updated: %s\n", result.Volume.ID)
	},
}

var nasVolumeDeleteCmd = &cobra.Command{
	Use:   "delete [volume-id]",
	Short: "Delete a NAS volume",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		if err := client.DeleteVolume(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete volume", err)
		}
		fmt.Printf("Volume %s deleted\n", args[0])
	},
}

var nasVolumeUsageCmd = &cobra.Command{
	Use:   "usage [volume-id]",
	Short: "Get volume usage information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		result, err := client.GetVolumeUsage(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get volume usage", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Volume ID:              %s\n", args[0])
		fmt.Printf("Used (GB):              %d\n", result.Usage.UsedGB)
		fmt.Printf("Snapshot Reserve (GB):  %d\n", result.Usage.SnapshotReserveGB)
	},
}

// ================================
// Snapshot Commands
// ================================

var nasSnapshotListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all snapshots for a volume",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		volumeID, _ := cmd.Flags().GetString("volume-id")

		result, err := client.ListSnapshots(context.Background(), volumeID)
		if err != nil {
			exitWithError("Failed to list snapshots", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var nasSnapshotGetCmd = &cobra.Command{
	Use:   "get [snapshot-id]",
	Short: "Get snapshot details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		volumeID, _ := cmd.Flags().GetString("volume-id")

		result, err := client.GetSnapshot(context.Background(), volumeID, args[0])
		if err != nil {
			exitWithError("Failed to get snapshot", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		s := result.Snapshot
		fmt.Printf("ID:        %s\n", s.ID)
		fmt.Printf("Name:      %s\n", s.Name)
		fmt.Printf("Size:      %d\n", s.Size)
		fmt.Printf("Preserved: %v\n", s.Preserved)
		fmt.Printf("Created:   %s\n", s.CreatedAt.Format("2006-01-02 15:04:05"))
		if s.ReclaimableSpace != nil {
			fmt.Printf("Reclaimable Space: %d\n", *s.ReclaimableSpace)
		}
	},
}

var nasSnapshotCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		volumeID, _ := cmd.Flags().GetString("volume-id")
		name, _ := cmd.Flags().GetString("name")

		input := &nas.CreateSnapshotInput{
			Name: name,
		}

		result, err := client.CreateSnapshot(context.Background(), volumeID, input)
		if err != nil {
			exitWithError("Failed to create snapshot", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Snapshot created: %s\n", result.Snapshot.ID)
		fmt.Printf("Name: %s\n", result.Snapshot.Name)
	},
}

var nasSnapshotDeleteCmd = &cobra.Command{
	Use:   "delete [snapshot-id]",
	Short: "Delete a snapshot",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		volumeID, _ := cmd.Flags().GetString("volume-id")

		if err := client.DeleteSnapshot(context.Background(), volumeID, args[0]); err != nil {
			exitWithError("Failed to delete snapshot", err)
		}
		fmt.Printf("Snapshot %s deleted\n", args[0])
	},
}

var nasSnapshotRestoreCmd = &cobra.Command{
	Use:   "restore [snapshot-id]",
	Short: "Restore volume from snapshot",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		volumeID, _ := cmd.Flags().GetString("volume-id")

		if err := client.RestoreSnapshot(context.Background(), volumeID, args[0]); err != nil {
			exitWithError("Failed to restore snapshot", err)
		}
		fmt.Printf("Volume restored from snapshot %s\n", args[0])
	},
}

// ================================
// Interface Commands
// ================================

var nasInterfaceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new interface for a volume",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		volumeID, _ := cmd.Flags().GetString("volume-id")
		subnetID, _ := cmd.Flags().GetString("subnet-id")

		input := &nas.CreateInterfaceInput{
			SubnetID: subnetID,
		}

		result, err := client.CreateInterface(context.Background(), volumeID, input)
		if err != nil {
			exitWithError("Failed to create interface", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Interface created: %s\n", result.Interface.ID)
		fmt.Printf("Path: %s\n", result.Interface.Path)
		fmt.Printf("Status: %s\n", result.Interface.Status)
		fmt.Printf("Subnet ID: %s\n", result.Interface.SubnetID)
	},
}

var nasInterfaceDeleteCmd = &cobra.Command{
	Use:   "delete [interface-id]",
	Short: "Delete an interface from a volume",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		volumeID, _ := cmd.Flags().GetString("volume-id")

		if err := client.DeleteInterface(context.Background(), volumeID, args[0]); err != nil {
			exitWithError("Failed to delete interface", err)
		}
		fmt.Printf("Interface %s deleted\n", args[0])
	},
}
