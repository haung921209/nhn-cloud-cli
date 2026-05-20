package flowlog

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

// Client represents a FlowLog service client
type Client struct {
	region        string
	credentials   credentials.IdentityCredentials
	httpClient    *client.Client
	tokenProvider *client.IdentityTokenProvider
	debug         bool
}

// NewClient creates a new FlowLog client
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
// Logger Operations
// ================================

// ListLoggers lists all flow log loggers
func (c *Client) ListLoggers(ctx context.Context) (*ListLoggersOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListLoggersOutput
	if err := c.httpClient.GET(ctx, "/v2.0/flowlog-loggers", &out); err != nil {
		return nil, fmt.Errorf("list loggers: %w", err)
	}
	return &out, nil
}

// GetLogger gets a logger by ID
func (c *Client) GetLogger(ctx context.Context, loggerID string) (*GetLoggerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetLoggerOutput
	if err := c.httpClient.GET(ctx, "/v2.0/flowlog-loggers/"+loggerID, &out); err != nil {
		return nil, fmt.Errorf("get logger %s: %w", loggerID, err)
	}
	return &out, nil
}

// CreateLogger creates a new logger
func (c *Client) CreateLogger(ctx context.Context, input *CreateLoggerInput) (*GetLoggerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateLoggerRequest{Logger: *input}
	var out GetLoggerOutput
	if err := c.httpClient.POST(ctx, "/v2.0/flowlog-loggers", req, &out); err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}
	return &out, nil
}

// UpdateLogger updates a logger
func (c *Client) UpdateLogger(ctx context.Context, loggerID string, input *UpdateLoggerInput) (*GetLoggerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &UpdateLoggerRequest{Logger: *input}
	var out GetLoggerOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/flowlog-loggers/"+loggerID, req, &out); err != nil {
		return nil, fmt.Errorf("update logger %s: %w", loggerID, err)
	}
	return &out, nil
}

// DeleteLogger deletes a logger
func (c *Client) DeleteLogger(ctx context.Context, loggerID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/flowlog-loggers/"+loggerID, nil); err != nil {
		return fmt.Errorf("delete logger %s: %w", loggerID, err)
	}
	return nil
}

// ================================
// Logging Port Operations
// ================================

// ListLoggingPorts lists logging ports
func (c *Client) ListLoggingPorts(ctx context.Context) (*ListLoggingPortsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListLoggingPortsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/flowlog-logging-ports", &out); err != nil {
		return nil, fmt.Errorf("list logging ports: %w", err)
	}
	return &out, nil
}

// GetLoggingPort gets a logging port by ID
func (c *Client) GetLoggingPort(ctx context.Context, portID string) (*GetLoggingPortOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetLoggingPortOutput
	if err := c.httpClient.GET(ctx, "/v2.0/flowlog-logging-ports/"+portID, &out); err != nil {
		return nil, fmt.Errorf("get logging port %s: %w", portID, err)
	}
	return &out, nil
}
