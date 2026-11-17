# 네트워크 서비스 가이드 (Network Service Guide)

이 문서는 **VPC**, **서브넷(Subnets)**, **Floating IP(공인 IP)** 관리를 다룹니다.

## 1. VPC & 서브넷 (VPC & Subnets)

### VPC 목록 조회 (List VPCs)
```bash
nhncloud network describe-vpcs
```

### 서브넷 목록 조회 (List Subnets)
인스턴스 및 RDS 생성 시 필요한 서브넷(Subnet) UUID를 확인할 수 있습니다.
```bash
nhncloud network describe-subnets
```

## 2. Floating IP (공인 IP)

### Floating IP 생성 (Create Floating IP)
새로운 공인 IP를 할당받습니다.
```bash
nhncloud network create-floating-ip --floating-ip <ip-address> # 또는 자동 할당
```

### Floating IP 목록 조회 (List Floating IPs)
```bash
nhncloud network describe-floating-ips
```

## 3. 라우팅 테이블 (Routing Tables)

VPC의 기본 라우팅 테이블을 찾고, 인터넷 게이트웨이와 서브넷을 연결합니다.
```bash
nhncloud network describe-routing-tables --vpc-id <vpc-id>
nhncloud network create-routing-table --name rt-a --vpc-id <vpc-id>
nhncloud network attach-routing-table-gateway --routing-table-id <id> --gateway-id <igw-id>
nhncloud network attach-subnet-routing-table --subnet-id <subnet-id> --routing-table-id <id>
nhncloud network set-default-routing-table --routing-table-id <id>
nhncloud network create-route --routing-table-id <id> --cidr 10.10.0.0/24 --gateway 192.168.0.5
```

`internet-gateway create-internet-gateway`에 `--routing-table-id`를 주면 생성한 뒤 그 라우팅 테이블에 연결합니다.

## 4. 포트 (Ports)

인스턴스나 로드 밸런서 VIP의 포트 ID를 찾습니다.
```bash
nhncloud network describe-ports --device-id <instance-id>
nhncloud network describe-ports --fixed-ip 192.168.0.10
nhncloud network associate-floating-ip --floating-ip-id <fip-id> --instance-id <instance-id>
```

## 5. 리소스 한도 (Quotas)

프로젝트가 쓰는 네트워크 리소스 수를 공개 리소스 제공 정책의 기본 한도와 함께 보여 줍니다.
한도를 올려 받은 프로젝트는 실제 한도가 더 큽니다.
```bash
nhncloud network describe-quotas
```

## 6. 로드 밸런서 리스너 옵션

```bash
nhncloud loadbalancer create-listener --name l1 --lb-id <lb-id> --protocol TCP --port 3306 \
  --keepalive-timeout 3600 --proxy-protocol
nhncloud loadbalancer update-listener --listener-id <id> --keepalive-timeout 600
nhncloud loadbalancer describe-listeners --listener-id <id> --all-fields
```

## 7. 서비스 게이트웨이 (Service Gateway)

```bash
nhncloud service-gateway create-service-gateway --name sgw --subnet-id <subnet-id> \
  --service-endpoint-id <endpoint-id> --include-gateway-identity
nhncloud service-gateway create-my-service-endpoint --name my-svc --port-id <lb-vip-port-id>
nhncloud service-gateway add-endpoint-allowed-project --endpoint-id <id> --target-project-id <tenant-id>
nhncloud service-gateway describe-endpoint-allowed-projects --endpoint-id <id>
```

## 8. 전용 명령이 없는 API

```bash
nhncloud api /v2.0/routingtables
nhncloud api -X PUT /v2.0/routingtables/<id>/set_as_default
```
