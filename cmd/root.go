package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	region   string
	appKey   string
	debug    bool
	output   string
	query    string
	username string
	password string
	tenantID string
)

var rootCmd = &cobra.Command{
	Use:   "nhncloud",
	Short: "NHN Cloud CLI - Command line interface for NHN Cloud services",
	Long: `NHN Cloud CLI provides a unified command line interface to manage
NHN Cloud services including RDS, Compute, Network, and more.

Configuration Priority (highest to lowest):
  1. Command-line flags (--region, --appkey, etc.)
  2. Environment variables (NHN_CLOUD_REGION, etc.)
  3. Config file (~/.nhncloud/credentials)

Config File Format (~/.nhncloud/credentials):
  [default]
  access_key_id = your-access-key
  secret_access_key = your-secret-key
  region = kr1
  username = your-email
  api_password = your-password
  tenant_id = your-tenant-id
  rds_app_key = your-rds-appkey`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", os.Getenv("NHN_CLOUD_REGION"), "NHN Cloud region (kr1, kr2, jp1)")
	rootCmd.PersistentFlags().StringVar(&appKey, "appkey", os.Getenv("NHN_CLOUD_APPKEY"), "Application key")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "Output format (table, json, yaml)")
	rootCmd.PersistentFlags().StringVar(&query, "query", "", "JMESPath query to filter output")

	rootCmd.PersistentFlags().StringVar(&username, "username", os.Getenv("NHN_CLOUD_USERNAME"), "API username (for Compute/Network)")
	rootCmd.PersistentFlags().StringVar(&password, "password", os.Getenv("NHN_CLOUD_PASSWORD"), "API password (for Compute/Network)")
	rootCmd.PersistentFlags().StringVar(&tenantID, "tenant-id", os.Getenv("NHN_CLOUD_TENANT_ID"), "Tenant ID (for Compute/Network)")
}

func getRegion() string {
	cfg := LoadConfig()
	if region != "" {
		return region
	}
	if r := os.Getenv("NHN_CLOUD_REGION"); r != "" {
		return r
	}
	if cfg.Region != "" {
		return strings.ToLower(cfg.Region)
	}
	return "kr1"
}

func getAppKey() string {
	cfg := LoadConfig()
	if appKey != "" {
		return appKey
	}
	if k := os.Getenv("NHN_CLOUD_APPKEY"); k != "" {
		return k
	}
	if cfg.AppKey != "" {
		return cfg.AppKey
	}
	return cfg.RDSAppKey
}

func getMariaDBAppKey() string {
	cfg := LoadConfig()
	if appKey != "" {
		return appKey
	}
	if k := os.Getenv("NHN_CLOUD_MARIADB_APPKEY"); k != "" {
		return k
	}
	if cfg.RDSMariaDBAppKey != "" {
		return cfg.RDSMariaDBAppKey
	}
	return cfg.AppKey
}

func getPostgreSQLAppKey() string {
	cfg := LoadConfig()
	if appKey != "" {
		return appKey
	}
	if k := os.Getenv("NHN_CLOUD_POSTGRESQL_APPKEY"); k != "" {
		return k
	}
	if cfg.RDSPostgreSQLAppKey != "" {
		return cfg.RDSPostgreSQLAppKey
	}
	return cfg.AppKey
}

func getAccessKey() string {
	cfg := LoadConfig()
	if k := os.Getenv("NHN_CLOUD_ACCESS_KEY"); k != "" {
		return k
	}
	return cfg.AccessKeyID
}

func getSecretKey() string {
	cfg := LoadConfig()
	if k := os.Getenv("NHN_CLOUD_SECRET_KEY"); k != "" {
		return k
	}
	return cfg.SecretAccessKey
}

func getUsername() string {
	cfg := LoadConfig()
	if username != "" {
		return username
	}
	if u := os.Getenv("NHN_CLOUD_USERNAME"); u != "" {
		return u
	}
	return cfg.Username
}

func getPassword() string {
	cfg := LoadConfig()
	if password != "" {
		return password
	}
	if p := os.Getenv("NHN_CLOUD_PASSWORD"); p != "" {
		return p
	}
	return cfg.APIPassword
}

func getTenantID() string {
	cfg := LoadConfig()
	if tenantID != "" {
		return tenantID
	}
	if t := os.Getenv("NHN_CLOUD_TENANT_ID"); t != "" {
		return t
	}
	return cfg.TenantID
}

func exitWithError(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s: %v\n", msg, err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	}
	os.Exit(1)
}

func getNCRAppKey() string {
	cfg := LoadConfig()
	if appKey != "" {
		return appKey
	}
	if k := os.Getenv("NHN_CLOUD_NCR_APPKEY"); k != "" {
		return k
	}
	if cfg.NCRAppKey != "" {
		return cfg.NCRAppKey
	}
	return cfg.AppKey
}

func getNCSAppKey() string {
	cfg := LoadConfig()
	if appKey != "" {
		return appKey
	}
	if k := os.Getenv("NHN_CLOUD_NCS_APPKEY"); k != "" {
		return k
	}
	return cfg.AppKey
}

func getRDSAppKey() string {
	cfg := LoadConfig()
	if appKey != "" {
		return appKey
	}
	if k := os.Getenv("NHN_CLOUD_MYSQL_APPKEY"); k != "" {
		return k
	}
	if cfg.RDSAppKey != "" {
		return cfg.RDSAppKey
	}
	return cfg.AppKey
}
