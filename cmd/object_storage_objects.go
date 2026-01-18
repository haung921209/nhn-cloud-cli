package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/object"
	"github.com/spf13/cobra"
)

func init() {
	objectStorageCmd.AddCommand(objectStorageListObjectsCmd)
	objectStorageCmd.AddCommand(objectStorageDescribeObjectCmd)
	objectStorageCmd.AddCommand(objectStoragePutObjectCmd)
	objectStorageCmd.AddCommand(objectStorageGetObjectCmd)
	objectStorageCmd.AddCommand(objectStorageCopyObjectCmd)
	objectStorageCmd.AddCommand(objectStorageDeleteObjectCmd)

	objectStorageListObjectsCmd.Flags().String("container-name", "", "Container Name (required)")
	objectStorageListObjectsCmd.Flags().String("prefix", "", "Filter by prefix")
	objectStorageListObjectsCmd.Flags().String("delimiter", "", "Delimiter for pseudo-directories")
	objectStorageListObjectsCmd.Flags().String("marker", "", "Pagination marker")
	objectStorageListObjectsCmd.Flags().Int("limit", 0, "Maximum number of objects")
	objectStorageListObjectsCmd.MarkFlagRequired("container-name")

	objectStoragePutObjectCmd.Flags().String("container-name", "", "Container Name (required)")
	objectStoragePutObjectCmd.Flags().String("object-name", "", "Object Name (required)")
	objectStoragePutObjectCmd.Flags().String("file", "", "Local file path (required)")
	objectStoragePutObjectCmd.Flags().String("content-type", "", "Content type (auto-detected if not set)")
	objectStoragePutObjectCmd.Flags().Int64("delete-after", 0, "Object TTL in seconds")
	objectStoragePutObjectCmd.MarkFlagRequired("container-name")
	objectStoragePutObjectCmd.MarkFlagRequired("object-name")
	objectStoragePutObjectCmd.MarkFlagRequired("file")

	objectStorageGetObjectCmd.Flags().String("container-name", "", "Container Name (required)")
	objectStorageGetObjectCmd.Flags().String("object-name", "", "Object Name (required)")
	objectStorageGetObjectCmd.Flags().String("output", "", "Output file path (defaults to object name)")
	objectStorageGetObjectCmd.MarkFlagRequired("container-name")
	objectStorageGetObjectCmd.MarkFlagRequired("object-name")

	objectStorageCopyObjectCmd.Flags().String("container-name", "", "Source Container Name (required)")
	objectStorageCopyObjectCmd.Flags().String("object-name", "", "Source Object Name (required)")
	objectStorageCopyObjectCmd.Flags().String("dest-container", "", "Destination container (required)")
	objectStorageCopyObjectCmd.Flags().String("dest-object", "", "Destination object name (defaults to source name)")
	objectStorageCopyObjectCmd.MarkFlagRequired("container-name")
	objectStorageCopyObjectCmd.MarkFlagRequired("object-name")
	objectStorageCopyObjectCmd.MarkFlagRequired("dest-container")

	objectStorageDeleteObjectCmd.Flags().String("container-name", "", "Container Name (required)")
	objectStorageDeleteObjectCmd.Flags().String("object-name", "", "Object Name (required)")
	objectStorageDeleteObjectCmd.MarkFlagRequired("container-name")
	objectStorageDeleteObjectCmd.MarkFlagRequired("object-name")

	objectStorageDescribeObjectCmd.Flags().String("container-name", "", "Container Name (required)")
	objectStorageDescribeObjectCmd.Flags().String("object-name", "", "Object Name (required)")
	objectStorageDescribeObjectCmd.MarkFlagRequired("container-name")
	objectStorageDescribeObjectCmd.MarkFlagRequired("object-name")
}

var objectStorageListObjectsCmd = &cobra.Command{
	Use:   "list-objects",
	Short: "List objects in a container",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()
		containerName, _ := cmd.Flags().GetString("container-name")

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

		result, err := client.ListObjects(ctx, containerName, input)
		if err != nil {
			exitWithError("Failed to list objects", err)
		}

		if output == "json" {
			printJSON(result)
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

var objectStorageDescribeObjectCmd = &cobra.Command{
	Use:   "describe-object",
	Short: "Show object information (Head Object)",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()
		containerName, _ := cmd.Flags().GetString("container-name")
		objectName, _ := cmd.Flags().GetString("object-name")

		info, err := client.GetObjectInfo(ctx, containerName, objectName)
		if err != nil {
			exitWithError("Failed to get object info", err)
		}

		if output == "json" {
			printJSON(info)
			return
		}

		fmt.Printf("Content-Type:   %s\n", info.ContentType)
		fmt.Printf("Content-Length: %s\n", formatSize(info.ContentLength))
		fmt.Printf("ETag:           %s\n", info.ETag)
		// ... (Same as original but concise)
	},
}

var objectStoragePutObjectCmd = &cobra.Command{
	Use:   "put-object",
	Short: "Upload an object",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		containerName, _ := cmd.Flags().GetString("container-name")
		objectName, _ := cmd.Flags().GetString("object-name")
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
			Container:   containerName,
			ObjectName:  objectName,
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

		fmt.Printf("Object %s/%s uploaded successfully\n", containerName, objectName)
		fmt.Printf("ETag: %s\n", result.ETag)
	},
}

var objectStorageGetObjectCmd = &cobra.Command{
	Use:   "get-object",
	Short: "Download an object",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		containerName, _ := cmd.Flags().GetString("container-name")
		objectName, _ := cmd.Flags().GetString("object-name")
		outputPath, _ := cmd.Flags().GetString("output")

		if outputPath == "" {
			parts := strings.Split(objectName, "/")
			outputPath = parts[len(parts)-1]
		}

		result, err := client.GetObject(ctx, containerName, objectName)
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

		fmt.Printf("Downloaded %s to %s (%s)\n", objectName, outputPath, formatSize(written))
	},
}

var objectStorageCopyObjectCmd = &cobra.Command{
	Use:   "copy-object",
	Short: "Copy an object",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		containerName, _ := cmd.Flags().GetString("container-name")
		objectName, _ := cmd.Flags().GetString("object-name")
		destContainer, _ := cmd.Flags().GetString("dest-container")
		destObject, _ := cmd.Flags().GetString("dest-object")
		if destObject == "" {
			destObject = objectName
		}

		input := &object.CopyObjectInput{
			SourceContainer:       containerName,
			SourceObjectName:      objectName,
			DestinationContainer:  destContainer,
			DestinationObjectName: destObject,
		}

		if err := client.CopyObject(ctx, input); err != nil {
			exitWithError("Failed to copy object", err)
		}

		fmt.Printf("Copied %s/%s to %s/%s\n", containerName, objectName, destContainer, destObject)
	},
}

var objectStorageDeleteObjectCmd = &cobra.Command{
	Use:   "delete-object",
	Short: "Delete an object",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()
		containerName, _ := cmd.Flags().GetString("container-name")
		objectName, _ := cmd.Flags().GetString("object-name")

		if err := client.DeleteObject(ctx, containerName, objectName); err != nil {
			exitWithError("Failed to delete object", err)
		}

		fmt.Printf("Object %s/%s deleted successfully\n", containerName, objectName)
	},
}
