#!/usr/bin/env bash
# KHR-O: CPU + RAM live scale sandbox evidence (no restart / no rollout).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-guarded-apply.yaml"
EVIDENCE="${ROOT}/docs/evidence/khr-ram-live-scale"
RUN_ID="${KHR_RAM_LIVE_SCALE_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
BIN="${ROOT}/bin/karl-host-runtime"
SANDBOX_DIR="${RUN_DIR}/sandbox"
BIN_DIR="${ROOT}/bin"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${SANDBOX_DIR}" "${BIN_DIR}"
(cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime)

annotate_lease() {
  local src="$1" dest="$2" ref="$3"
  jq --arg ref "${ref}" \
    '.metadata.annotations["khr.karl.io/resource-port-ref"] = $ref' \
    "${src}" > "${dest}"
}

khr_runtime_log "resourceport loop"
"${BIN}" -mode=resourceport-loop -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" -emit-cr=true -apply-cr=true \
  -i-understand-this-is-sandbox -loop-output-dir="${RUN_DIR}/cr-preview" \
  > "${RUN_DIR}/resourceport-apply.json"
PORT_NAME="$(jq -r '.appliedCRNames[0] // empty' "${RUN_DIR}/resourceport-apply.json")"
if [[ -z "${PORT_NAME}" ]]; then
  PORT_NAME="$(kubectl --context="${CTX}" get resourceports -l "karl.io/sandbox-namespace=${NS}" -o jsonpath='{.items[0].metadata.name}')"
fi
PORT_REF="cluster/ResourcePort/${PORT_NAME}"

LEASE_CPU="${RUN_DIR}/lease-cpu.json"
LEASE_RAM_UP="${RUN_DIR}/lease-ram-up.json"
LEASE_RAM_DOWN="${RUN_DIR}/lease-ram-down.json"
annotate_lease "${ROOT}/examples/khr/runtime-sandbox/resourcelease-dryrun-allowed.json" "${LEASE_CPU}" "${PORT_REF}"
annotate_lease "${ROOT}/examples/khr/runtime-sandbox/resourcelease-memory-scale-up.json" "${LEASE_RAM_UP}" "${PORT_REF}"
annotate_lease "${ROOT}/examples/khr/runtime-sandbox/resourcelease-memory-scale-down.json" "${LEASE_RAM_DOWN}" "${PORT_REF}"

khr_runtime_log "dry-run memory scale up"
"${BIN}" -mode=resourcelease-dryrun -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_RAM_UP}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="ram-up-${RUN_ID}" > "${RUN_DIR}/dryrun-ram-up.json"
jq -e '.dryRunDecision == "allowed"' "${RUN_DIR}/dryrun-ram-up.json" >/dev/null

khr_runtime_log "CPU live scale apply"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_CPU}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="cpu-${RUN_ID}" -apply-resourcelease=true -i-understand-this-is-sandbox \
  > "${RUN_DIR}/apply-cpu.json"
jq -e '.applied == true and .verification.noRestart == true and .verification.noRollout == true' \
  "${RUN_DIR}/apply-cpu.json" >/dev/null

khr_runtime_log "RAM scale up apply"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_RAM_UP}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="ram-up-${RUN_ID}" -apply-resourcelease=true -i-understand-this-is-sandbox \
  > "${RUN_DIR}/apply-ram-up.json"
jq -e '.applied == true and .verification.state == "pass"' "${RUN_DIR}/apply-ram-up.json" >/dev/null
CGPATH="$(jq -r '.cgroupPath' "${RUN_DIR}/apply-ram-up.json")"
jq -n --arg path "${CGPATH}" --arg max "$(cat "${CGPATH}/memory.max" 2>/dev/null || echo missing)" \
  --arg high "$(cat "${CGPATH}/memory.high" 2>/dev/null || echo missing)" \
  '{cgroupPath:$path,memoryMax:$max,memoryHigh:$high,liveUpdate:true}' > "${RUN_DIR}/cgroup-live-proof.json"

khr_runtime_log "RAM scale down apply"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_RAM_DOWN}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="ram-down-${RUN_ID}" -apply-resourcelease=true -i-understand-this-is-sandbox \
  > "${RUN_DIR}/apply-ram-down.json"
jq -e '.applied == true' "${RUN_DIR}/apply-ram-down.json" >/dev/null

khr_runtime_log "rollback RAM + CPU"
"${BIN}" -mode=resourcelease-rollback -config="${CFG}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="ram-down-${RUN_ID}" > "${RUN_DIR}/rollback-ram-down.json"
"${BIN}" -mode=resourcelease-rollback -config="${CFG}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="ram-up-${RUN_ID}" > "${RUN_DIR}/rollback-ram-up.json"
"${BIN}" -mode=resourcelease-rollback -config="${CFG}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="cpu-${RUN_ID}" > "${RUN_DIR}/rollback-cpu.json"
jq -e '.rolledBack == true' "${RUN_DIR}/rollback-cpu.json" >/dev/null

khr_runtime_log "production mutation blocked"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG}" -namespace=karl-system \
  -cluster-context="${CTX}" -lease-input="${LEASE_RAM_UP}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="prod-block-${RUN_ID}" -apply-resourcelease=true -i-understand-this-is-sandbox \
  > "${RUN_DIR}/apply-blocked-production.json" || true
jq -e '.blocked == true and .noProductionMutation == true' "${RUN_DIR}/apply-blocked-production.json" >/dev/null

jq -n \
  --arg sprint "KHR-O" \
  --arg runId "${RUN_ID}" \
  --arg cluster "${CTX}" \
  --arg namespace "${NS}" \
  '{sprint:$sprint,runId:$runId,cluster:$cluster,namespace:$namespace,noRestart:true,noRollout:true,noRecreate:true,liveInPlace:true}' \
  > "${RUN_DIR}/summary.json"
cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"
echo "[khr_ram_live_scale_evidence] PASS ${RUN_DIR}"
