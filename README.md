# NHN Cloud CLI

Command line interface for NHN Cloud services.

## Installation

```bash
go install github.com/haung921209/nhn-cloud-cli@latest
```

Or build from source:

```bash
git clone https://github.com/haung921209/nhn-cloud-cli.git
cd nhn-cloud-cli
go build -o nhncloud .
```

## Quick Start

### Environment Setup

```bash
export NHN_CLOUD_REGION=kr1
export NHN_CLOUD_APPKEY=your-appkey
export NHN_CLOUD_ACCESS_KEY=your-access-key
export NHN_CLOUD_SECRET_KEY=your-secret-key
```

### Basic Usage

```bash
# List MySQL instances
nhncloud rds-mysql list

# Get instance details
nhncloud rds-mysql get INSTANCE_ID

# List flavors
nhncloud rds-mysql flavors

# List MariaDB instances
nhncloud rds-mariadb list

# List PostgreSQL instances
nhncloud rds-postgresql list

# List PostgreSQL databases
nhncloud rds-postgresql database list INSTANCE_ID
```

## Commands

### RDS MySQL (28 commands)

```bash
nhncloud rds-mysql --help
```

| Command | Description |
|---------|-------------|
| `list` | List all MySQL instances |
| `get` | Get instance details |
| `create` | Create new instance |
| `delete` | Delete instance |
| `modify` | Modify instance |
| `start` | Start instance |
| `stop` | Stop instance |
| `restart` | Restart instance |
| `promote-master` | Promote replica to master |
| `flavors` | List available flavors |
| `versions` | List available versions |
| `backup-list` | List backups |
| `backup-create` | Create backup |
| `backup-delete` | Delete backup |
| `backup-export` | Export backup to Object Storage |
| `restore-point-in-time` | Restore to point in time |
| `security-group-list` | List security groups |
| `security-group-get` | Get security group details |
| `security-group-create` | Create security group |
| `security-group-update` | Update security group |
| `security-group-delete` | Delete security group |
| `parameter-group-list` | List parameter groups |
| `parameter-group-get` | Get parameter group details |
| `parameter-group-create` | Create parameter group |
| `parameter-group-update` | Update parameter group |
| `parameter-group-delete` | Delete parameter group |
| `user-list` | List DB users |
| `schema-list` | List schemas |

### RDS MariaDB (10 commands)

```bash
nhncloud rds-mariadb --help
```

| Command | Description |
|---------|-------------|
| `list` | List all MariaDB instances |
| `get` | Get instance details |
| `create` | Create new instance |
| `delete` | Delete instance |
| `modify` | Modify instance |
| `start` | Start instance |
| `stop` | Stop instance |
| `restart` | Restart instance |
| `flavors` | List available flavors |
| `versions` | List available versions |

### RDS PostgreSQL (13 commands)

```bash
nhncloud rds-postgresql --help
```

| Command | Description |
|---------|-------------|
| `list` | List all PostgreSQL instances |
| `get` | Get instance details |
| `create` | Create new instance |
| `delete` | Delete instance |
| `modify` | Modify instance |
| `start` | Start instance |
| `stop` | Stop instance |
| `restart` | Restart instance |
| `flavors` | List available flavors |
| `versions` | List available versions |
| `database list` | List databases in instance |
| `database create` | Create database |
| `database delete` | Delete database |

## Output Formats

```bash
# Table (default)
nhncloud rds-mysql list

# JSON
nhncloud rds-mysql list -o json

# YAML
nhncloud rds-mysql list -o yaml
```

## Global Flags

| Flag | Environment Variable | Description |
|------|---------------------|-------------|
| `--region` | `NHN_CLOUD_REGION` | NHN Cloud region (kr1, kr2, jp1) |
| `--appkey` | `NHN_CLOUD_APPKEY` | Application key |
| `--debug` | `NHN_CLOUD_DEBUG` | Enable debug output |
| `--output` | - | Output format (table, json, yaml) |

## Examples

### Create MySQL Instance

```bash
nhncloud rds-mysql create \
  --name "my-mysql" \
  --flavor-id "FLAVOR_ID" \
  --version "MYSQL_V8033" \
  --storage-type "General SSD" \
  --storage-size 20 \
  --subnet-id "SUBNET_ID" \
  --username "admin" \
  --password "SecurePassword123!" \
  --parameter-group-id "PARAM_GROUP_ID"
```

### Backup and Restore

```bash
# Create backup
nhncloud rds-mysql backup-create INSTANCE_ID --name "manual-backup"

# List backups
nhncloud rds-mysql backup-list INSTANCE_ID

# Restore to point in time
nhncloud rds-mysql restore-point-in-time INSTANCE_ID \
  --name "restored-db" \
  --restore-time "2024-01-15T10:30:00"
```

### Security Groups

```bash
# List security groups
nhncloud rds-mysql security-group-list

# Create security group
nhncloud rds-mysql security-group-create \
  --name "web-servers" \
  --description "Access from web servers"

# Add rule
nhncloud rds-mysql security-group-update SECURITY_GROUP_ID \
  --add-cidr "10.0.0.0/24"
```

## SDK

This CLI uses the [NHN Cloud SDK for Go](https://github.com/haung921209/nhn-cloud-sdk-go).

## License

Apache License 2.0
