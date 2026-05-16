#!/usr/bin/env bash
# KHR-R: ResourceFuture simulation evidence on karl-metal-01@ovh (read-only).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
EVIDENCE="${ROOT}/docs/evidence/khr-resourcefuture"
RUN_ID="${KHR_RESOURCEFUTURE_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
BIN="${ROOT}/bin/karl-host-runtime"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-resourcefuture-simulate.yaml"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${ROOT}/bin"
(cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime)

khr_runtime_log "resourcefuture-simulate (read-only)"
"${BIN}" -mode=resourcefuture-simulate -config="${CFG}" -cluster-context="${CTX}" \
  > "${RUN_DIR}/simulation.json"

jq -e '.mode == "resourcefuture-simulate"' "${RUN_DIR}/simulation.json" >/dev/null
jq -e '.safety.readOnly == true and .safety.noApply == true and .safety.simulationOnly == true' \
  "${RUN_DIR}/simulation.json" >/dev/null
jq -e '.safety.noMutation == true and .safety.noRestart == true and .safety.noAutonomousOrchestration == true' \
  "${RUN_DIR}/simulation.json" >/dev/null
jq -e '(.candidateScalePlans | length) >= 1' "${RUN_DIR}/simulation.json" >/dev/null
jq -e '(.saturationForecast | length) >= 1' "${RUN_DIR}/simulation.json" >/dev/null
jq -e '(.forecasts.cpuScale | length) >= 1 and (.forecasts.ramScale | length) >= 1' \
  "${RUN_DIR}/simulation.json" >/dev/null

PLANS="$(jq '.candidateScalePlans | length' "${RUN_DIR}/simulation.json")"
LIVE="$(jq '[.liveInPlaceEligibility[] | select(.eligible)] | length' "${RUN_DIR}/simulation.json")"
RESTART="$(jq '[.restartRequiredPrediction[] | select(.required)] | length' "${RUN_DIR}/simulation.json")"

echo "no-apply: resourcefuture-simulate only" > "${RUN_DIR}/mutation-check.txt"
echo "no-patch: no kubectl apply in script" >> "${RUN_DIR}/mutation-check.txt"

jq -n \
  --arg sprint "KHR-R" \
  --arg runId "${RUN_ID}" \
  --arg cluster "${CTX}" \
  --argjson plans "${PLANS}" \
  --argjson live "${LIVE}" \
  --argjson restart "${RESTART}" \
  '{sprint:$sprint,runId:$runId,cluster:$cluster,candidatePlans:$plans,liveInPlaceEligible:$live,restartRequired:$restart,readOnly:true,noApply:true,simulationOnly:true}' \
  > "${RUN_DIR}/summary.json"
cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"

echo "[khr_resourcefuture_evidence] PASS ${RUN_DIR} plans=${PLANS} live=${LIVE} restart=${RESTART}"
