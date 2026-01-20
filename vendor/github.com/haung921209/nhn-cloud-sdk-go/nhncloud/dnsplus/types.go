// Package dnsplus provides DNS Plus service types and client
package dnsplus

import "time"

// Header represents the response header
type Header struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// APIResponse represents the common API response wrapper
type APIResponse struct {
	Header Header `json:"header"`
}

// ================================
// DNS Zone Types
// ================================

// Zone represents a DNS zone
type Zone struct {
	ZoneID         string    `json:"zoneId"`
	ZoneName       string    `json:"zoneName"`
	ZoneStatus     string    `json:"zoneStatus"` // USE, STOP
	Description    string    `json:"description,omitempty"`
	RecordSetCount int       `json:"recordsetCount"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	NameServers    []string  `json:"nameServers,omitempty"`
}

// ListZonesOutput represents the response from listing zones
type ListZonesOutput struct {
	Header     Header `json:"header"`
	TotalCount int    `json:"totalCount"`
	ZoneList   []Zone `json:"zoneList"`
}

// ZoneOutput represents a single zone response
type ZoneOutput struct {
	Header Header `json:"header"`
	Zone   *Zone  `json:"zone,omitempty"`
}

// CreateZoneInput represents a request to create a zone
type CreateZoneInput struct {
	ZoneName    string `json:"zoneName"`
	Description string `json:"description,omitempty"`
}

// UpdateZoneInput represents a request to update a zone
type UpdateZoneInput struct {
	Description string `json:"description,omitempty"`
	ZoneStatus  string `json:"zoneStatus,omitempty"` // USE, STOP
}

// ================================
// Record Set Types
// ================================

// RecordSet represents a DNS record set
type RecordSet struct {
	RecordSetID   string    `json:"recordsetId"`
	RecordSetName string    `json:"recordsetName"`
	RecordSetType string    `json:"recordsetType"` // A, AAAA, CAA, CNAME, MX, NAPTR, PTR, TXT, SRV, NS, SOA
	TTL           int       `json:"ttl"`
	RecordList    []Record  `json:"recordList"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// Record represents a single DNS record
type Record struct {
	RecordDisabled bool   `json:"recordDisabled"`
	RecordContent  string `json:"recordContent"`
}

// ListRecordSetsOutput represents the response from listing record sets
type ListRecordSetsOutput struct {
	Header        Header      `json:"header"`
	TotalCount    int         `json:"totalCount"`
	RecordSetList []RecordSet `json:"recordsetList"`
}

// RecordSetOutput represents a single record set response
type RecordSetOutput struct {
	Header    Header     `json:"header"`
	RecordSet *RecordSet `json:"recordset,omitempty"`
}

// CreateRecordSetInput represents a request to create a record set
type CreateRecordSetInput struct {
	RecordSetName string   `json:"recordsetName"`
	RecordSetType string   `json:"recordsetType"`
	TTL           int      `json:"ttl"`
	RecordList    []Record `json:"recordList"`
}

// UpdateRecordSetInput represents a request to update a record set
type UpdateRecordSetInput struct {
	TTL        int      `json:"ttl,omitempty"`
	RecordList []Record `json:"recordList,omitempty"`
}

// ================================
// GSLB Types
// ================================

// GSLB represents a Global Server Load Balancing configuration
type GSLB struct {
	GslbID        string    `json:"gslbId"`
	GslbName      string    `json:"gslbName"`
	GslbDomain    string    `json:"gslbDomain"`
	GslbStatus    string    `json:"gslbStatus"` // USE, STOP
	Description   string    `json:"description,omitempty"`
	RoutingType   string    `json:"routingType"` // FAILOVER, RANDOM, GEOLOCATION
	TTL           int       `json:"ttl"`
	PoolCount     int       `json:"poolCount"`
	HealthCheckID string    `json:"healthCheckId,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ListGSLBsOutput represents the response from listing GSLBs
type ListGSLBsOutput struct {
	Header     Header `json:"header"`
	TotalCount int    `json:"totalCount"`
	GslbList   []GSLB `json:"gslbList"`
}

// GSLBOutput represents a single GSLB response
type GSLBOutput struct {
	Header Header `json:"header"`
	Gslb   *GSLB  `json:"gslb,omitempty"`
}

// CreateGSLBInput represents a request to create a GSLB
type CreateGSLBInput struct {
	GslbName      string `json:"gslbName"`
	Description   string `json:"description,omitempty"`
	RoutingType   string `json:"routingType"`
	TTL           int    `json:"ttl"`
	HealthCheckID string `json:"healthCheckId,omitempty"`
}

// UpdateGSLBInput represents a request to update a GSLB
type UpdateGSLBInput struct {
	GslbName      string `json:"gslbName,omitempty"`
	Description   string `json:"description,omitempty"`
	GslbStatus    string `json:"gslbStatus,omitempty"`
	TTL           int    `json:"ttl,omitempty"`
	HealthCheckID string `json:"healthCheckId,omitempty"`
}

// ================================
// GSLB Pool Types
// ================================

// Pool represents a GSLB pool
type Pool struct {
	PoolID        string    `json:"poolId"`
	PoolName      string    `json:"poolName"`
	PoolStatus    string    `json:"poolStatus"` // USE, STOP
	Description   string    `json:"description,omitempty"`
	Priority      int       `json:"priority"`
	Weight        int       `json:"weight"`
	Region        string    `json:"region,omitempty"`
	EndpointCount int       `json:"endpointCount"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ListPoolsOutput represents the response from listing pools
type ListPoolsOutput struct {
	Header     Header `json:"header"`
	TotalCount int    `json:"totalCount"`
	PoolList   []Pool `json:"poolList"`
}

// PoolOutput represents a single pool response
type PoolOutput struct {
	Header Header `json:"header"`
	Pool   *Pool  `json:"pool,omitempty"`
}

// CreatePoolInput represents a request to create a pool
type CreatePoolInput struct {
	PoolName    string `json:"poolName"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority"`
	Weight      int    `json:"weight,omitempty"`
	Region      string `json:"region,omitempty"`
}

// UpdatePoolInput represents a request to update a pool
type UpdatePoolInput struct {
	PoolName    string `json:"poolName,omitempty"`
	Description string `json:"description,omitempty"`
	PoolStatus  string `json:"poolStatus,omitempty"`
	Priority    int    `json:"priority,omitempty"`
	Weight      int    `json:"weight,omitempty"`
	Region      string `json:"region,omitempty"`
}

// ================================
// GSLB Endpoint Types
// ================================

// Endpoint represents a GSLB pool endpoint
type Endpoint struct {
	EndpointID      string    `json:"endpointId"`
	EndpointAddress string    `json:"endpointAddress"`
	EndpointStatus  string    `json:"endpointStatus"` // USE, STOP
	Weight          int       `json:"weight"`
	Description     string    `json:"description,omitempty"`
	HealthStatus    string    `json:"healthStatus,omitempty"` // HEALTHY, UNHEALTHY
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// ListEndpointsOutput represents the response from listing endpoints
type ListEndpointsOutput struct {
	Header       Header     `json:"header"`
	TotalCount   int        `json:"totalCount"`
	EndpointList []Endpoint `json:"endpointList"`
}

// EndpointOutput represents a single endpoint response
type EndpointOutput struct {
	Header   Header    `json:"header"`
	Endpoint *Endpoint `json:"endpoint,omitempty"`
}

// CreateEndpointInput represents a request to create an endpoint
type CreateEndpointInput struct {
	EndpointAddress string `json:"endpointAddress"`
	Weight          int    `json:"weight,omitempty"`
	Description     string `json:"description,omitempty"`
}

// UpdateEndpointInput represents a request to update an endpoint
type UpdateEndpointInput struct {
	EndpointAddress string `json:"endpointAddress,omitempty"`
	EndpointStatus  string `json:"endpointStatus,omitempty"`
	Weight          int    `json:"weight,omitempty"`
	Description     string `json:"description,omitempty"`
}

// ================================
// Health Check Types
// ================================

// HealthCheck represents a health check configuration
type HealthCheck struct {
	HealthCheckID   string    `json:"healthCheckId"`
	HealthCheckName string    `json:"healthCheckName"`
	Description     string    `json:"description,omitempty"`
	Protocol        string    `json:"protocol"` // HTTP, HTTPS, TCP
	Port            int       `json:"port"`
	Path            string    `json:"path,omitempty"`          // For HTTP/HTTPS
	Host            string    `json:"host,omitempty"`          // Host header for HTTP/HTTPS
	Interval        int       `json:"interval"`                // Seconds between checks
	Timeout         int       `json:"timeout"`                 // Seconds to wait for response
	Retries         int       `json:"retries"`                 // Number of retries before marking unhealthy
	ExpectedCodes   string    `json:"expectedCodes,omitempty"` // e.g., "200,201,202"
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// ListHealthChecksOutput represents the response from listing health checks
type ListHealthChecksOutput struct {
	Header          Header        `json:"header"`
	TotalCount      int           `json:"totalCount"`
	HealthCheckList []HealthCheck `json:"healthCheckList"`
}

// HealthCheckOutput represents a single health check response
type HealthCheckOutput struct {
	Header      Header       `json:"header"`
	HealthCheck *HealthCheck `json:"healthCheck,omitempty"`
}

// CreateHealthCheckInput represents a request to create a health check
type CreateHealthCheckInput struct {
	HealthCheckName string `json:"healthCheckName"`
	Description     string `json:"description,omitempty"`
	Protocol        string `json:"protocol"`
	Port            int    `json:"port"`
	Path            string `json:"path,omitempty"`
	Host            string `json:"host,omitempty"`
	Interval        int    `json:"interval"`
	Timeout         int    `json:"timeout"`
	Retries         int    `json:"retries"`
	ExpectedCodes   string `json:"expectedCodes,omitempty"`
}

// UpdateHealthCheckInput represents a request to update a health check
type UpdateHealthCheckInput struct {
	HealthCheckName string `json:"healthCheckName,omitempty"`
	Description     string `json:"description,omitempty"`
	Port            int    `json:"port,omitempty"`
	Path            string `json:"path,omitempty"`
	Host            string `json:"host,omitempty"`
	Interval        int    `json:"interval,omitempty"`
	Timeout         int    `json:"timeout,omitempty"`
	Retries         int    `json:"retries,omitempty"`
	ExpectedCodes   string `json:"expectedCodes,omitempty"`
}
