package cert

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The DB password must not appear in the client's argv (visible in ps and in
// the printed command); it goes through the environment instead.
func TestConnectionCommandKeepsPasswordOutOfArgv(t *testing.T) {
	h := &Helper{store: &CertificateStore{}}
	for _, svc := range []string{"rds-mysql", "rds-mariadb", "rds-postgresql"} {
		args, err := h.GetConnectionCommand(svc, "db.example", "3306", "app", "user", "S3cretPW", "kr1", "id", "", false, "")
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(strings.Join(args, " "), "S3cretPW") {
			t.Errorf("%s: password in argv: %v", svc, args)
		}
		env := strings.Join(ConnectionEnv(svc, "S3cretPW"), " ")
		if !strings.Contains(env, "S3cretPW") {
			t.Errorf("%s: password not passed through the environment", svc)
		}
	}
	if ConnectionEnv("rds-mysql", "") != nil {
		t.Error("no password means no environment entry")
	}
}

// The store holds client private keys: directory 0700, files 0600, including
// a directory left world-readable by an older version.
func TestCertificateStoreIsPrivate(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	base := filepath.Join(home, DefaultCertDir)
	if err := os.MkdirAll(base, 0755); err != nil {
		t.Fatal(err)
	}
	cs, err := NewCertificateStore()
	if err != nil {
		t.Fatal(err)
	}
	info, err := cs.StoreCertificate(&CertificateRequest{Type: "CLIENT-KEY", ServiceType: "rds-mysql", Region: "kr1", CertData: []byte("-----BEGIN PRIVATE KEY-----\nx\n")})
	if err != nil {
		t.Fatal(err)
	}
	check := func(path string, want os.FileMode) {
		t.Helper()
		st, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := st.Mode().Perm(); got != want {
			t.Errorf("%s: mode %o, want %o", filepath.Base(path), got, want)
		}
	}
	check(base, 0700)
	check(info.FilePath, 0600)
	check(cs.MetadataFile, 0600)
}

// Opening the store also fixes files an older version left readable.
func TestCertificateStoreFixesExistingFiles(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	base := filepath.Join(home, DefaultCertDir)
	if err := os.MkdirAll(base, 0755); err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(base, "rds-mysql-kr1-CLIENT-KEY-old.pem")
	if err := os.WriteFile(old, []byte("key"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := NewCertificateStore(); err != nil {
		t.Fatal(err)
	}
	if st, _ := os.Stat(old); st.Mode().Perm() != 0600 {
		t.Errorf("existing key file mode %o, want 600", st.Mode().Perm())
	}
}
