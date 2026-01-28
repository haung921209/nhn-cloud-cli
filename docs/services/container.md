# 컨테이너 서비스 가이드 (Container Service Guide)

이 문서는 **NCS (NHN Container Service)**와 **NCR (NHN Container Registry)** 사용법을 다룹니다.
*> NKS (Kubernetes)는 [NKS 가이드](nks.md)를 참고하세요.*

## 1. NCR (Container Registry)

프라이빗 도커 레지스트리(Private Docker Registry)를 제공합니다.

### 레지스트리 생성
```bash
nhncloud ncr create-registry --name private-repo
```

### 레지스트리 목록 조회
```bash
nhncloud ncr list-registries
```

### 레지스트리 삭제
```bash
nhncloud ncr delete-registry --registry-name private-repo
```

---

## 2. NCS (Container Service)

클러스터 관리 없이 컨테이너를 손쉽게 실행할 수 있는 서비스입니다.

### 워크로드(Workload) 생성
컨테이너 이미지를 배포합니다.
```bash
nhncloud ncs create-workload \
  --name my-app \
  --image nginx:latest \
  --container-port 80 \
  --cpu 1 \
  --memory 1Gi \
  --vpc-subnet-id <subnet-uuid>
```

### 워크로드 목록 조회
```bash
nhncloud ncs describe-workloads
```

### 워크로드 삭제
```bash
nhncloud ncs delete-workload --workload-id <uuid>
```
