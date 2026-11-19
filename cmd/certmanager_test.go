package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// An API failure is reported with its result code, without the usage text,
// and not as "certificate not found".
func TestCertManagerReportsAPIFailureWithoutUsage(t *testing.T) {
	old := loadedConfig
	loadedConfig = &Config{}
	defer func() { loadedConfig = old }()
	t.Setenv("NHN_CLOUD_CERTMANAGER_APPKEY", "CMKEY")
	t.Setenv("NHN_CLOUD_ACCESS_KEY", "")
	t.Setenv("NHN_CLOUD_SECRET_KEY", "")
	useTransport(t, &recordTransport{body: `{"header":{"isSuccessful":false,"resultCode":30000,"resultMessage":"external api error"}}`})

	var errOut bytes.Buffer
	rootCmd.SetErr(&errOut)
	rootCmd.SetOut(&errOut)
	rootCmd.SetArgs([]string{"certmanager", "describe-certificate", "--certificate-name", "x"})
	defer func() { rootCmd.SetErr(nil); rootCmd.SetOut(nil); rootCmd.SetArgs(nil) }()
	err := rootCmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "30000") || strings.Contains(err.Error(), "not found") {
		t.Errorf("error %v, want the API result code", err)
	}
	if strings.Contains(errOut.String(), "Usage:") {
		t.Errorf("usage printed on an API error:\n%s", errOut.String())
	}
}
