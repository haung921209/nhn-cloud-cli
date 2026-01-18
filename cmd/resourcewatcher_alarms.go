package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/resourcewatcher"
	"github.com/spf13/cobra"
)

func init() {
	resourceWatcherCmd.AddCommand(rwDescribeAlarmsCmd)
	resourceWatcherCmd.AddCommand(rwGetAlarmCmd)
	resourceWatcherCmd.AddCommand(rwCreateAlarmCmd)
	resourceWatcherCmd.AddCommand(rwUpdateAlarmCmd)
	resourceWatcherCmd.AddCommand(rwDeleteAlarmCmd)

	// List flags
	rwDescribeAlarmsCmd.Flags().String("name", "", "Filter by alarm name")
	rwDescribeAlarmsCmd.Flags().String("status", "", "Filter by status (STABLE, DISABLED, CLOSED)")
	rwDescribeAlarmsCmd.Flags().Int("page", 0, "Page number")
	rwDescribeAlarmsCmd.Flags().Int("size", 20, "Page size")

	// Get flags
	rwGetAlarmCmd.Flags().String("alarm-id", "", "Alarm ID (required)")
	rwGetAlarmCmd.MarkFlagRequired("alarm-id")

	// Create flags
	rwCreateAlarmCmd.Flags().String("name", "", "Alarm name (required)")
	rwCreateAlarmCmd.Flags().String("description", "", "Alarm description")
	rwCreateAlarmCmd.Flags().String("event-rule-id", "", "Event rule ID")
	rwCreateAlarmCmd.Flags().String("resource-group-id", "", "Resource group ID")
	rwCreateAlarmCmd.Flags().String("resource-tag-id", "", "Resource tag ID")
	rwCreateAlarmCmd.Flags().StringSlice("target", nil, "Target in format TYPE:ID (e.g., UUID:user-id, WEBHOOK:webhook-url)")
	rwCreateAlarmCmd.MarkFlagRequired("name")

	// Update flags
	rwUpdateAlarmCmd.Flags().String("alarm-id", "", "Alarm ID (required)")
	rwUpdateAlarmCmd.Flags().String("name", "", "New alarm name")
	rwUpdateAlarmCmd.Flags().String("description", "", "New alarm description")
	rwUpdateAlarmCmd.Flags().String("status", "", "New status (STABLE, DISABLED, CLOSED)")
	rwUpdateAlarmCmd.Flags().String("event-rule-id", "", "New event rule ID")
	rwUpdateAlarmCmd.Flags().String("resource-group-id", "", "New resource group ID")
	rwUpdateAlarmCmd.Flags().String("resource-tag-id", "", "New resource tag ID")
	rwUpdateAlarmCmd.Flags().StringSlice("target", nil, "Target in format TYPE:ID")
	rwUpdateAlarmCmd.MarkFlagRequired("alarm-id")

	// Delete flags
	rwDeleteAlarmCmd.Flags().String("alarm-id", "", "Alarm ID (required)")
	rwDeleteAlarmCmd.MarkFlagRequired("alarm-id")
}

var rwDescribeAlarmsCmd = &cobra.Command{
	Use:     "describe-alarms",
	Aliases: []string{"list-alarms", "list"},
	Short:   "List event alarms",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		status, _ := cmd.Flags().GetString("status")
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")

		input := &resourcewatcher.SearchEventAlarmsInput{
			AlarmName:       name,
			AlarmStatusCode: status,
			Page:            page,
			Size:            size,
		}

		result, err := client.SearchEventAlarms(ctx, input)
		if err != nil {
			exitWithError("Failed to list alarms", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ALARM_ID\tNAME\tSTATUS\tCREATED")
		for _, a := range result.Alarms {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				a.AlarmID, a.AlarmName, a.AlarmStatusCode, a.CreatedDateTime)
		}
		w.Flush()
		fmt.Printf("\nTotal: %d\n", result.TotalCount)
	},
}

var rwGetAlarmCmd = &cobra.Command{
	Use:     "describe-alarm",
	Aliases: []string{"get-alarm", "get"},
	Short:   "Get alarm details",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("alarm-id")

		result, err := client.GetEventAlarm(ctx, id)
		if err != nil {
			exitWithError("Failed to get alarm", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		a := result.Alarm
		fmt.Printf("Alarm ID:         %s\n", a.AlarmID)
		fmt.Printf("Name:             %s\n", a.AlarmName)
		fmt.Printf("Description:      %s\n", a.AlarmDescription)
		fmt.Printf("Status:           %s\n", a.AlarmStatusCode)
		fmt.Printf("Event Rule ID:    %s\n", a.EventRuleID)
		fmt.Printf("Resource Group:   %s\n", a.ResourceGroupID)
		fmt.Printf("Resource Tag:     %s\n", a.ResourceTagID)
		fmt.Printf("Created:          %s\n", a.CreatedDateTime)
		fmt.Printf("Updated:          %s\n", a.UpdatedDateTime)

		if len(a.Targets) > 0 {
			fmt.Println("\nTargets:")
			for _, t := range a.Targets {
				fmt.Printf("  - %s: %s\n", t.TargetType, t.TargetID)
			}
		}
	},
}

var rwCreateAlarmCmd = &cobra.Command{
	Use:   "create-alarm",
	Short: "Create a new event alarm",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		eventRuleID, _ := cmd.Flags().GetString("event-rule-id")
		resourceGroupID, _ := cmd.Flags().GetString("resource-group-id")
		resourceTagID, _ := cmd.Flags().GetString("resource-tag-id")
		targets, _ := cmd.Flags().GetStringSlice("target")

		// Parse targets
		var alarmTargets []resourcewatcher.AlarmTarget
		for _, t := range targets {
			var targetType, targetID string
			if _, err := fmt.Sscanf(t, "%[^:]:%s", &targetType, &targetID); err == nil {
				alarmTargets = append(alarmTargets, resourcewatcher.AlarmTarget{
					TargetType: targetType,
					TargetID:   targetID,
				})
			}
		}

		input := &resourcewatcher.CreateEventAlarmInput{
			AlarmName:        name,
			AlarmDescription: description,
			EventRuleID:      eventRuleID,
			ResourceGroupID:  resourceGroupID,
			ResourceTagID:    resourceTagID,
			Targets:          alarmTargets,
		}

		result, err := client.CreateEventAlarm(ctx, input)
		if err != nil {
			exitWithError("Failed to create alarm", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Alarm created successfully!\n")
		fmt.Printf("Alarm ID: %s\n", result.AlarmID)
	},
}

var rwUpdateAlarmCmd = &cobra.Command{
	Use:   "update-alarm",
	Short: "Update an event alarm",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		id, _ := cmd.Flags().GetString("alarm-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")
		eventRuleID, _ := cmd.Flags().GetString("event-rule-id")
		resourceGroupID, _ := cmd.Flags().GetString("resource-group-id")
		resourceTagID, _ := cmd.Flags().GetString("resource-tag-id")
		targets, _ := cmd.Flags().GetStringSlice("target")

		// Parse targets
		var alarmTargets []resourcewatcher.AlarmTarget
		for _, t := range targets {
			var targetType, targetID string
			if _, err := fmt.Sscanf(t, "%[^:]:%s", &targetType, &targetID); err == nil {
				alarmTargets = append(alarmTargets, resourcewatcher.AlarmTarget{
					TargetType: targetType,
					TargetID:   targetID,
				})
			}
		}

		input := &resourcewatcher.UpdateEventAlarmInput{
			AlarmName:        name,
			AlarmDescription: description,
			AlarmStatusCode:  status,
			EventRuleID:      eventRuleID,
			ResourceGroupID:  resourceGroupID,
			ResourceTagID:    resourceTagID,
		}
		if len(alarmTargets) > 0 {
			input.Targets = alarmTargets
		}

		_, err := client.UpdateEventAlarm(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update alarm", err)
		}

		fmt.Printf("Alarm %s updated successfully\n", id)
	},
}

var rwDeleteAlarmCmd = &cobra.Command{
	Use:   "delete-alarm",
	Short: "Delete an event alarm",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("alarm-id")

		_, err := client.DeleteEventAlarm(ctx, id)
		if err != nil {
			exitWithError("Failed to delete alarm", err)
		}

		fmt.Printf("Alarm %s deleted successfully\n", id)
	},
}
