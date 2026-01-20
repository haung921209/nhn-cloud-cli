package floatingip

import "errors"

var (
	ErrNoCredentials = errors.New("identity credentials required for floating ip service")
)
