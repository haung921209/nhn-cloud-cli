package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mariadb"
	"github.com/spf13/cobra"
)

// ============================================================================
// User Group Commands
// ============================================================================

var mariadbDescribeUserGroupsCmd = &cobra.Command{
	Use:   "describe-user-groups",
	Short: "Describe user groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		groupID, _ := cmd.Flags().GetString("user-group-id")

		if groupID != "" {
			result, err := client.GetUserGroup(context.Background(), groupID)
			if err != nil {
				exitWithError("failed to get user group", err)
			}
			if output == "json" {
				mariadbPrintJSON(result)
			} else {
				mariadbPrintUserGroupDetail(result)
			}
		} else {
			result, err := client.ListUserGroups(context.Background())
			if err != nil {
				exitWithError("failed to list user groups", err)
			}
			if output == "json" {
				mariadbPrintJSON(result)
			} else {
				mariadbPrintUserGroupList(result)
			}
		}
	},
}

var mariadbCreateUserGroupCmd = &cobra.Command{
	Use:   "create-user-group",
	Short: "Create a user group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		name, _ := cmd.Flags().GetString("name")
		memberIDs, _ := cmd.Flags().GetStringSlice("member-ids")
		selectAll, _ := cmd.Flags().GetBool("select-all")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		req := &mariadb.CreateUserGroupRequest{
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

var mariadbDeleteUserGroupCmd = &cobra.Command{
	Use:   "delete-user-group",
	Short: "Delete a user group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

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

var mariadbDescribeMetricsCmd = &cobra.Command{
	Use:   "describe-metrics",
	Short: "Describe available monitoring metrics",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		result, err := client.ListMetrics(context.Background())
		if err != nil {
			exitWithError("failed to list metrics", err)
		}

		if output == "json" {
			mariadbPrintJSON(result)
		} else {
			mariadbPrintMetricList(result)
		}
	},
}

var mariadbGetMetricStatisticsCmd = &cobra.Command{
	Use:   "get-metric-statistics",
	Short: "Get metric statistics for an instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		instanceID, err := getResolvedMariaDBInstanceID(cmd, client)
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
			mariadbPrintJSON(result)
		} else {
			mariadbPrintMetricStatistics(result)
		}
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func mariadbPrintUserGroupList(result *mariadb.ListUserGroupsResponse) {
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

func mariadbPrintUserGroupDetail(result *mariadb.GetUserGroupResponse) {
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

func mariadbPrintMetricList(result *mariadb.ListMetricsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "METRIC_NAME\tUNIT")
	for _, m := range result.Metrics {
		fmt.Fprintf(w, "%s\t%s\n", m.MetricName, m.Unit)
	}
	w.Flush()
}

func mariadbPrintMetricStatistics(result *mariadb.GetMetricStatisticsResponse) {
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
	rdsMariaDBCmd.AddCommand(mariadbDescribeUserGroupsCmd)
	rdsMariaDBCmd.AddCommand(mariadbCreateUserGroupCmd)
	rdsMariaDBCmd.AddCommand(mariadbDeleteUserGroupCmd)

	mariadbDescribeUserGroupsCmd.Flags().String("user-group-id", "", "Specific user group ID")

	mariadbCreateUserGroupCmd.Flags().String("name", "", "User group name (required)")
	mariadbCreateUserGroupCmd.Flags().StringSlice("member-ids", nil, "Member IDs to add")
	mariadbCreateUserGroupCmd.Flags().Bool("select-all", false, "Select all project members")

	mariadbDeleteUserGroupCmd.Flags().String("user-group-id", "", "User group ID (required)")

	// Monitoring commands
	rdsMariaDBCmd.AddCommand(mariadbDescribeMetricsCmd)
	rdsMariaDBCmd.AddCommand(mariadbGetMetricStatisticsCmd)

	mariadbGetMetricStatisticsCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	mariadbGetMetricStatisticsCmd.Flags().String("from", "", "Start time (ISO8601 format, required)")
	mariadbGetMetricStatisticsCmd.Flags().String("to", "", "End time (ISO8601 format, required)")
	mariadbGetMetricStatisticsCmd.Flags().Int("interval", 60, "Interval in seconds")
}
