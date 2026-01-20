package ncr

type ResponseHeader struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

type Registry struct {
	ID        int64  `json:"project_id"`
	Name      string `json:"name"`
	URI       string `json:"uri"`
	IsPublic  bool   `json:"isPublic"`
	Status    string `json:"status"`
	CreatedAt string `json:"creation_time"`
	UpdatedAt string `json:"update_time"`
}

type ListRegistriesOutput struct {
	Header     *ResponseHeader `json:"header"`
	Registries []Registry      `json:"registries"`
}

type GetRegistryOutput struct {
	Header *ResponseHeader `json:"header"`
	Registry
}

type CreateRegistryInput struct {
	Name        string `json:"project_name"`
	Description string `json:"description,omitempty"`
	IsPublic    bool   `json:"isPublic,omitempty"`
}

type CreateRegistryOutput struct {
	Header *ResponseHeader `json:"header"`
	Registry
}

type UpdateRegistryInput struct {
	Description string `json:"description,omitempty"`
	IsPublic    *bool  `json:"isPublic,omitempty"`
}

type Image struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	RegistryID string   `json:"registryId"`
	PullCount  int64    `json:"pullCount"`
	Tags       []string `json:"tags,omitempty"`
	Digest     string   `json:"digest,omitempty"`
	Size       int64    `json:"size"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
}

type ListImagesOutput struct {
	Header *ResponseHeader `json:"header"`
	Images []Image         `json:"images"`
}

type GetImageOutput struct {
	Header *ResponseHeader `json:"header"`
	Image
}

type Tag struct {
	Name         string `json:"name"`
	Digest       string `json:"digest"`
	Size         int64  `json:"size"`
	CreatedAt    string `json:"createdAt"`
	LastPulledAt string `json:"lastPulledAt,omitempty"`
}

type ListTagsOutput struct {
	Header *ResponseHeader `json:"header"`
	Tags   []Tag           `json:"tags"`
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

type Webhook struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	RegistryID string   `json:"registryId"`
	TargetURL  string   `json:"targetUrl"`
	Enabled    bool     `json:"enabled"`
	Events     []string `json:"events"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
}

type ListWebhooksOutput struct {
	Header   *ResponseHeader `json:"header"`
	Webhooks []Webhook       `json:"webhooks"`
}

type CreateWebhookInput struct {
	Name      string   `json:"name"`
	TargetURL string   `json:"targetUrl"`
	Enabled   bool     `json:"enabled,omitempty"`
	Events    []string `json:"events"`
}

type CreateWebhookOutput struct {
	Header *ResponseHeader `json:"header"`
	Webhook
}
