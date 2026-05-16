#!/usr/bin/env bash
# KHR-L: ResourceLease dry-run against sandbox ResourcePort CRs on karl-metal-01@ovh.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG_LOOP="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-loop.yaml"
CRD="${ROOT}/api/crds/runtime.karl.io/resourceport.yaml"
EVIDENCE="${ROOT}/docs/evidence/khr-resourcelease-dryrun"
RUN_ID="${KHR_RESOURCELEASE_DRYRUN_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
BIN="${ROOT}/bin/karl-host-runtime"
PREVIEW_DIR="${RUN_DIR}/cr-preview"
LEASE_ALLOWED="${RUN_DIR}/lease-allowed.json"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${PREVIEW_DIR}" "${ROOT}/bin"
(cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime)

khr_runtime_log "ensure ResourcePort CRD"
kubectl --context="${CTX}" apply -f "${CRD}" > "${RUN_DIR}/crd-apply.txt"

khr_runtime_log "apply sandbox ResourcePort CRs"
"${BIN}" -mode=resourceport-loop -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" -emit-cr=true -apply-cr=true \
  -i-understand-this-is-sandbox -loop-output-dir="${PREVIEW_DIR}" \
  > "${RUN_DIR}/resourceport-apply.json"

PORT_NAME="$(jq -r '.appliedCRNames[0] // empty' "${RUN_DIR}/resourceport-apply.json")"
if [[ -z "${PORT_NAME}" ]]; then
  PORT_NAME="$(kubectl --context="${CTX}" get resourceports -l "karl.io/sandbox-namespace=${NS}" -o jsonpath='{.items[0].metadata.name}')"
fi
if [[ -z "${PORT_NAME}" ]]; then
  echo "BLOCKED: no sandbox ResourcePort found" >&2
  exit 2
fi

jq --arg ref "cluster/ResourcePort/${PORT_NAME}" \
  '.metadata.annotations["khr.karl.io/resource-port-ref"] = $ref' \
  "${ROOT}/examples/khr/runtime-sandbox/resourcelease-dryrun-allowed.json" > "${LEASE_ALLOWED}"

khr_runtime_log "ResourceLease dry-run allowed"
"${BIN}" -mode=resourcelease-dryrun -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_ALLOWED}" \
  > "${RUN_DIR}/dryrun-allowed.json"
jq -e '.dryRunDecision == "allowed" and .noApply == true' "${RUN_DIR}/dryrun-allowed.json" >/dev/null

khr_runtime_log "blocked: production namespace"
"${BIN}" -mode=resourcelease-dryrun -config="${CFG_LOOP}" -namespace=karl-system \
  -cluster-context="${CTX}" -lease-input="${LEASE_ALLOWED}" \
  > "${RUN_DIR}/dryrun-blocked-namespace.json" || true
jq -e '.blocked == true' "${RUN_DIR}/dryrun-blocked-namespace.json" >/dev/null

khr_runtime_log "blocked: label missing"
"${BIN}" -mode=resourcelease-dryrun -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" \
  -lease-input="${ROOT}/examples/khr/runtime-sandbox/resourcelease-dryrun-blocked-label.json" \
  > "${RUN_DIR}/dryrun-blocked-label.json"
jq -e '.blocked == true' "${RUN_DIR}/dryrun-blocked-label.json" >/dev/null

khr_runtime_log "blocked: over limit"
"${BIN}" -mode=resourcelease-dryrun -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" \
  -lease-input="${ROOT}/examples/khr/runtime-sandbox/resourcelease-dryrun-blocked-over-limit.json" \
  > "${RUN_DIR}/dryrun-blocked-over-limit.json"
jq -e '.blocked == true' "${RUN_DIR}/dryrun-blocked-over-limit.json" >/dev/null

khr_runtime_log "blocked: missing ResourcePort"
"${BIN}" -mode=resourcelease-dryrun -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" \
  -lease-input="${ROOT}/examples/khr/runtime-sandbox/resourcelease-dryrun-blocked-missing-port.json" \
  > "${RUN_DIR}/dryrun-blocked-missing-port.json"
jq -e '.blocked == true' "${RUN_DIR}/dryrun-blocked-missing-port.json" >/dev/null

khr_runtime_log "cleanup ResourcePorts"
"${BIN}" -mode=resourceport-cleanup -config="${CFG_LOOP}" -namespace="${NS}" \
  -cluster-context="${CTX}" > "${RUN_DIR}/resourceport-cleanup.json"

khr_runtime_log "production mutation proof"
khr_runtime_production_mutation_proof > "${RUN_DIR}/production-mutation-proof.json"

jq -n \
  --arg sprint "KHR-L" \
  --arg status "PASS" \
  --arg ctx "${CTX}" \
  --arg ns "${NS}" \
  --arg port "${PORT_NAME}" \
  --arg runId "${RUN_ID}" \
  --arg at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  '{
    sprint: $sprint,
    status: $status,
    clusterContext: $ctx,
    sandboxNamespace: $ns,
    matchedResourcePort: $port,
    noResourceLeaseApply: true,
    noProductionMutation: true,
    runId: $runId,
    at: $at
  }' > "${RUN_DIR}/summary.json"

cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"
echo "[khr_resourcelease_dryrun_evidence] PASS -> ${EVIDENCE}/summary.json"
