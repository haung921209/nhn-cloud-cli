package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Endpoint overrides send requests for one API host to another host, so the
// CLI can talk to a non-public environment (for example a staging identity
// and API set) without code changes. Every SDK client and the CLI's own API
// calls go through http.DefaultTransport, which is wrapped here.
//
// Format: "<from-host>=<to-host or scheme://host[:port]>", comma separated.
// Sources (first non-empty wins): --endpoint-override flags,
// NHN_CLOUD_ENDPOINT_OVERRIDES, endpoint_overrides in the selected profile.

var endpointOverrideFlags []string

func init() {
	rootCmd.PersistentFlags().StringArrayVar(&endpointOverrideFlags, "endpoint-override", nil,
		"Send requests for one API host to another: from-host=to-host (repeatable; also NHN_CLOUD_ENDPOINT_OVERRIDES or endpoint_overrides in the profile)")
	cobra.OnInitialize(installEndpointOverrides)
}

func installEndpointOverrides() {
	raw := strings.Join(endpointOverrideFlags, ",")
	if raw == "" {
		raw = os.Getenv("NHN_CLOUD_ENDPOINT_OVERRIDES")
	}
	if raw == "" {
		raw = LoadConfig().EndpointOverrides
	}
	if raw == "" {
		return
	}
	m, err := parseEndpointOverrides(raw)
	if err != nil {
		exitWithError("Invalid endpoint override", err)
	}
	http.DefaultTransport = &hostRewriteTransport{inner: http.DefaultTransport, hosts: m}
}

// parseEndpointOverrides maps a lower-case source host to the target URL
// (scheme and host only).
func parseEndpointOverrides(raw string) (map[string]*url.URL, error) {
	m := map[string]*url.URL{}
	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		from, to, ok := strings.Cut(pair, "=")
		from, to = strings.ToLower(strings.TrimSpace(from)), strings.TrimSpace(to)
		if !ok || from == "" || to == "" {
			return nil, fmt.Errorf("%q: want from-host=to-host", pair)
		}
		if !strings.Contains(to, "://") {
			to = "https://" + to
		}
		u, err := url.Parse(to)
		if err != nil || u.Host == "" || (u.Path != "" && u.Path != "/") {
			return nil, fmt.Errorf("%q: target must be a host or scheme://host[:port] without a path", pair)
		}
		m[from] = &url.URL{Scheme: u.Scheme, Host: u.Host}
	}
	return m, nil
}

type hostRewriteTransport struct {
	inner http.RoundTripper
	hosts map[string]*url.URL
}

func (t *hostRewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	target, ok := t.hosts[strings.ToLower(req.URL.Host)]
	if !ok {
		return t.inner.RoundTrip(req)
	}
	r := req.Clone(req.Context())
	r.URL.Scheme = target.Scheme
	r.URL.Host = target.Host
	r.Host = ""
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] endpoint override: %s -> %s\n", req.URL.Host, target.Host)
	}
	return t.inner.RoundTrip(r)
}
