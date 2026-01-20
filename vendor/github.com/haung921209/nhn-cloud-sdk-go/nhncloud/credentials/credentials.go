package credentials

// Credentials provides OAuth credentials for NHN Cloud API authentication.
// Used for IAM and App Key based services.
type Credentials interface {
	GetAccessKeyID() string
	GetSecretAccessKey() string
}

// IdentityCredentials provides OpenStack Identity credentials.
// Used for Compute, Network, and Storage services.
type IdentityCredentials interface {
	GetUsername() string
	GetPassword() string
	GetTenantID() string
}
