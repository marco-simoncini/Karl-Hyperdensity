#!/usr/bin/env bash
# KHR-M: ResourceLease guarded apply + rollback evidence on karl-metal-01@ovh.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-guarded-apply.yaml"
CRD="${ROOT}/api/crds/runtime.karl.io/resourceport.yaml"
EVIDENCE="${ROOT}/docs/evidence/khr-resourcelease-guarded-apply"
RUN_ID="${KHR_RESOURCELEASE_GUARDED_APPLY_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
BIN="${ROOT}/bin/karl-host-runtime"
SANDBOX_DIR="${RUN_DIR}/sandbox"
PREVIEW_DIR="${RUN_DIR}/cr-preview"
LEASE_ALLOWED="${RUN_DIR}/lease-allowed.json"
BASELINE_ID="khr-m-${RUN_ID}"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${SANDBOX_DIR}" "${PREVIEW_DIR}" "${ROOT}/bin"
(cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime)

khr_runtime_log "ensure ResourcePort CRD"
kubectl --context="${CTX}" apply -f "${CRD}" > "${RUN_DIR}/crd-apply.txt"

khr_runtime_log "apply sandbox ResourcePort CRs"
"${BIN}" -mode=resourceport-loop -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" -emit-cr=true -apply-cr=true \
  -i-understand-this-is-sandbox -loop-output-dir="${PREVIEW_DIR}" \
  > "${RUN_DIR}/resourceport-apply.json"

PORT_NAME="$(jq -r '.appliedCRNames[0] // empty' "${RUN_DIR}/resourceport-apply.json")"
if [[ -z "${PORT_NAME}" ]]; then
  PORT_NAME="$(kubectl --context="${CTX}" get resourceports -l "karl.io/sandbox-namespace=${NS}" -o jsonpath='{.items[0].metadata.name}')"
fi
jq --arg ref "cluster/ResourcePort/${PORT_NAME}" \
  '.metadata.annotations["khr.karl.io/resource-port-ref"] = $ref' \
  "${ROOT}/examples/khr/runtime-sandbox/resourcelease-dryrun-allowed.json" > "${LEASE_ALLOWED}"

khr_runtime_log "dry-run allowed"
"${BIN}" -mode=resourcelease-dryrun -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_ALLOWED}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="${BASELINE_ID}" > "${RUN_DIR}/dryrun-allowed.json"
jq -e '.dryRunDecision == "allowed"' "${RUN_DIR}/dryrun-allowed.json" >/dev/null

khr_runtime_log "guarded apply blocked without confirmation"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_ALLOWED}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="${BASELINE_ID}" -apply-resourcelease=true \
  > "${RUN_DIR}/apply-blocked-no-confirm.json"
jq -e '.blocked == true' "${RUN_DIR}/apply-blocked-no-confirm.json" >/dev/null

khr_runtime_log "guarded apply blocked missing rollbackPlanRef"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" \
  -lease-input="${ROOT}/examples/khr/runtime-sandbox/resourcelease-guarded-apply-blocked-no-rollback.json" \
  -sandbox-dir="${SANDBOX_DIR}" -baseline-id="${BASELINE_ID}" \
  -apply-resourcelease=true -i-understand-this-is-sandbox \
  > "${RUN_DIR}/apply-blocked-no-rollback.json"
jq -e '.blocked == true' "${RUN_DIR}/apply-blocked-no-rollback.json" >/dev/null

khr_runtime_log "guarded apply blocked over cap"
OVER="${RUN_DIR}/lease-over.json"
jq '.spec.transfer.amount.milliCpu = 5000' "${LEASE_ALLOWED}" > "${OVER}"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${OVER}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="${BASELINE_ID}" -apply-resourcelease=true -i-understand-this-is-sandbox \
  > "${RUN_DIR}/apply-blocked-over-cap.json"
jq -e '.blocked == true' "${RUN_DIR}/apply-blocked-over-cap.json" >/dev/null

khr_runtime_log "guarded apply blocked production namespace"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG}" -namespace=karl-system \
  -cluster-context="${CTX}" -lease-input="${LEASE_ALLOWED}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="${BASELINE_ID}" -apply-resourcelease=true -i-understand-this-is-sandbox \
  > "${RUN_DIR}/apply-blocked-namespace.json" || true
jq -e '.blocked == true' "${RUN_DIR}/apply-blocked-namespace.json" >/dev/null

khr_runtime_log "guarded apply CPU sandbox"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_ALLOWED}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="${BASELINE_ID}" -apply-resourcelease=true -i-understand-this-is-sandbox \
  > "${RUN_DIR}/apply-pass.json"
jq -e '.applied == true and .verification.state == "pass"' "${RUN_DIR}/apply-pass.json" >/dev/null

khr_runtime_log "rollback baseline"
"${BIN}" -mode=resourcelease-rollback -config="${CFG}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="${BASELINE_ID}" > "${RUN_DIR}/rollback.json"
jq -e '.rolledBack == true and .verification.state == "pass"' "${RUN_DIR}/rollback.json" >/dev/null

khr_runtime_log "cleanup ResourcePorts"
"${BIN}" -mode=resourceport-cleanup -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" > "${RUN_DIR}/resourceport-cleanup.json"

khr_runtime_log "production mutation proof"
khr_runtime_production_mutation_proof > "${RUN_DIR}/production-mutation-proof.json"

jq -n \
  --arg sprint "KHR-M" \
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
    noResourceLeaseClusterApply: true,
    noProductionMutation: true,
    runId: $runId,
    at: $at
  }' > "${RUN_DIR}/summary.json"

cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"
echo "[khr_resourcelease_guarded_apply_evidence] PASS -> ${EVIDENCE}/summary.json"
