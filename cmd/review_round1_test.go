package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// A missing required flag is a usage mistake: the usage text stays.
func TestRequiredFlagErrorKeepsUsage(t *testing.T) {
	c := &cobra.Command{Use: "x", RunE: func(*cobra.Command, []string) error { return nil }}
	c.Flags().String("need", "", "")
	c.MarkFlagRequired("need")
	c.PersistentPreRunE = rootCmd.PersistentPreRunE
	c.SetArgs([]string{})
	c.SilenceErrors = true
	if err := c.Execute(); err == nil || !strings.Contains(err.Error(), "need") {
		t.Fatalf("want required flag error, got %v", err)
	}
	if c.SilenceUsage {
		t.Error("usage was switched off for a missing required flag")
	}
}

// Completion lines stay on stdout in text output and leave it for -o json.
func TestResultLineGoesToStdoutOnlyInTextOutput(t *testing.T) {
	oldOut, oldQ := output, query
	defer func() { output, query = oldOut, oldQ }()
	output, query = "table", ""
	if got := captureStdout(t, func() { resultLinef("Upload complete: a\n") }); got != "Upload complete: a\n" {
		t.Errorf("text output: %q", got)
	}
	output = "json"
	if got := captureStdout(t, func() { resultLinef("Upload complete: a\n") }); got != "" {
		t.Errorf("-o json stdout: %q", got)
	}
}

// The old list form "--notify-email ops@example.com" must fail, not create a
// group with email on and the address dropped.
func TestPostgreSQLNotificationGroupRejectsStrayArgs(t *testing.T) {
	if err := postgresqlCreateNotificationGroupCmd.Args(postgresqlCreateNotificationGroupCmd, []string{"ops@example.com"}); err == nil {
		t.Fatal("a stray positional value was accepted")
	}
}
