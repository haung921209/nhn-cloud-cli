package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/postgresql"
	"github.com/spf13/cobra"
)

// ============================================================================
// PostgreSQL Extended Commands - Users, Security Groups, Parameter Groups, etc.
// ============================================================================

// DB User Commands
var pgUserCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "Manage database users",
}

var pgListDBUsersCmd = &cobra.Command{
	Use:   "list [instance-id]",
	Short: "List all database users",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListDBUsers(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list users", err)
		}
		printPGUsers(result)
	},
}

var pgCreateDBUserCmd = &cobra.Command{
	Use:   "create [instance-id]",
	Short: "Create a database user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userName, _ := cmd.Flags().GetString("name")
		password, _ := cmd.Flags().GetString("password")
		authorityType, _ := cmd.Flags().GetString("authority-type")

		if userName == "" || password == "" {
			exitWithError("--name and --password are required", nil)
		}

		input := &postgresql.CreateDBUserInput{
			DBUserName:    userName,
			DBPassword:    password,
			AuthorityType: authorityType,
		}

		client := newPostgreSQLClient()
		result, err := client.CreateDBUser(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create user", err)
		}
		fmt.Printf("User creation initiated. Job ID: %s\n", result.JobID)
	},
}

var pgUpdateDBUserCmd = &cobra.Command{
	Use:   "update [instance-id] [user-id]",
	Short: "Update a database user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")
		if password == "" {
			exitWithError("--password is required", nil)
		}

		input := &postgresql.UpdateDBUserInput{
			DBPassword: password,
		}

		client := newPostgreSQLClient()
		result, err := client.UpdateDBUser(context.Background(), args[0], args[1], input)
		if err != nil {
			exitWithError("failed to update user", err)
		}
		fmt.Printf("User update initiated. Job ID: %s\n", result.JobID)
	},
}

var pgDeleteDBUserCmd = &cobra.Command{
	Use:   "delete [instance-id] [user-id]",
	Short: "Delete a database user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.DeleteDBUser(context.Background(), args[0], args[1])
		if err != nil {
			exitWithError("failed to delete user", err)
		}
		fmt.Printf("User deletion initiated. Job ID: %s\n", result.JobID)
	},
}

// Security Group Commands
var pgSGCmd = &cobra.Command{
	Use:     "security-group",
	Aliases: []string{"sg"},
	Short:   "Manage DB security groups",
}

var pgListSGCmd = &cobra.Command{
	Use:   "list",
	Short: "List all security groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListSecurityGroups(context.Background())
		if err != nil {
			exitWithError("failed to list security groups", err)
		}
		printPGSecurityGroups(result)
	},
}

var pgGetSGCmd = &cobra.Command{
	Use:   "get [security-group-id]",
	Short: "Get security group details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.GetSecurityGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get security group", err)
		}
		printPGSecurityGroupDetail(result)
	},
}

var pgCreateSGCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a security group",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		if name == "" {
			exitWithError("--name is required", nil)
		}

		input := &postgresql.CreateSecurityGroupInput{
			DBSecurityGroupName: name,
			Description:         desc,
		}
		client := newPostgreSQLClient()
		result, err := client.CreateSecurityGroup(context.Background(), input)
		if err != nil {
			exitWithError("failed to create security group", err)
		}
		fmt.Printf("Security group created. Job ID: %s\n", result.JobID)
	},
}

var pgDeleteSGCmd = &cobra.Command{
	Use:   "delete [security-group-id]",
	Short: "Delete a security group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		_, err := client.DeleteSecurityGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete security group", err)
		}
		fmt.Println("Security group deleted.")
	},
}

var pgCreateSGRuleCmd = &cobra.Command{
	Use:   "rule-create [security-group-id]",
	Short: "Create a security group rule",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		direction, _ := cmd.Flags().GetString("direction")
		etherType, _ := cmd.Flags().GetString("ether-type")
		cidr, _ := cmd.Flags().GetString("cidr")
		desc, _ := cmd.Flags().GetString("description")

		if direction == "" || cidr == "" {
			exitWithError("--direction and --cidr are required", nil)
		}
		if etherType == "" {
			etherType = "IPV4"
		}

		input := &postgresql.CreateSecurityGroupRuleInput{
			Direction:   direction,
			EtherType:   etherType,
			CIDR:        cidr,
			Description: desc,
		}
		client := newPostgreSQLClient()
		result, err := client.CreateSecurityGroupRule(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to create rule", err)
		}
		fmt.Printf("Security group rule created. Job ID: %s\n", result.JobID)
	},
}

var pgDeleteSGRuleCmd = &cobra.Command{
	Use:   "rule-delete [security-group-id] [rule-id]",
	Short: "Delete a security group rule",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		_, err := client.DeleteSecurityGroupRule(context.Background(), args[0], args[1])
		if err != nil {
			exitWithError("failed to delete rule", err)
		}
		fmt.Println("Security group rule deleted.")
	},
}

// Parameter Group Commands
var pgParamGroupCmd = &cobra.Command{
	Use:     "parameter-group",
	Aliases: []string{"pg", "param-group"},
	Short:   "Manage parameter groups",
}

var pgListParamGroupCmd = &cobra.Command{
	Use:   "list",
	Short: "List all parameter groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListParameterGroups(context.Background())
		if err != nil {
			exitWithError("failed to list parameter groups", err)
		}
		printPGParameterGroups(result)
	},
}

var pgGetParamGroupCmd = &cobra.Command{
	Use:   "get [parameter-group-id]",
	Short: "Get parameter group details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.GetParameterGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get parameter group", err)
		}
		printPGParameterGroupDetail(result)
	},
}

var pgCreateParamGroupCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a parameter group",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		version, _ := cmd.Flags().GetString("version")
		if name == "" || version == "" {
			exitWithError("--name and --version are required", nil)
		}

		input := &postgresql.CreateParameterGroupInput{
			ParameterGroupName: name,
			Description:        desc,
			DBVersion:          version,
		}
		client := newPostgreSQLClient()
		result, err := client.CreateParameterGroup(context.Background(), input)
		if err != nil {
			exitWithError("failed to create parameter group", err)
		}
		fmt.Printf("Parameter group created. Job ID: %s\n", result.JobID)
	},
}

var pgDeleteParamGroupCmd = &cobra.Command{
	Use:   "delete [parameter-group-id]",
	Short: "Delete a parameter group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		_, err := client.DeleteParameterGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete parameter group", err)
		}
		fmt.Println("Parameter group deleted.")
	},
}

// Notification Group Commands
var pgNotifGroupCmd = &cobra.Command{
	Use:     "notification-group",
	Aliases: []string{"ng"},
	Short:   "Manage notification groups",
}

var pgListNotifGroupCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notification groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListNotificationGroups(context.Background())
		if err != nil {
			exitWithError("failed to list notification groups", err)
		}
		printPGNotificationGroups(result)
	},
}

var pgGetNotifGroupCmd = &cobra.Command{
	Use:   "get [notification-group-id]",
	Short: "Get notification group details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.GetNotificationGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get notification group", err)
		}
		printPGNotificationGroupDetail(result)
	},
}

var pgDeleteNotifGroupCmd = &cobra.Command{
	Use:   "delete [notification-group-id]",
	Short: "Delete a notification group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		_, err := client.DeleteNotificationGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to delete notification group", err)
		}
		fmt.Println("Notification group deleted.")
	},
}

// Log Commands
var pgLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Manage logs",
}

var pgListLogsCmd = &cobra.Command{
	Use:   "list [instance-id]",
	Short: "List log files",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListLogFiles(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to list log files", err)
		}
		printPGLogFiles(result)
	},
}

// Resource Commands
var pgStorageTypesCmd = &cobra.Command{
	Use:   "storage-types",
	Short: "List available storage types",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListStorageTypes(context.Background())
		if err != nil {
			exitWithError("failed to list storage types", err)
		}
		printPGStorageTypes(result)
	},
}

var pgSubnetsCmd = &cobra.Command{
	Use:   "subnets",
	Short: "List available subnets",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.ListSubnets(context.Background())
		if err != nil {
			exitWithError("failed to list subnets", err)
		}
		printPGSubnets(result)
	},
}

// Network Commands
var pgNetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage network settings",
}

var pgGetNetworkCmd = &cobra.Command{
	Use:   "get [instance-id]",
	Short: "Get network information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		result, err := client.GetNetworkInfo(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get network info", err)
		}
		printPGNetworkInfo(result)
	},
}

// Storage & Protection Commands
var pgResizeStorageCmd = &cobra.Command{
	Use:   "resize-storage [instance-id]",
	Short: "Resize storage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		size, _ := cmd.Flags().GetInt("size")
		if size <= 0 {
			exitWithError("--size is required and must be positive", nil)
		}

		input := &postgresql.ModifyStorageInfoInput{StorageSize: size}
		client := newPostgreSQLClient()
		result, err := client.ModifyStorageInfo(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to resize storage", err)
		}
		fmt.Printf("Storage resize initiated. Job ID: %s\n", result.JobID)
	},
}

var pgDeletionProtectionCmd = &cobra.Command{
	Use:   "deletion-protection [instance-id]",
	Short: "Enable or disable deletion protection",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		enable, _ := cmd.Flags().GetBool("enable")
		disable, _ := cmd.Flags().GetBool("disable")
		if enable == disable {
			exitWithError("specify either --enable or --disable", nil)
		}

		input := &postgresql.ModifyDeletionProtectionInput{UseDeletionProtection: enable}
		client := newPostgreSQLClient()
		result, err := client.ModifyDeletionProtection(context.Background(), args[0], input)
		if err != nil {
			exitWithError("failed to modify deletion protection", err)
		}
		fmt.Printf("Deletion protection modified. Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Init
// ============================================================================

func init() {
	// User commands
	rdsPostgreSQLCmd.AddCommand(pgUserCmd)
	pgUserCmd.AddCommand(pgListDBUsersCmd)
	pgUserCmd.AddCommand(pgCreateDBUserCmd)
	pgUserCmd.AddCommand(pgUpdateDBUserCmd)
	pgUserCmd.AddCommand(pgDeleteDBUserCmd)
	pgCreateDBUserCmd.Flags().String("name", "", "User name (required)")
	pgCreateDBUserCmd.Flags().String("password", "", "Password (required)")
	pgCreateDBUserCmd.Flags().String("authority-type", "CRUD", "Authority type: READ, CRUD, or DDL")
	pgUpdateDBUserCmd.Flags().String("password", "", "New password (required)")

	// Security Group commands
	rdsPostgreSQLCmd.AddCommand(pgSGCmd)
	pgSGCmd.AddCommand(pgListSGCmd)
	pgSGCmd.AddCommand(pgGetSGCmd)
	pgSGCmd.AddCommand(pgCreateSGCmd)
	pgSGCmd.AddCommand(pgDeleteSGCmd)
	pgSGCmd.AddCommand(pgCreateSGRuleCmd)
	pgSGCmd.AddCommand(pgDeleteSGRuleCmd)
	pgCreateSGCmd.Flags().String("name", "", "Name (required)")
	pgCreateSGCmd.Flags().String("description", "", "Description")
	pgCreateSGRuleCmd.Flags().String("direction", "", "Direction: INGRESS or EGRESS (required)")
	pgCreateSGRuleCmd.Flags().String("ether-type", "IPV4", "Ether type: IPV4 or IPV6")
	pgCreateSGRuleCmd.Flags().String("cidr", "", "CIDR (required)")
	pgCreateSGRuleCmd.Flags().String("description", "", "Rule description")

	// Parameter Group commands
	rdsPostgreSQLCmd.AddCommand(pgParamGroupCmd)
	pgParamGroupCmd.AddCommand(pgListParamGroupCmd)
	pgParamGroupCmd.AddCommand(pgGetParamGroupCmd)
	pgParamGroupCmd.AddCommand(pgCreateParamGroupCmd)
	pgParamGroupCmd.AddCommand(pgDeleteParamGroupCmd)
	pgCreateParamGroupCmd.Flags().String("name", "", "Name (required)")
	pgCreateParamGroupCmd.Flags().String("description", "", "Description")
	pgCreateParamGroupCmd.Flags().String("version", "", "DB version (required)")

	// Notification Group commands
	rdsPostgreSQLCmd.AddCommand(pgNotifGroupCmd)
	pgNotifGroupCmd.AddCommand(pgListNotifGroupCmd)
	pgNotifGroupCmd.AddCommand(pgGetNotifGroupCmd)
	pgNotifGroupCmd.AddCommand(pgDeleteNotifGroupCmd)

	// Log commands
	rdsPostgreSQLCmd.AddCommand(pgLogCmd)
	pgLogCmd.AddCommand(pgListLogsCmd)

	// Resource commands
	rdsPostgreSQLCmd.AddCommand(pgStorageTypesCmd)
	rdsPostgreSQLCmd.AddCommand(pgSubnetsCmd)

	// Network commands
	rdsPostgreSQLCmd.AddCommand(pgNetworkCmd)
	pgNetworkCmd.AddCommand(pgGetNetworkCmd)

	// Storage & Protection commands
	rdsPostgreSQLCmd.AddCommand(pgResizeStorageCmd)
	pgResizeStorageCmd.Flags().Int("size", 0, "New storage size in GB (required)")
	rdsPostgreSQLCmd.AddCommand(pgDeletionProtectionCmd)
	pgDeletionProtectionCmd.Flags().Bool("enable", false, "Enable deletion protection")
	pgDeletionProtectionCmd.Flags().Bool("disable", false, "Disable deletion protection")
}

// ============================================================================
// Print Functions
// ============================================================================

func printPGUsers(result *postgresql.ListDBUsersOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tAUTHORITY\tSTATUS\tCREATED")
	for _, u := range result.DBUsers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", u.DBUserID, u.DBUserName, u.AuthorityType, u.DBUserStatus, u.CreatedYmdt.Format("2006-01-02 15:04:05"))
	}
	w.Flush()
}

func printPGSecurityGroups(result *postgresql.ListSecurityGroupsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tRULES")
	for _, sg := range result.DBSecurityGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", sg.DBSecurityGroupID, sg.DBSecurityGroupName, sg.Description, len(sg.Rules))
	}
	w.Flush()
}

func printPGSecurityGroupDetail(result *postgresql.SecurityGroupOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	fmt.Printf("ID:          %s\n", result.DBSecurityGroup.DBSecurityGroupID)
	fmt.Printf("Name:        %s\n", result.DBSecurityGroup.DBSecurityGroupName)
	fmt.Printf("Description: %s\n", result.DBSecurityGroup.Description)
	fmt.Println("\nRules:")
	for _, r := range result.DBSecurityGroup.Rules {
		fmt.Printf("  - %s: %s %s\n", r.RuleID, r.Direction, r.CIDR)
	}
}

func printPGParameterGroups(result *postgresql.ListParameterGroupsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tVERSION\tDEFAULT")
	for _, pg := range result.ParameterGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", pg.ParameterGroupID, pg.ParameterGroupName, pg.DBVersion, pg.IsDefault)
	}
	w.Flush()
}

func printPGParameterGroupDetail(result *postgresql.ParameterGroupOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	fmt.Printf("ID:          %s\n", result.ParameterGroupID)
	fmt.Printf("Name:        %s\n", result.ParameterGroupName)
	fmt.Printf("Version:     %s\n", result.DBVersion)
	fmt.Printf("Default:     %v\n", result.IsDefault)
	fmt.Printf("Parameters:  %d\n", len(result.Parameters))
}

func printPGNotificationGroups(result *postgresql.ListNotificationGroupsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tTYPE\tENABLED")
	for _, ng := range result.NotificationGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", ng.NotificationGroupID, ng.NotificationGroupName, ng.NotificationType, ng.IsEnabled)
	}
	w.Flush()
}

func printPGNotificationGroupDetail(result *postgresql.NotificationGroupOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	fmt.Printf("ID:      %s\n", result.NotificationGroupID)
	fmt.Printf("Name:    %s\n", result.NotificationGroupName)
	fmt.Printf("Type:    %s\n", result.NotificationType)
	fmt.Printf("Enabled: %v\n", result.IsEnabled)
}

func printPGLogFiles(result *postgresql.ListLogFilesOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "FILE NAME\tSIZE\tCREATED")
	for _, l := range result.LogFiles {
		fmt.Fprintf(w, "%s\t%d\t%s\n", l.LogFileName, l.LogFileSize, l.CreatedYmdt.Format("2006-01-02 15:04:05"))
	}
	w.Flush()
}

func printPGStorageTypes(result *postgresql.ListStorageTypesOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE")
	for _, t := range result.StorageTypes {
		fmt.Fprintf(w, "%s\n", t)
	}
	w.Flush()
}

func printPGSubnets(result *postgresql.ListSubnetsOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCIDR")
	for _, s := range result.Subnets {
		fmt.Fprintf(w, "%s\t%s\t%s\n", s.SubnetID, s.SubnetName, s.SubnetCIDR)
	}
	w.Flush()
}

func printPGNetworkInfo(result *postgresql.NetworkInfoOutput) {
	if output == "json" {
		printJSON(result)
		return
	}
	fmt.Println("Endpoints:")
	for _, ep := range result.EndPoints {
		publicStr := ""
		if ep.IsPublicAccess {
			publicStr = " (public)"
		}
		fmt.Printf("  - %s: %s%s\n", ep.Domain, ep.IPAddress, publicStr)
	}
}
