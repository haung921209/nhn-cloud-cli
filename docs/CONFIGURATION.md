# Service Configuration Guide

This guide details the **exact requirements** to configure and authenticate each NHN Cloud service supported by the SDK and CLI.

> **Precedence Rule**:
> 1. CLI Flags (`--region`, `--appkey`)
> 2. Environment Variables (`NHN_CLOUD_APPKEY`)
> 3. Config File (`~/.nhncloud/credentials`)

## 1. Global Authentication
Most services connect via the **NHN Cloud API Gateway**. Authentication requires either an **AppKey** or **Identity Token** (Tenant Credentials).

| Credential | Env Variable | Config File Key | Description |
|------------|--------------|-----------------|-------------|
| **Region** | `NHN_CLOUD_REGION` | `region` | `kr1`, `kr2`, `jp1` (Required) |
| **AppKey** | `NHN_CLOUD_APPKEY` | `appkey` | Default AppKey for the project (Required for most) |

---

## 2. Service-Specific Requirements

### A. Compute & Network (Nova/Neutron)
**Services**: `Compute`, `VPC`, `NKS` (Partial), `Cloud Monitoring` (Project Info)
**Requirement**: **Identity Credentials** (Username, Password, TenantID) are mandatory to generate an `X-Auth-Token`.

| Env Variable | Config File Key | Notes |
|--------------|-----------------|-------|
| `NHN_CLOUD_TENANT_ID` | `tenant_id` | Project/Tenant UUID |
| `NHN_CLOUD_USERNAME` | `username` | Email ID (e.g. `user@example.com`) |
| `NHN_CLOUD_PASSWORD` | `api_password` | API Password (Not Console Login PW) |

### B. Database (RDS)
**Services**: `RDS MySQL`, `RDS MariaDB`, `RDS PostgreSQL`
**Requirement**: **AppKey**. Each DB instance type often has a **separate AppKey**.

| Env Variable | Config File Key | Priority |
|--------------|-----------------|----------|
| `NHN_CLOUD_MYSQL_APPKEY` | `mysql_appkey` | Overrides Default AppKey for MySQL |
| `NHN_CLOUD_MARIADB_APPKEY` | `mariadb_appkey` | Overrides Default AppKey for MariaDB |
| `NHN_CLOUD_POSTGRESQL_APPKEY`| `postgresql_appkey`| Overrides Default AppKey for PG |

### C. Object Storage (OBS)
**Services**: `Object Storage`
**Requirement**: **API Password** (Tenant Credentials) for Token Auth **OR** Access Key/Secret Key for S3-Compatible API.
*The CLI uses Token Auth (Swift) by default.*

| Env Variable | Config File Key |
|--------------|-----------------|
| `NHN_CLOUD_TENANT_ID` | `tenant_id` |
| `NHN_CLOUD_USERNAME` | `username` |
| `NHN_CLOUD_PASSWORD` | `api_password` |

### D. NKS (Kubernetes)
**Services**: `NKS`
**Requirement**: Identity Credentials (for API access) + **AppKey** (for some operations).
It uses `NHN_CLOUD_TENANT_ID`, `USERNAME`, `PASSWORD`.

---

## 3. Configuration File Example
Create `~/.nhncloud/credentials`:

```ini
[default]
region = kr1
tenant_id = 3123...
username = email@nhn.com
api_password = secret...
appkey = DefaultAppKey...

# Optional Overrides
mysql_appkey = MySQLSpecific...
mariadb_appkey = MariaDBSpecific...
```

## 4. Database Connection Setup (SSL/TLS)
To connect to the **Data Plane** (SQL connection) of an RDS instance, you **MUST** use the NHN Cloud CA Certificate.

### Step 1: Download CA
[Download Root CA](https://static.toastoven.net/toastcloud/sdk_download/rds/ca-certificate.crt)

### Step 2: Configure Helper (Go)
```go
import (
    "crypto/tls"
    "crypto/x509"
    "database/sql"
    "github.com/go-sql-driver/mysql"
)

func RegisterTLS() {
    rootCertPool := x509.NewCertPool()
    pem, _ := os.ReadFile("/path/to/ca-certificate.crt")
    if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
        panic("Failed to append PEM")
    }
    mysql.RegisterTLSConfig("custom-tls", &tls.Config{RootCAs: rootCertPool})
}

// Connect
db, err := sql.Open("mysql", "user:pass@tcp(host:port)/dbname?tls=custom-tls")
```

### Step 3: CLI Usage
The CLI manages **Control Plane** (Create/Delete Instance). It does not execute SQL queries. Use standard clients (`mysql`, `psql`) with the CA Cert for data access.
