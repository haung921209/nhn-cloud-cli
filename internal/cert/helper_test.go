package cert

import (
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
