package object

import "errors"

var (
	ErrNoCredentials = errors.New("identity credentials required for object storage service")
)
