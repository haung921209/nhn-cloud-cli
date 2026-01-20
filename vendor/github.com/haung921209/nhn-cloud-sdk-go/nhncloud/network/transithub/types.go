package transithub

// =============================================================================
// Transit Hub Types
// =============================================================================

// TransitHub represents a transit hub for multi-VPC networking
type TransitHub struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description,omitempty"`
	TenantID            string `json:"tenant_id"`
	State               string `json:"state,omitempty"`
	Status              string `json:"status,omitempty"`
	DefaultRoutingTable string `json:"default_routingtable_id,omitempty"`
	CreatedAt           string `json:"created_at,omitempty"`
	UpdatedAt           string `json:"updated_at,omitempty"`
}

type ListTransitHubsOutput struct {
	TransitHubs []TransitHub `json:"transithubs"`
}

type GetTransitHubOutput struct {
	TransitHub TransitHub `json:"transithub"`
}

type CreateTransitHubInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CreateTransitHubRequest struct {
	TransitHub CreateTransitHubInput `json:"transithub"`
}

type UpdateTransitHubInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateTransitHubRequest struct {
	TransitHub UpdateTransitHubInput `json:"transithub"`
}

// =============================================================================
// Attachment Types
// =============================================================================

// Attachment represents a transit hub attachment (VPC connection)
type Attachment struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	TenantID     string `json:"tenant_id"`
	TransitHubID string `json:"transithub_id"`
	ResourceType string `json:"resource_type"` // "VPC", "VPN", etc.
	ResourceID   string `json:"resource_id"`
	ResourceName string `json:"resource_name,omitempty"`
	State        string `json:"state,omitempty"`
	Status       string `json:"status,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type ListAttachmentsOutput struct {
	Attachments []Attachment `json:"transithub_attachments"`
}

type GetAttachmentOutput struct {
	Attachment Attachment `json:"transithub_attachment"`
}

type CreateAttachmentInput struct {
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	TransitHubID string `json:"transithub_id"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
}

type CreateAttachmentRequest struct {
	Attachment CreateAttachmentInput `json:"transithub_attachment"`
}

type UpdateAttachmentInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateAttachmentRequest struct {
	Attachment UpdateAttachmentInput `json:"transithub_attachment"`
}

// =============================================================================
// Routing Table Types
// =============================================================================

// RoutingTable represents a transit hub routing table
type RoutingTable struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	TenantID     string `json:"tenant_id"`
	TransitHubID string `json:"transithub_id"`
	DefaultTable bool   `json:"default_table,omitempty"`
	State        string `json:"state,omitempty"`
	Status       string `json:"status,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type ListRoutingTablesOutput struct {
	RoutingTables []RoutingTable `json:"transithub_routing_tables"`
}

type GetRoutingTableOutput struct {
	RoutingTable RoutingTable `json:"transithub_routing_table"`
}

type CreateRoutingTableInput struct {
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	TransitHubID string `json:"transithub_id"`
}

type CreateRoutingTableRequest struct {
	RoutingTable CreateRoutingTableInput `json:"transithub_routing_table"`
}

type UpdateRoutingTableInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateRoutingTableRequest struct {
	RoutingTable UpdateRoutingTableInput `json:"transithub_routing_table"`
}

// =============================================================================
// Routing Association Types
// =============================================================================

// RoutingAssociation represents an association between attachment and routing table
type RoutingAssociation struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenant_id"`
	RoutingTableID string `json:"routing_table_id"`
	AttachmentID   string `json:"attachment_id"`
	State          string `json:"state,omitempty"`
	Status         string `json:"status,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type ListRoutingAssociationsOutput struct {
	Associations []RoutingAssociation `json:"transithub_routing_associations"`
}

type GetRoutingAssociationOutput struct {
	Association RoutingAssociation `json:"transithub_routing_association"`
}

type CreateRoutingAssociationInput struct {
	RoutingTableID string `json:"routing_table_id"`
	AttachmentID   string `json:"attachment_id"`
}

type CreateRoutingAssociationRequest struct {
	Association CreateRoutingAssociationInput `json:"transithub_routing_association"`
}

// =============================================================================
// Routing Propagation Types
// =============================================================================

// RoutingPropagation represents route propagation from attachment to routing table
type RoutingPropagation struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenant_id"`
	RoutingTableID string `json:"routing_table_id"`
	AttachmentID   string `json:"attachment_id"`
	State          string `json:"state,omitempty"`
	Status         string `json:"status,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type ListRoutingPropagationsOutput struct {
	Propagations []RoutingPropagation `json:"transithub_routing_propagations"`
}

type GetRoutingPropagationOutput struct {
	Propagation RoutingPropagation `json:"transithub_routing_propagation"`
}

type CreateRoutingPropagationInput struct {
	RoutingTableID string `json:"routing_table_id"`
	AttachmentID   string `json:"attachment_id"`
}

type CreateRoutingPropagationRequest struct {
	Propagation CreateRoutingPropagationInput `json:"transithub_routing_propagation"`
}

// =============================================================================
// Routing Rule Types
// =============================================================================

// RoutingRule represents a static routing rule in a routing table
type RoutingRule struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenant_id"`
	RoutingTableID string `json:"routing_table_id"`
	Destination    string `json:"destination"`
	TargetType     string `json:"target_type"` // "ATTACHMENT", "BLACKHOLE"
	TargetID       string `json:"target_id,omitempty"`
	State          string `json:"state,omitempty"`
	Status         string `json:"status,omitempty"`
	Propagated     bool   `json:"propagated,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type ListRoutingRulesOutput struct {
	Rules []RoutingRule `json:"transithub_routing_rules"`
}

type GetRoutingRuleOutput struct {
	Rule RoutingRule `json:"transithub_routing_rule"`
}

type CreateRoutingRuleInput struct {
	RoutingTableID string `json:"routing_table_id"`
	Destination    string `json:"destination"`
	TargetType     string `json:"target_type"`
	TargetID       string `json:"target_id,omitempty"`
}

type CreateRoutingRuleRequest struct {
	Rule CreateRoutingRuleInput `json:"transithub_routing_rule"`
}

type UpdateRoutingRuleInput struct {
	Destination string `json:"destination,omitempty"`
	TargetType  string `json:"target_type,omitempty"`
	TargetID    string `json:"target_id,omitempty"`
}

type UpdateRoutingRuleRequest struct {
	Rule UpdateRoutingRuleInput `json:"transithub_routing_rule"`
}

// =============================================================================
// Multicast Domain Types
// =============================================================================

// MulticastDomain represents a multicast domain for transit hub
type MulticastDomain struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	TenantID     string `json:"tenant_id"`
	TransitHubID string `json:"transithub_id"`
	State        string `json:"state,omitempty"`
	Status       string `json:"status,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type ListMulticastDomainsOutput struct {
	Domains []MulticastDomain `json:"transithub_multicast_domains"`
}

type GetMulticastDomainOutput struct {
	Domain MulticastDomain `json:"transithub_multicast_domain"`
}

type CreateMulticastDomainInput struct {
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	TransitHubID string `json:"transithub_id"`
}

type CreateMulticastDomainRequest struct {
	Domain CreateMulticastDomainInput `json:"transithub_multicast_domain"`
}

type UpdateMulticastDomainInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateMulticastDomainRequest struct {
	Domain UpdateMulticastDomainInput `json:"transithub_multicast_domain"`
}
