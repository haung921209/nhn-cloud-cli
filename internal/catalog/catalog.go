// Package catalog holds machine-readable descriptions of NHN Cloud public
// APIs, one JSON file per service under data/. The CLI turns every operation
// into a command, so covering an API is a catalog entry, not code.
package catalog

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strings"
)

//go:embed data/*.json
var dataFS embed.FS

// Service describes one service API guide.
type Service struct {
	// Name is the CLI command the operations go under. Several catalogs may
	// share a name (e.g. "network"); their operation names must not clash.
	Name    string   `json:"name"`
	Aliases []string `json:"aliases,omitempty"`
	Title   string   `json:"title"`
	Docs    string   `json:"docs"`
	Version string   `json:"version,omitempty"`

	// Exactly one base URL source:
	//   iaasServiceType: endpoint from the IaaS token's service catalog
	//   regionBaseURLs:  region -> URL
	//   baseURL:         one URL; "{region}" is replaced by the region.
	//                    "{endpoint}" means the guide gives no public URL (a
	//                    per-account host); users set it per service.
	// Any service's base URL can be overridden with NHN_CLOUD_<SERVICE>_ENDPOINT
	// or <service>_endpoint in the profile.
	IaaSServiceType string            `json:"iaasServiceType,omitempty"`
	RegionBaseURLs  map[string]string `json:"regionBaseURLs,omitempty"`
	BaseURL         string            `json:"baseURL,omitempty"`

	Auth       Auth        `json:"auth"`
	Operations []Operation `json:"operations"`
}

// Auth lists request headers and their value templates. Templates name
// credentials in braces, e.g. "Bearer {userAccessToken}".
type Auth struct {
	Headers map[string]string `json:"headers,omitempty"`
}

// Operation is one API call.
type Operation struct {
	Name    string  `json:"name"`
	Summary string  `json:"summary"`
	Method  string  `json:"method"`
	Path    string  `json:"path"`
	Params  []Param `json:"params,omitempty"`
	Body    *Body   `json:"body,omitempty"`
}

// Param is a path, query or header parameter.
type Param struct {
	Name        string `json:"name"`
	In          string `json:"in"` // path, query, header
	Type        string `json:"type,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
	// Credential fills the value from a credential (no flag), e.g. appKey.
	Credential string `json:"credential,omitempty"`
	// Flag overrides the flag name derived from Name.
	Flag string `json:"flag,omitempty"`
	// Secret params accept "-" to read the value from stdin.
	Secret bool `json:"secret,omitempty"`
}

// Body describes the request body.
type Body struct {
	// Encoding is json (default), multipart (form-data, "file" fields are
	// uploaded), form (x-www-form-urlencoded) or raw (the --data-file bytes
	// are sent as is with ContentType).
	Encoding    string `json:"encoding,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	// Wrapper nests JSON fields under one key, e.g. {"vpc": {...}}.
	Wrapper string  `json:"wrapper,omitempty"`
	Fields  []Field `json:"fields,omitempty"`
}

// Field is a top-level body field. Scalars (and arrays of strings) become
// flags; objects and other arrays are passed with --body.
type Field struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // string, integer, number, boolean, array, object, file
	Items       string `json:"items,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
	Flag        string `json:"flag,omitempty"`
	// Secret fields accept "-" to read the value from stdin.
	Secret bool `json:"secret,omitempty"`
}

// Credentials the CLI knows how to fill.
var Credentials = map[string]bool{
	"iaasToken": true, "userAccessToken": true, "accessKeyId": true,
	"secretAccessKey": true, "appKey": true, "secretKey": true, "tenantId": true,
}

// ReservedFlags are global CLI flags or flags the generator adds itself.
var ReservedFlags = map[string]bool{
	"region": true, "appkey": true, "debug": true, "output": true, "query": true,
	"username": true, "password": true, "tenant-id": true, "profile": true,
	"endpoint-override": true, "help": true, "body": true, "secret-key": true,
	"data-file": true, "header": true,
}

var (
	kebabRe    = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
	pathVarRe  = regexp.MustCompile(`\{([^{}]+)\}`)
	templateRe = regexp.MustCompile(`\{([A-Za-z]+)\}`)
	upperRe    = regexp.MustCompile(`([a-z0-9])([A-Z])`)
)

// FlagName turns a parameter name into a kebab-case flag name.
func FlagName(name string) string {
	s := upperRe.ReplaceAllString(name, "${1}-${2}")
	s = strings.NewReplacer("_", "-", ".", "-", " ", "-", ":", "-").Replace(s)
	return strings.ToLower(strings.Trim(s, "-"))
}

func (p Param) FlagName() string {
	if p.Flag != "" {
		return p.Flag
	}
	return FlagName(p.Name)
}

func (f Field) FlagName() string {
	if f.Flag != "" {
		return f.Flag
	}
	return FlagName(f.Name)
}

// Scalar reports whether the field gets its own flag.
func (f Field) Scalar() bool {
	switch f.Type {
	case "string", "integer", "number", "boolean", "file":
		return true
	case "array":
		return f.Items == "string"
	}
	return false
}

// Load reads every embedded catalog.
func Load() ([]Service, error) { return LoadFS(dataFS, "data") }

// LoadFS reads the *.json catalogs in dir of fsys, sorted by file name.
func LoadFS(fsys fs.FS, dir string) ([]Service, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}
	var out []Service
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := fs.ReadFile(fsys, dir+"/"+e.Name())
		if err != nil {
			return nil, err
		}
		var s Service
		dec := json.NewDecoder(strings.NewReader(string(data)))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&s); err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		out = append(out, s)
	}
	return out, nil
}

// Validate checks one catalog and returns every problem found.
func Validate(s Service) []string {
	var errs []string
	bad := func(format string, a ...interface{}) { errs = append(errs, fmt.Sprintf(format, a...)) }

	if !kebabRe.MatchString(s.Name) {
		bad("service name %q is not kebab-case", s.Name)
	}
	if s.Title == "" || !strings.HasPrefix(s.Docs, "https://docs.nhncloud.com/") {
		bad("title and a docs.nhncloud.com docs URL are required")
	}
	sources := 0
	if s.IaaSServiceType != "" {
		sources++
	}
	if len(s.RegionBaseURLs) > 0 {
		sources++
	}
	if s.BaseURL != "" {
		sources++
	}
	if sources != 1 {
		bad("exactly one of iaasServiceType, regionBaseURLs, baseURL is required")
	}
	for r, u := range s.RegionBaseURLs {
		if !strings.HasPrefix(u, "https://") {
			bad("regionBaseURLs[%s] must be https", r)
		}
	}
	if s.BaseURL != "" && s.BaseURL != "{endpoint}" && !strings.HasPrefix(s.BaseURL, "https://") {
		bad("baseURL must be https or {endpoint}")
	}
	for h, tmpl := range s.Auth.Headers {
		for _, m := range templateRe.FindAllStringSubmatch(tmpl, -1) {
			if !Credentials[m[1]] {
				bad("auth header %s uses unknown credential %q", h, m[1])
			}
		}
	}
	if len(s.Operations) == 0 {
		bad("no operations")
	}
	seenOp := map[string]bool{}
	for _, op := range s.Operations {
		where := s.Name + " " + op.Name
		if !kebabRe.MatchString(op.Name) {
			bad("%s: operation name is not kebab-case", where)
		}
		if seenOp[op.Name] {
			bad("%s: duplicate operation name", where)
		}
		seenOp[op.Name] = true
		switch op.Method {
		case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "COPY":
		default:
			bad("%s: bad method %q", where, op.Method)
		}
		if !strings.HasPrefix(op.Path, "/") || strings.Contains(op.Path, "?") {
			bad("%s: path must start with / and carry no query", where)
		}
		// Guides sometimes show a trailing slash the servers reject (404).
		if len(op.Path) > 1 && strings.HasSuffix(op.Path, "/") {
			bad("%s: path must not end with /", where)
		}
		if op.Summary == "" {
			bad("%s: summary is required", where)
		}
		pathParams := map[string]bool{}
		flags := map[string]string{}
		addFlag := func(name, what string) {
			if ReservedFlags[name] {
				bad("%s: %s flag %q clashes with a global flag; set \"flag\"", where, what, name)
			}
			if prev, ok := flags[name]; ok {
				bad("%s: flag %q used by %s and %s; set \"flag\"", where, name, prev, what)
			}
			flags[name] = what
		}
		for _, p := range op.Params {
			switch p.In {
			case "path", "query", "header":
			default:
				bad("%s: param %s has bad in %q", where, p.Name, p.In)
			}
			if p.Credential != "" && !Credentials[p.Credential] {
				bad("%s: param %s uses unknown credential %q", where, p.Name, p.Credential)
			}
			if p.In == "path" {
				pathParams[p.Name] = true
			}
			if p.Credential == "" {
				addFlag(p.FlagName(), "param "+p.Name)
			}
		}
		for _, m := range pathVarRe.FindAllStringSubmatch(op.Path, -1) {
			if !pathParams[m[1]] {
				bad("%s: path variable {%s} has no path param", where, m[1])
			}
		}
		if op.Body != nil {
			if op.Method == "GET" {
				bad("%s: GET with a body", where)
			}
			switch op.Body.Encoding {
			case "", "json", "multipart", "form":
			case "raw":
				if op.Body.ContentType == "" || len(op.Body.Fields) > 0 {
					bad("%s: raw body needs contentType and no fields", where)
				}
			default:
				bad("%s: bad body encoding %q", where, op.Body.Encoding)
			}
			if op.Body.Wrapper != "" && op.Body.Encoding != "" && op.Body.Encoding != "json" {
				bad("%s: wrapper is only for json bodies", where)
			}
			for _, f := range op.Body.Fields {
				switch f.Type {
				case "string", "integer", "number", "boolean", "array", "object":
				case "file":
					if op.Body.Encoding != "multipart" {
						bad("%s: file field %s needs multipart encoding", where, f.Name)
					}
				default:
					bad("%s: body field %s has bad type %q", where, f.Name, f.Type)
				}
				if f.Scalar() {
					addFlag(f.FlagName(), "body field "+f.Name)
				}
			}
		}
	}
	return errs
}

// ValidateAll validates every catalog and the operation names shared by
// catalogs under the same command name.
func ValidateAll(services []Service) []string {
	var errs []string
	owner := map[string]string{}
	for _, s := range services {
		for _, e := range Validate(s) {
			errs = append(errs, s.Title+": "+e)
		}
		for _, op := range s.Operations {
			key := s.Name + " " + op.Name
			if prev, ok := owner[key]; ok && prev != s.Title {
				errs = append(errs, fmt.Sprintf("operation %q defined by both %s and %s", key, prev, s.Title))
			}
			owner[key] = s.Title
		}
	}
	sort.Strings(errs)
	return errs
}
