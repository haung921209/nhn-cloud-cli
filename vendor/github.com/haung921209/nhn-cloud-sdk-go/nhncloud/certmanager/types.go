// Package certmanager provides Certificate Manager service types and client
package certmanager

import "time"

// Header represents the response header
type Header struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// Certificate represents a certificate
type Certificate struct {
	CertificateName         string    `json:"certificateName"`
	CertificateType         string    `json:"certificateType"` // SINGLE_DOMAIN, WILDCARD, MULTI_DOMAIN
	Status                  string    `json:"status"`          // ACTIVE, EXPIRED, REVOKED
	DomainName              string    `json:"domainName"`
	SubjectAlternativeNames []string  `json:"subjectAlternativeNames,omitempty"`
	Issuer                  string    `json:"issuer"`
	SerialNumber            string    `json:"serialNumber"`
	NotBefore               time.Time `json:"notBefore"`
	NotAfter                time.Time `json:"notAfter"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdatedAt               time.Time `json:"updatedAt"`
	KeyAlgorithm            string    `json:"keyAlgorithm,omitempty"` // RSA, ECDSA
	KeySize                 int       `json:"keySize,omitempty"`
	SignatureAlgorithm      string    `json:"signatureAlgorithm,omitempty"`
}

// ListCertificatesOutput represents the response from listing certificates
type ListCertificatesOutput struct {
	Header Header `json:"header"`
	Body   struct {
		TotalCount   int           `json:"totalCount"`
		Certificates []Certificate `json:"certificates"`
	} `json:"body"`
}

// CertificateFiles represents certificate file data
type CertificateFiles struct {
	Certificate      string `json:"certificate"`      // PEM format certificate
	PrivateKey       string `json:"privateKey"`       // PEM format private key (if available)
	CertificateChain string `json:"certificateChain"` // PEM format certificate chain
}

// DownloadCertificateFilesOutput represents the response from downloading certificate files
type DownloadCertificateFilesOutput struct {
	Header Header           `json:"header"`
	Body   CertificateFiles `json:"body"`
}
