package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/mirroring"
	"github.com/spf13/cobra"
)

func init() {
	mirroringCmd.AddCommand(mirroringDescribeSessionsCmd)
	mirroringCmd.AddCommand(mirroringGetSessionCmd)
	mirroringCmd.AddCommand(mirroringCreateSessionCmd)
	mirroringCmd.AddCommand(mirroringUpdateSessionCmd)
	mirroringCmd.AddCommand(mirroringDeleteSessionCmd)

	mirroringGetSessionCmd.Flags().String("session-id", "", "Session ID (required)")
	mirroringGetSessionCmd.MarkFlagRequired("session-id")

	mirroringCreateSessionCmd.Flags().String("name", "", "Session name (required)")
	mirroringCreateSessionCmd.Flags().String("source-type", "PORT", "Source type: PORT, LOADBALANCER_MEMBER")
	mirroringCreateSessionCmd.Flags().String("source-id", "", "Source ID (required)")
	mirroringCreateSessionCmd.Flags().String("target-type", "PORT", "Target type: PORT")
	mirroringCreateSessionCmd.Flags().String("target-id", "", "Target ID (required)")
	mirroringCreateSessionCmd.Flags().String("direction", "both", "Direction: in, out, both")
	mirroringCreateSessionCmd.Flags().String("filter-group-id", "", "Filter group ID")
	mirroringCreateSessionCmd.Flags().String("description", "", "Description")
	mirroringCreateSessionCmd.MarkFlagRequired("name")
	mirroringCreateSessionCmd.MarkFlagRequired("source-id")
	mirroringCreateSessionCmd.MarkFlagRequired("target-id")

	mirroringUpdateSessionCmd.Flags().String("session-id", "", "Session ID (required)")
	mirroringUpdateSessionCmd.Flags().String("name", "", "Session name")
	mirroringUpdateSessionCmd.Flags().String("direction", "", "Direction: in, out, both")
	mirroringUpdateSessionCmd.Flags().String("filter-group-id", "", "Filter group ID")
	mirroringUpdateSessionCmd.Flags().String("description", "", "Description")
	mirroringUpdateSessionCmd.MarkFlagRequired("session-id")

	mirroringDeleteSessionCmd.Flags().String("session-id", "", "Session ID (required)")
	mirroringDeleteSessionCmd.MarkFlagRequired("session-id")
}

var mirroringDescribeSessionsCmd = &cobra.Command{
	Use:     "describe-sessions",
	Aliases: []string{"list-sessions"},
	Short:   "List all mirroring sessions",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		result, err := client.ListSessions(context.Background())
		if err != nil {
			exitWithError("Failed to list mirroring sessions", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSOURCE_TYPE\tDIRECTION\tSTATE")
		for _, s := range result.MirroringSessions {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				s.ID, s.Name, s.SourceType, s.Direction, s.State)
		}
		w.Flush()
	},
}

var mirroringGetSessionCmd = &cobra.Command{
	Use:     "describe-session",
	Aliases: []string{"get-session"},
	Short:   "Get mirroring session details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("session-id")

		result, err := client.GetSession(ctx, id)
		if err != nil {
			exitWithError("Failed to get mirroring session", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		s := result.MirroringSession
		fmt.Printf("ID:              %s\n", s.ID)
		fmt.Printf("Name:            %s\n", s.Name)
		fmt.Printf("Description:     %s\n", s.Description)
		fmt.Printf("Source Type:     %s\n", s.SourceType)
		fmt.Printf("Source ID:       %s\n", s.SourceID)
		fmt.Printf("Target Type:     %s\n", s.TargetType)
		fmt.Printf("Target ID:       %s\n", s.TargetID)
		fmt.Printf("Direction:       %s\n", s.Direction)
		fmt.Printf("Filter Group ID: %s\n", s.FilterGroupID)
		fmt.Printf("Admin State Up:  %v\n", s.AdminStateUp)
		fmt.Printf("State:           %s\n", s.State)
		fmt.Printf("Created At:      %s\n", s.CreatedAt.Format("2006-01-02 15:04:05"))
	},
}

var mirroringCreateSessionCmd = &cobra.Command{
	Use:   "create-session",
	Short: "Create a new mirroring session",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		name, _ := cmd.Flags().GetString("name")
		sourceType, _ := cmd.Flags().GetString("source-type")
		sourceID, _ := cmd.Flags().GetString("source-id")
		targetType, _ := cmd.Flags().GetString("target-type")
		targetID, _ := cmd.Flags().GetString("target-id")
		direction, _ := cmd.Flags().GetString("direction")
		filterGroupID, _ := cmd.Flags().GetString("filter-group-id")
		description, _ := cmd.Flags().GetString("description")

		input := &mirroring.CreateSessionInput{
			Name:          name,
			SourceType:    sourceType,
			SourceID:      sourceID,
			TargetType:    targetType,
			TargetID:      targetID,
			Direction:     direction,
			FilterGroupID: filterGroupID,
			Description:   description,
			AdminStateUp:  true,
		}

		result, err := client.CreateSession(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create mirroring session", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Mirroring session created: %s\n", result.MirroringSession.ID)
		fmt.Printf("Name: %s\n", result.MirroringSession.Name)
	},
}

var mirroringUpdateSessionCmd = &cobra.Command{
	Use:   "update-session",
	Short: "Update a mirroring session",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		id, _ := cmd.Flags().GetString("session-id")
		name, _ := cmd.Flags().GetString("name")
		direction, _ := cmd.Flags().GetString("direction")
		filterGroupID, _ := cmd.Flags().GetString("filter-group-id")
		description, _ := cmd.Flags().GetString("description")

		input := &mirroring.UpdateSessionInput{
			Name:          name,
			Direction:     direction,
			FilterGroupID: filterGroupID,
			Description:   description,
		}

		result, err := client.UpdateSession(context.Background(), id, input)
		if err != nil {
			exitWithError("Failed to update mirroring session", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Mirroring session updated: %s\n", result.MirroringSession.ID)
	},
}

var mirroringDeleteSessionCmd = &cobra.Command{
	Use:   "delete-session",
	Short: "Delete a mirroring session",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		id, _ := cmd.Flags().GetString("session-id")
		if err := client.DeleteSession(context.Background(), id); err != nil {
			exitWithError("Failed to delete mirroring session", err)
		}
		fmt.Printf("Mirroring session %s deleted\n", id)
	},
}
