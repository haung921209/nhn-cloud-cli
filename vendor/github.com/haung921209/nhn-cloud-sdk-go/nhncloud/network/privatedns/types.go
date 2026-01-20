package privatedns

import "time"

// Zone represents a private DNS zone
type Zone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	VPCID       string    `json:"vpc_id"`
	State       string    `json:"state"`
	RecordCount int       `json:"record_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// RRSet represents a DNS record set
type RRSet struct {
	ID        string    `json:"id"`
	ZoneID    string    `json:"zone_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // A, AAAA, CNAME, MX, TXT, etc.
	TTL       int       `json:"ttl"`
	Records   []string  `json:"records"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// ListZonesOutput represents the response from listing zones
type ListZonesOutput struct {
	Zones []Zone `json:"privatedns_zones"`
}

// GetZoneOutput represents the response containing a single zone
type GetZoneOutput struct {
	Zone *Zone `json:"privatedns_zone"`
}

// CreateZoneInput contains the zone creation data
type CreateZoneInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	VPCID       string `json:"vpc_id"`
}

// CreateZoneRequest wraps the create input
type CreateZoneRequest struct {
	Zone CreateZoneInput `json:"privatedns_zone"`
}

// UpdateZoneInput contains the zone update data
type UpdateZoneInput struct {
	Description string `json:"description,omitempty"`
}

// UpdateZoneRequest wraps the update input
type UpdateZoneRequest struct {
	Zone UpdateZoneInput `json:"privatedns_zone"`
}

// ListRRSetsOutput represents the response from listing record sets
type ListRRSetsOutput struct {
	RRSets []RRSet `json:"rrsets"`
}

// GetRRSetOutput represents the response containing a single record set
type GetRRSetOutput struct {
	RRSet *RRSet `json:"rrset"`
}

// CreateRRSetInput contains the record set creation data
type CreateRRSetInput struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	TTL     int      `json:"ttl"`
	Records []string `json:"records"`
}

// CreateRRSetRequest wraps the create input
type CreateRRSetRequest struct {
	RRSet CreateRRSetInput `json:"rrset"`
}

// UpdateRRSetInput contains the record set update data
type UpdateRRSetInput struct {
	TTL     int      `json:"ttl,omitempty"`
	Records []string `json:"records,omitempty"`
}

// UpdateRRSetRequest wraps the update input
type UpdateRRSetRequest struct {
	RRSet UpdateRRSetInput `json:"rrset"`
}
