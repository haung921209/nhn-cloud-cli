// Package mirroring provides Traffic Mirroring service types and client
package mirroring

import "time"

// ================================
// Session Types
// ================================

// Session represents a mirroring session
type Session struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	SourceType    string    `json:"source_type"` // PORT, LOADBALANCER_MEMBER
	SourceID      string    `json:"source_id"`
	TargetType    string    `json:"target_type"` // PORT
	TargetID      string    `json:"target_id"`
	Direction     string    `json:"direction"` // in, out, both
	FilterGroupID string    `json:"filter_group_id,omitempty"`
	AdminStateUp  bool      `json:"admin_state_up"`
	State         string    `json:"state"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

// ListSessionsOutput represents the response from listing sessions
type ListSessionsOutput struct {
	MirroringSessions []Session `json:"mirroring_sessions"`
}

// SessionOutput represents the response containing a single session
type SessionOutput struct {
	MirroringSession *Session `json:"mirroring_session"`
}

// CreateSessionInput represents session creation data
type CreateSessionInput struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	SourceType    string `json:"source_type"`
	SourceID      string `json:"source_id"`
	TargetType    string `json:"target_type"`
	TargetID      string `json:"target_id"`
	Direction     string `json:"direction"`
	FilterGroupID string `json:"filter_group_id,omitempty"`
	AdminStateUp  bool   `json:"admin_state_up,omitempty"`
}

// UpdateSessionInput represents session update data
type UpdateSessionInput struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	Direction     string `json:"direction,omitempty"`
	FilterGroupID string `json:"filter_group_id,omitempty"`
	AdminStateUp  *bool  `json:"admin_state_up,omitempty"`
}

// ================================
// Filter Group Types
// ================================

// FilterGroup represents a mirroring filter group
type FilterGroup struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	FilterIDs   []string  `json:"filter_ids,omitempty"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// ListFilterGroupsOutput represents the response from listing filter groups
type ListFilterGroupsOutput struct {
	MirroringFilterGroups []FilterGroup `json:"mirroring_filtergroups"`
}

// FilterGroupOutput represents the response containing a single filter group
type FilterGroupOutput struct {
	MirroringFilterGroup *FilterGroup `json:"mirroring_filtergroup"`
}

// CreateFilterGroupInput represents filter group creation data
type CreateFilterGroupInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
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
	ID            string    `json:"id"`
	FilterGroupID string    `json:"filter_group_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	Protocol      string    `json:"protocol,omitempty"`
	SourceCIDR    string    `json:"source_cidr,omitempty"`
	DestCIDR      string    `json:"destination_cidr,omitempty"`
	SourcePortMin int       `json:"source_port_min,omitempty"`
	SourcePortMax int       `json:"source_port_max,omitempty"`
	DestPortMin   int       `json:"destination_port_min,omitempty"`
	DestPortMax   int       `json:"destination_port_max,omitempty"`
	Action        string    `json:"action"` // accept, drop
	State         string    `json:"state"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

// ListFiltersOutput represents the response from listing filters
type ListFiltersOutput struct {
	MirroringFilters []Filter `json:"mirroring_filters"`
}

// FilterOutput represents the response containing a single filter
type FilterOutput struct {
	MirroringFilter *Filter `json:"mirroring_filter"`
}

// CreateFilterInput represents filter creation data
type CreateFilterInput struct {
	FilterGroupID string `json:"filter_group_id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	SourceCIDR    string `json:"source_cidr,omitempty"`
	DestCIDR      string `json:"destination_cidr,omitempty"`
	SourcePortMin int    `json:"source_port_min,omitempty"`
	SourcePortMax int    `json:"source_port_max,omitempty"`
	DestPortMin   int    `json:"destination_port_min,omitempty"`
	DestPortMax   int    `json:"destination_port_max,omitempty"`
	Action        string `json:"action"`
}

// UpdateFilterInput represents filter update data
type UpdateFilterInput struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	SourceCIDR    string `json:"source_cidr,omitempty"`
	DestCIDR      string `json:"destination_cidr,omitempty"`
	SourcePortMin int    `json:"source_port_min,omitempty"`
	SourcePortMax int    `json:"source_port_max,omitempty"`
	DestPortMin   int    `json:"destination_port_min,omitempty"`
	DestPortMax   int    `json:"destination_port_max,omitempty"`
	Action        string `json:"action,omitempty"`
}
