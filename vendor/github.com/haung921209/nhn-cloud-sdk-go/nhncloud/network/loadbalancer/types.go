package loadbalancer

type IDReference struct {
	ID string `json:"id"`
}

type LoadBalancer struct {
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	Description        string        `json:"description,omitempty"`
	TenantID           string        `json:"tenant_id"`
	VIPAddress         string        `json:"vip_address"`
	VIPPortID          string        `json:"vip_port_id"`
	VIPSubnetID        string        `json:"vip_subnet_id"`
	VIPNetworkID       string        `json:"vip_network_id"`
	ProvisioningStatus string        `json:"provisioning_status"`
	OperatingStatus    string        `json:"operating_status"`
	AdminStateUp       bool          `json:"admin_state_up"`
	Provider           string        `json:"provider"`
	Listeners          []IDReference `json:"listeners,omitempty"`
	Pools              []IDReference `json:"pools,omitempty"`
	CreatedAt          string        `json:"created_at"`
	UpdatedAt          string        `json:"updated_at"`
}

// ListLoadBalancersOutput represents the response for listing load balancers
type ListLoadBalancersOutput struct {
	LoadBalancers []LoadBalancer `json:"loadbalancers"`
}

// GetLoadBalancerOutput represents the response for getting a single load balancer
type GetLoadBalancerOutput struct {
	LoadBalancer LoadBalancer `json:"loadbalancer"`
}

// CreateLoadBalancerInput represents the input for creating a load balancer
type CreateLoadBalancerInput struct {
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	VIPSubnetID  string `json:"vip_subnet_id"`
	VIPAddress   string `json:"vip_address,omitempty"`
	AdminStateUp *bool  `json:"admin_state_up,omitempty"`
	Provider     string `json:"provider,omitempty"`
}

// CreateLoadBalancerRequest wraps the create input
type CreateLoadBalancerRequest struct {
	LoadBalancer CreateLoadBalancerInput `json:"loadbalancer"`
}

// CreateLoadBalancerOutput represents the response for creating a load balancer
type CreateLoadBalancerOutput struct {
	LoadBalancer LoadBalancer `json:"loadbalancer"`
}

// UpdateLoadBalancerInput represents the input for updating a load balancer
type UpdateLoadBalancerInput struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	AdminStateUp *bool  `json:"admin_state_up,omitempty"`
}

// UpdateLoadBalancerRequest wraps the update input
type UpdateLoadBalancerRequest struct {
	LoadBalancer UpdateLoadBalancerInput `json:"loadbalancer"`
}

// Listener represents a load balancer listener
type Listener struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	Description            string   `json:"description,omitempty"`
	TenantID               string   `json:"tenant_id"`
	LoadBalancerID         string   `json:"loadbalancer_id"`
	Protocol               string   `json:"protocol"`
	ProtocolPort           int      `json:"protocol_port"`
	DefaultPoolID          string   `json:"default_pool_id,omitempty"`
	ConnectionLimit        int      `json:"connection_limit"`
	AdminStateUp           bool     `json:"admin_state_up"`
	ProvisioningStatus     string   `json:"provisioning_status"`
	OperatingStatus        string   `json:"operating_status"`
	DefaultTLSContainerRef string   `json:"default_tls_container_ref,omitempty"`
	SNIContainerRefs       []string `json:"sni_container_refs,omitempty"`
	CreatedAt              string   `json:"created_at"`
	UpdatedAt              string   `json:"updated_at"`
}

// ListListenersOutput represents the response for listing listeners
type ListListenersOutput struct {
	Listeners []Listener `json:"listeners"`
}

// GetListenerOutput represents the response for getting a single listener
type GetListenerOutput struct {
	Listener Listener `json:"listener"`
}

// CreateListenerInput represents the input for creating a listener
type CreateListenerInput struct {
	Name                   string `json:"name"`
	Description            string `json:"description,omitempty"`
	LoadBalancerID         string `json:"loadbalancer_id"`
	Protocol               string `json:"protocol"`
	ProtocolPort           int    `json:"protocol_port"`
	DefaultPoolID          string `json:"default_pool_id,omitempty"`
	ConnectionLimit        int    `json:"connection_limit,omitempty"`
	AdminStateUp           *bool  `json:"admin_state_up,omitempty"`
	DefaultTLSContainerRef string `json:"default_tls_container_ref,omitempty"`
}

// CreateListenerRequest wraps the create input
type CreateListenerRequest struct {
	Listener CreateListenerInput `json:"listener"`
}

// Pool represents a load balancer pool
type Pool struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Description        string              `json:"description,omitempty"`
	TenantID           string              `json:"tenant_id"`
	Protocol           string              `json:"protocol"`
	LBAlgorithm        string              `json:"lb_algorithm"`
	AdminStateUp       bool                `json:"admin_state_up"`
	ProvisioningStatus string              `json:"provisioning_status"`
	OperatingStatus    string              `json:"operating_status"`
	LoadBalancerID     string              `json:"loadbalancer_id,omitempty"`
	ListenerID         string              `json:"listener_id,omitempty"`
	Members            []Member            `json:"members,omitempty"`
	HealthMonitorID    string              `json:"healthmonitor_id,omitempty"`
	SessionPersistence *SessionPersistence `json:"session_persistence,omitempty"`
	CreatedAt          string              `json:"created_at"`
	UpdatedAt          string              `json:"updated_at"`
}

// SessionPersistence represents session persistence configuration
type SessionPersistence struct {
	Type       string `json:"type"`
	CookieName string `json:"cookie_name,omitempty"`
}

// ListPoolsOutput represents the response for listing pools
type ListPoolsOutput struct {
	Pools []Pool `json:"pools"`
}

// GetPoolOutput represents the response for getting a single pool
type GetPoolOutput struct {
	Pool Pool `json:"pool"`
}

// CreatePoolInput represents the input for creating a pool
type CreatePoolInput struct {
	Name               string              `json:"name"`
	Description        string              `json:"description,omitempty"`
	Protocol           string              `json:"protocol"`
	LBAlgorithm        string              `json:"lb_algorithm"`
	LoadBalancerID     string              `json:"loadbalancer_id,omitempty"`
	ListenerID         string              `json:"listener_id,omitempty"`
	AdminStateUp       *bool               `json:"admin_state_up,omitempty"`
	SessionPersistence *SessionPersistence `json:"session_persistence,omitempty"`
}

// CreatePoolRequest wraps the create input
type CreatePoolRequest struct {
	Pool CreatePoolInput `json:"pool"`
}

// Member represents a pool member
type Member struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	TenantID           string `json:"tenant_id"`
	Address            string `json:"address"`
	ProtocolPort       int    `json:"protocol_port"`
	Weight             int    `json:"weight"`
	SubnetID           string `json:"subnet_id"`
	AdminStateUp       bool   `json:"admin_state_up"`
	ProvisioningStatus string `json:"provisioning_status"`
	OperatingStatus    string `json:"operating_status"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// ListMembersOutput represents the response for listing members
type ListMembersOutput struct {
	Members []Member `json:"members"`
}

// GetMemberOutput represents the response for getting a single member
type GetMemberOutput struct {
	Member Member `json:"member"`
}

// CreateMemberInput represents the input for creating a member
type CreateMemberInput struct {
	Name         string `json:"name,omitempty"`
	Address      string `json:"address"`
	ProtocolPort int    `json:"protocol_port"`
	Weight       int    `json:"weight,omitempty"`
	SubnetID     string `json:"subnet_id,omitempty"`
	AdminStateUp *bool  `json:"admin_state_up,omitempty"`
}

// CreateMemberRequest wraps the create input
type CreateMemberRequest struct {
	Member CreateMemberInput `json:"member"`
}

// HealthMonitor represents a health monitor
type HealthMonitor struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	TenantID           string `json:"tenant_id"`
	PoolID             string `json:"pool_id"`
	Type               string `json:"type"`
	Delay              int    `json:"delay"`
	Timeout            int    `json:"timeout"`
	MaxRetries         int    `json:"max_retries"`
	MaxRetriesDown     int    `json:"max_retries_down"`
	HTTPMethod         string `json:"http_method,omitempty"`
	URLPath            string `json:"url_path,omitempty"`
	ExpectedCodes      string `json:"expected_codes,omitempty"`
	AdminStateUp       bool   `json:"admin_state_up"`
	ProvisioningStatus string `json:"provisioning_status"`
	OperatingStatus    string `json:"operating_status"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// ListHealthMonitorsOutput represents the response for listing health monitors
type ListHealthMonitorsOutput struct {
	HealthMonitors []HealthMonitor `json:"healthmonitors"`
}

// GetHealthMonitorOutput represents the response for getting a single health monitor
type GetHealthMonitorOutput struct {
	HealthMonitor HealthMonitor `json:"healthmonitor"`
}

// CreateHealthMonitorInput represents the input for creating a health monitor
type CreateHealthMonitorInput struct {
	Name           string `json:"name,omitempty"`
	PoolID         string `json:"pool_id"`
	Type           string `json:"type"`
	Delay          int    `json:"delay"`
	Timeout        int    `json:"timeout"`
	MaxRetries     int    `json:"max_retries"`
	MaxRetriesDown int    `json:"max_retries_down,omitempty"`
	HTTPMethod     string `json:"http_method,omitempty"`
	URLPath        string `json:"url_path,omitempty"`
	ExpectedCodes  string `json:"expected_codes,omitempty"`
	AdminStateUp   *bool  `json:"admin_state_up,omitempty"`
}

type CreateHealthMonitorRequest struct {
	HealthMonitor CreateHealthMonitorInput `json:"healthmonitor"`
}

type UpdateListenerInput struct {
	Name                   string `json:"name,omitempty"`
	Description            string `json:"description,omitempty"`
	DefaultPoolID          string `json:"default_pool_id,omitempty"`
	ConnectionLimit        *int   `json:"connection_limit,omitempty"`
	AdminStateUp           *bool  `json:"admin_state_up,omitempty"`
	DefaultTLSContainerRef string `json:"default_tls_container_ref,omitempty"`
}

type UpdateListenerRequest struct {
	Listener UpdateListenerInput `json:"listener"`
}

type UpdatePoolInput struct {
	Name               string              `json:"name,omitempty"`
	Description        string              `json:"description,omitempty"`
	LBAlgorithm        string              `json:"lb_algorithm,omitempty"`
	AdminStateUp       *bool               `json:"admin_state_up,omitempty"`
	SessionPersistence *SessionPersistence `json:"session_persistence,omitempty"`
}

type UpdatePoolRequest struct {
	Pool UpdatePoolInput `json:"pool"`
}

type UpdateMemberInput struct {
	Name         string `json:"name,omitempty"`
	Weight       *int   `json:"weight,omitempty"`
	AdminStateUp *bool  `json:"admin_state_up,omitempty"`
}

type UpdateMemberRequest struct {
	Member UpdateMemberInput `json:"member"`
}

type UpdateHealthMonitorInput struct {
	Name           string `json:"name,omitempty"`
	Delay          *int   `json:"delay,omitempty"`
	Timeout        *int   `json:"timeout,omitempty"`
	MaxRetries     *int   `json:"max_retries,omitempty"`
	MaxRetriesDown *int   `json:"max_retries_down,omitempty"`
	HTTPMethod     string `json:"http_method,omitempty"`
	URLPath        string `json:"url_path,omitempty"`
	ExpectedCodes  string `json:"expected_codes,omitempty"`
	AdminStateUp   *bool  `json:"admin_state_up,omitempty"`
}

type UpdateHealthMonitorRequest struct {
	HealthMonitor UpdateHealthMonitorInput `json:"healthmonitor"`
}

type L7Policy struct {
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	Description        string        `json:"description,omitempty"`
	TenantID           string        `json:"tenant_id"`
	ListenerID         string        `json:"listener_id"`
	Action             string        `json:"action"`
	Position           int           `json:"position"`
	RedirectPoolID     string        `json:"redirect_pool_id,omitempty"`
	RedirectURL        string        `json:"redirect_url,omitempty"`
	RedirectPrefix     string        `json:"redirect_prefix,omitempty"`
	RedirectHTTPCode   int           `json:"redirect_http_code,omitempty"`
	AdminStateUp       bool          `json:"admin_state_up"`
	ProvisioningStatus string        `json:"provisioning_status"`
	OperatingStatus    string        `json:"operating_status"`
	Rules              []IDReference `json:"rules,omitempty"`
	CreatedAt          string        `json:"created_at"`
	UpdatedAt          string        `json:"updated_at"`
}

type ListL7PoliciesOutput struct {
	L7Policies []L7Policy `json:"l7policies"`
}

type GetL7PolicyOutput struct {
	L7Policy L7Policy `json:"l7policy"`
}

type CreateL7PolicyInput struct {
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	ListenerID       string `json:"listener_id"`
	Action           string `json:"action"`
	Position         *int   `json:"position,omitempty"`
	RedirectPoolID   string `json:"redirect_pool_id,omitempty"`
	RedirectURL      string `json:"redirect_url,omitempty"`
	RedirectPrefix   string `json:"redirect_prefix,omitempty"`
	RedirectHTTPCode *int   `json:"redirect_http_code,omitempty"`
	AdminStateUp     *bool  `json:"admin_state_up,omitempty"`
}

type CreateL7PolicyRequest struct {
	L7Policy CreateL7PolicyInput `json:"l7policy"`
}

type UpdateL7PolicyInput struct {
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	Action           string `json:"action,omitempty"`
	Position         *int   `json:"position,omitempty"`
	RedirectPoolID   string `json:"redirect_pool_id,omitempty"`
	RedirectURL      string `json:"redirect_url,omitempty"`
	RedirectPrefix   string `json:"redirect_prefix,omitempty"`
	RedirectHTTPCode *int   `json:"redirect_http_code,omitempty"`
	AdminStateUp     *bool  `json:"admin_state_up,omitempty"`
}

type UpdateL7PolicyRequest struct {
	L7Policy UpdateL7PolicyInput `json:"l7policy"`
}

type L7Rule struct {
	ID                 string `json:"id"`
	TenantID           string `json:"tenant_id"`
	Type               string `json:"type"`
	CompareType        string `json:"compare_type"`
	Key                string `json:"key,omitempty"`
	Value              string `json:"value"`
	Invert             bool   `json:"invert"`
	AdminStateUp       bool   `json:"admin_state_up"`
	ProvisioningStatus string `json:"provisioning_status"`
	OperatingStatus    string `json:"operating_status"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type ListL7RulesOutput struct {
	Rules []L7Rule `json:"rules"`
}

type GetL7RuleOutput struct {
	Rule L7Rule `json:"rule"`
}

type CreateL7RuleInput struct {
	Type         string `json:"type"`
	CompareType  string `json:"compare_type"`
	Key          string `json:"key,omitempty"`
	Value        string `json:"value"`
	Invert       *bool  `json:"invert,omitempty"`
	AdminStateUp *bool  `json:"admin_state_up,omitempty"`
}

type CreateL7RuleRequest struct {
	Rule CreateL7RuleInput `json:"rule"`
}

type UpdateL7RuleInput struct {
	Type         string `json:"type,omitempty"`
	CompareType  string `json:"compare_type,omitempty"`
	Key          string `json:"key,omitempty"`
	Value        string `json:"value,omitempty"`
	Invert       *bool  `json:"invert,omitempty"`
	AdminStateUp *bool  `json:"admin_state_up,omitempty"`
}

type UpdateL7RuleRequest struct {
	Rule UpdateL7RuleInput `json:"rule"`
}
