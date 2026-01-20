package credentials

// Static implements Credentials with static values.
type Static struct {
	accessKeyID     string
	secretAccessKey string
}

// NewStatic creates credentials with static access key and secret.
func NewStatic(accessKeyID, secretAccessKey string) *Static {
	return &Static{
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
	}
}

func (s *Static) GetAccessKeyID() string {
	return s.accessKeyID
}

func (s *Static) GetSecretAccessKey() string {
	return s.secretAccessKey
}

// StaticIdentity implements IdentityCredentials with static values.
type StaticIdentity struct {
	username string
	password string
	tenantID string
}

// NewStaticIdentity creates identity credentials with static values.
func NewStaticIdentity(username, password, tenantID string) *StaticIdentity {
	return &StaticIdentity{
		username: username,
		password: password,
		tenantID: tenantID,
	}
}

func (s *StaticIdentity) GetUsername() string {
	return s.username
}

func (s *StaticIdentity) GetPassword() string {
	return s.password
}

func (s *StaticIdentity) GetTenantID() string {
	return s.tenantID
}
