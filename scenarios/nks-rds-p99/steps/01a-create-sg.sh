#!/usr/bin/env bash
# steps/01a-create-sg.sh
#
# Phase B variant — creates TWO per-run security groups:
#
#   1. bench-mysql-<RUN_TS>   (RDS-MySQL DB security group)
#        ingress 3306        from <subnet_cidr>   (Compute → MySQL inside VPC)
#        ingress 3306-43306  from <subnet_cidr>   (NHN's in-VPC plumbing — same
#                                                  pattern as Phase A; required
#                                                  or create-db-instance fails
#                                                  validation asynchronously)
#       NO public-IP /32 ingress — DB is internal-only per scenario.yaml
#       use_public_access=false (B.6).
#
#   2. bench-compute-<RUN_TS> (Network security group attached to Compute port)
#        ingress 22          from <runner-public-ip>/32   (SSH from runner)
#
# Persists: BENCH_MYSQL_SG_ID, BENCH_MYSQL_SG_NAME,
#           BENCH_COMPUTE_SG_ID, BENCH_COMPUTE_SG_NAME,
#           BENCH_PUBLIC_IP, BENCH_SUBNET_CIDR.
#
# Ref: https://docs.nhncloud.com/en/Database/RDS%20for%20MySQL/en/public-api/ (Korean: rds-mysql-v4.0 spec)#db-보안-그룹-생성하기
# Ref: https://docs.nhncloud.com/en/Database/RDS%20for%20MySQL/en/public-api/ (Korean: rds-mysql-v4.0 spec)#db-보안-그룹-규칙-생성하기
# Ref: https://docs.nhncloud.com/en/Network/Network/en/public-api/#security-group (network create-security-group)

set -euo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"
: "${RUN_TS:?must be set by run.sh}"

# shellcheck disable=SC1091
source "$SCENARIO_DIR/state.env" 2>/dev/null || true
: "${ENGINE:=mysql}"
CLI_PREFIX=$(yq -r ".parameters.${ENGINE}.cli_prefix // \"rds-mysql\"" "$SCENARIO_DIR/scenario.yaml")
# Brew CLI v0.7.18 lacks create-db-security-group / delete-db-security-group
# for rds-postgresql. run.sh builds the workspace CLI into bin/nhncloud-ws and
# exports NHNCLOUD_WS_BIN when ENGINE=postgresql; route only the SG-verbs
# through it. Everything else continues to use the brew nhncloud on PATH.
if [[ "$ENGINE" == "postgresql" ]]; then
  : "${NHNCLOUD_WS_BIN:?postgresql requires NHNCLOUD_WS_BIN from run.sh — workspace CLI build missing}"
  SG_CLI="$NHNCLOUD_WS_BIN"
else
  SG_CLI="nhncloud"
fi
echo "[sg] engine=$ENGINE cli_prefix=$CLI_PREFIX sg_cli=$SG_CLI"

# ─── Detect runner public IP (needed only for Compute SG) ───────────────────
PUBLIC_IP=""
for url in https://api.ipify.org https://ifconfig.me; do
  PUBLIC_IP=$(curl -fsS --max-time 5 "$url" 2>/dev/null || true)
  if [[ "$PUBLIC_IP" =~ ^[0-9.]+$ ]]; then
    break
  fi
  PUBLIC_IP=""
done
if [[ -z "$PUBLIC_IP" ]]; then
  echo "[sg] ERROR: could not resolve runner public IP" >&2
  exit 1
fi
echo "[sg] runner public IP: $PUBLIC_IP"

# ─── Resolve subnet CIDR (needed for MySQL SG) ──────────────────────────────
SUBNET_CIDR=$(yq -r '.parameters.subnet_cidr // ""' "$SCENARIO_DIR/scenario.yaml")
if [[ -z "$SUBNET_CIDR" ]]; then
  SUBNET_ID=$(yq -r '.parameters.subnet_id' "$SCENARIO_DIR/scenario.yaml")
  SUBNET_CIDR=$(nhncloud "$CLI_PREFIX" describe-subnets -o json 2>/dev/null \
                 | jq -r --arg id "$SUBNET_ID" \
                     '.subnets[]? | select(.subnetId==$id) | .subnetCidr // ""' \
                 | head -1)
fi
if [[ -z "$SUBNET_CIDR" ]]; then
  echo "[sg] ERROR: could not resolve subnet CIDR — MySQL SG cannot be scoped." >&2
  exit 1
fi
echo "[sg] subnet CIDR: $SUBNET_CIDR"

TS_SHORT="${RUN_TS:0:14}"   # YYYYMMDDTHHMM (13 chars)

# ────────────────────────────────────────────────────────────────────────────
# 1. MySQL DB security group (RDS-MySQL service, internal-only)
# ────────────────────────────────────────────────────────────────────────────
MYSQL_SG_NAME="bench-mysql-${TS_SHORT}"

# CLI emits human-readable text "Security group created: <uuid>" —
# extract via regex (also tolerates JSON-of-day in case the verb gains -o json).
# MariaDB / PostgreSQL `create-db-security-group` REQUIRES --cidr at creation
# time (initial ingress rule). MySQL's signature doesn't accept --cidr and
# emits "unknown flag" — branch on engine. PG defaults to port 5432 when
# --port is omitted (CLI source verified), matching the loadgen path.
SG_CREATE_ARGS=(create-db-security-group
  --db-security-group-name "$MYSQL_SG_NAME"
  --description "scenario nks-rds-p99 engine=$ENGINE RUN_TS=$RUN_TS"
)
if [[ "$ENGINE" == "mariadb" || "$ENGINE" == "postgresql" ]]; then
  SG_CREATE_ARGS+=(--cidr "$SUBNET_CIDR")
fi
set +e
CREATE_OUT=$("$SG_CLI" "$CLI_PREFIX" "${SG_CREATE_ARGS[@]}" 2>&1)
SG_RC=$?
set -e
if [[ $SG_RC -ne 0 ]]; then
  echo "[sg] ERROR: create-db-security-group exited $SG_RC. Output:" >&2
  printf '%s\n' "$CREATE_OUT" >&2
  exit 1
fi
MYSQL_SG_ID=$(printf '%s' "$CREATE_OUT" \
              | grep -oE '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' \
              | head -1)
if [[ -z "$MYSQL_SG_ID" ]]; then
  echo "[sg] could not parse MySQL SG id from create output; resolving by name..."
  MYSQL_SG_ID=$("$SG_CLI" "$CLI_PREFIX" describe-db-security-groups -o json 2>/dev/null \
                 | jq -r --arg n "$MYSQL_SG_NAME" \
                     '.dbSecurityGroups[] | select(.dbSecurityGroupName==$n) | .dbSecurityGroupId' \
                 | head -1)
fi
[[ -n "$MYSQL_SG_ID" ]] || {
  echo "[sg] ERROR: could not resolve MySQL SG id; create output was:" >&2
  printf '%s\n' "$CREATE_OUT" >&2
  exit 1
}
echo "[sg] created MySQL SG $MYSQL_SG_NAME ($MYSQL_SG_ID)"

# Ingress 3306 from subnet CIDR (NKS Pods → DB in-VPC).
# MariaDB's create-db-security-group already attached an initial rule with
# the --cidr we just passed (default port 3306) — adding another rule with
# the same CIDR+port returns API 500. Only run the explicit authorize for
# MySQL, where create-db-security-group leaves the SG empty.
if [[ "$ENGINE" == "mysql" ]]; then
  nhncloud "$CLI_PREFIX" authorize-db-security-group-ingress \
    --db-security-group-identifier "$MYSQL_SG_ID" \
    --cidr "$SUBNET_CIDR" \
    --port 3306 \
    --description "NKS Pods → MySQL in-VPC" >/dev/null
  echo "[sg] mysql ingress added: $SUBNET_CIDR → 3306"
else
  echo "[sg] $ENGINE ingress already attached at create time (--cidr $SUBNET_CIDR)"
fi

# ────────────────────────────────────────────────────────────────────────────
# 2. Compute security group (Network service)
# ────────────────────────────────────────────────────────────────────────────
COMPUTE_SG_NAME="bench-compute-${TS_SHORT}"

# `network create-security-group` returns JSON when -o json is passed.
COMPUTE_SG_JSON=$(nhncloud network create-security-group -o json \
                    --name "$COMPUTE_SG_NAME" \
                    --description "scenario rds-mysql-p99-internal RUN_TS=$RUN_TS" 2>&1)
COMPUTE_SG_ID=$(printf '%s' "$COMPUTE_SG_JSON" \
                | jq -r '.securityGroup.id // .id // empty' 2>/dev/null)
if [[ -z "$COMPUTE_SG_ID" ]]; then
  # Fallback: regex over text output.
  COMPUTE_SG_ID=$(printf '%s' "$COMPUTE_SG_JSON" \
                  | grep -oE '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' \
                  | head -1)
fi
if [[ -z "$COMPUTE_SG_ID" ]]; then
  echo "[sg] could not parse Compute SG id; resolving by name..."
  COMPUTE_SG_ID=$(nhncloud network describe-security-groups -o json 2>/dev/null \
                  | jq -r --arg n "$COMPUTE_SG_NAME" \
                      '.. | objects | select(.name==$n) | .id' \
                  | head -1)
fi
[[ -n "$COMPUTE_SG_ID" ]] || {
  echo "[sg] ERROR: could not resolve Compute SG id; create output was:" >&2
  printf '%s\n' "$COMPUTE_SG_JSON" >&2
  exit 1
}
echo "[sg] created Compute SG $COMPUTE_SG_NAME ($COMPUTE_SG_ID)"

# Ingress 22 from runner /32 — SSH from this host only.
nhncloud network authorize-security-group-ingress \
  --group-id "$COMPUTE_SG_ID" \
  --direction ingress \
  --ethertype IPv4 \
  --protocol tcp \
  --port-min 22 --port-max 22 \
  --remote-ip "$PUBLIC_IP/32" >/dev/null
echo "[sg] compute ingress added: $PUBLIC_IP/32 → 22"

# ────────────────────────────────────────────────────────────────────────────
# Persist state
# ────────────────────────────────────────────────────────────────────────────
{
  echo "BENCH_MYSQL_SG_ID=$MYSQL_SG_ID"
  echo "BENCH_MYSQL_SG_NAME=$MYSQL_SG_NAME"
  echo "BENCH_COMPUTE_SG_ID=$COMPUTE_SG_ID"
  echo "BENCH_COMPUTE_SG_NAME=$COMPUTE_SG_NAME"
  echo "BENCH_PUBLIC_IP=$PUBLIC_IP"
  echo "BENCH_SUBNET_CIDR=$SUBNET_CIDR"
} >> "$SCENARIO_DIR/state.env"

# Empirically, NHN's create-db-instance fails with FAIL_TO_READY when
# the supplied SG was created seconds earlier — the SG/rule write
# hasn't propagated to the create-validator yet. 30s is enough.
echo "[sg] waiting 30s for SG propagation..."
sleep 30
