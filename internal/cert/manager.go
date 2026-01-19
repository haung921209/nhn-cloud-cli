package cert

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	DefaultCertDir    = ".nhncloud/certs"
	MetadataFileName  = "certificates.json"
	CertFileExtension = ".pem"
)

// NewCertificateStore creates a new certificate store
func NewCertificateStore() (*CertificateStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	baseDir := filepath.Join(homeDir, DefaultCertDir)
	metadataFile := filepath.Join(baseDir, MetadataFileName)

	store := &CertificateStore{
		BaseDir:      baseDir,
		MetadataFile: metadataFile,
		Certificates: []CertificateInfo{},
	}

	// Ensure directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create certificate directory: %w", err)
	}

	// Load existing metadata if it exists
	if err := store.loadMetadata(); err != nil {
		// If file doesn't exist, that's okay - we'll create it later
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load certificate metadata: %w", err)
		}
	}

	return store, nil
}

// StoreCertificate stores a certificate and updates metadata
func (cs *CertificateStore) StoreCertificate(req *CertificateRequest) (*CertificateInfo, error) {
	// Generate unique ID based on content hash
	hash := sha256.Sum256(req.CertData)
	certID := fmt.Sprintf("%x", hash[:8])

	// Default type to CA if empty
	certType := req.Type
	if certType == "" {
		certType = "CA"
	}
	certType = strings.ToUpper(certType)

	// Create filename with instance ID if provided
	var filename string
	if req.InstanceID != "" {
		if req.Version != "" {
			filename = fmt.Sprintf("%s-%s-%s-%s-%s-%s%s", req.ServiceType, req.Region, req.InstanceID, req.Version, certType, certID, CertFileExtension)
		} else {
			filename = fmt.Sprintf("%s-%s-%s-%s-%s%s", req.ServiceType, req.Region, req.InstanceID, certType, certID, CertFileExtension)
		}
	} else {
		if req.Version != "" {
			filename = fmt.Sprintf("%s-%s-%s-%s-%s%s", req.ServiceType, req.Region, req.Version, certType, certID, CertFileExtension)
		} else {
			filename = fmt.Sprintf("%s-%s-%s-%s%s", req.ServiceType, req.Region, certType, certID, CertFileExtension)
		}
	}

	certPath := filepath.Join(cs.BaseDir, filename)

	// Check if certificate already exists
	for _, cert := range cs.Certificates {
		if cert.ID == certID && cert.Type == certType && cert.InstanceID == req.InstanceID {
			return &cert, fmt.Errorf("certificate already exists with ID: %s", certID)
		}
	}

	// Write certificate file
	if err := os.WriteFile(certPath, req.CertData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write certificate file: %w", err)
	}

	// Create certificate info
	certInfo := CertificateInfo{
		ID:          certID,
		Type:        certType,
		ServiceType: req.ServiceType,
		Region:      req.Region,
		InstanceID:  req.InstanceID,
		Version:     req.Version,
		FilePath:    certPath,
		StoredAt:    time.Now(),
		Source:      req.Source,
		Description: req.Description,
	}

	// Add to store and save metadata
	cs.Certificates = append(cs.Certificates, certInfo)
	if err := cs.saveMetadata(); err != nil {
		// Try to remove the certificate file if metadata save fails
		os.Remove(certPath)
		return nil, fmt.Errorf("failed to save certificate metadata: %w", err)
	}

	return &certInfo, nil
}

// GetCertificate retrieves a certificate by ID
func (cs *CertificateStore) GetCertificate(certID string) (*CertificateInfo, error) {
	for _, cert := range cs.Certificates {
		if cert.ID == certID {
			// Verify file still exists
			if _, err := os.Stat(cert.FilePath); os.IsNotExist(err) {
				return nil, fmt.Errorf("certificate file not found: %s", cert.FilePath)
			}
			return &cert, nil
		}
	}
	return nil, fmt.Errorf("certificate not found: %s", certID)
}

// ListCertificates returns all stored certificates, optionally filtered
func (cs *CertificateStore) ListCertificates(serviceType, region, instanceID string) ([]CertificateInfo, error) {
	var filtered []CertificateInfo

	for _, cert := range cs.Certificates {
		// Apply filters if specified
		if serviceType != "" && cert.ServiceType != serviceType {
			continue
		}
		if region != "" && cert.Region != region {
			continue
		}
		if instanceID != "" && cert.InstanceID != instanceID {
			continue
		}

		// Verify file still exists
		if _, err := os.Stat(cert.FilePath); os.IsNotExist(err) {
			// File is missing, skip this certificate
			continue
		}

		filtered = append(filtered, cert)
	}

	return filtered, nil
}

// RemoveCertificate removes a certificate by ID
func (cs *CertificateStore) RemoveCertificate(certID string) error {
	for i, cert := range cs.Certificates {
		if cert.ID == certID {
			// Remove file
			if err := os.Remove(cert.FilePath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove certificate file: %w", err)
			}

			// Remove from slice
			cs.Certificates = append(cs.Certificates[:i], cs.Certificates[i+1:]...)

			// Save updated metadata
			if err := cs.saveMetadata(); err != nil {
				return fmt.Errorf("failed to save updated metadata: %w", err)
			}

			return nil
		}
	}
	return fmt.Errorf("certificate not found: %s", certID)
}

// FindCertificateForConnection finds the best certificate for a database connection
func (cs *CertificateStore) FindCertificateForConnection(serviceType, region, instanceID, version, certType string) (*CertificateInfo, error) {
	var candidates []CertificateInfo

	if certType == "" {
		certType = "CA"
	}
	certType = strings.ToUpper(certType)

	// Find all certificates matching service type, region, and type
	for _, cert := range cs.Certificates {
		match := (cert.ServiceType == serviceType && cert.Region == region && cert.Type == certType)
		if match {
			// If instance ID is specified, prefer exact instance matches
			if instanceID != "" && cert.InstanceID != "" && cert.InstanceID != instanceID {
				continue
			}
			// Verify file still exists
			if _, err := os.Stat(cert.FilePath); os.IsNotExist(err) {
				continue
			}
			candidates = append(candidates, cert)
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no certificate found for %s in region %s", serviceType, region)
	}

	// Priority 1: Instance-specific certificates with matching version
	if instanceID != "" && version != "" {
		for _, cert := range candidates {
			if cert.InstanceID == instanceID && cert.Version == version {
				return &cert, nil
			}
		}
	}

	// Priority 2: Instance-specific certificates (any version)
	if instanceID != "" {
		for _, cert := range candidates {
			if cert.InstanceID == instanceID {
				return &cert, nil
			}
		}
	}

	// Priority 3: General certificates with matching version
	if version != "" {
		for _, cert := range candidates {
			if cert.Version == version {
				return &cert, nil
			}
		}
	}

	// Priority 4: Most recently stored certificate
	newest := candidates[0]
	for _, cert := range candidates[1:] {
		if cert.StoredAt.After(newest.StoredAt) {
			newest = cert
		}
	}

	return &newest, nil
}

// GetCertificatePath returns the file path for a certificate
func (cs *CertificateStore) GetCertificatePath(certID string) (string, error) {
	cert, err := cs.GetCertificate(certID)
	if err != nil {
		return "", err
	}
	return cert.FilePath, nil
}

// loadMetadata loads certificate metadata from file
func (cs *CertificateStore) loadMetadata() error {
	data, err := os.ReadFile(cs.MetadataFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cs.Certificates)
}

// saveMetadata saves certificate metadata to file
func (cs *CertificateStore) saveMetadata() error {
	data, err := json.MarshalIndent(cs.Certificates, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cs.MetadataFile, data, 0644)
}

// CleanupOrphanedFiles removes certificate files that are not in metadata
func (cs *CertificateStore) CleanupOrphanedFiles() error {
	// Get all .pem files in the directory
	files, err := filepath.Glob(filepath.Join(cs.BaseDir, "*"+CertFileExtension))
	if err != nil {
		return fmt.Errorf("failed to list certificate files: %w", err)
	}

	// Create map of known file paths
	knownFiles := make(map[string]bool)
	for _, cert := range cs.Certificates {
		knownFiles[cert.FilePath] = true
	}

	// Remove orphaned files
	var removed []string
	for _, file := range files {
		if !knownFiles[file] {
			if err := os.Remove(file); err != nil {
				fmt.Printf("Warning: failed to remove orphaned file %s: %v\n", file, err)
			} else {
				removed = append(removed, file)
			}
		}
	}

	if len(removed) > 0 {
		fmt.Printf("Cleaned up %d orphaned certificate files\n", len(removed))
	}

	return nil
}

// ValidateStore checks the integrity of the certificate store
func (cs *CertificateStore) ValidateStore() error {
	var issues []string

	// Check if metadata file is readable
	if _, err := os.Stat(cs.MetadataFile); err != nil {
		if os.IsNotExist(err) {
			// This is okay, we'll create it when needed
		} else {
			issues = append(issues, fmt.Sprintf("metadata file error: %v", err))
		}
	}

	// Check each certificate file
	for i, cert := range cs.Certificates {
		if _, err := os.Stat(cert.FilePath); os.IsNotExist(err) {
			issues = append(issues, fmt.Sprintf("certificate %d (ID: %s) file missing: %s", i+1, cert.ID, cert.FilePath))
		}
	}

	if len(issues) > 0 {
		return fmt.Errorf("certificate store validation failed:\n%s", strings.Join(issues, "\n"))
	}

	return nil
}
