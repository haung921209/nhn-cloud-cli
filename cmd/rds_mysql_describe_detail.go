package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rdsMySQLCmd.AddCommand(describeDBInstanceCmd)
	rdsMySQLCmd.AddCommand(describeDBInstanceNetworkCmd)
	rdsMySQLCmd.AddCommand(describeDBInstanceStorageCmd)

	describeDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	describeDBInstanceCmd.MarkFlagRequired("db-instance-identifier")

	describeDBInstanceNetworkCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	describeDBInstanceNetworkCmd.MarkFlagRequired("db-instance-identifier")

	describeDBInstanceStorageCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	describeDBInstanceStorageCmd.MarkFlagRequired("db-instance-identifier")
}

var describeDBInstanceCmd = &cobra.Command{
	Use:   "describe-db-instance",
	Short: "Describe details of a specific MySQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()

		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")

		fmt.Printf("Fetching details for instance %s...\n", instanceID)

		ctx := context.Background()
		resp, err := client.GetInstance(ctx, instanceID)
		if err != nil {
			exitWithError("failed to get instance details", err)
		}

		printJSON(resp)
	},
}

var describeDBInstanceNetworkCmd = &cobra.Command{
	Use:   "describe-db-instance-network",
	Short: "Describe network information of a MySQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()

		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")

		fmt.Printf("Fetching network info for instance %s...\n", instanceID)

		ctx := context.Background()
		resp, err := client.GetNetworkInfo(ctx, instanceID)
		if err != nil {
			exitWithError("failed to get network info", err)
		}

		printJSON(resp)
	},
}

var describeDBInstanceStorageCmd = &cobra.Command{
	Use:   "describe-db-instance-storage",
	Short: "Describe storage information of a MySQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()

		instanceID, _ := cmd.Flags().GetString("db-instance-identifier")

		fmt.Printf("Fetching storage info for instance %s...\n", instanceID)

		ctx := context.Background()
		resp, err := client.GetStorageInfo(ctx, instanceID)
		if err != nil {
			exitWithError("failed to get storage info", err)
		}

		printJSON(resp)
	},
}
