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

var resourceWatcherCmd = &cobra.Command{
	Use:     "resource-watcher",
	Aliases: []string{"rw", "watcher"},
	Short:   "Manage Resource Watcher (Governance) alarms and events",
	Long:    `Manage event alarms, alarm history, events, resource groups, and resource tags.`,
}

var rwAlarmCmd = &cobra.Command{
	Use:   "alarm",
	Short: "Manage event alarms",
}

var rwHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "View alarm history",
}

var rwEventCmd = &cobra.Command{
	Use:   "event",
	Short: "View events",
}

var rwResourceGroupCmd = &cobra.Command{
	Use:   "resource-group",
	Short: "View resource groups",
}

var rwResourceTagCmd = &cobra.Command{
	Use:   "resource-tag",
	Short: "View resource tags",
}

func init() {
	rootCmd.AddCommand(resourceWatcherCmd)

	// Alarm subcommands
	resourceWatcherCmd.AddCommand(rwAlarmCmd)
	rwAlarmCmd.AddCommand(rwAlarmListCmd)
	rwAlarmCmd.AddCommand(rwAlarmGetCmd)
	rwAlarmCmd.AddCommand(rwAlarmCreateCmd)
	rwAlarmCmd.AddCommand(rwAlarmUpdateCmd)
	rwAlarmCmd.AddCommand(rwAlarmDeleteCmd)

	// History subcommands
	resourceWatcherCmd.AddCommand(rwHistoryCmd)
	rwHistoryCmd.AddCommand(rwHistoryListCmd)
	rwHistoryCmd.AddCommand(rwHistoryGetCmd)

	// Event subcommands
	resourceWatcherCmd.AddCommand(rwEventCmd)
	rwEventCmd.AddCommand(rwEventListCmd)
	rwEventCmd.AddCommand(rwEventGetCmd)

	// Resource group subcommands
	resourceWatcherCmd.AddCommand(rwResourceGroupCmd)
	rwResourceGroupCmd.AddCommand(rwResourceGroupListCmd)

	// Resource tag subcommands
	resourceWatcherCmd.AddCommand(rwResourceTagCmd)
	rwResourceTagCmd.AddCommand(rwResourceTagListCmd)

	// Alarm list flags
	rwAlarmListCmd.Flags().String("name", "", "Filter by alarm name")
	rwAlarmListCmd.Flags().String("status", "", "Filter by status (STABLE, DISABLED, CLOSED)")
	rwAlarmListCmd.Flags().Int("page", 0, "Page number")
	rwAlarmListCmd.Flags().Int("size", 20, "Page size")

	// Alarm create flags
	rwAlarmCreateCmd.Flags().String("name", "", "Alarm name (required)")
	rwAlarmCreateCmd.Flags().String("description", "", "Alarm description")
	rwAlarmCreateCmd.Flags().String("event-rule-id", "", "Event rule ID")
	rwAlarmCreateCmd.Flags().String("resource-group-id", "", "Resource group ID")
	rwAlarmCreateCmd.Flags().String("resource-tag-id", "", "Resource tag ID")
	rwAlarmCreateCmd.Flags().StringSlice("target", nil, "Target in format TYPE:ID (e.g., UUID:user-id, WEBHOOK:webhook-url)")
	rwAlarmCreateCmd.MarkFlagRequired("name")

	// Alarm update flags
	rwAlarmUpdateCmd.Flags().String("name", "", "New alarm name")
	rwAlarmUpdateCmd.Flags().String("description", "", "New alarm description")
	rwAlarmUpdateCmd.Flags().String("status", "", "New status (STABLE, DISABLED, CLOSED)")
	rwAlarmUpdateCmd.Flags().String("event-rule-id", "", "New event rule ID")
	rwAlarmUpdateCmd.Flags().String("resource-group-id", "", "New resource group ID")
	rwAlarmUpdateCmd.Flags().String("resource-tag-id", "", "New resource tag ID")
	rwAlarmUpdateCmd.Flags().StringSlice("target", nil, "Target in format TYPE:ID")

	// History list flags
	rwHistoryListCmd.Flags().String("start", "", "Start datetime (ISO8601)")
	rwHistoryListCmd.Flags().String("end", "", "End datetime (ISO8601)")
	rwHistoryListCmd.Flags().Int("page", 0, "Page number")
	rwHistoryListCmd.Flags().Int("size", 20, "Page size")
}

func getResourceWatcherClient() *resourcewatcher.Client {
	return resourcewatcher.NewClient(getAppKey(), getAccessKey(), getSecretKey(), nil, debug)
}

// ============ Alarm Commands ============

var rwAlarmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List event alarms",
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

var rwAlarmGetCmd = &cobra.Command{
	Use:   "get [alarm-id]",
	Short: "Get alarm details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		result, err := client.GetEventAlarm(ctx, args[0])
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

var rwAlarmCreateCmd = &cobra.Command{
	Use:   "create",
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

var rwAlarmUpdateCmd = &cobra.Command{
	Use:   "update [alarm-id]",
	Short: "Update an event alarm",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

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

		_, err := client.UpdateEventAlarm(ctx, args[0], input)
		if err != nil {
			exitWithError("Failed to update alarm", err)
		}

		fmt.Printf("Alarm %s updated successfully\n", args[0])
	},
}

var rwAlarmDeleteCmd = &cobra.Command{
	Use:   "delete [alarm-id]",
	Short: "Delete an event alarm",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		_, err := client.DeleteEventAlarm(ctx, args[0])
		if err != nil {
			exitWithError("Failed to delete alarm", err)
		}

		fmt.Printf("Alarm %s deleted successfully\n", args[0])
	},
}

// ============ History Commands ============

var rwHistoryListCmd = &cobra.Command{
	Use:   "list [alarm-id]",
	Short: "List alarm history",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")

		input := &resourcewatcher.SearchAlarmHistoryInput{
			StartDateTime: start,
			EndDateTime:   end,
			Page:          page,
			Size:          size,
		}

		result, err := client.SearchAlarmHistory(ctx, args[0], input)
		if err != nil {
			exitWithError("Failed to list alarm history", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "HISTORY_ID\tEVENT_NAME\tRESOURCE\tCREATED")
		for _, h := range result.Histories {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				h.AlarmHistoryID, h.EventName, h.ResourceName, h.CreatedDateTime)
		}
		w.Flush()
		fmt.Printf("\nTotal: %d\n", result.TotalCount)
	},
}

var rwHistoryGetCmd = &cobra.Command{
	Use:   "get [alarm-id] [history-id]",
	Short: "Get alarm history details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		result, err := client.GetAlarmHistory(ctx, args[0], args[1])
		if err != nil {
			exitWithError("Failed to get alarm history", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		h := result.History
		fmt.Printf("History ID:    %s\n", h.AlarmHistoryID)
		fmt.Printf("Alarm ID:      %s\n", h.AlarmID)
		fmt.Printf("Event ID:      %s\n", h.EventID)
		fmt.Printf("Event Name:    %s\n", h.EventName)
		fmt.Printf("Product ID:    %s\n", h.ProductID)
		fmt.Printf("Resource ID:   %s\n", h.ResourceID)
		fmt.Printf("Resource Name: %s\n", h.ResourceName)
		fmt.Printf("Created:       %s\n", h.CreatedDateTime)

		if len(h.AlarmSendResults) > 0 {
			fmt.Println("\nSend Results:")
			for _, r := range h.AlarmSendResults {
				fmt.Printf("  - %s (%s): %s at %s\n", r.TargetType, r.TargetID, r.SendStatus, r.SentDateTime)
			}
		}
	},
}

// ============ Event Commands ============

var rwEventListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available events",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		result, err := client.ListEvents(ctx)
		if err != nil {
			exitWithError("Failed to list events", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "EVENT_ID\tPRODUCT\tNAME\tTYPE")
		for _, e := range result.Events {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				e.EventID, e.ProductID, e.EventName, e.EventType)
		}
		w.Flush()
	},
}

var rwEventGetCmd = &cobra.Command{
	Use:   "get [product-id] [event-id]",
	Short: "Get event details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		result, err := client.GetEvent(ctx, args[0], args[1])
		if err != nil {
			exitWithError("Failed to get event", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		e := result.Event
		fmt.Printf("Event ID:    %s\n", e.EventID)
		fmt.Printf("Product ID:  %s\n", e.ProductID)
		fmt.Printf("Event Name:  %s\n", e.EventName)
		fmt.Printf("Event Type:  %s\n", e.EventType)
		fmt.Printf("Description: %s\n", e.Description)
	},
}

// ============ Resource Group Commands ============

var rwResourceGroupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List resource groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		result, err := client.ListResourceGroups(ctx)
		if err != nil {
			exitWithError("Failed to list resource groups", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "GROUP_ID\tNAME\tDESCRIPTION\tCREATED")
		for _, g := range result.ResourceGroups {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				g.ResourceGroupID, g.ResourceGroupName, g.Description, g.CreatedDateTime)
		}
		w.Flush()
	},
}

// ============ Resource Tag Commands ============

var rwResourceTagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List resource tags",
	Run: func(cmd *cobra.Command, args []string) {
		client := getResourceWatcherClient()
		ctx := context.Background()

		result, err := client.ListResourceTags(ctx)
		if err != nil {
			exitWithError("Failed to list resource tags", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TAG_ID\tNAME\tKEY\tVALUE\tCREATED")
		for _, t := range result.ResourceTags {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				t.ResourceTagID, t.ResourceTagName, t.TagKey, t.TagValue, t.CreatedDateTime)
		}
		w.Flush()
	},
}
