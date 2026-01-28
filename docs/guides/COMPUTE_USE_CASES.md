# Compute & Network 활용 사례 (Use Cases)

이 문서는 NHN Cloud CLI를 사용하여 Compute 인스턴스와 네트워크 리소스를 관리하는 일반적인 시나리오를 다룹니다.

## 1. 인스턴스 접속 (SSH)

CLI는 인스턴스의 공인 IP와 SSH 키를 자동으로 감지하여 접속을 도와주는 편리한 `connect` 명령어를 제공합니다.

### 사전 준비: 키 가져오기 (Import Key)
키 페어당 한 번만 수행하면 됩니다. 다운로드 받은 PEM 파일을 안전하게 로컬에 등록합니다.

```bash
nhncloud compute import-key \
  --key-name my-key-pair \
  --private-key-file ~/Downloads/my-key-pair.pem
```

### 인스턴스 접속 (Connect to Instance)

```bash
# 1. 표준 접속 (자동으로 IP 및 키 검색)
nhncloud compute connect --instance-id <instance-id>

# 2. 사용자명 지정 (기본값은 메타데이터의 'login_username' 또는 'centos')
nhncloud compute connect --instance-id <instance-id> --username ubuntu

# 3. 키 파일 직접 지정
nhncloud compute connect \
  --instance-id <instance-id> \
  --identity-file ~/.ssh/my-special-key.pem
```

---

## 2. 인스턴스 관리 (Instance Management)

### 인스턴스 생성 (Create Instance)
```bash
nhncloud compute create-instance \
  --name my-web-server \
  --image-id <image-id> \
  --flavor-id <flavor-id> \
  --subnet-id <subnet-id> \
  --key-name my-key-pair
```

### 인스턴스 목록 조회 (List Instances)
```bash
nhncloud compute describe-instances
```

### 전원 관리 (Power Management)
```bash
# 시작 (Start)
nhncloud compute start-instances --instance-id <id>

# 정지 (Stop)
nhncloud compute stop-instances --instance-id <id>

# 재부팅 (Reboot - Soft)
nhncloud compute reboot-instances --instance-id <id>
```
