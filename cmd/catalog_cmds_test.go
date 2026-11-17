package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/haung921209/nhn-cloud-cli/internal/catalog"
	"github.com/spf13/cobra"
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
