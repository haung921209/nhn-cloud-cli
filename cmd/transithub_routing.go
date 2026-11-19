package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/transithub"
	"github.com/spf13/cobra"
)

func init() {
	// Routing Tables
	transitHubCmd.AddCommand(thDescribeRoutingTablesCmd)
	transitHubCmd.AddCommand(thCreateRoutingTableCmd)
	transitHubCmd.AddCommand(thUpdateRoutingTableCmd)
	transitHubCmd.AddCommand(thDeleteRoutingTableCmd)

	// Associations
	transitHubCmd.AddCommand(thDescribeAssociationsCmd)
	transitHubCmd.AddCommand(thCreateAssociationCmd)
	transitHubCmd.AddCommand(thDeleteAssociationCmd)

	// Propagations
	transitHubCmd.AddCommand(thDescribePropagationsCmd)
	transitHubCmd.AddCommand(thCreatePropagationCmd)
	transitHubCmd.AddCommand(thDeletePropagationCmd)

	// Rules
	transitHubCmd.AddCommand(thDescribeRulesCmd)
	transitHubCmd.AddCommand(thCreateRuleCmd)
	transitHubCmd.AddCommand(thDeleteRuleCmd)

	// -- Flags for Routing Tables --
	thCreateRoutingTableCmd.Flags().String("name", "", "Name (required)")
	thCreateRoutingTableCmd.Flags().String("description", "", "Description")
	thCreateRoutingTableCmd.Flags().String("transit-hub-id", "", "Transit Hub ID (required)")
	thCreateRoutingTableCmd.MarkFlagRequired("name")
	thCreateRoutingTableCmd.MarkFlagRequired("transit-hub-id")

	thUpdateRoutingTableCmd.Flags().String("routing-table-id", "", "Routing Table ID (required)")
	thUpdateRoutingTableCmd.Flags().String("name", "", "Name")
	thUpdateRoutingTableCmd.Flags().String("description", "", "Description")
	thUpdateRoutingTableCmd.MarkFlagRequired("routing-table-id")

	thDeleteRoutingTableCmd.Flags().String("routing-table-id", "", "Routing Table ID (required)")
	thDeleteRoutingTableCmd.MarkFlagRequired("routing-table-id")

	// -- Flags for Associations --
	thCreateAssociationCmd.Flags().String("routing-table-id", "", "Routing Table ID (required)")
	thCreateAssociationCmd.Flags().String("attachment-id", "", "Attachment ID (required)")
	thCreateAssociationCmd.MarkFlagRequired("routing-table-id")
	thCreateAssociationCmd.MarkFlagRequired("attachment-id")

	thDeleteAssociationCmd.Flags().String("association-id", "", "Association ID (required)")
	thDeleteAssociationCmd.MarkFlagRequired("association-id")

	// -- Flags for Propagations --
	thCreatePropagationCmd.Flags().String("routing-table-id", "", "Routing Table ID (required)")
	thCreatePropagationCmd.Flags().String("attachment-id", "", "Attachment ID (required)")
	thCreatePropagationCmd.MarkFlagRequired("routing-table-id")
	thCreatePropagationCmd.MarkFlagRequired("attachment-id")

	thDeletePropagationCmd.Flags().String("propagation-id", "", "Propagation ID (required)")
	thDeletePropagationCmd.MarkFlagRequired("propagation-id")

	// -- Flags for Rules --
	thCreateRuleCmd.Flags().String("routing-table-id", "", "Routing Table ID (required)")
	thCreateRuleCmd.Flags().String("cidr", "", "Destination CIDR (required)")
	thCreateRuleCmd.Flags().String("attachment-id", "", "Attachment the traffic goes to (required)")
	thCreateRuleCmd.Flags().String("action", "", "FORWARD (default) or BLACKHOLE")
	thCreateRuleCmd.MarkFlagRequired("routing-table-id")
	thCreateRuleCmd.MarkFlagRequired("cidr")
	thCreateRuleCmd.MarkFlagRequired("attachment-id")

	thDeleteRuleCmd.Flags().String("rule-id", "", "Rule ID (required)")
	thDeleteRuleCmd.MarkFlagRequired("rule-id")
}

// -----------------------------------------------------------------------------
// Routing Tables
// -----------------------------------------------------------------------------

var thDescribeRoutingTablesCmd = &cobra.Command{
	Use:     "describe-routing-tables",
	Aliases: []string{"list-routing-tables"},
	Short:   "List routing tables",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()

		result, err := client.ListRoutingTables(ctx)
		if err != nil {
			exitWithError("Failed to list routing tables", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTRANSIT_HUB_ID\tDEFAULT\tSTATUS")
		for _, rt := range result.RoutingTables {
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n", rt.ID, rt.Name, rt.TransitHubID, rt.DefaultTable, rt.Status)
		}
		w.Flush()
	},
}

var thCreateRoutingTableCmd = &cobra.Command{
	Use:   "create-routing-table",
	Short: "Create a new routing table",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		thID, _ := cmd.Flags().GetString("transit-hub-id")

		input := &transithub.CreateRoutingTableInput{
			Name:         name,
			Description:  desc,
			TransitHubID: thID,
		}

		result, err := client.CreateRoutingTable(ctx, input)
		if err != nil {
			exitWithError("Failed to create routing table", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Routing Table created: %s (%s)\n", result.RoutingTable.Name, result.RoutingTable.ID)
	},
}

var thUpdateRoutingTableCmd = &cobra.Command{
	Use:   "update-routing-table",
	Short: "Update a routing table",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("routing-table-id")
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")

		input := &transithub.UpdateRoutingTableInput{
			Name:        name,
			Description: desc,
		}

		result, err := client.UpdateRoutingTable(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update routing table", err)
		}

		if err := printWriteResult(result, "Routing Table updated: %s\n", result.RoutingTable.ID); err != nil {
			exitWithError("Failed to print output", err)
		}
	},
}

var thDeleteRoutingTableCmd = &cobra.Command{
	Use:   "delete-routing-table",
	Short: "Delete a routing table",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("routing-table-id")

		if err := client.DeleteRoutingTable(ctx, id); err != nil {
			exitWithError("Failed to delete routing table", err)
		}

		fmt.Printf("Routing Table %s deleted\n", id)
	},
}

// -----------------------------------------------------------------------------
// Associations
// -----------------------------------------------------------------------------

var thDescribeAssociationsCmd = &cobra.Command{
	Use:     "describe-associations",
	Aliases: []string{"list-associations"},
	Short:   "List routing associations",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()

		result, err := client.ListRoutingAssociations(ctx)
		if err != nil {
			exitWithError("Failed to list associations", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tROUTING_TABLE_ID\tATTACHMENT_ID\tSTATUS")
		for _, a := range result.Associations {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", a.ID, a.RoutingTableID, a.AttachmentID, a.Status)
		}
		w.Flush()
	},
}

var thCreateAssociationCmd = &cobra.Command{
	Use:   "create-association",
	Short: "Create a new routing association",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		rtID, _ := cmd.Flags().GetString("routing-table-id")
		attID, _ := cmd.Flags().GetString("attachment-id")

		input := &transithub.CreateRoutingAssociationInput{
			RoutingTableID: rtID,
			AttachmentID:   attID,
		}

		result, err := client.CreateRoutingAssociation(ctx, input)
		if err != nil {
			exitWithError("Failed to create association", err)
		}

		if err := printWriteResult(result, "Association created: %s\n", result.Association.ID); err != nil {
			exitWithError("Failed to print output", err)
		}
	},
}

var thDeleteAssociationCmd = &cobra.Command{
	Use:   "delete-association",
	Short: "Delete an association",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("association-id")

		if err := client.DeleteRoutingAssociation(ctx, id); err != nil {
			exitWithError("Failed to delete association", err)
		}

		fmt.Printf("Association %s deleted\n", id)
	},
}

// -----------------------------------------------------------------------------
// Propagations
// -----------------------------------------------------------------------------

var thDescribePropagationsCmd = &cobra.Command{
	Use:     "describe-propagations",
	Aliases: []string{"list-propagations"},
	Short:   "List routing propagations",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()

		result, err := client.ListRoutingPropagations(ctx)
		if err != nil {
			exitWithError("Failed to list propagations", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tROUTING_TABLE_ID\tATTACHMENT_ID\tSTATUS")
		for _, p := range result.Propagations {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.ID, p.RoutingTableID, p.AttachmentID, p.Status)
		}
		w.Flush()
	},
}

var thCreatePropagationCmd = &cobra.Command{
	Use:   "create-propagation",
	Short: "Create a new routing propagation",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		rtID, _ := cmd.Flags().GetString("routing-table-id")
		attID, _ := cmd.Flags().GetString("attachment-id")

		input := &transithub.CreateRoutingPropagationInput{
			RoutingTableID: rtID,
			AttachmentID:   attID,
		}

		result, err := client.CreateRoutingPropagation(ctx, input)
		if err != nil {
			exitWithError("Failed to create propagation", err)
		}

		if err := printWriteResult(result, "Propagation created: %s\n", result.Propagation.ID); err != nil {
			exitWithError("Failed to print output", err)
		}
	},
}

var thDeletePropagationCmd = &cobra.Command{
	Use:   "delete-propagation",
	Short: "Delete a propagation",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("propagation-id")

		if err := client.DeleteRoutingPropagation(ctx, id); err != nil {
			exitWithError("Failed to delete propagation", err)
		}

		fmt.Printf("Propagation %s deleted\n", id)
	},
}

// -----------------------------------------------------------------------------
// Rules
// -----------------------------------------------------------------------------

var thDescribeRulesCmd = &cobra.Command{
	Use:     "describe-rules",
	Aliases: []string{"list-rules"},
	Short:   "List routing rules",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()

		result, err := client.ListRoutingRules(ctx)
		if err != nil {
			exitWithError("Failed to list rules", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tCIDR\tACTION\tATTACHMENT_ID\tTYPE")
		for _, r := range result.Rules {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", r.ID, r.CIDR, r.Action, r.AttachmentID, r.RuleType)
		}
		w.Flush()
	},
}

var thCreateRuleCmd = &cobra.Command{
	Use:   "create-rule",
	Short: "Create a new routing rule",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		rtID, _ := cmd.Flags().GetString("routing-table-id")
		cidr, _ := cmd.Flags().GetString("cidr")
		attID, _ := cmd.Flags().GetString("attachment-id")
		action, _ := cmd.Flags().GetString("action")

		input := &transithub.CreateRoutingRuleInput{
			RoutingTableID: rtID,
			CIDR:           cidr,
			Action:         action,
			AttachmentID:   attID,
		}

		result, err := client.CreateRoutingRule(ctx, input)
		if err != nil {
			exitWithError("Failed to create rule", err)
		}

		if err := printWriteResult(result, "Rule created: %s\n", result.Rule.ID); err != nil {
			exitWithError("Failed to print output", err)
		}
	},
}

var thDeleteRuleCmd = &cobra.Command{
	Use:   "delete-rule",
	Short: "Delete a routing rule",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("rule-id")

		if err := client.DeleteRoutingRule(ctx, id); err != nil {
			exitWithError("Failed to delete rule", err)
		}

		fmt.Printf("Rule %s deleted\n", id)
	},
}
