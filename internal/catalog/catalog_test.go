package catalog

import (
	"strings"
	"testing"
)

func base() Service {
	return Service{
		Name: "svc", Title: "T", Docs: "https://docs.nhncloud.com/ko/x/", BaseURL: "https://api.example.com",
		Operations: []Operation{{Name: "get-x", Summary: "s", Method: "GET", Path: "/v1/x/{id}",
			Params: []Param{{Name: "id", In: "path", Required: true}}}},
	}
}

func TestValidateAcceptsAMinimalService(t *testing.T) {
	if errs := Validate(base()); len(errs) != 0 {
		t.Fatal(errs)
	}
}

func TestValidateRejects(t *testing.T) {
	cases := map[string]func(*Service){
		"must not end with /": func(s *Service) { s.Operations[0].Path = "/v1/x/{id}/" },
		"has no path param":   func(s *Service) { s.Operations[0].Params = nil },
		"clashes with a global flag": func(s *Service) {
			s.Operations[0].Params = append(s.Operations[0].Params, Param{Name: "region", In: "query"})
		},
		"unknown credential": func(s *Service) { s.Auth.Headers = map[string]string{"X": "{nope}"} },
		"GET with a body":    func(s *Service) { s.Operations[0].Body = &Body{} },
		"exactly one of":     func(s *Service) { s.IaaSServiceType = "network" },
		"needs multipart": func(s *Service) {
			s.Operations[0].Method = "POST"
			s.Operations[0].Body = &Body{Fields: []Field{{Name: "f", Type: "file"}}}
		},
	}
	for want, mutate := range cases {
		s := base()
		mutate(&s)
		errs := strings.Join(Validate(s), "\n")
		if !strings.Contains(errs, want) {
			t.Errorf("%s: got %q", want, errs)
		}
	}
}

func TestFlagName(t *testing.T) {
	for in, want := range map[string]string{"floating_network_id": "floating-network-id", "pageSize": "page-size", "X-Is-Resource-Name": "x-is-resource-name"} {
		if got := FlagName(in); got != want {
			t.Errorf("%s: %s", in, got)
		}
	}
}
