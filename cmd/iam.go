package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/iam"
	"github.com/spf13/cobra"
)

var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Manage IAM organizations, projects, and members",
	Long:  `Manage Identity and Access Management resources including organizations, projects, and members.`,
}

func init() {
	rootCmd.AddCommand(iamCmd)

	iamCmd.AddCommand(iamOrgsCmd)
	iamCmd.AddCommand(iamOrgGetCmd)
	iamCmd.AddCommand(iamProjectsCmd)
	iamCmd.AddCommand(iamProjectGetCmd)
	iamCmd.AddCommand(iamMembersCmd)
	iamCmd.AddCommand(iamMemberGetCmd)
	iamCmd.AddCommand(iamMemberInviteCmd)
	iamCmd.AddCommand(iamMemberRemoveCmd)

	iamProjectsCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamProjectsCmd.MarkFlagRequired("org-id")

	iamProjectGetCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamProjectGetCmd.MarkFlagRequired("org-id")

	iamMembersCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamMembersCmd.MarkFlagRequired("org-id")

	iamMemberGetCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamMemberGetCmd.MarkFlagRequired("org-id")

	iamMemberInviteCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamMemberInviteCmd.Flags().String("email", "", "Member email (required)")
	iamMemberInviteCmd.Flags().StringSlice("roles", []string{}, "Roles to assign")
	iamMemberInviteCmd.MarkFlagRequired("org-id")
	iamMemberInviteCmd.MarkFlagRequired("email")

	iamMemberRemoveCmd.Flags().String("org-id", "", "Organization ID (required)")
	iamMemberRemoveCmd.MarkFlagRequired("org-id")
}

func getIAMClient() *iam.Client {
	creds := credentials.NewStatic(
		os.Getenv("NHN_CLOUD_ACCESS_KEY"),
		os.Getenv("NHN_CLOUD_SECRET_KEY"),
	)
	return iam.NewClient(getRegion(), creds, nil, debug)
}

var iamOrgsCmd = &cobra.Command{
	Use:   "organizations",
	Short: "List all organizations",
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()

		result, err := client.ListOrganizations(ctx)
		if err != nil {
			exitWithError("Failed to list organizations", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tCREATED")
		for _, o := range result.Organizations {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				o.ID, o.Name, o.Status, o.CreatedAt)
		}
		w.Flush()
	},
}

var iamOrgGetCmd = &cobra.Command{
	Use:   "organization-get [org-id]",
	Short: "Get organization details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()

		result, err := client.GetOrganization(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get organization", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		o := result.Organization
		fmt.Printf("ID:          %s\n", o.ID)
		fmt.Printf("Name:        %s\n", o.Name)
		fmt.Printf("Status:      %s\n", o.Status)
		fmt.Printf("Description: %s\n", o.Description)
		fmt.Printf("Created:     %s\n", o.CreatedAt)
	},
}

var iamProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List projects in an organization",
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()

		orgID, _ := cmd.Flags().GetString("org-id")

		result, err := client.ListProjects(ctx, orgID)
		if err != nil {
			exitWithError("Failed to list projects", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tCREATED")
		for _, p := range result.Projects {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				p.ID, p.Name, p.Status, p.CreatedAt)
		}
		w.Flush()
	},
}

var iamProjectGetCmd = &cobra.Command{
	Use:   "project-get [project-id]",
	Short: "Get project details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getIAMClient()
		ctx := context.Background()

		orgID, _ := cmd.Flags().GetString("org-id")

		result, err := client.GetProject(ctx, orgID, args[0])
		if err != nil {
			exitWithError("Failed to get project", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		p := result.Project
		fmt.Printf("ID:          %s\n", p.ID)
		fmt.Printf("Name:        %s\n", p.Name)
		fmt.Printf("Status:      %s\n", p.Status)
		fmt.Printf("Description: %s\n", p.Description)
		fmt.Printf("Org ID:      %s\n", p.OrganizationID)
		fmt.Printf("Created:     %s\n", p.CreatedAt)
	},
}

var iamMembersCmd = &cobra.Command{
	Use:   "members",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var iamMemberGetCmd = &cobra.Command{
	Use:   "member-get [member-id]",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var iamMemberInviteCmd = &cobra.Command{
	Use:   "member-invite",
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

var iamMemberRemoveCmd = &cobra.Command{
	Use:   "member-remove [member-id]",
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
