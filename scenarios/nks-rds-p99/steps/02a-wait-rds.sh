#!/usr/bin/env bash
# steps/02a-wait-rds.sh  (scenario: nks-rds-p99)
#
# Polls `wait-db-instance --for-state AVAILABLE`, then resolves the
# INTERNAL_VIP endpoint.  In this scenario use_public_access is always
# false — cluster Pods reach the RDS via the INTERNAL_VIP domain, so
# the runner-side TCP probe block is intentionally removed.  The loadgen
# Pod will perform its own connectivity check when it first connects.
#
# Endpoint resolution path (internal-only):
#   describe-db-instance-network → INTERNAL_VIP domain + dbPort from
#   describe-db-instances → ENDPOINT written to state.env.
#
# After the endpoint is known:
#   1. create-db-schema  (idempotent is NOT guaranteed — provisioning fresh)
#   2. create-db-user    (app user with DDL authority, host=%)
#   3. re-confirm AVAILABLE after each mutation
#
# Ref: https://docs.nhncloud.com/en/Database/ (per-engine public API)#db-인스턴스-목록-보기
# Ref: nhn-cloud-cli/cmd/rds_wait_endpoint.go
# Ref: nhn-cloud-cli/cmd/rds_describe_detail.go (describe-db-instance-network)

set -euo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"

# shellcheck disable=SC1091
source "$SCENARIO_DIR/state.env"
: "${INSTANCE_ID:?state.env did not carry INSTANCE_ID — provision step missing}"
: "${ENGINE:?state.env missing ENGINE}"
CLI_PREFIX=$(yq -r ".parameters.${ENGINE}.cli_prefix" "$SCENARIO_DIR/scenario.yaml")

# Round-2 used `nhncloud <engine-cli> wait-db-instance --for-state AVAILABLE`,
# but that command only exists on the workspace CLI (round-2 commit 2bd9b7c)
# and is missing from the brew-published binary we run with. Inline an
# equivalent polling loop against describe-db-instances.
echo "[wait] waiting for $INSTANCE_ID to reach AVAILABLE (timeout 30m)..."
DEADLINE=$(( $(date +%s) + 30 * 60 ))
last_status=""
while :; do
  if [[ $(date +%s) -ge $DEADLINE ]]; then
    echo "[wait] ERROR: timeout waiting for AVAILABLE (last status=$last_status)" >&2
    exit 1
  fi
  raw_describe=$(nhncloud "$CLI_PREFIX" describe-db-instances -o json 2>/dev/null)
  if [[ "$ENGINE" == "mariadb" || "$ENGINE" == "postgresql" ]]; then
    raw_describe=$(printf '%s\n' "$raw_describe" | sed -n '/^{/,$p')
  fi
  status=$(printf '%s' "$raw_describe" | jq -r --arg id "$INSTANCE_ID" \
            '.dbInstances[] | select(.dbInstanceId==$id) | .dbInstanceStatus' | head -1)
  if [[ "$status" != "$last_status" ]]; then
    echo "[wait] status=$status"
    last_status="$status"
  else
    printf '.'
  fi
  if [[ "$status" == "AVAILABLE" ]]; then
    echo
    break
  fi
  if [[ "$status" == *"FAIL"* || "$status" == *"ERROR"* ]]; then
    echo
    echo "[wait] ERROR: instance reached terminal failure status=$status" >&2
    exit 1
  fi
  sleep 15
done

# ─── Resolve endpoint (internal-only) ───────────────────────────────────────
# Internal instance: describe-db-instance-network prints a
# "Fetching network info for instance <id>...\n" prefix line before the
# JSON body; strip everything before the first '{' before piping to jq.
NET_RAW=$(nhncloud "$CLI_PREFIX" describe-db-instance-network \
             --db-instance-identifier "$INSTANCE_ID" -o json 2>/dev/null)
NET_JSON=$(printf '%s' "$NET_RAW" | sed -n '/^{/,$p')

# Engine-specific endPointType naming. MySQL/MariaDB use INTERNAL_VIP;
# PostgreSQL uses plain INTERNAL (no VIP suffix). Match both.
EP_DOMAIN=$(printf '%s' "$NET_JSON" \
  | jq -r '.endPoints[]? | select(.endPointType=="INTERNAL_VIP" or .endPointType=="INTERNAL") | .domain' \
  | head -1)

if [[ -z "$EP_DOMAIN" ]]; then
  echo "[wait] ERROR: INTERNAL_VIP endpoint not found. describe-db-instance-network output:" >&2
  printf '%s\n' "$NET_JSON" >&2
  exit 1
fi

ENGINE_DEFAULT_PORT=$(yq -r ".parameters.${ENGINE}.db_port // 3306" "$SCENARIO_DIR/scenario.yaml")
DB_PORT=$(nhncloud "$CLI_PREFIX" describe-db-instances -o json 2>/dev/null \
          | jq -r --arg id "$INSTANCE_ID" --arg fb "$ENGINE_DEFAULT_PORT" \
              '.dbInstances[] | select(.dbInstanceId==$id) | (.dbPort // ($fb | tonumber))')
[[ -z "$DB_PORT" || "$DB_PORT" == "null" ]] && DB_PORT="$ENGINE_DEFAULT_PORT"

ENDPOINT="$EP_DOMAIN:$DB_PORT"

if [[ -z "$ENDPOINT" || "$ENDPOINT" != *:* ]]; then
  echo "[wait] ERROR: endpoint resolution returned unexpected output: '$ENDPOINT'" >&2
  exit 1
fi

echo "ENDPOINT=$ENDPOINT" >> "$SCENARIO_DIR/state.env"
echo "[wait] $INSTANCE_ID → AVAILABLE @ $ENDPOINT"

# Create the sysbench schema. NHN Cloud RDS does not auto-create a
# default database from the master user name; sysbench prepare needs
# the schema to exist with master-user grants. Idempotent enough — if
# this fails because the schema already exists, treat as fatal (we
# just provisioned a fresh instance, so the schema must not exist).
# Ref: https://docs.nhncloud.com/en/Database/ (per-engine public API)#db-스키마-생성
# Helper: poll describe-db-instances until the instance's progressStatus
# returns to NONE (no in-flight job). NHN rejects subsequent mutations
# with `API error 7008: Job is fail to register` while a job is in flight.
wait_progress_none() {
  local label="$1"
  for i in $(seq 1 60); do  # 60 * 5s = 5 min
    local ps
    ps=$(nhncloud "$CLI_PREFIX" describe-db-instances -o json 2>/dev/null \
         | jq -r --arg id "$INSTANCE_ID" \
             '.dbInstances[] | select(.dbInstanceId==$id) | .progressStatus // "NONE"')
    if [[ "$ps" == "NONE" || -z "$ps" ]]; then
      echo "[wait] $label progressStatus=NONE (idle)"
      return 0
    fi
    printf '.'
    sleep 5
  done
  echo
  echo "[wait] WARNING: $label progressStatus never returned to NONE; proceeding anyway." >&2
}

# Retry helper: NHN's RDS rejects mutations with API error 7008
# ("Job is fail to register") if a previous job is still resolving on
# their side. wait_progress_none above waits for progressStatus=NONE
# but there's still a brief window where the next call races and
# 7008s. Retry up to 6 × 15s before giving up.
rds_retry_7008() {
  local label="$1"; shift
  for i in $(seq 1 6); do
    local out
    if out=$("$@" 2>&1); then
      echo "[wait] $label OK"
      return 0
    fi
    if echo "$out" | grep -q '7008'; then
      echo "[wait] $label hit 7008 — retry $i/6 in 15s..."
      sleep 15
      continue
    fi
    echo "[wait] $label failed: $out" >&2
    return 1
  done
  echo "[wait] $label still 7008-ing after 6 retries" >&2
  return 1
}

DB_SCHEMA=$(yq '.parameters.database' "$SCENARIO_DIR/scenario.yaml")
if [[ "$ENGINE" == "postgresql" ]]; then
  # PostgreSQL: --database-name was passed at create-db-instance time, so
  # the schema already exists at provisioning. PG's CLI also exposes
  # `create-database` (not `create-db-schema`), and the loadgen connects
  # to the existing DB via DB_NAME env. No schema-create call needed.
  echo "[wait] postgresql: schema '$DB_SCHEMA' was created at provisioning (--database-name) — skipping create-db-schema"
else
  echo "[wait] creating db schema '$DB_SCHEMA' on $INSTANCE_ID..."
  rds_retry_7008 "schema-create" \
    nhncloud "$CLI_PREFIX" create-db-schema \
      --db-instance-identifier "$INSTANCE_ID" \
      --db-schema-name "$DB_SCHEMA"
  echo "[wait] schema '$DB_SCHEMA' created; waiting for idle..."
  wait_progress_none "post-schema"
fi

# Create an application DB user with DDL authority. The master user
# created at instance-creation time is host-locked (likely VPC-only)
# and gets `Access denied for user 'X'@'%' to database 'X'` from
# external clients even after the schema exists. NHN's intended model
# is to add a separate dbUser with explicit host pattern + grants.
# Ref: https://docs.nhncloud.com/en/Database/ (per-engine public API)#db-사용자-생성하기
APP_DB_USER=$(yq -r '.parameters.app_db_user // "appbench"' "$SCENARIO_DIR/scenario.yaml")
: "${APP_DB_PASSWORD:?must be set by run.sh}"
if [[ "$ENGINE" == "postgresql" ]]; then
  # NHN's rds-postgresql CLI doesn't expose --authority-type, and round-2
  # confirmed create-db-user is unusable end-to-end. The loadgen's
  # PHASE=grants step does CREATE USER + GRANT directly from inside the
  # cluster via pgx. Persist APP_DB_USER so 04-deploy can pass it through.
  echo "[wait] postgresql: skipping create-db-user (loadgen grants phase will do it)"
  echo "APP_DB_USER=$APP_DB_USER" >> "$SCENARIO_DIR/state.env"
else
  echo "[wait] creating db user '$APP_DB_USER@%' with DDL authority..."
  rds_retry_7008 "user-create" \
    nhncloud "$CLI_PREFIX" create-db-user \
      --db-instance-identifier "$INSTANCE_ID" \
      --db-user-name "$APP_DB_USER" \
      --db-password "$APP_DB_PASSWORD" \
      --host "%" \
      --authority-type DDL
  echo "APP_DB_USER=$APP_DB_USER" >> "$SCENARIO_DIR/state.env"
fi
if [[ "$ENGINE" != "postgresql" ]]; then
  echo "[wait] db user '$APP_DB_USER' created; waiting for idle..."
  wait_progress_none "post-user"
fi

# Re-verify the instance is back to AVAILABLE — create-db-user briefly
# transitions the instance through a CREATING_USER progressStatus on
# NHN's side.  The loadgen Pod inside the cluster will perform its own
# TCP probe; no runner-side probe needed for internal-only mode.
echo "[wait] re-confirming AVAILABLE after user creation..."
DEADLINE=$(( $(date +%s) + 5 * 60 ))
last_status=""
while :; do
  if [[ $(date +%s) -ge $DEADLINE ]]; then
    echo "[wait] ERROR: re-confirm timeout (last status=$last_status)" >&2
    exit 1
  fi
  status=$(nhncloud "$CLI_PREFIX" describe-db-instances -o json 2>/dev/null \
    | jq -r --arg id "$INSTANCE_ID" \
        '.dbInstances[] | select(.dbInstanceId==$id) | .dbInstanceStatus' \
    | head -1)
  if [[ "$status" != "$last_status" ]]; then
    echo "[wait] status=$status"
    last_status="$status"
  fi
  [[ "$status" == "AVAILABLE" ]] && break
  sleep 10
done

echo "[wait] use_public_access=false (internal-only mode) — skipping runner-side TCP probe; loadgen Pod will probe from inside the cluster"
