package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Version information - set at build time with ldflags
var (
	Version   = "v0.7.18"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// Version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("NHN Cloud CLI %s\n", Version)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Build Date: %s\n", BuildDate)
	},
}

// Configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure NHN Cloud CLI credentials",
	Long: `Interactive setup for NHN Cloud CLI credentials.

This will create/update the config file at ~/.nhncloud/credentials

Required for OAuth services (RDS, NCR, NCS, IAM):
  - Access Key ID
  - Secret Access Key

Required for Identity services (Compute, Network, Block Storage):
  - Username (email)
  - API Password
  - Tenant ID

Optional App Keys for specific services:
  - RDS MySQL App Key
  - RDS MariaDB App Key
  - RDS PostgreSQL App Key`,
	Run: func(cmd *cobra.Command, args []string) {
		runConfigure()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configureCmd)
}

func runConfigure() {
	reader := bufio.NewReader(os.Stdin)
	cfg := LoadConfig()

	fmt.Println("NHN Cloud CLI Configuration")
	fmt.Println("============================")
	fmt.Println()

	// OAuth credentials
	fmt.Println("OAuth Credentials (for RDS, NCR, NCS, IAM services):")
	cfg.AccessKeyID = promptWithDefault(reader, "Access Key ID", cfg.AccessKeyID)
	cfg.SecretAccessKey = promptWithDefault(reader, "Secret Access Key", maskSecret(cfg.SecretAccessKey))
	if !strings.HasPrefix(cfg.SecretAccessKey, "***") {
		// User entered a new value
	} else {
		// Keep existing value - reload
		cfg.SecretAccessKey = LoadConfig().SecretAccessKey
	}
	fmt.Println()

	// Identity credentials
	fmt.Println("Identity Credentials (for Compute, Network, Block Storage):")
	cfg.Username = promptWithDefault(reader, "Username (email)", cfg.Username)
	cfg.APIPassword = promptWithDefault(reader, "API Password", maskSecret(cfg.APIPassword))
	if strings.HasPrefix(cfg.APIPassword, "***") {
		cfg.APIPassword = LoadConfig().APIPassword
	}
	cfg.TenantID = promptWithDefault(reader, "Tenant ID", cfg.TenantID)
	fmt.Println()

	// Region
	fmt.Println("Default Region:")
	cfg.Region = promptWithDefault(reader, "Region (kr1, kr2, jp1)", cfg.Region)
	fmt.Println()

	// App Keys
	fmt.Println("App Keys (optional, for RDS services):")
	cfg.RDSAppKey = promptWithDefault(reader, "RDS MySQL App Key", cfg.RDSAppKey)
	cfg.RDSMariaDBAppKey = promptWithDefault(reader, "RDS MariaDB App Key", cfg.RDSMariaDBAppKey)
	cfg.RDSPostgreSQLAppKey = promptWithDefault(reader, "RDS PostgreSQL App Key", cfg.RDSPostgreSQLAppKey)
	fmt.Println()

	// Optional tenant IDs
	fmt.Println("Optional Tenant IDs (if different from main tenant):")
	cfg.NKSTenantID = promptWithDefault(reader, "NKS Tenant ID", cfg.NKSTenantID)
	cfg.OBSTenantID = promptWithDefault(reader, "Object Storage Tenant ID", cfg.OBSTenantID)
	fmt.Println()

	// Save config
	if err := saveConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration saved to ~/.nhncloud/credentials")
}

func promptWithDefault(reader *bufio.Reader, prompt, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultVal)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultVal
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

func maskSecret(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 4 {
		return "****"
	}
	return "***" + s[len(s)-4:]
}

func saveConfig(cfg *Config) error {
	configDir := filepath.Join(os.Getenv("HOME"), ".nhncloud")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "credentials")
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	lines := []string{
		"[default]",
		fmt.Sprintf("access_key_id = %s", cfg.AccessKeyID),
		fmt.Sprintf("secret_access_key = %s", cfg.SecretAccessKey),
		fmt.Sprintf("region = %s", cfg.Region),
		fmt.Sprintf("username = %s", cfg.Username),
		fmt.Sprintf("api_password = %s", cfg.APIPassword),
		fmt.Sprintf("tenant_id = %s", cfg.TenantID),
	}

	// Only add optional fields if they have values
	if cfg.NKSTenantID != "" {
		lines = append(lines, fmt.Sprintf("nks_tenant_id = %s", cfg.NKSTenantID))
	}
	if cfg.OBSTenantID != "" {
		lines = append(lines, fmt.Sprintf("obs_tenant_id = %s", cfg.OBSTenantID))
	}
	if cfg.RDSAppKey != "" {
		lines = append(lines, fmt.Sprintf("rds_app_key = %s", cfg.RDSAppKey))
	}
	if cfg.RDSMariaDBAppKey != "" {
		lines = append(lines, fmt.Sprintf("rds_mariadb_app_key = %s", cfg.RDSMariaDBAppKey))
	}
	if cfg.RDSPostgreSQLAppKey != "" {
		lines = append(lines, fmt.Sprintf("rds_postgresql_app_key = %s", cfg.RDSPostgreSQLAppKey))
	}

	for _, line := range lines {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}

	return nil
}
