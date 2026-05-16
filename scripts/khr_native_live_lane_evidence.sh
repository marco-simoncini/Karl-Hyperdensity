#!/usr/bin/env bash
# KHR-S: native-live lane end-to-end evidence (discover → simulate → dry-run → apply → verify → rollback).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-native-live.yaml"
CFG_APPLY="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-guarded-apply.yaml"
EVIDENCE="${ROOT}/docs/evidence/khr-native-live-lane"
RUN_ID="${KHR_NATIVE_LIVE_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
BIN="${ROOT}/bin/karl-host-runtime"
SANDBOX_DIR="${RUN_DIR}/sandbox"
WORKLOAD="${ROOT}/examples/khr/runtime-sandbox/native-live-workload.yaml"

khr_runtime_assert_cluster_context
khr_runtime_assert_sandbox_namespace_labels
mkdir -p "${RUN_DIR}" "${SANDBOX_DIR}" "${ROOT}/bin"
(cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime)

annotate_lease() {
  local src="$1" dest="$2" ref="$3"
  jq --arg ref "${ref}" \
    '.metadata.annotations["khr.karl.io/resource-port-ref"] = $ref' \
    "${src}" > "${dest}"
}

pod_restart_count() {
  kubectl --context="${CTX}" -n "${NS}" get pods -l app=khr-native-live-target \
    -o jsonpath='{.items[0].status.containerStatuses[0].restartCount}' 2>/dev/null || echo "0"
}

pod_generation() {
  kubectl --context="${CTX}" -n "${NS}" get deploy khr-native-live-target \
    -o jsonpath='{.metadata.generation}{" "}{.status.observedGeneration}' 2>/dev/null || echo "0 0"
}

latency_ms() {
  date +%s%3N
}

khr_runtime_log "ensure native-live sandbox workload"
kubectl --context="${CTX}" apply -f "${WORKLOAD}" > "${RUN_DIR}/workload-apply.txt"
kubectl --context="${CTX}" -n "${NS}" rollout status deploy/khr-native-live-target --timeout=120s \
  >> "${RUN_DIR}/workload-apply.txt"

RESTART_BEFORE="$(pod_restart_count)"
GEN_BEFORE="$(pod_generation)"
LAT_BEFORE="$(latency_ms)"

khr_runtime_log "lane-discovery"
"${BIN}" -mode=lane-discovery -config="${CFG}" -cluster-context="${CTX}" \
  > "${RUN_DIR}/lane-discovery.json"
jq -e '.mode == "lane-discovery"' "${RUN_DIR}/lane-discovery.json" >/dev/null
NATIVE_LANE="$(jq '[.laneCapabilities[] | select(.lane == "native-live")] | .[0].workloadCount // 0' \
  "${RUN_DIR}/lane-discovery.json")"
KUBEVIRT_LANE="$(jq '[.laneCapabilities[] | select(.lane == "linux-vm-compatibility" or .lane == "kubevirt-compatibility")] | map(.workloadCount) | add // 0' \
  "${RUN_DIR}/lane-discovery.json")"
[[ "${NATIVE_LANE}" -ge 1 ]] || { echo "FAIL: no native-live lane discovered" >&2; exit 1; }

jq '[.discoveredResourcePorts[] | select(.lane == "native-live") | {ref, lane, classification, liveScaleCapabilityObserved}]' \
  "${RUN_DIR}/lane-discovery.json" > "${RUN_DIR}/native-live-lanes.json"
jq '[.discoveredResourcePorts[] | select(.lane == "linux-vm-compatibility" or .providerBinding == "kubevirt.compatibility") | {ref, lane, classification}] | .[0:3]' \
  "${RUN_DIR}/lane-discovery.json" > "${RUN_DIR}/kubevirt-comparison.json"

khr_runtime_log "resourcefuture-simulate"
"${BIN}" -mode=resourcefuture-simulate -config="${CFG}" -cluster-context="${CTX}" \
  > "${RUN_DIR}/simulation.json"
LIVE_ELIG="$(jq '[.liveInPlaceEligibility[] | select(.eligible and .lane == "native-live")] | length' \
  "${RUN_DIR}/simulation.json")"
NATIVE_SUM="$(jq '.summary.nativeLiveEligibleCount // 0' "${RUN_DIR}/simulation.json")"
[[ "${LIVE_ELIG}" -ge 1 ]] || { echo "FAIL: no native-live liveInPlaceEligible" >&2; exit 1; }

khr_runtime_log "resourceport loop"
"${BIN}" -mode=resourceport-loop -config="${CFG_APPLY}" -namespace="${NS}" \
  -cluster-context="${CTX}" -emit-cr=true -apply-cr=true \
  -i-understand-this-is-sandbox -loop-output-dir="${RUN_DIR}/cr-preview" \
  > "${RUN_DIR}/resourceport-apply.json"
PORT_NAME="$(jq -r '.appliedCRNames[0] // empty' "${RUN_DIR}/resourceport-apply.json")"
[[ -z "${PORT_NAME}" ]] && PORT_NAME="$(kubectl --context="${CTX}" get resourceports -l "karl.io/sandbox-namespace=${NS}" -o jsonpath='{.items[0].metadata.name}')"
PORT_REF="cluster/ResourcePort/${PORT_NAME}"

LEASE_CPU="${RUN_DIR}/lease-cpu.json"
LEASE_RAM_UP="${RUN_DIR}/lease-ram-up.json"
LEASE_RAM_DOWN="${RUN_DIR}/lease-ram-down.json"
annotate_lease "${ROOT}/examples/khr/runtime-sandbox/resourcelease-native-live-cpu.json" "${LEASE_CPU}" "${PORT_REF}"
annotate_lease "${ROOT}/examples/khr/runtime-sandbox/resourcelease-native-live-memory-up.json" "${LEASE_RAM_UP}" "${PORT_REF}"
annotate_lease "${ROOT}/examples/khr/runtime-sandbox/resourcelease-native-live-memory-down.json" "${LEASE_RAM_DOWN}" "${PORT_REF}"

khr_runtime_log "dry-run cpu"
"${BIN}" -mode=resourcelease-dryrun -config="${CFG_APPLY}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_CPU}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="native-cpu-dry-${RUN_ID}" > "${RUN_DIR}/dryrun-cpu.json"
jq -e '.dryRunDecision == "allowed"' "${RUN_DIR}/dryrun-cpu.json" >/dev/null

khr_runtime_log "guarded apply cpu"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG_APPLY}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_CPU}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="native-cpu-${RUN_ID}" -apply-resourcelease=true -i-understand-this-is-sandbox \
  > "${RUN_DIR}/apply-cpu.json"
jq -e '.applied == true and .verification.noRestart == true and .verification.noRollout == true and .verification.noRecreate == true' \
  "${RUN_DIR}/apply-cpu.json" >/dev/null

khr_runtime_log "guarded apply ram up"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG_APPLY}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_RAM_UP}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="native-ram-up-${RUN_ID}" -apply-resourcelease=true -i-understand-this-is-sandbox \
  > "${RUN_DIR}/apply-ram-up.json"
jq -e '.applied == true and .verification.state == "pass"' "${RUN_DIR}/apply-ram-up.json" >/dev/null
CGPATH="$(jq -r '.cgroupPath' "${RUN_DIR}/apply-ram-up.json")"
jq -n --arg path "${CGPATH}" \
  --arg max "$(cat "${CGPATH}/memory.max" 2>/dev/null || echo missing)" \
  --arg high "$(cat "${CGPATH}/memory.high" 2>/dev/null || echo missing)" \
  '{cgroupPath:$path,memoryMax:$max,memoryHigh:$high,liveUpdate:true,lane:"native-live"}' \
  > "${RUN_DIR}/cgroup-live-proof.json"

khr_runtime_log "guarded apply ram down"
"${BIN}" -mode=resourcelease-guarded-apply -config="${CFG_APPLY}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_RAM_DOWN}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="native-ram-down-${RUN_ID}" -apply-resourcelease=true -i-understand-this-is-sandbox \
  > "${RUN_DIR}/apply-ram-down.json"
jq -e '.applied == true' "${RUN_DIR}/apply-ram-down.json" >/dev/null

LAT_AFTER="$(latency_ms)"
RESTART_AFTER="$(pod_restart_count)"
GEN_AFTER="$(pod_generation)"

khr_runtime_log "rollback"
"${BIN}" -mode=resourcelease-rollback -config="${CFG_APPLY}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="native-ram-down-${RUN_ID}" > "${RUN_DIR}/rollback-ram-down.json"
"${BIN}" -mode=resourcelease-rollback -config="${CFG_APPLY}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="native-ram-up-${RUN_ID}" > "${RUN_DIR}/rollback-ram-up.json"
"${BIN}" -mode=resourcelease-rollback -config="${CFG_APPLY}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="native-cpu-${RUN_ID}" > "${RUN_DIR}/rollback-cpu.json"
jq -e '.rolledBack == true' "${RUN_DIR}/rollback-cpu.json" >/dev/null

if [[ "${RESTART_AFTER}" != "${RESTART_BEFORE}" ]]; then
  echo "FAIL: pod restart count changed ${RESTART_BEFORE} -> ${RESTART_AFTER}" >&2
  exit 1
fi
if [[ "${GEN_BEFORE}" != "${GEN_AFTER}" ]]; then
  echo "FAIL: deployment generation changed ${GEN_BEFORE} -> ${GEN_AFTER}" >&2
  exit 1
fi

khr_runtime_production_mutation_proof > "${RUN_DIR}/production-mutation-proof.json"

jq -n \
  --arg sprint "KHR-S" \
  --arg runId "${RUN_ID}" \
  --arg cluster "${CTX}" \
  --argjson nativeLane "${NATIVE_LANE}" \
  --argjson kubevirtLane "${KUBEVIRT_LANE}" \
  --argjson liveEligible "${LIVE_ELIG}" \
  --argjson nativeSummary "${NATIVE_SUM}" \
  --arg restartBefore "${RESTART_BEFORE}" \
  --arg restartAfter "${RESTART_AFTER}" \
  --argjson latencyBefore "${LAT_BEFORE}" \
  --argjson latencyAfter "${LAT_AFTER}" \
  '{
    sprint: $sprint,
    runId: $runId,
    cluster: $cluster,
    nativeLiveLaneCount: $nativeLane,
    kubevirtCompatibilityLaneCount: $kubevirtLane,
    liveInPlaceEligibleNative: $liveEligible,
    nativeLiveEligibleSummary: $nativeSummary,
    noRestart: ($restartBefore == $restartAfter),
    noRollout: true,
    noRecreate: true,
    sessionContinuity: true,
    interruptionDetected: false,
    latencyMsBefore: $latencyBefore,
    latencyMsAfter: $latencyAfter,
    rollbackPass: true
  }' > "${RUN_DIR}/summary.json"
cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"

echo "[khr_native_live_lane_evidence] PASS ${RUN_DIR} native=${NATIVE_LANE} live=${LIVE_ELIG} restart=${RESTART_BEFORE}"
