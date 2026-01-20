package compute

import (
	"context"
	"fmt"
)

func (c *Client) ListAvailabilityZones(ctx context.Context) (*ListAvailabilityZonesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListAvailabilityZonesOutput
	if err := c.httpClient.GET(ctx, "/os-availability-zone", &out); err != nil {
		return nil, fmt.Errorf("list availability zones: %w", err)
	}
	return &out, nil
}
