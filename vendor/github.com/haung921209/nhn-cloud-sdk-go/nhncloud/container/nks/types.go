package nks

type Cluster struct {
	ID                string            `json:"uuid"`
	Name              string            `json:"name"`
	TenantID          string            `json:"tenant_id"`
	Status            string            `json:"status"`
	StatusReason      string            `json:"status_reason,omitempty"`
	HealthStatus      string            `json:"health_status,omitempty"`
	APIAddress        string            `json:"api_address"`
	MasterAddresses   []string          `json:"master_addresses"`
	NodeAddresses     []string          `json:"node_addresses,omitempty"`
	K8sVersion        string            `json:"coe_version"`
	NodeCount         int               `json:"node_count"`
	MasterCount       int               `json:"master_count"`
	FlavorID          string            `json:"flavor_id"`
	MasterFlavorID    string            `json:"master_flavor_id"`
	KeyPair           string            `json:"keypair"`
	ClusterTemplateID string            `json:"cluster_template_id"`
	NetworkID         string            `json:"fixed_network,omitempty"`
	SubnetID          string            `json:"fixed_subnet,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	CreatedAt         string            `json:"created_at"`
	UpdatedAt         string            `json:"updated_at"`
}

type ListClustersOutput struct {
	Clusters []Cluster `json:"clusters"`
}

type GetClusterOutput struct {
	Cluster
}

type CreateClusterInput struct {
	Name              string `json:"name"`
	ClusterTemplateID string `json:"cluster_template_id,omitempty"`
	// K8sVersion        string            `json:"coe_version,omitempty"` // Read-only field
	MasterCount    int               `json:"master_count,omitempty"`
	NodeCount      int               `json:"node_count,omitempty"`
	MasterFlavorID string            `json:"master_flavor_id,omitempty"`
	FlavorID       string            `json:"flavor_id,omitempty"`
	KeyPair        string            `json:"keypair,omitempty"`
	NetworkID      string            `json:"fixed_network,omitempty"`
	SubnetID       string            `json:"fixed_subnet,omitempty"`
	Labels         map[string]string `json:"labels,omitempty"`
}

type CreateClusterOutput struct {
	Cluster
}

type UpdateClusterInput struct {
	NodeCount int `json:"node_count,omitempty"`
}

type NodeGroup struct {
	ID           string `json:"uuid"`
	Name         string `json:"name"`
	ClusterID    string `json:"cluster_id"`
	NodeCount    int    `json:"node_count"`
	MinNodeCount int    `json:"min_node_count"`
	MaxNodeCount int    `json:"max_node_count"`
	FlavorID     string `json:"flavor_id"`
	ImageID      string `json:"image_id"`
	Role         string `json:"role"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ListNodeGroupsOutput struct {
	NodeGroups []NodeGroup `json:"nodegroups"`
}

type GetNodeGroupOutput struct {
	NodeGroup
}

type CreateNodeGroupInput struct {
	Name         string `json:"name"`
	FlavorID     string `json:"flavor_id"`
	ImageID      string `json:"image_id,omitempty"`
	NodeCount    int    `json:"node_count"`
	MinNodeCount int    `json:"min_node_count,omitempty"`
	MaxNodeCount int    `json:"max_node_count,omitempty"`
}

type CreateNodeGroupOutput struct {
	NodeGroup
}

type UpdateNodeGroupInput struct {
	NodeCount    int `json:"node_count,omitempty"`
	MinNodeCount int `json:"min_node_count,omitempty"`
	MaxNodeCount int `json:"max_node_count,omitempty"`
}

type ClusterTemplate struct {
	ID                  string            `json:"uuid"`
	Name                string            `json:"name"`
	TenantID            string            `json:"tenant_id"`
	Public              bool              `json:"public"`
	COE                 string            `json:"coe"`
	ServerType          string            `json:"server_type"`
	NetworkDriver       string            `json:"network_driver"`
	VolumeDriver        string            `json:"volume_driver"`
	DockerStorageDriver string            `json:"docker_storage_driver"`
	DockerVolumeSize    int               `json:"docker_volume_size"`
	ExternalNetworkID   string            `json:"external_network_id"`
	FixedNetwork        string            `json:"fixed_network"`
	FixedSubnet         string            `json:"fixed_subnet"`
	DNSNameserver       string            `json:"dns_nameserver"`
	MasterFlavor        string            `json:"master_flavor_id"`
	Flavor              string            `json:"flavor_id"`
	Labels              map[string]string `json:"labels,omitempty"`
	CreatedAt           string            `json:"created_at"`
	UpdatedAt           string            `json:"updated_at"`
}

type ListClusterTemplatesOutput struct {
	ClusterTemplates []ClusterTemplate `json:"clustertemplates"`
}

type GetKubeconfigOutput struct {
	Kubeconfig string
}

type GetSupportedVersionsOutput struct {
	SupportedK8s map[string]bool `json:"supported_k8s"`
}
