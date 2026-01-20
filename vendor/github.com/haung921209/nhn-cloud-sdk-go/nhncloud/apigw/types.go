// Package apigw provides API Gateway service client for NHN Cloud.
package apigw

import "time"

// Header represents the API response header
type Header struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// --- Service Types ---

// Service represents an API Gateway service
type Service struct {
	ID          string     `json:"apigwServiceId"`
	Name        string     `json:"apigwServiceName"`
	Description string     `json:"apigwServiceDescription,omitempty"`
	TypeCode    string     `json:"apigwServiceTypeCode"`
	AppKey      string     `json:"appKey"`
	RegionCode  string     `json:"regionCode"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// ListServicesOutput represents the response for services list
type ListServicesOutput struct {
	Header   Header    `json:"header"`
	Services []Service `json:"serviceList"`
}

// GetServiceOutput represents the response for single service
type GetServiceOutput struct {
	Header  Header  `json:"header"`
	Service Service `json:"apigwService"`
}

// CreateServiceInput represents the request to create a service
type CreateServiceInput struct {
	Name        string `json:"apigwServiceName"`
	Description string `json:"apigwServiceDescription,omitempty"`
}

// UpdateServiceInput represents the request to update a service
type UpdateServiceInput struct {
	Name        string `json:"apigwServiceName,omitempty"`
	Description string `json:"apigwServiceDescription,omitempty"`
}

// --- Stage Types ---

// Stage represents a deployment stage
type Stage struct {
	ID                 string     `json:"stageId"`
	ServiceID          string     `json:"apigwServiceId"`
	RegionCode         string     `json:"regionCode"`
	Name               string     `json:"stageName"`
	Description        string     `json:"stageDescription,omitempty"`
	URL                string     `json:"stageUrl,omitempty"`
	BackendEndpointURL string     `json:"backendEndpointUrl,omitempty"`
	ResourceUpdatedAt  *time.Time `json:"resourceUpdatedAt,omitempty"`
	CreatedAt          *time.Time `json:"createdAt,omitempty"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty"`
}

// ListStagesOutput represents the response for stages list
type ListStagesOutput struct {
	Header Header  `json:"header"`
	Stages []Stage `json:"stageList"`
}

// GetStageOutput represents the response for single stage
type GetStageOutput struct {
	Header Header `json:"header"`
	Stage  Stage  `json:"stage"`
}

// CreateStageInput represents the request to create a stage
type CreateStageInput struct {
	Name               string `json:"stageName"`
	Description        string `json:"stageDescription,omitempty"`
	BackendEndpointURL string `json:"backendEndpointUrl,omitempty"`
}

// UpdateStageInput represents the request to update a stage
type UpdateStageInput struct {
	Name               string `json:"stageName,omitempty"`
	Description        string `json:"stageDescription,omitempty"`
	BackendEndpointURL string `json:"backendEndpointUrl,omitempty"`
}

// --- Resource Types ---

// ResourcePlugin represents resource plugin configuration
type ResourcePlugin struct {
	PluginType       string                 `json:"pluginType,omitempty"`
	PluginConfigJSON map[string]interface{} `json:"pluginConfigJson,omitempty"`
}

// Resource represents an API resource
type Resource struct {
	ID                 string          `json:"resourceId"`
	ServiceID          string          `json:"apigwServiceId"`
	ParentID           string          `json:"parentId,omitempty"`
	Path               string          `json:"path"`
	PathParameter      string          `json:"pathParameter,omitempty"`
	MethodType         string          `json:"methodType,omitempty"`
	MethodName         string          `json:"methodName,omitempty"`
	MethodDescription  string          `json:"methodDescription,omitempty"`
	BackendEndpointURL string          `json:"backendEndpointUrl,omitempty"`
	BackendServiceType string          `json:"backendServiceType,omitempty"`
	Plugin             *ResourcePlugin `json:"plugin,omitempty"`
	ResourceType       string          `json:"resourceType"`
	CreatedAt          *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt          *time.Time      `json:"updatedAt,omitempty"`
	Children           []Resource      `json:"children,omitempty"`
}

// ListResourcesOutput represents the response for resources list
type ListResourcesOutput struct {
	Header    Header     `json:"header"`
	Resources []Resource `json:"resourceList"`
}

// GetResourceOutput represents the response for single resource
type GetResourceOutput struct {
	Header   Header   `json:"header"`
	Resource Resource `json:"resource"`
}

// CreateResourceInput represents the request to create a resource
type CreateResourceInput struct {
	Path               string          `json:"path"`
	PathParameter      string          `json:"pathParameter,omitempty"`
	MethodType         string          `json:"methodType,omitempty"`
	MethodName         string          `json:"methodName,omitempty"`
	MethodDescription  string          `json:"methodDescription,omitempty"`
	BackendEndpointURL string          `json:"backendEndpointUrl,omitempty"`
	BackendServiceType string          `json:"backendServiceType,omitempty"`
	Plugin             *ResourcePlugin `json:"plugin,omitempty"`
}

// --- Deploy Types ---

// Deploy represents a deployment
type Deploy struct {
	ID          string     `json:"deployId"`
	ServiceID   string     `json:"apigwServiceId"`
	StageID     string     `json:"stageId"`
	StatusCode  string     `json:"deployStatusCode"`
	Description string     `json:"deployDescription,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// ListDeploysOutput represents the response for deployments list
type ListDeploysOutput struct {
	Header  Header   `json:"header"`
	Deploys []Deploy `json:"deployList"`
}

// GetDeployOutput represents the response for single deployment
type GetDeployOutput struct {
	Header Header `json:"header"`
	Deploy Deploy `json:"deploy"`
}

// CreateDeployInput represents the request to create a deployment
type CreateDeployInput struct {
	Description string `json:"deployDescription,omitempty"`
}

// --- API Key Types ---

// APIKey represents an API key
type APIKey struct {
	ID           string     `json:"apiKeyId"`
	Name         string     `json:"apiKeyName"`
	Description  string     `json:"apiKeyDescription,omitempty"`
	PrimaryKey   string     `json:"primaryKey,omitempty"`
	SecondaryKey string     `json:"secondaryKey,omitempty"`
	StatusCode   string     `json:"apiKeyStatusCode"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
}

// ListAPIKeysOutput represents the response for API keys list
type ListAPIKeysOutput struct {
	Header  Header   `json:"header"`
	APIKeys []APIKey `json:"apiKeyList"`
}

// GetAPIKeyOutput represents the response for single API key
type GetAPIKeyOutput struct {
	Header Header `json:"header"`
	APIKey APIKey `json:"apiKey"`
}

// CreateAPIKeyInput represents the request to create an API key
type CreateAPIKeyInput struct {
	Name        string `json:"apiKeyName"`
	Description string `json:"apiKeyDescription,omitempty"`
}

// UpdateAPIKeyInput represents the request to update an API key
type UpdateAPIKeyInput struct {
	Name        string `json:"apiKeyName,omitempty"`
	Description string `json:"apiKeyDescription,omitempty"`
	StatusCode  string `json:"apiKeyStatusCode,omitempty"`
}

// RegenerateAPIKeyInput represents the request to regenerate an API key
type RegenerateAPIKeyInput struct {
	KeyType string `json:"keyType"` // PRIMARY or SECONDARY
}

// --- Usage Plan Types ---

// UsagePlan represents a usage plan
type UsagePlan struct {
	ID                        string     `json:"usagePlanId"`
	Name                      string     `json:"usagePlanName"`
	Description               string     `json:"usagePlanDescription,omitempty"`
	RateLimitRequestPerSecond int        `json:"rateLimitRequestPerSecond,omitempty"`
	QuotaLimitRequestCount    int        `json:"quotaLimitRequestCount,omitempty"`
	QuotaPeriodUnitCode       string     `json:"quotaPeriodUnitCode,omitempty"` // DAY, MONTH
	CreatedAt                 *time.Time `json:"createdAt,omitempty"`
	UpdatedAt                 *time.Time `json:"updatedAt,omitempty"`
}

// ListUsagePlansOutput represents the response for usage plans list
type ListUsagePlansOutput struct {
	Header     Header      `json:"header"`
	UsagePlans []UsagePlan `json:"usagePlanList"`
}

// GetUsagePlanOutput represents the response for single usage plan
type GetUsagePlanOutput struct {
	Header    Header    `json:"header"`
	UsagePlan UsagePlan `json:"usagePlan"`
}

// CreateUsagePlanInput represents the request to create a usage plan
type CreateUsagePlanInput struct {
	Name                      string `json:"usagePlanName"`
	Description               string `json:"usagePlanDescription,omitempty"`
	RateLimitRequestPerSecond int    `json:"rateLimitRequestPerSecond,omitempty"`
	QuotaLimitRequestCount    int    `json:"quotaLimitRequestCount,omitempty"`
	QuotaPeriodUnitCode       string `json:"quotaPeriodUnitCode,omitempty"`
}

// UpdateUsagePlanInput represents the request to update a usage plan
type UpdateUsagePlanInput struct {
	Name                      string `json:"usagePlanName,omitempty"`
	Description               string `json:"usagePlanDescription,omitempty"`
	RateLimitRequestPerSecond int    `json:"rateLimitRequestPerSecond,omitempty"`
	QuotaLimitRequestCount    int    `json:"quotaLimitRequestCount,omitempty"`
	QuotaPeriodUnitCode       string `json:"quotaPeriodUnitCode,omitempty"`
}

// UsagePlanStage represents a stage connected to a usage plan
type UsagePlanStage struct {
	StageID   string `json:"stageId"`
	StageName string `json:"stageName"`
	StageURL  string `json:"stageUrl"`
}

// ListUsagePlanStagesOutput represents the response for usage plan stages list
type ListUsagePlanStagesOutput struct {
	Header Header           `json:"header"`
	Stages []UsagePlanStage `json:"stageList"`
}

// --- Subscription Types ---

// Subscription represents a subscription
type Subscription struct {
	ID          string     `json:"subscriptionId"`
	UsagePlanID string     `json:"usagePlanId"`
	StageID     string     `json:"stageId"`
	APIKeyID    string     `json:"apiKeyId"`
	APIKeyName  string     `json:"apiKeyName,omitempty"`
	StatusCode  string     `json:"subscriptionStatusCode"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// ListSubscriptionsOutput represents the response for subscriptions list
type ListSubscriptionsOutput struct {
	Header        Header         `json:"header"`
	Subscriptions []Subscription `json:"subscriptionList"`
}

// GetSubscriptionOutput represents the response for single subscription
type GetSubscriptionOutput struct {
	Header       Header       `json:"header"`
	Subscription Subscription `json:"subscription"`
}

// CreateSubscriptionInput represents the request to create a subscription
type CreateSubscriptionInput struct {
	APIKeyID string `json:"apiKeyId"`
}

// --- Model Types ---

// Model represents an API model
type Model struct {
	ID          string     `json:"modelId"`
	ServiceID   string     `json:"apigwServiceId"`
	Name        string     `json:"modelName"`
	Description string     `json:"modelDescription,omitempty"`
	Schema      string     `json:"modelSchema"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// ListModelsOutput represents the response for models list
type ListModelsOutput struct {
	Header Header  `json:"header"`
	Models []Model `json:"modelList"`
}

// GetModelOutput represents the response for single model
type GetModelOutput struct {
	Header Header `json:"header"`
	Model  Model  `json:"model"`
}

// CreateModelInput represents the request to create a model
type CreateModelInput struct {
	Name        string `json:"modelName"`
	Description string `json:"modelDescription,omitempty"`
	Schema      string `json:"modelSchema"`
}

// UpdateModelInput represents the request to update a model
type UpdateModelInput struct {
	Name        string `json:"modelName,omitempty"`
	Description string `json:"modelDescription,omitempty"`
	Schema      string `json:"modelSchema,omitempty"`
}

// --- Statistics Types ---

// MetricPoint represents a single metric data point
type MetricPoint struct {
	Timestamp    string  `json:"timestamp"`
	RequestCount int64   `json:"requestCount"`
	SuccessCount int64   `json:"successCount"`
	FailCount    int64   `json:"failCount"`
	AvgResponse  float64 `json:"avgResponseTimeMs"`
}

// StageMetrics represents stage metrics
type StageMetrics struct {
	TotalCount      int64         `json:"totalCount"`
	SuccessCount    int64         `json:"successCount"`
	FailCount       int64         `json:"failCount"`
	AvgResponseTime float64       `json:"avgResponseTimeMs"`
	MaxResponseTime int64         `json:"maxResponseTimeMs"`
	Data            []MetricPoint `json:"data,omitempty"`
}

// GetStageMetricsOutput represents the response for stage metrics
type GetStageMetricsOutput struct {
	Header  Header       `json:"header"`
	Metrics StageMetrics `json:"metrics"`
}

// MetricsInput represents query options for metrics
type MetricsInput struct {
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	TimeUnit  string `json:"timeUnit,omitempty"` // MINUTE, HOUR, DAY
}

// --- Gateway Response Types ---

// GatewayResponse represents a gateway response
type GatewayResponse struct {
	ID           string            `json:"gatewayResponseId"`
	ServiceID    string            `json:"apigwServiceId"`
	ResponseType string            `json:"responseType"`
	StatusCode   int               `json:"statusCode"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         string            `json:"body,omitempty"`
	CreatedAt    *time.Time        `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time        `json:"updatedAt,omitempty"`
}

// ListGatewayResponsesOutput represents the response for gateway responses list
type ListGatewayResponsesOutput struct {
	Header           Header            `json:"header"`
	GatewayResponses []GatewayResponse `json:"gatewayResponseList"`
}
