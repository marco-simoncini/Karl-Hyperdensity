#!/usr/bin/env bash
# KHR-AC: read-only Technical Preview operator bundle index (no mutating apply).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)"
OUT_DIR="${ROOT}/docs/evidence/khr-tp-operator-bundle/${RUN_ID}"
mkdir -p "${OUT_DIR}"

FAIL=0
CHECKS=()

record_check() {
  local name="$1"
  local status="$2"
  CHECKS+=("${name}:${status}")
  if [[ "${status}" != "PASS" ]]; then
    FAIL=1
  fi
}

echo "[khr_tp_operator_bundle] runId=${RUN_ID}"

if [[ -x "${ROOT}/scripts/guard_khr_docs_scope.sh" ]]; then
  if "${ROOT}/scripts/guard_khr_docs_scope.sh" >/dev/null 2>&1; then
    record_check "guard_khr_docs_scope" "PASS"
  else
    record_check "guard_khr_docs_scope" "FAIL"
  fi
else
  record_check "guard_khr_docs_scope" "SKIP"
fi

if [[ -x "${ROOT}/scripts/khr_tp_package_check.sh" ]]; then
  if "${ROOT}/scripts/khr_tp_package_check.sh" >/dev/null 2>&1; then
    record_check "khr_tp_package_check" "PASS"
  else
    record_check "khr_tp_package_check" "FAIL"
  fi
else
  record_check "khr_tp_package_check" "SKIP"
fi

require_evidence() {
  local path="$1"
  local name="$2"
  if [[ -f "${path}" ]]; then
    record_check "${name}" "PASS"
  else
    record_check "${name}" "FAIL"
  fi
}

require_evidence "docs/khr/TECHNICAL_PREVIEW_READINESS.md" "tp_readiness_doc"
require_evidence "docs/khr/TECHNICAL_PREVIEW_PACKAGE.md" "tp_package_doc"
require_evidence "docs/khr/TECHNICAL_PREVIEW_OPERATOR_RUNBOOK.md" "tp_operator_runbook"
require_evidence "docs/evidence/khr-native-live-lane/certification-summary.json" "native_live_certification"
require_evidence "docs/evidence/khr-certification-registry/summary.json" "certification_registry_summary"
require_evidence "docs/evidence/khr-provenance/summary.json" "provenance_summary"

# Bundle index (read-only inventory).
INDEX_JSON="${OUT_DIR}/bundle-index.json"
{
  echo '{'
  echo '  "sprint": "KHR-AC",'
  echo '  "runId": "'"${RUN_ID}"'",'
  echo '  "readOnly": true,'
  echo '  "productionReady": false,'
  echo '  "autonomousOrchestration": false,'
  echo '  "sandboxManualOnly": true,'
  echo '  "evidenceAnchors": ['
  first=1
  for f in \
    docs/evidence/khr-native-live-lane/certification-summary.json \
    docs/evidence/khr-certification-registry/summary.json \
    docs/evidence/khr-certification-registry/registry.json \
    docs/evidence/khr-provenance/summary.json \
    docs/evidence/khr-control-graph/control-graph.json \
    docs/evidence/khr-action-approval/summary.json; do
    if [[ -f "${ROOT}/${f}" ]]; then
      [[ "${first}" -eq 1 ]] || echo ','
      first=0
      echo -n '    "'"${f}"'"'
    fi
  done
  echo ''
  echo '  ]'
  echo '}'
} > "${INDEX_JSON}"

# Run summary.
RUN_SUMMARY="${OUT_DIR}/run-summary.json"
{
  echo '{'
  echo '  "runId": "'"${RUN_ID}"'",'
  echo '  "status": "'"$( [[ "${FAIL}" -eq 0 ]] && echo PASS || echo FAIL )"'",'
  echo '  "readOnly": true,'
  echo '  "checks": ['
  cfirst=1
  for c in "${CHECKS[@]}"; do
    name="${c%%:*}"
    status="${c#*:}"
    [[ "${cfirst}" -eq 1 ]] || echo ','
    cfirst=0
    echo -n '    {"name":"'"${name}"'","status":"'"${status}"'"}'
  done
  echo ''
  echo '  ]'
  echo '}'
} > "${RUN_SUMMARY}"

# Blocker summary (static TP boundaries + check failures).
BLOCKER_JSON="${OUT_DIR}/blocker-summary.json"
if [[ "${FAIL}" -ne 0 ]]; then
  jq -n '{
    productionReady: false,
    autonomousOrchestration: false,
    isoRuntimeDisabledByDefault: true,
    blockers: [
      {id:"AC-01",severity:"P0",message:"NOT production ready"},
      {id:"AC-02",severity:"P0",message:"No autonomous orchestration"},
      {id:"AC-03",severity:"P0",message:"Sandbox/manual evidence only"},
      {id:"AC-04",severity:"P1",message:"Bundle checks failed — review run-summary.json"}
    ]
  }' > "${BLOCKER_JSON}"
else
  jq -n '{
    productionReady: false,
    autonomousOrchestration: false,
    isoRuntimeDisabledByDefault: true,
    blockers: [
      {id:"AC-01",severity:"P0",message:"NOT production ready"},
      {id:"AC-02",severity:"P0",message:"No autonomous orchestration"},
      {id:"AC-03",severity:"P0",message:"Sandbox/manual evidence only"}
    ]
  }' > "${BLOCKER_JSON}"
fi

# Next-action hints.
NEXT_JSON="${OUT_DIR}/next-actions.json"
cat > "${NEXT_JSON}" <<EOF
{
  "readOnly": true,
  "hints": [
    "Review docs/khr/TECHNICAL_PREVIEW_OPERATOR_RUNBOOK.md",
    "Run sandbox evidence only on khr-runtime-sandbox with KHR_RUNTIME_SANDBOX_LIVE=1 if live proof required",
    "Enable Dashboard GET /api/hyperdensity/tp-readiness with HYPERDENSITY_KHR_TP_READINESS_ENABLED=true",
    "Export Inventory posture via scripts/khr_tp_observation_export.sh (stub)",
    "Do not enable karl-host-runtime systemd on ISO for TP",
    "Beta: wire Inventory periodic export and Dashboard approval UX (non-autonomous)"
  ]
}
EOF

# Latest summary symlink-style JSON at bundle root.
cp "${RUN_SUMMARY}" "${ROOT}/docs/evidence/khr-tp-operator-bundle/summary.json"

if [[ "${FAIL}" -ne 0 ]]; then
  echo "[khr_tp_operator_bundle] FAIL ${OUT_DIR}" >&2
  exit 1
fi

echo "[khr_tp_operator_bundle] PASS ${OUT_DIR}"
