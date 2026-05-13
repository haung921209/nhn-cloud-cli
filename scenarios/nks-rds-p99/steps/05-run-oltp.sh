#!/usr/bin/env bash
# steps/05-run-oltp.sh
#
# Runs the loadgen "run" phase via `kubectl exec` against the runner Pod
# created in 04. Captures the final summary line — the loadgen prints
#   tps=… p50=… p95=… p99=… max=… errs=…
# on the last stdout line — into reports/$RUN_TS.oltp.txt for 06 to parse.
#
# Also persists a JSON sidecar (LOADGEN_REPORT_JSON) by reformatting the
# summary line into a structured object — the loadgen itself prints text,
# so we synthesize the JSON here.

set -euo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"

# shellcheck disable=SC1091
source "$SCENARIO_DIR/state.env"
: "${KUBECONFIG:?state.env did not carry KUBECONFIG}"
: "${POD_NAME:?state.env did not carry POD_NAME — deploy step missing}"
: "${POD_NS:?state.env did not carry POD_NS}"
: "${PRESET:?run.sh must export PRESET}"

export KUBECONFIG

THREADS=$(yq -r ".presets.${PRESET}.sysbench_threads" "$SCENARIO_DIR/scenario.yaml")
DURATION_S=$(yq -r ".presets.${PRESET}.sysbench_time_s" "$SCENARIO_DIR/scenario.yaml")
TABLES=$(yq -r ".presets.${PRESET}.sysbench_tables" "$SCENARIO_DIR/scenario.yaml")
ROWS=$(yq -r ".presets.${PRESET}.sysbench_table_size" "$SCENARIO_DIR/scenario.yaml")

# Match 04-deploy-loadgen.sh — only inject CA_FILE if the local cert export
# produced a file. Otherwise loadgen falls back to InsecureSkipVerify.
CA_FILE_ENV=""
if [[ -n "${CA_FILE:-}" && -s "${CA_FILE:-}" ]]; then
  CA_FILE_ENV="CA_FILE=/etc/rds-ca/ca.pem "
fi

REPORTS="$SCENARIO_DIR/reports"
mkdir -p "$REPORTS"
OLTP_TXT="$REPORTS/$RUN_TS.oltp.txt"
OLTP_JSON="$REPORTS/$RUN_TS.oltp.json"

echo "[run-oltp] preset=$PRESET threads=$THREADS duration=${DURATION_S}s tables=$TABLES rows=$ROWS"
echo "[run-oltp] capturing pod stdout → $OLTP_TXT"

set -o pipefail
kubectl -n "$POD_NS" exec "$POD_NAME" -- /bin/sh -c "
  ${CA_FILE_ENV}PHASE=run THREADS=$THREADS DURATION_S=$DURATION_S TABLES=$TABLES ROWS=$ROWS /tmp/loadgen
" | tee "$OLTP_TXT"

# Last non-blank line is the summary, e.g.:
#   tps=4453.5 p50=2 p95=8 p99=10 max=42 errs=0
SUMMARY=$(grep -E '^tps=' "$OLTP_TXT" | tail -1)
if [[ -z "$SUMMARY" ]]; then
  echo "[run-oltp] ERROR: no 'tps=…' summary line found in pod output" >&2
  exit 1
fi

echo "[run-oltp] summary: $SUMMARY"

# Convert the space-separated k=v summary into JSON.
{
  echo '{'
  echo "  \"run_ts\": \"$RUN_TS\","
  echo "  \"preset\": \"$PRESET\","
  echo "  \"threads\": $THREADS,"
  echo "  \"duration_s\": $DURATION_S,"
  echo "  \"tables\": $TABLES,"
  echo "  \"rows\": $ROWS,"
  IFS=' ' read -r -a kvs <<< "$SUMMARY"
  total=${#kvs[@]}
  for ((i=0; i<total; i++)); do
    pair="${kvs[$i]}"
    k="${pair%%=*}"
    v="${pair##*=}"
    sep=","; [[ $i -eq $((total-1)) ]] && sep=""
    # All summary values are numeric — emit unquoted.
    echo "  \"$k\": $v$sep"
  done
  echo '}'
} > "$OLTP_JSON"

{
  echo "LOADGEN_REPORT_TXT=$OLTP_TXT"
  echo "LOADGEN_REPORT_JSON=$OLTP_JSON"
  echo "LOADGEN_SUMMARY=$SUMMARY"
} >> "$SCENARIO_DIR/state.env"

echo "[run-oltp] OK — JSON sidecar at $OLTP_JSON"
