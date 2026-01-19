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

## Configuration

For details on setting up credentials, appkeys, and environment variables, see [Configuration Guide](docs/CONFIGURATION.md).

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
