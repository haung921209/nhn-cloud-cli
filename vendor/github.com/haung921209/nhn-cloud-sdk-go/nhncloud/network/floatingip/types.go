package floatingip

// FloatingIP represents a floating IP address
type FloatingIP struct {
	ID                string  `json:"id"`
	FloatingNetworkID string  `json:"floating_network_id"`
	FloatingIPAddress string  `json:"floating_ip_address"`
	FixedIPAddress    string  `json:"fixed_ip_address,omitempty"`
	PortID            *string `json:"port_id,omitempty"`
	TenantID          string  `json:"tenant_id"`
	Status            string  `json:"status"`
	Description       string  `json:"description,omitempty"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// ListFloatingIPsOutput represents the response for listing floating IPs
type ListFloatingIPsOutput struct {
	FloatingIPs []FloatingIP `json:"floatingips"`
}

// GetFloatingIPOutput represents the response for getting a single floating IP
type GetFloatingIPOutput struct {
	FloatingIP FloatingIP `json:"floatingip"`
}

// CreateFloatingIPInput represents the input for creating a floating IP
type CreateFloatingIPInput struct {
	FloatingNetworkID string `json:"floating_network_id"`
	PortID            string `json:"port_id,omitempty"`
	FixedIPAddress    string `json:"fixed_ip_address,omitempty"`
	SubnetID          string `json:"subnet_id,omitempty"`
	Description       string `json:"description,omitempty"`
}

// CreateFloatingIPRequest wraps the create input
type CreateFloatingIPRequest struct {
	FloatingIP CreateFloatingIPInput `json:"floatingip"`
}

// CreateFloatingIPOutput represents the response for creating a floating IP
type CreateFloatingIPOutput struct {
	FloatingIP FloatingIP `json:"floatingip"`
}

// UpdateFloatingIPInput represents the input for updating a floating IP
type UpdateFloatingIPInput struct {
	PortID         *string `json:"port_id"`
	FixedIPAddress string  `json:"fixed_ip_address,omitempty"`
	Description    string  `json:"description,omitempty"`
}

// UpdateFloatingIPRequest wraps the update input
type UpdateFloatingIPRequest struct {
	FloatingIP UpdateFloatingIPInput `json:"floatingip"`
}

// UpdateFloatingIPOutput represents the response for updating a floating IP
type UpdateFloatingIPOutput struct {
	FloatingIP FloatingIP `json:"floatingip"`
}

// AssociateFloatingIPInput represents the input for associating a floating IP
type AssociateFloatingIPInput struct {
	PortID         string `json:"port_id"`
	FixedIPAddress string `json:"fixed_ip_address,omitempty"`
}

// DisassociateFloatingIPInput represents the input for disassociating a floating IP
type DisassociateFloatingIPInput struct {
	PortID *string `json:"port_id"` // nil to disassociate
}
