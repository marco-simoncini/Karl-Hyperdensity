#!/usr/bin/env bash
# KHR-T/U: single native-live lane evidence run; writes run-metrics.json (resource + shell continuity).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-native-live.yaml"
CFG_APPLY="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-guarded-apply.yaml"
BIN="${ROOT}/bin/karl-host-runtime"
CERTIFY_BIN="${ROOT}/bin/khr-native-live-certify"
CONTINUITY_BIN="${ROOT}/bin/khr-continuity-proof"
SNAPSHOT_BIN="${ROOT}/bin/khr-continuity-snapshot"
WORKLOAD="${ROOT}/examples/khr/runtime-sandbox/native-live-workload.yaml"

RUN_DIR="${1:?run output directory required}"
SANDBOX_DIR="${RUN_DIR}/sandbox"
mkdir -p "${RUN_DIR}" "${SANDBOX_DIR}" "${ROOT}/bin"
(cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime)
(cd "${ROOT}" && go build -o "${CERTIFY_BIN}" ./cmd/khr-native-live-certify)
(cd "${ROOT}" && go build -o "${CONTINUITY_BIN}" ./cmd/khr-continuity-proof)
(cd "${ROOT}" && go build -o "${SNAPSHOT_BIN}" ./cmd/khr-continuity-snapshot)

capture_continuity_snapshot() {
  "${SNAPSHOT_BIN}" -cluster-context="${CTX}" -namespace="${NS}" -out="$1"
}

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

deploy_rollout_detected() {
  local gen obs
  gen="$(kubectl --context="${CTX}" -n "${NS}" get deploy khr-native-live-target \
    -o jsonpath='{.metadata.generation}' 2>/dev/null || echo "0")"
  obs="$(kubectl --context="${CTX}" -n "${NS}" get deploy khr-native-live-target \
    -o jsonpath='{.status.observedGeneration}' 2>/dev/null || echo "0")"
  if [[ "${gen}" != "${obs}" ]]; then
    echo "1"
    return
  fi
  echo "0"
}

pod_recreate_detected() {
  local uid_before uid_after
  uid_before="${1:-}"
  uid_after="$(kubectl --context="${CTX}" -n "${NS}" get pods -l app=khr-native-live-target \
    -o jsonpath='{.items[0].metadata.uid}' 2>/dev/null || echo "")"
  if [[ -n "${uid_before}" && -n "${uid_after}" && "${uid_before}" != "${uid_after}" ]]; then
    echo "true"
    return
  fi
  echo "false"
}

now_ms() { date +%s%3N; }

timed_apply() {
  local name="$1" lease="$2" baseline="$3" out="$4"
  local start end
  start="$(now_ms)"
  "${BIN}" -mode=resourcelease-guarded-apply -config="${CFG_APPLY}" -namespace="${NS}" \
    -cluster-context="${CTX}" -lease-input="${lease}" -sandbox-dir="${SANDBOX_DIR}" \
    -baseline-id="${baseline}" -apply-resourcelease=true -i-understand-this-is-sandbox \
    > "${out}"
  end="$(now_ms)"
  echo $((end - start))
}

timed_rollback() {
  local baseline="$1" out="$2"
  local start end
  start="$(now_ms)"
  "${BIN}" -mode=resourcelease-rollback -config="${CFG_APPLY}" -sandbox-dir="${SANDBOX_DIR}" \
    -baseline-id="${baseline}" > "${out}"
  end="$(now_ms)"
  echo $((end - start))
}

khr_runtime_assert_cluster_context
khr_runtime_assert_sandbox_namespace_labels

kubectl --context="${CTX}" apply -f "${WORKLOAD}" > "${RUN_DIR}/workload-apply.txt"
kubectl --context="${CTX}" -n "${NS}" rollout status deploy/khr-native-live-target --timeout=120s \
  >> "${RUN_DIR}/workload-apply.txt" 2>&1 || true

RESTART_BEFORE="$(pod_restart_count)"
POD_UID_BEFORE="$(kubectl --context="${CTX}" -n "${NS}" get pods -l app=khr-native-live-target \
  -o jsonpath='{.items[0].metadata.uid}' 2>/dev/null || echo "")"
GEN_BEFORE="$(kubectl --context="${CTX}" -n "${NS}" get deploy khr-native-live-target \
  -o jsonpath='{.metadata.generation}' 2>/dev/null || echo "0")"

"${BIN}" -mode=lane-discovery -config="${CFG}" -cluster-context="${CTX}" \
  > "${RUN_DIR}/lane-discovery.json"
NATIVE_LANE="$(jq '[.laneCapabilities[] | select(.lane == "native-live")] | .[0].workloadCount // 0' \
  "${RUN_DIR}/lane-discovery.json")"
[[ "${NATIVE_LANE}" -ge 1 ]] || { echo "FAIL: no native-live lane" >&2; exit 1; }

"${BIN}" -mode=resourcefuture-simulate -config="${CFG}" -cluster-context="${CTX}" \
  > "${RUN_DIR}/simulation.json"
LIVE_ELIG="$(jq '[.liveInPlaceEligibility[] | select(.eligible and .lane == "native-live")] | length' \
  "${RUN_DIR}/simulation.json")"
if [[ "${LIVE_ELIG}" -lt 1 ]]; then
  LIVE_ELIG="$(jq '[.discoveredResourcePorts[] | select(.lane == "native-live" and .liveScaleCapabilityObserved)] | length' \
    "${RUN_DIR}/lane-discovery.json")"
fi
[[ "${LIVE_ELIG}" -ge 1 ]] || { echo "FAIL: no live-in-place eligible native-live" >&2; exit 1; }

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

"${BIN}" -mode=resourcelease-dryrun -config="${CFG_APPLY}" -namespace="${NS}" \
  -cluster-context="${CTX}" -lease-input="${LEASE_CPU}" -sandbox-dir="${SANDBOX_DIR}" \
  -baseline-id="dry-${RUN_DIR##*/}" > "${RUN_DIR}/dryrun-cpu.json"

capture_continuity_snapshot "${RUN_DIR}/continuity-before.json"

WINDOW_START="$(now_ms)"
LAT_CPU="$(timed_apply cpu "${LEASE_CPU}" "cpu-${RUN_DIR##*/}" "${RUN_DIR}/apply-cpu.json")"
jq -e '.applied == true and .verification.noRestart == true and .verification.noRollout == true and .verification.noRecreate == true' \
  "${RUN_DIR}/apply-cpu.json" >/dev/null

LAT_RAM_UP="$(timed_apply ram-up "${LEASE_RAM_UP}" "ram-up-${RUN_DIR##*/}" "${RUN_DIR}/apply-ram-up.json")"
jq -e '.applied == true and .verification.state == "pass"' "${RUN_DIR}/apply-ram-up.json" >/dev/null

LAT_RAM_DOWN="$(timed_apply ram-down "${LEASE_RAM_DOWN}" "ram-down-${RUN_DIR##*/}" "${RUN_DIR}/apply-ram-down.json")"
jq -e '.applied == true' "${RUN_DIR}/apply-ram-down.json" >/dev/null

LAT_RB_RAM_DOWN="$(timed_rollback "ram-down-${RUN_DIR##*/}" "${RUN_DIR}/rollback-ram-down.json")"
LAT_RB_RAM_UP="$(timed_rollback "ram-up-${RUN_DIR##*/}" "${RUN_DIR}/rollback-ram-up.json")"
LAT_RB_CPU="$(timed_rollback "cpu-${RUN_DIR##*/}" "${RUN_DIR}/rollback-cpu.json")"
jq -e '.rolledBack == true' "${RUN_DIR}/rollback-cpu.json" >/dev/null
WINDOW_END="$(now_ms)"

RESTART_AFTER="$(pod_restart_count)"
ROLLOUT_DETECTED="false"
if [[ "$(deploy_rollout_detected)" == "1" ]] || [[ "${GEN_BEFORE}" != "$(kubectl --context="${CTX}" -n "${NS}" get deploy khr-native-live-target -o jsonpath='{.metadata.generation}' 2>/dev/null || echo "0")" ]]; then
  ROLLOUT_DETECTED="true"
fi
RECREATE_DETECTED="$(pod_recreate_detected "${POD_UID_BEFORE}")"
RESTART_DELTA=$((RESTART_AFTER - RESTART_BEFORE))
INTERRUPTION_MS=0
if [[ "${RESTART_DELTA}" -gt 0 || "${ROLLOUT_DETECTED}" == "true" || "${RECREATE_DETECTED}" == "true" ]]; then
  INTERRUPTION_MS=$((WINDOW_END - WINDOW_START))
fi
INTERRUPTION_DETECTED="false"
if [[ "${INTERRUPTION_MS}" -gt 0 ]]; then
  INTERRUPTION_DETECTED="true"
fi

capture_continuity_snapshot "${RUN_DIR}/continuity-after.json"
"${CONTINUITY_BIN}" -before="${RUN_DIR}/continuity-before.json" -after="${RUN_DIR}/continuity-after.json" \
  -out="${RUN_DIR}/continuity-proof.json"
cp "${RUN_DIR}/continuity-proof.json" "${RUN_DIR}/shell-continuity-proof.json"
cp "${RUN_DIR}/continuity-proof.json" "${RUN_DIR}/session-continuity-proof.json"

SHELL_CONT="$(jq -r '.shellContinuityPreserved' "${RUN_DIR}/continuity-proof.json")"
APP_CONT="$(jq -r '.appContinuityPreserved' "${RUN_DIR}/continuity-proof.json")"
SESSION_CONT="$(jq -r '.sessionContinuityPreserved' "${RUN_DIR}/continuity-proof.json")"
CONT_STATE="$(jq -r '.continuityEvidence.continuityState' "${RUN_DIR}/continuity-proof.json")"
CONT_EVIDENCE="$(jq '.continuityEvidence' "${RUN_DIR}/continuity-proof.json")"

jq -n \
  --argjson restartBefore "${RESTART_BEFORE}" \
  --argjson restartAfter "${RESTART_AFTER}" \
  --argjson restartDelta "${RESTART_DELTA}" \
  --argjson rolloutCount "$( [[ "${ROLLOUT_DETECTED}" == "true" ]] && echo 1 || echo 0 )" \
  --argjson rolloutDetected "${ROLLOUT_DETECTED}" \
  --argjson recreateDetected "${RECREATE_DETECTED}" \
  --argjson interruptionDetected "${INTERRUPTION_DETECTED}" \
  --argjson interruptionWindowMs "${INTERRUPTION_MS}" \
  --argjson nativeLane "${NATIVE_LANE}" \
  --argjson liveEligible "$( [[ "${LIVE_ELIG}" -ge 1 ]] && echo true || echo false )" \
  --argjson rollbackPass true \
  --argjson latCpu "${LAT_CPU}" \
  --argjson latRamUp "${LAT_RAM_UP}" \
  --argjson latRamDown "${LAT_RAM_DOWN}" \
  --argjson latRbRamDown "${LAT_RB_RAM_DOWN}" \
  --argjson latRbRamUp "${LAT_RB_RAM_UP}" \
  --argjson latRbCpu "${LAT_RB_CPU}" \
  --argjson shellCont "${SHELL_CONT}" \
  --argjson appCont "${APP_CONT}" \
  --argjson sessionCont "${SESSION_CONT}" \
  --arg continState "${CONT_STATE}" \
  --argjson contEvidence "${CONT_EVIDENCE}" \
  '{
    restartCountBefore: $restartBefore,
    restartCountAfter: $restartAfter,
    restartCountDelta: $restartDelta,
    rolloutCount: $rolloutCount,
    rolloutDetected: $rolloutDetected,
    recreateDetected: $recreateDetected,
    interruptionDetected: $interruptionDetected,
    interruptionWindowMs: $interruptionWindowMs,
    applyLatencyMs: { cpu: $latCpu, ramUp: $latRamUp, ramDown: $latRamDown },
    rollbackLatencyMs: { ramDown: $latRbRamDown, ramUp: $latRbRamUp, cpu: $latRbCpu },
    nativeLiveLaneCount: $nativeLane,
    liveInPlaceEligible: $liveEligible,
    rollbackPass: $rollbackPass,
    shellContinuityPreserved: $shellCont,
    appContinuityPreserved: $appCont,
    userSessionContinuityObserved: $sessionCont,
    sessionContinuityPreserved: $sessionCont,
    continuityState: $continState,
    continuityEvidence: $contEvidence
  }' > "${RUN_DIR}/run-metrics.json"
