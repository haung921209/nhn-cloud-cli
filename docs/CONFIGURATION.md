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

### A. Compute & Network
Start with Identity Credentials...

### B. Service-Specific Tenants
Some services reside in a different Project (Tenant) than your main Compute resource.

| Service | Env Variable | Config File Key | Description |
|---------|--------------|-----------------|-------------|
| **NKS** | `NHN_CLOUD_NKS_TENANT_ID` | `nks_tenant_id` | Separate Tenant for Kubernetes |
| **Object Storage**| `NHN_CLOUD_OBS_TENANT_ID` | `obs_tenant_id` | Separate Tenant for Storage |

### C. Database (RDS) & Others
**AppKeys** are specific to the Service Instance Type.

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
# Global Defaults
region = kr1
tenant_id = 3123... (Main Compute Tenant)
username = email@nhn.com
api_password = secret...
appkey = DefaultAppKey...

# Service-Specific Overrides
nks_tenant_id = 8f31... (Kubernetes)
obs_tenant_id = cfcb... (Object Storage)

# AppKey Overrides
mysql_appkey = ...
mariadb_appkey = ...
postgresql_appkey = ...
ncr_app_key = ...
```

## 4. Database Connection Setup (SSL/TLS)
To connect to the **Data Plane** (SQL connection) of an RDS instance, you should use the NHN Cloud CA Certificate. The CLI provides built-in tools to manage these certificates and streamline connections.

### CA Certificate Management

#### Import Certificates
You can import CA certificates, Client Certificates, and Client Keys into the CLI's secure store.

```bash
# Import Root CA
nhncloud config ca import \
  --service rds-mysql \
  --region kr1 \
  --file ./ca.pem \
  --description "Root CA"

# Import Client Certificate (for Mutual TLS)
nhncloud config ca import \
  --service rds-mysql \
  --region kr1 \
  --type CLIENT-CERT \
  --instance-id <instance-uuid> \
  --file ./client-cert.pem

# Import Client Key
nhncloud config ca import \
  --service rds-mysql \
  --region kr1 \
  --type CLIENT-KEY \
  --instance-id <instance-uuid> \
  --file ./client-key.pem
```

#### List Certificates
```bash
nhncloud config ca list --service rds-mysql
```

### Automatic Connection Helper
The CLI can launch your local database client (`mysql`, `psql`) with the correct SSL configuration automatically applied.

```bash
# Connect to MySQL (Auto-injects --ssl-ca, --ssl-cert, --ssl-key)
nhncloud rds-mysql connect \
  --db-instance-identifier <instance-uuid> \
  --username <user> \
  --password <pass> \
  --database <db>

# Pass extra arguments (e.g. execute query)
nhncloud rds-mysql connect ... -- -e "SELECT 1;"
```

> **Note**: This requires the `mysql` or `psql` client to be installed in your system PATH.

