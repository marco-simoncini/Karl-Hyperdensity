#!/usr/bin/env bash
# KHR-T: repeatable native-live certification (multi-run, baseline compare, regression guard).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
EVIDENCE="${ROOT}/docs/evidence/khr-native-live-lane"
BASELINE="${ROOT}/examples/khr/native-live/baseline-certification.json"
CERT_ID="${KHR_NATIVE_LIVE_CERT_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
CERT_DIR="${EVIDENCE}/certification/${CERT_ID}"
RUNS="${KHR_NATIVE_LIVE_CERT_RUNS:-2}"
RUN_SCRIPT="${ROOT}/scripts/khr_native_live_lane_run.sh"
CERTIFY_BIN="${ROOT}/bin/khr-native-live-certify"

khr_runtime_assert_cluster_context
mkdir -p "${CERT_DIR}"
chmod +x "${RUN_SCRIPT}"

METRIC_FILES=()
for i in $(seq 1 "${RUNS}"); do
  run_dir="${CERT_DIR}/run-${i}"
  khr_runtime_log "certification run ${i}/${RUNS}"
  "${RUN_SCRIPT}" "${run_dir}"
  METRIC_FILES+=("${run_dir}/run-metrics.json")
done

(cd "${ROOT}" && go build -o "${CERTIFY_BIN}" ./cmd/khr-native-live-certify)

CERTIFY_ARGS=(-sprint=KHR-T -baseline="${BASELINE}" -out="${CERT_DIR}/certification-summary.json")
if [[ "${KHR_NATIVE_LIVE_REQUIRE_BASELINE:-1}" == "1" ]]; then
  CERTIFY_ARGS+=(-require-baseline-match)
fi
"${CERTIFY_BIN}" "${CERTIFY_ARGS[@]}" "${METRIC_FILES[@]}"

cp "${CERT_DIR}/certification-summary.json" "${EVIDENCE}/certification-summary.json"

jq -n \
  --arg certId "${CERT_ID}" \
  --arg cluster "${CTX}" \
  --argjson runs "${RUNS}" \
  --slurpfile summary "${CERT_DIR}/certification-summary.json" \
  '{
    sprint: "KHR-T",
    certificationId: $summary[0].certificationId,
    certificationRunId: $certId,
    cluster: $cluster,
    runCount: $runs,
    status: $summary[0].status,
    regressionDetected: $summary[0].regressionDetected,
    baselineMatch: $summary[0].baselineMatch,
    continuityScore: $summary[0].scores.continuityScore,
    liveScaleConfidence: $summary[0].scores.liveScaleConfidence,
    invariants: $summary[0].invariants
  }' > "${CERT_DIR}/summary.json"
cp "${CERT_DIR}/summary.json" "${EVIDENCE}/certification-run-summary.json"

echo "[khr_native_live_certify] PASS ${CERT_DIR} runs=${RUNS} status=$(jq -r .status "${CERT_DIR}/certification-summary.json")"
