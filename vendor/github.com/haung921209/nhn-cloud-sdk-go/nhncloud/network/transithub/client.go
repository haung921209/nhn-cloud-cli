package transithub

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

// Client represents a Transit Hub service client
type Client struct {
	region        string
	credentials   credentials.IdentityCredentials
	httpClient    *client.Client
	tokenProvider *client.IdentityTokenProvider
	debug         bool
}

// NewClient creates a new Transit Hub client
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

// =============================================================================
// Transit Hub APIs
// =============================================================================

func (c *Client) ListTransitHubs(ctx context.Context) (*ListTransitHubsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out ListTransitHubsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithubs", &out); err != nil {
		return nil, fmt.Errorf("list transit hubs: %w", err)
	}
	return &out, nil
}

func (c *Client) GetTransitHub(ctx context.Context, id string) (*GetTransitHubOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out GetTransitHubOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithubs/"+id, &out); err != nil {
		return nil, fmt.Errorf("get transit hub %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateTransitHub(ctx context.Context, input *CreateTransitHubInput) (*GetTransitHubOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &CreateTransitHubRequest{TransitHub: *input}
	var out GetTransitHubOutput
	if err := c.httpClient.POST(ctx, "/v2.0/gateways/transithubs", req, &out); err != nil {
		return nil, fmt.Errorf("create transit hub: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateTransitHub(ctx context.Context, id string, input *UpdateTransitHubInput) (*GetTransitHubOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &UpdateTransitHubRequest{TransitHub: *input}
	var out GetTransitHubOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/gateways/transithubs/"+id, req, &out); err != nil {
		return nil, fmt.Errorf("update transit hub %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) DeleteTransitHub(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}
	if err := c.httpClient.DELETE(ctx, "/v2.0/gateways/transithubs/"+id, nil); err != nil {
		return fmt.Errorf("delete transit hub %s: %w", id, err)
	}
	return nil
}

// =============================================================================
// Attachment APIs
// =============================================================================

func (c *Client) ListAttachments(ctx context.Context) (*ListAttachmentsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out ListAttachmentsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_attachments", &out); err != nil {
		return nil, fmt.Errorf("list attachments: %w", err)
	}
	return &out, nil
}

func (c *Client) GetAttachment(ctx context.Context, id string) (*GetAttachmentOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out GetAttachmentOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_attachments/"+id, &out); err != nil {
		return nil, fmt.Errorf("get attachment %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateAttachment(ctx context.Context, input *CreateAttachmentInput) (*GetAttachmentOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &CreateAttachmentRequest{Attachment: *input}
	var out GetAttachmentOutput
	if err := c.httpClient.POST(ctx, "/v2.0/gateways/transithub_attachments", req, &out); err != nil {
		return nil, fmt.Errorf("create attachment: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateAttachment(ctx context.Context, id string, input *UpdateAttachmentInput) (*GetAttachmentOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &UpdateAttachmentRequest{Attachment: *input}
	var out GetAttachmentOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/gateways/transithub_attachments/"+id, req, &out); err != nil {
		return nil, fmt.Errorf("update attachment %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) DeleteAttachment(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}
	if err := c.httpClient.DELETE(ctx, "/v2.0/gateways/transithub_attachments/"+id, nil); err != nil {
		return fmt.Errorf("delete attachment %s: %w", id, err)
	}
	return nil
}

// =============================================================================
// Routing Table APIs
// =============================================================================

func (c *Client) ListRoutingTables(ctx context.Context) (*ListRoutingTablesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out ListRoutingTablesOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_routing_tables", &out); err != nil {
		return nil, fmt.Errorf("list routing tables: %w", err)
	}
	return &out, nil
}

func (c *Client) GetRoutingTable(ctx context.Context, id string) (*GetRoutingTableOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out GetRoutingTableOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_routing_tables/"+id, &out); err != nil {
		return nil, fmt.Errorf("get routing table %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateRoutingTable(ctx context.Context, input *CreateRoutingTableInput) (*GetRoutingTableOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &CreateRoutingTableRequest{RoutingTable: *input}
	var out GetRoutingTableOutput
	if err := c.httpClient.POST(ctx, "/v2.0/gateways/transithub_routing_tables", req, &out); err != nil {
		return nil, fmt.Errorf("create routing table: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateRoutingTable(ctx context.Context, id string, input *UpdateRoutingTableInput) (*GetRoutingTableOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &UpdateRoutingTableRequest{RoutingTable: *input}
	var out GetRoutingTableOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/gateways/transithub_routing_tables/"+id, req, &out); err != nil {
		return nil, fmt.Errorf("update routing table %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) DeleteRoutingTable(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}
	if err := c.httpClient.DELETE(ctx, "/v2.0/gateways/transithub_routing_tables/"+id, nil); err != nil {
		return fmt.Errorf("delete routing table %s: %w", id, err)
	}
	return nil
}

// =============================================================================
// Routing Association APIs
// =============================================================================

func (c *Client) ListRoutingAssociations(ctx context.Context) (*ListRoutingAssociationsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out ListRoutingAssociationsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_routing_associations", &out); err != nil {
		return nil, fmt.Errorf("list routing associations: %w", err)
	}
	return &out, nil
}

func (c *Client) GetRoutingAssociation(ctx context.Context, id string) (*GetRoutingAssociationOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out GetRoutingAssociationOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_routing_associations/"+id, &out); err != nil {
		return nil, fmt.Errorf("get routing association %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateRoutingAssociation(ctx context.Context, input *CreateRoutingAssociationInput) (*GetRoutingAssociationOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &CreateRoutingAssociationRequest{Association: *input}
	var out GetRoutingAssociationOutput
	if err := c.httpClient.POST(ctx, "/v2.0/gateways/transithub_routing_associations", req, &out); err != nil {
		return nil, fmt.Errorf("create routing association: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteRoutingAssociation(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}
	if err := c.httpClient.DELETE(ctx, "/v2.0/gateways/transithub_routing_associations/"+id, nil); err != nil {
		return fmt.Errorf("delete routing association %s: %w", id, err)
	}
	return nil
}

// =============================================================================
// Routing Propagation APIs
// =============================================================================

func (c *Client) ListRoutingPropagations(ctx context.Context) (*ListRoutingPropagationsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out ListRoutingPropagationsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_routing_propagations", &out); err != nil {
		return nil, fmt.Errorf("list routing propagations: %w", err)
	}
	return &out, nil
}

func (c *Client) GetRoutingPropagation(ctx context.Context, id string) (*GetRoutingPropagationOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out GetRoutingPropagationOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_routing_propagations/"+id, &out); err != nil {
		return nil, fmt.Errorf("get routing propagation %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateRoutingPropagation(ctx context.Context, input *CreateRoutingPropagationInput) (*GetRoutingPropagationOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &CreateRoutingPropagationRequest{Propagation: *input}
	var out GetRoutingPropagationOutput
	if err := c.httpClient.POST(ctx, "/v2.0/gateways/transithub_routing_propagations", req, &out); err != nil {
		return nil, fmt.Errorf("create routing propagation: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteRoutingPropagation(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}
	if err := c.httpClient.DELETE(ctx, "/v2.0/gateways/transithub_routing_propagations/"+id, nil); err != nil {
		return fmt.Errorf("delete routing propagation %s: %w", id, err)
	}
	return nil
}

// =============================================================================
// Routing Rule APIs
// =============================================================================

func (c *Client) ListRoutingRules(ctx context.Context) (*ListRoutingRulesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out ListRoutingRulesOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_routing_rules", &out); err != nil {
		return nil, fmt.Errorf("list routing rules: %w", err)
	}
	return &out, nil
}

func (c *Client) GetRoutingRule(ctx context.Context, id string) (*GetRoutingRuleOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out GetRoutingRuleOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_routing_rules/"+id, &out); err != nil {
		return nil, fmt.Errorf("get routing rule %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateRoutingRule(ctx context.Context, input *CreateRoutingRuleInput) (*GetRoutingRuleOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &CreateRoutingRuleRequest{Rule: *input}
	var out GetRoutingRuleOutput
	if err := c.httpClient.POST(ctx, "/v2.0/gateways/transithub_routing_rules", req, &out); err != nil {
		return nil, fmt.Errorf("create routing rule: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateRoutingRule(ctx context.Context, id string, input *UpdateRoutingRuleInput) (*GetRoutingRuleOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &UpdateRoutingRuleRequest{Rule: *input}
	var out GetRoutingRuleOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/gateways/transithub_routing_rules/"+id, req, &out); err != nil {
		return nil, fmt.Errorf("update routing rule %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) DeleteRoutingRule(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}
	if err := c.httpClient.DELETE(ctx, "/v2.0/gateways/transithub_routing_rules/"+id, nil); err != nil {
		return fmt.Errorf("delete routing rule %s: %w", id, err)
	}
	return nil
}

// =============================================================================
// Multicast Domain APIs
// =============================================================================

func (c *Client) ListMulticastDomains(ctx context.Context) (*ListMulticastDomainsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out ListMulticastDomainsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_multicast_domains", &out); err != nil {
		return nil, fmt.Errorf("list multicast domains: %w", err)
	}
	return &out, nil
}

func (c *Client) GetMulticastDomain(ctx context.Context, id string) (*GetMulticastDomainOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	var out GetMulticastDomainOutput
	if err := c.httpClient.GET(ctx, "/v2.0/gateways/transithub_multicast_domains/"+id, &out); err != nil {
		return nil, fmt.Errorf("get multicast domain %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) CreateMulticastDomain(ctx context.Context, input *CreateMulticastDomainInput) (*GetMulticastDomainOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &CreateMulticastDomainRequest{Domain: *input}
	var out GetMulticastDomainOutput
	if err := c.httpClient.POST(ctx, "/v2.0/gateways/transithub_multicast_domains", req, &out); err != nil {
		return nil, fmt.Errorf("create multicast domain: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateMulticastDomain(ctx context.Context, id string, input *UpdateMulticastDomainInput) (*GetMulticastDomainOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}
	req := &UpdateMulticastDomainRequest{Domain: *input}
	var out GetMulticastDomainOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/gateways/transithub_multicast_domains/"+id, req, &out); err != nil {
		return nil, fmt.Errorf("update multicast domain %s: %w", id, err)
	}
	return &out, nil
}

func (c *Client) DeleteMulticastDomain(ctx context.Context, id string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}
	if err := c.httpClient.DELETE(ctx, "/v2.0/gateways/transithub_multicast_domains/"+id, nil); err != nil {
		return fmt.Errorf("delete multicast domain %s: %w", id, err)
	}
	return nil
}
