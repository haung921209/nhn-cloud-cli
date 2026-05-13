#!/usr/bin/env bash
# steps/01b-create-keypair.sh
#
# Phase B — provisions a per-run SSH keypair via NHN Cloud Compute (server-side
# generation) and saves both halves locally for later SSH access.
#
# WHY SERVER-GENERATED:
#   `nhncloud compute create-key-pair --key-name X --public-key "<content>"`
#   returns API Error 400 / Bad Request in our tenant regardless of the
#   public-key content (ed25519, rsa-2048, rsa-4096 — all fail).  The working
#   pattern is to omit --public-key entirely: NHN generates the RSA pair on its
#   side and returns both halves in the JSON response.  The `--public-key` flag
#   in the CLI help is therefore misleading / broken upstream; skip it.
#
# OUTPUT FILES:
#   $SCENARIO_DIR/keypairs/$RUN_TS/id_rsa      private key (mode 600)
#   $SCENARIO_DIR/keypairs/$RUN_TS/id_rsa.pub  public key  (mode 644)
#
# Persists: BENCH_KEYPAIR_NAME, BENCH_KEYPAIR_PRIV, BENCH_KEYPAIR_PUB

set -euo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"
: "${RUN_TS:?must be set by run.sh}"
: "${INSTANCE_NAME:?must be set by run.sh}"

# kp-HHMMSS-RR  →  2 + 1 + 6 + 1 + 2 = 12 chars (≤15 limit is safe)
KEYPAIR_NAME="kp-${INSTANCE_NAME#bench-}"

KEYDIR="$SCENARIO_DIR/keypairs/$RUN_TS"
mkdir -p "$KEYDIR"
chmod 700 "$KEYDIR"
PRIV="$KEYDIR/id_rsa"
PUB="$KEYDIR/id_rsa.pub"

echo "[keypair] provisioning server-generated keypair: $KEYPAIR_NAME"

# Use set +e so we can capture both exit-code and stdout/stderr together.
# NHN returns the full JSON (including private key) on stdout when creation
# succeeds; on failure it writes an error message to stdout as well.
set +e
CREATE_JSON=$(nhncloud compute create-key-pair \
                --key-name "$KEYPAIR_NAME" \
                -o json 2>&1)
CREATE_RC=$?
set -e

if [[ $CREATE_RC -ne 0 ]]; then
  echo "[keypair] ERROR: compute create-key-pair exited $CREATE_RC. Output:" >&2
  printf '%s\n' "$CREATE_JSON" >&2
  exit "$CREATE_RC"
fi

# NHN encodes the private key with literal "\n" in the JSON string; jq -r
# decodes those into real newlines, yielding a valid PEM block.
PRIV_KEY=$(printf '%s' "$CREATE_JSON" | jq -r '.keypair.private_key // empty')
PUB_KEY=$(printf '%s'  "$CREATE_JSON" | jq -r '.keypair.public_key  // empty')

if [[ -z "$PRIV_KEY" ]]; then
  echo "[keypair] ERROR: no private_key found in response" >&2
  printf '%s\n' "$CREATE_JSON" >&2
  exit 1
fi

printf '%s' "$PRIV_KEY" > "$PRIV"
chmod 600 "$PRIV"

if [[ -n "$PUB_KEY" ]]; then
  printf '%s\n' "$PUB_KEY" > "$PUB"
  chmod 644 "$PUB"
  echo "[keypair] private key → $PRIV"
  echo "[keypair] public  key → $PUB"
else
  echo "[keypair] private key → $PRIV  (no public key in response)"
fi

{
  echo "BENCH_KEYPAIR_NAME=$KEYPAIR_NAME"
  echo "BENCH_KEYPAIR_PRIV=$PRIV"
  echo "BENCH_KEYPAIR_PUB=$PUB"
} >> "$SCENARIO_DIR/state.env"

echo "[keypair] $KEYPAIR_NAME created server-side; private key saved at $PRIV"
