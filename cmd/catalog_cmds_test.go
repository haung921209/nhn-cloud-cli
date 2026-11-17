package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/haung921209/nhn-cloud-cli/internal/catalog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Every embedded catalog must pass validation; this is the gate for adding
// catalogs.
func TestEmbeddedCatalogsAreValid(t *testing.T) {
	services, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(services) == 0 {
		t.Fatal("no catalogs embedded")
	}
	for _, e := range catalog.ValidateAll(services) {
		t.Error(e)
	}
}

type seenRequest struct {
	method, path, query, contentType string
	header                           http.Header
	body                             []byte
}

// runFake builds the commands of svc against a fake server and runs args.
func runFake(t *testing.T, svc catalog.Service, reply string, args ...string) (seenRequest, string, error) {
	t.Helper()
	var seen seenRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = seenRequest{r.Method, r.URL.EscapedPath(), r.URL.RawQuery, r.Header.Get("Content-Type"), r.Header.Clone(), nil}
		seen.body, _ = io.ReadAll(r.Body)
		io.WriteString(w, reply)
	}))
	defer srv.Close()
	svc.BaseURL = srv.URL

	root := &cobra.Command{Use: "nhncloud"}
	addCatalogCommands(root, []catalog.Service{svc})
	var out strings.Builder
	root.SetOut(&out)
	root.SetArgs(args)
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := root.Execute()
	w.Close()
	os.Stdout = oldStdout
	printed, _ := io.ReadAll(r)
	return seen, string(printed), err
}

func fakeService() catalog.Service {
	return catalog.Service{
		Name:  "fake-svc",
		Title: "Test > Fake",
		Docs:  "https://docs.nhncloud.com/ko/test/",
		Auth:  catalog.Auth{Headers: map[string]string{"X-Secret-Key": "{secretKey}"}},
		Operations: []catalog.Operation{
			{
				Name: "create-thing", Summary: "Create", Method: "POST",
				Path: "/v1/appkeys/{appKey}/things/{thingName}",
				Params: []catalog.Param{
					{Name: "appKey", In: "path", Required: true, Credential: "appKey"},
					{Name: "thingName", In: "path", Required: true},
					{Name: "dryRun", In: "query", Type: "boolean"},
				},
				Body: &catalog.Body{Wrapper: "thing", Fields: []catalog.Field{
					{Name: "size", Type: "integer", Required: true},
					{Name: "tags", Type: "array", Items: "string"},
					{Name: "spec", Type: "object"},
				}},
			},
			{
				Name: "upload-file", Summary: "Upload", Method: "PUT", Path: "/v1/files",
				Body: &catalog.Body{Encoding: "multipart", Fields: []catalog.Field{
					{Name: "file", Type: "file", Required: true},
					{Name: "note", Type: "string"},
				}},
			},
			{Name: "get-thing", Summary: "Get", Method: "GET", Path: "/v1/thing"},
		},
	}
}

func TestCatalogCommandBuildsRequest(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "ak-from-env")
	old := loadedConfig
	loadedConfig = &Config{Extra: map[string]string{"fake_svc_secret_key": "sk-from-profile"}}
	defer func() { loadedConfig = old }()

	seen, out, err := runFake(t, fakeService(), `{"ok":true}`,
		"fake-svc", "create-thing", "--thing-name", "a b/c", "--dry-run", "--size", "3",
		"--tags", "x,y", "--body", `{"spec":{"k":"v"},"size":1}`)
	if err != nil {
		t.Fatal(err)
	}
	if seen.method != "POST" || seen.path != "/v1/appkeys/ak-from-env/things/a%20b%2Fc" || seen.query != "dryRun=true" {
		t.Errorf("request line: %s %s ? %s", seen.method, seen.path, seen.query)
	}
	if seen.header.Get("X-Secret-Key") != "sk-from-profile" {
		t.Errorf("secret key header: %q", seen.header.Get("X-Secret-Key"))
	}
	var body map[string]map[string]interface{}
	if err := json.Unmarshal(seen.body, &body); err != nil {
		t.Fatal(err)
	}
	thing := body["thing"]
	if thing["size"] != float64(3) || thing["spec"] == nil || len(thing["tags"].([]interface{})) != 2 {
		t.Errorf("body: %s", seen.body)
	}
	if !strings.Contains(out, `"ok": true`) {
		t.Errorf("output: %q", out)
	}
}

func TestCatalogCommandRequiresBodyFields(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "ak")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	_, _, err := runFake(t, fakeService(), `{}`, "fake-svc", "create-thing", "--thing-name", "a")
	if err == nil || !strings.Contains(err.Error(), "--size is required") {
		t.Errorf("got %v", err)
	}
}

func TestCatalogCommandMultipart(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	f := filepath.Join(t.TempDir(), "cert.pem")
	os.WriteFile(f, []byte("PEMDATA"), 0600)
	seen, _, err := runFake(t, fakeService(), `{}`, "fake-svc", "upload-file", "--file", f, "--note", "hi")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(seen.contentType, "multipart/form-data") || !strings.Contains(string(seen.body), "PEMDATA") ||
		!strings.Contains(string(seen.body), `filename="cert.pem"`) {
		t.Errorf("multipart: %s %q", seen.contentType, seen.body)
	}
}

// 200 with {"header":{"isSuccessful":false}} is a failure.
func TestCatalogCommandHeaderFailure(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	_, _, err := runFake(t, fakeService(), `{"header":{"isSuccessful":false,"resultCode":4003,"resultMessage":"Invalid appKey"}}`,
		"fake-svc", "get-thing")
	if err == nil || !strings.Contains(err.Error(), "4003") {
		t.Errorf("got %v", err)
	}
}

func TestCatalogCredentialErrorsNameTheSources(t *testing.T) {
	old := loadedConfig
	loadedConfig = &Config{Extra: map[string]string{}}
	defer func() { loadedConfig = old }()
	_, err := resolveCredential(t.Context(), "fake-svc", "secretKey")
	if err == nil || !strings.Contains(err.Error(), "NHN_CLOUD_FAKE_SVC_SECRET_KEY") || !strings.Contains(err.Error(), "fake_svc_secret_key") {
		t.Errorf("got %v", err)
	}
}

// Hand-written commands win over catalog operations with the same name.
func TestCatalogDoesNotReplaceHandWrittenCommands(t *testing.T) {
	root := &cobra.Command{Use: "nhncloud"}
	parent := &cobra.Command{Use: "fake-svc"}
	hand := &cobra.Command{Use: "get-thing", Short: "hand-written", Run: func(*cobra.Command, []string) {}}
	parent.AddCommand(hand)
	root.AddCommand(parent)
	addCatalogCommands(root, []catalog.Service{fakeService()})
	if c := findChild(parent, "get-thing"); c != hand {
		t.Error("catalog replaced a hand-written command")
	}
	if findChild(parent, "create-thing") == nil {
		t.Error("catalog operation not added")
	}
}

func TestCatalogEndpointOverrideHeaderFlagAndEmptyBody(t *testing.T) {
	var gotHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Object-Meta-Color")
		w.Header().Set("X-Account-Object-Count", "7")
		w.Header().Set("Set-Cookie", "session=SECRET")
	}))
	defer srv.Close()
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	t.Setenv("NHN_CLOUD_FAKE_SVC_ENDPOINT", srv.URL)

	svc := fakeService()
	svc.BaseURL = "{endpoint}"
	root := &cobra.Command{Use: "nhncloud"}
	addCatalogCommands(root, []catalog.Service{svc})
	root.SetArgs([]string{"fake-svc", "get-thing", "--header", "X-Object-Meta-Color: blue"})
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := root.Execute()
	w.Close()
	os.Stdout = oldStdout
	out, _ := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if gotHeader != "blue" {
		t.Errorf("--header not sent: %q", gotHeader)
	}
	if !strings.Contains(string(out), "X-Account-Object-Count") || strings.Contains(string(out), "SECRET") {
		t.Errorf("header output: %s", out)
	}

	os.Unsetenv("NHN_CLOUD_FAKE_SVC_ENDPOINT")
	old := loadedConfig
	loadedConfig = &Config{Extra: map[string]string{}}
	defer func() { loadedConfig = old }()
	if _, err := catalogBaseURL(t.Context(), &svc); err == nil || !strings.Contains(err.Error(), "NHN_CLOUD_FAKE_SVC_ENDPOINT") {
		t.Errorf("{endpoint} without configuration: %v", err)
	}
}

// --header Content-Type overrides the catalog's body content type.
func TestCatalogHeaderOverridesContentType(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	svc := fakeService()
	svc.Operations = append(svc.Operations, catalog.Operation{Name: "put-raw", Summary: "Raw", Method: "PUT", Path: "/v1/raw",
		Body: &catalog.Body{Encoding: "raw", ContentType: "text/plain"}})
	f := filepath.Join(t.TempDir(), "p.bin")
	os.WriteFile(f, []byte{1, 2}, 0600)
	seen, _, err := runFake(t, svc, `{}`, "fake-svc", "put-raw", "--data-file", f, "--header", "Content-Type: application/octet-stream")
	if err != nil {
		t.Fatal(err)
	}
	if seen.contentType != "application/octet-stream" {
		t.Errorf("content type: %q", seen.contentType)
	}
}

// --body null must be an error, not a panic.
func TestCatalogBodyNullIsAnError(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "ak")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	_, _, err := runFake(t, fakeService(), `{}`, "fake-svc", "create-thing", "--thing-name", "a", "--size", "1", "--body", "null")
	if err == nil || !strings.Contains(err.Error(), "JSON object") {
		t.Errorf("got %v", err)
	}
}

// Catalog names must be the canonical name of the command they attach to,
// so credentials and endpoints are looked up under the name users run.
func TestCatalogNamesAreCanonical(t *testing.T) {
	services, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range services {
		if c := findChild(rootCmd, s.Name); c != nil && c.Name() != s.Name {
			t.Errorf("catalog %q attaches to %q through an alias", s.Name, c.Name())
		}
	}
}

func TestCatalogAliasesNeverTakeOverCommands(t *testing.T) {
	root := &cobra.Command{Use: "nhncloud"}
	existing := &cobra.Command{Use: "network", Run: func(*cobra.Command, []string) {}}
	root.AddCommand(existing)
	svc := fakeService()
	svc.Aliases = []string{"network", "fs"}
	addCatalogCommands(root, []catalog.Service{svc})
	c := findChild(root, "fake-svc")
	if c == nil || c.HasAlias("network") || !c.HasAlias("fs") {
		t.Errorf("aliases: %v", c.Aliases)
	}
	if findChild(root, "network") != existing {
		t.Error("an alias took over an existing command")
	}
}

func TestCatalogRejectsDotPathValues(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "ak")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	for _, v := range []string{"..", ".", ""} {
		_, _, err := runFake(t, fakeService(), `{}`, "fake-svc", "create-thing", "--thing-name", v, "--size", "1")
		if err == nil || !strings.Contains(err.Error(), "not a valid path value") {
			t.Errorf("%q: got %v", v, err)
		}
	}
}

// Large integers in --body and in responses stay exact.
func TestCatalogKeepsLargeIntegers(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "ak")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	seen, out, err := runFake(t, fakeService(), `{"id": 9007199254740993}`,
		"fake-svc", "create-thing", "--thing-name", "a", "--size", "1", "--body", `{"spec": {"n": 9007199254740993}}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(seen.body), "9007199254740993") || !strings.Contains(out, "9007199254740993") {
		t.Errorf("body %s, output %s", seen.body, out)
	}
}

// Header-only responses hide sensitive headers; errors and debug output hide
// credentials in the path.
func TestCatalogHidesSecretsInHeadersAndPaths(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("X-Account-Meta-Temp-Url-Key", "TEMPKEY")
		w.Header().Set("X-Account-Object-Count", "7")
	}))
	defer srv.Close()
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "APPKEYVALUE")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	t.Setenv("NHN_CLOUD_FAKE_SVC_ENDPOINT", srv.URL)
	run := func(args ...string) (string, error) {
		root := &cobra.Command{Use: "nhncloud"}
		addCatalogCommands(root, []catalog.Service{fakeService()})
		root.SetArgs(args)
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		err := root.Execute()
		w.Close()
		os.Stdout = old
		out, _ := io.ReadAll(r)
		return string(out), err
	}
	out, err := run("fake-svc", "get-thing")
	if err != nil || strings.Contains(out, "TEMPKEY") || !strings.Contains(out, "X-Account-Object-Count") {
		t.Errorf("headers: %s %v", out, err)
	}
	_, err = run("fake-svc", "create-thing", "--thing-name", "a", "--size", "1")
	if err == nil || strings.Contains(err.Error(), "APPKEYVALUE") {
		t.Errorf("error shows the appkey: %v", err)
	}
	_, err = run("fake-svc", "get-thing", "--header", "X-Secret=VALUE")
	if err == nil || strings.Contains(err.Error(), "VALUE") {
		t.Errorf("bad --header error echoes the value: %v", err)
	}
}

func TestCatalogFormSendsArrayValues(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	svc := fakeService()
	svc.Operations = append(svc.Operations, catalog.Operation{Name: "post-form", Summary: "F", Method: "POST", Path: "/v1/form",
		Body: &catalog.Body{Encoding: "form", Fields: []catalog.Field{{Name: "tag", Type: "array", Items: "string"}}}})
	seen, _, err := runFake(t, svc, `{}`, "fake-svc", "post-form", "--tag", "x,y")
	if err != nil || string(seen.body) != "tag=x&tag=y" {
		t.Errorf("form body %q, %v", seen.body, err)
	}
}

// A catalog secret mark applies to that flag only, not to every flag with
// the same name.
func TestCatalogSecretMarkIsPerFlag(t *testing.T) {
	svc := fakeService()
	svc.Operations[0].Body.Fields = append(svc.Operations[0].Body.Fields, catalog.Field{Name: "name", Type: "string", Secret: true})
	root := &cobra.Command{Use: "nhncloud"}
	other := &cobra.Command{Use: "other", Run: func(*cobra.Command, []string) {}}
	other.Flags().String("name", "", "")
	root.AddCommand(other)
	addCatalogCommands(root, []catalog.Service{svc})
	if isSecretFlag(other.Flags().Lookup("name")) {
		t.Error("a catalog secret mark leaked to another command's flag")
	}
	op := findChild(findChild(root, "fake-svc"), "create-thing")
	if !isSecretFlag(op.Flags().Lookup("name")) {
		t.Error("the catalog flag is not marked secret")
	}
}

func TestCatalogEndpointOverrideIsChecked(t *testing.T) {
	svc := fakeService()
	for _, v := range []string{"https://h.example.com?token=SECRET", "https://u:p@h.example.com", "ftp://h.example.com", "h.example.com"} {
		t.Setenv("NHN_CLOUD_FAKE_SVC_ENDPOINT", v)
		_, err := catalogBaseURL(t.Context(), &svc)
		if err == nil || strings.Contains(err.Error(), "SECRET") {
			t.Errorf("%s: %v", v, err)
		}
	}
	t.Setenv("NHN_CLOUD_FAKE_SVC_ENDPOINT", "https://h.example.com/api/")
	if u, err := catalogBaseURL(t.Context(), &svc); err != nil || u != "https://h.example.com/api" {
		t.Errorf("valid override: %q %v", u, err)
	}
}

func TestCatalogFormChecksNumbers(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	svc := fakeService()
	svc.Operations = append(svc.Operations, catalog.Operation{Name: "post-form", Summary: "F", Method: "POST", Path: "/v1/form",
		Body: &catalog.Body{Encoding: "form", Fields: []catalog.Field{{Name: "size", Type: "integer"}}}})
	if _, _, err := runFake(t, svc, `{}`, "fake-svc", "post-form", "--size", "abc"); err == nil || !strings.Contains(err.Error(), "integer") {
		t.Errorf("got %v", err)
	}
}

func TestCatalogRoundTwoInputs(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "ak")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	svc := fakeService()
	svc.Operations = append(svc.Operations,
		catalog.Operation{Name: "list-things", Summary: "L", Method: "GET", Path: "/v1/things",
			Params: []catalog.Param{{Name: "id", In: "query", Type: "array"}, {Name: "ratio", In: "query", Type: "number"}}})

	seen, _, err := runFake(t, svc, `{}`, "fake-svc", "list-things", "--id", "a,b", "--id", "c")
	if err != nil || seen.query != "id=a&id=b&id=c" {
		t.Errorf("array query %q %v", seen.query, err)
	}
	if _, _, err := runFake(t, svc, `{}`, "fake-svc", "list-things", "--ratio", "NaN"); err == nil {
		t.Error("NaN must be rejected")
	}
	if _, _, err := runFake(t, svc, `{}`, "fake-svc", "create-thing", "--thing-name", "a",
		"--body", `{"thing": {"size": 1}}`); err == nil || !strings.Contains(err.Error(), "wrapped") {
		t.Errorf("double wrap: %v", err)
	}
}

// Local input errors come before credentials are resolved.
func TestCatalogInputErrorsComeBeforeCredentials(t *testing.T) {
	old := loadedConfig
	loadedConfig = &Config{Extra: map[string]string{}}
	defer func() { loadedConfig = old }()
	_, _, err := runFake(t, fakeService(), `{}`, "fake-svc", "upload-file", "--file", "/no/such/file")
	if err == nil || !strings.Contains(err.Error(), "read --file") {
		t.Errorf("want the file error first, got %v", err)
	}
}

func TestStripURLQueriesInHeaders(t *testing.T) {
	if got := stripURLQueries("https://h.example.com/o?temp_url_sig=SIG&x=1"); strings.Contains(got, "SIG") {
		t.Errorf("got %q", got)
	}
}

func TestCatalogFailureMessagesDropURLQueries(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "ak")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	_, _, err := runFake(t, fakeService(), `{"header":{"isSuccessful":false,"resultCode":1,"resultMessage":"see https://h.example.com/x?token=SECRET"}}`,
		"fake-svc", "create-thing", "--thing-name", "a", "--size", "1")
	if err == nil || strings.Contains(err.Error(), "SECRET") || !strings.Contains(err.Error(), "https://h.example.com/x") {
		t.Errorf("got %v", err)
	}
}

func TestIaaSClientRefusesCrossHostRedirects(t *testing.T) {
	var leaked bool
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		leaked = r.Header.Get("X-Auth-Token") != ""
	}))
	defer other.Close()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/same":
			w.WriteHeader(204)
		case "/hop":
			http.Redirect(w, r, "/same", http.StatusFound)
		default:
			http.Redirect(w, r, strings.Replace(other.URL, "127.0.0.1", "localhost", 1)+"/x", http.StatusFound)
		}
	}))
	defer srv.Close()
	req, _ := http.NewRequest("GET", srv.URL+"/away", nil)
	req.Header.Set("X-Auth-Token", "TOKEN")
	if resp, err := iaasHTTPClient.Do(req); err == nil {
		resp.Body.Close()
		t.Error("a redirect to another host must be refused")
	}
	if leaked {
		t.Error("the token followed the redirect")
	}
	req, _ = http.NewRequest("GET", srv.URL+"/hop", nil)
	if resp, err := iaasHTTPClient.Do(req); err != nil || resp.StatusCode != 204 {
		t.Errorf("same-host redirect: %v", err)
	}
}

// Param and header input errors come before the endpoint lookup (which
// issues an IaaS token for IaaS catalogs).
func TestCatalogParamErrorsComeBeforeEndpoint(t *testing.T) {
	svc := fakeService()
	svc.BaseURL = "{endpoint}"
	svc.Operations = append(svc.Operations, catalog.Operation{Name: "list-things", Summary: "L", Method: "GET", Path: "/v1/things",
		Params: []catalog.Param{{Name: "limit", In: "query", Type: "integer"}}})
	t.Setenv("NHN_CLOUD_FAKE_SVC_ENDPOINT", "")
	for _, args := range [][]string{{"--limit", "abc"}, {"--header", "bad"}} {
		root := &cobra.Command{Use: "nhncloud"}
		addCatalogCommands(root, []catalog.Service{svc})
		root.SetArgs(append([]string{"fake-svc", "list-things"}, args...))
		root.SetOut(io.Discard)
		root.SetErr(io.Discard)
		err := root.Execute()
		if err == nil || strings.Contains(err.Error(), "no public endpoint") {
			t.Errorf("%v: got %v", args, err)
		}
	}
}

func TestCatalogRejectsMultipartContentTypeAndEmptyAppkey(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	f := filepath.Join(t.TempDir(), "f")
	os.WriteFile(f, []byte("x"), 0o600)
	if _, _, err := runFake(t, fakeService(), `{}`, "fake-svc", "upload-file", "--file", f,
		"--header", "Content-Type: multipart/form-data"); err == nil || !strings.Contains(err.Error(), "boundary") {
		t.Errorf("multipart Content-Type: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("appkey", ""); err != nil {
		t.Fatal(err)
	}
	defer func() { rootCmd.PersistentFlags().Lookup("appkey").Changed = false }()
	c := &credentialCache{svc: &catalog.Service{Name: "fake-svc"}, values: map[string]string{}}
	if _, err := c.get("appKey"); err == nil || !strings.Contains(err.Error(), "empty") {
		t.Errorf("empty --appkey: %v", err)
	}
}

func TestCatalogAliasesFromLaterCatalog(t *testing.T) {
	root := &cobra.Command{Use: "nhncloud"}
	hand := &cobra.Command{Use: "hand", Run: func(*cobra.Command, []string) {}}
	root.AddCommand(hand)
	a, b := fakeService(), fakeService()
	b.Aliases = []string{"fs"}
	b.Operations = nil
	h := fakeService()
	h.Name, h.Aliases, h.Operations = "hand", []string{"hd"}, nil
	addCatalogCommands(root, []catalog.Service{a, b, h})
	if !findChild(root, "fake-svc").HasAlias("fs") {
		t.Error("aliases of a later catalog with the same name are dropped")
	}
	if hand.HasAlias("hd") {
		t.Error("aliases went onto a hand-written command")
	}
}

// Failures echoing credentials or signed URLs do not show them.
func TestCatalogFailuresHideEchoedSecrets(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "APPKEYVALUE")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "SECRETVALUE")
	_, out, err := runFake(t, fakeService(), `{"header":{"isSuccessful":false,"resultCode":1,"resultMessage":"bad key SECRETVALUE at https://h.example.com/x?a=1,token=TOK"}}`,
		"fake-svc", "create-thing", "--thing-name", "a", "--size", "1")
	for _, leak := range []string{"SECRETVALUE", "TOK"} {
		if err == nil || strings.Contains(err.Error(), leak) || strings.Contains(out, leak) {
			t.Errorf("%s leaked: %v / %s", leak, err, out)
		}
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"error":"path %s rejected"}`, r.URL.Path)
	}))
	defer srv.Close()
	t.Setenv("NHN_CLOUD_FAKE_SVC_ENDPOINT", srv.URL)
	root := &cobra.Command{Use: "nhncloud"}
	addCatalogCommands(root, []catalog.Service{fakeService()})
	root.SetArgs([]string{"fake-svc", "create-thing", "--thing-name", "a", "--size", "1"})
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	if err := root.Execute(); err == nil || strings.Contains(err.Error(), "APPKEYVALUE") {
		t.Errorf("400 body echoes the appkey: %v", err)
	}
}

// A catalog secret read from stdin keeps every line (PEM keys); a catalog
// flag named like a global secret is not secret unless the catalog says so.
func TestCatalogSecretStdinAndMarks(t *testing.T) {
	svc := fakeService()
	svc.Operations[0].Body.Fields = append(svc.Operations[0].Body.Fields,
		catalog.Field{Name: "payload", Type: "string", Secret: true},
		catalog.Field{Name: "password", Type: "string"})
	root := &cobra.Command{Use: "nhncloud"}
	addCatalogCommands(root, []catalog.Service{svc})
	op := findChild(findChild(root, "fake-svc"), "create-thing")
	if isSecretFlag(op.Flags().Lookup("password")) {
		t.Error("an unmarked catalog flag named password is treated as secret")
	}
	op.Flags().Set("payload", "-")
	pem := "-----BEGIN KEY-----\nAAAA\n-----END KEY-----"
	if err := readSecretFlags(op, strings.NewReader(pem+"\n")); err != nil {
		t.Fatal(err)
	}
	if v, _ := op.Flags().GetString("payload"); v != pem {
		t.Errorf("payload %q", v)
	}
}

func TestCatalogLegacyAppKeyNamesAndShownMask(t *testing.T) {
	t.Setenv("NHN_CLOUD_RDS_MYSQL_APPKEY", "")
	t.Setenv("NHN_CLOUD_MYSQL_APPKEY", "LEGACYKEY")
	t.Setenv("NHN_CLOUD_APPKEY", "GLOBALKEY")
	old := loadedConfig
	loadedConfig = &Config{Extra: map[string]string{}}
	defer func() { loadedConfig = old }()
	if v, err := resolveCredential(t.Context(), "rds-mysql", "appKey"); err != nil || v != "LEGACYKEY" {
		t.Errorf("got %q %v", v, err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) }))
	defer srv.Close()
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "APPKEYVALUE")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	t.Setenv("NHN_CLOUD_FAKE_SVC_ENDPOINT", srv.URL)
	root := &cobra.Command{Use: "nhncloud"}
	addCatalogCommands(root, []catalog.Service{fakeService()})
	root.SetArgs([]string{"fake-svc", "create-thing", "--thing-name", "a", "--size", "1"})
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "/appkeys/***/") {
		t.Errorf("masked path: %v", err)
	}
}

// --header values, secret fields given in --body, JSON-escaped forms and
// header-only responses are hidden too.
func TestCatalogRedactsEveryInputSecret(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "ak")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	svc := fakeService()
	svc.Operations[0].Body.Fields = append(svc.Operations[0].Body.Fields, catalog.Field{Name: "pass", Type: "string", Secret: true})
	echo := `{"header":{"isSuccessful":false,"resultCode":1,"resultMessage":"HDRSECRET BODYS&ECRET"}}`
	_, out, err := runFake(t, svc, echo, "fake-svc", "create-thing", "--thing-name", "a", "--size", "1",
		"--header", "X-Custom: HDRSECRET", "--body", `{"pass":"BODYS&ECRET"}`)
	for _, leak := range []string{"HDRSECRET", "BODYS"} {
		if err == nil || strings.Contains(err.Error(), leak) || strings.Contains(out, leak) {
			t.Errorf("%s leaked: %v / %s", leak, err, out)
		}
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Trace", "trace "+r.Header.Get("X-Secret-Key"))
	}))
	defer srv.Close()
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "SECRETVALUE")
	t.Setenv("NHN_CLOUD_FAKE_SVC_ENDPOINT", srv.URL)
	root := &cobra.Command{Use: "nhncloud"}
	addCatalogCommands(root, []catalog.Service{fakeService()})
	root.SetArgs([]string{"fake-svc", "get-thing"})
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	root.Execute()
	w.Close()
	os.Stdout = oldOut
	got, _ := io.ReadAll(r)
	if strings.Contains(string(got), "SECRETVALUE") {
		t.Errorf("header output leaks the secret key: %s", got)
	}
}

func TestCatalogEndpointOverrideRejectsEmptyQueryOrFragment(t *testing.T) {
	svc := fakeService()
	for _, v := range []string{"https://h.example.com?", "https://h.example.com#"} {
		t.Setenv("NHN_CLOUD_FAKE_SVC_ENDPOINT", v)
		if _, err := catalogBaseURL(t.Context(), &svc); err == nil {
			t.Errorf("%s accepted", v)
		}
	}
}

// Secrets echoed with other JSON escapes (\/, \u000a) are hidden too, and so
// are multi-line secrets.
func TestCatalogRedactsDecodedStrings(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "ak")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	svc := fakeService()
	svc.Operations[0].Body.Fields = append(svc.Operations[0].Body.Fields, catalog.Field{Name: "pass", Type: "string", Secret: true})
	echo := `{"header":{"isSuccessful":false,"resultCode":1,"resultMessage":"bad a\/bSECRET line1\u000aline2"}}`
	_, out, err := runFake(t, svc, echo, "fake-svc", "create-thing", "--thing-name", "a", "--size", "1",
		"--body", `{"pass":"a/bSECRET line1\nline2"}`)
	for _, leak := range []string{"bSECRET", "line2"} {
		if err == nil || strings.Contains(err.Error(), leak) || strings.Contains(out, leak) {
			t.Errorf("%s leaked: %v / %s", leak, err, out)
		}
	}
}

func TestStdinUsersIgnoresLiteralDash(t *testing.T) {
	c := &cobra.Command{Use: "x"}
	c.Flags().String("payload", "", "")
	c.Flags().SetAnnotation("payload", secretAnnotation, []string{"true"})
	c.Flags().String("name", "", "")
	c.Flags().SetAnnotation("name", secretAnnotation, []string{"false"})
	c.Flags().Set("payload", "-")
	c.Flags().Set("name", "-")
	if err := readSecretFlags(c, strings.NewReader("l1\nl2\n")); err != nil {
		t.Fatal(err)
	}
	if v, _ := c.Flags().GetString("payload"); v != "l1\nl2" {
		t.Errorf("payload %q", v)
	}
}

func TestRedirectToRewrittenHostIsSameHost(t *testing.T) {
	old := http.DefaultTransport
	defer func() { http.DefaultTransport = old }()
	target, _ := neturl.Parse("http://127.0.0.1:1")
	http.DefaultTransport = &hostRewriteTransport{inner: old, hosts: map[string]*neturl.URL{hostKey("https", "api.example.com"): target}}
	from, _ := http.NewRequest("GET", "https://api.example.com/a", nil)
	to, _ := http.NewRequest("GET", "http://127.0.0.1:1/b", nil)
	if err := sameHostRedirect(to, []*http.Request{from}); err != nil {
		t.Errorf("redirect to the rewritten host refused: %v", err)
	}
}

func TestCatalogHelpAndEmptySecretKey(t *testing.T) {
	root := &cobra.Command{Use: "nhncloud"}
	addCatalogCommands(root, []catalog.Service{fakeService()})
	up := findChild(findChild(root, "fake-svc"), "upload-file")
	if u := up.Flags().Lookup("file").Usage; strings.Contains(u, "--body") {
		t.Errorf("multipart help mentions --body: %q", u)
	}
	if err := rootCmd.PersistentFlags().Set("secret-key", ""); err != nil {
		t.Fatal(err)
	}
	defer func() { rootCmd.PersistentFlags().Lookup("secret-key").Changed = false }()
	if _, err := resolveCredential(t.Context(), "fake-svc", "secretKey"); err == nil || !strings.Contains(err.Error(), "empty") {
		t.Errorf("empty --secret-key: %v", err)
	}
	t.Setenv("NHN_CLOUD_APPKEY", "")
	old := loadedConfig
	loadedConfig = &Config{Extra: map[string]string{}}
	defer func() { loadedConfig = old }()
	if _, err := resolveCredential(t.Context(), "rds-mysql", "appKey"); err == nil || !strings.Contains(err.Error(), "NHN_CLOUD_MYSQL_APPKEY") {
		t.Errorf("appkey error misses legacy names: %v", err)
	}
}

func TestCatalogRedactsNumericSecretFromBody(t *testing.T) {
	t.Setenv("NHN_CLOUD_FAKE_SVC_APPKEY", "ak")
	t.Setenv("NHN_CLOUD_FAKE_SVC_SECRET_KEY", "sk")
	svc := fakeService()
	svc.Operations[0].Body.Fields = append(svc.Operations[0].Body.Fields, catalog.Field{Name: "pin", Type: "integer", Secret: true})
	echo := `{"header":{"isSuccessful":false,"resultCode":1,"resultMessage":"pin 482913 rejected"},"pin":482913}`
	_, out, err := runFake(t, svc, echo, "fake-svc", "create-thing", "--thing-name", "a", "--size", "1", "--body", `{"pin":482913}`)
	if err == nil || strings.Contains(err.Error(), "482913") || strings.Contains(out, "482913") {
		t.Errorf("numeric secret leaked: %v / %s", err, out)
	}
}

// Catalog flags must not hide the persistent flags of hand-written parents,
// and catalogs never attach under a hand-written leaf command.
func TestCatalogFitsHandWrittenCommands(t *testing.T) {
	services, err := catalog.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range services {
		c := findChild(rootCmd, s.Name)
		if c == nil {
			continue
		}
		if c.Runnable() && !c.HasSubCommands() {
			t.Errorf("catalog %s attaches under the leaf command %s", s.Title, c.Name())
		}
		inherited := map[string]bool{}
		for p := c; p != nil; p = p.Parent() {
			p.PersistentFlags().VisitAll(func(f *pflag.Flag) { inherited[f.Name] = true })
		}
		for _, op := range s.Operations {
			for _, p := range op.Params {
				if p.Credential == "" && inherited[p.FlagName()] {
					t.Errorf("%s %s: --%s hides a persistent flag", s.Name, op.Name, p.FlagName())
				}
			}
			if op.Body != nil {
				for _, f := range op.Body.Fields {
					if f.Scalar() && inherited[f.FlagName()] {
						t.Errorf("%s %s: --%s hides a persistent flag", s.Name, op.Name, f.FlagName())
					}
				}
			}
		}
	}
}

func TestCatalogRedactsKeysAndEncodedForms(t *testing.T) {
	got := redactValue(map[string]interface{}{"p@ss/word": "bad p%40ss%2Fword"}, []string{"p@ss/word"})
	b, _ := json.Marshal(got)
	if strings.Contains(string(b), "ss") {
		t.Errorf("got %s", b)
	}
}

func TestCatalogOverrideNeedsPlainPath(t *testing.T) {
	svc := fakeService()
	t.Setenv("NHN_CLOUD_FAKE_SVC_ENDPOINT", "https://h.example.com/my api")
	if _, err := catalogBaseURL(t.Context(), &svc); err == nil {
		t.Error("an override path that needs escaping is accepted")
	}
}
