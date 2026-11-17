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
		`unknown credential "app_key"`: func(s *Service) {
			s.Auth.Headers = map[string]string{"Authorization": "Bearer {app_key}"}
		},
		"GET with a body": func(s *Service) { s.Operations[0].Body = &Body{} },
		"exactly one of":  func(s *Service) { s.IaaSServiceType = "network" },
		"needs multipart": func(s *Service) {
			s.Operations[0].Method = "POST"
			s.Operations[0].Body = &Body{Fields: []Field{{Name: "f", Type: "file"}}}
		},
		"must be required":     func(s *Service) { s.Operations[0].Params[0].Required = false },
		"empty path variable":  func(s *Service) { s.Operations[0].Path = "/v1/x/{id}/{}" },
		"no query or fragment": func(s *Service) { s.Operations[0].Path = "/v1/x/{id}#tail" },
		"is not in the path": func(s *Service) {
			s.Operations[0].Params = append(s.Operations[0].Params, Param{Name: "other", In: "path", Required: true})
		},
		"has bad type":       func(s *Service) { s.Operations[0].Params[0].Type = "bool" },
		"must be lower case": func(s *Service) { s.BaseURL = ""; s.RegionBaseURLs = map[string]string{"KR1": "https://a.example.com"} },
		"is reserved":        func(s *Service) { s.Name = "api" },
		"only for raw bodies": func(s *Service) {
			s.Operations[0].Method = "POST"
			s.Operations[0].Body = &Body{ContentType: "text/plain"}
		},
		"only be sent in a json body": func(s *Service) {
			s.Operations[0].Method = "POST"
			s.Operations[0].Body = &Body{Encoding: "form", Fields: []Field{{Name: "o", Type: "object", Required: true}}}
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

func TestValidateAllCatchesSameTitleDuplicates(t *testing.T) {
	if errs := ValidateAll([]Service{base(), base()}); len(errs) == 0 {
		t.Error("two files defining the same operation must be rejected even with the same title")
	}
}

func TestDecodeRejectsTrailingData(t *testing.T) {
	if _, err := Decode([]byte(`{"name":"a"} {"name":"b"}`)); err == nil {
		t.Error("trailing data must be rejected")
	}
	if _, err := Decode([]byte("{\"name\":\"a\"}\n")); err != nil {
		t.Errorf("a trailing newline is fine: %v", err)
	}
}

func TestFlagName(t *testing.T) {
	for in, want := range map[string]string{"floating_network_id": "floating-network-id", "pageSize": "page-size", "X-Is-Resource-Name": "x-is-resource-name"} {
		if got := FlagName(in); got != want {
			t.Errorf("%s: %s", in, got)
		}
	}
}

func TestValidateAllNameAndEndpointClashes(t *testing.T) {
	a := base()
	a.Aliases = []string{"other"}
	b := base()
	b.Name, b.Title = "other", "B"
	b.Operations[0].Name = "get-y"
	if errs := ValidateAll([]Service{a, b}); !strings.Contains(strings.Join(errs, "\n"), `alias "other"`) {
		t.Errorf("alias clash: %v", errs)
	}

	c := base()
	d := base()
	d.Title = "D"
	d.Operations[0].Name = "get-z"
	d.BaseURL = "https://other.example.com"
	if errs := ValidateAll([]Service{c, d}); !strings.Contains(strings.Join(errs, "\n"), "endpointName") {
		t.Errorf("endpoint clash: %v", errs)
	}
	d.EndpointName = "svc-other"
	if errs := ValidateAll([]Service{c, d}); len(errs) != 0 {
		t.Errorf("distinct endpointName: %v", errs)
	}
}
