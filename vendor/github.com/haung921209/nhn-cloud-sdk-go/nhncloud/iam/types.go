package iam

type Organization struct {
	ID          string `json:"orgId"`
	Name        string `json:"orgName"`
	Description string `json:"description,omitempty"`
	Status      string `json:"orgStatusCode"`
	CreatedAt   string `json:"regDateTime,omitempty"`
	UpdatedAt   string `json:"modDateTime,omitempty"`
}

type OrganizationWrapper struct {
	Org Organization `json:"org"`
}

type Project struct {
	ID             string `json:"projectId"`
	Name           string `json:"projectName"`
	Description    string `json:"description,omitempty"`
	OrganizationID string `json:"orgId"`
	Status         string `json:"projectStatusCode"`
	CreatedAt      string `json:"createdDateTime,omitempty"`
	UpdatedAt      string `json:"modifiedDateTime,omitempty"`
}

type Member struct {
	ID        string   `json:"uuid"`
	Email     string   `json:"emailAddress"`
	Name      string   `json:"memberName"`
	Status    string   `json:"memberStatus"`
	Roles     []string `json:"roles,omitempty"`
	CreatedAt string   `json:"createdDateTime,omitempty"`
	UpdatedAt string   `json:"modifiedDateTime,omitempty"`
}

type ListOrganizationsOutput struct {
	OrganizationWrappers []OrganizationWrapper `json:"orgList"`
}

func (o *ListOrganizationsOutput) Organizations() []Organization {
	orgs := make([]Organization, len(o.OrganizationWrappers))
	for i, w := range o.OrganizationWrappers {
		orgs[i] = w.Org
	}
	return orgs
}

type GetOrganizationOutput struct {
	Organization Organization `json:"org"`
}

type ListProjectsOutput struct {
	Projects []Project `json:"projectList"`
}

type GetProjectOutput struct {
	Project Project `json:"project"`
}

type ListMembersOutput struct {
	Members []Member `json:"memberList"`
}

type GetMemberOutput struct {
	Member Member `json:"member"`
}

type InviteMemberInput struct {
	Email string   `json:"emailAddress"`
	Roles []string `json:"roles,omitempty"`
}

type InviteMemberOutput struct {
	MemberID string `json:"uuid"`
}

// UpdateMemberInput defines the input for updating a member's role
type UpdateMemberInput struct {
	AssignRoles []RoleAssignment `json:"assignRoles,omitempty"`
}

// RoleAssignment defines a role assignment with optional conditions
type RoleAssignment struct {
	RoleID     string          `json:"roleId"`
	Conditions []RoleCondition `json:"conditions,omitempty"`
}

// RoleCondition defines a condition for role assignment
type RoleCondition struct {
	AttributeID           string   `json:"attributeId"`
	AttributeOperatorType string   `json:"attributeOperatorTypeCode"`
	AttributeValues       []string `json:"attributeValues"`
}

// UpdateMemberOutput defines the output for updating a member
type UpdateMemberOutput struct {
	Member Member `json:"member"`
}

// Role represents a role in the organization or project
type Role struct {
	RoleID          string `json:"roleId"`
	RoleName        string `json:"roleName"`
	RoleGroupID     string `json:"roleGroupId,omitempty"`
	Description     string `json:"description,omitempty"`
	CategoryKey     string `json:"categoryKey,omitempty"`
	ExposureOrder   int    `json:"exposureOrder,omitempty"`
	RegDateTime     string `json:"regDateTime,omitempty"`
	Assignable      bool   `json:"assignable,omitempty"`
	RoleApplyPolicy string `json:"roleApplyPolicyCode,omitempty"`
}

// ListRolesOutput defines the output for listing roles
type ListRolesOutput struct {
	Roles []Role `json:"roleList"`
}

// RoleGroup represents a role group
type RoleGroup struct {
	RoleGroupID   string `json:"roleGroupId"`
	RoleGroupName string `json:"roleGroupName"`
	Description   string `json:"description,omitempty"`
	RoleGroupType string `json:"roleGroupType,omitempty"`
	CreatedAt     string `json:"regDateTime,omitempty"`
	UpdatedAt     string `json:"modDateTime,omitempty"`
	Roles         []Role `json:"roles,omitempty"`
}

// ListRoleGroupsOutput defines the output for listing role groups
type ListRoleGroupsOutput struct {
	RoleGroups []RoleGroup `json:"roleGroupList"`
}

// GetRoleGroupOutput defines the output for getting a role group
type GetRoleGroupOutput struct {
	RoleGroup RoleGroup `json:"roleGroup"`
}

// CreateRoleGroupInput defines the input for creating a role group
type CreateRoleGroupInput struct {
	RoleGroupName string   `json:"roleGroupName"`
	Description   string   `json:"description,omitempty"`
	RoleIDs       []string `json:"roleIds,omitempty"`
}

// CreateRoleGroupOutput defines the output for creating a role group
type CreateRoleGroupOutput struct {
	RoleGroupID string `json:"roleGroupId"`
}

// UpdateRoleGroupInfoInput defines the input for updating role group info
type UpdateRoleGroupInfoInput struct {
	RoleGroupName string `json:"roleGroupName,omitempty"`
	Description   string `json:"description,omitempty"`
}

// UpdateRoleGroupRolesInput defines the input for updating role group roles
type UpdateRoleGroupRolesInput struct {
	RoleIDs []string `json:"roleIds"`
}

// CreateProjectInput defines the input for creating a project
type CreateProjectInput struct {
	ProjectName string `json:"projectName"`
	Description string `json:"description,omitempty"`
}

// CreateProjectOutput defines the output for creating a project
type CreateProjectOutput struct {
	ProjectID string `json:"projectId"`
}

// CreateProjectMemberInput defines the input for creating a project member
type CreateProjectMemberInput struct {
	MemberUUID  string           `json:"memberUuid,omitempty"`
	Email       string           `json:"email,omitempty"`
	UserCode    string           `json:"userCode,omitempty"`
	AssignRoles []RoleAssignment `json:"assignRoles,omitempty"`
}

// CreateProjectMemberOutput defines the output for creating a project member
type CreateProjectMemberOutput struct {
	MemberUUID string `json:"uuid"`
}

// ListProjectMembersInput defines the input for listing project members
type ListProjectMembersInput struct {
	MemberStatusCodes []string `json:"memberStatusCodes,omitempty"`
	RoleIDs           []string `json:"roleIds,omitempty"`
	Page              int      `json:"paging.page,omitempty"`
	Limit             int      `json:"paging.limit,omitempty"`
}

// ListProjectMembersOutput defines the output for listing project members
type ListProjectMembersOutput struct {
	Members    []Member `json:"memberList"`
	TotalCount int      `json:"totalCount"`
}

// UserAccessKey represents a user access key
type UserAccessKey struct {
	UserAccessKeyID string `json:"userAccessKeyId"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"` // Only returned on create
	Status          string `json:"userAccessKeyStatusCode"`
	CreatedAt       string `json:"regDateTime,omitempty"`
	UpdatedAt       string `json:"modDateTime,omitempty"`
}

// ListUserAccessKeysOutput defines the output for listing user access keys
type ListUserAccessKeysOutput struct {
	UserAccessKeys []UserAccessKey `json:"userAccessKeyList"`
}

// CreateUserAccessKeyOutput defines the output for creating a user access key
type CreateUserAccessKeyOutput struct {
	UserAccessKey UserAccessKey `json:"userAccessKey"`
}

// UpdateUserAccessKeyInput defines the input for updating a user access key
type UpdateUserAccessKeyInput struct {
	Status string `json:"userAccessKeyStatusCode"` // ACTIVE, INACTIVE
}

// ProjectAppKey represents a project app key
type ProjectAppKey struct {
	AppKey      string `json:"appKey"`
	AppKeyName  string `json:"appKeyName,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"regDateTime,omitempty"`
}

// ListProjectAppKeysOutput defines the output for listing project app keys
type ListProjectAppKeysOutput struct {
	ProjectAppKeys []ProjectAppKey `json:"projectAppKeyList"`
}

// CreateProjectAppKeyInput defines the input for creating a project app key
type CreateProjectAppKeyInput struct {
	AppKeyName  string `json:"appKeyName"`
	Description string `json:"description,omitempty"`
}

// CreateProjectAppKeyOutput defines the output for creating a project app key
type CreateProjectAppKeyOutput struct {
	AppKey string `json:"appKey"`
}

// Product represents a service/product
type Product struct {
	ProductID   string `json:"productId"`
	ProductName string `json:"productName"`
	StatusCode  string `json:"statusCode"`
	AppKey      string `json:"appKey,omitempty"`
	SecretKey   string `json:"secretKey,omitempty"`
}

// ListProductsOutput defines the output for listing products
type ListProductsOutput struct {
	Products []Product `json:"productList"`
}

// EnableProductOutput defines the output for enabling a product
type EnableProductOutput struct {
	AppKey        string   `json:"appKey"`
	SecretKey     string   `json:"secretKey,omitempty"`
	ParentProduct *Product `json:"parentProduct,omitempty"`
}

// Governance represents governance settings
type Governance struct {
	GovernanceID   string `json:"governanceId"`
	GovernanceName string `json:"governanceName"`
	Description    string `json:"description,omitempty"`
	Enabled        bool   `json:"enabled"`
}

// ListGovernancesOutput defines the output for listing governances
type ListGovernancesOutput struct {
	Governances []Governance `json:"governanceList"`
}

// Domain represents an organization domain
type Domain struct {
	DomainID   string `json:"domainId"`
	DomainName string `json:"domainName"`
	Status     string `json:"domainStatusCode"`
}

// ListDomainsOutput defines the output for listing domains
type ListDomainsOutput struct {
	Domains []Domain `json:"domainList"`
}

// IPACL represents an IP ACL entry
type IPACL struct {
	ACLID       string `json:"aclId"`
	IPAddress   string `json:"ipAddress"`
	Description string `json:"description,omitempty"`
	TypeCode    string `json:"typeCode"` // ALLOW, DENY
	CreatedAt   string `json:"regDateTime,omitempty"`
}

// ListIPACLOutput defines the output for listing IP ACL
type ListIPACLOutput struct {
	IPACLList []IPACL `json:"ipAclList"`
}

// IAMSessionSettings represents IAM session settings
type IAMSessionSettings struct {
	SessionTimeoutSeconds int  `json:"sessionTimeoutSeconds"`
	MfaEnabled            bool `json:"mfaEnabled"`
}

// GetIAMSessionSettingsOutput defines the output for getting IAM session settings
type GetIAMSessionSettingsOutput struct {
	Settings IAMSessionSettings `json:"sessionSetting"`
}

// IAMSecurityMFASettings represents IAM MFA settings
type IAMSecurityMFASettings struct {
	MfaEnabled   bool   `json:"mfaEnabled"`
	MfaType      string `json:"mfaType,omitempty"`
	ExemptPeriod int    `json:"exemptPeriod,omitempty"`
}

// GetIAMSecurityMFASettingsOutput defines the output for getting IAM MFA settings
type GetIAMSecurityMFASettingsOutput struct {
	Settings IAMSecurityMFASettings `json:"securityMfaSetting"`
}

// IAMLoginFailSettings represents IAM login failure settings
type IAMLoginFailSettings struct {
	MaxLoginFailCount  int `json:"maxLoginFailCount"`
	LockoutDurationSec int `json:"lockoutDurationSeconds"`
}

// GetIAMLoginFailSettingsOutput defines the output for getting IAM login failure settings
type GetIAMLoginFailSettingsOutput struct {
	Settings IAMLoginFailSettings `json:"securityLoginFailSetting"`
}

// IAMPasswordRule represents IAM password policy
type IAMPasswordRule struct {
	MinLength        int  `json:"minLength"`
	MaxLength        int  `json:"maxLength"`
	RequireUppercase bool `json:"requireUppercase"`
	RequireLowercase bool `json:"requireLowercase"`
	RequireDigit     bool `json:"requireDigit"`
	RequireSpecial   bool `json:"requireSpecial"`
	ExpirationDays   int  `json:"expirationDays"`
	HistoryCount     int  `json:"historyCount"`
}

// GetIAMPasswordRuleOutput defines the output for getting IAM password rule
type GetIAMPasswordRuleOutput struct {
	Rule IAMPasswordRule `json:"passwordRule"`
}

// IAMMember represents an IAM account member
type IAMMember struct {
	UUID         string `json:"uuid"`
	MemberID     string `json:"memberId"`
	MemberName   string `json:"memberName"`
	Email        string `json:"emailAddress,omitempty"`
	Status       string `json:"memberStatusCode"`
	CountryCode  string `json:"countryCode,omitempty"`
	LanguageCode string `json:"languageCode,omitempty"`
	Description  string `json:"description,omitempty"`
	CreatedAt    string `json:"regDateTime,omitempty"`
	LastLoginAt  string `json:"lastLoginDateTime,omitempty"`
}

// ListIAMMembersOutput defines the output for listing IAM members
type ListIAMMembersOutput struct {
	Members    []IAMMember `json:"memberList"`
	TotalCount int         `json:"totalCount"`
}

// GetIAMMemberOutput defines the output for getting an IAM member
type GetIAMMemberOutput struct {
	Member IAMMember `json:"member"`
}

// CreateIAMMemberInput defines the input for creating an IAM member
type CreateIAMMemberInput struct {
	MemberID     string           `json:"memberId"`
	MemberName   string           `json:"memberName,omitempty"`
	Email        string           `json:"emailAddress,omitempty"`
	Password     string           `json:"password,omitempty"`
	CountryCode  string           `json:"countryCode,omitempty"`
	LanguageCode string           `json:"languageCode,omitempty"`
	Description  string           `json:"description,omitempty"`
	AssignRoles  []RoleAssignment `json:"assignRoles,omitempty"`
}

// CreateIAMMemberOutput defines the output for creating an IAM member
type CreateIAMMemberOutput struct {
	UUID string `json:"uuid"`
}

// UpdateIAMMemberInput defines the input for updating an IAM member
type UpdateIAMMemberInput struct {
	MemberName   string           `json:"memberName,omitempty"`
	Email        string           `json:"emailAddress,omitempty"`
	CountryCode  string           `json:"countryCode,omitempty"`
	LanguageCode string           `json:"languageCode,omitempty"`
	Description  string           `json:"description,omitempty"`
	AssignRoles  []RoleAssignment `json:"assignRoles,omitempty"`
}
