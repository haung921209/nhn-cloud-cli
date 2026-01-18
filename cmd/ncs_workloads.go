package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncs"
	"github.com/spf13/cobra"
)

func init() {
	ncsCmd.AddCommand(ncsDescribeWorkloadsCmd)
	ncsCmd.AddCommand(ncsGetWorkloadCmd)
	ncsCmd.AddCommand(ncsCreateWorkloadCmd)
	ncsCmd.AddCommand(ncsDeleteWorkloadCmd)

	ncsDescribeWorkloadsCmd.Flags().String("namespace", "", "Filter by namespace")

	ncsCreateWorkloadCmd.Flags().String("name", "", "Workload name (required)")
	ncsCreateWorkloadCmd.Flags().String("namespace", "default", "Namespace")
	ncsCreateWorkloadCmd.Flags().String("image", "", "Container image (required)")
	ncsCreateWorkloadCmd.Flags().Int("replicas", 1, "Number of replicas")
	ncsCreateWorkloadCmd.Flags().String("cpu", "1", "CPU request (e.g., 1, 500m)")
	ncsCreateWorkloadCmd.Flags().String("memory", "2Gi", "Memory request (e.g., 2Gi, 512Mi)")
	ncsCreateWorkloadCmd.Flags().Int("port", 0, "Container port")
	ncsCreateWorkloadCmd.MarkFlagRequired("name")
	ncsCreateWorkloadCmd.MarkFlagRequired("image")

	ncsDeleteWorkloadCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsDeleteWorkloadCmd.MarkFlagRequired("workload-id")
}

var ncsDescribeWorkloadsCmd = &cobra.Command{
	Use:     "describe-workloads",
	Aliases: []string{"list-workloads", "workloads"},
	Short:   "List all workloads",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		namespace, _ := cmd.Flags().GetString("namespace")

		result, err := client.ListWorkloads(ctx, namespace)
		if err != nil {
			exitWithError("Failed to list workloads", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tNAMESPACE\tSTATUS\tREPLICAS\tCREATED")
		for _, wl := range result.Workloads {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d/%d\t%s\n",
				wl.ID, wl.Name, wl.Namespace, wl.Status,
				wl.AvailableReplicas, wl.Replicas, wl.CreatedAt)
		}
		w.Flush()
	},
}

var ncsGetWorkloadCmd = &cobra.Command{
	Use:     "describe-workload",
	Aliases: []string{"get-workload", "workload-get"},
	Short:   "Get workload details",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")
		// Support positional arg if flag not set
		if workloadID == "" && len(args) > 0 {
			workloadID = args[0]
		}
		if workloadID == "" {
			exitWithError("Workload ID required (use --workload-id or positional argument)", nil)
		}

		result, err := client.GetWorkload(ctx, workloadID)
		if err != nil {
			exitWithError("Failed to get workload", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ID:        %s\n", result.ID)
		fmt.Printf("Name:      %s\n", result.Name)
		fmt.Printf("Namespace: %s\n", result.Namespace)
		fmt.Printf("Status:    %s\n", result.Status)
		fmt.Printf("Replicas:  %d/%d\n", result.AvailableReplicas, result.Replicas)
		fmt.Printf("Created:   %s\n", result.CreatedAt)
		if len(result.Containers) > 0 {
			fmt.Printf("\nContainers:\n")
			for _, c := range result.Containers {
				fmt.Printf("  - Name:  %s\n", c.Name)
				fmt.Printf("    Image: %s\n", c.Image)
				if c.Resources != nil {
					if c.Resources.Requests.CPU != "" {
						fmt.Printf("    CPU:   %s\n", c.Resources.Requests.CPU)
					}
					if c.Resources.Requests.Memory != "" {
						fmt.Printf("    Memory: %s\n", c.Resources.Requests.Memory)
					}
				}
			}
		}
	},
}

var ncsCreateWorkloadCmd = &cobra.Command{
	Use:   "create-workload",
	Short: "Create a new workload",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		namespace, _ := cmd.Flags().GetString("namespace")
		image, _ := cmd.Flags().GetString("image")
		replicas, _ := cmd.Flags().GetInt("replicas")
		cpu, _ := cmd.Flags().GetString("cpu")
		memory, _ := cmd.Flags().GetString("memory")
		port, _ := cmd.Flags().GetInt("port")

		container := ncs.Container{
			Name:  name,
			Image: image,
			Resources: &ncs.ResourceRequirements{
				Requests: ncs.ResourceList{
					CPU:    cpu,
					Memory: memory,
				},
				Limits: ncs.ResourceList{
					CPU:    cpu,
					Memory: memory,
				},
			},
		}

		if port > 0 {
			container.Ports = []ncs.ContainerPort{
				{
					ContainerPort: port,
					Protocol:      "TCP",
				},
			}
		}

		input := &ncs.CreateWorkloadInput{
			Name:       name,
			Namespace:  namespace,
			Replicas:   replicas,
			Containers: []ncs.Container{container},
		}

		result, err := client.CreateWorkload(ctx, input)
		if err != nil {
			exitWithError("Failed to create workload", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Workload created successfully!\n")
		fmt.Printf("ID:   %s\n", result.ID)
		fmt.Printf("Name: %s\n", result.Name)
	},
}

var ncsDeleteWorkloadCmd = &cobra.Command{
	Use:   "delete-workload",
	Short: "Delete a workload",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")
		// Support positional arg
		if workloadID == "" && len(args) > 0 {
			workloadID = args[0]
		}
		if workloadID == "" {
			exitWithError("Workload ID required", nil)
		}

		if err := client.DeleteWorkload(ctx, workloadID); err != nil {
			exitWithError("Failed to delete workload", err)
		}

		fmt.Printf("Workload %s deleted successfully\n", workloadID)
	},
}
