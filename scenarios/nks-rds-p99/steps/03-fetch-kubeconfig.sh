#!/usr/bin/env bash
# steps/03-fetch-kubeconfig.sh
#
# Pulls the cluster's kubeconfig via helpers/nks-control and writes it to a
# scenario-scoped path. Subsequent steps (04, 05) export KUBECONFIG to that
# path so kubectl never touches the user's default ~/.kube/config.

set -euo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"

# shellcheck disable=SC1091
source "$SCENARIO_DIR/state.env"
: "${CLUSTER_ID:?state.env did not carry CLUSTER_ID}"

CRED="${HOME}/.nhncloud/credentials"
read_cred() { grep -E "^${1}[[:space:]]*=" "$CRED" | head -1 | sed -E 's/^[^=]*=[[:space:]]*//' | tr -d '"'; }
export NHNCLOUD_USERNAME=$(read_cred username)
export NHNCLOUD_PASSWORD=$(read_cred api_password)
export NHNCLOUD_TENANT_ID=$(read_cred tenant_id)
NHNCLOUD_REGION=$(read_cred region); export NHNCLOUD_REGION="${NHNCLOUD_REGION:-kr1}"

NKS_BIN="$SCENARIO_DIR/helpers/nks-control/nks-control"

KUBECONFIG_DIR="$SCENARIO_DIR/kubeconfigs/$RUN_TS"
mkdir -p "$KUBECONFIG_DIR"
chmod 700 "$KUBECONFIG_DIR"
KCFG="$KUBECONFIG_DIR/config"

echo "[fetch-kubeconfig] fetching kubeconfig for $CLUSTER_ID..."
"$NKS_BIN" kubeconfig --cluster-id "$CLUSTER_ID" --out "$KCFG"

# Sanity: file is non-empty YAML and the API server responds.
[[ -s "$KCFG" ]] || { echo "[fetch-kubeconfig] ERROR: empty kubeconfig" >&2; exit 1; }

echo "KUBECONFIG=$KCFG" >> "$SCENARIO_DIR/state.env"
export KUBECONFIG="$KCFG"

echo "[fetch-kubeconfig] probing API server..."
# kubectl cluster-info is a low-overhead reachability check. Allow up to
# 90s for the apiserver to be ready after cluster reports ACTIVE — NHN
# sometimes signals ACTIVE a few seconds before the apiserver TLS is up.
for i in $(seq 1 18); do
  if kubectl --kubeconfig "$KCFG" cluster-info >/dev/null 2>&1; then
    echo "[fetch-kubeconfig] apiserver reachable after ${i} probe(s)"
    break
  fi
  printf '.'; sleep 5
  if [[ $i -eq 18 ]]; then
    echo
    echo "[fetch-kubeconfig] WARNING: apiserver still not reachable after 90s; continuing anyway" >&2
  fi
done
