# NKS 활용 사례 (Kubernetes Use Cases)

이 가이드는 NHN Kubernetes Service (NKS) 클러스터의 전체 라이프사이클을 다룹니다.

## 전체 라이프사이클 (End-to-End Lifecycle)

### 1단계: 클러스터 생성 (Create Cluster)
Control Plane을 생성합니다. 워커 노드는 자동으로 생성되지 않으며, 별도의 노드 그룹(Node Group)으로 관리됩니다.

```bash
nhncloud nks create-cluster \
  --name my-k8s-cluster \
  --network-id <network-id> \
  --subnet-id <subnet-id> \
  --k8s-version v1.31.4 \
  --node-count 0 
```
*참고: `create-cluster`의 `node-count`는 기본 노드 그룹을 초기화하는 데 사용되거나 더 이상 사용되지 않는 경우가 많습니다. 명확한 관리를 위해 노드 그룹을 명시적으로 생성하는 것을 권장합니다.*

### 2단계: 노드 그룹 생성 (Worker Nodes)
클러스터가 ACTIVE 상태가 된 후 워커 노드를 추가합니다.

```bash
nhncloud nks create-node-group \
  --cluster-id <cluster-id> \
  --name default-worker-group \
  --flavor-id m2.c4m8 \
  --node-count 2
```

### 3단계: `kubectl` 설정 (Configure kubectl)
클러스터에 접근하기 위해 kubeconfig 파일을 다운로드합니다.

```bash
# 기본적으로 ~/.kube/config에 저장되거나 병합됩니다.
nhncloud nks update-kubeconfig --cluster-id <cluster-id>
```

### 4단계: 접속 확인 (Verify Access)
```bash
kubectl get nodes
```

### 5단계: 스케일링 (Scaling)
노드 그룹의 크기를 조정합니다.

```bash
nhncloud nks update-node-group \
  --cluster-id <cluster-id> \
  --node-group-id <group-id> \
  --node-count 4
```

### 6단계: 정리 (Cleanup)
반드시 노드 그룹을 먼저 삭제한 후, 클러스터를 삭제해야 합니다.

```bash
# 1. 노드 그룹 삭제
nhncloud nks delete-node-group --cluster-id <cluster-id> --node-group-id <group-id>

# 2. 클러스터 삭제
nhncloud nks delete-cluster --cluster-id <cluster-id>
```

---

## 명령어 참조 (Command Reference)

| 작업 (Operation) | 명령어 (Command) |
|-----------|---------|
| 클러스터 목록 조회 | `nhncloud nks describe-clusters` |
| 클러스터 상세 조회 | `nhncloud nks describe-clusters --cluster-id <id>` |
| Kubeconfig 가져오기 | `nhncloud nks update-kubeconfig --cluster-id <id>` |
| 노드 그룹 목록 조회 | `nhncloud nks describe-node-groups --cluster-id <id>` |
