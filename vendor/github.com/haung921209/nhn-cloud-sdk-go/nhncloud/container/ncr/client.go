package ncr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

type Client struct {
	region        string
	appKey        string
	credentials   credentials.Credentials
	httpClient    *client.Client
	tokenProvider *client.OAuthTokenProvider
	debug         bool
}

func NewClient(region, appKey string, creds credentials.Credentials, hc *http.Client, debug bool) *Client {
	c := &Client{
		region:      region,
		appKey:      appKey,
		credentials: creds,
		debug:       debug,
	}

	if creds != nil {
		c.tokenProvider = client.NewOAuthTokenProvider(
			creds.GetAccessKeyID(),
			creds.GetSecretAccessKey(),
		)
		c.initHTTPClient()
	}

	return c
}

func (c *Client) initHTTPClient() {
	baseURL := fmt.Sprintf("https://%s-ncr.api.nhncloudservice.com/ncr/v2.0/appkeys/%s", c.region, c.appKey)
	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}
	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
}

func (c *Client) ListRegistries(ctx context.Context) (*ListRegistriesOutput, error) {
	var out ListRegistriesOutput
	if err := c.httpClient.GET(ctx, "/registries", &out); err != nil {
		return nil, fmt.Errorf("list registries: %w", err)
	}
	return &out, nil
}

func (c *Client) GetRegistry(ctx context.Context, registryID string) (*GetRegistryOutput, error) {
	var out GetRegistryOutput
	if err := c.httpClient.GET(ctx, "/registries/"+registryID, &out); err != nil {
		return nil, fmt.Errorf("get registry %s: %w", registryID, err)
	}
	return &out, nil
}

func (c *Client) CreateRegistry(ctx context.Context, input *CreateRegistryInput) (*CreateRegistryOutput, error) {
	var out CreateRegistryOutput
	if err := c.httpClient.POST(ctx, "/registries", input, &out); err != nil {
		return nil, fmt.Errorf("create registry: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateRegistry(ctx context.Context, registryID string, input *UpdateRegistryInput) (*GetRegistryOutput, error) {
	var out GetRegistryOutput
	if err := c.httpClient.PUT(ctx, "/registries/"+registryID, input, &out); err != nil {
		return nil, fmt.Errorf("update registry %s: %w", registryID, err)
	}
	return &out, nil
}

func (c *Client) DeleteRegistry(ctx context.Context, registryID string) error {
	if err := c.httpClient.DELETE(ctx, "/registries/"+registryID, nil); err != nil {
		return fmt.Errorf("delete registry %s: %w", registryID, err)
	}
	return nil
}

func (c *Client) ListImages(ctx context.Context, registryID string) (*ListImagesOutput, error) {
	var out ListImagesOutput
	if err := c.httpClient.GET(ctx, "/registries/"+registryID+"/images", &out); err != nil {
		return nil, fmt.Errorf("list images in registry %s: %w", registryID, err)
	}
	return &out, nil
}

func (c *Client) GetImage(ctx context.Context, registryID, imageName string) (*GetImageOutput, error) {
	var out GetImageOutput
	if err := c.httpClient.GET(ctx, "/registries/"+registryID+"/images/"+imageName, &out); err != nil {
		return nil, fmt.Errorf("get image %s: %w", imageName, err)
	}
	return &out, nil
}

func (c *Client) DeleteImage(ctx context.Context, registryID, imageName string) error {
	if err := c.httpClient.DELETE(ctx, "/registries/"+registryID+"/images/"+imageName, nil); err != nil {
		return fmt.Errorf("delete image %s: %w", imageName, err)
	}
	return nil
}

func (c *Client) ListTags(ctx context.Context, registryID, imageName string) (*ListTagsOutput, error) {
	var out ListTagsOutput
	if err := c.httpClient.GET(ctx, "/registries/"+registryID+"/images/"+imageName+"/tags", &out); err != nil {
		return nil, fmt.Errorf("list tags for image %s: %w", imageName, err)
	}
	return &out, nil
}

func (c *Client) DeleteTag(ctx context.Context, registryID, imageName, tagName string) error {
	if err := c.httpClient.DELETE(ctx, "/registries/"+registryID+"/images/"+imageName+"/tags/"+tagName, nil); err != nil {
		return fmt.Errorf("delete tag %s: %w", tagName, err)
	}
	return nil
}

func (c *Client) ScanImage(ctx context.Context, registryID, imageName, tag string) error {
	req := map[string]string{"tag": tag}
	if err := c.httpClient.POST(ctx, "/registries/"+registryID+"/images/"+imageName+"/scan", req, nil); err != nil {
		return fmt.Errorf("scan image %s:%s: %w", imageName, tag, err)
	}
	return nil
}

func (c *Client) GetImageScanResult(ctx context.Context, registryID, imageName, tag string) (*GetImageScanResultOutput, error) {
	var out GetImageScanResultOutput
	if err := c.httpClient.GET(ctx, "/registries/"+registryID+"/images/"+imageName+"/scan/"+tag, &out); err != nil {
		return nil, fmt.Errorf("get scan result for %s:%s: %w", imageName, tag, err)
	}
	return &out, nil
}

func (c *Client) ListWebhooks(ctx context.Context, registryID string) (*ListWebhooksOutput, error) {
	var out ListWebhooksOutput
	if err := c.httpClient.GET(ctx, "/registries/"+registryID+"/webhooks", &out); err != nil {
		return nil, fmt.Errorf("list webhooks for registry %s: %w", registryID, err)
	}
	return &out, nil
}

func (c *Client) CreateWebhook(ctx context.Context, registryID string, input *CreateWebhookInput) (*CreateWebhookOutput, error) {
	var out CreateWebhookOutput
	if err := c.httpClient.POST(ctx, "/registries/"+registryID+"/webhooks", input, &out); err != nil {
		return nil, fmt.Errorf("create webhook: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteWebhook(ctx context.Context, registryID, webhookID string) error {
	if err := c.httpClient.DELETE(ctx, "/registries/"+registryID+"/webhooks/"+webhookID, nil); err != nil {
		return fmt.Errorf("delete webhook %s: %w", webhookID, err)
	}
	return nil
}
