package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// Servers echo credentials back in error messages; the CLI must not print them.
func TestErrorTextHidesCredentials(t *testing.T) {
	old := loadedConfig
	defer func() { loadedConfig = old }()
	loadedConfig = &Config{AppKey: "PROFILEAPPKEY1", Extra: map[string]string{"dnsplus_app_key": "DNSAPPKEY99", "region": "kr1"}}
	t.Setenv("NHN_CLOUD_RESOURCE_WATCHER_APPKEY", "RWAPPKEY77")
	t.Setenv("NHN_MYSQL_APP_KEY", "MYSQLAPPKEY55")
	t.Setenv("NHN_SECRET_ACCESS_KEY", "RDSSECRET33")
	for _, leak := range []string{"PROFILEAPPKEY1", "DNSAPPKEY99", "RWAPPKEY77", "MYSQLAPPKEY55", "RDSSECRET33"} {
		got := errorText("list events: Invalid appKey: " + leak)
		if strings.Contains(got, leak) {
			t.Errorf("%s leaked: %s", leak, got)
		}
	}
	if got := errorText("region kr1 failed"); got != "region kr1 failed" {
		t.Errorf("non-secret text changed: %s", got)
	}
}

func TestGroupCommandRejectsUnknownSubcommand(t *testing.T) {
	root := &cobra.Command{Use: "nhncloud"}
	group := &cobra.Command{Use: "svc"}
	group.AddCommand(&cobra.Command{Use: "list", Run: func(*cobra.Command, []string) {}})
	root.AddCommand(group)
	rejectUnknownSubcommands(root)
	root.SetArgs([]string{"svc", "nosuch"})
	root.SilenceUsage, root.SilenceErrors = true, true
	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), `unknown command "nosuch"`) {
		t.Fatalf("want unknown command error, got %v", err)
	}
}

// A command's own secret flag (here --db-password) is hidden too.
func TestErrorTextHidesSecretFlagValues(t *testing.T) {
	old := loadedConfig
	defer func() { loadedConfig = old }()
	loadedConfig = &Config{Extra: map[string]string{}}
	sub := &cobra.Command{Use: "create-db-user-test", Run: func(*cobra.Command, []string) {}}
	sub.Flags().String("db-password", "", "")
	sub.Flags().SetAnnotation("db-password", secretAnnotation, []string{"true"})
	rootCmd.AddCommand(sub)
	defer rootCmd.RemoveCommand(sub)
	if err := sub.Flags().Set("db-password", "DBPASSWORD42"); err != nil {
		t.Fatal(err)
	}
	if got := errorText("create user: bad password DBPASSWORD42"); strings.Contains(got, "DBPASSWORD42") {
		t.Errorf("secret flag value leaked: %s", got)
	}
}
