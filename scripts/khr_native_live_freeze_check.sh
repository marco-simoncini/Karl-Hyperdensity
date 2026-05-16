#!/usr/bin/env bash
# KHR-AF: Native-live + ResourcePort TP freeze candidate checks (evidence invariants; no runtime mutation).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

FAIL=0
CERT="${ROOT}/docs/evidence/khr-native-live-lane/certification-summary.json"

require_file() {
  local path="$1"
  local label="$2"
  if [[ -f "${path}" ]]; then
    echo "[khr_native_live_freeze_check] OK: ${label}"
  else
    echo "[khr_native_live_freeze_check] FAIL: missing ${label}: ${path}" >&2
    FAIL=1
  fi
}

require_jq_eq() {
  local file="$1"
  local filter="$2"
  local want="$3"
  local label="$4"
  local got
  got="$(jq -r "${filter}" "${file}" 2>/dev/null || echo "")"
  if [[ "${got}" == "${want}" ]]; then
    echo "[khr_native_live_freeze_check] OK: ${label}"
  else
    echo "[khr_native_live_freeze_check] FAIL: ${label} (got=${got} want=${want})" >&2
    FAIL=1
  fi
}

echo "[khr_native_live_freeze_check] Native-live / ResourcePort TP freeze audit..."

# Docs.
require_file "${ROOT}/docs/khr/NATIVE_LIVE_TP_FREEZE_CANDIDATE.md" "native-live freeze doc"
require_file "${ROOT}/docs/khr/RESOURCEPORT_TP_FREEZE_CANDIDATE.md" "resourceport freeze doc"
require_file "${ROOT}/docs/khr/RESOURCEPORT_CONTRACT.md" "resourceport contract"
require_file "${ROOT}/docs/khr/NATIVE_LIVE_CERTIFICATION.md" "native-live certification doc"
require_file "${ROOT}/api/crds/runtime.karl.io/resourceport.yaml" "ResourcePort CRD"

# Native-live evidence anchor.
require_file "${CERT}" "native-live certification summary"

if [[ -f "${CERT}" ]]; then
  require_jq_eq "${CERT}" '.lane' 'native-live' 'lane native-live'
  require_jq_eq "${CERT}" '.scores.continuityScore' '1' 'continuityScore == 1'
  require_jq_eq "${CERT}" '.scores.liveScaleConfidence' 'high' 'liveScaleConfidence high'
  require_jq_eq "${CERT}" '.invariants.noRestart' 'true' 'noRestart'
  require_jq_eq "${CERT}" '.invariants.noRollout' 'true' 'noRollout'
  require_jq_eq "${CERT}" '.invariants.noRecreate' 'true' 'noRecreate'
  require_jq_eq "${CERT}" '.invariants.interruptionWindowMs' '0' 'invariants interruptionWindowMs == 0'
  require_jq_eq "${CERT}" '.metrics.interruptionWindowMs' '0' 'metrics interruptionWindowMs == 0'
  require_jq_eq "${CERT}" '.invariants.interruptionDetected' 'false' 'interruptionDetected false'
  require_jq_eq "${CERT}" '.metrics.restartCountDelta' '0' 'restartCountDelta == 0'
  require_jq_eq "${CERT}" '.metrics.rolloutCount' '0' 'rolloutCount == 0'
  require_jq_eq "${CERT}" '.metrics.recreateDetected' 'false' 'recreateDetected false'
  require_jq_eq "${CERT}" '.readOnly' 'true' 'certification readOnly'
  require_jq_eq "${CERT}" '.noAutonomousOrchestration' 'true' 'noAutonomousOrchestration'
fi

# Certification run dirs exist.
if ! compgen -G "${ROOT}/docs/evidence/khr-native-live-lane/certification/*" >/dev/null; then
  echo "[khr_native_live_freeze_check] FAIL: no certification run directories" >&2
  FAIL=1
else
  echo "[khr_native_live_freeze_check] OK: native-live certification runs"
fi

# ResourceFuture eligibility on native-live run metrics.
live_eligible=false
metrics_found=false
while IFS= read -r -d '' mf; do
  metrics_found=true
  if jq -e '.liveInPlaceEligible == true' "${mf}" >/dev/null 2>&1; then
    live_eligible=true
    break
  fi
done < <(find "${ROOT}/docs/evidence/khr-native-live-lane/certification" -name run-metrics.json -print0 2>/dev/null || true)
if ! "${metrics_found}"; then
  echo "[khr_native_live_freeze_check] FAIL: missing run-metrics.json under certification" >&2
  FAIL=1
elif ! "${live_eligible}"; then
  echo "[khr_native_live_freeze_check] FAIL: no run-metrics.json with liveInPlaceEligible true" >&2
  FAIL=1
else
  echo "[khr_native_live_freeze_check] OK: ResourceFuture liveInPlaceEligible on native-live run"
fi

# ResourcePort native-live CR preview in certification evidence.
if compgen -G "${ROOT}/docs/evidence/khr-native-live-lane/certification/*/run-*/cr-preview/resourceport-*native-live*" >/dev/null; then
  echo "[khr_native_live_freeze_check] OK: native-live ResourcePort CR preview evidence"
else
  echo "[khr_native_live_freeze_check] FAIL: missing native-live ResourcePort cr-preview" >&2
  FAIL=1
fi

# Provenance + registry + policy gate evidence (read-only program).
require_file "${ROOT}/docs/evidence/khr-provenance/summary.json" "provenance summary"
require_file "${ROOT}/docs/evidence/khr-certification-registry/summary.json" "certification registry summary"
require_jq_eq "${ROOT}/docs/evidence/khr-provenance/summary.json" '.noAutonomousOrchestration' 'true' 'provenance noAutonomousOrchestration'
require_jq_eq "${ROOT}/docs/evidence/khr-certification-registry/summary.json" '.readOnly' 'true' 'registry readOnly'

# ResourceFuture program evidence (simulation-only).
require_file "${ROOT}/docs/evidence/khr-resourcefuture/summary.json" "resourcefuture summary"
require_jq_eq "${ROOT}/docs/evidence/khr-resourcefuture/summary.json" '.simulationOnly' 'true' 'resourcefuture simulationOnly'
require_jq_eq "${ROOT}/docs/evidence/khr-resourcefuture/summary.json" '.noApply' 'true' 'resourcefuture noApply'

# Dashboard projection (sibling repo).
PROJ="${ROOT}/../Karl-Dashboard/Karl-Dashboard-dashboard/kubernetes-console/pkg/server/hyperdensity_parent_fabric_khr_projection_v1.go"
if [[ -f "${PROJ}" ]]; then
  if rg -q 'NativeLiveEligible|native-live|LiveScaleCapabilities' "${PROJ}"; then
    echo "[khr_native_live_freeze_check] OK: Dashboard native-live projection fields"
  else
    echo "[khr_native_live_freeze_check] FAIL: missing native-live projection in Dashboard" >&2
    FAIL=1
  fi
else
  echo "[khr_native_live_freeze_check] WARN: Dashboard sibling not found; skip projection source"
fi

if [[ -x "${ROOT}/scripts/guard_khr_docs_scope.sh" ]]; then
  "${ROOT}/scripts/guard_khr_docs_scope.sh"
fi

if [[ "${FAIL}" -ne 0 ]]; then
  echo "[khr_native_live_freeze_check] See docs/khr/NATIVE_LIVE_TP_FREEZE_CANDIDATE.md" >&2
  exit 1
fi

echo "[khr_native_live_freeze_check] PASS — Native-live / ResourcePort TP freeze candidate"
