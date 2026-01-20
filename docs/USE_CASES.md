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

### 1. End-to-End Lifecycle (Walkthrough)

This scenario covers creating an instance, configuring network access, connecting via the native CLI shell, and cleaning up.

#### Step 1: Create Instance
Create a MySQL 8.0 instance. Note the `db-instance-identifier` for future steps.

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
  --db-parameter-group-id default-mysql8
```

#### Step 2: Configure Network Access
By default, instances are isolated. You must create a Security Group and allow ingress traffic.

```bash
# 1. Create Security Group
nhncloud rds-mysql create-db-security-group \
  --db-security-group-name mysq-prod-sg \
  --description "Allow Access"

# 2. Authorize Port 3306 (e.g., from anywhere or specific IP)
nhncloud rds-mysql authorize-db-security-group-ingress \
  --db-security-group-identifier mysq-prod-sg \
  --cidr 0.0.0.0/0 \
  --port 3306

# 3. Attach to Instance
nhncloud rds-mysql modify-db-instance \
  --db-instance-identifier my-mysql-prod \
  --db-security-group-ids mysq-prod-sg
```

#### Step 3: Connect & Query
You can connect directly using the CLI's native driver (no external `mysql` client required).

```bash
# Enter Interactive SQL Shell
nhncloud rds-mysql connect \
  --db-instance-identifier my-mysql-prod \
  --username admin \
  --password 'SecurePass123!'

# Or execute a single query
nhncloud rds-mysql connect \
  --db-instance-identifier my-mysql-prod \
  --username admin \
  --password 'SecurePass123!' \
  -e "SELECT VERSION();"
```

#### Step 4: Cleanup
Delete the instance and security group when done.

```bash
# Delete Instance
nhncloud rds-mysql delete-db-instance --db-instance-identifier my-mysql-prod

# Delete Security Group (after instance is deleted)
nhncloud rds-mysql delete-db-security-group --db-security-group-name mysq-prod-sg
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

### 1. End-to-End Lifecycle (Walkthrough)

#### Step 1: Create Instance
Create a MariaDB 10.11 instance.

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

#### Step 2: Configure Network Access
MariaDB security groups require an initial CIDR rule upon creation.

```bash
# 1. Create Security Group with Initial Rule
nhncloud rds-mariadb create-db-security-group \
  --db-security-group-name mariadb-sg \
  --cidr 0.0.0.0/0
  
# (Optional) Add more rules
# nhncloud rds-mariadb authorize-db-security-group-ingress ...

# 2. Attach to Instance
nhncloud rds-mariadb modify-db-instance \
  --db-instance-identifier my-mariadb-prod \
  --db-security-group-ids mariadb-sg
```

#### Step 3: Connect & Query
Use the native CLI shell (fallback to built-in MySQL driver).

```bash
# Connect (Interactive)
nhncloud rds-mariadb connect \
  --db-instance-identifier my-mariadb-prod \
  --username admin \
  --password 'SecurePass123!'

# One-liner Query
nhncloud rds-mariadb connect \
  --db-instance-identifier my-mariadb-prod \
  --username admin \
  --password 'SecurePass123!' \
  -e "SHOW DATABASES;"
```

#### Step 4: Cleanup

```bash
nhncloud rds-mariadb delete-db-instance --db-instance-identifier my-mariadb-prod
nhncloud rds-mariadb delete-db-security-group --db-security-group-name mariadb-sg
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

### 7. User Groups & Monitoring (New)

```bash
# User Groups
nhncloud rds-mysql create-user-group --name dev-team --member-ids user1,user2
nhncloud rds-mysql describe-user-groups
nhncloud rds-mysql delete-user-group --user-group-id ug-xxxxxxxx

# Monitoring
nhncloud rds-mysql describe-metrics
nhncloud rds-mysql get-metric-statistics \
  --db-instance-identifier my-mysql-prod \
  --from 2026-01-16T10:00:00Z \
  --to 2026-01-16T11:00:00Z \
  --interval 60
```

---

## MariaDB Use Cases

...

### 7. User Groups & Monitoring (New)

```bash
# User Groups
nhncloud rds-mariadb create-user-group --name dev-team
nhncloud rds-mariadb describe-user-groups

# Monitoring
nhncloud rds-mariadb describe-metrics
nhncloud rds-mariadb get-metric-statistics \
  --db-instance-identifier my-mariadb-prod \
  --from 2026-01-16T10:00:00Z \
  --to 2026-01-16T11:00:00Z
```

---

## PostgreSQL Use Cases

### 1. End-to-End Lifecycle (Walkthrough)

#### Step 1: Create Instance
Create a PostgreSQL instance (v13+).

```bash
nhncloud rds-postgresql create-db-instance \
  --db-instance-identifier my-pg-prod \
  --db-flavor-id m2.c4m8 \
  --engine-version POSTGRESQL_V13 \
  --master-username admin \
  --master-user-password 'SecurePass123!' \
  --allocated-storage 20 \
  --subnet-id subnet-xxxxxxxx \
  --availability-zone kr-pub-a
```

#### Step 2: Configure Network Access
**Important**: PostgreSQL uses a separate set of security groups from MySQL/MariaDB.

```bash
# 1. List Available PostgreSQL Security Groups
nhncloud rds-postgresql describe-db-security-groups

# 2. Get Details (Check Rules)
nhncloud rds-postgresql get-db-security-group \
  --db-security-group-identifier <pg-sg-id>

# 3. Add Ingress Rule (Port 5432)
nhncloud rds-postgresql authorize-db-security-group-ingress \
  --db-security-group-identifier <pg-sg-id> \
  --cidr 0.0.0.0/0 \
  --port 5432

# 4. Attach Security Group to Instance
nhncloud rds-postgresql modify-db-instance \
  --db-instance-identifier my-pg-prod \
  --db-security-group-ids <pg-sg-id>
```

#### Step 3: Connect & Query
Connect using the built-in native Go driver (supports `psql` fallback).

```bash
# Interactive Shell (REPL)
nhncloud rds-postgresql connect \
  --db-instance-identifier my-pg-prod \
  --username admin \
  --password 'SecurePass123!'

# One-liner Query
nhncloud rds-postgresql connect \
  --db-instance-identifier my-pg-prod \
  --username admin \
  --password 'SecurePass123!' \
  -e "SELECT version();"
```

#### Step 4: Cleanup

```bash

### 4. Compute Access (SSH)

The CLI provides a convenient `connect` command that automatically detects the instance's Public IP and SSH Key.

**Prerequisite (One-time setup per key)**: Import your downloaded key.
```bash
nhncloud compute import-key \
  --key-name my-key-pair \
  --private-key-file ~/Downloads/my-key-pair.pem
```

**Connect**:
```bash
# 1. Standard Connect (Auto-IP & Key Lookup)
nhncloud compute connect --instance-id <instance-id>

# 2. Specify Username (defaults to metadata 'login_username' or 'centos')
nhncloud compute connect --instance-id <instance-id> --username ubuntu

# 3. Explicit Key File
nhncloud compute connect \
  --instance-id <instance-id> \
  --identity-file ~/.ssh/my-special-key.pem
```

### 2. High Availability

```bash
# Enable
nhncloud rds-postgresql enable-multi-az \
  --db-instance-identifier my-pg-prod \
  --ping-interval 10

# Disable
nhncloud rds-postgresql disable-multi-az --db-instance-identifier my-pg-prod
```

### 3. User Groups & Monitoring

```bash
# User Groups
nhncloud rds-postgresql create-user-group --name dev-team
nhncloud rds-postgresql describe-user-groups

# Monitoring
nhncloud rds-postgresql describe-metrics
nhncloud rds-postgresql get-metric-statistics \
  --db-instance-identifier my-pg-prod \
  --from 2026-01-16T10:00:00Z \
  --to 2026-01-16T11:00:00Z
```

---

## Common Workflows

...

## Command Reference Summary

| Operation | MySQL | MariaDB | PostgreSQL |
|-----------|-------|---------|------------|
| List Instances | `describe-db-instances` | `describe-db-instances` | `describe-db-instances` |
| Create Instance | `create-db-instance` | `create-db-instance` | `create-db-instance` |
| Enable HA | `enable-multi-az` | `enable-multi-az` | `enable-multi-az` |
| Create Snapshot | `create-db-snapshot` | `create-db-snapshot` | `create-db-snapshot` |
| Create User | `create-db-user` | `create-db-user` | `create-db-user` |
| Security Group | `create-db-security-group` | `create-db-security-group` | `create-db-security-group` |
| User Groups | `create-user-group` | `create-user-group` | `create-user-group` |
| Monitoring | `get-metric-statistics` | `get-metric-statistics` | `get-metric-statistics` |

---

## Known Limitations

1. **Snapshot Listing**: `describe-db-snapshots` requires `--db-instance-identifier` (API constraint)
2. **HA Pre-requisite**: `enable-multi-az` requires automatic backup to be configured first
3. **MariaDB Security Groups**: Must include initial `--cidr` rule when creating
