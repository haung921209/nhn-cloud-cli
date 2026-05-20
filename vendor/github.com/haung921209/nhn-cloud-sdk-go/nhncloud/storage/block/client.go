package block

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

	baseURL, err := c.tokenProvider.GetServiceEndpoint("volumev2", c.region)
	if err != nil {
		return fmt.Errorf("resolve endpoint: %w", err)
	}

	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}

	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
	return nil
}

func (c *Client) ListVolumes(ctx context.Context) (*ListVolumesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListVolumesOutput
	if err := c.httpClient.GET(ctx, "/volumes/detail", &out); err != nil {
		return nil, fmt.Errorf("list volumes: %w", err)
	}
	return &out, nil
}

func (c *Client) GetVolume(ctx context.Context, volumeID string) (*GetVolumeOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetVolumeOutput
	if err := c.httpClient.GET(ctx, "/volumes/"+volumeID, &out); err != nil {
		return nil, fmt.Errorf("get volume %s: %w", volumeID, err)
	}
	return &out, nil
}

func (c *Client) CreateVolume(ctx context.Context, input *CreateVolumeInput) (*CreateVolumeOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"volume": input}
	var out CreateVolumeOutput
	if err := c.httpClient.POST(ctx, "/volumes", req, &out); err != nil {
		return nil, fmt.Errorf("create volume: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateVolume(ctx context.Context, volumeID string, input *UpdateVolumeInput) (*GetVolumeOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"volume": input}
	var out GetVolumeOutput
	if err := c.httpClient.PUT(ctx, "/volumes/"+volumeID, req, &out); err != nil {
		return nil, fmt.Errorf("update volume %s: %w", volumeID, err)
	}
	return &out, nil
}

func (c *Client) DeleteVolume(ctx context.Context, volumeID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/volumes/"+volumeID, nil); err != nil {
		return fmt.Errorf("delete volume %s: %w", volumeID, err)
	}
	return nil
}

func (c *Client) ExtendVolume(ctx context.Context, volumeID string, newSize int) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	req := map[string]interface{}{
		"os-extend": map[string]int{"new_size": newSize},
	}
	if err := c.httpClient.POST(ctx, "/volumes/"+volumeID+"/action", req, nil); err != nil {
		return fmt.Errorf("extend volume %s: %w", volumeID, err)
	}
	return nil
}

func (c *Client) AttachVolume(ctx context.Context, volumeID, serverID, device string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	req := map[string]interface{}{
		"os-attach": map[string]string{
			"instance_uuid": serverID,
			"mountpoint":    device,
		},
	}
	if err := c.httpClient.POST(ctx, "/volumes/"+volumeID+"/action", req, nil); err != nil {
		return fmt.Errorf("attach volume %s: %w", volumeID, err)
	}
	return nil
}

func (c *Client) DetachVolume(ctx context.Context, volumeID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	req := map[string]interface{}{"os-detach": map[string]interface{}{}}
	if err := c.httpClient.POST(ctx, "/volumes/"+volumeID+"/action", req, nil); err != nil {
		return fmt.Errorf("detach volume %s: %w", volumeID, err)
	}
	return nil
}

func (c *Client) ListSnapshots(ctx context.Context) (*ListSnapshotsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListSnapshotsOutput
	if err := c.httpClient.GET(ctx, "/snapshots/detail", &out); err != nil {
		return nil, fmt.Errorf("list snapshots: %w", err)
	}
	return &out, nil
}

func (c *Client) GetSnapshot(ctx context.Context, snapshotID string) (*GetSnapshotOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetSnapshotOutput
	if err := c.httpClient.GET(ctx, "/snapshots/"+snapshotID, &out); err != nil {
		return nil, fmt.Errorf("get snapshot %s: %w", snapshotID, err)
	}
	return &out, nil
}

func (c *Client) CreateSnapshot(ctx context.Context, input *CreateSnapshotInput) (*CreateSnapshotOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"snapshot": input}
	var out CreateSnapshotOutput
	if err := c.httpClient.POST(ctx, "/snapshots", req, &out); err != nil {
		return nil, fmt.Errorf("create snapshot: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteSnapshot(ctx context.Context, snapshotID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/snapshots/"+snapshotID, nil); err != nil {
		return fmt.Errorf("delete snapshot %s: %w", snapshotID, err)
	}
	return nil
}

func (c *Client) ListVolumeTypes(ctx context.Context) (*ListVolumeTypesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListVolumeTypesOutput
	if err := c.httpClient.GET(ctx, "/types", &out); err != nil {
		return nil, fmt.Errorf("list volume types: %w", err)
	}
	return &out, nil
}
