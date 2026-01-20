package compute

import (
	"context"
	"fmt"
)

func (c *Client) ListServers(ctx context.Context) (*ListServersOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListServersOutput
	if err := c.httpClient.GET(ctx, "/servers/detail", &out); err != nil {
		return nil, fmt.Errorf("list servers: %w", err)
	}
	return &out, nil
}

func (c *Client) GetServer(ctx context.Context, serverID string) (*GetServerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetServerOutput
	if err := c.httpClient.GET(ctx, "/servers/"+serverID, &out); err != nil {
		return nil, fmt.Errorf("get server %s: %w", serverID, err)
	}
	return &out, nil
}

func (c *Client) CreateServer(ctx context.Context, input *CreateServerInput) (*CreateServerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"server": input}
	var out CreateServerOutput
	if err := c.httpClient.POST(ctx, "/servers", req, &out); err != nil {
		return nil, fmt.Errorf("create server: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteServer(ctx context.Context, serverID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/servers/"+serverID, nil); err != nil {
		return fmt.Errorf("delete server %s: %w", serverID, err)
	}
	return nil
}

func (c *Client) StartServer(ctx context.Context, serverID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	req := map[string]interface{}{"os-start": nil}
	if err := c.httpClient.POST(ctx, "/servers/"+serverID+"/action", req, nil); err != nil {
		return fmt.Errorf("start server %s: %w", serverID, err)
	}
	return nil
}

func (c *Client) StopServer(ctx context.Context, serverID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	req := map[string]interface{}{"os-stop": nil}
	if err := c.httpClient.POST(ctx, "/servers/"+serverID+"/action", req, nil); err != nil {
		return fmt.Errorf("stop server %s: %w", serverID, err)
	}
	return nil
}

func (c *Client) RebootServer(ctx context.Context, serverID string, hard bool) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	rebootType := "SOFT"
	if hard {
		rebootType = "HARD"
	}
	req := map[string]interface{}{
		"reboot": map[string]string{"type": rebootType},
	}
	if err := c.httpClient.POST(ctx, "/servers/"+serverID+"/action", req, nil); err != nil {
		return fmt.Errorf("reboot server %s: %w", serverID, err)
	}
	return nil
}

func (c *Client) ResizeServer(ctx context.Context, serverID, flavorRef string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	req := map[string]interface{}{
		"resize": map[string]string{"flavorRef": flavorRef},
	}
	if err := c.httpClient.POST(ctx, "/servers/"+serverID+"/action", req, nil); err != nil {
		return fmt.Errorf("resize server %s: %w", serverID, err)
	}
	return nil
}

func (c *Client) ConfirmResize(ctx context.Context, serverID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	req := map[string]interface{}{"confirmResize": nil}
	if err := c.httpClient.POST(ctx, "/servers/"+serverID+"/action", req, nil); err != nil {
		return fmt.Errorf("confirm resize %s: %w", serverID, err)
	}
	return nil
}
