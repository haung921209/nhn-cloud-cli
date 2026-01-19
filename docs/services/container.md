# Container Service Guide

This guide covers **NCS (NHN Container Service)** and **NCR (NHN Container Registry)**.
*> For NKS (Kubernetes), see [NKS Guide](nks.md).*

## 1. NCR (Container Registry)

Private Docker Registry hosting.

### Create Registry
```bash
nhncloud ncr create-registry --name private-repo
```

### List Registries
```bash
nhncloud ncr list-registries
```

### Delete Registry
```bash
nhncloud ncr delete-registry --registry-name private-repo
```

---

## 2. NCS (Container Service)

Run containers without managing clusters.

### Create Workload
Deploy a container image.
```bash
nhncloud ncs create-workload \
  --name my-app \
  --image nginx:latest \
  --container-port 80 \
  --cpu 1 \
  --memory 1Gi \
  --vpc-subnet-id <subnet-uuid>
```

### List Workloads
```bash
nhncloud ncs describe-workloads
```

### Delete Workload
```bash
nhncloud ncs delete-workload --workload-id <uuid>
```
