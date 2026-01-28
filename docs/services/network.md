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
