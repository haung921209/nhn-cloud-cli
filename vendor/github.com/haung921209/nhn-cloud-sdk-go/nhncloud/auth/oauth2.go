// Package auth provides OAuth2 authentication for RDS services.
package auth

import "net/http"

// OAuth2Auth implements OAuth2 authentication for RDS services (MySQL, MariaDB, PostgreSQL)
type OAuth2Auth struct {
	appKey    string
	accessKey string
	secretKey string
}

// NewOAuth2Auth creates a new OAuth2 authenticator
func NewOAuth2Auth(appKey, accessKey, secretKey string) *OAuth2Auth {
	return &OAuth2Auth{
		appKey:    appKey,
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

// Authenticate adds OAuth2 headers to the request
func (a *OAuth2Auth) Authenticate(req *http.Request) error {
	req.Header.Set("X-TC-APP-KEY", a.appKey)
	req.Header.Set("X-TC-AUTHENTICATION-ID", a.accessKey)
	req.Header.Set("X-TC-AUTHENTICATION-SECRET", a.secretKey)
	return nil
}
