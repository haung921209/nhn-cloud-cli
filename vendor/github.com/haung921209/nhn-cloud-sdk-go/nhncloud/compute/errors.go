package compute

import "errors"

var (
	ErrNoCredentials = errors.New("identity credentials required for compute service")
)
