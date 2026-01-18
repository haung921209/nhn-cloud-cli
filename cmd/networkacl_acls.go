package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/networkacl"
	"github.com/spf13/cobra"
)

func init() {
	networkACLCmd.AddCommand(aclDescribeCmd)
	networkACLCmd.AddCommand(aclGetCmd)
	networkACLCmd.AddCommand(aclCreateCmd)
	networkACLCmd.AddCommand(aclUpdateCmd)
	networkACLCmd.AddCommand(aclDeleteCmd)

	aclGetCmd.Flags().String("acl-id", "", "ACL ID (required)")
	aclGetCmd.MarkFlagRequired("acl-id")

	aclCreateCmd.Flags().String("name", "", "ACL name (required)")
	aclCreateCmd.Flags().String("description", "", "ACL description")
	aclCreateCmd.MarkFlagRequired("name")

	aclUpdateCmd.Flags().String("acl-id", "", "ACL ID (required)")
	aclUpdateCmd.Flags().String("name", "", "ACL name")
	aclUpdateCmd.Flags().String("description", "", "ACL description")
	aclUpdateCmd.MarkFlagRequired("acl-id")

	aclDeleteCmd.Flags().String("acl-id", "", "ACL ID (required)")
	aclDeleteCmd.MarkFlagRequired("acl-id")
}

var aclDescribeCmd = &cobra.Command{
	Use:     "describe-network-acls",
	Aliases: []string{"list-network-acls", "list-acls"},
	Short:   "List all Network ACLs",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()

		result, err := client.ListACLs(ctx)
		if err != nil {
			exitWithError("Failed to list ACLs", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tSHARED")
		for _, acl := range result.ACLs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n",
				acl.ID, acl.Name, acl.Description, acl.Shared)
		}
		w.Flush()
	},
}

var aclGetCmd = &cobra.Command{
	Use:     "describe-network-acl",
	Aliases: []string{"get-network-acl", "get-acl"},
	Short:   "Get Network ACL details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("acl-id")

		result, err := client.GetACL(ctx, id)
		if err != nil {
			exitWithError("Failed to get ACL", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		acl := result.ACL
		fmt.Printf("ID:          %s\n", acl.ID)
		fmt.Printf("Name:        %s\n", acl.Name)
		fmt.Printf("Description: %s\n", acl.Description)
		fmt.Printf("Shared:      %v\n", acl.Shared)
		fmt.Printf("Tenant ID:   %s\n", acl.TenantID)
		fmt.Printf("Created:     %s\n", acl.CreateTime)
		fmt.Printf("Updated:     %s\n", acl.UpdateTime)
	},
}

var aclCreateCmd = &cobra.Command{
	Use:   "create-network-acl",
	Short: "Create a new Network ACL",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &networkacl.CreateACLInput{
			Name:        name,
			Description: description,
		}

		result, err := client.CreateACL(ctx, input)
		if err != nil {
			exitWithError("Failed to create ACL", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("ACL created: %s (%s)\n", result.ACL.Name, result.ACL.ID)
	},
}

var aclUpdateCmd = &cobra.Command{
	Use:   "update-network-acl",
	Short: "Update a Network ACL",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("acl-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &networkacl.UpdateACLInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateACL(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update ACL", err)
		}

		fmt.Printf("ACL updated: %s\n", result.ACL.ID)
	},
}

var aclDeleteCmd = &cobra.Command{
	Use:   "delete-network-acl",
	Short: "Delete a Network ACL",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("acl-id")

		if err := client.DeleteACL(ctx, id); err != nil {
			exitWithError("Failed to delete ACL", err)
		}

		fmt.Printf("ACL %s deleted\n", id)
	},
}
