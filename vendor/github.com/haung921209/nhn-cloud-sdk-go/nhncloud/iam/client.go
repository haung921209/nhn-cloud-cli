package iam

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
	credentials   credentials.Credentials
	httpClient    *client.Client
	tokenProvider *client.OAuthTokenProvider
	debug         bool
}

func NewClient(region string, creds credentials.Credentials, hc *http.Client, debug bool) *Client {
	c := &Client{
		region:      region,
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
	baseURL := endpoint.Resolve(endpoint.ServiceIAM, c.region)
	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}
	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
}

func (c *Client) ListOrganizations(ctx context.Context) (*ListOrganizationsOutput, error) {
	var out ListOrganizationsOutput
	if err := c.httpClient.GET(ctx, "/v1/organizations", &out); err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}
	return &out, nil
}

func (c *Client) GetOrganization(ctx context.Context, orgID string) (*GetOrganizationOutput, error) {
	var out GetOrganizationOutput
	if err := c.httpClient.GET(ctx, "/v1/organizations/"+orgID, &out); err != nil {
		return nil, fmt.Errorf("get organization %s: %w", orgID, err)
	}
	return &out, nil
}

func (c *Client) ListProjects(ctx context.Context, orgID string) (*ListProjectsOutput, error) {
	var out ListProjectsOutput
	path := fmt.Sprintf("/v1/organizations/%s/projects", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return &out, nil
}

func (c *Client) GetProject(ctx context.Context, orgID, projectID string) (*GetProjectOutput, error) {
	var out GetProjectOutput
	path := fmt.Sprintf("/v1/organizations/%s/projects/%s", orgID, projectID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get project %s: %w", projectID, err)
	}
	return &out, nil
}

func (c *Client) ListMembers(ctx context.Context, orgID string) (*ListMembersOutput, error) {
	var out ListMembersOutput
	path := fmt.Sprintf("/v1/organizations/%s/members", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	return &out, nil
}

func (c *Client) GetMember(ctx context.Context, orgID, memberID string) (*GetMemberOutput, error) {
	var out GetMemberOutput
	path := fmt.Sprintf("/v1/organizations/%s/members/%s", orgID, memberID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get member %s: %w", memberID, err)
	}
	return &out, nil
}

func (c *Client) InviteMember(ctx context.Context, orgID string, input *InviteMemberInput) (*InviteMemberOutput, error) {
	var out InviteMemberOutput
	path := fmt.Sprintf("/v1/organizations/%s/members", orgID)
	if err := c.httpClient.POST(ctx, path, input, &out); err != nil {
		return nil, fmt.Errorf("invite member: %w", err)
	}
	return &out, nil
}

func (c *Client) RemoveMember(ctx context.Context, orgID, memberID string) error {
	path := fmt.Sprintf("/v1/organizations/%s/members/%s", orgID, memberID)
	if err := c.httpClient.DELETE(ctx, path, nil); err != nil {
		return fmt.Errorf("remove member %s: %w", memberID, err)
	}
	return nil
}

func (c *Client) UpdateMember(ctx context.Context, orgID, memberID string, input *UpdateMemberInput) (*UpdateMemberOutput, error) {
	var out UpdateMemberOutput
	path := fmt.Sprintf("/v1/organizations/%s/members/%s", orgID, memberID)
	if err := c.httpClient.PUT(ctx, path, input, &out); err != nil {
		return nil, fmt.Errorf("update member %s: %w", memberID, err)
	}
	return &out, nil
}

func (c *Client) ListOrganizationRoles(ctx context.Context, orgID string) (*ListRolesOutput, error) {
	var out ListRolesOutput
	path := fmt.Sprintf("/v1/organizations/%s/roles", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list organization roles: %w", err)
	}
	return &out, nil
}

func (c *Client) ListOrganizationRoleGroups(ctx context.Context, orgID string) (*ListRoleGroupsOutput, error) {
	var out ListRoleGroupsOutput
	path := fmt.Sprintf("/v1/organizations/%s/org-role-groups", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list organization role groups: %w", err)
	}
	return &out, nil
}

func (c *Client) GetOrganizationRoleGroup(ctx context.Context, orgID, roleGroupID string) (*GetRoleGroupOutput, error) {
	var out GetRoleGroupOutput
	path := fmt.Sprintf("/v1/organizations/%s/org-role-groups/%s", orgID, roleGroupID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get organization role group %s: %w", roleGroupID, err)
	}
	return &out, nil
}

func (c *Client) CreateOrganizationRoleGroup(ctx context.Context, orgID string, input *CreateRoleGroupInput) (*CreateRoleGroupOutput, error) {
	var out CreateRoleGroupOutput
	path := fmt.Sprintf("/v1/organizations/%s/org-role-groups", orgID)
	if err := c.httpClient.POST(ctx, path, input, &out); err != nil {
		return nil, fmt.Errorf("create organization role group: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteOrganizationRoleGroups(ctx context.Context, orgID string, roleGroupIDs []string) error {
	path := fmt.Sprintf("/v1/organizations/%s/org-role-groups", orgID)
	body := map[string][]string{"roleGroupIds": roleGroupIDs}
	if err := c.httpClient.DELETE(ctx, path, body); err != nil {
		return fmt.Errorf("delete organization role groups: %w", err)
	}
	return nil
}

func (c *Client) UpdateOrganizationRoleGroupInfo(ctx context.Context, orgID, roleGroupID string, input *UpdateRoleGroupInfoInput) error {
	path := fmt.Sprintf("/v1/organizations/%s/org-role-groups/%s/infos", orgID, roleGroupID)
	if err := c.httpClient.PUT(ctx, path, input, nil); err != nil {
		return fmt.Errorf("update organization role group info: %w", err)
	}
	return nil
}

func (c *Client) UpdateOrganizationRoleGroupRoles(ctx context.Context, orgID, roleGroupID string, input *UpdateRoleGroupRolesInput) error {
	path := fmt.Sprintf("/v1/organizations/%s/org-role-groups/%s/roles", orgID, roleGroupID)
	if err := c.httpClient.PUT(ctx, path, input, nil); err != nil {
		return fmt.Errorf("update organization role group roles: %w", err)
	}
	return nil
}

func (c *Client) CreateProject(ctx context.Context, orgID string, input *CreateProjectInput) (*CreateProjectOutput, error) {
	var out CreateProjectOutput
	path := fmt.Sprintf("/v1/organizations/%s/projects", orgID)
	if err := c.httpClient.POST(ctx, path, input, &out); err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteProject(ctx context.Context, projectID string) error {
	path := fmt.Sprintf("/v1/projects/%s", projectID)
	if err := c.httpClient.DELETE(ctx, path, nil); err != nil {
		return fmt.Errorf("delete project %s: %w", projectID, err)
	}
	return nil
}

func (c *Client) ListProjectRoles(ctx context.Context, projectID string) (*ListRolesOutput, error) {
	var out ListRolesOutput
	path := fmt.Sprintf("/v1/projects/%s/roles", projectID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list project roles: %w", err)
	}
	return &out, nil
}

func (c *Client) CreateProjectMember(ctx context.Context, projectID string, input *CreateProjectMemberInput) (*CreateProjectMemberOutput, error) {
	var out CreateProjectMemberOutput
	path := fmt.Sprintf("/v1/projects/%s/members", projectID)
	if err := c.httpClient.POST(ctx, path, input, &out); err != nil {
		return nil, fmt.Errorf("create project member: %w", err)
	}
	return &out, nil
}

func (c *Client) GetProjectMember(ctx context.Context, projectID, memberUUID string) (*GetMemberOutput, error) {
	var out GetMemberOutput
	path := fmt.Sprintf("/v1/projects/%s/members/%s", projectID, memberUUID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get project member %s: %w", memberUUID, err)
	}
	return &out, nil
}

func (c *Client) UpdateProjectMember(ctx context.Context, projectID, memberUUID string, input *UpdateMemberInput) (*UpdateMemberOutput, error) {
	var out UpdateMemberOutput
	path := fmt.Sprintf("/v1/projects/%s/members/%s", projectID, memberUUID)
	if err := c.httpClient.PUT(ctx, path, input, &out); err != nil {
		return nil, fmt.Errorf("update project member %s: %w", memberUUID, err)
	}
	return &out, nil
}

func (c *Client) DeleteProjectMember(ctx context.Context, projectID, memberUUID string) error {
	path := fmt.Sprintf("/v1/projects/%s/members/%s", projectID, memberUUID)
	if err := c.httpClient.DELETE(ctx, path, nil); err != nil {
		return fmt.Errorf("delete project member %s: %w", memberUUID, err)
	}
	return nil
}

func (c *Client) ListProjectRoleGroups(ctx context.Context, projectID string) (*ListRoleGroupsOutput, error) {
	var out ListRoleGroupsOutput
	path := fmt.Sprintf("/v1/projects/%s/project-role-groups", projectID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list project role groups: %w", err)
	}
	return &out, nil
}

func (c *Client) GetProjectRoleGroup(ctx context.Context, projectID, roleGroupID string) (*GetRoleGroupOutput, error) {
	var out GetRoleGroupOutput
	path := fmt.Sprintf("/v1/projects/%s/project-role-groups/%s", projectID, roleGroupID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get project role group %s: %w", roleGroupID, err)
	}
	return &out, nil
}

func (c *Client) CreateProjectRoleGroup(ctx context.Context, projectID string, input *CreateRoleGroupInput) (*CreateRoleGroupOutput, error) {
	var out CreateRoleGroupOutput
	path := fmt.Sprintf("/v1/projects/%s/project-role-groups", projectID)
	if err := c.httpClient.POST(ctx, path, input, &out); err != nil {
		return nil, fmt.Errorf("create project role group: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteProjectRoleGroups(ctx context.Context, projectID string, roleGroupIDs []string) error {
	path := fmt.Sprintf("/v1/projects/%s/project-role-groups", projectID)
	body := map[string][]string{"roleGroupIds": roleGroupIDs}
	if err := c.httpClient.DELETE(ctx, path, body); err != nil {
		return fmt.Errorf("delete project role groups: %w", err)
	}
	return nil
}

func (c *Client) ListUserAccessKeys(ctx context.Context) (*ListUserAccessKeysOutput, error) {
	var out ListUserAccessKeysOutput
	if err := c.httpClient.GET(ctx, "/v1/authentications/user-access-keys", &out); err != nil {
		return nil, fmt.Errorf("list user access keys: %w", err)
	}
	return &out, nil
}

func (c *Client) CreateUserAccessKey(ctx context.Context) (*CreateUserAccessKeyOutput, error) {
	var out CreateUserAccessKeyOutput
	if err := c.httpClient.POST(ctx, "/v1/authentications/user-access-keys", nil, &out); err != nil {
		return nil, fmt.Errorf("create user access key: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateUserAccessKey(ctx context.Context, keyID string, input *UpdateUserAccessKeyInput) error {
	path := fmt.Sprintf("/v1/authentications/user-access-keys/%s", keyID)
	if err := c.httpClient.PUT(ctx, path, input, nil); err != nil {
		return fmt.Errorf("update user access key %s: %w", keyID, err)
	}
	return nil
}

func (c *Client) DeleteUserAccessKey(ctx context.Context, keyID string) error {
	path := fmt.Sprintf("/v1/authentications/user-access-keys/%s", keyID)
	if err := c.httpClient.DELETE(ctx, path, nil); err != nil {
		return fmt.Errorf("delete user access key %s: %w", keyID, err)
	}
	return nil
}

func (c *Client) ReissueSecretKey(ctx context.Context, keyID string) (*CreateUserAccessKeyOutput, error) {
	var out CreateUserAccessKeyOutput
	path := fmt.Sprintf("/v1/authentications/user-access-keys/%s/secretkey-reissue", keyID)
	if err := c.httpClient.PUT(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("reissue secret key %s: %w", keyID, err)
	}
	return &out, nil
}

func (c *Client) ListProjectAppKeys(ctx context.Context, projectID string) (*ListProjectAppKeysOutput, error) {
	var out ListProjectAppKeysOutput
	path := fmt.Sprintf("/v1/authentications/projects/%s/project-appkeys", projectID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list project app keys: %w", err)
	}
	return &out, nil
}

func (c *Client) CreateProjectAppKey(ctx context.Context, projectID string, input *CreateProjectAppKeyInput) (*CreateProjectAppKeyOutput, error) {
	var out CreateProjectAppKeyOutput
	path := fmt.Sprintf("/v1/authentications/projects/%s/project-appkeys", projectID)
	if err := c.httpClient.POST(ctx, path, input, &out); err != nil {
		return nil, fmt.Errorf("create project app key: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteProjectAppKey(ctx context.Context, projectID, appKey string) error {
	path := fmt.Sprintf("/v1/authentications/projects/%s/project-appkeys/%s", projectID, appKey)
	if err := c.httpClient.DELETE(ctx, path, nil); err != nil {
		return fmt.Errorf("delete project app key %s: %w", appKey, err)
	}
	return nil
}

func (c *Client) EnableProjectProduct(ctx context.Context, projectID, productID string) (*EnableProductOutput, error) {
	var out EnableProductOutput
	path := fmt.Sprintf("/v1/projects/%s/products/%s/enable", projectID, productID)
	if err := c.httpClient.POST(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("enable project product %s: %w", productID, err)
	}
	return &out, nil
}

func (c *Client) DisableProjectProduct(ctx context.Context, projectID, productID string) error {
	path := fmt.Sprintf("/v1/projects/%s/products/%s/disable", projectID, productID)
	if err := c.httpClient.DELETE(ctx, path, nil); err != nil {
		return fmt.Errorf("disable project product %s: %w", productID, err)
	}
	return nil
}

func (c *Client) ListOrganizationGovernances(ctx context.Context, orgID string) (*ListGovernancesOutput, error) {
	var out ListGovernancesOutput
	path := fmt.Sprintf("/v1/organizations/%s/governances", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list organization governances: %w", err)
	}
	return &out, nil
}

func (c *Client) ListOrganizationDomains(ctx context.Context, orgID string) (*ListDomainsOutput, error) {
	var out ListDomainsOutput
	path := fmt.Sprintf("/v1/organizations/%s/domains", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list organization domains: %w", err)
	}
	return &out, nil
}

func (c *Client) ListOrganizationIPACL(ctx context.Context, orgID string) (*ListIPACLOutput, error) {
	var out ListIPACLOutput
	path := fmt.Sprintf("/v1/organizations/%s/products/ip-acl", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list organization IP ACL: %w", err)
	}
	return &out, nil
}

func (c *Client) ListProducts(ctx context.Context) (*ListProductsOutput, error) {
	var out ListProductsOutput
	if err := c.httpClient.GET(ctx, "/v1/products", &out); err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	return &out, nil
}

func (c *Client) ListIAMOrganizationMembers(ctx context.Context, orgID string) (*ListIAMMembersOutput, error) {
	var out ListIAMMembersOutput
	path := fmt.Sprintf("/v1/iam/organizations/%s/members", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list IAM organization members: %w", err)
	}
	return &out, nil
}

func (c *Client) GetIAMOrganizationMember(ctx context.Context, orgID, memberUUID string) (*GetIAMMemberOutput, error) {
	var out GetIAMMemberOutput
	path := fmt.Sprintf("/v1/iam/organizations/%s/members/%s", orgID, memberUUID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get IAM organization member %s: %w", memberUUID, err)
	}
	return &out, nil
}

func (c *Client) CreateIAMOrganizationMember(ctx context.Context, orgID string, input *CreateIAMMemberInput) (*CreateIAMMemberOutput, error) {
	var out CreateIAMMemberOutput
	path := fmt.Sprintf("/v1/iam/organizations/%s/members", orgID)
	if err := c.httpClient.POST(ctx, path, input, &out); err != nil {
		return nil, fmt.Errorf("create IAM organization member: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateIAMOrganizationMember(ctx context.Context, orgID, memberUUID string, input *UpdateIAMMemberInput) error {
	path := fmt.Sprintf("/v1/iam/organizations/%s/members/%s", orgID, memberUUID)
	if err := c.httpClient.PUT(ctx, path, input, nil); err != nil {
		return fmt.Errorf("update IAM organization member %s: %w", memberUUID, err)
	}
	return nil
}

func (c *Client) GetIAMSessionSettings(ctx context.Context, orgID string) (*GetIAMSessionSettingsOutput, error) {
	var out GetIAMSessionSettingsOutput
	path := fmt.Sprintf("/v1/iam/organizations/%s/settings/session", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get IAM session settings: %w", err)
	}
	return &out, nil
}

func (c *Client) GetIAMSecurityMFASettings(ctx context.Context, orgID string) (*GetIAMSecurityMFASettingsOutput, error) {
	var out GetIAMSecurityMFASettingsOutput
	path := fmt.Sprintf("/v1/iam/organizations/%s/settings/security-mfa", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get IAM MFA settings: %w", err)
	}
	return &out, nil
}

func (c *Client) GetIAMLoginFailSettings(ctx context.Context, orgID string) (*GetIAMLoginFailSettingsOutput, error) {
	var out GetIAMLoginFailSettingsOutput
	path := fmt.Sprintf("/v1/iam/organizations/%s/settings/security-login-fail", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get IAM login fail settings: %w", err)
	}
	return &out, nil
}

func (c *Client) GetIAMPasswordRule(ctx context.Context, orgID string) (*GetIAMPasswordRuleOutput, error) {
	var out GetIAMPasswordRuleOutput
	path := fmt.Sprintf("/v1/iam/organizations/%s/settings/password-rule", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get IAM password rule: %w", err)
	}
	return &out, nil
}

func (c *Client) ListIAMProjectMembers(ctx context.Context, projectID string) (*ListIAMMembersOutput, error) {
	var out ListIAMMembersOutput
	path := fmt.Sprintf("/v1/iam/projects/%s/members", projectID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list IAM project members: %w", err)
	}
	return &out, nil
}

func (c *Client) GetIAMProjectMember(ctx context.Context, projectID, memberUUID string) (*GetIAMMemberOutput, error) {
	var out GetIAMMemberOutput
	path := fmt.Sprintf("/v1/iam/projects/%s/members/%s", projectID, memberUUID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get IAM project member %s: %w", memberUUID, err)
	}
	return &out, nil
}

func (c *Client) CreateIAMProjectMember(ctx context.Context, projectID string, input *CreateIAMMemberInput) (*CreateIAMMemberOutput, error) {
	var out CreateIAMMemberOutput
	path := fmt.Sprintf("/v1/iam/projects/%s/members", projectID)
	if err := c.httpClient.POST(ctx, path, input, &out); err != nil {
		return nil, fmt.Errorf("create IAM project member: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateIAMProjectMember(ctx context.Context, projectID, memberUUID string, input *UpdateIAMMemberInput) error {
	path := fmt.Sprintf("/v1/iam/projects/%s/members/%s", projectID, memberUUID)
	if err := c.httpClient.PUT(ctx, path, input, nil); err != nil {
		return fmt.Errorf("update IAM project member %s: %w", memberUUID, err)
	}
	return nil
}

func (c *Client) DeleteIAMProjectMembers(ctx context.Context, projectID string, memberUUIDs []string) error {
	path := fmt.Sprintf("/v1/iam/projects/%s/members", projectID)
	body := map[string][]string{"memberUuids": memberUUIDs}
	if err := c.httpClient.DELETE(ctx, path, body); err != nil {
		return fmt.Errorf("delete IAM project members: %w", err)
	}
	return nil
}
