// Package cloudtrail provides CloudTrail service types and client
package cloudtrail

import "time"

// Header represents the response header
type Header struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// SearchEventsInput represents a request to search events
type SearchEventsInput struct {
	From                time.Time `json:"from"`
	To                  time.Time `json:"to"`
	EventSourceTypeList []string  `json:"eventSourceTypeList,omitempty"` // CONSOLE, API
	MemberTypeList      []string  `json:"memberTypeList,omitempty"`      // TOAST, IAM
	MemberIDList        []string  `json:"memberIdList,omitempty"`
	EventIDList         []string  `json:"eventIdList,omitempty"`
	Page                int       `json:"page,omitempty"`
	Size                int       `json:"size,omitempty"`
}

// SearchEventsOutput represents the response from event search
type SearchEventsOutput struct {
	Header Header             `json:"header"`
	Body   SearchEventsResult `json:"body,omitempty"`
}

// SearchEventsResult represents the search result data
type SearchEventsResult struct {
	TotalCount int     `json:"totalCount"`
	Events     []Event `json:"events"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
}

// Event represents a CloudTrail event
type Event struct {
	EventTime       time.Time  `json:"eventTime"`
	EventSourceType string     `json:"eventSourceType"` // CONSOLE, API
	EventType       string     `json:"eventType"`       // API, SIGNIN, SIGNOUT, etc
	MemberType      string     `json:"memberType"`      // TOAST, IAM
	MemberID        string     `json:"memberId"`
	EventID         string     `json:"eventId"`
	SourceIP        string     `json:"sourceIp"`
	UserAgent       string     `json:"userAgent"`
	OrgID           string     `json:"orgId"`
	ProjectID       string     `json:"projectId"`
	ProductID       string     `json:"productId"`
	Region          string     `json:"region"`
	Resources       []Resource `json:"resources,omitempty"`
	RequestID       string     `json:"requestId,omitempty"`
	Request         string     `json:"request,omitempty"`
	Response        string     `json:"response,omitempty"`
}

// Resource represents a resource affected by an event
type Resource struct {
	ResourceType string `json:"resourceType"`
	ResourceID   string `json:"resourceId"`
	ResourceName string `json:"resourceName,omitempty"`
}
