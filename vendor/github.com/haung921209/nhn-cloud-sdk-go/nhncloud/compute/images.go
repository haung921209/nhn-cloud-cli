package compute

import (
	"context"
	"fmt"
)

func (c *Client) ListImages(ctx context.Context) (*ListImagesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListImagesOutput
	if err := c.httpClient.GET(ctx, "/images/detail", &out); err != nil {
		return nil, fmt.Errorf("list images: %w", err)
	}
	return &out, nil
}
