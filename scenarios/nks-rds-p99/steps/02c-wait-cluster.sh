#!/usr/bin/env bash
# steps/02c-wait-cluster.sh
#
# Polls helpers/nks-control wait until the cluster status reaches an
# "ACTIVE / *_COMPLETE" terminal state. NHN's NKS provisioning typically
# takes 25-40 minutes for a 1-node cluster; the helper's default timeout
# is 60m which we keep.

set -euo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"

# shellcheck disable=SC1091
source "$SCENARIO_DIR/state.env"
: "${CLUSTER_ID:?state.env did not carry CLUSTER_ID — create step missing}"

CRED="${HOME}/.nhncloud/credentials"
read_cred() { grep -E "^${1}[[:space:]]*=" "$CRED" | head -1 | sed -E 's/^[^=]*=[[:space:]]*//' | tr -d '"'; }
export NHNCLOUD_USERNAME=$(read_cred username)
export NHNCLOUD_PASSWORD=$(read_cred api_password)
export NHNCLOUD_TENANT_ID=$(read_cred tenant_id)
NHNCLOUD_REGION=$(read_cred region); export NHNCLOUD_REGION="${NHNCLOUD_REGION:-kr1}"

NKS_BIN="$SCENARIO_DIR/helpers/nks-control/nks-control"

echo "[wait-cluster] waiting for $CLUSTER_ID to reach ACTIVE (timeout 60m)..."
"$NKS_BIN" wait --cluster-id "$CLUSTER_ID" --for-state ACTIVE --timeout 60m
echo "[wait-cluster] $CLUSTER_ID is ACTIVE"
