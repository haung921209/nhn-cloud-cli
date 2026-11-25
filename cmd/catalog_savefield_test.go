package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/haung921209/nhn-cloud-cli/internal/catalog"
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
	if string(got) != "apiVersion: v1\nusers:\n- name: SECRET-TOKEN\n" {
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
