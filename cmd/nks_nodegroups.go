package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/nks"
	"github.com/spf13/cobra"
)

func init() {
	nksCmd.AddCommand(nksDescribeNodeGroupsCmd)
	nksCmd.AddCommand(nksCreateNodeGroupCmd)
	nksCmd.AddCommand(nksDeleteNodeGroupCmd)
	nksCmd.AddCommand(nksUpdateNodeGroupCmd)

	nksDescribeNodeGroupsCmd.Flags().String("cluster-id", "", "Cluster ID (required)")
	nksDescribeNodeGroupsCmd.Flags().String("node-group-id", "", "Node Group ID (optional)")
	nksDescribeNodeGroupsCmd.MarkFlagRequired("cluster-id")

	nksCreateNodeGroupCmd.Flags().String("cluster-id", "", "Cluster ID (required)")
	nksCreateNodeGroupCmd.Flags().String("name", "", "Node group name (required)")
	nksCreateNodeGroupCmd.Flags().String("flavor-id", "", "Flavor ID (required)")
	nksCreateNodeGroupCmd.Flags().Int("node-count", 1, "Number of nodes")
	nksCreateNodeGroupCmd.Flags().String("image-id", "", "Node base image ID (required)")
	nksCreateNodeGroupCmd.Flags().String("availability-zone", "", "Availability zone of the nodes (required)")
	nksCreateNodeGroupCmd.Flags().String("boot-volume-type", "General HDD", "Boot volume type")
	nksCreateNodeGroupCmd.Flags().Int("boot-volume-size", 0, "Boot volume size in GB (required)")
	nksCreateNodeGroupCmd.Flags().StringToString("label", nil, "Extra node group label key=value, e.g. ca_enable=True,user_script_v2=... (repeatable)")
	nksCreateNodeGroupCmd.MarkFlagRequired("cluster-id")
	nksCreateNodeGroupCmd.MarkFlagRequired("name")
	nksCreateNodeGroupCmd.MarkFlagRequired("flavor-id")
	nksCreateNodeGroupCmd.MarkFlagRequired("image-id")
	nksCreateNodeGroupCmd.MarkFlagRequired("availability-zone")
	nksCreateNodeGroupCmd.MarkFlagRequired("boot-volume-size")

	nksDeleteNodeGroupCmd.Flags().String("cluster-id", "", "Cluster ID (required)")
	nksDeleteNodeGroupCmd.Flags().String("node-group-id", "", "Node Group ID (required)")
	nksDeleteNodeGroupCmd.MarkFlagRequired("cluster-id")
	nksDeleteNodeGroupCmd.MarkFlagRequired("node-group-id")

	nksUpdateNodeGroupCmd.Flags().String("cluster-id", "", "Cluster ID (required)")
	nksUpdateNodeGroupCmd.Flags().String("node-group-id", "", "Node Group ID (required)")
	nksUpdateNodeGroupCmd.Flags().Int("node-count", 0, "New node count")
	nksUpdateNodeGroupCmd.MarkFlagRequired("cluster-id")
	nksUpdateNodeGroupCmd.MarkFlagRequired("node-group-id")
}

var nksDescribeNodeGroupsCmd = &cobra.Command{
	Use:   "describe-node-groups",
	Short: "Describe node groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()
		clusterID, _ := cmd.Flags().GetString("cluster-id")
		groupID, _ := cmd.Flags().GetString("node-group-id")

		if groupID != "" {
			result, err := client.GetNodeGroup(ctx, clusterID, groupID)
			if err != nil {
				exitWithError("Failed to get node group", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			fmt.Printf("ID:         %s\n", result.ID)
			fmt.Printf("Name:       %s\n", result.Name)
			fmt.Printf("Status:     %s\n", result.Status)
			fmt.Printf("Node Count: %d\n", result.NodeCount)
			fmt.Printf("Flavor:     %s\n", result.FlavorID)
			fmt.Printf("Created:    %s\n", result.CreatedAt)
		} else {
			result, err := client.ListNodeGroups(ctx, clusterID)
			if err != nil {
				exitWithError("Failed to list node groups", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tSTATUS\tNODE_COUNT\tFLAVOR")
			for _, ng := range result.NodeGroups {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
					ng.ID, ng.Name, ng.Status, ng.NodeCount, ng.FlavorID)
			}
			w.Flush()
		}
	},
}

var nksCreateNodeGroupCmd = &cobra.Command{
	Use:   "create-node-group",
	Short: "Create a new node group",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		clusterID, _ := cmd.Flags().GetString("cluster-id")
		name, _ := cmd.Flags().GetString("name")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		nodeCount, _ := cmd.Flags().GetInt("node-count")

		result, err := client.CreateNodeGroup(ctx, clusterID, nodeGroupInput(cmd, name, flavorID, nodeCount))
		if err != nil {
			exitWithError("Failed to create node group", err)
		}

		printWriteResult(result, "Node group creation initiated!\nID:   %s\nName: %s\n", result.ID, result.Name)
	},
}

// nodeGroupInput builds the create body. The API requires image_id and the
// labels availability_zone, boot_volume_type, boot_volume_size and ca_enable.
func nodeGroupInput(cmd *cobra.Command, name, flavorID string, nodeCount int) *nks.CreateNodeGroupInput {
	imageID, _ := cmd.Flags().GetString("image-id")
	az, _ := cmd.Flags().GetString("availability-zone")
	volType, _ := cmd.Flags().GetString("boot-volume-type")
	volSize, _ := cmd.Flags().GetInt("boot-volume-size")
	extra, _ := cmd.Flags().GetStringToString("label")

	labels := map[string]string{"ca_enable": "False"}
	for k, v := range extra {
		labels[k] = v
	}
	labels["availability_zone"] = az
	labels["boot_volume_type"] = volType
	labels["boot_volume_size"] = strconv.Itoa(volSize)

	return &nks.CreateNodeGroupInput{
		Name:      name,
		FlavorID:  flavorID,
		ImageID:   imageID,
		NodeCount: nodeCount,
		Labels:    labels,
	}
}

var nksUpdateNodeGroupCmd = &cobra.Command{
	Use:   "update-node-group",
	Short: "Update a node group (scale)",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		clusterID, _ := cmd.Flags().GetString("cluster-id")
		groupID, _ := cmd.Flags().GetString("node-group-id")
		nodeCount, _ := cmd.Flags().GetInt("node-count")

		input := &nks.UpdateNodeGroupInput{
			NodeCount: nodeCount,
		}

		if err := client.UpdateNodeGroup(ctx, clusterID, groupID, input); err != nil {
			exitWithError("Failed to update node group", err)
		}

		fmt.Printf("Node group %s update initiated\n", groupID)
	},
}

var nksDeleteNodeGroupCmd = &cobra.Command{
	Use:   "delete-node-group",
	Short: "Delete a node group",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()
		clusterID, _ := cmd.Flags().GetString("cluster-id")
		groupID, _ := cmd.Flags().GetString("node-group-id")

		if err := client.DeleteNodeGroup(ctx, clusterID, groupID); err != nil {
			exitWithError("Failed to delete node group", err)
		}

		fmt.Printf("Node group %s deletion initiated\n", groupID)
	},
}
