#!/usr/bin/env bash
# KHR-W: operator action approval workflow evidence (read-only; no apply).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
EVIDENCE="${ROOT}/docs/evidence/khr-action-approval"
RUN_ID="${KHR_ACTION_APPROVAL_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
REGISTRY="${ROOT}/docs/evidence/khr-certification-registry/registry.json"
CERT_SUMMARY="${ROOT}/docs/evidence/khr-native-live-lane/certification-summary.json"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-resourcefuture-simulate.yaml"
BIN_HOST="${ROOT}/bin/karl-host-runtime"
BIN_APPROVAL="${ROOT}/bin/khr-action-approval"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${ROOT}/bin"

if [[ ! -f "${REGISTRY}" ]]; then
  echo "[khr_action_approval_evidence] run khr_cert_registry_policy_gates.sh first" >&2
  exit 1
fi

(cd "${ROOT}" && go build -o "${BIN_APPROVAL}" ./cmd/khr-action-approval)
(cd "${ROOT}" && go build -o "${BIN_HOST}" ./cmd/karl-host-runtime)

khr_runtime_log "resourcefuture-simulate (certified native-live eligible)"
"${BIN_HOST}" -mode=resourcefuture-simulate -config="${CFG}" -cluster-context="${CTX}" \
  -cert-registry="${REGISTRY}" > "${RUN_DIR}/simulation-gated.json"

"${BIN_APPROVAL}" -cmd=generate \
  -simulation="${RUN_DIR}/simulation-gated.json" \
  -registry="${REGISTRY}" \
  -cert-ref="${CERT_SUMMARY}" \
  -out="${RUN_DIR}/pending-bundle.json"

PENDING="$(jq -r '.approvals[0]' "${RUN_DIR}/pending-bundle.json")"
echo "${PENDING}" > "${RUN_DIR}/pending.json"

"${BIN_APPROVAL}" -cmd=approve \
  -approval="${RUN_DIR}/pending.json" \
  -registry="${REGISTRY}" \
  -by="operator-sandbox-a" \
  -out="${RUN_DIR}/approved.json"

# second pending for reject path
PENDING2="$(echo "${PENDING}" | jq '.actionId = .actionId + "-reject"')"
echo "${PENDING2}" > "${RUN_DIR}/pending-reject.json"
"${BIN_APPROVAL}" -cmd=reject \
  -approval="${RUN_DIR}/pending-reject.json" \
  -by="operator-sandbox-b" \
  -reason="sandbox reject evidence" \
  -out="${RUN_DIR}/rejected.json"

PENDING3="$(echo "${PENDING}" | jq '.actionId = .actionId + "-expire"')"
echo "${PENDING3}" > "${RUN_DIR}/pending-expire.json"
"${BIN_APPROVAL}" -cmd=expire \
  -approval="${RUN_DIR}/pending-expire.json" \
  -out="${RUN_DIR}/expired.json"

STALE_REG="${RUN_DIR}/registry-stale.json"
cp "${REGISTRY}" "${STALE_REG}"
jq '.entries[0].lastCertifiedAt = "2019-01-01T00:00:00Z"' "${STALE_REG}" > "${STALE_REG}.tmp" && mv "${STALE_REG}.tmp" "${STALE_REG}"

PENDING4="$(echo "${PENDING}" | jq '.actionId = .actionId + "-stale"')"
echo "${PENDING4}" > "${RUN_DIR}/pending-stale.json"
if "${BIN_APPROVAL}" -cmd=approve -approval="${RUN_DIR}/pending-stale.json" \
  -registry="${STALE_REG}" -by="operator-sandbox-c" -out="${RUN_DIR}/stale-approved.json" 2>/dev/null; then
  echo "[khr_action_approval_evidence] FAIL: stale approval should be blocked" >&2
  exit 1
fi
echo '{"blocked":true,"reason":"stale certification blocks approval"}' > "${RUN_DIR}/stale-blocked.json"

echo "no-apply: action approval evidence only" > "${RUN_DIR}/mutation-check.txt"
echo "no-kubectl-apply: no cluster mutation" >> "${RUN_DIR}/mutation-check.txt"

jq -n \
  --arg sprint "KHR-W" \
  --arg runId "${RUN_ID}" \
  --arg cluster "${CTX}" \
  '{sprint:$sprint,runId:$runId,cluster:$cluster,pending:true,approved:true,rejected:true,expired:true,staleBlocked:true,readOnly:true,noApply:true,noAutonomousOrchestration:true}' \
  > "${RUN_DIR}/summary.json"
cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"

echo "[khr_action_approval_evidence] PASS ${RUN_DIR}"
