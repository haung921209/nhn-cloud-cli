package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/haung921209/nhn-cloud-cli/internal/auth"
	"github.com/haung921209/nhn-cloud-cli/internal/catalog"
	"github.com/spf13/cobra"
)

// Catalog commands: every operation of internal/catalog/data/*.json becomes
// "nhncloud <service> <operation>". Hand-written commands win: an operation
// whose name an existing command already uses is not added.

var catalogAdded bool

// ensureCatalogCommands adds the embedded catalogs once. Invalid catalogs are
// a build error caught by tests, so they are skipped here without noise.
func ensureCatalogCommands() {
	if catalogAdded {
		return
	}
	catalogAdded = true
	services, err := catalog.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: API catalog not loaded: %v\n", err)
		return
	}
	addCatalogCommands(rootCmd, services)
}

// addCatalogCommands attaches catalog operations to root. It runs from
// Execute, after every package init has registered the hand-written commands.
func addCatalogCommands(root *cobra.Command, services []catalog.Service) {
	for i := range services {
		svc := &services[i]
		parent := findChild(root, svc.Name)
		if parent == nil {
			parent = &cobra.Command{
				Use:     svc.Name,
				Aliases: svc.Aliases,
				Short:   svc.Title + " (API catalog)",
			}
			root.AddCommand(parent)
		}
		for j := range svc.Operations {
			op := &svc.Operations[j]
			if findChild(parent, op.Name) != nil {
				continue
			}
			parent.AddCommand(newCatalogCommand(svc, op))
		}
	}
}

func findChild(parent *cobra.Command, name string) *cobra.Command {
	for _, c := range parent.Commands() {
		if c.Name() == name || c.HasAlias(name) {
			return c
		}
	}
	return nil
}

func newCatalogCommand(svc *catalog.Service, op *catalog.Operation) *cobra.Command {
	c := &cobra.Command{
		Use:   op.Name,
		Short: op.Summary,
		Long: fmt.Sprintf("%s\n\n%s %s\nAPI guide: %s\n\nOutput is the API response (JSON); use -o yaml or --query to reshape it.",
			op.Summary, op.Method, op.Path, svc.Docs),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCatalogOperation(cmd, svc, op)
		},
	}
	f := c.Flags()
	f.StringArray("header", nil, "Extra request header Name:Value (repeatable)")
	for _, p := range op.Params {
		if p.Credential != "" {
			continue
		}
		usage := p.Description
		if p.Type == "boolean" {
			f.Bool(p.FlagName(), false, usage)
		} else {
			f.String(p.FlagName(), "", usage)
		}
		if p.Required {
			c.MarkFlagRequired(p.FlagName())
		}
		if p.Secret {
			secretFlagNames[p.FlagName()] = true
		}
	}
	if op.Body != nil {
		for _, fd := range op.Body.Fields {
			if !fd.Scalar() {
				continue
			}
			usage := fd.Description
			if fd.Required {
				usage += " (required unless given in --body)"
			}
			switch {
			case fd.Type == "boolean":
				f.Bool(fd.FlagName(), false, usage)
			case fd.Type == "array":
				f.StringSlice(fd.FlagName(), nil, usage)
			case fd.Type == "file":
				f.String(fd.FlagName(), "", usage+" (file path)")
			default:
				f.String(fd.FlagName(), "", usage)
			}
			if fd.Secret {
				secretFlagNames[fd.FlagName()] = true
			}
		}
		switch op.Body.Encoding {
		case "raw":
			f.String("data-file", "", "File to send as the request body (- for stdin)")
			c.MarkFlagRequired("data-file")
		case "", "json":
			f.String("body", "", "JSON request body (fields not covered by flags), or @file")
		}
	}
	return c
}

func runCatalogOperation(cmd *cobra.Command, svc *catalog.Service, op *catalog.Operation) error {
	ctx := context.Background()
	creds := &credentialCache{svc: svc, ctx: ctx, values: map[string]string{}}

	base, err := catalogBaseURL(ctx, svc)
	if err != nil {
		return err
	}

	path := op.Path
	query := url.Values{}
	headers := http.Header{}
	for _, p := range op.Params {
		var value string
		var set bool
		if p.Credential != "" {
			value, err = creds.get(p.Credential)
			if err != nil {
				return err
			}
			set = true
		} else if cmd.Flags().Changed(p.FlagName()) {
			value = cmd.Flags().Lookup(p.FlagName()).Value.String()
			set = true
		}
		if !set {
			continue
		}
		switch p.In {
		case "path":
			path = strings.ReplaceAll(path, "{"+p.Name+"}", url.PathEscape(value))
		case "query":
			query.Set(p.Name, value)
		case "header":
			headers.Set(p.Name, value)
		}
	}
	extra, _ := cmd.Flags().GetStringArray("header")
	for _, h := range extra {
		name, value, ok := strings.Cut(h, ":")
		if !ok || strings.TrimSpace(name) == "" {
			return fmt.Errorf("--header %q: want Name:Value", h)
		}
		headers.Add(strings.TrimSpace(name), strings.TrimSpace(value))
	}
	for name, tmpl := range svc.Auth.Headers {
		value, err := creds.expand(tmpl)
		if err != nil {
			return err
		}
		headers.Set(name, value)
	}

	body, contentType, err := catalogBody(cmd, op)
	if err != nil {
		return err
	}
	if contentType != "" {
		headers.Set("Content-Type", contentType)
	}

	full := base + path
	if len(query) > 0 {
		full += "?" + query.Encode()
	}
	data, respHeader, err := catalogDo(ctx, op.Method, full, headers, body)
	if err != nil {
		return err
	}
	return printCatalogResponse(data, respHeader)
}

// credentialCache resolves credentials once per command.
type credentialCache struct {
	svc    *catalog.Service
	ctx    context.Context
	values map[string]string
}

var templateVar = regexp.MustCompile(`\{([A-Za-z]+)\}`)

func (c *credentialCache) expand(tmpl string) (string, error) {
	var firstErr error
	out := templateVar.ReplaceAllStringFunc(tmpl, func(m string) string {
		v, err := c.get(m[1 : len(m)-1])
		if err != nil && firstErr == nil {
			firstErr = err
		}
		return v
	})
	return out, firstErr
}

func (c *credentialCache) get(name string) (string, error) {
	if v, ok := c.values[name]; ok {
		return v, nil
	}
	v, err := resolveCredential(c.ctx, c.svc.Name, name)
	if err != nil {
		return "", err
	}
	c.values[name] = v
	return v, nil
}

// serviceEnvPrefix turns "rds-mysql" into "NHN_CLOUD_RDS_MYSQL_".
func serviceEnvPrefix(service string) string {
	return "NHN_CLOUD_" + strings.ToUpper(strings.ReplaceAll(service, "-", "_")) + "_"
}

// serviceConfigKey turns ("rds-mysql", "app_key") into "rds_mysql_app_key".
func serviceConfigKey(service, key string) string {
	return strings.ReplaceAll(service, "-", "_") + "_" + key
}

func resolveCredential(ctx context.Context, service, name string) (string, error) {
	cfg := LoadConfig()
	lookup := func(envSuffix, cfgKey string) string {
		if v := os.Getenv(serviceEnvPrefix(service) + envSuffix); v != "" {
			return v
		}
		return cfg.Extra[serviceConfigKey(service, cfgKey)]
	}
	switch name {
	case "iaasToken":
		api, err := getIaaSAPI(ctx)
		if err != nil {
			return "", err
		}
		return api.token, nil
	case "tenantId":
		if t := getTenantID(); t != "" {
			return t, nil
		}
	case "accessKeyId":
		if v := getAccessKey(); v != "" {
			return v, nil
		}
	case "secretAccessKey":
		if v := getSecretKey(); v != "" {
			return v, nil
		}
	case "userAccessToken":
		ak, sk := getAccessKey(), getSecretKey()
		if ak == "" || sk == "" {
			return "", fmt.Errorf("a User Access Key is required: access_key_id and secret_access_key in the profile, or NHN_CLOUD_ACCESS_KEY and NHN_CLOUD_SECRET_KEY")
		}
		tok, err := auth.NewTokenManager(getRegion(), ak, sk).GetToken()
		if err != nil {
			return "", fmt.Errorf("issue User Access Key token: %w", err)
		}
		return tok.AccessToken, nil
	case "appKey":
		if rootCmd.PersistentFlags().Changed("appkey") {
			return appKey, nil
		}
		if v := lookup("APPKEY", "app_key"); v != "" {
			return v, nil
		}
		if v := os.Getenv("NHN_CLOUD_APPKEY"); v != "" {
			return v, nil
		}
		if cfg.AppKey != "" {
			return cfg.AppKey, nil
		}
		return "", fmt.Errorf("an appkey is required: --appkey, %sAPPKEY, or %s in the profile", serviceEnvPrefix(service), serviceConfigKey(service, "app_key"))
	case "secretKey":
		if secretKeyFlag != "" {
			return secretKeyFlag, nil
		}
		if v := lookup("SECRET_KEY", "secret_key"); v != "" {
			return v, nil
		}
		return "", fmt.Errorf("a secret key is required: --secret-key, %sSECRET_KEY, or %s in the profile", serviceEnvPrefix(service), serviceConfigKey(service, "secret_key"))
	}
	return "", fmt.Errorf("credential %s is not configured", name)
}

func catalogBaseURL(ctx context.Context, svc *catalog.Service) (string, error) {
	region := getRegion()
	if v := os.Getenv(serviceEnvPrefix(svc.Name) + "ENDPOINT"); v != "" {
		return strings.TrimSuffix(v, "/"), nil
	}
	if v := LoadConfig().Extra[serviceConfigKey(svc.Name, "endpoint")]; v != "" {
		return strings.TrimSuffix(v, "/"), nil
	}
	if svc.BaseURL == "{endpoint}" {
		return "", fmt.Errorf("%s has no public endpoint; set the one from the console in %sENDPOINT or %s in the profile",
			svc.Title, serviceEnvPrefix(svc.Name), serviceConfigKey(svc.Name, "endpoint"))
	}
	switch {
	case svc.IaaSServiceType != "":
		api, err := getIaaSAPI(ctx)
		if err != nil {
			return "", err
		}
		return api.endpoint(svc.IaaSServiceType)
	case len(svc.RegionBaseURLs) > 0:
		if u, ok := svc.RegionBaseURLs[strings.ToLower(region)]; ok {
			return strings.TrimSuffix(u, "/"), nil
		}
		regions := make([]string, 0, len(svc.RegionBaseURLs))
		for r := range svc.RegionBaseURLs {
			regions = append(regions, r)
		}
		sort.Strings(regions)
		return "", fmt.Errorf("%s has no endpoint for region %q (available: %s)", svc.Title, region, strings.Join(regions, ", "))
	default:
		return strings.TrimSuffix(strings.ReplaceAll(svc.BaseURL, "{region}", strings.ToLower(region)), "/"), nil
	}
}

// catalogBody builds the request body from --body and the field flags.
func catalogBody(cmd *cobra.Command, op *catalog.Operation) (io.Reader, string, error) {
	if op.Body == nil {
		return nil, "", nil
	}
	flags := cmd.Flags()
	switch op.Body.Encoding {
	case "raw":
		name, _ := flags.GetString("data-file")
		var data []byte
		var err error
		if name == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(name)
		}
		if err != nil {
			return nil, "", fmt.Errorf("read --data-file: %w", err)
		}
		return bytes.NewReader(data), op.Body.ContentType, nil

	case "multipart":
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		for _, fd := range op.Body.Fields {
			if !flags.Changed(fd.FlagName()) {
				if fd.Required {
					return nil, "", fmt.Errorf("--%s is required", fd.FlagName())
				}
				continue
			}
			value := flags.Lookup(fd.FlagName()).Value.String()
			if fd.Type == "file" {
				data, err := os.ReadFile(value)
				if err != nil {
					return nil, "", fmt.Errorf("read --%s: %w", fd.FlagName(), err)
				}
				part, err := w.CreateFormFile(fd.Name, filepath.Base(value))
				if err != nil {
					return nil, "", err
				}
				part.Write(data)
				continue
			}
			w.WriteField(fd.Name, value)
		}
		if err := w.Close(); err != nil {
			return nil, "", err
		}
		return &buf, w.FormDataContentType(), nil

	case "form":
		form := url.Values{}
		for _, fd := range op.Body.Fields {
			if flags.Changed(fd.FlagName()) {
				form.Set(fd.Name, flags.Lookup(fd.FlagName()).Value.String())
			} else if fd.Required {
				return nil, "", fmt.Errorf("--%s is required", fd.FlagName())
			}
		}
		return strings.NewReader(form.Encode()), "application/x-www-form-urlencoded", nil
	}

	obj := map[string]interface{}{}
	if raw, _ := flags.GetString("body"); raw != "" {
		data, err := readValueOrFile(raw)
		if err != nil {
			return nil, "", fmt.Errorf("read --body: %w", err)
		}
		if err := json.Unmarshal(data, &obj); err != nil {
			return nil, "", fmt.Errorf("--body is not a JSON object: %w", err)
		}
	}
	for _, fd := range op.Body.Fields {
		if fd.Scalar() && flags.Changed(fd.FlagName()) {
			v, err := fieldValue(cmd, fd)
			if err != nil {
				return nil, "", err
			}
			obj[fd.Name] = v
		}
		if _, ok := obj[fd.Name]; fd.Required && !ok {
			if fd.Scalar() {
				return nil, "", fmt.Errorf("--%s is required (or set %q in --body)", fd.FlagName(), fd.Name)
			}
			return nil, "", fmt.Errorf("%q is required in --body", fd.Name)
		}
	}
	var payload interface{} = obj
	if op.Body.Wrapper != "" {
		payload = map[string]interface{}{op.Body.Wrapper: obj}
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, "", err
	}
	return bytes.NewReader(data), "application/json", nil
}

func fieldValue(cmd *cobra.Command, fd catalog.Field) (interface{}, error) {
	flags := cmd.Flags()
	name := fd.FlagName()
	switch fd.Type {
	case "boolean":
		return flags.GetBool(name)
	case "array":
		return flags.GetStringSlice(name)
	case "integer":
		s, _ := flags.GetString(name)
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("--%s must be an integer", name)
		}
		return n, nil
	case "number":
		s, _ := flags.GetString(name)
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("--%s must be a number", name)
		}
		return n, nil
	}
	s, _ := flags.GetString(name)
	return s, nil
}

func catalogDo(ctx context.Context, method, rawURL string, headers http.Header, body io.Reader) ([]byte, http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %s: invalid request URL", method, stripQuery(rawURL))
	}
	req.Header = headers
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s %s\n", method, stripQuery(rawURL))
	}
	resp, err := iaasHTTPClient.Do(req)
	if err != nil {
		var ue *url.Error
		if errors.As(err, &ue) {
			return nil, nil, fmt.Errorf("%s %s: %w", method, stripQuery(rawURL), ue.Err)
		}
		return nil, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, nil, &apiError{Method: method, URL: stripQuery(rawURL), Status: resp.StatusCode, Body: string(data)}
	}
	return data, resp.Header, nil
}

func printCatalogResponse(data []byte, header http.Header) error {
	if len(bytes.TrimSpace(data)) == 0 {
		// Some APIs (e.g. HEAD on object storage) answer only in headers.
		h := map[string]string{}
		for k, v := range header {
			if k != "Date" && k != "Content-Length" && k != "Connection" && k != "Set-Cookie" {
				h[k] = strings.Join(v, ", ")
			}
		}
		emit(map[string]interface{}{"headers": h}, nil)
		return nil
	}
	var v interface{}
	if json.Unmarshal(data, &v) == nil {
		emit(v, nil)
		// Many NHN Cloud APIs answer 200 with {"header": {"isSuccessful": false}}.
		if m, ok := v.(map[string]interface{}); ok {
			if h, ok := m["header"].(map[string]interface{}); ok && h["isSuccessful"] == false {
				return fmt.Errorf("the API reported a failure: resultCode %v, %v", h["resultCode"], h["resultMessage"])
			}
		}
		return nil
	}
	_, err := os.Stdout.Write(data)
	return err
}
