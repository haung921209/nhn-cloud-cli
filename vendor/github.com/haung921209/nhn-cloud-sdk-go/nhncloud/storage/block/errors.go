package block

import "errors"

var (
	ErrNoCredentials = errors.New("identity credentials required for block storage service")
)
