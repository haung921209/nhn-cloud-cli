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
	ncsCreateWorkloadCmd.Flags().String("template-id", "", "Template ID; the template holds the containers and network (required)")
	ncsCreateWorkloadCmd.Flags().String("template-version", "", "Template version (default: latest)")
	ncsCreateWorkloadCmd.Flags().Int("replicas", 1, "Number of tasks (desired)")
	ncsCreateWorkloadCmd.Flags().String("type", "", "Deployment controller: deployment (default) or statefulset")
	ncsCreateWorkloadCmd.Flags().String("description", "", "Workload description")
	ncsCreateWorkloadCmd.Flags().Bool("load-balancer", false, "Put a load balancer in front of the workload")
	ncsCreateWorkloadCmd.Flags().Bool("floating-ip", false, "Give the load balancer a floating IP")
	ncsCreateWorkloadCmd.MarkFlagRequired("name")
	ncsCreateWorkloadCmd.MarkFlagRequired("template-id")

	ncsGetWorkloadCmd.Flags().String("workload-id", "", "Workload ID (or pass it as an argument)")

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
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tSTATUS\tAVAILABLE\tTEMPLATE\tCREATED")
		for _, wl := range result.Workloads {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d/%d\t%s\t%s\n",
				wl.ID, wl.Name, wl.Type, wl.Status,
				wl.Available, wl.Desired, wl.TemplateID, wl.CreatedAt)
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

		if structuredOutput() {
			printOutput(result)
			return
		}

		wl := result.Workload
		fmt.Printf("ID:        %s\n", wl.ID)
		fmt.Printf("Name:      %s\n", wl.Name)
		fmt.Printf("Type:      %s\n", wl.Type)
		fmt.Printf("Status:    %s\n", wl.Status)
		fmt.Printf("Available: %d/%d\n", wl.Available, wl.Desired)
		fmt.Printf("Template:  %s (%s)\n", wl.TemplateID, wl.TemplateVersion)
		fmt.Printf("Created:   %s\n", wl.CreatedAt)
	},
}

var ncsCreateWorkloadCmd = &cobra.Command{
	Use:   "create-workload",
	Short: "Create a new workload",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		templateID, _ := cmd.Flags().GetString("template-id")
		templateVersion, _ := cmd.Flags().GetString("template-version")
		replicas, _ := cmd.Flags().GetInt("replicas")
		typ, _ := cmd.Flags().GetString("type")
		description, _ := cmd.Flags().GetString("description")
		lb, _ := cmd.Flags().GetBool("load-balancer")
		fip, _ := cmd.Flags().GetBool("floating-ip")

		result, err := client.CreateWorkload(ctx, &ncs.CreateWorkloadInput{
			Name:            name,
			Description:     description,
			Type:            typ,
			TemplateID:      templateID,
			TemplateVersion: templateVersion,
			Desired:         replicas,
			LoadBalancing:   &ncs.LoadBalancing{Enabled: lb, FloatingIP: fip},
		})
		if err != nil {
			exitWithError("Failed to create workload", err)
		}

		printWriteResult(result, "Workload created successfully!\nID:   %s\nName: %s\n", result.Workload.ID, result.Workload.Name)
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
