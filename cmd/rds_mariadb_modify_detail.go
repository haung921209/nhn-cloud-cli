package cmd

import (
	"context"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mariadb"
	"github.com/spf13/cobra"
)

func init() {
	rdsMariaDBCmd.AddCommand(modifyMariaDBInstanceNetworkCmd)
	rdsMariaDBCmd.AddCommand(modifyMariaDBInstanceStorageCmd)

	modifyMariaDBInstanceNetworkCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	modifyMariaDBInstanceNetworkCmd.Flags().Bool("enable-public-access", false, "Enable/Disable public access")
	modifyMariaDBInstanceNetworkCmd.MarkFlagRequired("db-instance-identifier")

	modifyMariaDBInstanceStorageCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	modifyMariaDBInstanceStorageCmd.Flags().Int("storage-size", 0, "New storage size in GB (Required)")
	modifyMariaDBInstanceStorageCmd.MarkFlagRequired("db-instance-identifier")
	modifyMariaDBInstanceStorageCmd.MarkFlagRequired("storage-size")
}

var modifyMariaDBInstanceNetworkCmd = &cobra.Command{
	Use:   "modify-db-instance-network",
	Short: "Modify network settings (Public Access) for a MariaDB DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		identifier, _ := cmd.Flags().GetString("db-instance-identifier")
		instanceID, err := resolveMariaDBInstanceIdentifier(client, identifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		usePublicAccess, _ := cmd.Flags().GetBool("enable-public-access")

		progressf("Modifying network info for instance %s (Public Access: %v)...\n", instanceID, usePublicAccess)

		req := &mariadb.ModifyNetworkInfoRequest{
			UsePublicAccess: usePublicAccess,
		}

		ctx := context.Background()
		resp, err := client.ModifyNetworkInfo(ctx, instanceID, req)
		if err != nil {
			exitWithError("failed to modify network info", err)
		}

		if err := printWriteResult(resp, "Job ID: %s\n", resp.JobID); err != nil {
			exitWithError("failed to print result", err)
		}
	},
}

var modifyMariaDBInstanceStorageCmd = &cobra.Command{
	Use:   "modify-db-instance-storage",
	Short: "Modify storage size for a MariaDB DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		identifier, _ := cmd.Flags().GetString("db-instance-identifier")
		instanceID, err := resolveMariaDBInstanceIdentifier(client, identifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		size, _ := cmd.Flags().GetInt("storage-size")

		progressf("Modifying storage size for instance %s to %d GB...\n", instanceID, size)

		req := &mariadb.ModifyStorageInfoRequest{
			StorageSize: size,
		}

		ctx := context.Background()
		resp, err := client.ModifyStorageInfo(ctx, instanceID, req)
		if err != nil {
			exitWithError("failed to modify storage info", err)
		}

		if err := printWriteResult(resp, "Job ID: %s\n", resp.JobID); err != nil {
			exitWithError("failed to print result", err)
		}
	},
}
