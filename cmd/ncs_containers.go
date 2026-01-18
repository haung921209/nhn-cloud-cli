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
	ncsCmd.AddCommand(ncsGetContainerLogsCmd)
	ncsCmd.AddCommand(ncsExecCommandCmd)
	ncsCmd.AddCommand(ncsGetContainerStatusCmd)

	ncsGetContainerLogsCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsGetContainerLogsCmd.Flags().Int("tail", 100, "Number of lines to show from end of logs")
	ncsGetContainerLogsCmd.Flags().Int("since", 0, "Show logs since N seconds ago")
	ncsGetContainerLogsCmd.MarkFlagRequired("workload-id")

	ncsExecCommandCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsExecCommandCmd.Flags().String("container", "", "Container name")
	ncsExecCommandCmd.Flags().BoolP("stdin", "i", false, "Pass stdin to the container")
	ncsExecCommandCmd.Flags().BoolP("tty", "t", false, "Allocate a pseudo-TTY")
	ncsExecCommandCmd.MarkFlagRequired("workload-id")

	ncsGetContainerStatusCmd.Flags().String("workload-id", "", "Workload ID (required)")
	ncsGetContainerStatusCmd.MarkFlagRequired("workload-id")
}

var ncsGetContainerLogsCmd = &cobra.Command{
	Use:     "get-container-logs",
	Aliases: []string{"logs"},
	Short:   "Get container logs",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")
		tail, _ := cmd.Flags().GetInt("tail")
		since, _ := cmd.Flags().GetInt("since")

		result, err := client.GetWorkloadLogs(ctx, workloadID, tail, since)
		if err != nil {
			exitWithError("Failed to get logs", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		for _, log := range result.Logs {
			fmt.Printf("[%s] %s\n", log.Timestamp, log.Message)
		}
	},
}

var ncsExecCommandCmd = &cobra.Command{
	Use:     "exec-command [command...]",
	Aliases: []string{"exec"},
	Short:   "Execute a command in a container",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")
		container, _ := cmd.Flags().GetString("container")
		stdin, _ := cmd.Flags().GetBool("stdin")
		tty, _ := cmd.Flags().GetBool("tty")

		if len(args) == 0 {
			exitWithError("Command required", nil)
		}

		input := &ncs.ExecInput{
			ContainerName: container,
			Command:       args,
			Stdin:         stdin,
			Stdout:        true,
			Stderr:        true,
			TTY:           tty,
		}

		result, err := client.ExecWorkloadContainer(ctx, workloadID, input)
		if err != nil {
			exitWithError("Failed to exec command", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		if result.Output != "" {
			fmt.Println(result.Output)
		}
		if result.Error != "" {
			fmt.Fprintf(os.Stderr, "Error: %s\n", result.Error)
		}
	},
}

var ncsGetContainerStatusCmd = &cobra.Command{
	Use:     "describe-container-status",
	Aliases: []string{"container-status"},
	Short:   "Get container status",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		workloadID, _ := cmd.Flags().GetString("workload-id")

		result, err := client.GetContainerStatus(ctx, workloadID)
		if err != nil {
			exitWithError("Failed to get container status", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CONTAINER\tSTATE\tREADY\tRESTARTS\tIMAGE")
		for _, s := range result.Containers {
			fmt.Fprintf(w, "%s\t%s\t%v\t%d\t%s\n",
				s.ContainerName, s.State, s.Ready, s.RestartCount, s.Image)
		}
		w.Flush()
	},
}
