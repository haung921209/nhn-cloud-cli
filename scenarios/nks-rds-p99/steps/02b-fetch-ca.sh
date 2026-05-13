#!/usr/bin/env bash
# steps/02b-fetch-ca.sh
#
# Triggers the v4.0 cert-export API (POST /v4.0/db-instances/{id}/
# certificates/upload) which writes ca.pem to NHN Object Storage,
# then pulls it down via `nhncloud object-storage cp` so the loadgen
# can do real TLS verification (not InsecureSkipVerify).
#
# Ref: https://docs.nhncloud.com/en/Database/RDS%20for%20MySQL/en/public-api/ (Korean: rds-mysql-v4.0 spec)#인증서-파일-내보내기 (line 2859)
#
# This is necessary because:
#   1. The brew sysbench / libmysqlclient on macOS doesn't trust the
#      NHN private CA, so without skip-verify it always fails
#      with error 2026.
#   2. MySQL 8.0+ in TLS-required mode rejects plaintext, so disabling
#      TLS at the client is not always possible.
#
# Output:
#   $SCENARIO_DIR/certs/$RUN_TS-ca.pem
#   state.env: CA_FILE=<absolute path>

set -euo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"
: "${RUN_TS:?must be set by run.sh}"

# shellcheck disable=SC1091
source "$SCENARIO_DIR/state.env"
: "${INSTANCE_ID:?state.env missing INSTANCE_ID}"

ENGINE="${ENGINE:-$(yq -r '.parameters.engine // "mysql"' "$SCENARIO_DIR/scenario.yaml")}"
CLI_PREFIX=$(yq -r ".parameters.${ENGINE}.cli_prefix" "$SCENARIO_DIR/scenario.yaml")
CERT_TIMEOUT_S=$(yq -r ".parameters.${ENGINE}.cert_export_timeout_s // 300" "$SCENARIO_DIR/scenario.yaml")
echo "[ca] engine=$ENGINE timeout=${CERT_TIMEOUT_S}s"

# PostgreSQL: the public API spec (https://docs.nhncloud.com/en/Database/RDS%20for%20PostgreSQL/en/public-api/)
# does not document a /certificates/upload endpoint, AND the API host is
# `kr1-rds-postgres.api.nhncloudservice.com` (not -postgresql) so the v4.0-
# style URL constructed below resolves to NXDOMAIN. Short-circuit to the
# same soft-fall path MariaDB uses: emit empty CA_FILE so the loadgen's
# pgx driver falls back to InsecureSkipVerify+VerifyPeerCertificate (the
# round-2 PG TLS workaround). One WARNING line, exit 0.
if [[ "$ENGINE" == "postgresql" ]]; then
  echo "[ca] WARNING: postgresql has no public cert-export endpoint — falling back to skip-verify TLS" >&2
  echo 'CA_FILE=""' >> "$SCENARIO_DIR/state.env"
  exit 0
fi

# Read tenant + OS credentials from ~/.nhncloud/credentials.
CRED="$HOME/.nhncloud/credentials"
[[ -r "$CRED" ]] || { echo "[ca] ERROR: $CRED not readable" >&2; exit 1; }

read_cred() { grep -E "^${1}[[:space:]]*=" "$CRED" | head -1 | sed -E 's/^[^=]*=[[:space:]]*//' | tr -d '"'; }
ACCESS_KEY=$(read_cred access_key_id)
SECRET=$(read_cred secret_access_key)
APPKEY_PRIMARY_KEY=$(yq -r ".parameters.${ENGINE}.cred_file_appkey // \"rds_app_key\"" "$SCENARIO_DIR/scenario.yaml")
APPKEY_ALT_KEY=$(yq -r ".parameters.${ENGINE}.cred_file_appkey_alt // \"\"" "$SCENARIO_DIR/scenario.yaml")
APPKEY=$(read_cred "$APPKEY_PRIMARY_KEY")
if [[ -z "$APPKEY" && -n "$APPKEY_ALT_KEY" ]]; then
  APPKEY=$(read_cred "$APPKEY_ALT_KEY")
fi
OBS_TENANT=$(read_cred obs_tenant_id)
OS_USER=$(read_cred username)
OS_PWD=$(read_cred api_password)
REGION=$(read_cred region)
[[ -z "$REGION" ]] && REGION=kr1
REGION=$(echo "$REGION" | tr '[:upper:]' '[:lower:]')

for v in ACCESS_KEY SECRET APPKEY OBS_TENANT OS_USER OS_PWD; do
  [[ -n "${!v}" ]] || { echo "[ca] ERROR: missing $v in credentials" >&2; exit 1; }
done

# Step 1: get a Bearer token from /oauth2/token/create. Auth uses
# HTTP Basic with `<access_key_id>:<secret_access_key>` (matches the
# SDK's nhncloud/auth/bearer_auto.go); form-data for credentials does
# not work and returns 400.
echo "[ca] requesting Bearer token..."
BASIC=$(printf '%s:%s' "$ACCESS_KEY" "$SECRET" | base64 | tr -d '\n')
TOKEN_JSON=$(curl -fsS -X POST https://oauth.api.nhncloudservice.com/oauth2/token/create \
  -H "Authorization: Basic $BASIC" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'grant_type=client_credentials')
TOKEN=$(echo "$TOKEN_JSON" | jq -r '.access_token // empty')
[[ -n "$TOKEN" ]] || { echo "[ca] ERROR: token issue failed: $TOKEN_JSON" >&2; exit 1; }

# Step 2: ensure the OS container exists (mb requires obs:// prefix
# and is idempotent — exit 0 on existing container).
CONTAINER="${BENCH_CERT_CONTAINER:-bench-rds-certs}"
nhncloud object-storage mb "obs://$CONTAINER" >/dev/null 2>&1 || true
# Verify by listing — if container truly doesn't exist, fail loud now
# rather than letting cert-export rejection surface it.
if ! nhncloud object-storage ls 2>/dev/null | grep -qE "^${CONTAINER}\b"; then
  echo "[ca] ERROR: OS container '$CONTAINER' could not be created/found" >&2
  exit 1
fi

# Step 3: trigger cert export to OS.
# NB: NHN's cert-export concatenates objectPath + filename, but with
# a trailing slash on objectPath the result becomes "<path>//ca.pem"
# (double slash). Drop the trailing slash so the file lands cleanly
# at "<path>/ca.pem".
OBJ_PATH="$CLI_PREFIX-p99/$RUN_TS"
echo "[ca] requesting cert export to obs://$CONTAINER/$OBJ_PATH ..."
EXPORT_RESP=$(curl -fsS -X POST \
  "https://${REGION}-${CLI_PREFIX}.api.nhncloudservice.com/v4.0/db-instances/$INSTANCE_ID/certificates/upload" \
  -H "X-TC-APP-KEY: $APPKEY" \
  -H "X-NHN-AUTHORIZATION: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "$(jq -n --arg t "$OBS_TENANT" --arg u "$OS_USER" --arg p "$OS_PWD" \
              --arg c "$CONTAINER" --arg o "$OBJ_PATH" \
        '{certificateTypes: ["CA_FILE"], tenantId: $t, username: $u,
          password: $p, targetContainer: $c, objectPath: $o}')")
EXPORT_JOB=$(echo "$EXPORT_RESP" | jq -r '.jobId // empty')
if [[ -z "$EXPORT_JOB" ]]; then
  echo "[ca] WARNING: ${ENGINE} cert-export rejected — falling back to skip-verify TLS" >&2
  echo "[ca] WARNING:   response: $EXPORT_RESP" >&2
  CA_FILE=""
  echo "CA_FILE=$CA_FILE" >> "$SCENARIO_DIR/state.env"
  exit 0
fi
echo "[ca] export job $EXPORT_JOB submitted; polling Object Storage for ca.pem..."

# Step 4: poll the OS listing for ca.pem to land. NHN's cert export
# is async; in practice it lands within ~10s but we allow up to 2min.
LOCAL_DIR="$SCENARIO_DIR/certs"
mkdir -p "$LOCAL_DIR"
LOCAL_CA="$LOCAL_DIR/$RUN_TS-ca.pem"

CA_OS_KEY=""
CA_FILE=""
# pipefail+grep would abort the loop on the first empty-list iteration;
# disable pipefail just for the polling block.
set +o pipefail
PROBE=0
DEADLINE=$(( $(date +%s) + CERT_TIMEOUT_S ))
while :; do
  if [[ $(date +%s) -ge $DEADLINE ]]; then
    echo
    echo "[ca] WARNING: ${ENGINE} cert-export did not produce a file in ${CERT_TIMEOUT_S}s — falling back to skip-verify TLS" >&2
    CA_FILE=""
    break
  fi
  PROBE=$(( PROBE + 1 ))
  # Recursive list — NHN sometimes nests ca.pem under multipart segments,
  # so we look for the top-level ca.pem (NOT the .../ca.pem/0001 segment).
  CA_OS_KEY=$(nhncloud object-storage ls "obs://$CONTAINER/${OBJ_PATH}" -r 2>/dev/null \
              | awk '{print $1}' \
              | grep -E '/ca\.pem$' \
              | head -1)
  if [[ -n "$CA_OS_KEY" ]]; then
    echo "[ca] ca.pem present at obs://$CONTAINER/$CA_OS_KEY (after ${PROBE} probe(s))"
    break
  fi
  printf '.'; sleep 5
done
set -o pipefail

if [[ -n "$CA_OS_KEY" ]]; then
  # Step 5: download ca.pem to local using the exact key we just found.
  nhncloud object-storage cp "obs://$CONTAINER/$CA_OS_KEY" "$LOCAL_CA" >/dev/null
  if [[ -s "$LOCAL_CA" ]]; then
    CA_FILE="$LOCAL_CA"
    echo "[ca] ca.pem -> $CA_FILE ($(wc -c < "$CA_FILE") bytes)"
  else
    echo "[ca] WARNING: download produced empty file $LOCAL_CA — falling back to skip-verify TLS" >&2
    CA_FILE=""
  fi
fi

{
  echo "CA_FILE=$CA_FILE"
  echo "BENCH_CERT_CONTAINER=$CONTAINER"
  echo "BENCH_CERT_OBJ_KEY=$CA_OS_KEY"
} >> "$SCENARIO_DIR/state.env"
