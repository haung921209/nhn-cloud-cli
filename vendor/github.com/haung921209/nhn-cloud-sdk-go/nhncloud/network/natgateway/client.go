package natgateway

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

func (c *Client) ListNATGateways(ctx context.Context) (*ListNATGatewaysOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListNATGatewaysOutput
	if err := c.httpClient.GET(ctx, "/v2.0/natgateways", &out); err != nil {
		return nil, fmt.Errorf("list NAT gateways: %w", err)
	}
	return &out, nil
}

func (c *Client) GetNATGateway(ctx context.Context, id string) (*GetNATGatewayOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetNATGatewayOutput
	if err := c.httpClient.GET(ctx, "/v2.0/natgateways/"+id, &out); err != nil {
		return nil, fmt.Errorf("get NAT gateway %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateNATGateway(ctx context.Context, input *CreateNATGatewayInput) (*GetNATGatewayOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateNATGatewayRequest{NATGateway: *input}
	var out GetNATGatewayOutput
	if err := c.httpClient.POST(ctx, "/v2.0/natgateways", req, &out); err != nil {
		return nil, fmt.Errorf("create NAT gateway: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateNATGateway(ctx context.Context, id string, input *UpdateNATGatewayInput) (*GetNATGatewayOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &UpdateNATGatewayRequest{NATGateway: *input}
	var out GetNATGatewayOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/natgateways/"+id, req, &out); err != nil {
		return nil, fmt.Errorf("update NAT gateway %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) DeleteNATGateway(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/natgateways/"+id, nil); err != nil {
		return fmt.Errorf("delete NAT gateway %s: %w", id, err)
	}
	return nil
}
