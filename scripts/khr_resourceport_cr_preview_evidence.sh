#!/usr/bin/env bash
# KHR-K: ResourcePort CR preview + sandbox apply evidence on karl-metal-01@ovh.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG_LOOP="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-loop.yaml"
CRD="${ROOT}/api/crds/runtime.karl.io/resourceport.yaml"
EVIDENCE="${ROOT}/docs/evidence/khr-resourceport-cr-preview"
RUN_ID="${KHR_RESOURCEPORT_CR_PREVIEW_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
BIN="${ROOT}/bin/karl-host-runtime"
PREVIEW_DIR="${RUN_DIR}/cr-preview"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${PREVIEW_DIR}" "${ROOT}/bin"
(cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime)

khr_runtime_log "ensure ResourcePort CRD (cluster-scoped)"
kubectl --context="${CTX}" apply -f "${CRD}" > "${RUN_DIR}/crd-apply.txt"

khr_runtime_log "1) JSON-only loop"
"${BIN}" -mode=resourceport-loop -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" > "${RUN_DIR}/loop-json-only.json"
jq -e '.emissionMode == "observed-json" and .applyCR == false' "${RUN_DIR}/loop-json-only.json" >/dev/null

khr_runtime_log "2) CR preview (local files)"
"${BIN}" -mode=resourceport-loop -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" -emit-cr=true -loop-output-dir="${PREVIEW_DIR}" \
  > "${RUN_DIR}/loop-cr-preview.json"
jq -e '.emissionMode == "cr-preview" and .emitCRApplied == true' "${RUN_DIR}/loop-cr-preview.json" >/dev/null
test -n "$(find "${PREVIEW_DIR}" -name 'resourceport-*.json' | head -1)"

khr_runtime_log "3) apply blocked without confirmation"
"${BIN}" -mode=resourceport-loop -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" -emit-cr=true -apply-cr=true -loop-output-dir="${PREVIEW_DIR}" \
  > "${RUN_DIR}/loop-apply-blocked.json"
jq -e '.applyCRBlocked == true' "${RUN_DIR}/loop-apply-blocked.json" >/dev/null

khr_runtime_log "4) sandbox apply"
"${BIN}" -mode=resourceport-loop -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" -emit-cr=true -apply-cr=true \
  -i-understand-this-is-sandbox -loop-output-dir="${PREVIEW_DIR}" \
  > "${RUN_DIR}/loop-apply.json"
jq -e '.applyCRApplied == true and .emissionMode == "cr-applied-sandbox"' "${RUN_DIR}/loop-apply.json" >/dev/null

khr_runtime_log "5) verify resourceports"
kubectl --context="${CTX}" get resourceports -l karl.io/managed-by=karl-host-runtime \
  > "${RUN_DIR}/get-resourceports.txt"

khr_runtime_log "6) cleanup"
"${BIN}" -mode=resourceport-cleanup -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" > "${RUN_DIR}/cleanup.json"

khr_runtime_log "production namespace mutation proof"
khr_runtime_production_mutation_proof > "${RUN_DIR}/production-mutation-proof.json"

jq -n \
  --arg sprint "KHR-K" \
  --arg status "PASS" \
  --arg ctx "${CTX}" \
  --arg ns "${NS}" \
  --arg runId "${RUN_ID}" \
  --arg at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  '{
    sprint: $sprint,
    status: $status,
    clusterContext: $ctx,
    sandboxNamespace: $ns,
    emissionModes: ["observed-json", "cr-preview", "cr-applied-sandbox"],
    noProductionMutation: true,
    noResourceLeaseApply: true,
    runId: $runId,
    at: $at
  }' > "${RUN_DIR}/summary.json"

cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"
echo "[khr_resourceport_cr_preview_evidence] PASS -> ${EVIDENCE}/summary.json"
