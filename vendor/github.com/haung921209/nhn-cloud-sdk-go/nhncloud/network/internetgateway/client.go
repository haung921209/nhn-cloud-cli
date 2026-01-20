package internetgateway

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

func (c *Client) ListInternetGateways(ctx context.Context) (*ListInternetGatewaysOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListInternetGatewaysOutput
	if err := c.httpClient.GET(ctx, "/v2.0/internetgateways", &out); err != nil {
		return nil, fmt.Errorf("list internet gateways: %w", err)
	}
	return &out, nil
}

func (c *Client) GetInternetGateway(ctx context.Context, id string) (*GetInternetGatewayOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetInternetGatewayOutput
	if err := c.httpClient.GET(ctx, "/v2.0/internetgateways/"+id, &out); err != nil {
		return nil, fmt.Errorf("get internet gateway %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateInternetGateway(ctx context.Context, input *CreateInternetGatewayInput) (*GetInternetGatewayOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateInternetGatewayRequest{InternetGateway: *input}
	var out GetInternetGatewayOutput
	if err := c.httpClient.POST(ctx, "/v2.0/internetgateways", req, &out); err != nil {
		return nil, fmt.Errorf("create internet gateway: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteInternetGateway(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/internetgateways/"+id, nil); err != nil {
		return fmt.Errorf("delete internet gateway %s: %w", id, err)
	}
	return nil
}

func (c *Client) ListExternalNetworks(ctx context.Context) (*ListExternalNetworksOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListExternalNetworksOutput
	if err := c.httpClient.GET(ctx, "/v2.0/vpcs?router:external=true", &out); err != nil {
		return nil, fmt.Errorf("list external networks: %w", err)
	}
	return &out, nil
}
