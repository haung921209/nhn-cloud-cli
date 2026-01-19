package cert

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Helper provides utility functions for certificate management in database connections
type Helper struct {
	store *CertificateStore
}

// NewHelper creates a new certificate helper
func NewHelper() (*Helper, error) {
	store, err := NewCertificateStore()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize certificate store: %w", err)
	}

	return &Helper{
		store: store,
	}, nil
}

// GetCertificateForDatabase finds and returns the certificate path for a database connection
// Returns empty string if no certificate is found or needed
func (h *Helper) GetCertificateForDatabase(serviceType, region, instanceID, version string, autoFind bool, explicitPath string, certType string) (string, error) {
	// If explicit path is provided, validate and return it (Only if CA or primary cert)
	if explicitPath != "" && (certType == "" || certType == "CA") {
		if err := h.validateCertificateFile(explicitPath); err != nil {
			return "", fmt.Errorf("explicit certificate path invalid: %w", err)
		}
		return explicitPath, nil
	}

	// If auto-find is disabled, return empty (no certificate)
	if !autoFind {
		return "", nil
	}

	// Try to find certificate in store
	cert, err := h.store.FindCertificateForConnection(serviceType, region, instanceID, version, certType)
	if err != nil {
		// No certificate found is not an error for optional certificates
		return "", nil
	}

	// Validate the found certificate file
	if err := h.validateCertificateFile(cert.FilePath); err != nil {
		return "", fmt.Errorf("stored certificate invalid: %w", err)
	}

	return cert.FilePath, nil
}

// validateCertificateFile checks if a certificate file exists and is readable
func (h *Helper) validateCertificateFile(certPath string) error {
	// Expand path if it contains ~
	if strings.HasPrefix(certPath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		certPath = filepath.Join(homeDir, certPath[1:])
	}

	// Check if file exists and is readable
	info, err := os.Stat(certPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("certificate file does not exist: %s", certPath)
		}
		return fmt.Errorf("cannot access certificate file: %w", err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("certificate path is not a regular file: %s", certPath)
	}

	// Try to read the file to ensure it's accessible
	file, err := os.Open(certPath)
	if err != nil {
		return fmt.Errorf("cannot read certificate file: %w", err)
	}
	file.Close()

	return nil
}

// GenerateConnectionString builds a connection string with SSL/TLS certificate if available
// This is a helper for database CLI tools that need certificate configuration
func (h *Helper) GenerateConnectionString(serviceType, host, port, database, username, password, region, instanceID, version string, autoFind bool, explicitPath string) (string, map[string]string, error) {
	var connectionString string
	params := make(map[string]string)

	// Get certificate path if available
	certPath, err := h.GetCertificateForDatabase(serviceType, region, instanceID, version, autoFind, explicitPath, "CA")
	if err != nil {
		return "", nil, fmt.Errorf("certificate error: %w", err)
	}

	// Build connection string based on service type
	switch serviceType {
	case "mysql":
		connectionString = fmt.Sprintf("mysql://%s:%s@%s:%s/%s", username, password, host, port, database)
		if certPath != "" {
			params["sslmode"] = "required"
			params["sslcert"] = certPath
			params["sslrootcert"] = certPath
			connectionString += "?sslmode=required"
		}

	case "mariadb":
		connectionString = fmt.Sprintf("mysql://%s:%s@%s:%s/%s", username, password, host, port, database)
		if certPath != "" {
			params["sslmode"] = "required"
			params["sslcert"] = certPath
			params["sslrootcert"] = certPath
			connectionString += "?sslmode=required"
		}

	case "postgresql":
		connectionString = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, database)
		if certPath != "" {
			params["sslmode"] = "require"
			params["sslcert"] = certPath
			params["sslrootcert"] = certPath
			connectionString += "?sslmode=require"
		}

	default:
		return "", nil, fmt.Errorf("unsupported service type: %s", serviceType)
	}

	return connectionString, params, nil
}

// GetConnectionCommand generates the appropriate CLI command for connecting to a database
func (h *Helper) GetConnectionCommand(serviceType, host, port, database, username, password, region, instanceID, version string, autoFind bool, explicitPath string) ([]string, error) {
	var cmd []string

	// Get CA Path
	caPath, _ := h.GetCertificateForDatabase(serviceType, region, instanceID, version, autoFind, explicitPath, "CA")

	// Get Client Cert/Key (Auto-find only, explicit path override usually for CA only, or we'd need 3 expl flags)
	clientCertPath, _ := h.GetCertificateForDatabase(serviceType, region, instanceID, version, autoFind, "", "CLIENT-CERT")
	clientKeyPath, _ := h.GetCertificateForDatabase(serviceType, region, instanceID, version, autoFind, "", "CLIENT-KEY")

	// Build command based on service type
	// Build command based on service type
	switch serviceType {
	case "mysql", "mariadb", "rds-mysql", "rds-mariadb":
		cmd = []string{"mysql"}
		cmd = append(cmd, "-h", host)
		cmd = append(cmd, "-P", port)
		cmd = append(cmd, "-u", username)
		cmd = append(cmd, "-p"+password)
		cmd = append(cmd, database)

		if caPath != "" {
			cmd = append(cmd, "--ssl-mode=REQUIRED") // REQUIRED or VERIFY_CA/IDENTITY
			cmd = append(cmd, "--ssl-ca="+caPath)
		}
		if clientCertPath != "" {
			cmd = append(cmd, "--ssl-cert="+clientCertPath)
		}
		if clientKeyPath != "" {
			cmd = append(cmd, "--ssl-key="+clientKeyPath)
		}

	case "postgresql", "rds-postgresql":
		cmd = []string{"psql"}
		connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, database)

		params := []string{}
		if caPath != "" {
			params = append(params, "sslmode=verify-full") // or verify-ca
			params = append(params, "sslrootcert="+caPath)
		}
		if clientCertPath != "" {
			params = append(params, "sslcert="+clientCertPath)
		}
		if clientKeyPath != "" {
			params = append(params, "sslkey="+clientKeyPath)
		}

		if len(params) > 0 {
			connStr += "?" + strings.Join(params, "&")
		}

		cmd = append(cmd, connStr)

	default:
		return nil, fmt.Errorf("unsupported service type: %s", serviceType)
	}

	return cmd, nil
}

// ListCertificatesForService returns certificates available for a specific service type
func (h *Helper) ListCertificatesForService(serviceType, region, instanceID string) ([]CertificateInfo, error) {
	return h.store.ListCertificates(serviceType, region, instanceID)
}

// GetCertificateStore returns the underlying certificate store for advanced operations
func (h *Helper) GetCertificateStore() *CertificateStore {
	return h.store
}
