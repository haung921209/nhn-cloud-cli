# RDS (Database) 서비스 가이드

이 가이드는 **MySQL**, **MariaDB**, **PostgreSQL** 데이터베이스 인스턴스를 CLI로 관리하는 방법을 설명합니다.

## 사전 요구사항 (Prerequisites)
- **인증 (Authentication)**: RDS를 사용하려면 해당 DB 엔진에 맞는 **AppKey**가 설정되어 있어야 합니다.
- **AppKey 설정**: 환경 변수(`NHN_CLOUD_MYSQL_APPKEY` 등) 또는 설정 파일을 통해 등록하세요.
- 자세한 내용은 [설정 가이드](../CONFIGURATION.md)를 참고하세요.

---

## 1. 인스턴스 관리 (Instance Management)

### 인스턴스 생성 (Create Instance)
가장 기본적인 MySQL 인스턴스 생성 예제입니다.

```bash
nhncloud rds-mysql create-instance \
  --name my-db-01 \
  --flavor-id <flavor-uuid> \
  --version 8.0 \
  --username admin \
  --password <strong-password> \
  --subnet-id <subnet-uuid> \
  --multi-az false
```
*팁: MariaDB는 `rds-mariadb`, PostgreSQL은 `rds-postgresql` 명령어를 사용하세요.*

### 인스턴스 목록 조회 (List Instances)
현재 프로젝트에 생성된 인스턴스 목록을 확인합니다.
```bash
nhncloud rds-mysql describe-db-instances
```

### 인스턴스 삭제 (Delete Instance)
더 이상 필요하지 않은 인스턴스를 삭제합니다.
```bash
nhncloud rds-mysql delete-instance --instance-id <instance-id>
```

---

## 2. 고급 설정 (Advanced Configuration)

### 파라미터 그룹 (Parameter Groups)
DB 엔진의 상세 설정(Buffer Pool Size, Timeout 등)을 튜닝할 수 있습니다.

```bash
# 사용 가능한 파라미터 그룹 목록 조회
nhncloud rds-mysql describe-db-parameter-groups

# 사용자 정의 파라미터 그룹 생성
nhncloud rds-mysql create-db-parameter-group --name my-params --description "Production config"

# 인스턴스에 파라미터 그룹 적용 (재시작 필요)
nhncloud rds-mysql modify-db-instance --instance-id <id> --parameter-group <group-name>
```

### 보안 그룹 (Security Groups) - 접근 제어
데이터베이스 포트(3306/5432)에 대한 접근을 제어합니다.

```bash
# 보안 그룹 생성
nhncloud rds-mysql create-db-security-group --name my-db-sg

# 인바운드 규칙 추가 (특정 IP 허용)
nhncloud rds-mysql authorize-db-security-group-ingress \
  --db-security-group-identifier my-db-sg \
  --cidr 203.0.113.5/32

# 인스턴스에 보안 그룹 적용
nhncloud rds-mysql modify-db-instance --instance-id <id> --security-groups my-db-sg
```

---

## 3. 고가용성 관리 (High Availability)

### Multi-AZ 활성화 (Enable Multi-AZ)
단일 인스턴스를 고가용성(Master + Standby) 구조로 승격시킵니다. 장애 발생 시 자동으로 절체(Failover)됩니다.

```bash
nhncloud rds-mysql enable-multi-az --instance-id <id> --subnet-id <standby-subnet-id>
```

### 읽기 전용 복제본 생성 (Read Replicas)
읽기 분하 분산을 위한 Read Replica를 생성합니다.
```bash
nhncloud rds-mysql create-read-replica \
  --source-db-instance-identifier <master-id> \
  --name my-replica-01
```

---

## 4. 사용자 및 스키마 관리 (Users and Schemas)

### DB 사용자 생성
```bash
nhncloud rds-mysql create-db-user \
  --instance-id <id> \
  --username app_user \
  --password <password> \
  --host %
```

### 스키마(Database) 생성
```bash
nhncloud rds-mysql create-db-schema --instance-id <id> --name my_app_db
```

---

## 5. 데이터베이스 접속 (Secure Connection)

CLI는 로컬 DB 클라이언트와 연동하여 **SSL/TLS**가 적용된 안전한 접속을 자동으로 처리해줍니다.

### 사전 준비
1. **공인 IP (Floating IP)**: 인스턴스에 공인 IP가 연결되어 있어야 외부 접속이 가능합니다.
2. **보안 그룹**: 접속하려는 클라이언트의 IP가 허용되어 있어야 합니다.
3. **CA 인증서**: `nhncloud config ca import` 명령으로 CA 인증서가 등록되어 있어야 합니다.

### 사용법
```bash
# 기본 접속 (호스트 자동 감지 및 SSL 옵션 주입)
nhncloud rds-mysql connect \
  --db-instance-identifier <instance-id> \
  --username <user> \
  --password <password> \
  --database <db>

# 비대화형 쿼리 실행
nhncloud rds-mysql connect ... -- -e "SELECT version();"

# 인증 플러그인 지정 (필요 시)
# 기본값: caching_sha2_password (SSL 강제)
# 레거시: mysql_native_password (엄격한 SSL 체크 우회 가능)
nhncloud rds-mysql connect ... --auth-plugin mysql_native_password
```
