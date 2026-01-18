package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncs"
	"github.com/spf13/cobra"
)

func init() {
	ncsCmd.AddCommand(ncsRestartWorkloadCmd)
	ncsCmd.AddCommand(ncsScaleWorkloadCmd)
	ncsCmd.AddCommand(ncsUpdateResourcesCmd)
	ncsCmd.AddCommand(ncsConfigureHealthCheckCmd)
	ncsCmd.AddCommand(ncsGetHealthCheckStatusCmd)

	ncsRestartWorkloadCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsRestartWorkloadCmd.MarkFlagRequired("workload-id")

	ncsScaleWorkloadCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsScaleWorkloadCmd.Flags().Int("replicas", 1, "Number of replicas")
	ncsScaleWorkloadCmd.MarkFlagRequired("workload-id")
	ncsScaleWorkloadCmd.MarkFlagRequired("replicas")

	ncsUpdateResourcesCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsUpdateResourcesCmd.Flags().String("cpu-limit", "", "CPU limit (e.g., 2, 500m)")
	ncsUpdateResourcesCmd.Flags().String("memory-limit", "", "Memory limit (e.g., 4Gi, 512Mi)")
	ncsUpdateResourcesCmd.Flags().String("cpu-request", "", "CPU request (e.g., 1, 250m)")
	ncsUpdateResourcesCmd.Flags().String("memory-request", "", "Memory request (e.g., 2Gi, 256Mi)")
	ncsUpdateResourcesCmd.MarkFlagRequired("workload-id")

	ncsConfigureHealthCheckCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsConfigureHealthCheckCmd.Flags().String("liveness-path", "", "HTTP path for liveness probe")
	ncsConfigureHealthCheckCmd.Flags().String("readiness-path", "", "HTTP path for readiness probe")
	ncsConfigureHealthCheckCmd.Flags().Int("port", 8080, "Port for health check")
	ncsConfigureHealthCheckCmd.Flags().Int("initial-delay", 30, "Initial delay in seconds")
	ncsConfigureHealthCheckCmd.MarkFlagRequired("workload-id")

	ncsGetHealthCheckStatusCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsGetHealthCheckStatusCmd.MarkFlagRequired("workload-id")
}

var ncsRestartWorkloadCmd = &cobra.Command{
	Use:   "restart-workload",
	Short: "Restart a workload",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")

		if err := client.RestartWorkload(ctx, workloadID); err != nil {
			exitWithError("Failed to restart workload", err)
		}

		fmt.Printf("Workload %s restarted successfully\n", workloadID)
	},
}

var ncsScaleWorkloadCmd = &cobra.Command{
	Use:   "scale-workload",
	Short: "Scale a workload",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")
		replicas, _ := cmd.Flags().GetInt("replicas")

		if err := client.ScaleWorkload(ctx, workloadID, replicas); err != nil {
			exitWithError("Failed to scale workload", err)
		}

		fmt.Printf("Workload %s scaled to %d replicas\n", workloadID, replicas)
	},
}

var ncsUpdateResourcesCmd = &cobra.Command{
	Use:     "update-workload-resources",
	Aliases: []string{"update-resources"},
	Short:   "Update workload resources",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")

		cpuLimit, _ := cmd.Flags().GetString("cpu-limit")
		memLimit, _ := cmd.Flags().GetString("memory-limit")
		cpuReq, _ := cmd.Flags().GetString("cpu-request")
		memReq, _ := cmd.Flags().GetString("memory-request")

		input := &ncs.UpdateResourceLimitsInput{
			Resources: &ncs.ResourceRequirements{
				Limits: ncs.ResourceList{
					CPU:    cpuLimit,
					Memory: memLimit,
				},
				Requests: ncs.ResourceList{
					CPU:    cpuReq,
					Memory: memReq,
				},
			},
		}

		if err := client.UpdateResourceLimits(ctx, workloadID, input); err != nil {
			exitWithError("Failed to update resources", err)
		}

		fmt.Printf("Workload %s resources updated successfully\n", workloadID)
	},
}

var ncsConfigureHealthCheckCmd = &cobra.Command{
	Use:   "configure-health-check",
	Short: "Configure workload health check",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")

		livenessPath, _ := cmd.Flags().GetString("liveness-path")
		readinessPath, _ := cmd.Flags().GetString("readiness-path")
		port, _ := cmd.Flags().GetInt("port")
		delay, _ := cmd.Flags().GetInt("initial-delay")

		config := &ncs.HealthCheckConfig{
			LivenessProbe: &ncs.HealthProbe{
				Type: "HTTP",
				HTTPGet: &ncs.HTTPGetAction{
					Path: livenessPath,
					Port: port,
				},
				InitialDelaySeconds: delay,
			},
			ReadinessProbe: &ncs.HealthProbe{
				Type: "HTTP",
				HTTPGet: &ncs.HTTPGetAction{
					Path: readinessPath,
					Port: port,
				},
				InitialDelaySeconds: delay,
			},
		}

		if err := client.ConfigureHealthCheck(ctx, workloadID, config); err != nil {
			exitWithError("Failed to configure health check", err)
		}

		fmt.Printf("Health check configured for workload %s\n", workloadID)
	},
}

var ncsGetHealthCheckStatusCmd = &cobra.Command{
	Use:     "describe-health-check-status",
	Aliases: []string{"health-check-status"},
	Short:   "Get health check status",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")

		result, err := client.GetHealthCheckStatus(ctx, workloadID)
		if err != nil {
			exitWithError("Failed to get health check status", err)
		}

		fmt.Printf("Health Check Status for %s:\n", workloadID)
		for _, s := range result.Status {
			fmt.Printf("  - Container: %s\n", s.ContainerName)
			fmt.Printf("    Liveness:  %s\n", s.Liveness)
			fmt.Printf("    Readiness: %s\n", s.Readiness)
			fmt.Printf("    LastCheck: %s\n", s.LastCheck)
		}
	},
}
