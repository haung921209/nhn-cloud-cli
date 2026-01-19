package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLInstanceCmd)
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLInstanceNetworkCmd)
	rdsPostgreSQLCmd.AddCommand(describePostgreSQLInstanceStorageCmd)

	describePostgreSQLInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	describePostgreSQLInstanceCmd.MarkFlagRequired("db-instance-identifier")

	describePostgreSQLInstanceNetworkCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	describePostgreSQLInstanceNetworkCmd.MarkFlagRequired("db-instance-identifier")

	describePostgreSQLInstanceStorageCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	describePostgreSQLInstanceStorageCmd.MarkFlagRequired("db-instance-identifier")
}

var describePostgreSQLInstanceCmd = &cobra.Command{
	Use:   "describe-db-instance",
	Short: "Describe details of a specific PostgreSQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		identifier, _ := cmd.Flags().GetString("db-instance-identifier")
		instanceID, err := resolvePostgreSQLInstanceIdentifier(client, identifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		fmt.Printf("Fetching details for instance %s...\n", instanceID)

		ctx := context.Background()
		resp, err := client.GetInstance(ctx, instanceID)
		if err != nil {
			exitWithError("failed to get instance details", err)
		}

		postgresqlPrintJSON(resp)
	},
}

var describePostgreSQLInstanceNetworkCmd = &cobra.Command{
	Use:   "describe-db-instance-network",
	Short: "Describe network information of a PostgreSQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		identifier, _ := cmd.Flags().GetString("db-instance-identifier")
		instanceID, err := resolvePostgreSQLInstanceIdentifier(client, identifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		fmt.Printf("Fetching network info for instance %s...\n", instanceID)

		ctx := context.Background()
		resp, err := client.GetNetworkInfo(ctx, instanceID)
		if err != nil {
			exitWithError("failed to get network info", err)
		}

		postgresqlPrintJSON(resp)
	},
}

var describePostgreSQLInstanceStorageCmd = &cobra.Command{
	Use:   "describe-db-instance-storage",
	Short: "Describe storage information of a PostgreSQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		identifier, _ := cmd.Flags().GetString("db-instance-identifier")
		instanceID, err := resolvePostgreSQLInstanceIdentifier(client, identifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		fmt.Printf("Fetching storage info for instance %s...\n", instanceID)

		ctx := context.Background()
		resp, err := client.GetStorageInfo(ctx, instanceID)
		if err != nil {
			exitWithError("failed to get storage info", err)
		}

		postgresqlPrintJSON(resp)
	},
}
