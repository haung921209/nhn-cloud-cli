# RDS (Database) Service Guide

This guide covers usage for **MySQL**, **MariaDB**, and **PostgreSQL**.

## Prerequisites
- **Authentication**: Usage of RDS requires an **AppKey**.
- Each database engine (MySQL, MariaDB, PG) typically has its own AppKey.
- Set them via Env Vars: `NHN_CLOUD_MYSQL_APPKEY`, `NHN_CLOUD_MARIADB_APPKEY`, etc.
- See [Configuration Guide](../CONFIGURATION.md) for details.

---

## 1. Instance Management

### Create Instance (MySQL Example)
```bash
nhncloud rds-mysql create-instance \
  --name my-db-01 \
  --flavor-id <flavor-uuid> \
  --version 8.0 \
  --username admin \
  --password <strong-password> \
  --subnet-id <subnet-uuid> \
  --multi-az false
```
*For MariaDB, use `rds-mariadb`. For PostgreSQL, use `rds-postgresql`.*

### List Instances
```bash
nhncloud rds-mysql describe-db-instances
```

### Delete Instance
```bash
nhncloud rds-mysql delete-instance --instance-id <instance-id>
```

---

## 2. Advanced Configuration

### Parameter Groups
Customize DB settings (buffer pool size, timeouts).
```bash
# List available groups
nhncloud rds-mysql describe-db-parameter-groups

# Create custom group
nhncloud rds-mysql create-db-parameter-group --name my-params --description "Production config"

# Apply to instance (Requires Restart)
nhncloud rds-mysql modify-db-instance --instance-id <id> --parameter-group <group-name>
```

### Security Groups (Access Control)
Strictly control who can access port 3306/5432.
```bash
# Create Group
nhncloud rds-mysql create-db-security-group --name my-db-sg

# Add Rule (Allow Office IP)
nhncloud rds-mysql authorize-db-security-group-ingress \
  --db-security-group-identifier my-db-sg \
  --cidr 203.0.113.5/32

# Apply to Instance
nhncloud rds-mysql modify-db-instance --instance-id <id> --security-groups my-db-sg
```

---

## 3. High Availability (HA)

### Enable Multi-AZ
Promote a Single instance to High Availability (Master + Standby).
```bash
nhncloud rds-mysql enable-multi-az --instance-id <id> --subnet-id <standby-subnet-id>
```

### Read Replicas
Create read-only copies for scaling.
```bash
nhncloud rds-mysql create-read-replica \
  --source-db-instance-identifier <master-id> \
  --name my-replica-01
```

---

## 4. Users and Schemas

### Create DB User
```bash
nhncloud rds-mysql create-db-user \
  --instance-id <id> \
  --username app_user \
  --password <password> \
  --host %
```

### Create Schema (Database)
```bash
nhncloud rds-mysql create-db-schema \
  --instance-id <id> \
  --name production_db \
  --character-set utf8mb4
```
