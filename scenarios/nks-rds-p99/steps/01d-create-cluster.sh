#!/usr/bin/env bash
# steps/01d-create-cluster.sh
#
# Submits NKS cluster creation via helpers/nks-control. Does NOT block on the
# cluster reaching ACTIVE — that's 02c-wait-cluster.sh, which runs after the
# RDS wait so total wall time = max(RDS_provision, NKS_provision).
#
# Critical NHN spec quirks (https://docs.nhncloud.com/en/Container/NKS/en/public-api/ "클러스터 생성하기"):
#   - cluster_template_id MUST be literal "iaas_console" (helper default).
#   - node_count is a *string* in the API body — helper handles the conv.
#   - flavor_id, fixed_network, fixed_subnet, keypair are root-level required.
#   - labels.{kube_tag, node_image, boot_volume_type, boot_volume_size,
#     cert_manager_api, ca_enable, master_lb_floating_ip_enabled,
#     availability_zone} are all REQUIRED.
#
# Inputs from state.env:
#   CLUSTER_NAME            (run.sh)
#   BENCH_KEYPAIR_NAME      (01b-create-keypair.sh)
#   NKS_NETWORK_ID          (00-preflight.sh: VPC UUID resolved from subnet)
#   NKS_NODE_IMAGE          (00-preflight.sh: NKS-compatible image UUID, or
#                            empty if scenario.yaml pins an explicit value)
#
# Inputs from scenario.yaml:
#   parameters.subnet_id, parameters.node_flavor_id, parameters.node_count
#   parameters.nks_labels.* (see above)
#
# Persists: CLUSTER_ID

set -euo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"

# shellcheck disable=SC1091
source "$SCENARIO_DIR/state.env"
: "${CLUSTER_NAME:?state.env did not carry CLUSTER_NAME}"
: "${BENCH_KEYPAIR_NAME:?state.env did not carry BENCH_KEYPAIR_NAME}"
: "${NKS_NETWORK_ID:?state.env did not carry NKS_NETWORK_ID — preflight VPC resolution missing}"

CRED="${HOME}/.nhncloud/credentials"
[[ -f "$CRED" ]] || { echo "[create-cluster] ERROR: $CRED missing" >&2; exit 1; }
read_cred() { grep -E "^${1}[[:space:]]*=" "$CRED" | head -1 | sed -E 's/^[^=]*=[[:space:]]*//' | tr -d '"'; }
export NHNCLOUD_USERNAME=$(read_cred username)
export NHNCLOUD_PASSWORD=$(read_cred api_password)
export NHNCLOUD_TENANT_ID=$(read_cred tenant_id)
NHNCLOUD_REGION=$(read_cred region); export NHNCLOUD_REGION="${NHNCLOUD_REGION:-kr1}"

SUBNET_ID=$(yq -r '.parameters.subnet_id' "$SCENARIO_DIR/scenario.yaml")
NODE_FLAVOR=$(yq -r '.parameters.node_flavor_id' "$SCENARIO_DIR/scenario.yaml")
NODE_COUNT=$(yq -r '.parameters.node_count' "$SCENARIO_DIR/scenario.yaml")

# Required labels (read from scenario.yaml). node_image may be overridden by
# preflight if scenario.yaml pinned it as "".
KUBE_TAG=$(yq -r '.parameters.nks_labels.kube_tag' "$SCENARIO_DIR/scenario.yaml")
NODE_IMAGE_YAML=$(yq -r '.parameters.nks_labels.node_image' "$SCENARIO_DIR/scenario.yaml")
NODE_IMAGE="${NKS_NODE_IMAGE:-$NODE_IMAGE_YAML}"
BOOT_VOLUME_TYPE=$(yq -r '.parameters.nks_labels.boot_volume_type' "$SCENARIO_DIR/scenario.yaml")
BOOT_VOLUME_SIZE=$(yq -r '.parameters.nks_labels.boot_volume_size' "$SCENARIO_DIR/scenario.yaml")
CERT_MANAGER_API=$(yq -r '.parameters.nks_labels.cert_manager_api' "$SCENARIO_DIR/scenario.yaml")
CA_ENABLE=$(yq -r '.parameters.nks_labels.ca_enable' "$SCENARIO_DIR/scenario.yaml")
MASTER_LB_FIP=$(yq -r '.parameters.nks_labels.master_lb_floating_ip_enabled' "$SCENARIO_DIR/scenario.yaml")
AZ=$(yq -r '.parameters.nks_labels.availability_zone' "$SCENARIO_DIR/scenario.yaml")
EXT_NET=$(yq -r '.parameters.nks_labels.external_network_id // ""' "$SCENARIO_DIR/scenario.yaml")
EXT_SUB=$(yq -r '.parameters.nks_labels.external_subnet_id_list // ""' "$SCENARIO_DIR/scenario.yaml")
CLUSTERAUTOSCALE=$(yq -r '.parameters.nks_labels.clusterautoscale // ""' "$SCENARIO_DIR/scenario.yaml")
ETCD_VOL_SIZE=$(yq -r '.parameters.nks_labels.etcd_volume_size // ""' "$SCENARIO_DIR/scenario.yaml")

if [[ -z "$NODE_IMAGE" || "$NODE_IMAGE" == "null" ]]; then
  echo "[create-cluster] ERROR: nks_labels.node_image not set and preflight did not auto-resolve NKS_NODE_IMAGE" >&2
  exit 1
fi

NKS_BIN="$SCENARIO_DIR/helpers/nks-control/nks-control"
[[ -x "$NKS_BIN" ]] || { echo "[create-cluster] ERROR: $NKS_BIN not built — run.sh should have built it" >&2; exit 1; }

echo "[create-cluster] submitting CreateCluster: name=$CLUSTER_NAME flavor=$NODE_FLAVOR nodes=$NODE_COUNT keypair=$BENCH_KEYPAIR_NAME network=$NKS_NETWORK_ID subnet=$SUBNET_ID"
echo "[create-cluster] labels: kube_tag=$KUBE_TAG node_image=$NODE_IMAGE az=$AZ boot=${BOOT_VOLUME_TYPE}/${BOOT_VOLUME_SIZE}GB cert_mgr=$CERT_MANAGER_API ca=$CA_ENABLE master_fip=$MASTER_LB_FIP"

CREATE_ARGS=(create
  --name "$CLUSTER_NAME"
  --node-count "$NODE_COUNT"
  --node-flavor-id "$NODE_FLAVOR"
  --keypair "$BENCH_KEYPAIR_NAME"
  --subnet-id "$SUBNET_ID"
  --network-id "$NKS_NETWORK_ID"
  --label "availability_zone=$AZ"
  --label "kube_tag=$KUBE_TAG"
  --label "node_image=$NODE_IMAGE"
  --label "boot_volume_type=$BOOT_VOLUME_TYPE"
  --label "boot_volume_size=$BOOT_VOLUME_SIZE"
  --label "cert_manager_api=$CERT_MANAGER_API"
  --label "ca_enable=$CA_ENABLE"
  --label "master_lb_floating_ip_enabled=$MASTER_LB_FIP"
)
# external_*_id are required by NHN's CreateCluster when the subnet's
# router has an internet gateway. Empty values are dropped so the
# label string isn't sent at all.
[[ -n "$EXT_NET" ]] && CREATE_ARGS+=(--label "external_network_id=$EXT_NET")
[[ -n "$EXT_SUB" ]] && CREATE_ARGS+=(--label "external_subnet_id_list=$EXT_SUB")
[[ -n "$CLUSTERAUTOSCALE" ]] && CREATE_ARGS+=(--label "clusterautoscale=$CLUSTERAUTOSCALE")
[[ -n "$ETCD_VOL_SIZE" ]] && CREATE_ARGS+=(--label "etcd_volume_size=$ETCD_VOL_SIZE")

# NHN's NKS keypair lookup races against compute create-key-pair: in
# practice we've seen "API Error 403: Unable to find keypair X"
# immediately after the keypair was created in 01b. A short propagation
# wait + retry covers the race without hiding genuine errors.
echo "[create-cluster] waiting 30s for keypair to propagate to NKS..."
sleep 30

CREATE_OUT=""
for i in 1 2 3; do
  if CREATE_OUT=$("$NKS_BIN" "${CREATE_ARGS[@]}" 2>&1); then
    break
  fi
  if echo "$CREATE_OUT" | grep -qE 'Unable to find keypair|API Error 403'; then
    echo "[create-cluster] transient 403 / keypair race — retry $i/3 in 30s..."
    echo "$CREATE_OUT" | tail -5
    sleep 30
    continue
  fi
  echo "[create-cluster] ERROR (non-retryable):" >&2
  echo "$CREATE_OUT" >&2
  exit 1
done

CLUSTER_ID=$(printf '%s\n' "$CREATE_OUT" | awk -F= '/^CLUSTER_ID=/ {print $2}' | head -1)
if [[ -z "$CLUSTER_ID" ]]; then
  echo "[create-cluster] ERROR: helper did not emit CLUSTER_ID. Raw output:" >&2
  printf '%s\n' "$CREATE_OUT" >&2
  exit 1
fi

echo "CLUSTER_ID=$CLUSTER_ID" >> "$SCENARIO_DIR/state.env"
echo "[create-cluster] submitted; CLUSTER_ID=$CLUSTER_ID (wait phase 02c will block until ACTIVE)"
