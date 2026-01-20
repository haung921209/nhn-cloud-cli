package privatedns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

// Client represents a Private DNS service client
type Client struct {
	region        string
	credentials   credentials.IdentityCredentials
	httpClient    *client.Client
	tokenProvider *client.IdentityTokenProvider
	debug         bool
}

// NewClient creates a new Private DNS client
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
// Zone Operations
// ================================

// ListZones lists all private DNS zones
func (c *Client) ListZones(ctx context.Context) (*ListZonesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListZonesOutput
	if err := c.httpClient.GET(ctx, "/v2.0/privatedns/zones", &out); err != nil {
		return nil, fmt.Errorf("list zones: %w", err)
	}
	return &out, nil
}

// GetZone gets a private DNS zone by ID
func (c *Client) GetZone(ctx context.Context, zoneID string) (*GetZoneOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetZoneOutput
	if err := c.httpClient.GET(ctx, "/v2.0/privatedns/zones/"+zoneID, &out); err != nil {
		return nil, fmt.Errorf("get zone %s: %w", zoneID, err)
	}
	return &out, nil
}

// CreateZone creates a new private DNS zone
func (c *Client) CreateZone(ctx context.Context, input *CreateZoneInput) (*GetZoneOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateZoneRequest{Zone: *input}
	var out GetZoneOutput
	if err := c.httpClient.POST(ctx, "/v2.0/privatedns/zones", req, &out); err != nil {
		return nil, fmt.Errorf("create zone: %w", err)
	}
	return &out, nil
}

// UpdateZone updates a private DNS zone
func (c *Client) UpdateZone(ctx context.Context, zoneID string, input *UpdateZoneInput) (*GetZoneOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &UpdateZoneRequest{Zone: *input}
	var out GetZoneOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/privatedns/zones/"+zoneID, req, &out); err != nil {
		return nil, fmt.Errorf("update zone %s: %w", zoneID, err)
	}
	return &out, nil
}

// DeleteZone deletes a private DNS zone
func (c *Client) DeleteZone(ctx context.Context, zoneID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/privatedns/zones/"+zoneID, nil); err != nil {
		return fmt.Errorf("delete zone %s: %w", zoneID, err)
	}
	return nil
}

// ================================
// Record Set Operations
// ================================

// ListRRSets lists record sets in a zone
func (c *Client) ListRRSets(ctx context.Context, zoneID string) (*ListRRSetsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListRRSetsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/privatedns/zones/"+zoneID+"/rrsets", &out); err != nil {
		return nil, fmt.Errorf("list record sets: %w", err)
	}
	return &out, nil
}

// GetRRSet gets a record set by ID
func (c *Client) GetRRSet(ctx context.Context, zoneID, rrsetID string) (*GetRRSetOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetRRSetOutput
	if err := c.httpClient.GET(ctx, "/v2.0/privatedns/zones/"+zoneID+"/rrsets/"+rrsetID, &out); err != nil {
		return nil, fmt.Errorf("get record set %s: %w", rrsetID, err)
	}
	return &out, nil
}

// CreateRRSet creates a new record set
func (c *Client) CreateRRSet(ctx context.Context, zoneID string, input *CreateRRSetInput) (*GetRRSetOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateRRSetRequest{RRSet: *input}
	var out GetRRSetOutput
	if err := c.httpClient.POST(ctx, "/v2.0/privatedns/zones/"+zoneID+"/rrsets", req, &out); err != nil {
		return nil, fmt.Errorf("create record set: %w", err)
	}
	return &out, nil
}

// UpdateRRSet updates a record set
func (c *Client) UpdateRRSet(ctx context.Context, zoneID, rrsetID string, input *UpdateRRSetInput) (*GetRRSetOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &UpdateRRSetRequest{RRSet: *input}
	var out GetRRSetOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/privatedns/zones/"+zoneID+"/rrsets/"+rrsetID, req, &out); err != nil {
		return nil, fmt.Errorf("update record set %s: %w", rrsetID, err)
	}
	return &out, nil
}

// DeleteRRSet deletes a record set
func (c *Client) DeleteRRSet(ctx context.Context, zoneID, rrsetID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/privatedns/zones/"+zoneID+"/rrsets/"+rrsetID, nil); err != nil {
		return fmt.Errorf("delete record set %s: %w", rrsetID, err)
	}
	return nil
}
