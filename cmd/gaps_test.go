package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestParseEndpointOverrides(t *testing.T) {
	m, err := parseEndpointOverrides(" API.example.com = other.example.net , b.example.com=http://127.0.0.1:8080,")
	if err != nil {
		t.Fatal(err)
	}
	if got := m["api.example.com"].String(); got != "https://other.example.net" {
		t.Errorf("host-only target: got %q", got)
	}
	if got := m["b.example.com"].String(); got != "http://127.0.0.1:8080" {
		t.Errorf("scheme target: got %q", got)
	}
	for _, bad := range []string{"noequals", "a=", "=b", "a=https://h/path"} {
		if _, err := parseEndpointOverrides(bad); err == nil {
			t.Errorf("%q: want error", bad)
		}
	}
}

func TestHostRewriteTransport(t *testing.T) {
	var gotPath, gotHost string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotHost = r.URL.RequestURI(), r.Host
	}))
	defer srv.Close()
	target, _ := url.Parse(srv.URL)
	tr := &hostRewriteTransport{inner: http.DefaultTransport, hosts: map[string]*url.URL{"api.example.com": target}}

	req, _ := http.NewRequest(http.MethodGet, "https://API.example.com/v2.0/vpcs?x=1", nil)
	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if gotPath != "/v2.0/vpcs?x=1" || gotHost != target.Host {
		t.Errorf("rewritten request: path %q host %q", gotPath, gotHost)
	}
	if req.URL.Host != "API.example.com" {
		t.Errorf("original request was modified: %q", req.URL.Host)
	}
}

func TestDoJSONErrorKeepsBodyAndDropsQuery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Auth-Token") != "tok" {
			t.Errorf("token header missing")
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"a":1`) {
			t.Errorf("body not sent: %s", body)
		}
		w.WriteHeader(409)
		io.WriteString(w, `{"message":"OverQuota"}`)
	}))
	defer srv.Close()
	err := doJSON(t.Context(), http.MethodPost, srv.URL+"/v2.0/vpcs?secret=1", "tok", map[string]int{"a": 1}, nil)
	ae, ok := err.(*apiError)
	if !ok || ae.Status != 409 || !strings.Contains(ae.Error(), "OverQuota") || strings.Contains(ae.Error(), "secret") {
		t.Fatalf("got %v", err)
	}
}

func TestRoutingTableDecodingAndFilter(t *testing.T) {
	// The list API (detail=true) returns objects; the single API returns IDs.
	var list struct {
		RoutingTables []routingTable `json:"routingtables"`
	}
	data := `{"routingtables":[
		{"id":"rt1","vpcs":[{"id":"vpc-a","name":"a"}],"subnets":[{"id":"s1","name":"x"}]},
		{"id":"rt2","vpcs":[{"id":"vpc-b","name":"b"}]}]}`
	if err := json.Unmarshal([]byte(data), &list); err != nil {
		t.Fatal(err)
	}
	got := filterRoutingTablesByVPC(list.RoutingTables, "vpc-b")
	if len(got) != 1 || got[0].ID != "rt2" {
		t.Errorf("filter: %+v", got)
	}
	if len(filterRoutingTablesByVPC(list.RoutingTables, "")) != 2 {
		t.Error("empty filter must keep all tables")
	}

	var single struct {
		RoutingTable routingTable `json:"routingtable"`
	}
	if err := json.Unmarshal([]byte(`{"routingtable":{"id":"rt1","vpcs":["vpc-a"],"subnets":["s1","s2"]}}`), &single); err != nil {
		t.Fatal(err)
	}
	if single.RoutingTable.vpcID() != "vpc-a" || len(single.RoutingTable.Subnets) != 2 {
		t.Errorf("single: %+v", single.RoutingTable)
	}
}

func TestBuildQuotaRows(t *testing.T) {
	rows := buildQuotaRows(quotaInputs{
		tenantID: "me",
		vpcs: []map[string]interface{}{
			{"id": "v1", "name": "one", "tenant_id": "me"},
			{"id": "v2", "name": "two", "tenant_id": "me"},
			{"id": "pub", "tenant_id": "other", "router:external": true},
			{"id": "shared", "tenant_id": "other"},
		},
		subnets:         []map[string]interface{}{{"vpc_id": "v1"}, {"vpc_id": "v1"}, {"vpc_id": "pub"}},
		routingTables:   []routingTable{{VPCs: []idName{{ID: "v2"}}}},
		projectCounts:   map[string]int{"floating IP": 7},
		projectFailures: map[string]string{"NAT gateway": "list failed"},
	})
	byKey := map[string]quotaRow{}
	for _, r := range rows {
		byKey[r.Resource+"|"+r.Scope] = r
	}
	check := func(key, used string, limit int) {
		t.Helper()
		r, ok := byKey[key]
		if !ok || r.Used != used || r.Limit != limit {
			t.Errorf("%s: got %+v", key, r)
		}
	}
	check("VPC|project", "2", 3)
	check("subnet|VPC one (v1)", "2", 10)
	check("routing table|VPC two (v2)", "1", 10)
	check("floating IP|project", "7", 50)
	check("NAT gateway|project", "unknown (list failed)", 3)
	if _, ok := byKey["subnet|VPC  (pub)"]; ok {
		t.Error("the shared public network must not get per-VPC rows")
	}
}

func TestListenerOptionsOnlyWhenSet(t *testing.T) {
	c := &cobra.Command{}
	addListenerOptFlags(c, listenerExtraOpts)
	if got := changedListenerFields(c, listenerExtraOpts); len(got) != 0 {
		t.Fatalf("no flags set must keep the SDK path, got %v", got)
	}
	c.Flags().Set("keepalive-timeout", "3600")
	c.Flags().Set("proxy-protocol", "true")
	c.Flags().Set("keepalive-enable", "false")
	got := changedListenerFields(c, listenerExtraOpts)
	if got["keepalive_timeout"] != 3600 || got["proxy_protocol"] != true || got["keepalive_enable"] != false || len(got) != 3 {
		t.Errorf("got %v", got)
	}
}

func TestKeyTypeHintAndCSV(t *testing.T) {
	if keyTypeHint("ssh-ed25519 AAAA x") == "" || keyTypeHint("ssh-rsa AAAA x") != "" {
		t.Error("keyTypeHint")
	}
	if got := strings.Join(splitCSV(" a, ,b,"), "|"); got != "a|b" {
		t.Errorf("splitCSV: %q", got)
	}
}

// requiredIDFlagsWithoutHint lists "command --flag" for required *-id flags
// whose help does not say how to find the value.
func requiredIDFlagsWithoutHint(root *cobra.Command) []string {
	var out []string
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		c.LocalFlags().VisitAll(func(f *pflag.Flag) {
			if !strings.HasSuffix(f.Name, "-id") {
				return
			}
			if _, req := f.Annotations[cobra.BashCompOneRequiredFlag]; !req {
				return
			}
			if !strings.Contains(f.Usage, "find it with") && !strings.Contains(f.Usage, "find them with") {
				out = append(out, c.CommandPath()+" --"+f.Name)
			}
		})
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(root)
	sort.Strings(out)
	return out
}

// Commands on the network setup path must say where each required ID comes from.
func TestRequiredIDFlagsHaveLookupHint(t *testing.T) {
	covered := []string{
		"nhncloud network ", "nhncloud internet-gateway ", "nhncloud service-gateway ",
		"nhncloud loadbalancer ", "nhncloud nat-gateway ", "nhncloud compute ",
	}
	for _, miss := range requiredIDFlagsWithoutHint(rootCmd) {
		for _, prefix := range covered {
			if strings.HasPrefix(miss, prefix) {
				t.Errorf("required flag without lookup hint: %s", miss)
			}
		}
	}
}

// Run with -run TestHintInventory -v to list the remaining commands.
func TestHintInventory(t *testing.T) {
	for _, m := range requiredIDFlagsWithoutHint(rootCmd) {
		t.Log(m)
	}
}

// A failed VPC, subnet or routing table list must show "unknown", not a count.
func TestBuildQuotaRowsMarksFailedPrerequisiteLists(t *testing.T) {
	rows := buildQuotaRows(quotaInputs{
		tenantID:        "me",
		vpcs:            []map[string]interface{}{{"id": "v1", "name": "one", "tenant_id": "me"}},
		projectCounts:   map[string]int{},
		projectFailures: map[string]string{"subnet": "list failed", "routing table": "list failed"},
	})
	for _, r := range rows {
		if (r.Resource == "subnet" || r.Resource == "routing table") && !strings.HasPrefix(r.Used, "unknown") {
			t.Errorf("%s row: got %q, want unknown", r.Resource, r.Used)
		}
	}
	rows = buildQuotaRows(quotaInputs{projectCounts: map[string]int{}, projectFailures: map[string]string{"VPC": "list failed"}})
	if rows[0].Resource != "VPC" || !strings.HasPrefix(rows[0].Used, "unknown") {
		t.Errorf("VPC row: %+v", rows[0])
	}
}

// Transport errors repeat the URL; the query must not appear in them.
func TestDoJSONTransportErrorDropsQuery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	addr := srv.URL
	srv.Close() // connection refused
	err := doJSON(t.Context(), http.MethodGet, addr+"/v2.0/ports?device_id=SECRETVALUE", "", nil, nil)
	if err == nil || strings.Contains(err.Error(), "SECRETVALUE") {
		t.Fatalf("got %v", err)
	}
}

// A source written with a scheme or default port must still match; a path is
// rejected instead of silently never matching.
func TestEndpointOverrideSourceForms(t *testing.T) {
	hit := false
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { hit = true }))
	defer srv.Close()
	m, err := parseEndpointOverrides("https://API.example.invalid:443/=" + srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	tr := &hostRewriteTransport{inner: http.DefaultTransport, hosts: m}
	for _, u := range []string{"https://api.example.invalid/x", "https://api.example.invalid:443/x"} {
		hit = false
		req, _ := http.NewRequest(http.MethodGet, u, nil)
		if resp, err := tr.RoundTrip(req); err == nil {
			resp.Body.Close()
		}
		if !hit {
			t.Errorf("%s was not sent to the override target", u)
		}
	}
	if _, err := parseEndpointOverrides("api.example.invalid/v2=staging.example.net"); err == nil {
		t.Error("source with a path must be rejected")
	}
}

func TestKeepRawRoutingTables(t *testing.T) {
	raw := map[string]interface{}{"routingtables": []interface{}{
		map[string]interface{}{"id": "rt1", "extra": "kept"},
		map[string]interface{}{"id": "rt2"},
	}}
	out, _ := json.Marshal(keepRawRoutingTables(raw, []routingTable{{ID: "rt1"}}))
	if string(out) != `{"routingtables":[{"extra":"kept","id":"rt1"}]}` {
		t.Errorf("got %s", out)
	}
	out, _ = json.Marshal(keepRawRoutingTables(raw, nil))
	if string(out) != `{"routingtables":[]}` {
		t.Errorf("no match: got %s", out)
	}
}

func TestEmitNilTableFallsBackToJSON(t *testing.T) {
	old := output
	defer func() { output = old }()
	output = "text"
	emit(map[string]string{"a": "b"}, nil) // must not panic
}

// Every error path of doJSON drops the query, including URL parse errors.
func TestDoJSONParseErrorDropsQuery(t *testing.T) {
	err := doJSON(t.Context(), http.MethodGet, "https://h.invalid/v2.0/%ZZ?token=SECRETVALUE", "", nil, nil)
	if err == nil || strings.Contains(err.Error(), "SECRETVALUE") {
		t.Fatalf("got %v", err)
	}
}
