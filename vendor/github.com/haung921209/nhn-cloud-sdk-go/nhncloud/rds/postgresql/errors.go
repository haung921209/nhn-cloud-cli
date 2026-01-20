package postgresql

import "errors"

var (
	ErrNoCredentials    = errors.New("postgresql: credentials required")
	ErrNoAppKey         = errors.New("postgresql: app key required")
	ErrInstanceNotFound = errors.New("postgresql: instance not found")
	ErrInvalidInput     = errors.New("postgresql: invalid input")
)
