#!/usr/bin/env bash
# steps/99-teardown.sh  (scenario: nks-rds-p99)
#
# Deletes all billable resources in dependency order.
# NKS-specific additions vs rds-mysql-p99-internal:
#   - Delete the NKS cluster first (nodes are NKS-managed; no separate
#     Compute instance or Floating-IP to delete in this scenario).
#   - Wait for cluster to fully disappear before removing SGs/keypair
#     that the cluster nodes still reference.
#
# Order:
#   1. Delete the NKS cluster (helpers/nks-control delete)
#   2. Wait for cluster DELETED / gone (best-effort, 30m budget)
#   3. Delete the keypair (after nodes are gone — nodes reference it)
#   4. Delete the MySQL instance
#   5. Wait for MySQL to disappear
#   6. Delete the MySQL SG  (after MySQL — MySQL references it)
#   7. Delete the Nodes SG  (after cluster — nodes reference it)
#   8. Remove the cert artifact from Object Storage
#
# Idempotent: every section checks whether the resource ID is set and
# treats "already gone" (non-zero exit, 404-like) as success.
# Honours KEEP=1 to skip all deletions (useful for debugging).
#
# Ref: https://docs.nhncloud.com/en/Database/RDS%20for%20MySQL/en/public-api/ (Korean: rds-mysql-v4.0 spec)#db-인스턴스-삭제하기
# Ref: https://docs.nhncloud.com/en/Database/RDS%20for%20MySQL/en/public-api/ (Korean: rds-mysql-v4.0 spec)#db-보안-그룹-삭제하기
# Ref: https://docs.nhncloud.com/en/Container/NKS/en/public-api/  (cluster delete + GET polling)

set -uo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"

if [[ "${KEEP:-0}" == "1" ]]; then
  echo "[teardown] KEEP=1, skipping delete"
  exit 0
fi

# shellcheck disable=SC1091
source "$SCENARIO_DIR/state.env" 2>/dev/null || true

: "${ENGINE:=mysql}"   # default mysql so legacy state.env from round 3 still works
CLI_PREFIX=$(yq -r ".parameters.${ENGINE}.cli_prefix // \"rds-mysql\"" "$SCENARIO_DIR/scenario.yaml")
# Brew CLI v0.7.18 lacks delete-db-security-group for rds-postgresql; route
# only PG's SG-verbs through the workspace build (NHNCLOUD_WS_BIN is exported
# by run.sh and also persisted to state.env). Everything else uses brew.
if [[ "$ENGINE" == "postgresql" ]]; then
  SG_CLI="${NHNCLOUD_WS_BIN:-nhncloud}"   # tolerate missing on legacy state.env (no-op for non-PG runs)
else
  SG_CLI="nhncloud"
fi
echo "[teardown] engine=$ENGINE cli_prefix=$CLI_PREFIX sg_cli=$SG_CLI"

NKS_CTL="$SCENARIO_DIR/helpers/nks-control/nks-control"

# nks-control reads identity creds from env. Surface them from the
# credentials file the same way 01d/02c/03 do.
CRED="${HOME}/.nhncloud/credentials"
if [[ -f "$CRED" ]]; then
  read_cred() { grep -E "^${1}[[:space:]]*=" "$CRED" | head -1 | sed -E 's/^[^=]*=[[:space:]]*//' | tr -d '"'; }
  export NHNCLOUD_USERNAME=$(read_cred username)
  export NHNCLOUD_PASSWORD=$(read_cred api_password)
  export NHNCLOUD_TENANT_ID=$(read_cred tenant_id)
  _r=$(read_cred region); export NHNCLOUD_REGION="${_r:-kr1}"
fi

# ─── 1. Delete the NKS cluster ──────────────────────────────────────────────
# In existing-cluster mode (Phase 4) the cluster is owned outside the
# scenario's lifecycle (console-created, long-lived) — never delete it.
if [[ "${EXISTING_CLUSTER_MODE:-0}" == "1" ]]; then
  echo "[teardown] EXISTING_CLUSTER_MODE=1 — keeping cluster $CLUSTER_ID alive (owned outside scenario)"
elif [[ -n "${CLUSTER_ID:-}" ]]; then
  echo "[teardown] deleting NKS cluster $CLUSTER_ID ..."
  if ! out=$("$NKS_CTL" delete --cluster-id "$CLUSTER_ID" 2>&1); then
    echo "[teardown] cluster delete: $out" >&2
    echo "[teardown] WARNING: cluster $CLUSTER_ID may need manual cleanup" >&2
  else
    echo "[teardown] $out"
  fi
else
  echo "[teardown] no CLUSTER_ID — NKS cluster step skipped"
fi

# ─── 2. Wait for cluster to fully disappear (best-effort) ───────────────────
if [[ "${EXISTING_CLUSTER_MODE:-0}" == "1" ]]; then
  : # nothing to wait for — cluster stays
elif [[ -n "${CLUSTER_ID:-}" ]]; then
  echo "[teardown] waiting for cluster $CLUSTER_ID to reach DELETED (30m budget, best-effort)..."
  "$NKS_CTL" wait --cluster-id "$CLUSTER_ID" --for-state DELETED --timeout 30m \
    >/dev/null 2>&1 \
    || echo "[teardown] cluster wait-DELETED timed out or errored — proceeding anyway" >&2
  echo "[teardown] cluster wait done."
fi

# ─── 3. Delete the keypair (after cluster nodes are gone) ───────────────────
# Skipped in existing-cluster mode: the keypair was created outside the
# scenario and is referenced by long-lived nodes.
if [[ "${EXISTING_CLUSTER_MODE:-0}" == "1" ]]; then
  echo "[teardown] EXISTING_CLUSTER_MODE=1 — keypair kept (owned outside scenario)"
elif [[ -n "${BENCH_KEYPAIR_NAME:-}" ]]; then
  echo "[teardown] deleting keypair $BENCH_KEYPAIR_NAME ..."
  nhncloud compute delete-key-pair --key-name "$BENCH_KEYPAIR_NAME" \
    >/dev/null 2>&1 \
    || echo "[teardown] keypair delete failed (already gone?)" >&2
fi

# ─── 4. Delete the MySQL instance ───────────────────────────────────────────
# Use INSTANCE_ID only — INSTANCE_NAME is populated by run.sh BEFORE the
# provision step runs, so an early failure (e.g. SG-creation 500) would
# otherwise trip a delete against a non-existent instance and abort cleanup.
# CLI's `delete-db-instance` does not accept `--yes` (different from create).
if [[ -n "${INSTANCE_ID:-}" ]]; then
  echo "[teardown] deleting MySQL instance $INSTANCE_ID ..."
  delete_ok=0
  for i in $(seq 1 6); do
    if out=$(nhncloud "$CLI_PREFIX" delete-db-instance \
               --db-instance-identifier "$INSTANCE_ID" 2>&1); then
      echo "[teardown] $out"
      delete_ok=1
      break
    fi
    if echo "$out" | grep -q '7008'; then
      echo "[teardown] 7008 (in-flight job) — retry $i/6 in 15s..."
      sleep 15
      continue
    fi
    echo "[teardown] delete failed: $out" >&2
    break
  done
  [[ "$delete_ok" -eq 1 ]] || \
    echo "[teardown] WARNING: MySQL instance $INSTANCE_ID may need manual cleanup" >&2
else
  echo "[teardown] no INSTANCE_ID — MySQL step skipped (instance was never created)"
fi

# ─── 5. Wait for MySQL to disappear ─────────────────────────────────────────
if [[ -n "${INSTANCE_ID:-}" ]]; then
  echo -n "[teardown] waiting for MySQL $INSTANCE_ID to disappear (10min budget)..."
  for _ in $(seq 1 60); do
    if ! nhncloud "$CLI_PREFIX" describe-db-instances 2>/dev/null \
         | grep -q "$INSTANCE_ID"; then
      echo
      echo "[teardown] MySQL gone."
      break
    fi
    printf '.'; sleep 10
  done
fi

# ─── 6. Delete the MySQL SG (must be after MySQL) ───────────────────────────
# Accept legacy BENCH_SG_ID in case state.env predates the B.3 split.
MYSQL_SG_ID="${BENCH_MYSQL_SG_ID:-${BENCH_SG_ID:-}}"
if [[ -n "$MYSQL_SG_ID" ]]; then
  echo "[teardown] deleting MySQL SG $MYSQL_SG_ID ..."
  for i in $(seq 1 6); do
    if "$SG_CLI" "$CLI_PREFIX" delete-db-security-group \
         --db-security-group-identifier "$MYSQL_SG_ID" >/dev/null 2>&1; then
      echo "[teardown] MySQL SG $MYSQL_SG_ID deleted."
      break
    fi
    echo "[teardown] MySQL SG delete failed, retry $i/6 in 10s..."
    sleep 10
  done
fi

# ─── 7. Delete the Nodes SG (must be after cluster nodes are gone) ──────────
# BENCH_COMPUTE_SG_ID is written by 01a-create-sg.sh as the node-group SG.
if [[ -n "${BENCH_COMPUTE_SG_ID:-}" ]]; then
  echo "[teardown] deleting Nodes SG $BENCH_COMPUTE_SG_ID ..."
  for i in $(seq 1 6); do
    if nhncloud network delete-security-group \
         --group-id "$BENCH_COMPUTE_SG_ID" >/dev/null 2>&1; then
      echo "[teardown] Nodes SG $BENCH_COMPUTE_SG_ID deleted."
      break
    fi
    echo "[teardown] Nodes SG delete failed, retry $i/6 in 10s..."
    sleep 10
  done
fi

# ─── 8. Remove the cert object from Object Storage (keep container) ─────────
if [[ -n "${BENCH_CERT_CONTAINER:-}" && -n "${BENCH_CERT_OBJ_KEY:-}" ]]; then
  echo "[teardown] removing cert artifact obs://$BENCH_CERT_CONTAINER/$BENCH_CERT_OBJ_KEY ..."
  nhncloud object-storage rm \
    "obs://$BENCH_CERT_CONTAINER/$BENCH_CERT_OBJ_KEY" \
    >/dev/null 2>&1 || true
fi

exit 0
