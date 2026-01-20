package apigw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client handles API Gateway API operations
type Client struct {
	region      string
	appKey      string
	accessKeyID string
	secretKey   string
	httpClient  *http.Client
	debug       bool
}

// NewClient creates a new API Gateway client
func NewClient(region, appKey, accessKeyID, secretKey string, hc *http.Client, debug bool) *Client {
	c := &Client{
		region:      region,
		appKey:      appKey,
		accessKeyID: accessKeyID,
		secretKey:   secretKey,
		httpClient:  hc,
		debug:       debug,
	}

	if hc == nil {
		c.httpClient = http.DefaultClient
	}

	return c
}

func (c *Client) getBaseURL() string {
	return fmt.Sprintf("https://%s-apigateway.api.nhncloudservice.com", strings.ToLower(c.region))
}

func (c *Client) buildPath(path string) string {
	return fmt.Sprintf("/v1.0/appkeys/%s%s", c.appKey, path)
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, query url.Values) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)

		if c.debug {
			fmt.Printf("[DEBUG] APIGW Request Body: %s\n", string(jsonBody))
		}
	}

	fullURL := c.getBaseURL() + c.buildPath(path)
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TC-AUTHENTICATION-ID", c.accessKeyID)
	req.Header.Set("X-TC-AUTHENTICATION-SECRET", c.secretKey)

	if c.debug {
		fmt.Printf("[DEBUG] APIGW Request: %s %s\n", method, fullURL)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if c.debug {
		fmt.Printf("[DEBUG] APIGW Response: %d %s\n", resp.StatusCode, resp.Status)
		if len(respBody) > 0 {
			fmt.Printf("[DEBUG] APIGW Response Body: %s\n", string(respBody))
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %d %s - %s", resp.StatusCode, resp.Status, string(respBody))
	}

	return respBody, nil
}

// --- Service Operations ---

// ListServices lists all API Gateway services
func (c *Client) ListServices(ctx context.Context) (*ListServicesOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/services", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	var out ListServicesOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// GetService retrieves a specific service
func (c *Client) GetService(ctx context.Context, serviceID string) (*GetServiceOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/services/"+serviceID, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get service %s: %w", serviceID, err)
	}
	var out GetServiceOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// CreateService creates a new service
func (c *Client) CreateService(ctx context.Context, input *CreateServiceInput) (*GetServiceOutput, error) {
	respBody, err := c.doRequest(ctx, "POST", "/services", input, nil)
	if err != nil {
		return nil, fmt.Errorf("create service: %w", err)
	}
	var out GetServiceOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// UpdateService updates a service
func (c *Client) UpdateService(ctx context.Context, serviceID string, input *UpdateServiceInput) (*GetServiceOutput, error) {
	respBody, err := c.doRequest(ctx, "PUT", "/services/"+serviceID, input, nil)
	if err != nil {
		return nil, fmt.Errorf("update service %s: %w", serviceID, err)
	}
	var out GetServiceOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// DeleteService deletes a service
func (c *Client) DeleteService(ctx context.Context, serviceID string) error {
	if _, err := c.doRequest(ctx, "DELETE", "/services/"+serviceID, nil, nil); err != nil {
		return fmt.Errorf("delete service %s: %w", serviceID, err)
	}
	return nil
}

// --- Resource Operations ---

// ListResources lists resources for a service
func (c *Client) ListResources(ctx context.Context, serviceID string) (*ListResourcesOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/services/"+serviceID+"/resources", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list resources: %w", err)
	}
	var out ListResourcesOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// CreateResource creates a resource with path and optionally method
func (c *Client) CreateResource(ctx context.Context, serviceID string, input *CreateResourceInput) (*GetResourceOutput, error) {
	respBody, err := c.doRequest(ctx, "POST", "/services/"+serviceID+"/resources", input, nil)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}
	var out GetResourceOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// DeleteResource deletes a resource
func (c *Client) DeleteResource(ctx context.Context, serviceID, resourceID string) error {
	if _, err := c.doRequest(ctx, "DELETE", "/services/"+serviceID+"/resources/"+resourceID, nil, nil); err != nil {
		return fmt.Errorf("delete resource %s: %w", resourceID, err)
	}
	return nil
}

// --- Stage Operations ---

// ListStages lists stages for a service
func (c *Client) ListStages(ctx context.Context, serviceID string) (*ListStagesOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/services/"+serviceID+"/stages", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list stages: %w", err)
	}
	var out ListStagesOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// CreateStage creates a stage
func (c *Client) CreateStage(ctx context.Context, serviceID string, input *CreateStageInput) (*GetStageOutput, error) {
	respBody, err := c.doRequest(ctx, "POST", "/services/"+serviceID+"/stages", input, nil)
	if err != nil {
		return nil, fmt.Errorf("create stage: %w", err)
	}
	var out GetStageOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// UpdateStage updates a stage
func (c *Client) UpdateStage(ctx context.Context, serviceID, stageID string, input *UpdateStageInput) (*GetStageOutput, error) {
	respBody, err := c.doRequest(ctx, "PUT", "/services/"+serviceID+"/stages/"+stageID, input, nil)
	if err != nil {
		return nil, fmt.Errorf("update stage %s: %w", stageID, err)
	}
	var out GetStageOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// DeleteStage deletes a stage
func (c *Client) DeleteStage(ctx context.Context, serviceID, stageID string) error {
	if _, err := c.doRequest(ctx, "DELETE", "/services/"+serviceID+"/stages/"+stageID, nil, nil); err != nil {
		return fmt.Errorf("delete stage %s: %w", stageID, err)
	}
	return nil
}

// --- Deploy Operations ---

// ListDeploys lists deployments for a stage
func (c *Client) ListDeploys(ctx context.Context, serviceID, stageID string) (*ListDeploysOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/services/"+serviceID+"/stages/"+stageID+"/deploys", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list deploys: %w", err)
	}
	var out ListDeploysOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// DeployStage deploys a stage
func (c *Client) DeployStage(ctx context.Context, serviceID, stageID string, input *CreateDeployInput) (*GetDeployOutput, error) {
	respBody, err := c.doRequest(ctx, "POST", "/services/"+serviceID+"/stages/"+stageID+"/deploys", input, nil)
	if err != nil {
		return nil, fmt.Errorf("deploy stage: %w", err)
	}
	var out GetDeployOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// GetLatestDeploy gets the latest deployment for a stage
func (c *Client) GetLatestDeploy(ctx context.Context, serviceID, stageID string) (*GetDeployOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/services/"+serviceID+"/stages/"+stageID+"/deploys/latest", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get latest deploy: %w", err)
	}
	var out GetDeployOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// DeleteDeploy deletes a deployment
func (c *Client) DeleteDeploy(ctx context.Context, serviceID, stageID, deployID string) error {
	if _, err := c.doRequest(ctx, "DELETE", "/services/"+serviceID+"/stages/"+stageID+"/deploys/"+deployID, nil, nil); err != nil {
		return fmt.Errorf("delete deploy %s: %w", deployID, err)
	}
	return nil
}

// RollbackDeploy rolls back to a specific deployment
func (c *Client) RollbackDeploy(ctx context.Context, serviceID, stageID, deployID string) (*GetDeployOutput, error) {
	respBody, err := c.doRequest(ctx, "POST", "/services/"+serviceID+"/stages/"+stageID+"/deploys/"+deployID+"/rollback", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("rollback deploy: %w", err)
	}
	var out GetDeployOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// --- API Key Operations ---

// ListAPIKeys lists all API keys
func (c *Client) ListAPIKeys(ctx context.Context) (*ListAPIKeysOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/apikeys", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list API keys: %w", err)
	}
	var out ListAPIKeysOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// CreateAPIKey creates an API key
func (c *Client) CreateAPIKey(ctx context.Context, input *CreateAPIKeyInput) (*GetAPIKeyOutput, error) {
	respBody, err := c.doRequest(ctx, "POST", "/apikeys", input, nil)
	if err != nil {
		return nil, fmt.Errorf("create API key: %w", err)
	}
	var out GetAPIKeyOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// UpdateAPIKey updates an API key
func (c *Client) UpdateAPIKey(ctx context.Context, apiKeyID string, input *UpdateAPIKeyInput) (*GetAPIKeyOutput, error) {
	respBody, err := c.doRequest(ctx, "PUT", "/apikeys/"+apiKeyID, input, nil)
	if err != nil {
		return nil, fmt.Errorf("update API key %s: %w", apiKeyID, err)
	}
	var out GetAPIKeyOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// DeleteAPIKey deletes an API key
func (c *Client) DeleteAPIKey(ctx context.Context, apiKeyID string) error {
	if _, err := c.doRequest(ctx, "DELETE", "/apikeys/"+apiKeyID, nil, nil); err != nil {
		return fmt.Errorf("delete API key %s: %w", apiKeyID, err)
	}
	return nil
}

// RegenerateAPIKey regenerates an API key
func (c *Client) RegenerateAPIKey(ctx context.Context, apiKeyID string, input *RegenerateAPIKeyInput) (*GetAPIKeyOutput, error) {
	respBody, err := c.doRequest(ctx, "POST", "/apikeys/"+apiKeyID+"/regenerate", input, nil)
	if err != nil {
		return nil, fmt.Errorf("regenerate API key: %w", err)
	}
	var out GetAPIKeyOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// --- Usage Plan Operations ---

// ListUsagePlans lists all usage plans
func (c *Client) ListUsagePlans(ctx context.Context) (*ListUsagePlansOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/usage-plans", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list usage plans: %w", err)
	}
	var out ListUsagePlansOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// GetUsagePlan retrieves a specific usage plan
func (c *Client) GetUsagePlan(ctx context.Context, usagePlanID string) (*GetUsagePlanOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/usage-plans/"+usagePlanID, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get usage plan %s: %w", usagePlanID, err)
	}
	var out GetUsagePlanOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// CreateUsagePlan creates a usage plan
func (c *Client) CreateUsagePlan(ctx context.Context, input *CreateUsagePlanInput) (*GetUsagePlanOutput, error) {
	respBody, err := c.doRequest(ctx, "POST", "/usage-plans", input, nil)
	if err != nil {
		return nil, fmt.Errorf("create usage plan: %w", err)
	}
	var out GetUsagePlanOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// UpdateUsagePlan updates a usage plan
func (c *Client) UpdateUsagePlan(ctx context.Context, usagePlanID string, input *UpdateUsagePlanInput) (*GetUsagePlanOutput, error) {
	respBody, err := c.doRequest(ctx, "PUT", "/usage-plans/"+usagePlanID, input, nil)
	if err != nil {
		return nil, fmt.Errorf("update usage plan %s: %w", usagePlanID, err)
	}
	var out GetUsagePlanOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// DeleteUsagePlan deletes a usage plan
func (c *Client) DeleteUsagePlan(ctx context.Context, usagePlanID string) error {
	if _, err := c.doRequest(ctx, "DELETE", "/usage-plans/"+usagePlanID, nil, nil); err != nil {
		return fmt.Errorf("delete usage plan %s: %w", usagePlanID, err)
	}
	return nil
}

// ListUsagePlanStages lists stages connected to a usage plan
func (c *Client) ListUsagePlanStages(ctx context.Context, usagePlanID string) (*ListUsagePlanStagesOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/usage-plans/"+usagePlanID+"/stages", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list usage plan stages: %w", err)
	}
	var out ListUsagePlanStagesOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// ConnectStageToUsagePlan connects a stage to a usage plan
func (c *Client) ConnectStageToUsagePlan(ctx context.Context, usagePlanID, stageID string) error {
	if _, err := c.doRequest(ctx, "POST", "/usage-plans/"+usagePlanID+"/stages/"+stageID, nil, nil); err != nil {
		return fmt.Errorf("connect stage to usage plan: %w", err)
	}
	return nil
}

// DisconnectStageFromUsagePlan disconnects a stage from a usage plan
func (c *Client) DisconnectStageFromUsagePlan(ctx context.Context, usagePlanID, stageID string) error {
	if _, err := c.doRequest(ctx, "DELETE", "/usage-plans/"+usagePlanID+"/stages/"+stageID, nil, nil); err != nil {
		return fmt.Errorf("disconnect stage from usage plan: %w", err)
	}
	return nil
}

// --- Subscription Operations ---

// ListSubscriptions lists subscriptions for a usage plan and stage
func (c *Client) ListSubscriptions(ctx context.Context, usagePlanID, stageID string) (*ListSubscriptionsOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/usage-plans/"+usagePlanID+"/stages/"+stageID+"/subscriptions", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	var out ListSubscriptionsOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// CreateSubscription creates a subscription
func (c *Client) CreateSubscription(ctx context.Context, usagePlanID, stageID string, input *CreateSubscriptionInput) (*GetSubscriptionOutput, error) {
	respBody, err := c.doRequest(ctx, "POST", "/usage-plans/"+usagePlanID+"/stages/"+stageID+"/subscriptions", input, nil)
	if err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}
	var out GetSubscriptionOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// DeleteSubscription deletes a subscription
func (c *Client) DeleteSubscription(ctx context.Context, usagePlanID, stageID, apiKeyID string) error {
	req := map[string]string{"apiKeyId": apiKeyID}
	if _, err := c.doRequest(ctx, "DELETE", "/usage-plans/"+usagePlanID+"/stages/"+stageID+"/subscriptions", req, nil); err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	return nil
}

// --- Model Operations ---

// ListModels lists models for a service
func (c *Client) ListModels(ctx context.Context, serviceID string) (*ListModelsOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/services/"+serviceID+"/models", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list models: %w", err)
	}
	var out ListModelsOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// CreateModel creates a model
func (c *Client) CreateModel(ctx context.Context, serviceID string, input *CreateModelInput) (*GetModelOutput, error) {
	respBody, err := c.doRequest(ctx, "POST", "/services/"+serviceID+"/models", input, nil)
	if err != nil {
		return nil, fmt.Errorf("create model: %w", err)
	}
	var out GetModelOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// UpdateModel updates a model
func (c *Client) UpdateModel(ctx context.Context, serviceID, modelID string, input *UpdateModelInput) (*GetModelOutput, error) {
	respBody, err := c.doRequest(ctx, "PUT", "/services/"+serviceID+"/models/"+modelID, input, nil)
	if err != nil {
		return nil, fmt.Errorf("update model %s: %w", modelID, err)
	}
	var out GetModelOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// DeleteModel deletes a model
func (c *Client) DeleteModel(ctx context.Context, serviceID, modelID string) error {
	if _, err := c.doRequest(ctx, "DELETE", "/services/"+serviceID+"/models/"+modelID, nil, nil); err != nil {
		return fmt.Errorf("delete model %s: %w", modelID, err)
	}
	return nil
}

// --- Gateway Response Operations ---

// ListGatewayResponses lists gateway responses for a service
func (c *Client) ListGatewayResponses(ctx context.Context, serviceID string) (*ListGatewayResponsesOutput, error) {
	respBody, err := c.doRequest(ctx, "GET", "/services/"+serviceID+"/gateway-responses", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list gateway responses: %w", err)
	}
	var out ListGatewayResponsesOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

// DeleteGatewayResponse deletes a gateway response
func (c *Client) DeleteGatewayResponse(ctx context.Context, serviceID, gatewayResponseID string) error {
	if _, err := c.doRequest(ctx, "DELETE", "/services/"+serviceID+"/gateway-responses/"+gatewayResponseID, nil, nil); err != nil {
		return fmt.Errorf("delete gateway response %s: %w", gatewayResponseID, err)
	}
	return nil
}

// --- Metrics Operations ---

// GetStageMetrics retrieves metrics for a stage
func (c *Client) GetStageMetrics(ctx context.Context, serviceID, stageID string, input *MetricsInput) (*GetStageMetricsOutput, error) {
	query := url.Values{}
	if input != nil {
		if input.StartTime != "" {
			query.Set("startTime", input.StartTime)
		}
		if input.EndTime != "" {
			query.Set("endTime", input.EndTime)
		}
		if input.TimeUnit != "" {
			query.Set("timeUnit", input.TimeUnit)
		}
	}
	respBody, err := c.doRequest(ctx, "GET", "/services/"+serviceID+"/stages/"+stageID+"/metrics", nil, query)
	if err != nil {
		return nil, fmt.Errorf("get stage metrics: %w", err)
	}
	var out GetStageMetricsOutput
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}
