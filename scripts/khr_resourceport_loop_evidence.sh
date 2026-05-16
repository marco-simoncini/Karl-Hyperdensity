#!/usr/bin/env bash
# KHR-J: ResourcePort loop JSON-only evidence on karl-metal-01@ovh.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG_LOOP="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-loop.yaml"
EVIDENCE="${ROOT}/docs/evidence/khr-resourceport-loop"
RUN_ID="${KHR_RESOURCEPORT_LOOP_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
BIN="${ROOT}/bin/karl-host-runtime"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${ROOT}/bin"
(cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime)

khr_runtime_log "resourceport-loop JSON-only (enabled config)"
"${BIN}" -mode=resourceport-loop -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" -loop-iterations=2 -loop-interval-ms=300 \
  > "${RUN_DIR}/loop-pass.json"
jq -e '.blocked == false and .emissionMode == "observed-json"' "${RUN_DIR}/loop-pass.json" >/dev/null

khr_runtime_log "blocked namespace proof"
"${BIN}" -mode=resourceport-loop -config="${CFG_LOOP}" -namespace=karl-system \
  -cluster-context="${CTX}" > "${RUN_DIR}/loop-blocked-namespace.json" || true
jq -e '.blocked == true' "${RUN_DIR}/loop-blocked-namespace.json" >/dev/null

khr_runtime_log "loop disabled by default config"
CFG_OFF="${RUN_DIR}/config-loop-disabled.yaml"
cp "${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-default.yaml" "${CFG_OFF}"
"${BIN}" -mode=resourceport-loop -config="${CFG_OFF}" -namespace="${NS}" \
  -cluster-context="${CTX}" > "${RUN_DIR}/loop-disabled.json"
jq -e '.blocked == true' "${RUN_DIR}/loop-disabled.json" >/dev/null

jq -n \
  --arg sprint "KHR-J" \
  --arg status "PASS" \
  --arg ctx "${CTX}" \
  --arg ns "${NS}" \
  --arg runId "${RUN_ID}" \
  --arg at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  '{
    sprint: $sprint,
    status: $status,
    clusterContext: $ctx,
    namespace: $ns,
    emissionMode: "observed-json",
    emitCR: false,
    noProductionMutation: true,
    noResourceLeaseApply: true,
    runId: $runId,
    at: $at
  }' > "${RUN_DIR}/summary.json"

cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"
echo "[khr_resourceport_loop_evidence] PASS -> ${EVIDENCE}/summary.json"
