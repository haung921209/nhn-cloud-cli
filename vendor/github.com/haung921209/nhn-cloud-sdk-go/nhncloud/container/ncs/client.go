package ncs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/endpoint"
)

type Client struct {
	region        string
	appKey        string
	credentials   credentials.Credentials
	httpClient    *client.Client
	tokenProvider *client.OAuthTokenProvider
	debug         bool
}

func NewClient(region, appKey string, creds credentials.Credentials, hc *http.Client, debug bool) *Client {
	c := &Client{
		region:      region,
		appKey:      appKey,
		credentials: creds,
		debug:       debug,
	}

	if creds != nil {
		c.tokenProvider = client.NewOAuthTokenProvider(
			creds.GetAccessKeyID(),
			creds.GetSecretAccessKey(),
		)
		c.initHTTPClient()
	}

	return c
}

func (c *Client) initHTTPClient() {
	baseURL := endpoint.ResolveWithAppKey(endpoint.ServiceNCS, c.region, c.appKey)
	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}
	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
}

func (c *Client) ListWorkloads(ctx context.Context, namespace string) (*ListWorkloadsOutput, error) {
	path := "/workloads"
	if namespace != "" {
		path += "?namespace=" + namespace
	}
	var out ListWorkloadsOutput
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list workloads: %w", err)
	}
	return &out, nil
}

func (c *Client) GetWorkload(ctx context.Context, workloadID string) (*GetWorkloadOutput, error) {
	var out GetWorkloadOutput
	if err := c.httpClient.GET(ctx, "/workloads/"+workloadID, &out); err != nil {
		return nil, fmt.Errorf("get workload %s: %w", workloadID, err)
	}
	return &out, nil
}

func (c *Client) CreateWorkload(ctx context.Context, input *CreateWorkloadInput) (*CreateWorkloadOutput, error) {
	var out CreateWorkloadOutput
	if err := c.httpClient.POST(ctx, "/workloads", input, &out); err != nil {
		return nil, fmt.Errorf("create workload: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateWorkload(ctx context.Context, workloadID string, input *UpdateWorkloadInput) (*GetWorkloadOutput, error) {
	var out GetWorkloadOutput
	if err := c.httpClient.PUT(ctx, "/workloads/"+workloadID, input, &out); err != nil {
		return nil, fmt.Errorf("update workload %s: %w", workloadID, err)
	}
	return &out, nil
}

func (c *Client) DeleteWorkload(ctx context.Context, workloadID string) error {
	if err := c.httpClient.DELETE(ctx, "/workloads/"+workloadID, nil); err != nil {
		return fmt.Errorf("delete workload %s: %w", workloadID, err)
	}
	return nil
}

func (c *Client) RestartWorkload(ctx context.Context, workloadID string) error {
	if err := c.httpClient.POST(ctx, "/workloads/"+workloadID+"/restart", nil, nil); err != nil {
		return fmt.Errorf("restart workload %s: %w", workloadID, err)
	}
	return nil
}

func (c *Client) ScaleWorkload(ctx context.Context, workloadID string, replicas int) error {
	req := map[string]int{"replicas": replicas}
	if err := c.httpClient.POST(ctx, "/workloads/"+workloadID+"/scale", req, nil); err != nil {
		return fmt.Errorf("scale workload %s: %w", workloadID, err)
	}
	return nil
}

func (c *Client) ListTemplates(ctx context.Context) (*ListTemplatesOutput, error) {
	var out ListTemplatesOutput
	if err := c.httpClient.GET(ctx, "/templates", &out); err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	return &out, nil
}

func (c *Client) GetTemplate(ctx context.Context, templateID string) (*GetTemplateOutput, error) {
	var out GetTemplateOutput
	if err := c.httpClient.GET(ctx, "/templates/"+templateID, &out); err != nil {
		return nil, fmt.Errorf("get template %s: %w", templateID, err)
	}
	return &out, nil
}

func (c *Client) ListServices(ctx context.Context, namespace string) (*ListServicesOutput, error) {
	path := "/services"
	if namespace != "" {
		path += "?namespace=" + namespace
	}
	var out ListServicesOutput
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	return &out, nil
}

func (c *Client) GetService(ctx context.Context, serviceID string) (*GetServiceOutput, error) {
	var out GetServiceOutput
	if err := c.httpClient.GET(ctx, "/services/"+serviceID, &out); err != nil {
		return nil, fmt.Errorf("get service %s: %w", serviceID, err)
	}
	return &out, nil
}

func (c *Client) CreateService(ctx context.Context, input *CreateServiceInput) (*CreateServiceOutput, error) {
	var out CreateServiceOutput
	if err := c.httpClient.POST(ctx, "/services", input, &out); err != nil {
		return nil, fmt.Errorf("create service: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteService(ctx context.Context, serviceID string) error {
	if err := c.httpClient.DELETE(ctx, "/services/"+serviceID, nil); err != nil {
		return fmt.Errorf("delete service %s: %w", serviceID, err)
	}
	return nil
}

func (c *Client) GetWorkloadLogs(ctx context.Context, workloadID string, tailLines int, sinceSeconds int) (*GetWorkloadLogsOutput, error) {
	path := fmt.Sprintf("/workloads/%s/logs?tailLines=%d", workloadID, tailLines)
	if sinceSeconds > 0 {
		path += fmt.Sprintf("&sinceSeconds=%d", sinceSeconds)
	}
	var out GetWorkloadLogsOutput
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get workload logs %s: %w", workloadID, err)
	}
	return &out, nil
}

func (c *Client) ConfigureHealthCheck(ctx context.Context, workloadID string, config *HealthCheckConfig) error {
	if err := c.httpClient.PUT(ctx, "/workloads/"+workloadID+"/health-checks", config, nil); err != nil {
		return fmt.Errorf("configure health check %s: %w", workloadID, err)
	}
	return nil
}

func (c *Client) GetHealthCheckStatus(ctx context.Context, workloadID string) (*GetHealthCheckStatusOutput, error) {
	var out GetHealthCheckStatusOutput
	if err := c.httpClient.GET(ctx, "/workloads/"+workloadID+"/health-checks", &out); err != nil {
		return nil, fmt.Errorf("get health check status %s: %w", workloadID, err)
	}
	return &out, nil
}

func (c *Client) UpdateResourceLimits(ctx context.Context, workloadID string, input *UpdateResourceLimitsInput) error {
	if err := c.httpClient.PUT(ctx, "/workloads/"+workloadID+"/resources", input, nil); err != nil {
		return fmt.Errorf("update resource limits %s: %w", workloadID, err)
	}
	return nil
}

func (c *Client) GetWorkloadEvents(ctx context.Context, workloadID string) (*GetWorkloadEventsOutput, error) {
	var out GetWorkloadEventsOutput
	if err := c.httpClient.GET(ctx, "/workloads/"+workloadID+"/events", &out); err != nil {
		return nil, fmt.Errorf("get workload events %s: %w", workloadID, err)
	}
	return &out, nil
}

func (c *Client) ListVolumes(ctx context.Context) (*ListVolumesOutput, error) {
	var out ListVolumesOutput
	if err := c.httpClient.GET(ctx, "/volumes", &out); err != nil {
		return nil, fmt.Errorf("list volumes: %w", err)
	}
	return &out, nil
}

func (c *Client) AttachVolume(ctx context.Context, workloadID string, input *VolumeAttachInput) (*AttachVolumeOutput, error) {
	var out AttachVolumeOutput
	if err := c.httpClient.POST(ctx, "/workloads/"+workloadID+"/volumes", input, &out); err != nil {
		return nil, fmt.Errorf("attach volume to workload %s: %w", workloadID, err)
	}
	return &out, nil
}

func (c *Client) ExecWorkloadContainer(ctx context.Context, workloadID string, input *ExecInput) (*ExecOutput, error) {
	var out ExecOutput
	if err := c.httpClient.POST(ctx, "/workloads/"+workloadID+"/exec", input, &out); err != nil {
		return nil, fmt.Errorf("exec in workload container %s: %w", workloadID, err)
	}
	return &out, nil
}

func (c *Client) GetContainerStatus(ctx context.Context, workloadID string) (*GetContainerStatusOutput, error) {
	var out GetContainerStatusOutput
	if err := c.httpClient.GET(ctx, "/workloads/"+workloadID+"/containers/status", &out); err != nil {
		return nil, fmt.Errorf("get container status %s: %w", workloadID, err)
	}
	return &out, nil
}

func (c *Client) ConfigureAutoScaling(ctx context.Context, workloadID string, input *ConfigureAutoScalingInput) error {
	if err := c.httpClient.PUT(ctx, "/workloads/"+workloadID+"/autoscaling", input, nil); err != nil {
		return fmt.Errorf("configure autoscaling %s: %w", workloadID, err)
	}
	return nil
}

func (c *Client) GetAutoScalingStatus(ctx context.Context, workloadID string) (*GetAutoScalingStatusOutput, error) {
	var out GetAutoScalingStatusOutput
	if err := c.httpClient.GET(ctx, "/workloads/"+workloadID+"/autoscaling", &out); err != nil {
		return nil, fmt.Errorf("get autoscaling status %s: %w", workloadID, err)
	}
	return &out, nil
}

func (c *Client) ListNetworkPolicies(ctx context.Context) (*ListNetworkPoliciesOutput, error) {
	var out ListNetworkPoliciesOutput
	if err := c.httpClient.GET(ctx, "/network-policies", &out); err != nil {
		return nil, fmt.Errorf("list network policies: %w", err)
	}
	return &out, nil
}

func (c *Client) GetNetworkPolicy(ctx context.Context, policyID string) (*GetNetworkPolicyOutput, error) {
	var out GetNetworkPolicyOutput
	if err := c.httpClient.GET(ctx, "/network-policies/"+policyID, &out); err != nil {
		return nil, fmt.Errorf("get network policy %s: %w", policyID, err)
	}
	return &out, nil
}

func (c *Client) CreateNetworkPolicy(ctx context.Context, input *CreateNetworkPolicyInput) (*GetNetworkPolicyOutput, error) {
	var out GetNetworkPolicyOutput
	if err := c.httpClient.POST(ctx, "/network-policies", input, &out); err != nil {
		return nil, fmt.Errorf("create network policy: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateNetworkPolicy(ctx context.Context, policyID string, input *UpdateNetworkPolicyInput) (*GetNetworkPolicyOutput, error) {
	var out GetNetworkPolicyOutput
	if err := c.httpClient.PUT(ctx, "/network-policies/"+policyID, input, &out); err != nil {
		return nil, fmt.Errorf("update network policy %s: %w", policyID, err)
	}
	return &out, nil
}

func (c *Client) DeleteNetworkPolicy(ctx context.Context, policyID string) error {
	if err := c.httpClient.DELETE(ctx, "/network-policies/"+policyID, nil); err != nil {
		return fmt.Errorf("delete network policy %s: %w", policyID, err)
	}
	return nil
}
