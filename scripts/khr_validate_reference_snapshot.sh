#!/usr/bin/env bash
# KHR-BU: offline reference snapshot v1 validation (docs + aggregator + committed evidence).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

for doc in \
  docs/khr/KHR_VALIDATION_MODES.md \
  docs/khr/KHR_SNAPSHOT_V1_FREEZE_POLICY.md \
  docs/khr/KHR_BETA_READINESS_PLAN.md \
  docs/khr/KHR_TP_REFERENCE_SNAPSHOT_V1.md; do
  [[ -f "${ROOT}/${doc}" ]] || {
    echo "[khr_validate_reference_snapshot] FAIL: missing ${doc}" >&2
    exit 1
  }
done

KHR_TP_REFERENCE_SNAPSHOT_RUN_ID="${KHR_TP_REFERENCE_SNAPSHOT_RUN_ID:-committed-khr-bt-v1}" \
  "${ROOT}/scripts/validate_khr_tp_reference_snapshot_v1.sh"

"${ROOT}/scripts/khr_validate_committed_evidence.sh"

echo "[khr_validate_reference_snapshot] PASS"
