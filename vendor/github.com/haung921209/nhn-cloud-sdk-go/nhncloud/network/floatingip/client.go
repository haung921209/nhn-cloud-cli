package floatingip

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

func (c *Client) ListFloatingIPs(ctx context.Context) (*ListFloatingIPsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListFloatingIPsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/floatingips", &out); err != nil {
		return nil, fmt.Errorf("list floating IPs: %w", err)
	}
	return &out, nil
}

func (c *Client) GetFloatingIP(ctx context.Context, floatingIPID string) (*GetFloatingIPOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetFloatingIPOutput
	if err := c.httpClient.GET(ctx, "/v2.0/floatingips/"+floatingIPID, &out); err != nil {
		return nil, fmt.Errorf("get floating IP %s: %w", floatingIPID, err)
	}
	return &out, nil
}

func (c *Client) CreateFloatingIP(ctx context.Context, input *CreateFloatingIPInput) (*CreateFloatingIPOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateFloatingIPRequest{FloatingIP: *input}
	var out CreateFloatingIPOutput
	if err := c.httpClient.POST(ctx, "/v2.0/floatingips", req, &out); err != nil {
		return nil, fmt.Errorf("create floating IP: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateFloatingIP(ctx context.Context, floatingIPID string, input *UpdateFloatingIPInput) (*UpdateFloatingIPOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &UpdateFloatingIPRequest{FloatingIP: *input}
	var out UpdateFloatingIPOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/floatingips/"+floatingIPID, req, &out); err != nil {
		return nil, fmt.Errorf("update floating IP %s: %w", floatingIPID, err)
	}
	return &out, nil
}

func (c *Client) DeleteFloatingIP(ctx context.Context, floatingIPID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/floatingips/"+floatingIPID, nil); err != nil {
		return fmt.Errorf("delete floating IP %s: %w", floatingIPID, err)
	}
	return nil
}

func (c *Client) AssociateFloatingIP(ctx context.Context, floatingIPID string, portID string) (*UpdateFloatingIPOutput, error) {
	input := &UpdateFloatingIPInput{PortID: &portID}
	return c.UpdateFloatingIP(ctx, floatingIPID, input)
}

func (c *Client) DisassociateFloatingIP(ctx context.Context, floatingIPID string) (*UpdateFloatingIPOutput, error) {
	input := &UpdateFloatingIPInput{PortID: nil}
	return c.UpdateFloatingIP(ctx, floatingIPID, input)
}
