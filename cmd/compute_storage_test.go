package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/spf13/pflag"
)

type fakeCall struct{ method, path, body string }

// fakeIaaSServer answers the identity token request with a catalog that
// points every IaaS service at itself (under /<type>) and every other call
// with reply(method, path).
func fakeIaaSServer(t *testing.T, reply func(method, path string) string) (*httptest.Server, func() []fakeCall) {
	t.Helper()
	var mu sync.Mutex
	var calls []fakeCall
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		if r.URL.Path == "/v2.0/tokens" {
			var catalog []map[string]interface{}
			for _, typ := range []string{"compute", "volumev2", "network", "object-store"} {
				catalog = append(catalog, map[string]interface{}{"type": typ,
					"endpoints": []map[string]string{{"publicURL": srv.URL + "/" + typ, "region": "kr1"}}})
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"access": map[string]interface{}{
				"token":          map[string]interface{}{"id": "tok", "expires": time.Now().Add(time.Hour)},
				"serviceCatalog": catalog}})
			return
		}
		mu.Lock()
		calls = append(calls, fakeCall{r.Method, r.URL.Path, string(b)})
		mu.Unlock()
		body := reply(r.Method, r.URL.Path)
		switch {
		case r.Method == "DELETE":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/object-store/"):
			w.WriteHeader(http.StatusCreated)
		}
		io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv, func() []fakeCall {
		mu.Lock()
		defer mu.Unlock()
		return append([]fakeCall(nil), calls...)
	}
}

// runIaaS runs the CLI in-process against a fake IaaS server and returns the
// API calls it made and what it printed on stdout.
func runIaaS(t *testing.T, reply func(method, path string) string, args ...string) ([]fakeCall, string) {
	t.Helper()
	srv, calls := fakeIaaSServer(t, reply)
	target, _ := url.Parse(srv.URL)
	idURL, _ := url.Parse(defaultIdentityURL)
	t.Setenv("NHN_CLOUD_ENDPOINT_OVERRIDES", "")
	oldTr, oldCfg := http.DefaultTransport, loadedConfig
	oldU, oldP, oldT, oldOut, oldQ := username, password, tenantID, output, query
	http.DefaultTransport = &hostRewriteTransport{inner: oldTr, hosts: map[string]*url.URL{idURL.Host: target}}
	loadedConfig = &Config{Username: "u", APIPassword: "p", TenantID: "t", Region: "kr1", Extra: map[string]string{}}
	username, password, tenantID = "u", "p", "t"
	c, _, _ := rootCmd.Find(args)
	defer func() {
		http.DefaultTransport, loadedConfig = oldTr, oldCfg
		username, password, tenantID, output, query = oldU, oldP, oldT, oldOut, oldQ
		c.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				f.Value.Set(f.DefValue)
				f.Changed = false
			}
		})
	}()
	rootCmd.SetArgs(args)
	out := captureStdout(t, func() {
		if err := rootCmd.Execute(); err != nil {
			t.Errorf("%v: %v", args, err)
		}
	})
	return calls(), out
}

const volumeJSON = `{"volume":{"id":"v1","name":"n","status":"available","size":20,"attachments":[{"id":"v1","volume_id":"v1","server_id":"s1","device":"/dev/vdb"}]}}`

func TestAttachVolumeUsesNovaAttachments(t *testing.T) {
	calls, out := runIaaS(t, func(string, string) string {
		return `{"volumeAttachment":{"id":"v1","device":"/dev/vdb","serverId":"s1","volumeId":"v1"}}`
	}, "block-storage", "attach-volume", "--volume-id", "v1", "--server-id", "s1", "-o", "json")
	if len(calls) != 1 || calls[0].method != "POST" || calls[0].path != "/compute/servers/s1/os-volume_attachments" ||
		!strings.Contains(calls[0].body, `"volumeId":"v1"`) {
		t.Errorf("calls: %+v", calls)
	}
	if !strings.Contains(out, `"device": "/dev/vdb"`) {
		t.Errorf("-o json printed %q", out)
	}
}

func TestDetachVolumeFindsServerAndUsesNova(t *testing.T) {
	calls, _ := runIaaS(t, func(method, path string) string {
		if method == "GET" {
			return volumeJSON
		}
		return ""
	}, "block-storage", "detach-volume", "--volume-id", "v1")
	if len(calls) != 2 || calls[1].method != "DELETE" || calls[1].path != "/compute/servers/s1/os-volume_attachments/v1" {
		t.Errorf("calls: %+v", calls)
	}
}

// Write commands print JSON for -o json, not their text line.
func TestWriteCommandsPrintJSON(t *testing.T) {
	reply := func(method, path string) string {
		if method == "PUT" || method == "GET" {
			return volumeJSON
		}
		return ""
	}
	for _, args := range [][]string{
		{"block-storage", "update-volume", "--volume-id", "v1", "--name", "n"},
		{"block-storage", "extend-volume", "--volume-id", "v1", "--size", "30"},
		{"block-storage", "delete-volume", "--volume-id", "v1"},
		{"block-storage", "delete-snapshot", "--snapshot-id", "sn1"},
		{"block-storage", "detach-volume", "--volume-id", "v1", "--server-id", "s1"},
		{"compute", "start-instances", "--instance-id", "s1"},
		{"compute", "stop-instances", "--instance-id", "s1"},
		{"compute", "reboot-instances", "--instance-id", "s1"},
		{"compute", "delete-instance", "--instance-id", "s1"},
		{"compute", "delete-key-pair", "--key-name", "k1"},
	} {
		_, out := runIaaS(t, reply, append(args, "-o", "json")...)
		var v map[string]interface{}
		if err := json.Unmarshal([]byte(out), &v); err != nil {
			t.Errorf("%v printed %q", args[:2], out)
		}
	}
	_, out := runIaaS(t, reply, "block-storage", "update-volume", "--volume-id", "v1", "--name", "n", "-o", "json")
	if !strings.Contains(out, `"volume"`) {
		t.Errorf("update-volume dropped the API response: %q", out)
	}
}

func TestCreateInstanceWaitPollsUntilActive(t *testing.T) {
	old := instanceWaitInterval
	instanceWaitInterval = time.Millisecond
	defer func() { instanceWaitInterval = old }()
	gets := 0
	calls, out := runIaaS(t, func(method, path string) string {
		switch {
		case strings.HasPrefix(path, "/network/"):
			return `{"vpcsubnet":{"id":"sn","vpc_id":"vpc"}}`
		case method == "POST":
			return `{"server":{"id":"s1"}}`
		}
		gets++
		status := "BUILD"
		if gets > 2 {
			status = "ACTIVE"
		}
		return fmt.Sprintf(`{"server":{"id":"s1","name":"vm","status":%q}}`, status)
	}, "compute", "create-instance", "--name", "vm", "--image-id", "i", "--flavor-id", "f", "--subnet-id", "sn", "--wait", "-o", "json")
	if gets != 3 {
		t.Errorf("server polled %d times; calls %+v", gets, calls)
	}
	if !strings.Contains(out, `"status": "ACTIVE"`) {
		t.Errorf("output %q", out)
	}
}

// The SSH login user has its own flag, so the global --username (the API
// user) still applies to connect.
// Existing scripts call connect with --username (the SSH user) and
// --allow-cidr; both must still parse. -l is --ssh-user.
func TestConnectKeepsOldFlags(t *testing.T) {
	fs := computeConnectCmd.Flags()
	if f := fs.ShorthandLookup("l"); f == nil || f.Name != "ssh-user" {
		t.Errorf("-l is %v, want --ssh-user", f)
	}
	if err := fs.Parse([]string{"--instance-id", "x", "--username", "ubuntu", "--allow-cidr", "10.0.0.0/8"}); err != nil {
		t.Fatalf("old flags rejected: %v", err)
	}
	if u, _ := fs.GetString("username"); u != "ubuntu" {
		t.Errorf("--username = %q", u)
	}
}

// Without a floating IP, connect stops with guidance and changes nothing:
// no floating IP is taken or allocated and no security group is created.
func TestConnectWithoutFloatingIPMakesNoChanges(t *testing.T) {
	if os.Getenv("CONNECT_HELPER") == "1" {
		rootCmd.SetArgs([]string{"compute", "connect", "--instance-id", "s1"})
		rootCmd.Execute()
		return
	}
	srv, calls := fakeIaaSServer(t, func(method, path string) string {
		switch {
		case strings.Contains(path, "/servers/"):
			return `{"server":{"id":"s1","name":"vm","status":"ACTIVE","addresses":{"net":[{"addr":"10.0.0.5","OS-EXT-IPS:type":"fixed"}]}}}`
		case strings.Contains(path, "/ports"):
			return `{"ports":[{"id":"p1","security_groups":[]}]}`
		case strings.Contains(path, "/floatingips"):
			return `{"floatingips":[{"id":"f1","status":"DOWN","floating_ip_address":"203.0.113.9"}]}`
		}
		return `{}`
	})
	idURL, _ := url.Parse(defaultIdentityURL)
	c := exec.Command(os.Args[0], "-test.run=^TestConnectWithoutFloatingIPMakesNoChanges$")
	c.Env = append(os.Environ(), "CONNECT_HELPER=1", "HOME="+t.TempDir(), "NHN_CLOUD_PROFILE=",
		"NHN_CLOUD_ENDPOINT_OVERRIDES="+idURL.Host+"="+srv.URL,
		"NHN_CLOUD_USERNAME=u", "NHN_CLOUD_PASSWORD=p", "NHN_CLOUD_TENANT_ID=t", "NHN_CLOUD_REGION=kr1")
	out, err := c.CombinedOutput()
	if err == nil {
		t.Errorf("connect succeeded without a floating IP: %s", out)
	}
	if !strings.Contains(string(out), "floating IP") {
		t.Errorf("no guidance: %s", out)
	}
	for _, call := range calls() {
		if call.method != "GET" {
			t.Errorf("connect changed something: %s %s", call.method, call.path)
		}
	}
}

func TestObjectStorageJSONOutput(t *testing.T) {
	reply := func(method, path string) string {
		if method == "GET" && path == "/object-store" {
			return `[{"name":"c1","count":1,"bytes":5}]`
		}
		return ""
	}
	_, out := runIaaS(t, reply, "object-storage", "ls", "-o", "json")
	var list []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &list); err != nil || len(list) != 1 || list[0]["name"] != "c1" {
		t.Errorf("ls -o json printed %q", out)
	}
	for _, args := range [][]string{
		{"object-storage", "info", "obs://c1/o.txt", "--json"},
		{"object-storage", "info", "obs://c1/o.txt", "-o", "json"},
		{"object-storage", "info", "obs://c1", "-o", "json"},
		{"object-storage", "mb", "obs://c2", "-o", "json"},
		{"object-storage", "rm", "obs://c1/o.txt", "-o", "json"},
	} {
		_, out := runIaaS(t, reply, args...)
		var v map[string]interface{}
		if err := json.Unmarshal([]byte(out), &v); err != nil {
			t.Errorf("%v printed %q", args, out)
		}
	}
}
