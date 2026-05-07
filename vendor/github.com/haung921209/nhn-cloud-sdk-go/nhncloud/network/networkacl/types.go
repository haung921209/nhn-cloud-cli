package networkacl

// ACL represents a network access control list
type ACL struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenant_id"`
	CreateTime  string `json:"create_time,omitempty"`
	UpdateTime  string `json:"update_time,omitempty"`
	Shared      bool   `json:"shared"`
}

// ListACLsOutput represents the response for ACL list API
type ListACLsOutput struct {
	ACLs []ACL `json:"acls"`
}

// GetACLOutput represents the response for single ACL API
type GetACLOutput struct {
	ACL ACL `json:"acl"`
}

// CreateACLInput represents the ACL creation parameters
type CreateACLInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CreateACLRequest represents the request body for creating an ACL
type CreateACLRequest struct {
	ACL CreateACLInput `json:"acl"`
}

// UpdateACLInput represents the ACL update parameters
type UpdateACLInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateACLRequest represents the request body for updating an ACL
type UpdateACLRequest struct {
	ACL UpdateACLInput `json:"acl"`
}

// ACLRule represents a network ACL rule
type ACLRule struct {
	ID          string `json:"id"`
	ACLID       string `json:"acl_id"`
	TenantID    string `json:"tenant_id"`
	Description string `json:"description"`
	Protocol    string `json:"protocol"`      // tcp, udp, icmp, or null (any)
	EtherType   string `json:"ethertype"`     // IPv4, IPv6
	SrcIPPrefix string `json:"src_ip_prefix"` // Source IP CIDR
	DstIPPrefix string `json:"dst_ip_prefix"` // Destination IP CIDR
	SrcPortMin  *int   `json:"src_port_min"`
	SrcPortMax  *int   `json:"src_port_max"`
	DstPortMin  *int   `json:"dst_port_min"`
	DstPortMax  *int   `json:"dst_port_max"`
	Policy      string `json:"policy"` // allow, deny
	OrderNum    int    `json:"order"`  // Rule order/priority
	CreateTime  string `json:"create_time,omitempty"`
	UpdateTime  string `json:"update_time,omitempty"`
}

// ListACLRulesOutput represents the response for ACL rules list API
type ListACLRulesOutput struct {
	ACLRules []ACLRule `json:"acl_rules"`
}

// GetACLRuleOutput represents the response for single ACL rule API
type GetACLRuleOutput struct {
	ACLRule ACLRule `json:"acl_rule"`
}

// CreateACLRuleInput represents the ACL rule creation parameters
type CreateACLRuleInput struct {
	ACLID       string `json:"acl_id"`
	Description string `json:"description,omitempty"`
	Protocol    string `json:"protocol,omitempty"` // tcp, udp, icmp, or omit for any
	EtherType   string `json:"ethertype"`          // IPv4, IPv6
	SrcIPPrefix string `json:"src_ip_prefix,omitempty"`
	DstIPPrefix string `json:"dst_ip_prefix,omitempty"`
	SrcPortMin  *int   `json:"src_port_min,omitempty"`
	SrcPortMax  *int   `json:"src_port_max,omitempty"`
	DstPortMin  *int   `json:"dst_port_min,omitempty"`
	DstPortMax  *int   `json:"dst_port_max,omitempty"`
	Policy      string `json:"policy"` // allow, deny
	OrderNum    int    `json:"order"`  // Rule order/priority
}

// CreateACLRuleRequest represents the request body for creating an ACL rule
type CreateACLRuleRequest struct {
	ACLRule CreateACLRuleInput `json:"acl_rule"`
}

// UpdateACLRuleInput represents the ACL rule update parameters
type UpdateACLRuleInput struct {
	Description string `json:"description,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	SrcIPPrefix string `json:"src_ip_prefix,omitempty"`
	DstIPPrefix string `json:"dst_ip_prefix,omitempty"`
	SrcPortMin  *int   `json:"src_port_min,omitempty"`
	SrcPortMax  *int   `json:"src_port_max,omitempty"`
	DstPortMin  *int   `json:"dst_port_min,omitempty"`
	DstPortMax  *int   `json:"dst_port_max,omitempty"`
	Policy      string `json:"policy,omitempty"`
	OrderNum    *int   `json:"order,omitempty"`
}

// UpdateACLRuleRequest represents the request body for updating an ACL rule
type UpdateACLRuleRequest struct {
	ACLRule UpdateACLRuleInput `json:"acl_rule"`
}

// ACLBinding represents an ACL binding to a network resource (subnet)
type ACLBinding struct {
	ID         string `json:"id"`
	ACLID      string `json:"acl_id"`
	SubnetID   string `json:"subnet_id"`
	TenantID   string `json:"tenant_id"`
	CreateTime string `json:"create_time,omitempty"`
}

// ListACLBindingsOutput represents the response for ACL bindings list API
type ListACLBindingsOutput struct {
	ACLBindings []ACLBinding `json:"acl_bindings"`
}

// GetACLBindingOutput represents the response for single ACL binding API
type GetACLBindingOutput struct {
	ACLBinding ACLBinding `json:"acl_binding"`
}

// CreateACLBindingInput represents the ACL binding creation parameters
type CreateACLBindingInput struct {
	ACLID    string `json:"acl_id"`
	SubnetID string `json:"subnet_id"`
}

// CreateACLBindingRequest represents the request body for creating an ACL binding
type CreateACLBindingRequest struct {
	ACLBinding CreateACLBindingInput `json:"acl_binding"`
}
