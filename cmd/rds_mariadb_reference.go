package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mariadb"
	"github.com/spf13/cobra"
)

var describeMariaDBFlavorsCmd = &cobra.Command{
	Use:   "describe-db-flavors",
	Short: "Describe MariaDB DB flavors",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		result, err := client.ListFlavors(context.Background())
		if err != nil {
			exitWithError("failed to list flavors", err)
		}

		if output == "json" {
			mariadbPrintJSON(result)
		} else {
			mariadbPrintFlavorList(result)
		}
	},
}

var describeMariaDBEngineVersionsCmd = &cobra.Command{
	Use:   "describe-db-engine-versions",
	Short: "Describe MariaDB DB engine versions",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		result, err := client.ListVersions(context.Background())
		if err != nil {
			exitWithError("failed to list versions", err)
		}

		if output == "json" {
			mariadbPrintJSON(result)
		} else {
			mariadbPrintVersionList(result)
		}
	},
}

var describeMariaDBStorageTypesCmd = &cobra.Command{
	Use:   "describe-db-storage-types",
	Short: "Describe MariaDB storage types",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		result, err := client.ListStorageTypes(context.Background())
		if err != nil {
			exitWithError("failed to list storage types", err)
		}

		if output == "json" {
			mariadbPrintJSON(result)
		} else {
			mariadbPrintStorageTypeList(result)
		}
	},
}

func mariadbPrintFlavorList(result *mariadb.ListFlavorsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "FLAVOR_ID\tNAME\tVCPUS\tRAM(MB)")
	for _, f := range result.DBFlavors {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\n",
			f.DBFlavorID,
			f.DBFlavorName,
			f.Vcpus,
			f.Ram,
		)
	}
	w.Flush()
}

func mariadbPrintVersionList(result *mariadb.ListVersionsResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VERSION\tNAME")
	for _, v := range result.DBVersions {
		fmt.Fprintf(w, "%s\t%s\n",
			v.DBVersion,
			v.DBVersionName,
		)
	}
	w.Flush()
}

func mariadbPrintStorageTypeList(result *mariadb.ListStorageTypesResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tNAME")
	for _, st := range result.StorageTypes {
		fmt.Fprintf(w, "%s\t%s\n", st.StorageType, st.StorageTypeName)
	}
	w.Flush()
}

func init() {
	rdsMariaDBCmd.AddCommand(describeMariaDBFlavorsCmd)
	rdsMariaDBCmd.AddCommand(describeMariaDBEngineVersionsCmd)
	rdsMariaDBCmd.AddCommand(describeMariaDBStorageTypesCmd)
}
