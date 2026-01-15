// Package auth provides adapters for converting CLI credentials to SDK v2.0 authentication
package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mariadbsdk "github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mariadb"
	mysqlsdk "github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	postgresqlsdk "github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
)

// loadConfigFile reads credentials from ~/.nhncloud/credentials
func loadConfigFile() (map[string]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".nhncloud", "credentials")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err // File doesn't exist or can't be read
	}

	config := make(map[string]string)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			continue
		}

		// Key-value pair
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				config[key] = value
			}
		}
	}

	return config, nil
}

// GetMySQLConfig creates MySQL SDK config from environment variables or credentials file
func GetMySQLConfig() (mysqlsdk.Config, error) {
	// Load config file
	fileConfig, _ := loadConfigFile() // Ignore error, will fall back to env vars

	// Region
	region := os.Getenv("NHN_REGION")
	if region == "" && fileConfig != nil {
		region = fileConfig["region"]
	}
	if region == "" {
		region = "kr1" // default
	}
	region = strings.ToLower(region)

	// App Key
	appKey := os.Getenv("NHN_MYSQL_APP_KEY")
	if appKey == "" {
		appKey = os.Getenv("NHN_APP_KEY")
	}
	if appKey == "" && fileConfig != nil {
		if val, ok := fileConfig["rds_mysql_app_key"]; ok {
			appKey = val
		} else {
			appKey = fileConfig["rds_app_key"]
		}
	}

	// Access Key ID
	accessKeyID := os.Getenv("NHN_MYSQL_ACCESS_KEY_ID")
	if accessKeyID == "" {
		accessKeyID = os.Getenv("NHN_ACCESS_KEY_ID")
	}
	if accessKeyID == "" && fileConfig != nil {
		accessKeyID = fileConfig["access_key_id"]
	}

	// Secret Access Key
	secretAccessKey := os.Getenv("NHN_MYSQL_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		secretAccessKey = os.Getenv("NHN_SECRET_ACCESS_KEY")
	}
	if secretAccessKey == "" && fileConfig != nil {
		secretAccessKey = fileConfig["secret_access_key"]
	}

	if appKey == "" {
		return mysqlsdk.Config{}, fmt.Errorf("missing app key: set NHN_APP_KEY or rds_app_key in ~/.nhncloud/credentials")
	}

	// For RDS, access key and secret key might be optional (depending on NHN Cloud setup)
	// Set defaults if not provided
	if accessKeyID == "" {
		accessKeyID = "default"
	}
	if secretAccessKey == "" {
		secretAccessKey = "default"
	}

	return mysqlsdk.Config{
		Region:    region,
		AppKey:    appKey,
		AccessKey: accessKeyID,
		SecretKey: secretAccessKey,
	}, nil
}

// GetMariaDBConfig creates MariaDB SDK config from environment variables
func GetMariaDBConfig() (mariadbsdk.Config, error) {
	fileConfig, _ := loadConfigFile()

	region := os.Getenv("NHN_REGION")
	if region == "" && fileConfig != nil {
		region = fileConfig["region"]
	}
	if region == "" {
		region = "kr1"
	}
	region = strings.ToLower(region)

	appKey := os.Getenv("NHN_MARIADB_APP_KEY")
	if appKey == "" {
		appKey = os.Getenv("NHN_APP_KEY")
	}
	if appKey == "" && fileConfig != nil {
		if val, ok := fileConfig["rds_mariadb_app_key"]; ok {
			appKey = val
		} else {
			appKey = fileConfig["rds_app_key"]
		}
	}

	accessKeyID := os.Getenv("NHN_MARIADB_ACCESS_KEY_ID")
	if accessKeyID == "" {
		accessKeyID = os.Getenv("NHN_ACCESS_KEY_ID")
	}
	if accessKeyID == "" && fileConfig != nil {
		accessKeyID = fileConfig["access_key_id"]
	}
	if accessKeyID == "" {
		accessKeyID = "default"
	}

	secretAccessKey := os.Getenv("NHN_MARIADB_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		secretAccessKey = os.Getenv("NHN_SECRET_ACCESS_KEY")
	}
	if secretAccessKey == "" && fileConfig != nil {
		secretAccessKey = fileConfig["secret_access_key"]
	}
	if secretAccessKey == "" {
		secretAccessKey = "default"
	}

	if appKey == "" {
		return mariadbsdk.Config{}, fmt.Errorf("missing app key: set NHN_APP_KEY or rds_app_key in ~/.nhncloud/credentials")
	}

	return mariadbsdk.Config{
		Region:    region,
		AppKey:    appKey,
		AccessKey: accessKeyID,
		SecretKey: secretAccessKey,
	}, nil
}

// GetPostgreSQLConfig creates PostgreSQL SDK config from environment variables
// Token is automatically issued using AccessKey and SecretKey
func GetPostgreSQLConfig() (postgresqlsdk.Config, error) {
	fileConfig, _ := loadConfigFile()

	region := os.Getenv("NHN_REGION")
	if region == "" && fileConfig != nil {
		region = fileConfig["region"]
	}
	if region == "" {
		region = "kr1"
	}
	region = strings.ToLower(region)

	appKey := os.Getenv("NHN_POSTGRESQL_APP_KEY")
	if appKey == "" {
		appKey = os.Getenv("NHN_APP_KEY")
	}
	if appKey == "" && fileConfig != nil {
		if val, ok := fileConfig["rds_postgresql_app_key"]; ok {
			appKey = val
		} else {
			appKey = fileConfig["rds_app_key"]
		}
	}

	// Use AccessKey and SecretKey for automatic token issuance
	accessKey := os.Getenv("NHN_POSTGRESQL_ACCESS_KEY_ID")
	if accessKey == "" {
		accessKey = os.Getenv("NHN_ACCESS_KEY_ID")
	}
	if accessKey == "" && fileConfig != nil {
		accessKey = fileConfig["access_key_id"]
	}

	secretKey := os.Getenv("NHN_POSTGRESQL_SECRET_ACCESS_KEY")
	if secretKey == "" {
		secretKey = os.Getenv("NHN_SECRET_ACCESS_KEY")
	}
	if secretKey == "" && fileConfig != nil {
		secretKey = fileConfig["secret_access_key"]
	}

	if appKey == "" {
		return postgresqlsdk.Config{}, fmt.Errorf("missing app key: set NHN_APP_KEY or rds_postgresql_app_key in ~/.nhncloud/credentials")
	}
	if accessKey == "" {
		return postgresqlsdk.Config{}, fmt.Errorf("missing access key: set NHN_ACCESS_KEY_ID or access_key_id in ~/.nhncloud/credentials")
	}
	if secretKey == "" {
		return postgresqlsdk.Config{}, fmt.Errorf("missing secret key: set NHN_SECRET_ACCESS_KEY or secret_access_key in ~/.nhncloud/credentials")
	}

	return postgresqlsdk.Config{
		Region:    region,
		AppKey:    appKey,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, nil
}
