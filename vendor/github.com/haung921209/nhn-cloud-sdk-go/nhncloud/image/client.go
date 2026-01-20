package image

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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

	baseURL, err := c.tokenProvider.GetServiceEndpoint("image", c.region)
	if err != nil {
		return fmt.Errorf("resolve endpoint: %w", err)
	}

	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}
	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)

	return nil
}

func (c *Client) ListImages(ctx context.Context, input *ListImagesInput) (*ListImagesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	path := "/v2/images"
	if input != nil {
		params := url.Values{}
		if input.Limit > 0 {
			params.Add("limit", strconv.Itoa(input.Limit))
		}
		if input.Name != "" {
			params.Add("name", input.Name)
		}
		if input.Owner != "" {
			params.Add("owner", input.Owner)
		}
		if input.Status != "" {
			params.Add("status", input.Status)
		}
		if input.Visibility != "" {
			params.Add("visibility", input.Visibility)
		}
		if input.OSType != "" {
			params.Add("os_type", input.OSType)
		}
		if input.OSDistro != "" {
			params.Add("os_distro", input.OSDistro)
		}
		if input.Marker != "" {
			params.Add("marker", input.Marker)
		}
		if input.SortKey != "" {
			params.Add("sort_key", input.SortKey)
		}
		if input.SortDir != "" {
			params.Add("sort_dir", input.SortDir)
		}
		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	var result ListImagesOutput
	if err := c.httpClient.GET(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list images: %w", err)
	}
	return &result, nil
}

func (c *Client) GetImage(ctx context.Context, imageID string) (*Image, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result Image
	if err := c.httpClient.GET(ctx, "/v2/images/"+imageID, &result); err != nil {
		return nil, fmt.Errorf("get image %s: %w", imageID, err)
	}
	return &result, nil
}

func (c *Client) CreateImage(ctx context.Context, input *CreateImageInput) (*Image, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result Image
	if err := c.httpClient.POST(ctx, "/v2/images", input, &result); err != nil {
		return nil, fmt.Errorf("create image: %w", err)
	}
	return &result, nil
}

func (c *Client) UpdateImage(ctx context.Context, imageID string, ops []UpdateImageOp) (*Image, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result Image
	if err := c.httpClient.PATCH(ctx, "/v2/images/"+imageID, ops, &result); err != nil {
		return nil, fmt.Errorf("update image %s: %w", imageID, err)
	}
	return &result, nil
}

func (c *Client) DeleteImage(ctx context.Context, imageID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2/images/"+imageID, nil); err != nil {
		return fmt.Errorf("delete image %s: %w", imageID, err)
	}
	return nil
}

func (c *Client) AddTag(ctx context.Context, imageID, tag string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	path := fmt.Sprintf("/v2/images/%s/tags/%s", imageID, tag)
	if err := c.httpClient.PUT(ctx, path, nil, nil); err != nil {
		return fmt.Errorf("add tag: %w", err)
	}
	return nil
}

func (c *Client) RemoveTag(ctx context.Context, imageID, tag string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	path := fmt.Sprintf("/v2/images/%s/tags/%s", imageID, tag)
	if err := c.httpClient.DELETE(ctx, path, nil); err != nil {
		return fmt.Errorf("remove tag: %w", err)
	}
	return nil
}

func (c *Client) ListImageMembers(ctx context.Context, imageID string) (*ListImageMembersOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result ListImageMembersOutput
	if err := c.httpClient.GET(ctx, fmt.Sprintf("/v2/images/%s/members", imageID), &result); err != nil {
		return nil, fmt.Errorf("list image members: %w", err)
	}
	return &result, nil
}

func (c *Client) AddImageMember(ctx context.Context, imageID string, input *CreateImageMemberInput) (*ImageMember, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result ImageMember
	if err := c.httpClient.POST(ctx, fmt.Sprintf("/v2/images/%s/members", imageID), input, &result); err != nil {
		return nil, fmt.Errorf("add image member: %w", err)
	}
	return &result, nil
}

func (c *Client) UpdateImageMember(ctx context.Context, imageID, memberID string, input *UpdateImageMemberInput) (*ImageMember, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var result ImageMember
	path := fmt.Sprintf("/v2/images/%s/members/%s", imageID, memberID)
	if err := c.httpClient.PUT(ctx, path, input, &result); err != nil {
		return nil, fmt.Errorf("update image member: %w", err)
	}
	return &result, nil
}

func (c *Client) RemoveImageMember(ctx context.Context, imageID, memberID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	path := fmt.Sprintf("/v2/images/%s/members/%s", imageID, memberID)
	if err := c.httpClient.DELETE(ctx, path, nil); err != nil {
		return fmt.Errorf("remove image member: %w", err)
	}
	return nil
}
