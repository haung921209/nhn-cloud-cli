package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/nks"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/spf13/cobra"
)

var nksCmd = &cobra.Command{
	Use:     "nks",
	Aliases: []string{"kubernetes"},
	Short:   "Manage NHN Kubernetes Service (NKS) clusters",
	Long:    `Manage NKS clusters including create, delete, scale node groups, and get kubeconfig.`,
}

func init() {
	rootCmd.AddCommand(nksCmd)

	nksCmd.AddCommand(nksListCmd)
	nksCmd.AddCommand(nksGetCmd)
	nksCmd.AddCommand(nksCreateCmd)
	nksCmd.AddCommand(nksDeleteCmd)
	nksCmd.AddCommand(nksKubeconfigCmd)
	nksCmd.AddCommand(nksTemplatesCmd)

	nksCmd.AddCommand(nksNodeGroupsCmd)
	nksCmd.AddCommand(nksNodeGroupGetCmd)
	nksCmd.AddCommand(nksNodeGroupCreateCmd)
	nksCmd.AddCommand(nksNodeGroupDeleteCmd)
	nksCmd.AddCommand(nksNodeGroupUpdateCmd)

	nksCreateCmd.Flags().String("name", "", "Cluster name (required)")
	nksCreateCmd.Flags().String("template-id", "", "Cluster template ID (optional)")
	nksCreateCmd.Flags().String("k8s-version", "", "Kubernetes version")
	nksCreateCmd.Flags().String("network-id", "", "Network ID (required)")
	nksCreateCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	nksCreateCmd.Flags().String("keypair", "", "SSH keypair name")
	nksCreateCmd.Flags().String("flavor-id", "", "Node flavor ID")
	nksCreateCmd.Flags().Int("node-count", 1, "Number of nodes")
	nksCreateCmd.MarkFlagRequired("name")
	nksCreateCmd.MarkFlagRequired("network-id")
	nksCreateCmd.MarkFlagRequired("subnet-id")

	nksNodeGroupCreateCmd.Flags().String("cluster-id", "", "Cluster ID (required)")
	nksNodeGroupCreateCmd.Flags().String("name", "", "Node group name (required)")
	nksNodeGroupCreateCmd.Flags().String("flavor-id", "", "Flavor ID (required)")
	nksNodeGroupCreateCmd.Flags().Int("node-count", 1, "Number of nodes")
	nksNodeGroupCreateCmd.MarkFlagRequired("cluster-id")
	nksNodeGroupCreateCmd.MarkFlagRequired("name")
	nksNodeGroupCreateCmd.MarkFlagRequired("flavor-id")

	nksNodeGroupUpdateCmd.Flags().String("cluster-id", "", "Cluster ID (required)")
	nksNodeGroupUpdateCmd.Flags().Int("node-count", 0, "New node count")
	nksNodeGroupUpdateCmd.MarkFlagRequired("cluster-id")
}

func getNKSClient() *nks.Client {
	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	return nks.NewClient(getRegion(), creds, nil, debug)
}

var nksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all NKS clusters",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		result, err := client.ListClusters(ctx)
		if err != nil {
			exitWithError("Failed to list clusters", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tK8S_VERSION\tNODE_COUNT\tCREATED")
		for _, c := range result.Clusters {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\n",
				c.ID, c.Name, c.Status, c.K8sVersion, c.NodeCount, c.CreatedAt)
		}
		w.Flush()
	},
}

var nksGetCmd = &cobra.Command{
	Use:   "get [cluster-id]",
	Short: "Get NKS cluster details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		result, err := client.GetCluster(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get cluster", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ID:           %s\n", result.ID)
		fmt.Printf("Name:         %s\n", result.Name)
		fmt.Printf("Status:       %s\n", result.Status)
		fmt.Printf("K8s Version:  %s\n", result.K8sVersion)
		fmt.Printf("Node Count:   %d\n", result.NodeCount)
		fmt.Printf("Master Count: %d\n", result.MasterCount)
		fmt.Printf("Network ID:   %s\n", result.NetworkID)
		fmt.Printf("Subnet ID:    %s\n", result.SubnetID)
		fmt.Printf("API Address:  %s\n", result.APIAddress)
		fmt.Printf("Created:      %s\n", result.CreatedAt)
		fmt.Printf("Updated:      %s\n", result.UpdatedAt)
	},
}

var nksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new NKS cluster",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		templateID, _ := cmd.Flags().GetString("template-id")
		k8sVersion, _ := cmd.Flags().GetString("k8s-version")
		networkID, _ := cmd.Flags().GetString("network-id")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		keypair, _ := cmd.Flags().GetString("keypair")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		nodeCount, _ := cmd.Flags().GetInt("node-count")

		input := &nks.CreateClusterInput{
			Name:              name,
			ClusterTemplateID: templateID,
			K8sVersion:        k8sVersion,
			NetworkID:         networkID,
			SubnetID:          subnetID,
			KeyPair:           keypair,
			FlavorID:          flavorID,
			NodeCount:         nodeCount,
			MasterCount:       1,
			Labels:            map[string]string{},
		}

		result, err := client.CreateCluster(ctx, input)
		if err != nil {
			exitWithError("Failed to create cluster", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Cluster creation initiated!\n")
		fmt.Printf("ID:   %s\n", result.ID)
		fmt.Printf("Name: %s\n", result.Name)
	},
}

var nksDeleteCmd = &cobra.Command{
	Use:   "delete [cluster-id]",
	Short: "Delete an NKS cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		if err := client.DeleteCluster(ctx, args[0]); err != nil {
			exitWithError("Failed to delete cluster", err)
		}

		fmt.Printf("Cluster %s deletion initiated\n", args[0])
	},
}

var nksKubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig [cluster-id]",
	Short: "Get kubeconfig for an NKS cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		result, err := client.GetKubeconfig(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get kubeconfig", err)
		}

		fmt.Println(result.Kubeconfig)
	},
}

var nksTemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List available cluster templates",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		result, err := client.ListClusterTemplates(ctx)
		if err != nil {
			exitWithError("Failed to list cluster templates", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tCOE\tPUBLIC")
		for _, t := range result.ClusterTemplates {
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n",
				t.ID, t.Name, t.COE, t.Public)
		}
		w.Flush()
	},
}

var nksNodeGroupsCmd = &cobra.Command{
	Use:   "node-groups [cluster-id]",
	Short: "List node groups in a cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		result, err := client.ListNodeGroups(ctx, args[0])
		if err != nil {
			exitWithError("Failed to list node groups", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tNODE_COUNT\tFLAVOR")
		for _, ng := range result.NodeGroups {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				ng.ID, ng.Name, ng.Status, ng.NodeCount, ng.FlavorID)
		}
		w.Flush()
	},
}

var nksNodeGroupGetCmd = &cobra.Command{
	Use:   "node-group-get [cluster-id] [node-group-id]",
	Short: "Get node group details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		result, err := client.GetNodeGroup(ctx, args[0], args[1])
		if err != nil {
			exitWithError("Failed to get node group", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ID:         %s\n", result.ID)
		fmt.Printf("Name:       %s\n", result.Name)
		fmt.Printf("Status:     %s\n", result.Status)
		fmt.Printf("Node Count: %d\n", result.NodeCount)
		fmt.Printf("Flavor:     %s\n", result.FlavorID)
		fmt.Printf("Created:    %s\n", result.CreatedAt)
	},
}

var nksNodeGroupCreateCmd = &cobra.Command{
	Use:   "node-group-create",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Node group creation initiated!\n")
		fmt.Printf("ID:   %s\n", result.ID)
		fmt.Printf("Name: %s\n", result.Name)
	},
}

var nksNodeGroupUpdateCmd = &cobra.Command{
	Use:   "node-group-update [node-group-id]",
	Short: "Update a node group (scale)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		clusterID, _ := cmd.Flags().GetString("cluster-id")
		nodeCount, _ := cmd.Flags().GetInt("node-count")

		input := &nks.UpdateNodeGroupInput{
			NodeCount: nodeCount,
		}

		if err := client.UpdateNodeGroup(ctx, clusterID, args[0], input); err != nil {
			exitWithError("Failed to update node group", err)
		}

		fmt.Printf("Node group %s update initiated\n", args[0])
	},
}

var nksNodeGroupDeleteCmd = &cobra.Command{
	Use:   "node-group-delete [cluster-id] [node-group-id]",
	Short: "Delete a node group",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		if err := client.DeleteNodeGroup(ctx, args[0], args[1]); err != nil {
			exitWithError("Failed to delete node group", err)
		}

		fmt.Printf("Node group %s deletion initiated\n", args[1])
	},
}
