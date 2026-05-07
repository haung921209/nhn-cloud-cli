package port

// Port represents a Neutron network port
type Port struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	NetworkID      string    `json:"network_id"`
	TenantID       string    `json:"tenant_id"`
	MACAddress     string    `json:"mac_address"`
	Status         string    `json:"status"`
	DeviceID       string    `json:"device_id"`
	DeviceOwner    string    `json:"device_owner"`
	FixedIPs       []FixedIP `json:"fixed_ips"`
	SecurityGroups []string  `json:"security_groups"`
}

// FixedIP represents a fixed IP assigned to a port
type FixedIP struct {
	SubnetID  string `json:"subnet_id"`
	IPAddress string `json:"ip_address"`
}

// ListPortsOutput represents the response from listing ports
type ListPortsOutput struct {
	Ports []Port `json:"ports"`
}

// GetPortOutput represents the response from getting a single port
type GetPortOutput struct {
	Port Port `json:"port"`
}
