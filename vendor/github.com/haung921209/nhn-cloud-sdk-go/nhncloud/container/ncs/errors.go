package ncs

import "errors"

var (
	ErrNoCredentials    = errors.New("ncs: credentials required")
	ErrNoAppKey         = errors.New("ncs: app key required")
	ErrWorkloadNotFound = errors.New("ncs: workload not found")
	ErrServiceNotFound  = errors.New("ncs: service not found")
	ErrInvalidInput     = errors.New("ncs: invalid input")
)
