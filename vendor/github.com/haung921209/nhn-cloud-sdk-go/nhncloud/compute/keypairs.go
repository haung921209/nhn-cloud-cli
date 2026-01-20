package compute

import (
	"context"
	"fmt"
)

func (c *Client) ListKeyPairs(ctx context.Context) (*ListKeyPairsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListKeyPairsOutput
	if err := c.httpClient.GET(ctx, "/os-keypairs", &out); err != nil {
		return nil, fmt.Errorf("list keypairs: %w", err)
	}
	return &out, nil
}

func (c *Client) CreateKeyPair(ctx context.Context, input *CreateKeyPairInput) (*CreateKeyPairOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"keypair": input}
	var out CreateKeyPairOutput
	if err := c.httpClient.POST(ctx, "/os-keypairs", req, &out); err != nil {
		return nil, fmt.Errorf("create keypair: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteKeyPair(ctx context.Context, name string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/os-keypairs/"+name, nil); err != nil {
		return fmt.Errorf("delete keypair %s: %w", name, err)
	}
	return nil
}
