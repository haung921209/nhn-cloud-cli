// Package auth provides Bearer Token authentication for PostgreSQL v1.0.
package auth

import "net/http"

// BearerAuth implements Bearer Token authentication for PostgreSQL v1.0.
// PostgreSQL uses a different authentication method than MySQL/MariaDB v3.0.
//
// Authentication Flow:
// 1. Obtain Bearer token via OAuth2 using Access Key ID and Secret Access Key
// 2. Use the Bearer token in X-NHN-AUTHORIZATION header
//
// Example:
//
//	auth := auth.NewBearerAuth(appKey, token)
//	client := core.NewClient(baseURL, auth, nil)
type BearerAuth struct {
	appKey string
	token  string
}

// NewBearerAuth creates a new Bearer token authenticator for PostgreSQL.
//
// Parameters:
//   - appKey: NHN Cloud App Key
//   - token: Bearer token obtained via OAuth2 token exchange
func NewBearerAuth(appKey, token string) *BearerAuth {
	return &BearerAuth{
		appKey: appKey,
		token:  token,
	}
}

// Authenticate adds Bearer token authentication headers to the request.
//
// Headers added:
//   - X-TC-APP-KEY: {appKey}
//   - X-NHN-AUTHORIZATION: Bearer {token}
func (a *BearerAuth) Authenticate(req *http.Request) error {
	req.Header.Set("X-TC-APP-KEY", a.appKey)
	req.Header.Set("X-NHN-AUTHORIZATION", "Bearer "+a.token)
	return nil
}
