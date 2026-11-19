package ncs

type ResponseHeader struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// Workload is a workload as the API returns it.
type Workload struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description,omitempty"`
	Type            string         `json:"type"`
	TemplateID      string         `json:"templateId"`
	TemplateName    string         `json:"templateName,omitempty"`
	TemplateVersion string         `json:"templateVersion"`
	Desired         int            `json:"desired"`
	Available       int            `json:"available"`
	Status          string         `json:"status"`
	LoadBalancing   *LoadBalancing `json:"loadBalancing,omitempty"`
	CreatedAt       string         `json:"createdAt"`
	UpdatedAt       string         `json:"updatedAt"`
}

// LoadBalancing is the workload's load balancer setting.
type LoadBalancing struct {
	Enabled    bool   `json:"enabled"`
	FloatingIP bool   `json:"floatingIp"`
	IPAddress  string `json:"ipAddress,omitempty"`
}

type Container struct {
	Name         string                `json:"name"`
	Image        string                `json:"image"`
	Command      []string              `json:"command,omitempty"`
	Args         []string              `json:"args,omitempty"`
	Ports        []ContainerPort       `json:"ports,omitempty"`
	Env          []EnvVar              `json:"env,omitempty"`
	Resources    *ResourceRequirements `json:"resources,omitempty"`
	VolumeMounts []VolumeMount         `json:"volumeMounts,omitempty"`
}

type ContainerPort struct {
	Name          string `json:"name,omitempty"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type ResourceRequirements struct {
	Limits   ResourceList `json:"limits,omitempty"`
	Requests ResourceList `json:"requests,omitempty"`
}

type ResourceList struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
	GPU    string `json:"nvidia.com/gpu,omitempty"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

type ListWorkloadsOutput struct {
	Header    *ResponseHeader `json:"header"`
	Workloads []Workload      `json:"workloads"`
}

type GetWorkloadOutput struct {
	Header   *ResponseHeader `json:"header"`
	Workload Workload        `json:"workload"`
}

// CreateWorkloadInput is sent as {"workload": {...}}. The containers and
// network come from the template.
type CreateWorkloadInput struct {
	Name            string         `json:"name"`
	Description     string         `json:"description,omitempty"`
	Type            string         `json:"type,omitempty"`
	TemplateID      string         `json:"templateId"`
	TemplateVersion string         `json:"templateVersion,omitempty"`
	Desired         int            `json:"desired"`
	LoadBalancing   *LoadBalancing `json:"loadBalancing"`
}

type CreateWorkloadOutput struct {
	Header   *ResponseHeader `json:"header"`
	Workload Workload        `json:"workload"`
}

type UpdateWorkloadInput struct {
	Description string      `json:"description,omitempty"`
	Replicas    int         `json:"replicas,omitempty"`
	Containers  []Container `json:"containers,omitempty"`
}

type Template struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description,omitempty"`
	Version      string      `json:"version"`
	VersionCount int         `json:"versionCount,omitempty"`
	Containers   []Container `json:"containers,omitempty"`
	CreatedAt    string      `json:"createdAt"`
	UpdatedAt    string      `json:"updatedAt"`
}

type ListTemplatesOutput struct {
	Header    *ResponseHeader `json:"header"`
	Templates []Template      `json:"templates"`
}

type GetTemplateOutput struct {
	Header   *ResponseHeader `json:"header"`
	Template Template        `json:"template"`
}

type CreateTemplateInput struct {
	Template Template `json:"template"`
}

type CreateTemplateOutput struct {
	Header   *ResponseHeader `json:"header"`
	Template Template        `json:"template"`
}

// Event types for workload events
type Event struct {
	EventID    string `json:"eventId"`
	EventType  string `json:"eventType"` // Normal, Warning
	Reason     string `json:"reason"`
	Message    string `json:"message"`
	Count      int    `json:"count"`
	FirstTime  string `json:"firstTimestamp"`
	LastTime   string `json:"lastTimestamp"`
	Source     string `json:"source,omitempty"`
	ObjectKind string `json:"objectKind,omitempty"`
	ObjectName string `json:"objectName,omitempty"`
}

type GetWorkloadEventsOutput struct {
	Header *ResponseHeader `json:"header"`
	Events []Event         `json:"events"`
}

type NetworkPolicyRule struct {
	Direction string   `json:"direction"` // Ingress, Egress
	Protocol  string   `json:"protocol"`  // TCP, UDP, ICMP
	Port      int      `json:"port,omitempty"`
	FromCIDR  []string `json:"fromCidr,omitempty"`
	ToCIDR    []string `json:"toCidr,omitempty"`
	Action    string   `json:"action"` // Allow, Deny
}

type NetworkPolicy struct {
	PolicyID    string              `json:"policyId"`
	PolicyName  string              `json:"policyName"`
	Description string              `json:"description,omitempty"`
	Rules       []NetworkPolicyRule `json:"rules"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
}

type ListNetworkPoliciesOutput struct {
	Header   *ResponseHeader `json:"header"`
	Policies []NetworkPolicy `json:"policies"`
}

type GetNetworkPolicyOutput struct {
	Header *ResponseHeader `json:"header"`
	Policy NetworkPolicy   `json:"policy"`
}

type CreateNetworkPolicyInput struct {
	PolicyName  string              `json:"policyName"`
	Description string              `json:"description,omitempty"`
	Rules       []NetworkPolicyRule `json:"rules"`
}

type UpdateNetworkPolicyInput struct {
	PolicyName  string              `json:"policyName,omitempty"`
	Description string              `json:"description,omitempty"`
	Rules       []NetworkPolicyRule `json:"rules,omitempty"`
}
