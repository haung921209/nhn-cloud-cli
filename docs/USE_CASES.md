# NHN Cloud CLI - RDS Use Cases

This document provides production-ready use cases for MySQL and MariaDB services using the AWS-style v2.0 CLI.

## Table of Contents

- [Configuration](#configuration)
- [MySQL Use Cases](#mysql-use-cases)
- [MariaDB Use Cases](#mariadb-use-cases)
- [Common Workflows](#common-workflows)

---

## Configuration

### Credentials File

Create `~/.nhncloud/credentials`:

```ini
[default]
region = kr1
access_key_id = YOUR_ACCESS_KEY
secret_access_key = YOUR_SECRET_KEY

# Service-specific AppKeys (required)
rds_mysql_app_key = YOUR_MYSQL_APPKEY
rds_mariadb_app_key = YOUR_MARIADB_APPKEY
rds_postgresql_app_key = YOUR_POSTGRESQL_APPKEY

# Fallback (used if service-specific key not set)
rds_app_key = YOUR_DEFAULT_RDS_APPKEY
```

```bash
chmod 600 ~/.nhncloud/credentials
```

### Environment Variables

```bash
# General
export NHN_REGION=kr1
export NHN_ACCESS_KEY_ID=YOUR_ACCESS_KEY
export NHN_SECRET_ACCESS_KEY=YOUR_SECRET_KEY

# Service-Specific (override)
export NHN_MYSQL_APP_KEY=YOUR_MYSQL_APPKEY
export NHN_MARIADB_APP_KEY=YOUR_MARIADB_APPKEY
```

---

## MySQL Use Cases

### 1. Instance Lifecycle

#### Create Instance

```bash
nhncloud rds-mysql create-db-instance \
  --db-instance-identifier my-mysql-prod \
  --db-flavor-id m2.c4m8 \
  --engine-version MYSQL_V8032 \
  --master-username admin \
  --master-user-password 'SecurePass123!' \
  --allocated-storage 100 \
  --subnet-id subnet-xxxxxxxx \
  --availability-zone kr-pub-a \
  --db-parameter-group-id default-mysql8 \
  --db-security-group-ids sg-001
```

#### List/Describe Instances

```bash
# List all
nhncloud rds-mysql describe-db-instances

# Get specific instance (supports name OR ID)
nhncloud rds-mysql describe-db-instances --db-instance-identifier my-mysql-prod
```

#### Modify Instance

```bash
# Change flavor
nhncloud rds-mysql modify-db-instance \
  --db-instance-identifier my-mysql-prod \
  --db-flavor-id m2.c8m16
```

#### Delete Instance

```bash
nhncloud rds-mysql delete-db-instance --db-instance-identifier my-mysql-prod
```

---

### 2. High Availability (Multi-AZ)

#### Prerequisites
- Automatic backup must be enabled
- `useBackupLock: true` (CLI default)

#### Enable HA

```bash
# 1. Configure backup first
nhncloud rds-mysql modify-db-backup-info \
  --db-instance-identifier my-mysql-prod \
  --backup-retention-period 5

# 2. Enable Multi-AZ
nhncloud rds-mysql enable-multi-az \
  --db-instance-identifier my-mysql-prod \
  --ping-interval 10
```

#### Disable HA

```bash
nhncloud rds-mysql disable-multi-az --db-instance-identifier my-mysql-prod
```

---

### 3. Backups & Snapshots

#### Create Snapshot

```bash
nhncloud rds-mysql create-db-snapshot \
  --db-instance-identifier my-mysql-prod \
  --db-snapshot-identifier prod-snap-20260116
```

#### List Snapshots

```bash
nhncloud rds-mysql describe-db-snapshots \
  --db-instance-identifier my-mysql-prod
```

#### Restore from Snapshot

```bash
nhncloud rds-mysql restore-db-instance-from-db-snapshot \
  --db-snapshot-identifier prod-snap-20260116 \
  --db-instance-identifier my-mysql-restored
```

---

### 4. Security Groups

#### Create Security Group

```bash
nhncloud rds-mysql create-db-security-group \
  --db-security-group-name prod-app-sg \
  --description "Production app servers"
```

#### Add Ingress Rule

```bash
nhncloud rds-mysql authorize-db-security-group-ingress \
  --db-security-group-identifier sg-xxxxxxxx \
  --cidr 10.0.0.0/16 \
  --description "VPC internal"
```

#### List Security Groups

```bash
nhncloud rds-mysql describe-db-security-groups
```

---

### 5. Users & Schemas

#### Create DB User

```bash
nhncloud rds-mysql create-db-user \
  --db-instance-identifier my-mysql-prod \
  --db-user-name app_user \
  --db-password 'AppPass123!' \
  --host '%' \
  --authority-type READ
```

#### Create Schema

```bash
nhncloud rds-mysql create-db-schema \
  --db-instance-identifier my-mysql-prod \
  --db-schema-name app_database
```

---

## MariaDB Use Cases

### 1. Instance Lifecycle

#### Create Instance

```bash
nhncloud rds-mariadb create-db-instance \
  --db-instance-identifier my-mariadb-prod \
  --db-flavor-id m2.c2m4 \
  --engine-version MARIADB_V1011 \
  --master-username admin \
  --master-user-password 'SecurePass123!' \
  --allocated-storage 50 \
  --subnet-id subnet-xxxxxxxx \
  --availability-zone kr-pub-a
```

#### Describe Instances

```bash
nhncloud rds-mariadb describe-db-instances
nhncloud rds-mariadb describe-db-instances --db-instance-identifier my-mariadb-prod
```

---

### 2. High Availability

```bash
# Enable Multi-AZ
nhncloud rds-mariadb enable-multi-az \
  --db-instance-identifier my-mariadb-prod \
  --ping-interval 10

# Disable
nhncloud rds-mariadb disable-multi-az --db-instance-identifier my-mariadb-prod

# Pause/Resume HA Monitoring
nhncloud rds-mariadb pause-multi-az --db-instance-identifier my-mariadb-prod
nhncloud rds-mariadb resume-multi-az --db-instance-identifier my-mariadb-prod
```

---

### 3. Read Replicas

```bash
# Create Read Replica
nhncloud rds-mariadb create-read-replica \
  --db-instance-identifier my-mariadb-prod \
  --replica-identifier my-mariadb-replica-01

# Promote to Standalone
nhncloud rds-mariadb promote-read-replica \
  --db-instance-identifier my-mariadb-replica-01
```

---

### 4. Security Groups (MariaDB-Specific)

> **Note**: MariaDB requires at least one rule when creating a security group.

```bash
# Create with initial rule (required)
nhncloud rds-mariadb create-db-security-group \
  --db-security-group-name mariadb-app-sg \
  --cidr 0.0.0.0/0

# Add additional rules
nhncloud rds-mariadb authorize-db-security-group-ingress \
  --db-security-group-identifier sg-xxxxxxxx \
  --cidr 10.0.0.0/16
```

---

### 5. Parameter Groups

```bash
# List
nhncloud rds-mariadb describe-db-parameter-groups

# Create
nhncloud rds-mariadb create-db-parameter-group \
  --db-parameter-group-name my-custom-params \
  --db-parameter-group-family 10.11

# View with parameters
nhncloud rds-mariadb describe-db-parameter-groups \
  --db-parameter-group-id pg-xxxxxxxx

# Reset to defaults
nhncloud rds-mariadb reset-db-parameter-group \
  --db-parameter-group-id pg-xxxxxxxx
```

---

### 6. Users & Schemas

```bash
# List users
nhncloud rds-mariadb describe-db-users --db-instance-identifier my-mariadb-prod

# Create user
nhncloud rds-mariadb create-db-user \
  --db-instance-identifier my-mariadb-prod \
  --db-user-name app_user \
  --db-password 'AppPass123!' \
  --host '%' \
  --authority-type WRITE

# List schemas
nhncloud rds-mariadb describe-db-schemas --db-instance-identifier my-mariadb-prod

# Create schema
nhncloud rds-mariadb create-db-schema \
  --db-instance-identifier my-mariadb-prod \
  --db-schema-name app_database
```

---

## Common Workflows

### Production Database Setup

```bash
# 1. Create security group
nhncloud rds-mysql create-db-security-group \
  --db-security-group-name prod-sg \
  --description "Production DB access"

# 2. Add allowed CIDRs
nhncloud rds-mysql authorize-db-security-group-ingress \
  --db-security-group-identifier sg-xxx \
  --cidr 10.0.0.0/16

# 3. Create instance
nhncloud rds-mysql create-db-instance \
  --db-instance-identifier prod-db \
  --db-flavor-id m2.c4m8 \
  --engine-version MYSQL_V8032 \
  --master-username admin \
  --master-user-password 'SecurePwd!' \
  --allocated-storage 100 \
  --subnet-id subnet-xxx \
  --db-security-group-ids sg-xxx

# 4. Wait for AVAILABLE status
# (Manual polling or use --wait in future versions)

# 5. Enable backups
nhncloud rds-mysql modify-db-backup-info \
  --db-instance-identifier prod-db \
  --backup-retention-period 7

# 6. Enable HA
nhncloud rds-mysql enable-multi-az \
  --db-instance-identifier prod-db \
  --ping-interval 10

# 7. Create application user
nhncloud rds-mysql create-db-user \
  --db-instance-identifier prod-db \
  --db-user-name app_user \
  --db-password 'AppPwd!' \
  --host '%' \
  --authority-type WRITE
```

### Disaster Recovery Drill

```bash
# 1. Create snapshot
nhncloud rds-mysql create-db-snapshot \
  --db-instance-identifier prod-db \
  --db-snapshot-identifier dr-test-snap

# 2. Restore to new instance
nhncloud rds-mysql restore-db-instance-from-db-snapshot \
  --db-snapshot-identifier dr-test-snap \
  --db-instance-identifier dr-test-instance

# 3. Verify data integrity (manually)

# 4. Clean up
nhncloud rds-mysql delete-db-instance --db-instance-identifier dr-test-instance
```

---

## Command Reference Summary

| Operation | MySQL | MariaDB |
|-----------|-------|---------|
| List Instances | `describe-db-instances` | `describe-db-instances` |
| Create Instance | `create-db-instance` | `create-db-instance` |
| Enable HA | `enable-multi-az` | `enable-multi-az` |
| Create Snapshot | `create-db-snapshot` | `create-db-snapshot` |
| Create User | `create-db-user` | `create-db-user` |
| Create Schema | `create-db-schema` | `create-db-schema` |
| Create Security Group | `create-db-security-group` | `create-db-security-group --cidr` |

---

## Known Limitations

1. **Snapshot Listing**: `describe-db-snapshots` requires `--db-instance-identifier` (API constraint)
2. **HA Pre-requisite**: `enable-multi-az` requires automatic backup to be configured first
3. **MariaDB Security Groups**: Must include initial `--cidr` rule when creating
