package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// configure used to rewrite the file with [default] only: other profiles and
// keys it does not know (service appkeys) were lost, and an existing file kept
// its old mode.
func TestSaveConfigKeepsOtherProfilesAndKeys(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("NHN_CLOUD_PROFILE", "")
	oldProfile := profile
	profile = ""
	defer func() { profile = oldProfile }()
	dir := filepath.Join(home, ".nhncloud")
	os.MkdirAll(dir, 0700)
	path := filepath.Join(dir, "credentials")
	orig := "[default]\nregion = kr1\ndnsplus_app_key = DNSKEY\nrds_app_key = OLDRDS\n\n[staging]\nregion = kr2\napi_password = STAGEPW\n"
	if err := os.WriteFile(path, []byte(orig), 0644); err != nil {
		t.Fatal(err)
	}

	if err := saveConfig(&Config{Region: "kr1", Username: "u@example.com", RDSAppKey: "NEWRDS"}); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(path)
	got := string(b)
	for _, want := range []string{"[staging]", "api_password = STAGEPW", "region = kr2", "dnsplus_app_key = DNSKEY", "rds_app_key = NEWRDS", "username = u@example.com"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "OLDRDS") {
		t.Errorf("old value kept:\n%s", got)
	}
	if fi, _ := os.Stat(path); fi.Mode().Perm() != 0600 {
		t.Errorf("mode %v, want 0600", fi.Mode().Perm())
	}
}

func TestMergeProfileAddsMissingSection(t *testing.T) {
	got := mergeProfile("[default]\nregion = kr1\n", "dev", [][2]string{{"region", "kr2"}, {"tenant_id", ""}})
	if got != "[default]\nregion = kr1\n\n[dev]\nregion = kr2\n" {
		t.Errorf("got %q", got)
	}
}
