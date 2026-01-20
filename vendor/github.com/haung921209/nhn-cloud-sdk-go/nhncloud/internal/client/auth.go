package client

import (
	"context"
	"net/http"
)

type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
	SetAuthHeader(req *http.Request, token string)
}

type TokenType string

const (
	TokenTypeOAuth    TokenType = "oauth"
	TokenTypeIdentity TokenType = "identity"
)
