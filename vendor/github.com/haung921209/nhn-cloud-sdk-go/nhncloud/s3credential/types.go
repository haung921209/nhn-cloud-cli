// Package s3credential provides S3 Credential service types and client
package s3credential

import "time"

// S3Credential represents S3 API credential
type S3Credential struct {
	Access     string    `json:"access"`
	Secret     string    `json:"secret"`
	UserID     string    `json:"user_id"`
	TenantID   string    `json:"tenant_id"`
	CreatedAt  time.Time `json:"created_at"`
	AccessedAt time.Time `json:"accessed_at"`
}

// ListCredentialsOutput represents the response for list credentials
type ListCredentialsOutput struct {
	Credentials []S3Credential `json:"credentials"`
}

// CredentialOutput represents the response for single credential
type CredentialOutput struct {
	Credential S3Credential `json:"credential"`
}

// CreateCredentialInput represents request to create credential
type CreateCredentialInput struct {
	TenantID string `json:"tenant_id"`
}
