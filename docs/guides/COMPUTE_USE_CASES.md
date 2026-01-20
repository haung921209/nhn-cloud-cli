# Compute & Network Use Cases

This guide covers common scenarios for managing Compute instances and Network resources using the NHN Cloud CLI.

## 1. Compute Access (SSH)

The CLI provides a convenient `connect` command that automatically detects the instance's Public IP and SSH Key.

### Prerequisite: Import Key
One-time setup per key pair. Imports your downloaded PEM file securely.

```bash
nhncloud compute import-key \
  --key-name my-key-pair \
  --private-key-file ~/Downloads/my-key-pair.pem
```

### Connect to Instance

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

---

## 2. Instance Management

### Create Instance
```bash
nhncloud compute create-instance \
  --name my-web-server \
  --image-id <image-id> \
  --flavor-id <flavor-id> \
  --subnet-id <subnet-id> \
  --key-name my-key-pair
```

### List Instances
```bash
nhncloud compute describe-instances
```

### Power Management
```bash
# Start
nhncloud compute start-instances --instance-id <id>

# Stop
nhncloud compute stop-instances --instance-id <id>

# Reboot (Soft)
nhncloud compute reboot-instances --instance-id <id>
```
