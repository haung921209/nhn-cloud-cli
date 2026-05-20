// Package colocationgw provides Colocation Gateway service client
package colocationgw

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

// Client represents a Colocation Gateway API client
type Client struct {
	region        string
	credentials   credentials.IdentityCredentials
	httpClient    *client.Client
	tokenProvider *client.IdentityTokenProvider
	debug         bool
}

// NewClient creates a new Colocation Gateway client
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

	baseURL := fmt.Sprintf("https://%s-api-network-infrastructure.nhncloudservice.com", c.region)
	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}
	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)

	return nil
}

// List lists all colocation gateways
func (c *Client) List(ctx context.Context) (*ListOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result ListOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/colocationgateways", &result); err != nil {
		return nil, fmt.Errorf("list colocation gateways: %w", err)
	}

	return &result, nil
}

// Get gets a colocation gateway by ID
func (c *Client) Get(ctx context.Context, gatewayID string) (*GetOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v2.0/gateways/colocationgateways/%s", gatewayID)
	var result GetOutput
	if err := c.httpClient.GET(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get colocation gateway %s: %w", gatewayID, err)
	}

	return &result, nil
}
