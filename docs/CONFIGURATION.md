# 서비스 환경 설정 가이드 (Service Configuration Guide)

NHN Cloud CLI를 사용하기 위한 인증 설정을 안내합니다. 
**사용하려는 서비스에 따라 필요한 인증 정보가 다릅니다.**

> **설정 우선순위 (Precedence Rule)**:
> 1. **CLI 플래그** (예: `--region`) 
> 2. **환경 변수** (예: `NHN_CLOUD_PASSWORD`)
> 3. **설정 파일** (`~/.nhncloud/credentials`)

---

## 1. 인증 정보 구분 (Credential Types)

NHN Cloud는 크게 두 가지 인증 방식을 혼용하여 사용합니다.

### A. API Key / AppKey 방식 (주로 플랫폼 서비스)
**대상**: `RDS`, `Object Storage(S3 API)`, `NCR`, `NCS` 등
**필요 정보**: 
- `access_key_id` & `secret_access_key` (API 인증 토큰 발급용)
- `appkey` (서비스 인스턴스 식별용)

### B. Identity (OpenStack) 방식 (주로 인프라 서비스)
**대상**: `Compute (Nova)`, `Network (VPC)`, `Block Storage`
**필요 정보**: 
- `username` (이메일 아이디)
- `api_password` (API 비밀번호 - 콘솔 로그인 비밀번호와 다름)
- `tenant_id` (프로젝트 아이디)

> **주의**: `Compute` 명령어를 사용하려면 `api_password`가 반드시 설정되어야 합니다. `access_key_id`만으로는 컴퓨트 인스턴스를 제어할 수 없습니다.

---

## 2. 설정 파일 예시 (Configuration File Example)

`~/.nhncloud/credentials` 파일에 필요한 모든 정보를 한 번에 설정할 수 있습니다.

```ini
[default]
# 공통 설정
region = kr1

# 1. API Key 방식 (RDS, Storage 등)
access_key_id = <ACCESS_KEY_ID>
secret_access_key = <SECRET_ACCESS_KEY>
# AppKey는 필요에 따라 추가
rds_mysql_app_key = <APPKEY>

# 2. Identity 방식 (Compute, Network 등)
# Compute 서비스를 쓰지 않는다면 생략 가능하지만, 사용 시 필수입니다.
username = user@nhn.com
api_password = <API_PASSWORD>
tenant_id = <TENANT_ID>
```

## 3. 상세 항목 설명

| 키 (Key) | 설명 | 환경 변수 매핑 | 필수 서비스 |
|----------|------|----------------|-------------|
| `access_key_id` | API Access Key ID | `NHN_CLOUD_ACCESS_KEY` | RDS, S3 |
| `secret_access_key` | API Secret Key | `NHN_CLOUD_SECRET_KEY` | RDS, S3 |
| `api_password` | **API 전용 비밀번호** (콘솔 > 회원정보 > API보안설정) | `NHN_CLOUD_PASSWORD` | **Compute**, Network |
| `username` | NHN Cloud ID (이메일) | `NHN_CLOUD_USERNAME` | **Compute**, Network |
| `tenant_id` | 프로젝트(Tenant) ID | `NHN_CLOUD_TENANT_ID` | **Compute**, Network |
| `appkey` | 기본 AppKey | `NHN_CLOUD_APPKEY` | 공통 |

---

---

## 4. 데이터베이스 연결 설정 (SSL/TLS)

RDS 인스턴스의 **Data Plane** (실제 SQL 연결)에 접속하기 위해서는 보안을 위해 SSL 연결이 권장되며, NHN Cloud CA 인증서가 필요합니다. CLI는 이를 쉽게 관리할 수 있는 도구를 제공합니다.

### CA 인증서 관리 (CA Certificate Management)

#### 인증서 가져오기 (Import)
Root CA 인증서나 Client 인증서를 CLI 내부 보안 저장소에 등록합니다.

```bash
# Root CA 인증서 등록
nhncloud config ca import \
  --service rds-mysql \
  --region kr1 \
  --file ./ca.pem \
  --description "NHN Cloud Root CA"

# Client 인증서 등록 (Mutual TLS 사용 시)
nhncloud config ca import \
  --service rds-mysql \
  --region kr1 \
  --type CLIENT-CERT \
  --instance-id <instance-uuid> \
  --file ./client-cert.pem
```

#### 인증서 목록 확인
```bash
nhncloud config ca list --service rds-mysql
```

### 자동 연결 도우미 (Automatic Connection Helper)
CLI를 통해 로컬에 설치된 DB 클라이언트(`mysql`, `psql`)를 실행하면, 복잡한 SSL 옵션(`--ssl-ca`, `--ssl-cert` 등)을 자동으로 주입해줍니다.

```bash
# MySQL 접속 (자동으로 SSL 옵션 적용됨)
nhncloud rds-mysql connect \
  --db-instance-identifier <instance-uuid> \
  --username <user> \
  --password <pass> \
  --database <db>

# 비대화형 모드로 쿼리 실행 예제
nhncloud rds-mysql connect ... -- -e "SELECT 1;"
```

> **참고**: 이 기능을 사용하려면 시스템 PATH에 `mysql` 또는 `psql` 클라이언트가 설치되어 있어야 합니다.

