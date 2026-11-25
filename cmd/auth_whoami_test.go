package cmd

import "testing"

func TestTenantSourceNamesTheEnvironmentVariable(t *testing.T) {
	t.Setenv("NHN_CLOUD_TENANT_ID", "from-env")
	old := tenantID
	defer func() { tenantID = old }()
	tenantID = "from-env" // what the flag default holds when only the variable is set
	if got := tenantSource(); got != "NHN_CLOUD_TENANT_ID" {
		t.Errorf("env only: %q", got)
	}
	f := rootCmd.PersistentFlags().Lookup("tenant-id")
	defer func() { f.Changed = false }()
	rootCmd.PersistentFlags().Set("tenant-id", "typed")
	if got := tenantSource(); got != "--tenant-id flag" {
		t.Errorf("typed flag: %q", got)
	}
}
