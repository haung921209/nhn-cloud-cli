package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/loadbalancer"
	"github.com/spf13/cobra"
)

func init() {
	loadbalancerCmd.AddCommand(lbDescribeHealthMonitorsCmd)
	loadbalancerCmd.AddCommand(lbCreateHealthMonitorCmd)
	loadbalancerCmd.AddCommand(lbDeleteHealthMonitorCmd)

	lbDescribeHealthMonitorsCmd.Flags().String("monitor-id", "", "Health Monitor ID")

	lbCreateHealthMonitorCmd.Flags().String("pool-id", "", "Pool ID (required)")
	lbCreateHealthMonitorCmd.Flags().String("type", "TCP", "Monitor type (TCP/HTTP/HTTPS/PING)")
	lbCreateHealthMonitorCmd.Flags().Int("delay", 5, "Delay between checks (seconds)")
	lbCreateHealthMonitorCmd.Flags().Int("timeout", 5, "Check timeout (seconds)")
	lbCreateHealthMonitorCmd.Flags().Int("max-retries", 3, "Max retries before marking DOWN")
	lbCreateHealthMonitorCmd.Flags().String("http-method", "GET", "HTTP method (for HTTP/HTTPS)")
	lbCreateHealthMonitorCmd.Flags().String("url-path", "/", "URL path (for HTTP/HTTPS)")
	lbCreateHealthMonitorCmd.Flags().String("expected-codes", "200", "Expected HTTP codes")
	lbCreateHealthMonitorCmd.MarkFlagRequired("pool-id")

	lbDeleteHealthMonitorCmd.Flags().String("monitor-id", "", "Health Monitor ID (required)")
	lbDeleteHealthMonitorCmd.MarkFlagRequired("monitor-id")
}

var lbDescribeHealthMonitorsCmd = &cobra.Command{
	Use:   "describe-health-monitors",
	Short: "Describe health monitors",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("monitor-id")

		if id != "" {
			result, err := client.GetHealthMonitor(ctx, id)
			if err != nil {
				exitWithError("Failed to get health monitor", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			h := result.HealthMonitor
			fmt.Printf("ID:         %s\n", h.ID)
			fmt.Printf("Name:       %s\n", h.Name)
			fmt.Printf("Type:       %s\n", h.Type)
			fmt.Printf("Pool ID:    %s\n", h.PoolID)
			fmt.Printf("Delay:      %d\n", h.Delay)
			fmt.Printf("Timeout:    %d\n", h.Timeout)
			fmt.Printf("Max Retries:%d\n", h.MaxRetries)
			fmt.Printf("URL Path:   %s\n", h.URLPath)
			fmt.Printf("Status:     %s\n", h.OperatingStatus)
		} else {
			result, err := client.ListHealthMonitors(ctx)
			if err != nil {
				exitWithError("Failed to list health monitors", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tTYPE\tPOOL_ID\tSTATUS")
			for _, h := range result.HealthMonitors {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					h.ID, h.Name, h.Type, h.PoolID, h.OperatingStatus)
			}
			w.Flush()
		}
	},
}

var lbCreateHealthMonitorCmd = &cobra.Command{
	Use:   "create-health-monitor",
	Short: "Create a new health monitor",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		poolID, _ := cmd.Flags().GetString("pool-id")
		monitorType, _ := cmd.Flags().GetString("type")
		delay, _ := cmd.Flags().GetInt("delay")
		timeout, _ := cmd.Flags().GetInt("timeout")
		maxRetries, _ := cmd.Flags().GetInt("max-retries")
		httpMethod, _ := cmd.Flags().GetString("http-method")
		urlPath, _ := cmd.Flags().GetString("url-path")
		expectedCodes, _ := cmd.Flags().GetString("expected-codes")

		input := &loadbalancer.CreateHealthMonitorInput{
			PoolID:        poolID,
			Type:          monitorType,
			Delay:         delay,
			Timeout:       timeout,
			MaxRetries:    maxRetries,
			HTTPMethod:    httpMethod,
			URLPath:       urlPath,
			ExpectedCodes: expectedCodes,
		}

		result, err := client.CreateHealthMonitor(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create health monitor", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Health monitor created: %s\n", result.HealthMonitor.ID)
		fmt.Printf("Type: %s\n", result.HealthMonitor.Type)
	},
}

var lbDeleteHealthMonitorCmd = &cobra.Command{
	Use:   "delete-health-monitor",
	Short: "Delete a health monitor",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		id, _ := cmd.Flags().GetString("monitor-id")
		if err := client.DeleteHealthMonitor(context.Background(), id); err != nil {
			exitWithError("Failed to delete health monitor", err)
		}
		fmt.Printf("Health monitor %s deleted\n", id)
	},
}
