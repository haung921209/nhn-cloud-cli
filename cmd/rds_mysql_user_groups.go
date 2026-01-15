package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

// ============================================================================
// User Group Commands
// ============================================================================

var mysqlDescribeUserGroupsCmd = &cobra.Command{
	Use:   "describe-user-groups",
	Short: "Describe user groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		groupID, _ := cmd.Flags().GetString("user-group-id")

		if groupID != "" {
			result, err := client.GetUserGroup(context.Background(), groupID)
			if err != nil {
				exitWithError("failed to get user group", err)
			}
			if output == "json" {
				printJSON(result)
			} else {
				mysqlPrintUserGroupDetail(result)
			}
		} else {
			result, err := client.ListUserGroups(context.Background())
			if err != nil {
				exitWithError("failed to list user groups", err)
			}
			if output == "json" {
				printJSON(result)
			} else {
				mysqlPrintUserGroupList(result)
			}
		}
	},
}

var mysqlCreateUserGroupCmd = &cobra.Command{
	Use:   "create-user-group",
	Short: "Create a user group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()

		name, _ := cmd.Flags().GetString("name")
		memberIDs, _ := cmd.Flags().GetStringSlice("member-ids")
		selectAll, _ := cmd.Flags().GetBool("select-all")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		req := &mysql.CreateUserGroupRequest{
			UserGroupName: name,
			MemberIDs:     memberIDs,
			SelectAllYN:   selectAll,
		}

		result, err := client.CreateUserGroup(context.Background(), req)
		if err != nil {
			exitWithError("failed to create user group", err)
		}

		fmt.Printf("User group created.\n")
		fmt.Printf("ID: %s\n", result.UserGroupID)
	},
}

var mysqlDeleteUserGroupCmd = &cobra.Command{
	Use:   "delete-user-group",
	Short: "Delete a user group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()

		groupID, _ := cmd.Flags().GetString("user-group-id")
		if groupID == "" {
			exitWithError("--user-group-id is required", nil)
		}

		_, err := client.DeleteUserGroup(context.Background(), groupID)
		if err != nil {
			exitWithError("failed to delete user group", err)
		}

		fmt.Printf("User group deleted successfully.\n")
	},
}

// ============================================================================
// Monitoring Commands
// ============================================================================

var mysqlDescribeMetricsCmd = &cobra.Command{
	Use:   "describe-metrics",
	Short: "Describe available monitoring metrics",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()

		result, err := client.ListMetrics(context.Background())
		if err != nil {
			exitWithError("failed to list metrics", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			mysqlPrintMetricList(result)
		}
	},
}

var mysqlGetMetricStatisticsCmd = &cobra.Command{
	Use:   "get-metric-statistics",
	Short: "Get metric statistics for an instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()

		instanceID, err := getResolvedInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		interval, _ := cmd.Flags().GetInt("interval")

		if from == "" || to == "" {
			exitWithError("--from and --to are required (ISO8601 format)", nil)
		}

		result, err := client.GetMetricStatistics(context.Background(), instanceID, from, to, interval)
		if err != nil {
			exitWithError("failed to get metric statistics", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			mysqlPrintMetricStatistics(result)
		}
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func mysqlPrintUserGroupList(result *mysql.ListUserGroupsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "GROUP_ID\tNAME\tCREATED")
	for _, ug := range result.UserGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			ug.UserGroupID,
			ug.UserGroupName,
			ug.CreatedYmdt,
		)
	}
	w.Flush()
}

func mysqlPrintUserGroupDetail(result *mysql.GetUserGroupResponse) {
	fmt.Printf("ID: %s\n", result.UserGroupID)
	fmt.Printf("Name: %s\n", result.UserGroupName)
	fmt.Printf("Type: %s\n", result.UserGroupTypeCode)
	fmt.Printf("Created: %s\n", result.CreatedYmdt)
	fmt.Printf("Updated: %s\n", result.UpdatedYmdt)
	if len(result.Members) > 0 {
		fmt.Printf("Members:\n")
		for _, m := range result.Members {
			fmt.Printf("  - %s\n", m.MemberID)
		}
	}
}

func mysqlPrintMetricList(result *mysql.ListMetricsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "METRIC_NAME\tUNIT")
	for _, m := range result.Metrics {
		fmt.Fprintf(w, "%s\t%s\n", m.MetricName, m.Unit)
	}
	w.Flush()
}

func mysqlPrintMetricStatistics(result *mysql.GetMetricStatisticsResponse) {
	for _, stat := range result.MetricStatistics {
		fmt.Printf("Metric: %s (%s)\n", stat.MetricName, stat.Unit)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  TIMESTAMP\tVALUE")
		for _, v := range stat.Values {
			fmt.Fprintf(w, "  %s\t%.2f\n", v.Timestamp, v.Value)
		}
		w.Flush()
		fmt.Println()
	}
}

func init() {
	// User Group commands
	rdsMySQLCmd.AddCommand(mysqlDescribeUserGroupsCmd)
	rdsMySQLCmd.AddCommand(mysqlCreateUserGroupCmd)
	rdsMySQLCmd.AddCommand(mysqlDeleteUserGroupCmd)

	mysqlDescribeUserGroupsCmd.Flags().String("user-group-id", "", "Specific user group ID")

	mysqlCreateUserGroupCmd.Flags().String("name", "", "User group name (required)")
	mysqlCreateUserGroupCmd.Flags().StringSlice("member-ids", nil, "Member IDs to add")
	mysqlCreateUserGroupCmd.Flags().Bool("select-all", false, "Select all project members")

	mysqlDeleteUserGroupCmd.Flags().String("user-group-id", "", "User group ID (required)")

	// Monitoring commands
	rdsMySQLCmd.AddCommand(mysqlDescribeMetricsCmd)
	rdsMySQLCmd.AddCommand(mysqlGetMetricStatisticsCmd)

	mysqlGetMetricStatisticsCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	mysqlGetMetricStatisticsCmd.Flags().String("from", "", "Start time (ISO8601 format, required)")
	mysqlGetMetricStatisticsCmd.Flags().String("to", "", "End time (ISO8601 format, required)")
	mysqlGetMetricStatisticsCmd.Flags().Int("interval", 60, "Interval in seconds (1, 5, 30, 60)")
}
