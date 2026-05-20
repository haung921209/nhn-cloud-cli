// Package cloudtrail provides CloudTrail service types and client
package cloudtrail

// Header represents the response header
type Header struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// Member represents the event originator filter
type Member struct {
	MemberType   string `json:"memberType,omitempty"`   // TOAST or IAM
	UserCode     string `json:"userCode,omitempty"`     // IAM user identifier
	EmailAddress string `json:"emailAddress,omitempty"` // NHN Cloud user email
	IDNo         string `json:"idNo,omitempty"`         // User UUID
}

// PageInput represents the pagination sub-object in the request
type PageInput struct {
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
	SortBy string `json:"sortBy,omitempty"`
}

// SearchEventsInput represents a request to search events
type SearchEventsInput struct {
	StartDate string     `json:"startDate"`
	EndDate   string     `json:"endDate"`
	EventID   string     `json:"eventId"`
	IDNo      string     `json:"idNo,omitempty"`
	Member    *Member    `json:"member,omitempty"`
	Page      *PageInput `json:"page,omitempty"`
}

// SearchEventsOutput represents the response from event search
type SearchEventsOutput struct {
	Header Header             `json:"header"`
	Page   SearchEventsResult `json:"page,omitempty"`
}

// SearchEventsResult represents the search result data
type SearchEventsResult struct {
	TotalCount       int     `json:"totalElements"`
	TotalPages       int     `json:"totalPages"`
	NumberOfElements int     `json:"numberOfElements"`
	Events           []Event `json:"content"`
	First            bool    `json:"first"`
	Last             bool    `json:"last"`
	Empty            bool    `json:"empty"`
}

// TargetMember represents a user affected by an event
type TargetMember struct {
	IDNo         string `json:"idNo,omitempty"`
	Name         string `json:"name,omitempty"`
	UserCode     string `json:"userCode,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
}

// EventTarget represents the target of an event
type EventTarget struct {
	TargetMembers []TargetMember `json:"targetMembers,omitempty"`
}

// Event represents a CloudTrail event
type Event struct {
	EventTime       string      `json:"eventTime"`
	UserIDNo        string      `json:"userIdNo"`
	UserIP          string      `json:"userIp"`
	UserAgent       string      `json:"userAgent"`
	UserName        string      `json:"userName"`
	UserID          string      `json:"userId"`
	EventSourceType string      `json:"eventSourceType"`
	ProductID       string      `json:"productId"`
	Region          string      `json:"region"`
	OrgID           string      `json:"orgId"`
	ProjectID       string      `json:"projectId"`
	ProjectName     string      `json:"projectName"`
	AppKey          string      `json:"appKey"`
	TenantID        string      `json:"tenantId"`
	EventID         string      `json:"eventId"`
	EventLogUUID    string      `json:"eventLogUuid"`
	Request         string      `json:"request,omitempty"`
	Response        string      `json:"response,omitempty"`
	EventTarget     EventTarget `json:"eventTarget,omitempty"`
}
