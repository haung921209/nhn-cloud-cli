package cmd

import (
	"context"
	"fmt"
	"os"
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
	nksCreateNodeGroupCmd.MarkFlagRequired("cluster-id")
	nksCreateNodeGroupCmd.MarkFlagRequired("name")
	nksCreateNodeGroupCmd.MarkFlagRequired("flavor-id")

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

		input := &nks.CreateNodeGroupInput{
			Name:      name,
			FlavorID:  flavorID,
			NodeCount: nodeCount,
		}

		result, err := client.CreateNodeGroup(ctx, clusterID, input)
		if err != nil {
			exitWithError("Failed to create node group", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Node group creation initiated!\n")
		fmt.Printf("ID:   %s\n", result.ID)
		fmt.Printf("Name: %s\n", result.Name)
	},
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
