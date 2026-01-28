# Compute 서비스 가이드 (Compute Service Guide)

이 문서는 NHN Cloud CLI를 사용하여 **가상 머신(인스턴스)**, **키페어**, **플레이버(사양)**, **이미지(OS)**를 관리하는 방법을 설명합니다.

## 사전 요구사항 (Prerequisites)
- **인증 (Authentication)**: Compute 서비스를 사용하기 위해서는 **Identity 자격 증명** (Tenant ID, Username, Password) 설정이 필수입니다.
- **주의**: Compute API는 OpenStack Nova/Keystone을 기반으로 하므로, AppKey만으로는 인증되지 않습니다. 반드시 사용자/테넌트 정보를 설정해주세요.
- 자세한 내용은 [설정 가이드](../CONFIGURATION.md)를 참고하세요.

---

## 1. 인스턴스 (Instances)

### 인스턴스 목록 조회 (List Instances)
현재 리전과 프로젝트(Tenant)에 있는 모든 인스턴스를 조회합니다.
```bash
nhncloud compute list-instances
```

### 인스턴스 상세 조회 (Describe Instance)
Public IP, 보안 그룹, 상태 등 인스턴스의 상세 정보를 확인합니다.
```bash
nhncloud compute describe-instances --instance-id <uuid>
```

### 인스턴스 생성 (Create Instance)
새로운 인스턴스를 생성합니다.
> **참고**: 생성을 위해서는 `flavor-id`, `image-id`, `network-id` 정보를 미리 알고 있어야 합니다.

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

### 인스턴스 삭제 (Delete Instance)
```bash
nhncloud compute delete-instance --instance-id <uuid>
```

---

## 2. 키페어 (Keypairs)
리눅스 인스턴스에 SSH로 접속하기 위해 사용하는 키페어(Public/Private Key)를 관리합니다.

> **중요**: CLI를 통해 생성된 키페어는 NHN Cloud 콘솔의 [키페어] 메뉴에서도 확인할 수 있습니다.

### 키페어 목록 조회
```bash
nhncloud compute list-key-pairs
```

### 키페어 생성 (Create Keypair)
새로운 키페어를 생성하고 Private Key를 발급받습니다. 화면에 출력된 Private Key를 안전한 곳에 저장해야 합니다.
```bash
nhncloud compute create-key-pair --key-name my-key
```

### 키페어 가져오기 (Import Keypair)
로컬에 이미 존재하는 공개키(`id_rsa.pub` 등)를 클라우드에 등록합니다.
```bash
nhncloud compute create-key-pair --key-name my-imported-key --public-key "$(cat ~/.ssh/id_rsa.pub)"
```

### 키페어 삭제
```bash
nhncloud compute delete-key-pair --key-name my-key
```

---

## 3. 자원 조회 (Flavors & Images)

### 플레이버 목록 조회 (List Flavors)
생성 가능한 인스턴스의 사양(CPU/RAM/Disk) 목록을 확인합니다.
```bash
nhncloud compute list-flavors
```

### 이미지 목록 조회 (List Images)
사용 가능한 OS 이미지(Ubuntu, CentOS, Windows 등) 목록을 확인합니다.
```bash
nhncloud compute list-images
```
