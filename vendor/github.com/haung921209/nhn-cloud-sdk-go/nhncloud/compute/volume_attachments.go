package compute

import (
	"context"
	"fmt"
)

// VolumeAttachment is a block storage volume attached to a server.
type VolumeAttachment struct {
	ID       string `json:"id"`
	Device   string `json:"device"`
	ServerID string `json:"serverId"`
	VolumeID string `json:"volumeId"`
}

type VolumeAttachmentOutput struct {
	VolumeAttachment VolumeAttachment `json:"volumeAttachment"`
}

// AttachVolume attaches a volume to a server
// (POST /servers/{serverId}/os-volume_attachments). device is optional.
func (c *Client) AttachVolume(ctx context.Context, serverID, volumeID, device string) (*VolumeAttachmentOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	att := map[string]string{"volumeId": volumeID}
	if device != "" {
		att["device"] = device
	}
	req := map[string]interface{}{"volumeAttachment": att}
	var out VolumeAttachmentOutput
	if err := c.httpClient.POST(ctx, "/servers/"+serverID+"/os-volume_attachments", req, &out); err != nil {
		return nil, fmt.Errorf("attach volume %s to server %s: %w", volumeID, serverID, err)
	}
	return &out, nil
}

// DetachVolume detaches a volume from a server
// (DELETE /servers/{serverId}/os-volume_attachments/{volumeId}).
func (c *Client) DetachVolume(ctx context.Context, serverID, volumeID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/servers/"+serverID+"/os-volume_attachments/"+volumeID, nil); err != nil {
		return fmt.Errorf("detach volume %s from server %s: %w", volumeID, serverID, err)
	}
	return nil
}
