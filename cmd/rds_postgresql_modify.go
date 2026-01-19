package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
	"github.com/spf13/cobra"
)

func init() {
	rdsPostgreSQLCmd.AddCommand(modifyPostgreSQLInstanceCmd)

	modifyPostgreSQLInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	modifyPostgreSQLInstanceCmd.Flags().StringSlice("db-security-group-ids", nil, "List of security group IDs to attach")
	modifyPostgreSQLInstanceCmd.Flags().String("db-instance-name", "", "New name for the DB instance")
	modifyPostgreSQLInstanceCmd.Flags().String("db-flavor-id", "", "New flavor ID")
	modifyPostgreSQLInstanceCmd.Flags().String("description", "", "New description")

	modifyPostgreSQLInstanceCmd.MarkFlagRequired("db-instance-identifier")
}

var modifyPostgreSQLInstanceCmd = &cobra.Command{
	Use:   "modify-db-instance",
	Short: "Modify a PostgreSQL DB instance (Name, Security Groups, Flavor, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		identifier, _ := cmd.Flags().GetString("db-instance-identifier")
		instanceID, err := resolvePostgreSQLInstanceIdentifier(client, identifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		fmt.Printf("Modifying instance %s...\n", instanceID)

		req := &postgresql.ModifyInstanceRequest{}

		if cmd.Flags().Changed("db-security-group-ids") {
			sgs, _ := cmd.Flags().GetStringSlice("db-security-group-ids")
			req.DBSecurityGroupIDs = sgs
		}

		if cmd.Flags().Changed("db-instance-name") {
			name, _ := cmd.Flags().GetString("db-instance-name")
			req.DBInstanceName = &name
		}

		if cmd.Flags().Changed("db-flavor-id") {
			flavor, _ := cmd.Flags().GetString("db-flavor-id")
			req.DBFlavorID = &flavor
		}

		if cmd.Flags().Changed("description") {
			desc, _ := cmd.Flags().GetString("description")
			req.Description = &desc
		}

		ctx := context.Background()
		resp, err := client.ModifyInstance(ctx, instanceID, req)
		if err != nil {
			exitWithError("failed to modify instance", err)
		}

		fmt.Printf("Modification initiated.\n")
		fmt.Printf("Job ID: %s\n", resp.JobID)
	},
}
