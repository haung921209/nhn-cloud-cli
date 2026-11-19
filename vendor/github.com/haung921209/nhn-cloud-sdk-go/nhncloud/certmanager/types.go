// Package certmanager provides Certificate Manager service types and client
package certmanager

// Header represents the response header
type Header struct {
	IsSuccessful  bool   `json:"isSuccessful"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// Certificate represents a certificate in the list (API v1.3)
type Certificate struct {
	CertificateID      int64    `json:"certificateId"`
	CertificateName    string   `json:"certificateName"`
	Authority          string   `json:"authority"`
	Domains            []string `json:"domains"`
	SignatureAlgorithm string   `json:"signatureAlgorithm"`
	FileCreationDate   string   `json:"fileCreationDate"` // yyyy-mm-dd
	ExpirationDate     string   `json:"expirationDate"`   // yyyy-mm-dd
}

// ListCertificatesOutput represents the response from listing certificates
type ListCertificatesOutput struct {
	Header Header `json:"header"`
	Body   struct {
		TotalCount   int           `json:"totalCount"`
		Certificates []Certificate `json:"data"`
	} `json:"body"`
}

// CertificateFiles represents certificate file data
type CertificateFiles struct {
	Certificate      string `json:"certificate"`      // PEM format certificate
	PrivateKey       string `json:"privateKey"`       // PEM format private key (if available)
	CertificateChain string `json:"certificateChain"` // PEM format certificate chain
}

// DownloadCertificateFilesOutput represents the raw PEM binary response from downloading certificate files
type DownloadCertificateFilesOutput struct {
	Data []byte
}
