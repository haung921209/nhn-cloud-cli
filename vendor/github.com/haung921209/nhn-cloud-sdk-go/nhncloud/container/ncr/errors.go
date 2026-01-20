package ncr

import "errors"

var (
	ErrNoCredentials    = errors.New("ncr: credentials required")
	ErrNoAppKey         = errors.New("ncr: app key required")
	ErrRegistryNotFound = errors.New("ncr: registry not found")
	ErrImageNotFound    = errors.New("ncr: image not found")
	ErrInvalidInput     = errors.New("ncr: invalid input")
)
