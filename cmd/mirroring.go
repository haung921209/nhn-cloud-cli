package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/mirroring"
	"github.com/spf13/cobra"
)

var mirroringCmd = &cobra.Command{
	Use:     "mirroring",
	Aliases: []string{"mirror", "traffic-mirroring"},
	Short:   "Manage Traffic Mirroring",
}

var mirroringSessionCmd = &cobra.Command{
	Use:     "session",
	Aliases: []string{"sessions"},
	Short:   "Manage mirroring sessions",
}

var mirroringFilterGroupCmd = &cobra.Command{
	Use:     "filter-group",
	Aliases: []string{"filter-groups", "fg"},
	Short:   "Manage mirroring filter groups",
}

var mirroringFilterCmd = &cobra.Command{
	Use:     "filter",
	Aliases: []string{"filters"},
	Short:   "Manage mirroring filters",
}

func init() {
	rootCmd.AddCommand(mirroringCmd)

	mirroringCmd.AddCommand(mirroringSessionCmd)
	mirroringSessionCmd.AddCommand(mirroringSessionListCmd)
	mirroringSessionCmd.AddCommand(mirroringSessionGetCmd)
	mirroringSessionCmd.AddCommand(mirroringSessionCreateCmd)
	mirroringSessionCmd.AddCommand(mirroringSessionUpdateCmd)
	mirroringSessionCmd.AddCommand(mirroringSessionDeleteCmd)

	mirroringCmd.AddCommand(mirroringFilterGroupCmd)
	mirroringFilterGroupCmd.AddCommand(mirroringFilterGroupListCmd)
	mirroringFilterGroupCmd.AddCommand(mirroringFilterGroupGetCmd)
	mirroringFilterGroupCmd.AddCommand(mirroringFilterGroupCreateCmd)
	mirroringFilterGroupCmd.AddCommand(mirroringFilterGroupUpdateCmd)
	mirroringFilterGroupCmd.AddCommand(mirroringFilterGroupDeleteCmd)

	mirroringCmd.AddCommand(mirroringFilterCmd)
	mirroringFilterCmd.AddCommand(mirroringFilterListCmd)
	mirroringFilterCmd.AddCommand(mirroringFilterGetCmd)
	mirroringFilterCmd.AddCommand(mirroringFilterCreateCmd)
	mirroringFilterCmd.AddCommand(mirroringFilterUpdateCmd)
	mirroringFilterCmd.AddCommand(mirroringFilterDeleteCmd)

	mirroringSessionCreateCmd.Flags().String("name", "", "Session name (required)")
	mirroringSessionCreateCmd.Flags().String("source-type", "PORT", "Source type: PORT, LOADBALANCER_MEMBER")
	mirroringSessionCreateCmd.Flags().String("source-id", "", "Source ID (required)")
	mirroringSessionCreateCmd.Flags().String("target-type", "PORT", "Target type: PORT")
	mirroringSessionCreateCmd.Flags().String("target-id", "", "Target ID (required)")
	mirroringSessionCreateCmd.Flags().String("direction", "both", "Direction: in, out, both")
	mirroringSessionCreateCmd.Flags().String("filter-group-id", "", "Filter group ID")
	mirroringSessionCreateCmd.Flags().String("description", "", "Description")
	mirroringSessionCreateCmd.MarkFlagRequired("name")
	mirroringSessionCreateCmd.MarkFlagRequired("source-id")
	mirroringSessionCreateCmd.MarkFlagRequired("target-id")

	mirroringSessionUpdateCmd.Flags().String("name", "", "Session name")
	mirroringSessionUpdateCmd.Flags().String("direction", "", "Direction: in, out, both")
	mirroringSessionUpdateCmd.Flags().String("filter-group-id", "", "Filter group ID")
	mirroringSessionUpdateCmd.Flags().String("description", "", "Description")

	mirroringFilterGroupCreateCmd.Flags().String("name", "", "Filter group name (required)")
	mirroringFilterGroupCreateCmd.Flags().String("description", "", "Description")
	mirroringFilterGroupCreateCmd.MarkFlagRequired("name")

	mirroringFilterGroupUpdateCmd.Flags().String("name", "", "Filter group name")
	mirroringFilterGroupUpdateCmd.Flags().String("description", "", "Description")

	mirroringFilterCreateCmd.Flags().String("filter-group-id", "", "Filter group ID (required)")
	mirroringFilterCreateCmd.Flags().String("name", "", "Filter name (required)")
	mirroringFilterCreateCmd.Flags().String("protocol", "", "Protocol: tcp, udp, icmp")
	mirroringFilterCreateCmd.Flags().String("source-cidr", "", "Source CIDR")
	mirroringFilterCreateCmd.Flags().String("dest-cidr", "", "Destination CIDR")
	mirroringFilterCreateCmd.Flags().Int("source-port-min", 0, "Source port min")
	mirroringFilterCreateCmd.Flags().Int("source-port-max", 0, "Source port max")
	mirroringFilterCreateCmd.Flags().Int("dest-port-min", 0, "Destination port min")
	mirroringFilterCreateCmd.Flags().Int("dest-port-max", 0, "Destination port max")
	mirroringFilterCreateCmd.Flags().String("action", "accept", "Action: accept, drop")
	mirroringFilterCreateCmd.Flags().String("description", "", "Description")
	mirroringFilterCreateCmd.MarkFlagRequired("filter-group-id")
	mirroringFilterCreateCmd.MarkFlagRequired("name")

	mirroringFilterUpdateCmd.Flags().String("name", "", "Filter name")
	mirroringFilterUpdateCmd.Flags().String("protocol", "", "Protocol")
	mirroringFilterUpdateCmd.Flags().String("source-cidr", "", "Source CIDR")
	mirroringFilterUpdateCmd.Flags().String("dest-cidr", "", "Destination CIDR")
	mirroringFilterUpdateCmd.Flags().String("action", "", "Action: accept, drop")
	mirroringFilterUpdateCmd.Flags().String("description", "", "Description")
}

func newMirroringClient() *mirroring.Client {
	return mirroring.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

var mirroringSessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all mirroring sessions",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		result, err := client.ListSessions(context.Background())
		if err != nil {
			exitWithError("Failed to list mirroring sessions", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var mirroringSessionGetCmd = &cobra.Command{
	Use:   "get [session-id]",
	Short: "Get mirroring session details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		result, err := client.GetSession(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get mirroring session", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var mirroringSessionCreateCmd = &cobra.Command{
	Use:   "create",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Mirroring session created: %s\n", result.MirroringSession.ID)
		fmt.Printf("Name: %s\n", result.MirroringSession.Name)
	},
}

var mirroringSessionUpdateCmd = &cobra.Command{
	Use:   "update [session-id]",
	Short: "Update a mirroring session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()

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

		result, err := client.UpdateSession(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update mirroring session", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Mirroring session updated: %s\n", result.MirroringSession.ID)
	},
}

var mirroringSessionDeleteCmd = &cobra.Command{
	Use:   "delete [session-id]",
	Short: "Delete a mirroring session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		if err := client.DeleteSession(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete mirroring session", err)
		}
		fmt.Printf("Mirroring session %s deleted\n", args[0])
	},
}

var mirroringFilterGroupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all mirroring filter groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		result, err := client.ListFilterGroups(context.Background())
		if err != nil {
			exitWithError("Failed to list filter groups", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATE\tFILTER_COUNT")
		for _, fg := range result.MirroringFilterGroups {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
				fg.ID, fg.Name, fg.State, len(fg.FilterIDs))
		}
		w.Flush()
	},
}

var mirroringFilterGroupGetCmd = &cobra.Command{
	Use:   "get [filter-group-id]",
	Short: "Get filter group details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		result, err := client.GetFilterGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get filter group", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fg := result.MirroringFilterGroup
		fmt.Printf("ID:          %s\n", fg.ID)
		fmt.Printf("Name:        %s\n", fg.Name)
		fmt.Printf("Description: %s\n", fg.Description)
		fmt.Printf("State:       %s\n", fg.State)
		fmt.Printf("Filter IDs:  %v\n", fg.FilterIDs)
		fmt.Printf("Created At:  %s\n", fg.CreatedAt.Format("2006-01-02 15:04:05"))
	},
}

var mirroringFilterGroupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new filter group",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &mirroring.CreateFilterGroupInput{
			Name:        name,
			Description: description,
		}

		result, err := client.CreateFilterGroup(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create filter group", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Filter group created: %s\n", result.MirroringFilterGroup.ID)
		fmt.Printf("Name: %s\n", result.MirroringFilterGroup.Name)
	},
}

var mirroringFilterGroupUpdateCmd = &cobra.Command{
	Use:   "update [filter-group-id]",
	Short: "Update a filter group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &mirroring.UpdateFilterGroupInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateFilterGroup(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update filter group", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Filter group updated: %s\n", result.MirroringFilterGroup.ID)
	},
}

var mirroringFilterGroupDeleteCmd = &cobra.Command{
	Use:   "delete [filter-group-id]",
	Short: "Delete a filter group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		if err := client.DeleteFilterGroup(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete filter group", err)
		}
		fmt.Printf("Filter group %s deleted\n", args[0])
	},
}

var mirroringFilterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all mirroring filters",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		result, err := client.ListFilters(context.Background())
		if err != nil {
			exitWithError("Failed to list filters", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tPROTOCOL\tACTION\tSTATE")
		for _, f := range result.MirroringFilters {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				f.ID, f.Name, f.Protocol, f.Action, f.State)
		}
		w.Flush()
	},
}

var mirroringFilterGetCmd = &cobra.Command{
	Use:   "get [filter-id]",
	Short: "Get filter details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		result, err := client.GetFilter(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get filter", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		f := result.MirroringFilter
		fmt.Printf("ID:              %s\n", f.ID)
		fmt.Printf("Name:            %s\n", f.Name)
		fmt.Printf("Description:     %s\n", f.Description)
		fmt.Printf("Filter Group ID: %s\n", f.FilterGroupID)
		fmt.Printf("Protocol:        %s\n", f.Protocol)
		fmt.Printf("Source CIDR:     %s\n", f.SourceCIDR)
		fmt.Printf("Dest CIDR:       %s\n", f.DestCIDR)
		fmt.Printf("Source Port:     %d-%d\n", f.SourcePortMin, f.SourcePortMax)
		fmt.Printf("Dest Port:       %d-%d\n", f.DestPortMin, f.DestPortMax)
		fmt.Printf("Action:          %s\n", f.Action)
		fmt.Printf("State:           %s\n", f.State)
		fmt.Printf("Created At:      %s\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
	},
}

var mirroringFilterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new filter",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()

		filterGroupID, _ := cmd.Flags().GetString("filter-group-id")
		name, _ := cmd.Flags().GetString("name")
		protocol, _ := cmd.Flags().GetString("protocol")
		sourceCIDR, _ := cmd.Flags().GetString("source-cidr")
		destCIDR, _ := cmd.Flags().GetString("dest-cidr")
		sourcePortMin, _ := cmd.Flags().GetInt("source-port-min")
		sourcePortMax, _ := cmd.Flags().GetInt("source-port-max")
		destPortMin, _ := cmd.Flags().GetInt("dest-port-min")
		destPortMax, _ := cmd.Flags().GetInt("dest-port-max")
		action, _ := cmd.Flags().GetString("action")
		description, _ := cmd.Flags().GetString("description")

		input := &mirroring.CreateFilterInput{
			FilterGroupID: filterGroupID,
			Name:          name,
			Protocol:      protocol,
			SourceCIDR:    sourceCIDR,
			DestCIDR:      destCIDR,
			SourcePortMin: sourcePortMin,
			SourcePortMax: sourcePortMax,
			DestPortMin:   destPortMin,
			DestPortMax:   destPortMax,
			Action:        action,
			Description:   description,
		}

		result, err := client.CreateFilter(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create filter", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Filter created: %s\n", result.MirroringFilter.ID)
		fmt.Printf("Name: %s\n", result.MirroringFilter.Name)
	},
}

var mirroringFilterUpdateCmd = &cobra.Command{
	Use:   "update [filter-id]",
	Short: "Update a filter",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()

		name, _ := cmd.Flags().GetString("name")
		protocol, _ := cmd.Flags().GetString("protocol")
		sourceCIDR, _ := cmd.Flags().GetString("source-cidr")
		destCIDR, _ := cmd.Flags().GetString("dest-cidr")
		action, _ := cmd.Flags().GetString("action")
		description, _ := cmd.Flags().GetString("description")

		input := &mirroring.UpdateFilterInput{
			Name:        name,
			Protocol:    protocol,
			SourceCIDR:  sourceCIDR,
			DestCIDR:    destCIDR,
			Action:      action,
			Description: description,
		}

		result, err := client.UpdateFilter(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update filter", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Filter updated: %s\n", result.MirroringFilter.ID)
	},
}

var mirroringFilterDeleteCmd = &cobra.Command{
	Use:   "delete [filter-id]",
	Short: "Delete a filter",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMirroringClient()
		if err := client.DeleteFilter(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete filter", err)
		}
		fmt.Printf("Filter %s deleted\n", args[0])
	},
}
