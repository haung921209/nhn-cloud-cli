package servicegateway

// ServiceGateway represents a service gateway for accessing NHN Cloud services privately
type ServiceGateway struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	TenantID          string `json:"tenant_id"`
	SubnetID          string `json:"subnet_id"`
	ServiceEndpointID string `json:"service_endpoint_id"`
	IPAddress         string `json:"ip_address"`
	Status            string `json:"status"`
	CreateTime        string `json:"create_time,omitempty"`
	UpdateTime        string `json:"update_time,omitempty"`
}

// ServiceEndpoint represents a predefined NHN Cloud service endpoint
type ServiceEndpoint struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ServiceName  string `json:"service_name"`
	Description  string `json:"description"`
	Region       string `json:"region"`
	EndpointType string `json:"endpoint_type"`
	CreateTime   string `json:"create_time,omitempty"`
}

// ListServiceGatewaysOutput represents the response for service gateways list API
type ListServiceGatewaysOutput struct {
	ServiceGateways []ServiceGateway `json:"servicegateways"`
}

// GetServiceGatewayOutput represents the response for single service gateway API
type GetServiceGatewayOutput struct {
	ServiceGateway ServiceGateway `json:"servicegateway"`
}

// CreateServiceGatewayInput represents the service gateway creation parameters
type CreateServiceGatewayInput struct {
	Name              string `json:"name"`
	Description       string `json:"description,omitempty"`
	SubnetID          string `json:"subnet_id"`
	ServiceEndpointID string `json:"service_endpoint_id"`
}

// CreateServiceGatewayRequest wraps the create input
type CreateServiceGatewayRequest struct {
	ServiceGateway CreateServiceGatewayInput `json:"servicegateway"`
}

// UpdateServiceGatewayInput represents the service gateway update parameters
type UpdateServiceGatewayInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateServiceGatewayRequest wraps the update input
type UpdateServiceGatewayRequest struct {
	ServiceGateway UpdateServiceGatewayInput `json:"servicegateway"`
}

// ListServiceEndpointsOutput represents the response for service endpoints list API
type ListServiceEndpointsOutput struct {
	ServiceEndpoints []ServiceEndpoint `json:"serviceendpoints"`
}

// GetServiceEndpointOutput represents the response for single service endpoint API
type GetServiceEndpointOutput struct {
	ServiceEndpoint ServiceEndpoint `json:"serviceendpoint"`
}
