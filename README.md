# NHN Cloud CLI

A powerful, refactored Command Line Interface for NHN Cloud, built on top of the `nhn-cloud-sdk-go`.

> **Status**: ✅ **Verified** (January 2026)

## Features

- **Unified Interface**: Consistent verbs (`create`, `delete`, `describe`, `list`) across all services.
- **Smart Defaults**: Hardened types for complex inputs (e.g. NKS Labels).
- **Service Coverage**: Compute, Network, Container (NKS/NCR/NCS), Storage (NAS/OBS), Database (RDS).

## Installation

```bash
go build -o nhncloud .
```

## Service Implementation Status

| Service | Capability | Status | Notes |
|---------|------------|--------|-------|
| **NKS** | Kubernetes | 🟢 **Verified** | `create-cluster` (Verified Golden Config), `describe-versions` |
| **NCR** | Registry | 🟢 **Verified** | `create-registry`, `list-registries` |
| **NAS** | Storage | 🟢 **Verified** | `create-volume`, `delete-volume` |
| **OBS** | Storage | 🟢 **Verified** | `put-object`, `get-object` |
| **RDS** | Database | 🟢 **Verified** | MySQL/MariaDB/Postgres Instances |
| **NCS** | Container | 🟢 Logic Verified | `create-workload` (Blocked by Region Availability) |
| **Compute** | VM | 🟢 **Verified** | `list-instances`, Keypairs |
| **Network** | VPC | 🟢 **Verified** | Subnets, VPCs |

## Configuration

Valid credentials are **required** to use the CLI.

### 1. Environment Variables (Recommended)

```bash
export NHN_CLOUD_REGION="kr1"
export NHN_CLOUD_APPKEY="your-app-key"
export NHN_CLOUD_TENANT_ID="your-tenant-id"
export NHN_CLOUD_USERNAME="your-email"
export NHN_CLOUD_PASSWORD="your-password"
```

### 2. Config File (`~/.nhncloud/credentials`)

```ini
[default]
region = kr1
tenant_id = your-tenant-id
username = your-email
api_password = your-pw
appkey = your-app-key
```

### 3. Database SSL/TLS
To connect to RDS instances, use the [NHN Cloud CA Certificate](https://static.toastoven.net/toastcloud/sdk_download/rds/ca-certificate.crt).

## CRUD Usage Guide (Verified Services)

### 1. NKS (Kubernetes Service)

**Create Cluster (Golden Configuration)**
> Requires a specific combination of Image and Kuberentes Version Tag.
> See `nks_setup_guide.md` for details.

```bash
./nhncloud nks create-cluster \
  --name my-cluster \
  --network-id <network-id> \
  --subnet-id <subnet-id> \
  --flavor-id <flavor-id> \
  --keypair <nova-keypair-name> \
  --node-count 1
```

**Verify Creation**
```bash
./nhncloud nks describe-clusters
```

**Delete Cluster**
```bash
./nhncloud nks delete-cluster --cluster-id <cluster-id>
```

---

### 2. NCR (Container Registry)

**Create Registry**
```bash
./nhncloud ncr create-registry --name my-registry
```

**List Registries**
```bash
./nhncloud ncr list-registries
```

**Delete Registry**
```bash
./nhncloud ncr delete-registry --registry-name my-registry
```

---

### 3. NAS (Network Attached Storage)

**Create Volume**
> Minimum size is often 300GB or 500GB depending on the type.

```bash
./nhncloud nas create-volume --name my-volume --size 300
```

**Delete Volume**
```bash
./nhncloud nas delete-volume --volume-id <volume-id>
```

---

### 4. Object Storage (OBS)

**Upload File**
```bash
./nhncloud object-storage put-object --container my-container --object my-file.txt --file ./local-file.txt
```

**Download File**
```bash
./nhncloud object-storage get-object --container my-container --object my-file.txt --output-file ./downloaded.txt
```

---

### 5. RDS (Database)

**Create MySQL Instance**
```bash
./nhncloud rds-mysql create-instance \
  --name my-db \
  --flavor-id <flavor-id> \
  --version 8.0 \
  --username admin \
  --password <secure-password>
```

**Delete Instance**
```bash
./nhncloud rds-mysql delete-instance --instance-id <instance-id>
```
