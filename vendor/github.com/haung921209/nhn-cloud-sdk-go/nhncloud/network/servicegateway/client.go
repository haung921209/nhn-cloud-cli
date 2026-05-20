package servicegateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

// Client represents a Service Gateway service client
type Client struct {
	region        string
	credentials   credentials.IdentityCredentials
	httpClient    *client.Client
	tokenProvider *client.IdentityTokenProvider
	debug         bool
}

// NewClient creates a new Service Gateway client
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

// ================================
// Service Gateway Operations
// ================================

// ListServiceGateways retrieves a list of service gateways
func (c *Client) ListServiceGateways(ctx context.Context) (*ListServiceGatewaysOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListServiceGatewaysOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/servicegateways", &out); err != nil {
		return nil, fmt.Errorf("list service gateways: %w", err)
	}
	return &out, nil
}

// GetServiceGateway retrieves details of a specific service gateway
func (c *Client) GetServiceGateway(ctx context.Context, gatewayID string) (*GetServiceGatewayOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetServiceGatewayOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/servicegateways/"+gatewayID, &out); err != nil {
		return nil, fmt.Errorf("get service gateway %s: %w", gatewayID, err)
	}
	return &out, nil
}

// CreateServiceGateway creates a new service gateway
func (c *Client) CreateServiceGateway(ctx context.Context, input *CreateServiceGatewayInput) (*GetServiceGatewayOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateServiceGatewayRequest{ServiceGateway: *input}
	var out GetServiceGatewayOutput
	if err := c.httpClient.POST(ctx, "/v2.0/gateways/servicegateways", req, &out); err != nil {
		return nil, fmt.Errorf("create service gateway: %w", err)
	}
	return &out, nil
}

// UpdateServiceGateway updates an existing service gateway
func (c *Client) UpdateServiceGateway(ctx context.Context, gatewayID string, input *UpdateServiceGatewayInput) (*GetServiceGatewayOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &UpdateServiceGatewayRequest{ServiceGateway: *input}
	var out GetServiceGatewayOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/gateways/servicegateways/"+gatewayID, req, &out); err != nil {
		return nil, fmt.Errorf("update service gateway %s: %w", gatewayID, err)
	}
	return &out, nil
}

// DeleteServiceGateway deletes a service gateway
func (c *Client) DeleteServiceGateway(ctx context.Context, gatewayID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/gateways/servicegateways/"+gatewayID, nil); err != nil {
		return fmt.Errorf("delete service gateway %s: %w", gatewayID, err)
	}
	return nil
}

// ================================
// Service Endpoint Operations
// ================================

// ListServiceEndpoints retrieves a list of available service endpoints
func (c *Client) ListServiceEndpoints(ctx context.Context) (*ListServiceEndpointsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListServiceEndpointsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/serviceendpoints/", &out); err != nil {
		return nil, fmt.Errorf("list service endpoints: %w", err)
	}
	return &out, nil
}

// GetServiceEndpoint retrieves details of a specific service endpoint
func (c *Client) GetServiceEndpoint(ctx context.Context, endpointID string) (*GetServiceEndpointOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetServiceEndpointOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/serviceendpoints/"+endpointID, &out); err != nil {
		return nil, fmt.Errorf("get service endpoint %s: %w", endpointID, err)
	}
	return &out, nil
}
