package credentials

import "os"

const (
	EnvAccessKeyID     = "NHNCLOUD_ACCESS_KEY_ID"
	EnvSecretAccessKey = "NHNCLOUD_SECRET_ACCESS_KEY"
	EnvUsername        = "NHNCLOUD_USERNAME"
	EnvPassword        = "NHNCLOUD_PASSWORD"
	EnvTenantID        = "NHNCLOUD_TENANT_ID"
)

// Env implements Credentials by reading from environment variables.
type Env struct{}

// NewEnv creates credentials that read from environment variables.
func NewEnv() *Env {
	return &Env{}
}

func (e *Env) GetAccessKeyID() string {
	return os.Getenv(EnvAccessKeyID)
}

func (e *Env) GetSecretAccessKey() string {
	return os.Getenv(EnvSecretAccessKey)
}

// EnvIdentity implements IdentityCredentials by reading from environment variables.
type EnvIdentity struct{}

// NewEnvIdentity creates identity credentials that read from environment variables.
func NewEnvIdentity() *EnvIdentity {
	return &EnvIdentity{}
}

func (e *EnvIdentity) GetUsername() string {
	return os.Getenv(EnvUsername)
}

func (e *EnvIdentity) GetPassword() string {
	return os.Getenv(EnvPassword)
}

func (e *EnvIdentity) GetTenantID() string {
	return os.Getenv(EnvTenantID)
}
