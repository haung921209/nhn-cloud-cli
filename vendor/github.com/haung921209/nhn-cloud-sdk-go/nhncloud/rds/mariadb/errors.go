package mariadb

import "errors"

var (
	ErrNoCredentials    = errors.New("mariadb: credentials required")
	ErrNoAppKey         = errors.New("mariadb: app key required")
	ErrInstanceNotFound = errors.New("mariadb: instance not found")
	ErrInvalidInput     = errors.New("mariadb: invalid input")
)
