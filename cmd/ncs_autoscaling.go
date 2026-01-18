package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncs"
	"github.com/spf13/cobra"
)

func init() {
	ncsCmd.AddCommand(ncsConfigureAutoScalingCmd)
	ncsCmd.AddCommand(ncsGetAutoScalingStatusCmd)

	ncsConfigureAutoScalingCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsConfigureAutoScalingCmd.Flags().Bool("enabled", false, "Enable or disable auto-scaling")
	ncsConfigureAutoScalingCmd.Flags().Int("min-replicas", 1, "Minimum number of replicas")
	ncsConfigureAutoScalingCmd.Flags().Int("max-replicas", 10, "Maximum number of replicas")
	ncsConfigureAutoScalingCmd.Flags().Int("target-cpu", 80, "Target CPU utilization percentage")
	ncsConfigureAutoScalingCmd.Flags().Int("target-memory", 0, "Target memory utilization percentage (optional)")
	ncsConfigureAutoScalingCmd.MarkFlagRequired("workload-id")

	ncsGetAutoScalingStatusCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsGetAutoScalingStatusCmd.MarkFlagRequired("workload-id")
}

var ncsConfigureAutoScalingCmd = &cobra.Command{
	Use:   "configure-auto-scaling",
	Short: "Configure workload auto-scaling",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")
		enabled, _ := cmd.Flags().GetBool("enabled")
		min, _ := cmd.Flags().GetInt("min-replicas")
		max, _ := cmd.Flags().GetInt("max-replicas")
		cpu, _ := cmd.Flags().GetInt("target-cpu")
		mem, _ := cmd.Flags().GetInt("target-memory")

		input := &ncs.ConfigureAutoScalingInput{
			Enabled: enabled,
			Policy: &ncs.AutoScalingPolicy{
				MinReplicas:                       min,
				MaxReplicas:                       max,
				TargetCPUUtilizationPercentage:    cpu,
				TargetMemoryUtilizationPercentage: mem,
			},
		}

		if err := client.ConfigureAutoScaling(ctx, workloadID, input); err != nil {
			exitWithError("Failed to configure auto-scaling", err)
		}

		fmt.Printf("Auto-scaling configured for workload %s\n", workloadID)
	},
}

var ncsGetAutoScalingStatusCmd = &cobra.Command{
	Use:     "describe-auto-scaling-status",
	Aliases: []string{"auto-scaling-status"},
	Short:   "Get auto-scaling status",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")

		result, err := client.GetAutoScalingStatus(ctx, workloadID)
		if err != nil {
			exitWithError("Failed to get auto-scaling status", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Auto-Scaling Status for %s:\n", workloadID)
		fmt.Printf("  Enabled:      %v\n", result.Status.Enabled)
		fmt.Printf("  Current Reps: %d\n", result.Status.CurrentReplicas)
		fmt.Printf("  Desired Reps: %d\n", result.Status.DesiredReplicas)
		if result.Status.Policy != nil {
			fmt.Printf("  Min Replicas: %d\n", result.Status.Policy.MinReplicas)
			fmt.Printf("  Max Replicas: %d\n", result.Status.Policy.MaxReplicas)
		}
	},
}
