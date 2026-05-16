#!/usr/bin/env bash
# KHR-N: host heartbeat + runtime session evidence on karl-metal-01@ovh.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-loop.yaml"
EVIDENCE="${ROOT}/docs/evidence/khr-host-heartbeat"
RUN_ID="${KHR_HOST_HEARTBEAT_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
BIN="${ROOT}/bin/karl-host-runtime"
STATUS_OUT="${RUN_DIR}/host-status.json"
SANDBOX_DIR="${RUN_DIR}/sandbox"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${SANDBOX_DIR}" "${ROOT}/bin"
(cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime)

khr_runtime_log "heartbeat loop sandbox"
"${BIN}" -mode=host-heartbeat -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" -sandbox-dir="${SANDBOX_DIR}" \
  -heartbeat-iterations=3 -heartbeat-interval-ms=300 \
  -heartbeat-output="${STATUS_OUT}" > "${RUN_DIR}/heartbeat-loop.json"
jq -e '.noMutation == true and (.iterations | length) == 3' "${RUN_DIR}/heartbeat-loop.json" >/dev/null

SID="$(jq -r '.iterations[0].runtimeSession.runtimeSessionId' "${RUN_DIR}/heartbeat-loop.json")"
SID2="$(jq -r '.iterations[2].runtimeSession.runtimeSessionId' "${RUN_DIR}/heartbeat-loop.json")"
if [[ "${SID}" == "" || "${SID}" != "${SID2}" ]]; then
  echo "BLOCKED: runtimeSessionId not stable" >&2
  exit 2
fi

khr_runtime_log "flight recorder correlation"
jq -e '[.flightRecorder[] | select(.runtimeSessionId != "")] | length > 0' "${RUN_DIR}/heartbeat-loop.json" >/dev/null

khr_runtime_log "stale heartbeat simulation"
STALE_AT="$(date -u -d '10 minutes ago' +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -v-10M +%Y-%m-%dT%H:%M:%SZ)"
"${BIN}" -mode=host-heartbeat -config="${CFG}" -namespace="${NS}" \
  -prior-heartbeat-at="${STALE_AT}" > "${RUN_DIR}/heartbeat-stale.json"
jq -e '.staleDetected == true' "${RUN_DIR}/heartbeat-stale.json" >/dev/null

khr_runtime_log "inventory posture stub validation"
INV_ROOT="$(cd "${ROOT}/../Karl-Inventory" 2>/dev/null && pwd || echo /home/m.simoncini/GitHub/Karl-Inventory)"
python3 -c "
import json, pathlib
stub = pathlib.Path('${INV_ROOT}/docs/contracts/khr/examples/runtime-posture-sandbox-stub.json')
data = json.loads(stub.read_text())
assert data['hostPostureSummary']['sandbox'] is True
assert data['runtimeObservationSummary']['readOnly'] is True
print('inventory posture stub OK')
"

khr_runtime_log "production mutation proof"
khr_runtime_production_mutation_proof > "${RUN_DIR}/production-mutation-proof.json"

jq -n \
  --arg sprint "KHR-N" \
  --arg status "PASS" \
  --arg ctx "${CTX}" \
  --arg ns "${NS}" \
  --arg session "${SID}" \
  --arg runId "${RUN_ID}" \
  --arg at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  '{
    sprint: $sprint,
    status: $status,
    clusterContext: $ctx,
    sandboxNamespace: $ns,
    runtimeSessionId: $session,
    noMutation: true,
    noProductionMutation: true,
    runId: $runId,
    at: $at
  }' > "${RUN_DIR}/summary.json"

cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"
echo "[khr_host_heartbeat_evidence] PASS -> ${EVIDENCE}/summary.json"
