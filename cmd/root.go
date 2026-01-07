package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	region   string
	appKey   string
	debug    bool
	output   string
	username string
	password string
	tenantID string
)

var rootCmd = &cobra.Command{
	Use:   "nhncloud",
	Short: "NHN Cloud CLI - Command line interface for NHN Cloud services",
	Long: `NHN Cloud CLI provides a unified command line interface to manage
NHN Cloud services including RDS, Compute, Network, and more.

Environment Variables (OAuth - for RDS services):
  NHN_CLOUD_REGION       Default region (e.g., kr1, kr2, jp1)
  NHN_CLOUD_APPKEY       Application key
  NHN_CLOUD_ACCESS_KEY   Access key for authentication
  NHN_CLOUD_SECRET_KEY   Secret key for authentication

Environment Variables (Identity - for Compute/Network services):
  NHN_CLOUD_USERNAME     API username (email)
  NHN_CLOUD_PASSWORD     API password
  NHN_CLOUD_TENANT_ID    Project tenant ID`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", os.Getenv("NHN_CLOUD_REGION"), "NHN Cloud region (kr1, kr2, jp1)")
	rootCmd.PersistentFlags().StringVar(&appKey, "appkey", os.Getenv("NHN_CLOUD_APPKEY"), "Application key")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "Output format (table, json, yaml)")

	rootCmd.PersistentFlags().StringVar(&username, "username", os.Getenv("NHN_CLOUD_USERNAME"), "API username (for Compute/Network)")
	rootCmd.PersistentFlags().StringVar(&password, "password", os.Getenv("NHN_CLOUD_PASSWORD"), "API password (for Compute/Network)")
	rootCmd.PersistentFlags().StringVar(&tenantID, "tenant-id", os.Getenv("NHN_CLOUD_TENANT_ID"), "Tenant ID (for Compute/Network)")
}

func getRegion() string {
	if region != "" {
		return region
	}
	if r := os.Getenv("NHN_CLOUD_REGION"); r != "" {
		return r
	}
	return "kr1"
}

func getAppKey() string {
	if appKey != "" {
		return appKey
	}
	return os.Getenv("NHN_CLOUD_APPKEY")
}

func getUsername() string {
	if username != "" {
		return username
	}
	return os.Getenv("NHN_CLOUD_USERNAME")
}

func getPassword() string {
	if password != "" {
		return password
	}
	return os.Getenv("NHN_CLOUD_PASSWORD")
}

func getTenantID() string {
	if tenantID != "" {
		return tenantID
	}
	return os.Getenv("NHN_CLOUD_TENANT_ID")
}

func exitWithError(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s: %v\n", msg, err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	}
	os.Exit(1)
}
