// Package colocationgw provides Colocation Gateway service types and client
package colocationgw

// VPC represents the nested VPC information in a colocation gateway response
type VPC struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	CIDRv4 string `json:"cidrv4"`
}

// ColocationGateway represents a colocation gateway
type ColocationGateway struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	TenantID    string `json:"tenant_id"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	VPCID       string `json:"vpc_id"`
	VPC         VPC    `json:"vpc"`
}

// ListOutput represents the response from List operation
type ListOutput struct {
	ColocationGateways []ColocationGateway `json:"colocationgateways"`
}

// GetOutput represents the response from Get operation
type GetOutput struct {
	ColocationGateway ColocationGateway `json:"colocationgateway"`
}
