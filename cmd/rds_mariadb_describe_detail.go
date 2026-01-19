package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rdsMariaDBCmd.AddCommand(describeMariaDBInstanceCmd)
	rdsMariaDBCmd.AddCommand(describeMariaDBInstanceNetworkCmd)
	rdsMariaDBCmd.AddCommand(describeMariaDBInstanceStorageCmd)

	describeMariaDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	describeMariaDBInstanceCmd.MarkFlagRequired("db-instance-identifier")

	describeMariaDBInstanceNetworkCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	describeMariaDBInstanceNetworkCmd.MarkFlagRequired("db-instance-identifier")

	describeMariaDBInstanceStorageCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	describeMariaDBInstanceStorageCmd.MarkFlagRequired("db-instance-identifier")
}

var describeMariaDBInstanceCmd = &cobra.Command{
	Use:   "describe-db-instance",
	Short: "Describe details of a specific MariaDB DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		identifier, _ := cmd.Flags().GetString("db-instance-identifier")
		instanceID, err := resolveMariaDBInstanceIdentifier(client, identifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		fmt.Printf("Fetching details for instance %s...\n", instanceID)

		ctx := context.Background()
		resp, err := client.GetInstance(ctx, instanceID)
		if err != nil {
			exitWithError("failed to get instance details", err)
		}

		printJSON(resp)
	},
}

var describeMariaDBInstanceNetworkCmd = &cobra.Command{
	Use:   "describe-db-instance-network",
	Short: "Describe network information of a MariaDB DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		identifier, _ := cmd.Flags().GetString("db-instance-identifier")
		instanceID, err := resolveMariaDBInstanceIdentifier(client, identifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		fmt.Printf("Fetching network info for instance %s...\n", instanceID)

		ctx := context.Background()
		resp, err := client.GetNetworkInfo(ctx, instanceID)
		if err != nil {
			exitWithError("failed to get network info", err)
		}

		printJSON(resp)
	},
}

var describeMariaDBInstanceStorageCmd = &cobra.Command{
	Use:   "describe-db-instance-storage",
	Short: "Describe storage information of a MariaDB DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		identifier, _ := cmd.Flags().GetString("db-instance-identifier")
		instanceID, err := resolveMariaDBInstanceIdentifier(client, identifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		fmt.Printf("Fetching storage info for instance %s...\n", instanceID)

		ctx := context.Background()
		resp, err := client.GetStorageInfo(ctx, instanceID)
		if err != nil {
			exitWithError("failed to get storage info", err)
		}

		printJSON(resp)
	},
}
