package credentials

import (
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
	OAuthBaseURL   = "https://oauth.api.nhncloudservice.com"
	TokenCreateURL = "/oauth2/token/create"
)

type Token struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	IssuedAt    time.Time `json:"issued_at"`
}

func (t *Token) IsExpired() bool {
	if t.AccessToken == "" {
		return true
	}
	expirationTime := t.IssuedAt.Add(time.Duration(t.ExpiresIn) * time.Second)
	return time.Now().Add(5 * time.Minute).After(expirationTime)
}

func (t *Token) IsValid() bool {
	return t.AccessToken != "" && !t.IsExpired()
}

type TokenProvider struct {
	accessKeyID     string
	secretAccessKey string
	httpClient      *http.Client
	token           *Token
	mutex           sync.RWMutex
}

func NewTokenProvider(accessKeyID, secretAccessKey string) *TokenProvider {
	return &TokenProvider{
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *TokenProvider) GetToken() (*Token, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.token != nil && p.token.IsValid() {
		return p.token, nil
	}

	token, err := p.fetchNewToken()
	if err != nil {
		return nil, err
	}

	p.token = token
	return token, nil
}

func (p *TokenProvider) GetBearerToken() (string, error) {
	token, err := p.GetToken()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func (p *TokenProvider) fetchNewToken() (*Token, error) {
	tokenURL := OAuthBaseURL + TokenCreateURL

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", p.getBasicAuthHeader())

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}

	return &Token{
		AccessToken: tokenResponse.AccessToken,
		TokenType:   tokenResponse.TokenType,
		ExpiresIn:   tokenResponse.ExpiresIn,
		IssuedAt:    time.Now(),
	}, nil
}

func (p *TokenProvider) getBasicAuthHeader() string {
	credentials := fmt.Sprintf("%s:%s", p.accessKeyID, p.secretAccessKey)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return fmt.Sprintf("Basic %s", encoded)
}

func (p *TokenProvider) GetAccessKeyID() string {
	return p.accessKeyID
}

func (p *TokenProvider) GetSecretAccessKey() string {
	return p.secretAccessKey
}
