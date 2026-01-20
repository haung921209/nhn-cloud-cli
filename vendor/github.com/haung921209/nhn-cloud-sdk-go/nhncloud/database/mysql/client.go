// Package mysql provides a client for NHN Cloud RDS for MySQL v3.0 API.
//
// Official API Documentation:
// https://docs.nhncloud.com/ko/Database/RDS%20for%20MySQL/ko/api-guide-v3.0/
//
// Example:
//
//	cfg := mysql.Config{
//	    Region:    "kr1",
//	    AppKey:    "your-app-key",
//	    AccessKey: "your-access-key",
//	    SecretKey: "your-secret-key",
//	}
//
//	client, err := mysql.NewClient(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	instances, err := client.ListInstances(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, inst := range instances.DBInstances {
//	    fmt.Printf("%s: %s\n", inst.DBInstanceName, inst.DBInstanceStatus)
//	}
package mysql

import (
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/auth"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/core"
)

// Client is the MySQL API client
type Client struct {
	core *core.Client
}

// Config holds MySQL client configuration
type Config struct {
	Region    string
	AppKey    string
	AccessKey string
	SecretKey string
}

// NewClient creates a new MySQL client
func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	baseURL := fmt.Sprintf("%s-rds-mysql.api.nhncloudservice.com", cfg.Region)

	authenticator := auth.NewOAuth2Auth(cfg.AppKey, cfg.AccessKey, cfg.SecretKey)

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
		return &core.ValidationError{Field: "AccessKey", Message: "access key is required"}
	}
	if c.SecretKey == "" {
		return &core.ValidationError{Field: "SecretKey", Message: "secret key is required"}
	}
	return nil
}
