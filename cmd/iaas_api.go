package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
	"time"
)

// iaasAPI calls NHN Cloud IaaS (OpenStack-style) APIs that the SDK does not
// cover yet. It issues an IaaS token with the configured API credentials and
// resolves the service endpoint from the token's service catalog, the same
// way the SDK clients do.
// ponytail: duplicates the SDK's internal identity client; move into the SDK
// once these APIs are added there.
type iaasAPI struct {
	token      string
	tenantID   string
	tenantName string
	catalog    []iaasCatalogService
}

type iaasCatalogService struct {
	Type      string `json:"type"`
	Endpoints []struct {
		PublicURL string `json:"publicURL"`
		Region    string `json:"region"`
	} `json:"endpoints"`
}

// defaultIdentityURL is the public identity endpoint. Other environments are
// reached through endpoint overrides (see endpoint_override.go).
const defaultIdentityURL = "https://api-identity-infrastructure.nhncloudservice.com"

var iaasHTTPClient = &http.Client{Timeout: 60 * time.Second}

var (
	cachedIaaSAPI *iaasAPI
	// cachedIaaSErr keeps a failed login, so one command never retries the
	// password more than once (repeated failures can lock the account).
	cachedIaaSErr error
)

func getIaaSAPI(ctx context.Context) (*iaasAPI, error) {
	if cachedIaaSAPI != nil || cachedIaaSErr != nil {
		return cachedIaaSAPI, cachedIaaSErr
	}
	tenant, user, pass := getTenantID(), getUsername(), getPassword()
	if tenant == "" || user == "" || pass == "" {
		return nil, fmt.Errorf("tenant_id, username and api_password are required (flags, NHN_CLOUD_* env, or ~/.nhncloud/credentials)")
	}
	body := map[string]interface{}{
		"auth": map[string]interface{}{
			"tenantId": tenant,
			"passwordCredentials": map[string]string{
				"username": user,
				"password": pass,
			},
		},
	}
	var resp struct {
		Access struct {
			Token struct {
				ID     string `json:"id"`
				Tenant struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"tenant"`
			} `json:"token"`
			ServiceCatalog []iaasCatalogService `json:"serviceCatalog"`
		} `json:"access"`
	}
	if err := doJSON(ctx, http.MethodPost, defaultIdentityURL+"/v2.0/tokens", "", body, &resp); err != nil {
		cachedIaaSErr = fmt.Errorf("issue IaaS token: %w", err)
		return nil, cachedIaaSErr
	}
	t := resp.Access.Token
	cachedIaaSAPI = &iaasAPI{token: t.ID, tenantID: t.Tenant.ID, tenantName: t.Tenant.Name, catalog: resp.Access.ServiceCatalog}
	return cachedIaaSAPI, nil
}

// endpoint returns the catalog URL for serviceType in the current region.
func (a *iaasAPI) endpoint(serviceType string) (string, error) {
	region := getRegion()
	for _, svc := range a.catalog {
		if svc.Type != serviceType {
			continue
		}
		for _, ep := range svc.Endpoints {
			if strings.EqualFold(ep.Region, region) {
				return strings.TrimSuffix(ep.PublicURL, "/"), nil
			}
		}
	}
	return "", fmt.Errorf("no %q endpoint for region %q in the service catalog", serviceType, region)
}

// iaasRequest calls path on the serviceType endpoint ("network", "compute", ...).
// body and out may be nil.
func iaasRequest(ctx context.Context, serviceType, method, path string, body, out interface{}) error {
	api, err := getIaaSAPI(ctx)
	if err != nil {
		return err
	}
	base, err := api.endpoint(serviceType)
	if err != nil {
		return err
	}
	return doJSON(ctx, method, base+path, api.token, body, out)
}

func networkRequest(ctx context.Context, method, path string, body, out interface{}) error {
	return iaasRequest(ctx, "network", method, path, body, out)
}

// apiError is returned for HTTP status >= 400. Body holds the raw response.
type apiError struct {
	Method, URL string
	Status      int
	Body        string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s %s: HTTP %d: %s", e.Method, e.URL, e.Status, strings.TrimSpace(e.Body))
}

func doJSON(ctx context.Context, method, url, token string, body, out interface{}) error {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		// Parse errors repeat the URL; keep the query out of errors.
		return fmt.Errorf("%s %s: invalid request URL", method, stripQuery(url))
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("X-Auth-Token", token)
	}
	if debug {
		// Never print the request body: the token request carries the password.
		fmt.Fprintf(os.Stderr, "[DEBUG] %s %s\n", method, url)
	}
	resp, err := iaasHTTPClient.Do(req)
	if err != nil {
		// *url.Error repeats the full URL; keep the query out of errors.
		var ue *neturl.Error
		if errors.As(err, &ue) {
			return fmt.Errorf("%s %s: %w", method, stripQuery(url), ue.Err)
		}
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		// The URL never carries credentials; drop the query to keep errors short.
		return &apiError{Method: method, URL: stripQuery(url), Status: resp.StatusCode, Body: string(data)}
	}
	if out == nil || len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	return json.Unmarshal(data, out)
}

// remarshal converts a generic JSON value into a typed one.
func remarshal(in, out interface{}) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

// stripQuery drops the query string so errors do not repeat filter values.
func stripQuery(u string) string { return strings.SplitN(u, "?", 2)[0] }
