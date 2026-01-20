// Package postgresql provides a client for NHN Cloud RDS for PostgreSQL v1.0 API.
//
// Official API Documentation:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20PostgreSQL/ko/api-guide-v1.0/
//
// Key Differences from MySQL/MariaDB:
//   - API Version: v1.0 (not v3.0)
//   - Authentication: Bearer Token (automatically issued via OAuth2)
//   - Database Management: /databases (not /db-schemas)
//   - PostgreSQL-specific: Extensions, HBA Rules
//   - Port: 5432 (not 3306)
//   - Instance creation requires databaseName field
//
// Example:
//
//	cfg := postgresql.Config{
//	    Region:    "kr1",
//	    AppKey:    "your-app-key",
//	    AccessKey: "your-access-key-id",
//	    SecretKey: "your-secret-access-key",
//	}
//
//	client, err := postgresql.NewClient(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Bearer token is automatically issued and cached
//	instances, err := client.ListInstances(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, inst := range instances.DBInstances {
//	    fmt.Printf("%s: %s\n", inst.DBInstanceName, inst.DBInstanceStatus)
//	}
package postgresql

import (
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/auth"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// Client is the PostgreSQL API client
type Client struct {
	core *core.Client
}

// Config holds PostgreSQL client configuration
type Config struct {
	Region    string
	AppKey    string
	AccessKey string // User Access Key ID
	SecretKey string // Secret Access Key
}

// NewClient creates a new PostgreSQL client.
//
// Bearer token is automatically issued using Access Key ID and Secret Access Key,
// and cached locally at ~/.nhncloud/postgresql_token_cache.json
func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// PostgreSQL uses different base URL: rds-postgres (not rds-postgresql)
	baseURL := fmt.Sprintf("%s-rds-postgres.api.nhncloudservice.com", cfg.Region)

	// Use auto-refresh authenticator - token is issued automatically
	authenticator := auth.NewBearerAuthWithAutoRefresh(cfg.AppKey, cfg.AccessKey, cfg.SecretKey)

	coreClient := core.NewClient(baseURL, authenticator, nil)

	return &Client{
		core: coreClient,
	}, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Region == "" {
		return &core.ValidationError{Field: "Region", Message: "region is required"}
	}
	if c.AppKey == "" {
		return &core.ValidationError{Field: "AppKey", Message: "app key is required"}
	}
	if c.AccessKey == "" {
		return &core.ValidationError{Field: "AccessKey", Message: "access key ID is required"}
	}
	if c.SecretKey == "" {
		return &core.ValidationError{Field: "SecretKey", Message: "secret access key is required"}
	}
	return nil
}
