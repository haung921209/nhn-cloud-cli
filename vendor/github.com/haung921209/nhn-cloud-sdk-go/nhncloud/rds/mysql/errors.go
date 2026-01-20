package mysql

import "errors"

var (
	ErrNoCredentials    = errors.New("mysql: credentials required")
	ErrNoAppKey         = errors.New("mysql: app key required")
	ErrInstanceNotFound = errors.New("mysql: instance not found")
	ErrInvalidInput     = errors.New("mysql: invalid input")
)
