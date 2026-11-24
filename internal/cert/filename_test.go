package cert

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestFileNamePartDropsPathCharacters(t *testing.T) {
	got := fileNamePart("../../etc/x y")
	if strings.ContainsAny(got, "/ \\\\") {
		t.Errorf("got %q", got)
	}
	if fileNamePart("rds-mysql_8.0") != "rds-mysql_8.0" {
		t.Error("plain value changed")
	}
}

func TestStoreCertificateTypeStaysInStore(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	store, err := NewCertificateStore()
	if err != nil {
		t.Fatal(err)
	}
	info, err := store.StoreCertificate(&CertificateRequest{ServiceType: "rds-mysql", Region: "kr1", Type: "/../../x", CertData: []byte("pem")})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(info.FilePath) != store.BaseDir {
		t.Fatalf("written outside the store: %s", info.FilePath)
	}
}
