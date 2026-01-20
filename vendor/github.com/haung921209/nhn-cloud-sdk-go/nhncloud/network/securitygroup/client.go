package securitygroup

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
		return ErrNoCredentials
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

func (c *Client) ListSecurityGroups(ctx context.Context) (*ListSecurityGroupsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListSecurityGroupsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/security-groups", &out); err != nil {
		return nil, fmt.Errorf("list security groups: %w", err)
	}
	return &out, nil
}

func (c *Client) GetSecurityGroup(ctx context.Context, sgID string) (*GetSecurityGroupOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetSecurityGroupOutput
	if err := c.httpClient.GET(ctx, "/v2.0/security-groups/"+sgID, &out); err != nil {
		return nil, fmt.Errorf("get security group %s: %w", sgID, err)
	}
	return &out, nil
}

func (c *Client) CreateSecurityGroup(ctx context.Context, input *CreateSecurityGroupInput) (*CreateSecurityGroupOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"security_group": input}
	var out CreateSecurityGroupOutput
	if err := c.httpClient.POST(ctx, "/v2.0/security-groups", req, &out); err != nil {
		return nil, fmt.Errorf("create security group: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateSecurityGroup(ctx context.Context, sgID string, input *UpdateSecurityGroupInput) (*GetSecurityGroupOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"security_group": input}
	var out GetSecurityGroupOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/security-groups/"+sgID, req, &out); err != nil {
		return nil, fmt.Errorf("update security group %s: %w", sgID, err)
	}
	return &out, nil
}

func (c *Client) DeleteSecurityGroup(ctx context.Context, sgID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/security-groups/"+sgID, nil); err != nil {
		return fmt.Errorf("delete security group %s: %w", sgID, err)
	}
	return nil
}

func (c *Client) CreateRule(ctx context.Context, input *CreateRuleInput) (*CreateRuleOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"security_group_rule": input}
	var out CreateRuleOutput
	if err := c.httpClient.POST(ctx, "/v2.0/security-group-rules", req, &out); err != nil {
		return nil, fmt.Errorf("create security group rule: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteRule(ctx context.Context, ruleID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/security-group-rules/"+ruleID, nil); err != nil {
		return fmt.Errorf("delete security group rule %s: %w", ruleID, err)
	}
	return nil
}
