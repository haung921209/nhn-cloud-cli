// Package mirroring provides Traffic Mirroring service types and client
package mirroring

import "time"

// ================================
// Session Types
// ================================

// Session represents a mirroring session
type Session struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	TargetPortID    string    `json:"target_port_id"`
	SourcePortID    string    `json:"source_port_id"`
	FilterGroups    []string  `json:"filter_groups,omitempty"`
	Direction       string    `json:"direction"` // in, out, both
	TenantID        string    `json:"tenant_id"`
	ProjectID       string    `json:"project_id"`
	TargetNetworkID string    `json:"target_network_id,omitempty"`
	SourceNetworkID string    `json:"source_network_id,omitempty"`
	VNI             int       `json:"vni,omitempty"`
	Status          string    `json:"status"` // ACTIVE, BUILD, ERROR
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

// ListSessionsOutput represents the response from listing sessions
type ListSessionsOutput struct {
	Sessions []Session `json:"sessions"`
}

// SessionOutput represents the response containing a single session
type SessionOutput struct {
	Session *Session `json:"session"`
}

// CreateSessionInput represents session creation data
type CreateSessionInput struct {
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	TargetPortID  string   `json:"target_port_id"`
	SourcePortID  string   `json:"source_port_id"`
	FilterGroups  []string `json:"filter_groups,omitempty"`
	Direction     string   `json:"direction"`
	TenantID      string   `json:"tenant_id"`
	ProjectID     string   `json:"project_id"`
}

// UpdateSessionInput represents session update data
type UpdateSessionInput struct {
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	FilterGroups []string `json:"filter_groups,omitempty"`
	Direction    string   `json:"direction,omitempty"`
}

// ================================
// Filter Group Types
// ================================

// FilterGroup represents a mirroring filter group
type FilterGroup struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	TenantID    string    `json:"tenant_id"`
	ProjectID   string    `json:"project_id"`
	Filters     []string  `json:"filters,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// ListFilterGroupsOutput represents the response from listing filter groups
type ListFilterGroupsOutput struct {
	FilterGroups []FilterGroup `json:"filtergroups"`
}

// FilterGroupOutput represents the response containing a single filter group
type FilterGroupOutput struct {
	FilterGroup *FilterGroup `json:"filtergroup"`
}

// CreateFilterGroupInput represents filter group creation data
type CreateFilterGroupInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	TenantID    string `json:"tenant_id"`
	ProjectID   string `json:"project_id"`
}

// UpdateFilterGroupInput represents filter group update data
type UpdateFilterGroupInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ================================
// Filter Types
// ================================

// Filter represents a mirroring filter
type Filter struct {
	ID             string    `json:"id"`
	Description    string    `json:"description,omitempty"`
	FilterGroupID  string    `json:"filter_group_id"`
	SrcCIDR        string    `json:"src_cidr,omitempty"`
	DstCIDR        string    `json:"dst_cidr,omitempty"`
	SrcPortMin     int       `json:"src_port_range_min,omitempty"`
	SrcPortMax     int       `json:"src_port_range_max,omitempty"`
	DstPortMin     int       `json:"dst_port_range_min,omitempty"`
	DstPortMax     int       `json:"dst_port_range_max,omitempty"`
	Protocol       string    `json:"protocol,omitempty"`
	Action         string    `json:"action"` // accept, drop
	Priority       int       `json:"priority,omitempty"`
	TenantID       string    `json:"tenant_id"`
	ProjectID      string    `json:"project_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// ListFiltersOutput represents the response from listing filters
type ListFiltersOutput struct {
	Filters []Filter `json:"filters"`
}

// FilterOutput represents the response containing a single filter
type FilterOutput struct {
	Filter *Filter `json:"filter"`
}

// CreateFilterInput represents filter creation data
type CreateFilterInput struct {
	Description   string `json:"description,omitempty"`
	FilterGroupID string `json:"filter_group_id"`
	SrcCIDR       string `json:"src_cidr,omitempty"`
	DstCIDR       string `json:"dst_cidr,omitempty"`
	SrcPortMin    int    `json:"src_port_range_min,omitempty"`
	SrcPortMax    int    `json:"src_port_range_max,omitempty"`
	DstPortMin    int    `json:"dst_port_range_min,omitempty"`
	DstPortMax    int    `json:"dst_port_range_max,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	Action        string `json:"action"`
	Priority      int    `json:"priority,omitempty"`
	TenantID      string `json:"tenant_id"`
	ProjectID     string `json:"project_id"`
}

// UpdateFilterInput represents filter update data
type UpdateFilterInput struct {
	Description string `json:"description,omitempty"`
}
