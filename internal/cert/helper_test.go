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
	}
}

// MySQL and MariaDB get the password from a private option file that wins
// over ~/.my.cnf; PostgreSQL from PGPASSWORD. Nothing is left behind.
func TestClientSecret(t *testing.T) {
	const pw = `p"w\x#;`
	args, env, cleanup, err := ClientSecret("rds-mysql", pw)
	if err != nil || len(args) != 1 || env != nil || !strings.HasPrefix(args[0], "--defaults-extra-file=") {
		t.Fatalf("mysql: %v %v %v", args, env, err)
	}
	path := strings.TrimPrefix(args[0], "--defaults-extra-file=")
	st, err := os.Stat(path)
	if err != nil || st.Mode().Perm() != 0600 {
		t.Fatalf("option file: %v %v", st, err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "[client]\npassword=\"p\\\"w\\\\x#;\"\n" {
		t.Errorf("option file content: %q", data)
	}
	cleanup()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("cleanup must remove the option file")
	}

	args, env, cleanup, err = ClientSecret("rds-postgresql", pw)
	if err != nil || args != nil || len(env) != 1 || env[0] != "PGPASSWORD="+pw {
		t.Fatalf("postgresql: %v %v %v", args, env, err)
	}
	cleanup()
	if args, env, _, _ := ClientSecret("rds-mysql", ""); args != nil || env != nil {
		t.Error("no password means nothing to pass")
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

// A store the user cannot chmod (root owned, read-only mount) still opens.
func TestCertificateStoreOpensWhenChmodFails(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, DefaultCertDir), 0755); err != nil {
		t.Fatal(err)
	}
	old := chmod
	chmod = func(string, os.FileMode) error { return os.ErrPermission }
	defer func() { chmod = old }()
	if _, err := NewCertificateStore(); err != nil {
		t.Fatalf("store must open even when tightening fails: %v", err)
	}
}
