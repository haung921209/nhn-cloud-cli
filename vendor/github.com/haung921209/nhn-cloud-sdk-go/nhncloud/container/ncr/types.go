package ncr

type ResponseHeader struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// Registry is a registry as the API returns it. ID is the API's project_id,
// which the registryNameOrId path parameter accepts.
type Registry struct {
	ID        int64             `json:"project_id"`
	Name      string            `json:"name"`
	URI       string            `json:"uri"`
	RepoCount int               `json:"repo_count"`
	Usage     int64             `json:"usage"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt string            `json:"creation_time"`
	UpdatedAt string            `json:"update_time"`
}

type ListRegistriesOutput struct {
	Header     *ResponseHeader `json:"header"`
	Registries []Registry      `json:"registries"`
}

type GetRegistryOutput struct {
	Header   *ResponseHeader `json:"header"`
	Registry Registry        `json:"registry"`
}

type CreateRegistryInput struct {
	Name        string `json:"project_name"`
	Description string `json:"description,omitempty"`
	IsPublic    bool   `json:"isPublic,omitempty"`
}

// CreateRegistryOutput is the common header only; the API returns no registry.
type CreateRegistryOutput struct {
	Header *ResponseHeader `json:"header"`
}

type UpdateRegistryInput struct {
	Description string `json:"description,omitempty"`
	IsPublic    *bool  `json:"isPublic,omitempty"`
}

// Image is an image (repository) of a registry. Name has the registry
// prefix ("registry/image"); paths take the image name without it.
type Image struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	ProjectID     int64  `json:"project_id"`
	Description   string `json:"description,omitempty"`
	PullCount     int64  `json:"pull_count"`
	ArtifactCount int64  `json:"artifact_count"`
	CreatedAt     string `json:"creation_time"`
	UpdatedAt     string `json:"update_time"`
}

type ListImagesOutput struct {
	Header *ResponseHeader `json:"header"`
	Images []Image         `json:"images"`
}

type GetImageOutput struct {
	Header *ResponseHeader `json:"header"`
	Image  Image           `json:"image"`
}

// Artifact is one pushed manifest of an image; tags point at artifacts.
type Artifact struct {
	ID       int64  `json:"id"`
	Digest   string `json:"digest"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
	PushTime string `json:"push_time"`
	PullTime string `json:"pull_time"`
	Tags     []Tag  `json:"tags"`
}

type Tag struct {
	ID         int64  `json:"id"`
	ArtifactID int64  `json:"artifact_id"`
	Name       string `json:"name"`
	Immutable  bool   `json:"immutable"`
	PushTime   string `json:"push_time"`
	PullTime   string `json:"pull_time"`
}

type ListArtifactsOutput struct {
	Header    *ResponseHeader `json:"header"`
	Artifacts []Artifact      `json:"artifacts"`
}

type ImageScanResult struct {
	ID              string                `json:"id"`
	ImageID         string                `json:"imageId"`
	Tag             string                `json:"tag"`
	Digest          string                `json:"digest"`
	Status          string                `json:"status"`
	ScanStartedAt   string                `json:"scanStartedAt"`
	ScanCompletedAt string                `json:"scanCompletedAt,omitempty"`
	Vulnerabilities []Vulnerability       `json:"vulnerabilities,omitempty"`
	Summary         *VulnerabilitySummary `json:"summary,omitempty"`
}

type Vulnerability struct {
	ID          string `json:"id"`
	Package     string `json:"package"`
	Version     string `json:"version"`
	FixVersion  string `json:"fixVersion,omitempty"`
	Severity    string `json:"severity"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
}

type VulnerabilitySummary struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Unknown  int `json:"unknown"`
	Total    int `json:"total"`
}

type GetImageScanResultOutput struct {
	Header *ResponseHeader `json:"header"`
	ImageScanResult
}

// Webhook is a webhook policy of a registry.
type Webhook struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`
	ProjectID  int64           `json:"project_id"`
	Enabled    bool            `json:"enabled"`
	EventTypes []string        `json:"event_types"`
	Targets    []WebhookTarget `json:"targets"`
	CreatedAt  string          `json:"creation_time"`
	UpdatedAt  string          `json:"update_time"`
}

type WebhookTarget struct {
	Address        string `json:"address"`
	Type           string `json:"type"`
	AuthHeader     string `json:"auth_header,omitempty"`
	SkipCertVerify bool   `json:"skip_cert_verify,omitempty"`
}

type ListWebhooksOutput struct {
	Header   *ResponseHeader `json:"header"`
	Webhooks []Webhook       `json:"policies"`
}

type CreateWebhookInput struct {
	Name       string          `json:"name,omitempty"`
	Enabled    bool            `json:"enabled"`
	EventTypes []string        `json:"event_types"`
	Targets    []WebhookTarget `json:"targets"`
}

// CreateWebhookOutput is the common header only; the API returns no policy.
type CreateWebhookOutput struct {
	Header *ResponseHeader `json:"header"`
}
