package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const defaultIdentityURL = "https://api-identity-infrastructure.nhncloudservice.com"

type IdentityTokenProvider struct {
	identityURL string
	tenantID    string
	username    string
	password    string
	httpClient  *http.Client

	mu             sync.RWMutex
	token          string
	expiresAt      time.Time
	serviceCatalog []IdentityService
}

type identityTokenRequest struct {
	Auth identityAuth `json:"auth"`
}

type identityAuth struct {
	TenantID            string                      `json:"tenantId"`
	PasswordCredentials identityPasswordCredentials `json:"passwordCredentials"`
}

type identityPasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type identityTokenResponse struct {
	Access identityAccess `json:"access"`
}

type identityAccess struct {
	Token          identityToken     `json:"token"`
	ServiceCatalog []IdentityService `json:"serviceCatalog"`
}

type identityToken struct {
	ID      string    `json:"id"`
	Expires time.Time `json:"expires"`
}

type IdentityService struct {
	Name      string             `json:"name"`
	Type      string             `json:"type"`
	Endpoints []IdentityEndpoint `json:"endpoints"`
}

type IdentityEndpoint struct {
	PublicURL string `json:"publicURL"`
	Region    string `json:"region"`
}

func NewIdentityTokenProvider(tenantID, username, password string) *IdentityTokenProvider {
	return &IdentityTokenProvider{
		identityURL: defaultIdentityURL,
		tenantID:    tenantID,
		username:    username,
		password:    password,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func NewIdentityTokenProviderWithURL(identityURL, tenantID, username, password string) *IdentityTokenProvider {
	p := NewIdentityTokenProvider(tenantID, username, password)
	p.identityURL = strings.TrimSuffix(identityURL, "/")
	return p
}

func (p *IdentityTokenProvider) GetToken(ctx context.Context) (string, error) {
	p.mu.RLock()
	if p.token != "" && time.Now().Add(tokenBufferDuration).Before(p.expiresAt) {
		token := p.token
		p.mu.RUnlock()
		return token, nil
	}
	p.mu.RUnlock()

	return p.refreshToken(ctx)
}

func (p *IdentityTokenProvider) SetAuthHeader(req *http.Request, token string) {
	req.Header.Set("X-Auth-Token", token)
}

func (p *IdentityTokenProvider) refreshToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.token != "" && time.Now().Add(tokenBufferDuration).Before(p.expiresAt) {
		return p.token, nil
	}

	tokenURL := p.identityURL + "/v2.0/tokens"

	tokenReq := identityTokenRequest{
		Auth: identityAuth{
			TenantID: p.tenantID,
			PasswordCredentials: identityPasswordCredentials{
				Username: p.username,
				Password: p.password,
			},
		},
	}

	reqBody, err := json.Marshal(tokenReq)
	if err != nil {
		return "", fmt.Errorf("marshaling token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("creating token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

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
		return "", fmt.Errorf("identity token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp identityTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("parsing token response: %w", err)
	}

	p.token = tokenResp.Access.Token.ID
	p.expiresAt = tokenResp.Access.Token.Expires
	p.serviceCatalog = tokenResp.Access.ServiceCatalog

	return p.token, nil
}

func (p *IdentityTokenProvider) ServiceCatalog() []IdentityService {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.serviceCatalog
}

func (p *IdentityTokenProvider) GetServiceEndpoint(serviceType, region string) (string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, svc := range p.serviceCatalog {
		if svc.Type == serviceType {
			for _, ep := range svc.Endpoints {
				if strings.EqualFold(ep.Region, region) {
					return ep.PublicURL, nil
				}
			}
			if len(svc.Endpoints) > 0 {
				return svc.Endpoints[0].PublicURL, nil
			}
		}
	}

	return "", fmt.Errorf("service endpoint not found for type=%s, region=%s", serviceType, region)
}

func (p *IdentityTokenProvider) Invalidate() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.token = ""
	p.expiresAt = time.Time{}
	p.serviceCatalog = nil
}
