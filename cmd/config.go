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
	Username            string
	APIPassword         string
	TenantID            string
	NKSTenantID         string
	OBSTenantID         string
	RDSAppKey           string
	RDSPostgreSQLAppKey string
	RDSMariaDBAppKey    string
}

var loadedConfig *Config

func LoadConfig() *Config {
	if loadedConfig != nil {
		return loadedConfig
	}

	loadedConfig = &Config{}

	configPath := filepath.Join(os.Getenv("HOME"), ".nhncloud", "credentials")
	file, err := os.Open(configPath)
	if err != nil {
		return loadedConfig
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inDefaultProfile := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			profile := strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			inDefaultProfile = (profile == "default")
			continue
		}

		if !inDefaultProfile {
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
		}
	}

	return loadedConfig
}

func getConfigValue(flagVal, envKey, configVal string) string {
	if flagVal != "" {
		return flagVal
	}
	if envVal := os.Getenv(envKey); envVal != "" {
		return envVal
	}
	return configVal
}
