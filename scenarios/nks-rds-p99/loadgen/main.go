// Tiny Go OLTP load generator for the rds-postgresql-p99 scenario.
//
// Replaces sysbench (which can't disable TLS verification on macOS brew
// libmysqlclient builds — NHN Cloud's RDS server cert isn't in the local
// truststore). Uses Go's pgx driver for PostgreSQL with a custom TLSConfig
// injected via pgx.ParseConfig + stdlib.RegisterConnConfig.
//
// PostgreSQL TLS behaviour (NHN Cloud constraints):
//   - NHN's pg_hba.conf rejects sslmode=disable for external clients (run-7
//     failure: "no encryption").
//   - NHN's server cert SAN doesn't match the auto-generated hostname, so Go's
//     default TLS handshake fails with cert verification errors (run-6 failure).
//   - The sslmode=require/verify-ca DSN parameters cannot override the TLS
//     Config once pgx decides the cert is untrusted — only ConnConfig.TLSConfig
//     provides enough control.
//
// Solution: always use TLS (satisfies pg_hba) but control verification via
// pgx's native ConnConfig API:
//   CA_FILE set   → real chain verification with fetched CA; hostname check
//                   skipped (InsecureSkipVerify=true + VerifyPeerCertificate).
//   CA_FILE empty → TLS encrypted but fully unverified (InsecureSkipVerify=true
//                   only). Safe for smoke tests; not for production.
//   sslmode=disable path removed — NHN's pg_hba forbids plain TCP for external.
//
// LOADGEN_ENGINE=pg routes to PostgreSQL (pgx); default `mysql` keeps the
// MySQL/MariaDB code path so this file can be reused across engines.
//
// Operation mix is a tiny OLTP-ish point-update workload — not a sysbench
// equivalent, but enough to surface tail latency on a warm-cache row set.
//
// Args via env:
//   DB_HOST       (required) — DB host
//   DB_PORT       (default 3306 for mysql; 5432 for pg)
//   DB_USER       (required) — username
//   DB_PASSWORD / MYSQL_PWD  (required) — password (env, never argv)
//   DB_NAME       (required) — database name
//   PHASE         "grants" | "prepare" | "run"
//   LOADGEN_ENGINE "mysql" (default) | "pg"
//   THREADS       (default 8)         — used in run
//   DURATION_S    (default 120)       — used in run
//   ROWS          (default 10000)     — used in prepare
//   TABLES        (default 4)         — used in prepare
//   GRANT_USER    (required for grants phase) — user to GRANT to
//   GRANT_PASSWORD (required for grants phase) — password for GRANT_USER
//   CA_FILE       (optional)          — TLS CA file for verify-ca
//
// Outputs:
//   grants : prints "[loadgen] grants done: GRANT_USER=… DB_NAME=…"
//   prepare: prints "[loadgen] prepare done: tables=N rows=M"
//   run    : prints final line "tps=… p50=… p95=… p99=… max=… errs=…"
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	mysqld "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func main() {
	host := mustEnv("DB_HOST")
	user := mustEnv("DB_USER")
	dbname := mustEnv("DB_NAME")
	phase := mustEnv("PHASE")

	// PG grants phase requires connecting AS the master role to CREATE USER /
	// GRANT (the app user is the GRANT target — it doesn't exist yet, can't
	// be used to authenticate). Swap to MASTER_DB_USER/_PASSWORD when both
	// are set AND we're in the grants phase. mysql/mariadb hit the no-op
	// grants branch below so this override is also harmless for those.
	if phase == "grants" {
		if mu := os.Getenv("MASTER_DB_USER"); mu != "" {
			if mp := os.Getenv("MASTER_DB_PASSWORD"); mp != "" {
				user = mu
				os.Setenv("DB_PASSWORD", mp) // password resolved below from DB_PASSWORD
			}
		}
	}

	engine := strings.ToLower(envOr("LOADGEN_ENGINE", "mysql"))

	defaultPort := "3306"
	if engine == "pg" || engine == "postgres" || engine == "postgresql" {
		defaultPort = "5432"
		engine = "pg"
	}
	port := envOr("DB_PORT", defaultPort)

	// Password: prefer DB_PASSWORD; fall back to MYSQL_PWD for back-compat.
	pwd := os.Getenv("DB_PASSWORD")
	if pwd == "" {
		pwd = os.Getenv("MYSQL_PWD")
	}
	if pwd == "" {
		fmt.Fprintln(os.Stderr, "missing env: DB_PASSWORD (or MYSQL_PWD)")
		os.Exit(2)
	}

	threads := envInt("THREADS", 8)
	durSec := envInt("DURATION_S", 120)
	rows := envInt("ROWS", 10000)
	tables := envInt("TABLES", 4)

	caPath := os.Getenv("CA_FILE")

	var (
		driverName string
		dsn        string
	)

	if engine == "pg" {
		// SSL mode is environment-controlled. Round-2's external setup hit pg_hba
		// rules that required TLS so we used `require`; round-5's intra-VPC path
		// hits NHN PG instances that don't advertise SSL at all (server returns
		// "server refused TLS connection" on every TLS handshake attempt). Default
		// `prefer` makes pgx try TLS first and silently fall back to plain when
		// the server rejects — both scenarios succeed, and 04-deploy-loadgen can
		// override via PG_SSLMODE if a deployment ever needs the strict path back.
		sslmode := os.Getenv("PG_SSLMODE")
		if sslmode == "" {
			sslmode = "prefer"
		}
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&connect_timeout=30",
			url.QueryEscape(user), url.QueryEscape(pwd),
			host, port, url.QueryEscape(dbname), url.QueryEscape(sslmode),
		)
		pgxCfg, err := pgx.ParseConfig(connStr)
		must(err)

		// Custom TLS: NHN's pg_hba requires TLS (sslmode=disable is rejected), but
		// the server cert SAN doesn't match the auto-generated hostname, so standard
		// verify-full / verify-ca handshake also fails. Mirror MySQL/MariaDB loadgen:
		// always encrypt, skip hostname check.
		if caPath != "" {
			caPEM, err := os.ReadFile(caPath)
			must(err)
			pool := x509.NewCertPool()
			if !pool.AppendCertsFromPEM(caPEM) {
				fmt.Fprintf(os.Stderr, "[loadgen] ERROR: %s did not contain a valid PEM CA chain\n", caPath)
				os.Exit(1)
			}
			pgxCfg.TLSConfig = &tls.Config{
				InsecureSkipVerify:    true,
				VerifyPeerCertificate: makePGChainOnlyVerifier(pool),
			}
			fmt.Fprintf(os.Stderr,
				"[loadgen] using CA file %s for chain verification (hostname check skipped — NHN cert has no matching SAN)\n", caPath)
		} else {
			// No CA — TLS still required by NHN's pg_hba, skip ALL verification.
			pgxCfg.TLSConfig = &tls.Config{InsecureSkipVerify: true}
			fmt.Fprintln(os.Stderr,
				"[loadgen] WARNING: CA_FILE not set — TLS with skip-verify (unsafe outside smoke tests). "+
					"Run steps/02b-fetch-ca.sh to fetch the NHN CA via the v1.0 cert-export API.")
		}

		connStrName := stdlib.RegisterConnConfig(pgxCfg)
		db, err := sql.Open("pgx", connStrName)
		must(err)
		defer db.Close()
		db.SetMaxOpenConns(threads*2 + 4)
		db.SetMaxIdleConns(threads + 2)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		must(db.PingContext(ctx))
		cancel()

		switch phase {
		case "grants":
			grantUser := mustEnv("GRANT_USER")
			grantPwd := mustEnv("GRANT_PASSWORD")
			doGrants(db, engine, grantUser, grantPwd, dbname)
		case "prepare":
			doPrepare(db, engine, tables, rows)
		case "run":
			doRun(db, engine, threads, durSec, tables, rows)
		default:
			fmt.Fprintf(os.Stderr, "unknown PHASE %q\n", phase)
			os.Exit(2)
		}
		return
	} else {
		driverName = "mysql"
		tlsName := configureMysqlTLS(caPath)
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?tls=%s&parseTime=true&interpolateParams=true&timeout=30s&readTimeout=30s",
			user, pwd, host, port, dbname, tlsName,
		)
	}

	db, err := sql.Open(driverName, dsn)
	must(err)
	defer db.Close()
	db.SetMaxOpenConns(threads*2 + 4)
	db.SetMaxIdleConns(threads + 2)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	must(db.PingContext(ctx))
	cancel()

	switch phase {
	case "grants":
		grantUser := mustEnv("GRANT_USER")
		grantPwd := mustEnv("GRANT_PASSWORD")
		doGrants(db, engine, grantUser, grantPwd, dbname)
	case "prepare":
		doPrepare(db, engine, tables, rows)
	case "run":
		doRun(db, engine, threads, durSec, tables, rows)
	default:
		fmt.Fprintf(os.Stderr, "unknown PHASE %q\n", phase)
		os.Exit(2)
	}
}

// doGrants connects as master and creates the application user (if needed),
// then grants enough privileges to do everything the prepare + run phases need
// (CREATE TABLE, INSERT/UPDATE/SELECT, TRUNCATE).
func doGrants(db *sql.DB, engine, user, pwd, dbname string) {
	if engine != "pg" {
		// MySQL/MariaDB grants are typically handled by NHN's create-db-user
		// API with --authority-type DDL; nothing to do here.
		fmt.Printf("[loadgen] grants: engine=%s — no-op (handled by API)\n", engine)
		return
	}
	// PostgreSQL: First create the user (idempotent), then GRANT on database,
	// public schema, all existing tables, and default privileges so future
	// tables created by this user are owned by them.

	// Idempotent CREATE USER: PostgreSQL doesn't have CREATE USER IF NOT EXISTS,
	// so wrap in DO $$ ... EXCEPTION WHEN duplicate_object THEN NULL; END $$.
	// If the user already exists, we also ALTER to ensure the password matches
	// the caller's provided GRANT_PASSWORD.
	createUserSQL := fmt.Sprintf(`
DO $do$
BEGIN
    CREATE USER %q WITH PASSWORD '%s';
EXCEPTION WHEN duplicate_object THEN
    -- user already exists; reset password to ensure caller's APP_DB_PASSWORD matches
    ALTER USER %q WITH PASSWORD '%s';
END
$do$`,
		user, pwd, user, pwd)

	if _, err := db.Exec(createUserSQL); err != nil {
		fmt.Fprintf(os.Stderr, "[loadgen] grants ERROR: CREATE USER %q failed: %v\n", user, err)
		os.Exit(1)
	}
	fmt.Printf("[loadgen] grants: ensured user %q exists with password\n", user)

	stmts := []string{
		// Connect privilege on the bench database.
		`GRANT CONNECT ON DATABASE "` + dbname + `" TO "` + user + `"`,
		// CREATE on public schema → user can issue CREATE TABLE.
		`GRANT USAGE, CREATE ON SCHEMA public TO "` + user + `"`,
		// Existing-tables privileges (idempotent if no tables yet).
		`GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON ALL TABLES IN SCHEMA public TO "` + user + `"`,
		// Default privileges so any tables created later by master ALSO
		// reach the app user (defensive — prepare phase creates tables AS
		// the app user, so they own them and don't strictly need this).
		`ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON TABLES TO "` + user + `"`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			// 0LP01 (invalid_grant_operation) and 42501 (insufficient_privilege)
			// can happen if the master user is itself not a superuser. Surface
			// loud but don't abort — let the prepare phase confirm.
			fmt.Fprintf(os.Stderr, "[loadgen] grants WARN: %q: %v\n", s, err)
			continue
		}
	}
	fmt.Printf("[loadgen] grants done: GRANT_USER=%s DB_NAME=%s\n", user, dbname)
}

func doPrepare(db *sql.DB, engine string, tables, rows int) {
	for t := 1; t <= tables; t++ {
		tname := fmt.Sprintf("t_oltp_%d", t)

		var createDDL string
		if engine == "pg" {
			// PostgreSQL: no inline INDEX clause; use a separate CREATE INDEX.
			createDDL = `CREATE TABLE IF NOT EXISTS ` + tname + ` (
                id INT PRIMARY KEY,
                k INT NOT NULL,
                c VARCHAR(120) NOT NULL
            )`
		} else {
			createDDL = `CREATE TABLE IF NOT EXISTS ` + tname + ` (
                id INT PRIMARY KEY,
                k INT NOT NULL,
                c VARCHAR(120) NOT NULL,
                INDEX idx_k (k)
            )`
		}
		_, err := db.Exec(createDDL)
		must(err)

		if engine == "pg" {
			// Separate index DDL.
			_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_k_` + tname + ` ON ` + tname + ` (k)`)
		}

		// Truncate is supported in both engines.
		_, _ = db.Exec("TRUNCATE TABLE " + tname)

		// Bulk-load via batched inserts.
		const batchSize = 500
		for start := 1; start <= rows; start += batchSize {
			end := start + batchSize - 1
			if end > rows {
				end = rows
			}
			args := make([]interface{}, 0, (end-start+1)*3)
			var placeholders strings.Builder
			for i := start; i <= end; i++ {
				if placeholders.Len() > 0 {
					placeholders.WriteByte(',')
				}
				if engine == "pg" {
					// pgx requires $1, $2, $3 numeric placeholders.
					base := (i - start) * 3
					fmt.Fprintf(&placeholders, "($%d,$%d,$%d)", base+1, base+2, base+3)
				} else {
					placeholders.WriteString("(?,?,?)")
				}
				args = append(args, i, rand.Intn(rows*10), randStr(80))
			}

			var q string
			if engine == "pg" {
				q = "INSERT INTO " + tname + " (id, k, c) VALUES " +
					placeholders.String() + " ON CONFLICT (id) DO NOTHING"
			} else {
				q = "INSERT IGNORE INTO " + tname + " (id, k, c) VALUES " +
					placeholders.String()
			}
			_, err := db.Exec(q, args...)
			must(err)
		}
	}
	fmt.Printf("[loadgen] prepare done: tables=%d rows_per_table=%d\n", tables, rows)
}

func doRun(db *sql.DB, engine string, threads, durSec, tables, rows int) {
	var (
		latencies []time.Duration
		muLat     sync.Mutex
		ops       int64
		errs      int64
	)
	var wg sync.WaitGroup
	deadline := time.Now().Add(time.Duration(durSec) * time.Second)

	fmt.Printf("[loadgen] run engine=%s threads=%d duration=%ds tables=%d rows_per_table=%d\n",
		engine, threads, durSec, tables, rows)

	// Pre-build query strings (engine-specific placeholders).
	var selQ, updQ string
	if engine == "pg" {
		selQ = "SELECT c FROM %s WHERE id=$1"
		updQ = "UPDATE %s SET c=$1 WHERE id=$2"
	} else {
		selQ = "SELECT c FROM %s WHERE id=?"
		updQ = "UPDATE %s SET c=? WHERE id=?"
	}

	start := time.Now()
	for w := 0; w < threads; w++ {
		wg.Add(1)
		go func(wid int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(wid)))
			local := make([]time.Duration, 0, 1024)
			for time.Now().Before(deadline) {
				tname := fmt.Sprintf("t_oltp_%d", rng.Intn(tables)+1)
				rid := rng.Intn(rows) + 1
				t0 := time.Now()
				// 70R / 30W mix.
				var err error
				if rng.Intn(10) < 7 {
					var c string
					err = db.QueryRow(fmt.Sprintf(selQ, tname), rid).Scan(&c)
				} else {
					_, err = db.Exec(fmt.Sprintf(updQ, tname), randStr(80), rid)
				}
				d := time.Since(t0)
				if err != nil {
					atomic.AddInt64(&errs, 1)
					continue
				}
				atomic.AddInt64(&ops, 1)
				local = append(local, d)
				if len(local) >= 1024 {
					muLat.Lock()
					latencies = append(latencies, local...)
					muLat.Unlock()
					local = local[:0]
				}
			}
			if len(local) > 0 {
				muLat.Lock()
				latencies = append(latencies, local...)
				muLat.Unlock()
			}
		}(w)
	}
	wg.Wait()
	elapsed := time.Since(start).Seconds()

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	p := func(q float64) time.Duration {
		if len(latencies) == 0 {
			return 0
		}
		return latencies[int(float64(len(latencies)-1)*q)]
	}
	tps := float64(ops) / elapsed
	fmt.Printf(
		"tps=%.1f p50=%dms p95=%dms p99=%dms max=%dms ops=%d errs=%d duration=%.1fs\n",
		tps,
		p(0.50).Milliseconds(),
		p(0.95).Milliseconds(),
		p(0.99).Milliseconds(),
		p(1.00).Milliseconds(),
		ops, errs, elapsed,
	)
}

// makePGChainOnlyVerifier returns a VerifyPeerCertificate callback that checks
// the server certificate chain against the given CA pool without verifying the
// hostname (NHN Cloud's per-instance cert has no matching SAN). Used with
// pgxCfg.TLSConfig.InsecureSkipVerify=true so that Go's standard hostname
// check is suppressed while chain trust is still enforced.
func makePGChainOnlyVerifier(pool *x509.CertPool) func([][]byte, [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
		if len(rawCerts) == 0 {
			return fmt.Errorf("server presented no certificate")
		}
		leaf, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return fmt.Errorf("parse leaf cert: %w", err)
		}
		intermediates := x509.NewCertPool()
		for _, raw := range rawCerts[1:] {
			if c, perr := x509.ParseCertificate(raw); perr == nil {
				intermediates.AddCert(c)
			}
		}
		_, verr := leaf.Verify(x509.VerifyOptions{
			Roots:         pool,
			Intermediates: intermediates,
		})
		return verr
	}
}

// configureMysqlTLS registers a tls.Config with the MySQL driver and returns
// the name to put in the DSN's `tls=` parameter. Used only when LOADGEN_ENGINE
// is mysql/mariadb. PostgreSQL uses pgx's ConnConfig.TLSConfig API instead
// (see the engine=="pg" branch in main) and does not call this function.
func configureMysqlTLS(caPath string) string {
	if caPath == "" {
		fmt.Fprintln(os.Stderr,
			"[loadgen] WARNING: CA_FILE not set — falling back to TLS skip-verify. "+
				"This is unsafe outside smoke tests; run steps/02b-fetch-ca.sh "+
				"to fetch the NHN CA via the v4.0 cert-export API.")
		_ = mysqld.RegisterTLSConfig("nhn-skip-verify",
			&tls.Config{InsecureSkipVerify: true})
		return "nhn-skip-verify"
	}
	pem, err := os.ReadFile(caPath)
	must(err)
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(pem) {
		fmt.Fprintf(os.Stderr,
			"[loadgen] ERROR: %s did not contain a valid PEM CA chain\n", caPath)
		os.Exit(1)
	}
	_ = mysqld.RegisterTLSConfig("nhn-rds", &tls.Config{
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
			if len(rawCerts) == 0 {
				return fmt.Errorf("server presented no certificate")
			}
			leaf, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return fmt.Errorf("parse leaf cert: %w", err)
			}
			intermediates := x509.NewCertPool()
			for _, raw := range rawCerts[1:] {
				if c, perr := x509.ParseCertificate(raw); perr == nil {
					intermediates.AddCert(c)
				}
			}
			_, verr := leaf.Verify(x509.VerifyOptions{
				Roots:         pool,
				Intermediates: intermediates,
			})
			return verr
		},
	})
	fmt.Fprintf(os.Stderr,
		"[loadgen] using CA file %s for chain verification (hostname check skipped — NHN cert has no matching SAN)\n", caPath)
	return "nhn-rds"
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		fmt.Fprintf(os.Stderr, "missing env: %s\n", k)
		os.Exit(2)
	}
	return v
}

func envOr(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func envInt(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func randStr(n int) string {
	const cs = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = cs[rand.Intn(len(cs))]
	}
	return string(b)
}
