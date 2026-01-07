# NHN Cloud CLI Use Cases

> **Last Updated**: 2026-01-08  
> **CLI Version**: v0.5.0  
> **Security Focus**: Credential management, audit logging, network isolation

This document provides production-ready use cases with security best practices.

---

## Table of Contents

- [Security Best Practices](#security-best-practices)
- [Credential Management](#credential-management)
- [RDS Database Operations](#rds-database-operations)
- [Compute & Network Operations](#compute--network-operations)
- [Container Operations](#container-operations)
- [Monitoring & Audit](#monitoring--audit)
- [Automation Scripts](#automation-scripts)

---

## Security Best Practices

### 1. Credential Storage Priority

| Method | Security Level | Use Case |
|--------|---------------|----------|
| Config File (`~/.nhncloud/credentials`) | HIGH | Production, persistent setup |
| Environment Variables | MEDIUM | CI/CD pipelines, containers |
| Command-line Flags | LOW | Testing only (avoid in production) |

### 2. File Permissions

```bash
# Secure your credentials file
chmod 600 ~/.nhncloud/credentials
chmod 700 ~/.nhncloud/
```

### 3. Never Do This

```bash
# DANGEROUS: Credentials in command line (visible in process list, shell history)
nhncloud rds-mysql list --access-key "XXXX" --secret-key "YYYY"

# DANGEROUS: Credentials in scripts committed to git
export NHN_CLOUD_SECRET_KEY="hardcoded-secret"
```

### 4. Safe Practices

```bash
# SAFE: Use config file
nhncloud configure  # Interactive setup

# SAFE: Use environment variables from secure source
source ~/.nhncloud/env  # File with 600 permissions

# SAFE: Use credential manager
eval $(nhncloud-credential-helper)
```

---

## Credential Management

### Initial Setup (Recommended)

```bash
# Interactive configuration - credentials saved securely
nhncloud configure
```

This creates `~/.nhncloud/credentials`:

```ini
[default]
region = kr1
access_key_id = ********
secret_access_key = ********
rds_mysql_app_key = ********
rds_mariadb_app_key = ********
rds_postgresql_app_key = ********

# Identity credentials (for Compute/Network)
username = your-email@company.com
api_password = ********
tenant_id = ********
```

### Multiple Profiles

```ini
[default]
region = kr1
access_key_id = prod-access-key
secret_access_key = prod-secret-key

[development]
region = kr1
access_key_id = dev-access-key
secret_access_key = dev-secret-key

[staging]
region = kr2
access_key_id = staging-access-key
secret_access_key = staging-secret-key
```

Usage:

```bash
# Use default profile
nhncloud rds-mysql list

# Use specific profile
NHN_CLOUD_PROFILE=development nhncloud rds-mysql list
```

### CI/CD Environment Variables

```yaml
# GitHub Actions example
env:
  NHN_CLOUD_ACCESS_KEY: ${{ secrets.NHN_CLOUD_ACCESS_KEY }}
  NHN_CLOUD_SECRET_KEY: ${{ secrets.NHN_CLOUD_SECRET_KEY }}
  NHN_CLOUD_MYSQL_APPKEY: ${{ secrets.NHN_CLOUD_MYSQL_APPKEY }}
```

---

## RDS Database Operations

### Use Case 1: Database Inventory Check

Monitor all database instances across services:

```bash
#!/bin/bash
# db-inventory.sh - Daily database inventory report

echo "=== Database Inventory Report $(date) ==="

echo -e "\n--- MySQL Instances ---"
nhncloud rds-mysql list -o json | jq -r '.[] | "\(.name)\t\(.status)\t\(.flavor)"'

echo -e "\n--- MariaDB Instances ---"
nhncloud rds-mariadb list -o json | jq -r '.[] | "\(.name)\t\(.status)\t\(.flavor)"'

echo -e "\n--- PostgreSQL Instances ---"
nhncloud rds-postgresql list -o json | jq -r '.[] | "\(.name)\t\(.status)\t\(.flavor)"'
```

### Use Case 2: Secure Database Creation

Create a MySQL instance with security best practices:

```bash
#!/bin/bash
# create-secure-mysql.sh

# Configuration (use environment variables, not hardcoded)
DB_NAME="${DB_NAME:-my-production-db}"
DB_FLAVOR="${DB_FLAVOR:-m2.c2m4}"
DB_VERSION="${DB_VERSION:-MYSQL_V8033}"
DB_STORAGE_SIZE="${DB_STORAGE_SIZE:-50}"

# Get required IDs
SUBNET_ID=$(nhncloud vpc subnets -o json | jq -r '.[0].id')
PARAM_GROUP_ID=$(nhncloud rds-mysql parameter-group list -o json | jq -r '.[0].id')
FLAVOR_ID=$(nhncloud rds-mysql flavors -o json | jq -r ".[] | select(.name==\"${DB_FLAVOR}\") | .id")

# Generate secure password (don't echo or log)
DB_PASSWORD=$(openssl rand -base64 24 | tr -dc 'a-zA-Z0-9!@#$%' | head -c 20)

# Create instance
nhncloud rds-mysql create \
  --name "$DB_NAME" \
  --flavor-id "$FLAVOR_ID" \
  --version "$DB_VERSION" \
  --storage-type "General SSD" \
  --storage-size "$DB_STORAGE_SIZE" \
  --subnet-id "$SUBNET_ID" \
  --username "admin" \
  --password "$DB_PASSWORD" \
  --parameter-group-id "$PARAM_GROUP_ID"

# Store password securely (example: AWS Secrets Manager, HashiCorp Vault, etc.)
# NEVER echo password to stdout or log files
echo "Database created. Password stored in secure vault."
```

### Use Case 3: Backup Management

```bash
#!/bin/bash
# backup-management.sh

# List all backups
echo "=== Current Backups ==="
nhncloud rds-mysql backup list -o json | jq -r '.[] | "\(.id)\t\(.status)\t\(.created_at)"'

# Create manual backup before maintenance
INSTANCE_ID=$(nhncloud rds-mysql list -o json | jq -r '.[0].id')
nhncloud rds-mysql backup create \
  --instance-id "$INSTANCE_ID" \
  --name "pre-maintenance-$(date +%Y%m%d)"

# Verify backup completed
echo "Waiting for backup to complete..."
sleep 60
nhncloud rds-mysql backup list --instance-id "$INSTANCE_ID" -o json | jq -r '.[-1]'
```

### Use Case 4: Database User Management (Least Privilege)

```bash
#!/bin/bash
# create-db-user.sh - Create user with minimal permissions

INSTANCE_ID="$1"
USERNAME="$2"
DATABASE="$3"

if [ -z "$INSTANCE_ID" ] || [ -z "$USERNAME" ] || [ -z "$DATABASE" ]; then
  echo "Usage: $0 <instance-id> <username> <database>"
  exit 1
fi

# Generate secure password
PASSWORD=$(openssl rand -base64 16)

# Create user with read-only access (example)
nhncloud rds-mysql user-create \
  --instance-id "$INSTANCE_ID" \
  --username "$USERNAME" \
  --password "$PASSWORD" \
  --host "%" \
  --privileges "SELECT" \
  --database "$DATABASE"

echo "User $USERNAME created with read-only access to $DATABASE"
# Store password in secure vault
```

### Use Case 5: High Availability Setup

```bash
#!/bin/bash
# setup-ha.sh - Enable HA for production database

INSTANCE_ID="$1"

if [ -z "$INSTANCE_ID" ]; then
  echo "Usage: $0 <instance-id>"
  exit 1
fi

# Check current status
echo "Current instance status:"
nhncloud rds-mariadb get "$INSTANCE_ID" -o json | jq '{name, status, ha_enabled: .high_availability}'

# Enable HA
echo "Enabling High Availability..."
nhncloud rds-mariadb ha enable --instance-id "$INSTANCE_ID"

# Monitor HA status
echo "Monitoring HA setup (this may take several minutes)..."
while true; do
  STATUS=$(nhncloud rds-mariadb get "$INSTANCE_ID" -o json | jq -r '.status')
  echo "Status: $STATUS"
  if [ "$STATUS" = "AVAILABLE" ]; then
    break
  fi
  sleep 30
done

echo "HA enabled successfully!"
```

---

## Compute & Network Operations

### Use Case 6: Secure VM Deployment

```bash
#!/bin/bash
# deploy-secure-vm.sh

VM_NAME="${VM_NAME:-secure-vm}"
KEY_NAME="${KEY_NAME:-my-keypair}"

# 1. Create dedicated security group
SG_ID=$(nhncloud sg create \
  --name "${VM_NAME}-sg" \
  --description "Security group for ${VM_NAME}" \
  -o json | jq -r '.id')

# 2. Add only necessary rules (principle of least privilege)
# SSH from specific IP only
nhncloud sg rule-create \
  --security-group-id "$SG_ID" \
  --protocol tcp \
  --port-min 22 \
  --port-max 22 \
  --remote-ip "YOUR_OFFICE_IP/32"

# HTTPS for public access
nhncloud sg rule-create \
  --security-group-id "$SG_ID" \
  --protocol tcp \
  --port-min 443 \
  --port-max 443

# 3. Get private subnet (not public)
PRIVATE_SUBNET_ID=$(nhncloud vpc subnets -o json | \
  jq -r '.[] | select(.name | contains("private")) | .id' | head -1)

# 4. Create VM
IMAGE_ID=$(nhncloud compute images -o json | jq -r '.[0].id')
FLAVOR_ID=$(nhncloud compute flavors -o json | jq -r '.[] | select(.name=="m2.c2m4") | .id')

nhncloud compute create \
  --name "$VM_NAME" \
  --image "$IMAGE_ID" \
  --flavor "$FLAVOR_ID" \
  --network "$PRIVATE_SUBNET_ID" \
  --key-name "$KEY_NAME" \
  --security-group "$SG_ID"

echo "VM deployed in private subnet with restricted security group"
```

### Use Case 7: Network Isolation Setup

```bash
#!/bin/bash
# setup-network-isolation.sh

PROJECT_NAME="my-project"

# Create VPC (if not exists)
# Note: VPC creation may require console or API

# List current network topology
echo "=== Network Topology ==="
echo "VPCs:"
nhncloud vpc list

echo -e "\nSubnets:"
nhncloud vpc subnets

echo -e "\nSecurity Groups:"
nhncloud sg list

echo -e "\nFloating IPs:"
nhncloud fip list
```

---

## Container Operations

### Use Case 8: Secure Container Registry

```bash
#!/bin/bash
# container-security.sh

REGISTRY_NAME="my-registry"

# Create private registry
nhncloud ncr create --name "$REGISTRY_NAME" --is-public false

# Scan image for vulnerabilities before deployment
IMAGE_NAME="my-app"
nhncloud ncr scan --registry-id "$REGISTRY_ID" --image "$IMAGE_NAME"

# Check scan results
echo "Vulnerability scan results:"
nhncloud ncr scan-result --registry-id "$REGISTRY_ID" --image "$IMAGE_NAME" -o json | \
  jq '.vulnerabilities | group_by(.severity) | map({severity: .[0].severity, count: length})'
```

### Use Case 9: Kubernetes Cluster with Security

```bash
#!/bin/bash
# create-secure-k8s.sh

CLUSTER_NAME="production-cluster"

# Get required IDs
TEMPLATE_ID=$(nhncloud nks templates -o json | jq -r '.[0].id')
NETWORK_ID=$(nhncloud vpc list -o json | jq -r '.[0].id')
SUBNET_ID=$(nhncloud vpc subnets -o json | jq -r '.[0].id')

# Create cluster
nhncloud nks create \
  --name "$CLUSTER_NAME" \
  --template-id "$TEMPLATE_ID" \
  --network-id "$NETWORK_ID" \
  --subnet-id "$SUBNET_ID" \
  --keypair "my-keypair" \
  --node-count 3

# Wait for cluster to be ready
echo "Waiting for cluster to be ready..."
while true; do
  STATUS=$(nhncloud nks get "$CLUSTER_NAME" -o json | jq -r '.status')
  echo "Cluster status: $STATUS"
  if [ "$STATUS" = "RUNNING" ]; then
    break
  fi
  sleep 60
done

# Get kubeconfig (store securely)
CLUSTER_ID=$(nhncloud nks list -o json | jq -r ".[] | select(.name==\"$CLUSTER_NAME\") | .id")
nhncloud nks kubeconfig "$CLUSTER_ID" > ~/.kube/config
chmod 600 ~/.kube/config

echo "Kubernetes cluster ready. Kubeconfig saved to ~/.kube/config"
```

---

## Monitoring & Audit

### Use Case 10: CloudTrail Audit Logging

```bash
#!/bin/bash
# audit-check.sh - Check recent API activities

# List CloudTrail events (if enabled)
echo "=== Recent API Activities ==="
nhncloud cloudtrail events --limit 50 -o json | \
  jq -r '.[] | "\(.event_time)\t\(.event_name)\t\(.user_name)\t\(.source_ip)"'

# Check for suspicious activities
echo -e "\n=== Security-Related Events ==="
nhncloud cloudtrail events --limit 100 -o json | \
  jq -r '.[] | select(.event_name | test("Delete|Create|Modify")) | "\(.event_time)\t\(.event_name)\t\(.user_name)"'
```

### Use Case 11: Resource Monitoring

```bash
#!/bin/bash
# resource-monitor.sh

echo "=== Resource Summary ==="

echo -e "\n--- Compute Instances ---"
nhncloud compute list -o json | jq -r '.[] | "\(.name)\t\(.status)\t\(.flavor)"'

echo -e "\n--- Database Instances ---"
nhncloud rds-mysql list -o json | jq -r '.[] | "\(.name)\t\(.status)"'
nhncloud rds-mariadb list -o json | jq -r '.[] | "\(.name)\t\(.status)"'
nhncloud rds-postgresql list -o json | jq -r '.[] | "\(.name)\t\(.status)"'

echo -e "\n--- Kubernetes Clusters ---"
nhncloud nks list -o json | jq -r '.[] | "\(.name)\t\(.status)\t\(.node_count) nodes"'

echo -e "\n--- Storage Volumes ---"
nhncloud bs list -o json | jq -r '.[] | "\(.name)\t\(.status)\t\(.size)GB"'
```

---

## Automation Scripts

### Use Case 12: Daily Health Check

```bash
#!/bin/bash
# daily-health-check.sh
# Run via cron: 0 9 * * * /path/to/daily-health-check.sh

LOG_FILE="/var/log/nhncloud-health-$(date +%Y%m%d).log"

{
  echo "=== NHN Cloud Health Check - $(date) ==="
  
  # Check all RDS instances
  echo -e "\n[RDS Health]"
  for db in $(nhncloud rds-mysql list -o json | jq -r '.[].id'); do
    STATUS=$(nhncloud rds-mysql get "$db" -o json | jq -r '.status')
    NAME=$(nhncloud rds-mysql get "$db" -o json | jq -r '.name')
    if [ "$STATUS" != "AVAILABLE" ]; then
      echo "WARNING: $NAME is $STATUS"
    else
      echo "OK: $NAME"
    fi
  done
  
  # Check compute instances
  echo -e "\n[Compute Health]"
  nhncloud compute list -o json | jq -r '.[] | "\(.name): \(.status)"'
  
  # Check NKS clusters
  echo -e "\n[NKS Health]"
  nhncloud nks list -o json | jq -r '.[] | "\(.name): \(.status)"'
  
} | tee "$LOG_FILE"

# Send alert if issues found
if grep -q "WARNING\|ERROR" "$LOG_FILE"; then
  # Send notification (integrate with your alerting system)
  echo "Health check found issues. Check $LOG_FILE"
fi
```

### Use Case 13: Disaster Recovery - Backup Export

```bash
#!/bin/bash
# dr-backup-export.sh
# Export backups to Object Storage for DR

CONTAINER_NAME="db-backups-dr"
DATE=$(date +%Y%m%d)

# Ensure container exists
nhncloud os container-create --name "$CONTAINER_NAME" 2>/dev/null

# Export MySQL backups
for INSTANCE_ID in $(nhncloud rds-mysql list -o json | jq -r '.[].id'); do
  INSTANCE_NAME=$(nhncloud rds-mysql get "$INSTANCE_ID" -o json | jq -r '.name')
  LATEST_BACKUP=$(nhncloud rds-mysql backup list --instance-id "$INSTANCE_ID" -o json | \
    jq -r 'sort_by(.created_at) | .[-1].id')
  
  if [ -n "$LATEST_BACKUP" ]; then
    echo "Exporting backup for $INSTANCE_NAME..."
    nhncloud rds-mysql backup export \
      --backup-id "$LATEST_BACKUP" \
      --container "$CONTAINER_NAME" \
      --object-name "${INSTANCE_NAME}-${DATE}.sql"
  fi
done

echo "DR backup export completed"
```

---

## Security Checklist

Before deploying to production, verify:

- [ ] Credentials stored in `~/.nhncloud/credentials` with `600` permissions
- [ ] No credentials in shell history (`HISTCONTROL=ignorespace`)
- [ ] Security groups follow least privilege principle
- [ ] Databases deployed in private subnets
- [ ] HA enabled for critical databases
- [ ] Regular backup schedule configured
- [ ] CloudTrail enabled for audit logging
- [ ] Vulnerability scanning enabled for container images
- [ ] SSH keys rotated regularly
- [ ] API passwords meet complexity requirements

---

## Related Resources

- [NHN Cloud CLI README](../README.md)
- [NHN Cloud SDK for Go](https://github.com/haung921209/nhn-cloud-sdk-go)
- [NHN Cloud Documentation](https://docs.nhncloud.com/)
