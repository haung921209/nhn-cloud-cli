package vpc

import "errors"

var (
	ErrNoCredentials = errors.New("identity credentials required for vpc service")
)
