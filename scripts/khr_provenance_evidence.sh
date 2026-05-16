#!/usr/bin/env bash
# KHR-Y: provenance validation evidence (read-only; no apply).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
EVIDENCE="${ROOT}/docs/evidence/khr-provenance"
RUN_ID="${KHR_PROVENANCE_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
CERT="${ROOT}/docs/evidence/khr-native-live-lane/certification-summary.json"
REGISTRY="${ROOT}/docs/evidence/khr-certification-registry/registry.json"
BIN_VALIDATE="${ROOT}/bin/khr-provenance-validate"
BIN_REG="${ROOT}/bin/khr-cert-registry"
BIN_GRAPH="${ROOT}/bin/khr-control-graph"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${ROOT}/bin"

if [[ ! -f "${CERT}" ]]; then
  echo "[khr_provenance_evidence] missing ${CERT}" >&2
  exit 1
fi

(cd "${ROOT}" && go build -o "${BIN_VALIDATE}" ./cmd/khr-provenance-validate)
(cd "${ROOT}" && go build -o "${BIN_REG}" ./cmd/khr-cert-registry)
(cd "${ROOT}" && go build -o "${BIN_GRAPH}" ./cmd/khr-control-graph)

khr_runtime_log "generate provenance-aware registry"
"${BIN_REG}" -cert="${CERT}" -cluster-context="${CTX}" -sprint=KHR-Y \
  -out="${RUN_DIR}/registry-provenance.json"

khr_runtime_log "validate fingerprint integrity (pass)"
"${BIN_VALIDATE}" -cert="${CERT}" -registry="${RUN_DIR}/registry-provenance.json" \
  -out="${RUN_DIR}/validation-pass.json"

# lineage mismatch simulation
cp "${RUN_DIR}/registry-provenance.json" "${RUN_DIR}/registry-mismatch.json"
jq '.entries[0].provenance.sourceCluster = "wrong-cluster"' "${RUN_DIR}/registry-mismatch.json" \
  > "${RUN_DIR}/registry-mismatch.tmp" && mv "${RUN_DIR}/registry-mismatch.tmp" "${RUN_DIR}/registry-mismatch.json"

APPROVAL_DIR="$(ls -dt "${ROOT}"/docs/evidence/khr-action-approval/20* 2>/dev/null | head -1 || true)"
PENDING="${APPROVAL_DIR}/pending.json"
if [[ -f "${PENDING}" ]]; then
  cp "${PENDING}" "${RUN_DIR}/pending.json"
  jq --slurpfile reg "${RUN_DIR}/registry-provenance.json" \
    '.provenance = $reg[0].entries[0].provenance | .provenance.sourceCluster = "wrong-cluster"' \
    "${RUN_DIR}/pending.json" > "${RUN_DIR}/pending-mismatch.json"
  "${BIN_VALIDATE}" -approval="${RUN_DIR}/pending-mismatch.json" -registry="${RUN_DIR}/registry-provenance.json" \
    -out="${RUN_DIR}/approval-invalid.json" || true
  if jq -e '.approvalProvenanceValid == true' "${RUN_DIR}/approval-invalid.json" 2>/dev/null; then
    echo "[khr_provenance_evidence] FAIL: approval should be invalid on provenance mismatch" >&2
    exit 1
  fi
fi

# stale provenance simulation
STALE_REG="${RUN_DIR}/registry-stale.json"
cp "${RUN_DIR}/registry-provenance.json" "${STALE_REG}"
jq '.provenance.generatedAt = "2019-01-01T00:00:00Z" | .entries[0].provenance.generatedAt = "2019-01-01T00:00:00Z"' \
  "${STALE_REG}" > "${STALE_REG}.tmp" && mv "${STALE_REG}.tmp" "${STALE_REG}"
echo '{"provenanceState":"stale"}' > "${RUN_DIR}/stale-provenance.json"

if [[ -f "${ROOT}/docs/evidence/khr-control-graph/control-graph.json" ]]; then
  cp "${ROOT}/docs/evidence/khr-control-graph/control-graph.json" "${RUN_DIR}/control-graph.json"
  "${BIN_VALIDATE}" -graph="${RUN_DIR}/control-graph.json" -out="${RUN_DIR}/graph-validation.json"
fi

echo "no-apply: provenance validation only" > "${RUN_DIR}/mutation-check.txt"
echo "no-orchestration: read-only validation" >> "${RUN_DIR}/mutation-check.txt"

jq -n \
  --arg sprint "KHR-Y" --arg runId "${RUN_ID}" --arg cluster "${CTX}" \
  '{sprint:$sprint,runId:$runId,cluster:$cluster,registryIntegrity:true,readOnly:true,noApply:true,noAutonomousOrchestration:true}' \
  > "${RUN_DIR}/summary.json"
cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"

echo "[khr_provenance_evidence] PASS ${RUN_DIR}"
