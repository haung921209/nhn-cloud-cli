package cert

import (
	"time"
)

// CertificateInfo represents metadata about a stored certificate
type CertificateInfo struct {
	ID          string    `json:"id"`
	ServiceType string    `json:"serviceType"` // mysql, mariadb, postgresql
	Region      string    `json:"region"`
	InstanceID  string    `json:"instanceId,omitempty"` // Database instance ID
	Version     string    `json:"version,omitempty"`
	FilePath    string    `json:"filePath"`
	StoredAt    time.Time `json:"storedAt"`
	Source      string    `json:"source"` // manual, downloaded, generated
	Description string    `json:"description,omitempty"`
}

// CertificateStore manages certificate storage and retrieval
type CertificateStore struct {
	BaseDir      string
	MetadataFile string
	Certificates []CertificateInfo
}

// CertificateRequest represents a request to store a certificate
type CertificateRequest struct {
	ServiceType string `json:"serviceType"`
	Region      string `json:"region"`
	InstanceID  string `json:"instanceId,omitempty"` // Database instance ID
	Version     string `json:"version,omitempty"`
	CertData    []byte `json:"certData"`
	Source      string `json:"source"`
	Description string `json:"description,omitempty"`
}

// CertificateListResponse represents the response for listing certificates
type CertificateListResponse struct {
	Certificates []CertificateInfo `json:"certificates"`
	TotalCount   int               `json:"totalCount"`
}
