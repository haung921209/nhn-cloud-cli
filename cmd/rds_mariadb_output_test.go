package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// runMariaDBCmd runs a hand-written rds-mariadb command against a fake API
// that answers every call with reply (the token request included) and
// returns the last API request URL and what the command printed to stdout.
func runMariaDBCmd(t *testing.T, c *cobra.Command, flags map[string]string, reply string) (*url.URL, string) {
	t.Helper()
	var last url.URL
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/oauth2/") {
			io.WriteString(w, `{"access_token":"t","token_type":"Bearer","expires_in":3600}`)
			return
		}
		last = *r.URL
		io.WriteString(w, reply)
	}))
	defer srv.Close()
	target, _ := url.Parse(srv.URL)
	oldTr, oldOut := http.DefaultTransport, output
	defer func() { http.DefaultTransport, output = oldTr, oldOut }()
	http.DefaultTransport = &hostRewriteTransport{inner: oldTr, hosts: map[string]*url.URL{
		hostKey("https", "kr1-rds-mariadb.api.nhncloudservice.com"): target,
		hostKey("https", "oauth.api.nhncloudservice.com"):           target,
	}}
	t.Setenv("HOME", t.TempDir())
	t.Setenv("NHN_REGION", "kr1")
	t.Setenv("NHN_MARIADB_APP_KEY", "ak")
	t.Setenv("NHN_MARIADB_ACCESS_KEY_ID", "id")
	t.Setenv("NHN_MARIADB_SECRET_ACCESS_KEY", "sk")
	output = "json"
	for k, v := range flags {
		f := c.Flags().Lookup(k)
		old := f.Value.String()
		if sv, ok := f.Value.(interface{ Replace([]string) error }); ok {
			defer sv.Replace(nil)
		} else {
			defer f.Value.Set(old)
		}
		if err := f.Value.Set(v); err != nil {
			t.Fatal(err)
		}
		f.Changed = true
		defer func() { f.Changed = false }()
	}
	out := captureStdout(t, func() { c.Run(c, nil) })
	return &last, out
}

const mariadbOK = `"header":{"isSuccessful":true,"resultCode":0,"resultMessage":"SUCCESS"}`

// Write commands print the API response for -o json, not a text line.
func TestMariaDBWriteCommandsPrintJSON(t *testing.T) {
	_, out := runMariaDBCmd(t, createMariaDBParameterGroupCmd,
		map[string]string{"db-parameter-group-name": "n", "db-parameter-group-family": "MARIADB_V1011"},
		`{`+mariadbOK+`,"parameterGroupId":"pg-1"}`)
	var pg struct{ ParameterGroupID string }
	if err := json.Unmarshal([]byte(out), &pg); err != nil || pg.ParameterGroupID != "pg-1" {
		t.Errorf("create-db-parameter-group -o json printed %q", out)
	}

	_, out = runMariaDBCmd(t, createMariaDBInstanceCmd, map[string]string{
		"db-instance-identifier": "db1", "db-flavor-id": "f", "engine-version": "v", "master-username": "u",
		"master-user-password": "Pw-1234!", "subnet-id": "s", "availability-zone": "z", "db-parameter-group-id": "pg",
	}, `{`+mariadbOK+`,"jobId":"job-1"}`)
	var job struct{ JobID, DBInstanceName string }
	if err := json.Unmarshal([]byte(out), &job); err != nil || job.JobID != "job-1" || job.DBInstanceName != "db1" {
		t.Errorf("create-db-instance -o json printed %q", out)
	}
}

// Progress lines must not come before the JSON on stdout.
func TestMariaDBDescribeInstanceStdoutIsJSON(t *testing.T) {
	_, out := runMariaDBCmd(t, describeMariaDBInstanceCmd,
		map[string]string{"db-instance-identifier": "00000000-0000-0000-0000-000000000001"},
		`{`+mariadbOK+`,"dbInstanceId":"00000000-0000-0000-0000-000000000001"}`)
	if !json.Valid([]byte(out)) {
		t.Errorf("stdout is not JSON: %q", out)
	}
}

func TestMariaDBMetricStatisticsSendsMeasureNames(t *testing.T) {
	req, out := runMariaDBCmd(t, mariadbGetMetricStatisticsCmd, map[string]string{
		"db-instance-identifier": "00000000-0000-0000-0000-000000000001", "measure-names": "CPU_USAGE,MEMORY_USAGE",
		"from": "2026-11-18T00:00:00.000+09:00", "to": "2026-11-18T01:00:00.000+09:00",
	}, `{`+mariadbOK+`,"metricStatistics":[{"measureName":"CPU_USAGE","unit":"%","values":[[1679298540,"7.5%"]]}]}`)
	q := req.Query()
	if len(q["measureNames"]) != 2 || q.Get("from") != "2026-11-18T00:00:00.000+09:00" || q.Has("interval") {
		t.Errorf("query %q", req.RawQuery)
	}
	if !strings.Contains(out, `"7.5%"`) {
		t.Errorf("printed %q", out)
	}
}
