package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/transithub"
	"github.com/spf13/cobra"
)

var transitHubCmd = &cobra.Command{
	Use:     "transit-hub",
	Aliases: []string{"th", "transithub"},
	Short:   "Manage Transit Hubs for multi-VPC networking",
}

// Subcommand groups
var thAttachmentCmd = &cobra.Command{Use: "attachment", Aliases: []string{"attachments", "att"}, Short: "Manage attachments (VPC connections)"}
var thRoutingTableCmd = &cobra.Command{Use: "routing-table", Aliases: []string{"rt", "routing-tables"}, Short: "Manage routing tables"}
var thAssociationCmd = &cobra.Command{Use: "association", Aliases: []string{"associations", "assoc"}, Short: "Manage routing associations"}
var thPropagationCmd = &cobra.Command{Use: "propagation", Aliases: []string{"propagations", "prop"}, Short: "Manage routing propagations"}
var thRuleCmd = &cobra.Command{Use: "rule", Aliases: []string{"rules"}, Short: "Manage routing rules"}
var thMulticastCmd = &cobra.Command{Use: "multicast", Aliases: []string{"multicast-domain", "mc"}, Short: "Manage multicast domains"}

func init() {
	rootCmd.AddCommand(transitHubCmd)

	// Transit Hub main commands
	transitHubCmd.AddCommand(thListCmd, thGetCmd, thCreateCmd, thUpdateCmd, thDeleteCmd)

	// Subcommand groups
	transitHubCmd.AddCommand(thAttachmentCmd, thRoutingTableCmd, thAssociationCmd, thPropagationCmd, thRuleCmd, thMulticastCmd)

	// Attachment commands
	thAttachmentCmd.AddCommand(thAttListCmd, thAttGetCmd, thAttCreateCmd, thAttUpdateCmd, thAttDeleteCmd)

	// Routing Table commands
	thRoutingTableCmd.AddCommand(thRTListCmd, thRTGetCmd, thRTCreateCmd, thRTUpdateCmd, thRTDeleteCmd)

	// Association commands
	thAssociationCmd.AddCommand(thAssocListCmd, thAssocGetCmd, thAssocCreateCmd, thAssocDeleteCmd)

	// Propagation commands
	thPropagationCmd.AddCommand(thPropListCmd, thPropGetCmd, thPropCreateCmd, thPropDeleteCmd)

	// Rule commands
	thRuleCmd.AddCommand(thRuleListCmd, thRuleGetCmd, thRuleCreateCmd, thRuleUpdateCmd, thRuleDeleteCmd)

	// Multicast commands
	thMulticastCmd.AddCommand(thMCListCmd, thMCGetCmd, thMCCreateCmd, thMCUpdateCmd, thMCDeleteCmd)

	// Flags for Transit Hub
	thCreateCmd.Flags().String("name", "", "Name (required)")
	thCreateCmd.Flags().String("description", "", "Description")
	thCreateCmd.MarkFlagRequired("name")
	thUpdateCmd.Flags().String("name", "", "Name")
	thUpdateCmd.Flags().String("description", "", "Description")

	// Flags for Attachment
	thAttCreateCmd.Flags().String("name", "", "Name (required)")
	thAttCreateCmd.Flags().String("description", "", "Description")
	thAttCreateCmd.Flags().String("transit-hub-id", "", "Transit Hub ID (required)")
	thAttCreateCmd.Flags().String("resource-type", "VPC", "Resource type (VPC, VPN)")
	thAttCreateCmd.Flags().String("resource-id", "", "Resource ID (required)")
	thAttCreateCmd.MarkFlagRequired("name")
	thAttCreateCmd.MarkFlagRequired("transit-hub-id")
	thAttCreateCmd.MarkFlagRequired("resource-id")
	thAttUpdateCmd.Flags().String("name", "", "Name")
	thAttUpdateCmd.Flags().String("description", "", "Description")

	// Flags for Routing Table
	thRTCreateCmd.Flags().String("name", "", "Name (required)")
	thRTCreateCmd.Flags().String("description", "", "Description")
	thRTCreateCmd.Flags().String("transit-hub-id", "", "Transit Hub ID (required)")
	thRTCreateCmd.MarkFlagRequired("name")
	thRTCreateCmd.MarkFlagRequired("transit-hub-id")
	thRTUpdateCmd.Flags().String("name", "", "Name")
	thRTUpdateCmd.Flags().String("description", "", "Description")

	// Flags for Association
	thAssocCreateCmd.Flags().String("routing-table-id", "", "Routing Table ID (required)")
	thAssocCreateCmd.Flags().String("attachment-id", "", "Attachment ID (required)")
	thAssocCreateCmd.MarkFlagRequired("routing-table-id")
	thAssocCreateCmd.MarkFlagRequired("attachment-id")

	// Flags for Propagation
	thPropCreateCmd.Flags().String("routing-table-id", "", "Routing Table ID (required)")
	thPropCreateCmd.Flags().String("attachment-id", "", "Attachment ID (required)")
	thPropCreateCmd.MarkFlagRequired("routing-table-id")
	thPropCreateCmd.MarkFlagRequired("attachment-id")

	// Flags for Rule
	thRuleCreateCmd.Flags().String("routing-table-id", "", "Routing Table ID (required)")
	thRuleCreateCmd.Flags().String("destination", "", "Destination CIDR (required)")
	thRuleCreateCmd.Flags().String("target-type", "", "Target type: ATTACHMENT, BLACKHOLE (required)")
	thRuleCreateCmd.Flags().String("target-id", "", "Target ID (for ATTACHMENT type)")
	thRuleCreateCmd.MarkFlagRequired("routing-table-id")
	thRuleCreateCmd.MarkFlagRequired("destination")
	thRuleCreateCmd.MarkFlagRequired("target-type")
	thRuleUpdateCmd.Flags().String("destination", "", "Destination CIDR")
	thRuleUpdateCmd.Flags().String("target-type", "", "Target type")
	thRuleUpdateCmd.Flags().String("target-id", "", "Target ID")

	// Flags for Multicast
	thMCCreateCmd.Flags().String("name", "", "Name (required)")
	thMCCreateCmd.Flags().String("description", "", "Description")
	thMCCreateCmd.Flags().String("transit-hub-id", "", "Transit Hub ID (required)")
	thMCCreateCmd.MarkFlagRequired("name")
	thMCCreateCmd.MarkFlagRequired("transit-hub-id")
	thMCUpdateCmd.Flags().String("name", "", "Name")
	thMCUpdateCmd.Flags().String("description", "", "Description")
}

func newTransitHubClient() *transithub.Client {
	return transithub.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

// =============================================================================
// Transit Hub Commands
// =============================================================================

var thListCmd = &cobra.Command{
	Use: "list", Short: "List all transit hubs",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.ListTransitHubs(context.Background())
		if err != nil {
			exitWithError("Failed to list transit hubs", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tSTATE")
		for _, th := range result.TransitHubs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", th.ID, th.Name, th.Status, th.State)
		}
		w.Flush()
	},
}

var thGetCmd = &cobra.Command{
	Use: "get [id]", Short: "Get transit hub details", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.GetTransitHub(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get transit hub", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		th := result.TransitHub
		fmt.Printf("ID:          %s\nName:        %s\nDescription: %s\nStatus:      %s\nState:       %s\nDefault RT:  %s\nCreated:     %s\n",
			th.ID, th.Name, th.Description, th.Status, th.State, th.DefaultRoutingTable, th.CreatedAt)
	},
}

var thCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a new transit hub",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		result, err := client.CreateTransitHub(context.Background(), &transithub.CreateTransitHubInput{Name: name, Description: desc})
		if err != nil {
			exitWithError("Failed to create transit hub", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Transit hub created: %s (%s)\n", result.TransitHub.ID, result.TransitHub.Name)
	},
}

var thUpdateCmd = &cobra.Command{
	Use: "update [id]", Short: "Update a transit hub", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		result, err := client.UpdateTransitHub(context.Background(), args[0], &transithub.UpdateTransitHubInput{Name: name, Description: desc})
		if err != nil {
			exitWithError("Failed to update transit hub", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Transit hub updated: %s\n", result.TransitHub.ID)
	},
}

var thDeleteCmd = &cobra.Command{
	Use: "delete [id]", Short: "Delete a transit hub", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		if err := client.DeleteTransitHub(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete transit hub", err)
		}
		fmt.Printf("Transit hub %s deleted\n", args[0])
	},
}

// =============================================================================
// Attachment Commands
// =============================================================================

var thAttListCmd = &cobra.Command{
	Use: "list", Short: "List all attachments",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.ListAttachments(context.Background())
		if err != nil {
			exitWithError("Failed to list attachments", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTRANSIT_HUB_ID\tRESOURCE_TYPE\tSTATUS")
		for _, a := range result.Attachments {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", a.ID, a.Name, a.TransitHubID, a.ResourceType, a.Status)
		}
		w.Flush()
	},
}

var thAttGetCmd = &cobra.Command{
	Use: "get [id]", Short: "Get attachment details", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.GetAttachment(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get attachment", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		a := result.Attachment
		fmt.Printf("ID:            %s\nName:          %s\nTransit Hub:   %s\nResource Type: %s\nResource ID:   %s\nStatus:        %s\nState:         %s\n",
			a.ID, a.Name, a.TransitHubID, a.ResourceType, a.ResourceID, a.Status, a.State)
	},
}

var thAttCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a new attachment",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		thID, _ := cmd.Flags().GetString("transit-hub-id")
		rType, _ := cmd.Flags().GetString("resource-type")
		rID, _ := cmd.Flags().GetString("resource-id")
		result, err := client.CreateAttachment(context.Background(), &transithub.CreateAttachmentInput{Name: name, Description: desc, TransitHubID: thID, ResourceType: rType, ResourceID: rID})
		if err != nil {
			exitWithError("Failed to create attachment", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Attachment created: %s\n", result.Attachment.ID)
	},
}

var thAttUpdateCmd = &cobra.Command{
	Use: "update [id]", Short: "Update an attachment", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		result, err := client.UpdateAttachment(context.Background(), args[0], &transithub.UpdateAttachmentInput{Name: name, Description: desc})
		if err != nil {
			exitWithError("Failed to update attachment", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Attachment updated: %s\n", result.Attachment.ID)
	},
}

var thAttDeleteCmd = &cobra.Command{
	Use: "delete [id]", Short: "Delete an attachment", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		if err := client.DeleteAttachment(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete attachment", err)
		}
		fmt.Printf("Attachment %s deleted\n", args[0])
	},
}

// =============================================================================
// Routing Table Commands
// =============================================================================

var thRTListCmd = &cobra.Command{
	Use: "list", Short: "List all routing tables",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.ListRoutingTables(context.Background())
		if err != nil {
			exitWithError("Failed to list routing tables", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var thRTGetCmd = &cobra.Command{
	Use: "get [id]", Short: "Get routing table details", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.GetRoutingTable(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get routing table", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		rt := result.RoutingTable
		fmt.Printf("ID:          %s\nName:        %s\nTransit Hub: %s\nDefault:     %v\nStatus:      %s\n", rt.ID, rt.Name, rt.TransitHubID, rt.DefaultTable, rt.Status)
	},
}

var thRTCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a new routing table",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		thID, _ := cmd.Flags().GetString("transit-hub-id")
		result, err := client.CreateRoutingTable(context.Background(), &transithub.CreateRoutingTableInput{Name: name, Description: desc, TransitHubID: thID})
		if err != nil {
			exitWithError("Failed to create routing table", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Routing table created: %s\n", result.RoutingTable.ID)
	},
}

var thRTUpdateCmd = &cobra.Command{
	Use: "update [id]", Short: "Update a routing table", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		result, err := client.UpdateRoutingTable(context.Background(), args[0], &transithub.UpdateRoutingTableInput{Name: name, Description: desc})
		if err != nil {
			exitWithError("Failed to update routing table", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Routing table updated: %s\n", result.RoutingTable.ID)
	},
}

var thRTDeleteCmd = &cobra.Command{
	Use: "delete [id]", Short: "Delete a routing table", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		if err := client.DeleteRoutingTable(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete routing table", err)
		}
		fmt.Printf("Routing table %s deleted\n", args[0])
	},
}

// =============================================================================
// Routing Association Commands
// =============================================================================

var thAssocListCmd = &cobra.Command{
	Use: "list", Short: "List all routing associations",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.ListRoutingAssociations(context.Background())
		if err != nil {
			exitWithError("Failed to list associations", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var thAssocGetCmd = &cobra.Command{
	Use: "get [id]", Short: "Get association details", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.GetRoutingAssociation(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get association", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		a := result.Association
		fmt.Printf("ID:             %s\nRouting Table:  %s\nAttachment:     %s\nStatus:         %s\n", a.ID, a.RoutingTableID, a.AttachmentID, a.Status)
	},
}

var thAssocCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a new routing association",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		rtID, _ := cmd.Flags().GetString("routing-table-id")
		attID, _ := cmd.Flags().GetString("attachment-id")
		result, err := client.CreateRoutingAssociation(context.Background(), &transithub.CreateRoutingAssociationInput{RoutingTableID: rtID, AttachmentID: attID})
		if err != nil {
			exitWithError("Failed to create association", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Association created: %s\n", result.Association.ID)
	},
}

var thAssocDeleteCmd = &cobra.Command{
	Use: "delete [id]", Short: "Delete an association", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		if err := client.DeleteRoutingAssociation(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete association", err)
		}
		fmt.Printf("Association %s deleted\n", args[0])
	},
}

// =============================================================================
// Routing Propagation Commands
// =============================================================================

var thPropListCmd = &cobra.Command{
	Use: "list", Short: "List all routing propagations",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.ListRoutingPropagations(context.Background())
		if err != nil {
			exitWithError("Failed to list propagations", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var thPropGetCmd = &cobra.Command{
	Use: "get [id]", Short: "Get propagation details", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.GetRoutingPropagation(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get propagation", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		p := result.Propagation
		fmt.Printf("ID:             %s\nRouting Table:  %s\nAttachment:     %s\nStatus:         %s\n", p.ID, p.RoutingTableID, p.AttachmentID, p.Status)
	},
}

var thPropCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a new routing propagation",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		rtID, _ := cmd.Flags().GetString("routing-table-id")
		attID, _ := cmd.Flags().GetString("attachment-id")
		result, err := client.CreateRoutingPropagation(context.Background(), &transithub.CreateRoutingPropagationInput{RoutingTableID: rtID, AttachmentID: attID})
		if err != nil {
			exitWithError("Failed to create propagation", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Propagation created: %s\n", result.Propagation.ID)
	},
}

var thPropDeleteCmd = &cobra.Command{
	Use: "delete [id]", Short: "Delete a propagation", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		if err := client.DeleteRoutingPropagation(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete propagation", err)
		}
		fmt.Printf("Propagation %s deleted\n", args[0])
	},
}

// =============================================================================
// Routing Rule Commands
// =============================================================================

var thRuleListCmd = &cobra.Command{
	Use: "list", Short: "List all routing rules",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.ListRoutingRules(context.Background())
		if err != nil {
			exitWithError("Failed to list rules", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tDESTINATION\tTARGET_TYPE\tTARGET_ID\tPROPAGATED")
		for _, r := range result.Rules {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n", r.ID, r.Destination, r.TargetType, r.TargetID, r.Propagated)
		}
		w.Flush()
	},
}

var thRuleGetCmd = &cobra.Command{
	Use: "get [id]", Short: "Get rule details", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.GetRoutingRule(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get rule", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		r := result.Rule
		fmt.Printf("ID:            %s\nRouting Table: %s\nDestination:   %s\nTarget Type:   %s\nTarget ID:     %s\nPropagated:    %v\n",
			r.ID, r.RoutingTableID, r.Destination, r.TargetType, r.TargetID, r.Propagated)
	},
}

var thRuleCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a new routing rule",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		rtID, _ := cmd.Flags().GetString("routing-table-id")
		dest, _ := cmd.Flags().GetString("destination")
		tType, _ := cmd.Flags().GetString("target-type")
		tID, _ := cmd.Flags().GetString("target-id")
		result, err := client.CreateRoutingRule(context.Background(), &transithub.CreateRoutingRuleInput{RoutingTableID: rtID, Destination: dest, TargetType: tType, TargetID: tID})
		if err != nil {
			exitWithError("Failed to create rule", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Rule created: %s\n", result.Rule.ID)
	},
}

var thRuleUpdateCmd = &cobra.Command{
	Use: "update [id]", Short: "Update a routing rule", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		dest, _ := cmd.Flags().GetString("destination")
		tType, _ := cmd.Flags().GetString("target-type")
		tID, _ := cmd.Flags().GetString("target-id")
		result, err := client.UpdateRoutingRule(context.Background(), args[0], &transithub.UpdateRoutingRuleInput{Destination: dest, TargetType: tType, TargetID: tID})
		if err != nil {
			exitWithError("Failed to update rule", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Rule updated: %s\n", result.Rule.ID)
	},
}

var thRuleDeleteCmd = &cobra.Command{
	Use: "delete [id]", Short: "Delete a routing rule", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		if err := client.DeleteRoutingRule(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete rule", err)
		}
		fmt.Printf("Rule %s deleted\n", args[0])
	},
}

// =============================================================================
// Multicast Domain Commands
// =============================================================================

var thMCListCmd = &cobra.Command{
	Use: "list", Short: "List all multicast domains",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.ListMulticastDomains(context.Background())
		if err != nil {
			exitWithError("Failed to list multicast domains", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTRANSIT_HUB_ID\tSTATUS")
		for _, m := range result.Domains {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", m.ID, m.Name, m.TransitHubID, m.Status)
		}
		w.Flush()
	},
}

var thMCGetCmd = &cobra.Command{
	Use: "get [id]", Short: "Get multicast domain details", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		result, err := client.GetMulticastDomain(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get multicast domain", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		m := result.Domain
		fmt.Printf("ID:          %s\nName:        %s\nTransit Hub: %s\nStatus:      %s\n", m.ID, m.Name, m.TransitHubID, m.Status)
	},
}

var thMCCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a new multicast domain",
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		thID, _ := cmd.Flags().GetString("transit-hub-id")
		result, err := client.CreateMulticastDomain(context.Background(), &transithub.CreateMulticastDomainInput{Name: name, Description: desc, TransitHubID: thID})
		if err != nil {
			exitWithError("Failed to create multicast domain", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Multicast domain created: %s\n", result.Domain.ID)
	},
}

var thMCUpdateCmd = &cobra.Command{
	Use: "update [id]", Short: "Update a multicast domain", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		result, err := client.UpdateMulticastDomain(context.Background(), args[0], &transithub.UpdateMulticastDomainInput{Name: name, Description: desc})
		if err != nil {
			exitWithError("Failed to update multicast domain", err)
		}
		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Multicast domain updated: %s\n", result.Domain.ID)
	},
}

var thMCDeleteCmd = &cobra.Command{
	Use: "delete [id]", Short: "Delete a multicast domain", Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newTransitHubClient()
		if err := client.DeleteMulticastDomain(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete multicast domain", err)
		}
		fmt.Printf("Multicast domain %s deleted\n", args[0])
	},
}
