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
	networkACLCmd.AddCommand(aclDescribeBindingsCmd)
	networkACLCmd.AddCommand(aclGetBindingCmd)
	networkACLCmd.AddCommand(aclCreateBindingCmd)
	networkACLCmd.AddCommand(aclDeleteBindingCmd)

	aclDescribeBindingsCmd.Flags().String("acl-id", "", "Filter by ACL ID")

	aclGetBindingCmd.Flags().String("binding-id", "", "Binding ID (required)")
	aclGetBindingCmd.MarkFlagRequired("binding-id")

	aclCreateBindingCmd.Flags().String("acl-id", "", "ACL ID (required)")
	aclCreateBindingCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	aclCreateBindingCmd.MarkFlagRequired("acl-id")
	aclCreateBindingCmd.MarkFlagRequired("subnet-id")

	aclDeleteBindingCmd.Flags().String("binding-id", "", "Binding ID (required)")
	aclDeleteBindingCmd.MarkFlagRequired("binding-id")
}

var aclDescribeBindingsCmd = &cobra.Command{
	Use:     "describe-bindings",
	Aliases: []string{"list-bindings"},
	Short:   "List ACL bindings",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		aclID, _ := cmd.Flags().GetString("acl-id")

		var result *networkacl.ListACLBindingsOutput
		var err error

		if aclID != "" {
			result, err = client.ListACLBindingsByACL(ctx, aclID)
		} else {
			result, err = client.ListACLBindings(ctx)
		}

		if err != nil {
			exitWithError("Failed to list ACL bindings", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tACL_ID\tSUBNET_ID")
		for _, binding := range result.ACLBindings {
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				binding.ID, binding.ACLID, binding.SubnetID)
		}
		w.Flush()
	},
}

var aclGetBindingCmd = &cobra.Command{
	Use:     "describe-binding",
	Aliases: []string{"get-binding"},
	Short:   "Get ACL binding details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("binding-id")

		result, err := client.GetACLBinding(ctx, id)
		if err != nil {
			exitWithError("Failed to get ACL binding", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		binding := result.ACLBinding
		fmt.Printf("ID:        %s\n", binding.ID)
		fmt.Printf("ACL ID:    %s\n", binding.ACLID)
		fmt.Printf("Subnet ID: %s\n", binding.SubnetID)
		fmt.Printf("Tenant ID: %s\n", binding.TenantID)
		fmt.Printf("Created:   %s\n", binding.CreateTime)
	},
}

var aclCreateBindingCmd = &cobra.Command{
	Use:   "create-binding",
	Short: "Create a new ACL binding",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		aclID, _ := cmd.Flags().GetString("acl-id")
		subnetID, _ := cmd.Flags().GetString("subnet-id")

		input := &networkacl.CreateACLBindingInput{
			ACLID:    aclID,
			SubnetID: subnetID,
		}

		result, err := client.CreateACLBinding(ctx, input)
		if err != nil {
			exitWithError("Failed to create ACL binding", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("ACL binding created: %s\n", result.ACLBinding.ID)
	},
}

var aclDeleteBindingCmd = &cobra.Command{
	Use:   "delete-binding",
	Short: "Delete an ACL binding",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNetworkACLClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("binding-id")

		if err := client.DeleteACLBinding(ctx, id); err != nil {
			exitWithError("Failed to delete ACL binding", err)
		}
		fmt.Printf("ACL binding %s deleted\n", id)
	},
}
