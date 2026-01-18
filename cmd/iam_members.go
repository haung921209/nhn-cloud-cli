package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/iam"
	"github.com/spf13/cobra"
)

func init() {
	iamCmd.AddCommand(iamDescribeMembersCmd)
	iamCmd.AddCommand(iamDescribeMemberCmd)
	iamCmd.AddCommand(iamInviteMemberCmd)
	iamCmd.AddCommand(iamRemoveMemberCmd)

	iamDescribeMembersCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamDescribeMembersCmd.MarkFlagRequired("org-id")

	iamDescribeMemberCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamDescribeMemberCmd.MarkFlagRequired("org-id")

	iamInviteMemberCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamInviteMemberCmd.Flags().String("email", "", "Member email (required)")
	iamInviteMemberCmd.Flags().StringSlice("roles", []string{}, "Roles to assign")
	iamInviteMemberCmd.MarkFlagRequired("org-id")
	iamInviteMemberCmd.MarkFlagRequired("email")

	iamRemoveMemberCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamRemoveMemberCmd.MarkFlagRequired("org-id")
}

var iamDescribeMembersCmd = &cobra.Command{
	Use:   "describe-members",
	Short: "List members in an organization",
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()
		orgID, _ := cmd.Flags().GetString("org-id")

		result, err := client.ListMembers(ctx, orgID)
		if err != nil {
			exitWithError("Failed to list members", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tEMAIL\tNAME\tSTATUS")
		for _, m := range result.Members {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				m.ID, m.Email, m.Name, m.Status)
		}
		w.Flush()
	},
}

var iamDescribeMemberCmd = &cobra.Command{
	Use:   "describe-member [member-id]",
	Short: "Get member details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()
		orgID, _ := cmd.Flags().GetString("org-id")

		result, err := client.GetMember(ctx, orgID, args[0])
		if err != nil {
			exitWithError("Failed to get member", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		m := result.Member
		fmt.Printf("ID:      %s\n", m.ID)
		fmt.Printf("Email:   %s\n", m.Email)
		fmt.Printf("Name:    %s\n", m.Name)
		fmt.Printf("Status:  %s\n", m.Status)
		fmt.Printf("Created: %s\n", m.CreatedAt)
		if len(m.Roles) > 0 {
			fmt.Printf("Roles:   %v\n", m.Roles)
		}
	},
}

var iamInviteMemberCmd = &cobra.Command{
	Use:   "invite-member",
	Short: "Invite a new member to the organization",
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()

		orgID, _ := cmd.Flags().GetString("org-id")
		email, _ := cmd.Flags().GetString("email")
		roles, _ := cmd.Flags().GetStringSlice("roles")

		input := &iam.InviteMemberInput{
			Email: email,
			Roles: roles,
		}

		result, err := client.InviteMember(ctx, orgID, input)
		if err != nil {
			exitWithError("Failed to invite member", err)
		}

		fmt.Printf("Member invited successfully!\n")
		fmt.Printf("Member ID: %s\n", result.MemberID)
	},
}

var iamRemoveMemberCmd = &cobra.Command{
	Use:   "remove-member [member-id]",
	Short: "Remove a member from the organization",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()
		orgID, _ := cmd.Flags().GetString("org-id")

		if err := client.RemoveMember(ctx, orgID, args[0]); err != nil {
			exitWithError("Failed to remove member", err)
		}

		fmt.Printf("Member %s removed successfully\n", args[0])
	},
}
