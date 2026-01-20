package loadbalancer

import "errors"

var (
	ErrNoCredentials = errors.New("identity credentials required for load balancer service")
)
