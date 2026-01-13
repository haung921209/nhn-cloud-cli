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

	ncsCmd.AddCommand(ncsLogsCmd)
	ncsCmd.AddCommand(ncsHealthCheckConfigureCmd)
	ncsCmd.AddCommand(ncsHealthCheckStatusCmd)
	ncsCmd.AddCommand(ncsResourcesUpdateCmd)
	ncsCmd.AddCommand(ncsEventsCmd)
	ncsCmd.AddCommand(ncsVolumesCmd)
	ncsCmd.AddCommand(ncsVolumeAttachCmd)
	ncsCmd.AddCommand(ncsExecCmd)
	ncsCmd.AddCommand(ncsContainerStatusCmd)
	ncsCmd.AddCommand(ncsAutoScalingCmd)
	ncsAutoScalingCmd.AddCommand(ncsAutoScalingConfigureCmd)
	ncsAutoScalingCmd.AddCommand(ncsAutoScalingStatusCmd)

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

	ncsLogsCmd.Flags().Int("tail", 100, "Number of lines to show from end of logs")
	ncsLogsCmd.Flags().Int("since", 0, "Show logs since N seconds ago")

	ncsHealthCheckConfigureCmd.Flags().String("liveness-path", "", "HTTP path for liveness probe")
	ncsHealthCheckConfigureCmd.Flags().String("readiness-path", "", "HTTP path for readiness probe")
	ncsHealthCheckConfigureCmd.Flags().Int("port", 8080, "Port for health check")
	ncsHealthCheckConfigureCmd.Flags().Int("initial-delay", 30, "Initial delay in seconds")

	ncsResourcesUpdateCmd.Flags().String("cpu-limit", "", "CPU limit (e.g., 2, 500m)")
	ncsResourcesUpdateCmd.Flags().String("memory-limit", "", "Memory limit (e.g., 4Gi, 512Mi)")
	ncsResourcesUpdateCmd.Flags().String("cpu-request", "", "CPU request (e.g., 1, 250m)")
	ncsResourcesUpdateCmd.Flags().String("memory-request", "", "Memory request (e.g., 2Gi, 256Mi)")

	ncsVolumeAttachCmd.Flags().String("mount-path", "", "Mount path in container (required)")
	ncsVolumeAttachCmd.Flags().Bool("read-only", false, "Mount as read-only")
	ncsVolumeAttachCmd.MarkFlagRequired("mount-path")

	ncsExecCmd.Flags().String("container", "", "Container name (if workload has multiple containers)")
	ncsExecCmd.Flags().BoolP("stdin", "i", false, "Pass stdin to the container")
	ncsExecCmd.Flags().BoolP("tty", "t", false, "Allocate a pseudo-TTY")

	ncsAutoScalingConfigureCmd.Flags().Bool("enabled", false, "Enable or disable auto-scaling")
	ncsAutoScalingConfigureCmd.Flags().Int("min-replicas", 1, "Minimum number of replicas")
	ncsAutoScalingConfigureCmd.Flags().Int("max-replicas", 10, "Maximum number of replicas")
	ncsAutoScalingConfigureCmd.Flags().Int("target-cpu", 80, "Target CPU utilization percentage")
	ncsAutoScalingConfigureCmd.Flags().Int("target-memory", 0, "Target memory utilization percentage (optional)")
}

func getNCSClient() *ncs.Client {
	creds := credentials.NewStatic(getAccessKey(), getSecretKey())
	return ncs.NewClient(getRegion(), getNCSAppKey(), creds, nil, debug)
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

		if svcType != "ClusterIP" && svcType != "LoadBalancer" {
			exitWithError("service type must be ClusterIP or LoadBalancer", nil)
		}

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

var ncsLogsCmd = &cobra.Command{
	Use:   "logs [workload-id]",
	Short: "Get workload logs",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		tail, _ := cmd.Flags().GetInt("tail")
		since, _ := cmd.Flags().GetInt("since")

		result, err := client.GetWorkloadLogs(ctx, args[0], tail, since)
		if err != nil {
			exitWithError("Failed to get logs", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		for _, log := range result.Logs {
			fmt.Printf("[%s] %s: %s\n", log.Timestamp, log.Stream, log.Message)
		}
	},
}

var ncsHealthCheckConfigureCmd = &cobra.Command{
	Use:   "health-check-configure [workload-id]",
	Short: "Configure health checks for a workload",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		livenessPath, _ := cmd.Flags().GetString("liveness-path")
		readinessPath, _ := cmd.Flags().GetString("readiness-path")
		port, _ := cmd.Flags().GetInt("port")
		initialDelay, _ := cmd.Flags().GetInt("initial-delay")

		config := &ncs.HealthCheckConfig{}

		if livenessPath != "" {
			config.LivenessProbe = &ncs.HealthProbe{
				Type: "HTTP",
				HTTPGet: &ncs.HTTPGetAction{
					Path: livenessPath,
					Port: port,
				},
				InitialDelaySeconds: initialDelay,
				PeriodSeconds:       10,
			}
		}

		if readinessPath != "" {
			config.ReadinessProbe = &ncs.HealthProbe{
				Type: "HTTP",
				HTTPGet: &ncs.HTTPGetAction{
					Path: readinessPath,
					Port: port,
				},
				InitialDelaySeconds: initialDelay,
				PeriodSeconds:       5,
			}
		}

		if err := client.ConfigureHealthCheck(ctx, args[0], config); err != nil {
			exitWithError("Failed to configure health check", err)
		}

		fmt.Printf("Health check configured successfully for workload %s\n", args[0])
	},
}

var ncsHealthCheckStatusCmd = &cobra.Command{
	Use:   "health-check-status [workload-id]",
	Short: "Get health check status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		result, err := client.GetHealthCheckStatus(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get health check status", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CONTAINER\tLIVENESS\tREADINESS\tLAST_CHECK")
		for _, status := range result.Status {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				status.ContainerName, status.Liveness, status.Readiness, status.LastCheck)
		}
		w.Flush()
	},
}

var ncsResourcesUpdateCmd = &cobra.Command{
	Use:   "resources-update [workload-id]",
	Short: "Update resource limits for a workload",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		cpuLimit, _ := cmd.Flags().GetString("cpu-limit")
		memoryLimit, _ := cmd.Flags().GetString("memory-limit")
		cpuRequest, _ := cmd.Flags().GetString("cpu-request")
		memoryRequest, _ := cmd.Flags().GetString("memory-request")

		input := &ncs.UpdateResourceLimitsInput{
			Resources: &ncs.ResourceRequirements{},
		}

		if cpuLimit != "" || memoryLimit != "" {
			input.Resources.Limits = ncs.ResourceList{
				CPU:    cpuLimit,
				Memory: memoryLimit,
			}
		}

		if cpuRequest != "" || memoryRequest != "" {
			input.Resources.Requests = ncs.ResourceList{
				CPU:    cpuRequest,
				Memory: memoryRequest,
			}
		}

		if err := client.UpdateResourceLimits(ctx, args[0], input); err != nil {
			exitWithError("Failed to update resource limits", err)
		}

		fmt.Printf("Resource limits updated successfully for workload %s\n", args[0])
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

var ncsEventsCmd = &cobra.Command{
	Use:   "events [workload-id]",
	Short: "Get workload events for debugging",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		result, err := client.GetWorkloadEvents(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get workload events", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		if len(result.Events) == 0 {
			fmt.Println("No events found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "LAST SEEN\tTYPE\tREASON\tCOUNT\tMESSAGE")
		for _, event := range result.Events {
			msg := event.Message
			if len(msg) > 60 {
				msg = msg[:57] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				event.LastTime, event.EventType, event.Reason, event.Count, msg)
		}
		w.Flush()
	},
}

var ncsVolumesCmd = &cobra.Command{
	Use:   "volumes",
	Short: "List all persistent volumes",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		result, err := client.ListVolumes(ctx)
		if err != nil {
			exitWithError("Failed to list volumes", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		if len(result.Volumes) == 0 {
			fmt.Println("No volumes found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSIZE(GB)\tSTATUS\tTYPE\tATTACHED TO")
		for _, vol := range result.Volumes {
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\n",
				vol.VolumeID, vol.Name, vol.Size, vol.Status, vol.VolumeType, vol.AttachedTo)
		}
		w.Flush()
	},
}

var ncsVolumeAttachCmd = &cobra.Command{
	Use:   "volume-attach [workload-id] [volume-id]",
	Short: "Attach a volume to a workload",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		mountPath, _ := cmd.Flags().GetString("mount-path")
		readOnly, _ := cmd.Flags().GetBool("read-only")

		if mountPath == "" {
			exitWithError("--mount-path is required", nil)
		}

		input := &ncs.VolumeAttachInput{
			VolumeID:  args[1],
			MountPath: mountPath,
			ReadOnly:  readOnly,
		}

		_, err := client.AttachVolume(ctx, args[0], input)
		if err != nil {
			exitWithError("Failed to attach volume", err)
		}

		fmt.Printf("Volume %s attached to workload %s at %s\n", args[1], args[0], mountPath)
	},
}

var ncsExecCmd = &cobra.Command{
	Use:   "exec [workload-id] -- [command...]",
	Short: "Execute command in a workload container",
	Long: `Execute a command in a workload container.
	
Example:
  nhncloud ncs exec my-workload -- ls -la
  nhncloud ncs exec my-workload --container sidecar -- /bin/sh
  nhncloud ncs exec my-workload -it -- /bin/bash`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		workloadID := args[0]
		container, _ := cmd.Flags().GetString("container")
		stdin, _ := cmd.Flags().GetBool("stdin")
		tty, _ := cmd.Flags().GetBool("tty")

		dashIdx := -1
		for i, arg := range args {
			if arg == "--" {
				dashIdx = i
				break
			}
		}

		var command []string
		if dashIdx >= 0 {
			command = args[dashIdx+1:]
		} else {
			command = args[1:]
		}

		if len(command) == 0 {
			exitWithError("No command specified", nil)
		}

		input := &ncs.ExecInput{
			ContainerName: container,
			Command:       command,
			Stdin:         stdin,
			Stdout:        true,
			Stderr:        true,
			TTY:           tty,
		}

		result, err := client.ExecWorkloadContainer(ctx, workloadID, input)
		if err != nil {
			exitWithError("Failed to execute command", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			return
		}

		if result.Output != "" {
			fmt.Print(result.Output)
		}

		if result.Error != "" {
			fmt.Fprintf(os.Stderr, "Error: %s\n", result.Error)
		}

		if result.ExitCode != 0 {
			os.Exit(result.ExitCode)
		}
	},
}

var ncsContainerStatusCmd = &cobra.Command{
	Use:   "container-status [workload-id]",
	Short: "Get container runtime status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		result, err := client.GetContainerStatus(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get container status", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			return
		}

		if len(result.Containers) == 0 {
			fmt.Println("No containers found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tSTATE\tREADY\tRESTARTS\tIMAGE")
		for _, container := range result.Containers {
			ready := "No"
			if container.Ready {
				ready = "Yes"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				container.ContainerName, container.State, ready, container.RestartCount, container.Image)
		}
		w.Flush()

		for _, container := range result.Containers {
			if container.Reason != "" || container.Message != "" {
				fmt.Printf("\nContainer '%s':\n", container.ContainerName)
				if container.Reason != "" {
					fmt.Printf("  Reason: %s\n", container.Reason)
				}
				if container.Message != "" {
					fmt.Printf("  Message: %s\n", container.Message)
				}
			}
		}
	},
}

var ncsAutoScalingCmd = &cobra.Command{
	Use:   "autoscaling",
	Short: "Manage workload auto-scaling",
}

var ncsAutoScalingConfigureCmd = &cobra.Command{
	Use:   "configure [workload-id]",
	Short: "Configure auto-scaling policy",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		enabled, _ := cmd.Flags().GetBool("enabled")
		minReplicas, _ := cmd.Flags().GetInt("min-replicas")
		maxReplicas, _ := cmd.Flags().GetInt("max-replicas")
		targetCPU, _ := cmd.Flags().GetInt("target-cpu")
		targetMemory, _ := cmd.Flags().GetInt("target-memory")

		if enabled && (minReplicas <= 0 || maxReplicas <= 0) {
			exitWithError("--min-replicas and --max-replicas are required when enabling auto-scaling", nil)
		}

		if minReplicas > maxReplicas {
			exitWithError("--min-replicas cannot be greater than --max-replicas", nil)
		}

		input := &ncs.ConfigureAutoScalingInput{
			Enabled: enabled,
		}

		if enabled {
			input.Policy = &ncs.AutoScalingPolicy{
				MinReplicas:                       minReplicas,
				MaxReplicas:                       maxReplicas,
				TargetCPUUtilizationPercentage:    targetCPU,
				TargetMemoryUtilizationPercentage: targetMemory,
			}
		}

		err := client.ConfigureAutoScaling(ctx, args[0], input)
		if err != nil {
			exitWithError("Failed to configure auto-scaling", err)
		}

		if enabled {
			fmt.Printf("Auto-scaling enabled for workload %s (min: %d, max: %d)\n", args[0], minReplicas, maxReplicas)
		} else {
			fmt.Printf("Auto-scaling disabled for workload %s\n", args[0])
		}
	},
}

var ncsAutoScalingStatusCmd = &cobra.Command{
	Use:   "status [workload-id]",
	Short: "Get auto-scaling status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		result, err := client.GetAutoScalingStatus(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get auto-scaling status", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			return
		}

		status := &result.Status
		fmt.Printf("Enabled:          %v\n", status.Enabled)
		fmt.Printf("Current Replicas: %d\n", status.CurrentReplicas)
		fmt.Printf("Desired Replicas: %d\n", status.DesiredReplicas)

		if status.Policy != nil {
			fmt.Printf("\nPolicy:\n")
			fmt.Printf("  Min Replicas:    %d\n", status.Policy.MinReplicas)
			fmt.Printf("  Max Replicas:    %d\n", status.Policy.MaxReplicas)
			if status.Policy.TargetCPUUtilizationPercentage > 0 {
				fmt.Printf("  Target CPU:      %d%%\n", status.Policy.TargetCPUUtilizationPercentage)
			}
			if status.Policy.TargetMemoryUtilizationPercentage > 0 {
				fmt.Printf("  Target Memory:   %d%%\n", status.Policy.TargetMemoryUtilizationPercentage)
			}
		}

		if status.LastScaleTime != "" {
			fmt.Printf("\nLast Scale Time:  %s\n", status.LastScaleTime)
		}

		if len(status.Conditions) > 0 {
			fmt.Printf("\nConditions:\n")
			for _, cond := range status.Conditions {
				fmt.Printf("  - Type: %s, Status: %s\n", cond.Type, cond.Status)
				if cond.Reason != "" {
					fmt.Printf("    Reason: %s\n", cond.Reason)
				}
				if cond.Message != "" {
					fmt.Printf("    Message: %s\n", cond.Message)
				}
			}
		}
	},
}
