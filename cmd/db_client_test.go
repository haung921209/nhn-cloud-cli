package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// When this test binary runs as the fake DB client, it reports the option
// file it was given (first argument) and exits before the test framework
// parses flags.
func init() {
	if os.Getenv("GAPS_FAKE_DB_CLIENT") != "1" {
		return
	}
	for _, a := range os.Args[1:] {
		if p, ok := strings.CutPrefix(a, "--defaults-extra-file="); ok {
			data, err := os.ReadFile(p)
			fmt.Printf("first=%t path=%s content=%q err=%v\n", a == os.Args[1], p, data, err)
		}
	}
	os.Exit(3) // the client fails; the secret file must still be removed
}

func TestRunWithSecretPassesOptionFileAndCleansUp(t *testing.T) {
	var out bytes.Buffer
	c := exec.Command(os.Args[0], "-h", "db.example", "app")
	c.Env = append(os.Environ(), "GAPS_FAKE_DB_CLIENT=1")
	c.Stdout = &out
	err := runWithSecret(c, "rds-mysql", "S3cretPW")
	if err == nil {
		t.Fatal("the fake client exits 3; runWithSecret must return its error")
	}
	got := out.String()
	if !strings.Contains(got, "first=true") || !strings.Contains(got, "S3cretPW") {
		t.Fatalf("client output: %q", got)
	}
	for _, a := range c.Args {
		if strings.Contains(a, "S3cretPW") {
			t.Errorf("password in argv: %v", c.Args)
		}
	}
	path := got[strings.Index(got, "path=")+5 : strings.Index(got, " content=")]
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("option file %s left behind after a failed client", path)
	}
}
