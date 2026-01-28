package cmd

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	AccessKeyID         string
	SecretAccessKey     string
	Region              string
	AppKey              string // 공용 앱키 (Compute, NCS 등)
	Username            string
	APIPassword         string
	TenantID            string
	NKSTenantID         string
	OBSTenantID         string
	RDSAppKey           string // RDS MySQL 전용
	RDSPostgreSQLAppKey string // RDS PostgreSQL 전용
	RDSMariaDBAppKey    string // RDS MariaDB 전용
	NCRAppKey           string // NCR 전용
}

var loadedConfig *Config

func LoadConfig() *Config {
	// If config is already loaded and profile hasn't changed, return it.
	// But simple CLI run usually runs once.
	// For now, let's reload if needed or just load once.
	if loadedConfig != nil {
		return loadedConfig
	}

	loadedConfig = &Config{}
	targetProfile := "default"
	if profile != "" {
		targetProfile = profile
	} else if p := os.Getenv("NHN_CLOUD_PROFILE"); p != "" {
		targetProfile = p
	}

	configPath := filepath.Join(os.Getenv("HOME"), ".nhncloud", "credentials")
	file, err := os.Open(configPath)
	if err != nil {
		return loadedConfig
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inTargetProfile := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currProfile := strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			inTargetProfile = (currProfile == targetProfile)
			continue
		}

		if !inTargetProfile {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "access_key_id":
			loadedConfig.AccessKeyID = value
		case "secret_access_key":
			loadedConfig.SecretAccessKey = value
		case "region":
			loadedConfig.Region = value
		case "app_key":
			loadedConfig.AppKey = value
		case "username":
			loadedConfig.Username = value
		case "api_password":
			loadedConfig.APIPassword = value
		case "tenant_id":
			loadedConfig.TenantID = value
		case "nks_tenant_id":
			loadedConfig.NKSTenantID = value
		case "obs_tenant_id":
			loadedConfig.OBSTenantID = value
		case "rds_app_key":
			loadedConfig.RDSAppKey = value
		case "rds_postgresql_app_key":
			loadedConfig.RDSPostgreSQLAppKey = value
		case "rds_mariadb_app_key":
			loadedConfig.RDSMariaDBAppKey = value
		case "ncr_app_key":
			loadedConfig.NCRAppKey = value
		}
	}

	return loadedConfig
}
