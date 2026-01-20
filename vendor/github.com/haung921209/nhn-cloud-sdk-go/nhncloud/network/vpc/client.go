package vpc

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

func (c *Client) ListVPCs(ctx context.Context) (*ListVPCsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListVPCsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/vpcs", &out); err != nil {
		return nil, fmt.Errorf("list vpcs: %w", err)
	}
	return &out, nil
}

func (c *Client) GetVPC(ctx context.Context, vpcID string) (*GetVPCOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetVPCOutput
	if err := c.httpClient.GET(ctx, "/v2.0/vpcs/"+vpcID, &out); err != nil {
		return nil, fmt.Errorf("get vpc %s: %w", vpcID, err)
	}
	return &out, nil
}

func (c *Client) CreateVPC(ctx context.Context, input *CreateVPCInput) (*CreateVPCOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"vpc": input}
	var out CreateVPCOutput
	if err := c.httpClient.POST(ctx, "/v2.0/vpcs", req, &out); err != nil {
		return nil, fmt.Errorf("create vpc: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateVPC(ctx context.Context, vpcID string, input *UpdateVPCInput) (*GetVPCOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"vpc": input}
	var out GetVPCOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/vpcs/"+vpcID, req, &out); err != nil {
		return nil, fmt.Errorf("update vpc %s: %w", vpcID, err)
	}
	return &out, nil
}

func (c *Client) DeleteVPC(ctx context.Context, vpcID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/vpcs/"+vpcID, nil); err != nil {
		return fmt.Errorf("delete vpc %s: %w", vpcID, err)
	}
	return nil
}

func (c *Client) ListSubnets(ctx context.Context) (*ListSubnetsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListSubnetsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/subnets", &out); err != nil {
		return nil, fmt.Errorf("list subnets: %w", err)
	}
	return &out, nil
}

func (c *Client) GetSubnet(ctx context.Context, subnetID string) (*GetSubnetOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetSubnetOutput
	if err := c.httpClient.GET(ctx, "/v2.0/subnets/"+subnetID, &out); err != nil {
		return nil, fmt.Errorf("get subnet %s: %w", subnetID, err)
	}
	return &out, nil
}

func (c *Client) CreateSubnet(ctx context.Context, input *CreateSubnetInput) (*CreateSubnetOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"vpcsubnet": input}
	var out CreateSubnetOutput
	if err := c.httpClient.POST(ctx, "/v2.0/vpcsubnets", req, &out); err != nil {
		return nil, fmt.Errorf("create subnet: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteSubnet(ctx context.Context, subnetID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/vpcsubnets/"+subnetID, nil); err != nil {
		return fmt.Errorf("delete subnet %s: %w", subnetID, err)
	}
	return nil
}

func (c *Client) ListRoutingTables(ctx context.Context, vpcID string) (*ListRoutingTablesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	path := "/v2.0/routingtables"
	if vpcID != "" {
		path += "?vpc_id=" + vpcID
	}

	var out ListRoutingTablesOutput
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list routing tables: %w", err)
	}
	return &out, nil
}

func (c *Client) GetRoutingTable(ctx context.Context, tableID string) (*GetRoutingTableOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetRoutingTableOutput
	if err := c.httpClient.GET(ctx, "/v2.0/routingtables/"+tableID, &out); err != nil {
		return nil, fmt.Errorf("get routing table %s: %w", tableID, err)
	}
	return &out, nil
}
