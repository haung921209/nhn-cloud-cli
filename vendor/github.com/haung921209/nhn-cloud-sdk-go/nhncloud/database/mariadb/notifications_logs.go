package mariadb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// NotificationGroup represents a notification group
type NotificationGroup struct {
	NotificationGroupID   string   `json:"notificationGroupId"`
	NotificationGroupName string   `json:"notificationGroupName"`
	IsEnabled             bool     `json:"isEnabled"`
	NotifyEmail           []string `json:"notifyEmail,omitempty"`
	NotifySms             []string `json:"notifySms,omitempty"`
	CreatedAt             string   `json:"createdAt,omitempty"`
	UpdatedAt             string   `json:"updatedAt,omitempty"`
}

// ListNotificationGroupsResponse is the response for ListNotificationGroups
type ListNotificationGroupsResponse struct {
	MariaDBResponse
	NotificationGroups []NotificationGroup `json:"notificationGroups"`
}

// ListNotificationGroups retrieves all notification groups.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#_69
func (c *Client) ListNotificationGroups(ctx context.Context) (*ListNotificationGroupsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "/v3.0/notification-groups", nil)
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

// GetNotificationGroupResponse is the response for GetNotificationGroup
type GetNotificationGroupResponse struct {
	MariaDBResponse
	NotificationGroup NotificationGroup `json:"notificationGroup"`
}

// GetNotificationGroup retrieves a specific notification group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#_70
func (c *Client) GetNotificationGroup(ctx context.Context, groupID string) (*GetNotificationGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "notification group ID is required"}
	}

	path := fmt.Sprintf("/v3.0/notification-groups/%s", groupID)
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
	NotifyEmail           []string `json:"notifyEmail,omitempty"`
	NotifySms             []string `json:"notifySms,omitempty"`
}

// CreateNotificationGroupResponse is the response for CreateNotificationGroup
type CreateNotificationGroupResponse struct {
	MariaDBResponse
	NotificationGroupID string `json:"notificationGroupId"`
}

// CreateNotificationGroup creates a new notification group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#_71
func (c *Client) CreateNotificationGroup(ctx context.Context, req *CreateNotificationGroupRequest) (*CreateNotificationGroupResponse, error) {
	if req.NotificationGroupName == "" {
		return nil, &core.ValidationError{Field: "NotificationGroupName", Message: "notification group name is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "/v3.0/notification-groups", bytes.NewReader(body))
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
	NotificationGroupName *string  `json:"notificationGroupName,omitempty"`
	IsEnabled             *bool    `json:"isEnabled,omitempty"`
	NotifyEmail           []string `json:"notifyEmail,omitempty"`
	NotifySms             []string `json:"notifySms,omitempty"`
}

// UpdateNotificationGroupResponse is the response for UpdateNotificationGroup
type UpdateNotificationGroupResponse struct {
	MariaDBResponse
}

// UpdateNotificationGroup updates a notification group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#_72
func (c *Client) UpdateNotificationGroup(ctx context.Context, groupID string, req *UpdateNotificationGroupRequest) (*UpdateNotificationGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "notification group ID is required"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/v3.0/notification-groups/%s", groupID)
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
	MariaDBResponse
}

// DeleteNotificationGroup deletes a notification group.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#_73
func (c *Client) DeleteNotificationGroup(ctx context.Context, groupID string) (*DeleteNotificationGroupResponse, error) {
	if groupID == "" {
		return nil, &core.ValidationError{Field: "groupID", Message: "notification group ID is required"}
	}

	path := fmt.Sprintf("/v3.0/notification-groups/%s", groupID)
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

// LogFile represents a log file
type LogFile struct {
	LogFileName string `json:"logFileName"`
	LogFileSize int64  `json:"logFileSize"`
	ModifiedAt  string `json:"modifiedAt,omitempty"`
}

// ListLogFilesResponse is the response for ListLogFiles
type ListLogFilesResponse struct {
	MariaDBResponse
	LogFiles []LogFile `json:"logFiles"`
}

// ListLogFiles retrieves log files for an instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#_74
func (c *Client) ListLogFiles(ctx context.Context, instanceID string) (*ListLogFilesResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}

	path := fmt.Sprintf("/v3.0/db-instances/%s/log-files", instanceID)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.core.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var result ListLogFilesResponse
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
	MariaDBResponse
	Metrics []Metric `json:"metrics"`
}

// ListMetrics retrieves available metrics.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#_75
func (c *Client) ListMetrics(ctx context.Context) (*ListMetricsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "/v3.0/metrics", nil)
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
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

// GetMetricStatisticsResponse is the response for GetMetricStatistics
type GetMetricStatisticsResponse struct {
	MariaDBResponse
	MetricStatistics []MetricStatistic `json:"metricStatistics"`
}

// GetMetricStatistics retrieves metric statistics for an instance.
//
// API Reference:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MariaDB/ko/api-guide-v3.0/#_76
func (c *Client) GetMetricStatistics(ctx context.Context, instanceID, from, to string, interval int) (*GetMetricStatisticsResponse, error) {
	if instanceID == "" {
		return nil, &core.ValidationError{Field: "instanceID", Message: "instance ID is required"}
	}
	if from == "" {
		return nil, &core.ValidationError{Field: "from", Message: "from timestamp is required"}
	}
	if to == "" {
		return nil, &core.ValidationError{Field: "to", Message: "to timestamp is required"}
	}

	path := fmt.Sprintf("/v3.0/metric-statistics?dbInstanceId=%s&from=%s&to=%s&interval=%d",
		instanceID, from, to, interval)
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
