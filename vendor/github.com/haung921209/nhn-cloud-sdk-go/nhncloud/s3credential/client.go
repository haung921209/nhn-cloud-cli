// Package s3credential provides S3 Credential service client
package s3credential

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

// Client represents a S3 Credential API client
type Client struct {
	region        string
	credentials   credentials.IdentityCredentials
	httpClient    *client.Client
	tokenProvider *client.IdentityTokenProvider
	debug         bool
}

// NewClient creates a new S3 Credential client
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

	baseURL := fmt.Sprintf("https://api-identity-infrastructure.%s.nhncloudservice.com", c.region)
	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}
	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)

	return nil
}

// ListCredentials lists all S3 credentials for a user
func (c *Client) ListCredentials(ctx context.Context, userID string) (*ListCredentialsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v2.0/users/%s/credentials/OS-EC2", userID)
	var creds []S3Credential
	if err := c.httpClient.GET(ctx, path, &creds); err != nil {
		return nil, fmt.Errorf("list credentials: %w", err)
	}

	return &ListCredentialsOutput{Credentials: creds}, nil
}

// CreateCredential creates a new S3 credential
func (c *Client) CreateCredential(ctx context.Context, apiUserID string, input *CreateCredentialInput) (*CredentialOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v2.0/users/%s/credentials/OS-EC2", apiUserID)
	var result CredentialOutput
	if err := c.httpClient.POST(ctx, path, input, &result); err != nil {
		return nil, fmt.Errorf("create credential: %w", err)
	}

	return &result, nil
}

// DeleteCredential deletes an S3 credential
func (c *Client) DeleteCredential(ctx context.Context, userID, accessKey string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	path := fmt.Sprintf("/v2.0/users/%s/credentials/OS-EC2/%s", userID, accessKey)
	if err := c.httpClient.DELETE(ctx, path, nil); err != nil {
		return fmt.Errorf("delete credential: %w", err)
	}
	return nil
}
