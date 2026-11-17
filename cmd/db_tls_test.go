package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
)

func writePEM(t *testing.T, der []byte) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(p, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

// newCADER returns a fresh self-signed CA certificate, unrelated to the
// httptest certificate.
func newCADER(t *testing.T) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "other test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return der
}

// The registered CA must accept the server it signed and reject others.
func TestCAVerifyingTLSConfig(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()

	dial := func(caPath string) error {
		tc, err := caVerifyingTLSConfig(caPath)
		if err != nil {
			return err
		}
		c, err := tls.Dial("tcp", srv.Listener.Addr().String(), tc)
		if err == nil {
			c.Close()
		}
		return err
	}
	if err := dial(writePEM(t, srv.Certificate().Raw)); err != nil {
		t.Errorf("matching CA rejected: %v", err)
	}
	if err := dial(writePEM(t, newCADER(t))); err == nil {
		t.Error("a different CA must be rejected")
	}
}

func TestMySQLDSN(t *testing.T) {
	dsn, err := mysqlDSN("u", "p@ss:w/rd?&", "db.example", "3306", "app", "")
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil || cfg.Passwd != "p@ss:w/rd?&" || cfg.TLSConfig != "skip-verify" || cfg.Addr != "db.example:3306" {
		t.Fatalf("no CA: %+v %v", cfg, err)
	}
	srv := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()
	dsn, err = mysqlDSN("u", "p", "db.example", "3306", "app", writePEM(t, srv.Certificate().Raw))
	if err != nil {
		t.Fatal(err)
	}
	if cfg, _ := mysql.ParseDSN(dsn); cfg == nil || cfg.TLSConfig != mysqlTLSName {
		t.Errorf("with CA: %q", dsn)
	}
}

func TestPQConnStringQuotes(t *testing.T) {
	s := pqConnString("h", "5432", "u", `a b'c\d`, "db", "require", "")
	if !strings.Contains(s, `password='a b\'c\\d'`) {
		t.Errorf("got %s", s)
	}
}

// Without a CA the client asks for TLS first and drops to plain text only
// when the server says it does not support TLS.
func TestOpenPostgresFallsBackOnlyWhenServerRefusesTLS(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	var sslRequests, plainStarts int32
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			buf := make([]byte, 8)
			if _, err := io.ReadFull(c, buf); err == nil {
				if binary.BigEndian.Uint32(buf[4:]) == 80877103 { // SSLRequest
					atomic.AddInt32(&sslRequests, 1)
					c.Write([]byte("N"))
				} else {
					atomic.AddInt32(&plainStarts, 1)
				}
			}
			c.Close()
		}
	}()
	host, port, _ := net.SplitHostPort(ln.Addr().String())
	_, err = openPostgres(host, port, "u", "p", "db", "")
	if err == nil {
		t.Fatal("the fake server never completes a login")
	}
	if atomic.LoadInt32(&sslRequests) != 1 || atomic.LoadInt32(&plainStarts) != 1 {
		t.Errorf("want one TLS attempt then one plain attempt, got ssl=%d plain=%d", sslRequests, plainStarts)
	}
}
