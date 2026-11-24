package sshkeys

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportRejectsPathNames(t *testing.T) {
	m := &Manager{baseDir: t.TempDir()}
	for _, bad := range []string{"", "..", "../x", "a/b", "/tmp/x"} {
		if _, err := m.Import(bad, "/nonexistent"); err == nil || err.Error()[:8] != "key name" {
			t.Errorf("%q: want a name error, got %v", bad, err)
		}
	}
}

func TestWritePrivateResetsMode(t *testing.T) {
	p := filepath.Join(t.TempDir(), "k.pem")
	if err := os.WriteFile(p, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := writePrivate(p, []byte("new")); err != nil {
		t.Fatal(err)
	}
	if st, _ := os.Stat(p); st.Mode().Perm() != 0600 {
		t.Fatalf("mode %v, want 0600", st.Mode().Perm())
	}
}
