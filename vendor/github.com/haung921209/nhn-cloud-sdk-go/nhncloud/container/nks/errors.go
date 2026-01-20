package nks

import "errors"

var (
	ErrNoCredentials = errors.New("identity credentials required for nks service")
)
