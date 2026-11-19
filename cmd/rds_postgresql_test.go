package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

type pgCall struct {
	method, path string
	query        url.Values
	body         string
}

// runPGFake runs "nhncloud rds-postgresql <args>" against a fake API (and a
// fake token endpoint) and returns the API calls and stdout.
func runPGFake(t *testing.T, respBody string, args ...string) ([]pgCall, string) {
	t.Helper()
	var calls []pgCall
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/oauth2/") {
			w.Write([]byte(`{"access_token":"tok","token_type":"Bearer","expires_in":3600}`))
			return
		}
		b, _ := io.ReadAll(r.Body)
		calls = append(calls, pgCall{r.Method, r.URL.Path, r.URL.Query(), string(b)})
		w.Write([]byte(respBody))
	}))
	defer srv.Close()
	target, _ := url.Parse(srv.URL)
	oldTr := http.DefaultTransport
	http.DefaultTransport = &hostRewriteTransport{inner: oldTr, hosts: map[string]*url.URL{
		"kr1-rds-postgres.api.nhncloudservice.com": target,
		"oauth.api.nhncloudservice.com":            target,
	}}
	defer func() { http.DefaultTransport = oldTr }()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("NHN_REGION", "kr1")
	t.Setenv("NHN_POSTGRESQL_APP_KEY", "ak")
	t.Setenv("NHN_POSTGRESQL_ACCESS_KEY_ID", "id")
	t.Setenv("NHN_POSTGRESQL_SECRET_ACCESS_KEY", "sk")
	oldOut, oldQ := output, query
	defer func() { output, query = oldOut, oldQ }()

	rootCmd.SetArgs(append([]string{"rds-postgresql"}, args...))
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := rootCmd.Execute()
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	if err != nil {
		t.Fatalf("%v: %v", args, err)
	}
	return calls, string(out)
}

func mustJSON(t *testing.T, out string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out)
	}
	return m
}

const pgOK = `"header":{"resultCode":0,"resultMessage":"SUCCESS","isSuccessful":true}`

// Write commands print the API response with -o json (defects: text only).
func TestPGWriteCommandsPrintJSON(t *testing.T) {
	calls, out := runPGFake(t, `{`+pgOK+`,"jobId":"job-1"}`, "create-db-instance", "-o", "json",
		"--db-instance-name", "inst-a", "--database-name", "db", "--db-flavor-id", "f", "--db-version", "V",
		"--db-user-name", "u", "--db-password", "Pw-12345!", "--db-parameter-group-id", "pg",
		"--subnet-id", "s", "--availability-zone", "kr-pub-a")
	m := mustJSON(t, out)
	if m["jobId"] != "job-1" || m["dbInstanceName"] != "inst-a" {
		t.Errorf("create-db-instance: %v", m)
	}
	// The API's storage type names are those of describe-db-storage-types.
	if b := calls[len(calls)-1].body; !strings.Contains(b, `"storageType":"General SSD"`) {
		t.Errorf("create-db-instance default storage type: %s", b)
	}

	_, out = runPGFake(t, `{`+pgOK+`,"parameterGroupId":"pg-1"}`, "create-db-parameter-group", "-o", "json",
		"--db-parameter-group-name", "n", "--db-parameter-group-family", "V")
	if m := mustJSON(t, out); m["parameterGroupId"] != "pg-1" {
		t.Errorf("create-db-parameter-group: %v", m)
	}

	_, out = runPGFake(t, `{`+pgOK+`,"jobId":"job-2"}`, "stop-db-instance", "-o", "json",
		"--db-instance-identifier", "00000000-0000-0000-0000-000000000001")
	if m := mustJSON(t, out); m["jobId"] != "job-2" {
		t.Errorf("stop-db-instance: %v", m)
	}

	_, out = runPGFake(t, `{`+pgOK+`}`, "delete-user-group", "-o", "json", "--user-group-id", "ug-1")
	mustJSON(t, out)

	// Text output is unchanged.
	_, out = runPGFake(t, `{`+pgOK+`,"jobId":"job-3"}`, "start-db-instance", "-o", "table",
		"--db-instance-identifier", "00000000-0000-0000-0000-000000000001")
	if out != "Instance start initiated.\nJob ID: job-3\n" {
		t.Errorf("text output %q", out)
	}
}

// Progress lines go to stderr so -o json stays parseable.
func TestPGDescribeAndModifyKeepStdoutJSON(t *testing.T) {
	id := "00000000-0000-0000-0000-000000000001"
	for _, c := range []string{"describe-db-instance", "describe-db-instance-network", "describe-db-instance-storage"} {
		_, out := runPGFake(t, `{`+pgOK+`,"dbInstanceId":"x","createdYmdt":"2026-01-01T00:00:00+09:00"}`, c, "-o", "json", "--db-instance-identifier", id)
		mustJSON(t, out)
	}
	_, out := runPGFake(t, `{`+pgOK+`,"jobId":"j"}`, "modify-db-instance-storage", "-o", "json", "--db-instance-identifier", id, "--storage-size", "30")
	mustJSON(t, out)
}

func TestPGDescribeSnapshotsSendsPageAndSize(t *testing.T) {
	calls, out := runPGFake(t, `{`+pgOK+`,"totalCounts":0,"backups":[]}`, "describe-db-snapshots", "-o", "json",
		"--db-instance-identifier", "00000000-0000-0000-0000-000000000001", "--page", "2", "--size", "10")
	mustJSON(t, out)
	q := calls[len(calls)-1].query
	if q.Get("page") != "2" || q.Get("size") != "10" {
		t.Errorf("query %v", q)
	}
}

func TestPGGetMetricStatisticsQuery(t *testing.T) {
	calls, out := runPGFake(t, `{`+pgOK+`,"metricStatistics":[{"metricName":"CPU_USAGE","unit":"%","values":[{"timestamp":1,"value":1.5}]}]}`,
		"get-metric-statistics", "-o", "json", "--db-instance-identifier", "00000000-0000-0000-0000-000000000001",
		"--metric-names", "CPU_USAGE,MEMORY_USAGE", "--from", "2026-11-17T00:00:00.000+09:00", "--to", "2026-11-17T01:00:00.000+09:00", "--interval", "5")
	mustJSON(t, out)
	q := calls[len(calls)-1].query
	if q.Get("from") != "2026-11-17T00:00:00.000+09:00" || len(q["metricNames"]) != 2 || q.Get("interval") != "5" {
		t.Errorf("query %v", q)
	}
	if f := findChild(rdsPostgreSQLCmd, "get-metric-statistics").Flags().Lookup("interval"); !strings.Contains(f.Usage, "minutes") {
		t.Errorf("--interval usage %q", f.Usage)
	}
}

func TestPGStorageTypesAndNoLogFilesCommand(t *testing.T) {
	_, out := runPGFake(t, `{`+pgOK+`,"storageTypes":["General SSD","General HDD"]}`, "describe-db-storage-types", "-o", "table")
	if !strings.Contains(out, "General SSD") {
		t.Errorf("storage types %q", out)
	}
	// GET /v1.0/db-instances/{id}/log-files is not in the RDS for PostgreSQL API.
	if findChild(rdsPostgreSQLCmd, "describe-logs") != nil {
		t.Error("describe-logs calls an endpoint the API does not have")
	}
}
