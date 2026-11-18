package cmd

import "testing"

// Hand-written commands must use the service's own appkey before the generic
// one, the same lookup as the service's catalog commands.
func TestServiceAppKeyPrefersServiceNames(t *testing.T) {
	old := loadedConfig
	defer func() { loadedConfig = old }()
	t.Setenv("NHN_CLOUD_APPKEY", "GLOBAL")

	loadedConfig = &Config{AppKey: "PROFILEGLOBAL", Extra: map[string]string{"dnsplus_app_key": "DNSKEY", "keymanager_app_key": "KMLEGACY"}}
	t.Setenv("NHN_CLOUD_APIGW_APPKEY", "APIGWKEY")
	cases := map[string]string{"apigw": "APIGWKEY", "dnsplus": "DNSKEY", "key-manager": "KMLEGACY", "certmanager": "GLOBAL"}
	for svc, want := range cases {
		if got := serviceAppKey(svc); got != want {
			t.Errorf("%s: got %q, want %q", svc, got, want)
		}
	}
}
