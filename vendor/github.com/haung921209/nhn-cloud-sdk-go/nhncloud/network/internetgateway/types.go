package internetgateway

type InternetGateway struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	TenantID          string `json:"tenant_id"`
	RoutingTableID    string `json:"routingtable_id,omitempty"`
	ExternalNetworkID string `json:"external_network_id,omitempty"`
	State             string `json:"state,omitempty"`
	CreateTime        string `json:"create_time,omitempty"`
}

type ListInternetGatewaysOutput struct {
	InternetGateways []InternetGateway `json:"internetgateways"`
}

type GetInternetGatewayOutput struct {
	InternetGateway InternetGateway `json:"internetgateway"`
}

type CreateInternetGatewayInput struct {
	Name              string `json:"name"`
	RoutingTableID    string `json:"routingtable_id"`
	ExternalNetworkID string `json:"external_network_id,omitempty"`
}

type CreateInternetGatewayRequest struct {
	InternetGateway CreateInternetGatewayInput `json:"internetgateway"`
}

type ExternalNetwork struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	RouterExternal bool   `json:"router:external"`
}

type ListExternalNetworksOutput struct {
	Networks []ExternalNetwork `json:"vpcs"`
}
