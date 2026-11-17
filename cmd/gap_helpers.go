package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// splitCSV splits a comma separated flag value, dropping empty items.
func splitCSV(s string) []string {
	var out []string
	for _, v := range strings.Split(s, ",") {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// readValueOrFile returns v, or the contents of the file when v is "@path".
func readValueOrFile(v string) ([]byte, error) {
	if strings.HasPrefix(v, "@") {
		return os.ReadFile(v[1:])
	}
	return []byte(v), nil
}

// instancePortID returns the only port of an instance; it exits when the
// instance has no port or several ports.
func instancePortID(instanceID string) string {
	var res struct {
		Ports []netPort `json:"ports"`
	}
	apiCall(http.MethodGet, "/v2.0/ports?device_id="+url.QueryEscape(instanceID), nil, &res)
	switch len(res.Ports) {
	case 1:
		return res.Ports[0].ID
	case 0:
		exitWithError(fmt.Sprintf("instance %s has no port", instanceID), nil)
	default:
		exitWithError(fmt.Sprintf("instance %s has %d ports; pick one with --port-id (nhncloud network describe-ports --device-id %s)",
			instanceID, len(res.Ports), instanceID), nil)
	}
	return ""
}

// keyTypeHint explains a failed key pair import for key types the compute
// API is known to reject.
func keyTypeHint(publicKey string) string {
	if strings.HasPrefix(strings.TrimSpace(publicKey), "ssh-ed25519") {
		return " (the compute API rejects ssh-ed25519 public keys with 400; import an RSA key instead, e.g. ssh-keygen -t rsa -b 4096)"
	}
	return ""
}
