# RDS Use Cases

This guide covers end-to-end lifecycle examples for MySQL, MariaDB, and PostgreSQL.

## 1. MySQL

### End-to-End Lifecycle

This scenario covers creating an instance, configuring network access, connecting via the native CLI shell, and cleaning up.

#### Step 1: Create Instance
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

#### Step 2: Configure Network Access
```bash
# 1. Create Security Group
nhncloud rds-mysql create-db-security-group \
  --db-security-group-name mysq-prod-sg \
  --description "Allow Access"

# 2. Authorize Port 3306
nhncloud rds-mysql authorize-db-security-group-ingress \
  --db-security-group-identifier mysq-prod-sg \
  --cidr 0.0.0.0/0 \
  --port 3306

# 3. Attach to Instance
nhncloud rds-mysql modify-db-instance \
  --db-instance-identifier my-mysql-prod \
  --db-security-group-ids mysq-prod-sg
```

#### Step 3: Connect & Query
```bash
# Interactive SQL Shell
nhncloud rds-mysql connect \
  --db-instance-identifier my-mysql-prod \
  --username admin \
  --password 'SecurePass123!'

# One-liner Query
nhncloud rds-mysql connect \
  --db-instance-identifier my-mysql-prod \
  --username admin \
  --password 'SecurePass123!' \
  -e "SELECT VERSION();"
```

#### Step 4: Cleanup
```bash
nhncloud rds-mysql delete-db-instance --db-instance-identifier my-mysql-prod
nhncloud rds-mysql delete-db-security-group --db-security-group-name mysq-prod-sg
```

---

## 2. MariaDB

### End-to-End Lifecycle

#### Step 1: Create Instance
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

#### Step 2: Configure Network Access
MariaDB security groups require an initial CIDR rule upon creation.
```bash
# 1. Create Security Group with Initial Rule
nhncloud rds-mariadb create-db-security-group \
  --db-security-group-name mariadb-sg \
  --cidr 0.0.0.0/0
  
# 2. Attach to Instance
nhncloud rds-mariadb modify-db-instance \
  --db-instance-identifier my-mariadb-prod \
  --db-security-group-ids mariadb-sg
```

#### Step 3: Connect & Query
```bash
nhncloud rds-mariadb connect \
  --db-instance-identifier my-mariadb-prod \
  --username admin \
  --password 'SecurePass123!' \
  -e "SHOW DATABASES;"
```

#### Step 4: Cleanup
```bash
nhncloud rds-mariadb delete-db-instance --db-instance-identifier my-mariadb-prod
nhncloud rds-mariadb delete-db-security-group --db-security-group-name mariadb-sg
```

---

## 3. PostgreSQL

### End-to-End Lifecycle

#### Step 1: Create Instance
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

#### Step 2: Configure Network Access
**Important**: PostgreSQL uses a separate set of security groups.
```bash
# 1. Add Ingress Rule (Port 5432)
nhncloud rds-postgresql authorize-db-security-group-ingress \
  --db-security-group-identifier <pg-sg-id> \
  --cidr 0.0.0.0/0 \
  --port 5432

# 2. Attach Security Group
nhncloud rds-postgresql modify-db-instance \
  --db-instance-identifier my-pg-prod \
  --db-security-group-ids <pg-sg-id>
```

#### Step 3: Connect & Query
```bash
# Interactive Shell (REPL)
nhncloud rds-postgresql connect \
  --db-instance-identifier my-pg-prod \
  --username admin \
  --password 'SecurePass123!'
```

#### Step 4: Cleanup
```bash
nhncloud rds-postgresql delete-db-instance --db-instance-identifier my-pg-prod
```
