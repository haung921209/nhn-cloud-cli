package natgateway

type NATGateway struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description,omitempty"`
	TenantID          string `json:"tenant_id"`
	SubnetID          string `json:"subnet_id,omitempty"`
	VPCID             string `json:"vpc_id,omitempty"`
	FloatingIPID      string `json:"floatingip_id,omitempty"`
	FloatingIPAddress string `json:"floatingip_address,omitempty"`
	Status            string `json:"status,omitempty"`
	State             string `json:"state,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty"`
}

type ListNATGatewaysOutput struct {
	NATGateways []NATGateway `json:"natgateways"`
}

type GetNATGatewayOutput struct {
	NATGateway NATGateway `json:"natgateway"`
}

type CreateNATGatewayInput struct {
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	VPCID        string `json:"vpc_id"`
	SubnetID     string `json:"subnet_id"`
	FloatingIPID string `json:"floatingips_id,omitempty"`
}

type CreateNATGatewayRequest struct {
	NATGateway CreateNATGatewayInput `json:"natgateway"`
}

type UpdateNATGatewayInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateNATGatewayRequest struct {
	NATGateway UpdateNATGatewayInput `json:"natgateway"`
}
