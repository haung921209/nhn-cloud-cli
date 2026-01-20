// Package mirroring provides Traffic Mirroring service client
package mirroring

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

// Client represents a Traffic Mirroring API client
type Client struct {
	region        string
	credentials   credentials.IdentityCredentials
	httpClient    *client.Client
	tokenProvider *client.IdentityTokenProvider
	debug         bool
}

// NewClient creates a new Traffic Mirroring client
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

// ================================
// Session Operations
// ================================

// ListSessions lists all mirroring sessions
func (c *Client) ListSessions(ctx context.Context) (*ListSessionsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result ListSessionsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/mirroring/sessions", &result); err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	return &result, nil
}

// GetSession gets a session by ID
func (c *Client) GetSession(ctx context.Context, sessionID string) (*SessionOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result SessionOutput
	if err := c.httpClient.GET(ctx, "/v2.0/mirroring/sessions/"+sessionID, &result); err != nil {
		return nil, fmt.Errorf("get session %s: %w", sessionID, err)
	}
	return &result, nil
}

// CreateSession creates a new session
func (c *Client) CreateSession(ctx context.Context, input *CreateSessionInput) (*SessionOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"mirroring_session": input}
	var result SessionOutput
	if err := c.httpClient.POST(ctx, "/v2.0/mirroring/sessions", req, &result); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return &result, nil
}

// UpdateSession updates a session
func (c *Client) UpdateSession(ctx context.Context, sessionID string, input *UpdateSessionInput) (*SessionOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"mirroring_session": input}
	var result SessionOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/mirroring/sessions/"+sessionID, req, &result); err != nil {
		return nil, fmt.Errorf("update session %s: %w", sessionID, err)
	}
	return &result, nil
}

// DeleteSession deletes a session
func (c *Client) DeleteSession(ctx context.Context, sessionID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/mirroring/sessions/"+sessionID, nil); err != nil {
		return fmt.Errorf("delete session %s: %w", sessionID, err)
	}
	return nil
}

// ================================
// Filter Group Operations
// ================================

// ListFilterGroups lists all filter groups
func (c *Client) ListFilterGroups(ctx context.Context) (*ListFilterGroupsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result ListFilterGroupsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/mirroring/filtergroups", &result); err != nil {
		return nil, fmt.Errorf("list filter groups: %w", err)
	}
	return &result, nil
}

// GetFilterGroup gets a filter group by ID
func (c *Client) GetFilterGroup(ctx context.Context, filterGroupID string) (*FilterGroupOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result FilterGroupOutput
	if err := c.httpClient.GET(ctx, "/v2.0/mirroring/filtergroups/"+filterGroupID, &result); err != nil {
		return nil, fmt.Errorf("get filter group %s: %w", filterGroupID, err)
	}
	return &result, nil
}

// CreateFilterGroup creates a new filter group
func (c *Client) CreateFilterGroup(ctx context.Context, input *CreateFilterGroupInput) (*FilterGroupOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"mirroring_filtergroup": input}
	var result FilterGroupOutput
	if err := c.httpClient.POST(ctx, "/v2.0/mirroring/filtergroups", req, &result); err != nil {
		return nil, fmt.Errorf("create filter group: %w", err)
	}
	return &result, nil
}

// UpdateFilterGroup updates a filter group
func (c *Client) UpdateFilterGroup(ctx context.Context, filterGroupID string, input *UpdateFilterGroupInput) (*FilterGroupOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"mirroring_filtergroup": input}
	var result FilterGroupOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/mirroring/filtergroups/"+filterGroupID, req, &result); err != nil {
		return nil, fmt.Errorf("update filter group %s: %w", filterGroupID, err)
	}
	return &result, nil
}

// DeleteFilterGroup deletes a filter group
func (c *Client) DeleteFilterGroup(ctx context.Context, filterGroupID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/mirroring/filtergroups/"+filterGroupID, nil); err != nil {
		return fmt.Errorf("delete filter group %s: %w", filterGroupID, err)
	}
	return nil
}

// ================================
// Filter Operations
// ================================

// ListFilters lists all filters
func (c *Client) ListFilters(ctx context.Context) (*ListFiltersOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result ListFiltersOutput
	if err := c.httpClient.GET(ctx, "/v2.0/mirroring/filters", &result); err != nil {
		return nil, fmt.Errorf("list filters: %w", err)
	}
	return &result, nil
}

// GetFilter gets a filter by ID
func (c *Client) GetFilter(ctx context.Context, filterID string) (*FilterOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result FilterOutput
	if err := c.httpClient.GET(ctx, "/v2.0/mirroring/filters/"+filterID, &result); err != nil {
		return nil, fmt.Errorf("get filter %s: %w", filterID, err)
	}
	return &result, nil
}

// CreateFilter creates a new filter
func (c *Client) CreateFilter(ctx context.Context, input *CreateFilterInput) (*FilterOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"mirroring_filter": input}
	var result FilterOutput
	if err := c.httpClient.POST(ctx, "/v2.0/mirroring/filters", req, &result); err != nil {
		return nil, fmt.Errorf("create filter: %w", err)
	}
	return &result, nil
}

// UpdateFilter updates a filter
func (c *Client) UpdateFilter(ctx context.Context, filterID string, input *UpdateFilterInput) (*FilterOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"mirroring_filter": input}
	var result FilterOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/mirroring/filters/"+filterID, req, &result); err != nil {
		return nil, fmt.Errorf("update filter %s: %w", filterID, err)
	}
	return &result, nil
}

// DeleteFilter deletes a filter
func (c *Client) DeleteFilter(ctx context.Context, filterID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/mirroring/filters/"+filterID, nil); err != nil {
		return fmt.Errorf("delete filter %s: %w", filterID, err)
	}
	return nil
}
