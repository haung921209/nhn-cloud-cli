package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/object"
	"github.com/spf13/cobra"
)

var objectStorageCmd = &cobra.Command{
	Use:     "object-storage",
	Aliases: []string{"os", "swift"},
	Short:   "Manage Object Storage containers and objects",
	Long:    `Manage object storage containers and objects (Swift-compatible).`,
}

func init() {
	rootCmd.AddCommand(objectStorageCmd)

	objectStorageCmd.AddCommand(osAccountInfoCmd)
	objectStorageCmd.AddCommand(osContainersCmd)
	objectStorageCmd.AddCommand(osContainerInfoCmd)
	objectStorageCmd.AddCommand(osContainerCreateCmd)
	objectStorageCmd.AddCommand(osContainerUpdateCmd)
	objectStorageCmd.AddCommand(osContainerDeleteCmd)
	objectStorageCmd.AddCommand(osObjectsCmd)
	objectStorageCmd.AddCommand(osObjectInfoCmd)
	objectStorageCmd.AddCommand(osObjectUploadCmd)
	objectStorageCmd.AddCommand(osObjectDownloadCmd)
	objectStorageCmd.AddCommand(osObjectCopyCmd)
	objectStorageCmd.AddCommand(osObjectDeleteCmd)

	osContainersCmd.Flags().String("prefix", "", "Filter by prefix")
	osContainersCmd.Flags().String("marker", "", "Pagination marker")
	osContainersCmd.Flags().Int("limit", 0, "Maximum number of containers")

	osContainerCreateCmd.Flags().String("name", "", "Container name (required)")
	osContainerCreateCmd.Flags().String("storage-class", "", "Storage class: Standard or Economy")
	osContainerCreateCmd.Flags().Int("worm-retention-day", 0, "Object lock period in days (creates WORM container)")
	osContainerCreateCmd.Flags().String("read-acl", "", "Read access control list")
	osContainerCreateCmd.Flags().String("write-acl", "", "Write access control list")
	osContainerCreateCmd.MarkFlagRequired("name")

	osContainerUpdateCmd.Flags().String("read-acl", "", "Read access control list")
	osContainerUpdateCmd.Flags().String("write-acl", "", "Write access control list")
	osContainerUpdateCmd.Flags().String("view-acl", "", "View access control list")
	osContainerUpdateCmd.Flags().Int("object-lifecycle", -1, "Object lifecycle in days (-1 to skip, 0 to clear)")
	osContainerUpdateCmd.Flags().String("transfer-to", "", "Container for expired objects")
	osContainerUpdateCmd.Flags().String("history-location", "", "Archive container for versioning")
	osContainerUpdateCmd.Flags().Int("versions-retention", -1, "Version retention in days")
	osContainerUpdateCmd.Flags().String("web-index", "", "Static website index document")
	osContainerUpdateCmd.Flags().String("web-error", "", "Static website error document suffix")
	osContainerUpdateCmd.Flags().String("cors-allow-origin", "", "CORS allowed origins")
	osContainerUpdateCmd.Flags().Int("worm-retention-day", -1, "Object lock period (can only extend)")

	osObjectsCmd.Flags().String("prefix", "", "Filter by prefix")
	osObjectsCmd.Flags().String("delimiter", "", "Delimiter for pseudo-directories")
	osObjectsCmd.Flags().String("marker", "", "Pagination marker")
	osObjectsCmd.Flags().Int("limit", 0, "Maximum number of objects")

	osObjectUploadCmd.Flags().String("file", "", "Local file path (required)")
	osObjectUploadCmd.Flags().String("content-type", "", "Content type (auto-detected if not set)")
	osObjectUploadCmd.Flags().Int64("delete-after", 0, "Object TTL in seconds")
	osObjectUploadCmd.MarkFlagRequired("file")

	osObjectDownloadCmd.Flags().String("output", "", "Output file path (defaults to object name)")

	osObjectCopyCmd.Flags().String("dest-container", "", "Destination container (required)")
	osObjectCopyCmd.Flags().String("dest-object", "", "Destination object name (defaults to source name)")
	osObjectCopyCmd.MarkFlagRequired("dest-container")
}

func getObjectStorageClient() *object.Client {
	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	return object.NewClient(getRegion(), creds, nil, debug)
}

var osAccountInfoCmd = &cobra.Command{
	Use:   "account-info",
	Short: "Show storage account information",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		info, err := client.GetAccountInfo(ctx)
		if err != nil {
			exitWithError("Failed to get account info", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(info, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Container Count: %d\n", info.ContainerCount)
		fmt.Printf("Object Count:    %d\n", info.ObjectCount)
		fmt.Printf("Bytes Used:      %s\n", formatSize(info.BytesUsed))
	},
}

var osContainersCmd = &cobra.Command{
	Use:   "containers",
	Short: "List all containers",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var osContainerInfoCmd = &cobra.Command{
	Use:   "container-info [container-name]",
	Short: "Show container information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		info, err := client.GetContainerInfo(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get container info", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(info, "", "  ")
			fmt.Println(string(data))
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
		if info.ObjectLifecycle > 0 {
			fmt.Printf("Object Lifecycle:  %d days\n", info.ObjectLifecycle)
		}
		if info.ObjectTransferTo != "" {
			fmt.Printf("Transfer To:       %s\n", info.ObjectTransferTo)
		}
		if info.HistoryLocation != "" {
			fmt.Printf("History Location:  %s\n", info.HistoryLocation)
		}
		if info.VersionsRetention > 0 {
			fmt.Printf("Versions Retention: %d days\n", info.VersionsRetention)
		}
		if info.WormRetentionDay > 0 {
			fmt.Printf("WORM Retention:    %d days\n", info.WormRetentionDay)
		}
		if info.WebIndex != "" {
			fmt.Printf("Web Index:         %s\n", info.WebIndex)
		}
		if info.WebError != "" {
			fmt.Printf("Web Error:         %s\n", info.WebError)
		}
		if info.CORSAllowOrigin != "" {
			fmt.Printf("CORS Allow Origin: %s\n", info.CORSAllowOrigin)
		}
	},
}

var osContainerCreateCmd = &cobra.Command{
	Use:   "container-create",
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

var osContainerUpdateCmd = &cobra.Command{
	Use:   "container-update [container-name]",
	Short: "Update container settings",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		input := &object.UpdateContainerInput{
			Name: args[0],
		}

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

		fmt.Printf("Container %s updated successfully\n", args[0])
	},
}

var osContainerDeleteCmd = &cobra.Command{
	Use:   "container-delete [container-name]",
	Short: "Delete a container (must be empty)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		if err := client.DeleteContainer(ctx, args[0]); err != nil {
			exitWithError("Failed to delete container", err)
		}

		fmt.Printf("Container %s deleted successfully\n", args[0])
	},
}

var osObjectsCmd = &cobra.Command{
	Use:   "objects [container-name]",
	Short: "List objects in a container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		prefix, _ := cmd.Flags().GetString("prefix")
		delimiter, _ := cmd.Flags().GetString("delimiter")
		marker, _ := cmd.Flags().GetString("marker")
		limit, _ := cmd.Flags().GetInt("limit")

		input := &object.ListObjectsInput{
			Prefix:    prefix,
			Delimiter: delimiter,
			Marker:    marker,
			Limit:     limit,
		}

		result, err := client.ListObjects(ctx, args[0], input)
		if err != nil {
			exitWithError("Failed to list objects", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tSIZE\tCONTENT_TYPE\tLAST_MODIFIED")
		for _, obj := range result.Objects {
			if obj.Subdir != "" {
				fmt.Fprintf(w, "%s\t-\t(directory)\t-\n", obj.Subdir)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					obj.Name, formatSize(obj.Bytes), obj.ContentType, obj.LastModified)
			}
		}
		w.Flush()
	},
}

var osObjectInfoCmd = &cobra.Command{
	Use:   "object-info [container-name] [object-name]",
	Short: "Show object information",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		info, err := client.GetObjectInfo(ctx, args[0], args[1])
		if err != nil {
			exitWithError("Failed to get object info", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(info, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Content-Type:   %s\n", info.ContentType)
		fmt.Printf("Content-Length: %s\n", formatSize(info.ContentLength))
		fmt.Printf("ETag:           %s\n", info.ETag)
		if info.DeleteAt != nil {
			fmt.Printf("Delete-At:      %d\n", *info.DeleteAt)
		}
		if info.WormRetainUntil != nil {
			fmt.Printf("WORM Until:     %d\n", *info.WormRetainUntil)
		}
		if info.ObjectManifest != "" {
			fmt.Printf("DLO Manifest:   %s\n", info.ObjectManifest)
		}
		if info.StaticLargeObject {
			fmt.Printf("SLO:            true\n")
			fmt.Printf("Manifest ETag:  %s\n", info.ManifestETag)
		}
		if len(info.CustomMetadata) > 0 {
			fmt.Println("Custom Metadata:")
			for k, v := range info.CustomMetadata {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}
	},
}

var osObjectUploadCmd = &cobra.Command{
	Use:   "object-upload [container-name] [object-name]",
	Short: "Upload an object",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		filePath, _ := cmd.Flags().GetString("file")
		contentType, _ := cmd.Flags().GetString("content-type")
		deleteAfter, _ := cmd.Flags().GetInt64("delete-after")

		file, err := os.Open(filePath)
		if err != nil {
			exitWithError("Failed to open file", err)
		}
		defer file.Close()

		if contentType == "" {
			contentType = detectContentType(filePath)
		}

		input := &object.PutObjectInput{
			Container:   args[0],
			ObjectName:  args[1],
			Body:        file,
			ContentType: contentType,
		}

		if deleteAfter > 0 {
			input.DeleteAfter = &deleteAfter
		}

		result, err := client.PutObject(ctx, input)
		if err != nil {
			exitWithError("Failed to upload object", err)
		}

		fmt.Printf("Object %s/%s uploaded successfully\n", args[0], args[1])
		fmt.Printf("ETag: %s\n", result.ETag)
	},
}

var osObjectDownloadCmd = &cobra.Command{
	Use:   "object-download [container-name] [object-name]",
	Short: "Download an object",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		outputPath, _ := cmd.Flags().GetString("output")
		if outputPath == "" {
			parts := strings.Split(args[1], "/")
			outputPath = parts[len(parts)-1]
		}

		result, err := client.GetObject(ctx, args[0], args[1])
		if err != nil {
			exitWithError("Failed to download object", err)
		}
		defer result.Body.Close()

		outFile, err := os.Create(outputPath)
		if err != nil {
			exitWithError("Failed to create output file", err)
		}
		defer outFile.Close()

		written, err := outFile.ReadFrom(result.Body)
		if err != nil {
			exitWithError("Failed to write file", err)
		}

		fmt.Printf("Downloaded %s to %s (%s)\n", args[1], outputPath, formatSize(written))
	},
}

var osObjectCopyCmd = &cobra.Command{
	Use:   "object-copy [container-name] [object-name]",
	Short: "Copy an object to another container",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		destContainer, _ := cmd.Flags().GetString("dest-container")
		destObject, _ := cmd.Flags().GetString("dest-object")
		if destObject == "" {
			destObject = args[1]
		}

		input := &object.CopyObjectInput{
			SourceContainer:       args[0],
			SourceObjectName:      args[1],
			DestinationContainer:  destContainer,
			DestinationObjectName: destObject,
		}

		if err := client.CopyObject(ctx, input); err != nil {
			exitWithError("Failed to copy object", err)
		}

		fmt.Printf("Copied %s/%s to %s/%s\n", args[0], args[1], destContainer, destObject)
	},
}

var osObjectDeleteCmd = &cobra.Command{
	Use:   "object-delete [container-name] [object-name]",
	Short: "Delete an object",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		if err := client.DeleteObject(ctx, args[0], args[1]); err != nil {
			exitWithError("Failed to delete object", err)
		}

		fmt.Printf("Object %s/%s deleted successfully\n", args[0], args[1])
	},
}

func detectContentType(filePath string) string {
	ext := strings.ToLower(filePath)
	switch {
	case strings.HasSuffix(ext, ".html"), strings.HasSuffix(ext, ".htm"):
		return "text/html"
	case strings.HasSuffix(ext, ".css"):
		return "text/css"
	case strings.HasSuffix(ext, ".js"):
		return "application/javascript"
	case strings.HasSuffix(ext, ".json"):
		return "application/json"
	case strings.HasSuffix(ext, ".xml"):
		return "application/xml"
	case strings.HasSuffix(ext, ".txt"):
		return "text/plain"
	case strings.HasSuffix(ext, ".pdf"):
		return "application/pdf"
	case strings.HasSuffix(ext, ".png"):
		return "image/png"
	case strings.HasSuffix(ext, ".jpg"), strings.HasSuffix(ext, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(ext, ".gif"):
		return "image/gif"
	case strings.HasSuffix(ext, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(ext, ".zip"):
		return "application/zip"
	case strings.HasSuffix(ext, ".tar"):
		return "application/x-tar"
	case strings.HasSuffix(ext, ".gz"):
		return "application/gzip"
	default:
		return "application/octet-stream"
	}
}

func parseOBSTime(s string) time.Time {
	layouts := []string{
		"2006-01-02T15:04:05.000000",
		"2006-01-02T15:04:05",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
