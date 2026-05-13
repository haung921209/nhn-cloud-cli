#!/usr/bin/env bash
# steps/00-preflight.sh  (scenario: nks-rds-p99)
#
# SDK-level preflight validation — runs BEFORE any billable resource is created.
#
# Checks (read-only API calls only, no billing impact):
#   1. ~/.nhncloud/credentials has required fields
#      (access_key_id, secret_access_key, rds_app_key,
#       username, password, tenant_id — the last three are needed by nks-control)
#   2. Bearer token can be issued via /oauth2/token/create
#   3. RDS MySQL service reachable (describe-db-instances returns 200)
#   4. Compute service reachable (describe-flavors returns 200)
#      (Compute API is used for keypair create/delete even in NKS mode)
#   5. helpers/nks-control/nks-control binary is present and executable
#   6. kubectl is on PATH
#   7. scenario.yaml UUIDs exist in the live catalog:
#      a. flavor_id in describe-db-instance-classes  (RDS flavor)
#      b. subnet_id in describe-subnets
#      c. db_parameter_group_id in describe-db-parameter-groups
#   8. NKS prerequisites resolved & validated:
#      - kube_tag in scenario.yaml is one of nks-control list-versions output
#      - all required nks_labels are non-empty
#      - NKS_NETWORK_ID written to state.env (VPC UUID resolved from subnet)
#      - NKS_NODE_IMAGE written to state.env (if scenario.yaml left blank,
#        fall back to a literal that 01d will fail loud on)
#
# Note: NHN's `cluster_template_id` field must be the literal string
# "iaas_console" — there is no template-listing API on this account
# (GET /clustertemplates returns 404). The helper defaults to that string;
# nothing for preflight to resolve here.
#
# Ubuntu Compute image resolution is NOT needed for this scenario —
# NKS manages its own node images internally.
#
# Target wall time: <45 s when all checks pass.
#
# Exit codes:
#   0  all checks passed
#   1  one or more checks failed (details printed inline)

set -euo pipefail

# ─── Locate scenario dir (standalone or called from run.sh) ─────────────────
SCENARIO_DIR="${SCENARIO_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
RUN_TS="${RUN_TS:-$(date -u +%Y%m%dT%H%M%SZ)}"

# ─── Engine constants (engine-aware preflight — Task 4) ──────────────────────
ENGINE="${ENGINE:-$(yq -r '.parameters.engine // "mysql"' "$SCENARIO_DIR/scenario.yaml")}"
CLI_PREFIX=$(yq -r ".parameters.${ENGINE}.cli_prefix" "$SCENARIO_DIR/scenario.yaml")
CRED_APPKEY_PRIMARY=$(yq -r ".parameters.${ENGINE}.cred_file_appkey" "$SCENARIO_DIR/scenario.yaml")
CRED_APPKEY_ALT=$(yq -r ".parameters.${ENGINE}.cred_file_appkey_alt // \"\"" "$SCENARIO_DIR/scenario.yaml")
echo "[preflight] engine=$ENGINE cli_prefix=$CLI_PREFIX"

# ─── Helpers ─────────────────────────────────────────────────────────────────
TOTAL=8
FAILED=0

pass() { echo "[preflight] step $1/$TOTAL: $2... OK"; }
fail() { echo "[preflight] step $1/$TOTAL: $2... FAIL: $3"; FAILED=$((FAILED + 1)); }

# ─── Tool availability guard ─────────────────────────────────────────────────
for tool in jq yq curl; do
  if ! command -v "$tool" >/dev/null 2>&1; then
    echo "[preflight] ABORT: required tool '$tool' not found in PATH" >&2
    exit 1
  fi
done

# Use whichever nhncloud is on PATH. run.sh prepends $SCENARIO_DIR/bin
# (built fresh from ../../nhn-cloud-cli) so this resolves to the workspace
# build, not the older homebrew copy.
if ! command -v nhncloud >/dev/null 2>&1; then
  echo "[preflight] ABORT: nhncloud CLI not found in PATH" >&2
  exit 1
fi
NHNCLOUD=nhncloud

NKS_CTL="$SCENARIO_DIR/helpers/nks-control/nks-control"

# ─── Step 1: credentials file has required fields ────────────────────────────
CRED="$HOME/.nhncloud/credentials"
STEP=1
DESC="credentials file has required fields (access_key_id, secret_access_key, rds_app_key, username, password, tenant_id)"

if [[ ! -r "$CRED" ]]; then
  fail $STEP "$DESC" "$CRED not readable"
else
  read_cred() { grep -E "^${1}[[:space:]]*=" "$CRED" | head -1 | sed -E 's/^[^=]*=[[:space:]]*//' | tr -d '"'; }
  ACCESS_KEY=$(read_cred access_key_id)
  SECRET=$(read_cred secret_access_key)
  APPKEY=$(read_cred "$CRED_APPKEY_PRIMARY")
  if [[ -z "$APPKEY" && -n "$CRED_APPKEY_ALT" ]]; then
    APPKEY=$(read_cred "$CRED_APPKEY_ALT")
  fi
  REGION=$(read_cred region)
  NHN_USERNAME=$(read_cred username)
  NHN_PASSWORD=$(read_cred api_password)
  TENANT_ID=$(read_cred tenant_id)
  [[ -z "$REGION" ]] && REGION=kr1
  REGION=$(echo "$REGION" | tr '[:upper:]' '[:lower:]')

  MISSING=""
  [[ -z "$ACCESS_KEY"   ]] && MISSING="${MISSING}access_key_id "
  [[ -z "$SECRET"       ]] && MISSING="${MISSING}secret_access_key "
  [[ -z "$APPKEY"       ]] && MISSING="${MISSING}${CRED_APPKEY_PRIMARY}(or ${CRED_APPKEY_ALT}) "
  [[ -z "$NHN_USERNAME" ]] && MISSING="${MISSING}username "
  [[ -z "$NHN_PASSWORD" ]] && MISSING="${MISSING}password "
  [[ -z "$TENANT_ID"    ]] && MISSING="${MISSING}tenant_id "

  if [[ -n "$MISSING" ]]; then
    fail $STEP "$DESC" "missing fields: $MISSING"
  else
    pass $STEP "$DESC"
  fi
fi

# Abort early if credentials are not usable — remaining steps all need them.
if [[ $FAILED -gt 0 ]]; then
  echo "[preflight] ABORT: cannot proceed without valid credentials"
  exit 1
fi

# ─── Step 2: Bearer token issuance ───────────────────────────────────────────
STEP=2
DESC="Bearer token issuable via /oauth2/token/create"

BASIC=$(printf '%s:%s' "$ACCESS_KEY" "$SECRET" | base64 | tr -d '\n')
TOKEN_JSON=$(curl -fsS --max-time 15 -X POST \
  https://oauth.api.nhncloudservice.com/oauth2/token/create \
  -H "Authorization: Basic $BASIC" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'grant_type=client_credentials' 2>&1) || true

TOKEN=$(echo "$TOKEN_JSON" | jq -r '.access_token // empty' 2>/dev/null || true)
if [[ -n "$TOKEN" ]]; then
  pass $STEP "$DESC"
else
  fail $STEP "$DESC" "token response: $(echo "$TOKEN_JSON" | head -c 200)"
fi

# ─── Step 3: RDS service reachable (engine-aware) ────────────────────────────
STEP=3
DESC="RDS $(echo "$ENGINE" | tr '[:lower:]' '[:upper:]') service reachable (describe-db-instances)"

if "$NHNCLOUD" "$CLI_PREFIX" describe-db-instances -o json >/dev/null 2>&1; then
  pass $STEP "$DESC"
else
  fail $STEP "$DESC" "'nhncloud $CLI_PREFIX describe-db-instances' exited non-zero"
fi

# ─── Step 4: Compute service reachable ───────────────────────────────────────
# Compute API is used by 01b-create-keypair for NKS node keypair management.
STEP=4
DESC="Compute service reachable (describe-flavors)"

if "$NHNCLOUD" compute describe-flavors -o json >/dev/null 2>&1; then
  pass $STEP "$DESC"
else
  fail $STEP "$DESC" "'nhncloud compute describe-flavors' exited non-zero"
fi

# ─── Step 5: nks-control binary present and executable ───────────────────────
# run.sh builds the binary before calling preflight; if it's absent something
# went wrong in the build step.
STEP=5
DESC="helpers/nks-control/nks-control binary present and executable"

if [[ -x "$NKS_CTL" ]]; then
  pass $STEP "$DESC"
else
  fail $STEP "$DESC" "not found or not executable: $NKS_CTL"
fi

# ─── Step 6: kubectl on PATH ─────────────────────────────────────────────────
STEP=6
DESC="kubectl is on PATH"

if command -v kubectl >/dev/null 2>&1; then
  pass $STEP "$DESC"
else
  fail $STEP "$DESC" "kubectl not found — install kubectl before running this scenario"
fi

# ─── Step 7: scenario.yaml UUIDs exist in live catalogs ──────────────────────
STEP=7
DESC="scenario.yaml UUIDs present in live catalogs"

YAML="$SCENARIO_DIR/scenario.yaml"
UUID_FAILS=""

# 7a: flavor_id — engine-aware UUID source and describe-* verb
FLAVOR_ID=$(yq -r ".parameters.${ENGINE}.flavor_id // .parameters.flavor_id" "$YAML")
# MariaDB/PostgreSQL CLI expose describe-db-flavors (not -classes). MySQL is
# -classes (and -flavors raises an unknown-command). Pick the right verb per
# engine. NOTE: PG has its OWN flavor catalog (different UUIDs from MySQL even
# for same m2.c1m2 SKU name).
if [[ "$ENGINE" == "mariadb" || "$ENGINE" == "postgresql" ]]; then
  FLAVOR_JSON=$("$NHNCLOUD" "$CLI_PREFIX" describe-db-flavors -o json 2>/dev/null) || true
else
  FLAVOR_JSON=$("$NHNCLOUD" "$CLI_PREFIX" describe-db-instance-classes -o json 2>/dev/null) || true
fi
if [[ -n "$FLAVOR_ID" && "$FLAVOR_ID" != "FILL_ME_IN" ]]; then
  if echo "$FLAVOR_JSON" | jq -e --arg id "$FLAVOR_ID" \
      'if type == "array" then map(select(.dbInstanceClassId == $id or .dbFlavorId == $id or .id == $id)) | length > 0
       elif .dbInstanceClasses then .dbInstanceClasses | map(select(.dbInstanceClassId == $id or .id == $id)) | length > 0
       elif .dbFlavors then .dbFlavors | map(select(.dbFlavorId == $id or .id == $id)) | length > 0
       else tostring | contains($id) end' >/dev/null 2>&1; then
    echo "[preflight]   7a/3: flavor_id $FLAVOR_ID... OK"
  else
    echo "[preflight]   7a/3: flavor_id $FLAVOR_ID... FAIL: not found in describe-db-instance-classes / describe-db-flavors"
    UUID_FAILS="${UUID_FAILS}flavor_id "
  fi
else
  echo "[preflight]   7a/3: flavor_id is FILL_ME_IN or empty — skipped"
  UUID_FAILS="${UUID_FAILS}flavor_id(unset) "
fi

# 7b: subnet_id — engine-aware: MariaDB/PostgreSQL CLIs lack describe-subnets,
# fall back to the Network service. The subnet catalog is VPC-wide, so all
# RDS engines see the same one regardless of which CLI surfaces the verb.
SUBNET_ID=$(yq -r '.parameters.subnet_id' "$YAML")
if [[ "$ENGINE" == "mariadb" || "$ENGINE" == "postgresql" ]]; then
  SUBNET_JSON=$("$NHNCLOUD" network describe-subnets -o json 2>/dev/null) || true
else
  SUBNET_JSON=$("$NHNCLOUD" "$CLI_PREFIX" describe-subnets -o json 2>/dev/null) || true
fi
if [[ -n "$SUBNET_ID" && "$SUBNET_ID" != "FILL_ME_IN" ]]; then
  if echo "$SUBNET_JSON" | jq -e --arg id "$SUBNET_ID" \
      'if type == "array" then map(select(.subnetId == $id or .id == $id)) | length > 0
       elif .subnets then .subnets | map(select(.subnetId == $id or .id == $id)) | length > 0
       else tostring | contains($id) end' >/dev/null 2>&1; then
    echo "[preflight]   7b/3: subnet_id $SUBNET_ID... OK"
  else
    echo "[preflight]   7b/3: subnet_id $SUBNET_ID... FAIL: not found in describe-subnets"
    UUID_FAILS="${UUID_FAILS}subnet_id "
  fi
else
  echo "[preflight]   7b/3: subnet_id is FILL_ME_IN or empty — skipped"
  UUID_FAILS="${UUID_FAILS}subnet_id(unset) "
fi

# 7c: db_parameter_group_id — engine-aware UUID source, same describe verb for both engines
PARAM_GROUP_ID=$(yq -r ".parameters.${ENGINE}.db_parameter_group_id // .parameters.db_parameter_group_id" "$YAML")
PARAM_JSON=$("$NHNCLOUD" "$CLI_PREFIX" describe-db-parameter-groups -o json 2>/dev/null) || true
if [[ -n "$PARAM_GROUP_ID" && "$PARAM_GROUP_ID" != "FILL_ME_IN" ]]; then
  if echo "$PARAM_JSON" | jq -e --arg id "$PARAM_GROUP_ID" \
      'if type == "array" then map(select(.dbParameterGroupId == $id or .id == $id)) | length > 0
       elif .dbParameterGroups then .dbParameterGroups | map(select(.dbParameterGroupId == $id or .id == $id)) | length > 0
       else tostring | contains($id) end' >/dev/null 2>&1; then
    echo "[preflight]   7c/3: db_parameter_group_id $PARAM_GROUP_ID... OK"
  else
    echo "[preflight]   7c/3: db_parameter_group_id $PARAM_GROUP_ID... FAIL: not found in describe-db-parameter-groups"
    UUID_FAILS="${UUID_FAILS}db_parameter_group_id "
  fi
else
  echo "[preflight]   7c/3: db_parameter_group_id is FILL_ME_IN or empty — skipped"
  UUID_FAILS="${UUID_FAILS}db_parameter_group_id(unset) "
fi

if [[ -z "$UUID_FAILS" ]]; then
  pass $STEP "$DESC"
else
  fail $STEP "$DESC" "missing UUIDs: $UUID_FAILS"
fi

# ─── Step 8: NKS prerequisites resolved & validated ─────────────────────────
# (a) kube_tag in scenario.yaml is one of nks-control list-versions output
# (b) all required nks_labels keys present and non-empty
# (c) write NKS_NETWORK_ID to state.env (VPC UUID resolved from subnet via
#     `nhncloud network describe-subnet` — same trick CLI's compute
#     create-instance uses since round 2)
# (d) write NKS_NODE_IMAGE to state.env (yaml value if pinned; left blank
#     to surface in 01d if scenario.yaml didn't pin it — auto-resolve
#     from compute images is a follow-up)
STEP=8
DESC="NKS prerequisites resolved & validated"

NKS_FAILS=""

# (a) kube_tag valid
if [[ -x "$NKS_CTL" ]]; then
  KUBE_TAG=$(yq -r '.parameters.nks_labels.kube_tag // ""' "$YAML")
  VERSIONS=$(NHNCLOUD_USERNAME="$NHN_USERNAME" NHNCLOUD_PASSWORD="$NHN_PASSWORD" \
             NHNCLOUD_TENANT_ID="$TENANT_ID" NHNCLOUD_REGION="$REGION" \
             "$NKS_CTL" list-versions 2>/dev/null | tr -d ' ')
  if [[ -z "$KUBE_TAG" ]]; then
    NKS_FAILS="${NKS_FAILS}kube_tag(empty) "
  elif ! printf '%s\n' "$VERSIONS" | grep -qx "$KUBE_TAG"; then
    NKS_FAILS="${NKS_FAILS}kube_tag($KUBE_TAG not in list-versions: $(echo "$VERSIONS" | tr '\n' ',')) "
  else
    echo "[preflight]   8a: kube_tag $KUBE_TAG is supported"
  fi
else
  NKS_FAILS="${NKS_FAILS}nks-control-missing "
fi

# (b) required labels non-empty
for lbl in availability_zone node_image boot_volume_type boot_volume_size cert_manager_api ca_enable master_lb_floating_ip_enabled; do
  v=$(yq -r ".parameters.nks_labels.${lbl} // \"\"" "$YAML")
  if [[ -z "$v" || "$v" == "null" ]]; then
    # node_image empty is allowed — falls through to 01d which may use a
    # state.env override or fail loud. All other labels MUST be set.
    if [[ "$lbl" != "node_image" ]]; then
      NKS_FAILS="${NKS_FAILS}label.${lbl}(empty) "
    fi
  fi
done

# (c) NKS_NETWORK_ID via subnet→vpc resolution. The rds-mysql describe-subnets
# response (used in step 7b) only carries subnetId/Name/Cidr — it does NOT
# return vpc_id, so we ask the VPC service directly via the helper. SDK
# v0.1.32+ populates Subnet.VPCID from NHN's `vpc_id` key.
SUBNET_ID=$(yq -r '.parameters.subnet_id' "$YAML")
NKS_NETWORK_ID=""
if [[ -x "$NKS_CTL" ]]; then
  RESOLVE_OUT=$(NHNCLOUD_USERNAME="$NHN_USERNAME" NHNCLOUD_PASSWORD="$NHN_PASSWORD" \
                NHNCLOUD_TENANT_ID="$TENANT_ID" NHNCLOUD_REGION="$REGION" \
                "$NKS_CTL" resolve-vpc --subnet-id "$SUBNET_ID" 2>&1) || true
  NKS_NETWORK_ID=$(printf '%s\n' "$RESOLVE_OUT" | awk -F= '/^VPC_ID=/ {print $2}' | head -1)
fi

if [[ -n "$NKS_NETWORK_ID" ]]; then
  echo "[preflight]   8c: NKS_NETWORK_ID=$NKS_NETWORK_ID (resolved from subnet $SUBNET_ID)"
  STATE_ENV="$SCENARIO_DIR/state.env"
  if [[ -f "$STATE_ENV" ]]; then
    grep -v '^NKS_NETWORK_ID=' "$STATE_ENV" > "${STATE_ENV}.tmp" 2>/dev/null \
      && mv "${STATE_ENV}.tmp" "$STATE_ENV" || true
  fi
  echo "NKS_NETWORK_ID=$NKS_NETWORK_ID" >> "$STATE_ENV"
else
  NKS_FAILS="${NKS_FAILS}NKS_NETWORK_ID(unresolved-from-subnet) "
fi

# (d) NKS_NODE_IMAGE — auto-resolve via nks-control list-node-images. The
# helper hits the Glance API with the NKS-only filter (SDK ≥ v0.1.34 threads
# this through image.ListImagesInput.ExtraParams). If scenario.yaml pinned a
# UUID we honour it (and verify it's in the listing); otherwise we pick the
# newest Ubuntu LTS image returned, since that's what NHN's Container-tagged
# templates typically expect.
NODE_IMAGE_YAML=$(yq -r '.parameters.nks_labels.node_image // ""' "$YAML")
NKS_NODE_IMAGE=""
if [[ -x "$NKS_CTL" ]]; then
  IMAGES_TSV=$(NHNCLOUD_USERNAME="$NHN_USERNAME" NHNCLOUD_PASSWORD="$NHN_PASSWORD" \
               NHNCLOUD_TENANT_ID="$TENANT_ID" NHNCLOUD_REGION="$REGION" \
               "$NKS_CTL" list-node-images 2>/dev/null) || IMAGES_TSV=""
  if [[ -z "$IMAGES_TSV" ]]; then
    NKS_FAILS="${NKS_FAILS}node_image(list-node-images returned nothing) "
  else
    if [[ -n "$NODE_IMAGE_YAML" ]]; then
      # Verify pinned UUID exists in listing.
      if printf '%s\n' "$IMAGES_TSV" | awk -F'\t' -v id="$NODE_IMAGE_YAML" '$1==id{found=1} END{exit !found}'; then
        NKS_NODE_IMAGE="$NODE_IMAGE_YAML"
        echo "[preflight]   8d: pinned node_image $NKS_NODE_IMAGE found in NKS image listing"
      else
        NKS_FAILS="${NKS_FAILS}node_image($NODE_IMAGE_YAML not in NKS image listing) "
      fi
    else
      # Auto-pick: newest Ubuntu LTS image. Sort by version number lex desc.
      NKS_NODE_IMAGE=$(printf '%s\n' "$IMAGES_TSV" \
        | awk -F'\t' '$2 ~ /Ubuntu Server [0-9]+\.[0-9]+/ {print}' \
        | sort -t$'\t' -k2,2 -r | head -1 | awk -F'\t' '{print $1}')
      if [[ -n "$NKS_NODE_IMAGE" ]]; then
        IMG_NAME=$(printf '%s\n' "$IMAGES_TSV" | awk -F'\t' -v id="$NKS_NODE_IMAGE" '$1==id{print $2}')
        echo "[preflight]   8d: auto-selected node_image $NKS_NODE_IMAGE ($IMG_NAME)"
      else
        NKS_FAILS="${NKS_FAILS}node_image(no Ubuntu image in NKS listing) "
      fi
    fi
  fi
fi
if [[ -n "$NKS_NODE_IMAGE" ]]; then
  STATE_ENV="$SCENARIO_DIR/state.env"
  if [[ -f "$STATE_ENV" ]]; then
    grep -v '^NKS_NODE_IMAGE=' "$STATE_ENV" > "${STATE_ENV}.tmp" 2>/dev/null \
      && mv "${STATE_ENV}.tmp" "$STATE_ENV" || true
  fi
  echo "NKS_NODE_IMAGE=$NKS_NODE_IMAGE" >> "$STATE_ENV"
fi

if [[ -z "$NKS_FAILS" ]]; then
  pass $STEP "$DESC"
else
  fail $STEP "$DESC" "$NKS_FAILS"
fi

# ─── Summary ─────────────────────────────────────────────────────────────────
echo ""
if [[ $FAILED -eq 0 ]]; then
  echo "[preflight] ALL $TOTAL checks passed — safe to proceed with provisioning."
  exit 0
else
  echo "[preflight] $FAILED/$TOTAL checks FAILED — fix the above issues before re-running."
  exit 1
fi
