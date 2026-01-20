package securitygroup

import "errors"

var (
	ErrNoCredentials = errors.New("identity credentials required for security group service")
)
