package iam

import "errors"

var (
	ErrNoCredentials        = errors.New("iam: credentials required")
	ErrOrganizationNotFound = errors.New("iam: organization not found")
	ErrProjectNotFound      = errors.New("iam: project not found")
	ErrMemberNotFound       = errors.New("iam: member not found")
	ErrInvalidInput         = errors.New("iam: invalid input")
)
