// Package colocationgw provides Colocation Gateway service types and client
package colocationgw

// ColocationGateway represents a colocation gateway
type ColocationGateway struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	TenantID        string `json:"tenant_id"`
	Description     string `json:"description,omitempty"`
	Status          string `json:"status"`
	RouterID        string `json:"router_id,omitempty"`
	SubnetID        string `json:"subnet_id,omitempty"`
	NetworkID       string `json:"network_id,omitempty"`
	LocalIPAddress  string `json:"local_ip_address,omitempty"`
	RemoteIPAddress string `json:"remote_ip_address,omitempty"`
	VLANID          int    `json:"vlan_id,omitempty"`
	ConnectionType  string `json:"connection_type,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

// ListOutput represents the response from List operation
type ListOutput struct {
	ColocationGateways []ColocationGateway `json:"colocationgateways"`
}

// GetOutput represents the response from Get operation
type GetOutput struct {
	ColocationGateway ColocationGateway `json:"colocationgateway"`
}
