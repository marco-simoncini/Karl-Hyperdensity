#!/usr/bin/env bash
# KHR-S/T: single-run native-live evidence (delegates to certification run script).
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
EVIDENCE="${ROOT}/docs/evidence/khr-native-live-lane"
RUN_ID="${KHR_NATIVE_LIVE_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
exec "${ROOT}/scripts/khr_native_live_lane_run.sh" "${RUN_DIR}"
