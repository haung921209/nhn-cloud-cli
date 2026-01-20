package securitygroup

type SecurityGroup struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	TenantID    string         `json:"tenant_id"`
	Description string         `json:"description,omitempty"`
	Rules       []SecurityRule `json:"security_group_rules,omitempty"`
}

type SecurityRule struct {
	ID              string  `json:"id"`
	TenantID        string  `json:"tenant_id"`
	SecurityGroupID string  `json:"security_group_id"`
	Direction       string  `json:"direction"`
	EtherType       string  `json:"ethertype"`
	Protocol        *string `json:"protocol"`
	PortRangeMin    *int    `json:"port_range_min"`
	PortRangeMax    *int    `json:"port_range_max"`
	RemoteIPPrefix  string  `json:"remote_ip_prefix,omitempty"`
	RemoteGroupID   string  `json:"remote_group_id,omitempty"`
	Description     string  `json:"description,omitempty"`
}

type ListSecurityGroupsOutput struct {
	SecurityGroups []SecurityGroup `json:"security_groups"`
}

type GetSecurityGroupOutput struct {
	SecurityGroup SecurityGroup `json:"security_group"`
}

type CreateSecurityGroupInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CreateSecurityGroupOutput struct {
	SecurityGroup SecurityGroup `json:"security_group"`
}

type UpdateSecurityGroupInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type CreateRuleInput struct {
	SecurityGroupID string `json:"security_group_id"`
	Direction       string `json:"direction"`
	EtherType       string `json:"ethertype,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	PortRangeMin    *int   `json:"port_range_min,omitempty"`
	PortRangeMax    *int   `json:"port_range_max,omitempty"`
	RemoteIPPrefix  string `json:"remote_ip_prefix,omitempty"`
	RemoteGroupID   string `json:"remote_group_id,omitempty"`
	Description     string `json:"description,omitempty"`
}

type CreateRuleOutput struct {
	SecurityGroupRule SecurityRule `json:"security_group_rule"`
}
