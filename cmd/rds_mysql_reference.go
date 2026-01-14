package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// ============================================================================
// Reference Data Commands (AWS-style)
// ============================================================================

var describeDBInstanceClassesCmd = &cobra.Command{
	Use:   "describe-db-instance-classes",
	Short: "Describe available DB instance classes (flavors)",
	Long:  `Lists available database instance classes (flavors) with their specifications.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListFlavors(context.Background())
		if err != nil {
			exitWithError("failed to list flavors", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "FLAVOR_ID\tNAME\tVCPUS\tRAM_MB")
			for _, flavor := range result.DBFlavors {
				fmt.Fprintf(w, "%s\t%s\t%d\t%d\n",
					flavor.DBFlavorID,
					flavor.DBFlavorName,
					flavor.Vcpus,
					flavor.Ram,
				)
			}
			w.Flush()
		}
	},
}

var describeDBEngineVersionsCmd = &cobra.Command{
	Use:   "describe-db-engine-versions",
	Short: "Describe available MySQL engine versions",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListVersions(context.Background())
		if err != nil {
			exitWithError("failed to list versions", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "VERSION_ID\tVERSION_NAME")
			for _, version := range result.DBVersions {
				fmt.Fprintf(w, "%s\t%s\n",
					version.DBVersion,
					version.DBVersionName,
				)
			}
			w.Flush()
		}
	},
}

var describeSubnetsCmd = &cobra.Command{
	Use:   "describe-subnets",
	Short: "Describe available subnets",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListSubnets(context.Background())
		if err != nil {
			exitWithError("failed to list subnets", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SUBNET_ID\tNAME\tCIDR")
			for _, subnet := range result.Subnets {
				fmt.Fprintf(w, "%s\t%s\t%s\n",
					subnet.SubnetID,
					subnet.SubnetName,
					subnet.SubnetCIDR,
				)
			}
			w.Flush()
		}
	},
}

var describeStorageTypesCmd = &cobra.Command{
	Use:   "describe-storage-types",
	Short: "Describe available storage types",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMySQLClient()
		result, err := client.ListStorageTypes(context.Background())
		if err != nil {
			exitWithError("failed to list storage types", err)
		}

		if output == "json" {
			printJSON(result)
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "STORAGE_TYPE")
			for _, st := range result.StorageTypes {
				fmt.Fprintf(w, "%s\n", st)
			}
			w.Flush()
		}
	},
}

// ============================================================================
// Initialization
// ============================================================================

func init() {
	rdsMySQLCmd.AddCommand(describeDBInstanceClassesCmd)
	rdsMySQLCmd.AddCommand(describeDBEngineVersionsCmd)
	rdsMySQLCmd.AddCommand(describeSubnetsCmd)
	rdsMySQLCmd.AddCommand(describeStorageTypesCmd)
}
