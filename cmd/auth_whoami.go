package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	authCmd.AddCommand(authWhoamiCmd)
}

// tenantSource names where the effective tenant ID comes from.
func tenantSource() string {
	switch {
	case tenantID != "":
		return "--tenant-id flag"
	case os.Getenv("NHN_CLOUD_TENANT_ID") != "":
		return "NHN_CLOUD_TENANT_ID"
	case LoadConfig().TenantID != "":
		return "profile " + activeProfile()
	}
	return "not set"
}

func activeProfile() string {
	if profile != "" {
		return profile
	}
	return "default"
}

var authWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show which project (tenant) and region Compute/Network commands act on",
	Long: `Issue an IaaS token with the effective credentials and show the project
(tenant) it is scoped to. Use it before creating resources to confirm the
target project, especially when --tenant-id, NHN_CLOUD_TENANT_ID or
--profile override the default profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		api, err := getIaaSAPI(context.Background())
		if err != nil {
			exitWithError("Failed to authenticate", err)
		}
		info := map[string]string{
			"profile":       activeProfile(),
			"region":        getRegion(),
			"tenant_id":     api.tenantID,
			"tenant_source": tenantSource(),
			"project_name":  api.tenantName,
			"username":      getUsername(),
		}
		emit(info, func() {
			fmt.Printf("Profile:      %s\n", info["profile"])
			fmt.Printf("Region:       %s\n", info["region"])
			fmt.Printf("Tenant ID:    %s (from %s)\n", info["tenant_id"], info["tenant_source"])
			fmt.Printf("Project name: %s\n", info["project_name"])
			fmt.Printf("Username:     %s\n", info["username"])
		})
	},
}
