package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	defaultOAuthBaseURL = "https://oauth.api.nhncloudservice.com"
	tokenBufferDuration = 5 * time.Minute
)

type OAuthTokenProvider struct {
	baseURL         string
	accessKeyID     string
	secretAccessKey string
	httpClient      *http.Client

	mu        sync.RWMutex
	token     string
	expiresAt time.Time
}

type oauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewOAuthTokenProvider(accessKeyID, secretAccessKey string) *OAuthTokenProvider {
	return &OAuthTokenProvider{
		baseURL:         defaultOAuthBaseURL,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func NewOAuthTokenProviderWithBaseURL(baseURL, accessKeyID, secretAccessKey string) *OAuthTokenProvider {
	p := NewOAuthTokenProvider(accessKeyID, secretAccessKey)
	p.baseURL = strings.TrimSuffix(baseURL, "/")
	return p
}

func (p *OAuthTokenProvider) GetToken(ctx context.Context) (string, error) {
	p.mu.RLock()
	if p.token != "" && time.Now().Add(tokenBufferDuration).Before(p.expiresAt) {
		token := p.token
		p.mu.RUnlock()
		return token, nil
	}
	p.mu.RUnlock()

	return p.refreshToken(ctx)
}

func (p *OAuthTokenProvider) SetAuthHeader(req *http.Request, token string) {
	req.Header.Set("x-nhn-authorization", "Bearer "+token)
}

func (p *OAuthTokenProvider) refreshToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.token != "" && time.Now().Add(tokenBufferDuration).Before(p.expiresAt) {
		return p.token, nil
	}

	tokenURL := p.baseURL + "/oauth2/token/create"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", p.basicAuthHeader())

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp oauthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("parsing token response: %w", err)
	}

	p.token = tokenResp.AccessToken
	p.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return p.token, nil
}

func (p *OAuthTokenProvider) basicAuthHeader() string {
	credentials := p.accessKeyID + ":" + p.secretAccessKey
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return "Basic " + encoded
}

func (p *OAuthTokenProvider) Invalidate() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.token = ""
	p.expiresAt = time.Time{}
}
