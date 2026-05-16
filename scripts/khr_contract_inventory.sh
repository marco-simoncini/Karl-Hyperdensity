#!/usr/bin/env bash
# KHR-AD: read-only KHR contract inventory (docs/schemas/tests/evidence presence).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

FAIL=0
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)"
OUT_DIR="${ROOT}/docs/evidence/khr-contract-inventory/${RUN_ID}"
mkdir -p "${OUT_DIR}"

check_file() {
  local path="$1"
  [[ -f "${ROOT}/${path}" ]]
}

check_glob() {
  local pattern="$1"
  if [[ "${pattern}" == ../* ]]; then
    compgen -G "${pattern}" >/dev/null 2>&1
  else
    compgen -G "${ROOT}/${pattern}" >/dev/null 2>&1
  fi
}

resolve_path() {
  local path="$1"
  if [[ "${path}" == Karl-Dashboard/* ]]; then
    echo "../Karl-Dashboard/${path#Karl-Dashboard/}"
  else
    echo "${path}"
  fi
}

record_contract() {
  local id="$1" version="$2" doc="$3" schema="$4" tests="$5" evidence="$6" stability="$7"
  local doc_ok schema_ok tests_ok evidence_ok status
  doc_ok=false; schema_ok=false; tests_ok=false; evidence_ok=false
  [[ -z "${doc}" || "${doc}" == "-" ]] && doc_ok=true
  if [[ -n "${doc}" && "${doc}" != "-" ]]; then
    local resolved
    resolved="$(resolve_path "${doc}")"
    if [[ "${resolved}" == ../* ]]; then
      [[ -f "${resolved}" ]] && doc_ok=true
    else
      check_file "${resolved}" && doc_ok=true
    fi
  fi
  [[ -z "${schema}" || "${schema}" == "-" ]] && schema_ok=true
  if [[ -n "${schema}" && "${schema}" != "-" ]]; then
    if [[ "${schema}" == Karl-Dashboard/* ]]; then
      check_file "${schema/#Karl-Dashboard\//../Karl-Dashboard/}" && schema_ok=true
    else
      { check_file "${schema}" || check_glob "${schema}"; } && schema_ok=true
    fi
  fi
  [[ -z "${tests}" || "${tests}" == "-" ]] && tests_ok=true
  if [[ -n "${tests}" && "${tests}" != "-" ]]; then
    check_glob "$(resolve_path "${tests}")" && tests_ok=true
  fi
  [[ -z "${evidence}" || "${evidence}" == "-" ]] && evidence_ok=true
  [[ -n "${evidence}" && "${evidence}" != "-" ]] && { check_file "${evidence}" || check_glob "${evidence}"; } && evidence_ok=true
  status="PASS"
  if ! "${doc_ok}" || ! "${schema_ok}" || ! "${tests_ok}" || ! "${evidence_ok}"; then
    status="GAP"
    FAIL=1
  fi
  CONTRACT_ROWS+=("$(jq -nc \
    --arg id "${id}" --arg version "${version}" --arg stability "${stability}" \
    --arg doc "${doc}" --arg schema "${schema}" --arg tests "${tests}" --arg evidence "${evidence}" \
    --arg status "${status}" \
    --argjson docOk "${doc_ok}" --argjson schemaOk "${schema_ok}" \
    --argjson testsOk "${tests_ok}" --argjson evidenceOk "${evidence_ok}" \
    '{id:$id,version:$version,stability:$stability,documentation:$doc,schemaOrCrd:$schema,tests:$tests,evidence:$evidence,status:$status,presence:{documentation:$docOk,schemaOrCrd:$schemaOk,tests:$testsOk,evidence:$evidenceOk}}')")
}

CONTRACT_ROWS=()

echo "[khr_contract_inventory] scanning contracts..."

# Dashboard-facing (Karl-Dashboard sibling repo when present).
if [[ -d "${ROOT}/../Karl-Dashboard" ]]; then
  record_contract "khr-projection" "khr-projection-v1alpha1-readonly-y" "Karl-Dashboard/docs/hyperdensity/KHR_PROJECTION_V1.md" "-" "Karl-Dashboard/Karl-Dashboard-dashboard/kubernetes-console/pkg/server/hyperdensity_*khr*_test.go" "-" "freeze-tp"
  record_contract "tp-readiness-summary" "khr-tp-readiness-summary-v1alpha1" "Karl-Dashboard/docs/khr/TECHNICAL_PREVIEW_DASHBOARD_GUIDE.md" "-" "Karl-Dashboard/Karl-Dashboard-dashboard/kubernetes-console/pkg/server/hyperdensity_tp_readiness_v1_test.go" "-" "freeze-tp"
else
  record_contract "khr-projection" "khr-projection-v1alpha1-readonly-y" "-" "-" "-" "-" "freeze-tp"
  record_contract "tp-readiness-summary" "khr-tp-readiness-summary-v1alpha1" "-" "-" "-" "-" "freeze-tp"
fi

record_contract "host" "runtime.karl.io/v1alpha1" "docs/khr/HOST_CONTRACT.md" "api/crds/runtime.karl.io/host.yaml" "pkg/khr/host/*_test.go" "docs/evidence/khr-host-heartbeat/*" "freeze-tp"
record_contract "shell-cell" "runtime.karl.io/v1alpha1" "docs/khr/SHELL_CELL_CONTRACT.md" "api/crds/runtime.karl.io/shell.yaml" "pkg/khr/crdv1alpha1/*_test.go" "docs/evidence/khr-runtime-sandbox/*" "freeze-tp"
record_contract "resourceport" "runtime.karl.io/v1alpha1" "docs/khr/RESOURCEPORT_CONTRACT.md" "api/crds/runtime.karl.io/resourceport.yaml" "pkg/khr/resourceport/*_test.go" "docs/evidence/khr-resourceport-loop/*" "freeze-tp"
record_contract "resourcelease" "hyperdensity.karl.io/v1alpha1" "docs/khr/RESOURCELEASE_LIFECYCLE.md" "api/crds/hyperdensity.karl.io/resourcelease.yaml" "pkg/khr/resourcelease/*_test.go" "docs/evidence/khr-resourcelease-dryrun/*" "experimental"
record_contract "resourcefuture" "hyperdensity.karl.io/v1alpha1" "docs/khr/RESOURCEFUTURE_SIMULATION.md" "api/crds/hyperdensity.karl.io/resourcefuture.yaml" "pkg/khr/resourcefuture/*_test.go" "docs/evidence/khr-resourcefuture/*" "freeze-tp"
record_contract "shelllease-gatewayroute" "runtime.karl.io/v1alpha1 / gateway.karl.io/v1alpha1" "docs/khr/SHELLLEASE_GATEWAYROUTE_CONTRACT.md" "api/crds/gateway.karl.io/gatewayroute.yaml" "pkg/khr/crdv1alpha1/*_test.go" "-" "freeze-tp"
record_contract "certification-registry" "khr-cert-registry-v1" "docs/khr/CERTIFICATION_REGISTRY_AND_POLICY_GATES.md" "-" "pkg/khr/certregistry/*_test.go" "docs/evidence/khr-certification-registry/summary.json" "experimental"
record_contract "policy-gates" "khr-policy-gates-v1" "docs/khr/CERTIFICATION_REGISTRY_AND_POLICY_GATES.md" "-" "pkg/khr/policygates/*_test.go" "docs/evidence/khr-certification-registry/*" "experimental"
record_contract "action-approval" "khr-action-approval-v1" "docs/khr/OPERATOR_ACTION_APPROVAL_WORKFLOW.md" "-" "pkg/khr/actionapproval/*_test.go" "docs/evidence/khr-action-approval/summary.json" "experimental"
record_contract "control-graph" "khr-control-graph-v1" "docs/khr/KHR_CONTROL_GRAPH.md" "-" "pkg/khr/controlgraph/*_test.go" "docs/evidence/khr-control-graph/summary.json" "experimental"
record_contract "provenance" "khr-provenance-v1" "docs/khr/TRUST_AND_PROVENANCE_MODEL.md" "-" "pkg/khr/provenance/*_test.go" "docs/evidence/khr-provenance/summary.json" "experimental"

# Hyperdensity-local required docs for freeze plan.
REQUIRED_DOCS=(
  docs/khr/BETA_READINESS_GAP_ANALYSIS.md
  docs/khr/KHR_CONTRACT_FREEZE_PLAN.md
  docs/khr/TECHNICAL_PREVIEW_READINESS.md
  docs/khr/TECHNICAL_PREVIEW_PACKAGE.md
  docs/khr/TECHNICAL_PREVIEW_OPERATOR_RUNBOOK.md
)
for doc in "${REQUIRED_DOCS[@]}"; do
  if ! check_file "${doc}"; then
    echo "[khr_contract_inventory] FAIL: missing ${doc}" >&2
    FAIL=1
  fi
done

# Build JSON report.
rows_json="["
first=1
for row in "${CONTRACT_ROWS[@]}"; do
  [[ "${first}" -eq 1 ]] || rows_json+=","
  first=0
  rows_json+="${row}"
done
rows_json+="]"

report="${OUT_DIR}/contract-inventory.json"
jq -n \
  --arg runId "${RUN_ID}" \
  --arg sprint "KHR-AD" \
  --argjson contracts "${rows_json}" \
  --argjson pass $([[ "${FAIL}" -eq 0 ]] && echo true || echo false) \
  '{
    sprint: $sprint,
    runId: $runId,
    readOnly: true,
    productionReady: false,
    autonomousOrchestration: false,
    status: (if $pass then "PASS" else "GAP" end),
    contracts: $contracts
  }' > "${report}"

cp "${report}" "${ROOT}/docs/evidence/khr-contract-inventory/summary.json"

echo "[khr_contract_inventory] wrote ${report}"

if [[ "${FAIL}" -ne 0 ]]; then
  echo "[khr_contract_inventory] FAIL — see contract-inventory.json for GAP rows" >&2
  exit 1
fi

echo "[khr_contract_inventory] PASS"
