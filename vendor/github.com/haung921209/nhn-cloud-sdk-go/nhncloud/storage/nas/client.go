package nas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

// Client handles NAS API operations
type Client struct {
	region        string
	credentials   credentials.IdentityCredentials
	httpClient    *http.Client
	tokenProvider *client.IdentityTokenProvider
	debug         bool
}

// NewClient creates a new NAS client
func NewClient(region string, creds credentials.IdentityCredentials, hc *http.Client, debug bool) *Client {
	c := &Client{
		region:      region,
		credentials: creds,
		httpClient:  hc,
		debug:       debug,
	}

	if hc == nil {
		c.httpClient = http.DefaultClient
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

func (c *Client) getBaseURL() string {
	return fmt.Sprintf("https://%s-api-nas-infrastructure.nhncloudservice.com", strings.ToLower(c.region))
}

func (c *Client) ensureToken(ctx context.Context) (string, error) {
	if c.tokenProvider == nil {
		return "", fmt.Errorf("no credentials provided")
	}
	return c.tokenProvider.GetToken(ctx)
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	token, err := c.ensureToken(ctx)
	if err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}

	u := c.getBaseURL() + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)

		if c.debug {
			fmt.Printf("[DEBUG] NAS Request Body: %s\n", string(jsonBody))
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("Content-Type", "application/json")

	if c.debug {
		fmt.Printf("[DEBUG] NAS Request: %s %s\n", method, u)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if c.debug {
		fmt.Printf("[DEBUG] NAS Response: %d %s\n", resp.StatusCode, resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if c.debug && len(respBody) > 0 {
		fmt.Printf("[DEBUG] NAS Response Body: %s\n", string(respBody))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error: %d %s - %s", resp.StatusCode, resp.Status, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// --- Volume Operations ---

// ListVolumes retrieves the list of NAS volumes
func (c *Client) ListVolumes(ctx context.Context, input *ListVolumesInput) (*ListVolumesOutput, error) {
	path := "/v1/volumes"

	if input != nil {
		params := url.Values{}
		if input.SizeGB != nil {
			params.Set("sizeGb", strconv.Itoa(*input.SizeGB))
		}
		if input.MaxSizeGB != nil {
			params.Set("maxSizeGb", strconv.Itoa(*input.MaxSizeGB))
		}
		if input.MinSizeGB != nil {
			params.Set("minSizeGb", strconv.Itoa(*input.MinSizeGB))
		}
		if input.Name != "" {
			params.Set("name", input.Name)
		}
		if input.NameContains != "" {
			params.Set("nameContains", input.NameContains)
		}
		if input.SubnetID != "" {
			params.Set("subnetId", input.SubnetID)
		}
		if input.Limit != nil {
			params.Set("limit", strconv.Itoa(*input.Limit))
		}
		if input.Page != nil {
			params.Set("page", strconv.Itoa(*input.Page))
		}

		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	var out ListVolumesOutput
	if err := c.doRequest(ctx, "GET", path, nil, &out); err != nil {
		return nil, fmt.Errorf("list volumes: %w", err)
	}
	return &out, nil
}

// GetVolume retrieves detailed information about a specific volume
func (c *Client) GetVolume(ctx context.Context, volumeID string) (*GetVolumeOutput, error) {
	var out GetVolumeOutput
	if err := c.doRequest(ctx, "GET", "/v1/volumes/"+volumeID, nil, &out); err != nil {
		return nil, fmt.Errorf("get volume %s: %w", volumeID, err)
	}
	return &out, nil
}

// CreateVolume creates a new NAS volume
func (c *Client) CreateVolume(ctx context.Context, input *CreateVolumeInput) (*CreateVolumeOutput, error) {
	req := map[string]interface{}{"volume": input}
	var out CreateVolumeOutput
	if err := c.doRequest(ctx, "POST", "/v1/volumes", req, &out); err != nil {
		return nil, fmt.Errorf("create volume: %w", err)
	}
	return &out, nil
}

// UpdateVolume updates a NAS volume configuration
func (c *Client) UpdateVolume(ctx context.Context, volumeID string, input *UpdateVolumeInput) (*UpdateVolumeOutput, error) {
	req := map[string]interface{}{"volume": input}
	var out UpdateVolumeOutput
	if err := c.doRequest(ctx, "PATCH", "/v1/volumes/"+volumeID, req, &out); err != nil {
		return nil, fmt.Errorf("update volume %s: %w", volumeID, err)
	}
	return &out, nil
}

// DeleteVolume deletes a NAS volume
func (c *Client) DeleteVolume(ctx context.Context, volumeID string) error {
	if err := c.doRequest(ctx, "DELETE", "/v1/volumes/"+volumeID, nil, nil); err != nil {
		return fmt.Errorf("delete volume %s: %w", volumeID, err)
	}
	return nil
}

// GetVolumeUsage retrieves usage information for a volume
func (c *Client) GetVolumeUsage(ctx context.Context, volumeID string) (*GetUsageOutput, error) {
	var out GetUsageOutput
	if err := c.doRequest(ctx, "GET", "/v1/volumes/"+volumeID+"/usage", nil, &out); err != nil {
		return nil, fmt.Errorf("get volume usage %s: %w", volumeID, err)
	}
	return &out, nil
}

// --- Interface Operations ---

// CreateInterface creates an interface for a volume
func (c *Client) CreateInterface(ctx context.Context, volumeID string, input *CreateInterfaceInput) (*CreateInterfaceOutput, error) {
	req := map[string]interface{}{"interface": input}
	var out CreateInterfaceOutput
	if err := c.doRequest(ctx, "POST", "/v1/volumes/"+volumeID+"/interfaces", req, &out); err != nil {
		return nil, fmt.Errorf("create interface for volume %s: %w", volumeID, err)
	}
	return &out, nil
}

// DeleteInterface deletes an interface from a volume
func (c *Client) DeleteInterface(ctx context.Context, volumeID, interfaceID string) error {
	if err := c.doRequest(ctx, "DELETE", "/v1/volumes/"+volumeID+"/interfaces/"+interfaceID, nil, nil); err != nil {
		return fmt.Errorf("delete interface %s from volume %s: %w", interfaceID, volumeID, err)
	}
	return nil
}

// --- Snapshot Operations ---

// ListSnapshots retrieves the list of snapshots for a volume
func (c *Client) ListSnapshots(ctx context.Context, volumeID string) (*ListSnapshotsOutput, error) {
	var out ListSnapshotsOutput
	if err := c.doRequest(ctx, "GET", "/v1/volumes/"+volumeID+"/snapshots", nil, &out); err != nil {
		return nil, fmt.Errorf("list snapshots for volume %s: %w", volumeID, err)
	}
	return &out, nil
}

// GetSnapshot retrieves detailed information about a specific snapshot
func (c *Client) GetSnapshot(ctx context.Context, volumeID, snapshotID string) (*GetSnapshotOutput, error) {
	var out GetSnapshotOutput
	if err := c.doRequest(ctx, "GET", "/v1/volumes/"+volumeID+"/snapshots/"+snapshotID, nil, &out); err != nil {
		return nil, fmt.Errorf("get snapshot %s: %w", snapshotID, err)
	}
	return &out, nil
}

// CreateSnapshot creates a snapshot of a volume
func (c *Client) CreateSnapshot(ctx context.Context, volumeID string, input *CreateSnapshotInput) (*CreateSnapshotOutput, error) {
	req := map[string]interface{}{"snapshot": input}
	var out CreateSnapshotOutput
	if err := c.doRequest(ctx, "POST", "/v1/volumes/"+volumeID+"/snapshots", req, &out); err != nil {
		return nil, fmt.Errorf("create snapshot for volume %s: %w", volumeID, err)
	}
	return &out, nil
}

// DeleteSnapshot deletes a snapshot
func (c *Client) DeleteSnapshot(ctx context.Context, volumeID, snapshotID string) error {
	if err := c.doRequest(ctx, "DELETE", "/v1/volumes/"+volumeID+"/snapshots/"+snapshotID, nil, nil); err != nil {
		return fmt.Errorf("delete snapshot %s: %w", snapshotID, err)
	}
	return nil
}

// RestoreSnapshot restores a volume from a snapshot
func (c *Client) RestoreSnapshot(ctx context.Context, volumeID, snapshotID string) error {
	if err := c.doRequest(ctx, "POST", "/v1/volumes/"+volumeID+"/snapshots/"+snapshotID+"/restore", nil, nil); err != nil {
		return fmt.Errorf("restore snapshot %s: %w", snapshotID, err)
	}
	return nil
}

// ListRestoreHistories retrieves restore histories for a volume
func (c *Client) ListRestoreHistories(ctx context.Context, volumeID string, input *ListRestoreHistoriesInput) (*ListRestoreHistoriesOutput, error) {
	path := "/v1/volumes/" + volumeID + "/restore-histories"

	if input != nil {
		params := url.Values{}
		if input.Limit != nil {
			params.Set("limit", strconv.Itoa(*input.Limit))
		}
		if input.Page != nil {
			params.Set("page", strconv.Itoa(*input.Page))
		}
		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	var out ListRestoreHistoriesOutput
	if err := c.doRequest(ctx, "GET", path, nil, &out); err != nil {
		return nil, fmt.Errorf("list restore histories for volume %s: %w", volumeID, err)
	}
	return &out, nil
}

// --- Volume Mirror Operations ---

// CreateVolumeMirror creates a volume mirror (replication)
func (c *Client) CreateVolumeMirror(ctx context.Context, volumeID string, input *CreateVolumeMirrorInput) (*CreateVolumeMirrorOutput, error) {
	req := map[string]interface{}{"volumeMirror": input}
	var out CreateVolumeMirrorOutput
	if err := c.doRequest(ctx, "POST", "/v1/volumes/"+volumeID+"/volume-mirrors", req, &out); err != nil {
		return nil, fmt.Errorf("create volume mirror for volume %s: %w", volumeID, err)
	}
	return &out, nil
}

// DeleteVolumeMirror deletes a volume mirror
func (c *Client) DeleteVolumeMirror(ctx context.Context, volumeID, mirrorID string) error {
	if err := c.doRequest(ctx, "DELETE", "/v1/volumes/"+volumeID+"/volume-mirrors/"+mirrorID, nil, nil); err != nil {
		return fmt.Errorf("delete volume mirror %s: %w", mirrorID, err)
	}
	return nil
}

// InvertVolumeMirrorDirection changes the direction of volume mirror
func (c *Client) InvertVolumeMirrorDirection(ctx context.Context, volumeID, mirrorID string) error {
	if err := c.doRequest(ctx, "POST", "/v1/volumes/"+volumeID+"/volume-mirrors/"+mirrorID+"/invert-direction", nil, nil); err != nil {
		return fmt.Errorf("invert volume mirror direction %s: %w", mirrorID, err)
	}
	return nil
}

// StartVolumeMirror starts volume mirror replication
func (c *Client) StartVolumeMirror(ctx context.Context, volumeID, mirrorID string) error {
	if err := c.doRequest(ctx, "POST", "/v1/volumes/"+volumeID+"/volume-mirrors/"+mirrorID+"/start", nil, nil); err != nil {
		return fmt.Errorf("start volume mirror %s: %w", mirrorID, err)
	}
	return nil
}

// StopVolumeMirror stops volume mirror replication
func (c *Client) StopVolumeMirror(ctx context.Context, volumeID, mirrorID string) error {
	if err := c.doRequest(ctx, "POST", "/v1/volumes/"+volumeID+"/volume-mirrors/"+mirrorID+"/stop", nil, nil); err != nil {
		return fmt.Errorf("stop volume mirror %s: %w", mirrorID, err)
	}
	return nil
}

// GetVolumeMirrorStat retrieves volume mirror statistics
func (c *Client) GetVolumeMirrorStat(ctx context.Context, volumeID, mirrorID string) (*GetVolumeMirrorStatOutput, error) {
	var out GetVolumeMirrorStatOutput
	if err := c.doRequest(ctx, "GET", "/v1/volumes/"+volumeID+"/volume-mirrors/"+mirrorID+"/stat", nil, &out); err != nil {
		return nil, fmt.Errorf("get volume mirror stat %s: %w", mirrorID, err)
	}
	return &out, nil
}
