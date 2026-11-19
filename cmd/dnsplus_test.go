package cmd

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// recordTransport answers every request with body and records the paths.
type recordTransport struct {
	body  string
	paths []string
	last  *http.Request
}

func (r *recordTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r.paths = append(r.paths, req.URL.Path)
	r.last = req
	return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(r.body)), Request: req}, nil
}

// useTransport sends the default HTTP client's requests to rt for this test.
func useTransport(t *testing.T, rt http.RoundTripper) {
	old := http.DefaultTransport
	http.DefaultTransport = rt
	t.Cleanup(func() { http.DefaultTransport = old })
}

// The DNS Plus client must put the DNS Plus appkey in the path (it used to
// put the region there).
func TestDNSPlusClientUsesServiceAppKey(t *testing.T) {
	rec := &recordTransport{}
	useTransport(t, rec)
	t.Setenv("NHN_CLOUD_DNSPLUS_APPKEY", "DNSKEY")
	t.Setenv("NHN_CLOUD_REGION", "KR1")

	newDNSPlusClient().ListZones(context.Background())
	if len(rec.paths) != 1 || !strings.Contains(rec.paths[0], "/appkeys/DNSKEY/") {
		t.Fatalf("request paths %v, want the DNS Plus appkey in the path", rec.paths)
	}
}
