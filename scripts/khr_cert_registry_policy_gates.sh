#!/usr/bin/env bash
# KHR-V: certification registry + gated ResourceFuture simulation (read-only).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
EVIDENCE="${ROOT}/docs/evidence/khr-certification-registry"
RUN_ID="${KHR_CERT_REGISTRY_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
CERT_SUMMARY="${ROOT}/docs/evidence/khr-native-live-lane/certification-summary.json"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-resourcefuture-simulate.yaml"
BIN_HOST="${ROOT}/bin/karl-host-runtime"
BIN_REG="${ROOT}/bin/khr-cert-registry"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${ROOT}/bin"

if [[ ! -f "${CERT_SUMMARY}" ]]; then
  echo "[khr_cert_registry_policy_gates] missing ${CERT_SUMMARY}; run khr_native_live_certify.sh first" >&2
  exit 1
fi

(cd "${ROOT}" && go build -o "${BIN_REG}" ./cmd/khr-cert-registry)
(cd "${ROOT}" && go build -o "${BIN_HOST}" ./cmd/karl-host-runtime)

CERT_DIR="$(jq -r '.latestCertificationDir // empty' "${ROOT}/docs/evidence/khr-native-live-lane/certification-run-summary.json" 2>/dev/null || true)"
EVIDENCE_REF="${CERT_DIR:-docs/evidence/khr-native-live-lane/certification-summary.json}"

"${BIN_REG}" -cert="${CERT_SUMMARY}" -evidence-ref="${EVIDENCE_REF}" -sprint=KHR-V \
  -out="${RUN_DIR}/registry.json"

cp "${RUN_DIR}/registry.json" "${EVIDENCE}/registry.json"

khr_runtime_log "resourcefuture-simulate with certification registry (fresh)"
"${BIN_HOST}" -mode=resourcefuture-simulate -config="${CFG}" -cluster-context="${CTX}" \
  -cert-registry="${RUN_DIR}/registry.json" > "${RUN_DIR}/simulation-gated.json"

jq -e '.safety.noApply == true and .safety.noMutation == true' "${RUN_DIR}/simulation-gated.json" >/dev/null
jq -e '[.liveInPlaceEligibility[] | select(.lane=="native-live" and .eligible)] | length >= 1' \
  "${RUN_DIR}/simulation-gated.json" >/dev/null

STALE_REG="${RUN_DIR}/registry-stale.json"
cp "${RUN_DIR}/registry.json" "${STALE_REG}"
jq '.entries[0].lastCertifiedAt = "2019-01-01T00:00:00Z"' "${STALE_REG}" > "${STALE_REG}.tmp" && mv "${STALE_REG}.tmp" "${STALE_REG}"

khr_runtime_log "resourcefuture-simulate stale evidence (blocked)"
"${BIN_HOST}" -mode=resourcefuture-simulate -config="${CFG}" -cluster-context="${CTX}" \
  -cert-registry="${STALE_REG}" > "${RUN_DIR}/simulation-stale.json"

jq -e '[.liveInPlaceEligibility[] | select(.lane=="native-live" and .staleEvidence)] | length >= 1' \
  "${RUN_DIR}/simulation-stale.json" >/dev/null

khr_runtime_log "resourcefuture-simulate without registry (legacy uncertified compatibility)"
"${BIN_HOST}" -mode=resourcefuture-simulate -config="${CFG}" -cluster-context="${CTX}" \
  > "${RUN_DIR}/simulation-legacy.json"

jq -e '[.liveInPlaceEligibility[] | select(.uncertifiedLane==true)] | length >= 0' \
  "${RUN_DIR}/simulation-legacy.json" >/dev/null

echo "no-apply: resourcefuture-simulate only" > "${RUN_DIR}/mutation-check.txt"
echo "no-patch: no kubectl apply" >> "${RUN_DIR}/mutation-check.txt"

jq -n \
  --arg sprint "KHR-V" \
  --arg runId "${RUN_ID}" \
  --arg cluster "${CTX}" \
  --arg registry "${EVIDENCE_REF}" \
  '{sprint:$sprint,runId:$runId,cluster:$cluster,registryEvidenceRef:$registry,readOnly:true,noApply:true,noAutonomousOrchestration:true}' \
  > "${RUN_DIR}/summary.json"
cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"

echo "[khr_cert_registry_policy_gates] PASS ${RUN_DIR}"
