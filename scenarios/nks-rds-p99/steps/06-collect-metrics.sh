#!/usr/bin/env bash
# steps/06-collect-metrics.sh  (scenario: nks-rds-p99)
#
# Parses the loadgen JSON written by 05-run-oltp.sh (path from
# state.env LOADGEN_REPORT_JSON) and emits:
#   reports/<RUN_TS>.json  — machine-readable
#   reports/<RUN_TS>.md    — human-readable, with headline p99
#
# The loadgen JSON file has the same key=value summary shape as
# rds-mysql-p99-internal so reports compare 1-for-1.
#
# Loadgen final-line format (key=value), e.g.:
#   tps=38.0 p50=24ms p95=54ms p99=89ms max=125ms ops=4567 errs=0 duration=120.0s
#
# LOADGEN_REPORT_JSON is set by 05-run-oltp.sh to the path of a file
# whose last line matching ^tps= is the loadgen summary.

set -euo pipefail
: "${SCENARIO_DIR:?must be set by run.sh}"
: "${RUN_TS:?must be set by run.sh}"

# shellcheck disable=SC1091
source "$SCENARIO_DIR/state.env"
: "${INSTANCE_ID:?state.env missing INSTANCE_ID}"
: "${ENDPOINT:?state.env missing ENDPOINT}"
# LOADGEN_REPORT_TXT is the loadgen's raw stdout (`tps=… p99=… …` line);
# LOADGEN_REPORT_JSON is the structured sidecar 05 already builds from that
# same summary. Both are written to state.env by 05-run-oltp.sh. The
# original 06 mistakenly grep'd the JSON file for `^tps=` and bailed.
: "${LOADGEN_REPORT_TXT:?state.env missing LOADGEN_REPORT_TXT — 05-run-oltp.sh did not run}"

RAW="$LOADGEN_REPORT_TXT"
JSON="$SCENARIO_DIR/reports/$RUN_TS.json"
MD="$SCENARIO_DIR/reports/$RUN_TS.md"

if [[ ! -s "$RAW" ]]; then
  echo "[metrics] ERROR: raw loadgen output missing or empty: $RAW" >&2
  exit 1
fi

SUMMARY=$(grep -E '^tps=' "$RAW" | tail -1)
if [[ -z "$SUMMARY" ]]; then
  echo "[metrics] ERROR: no summary line in $RAW" >&2
  exit 1
fi

# Strip optional trailing units (ms / s) and pull each key=value pair.
extract() { echo "$SUMMARY" | grep -oE "$1=[^ ]+" | sed -E "s/^$1=//; s/(ms|s)\$//"; }

tps=$(extract tps)
p50=$(extract p50)
p95=$(extract p95)
p99=$(extract p99)
maxlat=$(extract max)
ops=$(extract ops)
errs=$(extract errs)
dur=$(extract duration)

tps=${tps:-0}; p50=${p50:-0}; p95=${p95:-0}; p99=${p99:-0}
maxlat=${maxlat:-0}; ops=${ops:-0}; errs=${errs:-0}; dur=${dur:-0}

jq -n \
  --arg ts        "$RUN_TS" \
  --arg id        "$INSTANCE_ID" \
  --arg endpoint  "$ENDPOINT" \
  --arg preset    "${PRESET:-smoke}" \
  --argjson tps   "$tps" \
  --argjson p50   "$p50" \
  --argjson p95   "$p95" \
  --argjson p99   "$p99" \
  --argjson max   "$maxlat" \
  --argjson ops   "$ops" \
  --argjson errs  "$errs" \
  --argjson dur   "$dur" \
  '{ts:$ts,
    instance_id:$id,
    endpoint:$endpoint,
    preset:$preset,
    tps:$tps,
    latency_ms:{p50:$p50, p95:$p95, p99:$p99, max:$max},
    ops:$ops, errors:$errs, duration_s:$dur}' \
  > "$JSON"

cat > "$MD" <<EOF
# nks-rds-p99 — $RUN_TS

- instance_id: \`$INSTANCE_ID\`
- endpoint: \`$ENDPOINT\`
- preset: ${PRESET:-smoke}
- duration: ${dur}s
- ops: $ops, errors: $errs

| metric | value |
|---|---|
| TPS | $tps |
| p50 (ms) | $p50 |
| p95 (ms) | $p95 |
| **p99 (ms)** | **$p99** |
| max (ms) | $maxlat |

Raw output: [\`$(basename "$RAW")\`]($(basename "$RAW"))
EOF

echo "[metrics] wrote $JSON"
echo "[metrics] wrote $MD"
echo "[metrics] headline → tps=$tps  p99=${p99}ms  errors=$errs"
