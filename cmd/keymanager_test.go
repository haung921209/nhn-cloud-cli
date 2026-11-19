package cmd

import (
	"strings"
	"testing"
)

// --client-mac reaches the request header, for key stores with MAC
// address authentication.
func TestKeyManagerSendsClientMAC(t *testing.T) {
	old := loadedConfig
	loadedConfig = &Config{}
	defer func() { loadedConfig = old }()
	t.Setenv("NHN_CLOUD_KEY_MANAGER_APPKEY", "KMKEY")
	rec := &recordTransport{body: `{"header":{"isSuccessful":true},"body":{"clientMacHeader":"02:00:00:00:00:01"}}`}
	useTransport(t, rec)

	rootCmd.SetArgs([]string{"key-manager", "get-client-info", "--client-mac", "02:00:00:00:00:01", "-o", "json"})
	defer func() { rootCmd.SetArgs(nil); kmClientMAC = "" }()
	captureStdout(t, func() {
		if err := rootCmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})
	if rec.last == nil || strings.Join(rec.last.Header["X-TOAST-CLIENT-MAC-ADDR"], "") != "02:00:00:00:00:01" || rec.last.URL.Path != "/keymanager/v1.2/appkey/KMKEY/confirm" {
		t.Fatalf("request %v", rec.last)
	}
}
