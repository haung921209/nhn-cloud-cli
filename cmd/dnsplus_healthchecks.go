package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/dnsplus"
	"github.com/spf13/cobra"
)

func init() {
	dnsplusCmd.AddCommand(dnsDescribeHealthChecksCmd)
	dnsplusCmd.AddCommand(dnsCreateHealthCheckCmd)
	dnsplusCmd.AddCommand(dnsDeleteHealthCheckCmd)

	dnsCreateHealthCheckCmd.Flags().String("name", "", "Health check name")
	dnsCreateHealthCheckCmd.Flags().String("description", "", "Health check description")
	dnsCreateHealthCheckCmd.Flags().String("protocol", "", "Protocol (HTTP, HTTPS, TCP, ICMP)")
	dnsCreateHealthCheckCmd.Flags().Int("port", 0, "Port number")
	dnsCreateHealthCheckCmd.Flags().String("path", "", "URL path")
	dnsCreateHealthCheckCmd.Flags().String("host", "", "Host header")
	dnsCreateHealthCheckCmd.Flags().Int("interval", 30, "Check interval")
	dnsCreateHealthCheckCmd.Flags().Int("timeout", 5, "Check timeout")
	dnsCreateHealthCheckCmd.Flags().Int("retries", 3, "Retry count")
	dnsCreateHealthCheckCmd.Flags().String("expected-codes", "200", "Expected HTTP codes")
	dnsCreateHealthCheckCmd.MarkFlagRequired("name")
	dnsCreateHealthCheckCmd.MarkFlagRequired("protocol")
	dnsCreateHealthCheckCmd.MarkFlagRequired("port")

	dnsDeleteHealthCheckCmd.Flags().StringSlice("health-check-ids", nil, "Health Check IDs (required)")
	dnsDeleteHealthCheckCmd.MarkFlagRequired("health-check-ids")
}

var dnsDescribeHealthChecksCmd = &cobra.Command{
	Use:   "describe-health-checks",
	Short: "List health checks",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListHealthChecks(ctx)
		if err != nil {
			return fmt.Errorf("failed to list health checks: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.HealthCheckList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tPROTOCOL\tPORT\tINTERVAL\tTIMEOUT")
		for _, hc := range result.HealthCheckList {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%d\n",
				hc.HealthCheckID,
				hc.HealthCheckName,
				hc.Protocol,
				hc.Port,
				hc.Interval,
				hc.Timeout,
			)
		}
		return w.Flush()
	},
}

var dnsCreateHealthCheckCmd = &cobra.Command{
	Use:   "create-health-check",
	Short: "Create a health check",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		protocol, _ := cmd.Flags().GetString("protocol")
		port, _ := cmd.Flags().GetInt("port")
		path, _ := cmd.Flags().GetString("path")
		host, _ := cmd.Flags().GetString("host")
		interval, _ := cmd.Flags().GetInt("interval")
		timeout, _ := cmd.Flags().GetInt("timeout")
		retries, _ := cmd.Flags().GetInt("retries")
		expectedCodes, _ := cmd.Flags().GetString("expected-codes")

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.CreateHealthCheckInput{
			HealthCheckName: name,
			Description:     description,
			Protocol:        protocol,
			Port:            port,
			Path:            path,
			Host:            host,
			Interval:        interval,
			Timeout:         timeout,
			Retries:         retries,
			ExpectedCodes:   expectedCodes,
		}

		result, err := client.CreateHealthCheck(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to create health check: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.HealthCheck)
		}

		fmt.Printf("Health check created successfully: %s (%s)\n", result.HealthCheck.HealthCheckName, result.HealthCheck.HealthCheckID)
		return nil
	},
}

var dnsDeleteHealthCheckCmd = &cobra.Command{
	Use:   "delete-health-checks",
	Short: "Delete health checks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, _ := cmd.Flags().GetStringSlice("health-check-ids")
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeleteHealthChecks(ctx, ids)
		if err != nil {
			return fmt.Errorf("failed to delete health checks: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Health check(s) deleted successfully\n")
		return nil
	},
}
