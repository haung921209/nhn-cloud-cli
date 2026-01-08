package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/haung921209/nhn-cloud-cli/internal/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication management commands",
	Long:  `Manage authentication tokens and credentials for NHN Cloud CLI.`,
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Display current authentication status and credential information.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := LoadConfig()

		status := map[string]interface{}{
			"region":               getRegion(),
			"oauth_credentials":    "not configured",
			"identity_credentials": "not configured",
			"token":                "not available",
		}

		if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
			status["oauth_credentials"] = "configured"
			if debug {
				status["access_key_id"] = cfg.AccessKeyID
			} else if len(cfg.AccessKeyID) > 8 {
				status["access_key_id"] = cfg.AccessKeyID[:8] + "***"
			}

			mgr := auth.NewTokenManager(getRegion(), cfg.AccessKeyID, cfg.SecretAccessKey)
			if token, err := mgr.GetToken(); err == nil && token.IsValid() {
				status["token"] = "valid"
				status["token_expires"] = token.ExpiresAt().Format("2006-01-02 15:04:05")
			} else {
				status["token"] = "invalid or expired"
			}
		}

		if cfg.Username != "" && cfg.APIPassword != "" && cfg.TenantID != "" {
			status["identity_credentials"] = "configured"
			if debug {
				status["username"] = cfg.Username
				status["tenant_id"] = cfg.TenantID
			} else {
				status["username"] = cfg.Username
				if len(cfg.TenantID) > 8 {
					status["tenant_id"] = cfg.TenantID[:8] + "***"
				}
			}
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(status)
			return
		}

		fmt.Println("Authentication Status")
		fmt.Println("=====================")
		fmt.Printf("Region: %s\n", status["region"])
		fmt.Println()
		fmt.Println("OAuth Credentials (RDS, IAM):")
		fmt.Printf("  Status: %s\n", status["oauth_credentials"])
		if status["access_key_id"] != nil {
			fmt.Printf("  Access Key ID: %s\n", status["access_key_id"])
		}
		fmt.Printf("  Token: %s\n", status["token"])
		if status["token_expires"] != nil {
			fmt.Printf("  Expires: %s\n", status["token_expires"])
		}
		fmt.Println()
		fmt.Println("Identity Credentials (Compute, Network):")
		fmt.Printf("  Status: %s\n", status["identity_credentials"])
		if status["username"] != nil {
			fmt.Printf("  Username: %s\n", status["username"])
		}
		if status["tenant_id"] != nil {
			fmt.Printf("  Tenant ID: %s\n", status["tenant_id"])
		}
	},
}

var authTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Token management commands",
	Long:  `Manage OAuth 2.0 access tokens for NHN Cloud API authentication.`,
}

var authTokenRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh access token",
	Long:  `Force refresh of the current OAuth access token.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := LoadConfig()

		if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
			exitWithError("OAuth credentials not configured. Run 'nhncloud configure' first", nil)
		}

		fmt.Println("Refreshing OAuth token...")

		mgr := auth.NewTokenManager(getRegion(), cfg.AccessKeyID, cfg.SecretAccessKey)
		token, err := mgr.RefreshToken()
		if err != nil {
			exitWithError("Failed to refresh token", err)
		}

		if output == "json" {
			result := map[string]interface{}{
				"token_type":  token.TokenType,
				"expires_in":  token.ExpiresIn,
				"issued_at":   token.IssuedAt.Format("2006-01-02 15:04:05"),
				"valid_until": token.ExpiresAt().Format("2006-01-02 15:04:05"),
			}
			if debug {
				result["access_token"] = token.AccessToken
			} else {
				result["access_token"] = "***REDACTED***"
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
			return
		}

		fmt.Println("Token refreshed successfully!")
		fmt.Printf("  Type: %s\n", token.TokenType)
		fmt.Printf("  Expires in: %d seconds\n", token.ExpiresIn)
		fmt.Printf("  Valid until: %s\n", token.ExpiresAt().Format("2006-01-02 15:04:05"))
	},
}

var authTokenRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke access token",
	Long:  `Revoke the current access token and clear token cache.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := LoadConfig()

		if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
			exitWithError("OAuth credentials not configured", nil)
		}

		mgr := auth.NewTokenManager(getRegion(), cfg.AccessKeyID, cfg.SecretAccessKey)
		if err := mgr.ClearToken(); err != nil {
			exitWithError("Failed to revoke token", err)
		}

		fmt.Println("Token revoked and cache cleared.")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authTokenCmd)

	authTokenCmd.AddCommand(authTokenRefreshCmd)
	authTokenCmd.AddCommand(authTokenRevokeCmd)
}
