// Package resourcewatcher provides Resource Watcher (Governance) service client
package resourcewatcher

// Header represents the standard NHN Cloud API response header
type Header struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// EventAlarm represents an event alarm configuration
type EventAlarm struct {
	AlarmID          string        `json:"alarmId"`
	AlarmName        string        `json:"alarmName"`
	AlarmDescription string        `json:"alarmDescription,omitempty"`
	AlarmStatusCode  string        `json:"alarmStatusCode"` // STABLE, DISABLED, CLOSED
	EventRuleID      string        `json:"eventRuleId,omitempty"`
	ResourceGroupID  string        `json:"resourceGroupId,omitempty"`
	ResourceTagID    string        `json:"resourceTagId,omitempty"`
	Targets          []AlarmTarget `json:"targets,omitempty"`
	CreatedDateTime  string        `json:"createdDateTime,omitempty"`
	UpdatedDateTime  string        `json:"updatedDateTime,omitempty"`
}

// AlarmTarget represents an alarm notification target
type AlarmTarget struct {
	TargetType   string `json:"alarmTargetTypeCode"` // UUID, ROLE, ALARM_KEY, WEBHOOK
	TargetID     string `json:"alarmTarget"`
	EmailAlarm   bool   `json:"emailAlarm,omitempty"`
	SMSAlarm     bool   `json:"smsAlarm,omitempty"`
	WebhookURL   string `json:"webhookUrl,omitempty"`
	WebhookSecret string `json:"webhookSecret,omitempty"`
}

// AlarmConfig holds the alarm metadata sub-object for create/update requests
type AlarmConfig struct {
	AlarmName   string `json:"alarmName"`
	Description string `json:"description,omitempty"`
}

// AlarmEvent represents an event to monitor in an alarm
type AlarmEvent struct {
	ProductID string `json:"productId,omitempty"`
	EventID   string `json:"eventId,omitempty"`
}

// AlarmTargetScope defines the resource scope for an alarm
type AlarmTargetScope struct {
	ResourceGroupIDs []string `json:"resourceGroupIds,omitempty"`
	ResourceTagIDs   []string `json:"resourceTagIds,omitempty"`
}

// CreateEventAlarmInput represents the request body for creating an event alarm
type CreateEventAlarmInput struct {
	Alarm        AlarmConfig      `json:"alarm"`
	AlarmTargets []AlarmTarget    `json:"alarmTargets"`
	Events       []AlarmEvent     `json:"events,omitempty"`
	Target       *AlarmTargetScope `json:"target,omitempty"`
}

// CreateEventAlarmOutput represents the response for creating an event alarm
type CreateEventAlarmOutput struct {
	Header  Header `json:"header"`
	AlarmID string `json:"alarmId"`
}

// GetEventAlarmOutput represents the response for single event alarm API
type GetEventAlarmOutput struct {
	Header Header     `json:"header"`
	Alarm  EventAlarm `json:"alarm"`
}

// SearchEventAlarmsInput represents the request body for searching event alarms
type SearchEventAlarmsInput struct {
	AlarmName       string `json:"alarmName,omitempty"`
	AlarmStatusCode string `json:"alarmStatusCode,omitempty"`
	Page            int    `json:"page,omitempty"`
	Size            int    `json:"size,omitempty"`
}

// SearchEventAlarmsOutput represents the response for event alarms search API
type SearchEventAlarmsOutput struct {
	Header     Header       `json:"header"`
	Alarms     []EventAlarm `json:"alarms"`
	TotalCount int          `json:"totalCount"`
}

// UpdateEventAlarmInput represents the request body for updating an event alarm
type UpdateEventAlarmInput struct {
	Alarm        AlarmConfig      `json:"alarm"`
	AlarmTargets []AlarmTarget    `json:"alarmTargets"`
	Events       []AlarmEvent     `json:"events,omitempty"`
	Target       *AlarmTargetScope `json:"target,omitempty"`
}

// DeleteEventAlarmsInput represents the request body for deleting multiple alarms
type DeleteEventAlarmsInput struct {
	AlarmIDs []string `json:"alarmIds"`
}

// SimpleOutput represents a simple success/failure response
type SimpleOutput struct {
	Header Header `json:"header"`
}

// AlarmHistory represents an alarm history record
type AlarmHistory struct {
	AlarmHistoryID   string            `json:"alarmHistoryId"`
	AlarmID          string            `json:"alarmId"`
	EventID          string            `json:"eventId"`
	ResourceID       string            `json:"resourceId"`
	ResourceName     string            `json:"resourceName"`
	ProductID        string            `json:"productId"`
	EventName        string            `json:"eventName"`
	AlarmSendResults []AlarmSendResult `json:"alarmSendResults,omitempty"`
	CreatedDateTime  string            `json:"createdDateTime"`
}

// AlarmSendResult represents the result of sending an alarm notification
type AlarmSendResult struct {
	TargetType   string `json:"targetType"`
	TargetID     string `json:"targetId"`
	SendStatus   string `json:"sendStatus"`
	SentDateTime string `json:"sentDateTime,omitempty"`
}

// GetAlarmHistoryOutput represents the response for alarm history API
type GetAlarmHistoryOutput struct {
	Header  Header       `json:"header"`
	History AlarmHistory `json:"alarmHistory"`
}

// SearchAlarmHistoryInput represents the request body for searching alarm history
type SearchAlarmHistoryInput struct {
	StartDateTime string `json:"startDateTime,omitempty"`
	EndDateTime   string `json:"endDateTime,omitempty"`
	Page          int    `json:"page,omitempty"`
	Size          int    `json:"size,omitempty"`
}

// SearchAlarmHistoryOutput represents the response for alarm history search API
type SearchAlarmHistoryOutput struct {
	Header     Header         `json:"header"`
	Histories  []AlarmHistory `json:"alarmHistories"`
	TotalCount int            `json:"totalCount"`
}

// Event represents a resource event
type Event struct {
	EventID     string `json:"eventId"`
	ProductID   string `json:"productId"`
	EventName   string `json:"eventName"`
	EventType   string `json:"eventType"`
	Description string `json:"description,omitempty"`
}

// ListEventsOutput represents the response for events list API
type ListEventsOutput struct {
	Header Header  `json:"header"`
	Events []Event `json:"events"`
}

// GetEventOutput represents the response for single event API
type GetEventOutput struct {
	Header Header `json:"header"`
	Event  Event  `json:"event"`
}

// ResourceGroup represents a resource group
type ResourceGroup struct {
	ResourceGroupID   string `json:"resourceGroupId"`
	ResourceGroupName string `json:"resourceGroupName"`
	Description       string `json:"description,omitempty"`
	CreatedDateTime   string `json:"createdDateTime,omitempty"`
}

// ListResourceGroupsOutput represents the response for resource groups list API
type ListResourceGroupsOutput struct {
	Header         Header          `json:"header"`
	ResourceGroups []ResourceGroup `json:"resourceGroups"`
}

// ResourceTag represents a resource tag
type ResourceTag struct {
	ResourceTagID   string `json:"resourceTagId"`
	ResourceTagName string `json:"tagName"`
	TagKey          string `json:"tagKey"`
	TagValue        string `json:"tagValue,omitempty"`
	CreatedDateTime string `json:"createdDateTime,omitempty"`
}

// ListResourceTagsOutput represents the response for resource tags list API
type ListResourceTagsOutput struct {
	Header       Header        `json:"header"`
	ResourceTags []ResourceTag `json:"resourceTags"`
}
