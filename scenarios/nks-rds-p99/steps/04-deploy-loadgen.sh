#!/usr/bin/env bash
# steps/04-deploy-loadgen.sh
#
# Sets up the in-cluster loadgen runner Pod:
#   1. Create namespace `bench` (idempotent).
#   2. Push DB credentials + connection metadata as a k8s Secret.
#   3. Push CA cert (from 02b-fetch-ca) as a k8s ConfigMap so the loadgen can
#      verify NHN's RDS TLS server cert.
#   4. Create a long-lived runner Pod from the public Debian image — Pod
#      sleeps `infinity`, we exec into it for prepare + run.
#   5. Wait for Pod Ready.
#   6. `kubectl cp` the loadgen-linux binary into the Pod.
#   7. Run the prepare phase (creates schema rows used by the benchmark).
#
# Why a single sleeping Pod (not a Job): we need two phases (prepare → run)
# with the binary preserved between them, plus interactive log collection.
# A Job would create a fresh Pod per phase and lose the warm dataset.

set -euo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"

# shellcheck disable=SC1091
source "$SCENARIO_DIR/state.env"
: "${KUBECONFIG:?state.env did not carry KUBECONFIG — fetch step missing}"
: "${ENDPOINT:?state.env did not carry ENDPOINT — wait-rds step missing}"
: "${APP_DB_USER:?state.env did not carry APP_DB_USER — wait-rds step missing}"
: "${APP_DB_PASSWORD:?must be exported by run.sh}"
: "${ENGINE:?state.env missing ENGINE}"

LOADGEN_ENGINE=$(yq -r ".parameters.${ENGINE}.loadgen_engine // \"mysql\"" "$SCENARIO_DIR/scenario.yaml")
DB_PORT_PARAM=$(yq -r ".parameters.${ENGINE}.db_port // 3306" "$SCENARIO_DIR/scenario.yaml")

export KUBECONFIG

NAMESPACE="bench"
POD_NAME="loadgen-runner"

DB_HOST="${ENDPOINT%:*}"
DB_PORT="${ENDPOINT##*:}"
DB_NAME=$(yq -r '.parameters.database' "$SCENARIO_DIR/scenario.yaml")
LOADGEN_IMAGE=$(yq -r '.parameters.loadgen_runner_image' "$SCENARIO_DIR/scenario.yaml")

# 1. Namespace ─────────────────────────────────────────────────────────────
echo "[deploy] ensuring namespace $NAMESPACE..."
kubectl get ns "$NAMESPACE" >/dev/null 2>&1 || kubectl create ns "$NAMESPACE"

# 2. Secret with creds + conn info ──────────────────────────────────────────
# kubectl create secret writes through stdin to avoid the password ever
# appearing as an argv (visible in `ps`).
#
# Round-5 PG quirk: NHN's master role has NO CREATEROLE / NO CREATEDB. The
# round-2 grants phase fails with `permission denied to create role (42501)`.
# But NHN's default pg_hba is open (`host all all 0.0.0.0/0 scram-sha-256`),
# so the master user works fine as the benchmark identity for prepare+run.
# Wire DB_USER/DB_PASSWORD straight to MASTER for engine=postgresql and skip
# the grants exec below. mysql/mariadb keep using the app DB user (their
# masters are host-locked, the app-user indirection is still required).
MASTER_DB_USER_VAL=$(yq -r '.parameters.user' "$SCENARIO_DIR/scenario.yaml")
if [[ "$LOADGEN_ENGINE" == "pg" ]]; then
  BENCH_DB_USER="$MASTER_DB_USER_VAL"
  BENCH_DB_PASSWORD="$DB_PASSWORD"
  echo "[deploy] postgresql: benchmark connects as master '$BENCH_DB_USER' (no CREATEROLE → no app-user indirection)"
else
  BENCH_DB_USER="$APP_DB_USER"
  BENCH_DB_PASSWORD="$APP_DB_PASSWORD"
fi

echo "[deploy] creating Secret rds-creds (user/password/database)..."
kubectl -n "$NAMESPACE" delete secret rds-creds --ignore-not-found >/dev/null
kubectl -n "$NAMESPACE" create secret generic rds-creds \
  --from-literal="DB_USER=$BENCH_DB_USER" \
  --from-literal="MYSQL_PWD=$BENCH_DB_PASSWORD" \
  --from-literal="DB_PASSWORD=$BENCH_DB_PASSWORD" \
  --from-literal="DB_NAME=$DB_NAME" \
  --from-literal="DB_HOST=$DB_HOST" \
  --from-literal="DB_PORT=$DB_PORT" \
  --from-literal="LOADGEN_ENGINE=$LOADGEN_ENGINE" \
  --from-literal="GRANT_USER=$APP_DB_USER" \
  --from-literal="GRANT_PASSWORD=$APP_DB_PASSWORD" \
  --from-literal="MASTER_DB_USER=$MASTER_DB_USER_VAL" \
  --from-literal="MASTER_DB_PASSWORD=$DB_PASSWORD" >/dev/null

# 3. CA cert ConfigMap ──────────────────────────────────────────────────────
if [[ -n "${CA_FILE:-}" && -s "$CA_FILE" ]]; then
  echo "[deploy] creating ConfigMap rds-ca from $CA_FILE..."
  kubectl -n "$NAMESPACE" delete configmap rds-ca --ignore-not-found >/dev/null
  kubectl -n "$NAMESPACE" create configmap rds-ca --from-file=ca.pem="$CA_FILE" >/dev/null
else
  echo "[deploy] WARNING: CA_FILE not set or empty — loadgen will fall back to skip-verify TLS" >&2
fi

# 4. Runner Pod ─────────────────────────────────────────────────────────────
# Idempotent: delete any prior pod from a previous run, then recreate.
echo "[deploy] (re)creating runner Pod $POD_NAME (image=$LOADGEN_IMAGE)..."
kubectl -n "$NAMESPACE" delete pod "$POD_NAME" --ignore-not-found --grace-period=0 --force >/dev/null 2>&1 || true

cat <<EOF | kubectl -n "$NAMESPACE" apply -f - >/dev/null
apiVersion: v1
kind: Pod
metadata:
  name: $POD_NAME
  labels:
    app: loadgen
spec:
  restartPolicy: Never
  terminationGracePeriodSeconds: 5
  containers:
  - name: runner
    image: $LOADGEN_IMAGE
    command: ["/bin/sh","-c","sleep infinity"]
    envFrom:
    - secretRef:
        name: rds-creds
    volumeMounts:
$( [[ -n "${CA_FILE:-}" && -s "${CA_FILE:-}" ]] && cat <<'VOL'
    - name: rds-ca
      mountPath: /etc/rds-ca
      readOnly: true
VOL
)
  volumes:
$( [[ -n "${CA_FILE:-}" && -s "${CA_FILE:-}" ]] && cat <<'VOL'
  - name: rds-ca
    configMap:
      name: rds-ca
VOL
)
EOF

# 5. Wait Ready ─────────────────────────────────────────────────────────────
echo "[deploy] waiting for Pod Ready (timeout 180s)..."
kubectl -n "$NAMESPACE" wait pod/"$POD_NAME" --for=condition=Ready --timeout=180s

# 6. kubectl cp loadgen binary ─────────────────────────────────────────────
LG_BIN="$SCENARIO_DIR/loadgen/loadgen-linux"
[[ -s "$LG_BIN" ]] || { echo "[deploy] ERROR: $LG_BIN missing — run.sh should have built it" >&2; exit 1; }

echo "[deploy] copying loadgen binary into Pod ($(wc -c < "$LG_BIN") bytes)..."
# Try kubectl cp first; fall back to streaming via exec --stdin if cp blows
# up under macOS sandbox (same workaround we used for scp in round-2).
if ! kubectl -n "$NAMESPACE" cp "$LG_BIN" "$POD_NAME:/tmp/loadgen" 2>/dev/null; then
  echo "[deploy] kubectl cp failed; falling back to stdin streaming..."
  kubectl -n "$NAMESPACE" exec -i "$POD_NAME" -- /bin/sh -c 'cat > /tmp/loadgen' < "$LG_BIN"
fi
kubectl -n "$NAMESPACE" exec "$POD_NAME" -- chmod +x /tmp/loadgen

# Plant CA path on Pod, if available, so the loadgen reads it as CA_FILE.
if kubectl -n "$NAMESPACE" exec "$POD_NAME" -- test -f /etc/rds-ca/ca.pem 2>/dev/null; then
  echo "[deploy] Pod sees /etc/rds-ca/ca.pem — TLS will verify"
fi

# 7. Prepare phase ──────────────────────────────────────────────────────────
TABLES=$(yq -r ".presets.${PRESET}.sysbench_tables" "$SCENARIO_DIR/scenario.yaml")
ROWS=$(yq -r ".presets.${PRESET}.sysbench_table_size" "$SCENARIO_DIR/scenario.yaml")

# Only inject CA_FILE if the cert export actually produced a file. Empty
# CA_FILE → loadgen falls back to InsecureSkipVerify (the round-2/3 path
# for MariaDB whose cert-export is broken NHN-side).
CA_FILE_ENV=""
if [[ -n "${CA_FILE:-}" && -s "${CA_FILE:-}" ]]; then
  CA_FILE_ENV="CA_FILE=/etc/rds-ca/ca.pem "
fi

if [[ "$LOADGEN_ENGINE" == "pg" ]]; then
  # Run-7 surfaced that NHN PG's master role lacks CREATEROLE — the grants
  # phase would 42501 here. Secret already wires DB_USER to master above, so
  # prepare/run go through master directly and we skip grants entirely.
  echo "[deploy] postgresql: skipping grants phase (master role lacks CREATEROLE; benchmark runs as master)"
fi

echo "[deploy] running prepare phase (tables=$TABLES rows=$ROWS, ca=${CA_FILE_ENV:+verify}${CA_FILE_ENV:-skip-verify})..."
kubectl -n "$NAMESPACE" exec "$POD_NAME" -- /bin/sh -c "
  ${CA_FILE_ENV}PHASE=prepare TABLES=$TABLES ROWS=$ROWS /tmp/loadgen
"

echo "[deploy] runner Pod ready, dataset prepared."
{
  echo "POD_NAME=$POD_NAME"
  echo "POD_NS=$NAMESPACE"
} >> "$SCENARIO_DIR/state.env"
