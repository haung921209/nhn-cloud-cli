# Configuration and Credentials

This guide explains how to configure authentication and credentials for the NHN Cloud SDK and CLI.

## Authentication Methods

The SDK and CLI support three levels of configuration precedence:

1.  **Command-Line Flags** (CLI only - Highest Priority)
2.  **Environment Variables**
3.  **Configuration File** (`~/.nhncloud/credentials` - Lowest Priority)

## 1. Environment Variables (Recommended)

Set these variables in your shell profile or CI/CD environment.

| Variable | Description | Services |
|----------|-------------|----------|
| `NHN_CLOUD_REGION` | Target Region (`kr1`, `kr2`, `jp1`) | All |
| `NHN_CLOUD_APPKEY` | Default Application AppKey | All |
| `NHN_CLOUD_TENANT_ID` | API Tenant ID | Compute, Network, NKS |
| `NHN_CLOUD_USERNAME` | API Username (Email) | Compute, Network, NKS |
| `NHN_CLOUD_PASSWORD` | API Password | Compute, Network, NKS |
| `NHN_CLOUD_ACCESS_KEY` | Object Storage Access Key | Storage |
| `NHN_CLOUD_SECRET_KEY` | Object Storage Secret Key | Storage |
| `NHN_CLOUD_MARIADB_APPKEY` | Specific AppKey for MariaDB | RDS MariaDB |
| `NHN_CLOUD_MYSQL_APPKEY` | Specific AppKey for MySQL | RDS MySQL |

### Usage Example
```bash
export NHN_CLOUD_REGION="kr1"
export NHN_CLOUD_TENANT_ID="your-tenant-id"
export NHN_CLOUD_USERNAME="your-email@example.com"
export NHN_CLOUD_PASSWORD="your-api-password"
```

## 2. Configuration File

Create a file at `~/.nhncloud/credentials` (Linux/Mac) or `%USERPROFILE%\.nhncloud\credentials` (Windows).

**Format (`INIt` style)**:
```ini
[default]
region = kr1
tenant_id = your-tenant-id
username = your-email@example.com
api_password = your-api-password
appkey = default-app-key
mariadb_appkey = specific-mariadb-key
mysql_appkey = specific-mysql-key
```

## Service Requirements

### Compute & Network (NKS, NCR, VPC)
Requires `tenant_id`, `username`, and `password` to generate an Identity Token (`X-Auth-Token`).

### RDS (MySQL, MariaDB, Postgres) & Cloud Monitoring
Requires only `appkey` validation. However, the SDK unifies auth, so providing Identity Credentials is safe.

### Object Storage
Requires `tenant_id` (used as Account ID) and often Identity Credentials for Token, or Access/Secret keys for S3-compatible API.
