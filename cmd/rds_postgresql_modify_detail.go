package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
	"github.com/spf13/cobra"
)

func init() {
	rdsPostgreSQLCmd.AddCommand(modifyPostgreSQLInstanceNetworkCmd)
	rdsPostgreSQLCmd.AddCommand(modifyPostgreSQLInstanceStorageCmd)

	modifyPostgreSQLInstanceNetworkCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	modifyPostgreSQLInstanceNetworkCmd.Flags().Bool("enable-public-access", false, "Enable or disable public access (true/false)")
	modifyPostgreSQLInstanceNetworkCmd.MarkFlagRequired("db-instance-identifier")
	modifyPostgreSQLInstanceNetworkCmd.MarkFlagRequired("enable-public-access")

	modifyPostgreSQLInstanceStorageCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	modifyPostgreSQLInstanceStorageCmd.Flags().Int("storage-size", 0, "New storage size in GB (Required)")
	modifyPostgreSQLInstanceStorageCmd.MarkFlagRequired("db-instance-identifier")
	modifyPostgreSQLInstanceStorageCmd.MarkFlagRequired("storage-size")
}

var modifyPostgreSQLInstanceNetworkCmd = &cobra.Command{
	Use:   "modify-db-instance-network",
	Short: "Modify network configuration (Public Access) of a PostgreSQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		identifier, _ := cmd.Flags().GetString("db-instance-identifier")
		instanceID, err := resolvePostgreSQLInstanceIdentifier(client, identifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		usePublicAccess, _ := cmd.Flags().GetBool("enable-public-access")

		fmt.Printf("Modifying network info for instance %s (Public Access: %v)...\n", instanceID, usePublicAccess)

		ctx := context.Background()
		req := &postgresql.ModifyNetworkInfoRequest{
			UsePublicAccess: usePublicAccess,
		}

		resp, err := client.ModifyNetworkInfo(ctx, instanceID, req)
		if err != nil {
			exitWithError("failed to modify network info", err)
		}

		fmt.Printf("Job ID: %s\n", resp.JobID)
	},
}

var modifyPostgreSQLInstanceStorageCmd = &cobra.Command{
	Use:   "modify-db-instance-storage",
	Short: "Modify storage size of a PostgreSQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		identifier, _ := cmd.Flags().GetString("db-instance-identifier")
		instanceID, err := resolvePostgreSQLInstanceIdentifier(client, identifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		size, _ := cmd.Flags().GetInt("storage-size")

		fmt.Printf("Modifying storage size for instance %s to %d GB...\n", instanceID, size)

		ctx := context.Background()
		req := &postgresql.ModifyStorageInfoRequest{
			StorageSize: size,
		}

		resp, err := client.ModifyStorageInfo(ctx, instanceID, req)
		if err != nil {
			exitWithError("failed to modify storage info", err)
		}

		fmt.Printf("Job ID: %s\n", resp.JobID)
	},
}
