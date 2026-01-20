# MariaDB API Constraints & Known Issues

This document details all constraints and known issues from the official MariaDB v3.0 API specification.

---

## API Constraints

### Instance Creation

| Field | Type | Required | Min | Max | Notes |
|-------|------|----------|-----|-----|-------|
| dbInstanceName | string | Y | 1 | 100 chars | |
| description | string | N | - | 255 chars | |
| **dbPassword** | string | Y | **4** | **16** | **API allows 4-16, NOT 8+ minimum** |
| dbPort | int | N | 3306 | 43306 | Default: 3306 |
| dbUserName | string | Y | 2 | 32 chars | |
| storage.storageSize | int | Y | 20 | 2048 | GB |
| backup.backupPeriod | int | Y | 0 | 730 | Days |

### DB User Creation

| Field | Type | Required | Min | Max | Notes |
|-------|------|----------|-----|-----|-------|
| dbUserName | string | Y | 2 | 32 chars | |
| **dbPassword** | string | Y | **4** | **16** | **Same as instance password** |
| host | string | Y | - | - | Use `%` for any host |
| authorityType | string | Y | - | - | READ, CRUD, DDL |

### Authentication Plugins

**MariaDB Supported**:
- `NATIVE`
- `SHA256`

**NOT Supported** (MySQL only):
- ~~CACHING_SHA2~~

### TLS Options

❌ **NOT SUPPORTED** in MariaDB

---

## Known Issues (CSP)

### CSP-004: Security Group Rules Required ⚠️

**Issue**: Unlike MySQL, MariaDB **requires** at least one rule when creating a security group.

**API**: `POST /v3.0/db-security-groups`

**Impact**: Medium

**Workaround**: Always include at least one rule in the `rules` array

```go
// ❌ FAILS in MariaDB (works in MySQL)
req := &CreateSecurityGroupRequest{
    DBSecurityGroupName: "my-sg",
    Rules: []SecurityRule{}, // Empty array NOT allowed
}

// ✅ WORKS
req := &CreateSecurityGroupRequest{
    DBSecurityGroupName: "my-sg",
    Rules: []SecurityRule{
        {
            Description: "Allow DB",
            Direction:   "INGRESS",
            EtherType:   "IPV4",
            Port:        RulePort{PortType: "DB_PORT"},
            CIDR:        "0.0.0.0/0",
        },
    },
}
```

**SDK Validation**: ✅ Enforced

---

### CSP-009: Replica Creation May Fail

**Issue**: `POST /v3.0/db-instances/{id}/replicate` may return 500 error

**Impact**: Critical

**Workaround**: Retry after ensuring source instance is healthy

**SDK Validation**: ⚠️ Warning in documentation

---

### CSP-011: HA Repair Fails on Healthy Instances

**Issue**: `POST /v3.0/db-instances/{id}/high-availability/repair` returns 500 error on healthy instances

**Impact**: Critical

**Workaround**: Only call when HA is actually broken

**SDK Validation**: ⚠️ Warning in documentation

---

## Port Constraints

### dbPort Range

- **API Constraint**: 3306 - 43306
- **Default**: 3306
- **SDK Validation**: ✅ Will add

---

## Comparison with MySQL

### MariaDB is MORE Strict

| Feature | MySQL | MariaDB | Impact |
|---------|-------|---------|--------|
| Security Group Rules | Optional | **Required** (CSP-004) | Breaking |
| Auth Plugins | 3 options | 2 options (no CACHING_SHA2) | Compatibility |
| TLS Options | Supported | **Not supported** | Feature missing |

### MariaDB is LESS Strict

| Feature | MySQL | MariaDB | Impact |
|---------|-------|---------|--------|
| Password Length | Same (4-16) | Same (4-16) | None |

---

## CLI vs API Constraint Gaps

These gaps exist in the **old CLI** and should be fixed in v2:

| Field | API Constraint | Old CLI Constraint | Gap |
|-------|---------------|-------------------|-----|
| dbPassword min | **4 chars** | 8 chars | **CLI too strict** |
| dbPassword max | **16 chars** | No limit | **CLI missing** |
| dbPassword pattern | **None** | Uppercase+Lowercase+Number | **CLI too strict** |
| dbPort range | **3306-43306** | 1024-65535 | **CLI too loose** |

**Action**: New CLI should follow API constraints exactly.

---

## Summary

MariaDB API differences from MySQL:
1. ✅ **CSP-004**: Security group rules REQUIRED (enforced in SDK)
2. ✅ **Auth Plugins**: Only NATIVE, SHA256 (documented)
3. ✅ **TLS**: Not supported (removed from SDK)
4. ⚠️ **CSP-009**: Replica 500 error (document in code)
5. ⚠️ **CSP-011**: HA Repair 500 error (document in code)
6. 🔲 **Port Range**: 3306-43306 (need to add validation)
7. 🔲 **Password Length**: 4-16 chars (already validated)

**Critical**: Ensure password validation is 4-16, NOT 8+
