package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Version information - set at build time with ldflags
var (
	Version   = "v0.7.19"
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
	cfg.SecretAccessKey = promptSecret(reader, "Secret Access Key", cfg.SecretAccessKey)
	fmt.Println()

	// Identity credentials
	fmt.Println("Identity Credentials (for Compute, Network, Block Storage):")
	cfg.Username = promptWithDefault(reader, "Username (email)", cfg.Username)
	cfg.APIPassword = promptSecret(reader, "API Password", cfg.APIPassword)
	cfg.TenantID = promptWithDefault(reader, "Tenant ID", cfg.TenantID)
	fmt.Println()

	// Region
	fmt.Println("Default Region:")
	cfg.Region = promptWithDefault(reader, "Region (kr1, kr2, jp1)", cfg.Region)
	fmt.Println()

	// App Keys
	fmt.Println("App Keys (optional, for RDS services):")
	cfg.RDSAppKey = promptSecret(reader, "RDS MySQL App Key", cfg.RDSAppKey)
	cfg.RDSMariaDBAppKey = promptSecret(reader, "RDS MariaDB App Key", cfg.RDSMariaDBAppKey)
	cfg.RDSPostgreSQLAppKey = promptSecret(reader, "RDS PostgreSQL App Key", cfg.RDSPostgreSQLAppKey)
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

// promptSecret reads a credential without echoing it when stdin is a
// terminal. The current value is shown masked; an empty answer keeps it.
func promptSecret(reader *bufio.Reader, prompt, current string) string {
	if current != "" {
		fmt.Printf("%s [%s]: ", prompt, maskSecret(current))
	} else {
		fmt.Printf("%s: ", prompt)
	}
	var input string
	if term.IsTerminal(int(os.Stdin.Fd())) {
		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return current
		}
		input = string(b)
	} else {
		line, err := reader.ReadString('\n')
		if err != nil && line == "" {
			return current
		}
		input = line
	}
	if input = strings.TrimSpace(input); input == "" {
		return current
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

	target := "default"
	if profile != "" {
		target = profile
	} else if p := os.Getenv("NHN_CLOUD_PROFILE"); p != "" {
		target = p
	}
	updates := [][2]string{
		{"access_key_id", cfg.AccessKeyID},
		{"secret_access_key", cfg.SecretAccessKey},
		{"region", cfg.Region},
		{"username", cfg.Username},
		{"api_password", cfg.APIPassword},
		{"tenant_id", cfg.TenantID},
		{"nks_tenant_id", cfg.NKSTenantID},
		{"obs_tenant_id", cfg.OBSTenantID},
		{"rds_app_key", cfg.RDSAppKey},
		{"rds_mariadb_app_key", cfg.RDSMariaDBAppKey},
		{"rds_postgresql_app_key", cfg.RDSPostgreSQLAppKey},
	}

	old, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	out := mergeProfile(string(old), target, updates)

	// Write a new file and rename it over the old one, so a failed write
	// never leaves a truncated credentials file, and the result is 0600
	// even when the old file was not.
	tmp, err := os.CreateTemp(configDir, ".credentials-*")
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	defer os.Remove(tmp.Name())
	if err := tmp.Chmod(0600); err != nil {
		tmp.Close()
		return fmt.Errorf("failed to write config: %w", err)
	}
	if _, err := tmp.WriteString(out); err != nil {
		tmp.Close()
		return fmt.Errorf("failed to write config: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return os.Rename(tmp.Name(), configPath)
}

// mergeProfile sets keys in one [profile] section of an INI text and keeps
// everything else: other profiles, other keys, comments. Empty values leave
// an existing key alone and are not added.
func mergeProfile(text, profile string, updates [][2]string) string {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if text == "" {
		lines = nil
	}
	pending := map[string]string{}
	for _, u := range updates {
		if u[1] != "" {
			pending[u[0]] = u[1]
		}
	}
	var out []string
	in, found := false, false
	flush := func() {
		for _, u := range updates {
			if v, ok := pending[u[0]]; ok {
				out = append(out, u[0]+" = "+v)
				delete(pending, u[0])
			}
		}
	}
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
			if in {
				flush()
			}
			in = strings.TrimSpace(t[1:len(t)-1]) == profile
			found = found || in
			out = append(out, line)
			continue
		}
		if in {
			if k, _, ok := strings.Cut(t, "="); ok {
				if v, ok := pending[strings.TrimSpace(k)]; ok {
					out = append(out, strings.TrimSpace(k)+" = "+v)
					delete(pending, strings.TrimSpace(k))
					continue
				}
			}
		}
		out = append(out, line)
	}
	if in {
		flush()
	}
	if !found {
		if len(out) > 0 {
			out = append(out, "")
		}
		out = append(out, "["+profile+"]")
		flush()
	}
	return strings.Join(out, "\n") + "\n"
}
