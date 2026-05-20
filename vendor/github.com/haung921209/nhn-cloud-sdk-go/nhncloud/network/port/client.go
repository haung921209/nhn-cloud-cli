package port

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

func (c *Client) ListPorts(ctx context.Context) (*ListPortsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListPortsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/ports", &out); err != nil {
		return nil, fmt.Errorf("list ports: %w", err)
	}
	return &out, nil
}

func (c *Client) ListPortsByDevice(ctx context.Context, deviceID string) (*ListPortsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListPortsOutput
	path := fmt.Sprintf("/v2.0/ports?device_id=%s", deviceID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list ports by device: %w", err)
	}
	return &out, nil
}

func (c *Client) GetPort(ctx context.Context, portID string) (*GetPortOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetPortOutput
	if err := c.httpClient.GET(ctx, "/v2.0/ports/"+portID, &out); err != nil {
		return nil, fmt.Errorf("get port %s: %w", portID, err)
	}
	return &out, nil
}
