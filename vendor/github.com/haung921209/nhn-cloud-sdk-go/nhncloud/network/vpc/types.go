package vpc

type VPC struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	TenantID       string `json:"tenant_id"`
	CIDRv4         string `json:"cidrv4"`
	State          string `json:"state"`
	Shared         bool   `json:"shared"`
	RouterExternal bool   `json:"router:external"`
	CreatedAt      string `json:"created_time,omitempty"`
	UpdatedAt      string `json:"updated_time,omitempty"`
}

type Subnet struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	TenantID        string           `json:"tenant_id"`
	NetworkID       string           `json:"network_id"`
	CIDR            string           `json:"cidr"`
	GatewayIP       string           `json:"gateway_ip,omitempty"`
	IPVersion       int              `json:"ip_version"`
	EnableDHCP      bool             `json:"enable_dhcp"`
	DNSNameservers  []string         `json:"dns_nameservers,omitempty"`
	AllocationPools []AllocationPool `json:"allocation_pools,omitempty"`
}

type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type RoutingTable struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	TenantID     string  `json:"tenant_id"`
	VPCID        string  `json:"vpc_id"`
	DefaultTable bool    `json:"default_table"`
	Routes       []Route `json:"routes,omitempty"`
}

type Route struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	Destination string `json:"cidr"`
	Gateway     string `json:"gateway"`
}

type ListVPCsOutput struct {
	VPCs []VPC `json:"vpcs"`
}

type GetVPCOutput struct {
	VPC VPC `json:"vpc"`
}

type CreateVPCInput struct {
	Name   string `json:"name"`
	CIDRv4 string `json:"cidrv4"`
}

type CreateVPCOutput struct {
	VPC VPC `json:"vpc"`
}

type UpdateVPCInput struct {
	Name string `json:"name,omitempty"`
}

type ListSubnetsOutput struct {
	Subnets []Subnet `json:"subnets"`
}

type GetSubnetOutput struct {
	Subnet Subnet `json:"subnet"`
}

type CreateSubnetInput struct {
	Name      string `json:"name"`
	VPCID     string `json:"vpc_id"`
	CIDR      string `json:"cidr"`
	GatewayIP string `json:"gateway,omitempty"`
}

type CreateSubnetOutput struct {
	VPCSubnet Subnet `json:"vpcsubnet"`
}

type ListRoutingTablesOutput struct {
	RoutingTables []RoutingTable `json:"routingtables"`
}

type GetRoutingTableOutput struct {
	RoutingTable RoutingTable `json:"routingtable"`
}
