# NKS Use Cases (Kubernetes)

This guide covers the lifecycle of NHN Kubernetes Service (NKS) clusters.

## End-to-End Lifecycle

### Step 1: Create Cluster
Creates a control plane. Note that this does not create worker nodes automatically (managed separately via Node Groups).

```bash
nhncloud nks create-cluster \
  --name my-k8s-cluster \
  --network-id <network-id> \
  --subnet-id <subnet-id> \
  --k8s-version v1.31.4 \
  --node-count 0 
```
*Note: `node-count` in `create-cluster` often refers to default node group initialization or is deprecated in favor of explicit Node Groups. For clarity, we recommend creating Node Groups explicitly.*

### Step 2: Create Node Group (Worker Nodes)
After the cluster is ACTIVE, add worker nodes.

```bash
nhncloud nks create-node-group \
  --cluster-id <cluster-id> \
  --name default-worker-group \
  --flavor-id m2.c4m8 \
  --node-count 2
```

### Step 3: Configure `kubectl`
Download the kubeconfig file to access the cluster.

```bash
# Saves to ~/.kube/config by default (or merges)
nhncloud nks update-kubeconfig --cluster-id <cluster-id>
```

### Step 4: Verify Access
```bash
kubectl get nodes
```

### Step 5: Scaling
Resize the node group.

```bash
nhncloud nks update-node-group \
  --cluster-id <cluster-id> \
  --node-group-id <group-id> \
  --node-count 4
```

### Step 6: Cleanup
Delete node groups first, then the cluster.

```bash
# 1. Delete Node Group
nhncloud nks delete-node-group --cluster-id <cluster-id> --node-group-id <group-id>

# 2. Delete Cluster
nhncloud nks delete-cluster --cluster-id <cluster-id>
```

---

## Command Reference

| Operation | Command |
|-----------|---------|
| List Clusters | `nhncloud nks describe-clusters` |
| Get Cluster Detail | `nhncloud nks describe-clusters --cluster-id <id>` |
| Get Kubeconfig | `nhncloud nks update-kubeconfig --cluster-id <id>` |
| List Node Groups | `nhncloud nks describe-node-groups --cluster-id <id>` |
