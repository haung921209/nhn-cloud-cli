package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
	"github.com/spf13/cobra"
)

// ============================================================================
// User Group Commands
// ============================================================================

var postgresqlDescribeUserGroupsCmd = &cobra.Command{
	Use:   "describe-user-groups",
	Short: "Describe user groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		groupID, _ := cmd.Flags().GetString("user-group-id")

		if groupID != "" {
			result, err := client.GetUserGroup(context.Background(), groupID)
			if err != nil {
				exitWithError("failed to get user group", err)
			}
			if output == "json" {
				postgresqlPrintJSON(result)
			} else {
				postgresqlPrintUserGroupDetail(result)
			}
		} else {
			result, err := client.ListUserGroups(context.Background())
			if err != nil {
				exitWithError("failed to list user groups", err)
			}
			if output == "json" {
				postgresqlPrintJSON(result)
			} else {
				postgresqlPrintUserGroupList(result)
			}
		}
	},
}

var postgresqlCreateUserGroupCmd = &cobra.Command{
	Use:   "create-user-group",
	Short: "Create a user group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		name, _ := cmd.Flags().GetString("name")
		memberIDs, _ := cmd.Flags().GetStringSlice("member-ids")
		selectAll, _ := cmd.Flags().GetBool("select-all")

		if name == "" {
			exitWithError("--name is required", nil)
		}

		req := &postgresql.CreateUserGroupRequest{
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

var postgresqlDeleteUserGroupCmd = &cobra.Command{
	Use:   "delete-user-group",
	Short: "Delete a user group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

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

var postgresqlDescribeMetricsCmd = &cobra.Command{
	Use:   "describe-metrics",
	Short: "Describe available monitoring metrics",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		result, err := client.ListMetrics(context.Background())
		if err != nil {
			exitWithError("failed to list metrics", err)
		}

		if output == "json" {
			postgresqlPrintJSON(result)
		} else {
			postgresqlPrintMetricList(result)
		}
	},
}

var postgresqlGetMetricStatisticsCmd = &cobra.Command{
	Use:   "get-metric-statistics",
	Short: "Get metric statistics for an instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
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
			postgresqlPrintJSON(result)
		} else {
			postgresqlPrintMetricStatistics(result)
		}
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func postgresqlPrintUserGroupList(result *postgresql.ListUserGroupsResponse) {
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

func postgresqlPrintUserGroupDetail(result *postgresql.GetUserGroupResponse) {
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

func postgresqlPrintMetricList(result *postgresql.ListMetricsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "METRIC_NAME\tUNIT")
	for _, m := range result.Metrics {
		fmt.Fprintf(w, "%s\t%s\n", m.MetricName, m.Unit)
	}
	w.Flush()
}

func postgresqlPrintMetricStatistics(result *postgresql.GetMetricStatisticsResponse) {
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
	rdsPostgreSQLCmd.AddCommand(postgresqlDescribeUserGroupsCmd)
	rdsPostgreSQLCmd.AddCommand(postgresqlCreateUserGroupCmd)
	rdsPostgreSQLCmd.AddCommand(postgresqlDeleteUserGroupCmd)

	postgresqlDescribeUserGroupsCmd.Flags().String("user-group-id", "", "Specific user group ID")

	postgresqlCreateUserGroupCmd.Flags().String("name", "", "User group name (required)")
	postgresqlCreateUserGroupCmd.Flags().StringSlice("member-ids", nil, "Member IDs to add")
	postgresqlCreateUserGroupCmd.Flags().Bool("select-all", false, "Select all project members")

	postgresqlDeleteUserGroupCmd.Flags().String("user-group-id", "", "User group ID (required)")

	// Monitoring commands
	rdsPostgreSQLCmd.AddCommand(postgresqlDescribeMetricsCmd)
	rdsPostgreSQLCmd.AddCommand(postgresqlGetMetricStatisticsCmd)

	postgresqlGetMetricStatisticsCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	postgresqlGetMetricStatisticsCmd.Flags().String("from", "", "Start time (ISO8601 format, required)")
	postgresqlGetMetricStatisticsCmd.Flags().String("to", "", "End time (ISO8601 format, required)")
	postgresqlGetMetricStatisticsCmd.Flags().Int("interval", 60, "Interval in seconds")
}
