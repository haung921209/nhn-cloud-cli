package compute

import (
	"context"
	"fmt"
)

func (c *Client) ListFlavors(ctx context.Context) (*ListFlavorsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListFlavorsOutput
	if err := c.httpClient.GET(ctx, "/flavors/detail", &out); err != nil {
		return nil, fmt.Errorf("list flavors: %w", err)
	}
	return &out, nil
}
