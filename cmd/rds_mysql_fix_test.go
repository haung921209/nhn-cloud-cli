package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const fakeMySQLInstanceID = "11111111-2222-3333-4444-555555555555"

type mysqlCall struct {
	method, path string
	query        url.Values
	body         map[string]interface{}
}

// runMySQLFake runs an rds-mysql command against a fake API that knows one
// instance named "mydb" and answers every other call with a job and the IDs
// of the create responses. It returns the last non-list call and stdout.
func runMySQLFake(t *testing.T, args ...string) (mysqlCall, string) {
	t.Helper()
	var last mysqlCall
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth2/token/create":
			io.WriteString(w, `{"access_token":"tok","expires_in":3600}`)
			return
		case r.Method == "GET" && r.URL.Path == "/v4.0/db-instances":
			io.WriteString(w, `{"header":{"resultCode":0,"isSuccessful":true},"dbInstances":[{"dbInstanceId":"`+fakeMySQLInstanceID+`","dbInstanceName":"mydb"}]}`)
			return
		}
		last = mysqlCall{method: r.Method, path: r.URL.Path, query: r.URL.Query()}
		if b, _ := io.ReadAll(r.Body); len(b) > 0 {
			json.Unmarshal(b, &last.body)
		}
		io.WriteString(w, `{"header":{"resultCode":0,"isSuccessful":true},"jobId":"job-1","dbSecurityGroupId":"sg-1",`+
			`"parameterGroupId":"pg-1","notificationGroupId":"ng-1","userGroupId":"ug-1"}`)
	}))
	defer srv.Close()
	target, _ := url.Parse(srv.URL)

	t.Setenv("HOME", t.TempDir())
	t.Setenv("NHN_REGION", "kr1")
	t.Setenv("NHN_MYSQL_APP_KEY", "ak")
	t.Setenv("NHN_MYSQL_ACCESS_KEY_ID", "id")
	t.Setenv("NHN_MYSQL_SECRET_ACCESS_KEY", "sk")
	oldTr, oldOut, oldQ := http.DefaultTransport, output, query
	defer func() { http.DefaultTransport, output, query = oldTr, oldOut, oldQ }()
	http.DefaultTransport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		r = r.Clone(r.Context())
		r.URL.Scheme, r.URL.Host = target.Scheme, target.Host
		return oldTr.RoundTrip(r)
	})

	var err error
	out := captureStdout(t, func() {
		rootCmd.SetArgs(append([]string{"rds-mysql"}, args...))
		err = rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("%v: %v", args, err)
	}
	return last, out
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestRDSMySQLCommandsResolveNamesAndPrintJSON(t *testing.T) {
	id := fakeMySQLInstanceID
	cases := []struct {
		args       []string
		method     string
		path       string
		wantKey    string
		wantBody   map[string]interface{}
		wantQuery  url.Values
		notInQuery string
	}{
		{args: []string{"force-reboot-db-instance", "--db-instance-identifier", "mydb"},
			method: "POST", path: "/v4.0/db-instances/" + id + "/force-restart", wantKey: "jobId"},
		{args: []string{"create-db-user", "--db-instance-identifier", "mydb", "--db-user-name", "u", "--db-password", "Pw-12345!", "--host", "%", "--authority-type", "READ"},
			method: "POST", path: "/v4.0/db-instances/" + id + "/db-users", wantKey: "jobId"},
		{args: []string{"delete-db-user", "--db-instance-identifier", "mydb", "--db-user-id", "u-1"},
			method: "DELETE", path: "/v4.0/db-instances/" + id + "/db-users/u-1", wantKey: "jobId"},
		{args: []string{"delete-db-schema", "--db-instance-identifier", "mydb", "--db-schema-id", "s-1"},
			method: "DELETE", path: "/v4.0/db-instances/" + id + "/db-schemas/s-1", wantKey: "jobId"},
		{args: []string{"create-db-schema", "--db-instance-identifier", "mydb", "--db-schema-name", "s"},
			method: "POST", path: "/v4.0/db-instances/" + id + "/db-schemas", wantKey: "jobId"},
		{args: []string{"create-read-replica", "--db-instance-identifier", "mydb", "--replica-name", "rr", "--availability-zone", "kr-pub-b"},
			method: "POST", path: "/v4.0/db-instances/" + id + "/replicate", wantKey: "dbInstanceName",
			wantBody: map[string]interface{}{"dbInstanceName": "rr"}},
		{args: []string{"describe-db-instance", "--db-instance-identifier", "mydb"},
			method: "GET", path: "/v4.0/db-instances/" + id, wantKey: "header"},
		{args: []string{"describe-db-instance-network", "--db-instance-identifier", "mydb"},
			method: "GET", path: "/v4.0/db-instances/" + id + "/network-info", wantKey: "header"},
		{args: []string{"describe-db-instance-storage", "--db-instance-identifier", "mydb"},
			method: "GET", path: "/v4.0/db-instances/" + id + "/storage-info", wantKey: "header"},
		{args: []string{"get-metric-statistics", "--db-instance-identifier", "mydb", "--measure-names", "CPU_USAGE,MEMORY_USAGE",
			"--from", "2026-11-18T00:00:00+09:00", "--to", "2026-11-18T01:00:00+09:00"},
			method: "GET", path: "/v4.0/metric-statistics", wantKey: "header",
			wantQuery: url.Values{"dbInstanceId": {id}, "measureNames": {"CPU_USAGE", "MEMORY_USAGE"},
				"from": {"2026-11-18T00:00:00+09:00"}, "to": {"2026-11-18T01:00:00+09:00"}},
			notInQuery: "interval"},
		{args: []string{"enable-multi-az", "--db-instance-identifier", "mydb", "--db-instance-candidate-name", "sb", "--ping-interval", "3"},
			method: "PUT", path: "/v4.0/db-instances/" + id + "/high-availability", wantKey: "jobId",
			wantBody: map[string]interface{}{"useHighAvailability": true, "dbInstanceCandidateName": "sb", "pingInterval": float64(3)}},
		{args: []string{"create-db-snapshot", "--db-instance-identifier", "mydb", "--db-snapshot-identifier", "bk"},
			method: "POST", path: "/v4.0/backups", wantKey: "jobId",
			wantBody: map[string]interface{}{"backupName": "bk", "backupMethodType": "FULL", "dbInstanceId": id}},
		{args: []string{"modify-db-instance", "--db-instance-identifier", "mydb", "--port", "3307"},
			method: "PUT", path: "/v4.0/db-instances/" + id, wantKey: "jobId"},
		{args: []string{"stop-db-instance", "--db-instance-identifier", "mydb"},
			method: "POST", path: "/v4.0/db-instances/" + id + "/stop", wantKey: "jobId"},
		{args: []string{"create-db-security-group", "--db-security-group-name", "sg"},
			method: "POST", path: "/v4.0/db-security-groups", wantKey: "dbSecurityGroupId"},
		{args: []string{"authorize-db-security-group-ingress", "--db-security-group-identifier", "sg-1", "--cidr", "10.0.0.0/24"},
			method: "POST", path: "/v4.0/db-security-groups/sg-1/rules", wantKey: "jobId"},
		{args: []string{"create-db-parameter-group", "--db-parameter-group-name", "pg", "--engine-version", "MYSQL_V8032"},
			method: "POST", path: "/v4.0/parameter-groups", wantKey: "parameterGroupId"},
		{args: []string{"copy-db-parameter-group", "--source-parameter-group-id", "pg-0", "--target-parameter-group-name", "pg"},
			method: "POST", path: "/v4.0/parameter-groups/pg-0/copy", wantKey: "parameterGroupId"},
		{args: []string{"create-notification-group", "--name", "ng"},
			method: "POST", path: "/v4.0/notification-groups", wantKey: "notificationGroupId"},
		{args: []string{"create-user-group", "--name", "ug"},
			method: "POST", path: "/v4.0/user-groups", wantKey: "userGroupId"},
	}
	for _, c := range cases {
		t.Run(c.args[0], func(t *testing.T) {
			call, out := runMySQLFake(t, append(c.args, "-o", "json")...)
			if call.method != c.method || call.path != c.path {
				t.Errorf("request %s %s, want %s %s", call.method, call.path, c.method, c.path)
			}
			var printed map[string]interface{}
			if err := json.Unmarshal([]byte(out), &printed); err != nil {
				t.Fatalf("stdout is not JSON: %q", out)
			}
			if _, ok := printed[c.wantKey]; !ok {
				t.Errorf("stdout has no %s: %s", c.wantKey, out)
			}
			for k, v := range c.wantBody {
				if call.body[k] != v {
					t.Errorf("body %s = %v, want %v", k, call.body[k], v)
				}
			}
			for k, v := range c.wantQuery {
				if strings.Join(call.query[k], ",") != strings.Join(v, ",") {
					t.Errorf("query %s = %v, want %v", k, call.query[k], v)
				}
			}
			if c.notInQuery != "" && call.query.Has(c.notInQuery) {
				t.Errorf("query has %s: %v", c.notInQuery, call.query)
			}
		})
	}
}

// Without -o json the write commands keep their text line.
func TestRDSMySQLWriteCommandTextOutput(t *testing.T) {
	_, out := runMySQLFake(t, "start-db-instance", "--db-instance-identifier", "mydb", "-o", "table")
	if !strings.Contains(out, "Job ID: job-1") {
		t.Errorf("text output %q", out)
	}
}
