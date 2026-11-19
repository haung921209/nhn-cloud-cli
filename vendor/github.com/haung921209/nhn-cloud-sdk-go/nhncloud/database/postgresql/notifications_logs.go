package postgresql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// NotificationGroup represents a notification group
type NotificationGroup struct {
	NotificationGroupID     string `json:"notificationGroupId"`
	NotificationGroupName   string `json:"notificationGroupName"`
	NotificationGroupStatus string `json:"notificationGroupStatus,omitempty"`
	IsEnabled               bool   `json:"isEnabled"`
	NotifyEmail             bool   `json:"notifyEmail"`
	NotifySms               bool   `json:"notifySms"`
	CreatedAt               string `json:"createdYmdt,omitempty"`
	UpdatedAt               string `json:"updatedYmdt,omitempty"`
}

// ListNotificationGroupsResponse is the response for ListNotificationGroups
type ListNotificationGroupsResponse struct {
	PostgreSQLResponse
	NotificationGroups []NotificationGroup `json:"notificationGroups"`
}

// ListNotificationGroups retrieves all notification groups.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_69
func (c *Client) ListNotificationGroups(ctx context.Context) (*ListNotificationGroupsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "/v1.0/notification-groups", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListNotificationGroupsResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetNotificationGroupResponse is the response for GetNotificationGroup.
// The API returns the group's fields at the top level, next to header.
type GetNotificationGroupResponse struct {
	PostgreSQLResponse
	NotificationGroup
}

// GetNotificationGroup retrieves a specific notification group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_70
func (c *Client) GetNotificationGroup(ctx context.Context, groupID string) (*GetNotificationGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "notification group ID is required"}
	}

	path := fmt.Sprintf("/v1.0/notification-groups/%s", groupID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result GetNotificationGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateNotificationGroupRequest is the request for creating a notification group
type CreateNotificationGroupRequest struct {
	NotificationGroupName string   `json:"notificationGroupName"`
	IsEnabled             bool     `json:"isEnabled"`
	NotifyEmail           bool     `json:"notifyEmail"`
	NotifySms             bool     `json:"notifySms"`
	DBInstanceIDs         []string `json:"dbInstanceIds"` // required by the API; may be empty
	UserGroupIDs          []string `json:"userGroupIds"`  // required by the API; may be empty
}

// CreateNotificationGroupResponse is the response for CreateNotificationGroup
type CreateNotificationGroupResponse struct {
	PostgreSQLResponse
	NotificationGroupID string `json:"notificationGroupId"`
}

// CreateNotificationGroup creates a new notification group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_71
func (c *Client) CreateNotificationGroup(ctx context.Context, req *CreateNotificationGroupRequest) (*CreateNotificationGroupResponse, error) {
	if req.NotificationGroupName == "" {
		return nil, &core.ValidationError{Field: "NotificationGroupName", Message: "notification group name is required"}
	}
	r := *req
	if r.DBInstanceIDs == nil {
		r.DBInstanceIDs = []string{}
	}
	if r.UserGroupIDs == nil {
		r.UserGroupIDs = []string{}
	}

	body, err := json.Marshal(&r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "/v1.0/notification-groups", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result CreateNotificationGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateNotificationGroupRequest is the request for updating a notification group
type UpdateNotificationGroupRequest struct {
	NotificationGroupName *string `json:"notificationGroupName,omitempty"`
	IsEnabled             *bool   `json:"isEnabled,omitempty"`
	NotifyEmail           *bool   `json:"notifyEmail,omitempty"`
	NotifySms             *bool   `json:"notifySms,omitempty"`
}

// UpdateNotificationGroupResponse is the response for UpdateNotificationGroup
type UpdateNotificationGroupResponse struct {
	PostgreSQLResponse
}

// UpdateNotificationGroup updates a notification group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_72
func (c *Client) UpdateNotificationGroup(ctx context.Context, groupID string, req *UpdateNotificationGroupRequest) (*UpdateNotificationGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "notification group ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1.0/notification-groups/%s", groupID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var result UpdateNotificationGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteNotificationGroupResponse is the response for DeleteNotificationGroup
type DeleteNotificationGroupResponse struct {
	PostgreSQLResponse
}

// DeleteNotificationGroup deletes a notification group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_73
func (c *Client) DeleteNotificationGroup(ctx context.Context, groupID string) (*DeleteNotificationGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "notification group ID is required"}
	}

	path := fmt.Sprintf("/v1.0/notification-groups/%s", groupID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result DeleteNotificationGroupResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Metric represents a metric
type Metric struct {
	MetricName string `json:"metricName"`
	Unit       string `json:"unit,omitempty"`
}

// ListMetricsResponse is the response for ListMetrics
type ListMetricsResponse struct {
	PostgreSQLResponse
	Metrics []Metric `json:"metrics"`
}

// ListMetrics retrieves available metrics.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v3.0/#_75
func (c *Client) ListMetrics(ctx context.Context) (*ListMetricsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "/v1.0/metrics", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListMetricsResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// MetricStatistic represents metric statistics
type MetricStatistic struct {
	MetricName string                 `json:"metricName"`
	Unit       string                 `json:"unit,omitempty"`
	Values     []MetricStatisticValue `json:"values"`
}

// MetricStatisticValue represents a metric value at a point in time
type MetricStatisticValue struct {
	Timestamp int64   `json:"timestamp"` // Unix seconds
	Value     float64 `json:"value"`     // a number on the live API (the guide says string)
}

// UnmarshalJSON accepts value as a number or as a numeric string.
func (v *MetricStatisticValue) UnmarshalJSON(b []byte) error {
	var raw struct {
		Timestamp int64       `json:"timestamp"`
		Value     json.Number `json:"value"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	v.Timestamp = raw.Timestamp
	if raw.Value == "" {
		v.Value = 0
		return nil
	}
	f, err := raw.Value.Float64()
	if err != nil {
		return fmt.Errorf("metric value %q: %w", raw.Value, err)
	}
	v.Value = f
	return nil
}

// GetMetricStatisticsResponse is the response for GetMetricStatistics
type GetMetricStatisticsResponse struct {
	PostgreSQLResponse
	MetricStatistics []MetricStatistic `json:"metricStatistics"`
}

// GetMetricStatistics retrieves metric statistics for an instance.
// metricNames is required (see ListMetrics); from and to are
// YYYY-MM-DDThh:mm:ss.SSSTZD; interval is in minutes, 0 lets the API choose.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/
func (c *Client) GetMetricStatistics(ctx context.Context, instanceID string, metricNames []string, from, to string, interval int) (*GetMetricStatisticsResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if len(metricNames) == 0 {
		return nil, &core.ValidationError{Field: "metricNames", Message: "at least one metric name is required"}
	}
	if from == "" {
		return nil, &core.ValidationError{Field: "from", Message: "from timestamp is required"}
	}
	if to == "" {
		return nil, &core.ValidationError{Field: "to", Message: "to timestamp is required"}
	}

	q := url.Values{}
	q.Set("dbInstanceId", instanceID)
	q["metricNames"] = metricNames
	q.Set("from", from)
	q.Set("to", to)
	if interval > 0 {
		q.Set("interval", strconv.Itoa(interval))
	}
	path := "/v1.0/metric-statistics?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result GetMetricStatisticsResponse
	if err := core.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
