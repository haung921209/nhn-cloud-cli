package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/haung921209/nhn-cloud-cli/internal/catalog"
	"github.com/spf13/cobra"
)

func saveFieldService() catalog.Service {
	svc := fakeService()
	svc.Auth = catalog.Auth{}
	svc.Operations = []catalog.Operation{
		{Name: "get-config", Summary: "Get", Method: "GET", Path: "/v1/config", SaveField: "config"},
		{Name: "get-plain", Summary: "Get", Method: "GET", Path: "/v1/plain"},
	}
	return svc
}

func TestCatalogSaveFieldWritesOwnerOnlyFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "kube", "config") // parent is missing
	reply := `{"config":"apiVersion: v1\nusers:\n- name: SECRET-TOKEN"}`

	_, out, err := runFake(t, saveFieldService(), reply, "fake-svc", "get-config", "--output-file", file)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "SECRET-TOKEN") || !strings.Contains(out, file) {
		t.Errorf("stdout must name the file and not print the field: %q", out)
	}
	got, _ := os.ReadFile(file)
	if string(got) != "apiVersion: v1\nusers:\n- name: SECRET-TOKEN" { // the field, byte for byte
		t.Errorf("file content: %q", got)
	}
	if fi, _ := os.Stat(file); fi.Mode().Perm() != 0600 {
		t.Errorf("file mode %v, want 0600", fi.Mode().Perm())
	}
	if fi, _ := os.Stat(filepath.Dir(file)); fi.Mode().Perm() != 0700 {
		t.Errorf("new parent mode %v, want 0700", fi.Mode().Perm())
	}

	// An existing world-readable file becomes 0600 and is replaced.
	old := filepath.Join(dir, "old")
	os.WriteFile(old, []byte("previous"), 0644)
	if _, _, err := runFake(t, saveFieldService(), reply, "fake-svc", "get-config", "--output-file", old); err != nil {
		t.Fatal(err)
	}
	if fi, _ := os.Stat(old); fi.Mode().Perm() != 0600 {
		t.Errorf("existing file mode %v, want 0600", fi.Mode().Perm())
	}
	if got, _ := os.ReadFile(old); strings.Contains(string(got), "previous") {
		t.Errorf("existing file was not replaced: %q", got)
	}
}

func TestCatalogSaveFieldMissingWritesNothing(t *testing.T) {
	file := filepath.Join(t.TempDir(), "config")
	_, _, err := runFake(t, saveFieldService(), `{"other":"x"}`, "fake-svc", "get-config", "--output-file", file)
	if err == nil || !strings.Contains(err.Error(), `"config"`) {
		t.Fatalf("want an error naming the field, got %v", err)
	}
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Errorf("file must not be created: %v", err)
	}
}

func TestCatalogSaveFieldRefusals(t *testing.T) {
	dir := t.TempDir()
	kube := filepath.Join(dir, "kube")
	os.Mkdir(kube, 0755)
	reply := `{"config":"SECRET-TOKEN"}`

	// An existing directory is refused and keeps its mode.
	if _, _, err := runFake(t, saveFieldService(), reply, "fake-svc", "get-config", "--output-file", kube); err == nil || !strings.Contains(err.Error(), "not a regular file") {
		t.Errorf("a directory target must be refused: %v", err)
	}
	if fi, _ := os.Stat(kube); fi.Mode().Perm() != 0755 {
		t.Errorf("directory mode changed to %v", fi.Mode().Perm())
	}

	// An empty value is an error, not a fallback to printing the credential.
	_, out, err := runFake(t, saveFieldService(), reply, "fake-svc", "get-config", "--output-file", "")
	if err == nil || strings.Contains(out, "SECRET-TOKEN") {
		t.Errorf("empty --output-file: err %v, stdout %q", err, out)
	}

	// A 200 failure envelope is reported as a failure even with the field present.
	// The field is not printed with the failure either.
	file := filepath.Join(dir, "config")
	oldOutput := output
	output = "json"
	defer func() { output = oldOutput }()
	_, out, err = runFake(t, saveFieldService(), `{"header":{"isSuccessful":false,"resultCode":7,"resultMessage":"denied"},"config":"STALE-SECRET"}`,
		"fake-svc", "get-config", "--output-file", file)
	if err == nil || !strings.Contains(err.Error(), "denied") || strings.Contains(out, "STALE-SECRET") || !strings.Contains(out, "denied") {
		t.Errorf("failure envelope: err %v, stdout %q", err, out)
	}
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Errorf("file must not be written on a failure envelope: %v", err)
	}
}

func TestCatalogSaveFieldRejectsBadInputBeforeWriting(t *testing.T) {
	dir := t.TempDir()
	cases := []struct{ name, reply, file, q, want string }{
		{"empty field", `{"config":""}`, "a", "", `"config"`},
		{"non-string field", `{"config":5}`, "b", "", `"config"`},
		{"array response", `["config"]`, "c", "", "JSON object"},
		{"non-JSON response", `apiVersion: v1`, "d", "", "JSON object"},
		{"trailing slash", `{"config":"c"}`, "e/", "", "directory"},
		{"ends in dot", `{"config":"c"}`, "h/.", "", "directory"},
		{"ends in dot dot", `{"config":"c"}`, "i/..", "", "directory"},
		{"bad query", `{"config":"c"}`, "f", "[[[", "--query"},
		{"query that fails when run", `{"config":"c"}`, "g", "length(`1`)", "--query"},
	}
	for _, c := range cases {
		oldQuery := query
		query = c.q
		file := dir + "/" + c.file // Join would drop a trailing slash
		_, _, err := runFake(t, saveFieldService(), c.reply, "fake-svc", "get-config", "--output-file", file)
		query = oldQuery
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: want error with %q, got %v", c.name, c.want, err)
		}
		if _, err := os.Stat(dir + "/" + strings.Split(c.file, "/")[0]); !os.IsNotExist(err) {
			t.Errorf("%s: nothing may be created: %v", c.name, err)
		}
	}
}

func TestCatalogSaveFieldJSONOutputShape(t *testing.T) {
	file := filepath.Join(t.TempDir(), "config")
	oldOutput := output
	output = "json"
	defer func() { output = oldOutput }()
	_, out, err := runFake(t, saveFieldService(), `{"config":"SECRET"}`, "fake-svc", "get-config", "--output-file", file)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]string
	if json.Unmarshal([]byte(out), &got) != nil || len(got) != 2 || got["file"] != file || got["result"] != "written" {
		t.Errorf("-o json output: %q", out)
	}
}

func TestCatalogSaveFieldOnlyWhereDeclared(t *testing.T) {
	_, out, err := runFake(t, saveFieldService(), `{"config":"c"}`, "fake-svc", "get-config")
	if err != nil || !strings.Contains(out, `"config"`) {
		t.Errorf("without --output-file the response is printed: %q %v", out, err)
	}
	if _, _, err := runFake(t, saveFieldService(), `{}`, "fake-svc", "get-plain", "--output-file", "x"); err == nil {
		t.Error("an operation without saveField must not take --output-file")
	}
}

func TestNKSKubeconfigTakesOutputFile(t *testing.T) {
	services, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range services {
		for _, op := range s.Operations {
			if s.Name == "nks" && op.Name == "get-cluster-kubeconfig" {
				if op.SaveField != "config" {
					t.Errorf("saveField %q, want config", op.SaveField)
				}
				return
			}
		}
	}
	t.Fatal("nks get-cluster-kubeconfig not in the catalog")
}

func TestCatalogSaveFieldCaseOfHeaderKeys(t *testing.T) {
	// Not the failure envelope (keys differ in case): saved like any response.
	file := filepath.Join(t.TempDir(), "config")
	_, _, err := runFake(t, saveFieldService(), `{"Header":{"IsSuccessful":false},"config":"S"}`, "fake-svc", "get-config", "--output-file", file)
	if got, _ := os.ReadFile(file); err != nil || string(got) != "S" {
		t.Errorf("err %v, file %q", err, got)
	}
}

func TestRegisteredKubeconfigCommandTakesOutputFile(t *testing.T) {
	// A hand-written command of the same name would hide the catalog one.
	// Catalog commands carry "API guide:" in Long, so this holds in any test order.
	if c := findChild(nksCmd, "get-cluster-kubeconfig"); c != nil && !strings.Contains(c.Long, "API guide:") {
		t.Fatal("a hand-written nks get-cluster-kubeconfig hides the catalog command")
	}
	// Built the way root.go adds catalogs, on a fresh root (rootCmd stays as other tests expect).
	services, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	root := &cobra.Command{Use: "nhncloud"}
	addCatalogCommands(root, services)
	c, _, err := root.Find([]string{"nks", "get-cluster-kubeconfig"})
	if err != nil || c.Name() != "get-cluster-kubeconfig" || c.Flags().Lookup("output-file") == nil {
		t.Errorf("nks get-cluster-kubeconfig with --output-file not built: %v", err)
	}
}

func TestCatalogSaveFieldParentThroughDirectoryLink(t *testing.T) {
	root := t.TempDir()
	safe, other := filepath.Join(root, "safe"), filepath.Join(root, "other")
	os.MkdirAll(filepath.Join(other, "sub"), 0700)
	os.Mkdir(safe, 0700)
	os.Symlink(filepath.Join(other, "sub"), filepath.Join(safe, "hop"))
	file := safe + "/hop/../new/config" // the kernel reads this as other/new/config
	if _, _, err := runFake(t, saveFieldService(), `{"config":"S"}`, "fake-svc", "get-config", "--output-file", file); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(filepath.Join(other, "new", "config")); string(got) != "S" {
		t.Errorf("not written where the kernel resolves the path: %q", got)
	}
	if _, err := os.Stat(filepath.Join(safe, "new")); !os.IsNotExist(err) {
		t.Errorf("a stray directory was created: %v", err)
	}
}
