package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncs"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/spf13/cobra"
)

var ncsCmd = &cobra.Command{
	Use:     "ncs",
	Aliases: []string{"container-service"},
	Short:   "Manage NHN Container Service (NCS)",
	Long:    `Manage serverless container workloads and services.`,
}

func init() {
	rootCmd.AddCommand(ncsCmd)

	ncsCmd.AddCommand(ncsWorkloadsCmd)
	ncsCmd.AddCommand(ncsWorkloadGetCmd)
	ncsCmd.AddCommand(ncsWorkloadCreateCmd)
	ncsCmd.AddCommand(ncsWorkloadDeleteCmd)
	ncsCmd.AddCommand(ncsWorkloadRestartCmd)
	ncsCmd.AddCommand(ncsWorkloadScaleCmd)

	ncsCmd.AddCommand(ncsTemplatesCmd)
	ncsCmd.AddCommand(ncsTemplateGetCmd)

	ncsCmd.AddCommand(ncsServicesCmd)
	ncsCmd.AddCommand(ncsServiceGetCmd)
	ncsCmd.AddCommand(ncsServiceCreateCmd)
	ncsCmd.AddCommand(ncsServiceDeleteCmd)

	ncsWorkloadsCmd.Flags().String("namespace", "", "Filter by namespace")

	ncsWorkloadCreateCmd.Flags().String("name", "", "Workload name (required)")
	ncsWorkloadCreateCmd.Flags().String("namespace", "default", "Namespace")
	ncsWorkloadCreateCmd.Flags().String("image", "", "Container image (required)")
	ncsWorkloadCreateCmd.Flags().Int("replicas", 1, "Number of replicas")
	ncsWorkloadCreateCmd.Flags().String("cpu", "1", "CPU request (e.g., 1, 500m)")
	ncsWorkloadCreateCmd.Flags().String("memory", "2Gi", "Memory request (e.g., 2Gi, 512Mi)")
	ncsWorkloadCreateCmd.Flags().Int("port", 0, "Container port")
	ncsWorkloadCreateCmd.MarkFlagRequired("name")
	ncsWorkloadCreateCmd.MarkFlagRequired("image")

	ncsWorkloadScaleCmd.Flags().Int("replicas", 1, "Number of replicas")
	ncsWorkloadScaleCmd.MarkFlagRequired("replicas")

	ncsServicesCmd.Flags().String("namespace", "", "Filter by namespace")

	ncsServiceCreateCmd.Flags().String("name", "", "Service name (required)")
	ncsServiceCreateCmd.Flags().String("namespace", "default", "Namespace")
	ncsServiceCreateCmd.Flags().String("selector", "", "Label selector (key=value)")
	ncsServiceCreateCmd.Flags().Int("port", 80, "Service port")
	ncsServiceCreateCmd.Flags().Int("target-port", 0, "Target container port")
	ncsServiceCreateCmd.Flags().String("type", "LoadBalancer", "Service type (ClusterIP, LoadBalancer)")
	ncsServiceCreateCmd.MarkFlagRequired("name")
}

func getNCSClient() *ncs.Client {
	creds := credentials.NewStatic(getAccessKey(), getSecretKey())
	return ncs.NewClient(getRegion(), getAppKey(), creds, nil, debug)
}

var ncsWorkloadsCmd = &cobra.Command{
	Use:   "workloads",
	Short: "List all workloads",
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

var ncsWorkloadGetCmd = &cobra.Command{
	Use:   "workload-get [workload-id]",
	Short: "Get workload details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		result, err := client.GetWorkload(ctx, args[0])
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

var ncsWorkloadCreateCmd = &cobra.Command{
	Use:   "workload-create",
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
			},
		}

		if port > 0 {
			container.Ports = []ncs.ContainerPort{{ContainerPort: port}}
		}

		input := &ncs.CreateWorkloadInput{
			Name:       name,
			Namespace:  namespace,
			Containers: []ncs.Container{container},
			Replicas:   replicas,
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
		fmt.Printf("ID:        %s\n", result.ID)
		fmt.Printf("Name:      %s\n", result.Name)
		fmt.Printf("Namespace: %s\n", result.Namespace)
	},
}

var ncsWorkloadDeleteCmd = &cobra.Command{
	Use:   "workload-delete [workload-id]",
	Short: "Delete a workload",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		if err := client.DeleteWorkload(ctx, args[0]); err != nil {
			exitWithError("Failed to delete workload", err)
		}

		fmt.Printf("Workload %s deleted successfully\n", args[0])
	},
}

var ncsWorkloadRestartCmd = &cobra.Command{
	Use:   "workload-restart [workload-id]",
	Short: "Restart a workload",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		if err := client.RestartWorkload(ctx, args[0]); err != nil {
			exitWithError("Failed to restart workload", err)
		}

		fmt.Printf("Workload %s restart initiated\n", args[0])
	},
}

var ncsWorkloadScaleCmd = &cobra.Command{
	Use:   "workload-scale [workload-id]",
	Short: "Scale a workload",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		replicas, _ := cmd.Flags().GetInt("replicas")

		if err := client.ScaleWorkload(ctx, args[0], replicas); err != nil {
			exitWithError("Failed to scale workload", err)
		}

		fmt.Printf("Workload %s scaled to %d replicas\n", args[0], replicas)
	},
}

var ncsTemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List available templates",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		result, err := client.ListTemplates(ctx)
		if err != nil {
			exitWithError("Failed to list templates", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tVERSION\tDESCRIPTION")
		for _, t := range result.Templates {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				t.ID, t.Name, t.Version, t.Description)
		}
		w.Flush()
	},
}

var ncsTemplateGetCmd = &cobra.Command{
	Use:   "template-get [template-id]",
	Short: "Get template details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		result, err := client.GetTemplate(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get template", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ID:          %s\n", result.ID)
		fmt.Printf("Name:        %s\n", result.Name)
		fmt.Printf("Version:     %s\n", result.Version)
		fmt.Printf("Description: %s\n", result.Description)
	},
}

var ncsServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List all services",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		namespace, _ := cmd.Flags().GetString("namespace")

		result, err := client.ListServices(ctx, namespace)
		if err != nil {
			exitWithError("Failed to list services", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tNAMESPACE\tTYPE\tCLUSTER_IP\tEXTERNAL_IP")
		for _, svc := range result.Services {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				svc.ID, svc.Name, svc.Namespace, svc.Type, svc.ClusterIP, svc.ExternalIP)
		}
		w.Flush()
	},
}

var ncsServiceGetCmd = &cobra.Command{
	Use:   "service-get [service-id]",
	Short: "Get service details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		result, err := client.GetService(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get service", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ID:          %s\n", result.ID)
		fmt.Printf("Name:        %s\n", result.Name)
		fmt.Printf("Namespace:   %s\n", result.Namespace)
		fmt.Printf("Type:        %s\n", result.Type)
		fmt.Printf("Cluster IP:  %s\n", result.ClusterIP)
		fmt.Printf("External IP: %s\n", result.ExternalIP)
		fmt.Printf("Created:     %s\n", result.CreatedAt)
		if len(result.Ports) > 0 {
			fmt.Printf("\nPorts:\n")
			for _, p := range result.Ports {
				fmt.Printf("  - Port: %d -> Target: %d (%s)\n", p.Port, p.TargetPort, p.Protocol)
			}
		}
	},
}

var ncsServiceCreateCmd = &cobra.Command{
	Use:   "service-create",
	Short: "Create a new service",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		namespace, _ := cmd.Flags().GetString("namespace")
		selectorStr, _ := cmd.Flags().GetString("selector")
		port, _ := cmd.Flags().GetInt("port")
		targetPort, _ := cmd.Flags().GetInt("target-port")
		svcType, _ := cmd.Flags().GetString("type")

		if targetPort == 0 {
			targetPort = port
		}

		input := &ncs.CreateServiceInput{
			Name:      name,
			Namespace: namespace,
			Type:      svcType,
			Ports: []ncs.ServicePort{
				{
					Port:       port,
					TargetPort: targetPort,
					Protocol:   "TCP",
				},
			},
		}

		if selectorStr != "" {
			input.Selector = parseSelector(selectorStr)
		}

		result, err := client.CreateService(ctx, input)
		if err != nil {
			exitWithError("Failed to create service", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Service created successfully!\n")
		fmt.Printf("ID:          %s\n", result.ID)
		fmt.Printf("Name:        %s\n", result.Name)
		fmt.Printf("External IP: %s\n", result.ExternalIP)
	},
}

var ncsServiceDeleteCmd = &cobra.Command{
	Use:   "service-delete [service-id]",
	Short: "Delete a service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		if err := client.DeleteService(ctx, args[0]); err != nil {
			exitWithError("Failed to delete service", err)
		}

		fmt.Printf("Service %s deleted successfully\n", args[0])
	},
}

func parseSelector(s string) map[string]string {
	result := make(map[string]string)
	if s == "" {
		return result
	}
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			result[s[:i]] = s[i+1:]
			break
		}
	}
	return result
}
