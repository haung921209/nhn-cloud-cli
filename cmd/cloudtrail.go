package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/cloudtrail"
	"github.com/spf13/cobra"
)

// cloudtrailCmd represents the cloudtrail command
var cloudtrailCmd = &cobra.Command{
	Use:     "cloudtrail",
	Aliases: []string{"trail", "audit"},
	Short:   "Manage CloudTrail audit events",
	Long:    `View and search CloudTrail audit events for your NHN Cloud account.`,
}

// cloudtrailSearchCmd searches CloudTrail events
var cloudtrailSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search CloudTrail events",
	Long:  `Search CloudTrail events with various filters including time range, event source, and member type.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newCloudTrailClient()
		ctx := context.Background()

		// Parse time flags
		fromStr, _ := cmd.Flags().GetString("from")
		toStr, _ := cmd.Flags().GetString("to")

		// Default to last 24 hours if not specified
		var from, to time.Time
		var err error
		if fromStr != "" {
			from, err = time.Parse(time.RFC3339, fromStr)
			if err != nil {
				return fmt.Errorf("invalid from time format (use RFC3339): %w", err)
			}
		} else {
			from = time.Now().Add(-24 * time.Hour)
		}

		if toStr != "" {
			to, err = time.Parse(time.RFC3339, toStr)
			if err != nil {
				return fmt.Errorf("invalid to time format (use RFC3339): %w", err)
			}
		} else {
			to = time.Now()
		}

		// Build request
		input := &cloudtrail.SearchEventsInput{
			From: from,
			To:   to,
		}

		// Optional filters
		eventSource, _ := cmd.Flags().GetStringSlice("event-source")
		if len(eventSource) > 0 {
			input.EventSourceTypeList = eventSource
		}

		memberType, _ := cmd.Flags().GetStringSlice("member-type")
		if len(memberType) > 0 {
			input.MemberTypeList = memberType
		}

		memberID, _ := cmd.Flags().GetStringSlice("member-id")
		if len(memberID) > 0 {
			input.MemberIDList = memberID
		}

		eventID, _ := cmd.Flags().GetStringSlice("event-id")
		if len(eventID) > 0 {
			input.EventIDList = eventID
		}

		page, _ := cmd.Flags().GetInt("page")
		if page > 0 {
			input.Page = page
		}

		size, _ := cmd.Flags().GetInt("size")
		if size > 0 {
			input.Size = size
		} else {
			input.Size = 100 // Default size
		}

		result, err := client.SearchEvents(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to search events: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Body)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Total: %d events (Page %d, Size %d)\n\n", result.Body.TotalCount, result.Body.Page, result.Body.Size)
		fmt.Fprintln(w, "TIME\tTYPE\tSOURCE\tMEMBER\tIP\tPRODUCT")
		for _, event := range result.Body.Events {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				event.EventTime.Format("2006-01-02 15:04:05"),
				event.EventType,
				event.EventSourceType,
				event.MemberID,
				event.SourceIP,
				event.ProductID,
			)
		}
		return w.Flush()
	},
}

// cloudtrailRecentCmd shows recent CloudTrail events
var cloudtrailRecentCmd = &cobra.Command{
	Use:   "recent",
	Short: "Show recent CloudTrail events",
	Long:  `Show recent CloudTrail events from the last specified hours (default: 24 hours).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newCloudTrailClient()
		ctx := context.Background()

		hours, _ := cmd.Flags().GetInt("hours")
		if hours <= 0 {
			hours = 24
		}

		size, _ := cmd.Flags().GetInt("size")
		if size <= 0 {
			size = 50
		}

		from := time.Now().Add(-time.Duration(hours) * time.Hour)
		to := time.Now()

		result, err := client.SearchEventsSimple(ctx, from, to, 0, size)
		if err != nil {
			return fmt.Errorf("failed to get recent events: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Body)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Recent events (last %d hours): %d total\n\n", hours, result.Body.TotalCount)
		fmt.Fprintln(w, "TIME\tTYPE\tSOURCE\tMEMBER\tIP\tPRODUCT")
		for _, event := range result.Body.Events {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				event.EventTime.Format("2006-01-02 15:04:05"),
				event.EventType,
				event.EventSourceType,
				event.MemberID,
				event.SourceIP,
				event.ProductID,
			)
		}
		return w.Flush()
	},
}

// cloudtrailGetCmd gets a specific event by ID
var cloudtrailGetCmd = &cobra.Command{
	Use:   "get <event-id>",
	Short: "Get a specific CloudTrail event",
	Long:  `Get detailed information about a specific CloudTrail event by its ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		eventID := args[0]
		client := newCloudTrailClient()
		ctx := context.Background()

		// Search for the specific event
		input := &cloudtrail.SearchEventsInput{
			From:        time.Now().Add(-30 * 24 * time.Hour), // Last 30 days
			To:          time.Now(),
			EventIDList: []string{eventID},
			Size:        1,
		}

		result, err := client.SearchEvents(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to get event: %w", err)
		}

		if len(result.Body.Events) == 0 {
			return fmt.Errorf("event not found: %s", eventID)
		}

		event := result.Body.Events[0]

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(event)
		}

		fmt.Printf("Event ID:        %s\n", event.EventID)
		fmt.Printf("Event Time:      %s\n", event.EventTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("Event Type:      %s\n", event.EventType)
		fmt.Printf("Event Source:    %s\n", event.EventSourceType)
		fmt.Printf("Member Type:     %s\n", event.MemberType)
		fmt.Printf("Member ID:       %s\n", event.MemberID)
		fmt.Printf("Source IP:       %s\n", event.SourceIP)
		fmt.Printf("User Agent:      %s\n", event.UserAgent)
		fmt.Printf("Organization:    %s\n", event.OrgID)
		fmt.Printf("Project:         %s\n", event.ProjectID)
		fmt.Printf("Product:         %s\n", event.ProductID)
		fmt.Printf("Region:          %s\n", event.Region)
		if event.RequestID != "" {
			fmt.Printf("Request ID:      %s\n", event.RequestID)
		}
		if len(event.Resources) > 0 {
			fmt.Printf("Resources:\n")
			for _, r := range event.Resources {
				fmt.Printf("  - Type: %s, ID: %s, Name: %s\n", r.ResourceType, r.ResourceID, r.ResourceName)
			}
		}
		if event.Request != "" {
			fmt.Printf("Request:\n%s\n", event.Request)
		}
		if event.Response != "" {
			fmt.Printf("Response:\n%s\n", event.Response)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloudtrailCmd)
	cloudtrailCmd.AddCommand(cloudtrailSearchCmd)
	cloudtrailCmd.AddCommand(cloudtrailRecentCmd)
	cloudtrailCmd.AddCommand(cloudtrailGetCmd)

	// Search flags
	cloudtrailSearchCmd.Flags().String("from", "", "Start time (RFC3339 format, e.g., 2024-01-01T00:00:00Z)")
	cloudtrailSearchCmd.Flags().String("to", "", "End time (RFC3339 format)")
	cloudtrailSearchCmd.Flags().StringSlice("event-source", nil, "Event source types (CONSOLE, API)")
	cloudtrailSearchCmd.Flags().StringSlice("member-type", nil, "Member types (TOAST, IAM)")
	cloudtrailSearchCmd.Flags().StringSlice("member-id", nil, "Member IDs to filter")
	cloudtrailSearchCmd.Flags().StringSlice("event-id", nil, "Event IDs to filter")
	cloudtrailSearchCmd.Flags().Int("page", 0, "Page number")
	cloudtrailSearchCmd.Flags().Int("size", 100, "Page size")

	// Recent flags
	cloudtrailRecentCmd.Flags().Int("hours", 24, "Number of hours to look back")
	cloudtrailRecentCmd.Flags().Int("size", 50, "Number of events to return")
}

func newCloudTrailClient() *cloudtrail.Client {
	return cloudtrail.NewClient(getAppKey(), getAccessKey(), getSecretKey(), nil, debug)
}
