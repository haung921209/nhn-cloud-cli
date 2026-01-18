package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/object"
	"github.com/spf13/cobra"
)

func init() {
	objectStorageCmd.AddCommand(objectStorageDescribeContainersCmd)
	objectStorageCmd.AddCommand(objectStorageCreateContainerCmd)
	objectStorageCmd.AddCommand(objectStorageUpdateContainerCmd)
	objectStorageCmd.AddCommand(objectStorageDeleteContainerCmd)
	objectStorageCmd.AddCommand(objectStorageDescribeContainerCmd) // Single container info

	objectStorageDescribeContainersCmd.Flags().String("prefix", "", "Filter by prefix")
	objectStorageDescribeContainersCmd.Flags().String("marker", "", "Pagination marker")
	objectStorageDescribeContainersCmd.Flags().Int("limit", 0, "Maximum number of containers")

	objectStorageCreateContainerCmd.Flags().String("name", "", "Container name (required)")
	objectStorageCreateContainerCmd.Flags().String("storage-class", "", "Storage class: Standard or Economy")
	objectStorageCreateContainerCmd.Flags().Int("worm-retention-day", 0, "Object lock period in days (creates WORM container)")
	objectStorageCreateContainerCmd.Flags().String("read-acl", "", "Read access control list")
	objectStorageCreateContainerCmd.Flags().String("write-acl", "", "Write access control list")
	objectStorageCreateContainerCmd.MarkFlagRequired("name")

	objectStorageUpdateContainerCmd.Flags().String("container-name", "", "Container Name (required)")
	objectStorageUpdateContainerCmd.Flags().String("read-acl", "", "Read access control list")
	objectStorageUpdateContainerCmd.Flags().String("write-acl", "", "Write access control list")
	objectStorageUpdateContainerCmd.Flags().String("view-acl", "", "View access control list")
	objectStorageUpdateContainerCmd.Flags().Int("object-lifecycle", -1, "Object lifecycle in days (-1 to skip, 0 to clear)")
	objectStorageUpdateContainerCmd.Flags().String("transfer-to", "", "Container for expired objects")
	objectStorageUpdateContainerCmd.Flags().String("history-location", "", "Archive container for versioning")
	objectStorageUpdateContainerCmd.Flags().Int("versions-retention", -1, "Version retention in days")
	objectStorageUpdateContainerCmd.Flags().String("web-index", "", "Static website index document")
	objectStorageUpdateContainerCmd.Flags().String("web-error", "", "Static website error document suffix")
	objectStorageUpdateContainerCmd.Flags().String("cors-allow-origin", "", "CORS allowed origins")
	objectStorageUpdateContainerCmd.Flags().Int("worm-retention-day", -1, "Object lock period (can only extend)")
	objectStorageUpdateContainerCmd.MarkFlagRequired("container-name")

	objectStorageDeleteContainerCmd.Flags().String("container-name", "", "Container Name (required)")
	objectStorageDeleteContainerCmd.MarkFlagRequired("container-name")

	objectStorageDescribeContainerCmd.Flags().String("container-name", "", "Container Name (required)")
	objectStorageDescribeContainerCmd.MarkFlagRequired("container-name")
}

var objectStorageDescribeContainersCmd = &cobra.Command{
	Use:     "describe-containers",
	Aliases: []string{"list-containers", "ls"},
	Short:   "List all containers",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		prefix, _ := cmd.Flags().GetString("prefix")
		marker, _ := cmd.Flags().GetString("marker")
		limit, _ := cmd.Flags().GetInt("limit")

		input := &object.ListContainersInput{
			Prefix: prefix,
			Marker: marker,
			Limit:  limit,
		}

		result, err := client.ListContainers(ctx, input)
		if err != nil {
			exitWithError("Failed to list containers", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tOBJECT_COUNT\tSIZE\tLAST_MODIFIED")
		for _, c := range result.Containers {
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
				c.Name, c.Count, formatSize(c.Bytes), c.LastModified)
		}
		w.Flush()
	},
}

var objectStorageDescribeContainerCmd = &cobra.Command{
	Use:   "describe-container",
	Short: "Show container information",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("container-name")

		info, err := client.GetContainerInfo(ctx, name)
		if err != nil {
			exitWithError("Failed to get container info", err)
		}

		if output == "json" {
			printJSON(info)
			return
		}

		fmt.Printf("Name:              %s\n", info.Name)
		fmt.Printf("Object Count:      %d\n", info.ObjectCount)
		fmt.Printf("Bytes Used:        %s\n", formatSize(info.BytesUsed))
		fmt.Printf("Storage Policy:    %s\n", info.StoragePolicy)
		if info.ReadACL != "" {
			fmt.Printf("Read ACL:          %s\n", info.ReadACL)
		}
		if info.WriteACL != "" {
			fmt.Printf("Write ACL:         %s\n", info.WriteACL)
		}
		// ... (Same as original)
	},
}

var objectStorageCreateContainerCmd = &cobra.Command{
	Use:   "create-container",
	Short: "Create a new container",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		storageClass, _ := cmd.Flags().GetString("storage-class")
		wormRetentionDay, _ := cmd.Flags().GetInt("worm-retention-day")
		readACL, _ := cmd.Flags().GetString("read-acl")
		writeACL, _ := cmd.Flags().GetString("write-acl")

		input := &object.CreateContainerInput{
			Name:             name,
			StoragePolicy:    storageClass,
			WormRetentionDay: wormRetentionDay,
			ReadACL:          readACL,
			WriteACL:         writeACL,
		}

		if err := client.CreateContainer(ctx, input); err != nil {
			exitWithError("Failed to create container", err)
		}

		fmt.Printf("Container %s created successfully\n", name)
	},
}

var objectStorageUpdateContainerCmd = &cobra.Command{
	Use:   "update-container",
	Short: "Update container settings",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("container-name")

		input := &object.UpdateContainerInput{
			Name: name,
		}

		// Copied logic from original
		if v, _ := cmd.Flags().GetString("read-acl"); v != "" {
			input.ReadACL = v
		}
		if v, _ := cmd.Flags().GetString("write-acl"); v != "" {
			input.WriteACL = v
		}
		if v, _ := cmd.Flags().GetString("view-acl"); v != "" {
			input.ViewACL = v
		}
		if v, _ := cmd.Flags().GetInt("object-lifecycle"); v >= 0 {
			input.ObjectLifecycle = &v
		}
		if v, _ := cmd.Flags().GetString("transfer-to"); v != "" {
			input.ObjectTransferTo = v
		}
		if v, _ := cmd.Flags().GetString("history-location"); v != "" {
			input.HistoryLocation = v
		}
		if v, _ := cmd.Flags().GetInt("versions-retention"); v >= 0 {
			input.VersionsRetention = &v
		}
		if v, _ := cmd.Flags().GetString("web-index"); v != "" {
			input.WebIndex = v
		}
		if v, _ := cmd.Flags().GetString("web-error"); v != "" {
			input.WebError = v
		}
		if v, _ := cmd.Flags().GetString("cors-allow-origin"); v != "" {
			input.CORSAllowOrigin = v
		}
		if v, _ := cmd.Flags().GetInt("worm-retention-day"); v >= 0 {
			input.WormRetentionDay = &v
		}

		if err := client.UpdateContainer(ctx, input); err != nil {
			exitWithError("Failed to update container", err)
		}

		fmt.Printf("Container %s updated successfully\n", name)
	},
}

var objectStorageDeleteContainerCmd = &cobra.Command{
	Use:   "delete-container",
	Short: "Delete a container (must be empty)",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("container-name")

		if err := client.DeleteContainer(ctx, name); err != nil {
			exitWithError("Failed to delete container", err)
		}

		fmt.Printf("Container %s deleted successfully\n", name)
	},
}
