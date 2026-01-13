package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Token struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	IssuedAt    time.Time `json:"issued_at"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenManager struct {
	region          string
	accessKeyID     string
	secretAccessKey string
	cachePath       string
	httpClient      *http.Client
}

func NewTokenManager(region, accessKeyID, secretAccessKey string) *TokenManager {
	home, _ := os.UserHomeDir()
	cacheDir := filepath.Join(home, ".nhncloud", "cache")
	os.MkdirAll(cacheDir, 0700)

	return &TokenManager{
		region:          region,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		cachePath:       filepath.Join(cacheDir, "token.json"),
		httpClient:      &http.Client{Timeout: 30 * time.Second},
	}
}

func (m *TokenManager) GetToken() (*Token, error) {
	token, err := m.loadCachedToken()
	if err == nil && token.IsValid() {
		return token, nil
	}

	return m.RefreshToken()
}

func (m *TokenManager) RefreshToken() (*Token, error) {
	tokenURL := "https://oauth.api.nhncloudservice.com/oauth2/token/create"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	authStr := base64.StdEncoding.EncodeToString([]byte(m.accessKeyID + ":" + m.secretAccessKey))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+authStr)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	token := &Token{
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		ExpiresIn:   tokenResp.ExpiresIn,
		IssuedAt:    time.Now(),
	}

	if err := m.saveToken(token); err != nil {
		return nil, fmt.Errorf("saving token: %w", err)
	}

	return token, nil
}

func (m *TokenManager) ClearToken() error {
	if err := os.Remove(m.cachePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (m *TokenManager) loadCachedToken() (*Token, error) {
	data, err := os.ReadFile(m.cachePath)
	if err != nil {
		return nil, err
	}

	var token Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (m *TokenManager) saveToken(token *Token) error {
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.cachePath, data, 0600)
}

func (t *Token) IsValid() bool {
	if t.AccessToken == "" {
		return false
	}
	expiresAt := t.IssuedAt.Add(time.Duration(t.ExpiresIn) * time.Second)
	return time.Now().Add(5 * time.Minute).Before(expiresAt)
}

func (t *Token) ExpiresAt() time.Time {
	return t.IssuedAt.Add(time.Duration(t.ExpiresIn) * time.Second)
}
