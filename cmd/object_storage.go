package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

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

	objectStorageCmd.AddCommand(osContainersCmd)
	objectStorageCmd.AddCommand(osContainerCreateCmd)
	objectStorageCmd.AddCommand(osContainerDeleteCmd)
	objectStorageCmd.AddCommand(osObjectsCmd)
	objectStorageCmd.AddCommand(osObjectDeleteCmd)

	osContainerCreateCmd.Flags().String("name", "", "Container name (required)")
	osContainerCreateCmd.MarkFlagRequired("name")

	osObjectsCmd.Flags().String("prefix", "", "Filter by prefix")
	osObjectsCmd.Flags().String("delimiter", "", "Delimiter for pseudo-directories")
	osObjectsCmd.Flags().Int("limit", 0, "Maximum number of objects to return")
}

func getObjectStorageClient() *object.Client {
	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	return object.NewClient(getRegion(), creds, nil, debug)
}

var osContainersCmd = &cobra.Command{
	Use:   "containers",
	Short: "List all containers",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		result, err := client.ListContainers(ctx)
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

var osContainerCreateCmd = &cobra.Command{
	Use:   "container-create",
	Short: "Create a new container",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")

		input := &object.CreateContainerInput{
			Name: name,
		}

		if err := client.CreateContainer(ctx, input); err != nil {
			exitWithError("Failed to create container", err)
		}

		fmt.Printf("Container %s created successfully\n", name)
	},
}

var osContainerDeleteCmd = &cobra.Command{
	Use:   "container-delete [container-name]",
	Short: "Delete a container",
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
		limit, _ := cmd.Flags().GetInt("limit")

		input := &object.ListObjectsInput{
			Prefix:    prefix,
			Delimiter: delimiter,
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
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				obj.Name, formatSize(obj.Bytes), obj.ContentType, obj.LastModified.Format("2006-01-02 15:04:05"))
		}
		w.Flush()
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
