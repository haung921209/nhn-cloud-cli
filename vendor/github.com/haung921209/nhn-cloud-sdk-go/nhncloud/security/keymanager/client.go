package keymanager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Client struct {
	baseURL         string
	httpClient      *http.Client
	appKey          string
	userAccessKeyID string
	secretAccessKey string
	debug           bool
}

func NewClient(region, appKey, userAccessKeyID, secretAccessKey string, debug bool) *Client {
	baseURL := "https://api-keymanager.nhncloudservice.com"
	return &Client{
		baseURL:         baseURL,
		appKey:          appKey,
		userAccessKeyID: userAccessKeyID,
		secretAccessKey: secretAccessKey,
		debug:           debug,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) request(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = path.Join(u.Path, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "nhn-cloud-sdk-go/1.0.0")
	req.Header["X-TC-AUTHENTICATION-ID"] = []string{c.userAccessKeyID}
	req.Header["X-TC-AUTHENTICATION-SECRET"] = []string{c.secretAccessKey}

	if c.debug {
		fmt.Printf("=== REQUEST ===\n")
		fmt.Printf("%s %s\n", req.Method, req.URL.String())
		fmt.Printf("===============\n")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if c.debug {
		fmt.Printf("=== RESPONSE ===\n")
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Body: %s\n", string(respBody))
		fmt.Printf("================\n")
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s - %s", resp.Status, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func (c *Client) buildPath(resource string) string {
	return fmt.Sprintf("/keymanager/v1.2/appkey/%s/%s", c.appKey, resource)
}

// ============== Client Info APIs ==============

func (c *Client) GetClientInfo(ctx context.Context) (*GetClientInfoOutput, error) {
	path := c.buildPath("confirm")

	var out GetClientInfoOutput
	if err := c.request(ctx, "GET", path, nil, &out); err != nil {
		return nil, fmt.Errorf("get client info: %w", err)
	}
	return &out, nil
}

// ============== Key Store APIs ==============

func (c *Client) ListKeyStores(ctx context.Context) (*ListKeyStoresOutput, error) {
	path := c.buildPath("keystores")

	var out ListKeyStoresOutput
	if err := c.request(ctx, "GET", path, nil, &out); err != nil {
		return nil, fmt.Errorf("list key stores: %w", err)
	}
	return &out, nil
}

func (c *Client) GetKeyStore(ctx context.Context, keyStoreID string) (*GetKeyStoreOutput, error) {
	path := c.buildPath(fmt.Sprintf("keystores/%s", keyStoreID))

	var out GetKeyStoreOutput
	if err := c.request(ctx, "GET", path, nil, &out); err != nil {
		return nil, fmt.Errorf("get key store %s: %w", keyStoreID, err)
	}
	return &out, nil
}

func (c *Client) ListKeys(ctx context.Context, keyStoreID string) (*ListKeysOutput, error) {
	path := c.buildPath(fmt.Sprintf("keystores/%s/keys", keyStoreID))

	var out ListKeysOutput
	if err := c.request(ctx, "GET", path, nil, &out); err != nil {
		return nil, fmt.Errorf("list keys in store %s: %w", keyStoreID, err)
	}
	return &out, nil
}

func (c *Client) GetKey(ctx context.Context, keyStoreID, keyID string) (*GetKeyOutput, error) {
	path := c.buildPath(fmt.Sprintf("keystores/%s/keys/%s", keyStoreID, keyID))

	var out GetKeyOutput
	if err := c.request(ctx, "GET", path, nil, &out); err != nil {
		return nil, fmt.Errorf("get key %s: %w", keyID, err)
	}
	return &out, nil
}

// ============== Secret APIs ==============

func (c *Client) GetSecret(ctx context.Context, keyID string) (*GetSecretOutput, error) {
	path := c.buildPath(fmt.Sprintf("secrets/%s", keyID))

	var out GetSecretOutput
	if err := c.request(ctx, "GET", path, nil, &out); err != nil {
		return nil, fmt.Errorf("get secret %s: %w", keyID, err)
	}
	return &out, nil
}

// ============== Symmetric Key APIs ==============

func (c *Client) GetSymmetricKey(ctx context.Context, keyID string) (*GetSymmetricKeyOutput, error) {
	path := c.buildPath(fmt.Sprintf("symmetric-keys/%s/symmetric-key", keyID))

	var out GetSymmetricKeyOutput
	if err := c.request(ctx, "GET", path, nil, &out); err != nil {
		return nil, fmt.Errorf("get symmetric key %s: %w", keyID, err)
	}
	return &out, nil
}

func (c *Client) Encrypt(ctx context.Context, keyID string, input *EncryptInput) (*EncryptOutput, error) {
	path := c.buildPath(fmt.Sprintf("symmetric-keys/%s/encrypt", keyID))

	var out EncryptOutput
	if err := c.request(ctx, "POST", path, input, &out); err != nil {
		return nil, fmt.Errorf("encrypt with key %s: %w", keyID, err)
	}
	return &out, nil
}

func (c *Client) Decrypt(ctx context.Context, keyID string, input *DecryptInput) (*DecryptOutput, error) {
	path := c.buildPath(fmt.Sprintf("symmetric-keys/%s/decrypt", keyID))

	var out DecryptOutput
	if err := c.request(ctx, "POST", path, input, &out); err != nil {
		return nil, fmt.Errorf("decrypt with key %s: %w", keyID, err)
	}
	return &out, nil
}

func (c *Client) CreateLocalKey(ctx context.Context, keyID string) (*CreateLocalKeyOutput, error) {
	path := c.buildPath(fmt.Sprintf("symmetric-keys/%s/create-local-key", keyID))

	var out CreateLocalKeyOutput
	if err := c.request(ctx, "POST", path, nil, &out); err != nil {
		return nil, fmt.Errorf("create local key for %s: %w", keyID, err)
	}
	return &out, nil
}

// ============== Asymmetric Key APIs ==============

func (c *Client) GetPrivateKey(ctx context.Context, keyID string) (*GetPrivateKeyOutput, error) {
	path := c.buildPath(fmt.Sprintf("asymmetric-keys/%s/privateKey", keyID))

	var out GetPrivateKeyOutput
	if err := c.request(ctx, "GET", path, nil, &out); err != nil {
		return nil, fmt.Errorf("get private key %s: %w", keyID, err)
	}
	return &out, nil
}

func (c *Client) GetPublicKey(ctx context.Context, keyID string) (*GetPublicKeyOutput, error) {
	path := c.buildPath(fmt.Sprintf("asymmetric-keys/%s/publicKey", keyID))

	var out GetPublicKeyOutput
	if err := c.request(ctx, "GET", path, nil, &out); err != nil {
		return nil, fmt.Errorf("get public key %s: %w", keyID, err)
	}
	return &out, nil
}

func (c *Client) Sign(ctx context.Context, keyID string, input *SignInput) (*SignOutput, error) {
	path := c.buildPath(fmt.Sprintf("asymmetric-keys/%s/sign", keyID))

	var out SignOutput
	if err := c.request(ctx, "POST", path, input, &out); err != nil {
		return nil, fmt.Errorf("sign with key %s: %w", keyID, err)
	}
	return &out, nil
}

func (c *Client) Verify(ctx context.Context, keyID string, input *VerifyInput) (*VerifyOutput, error) {
	path := c.buildPath(fmt.Sprintf("asymmetric-keys/%s/verify", keyID))

	var out VerifyOutput
	if err := c.request(ctx, "POST", path, input, &out); err != nil {
		return nil, fmt.Errorf("verify with key %s: %w", keyID, err)
	}
	return &out, nil
}

// ============== Key Management APIs ==============

func (c *Client) CreateKey(ctx context.Context, keyType string, input *CreateKeyInput) (*CreateKeyOutput, error) {
	path := c.buildPath(fmt.Sprintf("keys/%s/create", keyType))

	var out CreateKeyOutput
	if err := c.request(ctx, "POST", path, input, &out); err != nil {
		return nil, fmt.Errorf("create key: %w", err)
	}
	return &out, nil
}

func (c *Client) RequestKeyDeletion(ctx context.Context, keyID string, input *DeleteKeyInput) (*DeleteKeyOutput, error) {
	path := c.buildPath(fmt.Sprintf("keys/%s/delete", keyID))

	var out DeleteKeyOutput
	if err := c.request(ctx, "PUT", path, input, &out); err != nil {
		return nil, fmt.Errorf("request key deletion %s: %w", keyID, err)
	}
	return &out, nil
}

func (c *Client) DeleteKeyImmediately(ctx context.Context, keyID string) (*APIResponse, error) {
	path := c.buildPath(fmt.Sprintf("keys/%s", keyID))

	var out APIResponse
	if err := c.request(ctx, "DELETE", path, nil, &out); err != nil {
		return nil, fmt.Errorf("delete key %s: %w", keyID, err)
	}
	return &out, nil
}
