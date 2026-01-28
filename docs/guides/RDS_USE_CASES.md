# RDS 활용 사례 (Use Cases)

이 가이드는 MySQL, MariaDB, PostgreSQL의 전체 생성 및 관리 라이프사이클 예제를 다룹니다.

## 1. MySQL

### 전체 라이프사이클 (End-to-End Lifecycle)

인스턴스 생성부터 네트워크 설정, 네이티브 CLI 접속, 그리고 정리(삭제)까지의 과정을 설명합니다.

#### 1단계: 인스턴스 생성 (Create Instance)
```bash
nhncloud rds-mysql create-db-instance \
  --db-instance-identifier my-mysql-prod \
  --db-flavor-id m2.c4m8 \
  --engine-version MYSQL_V8032 \
  --master-username admin \
  --master-user-password 'SecurePass123!' \
  --allocated-storage 100 \
  --subnet-id subnet-xxxxxxxx \
  --availability-zone kr-pub-a \
  --db-parameter-group-id default-mysql8
```

#### 2단계: 네트워크 접근 설정 (Configure Network Access)
```bash
# 1. 보안 그룹 생성
nhncloud rds-mysql create-db-security-group \
  --db-security-group-name mysq-prod-sg \
  --description "Allow Access"

# 2. 3306 포트 허용 (Ingress Rule)
nhncloud rds-mysql authorize-db-security-group-ingress \
  --db-security-group-identifier mysq-prod-sg \
  --cidr 0.0.0.0/0 \
  --port 3306

# 3. 인스턴스에 보안 그룹 연결
nhncloud rds-mysql modify-db-instance \
  --db-instance-identifier my-mysql-prod \
  --db-security-group-ids mysq-prod-sg
```

#### 3단계: 접속 및 쿼리 (Connect & Query)
```bash
# 대화형 SQL 쉘 (Interactive SQL Shell)
nhncloud rds-mysql connect \
  --db-instance-identifier my-mysql-prod \
  --username admin \
  --password 'SecurePass123!'

# 단일 쿼리 실행 (One-liner)
nhncloud rds-mysql connect \
  --db-instance-identifier my-mysql-prod \
  --username admin \
  --password 'SecurePass123!' \
  -e "SELECT VERSION();"
```

#### 4단계: 정리 (Cleanup)
```bash
nhncloud rds-mysql delete-db-instance --db-instance-identifier my-mysql-prod
nhncloud rds-mysql delete-db-security-group --db-security-group-name mysq-prod-sg
```

---

## 2. MariaDB

### 전체 라이프사이클 (End-to-End Lifecycle)

#### 1단계: 인스턴스 생성
```bash
nhncloud rds-mariadb create-db-instance \
  --db-instance-identifier my-mariadb-prod \
  --db-flavor-id m2.c2m4 \
  --engine-version MARIADB_V1011 \
  --master-username admin \
  --master-user-password 'SecurePass123!' \
  --allocated-storage 50 \
  --subnet-id subnet-xxxxxxxx \
  --availability-zone kr-pub-a
```

#### 2단계: 네트워크 접근 설정
MariaDB 보안 그룹은 생성 시 초기 CIDR 규칙을 설정할 수 있습니다.
```bash
# 1. 보안 그룹 생성 및 초기 규칙 설정
nhncloud rds-mariadb create-db-security-group \
  --db-security-group-name mariadb-sg \
  --cidr 0.0.0.0/0
  
# 2. 인스턴스에 연결
nhncloud rds-mariadb modify-db-instance \
  --db-instance-identifier my-mariadb-prod \
  --db-security-group-ids mariadb-sg
```

#### 3단계: 접속 및 쿼리
```bash
nhncloud rds-mariadb connect \
  --db-instance-identifier my-mariadb-prod \
  --username admin \
  --password 'SecurePass123!' \
  -e "SHOW DATABASES;"
```

#### 4단계: 정리
```bash
nhncloud rds-mariadb delete-db-instance --db-instance-identifier my-mariadb-prod
nhncloud rds-mariadb delete-db-security-group --db-security-group-name mariadb-sg
```

---

## 3. PostgreSQL

### 전체 라이프사이클 (End-to-End Lifecycle)

#### 1단계: 인스턴스 생성
```bash
nhncloud rds-postgresql create-db-instance \
  --db-instance-identifier my-pg-prod \
  --db-flavor-id m2.c4m8 \
  --engine-version POSTGRESQL_V13 \
  --master-username admin \
  --master-user-password 'SecurePass123!' \
  --allocated-storage 20 \
  --subnet-id subnet-xxxxxxxx \
  --availability-zone kr-pub-a
```

#### 2단계: 네트워크 접근 설정
**중요**: PostgreSQL은 별도의 보안 그룹 체계를 사용합니다.
```bash
# 1. Ingress 규칙 추가 (Port 5432)
nhncloud rds-postgresql authorize-db-security-group-ingress \
  --db-security-group-identifier <pg-sg-id> \
  --cidr 0.0.0.0/0 \
  --port 5432

# 2. 보안 그룹 연결
nhncloud rds-postgresql modify-db-instance \
  --db-instance-identifier my-pg-prod \
  --db-security-group-ids <pg-sg-id>
```

#### 3단계: 접속 및 쿼리
```bash
# 대화형 쉘 접속 (REPL)
nhncloud rds-postgresql connect \
  --db-instance-identifier my-pg-prod \
  --username admin \
  --password 'SecurePass123!'
```

#### 4단계: 정리
```bash
nhncloud rds-postgresql delete-db-instance --db-instance-identifier my-pg-prod
```
