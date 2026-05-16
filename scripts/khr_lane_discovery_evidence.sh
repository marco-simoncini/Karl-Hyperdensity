#!/usr/bin/env bash
# KHR-Q: read-only multi-lane discovery evidence on karl-metal-01@ovh.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
EVIDENCE="${ROOT}/docs/evidence/khr-lane-discovery"
RUN_ID="${KHR_LANE_DISCOVERY_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
BIN="${ROOT}/bin/karl-host-runtime"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-lane-discovery.yaml"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${ROOT}/bin"
(cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime)

khr_runtime_log "lane-discovery (read-only)"
"${BIN}" -mode=lane-discovery -config="${CFG}" -cluster-context="${CTX}" \
  > "${RUN_DIR}/lane-discovery.json"

jq -e '.mode == "lane-discovery" and .safety.readOnly == true and .safety.noApply == true' \
  "${RUN_DIR}/lane-discovery.json" >/dev/null
jq -e '.safety.noRestart == true and .safety.noRollout == true and .safety.noRecreate == true' \
  "${RUN_DIR}/lane-discovery.json" >/dev/null
jq -e '(.discoveredHosts | length) >= 1' "${RUN_DIR}/lane-discovery.json" >/dev/null
jq -e '(.discoveredCells | length) >= 1' "${RUN_DIR}/lane-discovery.json" >/dev/null
jq -e '(.laneCapabilities | length) >= 1' "${RUN_DIR}/lane-discovery.json" >/dev/null

khr_runtime_log "attest no kubectl mutation in this run"
: > "${RUN_DIR}/mutation-check.txt"
echo "no-apply: lane-discovery mode only emits JSON" >> "${RUN_DIR}/mutation-check.txt"
echo "no-restart: no rollout restart commands executed" >> "${RUN_DIR}/mutation-check.txt"

# Snapshot cluster VM count for correlation (read-only get).
kubectl --context="${CTX}" get virtualmachines.kubevirt.io -A -o json \
  > "${RUN_DIR}/cluster-vms-snapshot.json"
VM_COUNT="$(jq '.items | length' "${RUN_DIR}/cluster-vms-snapshot.json")"
DISC_COUNT="$(jq '.discoveredCells | length' "${RUN_DIR}/lane-discovery.json")"

jq -n \
  --arg sprint "KHR-Q" \
  --arg runId "${RUN_ID}" \
  --arg cluster "${CTX}" \
  --argjson vmCount "${VM_COUNT}" \
  --argjson discoveredCells "${DISC_COUNT}" \
  '{sprint:$sprint,runId:$runId,cluster:$cluster,vmCount:$vmCount,discoveredCells:$discoveredCells,readOnly:true,noApply:true,noRestart:true}' \
  > "${RUN_DIR}/summary.json"
cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"

echo "[khr_lane_discovery_evidence] PASS ${RUN_DIR} vms=${VM_COUNT} cells=${DISC_COUNT}"
