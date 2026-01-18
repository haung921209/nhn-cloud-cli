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
	nksCmd.AddCommand(nksDescribeClustersCmd)
	nksCmd.AddCommand(nksCreateClusterCmd)
	nksCmd.AddCommand(nksDeleteClusterCmd)
	nksCmd.AddCommand(nksUpdateKubeconfigCmd)
	nksCmd.AddCommand(nksDescribeClusterTemplatesCmd)

	nksDescribeClustersCmd.Flags().String("cluster-id", "", "Cluster ID")

	nksCreateClusterCmd.Flags().String("name", "", "Cluster name (required)")
	nksCreateClusterCmd.Flags().String("template-id", "", "Cluster template ID (optional)")
	nksCreateClusterCmd.Flags().String("k8s-version", "", "Kubernetes version")
	nksCreateClusterCmd.Flags().String("network-id", "", "Network ID (required)")
	nksCreateClusterCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	nksCreateClusterCmd.Flags().String("keypair", "", "SSH keypair name")
	nksCreateClusterCmd.Flags().String("flavor-id", "", "Node flavor ID")
	nksCreateClusterCmd.Flags().Int("node-count", 1, "Number of nodes")
	nksCreateClusterCmd.MarkFlagRequired("name")
	nksCreateClusterCmd.MarkFlagRequired("network-id")
	nksCreateClusterCmd.MarkFlagRequired("subnet-id")

	nksDeleteClusterCmd.Flags().String("cluster-id", "", "Cluster ID (required)")
	nksDeleteClusterCmd.MarkFlagRequired("cluster-id")

	nksUpdateKubeconfigCmd.Flags().String("cluster-id", "", "Cluster ID (required)")
	nksUpdateKubeconfigCmd.MarkFlagRequired("cluster-id")
}

var nksDescribeClustersCmd = &cobra.Command{
	Use:   "describe-clusters",
	Short: "Describe NKS clusters",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("cluster-id")

		if id != "" {
			result, err := client.GetCluster(ctx, id)
			if err != nil {
				exitWithError("Failed to get cluster", err)
			}
			if output == "json" {
				printJSON(result)
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
		} else {
			result, err := client.ListClusters(ctx)
			if err != nil {
				exitWithError("Failed to list clusters", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tSTATUS\tK8S_VERSION\tNODE_COUNT\tCREATED")
			for _, c := range result.Clusters {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\n",
					c.ID, c.Name, c.Status, c.K8sVersion, c.NodeCount, c.CreatedAt)
			}
			w.Flush()
		}
	},
}

var nksCreateClusterCmd = &cobra.Command{
	Use:   "create-cluster",
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
			printJSON(result)
			return
		}

		fmt.Printf("Cluster creation initiated!\n")
		fmt.Printf("ID:   %s\n", result.ID)
		fmt.Printf("Name: %s\n", result.Name)
	},
}

var nksDeleteClusterCmd = &cobra.Command{
	Use:   "delete-cluster",
	Short: "Delete an NKS cluster",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("cluster-id")

		if err := client.DeleteCluster(ctx, id); err != nil {
			exitWithError("Failed to delete cluster", err)
		}

		fmt.Printf("Cluster %s deletion initiated\n", id)
	},
}

var nksUpdateKubeconfigCmd = &cobra.Command{
	Use:   "update-kubeconfig",
	Short: "Get kubeconfig for an NKS cluster",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("cluster-id")

		result, err := client.GetKubeconfig(ctx, id)
		if err != nil {
			exitWithError("Failed to get kubeconfig", err)
		}

		fmt.Println(result.Kubeconfig)
	},
}

var nksDescribeClusterTemplatesCmd = &cobra.Command{
	Use:   "describe-cluster-templates",
	Short: "List available cluster templates",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNKSClient()
		ctx := context.Background()

		result, err := client.ListClusterTemplates(ctx)
		if err != nil {
			exitWithError("Failed to list cluster templates", err)
		}

		if output == "json" {
			printJSON(result)
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
