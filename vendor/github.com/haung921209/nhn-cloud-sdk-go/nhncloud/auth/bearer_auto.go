// Package auth provides OAuth2 authentication and token management for RDS services.
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
	"sync"
	"time"
)

// TokenCache holds cached token information
type TokenCache struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt  time.Time `json:"issued_at"`
}

// TokenResponse represents OAuth2 token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // seconds
}

// BearerAuthWithAutoRefresh implements Bearer Token authentication with automatic token issuance
type BearerAuthWithAutoRefresh struct {
	appKey      string
	accessKeyID string
	secretKey   string
	token       string
	expiresAt   time.Time
	cacheFile   string
	mu          sync.RWMutex
}

// NewBearerAuthWithAutoRefresh creates a new Bearer token authenticator with auto-refresh
func NewBearerAuthWithAutoRefresh(appKey, accessKeyID, secretKey string) *BearerAuthWithAutoRefresh {
	homeDir, _ := os.UserHomeDir()
	cacheFile := filepath.Join(homeDir, ".nhncloud", "postgresql_token_cache.json")

	auth := &BearerAuthWithAutoRefresh{
		appKey:      appKey,
		accessKeyID: accessKeyID,
		secretKey:   secretKey,
		cacheFile:   cacheFile,
	}

	// Try to load cached token
	auth.loadCachedToken()

	return auth
}

// Authenticate adds Bearer token authentication headers to the request
// If token is expired or missing, it automatically issues a new one
func (a *BearerAuthWithAutoRefresh) Authenticate(req *http.Request) error {
	token, err := a.getValidToken()
	if err != nil {
		return fmt.Errorf("failed to get valid token: %w", err)
	}

	req.Header.Set("X-TC-APP-KEY", a.appKey)
	req.Header.Set("X-NHN-AUTHORIZATION", "Bearer "+token)
	return nil
}

// getValidToken returns a valid token, issuing a new one if necessary
func (a *BearerAuthWithAutoRefresh) getValidToken() (string, error) {
	a.mu.RLock()
	// Check if current token is valid (with 5 minute buffer)
	if a.token != "" && time.Now().Add(5*time.Minute).Before(a.expiresAt) {
		token := a.token
		a.mu.RUnlock()
		return token, nil
	}
	a.mu.RUnlock()

	// Need to issue new token
	return a.issueNewToken()
}

// issueNewToken requests a new Bearer token from OAuth2 server
func (a *BearerAuthWithAutoRefresh) issueNewToken() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Double-check after acquiring write lock
	if a.token != "" && time.Now().Add(5*time.Minute).Before(a.expiresAt) {
		return a.token, nil
	}

	// Build request
	tokenURL := "https://oauth.api.nhncloudservice.com/oauth2/token/create"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	// Basic Auth header: Base64(AccessKeyID:SecretAccessKey)
	credentials := base64.StdEncoding.EncodeToString([]byte(a.accessKeyID + ":" + a.secretKey))
	req.Header.Set("Authorization", "Basic "+credentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	// Store token
	a.token = tokenResp.AccessToken
	a.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// Cache token to disk
	a.saveCachedToken()

	return a.token, nil
}

// loadCachedToken attempts to load a cached token from disk
func (a *BearerAuthWithAutoRefresh) loadCachedToken() {
	data, err := os.ReadFile(a.cacheFile)
	if err != nil {
		return // No cache file, ignore
	}

	var cache TokenCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return // Invalid cache, ignore
	}

	// Only use if not expired (with 5 minute buffer)
	if time.Now().Add(5 * time.Minute).Before(cache.ExpiresAt) {
		a.token = cache.Token
		a.expiresAt = cache.ExpiresAt
	}
}

// saveCachedToken saves the current token to disk
func (a *BearerAuthWithAutoRefresh) saveCachedToken() {
	cache := TokenCache{
		Token:     a.token,
		ExpiresAt: a.expiresAt,
		IssuedAt:  time.Now(),
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return // Silently fail
	}

	// Ensure directory exists
	dir := filepath.Dir(a.cacheFile)
	os.MkdirAll(dir, 0700)

	// Write with restricted permissions
	os.WriteFile(a.cacheFile, data, 0600)
}
