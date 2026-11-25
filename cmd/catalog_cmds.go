package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
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
	"github.com/spf13/pflag"
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
	// Invalid catalogs could redefine flags and panic; skip them.
	var valid []catalog.Service
	for _, s := range services {
		if errs := catalog.Validate(s); len(errs) > 0 {
			// Flags are not parsed yet, so this cannot wait for --debug.
			fmt.Fprintf(os.Stderr, "Warning: skipping API catalog %s: %s\n", s.Title, errs[0])
			continue
		}
		valid = append(valid, s)
	}
	addCatalogCommands(rootCmd, valid)
}

// addCatalogCommands attaches catalog operations to root. It runs from
// Execute, after every package init has registered the hand-written commands.
func addCatalogCommands(root *cobra.Command, services []catalog.Service) {
	created := map[*cobra.Command]bool{}
	for i := range services {
		svc := &services[i]
		parent := findChild(root, svc.Name)
		if parent == nil {
			parent = &cobra.Command{Use: svc.Name, Short: svc.Title + " (API catalog)"}
			root.AddCommand(parent)
			created[parent] = true
		}
		// Aliases go on catalog-made commands only, never on hand-written ones,
		// and never take over an existing command.
		for _, a := range svc.Aliases {
			if created[parent] && findChild(root, a) == nil {
				parent.Aliases = append(parent.Aliases, a)
			}
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
		Long: fmt.Sprintf("%s\n\n%s %s\nAPI guide: %s\nEndpoint override: %sENDPOINT or %s in the profile\n\nOutput is the API response (JSON); use -o yaml or --query to reshape it.",
			op.Summary, op.Method, op.Path, svc.Docs, serviceEnvPrefix(svc.OverrideName()), serviceConfigKey(svc.OverrideName(), "endpoint")),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// API errors are not usage errors; flag errors still show usage.
			cmd.SilenceUsage = true
			return runCatalogOperation(cmd, svc, op)
		},
	}
	f := c.Flags()
	f.StringArray("header", nil, "Extra request header Name:Value (repeatable)")
	if op.SaveField != "" {
		f.String("output-file", "", fmt.Sprintf("Write the response field %q to this file with mode 0600 instead of printing the response (an existing file is replaced, not merged)", op.SaveField))
	}
	for _, p := range op.Params {
		if p.Credential != "" {
			continue
		}
		usage := p.Description
		switch {
		case p.Type == "boolean":
			f.Bool(p.FlagName(), false, usage)
		case p.Type == "array" && p.In == "query":
			f.StringSlice(p.FlagName(), nil, usage+" (repeatable or comma separated)")
		default:
			f.String(p.FlagName(), "", usage)
		}
		if p.Required {
			c.MarkFlagRequired(p.FlagName())
		}
		// Catalog flags follow the catalog's mark only, not the global names.
		f.SetAnnotation(p.FlagName(), secretAnnotation, []string{strconv.FormatBool(p.Secret)})
	}
	if op.Body != nil {
		for _, fd := range op.Body.Fields {
			if !fd.Scalar() {
				continue
			}
			usage := fd.Description
			if fd.Required && (op.Body.Encoding == "" || op.Body.Encoding == "json") {
				usage += " (required unless given in --body)"
			} else if fd.Required {
				usage += " (required)"
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
			f.SetAnnotation(fd.FlagName(), secretAnnotation, []string{strconv.FormatBool(fd.Secret)})
		}
		switch op.Body.Encoding {
		case "raw":
			f.String("data-file", "", "File to send as the request body (- for stdin)")
			c.MarkFlagRequired("data-file")
		case "", "json":
			usage := "JSON request body (fields not covered by flags), or @file"
			if op.Body.Wrapper != "" {
				usage = fmt.Sprintf("JSON object that goes inside %q (fields not covered by flags), or @file", op.Body.Wrapper)
			}
			f.String("body", "", usage)
		}
	}
	return c
}

func runCatalogOperation(cmd *cobra.Command, svc *catalog.Service, op *catalog.Operation) error {
	ctx := context.Background()
	creds := &credentialCache{svc: svc, ctx: ctx, values: map[string]string{}}

	// Local input errors come before any network call (token issue).
	body, contentType, bodySecrets, err := catalogBody(cmd, op)
	if err != nil {
		return err
	}
	for _, p := range op.Params {
		if p.Credential == "" && cmd.Flags().Changed(p.FlagName()) {
			if err := checkParamValue(p, cmd.Flags().Lookup(p.FlagName()).Value.String()); err != nil {
				return err
			}
		}
	}
	extraHeaders, err := parseExtraHeaders(cmd, op)
	if err != nil {
		return err
	}
	base, err := catalogBaseURL(ctx, svc)
	if err != nil {
		return err
	}

	path, shown := op.Path, op.Path // shown hides credentials and secrets
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
		if err := checkParamValue(p, value); err != nil {
			return err
		}
		switch p.In {
		case "path":
			path = strings.ReplaceAll(path, "{"+p.Name+"}", url.PathEscape(value))
			masked := url.PathEscape(value)
			if p.Credential != "" || p.Secret {
				masked = "***"
			}
			shown = strings.ReplaceAll(shown, "{"+p.Name+"}", masked)
		case "query":
			if p.Type == "array" && p.Credential == "" {
				values, _ := cmd.Flags().GetStringSlice(p.FlagName())
				for _, v := range values {
					query.Add(p.Name, v)
				}
			} else {
				query.Set(p.Name, value)
			}
		case "header":
			headers.Set(p.Name, value)
		}
	}
	for _, h := range extraHeaders {
		headers.Add(h[0], h[1])
	}
	if m := templateVar.FindString(path); m != "" {
		return fmt.Errorf("path variable %s has no value", m)
	}
	for name, tmpl := range svc.Auth.Headers {
		value, err := creds.expand(tmpl)
		if err != nil {
			return err
		}
		headers.Set(name, value)
	}

	// A Content-Type given with --header wins (e.g. raw payload types).
	if contentType != "" && headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", contentType)
	}

	full := base + path
	if len(query) > 0 {
		full += "?" + query.Encode()
	}
	data, respHeader, err := catalogDo(ctx, op.Method, full, base+shown, headers, body)
	secrets := usedSecrets(cmd, creds)
	secrets = append(secrets, bodySecrets...)
	for _, h := range extraHeaders {
		secrets = append(secrets, h[1]) // --header values may be credentials
	}
	if ae, ok := err.(*apiError); ok {
		ae.Body = redactBody(ae.Body, secrets) // servers may echo credentials
	}
	if err != nil {
		return err
	}
	if op.SaveField != "" {
		if file, _ := cmd.Flags().GetString("output-file"); file != "" {
			if err := saveResponseField(data, op.SaveField, file); err != nil {
				return err
			}
			emit(map[string]string{"file": file, "result": "written"}, func() {
				fmt.Printf("Wrote %s to %s (mode 0600)\n", op.SaveField, file)
			})
			return nil
		}
	}
	return printCatalogResponse(data, respHeader, secrets)
}

// saveResponseField writes one string field of a JSON response, such as a
// kubeconfig, to an owner-only file. The parent directory is created 0700.
func saveResponseField(data []byte, field, file string) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("the response is not a JSON object, nothing written to %s", file)
	}
	s, _ := m[field].(string)
	if s == "" {
		return fmt.Errorf("the response has no %q text, nothing written to %s", field, file)
	}
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
		return err
	}
	return writePrivateFile(file, []byte(s))
}

// credentialCache resolves credentials once per command.
type credentialCache struct {
	svc    *catalog.Service
	ctx    context.Context
	values map[string]string
}

var templateVar = regexp.MustCompile(`\{([^{}]*)\}`)

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
		return "", fmt.Errorf("a tenant ID is required: --tenant-id, NHN_CLOUD_TENANT_ID, or tenant_id in the profile")
	case "accessKeyId", "secretAccessKey":
		ak, sk := getAccessKey(), getSecretKey()
		if name == "accessKeyId" && ak != "" {
			return ak, nil
		}
		if name == "secretAccessKey" && sk != "" {
			return sk, nil
		}
		return "", fmt.Errorf("a User Access Key is required: access_key_id and secret_access_key in the profile, or NHN_CLOUD_ACCESS_KEY and NHN_CLOUD_SECRET_KEY")
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
			if appKey == "" {
				return "", fmt.Errorf("--appkey is empty")
			}
			return appKey, nil
		}
		if v := lookup("APPKEY", "app_key"); v != "" {
			return v, nil
		}
		if old, ok := legacyAppKeyNames[service]; ok {
			if v := os.Getenv(old[0]); v != "" {
				return v, nil
			}
			if v := cfg.Extra[old[1]]; v != "" {
				return v, nil
			}
		}
		if v := os.Getenv("NHN_CLOUD_APPKEY"); v != "" {
			return v, nil
		}
		if cfg.AppKey != "" {
			return cfg.AppKey, nil
		}
		where := fmt.Sprintf("%sAPPKEY, %s in the profile", serviceEnvPrefix(service), serviceConfigKey(service, "app_key"))
		if old, ok := legacyAppKeyNames[service]; ok {
			where += fmt.Sprintf(", %s, %s in the profile", old[0], old[1])
		}
		return "", fmt.Errorf("an appkey is required: --appkey, %s, NHN_CLOUD_APPKEY, or app_key in the profile", where)
	case "secretKey":
		if rootCmd.PersistentFlags().Changed("secret-key") && secretKeyFlag == "" {
			return "", fmt.Errorf("--secret-key is empty")
		}
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
	name := svc.OverrideName()
	override := os.Getenv(serviceEnvPrefix(name) + "ENDPOINT")
	if override == "" {
		override = LoadConfig().Extra[serviceConfigKey(name, "endpoint")]
	}
	if override != "" {
		u, err := url.Parse(override)
		if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || u.User != nil || strings.ContainsAny(override, "?#") || !catalog.PlainPath(u.Path) {
			return "", fmt.Errorf("%sENDPOINT must be scheme://host[:port][/path] with a plain path and no user info, query or fragment", serviceEnvPrefix(name))
		}
		return strings.TrimSuffix(override, "/"), nil
	}
	if svc.BaseURL == "{endpoint}" {
		return "", fmt.Errorf("%s has no public endpoint; set the one from the console in %sENDPOINT or %s in the profile",
			svc.Title, serviceEnvPrefix(name), serviceConfigKey(name, "endpoint"))
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
func catalogBody(cmd *cobra.Command, op *catalog.Operation) (io.Reader, string, []string, error) {
	if op.Body == nil {
		return nil, "", nil, nil
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
			return nil, "", nil, fmt.Errorf("read --data-file: %w", err)
		}
		return bytes.NewReader(data), op.Body.ContentType, nil, nil

	case "multipart":
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		for _, fd := range op.Body.Fields {
			if !flags.Changed(fd.FlagName()) {
				if fd.Required {
					return nil, "", nil, fmt.Errorf("--%s is required", fd.FlagName())
				}
				continue
			}
			if fd.Type == "array" {
				values, _ := flags.GetStringSlice(fd.FlagName())
				for _, v := range values {
					w.WriteField(fd.Name, v)
				}
				continue
			}
			value := flags.Lookup(fd.FlagName()).Value.String()
			if err := checkNumericField(fd, value); err != nil {
				return nil, "", nil, err
			}
			if fd.Type == "file" {
				data, err := os.ReadFile(value)
				if err != nil {
					return nil, "", nil, fmt.Errorf("read --%s: %w", fd.FlagName(), err)
				}
				part, err := w.CreateFormFile(fd.Name, filepath.Base(value))
				if err != nil {
					return nil, "", nil, err
				}
				part.Write(data)
				continue
			}
			w.WriteField(fd.Name, value)
		}
		if err := w.Close(); err != nil {
			return nil, "", nil, err
		}
		return &buf, w.FormDataContentType(), nil, nil

	case "form":
		form := url.Values{}
		for _, fd := range op.Body.Fields {
			if flags.Changed(fd.FlagName()) && fd.Type == "array" {
				values, _ := flags.GetStringSlice(fd.FlagName())
				for _, v := range values {
					form.Add(fd.Name, v)
				}
			} else if flags.Changed(fd.FlagName()) {
				value := flags.Lookup(fd.FlagName()).Value.String()
				if err := checkNumericField(fd, value); err != nil {
					return nil, "", nil, err
				}
				form.Set(fd.Name, value)
			} else if fd.Required {
				return nil, "", nil, fmt.Errorf("--%s is required", fd.FlagName())
			}
		}
		return strings.NewReader(form.Encode()), "application/x-www-form-urlencoded", nil, nil
	}

	obj := map[string]interface{}{}
	if raw, _ := flags.GetString("body"); raw != "" {
		data, err := readValueOrFile(raw)
		if err != nil {
			return nil, "", nil, fmt.Errorf("read --body: %w", err)
		}
		var parsed interface{}
		if err := decodeJSONNumbers(data, &parsed); err != nil {
			return nil, "", nil, fmt.Errorf("--body must be a JSON object")
		}
		m, ok := parsed.(map[string]interface{})
		if !ok || m == nil {
			return nil, "", nil, fmt.Errorf("--body must be a JSON object")
		}
		obj = m
		if w := op.Body.Wrapper; w != "" && len(obj) == 1 && obj[w] != nil && !hasField(op.Body.Fields, w) {
			return nil, "", nil, fmt.Errorf("--body is wrapped in %q already; pass the object inside it", w)
		}
	}
	for _, fd := range op.Body.Fields {
		if fd.Scalar() && flags.Changed(fd.FlagName()) {
			v, err := fieldValue(cmd, fd)
			if err != nil {
				return nil, "", nil, err
			}
			obj[fd.Name] = v
		}
		if _, ok := obj[fd.Name]; fd.Required && !ok {
			if fd.Scalar() {
				return nil, "", nil, fmt.Errorf("--%s is required (or set %q in --body)", fd.FlagName(), fd.Name)
			}
			return nil, "", nil, fmt.Errorf("%q is required in --body", fd.Name)
		}
	}
	var payload interface{} = obj
	if op.Body.Wrapper != "" {
		payload = map[string]interface{}{op.Body.Wrapper: obj}
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, "", nil, err
	}
	// Secret fields given in --body are hidden in failures like flag values.
	var secrets []string
	for _, fd := range op.Body.Fields {
		if v, ok := obj[fd.Name]; ok && v != nil && fd.Secret {
			secrets = append(secrets, fmt.Sprint(v)) // numbers too (PINs)
		}
	}
	return bytes.NewReader(data), "application/json", secrets, nil
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
		if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
			return nil, fmt.Errorf("--%s must be a number", name)
		}
		return n, nil
	}
	s, _ := flags.GetString(name)
	return s, nil
}

// catalogDo sends the request. shownURL is rawURL without query and with
// credentials and secrets in the path masked; only it appears in output.
func catalogDo(ctx context.Context, method, rawURL, shownURL string, headers http.Header, body io.Reader) ([]byte, http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %s: invalid request URL", method, shownURL)
	}
	req.Header = headers
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s %s\n", method, shownURL)
	}
	resp, err := iaasHTTPClient.Do(req)
	if err != nil {
		var ue *url.Error
		if errors.As(err, &ue) {
			return nil, nil, fmt.Errorf("%s %s: %w", method, shownURL, ue.Err)
		}
		return nil, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, nil, &apiError{Method: method, URL: shownURL, Status: resp.StatusCode, Body: stripURLQueries(string(data))}
	}
	return data, resp.Header, nil
}

func printCatalogResponse(data []byte, header http.Header, secrets []string) error {
	if len(bytes.TrimSpace(data)) == 0 {
		// Some APIs (e.g. HEAD on object storage) answer only in headers.
		h := map[string]string{}
		for k, v := range header {
			switch k {
			case "Date", "Content-Length", "Connection", "Set-Cookie":
				continue
			}
			if sensitiveHeaderName(k) {
				h[k] = "***"
			} else {
				h[k] = redact(strings.Join(v, ", "), secrets)
			}
		}
		emit(map[string]interface{}{"headers": h}, nil)
		return nil
	}
	var v interface{}
	if decodeJSONNumbers(data, &v) == nil {
		// Many NHN Cloud APIs answer 200 with {"header": {"isSuccessful": false}}.
		// A failure may echo the credentials sent; hide them before printing.
		if m, ok := v.(map[string]interface{}); ok {
			if h, ok := m["header"].(map[string]interface{}); ok && h["isSuccessful"] == false {
				// Redact decoded strings, so every JSON escape form is covered.
				v = redactValue(v, secrets)
				emit(v, nil)
				h, _ := v.(map[string]interface{})["header"].(map[string]interface{})
				return fmt.Errorf("the API reported a failure: resultCode %v, %s", h["resultCode"], stripURLQueries(fmt.Sprint(h["resultMessage"])))
			}
		}
		emit(v, nil)
		return nil
	}
	_, err := os.Stdout.Write(data)
	return err
}

// checkParamValue rejects path values that would leave their segment and
// numeric values that are not numbers.
func checkParamValue(p catalog.Param, value string) error {
	if p.In == "path" && (value == "" || value == "." || value == "..") {
		return fmt.Errorf("--%s: %q is not a valid path value", p.FlagName(), value)
	}
	switch p.Type {
	case "integer":
		if _, err := strconv.ParseInt(value, 10, 64); err != nil && p.Credential == "" {
			return fmt.Errorf("--%s must be an integer", p.FlagName())
		}
	case "number":
		if !finiteNumber(value) && p.Credential == "" {
			return fmt.Errorf("--%s must be a number", p.FlagName())
		}
	}
	return nil
}

// sensitiveHeaderName reports response headers whose values are not shown,
// e.g. X-Account-Meta-Temp-Url-Key or X-Subject-Token.
func sensitiveHeaderName(name string) bool {
	n := strings.ToLower(name)
	for _, part := range []string{"key", "token", "secret", "password", "auth", "cookie", "credential", "signature"} {
		if strings.Contains(n, part) {
			return true
		}
	}
	return false
}

// decodeJSONNumbers decodes JSON keeping integers exact: integers become
// int64 (or json.Number when larger), other numbers float64.
func decodeJSONNumbers(data []byte, out *interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(out); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return fmt.Errorf("unexpected data after JSON value")
	}
	*out = normalizeNumbers(*out)
	return nil
}

func normalizeNumbers(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		for k, e := range t {
			t[k] = normalizeNumbers(e)
		}
	case []interface{}:
		for i, e := range t {
			t[i] = normalizeNumbers(e)
		}
	case json.Number:
		if i, err := t.Int64(); err == nil {
			return i
		}
		if !strings.ContainsAny(string(t), ".eE") {
			return t // an integer too large for int64 stays exact
		}
		f, _ := t.Float64()
		return f
	}
	return v
}

func checkNumericField(fd catalog.Field, value string) error {
	switch fd.Type {
	case "integer":
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return fmt.Errorf("--%s must be an integer", fd.FlagName())
		}
	case "number":
		if !finiteNumber(value) {
			return fmt.Errorf("--%s must be a number", fd.FlagName())
		}
	}
	return nil
}

func finiteNumber(s string) bool {
	n, err := strconv.ParseFloat(s, 64)
	return err == nil && !math.IsNaN(n) && !math.IsInf(n, 0)
}

func hasField(fields []catalog.Field, name string) bool {
	for _, f := range fields {
		if f.Name == name {
			return true
		}
	}
	return false
}

var urlInText = regexp.MustCompile(`https?://[^\s"'<>]+`)

// stripURLQueries drops the query of URLs in a header value (signed URLs).
func stripURLQueries(v string) string {
	return urlInText.ReplaceAllStringFunc(v, stripQuery)
}

// parseExtraHeaders reads --header Name:Value flags.
func parseExtraHeaders(cmd *cobra.Command, op *catalog.Operation) ([][2]string, error) {
	var out [][2]string
	extra, _ := cmd.Flags().GetStringArray("header")
	for _, h := range extra {
		name, value, ok := strings.Cut(h, ":")
		name = strings.TrimSpace(name)
		if !ok || name == "" {
			return nil, fmt.Errorf("--header: want Name:Value") // the value may be a secret
		}
		if strings.EqualFold(name, "Content-Type") && op.Body != nil && op.Body.Encoding == "multipart" {
			return nil, fmt.Errorf("--header Content-Type would drop the multipart boundary")
		}
		out = append(out, [2]string{name, strings.TrimSpace(value)})
	}
	return out, nil
}

// usedSecrets lists the credentials and secret flag values of this call.
func usedSecrets(cmd *cobra.Command, creds *credentialCache) []string {
	var out []string
	for _, v := range creds.values {
		out = append(out, v)
	}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed && isSecretFlag(f) {
			out = append(out, f.Value.String())
		}
	})
	return out
}

// redact replaces each secret (4 characters or more) in s with ***.
// Shorter values are left alone: replacing them would mangle ordinary text.
func redact(s string, secrets []string) string {
	for _, v := range secrets {
		if len(v) < 4 {
			continue
		}
		s = strings.ReplaceAll(s, v, "***")
		// The value as sent in a path or query.
		for _, enc := range []string{url.PathEscape(v), url.QueryEscape(v)} {
			if enc != v {
				s = strings.ReplaceAll(s, enc, "***")
			}
		}
		// The value as it appears inside a JSON string (\u0026, \", ...).
		if b, err := json.Marshal(v); err == nil {
			if esc := string(b[1 : len(b)-1]); esc != v {
				s = strings.ReplaceAll(s, esc, "***")
			}
		}
	}
	return stripURLQueries(s)
}

// legacyAppKeyNames are the env and profile names the hand-written commands
// already read for a service's appkey, checked after the catalog names.
var legacyAppKeyNames = map[string][2]string{
	"rds-mysql":      {"NHN_CLOUD_MYSQL_APPKEY", "rds_app_key"},
	"rds-mariadb":    {"NHN_CLOUD_MARIADB_APPKEY", "rds_mariadb_app_key"},
	"rds-postgresql": {"NHN_CLOUD_POSTGRESQL_APPKEY", "rds_postgresql_app_key"},
	"key-manager":    {"NHN_CLOUD_KEYMANAGER_APPKEY", "keymanager_app_key"},
}

// serviceAppKey is the appkey a hand-written command of a service uses: the
// same lookup as the service's catalog commands (--appkey, the service's own
// env and profile names, then the generic appkey).
func serviceAppKey(service string) string {
	k, _ := resolveCredential(context.Background(), service, "appKey")
	return k
}

// redactValue hides secrets in every string of a decoded JSON value.
func redactValue(v interface{}, secrets []string) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, e := range t {
			out[redact(k, secrets)] = redactValue(e, secrets) // keys may echo secrets too
		}
		return out
	case []interface{}:
		for i, e := range t {
			t[i] = redactValue(e, secrets)
		}
	case string:
		return redact(t, secrets)
	case int64, float64, json.Number:
		for _, s := range secrets {
			if len(s) >= 4 && fmt.Sprint(t) == s {
				return "***"
			}
		}
	}
	return v
}

// redactBody redacts an error body: through its decoded strings when it is
// JSON, as text otherwise.
func redactBody(body string, secrets []string) string {
	var v interface{}
	if decodeJSONNumbers([]byte(body), &v) == nil {
		if b, err := json.Marshal(redactValue(v, secrets)); err == nil {
			return string(b)
		}
	}
	return redact(body, secrets)
}
