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

// A key file that already exists with a wider mode ends up 0600.
func TestWritePrivateFileResetsMode(t *testing.T) {
	p := filepath.Join(t.TempDir(), "k.pem")
	os.WriteFile(p, []byte("old"), 0644)
	if err := writePrivateFile(p, []byte("new")); err != nil {
		t.Fatal(err)
	}
	if fi, _ := os.Stat(p); fi.Mode().Perm() != 0600 {
		t.Errorf("mode %v", fi.Mode().Perm())
	}
}

// Commands whose API returns no body still print JSON with -o json.
func TestNoBodyWritesPrintJSON(t *testing.T) {
	oldOut, oldQ := output, query
	defer func() { output, query = oldOut, oldQ }()
	output, query = "json", ""
	got := captureStdout(t, func() {
		printWriteResult(map[string]string{"id": "img-1", "result": "deleted"}, "Image %s deleted successfully\n", "img-1")
	})
	if !strings.Contains(got, `"result": "deleted"`) || strings.Contains(got, "successfully") {
		t.Errorf("got %q", got)
	}
}

func TestPlainFileName(t *testing.T) {
	for _, ok := range []string{"web", "my-key.v2", "*.example.com"} {
		if err := plainFileName("--name", ok); err != nil {
			t.Errorf("%q rejected: %v", ok, err)
		}
	}
	for _, bad := range []string{"", "..", "../x", "a/b", `a\b`, "/etc/x"} {
		if plainFileName("--name", bad) == nil {
			t.Errorf("%q accepted", bad)
		}
	}
}
