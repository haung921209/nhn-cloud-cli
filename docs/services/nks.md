# NKS (Kubernetes) Service Guide

This guide details how to create and manage Kubernetes clusters using the NHN Cloud CLI.

## Prerequisites
- **Authentication**: Requires **Identity Credentials** (Tenant ID, Username, Password).
- see [Configuration Guide](../CONFIGURATION.md).

---

## 🏗️ Create Cluster (Verified Golden Config)

Creating an NKS cluster is sensitive to version and image compatibility. Use the verified configuration below to ensure success.

### 1. Prepare Resources
You need the following IDs. Use `compute` and `network` commands to find them if needed.

- **Network ID**: VPC ID (e.g. `6201e913...`)
- **Subnet ID**: Subnet ID (e.g. `dd9c5a60...`)
- **Keypair**: Must be a Nova Keypair (created via `nhncloud compute create-key-pair`).
- **Flavor**: `m2.c2m4` (Standard) is recommended.

### 2. Choose Image & Version (Crucial)

To avoid "Invalid Tag" errors, use this combination (Verified Jan 2026):

| Component | Value | Notes |
|-----------|-------|-------|
| **Image** | **Ubuntu 22.04 Container** | ID ending in `...384281d64e67` |
| **Kube Tag** | `v1.31.4` | Exact string. Do NOT append `+nhn.1`. |

> **Tip**: Run `nhncloud nks describe-versions` to see valid tags for your region.

### 3. Run Command

```bash
nhncloud nks create-cluster \
  --name my-production-cluster \
  --cluster-template-id iaas_console \
  --network-id <network-uuid> \
  --subnet-id <subnet-uuid> \
  --flavor-id <flavor-uuid> \
  --keypair <keypair-name> \
  --node-count 1 \
  --debug
```

Status `202 Accepted` means the request is valid. Provisioning takes 10-15 minutes.

---

## Manage Cluster

### List Clusters
```bash
nhncloud nks describe-clusters
```

### Get Kubeconfig
Download the configuration file to access your cluster via `kubectl`.
```bash
nhncloud nks update-kubeconfig --cluster-id <cluster-uuid> --file ./kubeconfig.yaml
```

### Delete Cluster
```bash
nhncloud nks delete-cluster --cluster-id <cluster-uuid>
```
