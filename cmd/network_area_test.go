package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type networkCall struct {
	Method, Path string
	Body         map[string]map[string]interface{}
}

// runNetworkCmd runs args against one fake server that answers the identity
// token request (with a catalog whose network endpoint is the server itself)
// and replies to every other request with replies[method+" "+path].
func runNetworkCmd(t *testing.T, replies map[string]string, args ...string) (string, []networkCall) {
	t.Helper()
	var calls []networkCall
	var srvURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2.0/tokens" {
			fmt.Fprintf(w, `{"access":{"token":{"id":"tok","expires":%q,"tenant":{"id":"t"}},
				"serviceCatalog":[{"type":"network","endpoints":[{"publicURL":%q,"region":"kr1"}]}]}}`,
				time.Now().Add(time.Hour).Format(time.RFC3339), srvURL)
			return
		}
		c := networkCall{Method: r.Method, Path: r.URL.Path}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &c.Body)
		calls = append(calls, c)
		reply, ok := replies[r.Method+" "+r.URL.Path]
		if !ok {
			w.WriteHeader(404)
			io.WriteString(w, `{"NeutronError":{"message":"not in fake"}}`)
			return
		}
		io.WriteString(w, reply)
	}))
	defer srv.Close()
	srvURL = srv.URL
	target, _ := url.Parse(srv.URL)
	idHost, _ := url.Parse(defaultIdentityURL)

	oldTr, oldT, oldU, oldP, oldR, oldOut, oldQ := http.DefaultTransport, tenantID, username, password, region, output, query
	defer func() {
		http.DefaultTransport, tenantID, username, password, region, output, query = oldTr, oldT, oldU, oldP, oldR, oldOut, oldQ
		cachedIaaSAPI, cachedIaaSErr = nil, nil
	}()
	http.DefaultTransport = &hostRewriteTransport{inner: oldTr, hosts: map[string]*url.URL{idHost.Host: target}}
	tenantID, username, password = "t", "u", "p"
	cachedIaaSAPI, cachedIaaSErr = nil, nil

	rootCmd.SetArgs(append(args, "--region", "kr1"))
	var err error
	out := captureStdout(t, func() { err = rootCmd.Execute() })
	if err != nil {
		t.Fatalf("%v: %v", args, err)
	}
	return out, calls
}

func decodeJSON(t *testing.T, out string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	return m
}

// Write commands print the API response for -o json, not the text line.
func TestNetworkWriteCommandsFollowOutputFlag(t *testing.T) {
	cases := []struct {
		args    []string
		key     string // reply method+path
		reply   string
		wrapper string
	}{
		{[]string{"network", "disassociate-floating-ip", "--floating-ip-id", "f1"},
			"PUT /v2.0/floatingips/f1", `{"floatingip":{"id":"f1","port_id":null}}`, "floatingip"},
		{[]string{"network-acl", "update-network-acl", "--acl-id", "a1", "--name", "n"},
			"PUT /v2.0/acls/a1", `{"acl":{"id":"a1","name":"n"}}`, "acl"},
		{[]string{"network-acl", "update-rule", "--rule-id", "r1", "--description", "d"},
			"PUT /v2.0/acl_rules/r1", `{"acl_rule":{"id":"r1","description":"d"}}`, "acl_rule"},
		{[]string{"transit-hub", "update-transit-hub", "--hub-id", "h1", "--name", "n"},
			"PUT /v2.0/gateways/transithubs/h1", `{"transithub":{"id":"h1","name":"n"}}`, "transithub"},
		{[]string{"transit-hub", "update-routing-table", "--routing-table-id", "t1", "--name", "n"},
			"PUT /v2.0/gateways/transithub_routing_tables/t1", `{"transithub_routing_table":{"id":"t1"}}`, "transithub_routing_table"},
		{[]string{"transit-hub", "update-attachment", "--attachment-id", "a1", "--name", "n"},
			"PUT /v2.0/gateways/transithub_attachments/a1", `{"transithub_attachment":{"id":"a1"}}`, "transithub_attachment"},
		{[]string{"transit-hub", "update-multicast-domain", "--domain-id", "d1", "--name", "n"},
			"PUT /v2.0/gateways/transithub_multicast_domains/d1", `{"transithub_multicast_domain":{"id":"d1"}}`, "transithub_multicast_domain"},
		{[]string{"transit-hub", "create-association", "--routing-table-id", "t1", "--attachment-id", "a1"},
			"POST /v2.0/gateways/transithub_routing_associations", `{"transithub_routing_association":{"id":"as1"}}`, "transithub_routing_association"},
		{[]string{"transit-hub", "create-propagation", "--routing-table-id", "t1", "--attachment-id", "a1"},
			"POST /v2.0/gateways/transithub_routing_propagations", `{"transithub_routing_propagation":{"id":"pr1"}}`, "transithub_routing_propagation"},
		{[]string{"transit-hub", "create-rule", "--routing-table-id", "t1", "--cidr", "1.1.1.0/24", "--attachment-id", "a1"},
			"POST /v2.0/gateways/transithub_routing_rules", `{"transithub_routing_rule":{"id":"ru1"}}`, "transithub_routing_rule"},
	}
	for _, c := range cases {
		out, _ := runNetworkCmd(t, map[string]string{c.key: c.reply}, append(c.args, "-o", "json")...)
		if _, ok := decodeJSON(t, out)[c.wrapper]; !ok {
			t.Errorf("%v -o json: no %q in %s", c.args, c.wrapper, out)
		}
	}
}

// describe-load-balancers -o json prints the API object, with the fields the
// SDK struct does not have.
func TestDescribeLoadBalancersPrintsAPIObject(t *testing.T) {
	out, _ := runNetworkCmd(t, map[string]string{
		"GET /v2.0/lbaas/loadbalancers/lb1": `{"loadbalancer":{"id":"lb1","loadbalancer_type":"shared","engine_version":"v2"}}`,
	}, "loadbalancer", "describe-load-balancers", "--lb-id", "lb1", "-o", "json")
	lb, _ := decodeJSON(t, out)["loadbalancer"].(map[string]interface{})
	if lb["loadbalancer_type"] != "shared" || lb["engine_version"] != "v2" {
		t.Errorf("output: %s", out)
	}
}

// The hand-written stub must not hide the catalog op of the same name.
func TestGetLoadBalancerIsTheCatalogOp(t *testing.T) {
	ensureCatalogCommands()
	c, _, err := rootCmd.Find([]string{"loadbalancer", "get-load-balancer"})
	if err != nil || c.Flags().Lookup("loadbalancer-id") == nil {
		t.Errorf("get-load-balancer is not the catalog op (err %v)", err)
	}
}

func TestCreateTransitHubSendsMulticastFlag(t *testing.T) {
	_, calls := runNetworkCmd(t, map[string]string{
		"POST /v2.0/gateways/transithubs": `{"transithub":{"id":"h1"}}`,
	}, "transit-hub", "create-transit-hub", "--name", "h", "--multicast-enable", "-o", "json")
	hub := calls[0].Body["transithub"]
	if hub["multicast_enable"] != true || hub["default_association_enable"] != true || hub["default_propagation_enable"] != true {
		t.Errorf("request: %v", hub)
	}
}

func TestCreateAttachmentSendsSubnet(t *testing.T) {
	_, calls := runNetworkCmd(t, map[string]string{
		"POST /v2.0/gateways/transithub_attachments": `{"transithub_attachment":{"id":"a1"}}`,
	}, "transit-hub", "create-attachment", "--name", "a", "--transit-hub-id", "h1",
		"--resource-id", "vpc1", "--subnet-id", "sub1", "-o", "json")
	a := calls[0].Body["transithub_attachment"]
	if a["subnet_id"] != "sub1" || a["resource_id"] != "vpc1" {
		t.Errorf("request: %v", a)
	}
}

// The routing rule API has create, get and delete only.
func TestTransitHubHasNoUpdateRule(t *testing.T) {
	c, _, err := rootCmd.Find([]string{"transit-hub", "update-rule"})
	if err == nil && c.Name() == "update-rule" {
		t.Error("transit-hub update-rule calls an API that does not exist")
	}
}

func TestDescribeServiceGatewayShowsFixedIP(t *testing.T) {
	out, _ := runNetworkCmd(t, map[string]string{
		"GET /v2.0/gateways/servicegateways/g1": `{"servicegateway":{"id":"g1","fixed_ip":"192.168.0.82","port_id":"p1"}}`,
	}, "service-gateway", "describe-service-gateway", "--gateway-id", "g1", "-o", "table")
	if !strings.Contains(out, "192.168.0.82") || !strings.Contains(out, "p1") {
		t.Errorf("output: %s", out)
	}
}
