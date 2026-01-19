# Compute Service Guide

This guide covers the management of virtual machines (instances), keypairs, flavors, and images using the NHN Cloud CLI.

## Prerequisites
- **Authentication**: Usage of Compute services requires **Identity Credentials** (Tenant ID, Username, Password). AppKey alone is **not** sufficient because Compute APIs use OpenStack Nova which relies on Keystone Tokens.
- See [Configuration Guide](../CONFIGURATION.md) for credential setup.

## 1. Instances

### List Instances
List all instances in the current region.
```bash
nhncloud compute list-instances
```

### Describe Instance with Image/Flavor Details
Get detailed information including Public IP, Security Groups, and Status.
```bash
nhncloud compute describe-instances --instance-id <uuid>
```

### Create Instance
Create a new instance.
> **Note**: You need `flavor-id`, `image-id`, `network-id`, and `keypair-name` first.

```bash
nhncloud compute create-instance \
  --name my-web-server \
  --flavor-id <flavor-uuid> \
  --image-id <image-uuid> \
  --network-id <network-uuid> \
  --keypair <keypair-name> \
  --security-groups <sg-name> \
  --user-data ./cloud-init.sh
```

### Delete Instance
```bash
nhncloud compute delete-instance --instance-id <uuid>
```

---

## 2. Keypairs
Keypairs are used for SSH access to Linux instances.

> **Important**: Keypairs created via API (here) are stored in NHN Cloud. Keypairs created via `ssh-keys` command are legacy/local management. Use `compute` commands for NKS and VMs.

### List Keypairs
```bash
nhncloud compute list-key-pairs
```

### Create Keypair
This will generate a new Private Key and save it output (if JSON) or display it.
```bash
nhncloud compute create-key-pair --key-name my-key
```

### Import Keypair
Upload an existing public key (`id_rsa.pub`).
```bash
nhncloud compute create-key-pair --key-name my-imported-key --public-key "$(cat ~/.ssh/id_rsa.pub)"
```

### Delete Keypair
```bash
nhncloud compute delete-key-pair --key-name my-key
```

---

## 3. Flavors & Images

### List Flavors
Find the right CPU/RAM sizing.
```bash
nhncloud compute list-flavors
```

### List Images
Find available OS images (Ubuntu, CentOS, Windows).
```bash
nhncloud compute list-images
```
