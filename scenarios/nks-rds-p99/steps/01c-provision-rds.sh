#!/usr/bin/env bash
# steps/01c-provision-rds.sh  (scenario: nks-rds-p99)
#
# Calls `nhncloud "$CLI_PREFIX" create-db-instance` with values from
# scenario.yaml + env vars set by run.sh.
#
# Nearly identical to rds-mysql-p99-internal/steps/01-provision.sh.
# use_public_access is always false in this scenario — the cluster
# nodes reach the INTERNAL_VIP, so no public endpoint is needed.
#
# Ref: https://docs.nhncloud.com/en/Database/RDS%20for%20MySQL/en/public-api/ (Korean: rds-mysql-v4.0 spec)#db-인스턴스-생성하기
# Ref: nhn-cloud-cli/cmd/rds_mysql.go (createDBInstanceCmd flags)

set -euo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"
: "${INSTANCE_NAME:?must be set by run.sh}"
: "${DB_PASSWORD:?must be set by run.sh}"

# shellcheck disable=SC1091
source "$SCENARIO_DIR/state.env" 2>/dev/null || true

: "${ENGINE:?state.env missing ENGINE — run.sh did not export it}"
CLI_PREFIX=$(yq -r ".parameters.${ENGINE}.cli_prefix" "$SCENARIO_DIR/scenario.yaml")

# Pull all parameters in one yq call to avoid repeated parses.
PARAMS=$(yq '.parameters' "$SCENARIO_DIR/scenario.yaml" -o=json)
flavor_id=$(yq -r ".parameters.${ENGINE}.flavor_id // .parameters.flavor_id" "$SCENARIO_DIR/scenario.yaml")
storage=$(echo "$PARAMS"              | jq -r .storage_size_gb)
db_version=$(yq -r ".parameters.${ENGINE}.db_version // .parameters.db_version" "$SCENARIO_DIR/scenario.yaml")
availability_zone=$(echo "$PARAMS"    | jq -r .availability_zone)
subnet_id=$(echo "$PARAMS"            | jq -r .subnet_id)
parameter_group_id=$(yq -r ".parameters.${ENGINE}.db_parameter_group_id // .parameters.db_parameter_group_id" "$SCENARIO_DIR/scenario.yaml")
dbuser=$(echo "$PARAMS"               | jq -r .user)
dbname=$(echo "$PARAMS"               | jq -r .database)
use_public_access=$(echo "$PARAMS"    | jq -r '.use_public_access // false')

# Prefer the per-run MySQL SG created by 01a-create-sg.sh (scoped to
# the subnet CIDR — NKS nodes reach RDS via INTERNAL_VIP, no public /32
# ingress needed).  Fallback to a yaml-pinned SG only for ad-hoc manual
# runs that skip 01a.  BENCH_SG_ID is the legacy single-SG name; kept
# as a fallback for any stale state.env from before the B.3 split.
db_security_group_id="${BENCH_MYSQL_SG_ID:-${BENCH_SG_ID:-$(echo "$PARAMS" | jq -r '.db_security_group_id // empty')}}"

echo "[provision] creating $INSTANCE_NAME (flavor=$flavor_id storage=${storage}GB az=$availability_zone)"

# Note: --master-user-password is passed via env-substituted argv. The
# password is exported by run.sh; we never `echo` it and never log it.
# The CLI does not currently support reading the password from stdin or
# a file, so argv is the only available path. If/when --master-user-
# password-file is added upstream, switch to that.
# Use -o json for scriptable parsing. CLI emits {"jobId":"...","dbInstanceName":"..."}.
# Provisioning is async — the instance is not yet listable when create returns.
SG_ARGS=()
if [[ -n "$db_security_group_id" ]]; then
  SG_ARGS=(--db-security-group-ids "$db_security_group_id")
fi

# Engine-specific flag mapping. PostgreSQL CLI renamed almost every flag
# (--db-instance-name instead of -identifier, --db-version instead of
# --engine-version, --db-user-name / --db-password instead of --master-*),
# and requires --database-name at create time (mysql/mariadb create DB
# lazily on first connect). Pre-build the argv per engine so the retry
# loop body stays one nhncloud invocation.
if [[ "$ENGINE" == "postgresql" ]]; then
  # PG default --storage-type is "SSD" but NHN PG's storage catalog only
  # exposes "General SSD" / "General HDD" — passing the CLI default lands
  # an opaque API 500. Always set --storage-type from scenario.yaml.
  pg_storage_type=$(yq -r '.parameters.postgresql.storage_type // "General SSD"' "$SCENARIO_DIR/scenario.yaml")
  CREATE_ARGS=(
    --db-instance-name        "$INSTANCE_NAME"
    --db-flavor-id            "$flavor_id"
    --db-version              "$db_version"
    --db-user-name            "$dbuser"
    --db-password             "$DB_PASSWORD"
    --database-name           "$dbname"
    --subnet-id               "$subnet_id"
    --availability-zone       "$availability_zone"
    --db-parameter-group-id   "$parameter_group_id"
    --allocated-storage       "$storage"
    --storage-type            "$pg_storage_type"
  )
else
  CREATE_ARGS=(
    --db-instance-identifier "$INSTANCE_NAME"
    --db-flavor-id           "$flavor_id"
    --engine-version         "$db_version"
    --master-username        "$dbuser"
    --master-user-password   "$DB_PASSWORD"
    --subnet-id              "$subnet_id"
    --availability-zone      "$availability_zone"
    --db-parameter-group-id  "$parameter_group_id"
    --allocated-storage      "$storage"
  )
fi

# Disable -e around the create call so we can surface the CLI's stderr
# on failure. 7008 ("Job is fail to register") happens when an SG rule
# we just added is still propagating; retry up to 6 * 15s.
set +e
CREATE_JSON=""
CREATE_RC=1
for i in $(seq 1 6); do
  CREATE_JSON=$(nhncloud "$CLI_PREFIX" create-db-instance \
    "${CREATE_ARGS[@]}" \
    "${SG_ARGS[@]}" \
    -o json 2>&1)
  CREATE_RC=$?
  if [[ ( "$ENGINE" == "mariadb" || "$ENGINE" == "postgresql" ) && $CREATE_RC -eq 0 ]]; then
    # rds-mariadb / rds-postgresql CLI prefixes some `-o json` SUCCESS
    # outputs with a "Fetching ..." line before the JSON body. Strip only
    # on success — on failure the same `2>&1` capture holds the actual
    # error message whose first character is not '{' or 'Job ID:'.
    CREATE_JSON=$(printf '%s\n' "$CREATE_JSON" | sed -n '/^{/,$p; /^Job ID:/,$p')
  fi
  if [[ $CREATE_RC -eq 0 ]]; then
    break
  fi
  if echo "$CREATE_JSON" | grep -q '7008'; then
    echo "[provision] 7008 (in-flight job) — retry $i/6 in 15s..." >&2
    sleep 15
    continue
  fi
  break
done
set -e

# Always dump the raw create response BEFORE any error gate — failure
# diagnostics need the original stderr-merged body. Round-3 originally
# dumped after the gate which left silent failures with no trail.
DEBUG_DIR="$SCENARIO_DIR/reports"; mkdir -p "$DEBUG_DIR"
printf '%s\n' "$CREATE_JSON" > "$DEBUG_DIR/${RUN_TS}.create-db-instance.raw"

if [[ $CREATE_RC -ne 0 ]]; then
  echo "[provision] ERROR: create-db-instance exited $CREATE_RC. Raw response:" >&2
  echo "           $DEBUG_DIR/${RUN_TS}.create-db-instance.raw" >&2
  printf '%s\n' "$CREATE_JSON" >&2
  exit "$CREATE_RC"
fi

JOB_ID=$(printf '%s' "$CREATE_JSON" | jq -r '.jobId // empty' 2>/dev/null || true)
if [[ -z "$JOB_ID" ]]; then
  # JSON style (workspace CLI w/ commit caaef5d): {"jobId":"…"}.
  # Trailing `|| true` because under `set -o pipefail`, a grep with no
  # match propagates exit 1 → set -e kills the script before reaching
  # the banner-style fallback below.
  JOB_ID=$(printf '%s' "$CREATE_JSON" \
    | { grep -oE '"jobId"[[:space:]]*:[[:space:]]*"[^"]+"' \
        || true; } \
    | head -1 | sed -E 's/.*"jobId"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/')
fi
if [[ -z "$JOB_ID" ]]; then
  # Banner style (brew CLI lags caaef5d, ignores -o json on create):
  #   "DB instance creation initiated."
  #   "Job ID: <uuid>"
  # awk doesn't return non-zero on no-match, so this pipeline is
  # safe under set -e even when the banner shape changes.
  JOB_ID=$(printf '%s\n' "$CREATE_JSON" \
    | awk '/^Job ID:/ { print $NF; exit }')
fi
if [[ -z "$JOB_ID" ]]; then
  echo "[provision] ERROR: create did not return a jobId. Raw response saved at:" >&2
  echo "           $DEBUG_DIR/${RUN_TS}.create-db-instance.raw" >&2
  printf '%s\n' "$CREATE_JSON" >&2
  exit 1
fi
echo "JOB_ID=$JOB_ID" >> "$SCENARIO_DIR/state.env"
echo "[provision] create accepted — JOB_ID=$JOB_ID; polling for instance row to appear..."

# Retry loop: poll describe-db-instances until $INSTANCE_NAME appears.
# The instance enters BEFORE_CREATE → STORAGE_CREATING → ... → AVAILABLE
# (per spec dbInstanceStatus enum). We just wait until the row is
# listable (any status) so 02a-wait-rds can pick it up cleanly.
INSTANCE_ID=""
for i in $(seq 1 60); do  # 60 × 5s = 5 min
  INSTANCE_ID=$(nhncloud "$CLI_PREFIX" describe-db-instances -o json 2>/dev/null \
    | jq -r --arg n "$INSTANCE_NAME" \
        '.dbInstances[] | select(.dbInstanceName==$n) | .dbInstanceId' \
    | head -1)
  if [[ -n "$INSTANCE_ID" ]]; then
    break
  fi
  printf '.'; sleep 5
done
echo

if [[ -z "$INSTANCE_ID" ]]; then
  echo "[provision] ERROR: instance $INSTANCE_NAME never appeared in describe-db-instances after 5min" >&2
  echo "           (Job $JOB_ID may have failed validation; inspect via API directly)" >&2
  exit 1
fi

echo "INSTANCE_ID=$INSTANCE_ID" >> "$SCENARIO_DIR/state.env"
echo "[provision] resolved — INSTANCE_ID=$INSTANCE_ID"
