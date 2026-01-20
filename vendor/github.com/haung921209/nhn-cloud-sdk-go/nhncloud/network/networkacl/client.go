package networkacl

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

type Client struct {
	region        string
	credentials   credentials.IdentityCredentials
	httpClient    *client.Client
	tokenProvider *client.IdentityTokenProvider
	debug         bool
}

func NewClient(region string, creds credentials.IdentityCredentials, hc *http.Client, debug bool) *Client {
	c := &Client{
		region:      region,
		credentials: creds,
		debug:       debug,
	}

	if creds != nil {
		c.tokenProvider = client.NewIdentityTokenProvider(
			creds.GetTenantID(),
			creds.GetUsername(),
			creds.GetPassword(),
		)
	}

	return c
}

func (c *Client) ensureClient(ctx context.Context) error {
	if c.httpClient != nil {
		return nil
	}

	if c.tokenProvider == nil {
		return fmt.Errorf("no credentials provided")
	}

	if _, err := c.tokenProvider.GetToken(ctx); err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}

	baseURL, err := c.tokenProvider.GetServiceEndpoint("network", c.region)
	if err != nil {
		return fmt.Errorf("resolve endpoint: %w", err)
	}

	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}

	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
	return nil
}

// ============== ACL Management APIs ==============

func (c *Client) ListACLs(ctx context.Context) (*ListACLsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListACLsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/acls", &out); err != nil {
		return nil, fmt.Errorf("list ACLs: %w", err)
	}
	return &out, nil
}

func (c *Client) GetACL(ctx context.Context, id string) (*GetACLOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetACLOutput
	if err := c.httpClient.GET(ctx, "/v2.0/acls/"+id, &out); err != nil {
		return nil, fmt.Errorf("get ACL %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateACL(ctx context.Context, input *CreateACLInput) (*GetACLOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateACLRequest{ACL: *input}
	var out GetACLOutput
	if err := c.httpClient.POST(ctx, "/v2.0/acls", req, &out); err != nil {
		return nil, fmt.Errorf("create ACL: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateACL(ctx context.Context, id string, input *UpdateACLInput) (*GetACLOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &UpdateACLRequest{ACL: *input}
	var out GetACLOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/acls/"+id, req, &out); err != nil {
		return nil, fmt.Errorf("update ACL %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) DeleteACL(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/acls/"+id, nil); err != nil {
		return fmt.Errorf("delete ACL %s: %w", id, err)
	}
	return nil
}

// ============== ACL Rule Management APIs ==============

func (c *Client) ListACLRules(ctx context.Context) (*ListACLRulesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListACLRulesOutput
	if err := c.httpClient.GET(ctx, "/v2.0/acl_rules", &out); err != nil {
		return nil, fmt.Errorf("list ACL rules: %w", err)
	}
	return &out, nil
}

func (c *Client) ListACLRulesByACL(ctx context.Context, aclID string) (*ListACLRulesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListACLRulesOutput
	if err := c.httpClient.GET(ctx, "/v2.0/acl_rules?acl_id="+aclID, &out); err != nil {
		return nil, fmt.Errorf("list ACL rules for ACL %s: %w", aclID, err)
	}
	return &out, nil
}

func (c *Client) GetACLRule(ctx context.Context, id string) (*GetACLRuleOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetACLRuleOutput
	if err := c.httpClient.GET(ctx, "/v2.0/acl_rules/"+id, &out); err != nil {
		return nil, fmt.Errorf("get ACL rule %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateACLRule(ctx context.Context, input *CreateACLRuleInput) (*GetACLRuleOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateACLRuleRequest{ACLRule: *input}
	var out GetACLRuleOutput
	if err := c.httpClient.POST(ctx, "/v2.0/acl_rules", req, &out); err != nil {
		return nil, fmt.Errorf("create ACL rule: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateACLRule(ctx context.Context, id string, input *UpdateACLRuleInput) (*GetACLRuleOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &UpdateACLRuleRequest{ACLRule: *input}
	var out GetACLRuleOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/acl_rules/"+id, req, &out); err != nil {
		return nil, fmt.Errorf("update ACL rule %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) DeleteACLRule(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/acl_rules/"+id, nil); err != nil {
		return fmt.Errorf("delete ACL rule %s: %w", id, err)
	}
	return nil
}

// ============== ACL Binding Management APIs ==============

func (c *Client) ListACLBindings(ctx context.Context) (*ListACLBindingsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListACLBindingsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/acl_bindings", &out); err != nil {
		return nil, fmt.Errorf("list ACL bindings: %w", err)
	}
	return &out, nil
}

func (c *Client) ListACLBindingsByACL(ctx context.Context, aclID string) (*ListACLBindingsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListACLBindingsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/acl_bindings?acl_id="+aclID, &out); err != nil {
		return nil, fmt.Errorf("list ACL bindings for ACL %s: %w", aclID, err)
	}
	return &out, nil
}

func (c *Client) GetACLBinding(ctx context.Context, id string) (*GetACLBindingOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetACLBindingOutput
	if err := c.httpClient.GET(ctx, "/v2.0/acl_bindings/"+id, &out); err != nil {
		return nil, fmt.Errorf("get ACL binding %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateACLBinding(ctx context.Context, input *CreateACLBindingInput) (*GetACLBindingOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateACLBindingRequest{ACLBinding: *input}
	var out GetACLBindingOutput
	if err := c.httpClient.POST(ctx, "/v2.0/acl_bindings", req, &out); err != nil {
		return nil, fmt.Errorf("create ACL binding: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteACLBinding(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/acl_bindings/"+id, nil); err != nil {
		return fmt.Errorf("delete ACL binding %s: %w", id, err)
	}
	return nil
}
