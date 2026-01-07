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
# For RDS services (OAuth)
export NHN_CLOUD_REGION=kr1
export NHN_CLOUD_APPKEY=your-appkey
export NHN_CLOUD_ACCESS_KEY=your-access-key
export NHN_CLOUD_SECRET_KEY=your-secret-key

# For Compute/Network services (Identity)
export NHN_CLOUD_USERNAME=your-email
export NHN_CLOUD_PASSWORD=your-api-password
export NHN_CLOUD_TENANT_ID=your-tenant-id
```

### Basic Usage

```bash
# Compute
nhncloud compute list                    # List VM instances
nhncloud compute create --name my-vm --image IMG_ID --flavor FLV_ID --network NET_ID
nhncloud compute flavors                 # List available flavors
nhncloud compute images                  # List available images

# Network
nhncloud vpc list                        # List VPCs
nhncloud vpc subnets                     # List subnets
nhncloud security-group list             # List security groups
nhncloud floating-ip list                # List floating IPs

# RDS MySQL
nhncloud rds-mysql list                  # List MySQL instances
nhncloud rds-mysql flavors               # List available flavors
nhncloud rds-mysql create --name my-db --flavor-id FLV_ID ...
```

## Commands

### Compute (12 commands)

```bash
nhncloud compute --help
```

| Command | Description |
|---------|-------------|
| `list` | List all VM instances |
| `get` | Get instance details |
| `create` | Create new instance |
| `delete` | Delete instance |
| `start` | Start instance |
| `stop` | Stop instance |
| `reboot` | Reboot instance (--hard for hard reboot) |
| `flavors` | List available flavors |
| `images` | List available images |
| `keypairs` | List SSH keypairs |
| `keypair-create` | Create new SSH keypair |
| `keypair-delete` | Delete SSH keypair |

### VPC (3 commands)

```bash
nhncloud vpc --help
```

| Command | Description |
|---------|-------------|
| `list` | List all VPCs |
| `get` | Get VPC details |
| `subnets` | List subnets |

### Security Group (5 commands)

```bash
nhncloud security-group --help
nhncloud sg --help  # alias
```

| Command | Description |
|---------|-------------|
| `list` | List all security groups |
| `get` | Get security group details with rules |
| `create` | Create new security group |
| `delete` | Delete security group |
| `rule-create` | Add rule to security group |

### Floating IP (6 commands)

```bash
nhncloud floating-ip --help
nhncloud fip --help  # alias
```

| Command | Description |
|---------|-------------|
| `list` | List all floating IPs |
| `get` | Get floating IP details |
| `create` | Create new floating IP |
| `delete` | Delete floating IP |
| `associate` | Associate floating IP with port |
| `disassociate` | Disassociate floating IP |

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
| `flavors` | List available flavors |
| `versions` | List available versions |
| `backup` | Backup management (list/create/delete/export) |
| `security-group` | DB security group management |
| `parameter-group` | Parameter group management |
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
| `database list` | List databases |
| `database create` | Create database |
| `database delete` | Delete database |

## Output Formats

```bash
# Table (default)
nhncloud compute list

# JSON
nhncloud compute list -o json

# YAML
nhncloud compute list -o yaml
```

## Global Flags

| Flag | Environment Variable | Description |
|------|---------------------|-------------|
| `--region` | `NHN_CLOUD_REGION` | NHN Cloud region (kr1, kr2, jp1) |
| `--appkey` | `NHN_CLOUD_APPKEY` | Application key (RDS) |
| `--username` | `NHN_CLOUD_USERNAME` | API username (Compute/Network) |
| `--password` | `NHN_CLOUD_PASSWORD` | API password (Compute/Network) |
| `--tenant-id` | `NHN_CLOUD_TENANT_ID` | Tenant ID (Compute/Network) |
| `--debug` | - | Enable debug output |
| `--output` | - | Output format (table, json, yaml) |

## Examples

### Create VM with Floating IP

```bash
# Create security group
nhncloud sg create --name my-sg --description "My security group"

# Add SSH rule
nhncloud sg rule-create --security-group-id SG_ID \
  --protocol tcp --port-min 22 --port-max 22

# Create VM
nhncloud compute create \
  --name my-vm \
  --image IMAGE_ID \
  --flavor FLAVOR_ID \
  --network SUBNET_ID \
  --key-name my-keypair \
  --security-group my-sg

# Associate floating IP
nhncloud fip associate FIP_ID --port-id PORT_ID
```

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

## E2E Test

Run the integration test script:

```bash
# Set environment variables first
./scripts/e2e-test.sh
```

## SDK

This CLI uses the [NHN Cloud SDK for Go](https://github.com/haung921209/nhn-cloud-sdk-go).

## License

Apache License 2.0
