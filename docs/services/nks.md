# NKS (Kubernetes) 서비스 가이드 (Service Guide)

이 문서는 NHN Cloud CLI를 사용하여 **NKS (NHN Kubernetes Service)** 클러스터를 관리하는 방법을 상세히 설명합니다.

## 사전 요구사항 (Prerequisites)
- **인증 (Authentication)**: NKS 관리를 위해서는 **Identity 자격 증명** (Tenant ID, Username, Password)이 필요합니다.
- 자세한 내용은 [설정 가이드](../CONFIGURATION.md)를 참고하세요.

---

## 🏗️ 클러스터 생성 (검증된 Golden Config)

NKS 클러스터 생성은 버전과 이미지의 호환성에 민감합니다. 성공적인 생성을 위해 아래의 검증된 설정을 사용하는 것을 권장합니다.

### 1. 리소스 준비 (Prepare Resources)
다음 ID 정보들이 필요합니다. `compute` 및 `network` 명령어를 사용하여 확인할 수 있습니다.

- **Network ID**: VPC ID (예: `6201e913...`)
- **Subnet ID**: 서브넷 ID (예: `dd9c5a60...`)
- **Keypair**: Nova 키페어 (`nhncloud compute create-key-pair`로 생성된 키)
- **Flavor**: `m2.c2m4` (Standard) 사양 권장

### 2. 이미지 및 버전 선택 (Image & Version)

"Invalid Tag" 오류를 방지하기 위해 다음 조합을 사용하세요. (2026년 1월 검증됨)

| 컴포넌트 | 값 | 비고 |
|-----------|-------|-------|
| **Image** | **Ubuntu 22.04 Container** | ID가 `...384281d64e67`로 끝나는 이미지 |
| **Kube Tag** | `v1.31.4` | 정확한 문자열. `+nhn.1` 접미사를 붙이지 마세요. |

> **팁**: `nhncloud nks describe-versions` 명령어로 현재 리전에서 사용 가능한 유효 태그를 확인할 수 있습니다.

### 3. 생성 명령어 실행

```bash
nhncloud nks create-cluster \
  --name my-production-cluster \
  --cluster-template-id iaas_console \
  --network-id <network-uuid> \
  --subnet-id <subnet-uuid> \
  --flavor-id <flavor-uuid> \
  --keypair <keypair-name> \
  --node-count 1 \
  --debug
```

Status가 `202 Accepted`라면 요청이 정상적으로 접수된 것입니다. 프로비저닝에는 약 10~15분이 소요됩니다.

---

## 클러스터 관리 (Manage Cluster)

### 클러스터 목록 조회 (List Clusters)
```bash
nhncloud nks describe-clusters
```

### Kubeconfig 설정 (Get Kubeconfig)
`kubectl`로 클러스터에 접속하기 위한 설정 파일을 다운로드합니다.
```bash
nhncloud nks update-kubeconfig --cluster-id <cluster-uuid> --file ./kubeconfig.yaml
```

### 클러스터 삭제 (Delete Cluster)
```bash
nhncloud nks delete-cluster --cluster-id <cluster-uuid>
```
