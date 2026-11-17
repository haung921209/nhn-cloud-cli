package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
)

// TLS for the built-in (native) database connections.
// With a registered CA the server certificate chain is verified against it
// (host names are not compared, like VERIFY_CA). Without one the connection
// is still encrypted, but unverified, and the user is told so.

const mysqlTLSName = "nhncloud-ca"

const noCAWarning = "Warning: no CA certificate registered; the server certificate is not verified. Import it with: nhncloud config ca import"

// caVerifyingTLSConfig verifies the server chain against the CA in caPath.
func caVerifyingTLSConfig(caPath string) (*tls.Config, error) {
	pemData, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("read CA %s: %w", caPath, err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(pemData) {
		return nil, fmt.Errorf("no certificate found in CA file %s", caPath)
	}
	return &tls.Config{
		// The chain is verified below; the endpoint name often differs from
		// the certificate's names, so the host name check is skipped.
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(raw [][]byte, _ [][]*x509.Certificate) error {
			if len(raw) == 0 {
				return errors.New("server sent no certificate")
			}
			certs := make([]*x509.Certificate, len(raw))
			for i, r := range raw {
				c, err := x509.ParseCertificate(r)
				if err != nil {
					return err
				}
				certs[i] = c
			}
			inter := x509.NewCertPool()
			for _, c := range certs[1:] {
				inter.AddCert(c)
			}
			_, err := certs[0].Verify(x509.VerifyOptions{Roots: roots, Intermediates: inter})
			if err != nil {
				return fmt.Errorf("server certificate does not match the registered CA (re-import it with nhncloud config ca import): %w", err)
			}
			return nil
		},
	}, nil
}

// mysqlDSN builds the DSN for MySQL and MariaDB native connections.
func mysqlDSN(user, password, host, port, dbName, caPath string) (string, error) {
	cfg := mysql.NewConfig()
	cfg.User, cfg.Passwd = user, password
	cfg.Net, cfg.Addr, cfg.DBName = "tcp", net.JoinHostPort(host, port), dbName
	cfg.AllowNativePasswords = true
	if caPath == "" {
		fmt.Fprintln(os.Stderr, noCAWarning)
		cfg.TLSConfig = "skip-verify"
		return cfg.FormatDSN(), nil
	}
	tc, err := caVerifyingTLSConfig(caPath)
	if err != nil {
		return "", err
	}
	if err := mysql.RegisterTLSConfig(mysqlTLSName, tc); err != nil {
		return "", err
	}
	cfg.TLSConfig = mysqlTLSName
	return cfg.FormatDSN(), nil
}

// pqQuote quotes a value for a lib/pq key=value connection string.
func pqQuote(v string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(v) + "'"
}

func pqConnString(host, port, user, password, dbName, sslmode, caPath string) string {
	s := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pqQuote(host), pqQuote(port), pqQuote(user), pqQuote(password), pqQuote(dbName), sslmode)
	if caPath != "" {
		s += " sslrootcert=" + pqQuote(caPath)
	}
	return s
}

// openPostgres opens and pings a native PostgreSQL connection. With a CA it
// requires a verified chain. Without one it asks for TLS and falls back to
// plain text only when the server does not support TLS, with a warning.
func openPostgres(host, port, user, password, dbName, caPath string) (*sql.DB, error) {
	try := func(sslmode, ca string) (*sql.DB, error) {
		db, err := sql.Open("postgres", pqConnString(host, port, user, password, dbName, sslmode, ca))
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err != nil {
			db.Close()
			return nil, err
		}
		return db, nil
	}
	if caPath != "" {
		return try("verify-ca", caPath)
	}
	fmt.Fprintln(os.Stderr, noCAWarning)
	db, err := try("require", "")
	if errors.Is(err, pq.ErrSSLNotSupported) {
		fmt.Fprintln(os.Stderr, "Warning: the server does not support TLS; the connection is not encrypted.")
		return try("disable", "")
	}
	return db, err
}
