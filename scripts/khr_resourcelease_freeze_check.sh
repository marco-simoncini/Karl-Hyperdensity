#!/usr/bin/env bash
# KHR-AE: ResourceLease TP freeze candidate checks (docs/evidence/projection; no runtime mutation).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

FAIL=0

require_file() {
  local path="$1"
  local label="$2"
  if [[ -f "${path}" ]]; then
    echo "[khr_resourcelease_freeze_check] OK: ${label}"
  else
    echo "[khr_resourcelease_freeze_check] FAIL: missing ${label}: ${path}" >&2
    FAIL=1
  fi
}

require_json_field() {
  local file="$1"
  local field="$2"
  local expected="$3"
  local label="$4"
  if [[ ! -f "${file}" ]]; then
    echo "[khr_resourcelease_freeze_check] FAIL: missing ${label}" >&2
    FAIL=1
    return
  fi
  local actual
  actual="$(jq -r "${field}" "${file}" 2>/dev/null || echo "")"
  if [[ "${actual}" == "${expected}" ]]; then
    echo "[khr_resourcelease_freeze_check] OK: ${label}"
  else
    echo "[khr_resourcelease_freeze_check] FAIL: ${label} (${field}=${actual} want ${expected})" >&2
    FAIL=1
  fi
}

echo "[khr_resourcelease_freeze_check] ResourceLease TP freeze candidate audit..."

# Contract artifacts.
require_file "api/crds/hyperdensity.karl.io/resourcelease.yaml" "ResourceLease CRD"
require_file "docs/contracts/khr/resourcelease.schema.json" "ResourceLease JSON schema"
require_file "docs/contracts/khr/resourcelease.schema.manifest.json" "schema manifest"
require_file "docs/khr/RESOURCELEASE_LIFECYCLE.md" "lifecycle doc"
require_file "docs/khr/RESOURCELEASE_TP_FREEZE_CANDIDATE.md" "freeze candidate doc"
require_file "docs/khr/RESOURCELEASE_DRYRUN_AGAINST_RESOURCEPORT.md" "dry-run doc"
require_file "docs/khr/RESOURCELEASE_GUARDED_APPLY_SANDBOX.md" "guarded apply doc"

# Examples (at least one contract example + sandbox lease fixture).
if ! compgen -G "${ROOT}/docs/contracts/khr/examples/resourcelease-*.json" >/dev/null; then
  echo "[khr_resourcelease_freeze_check] FAIL: no docs/contracts/khr/examples/resourcelease-*.json" >&2
  FAIL=1
else
  echo "[khr_resourcelease_freeze_check] OK: contract examples"
fi
if ! compgen -G "${ROOT}/examples/khr/runtime-sandbox/resourcelease-*.json" >/dev/null; then
  echo "[khr_resourcelease_freeze_check] FAIL: no examples/khr/runtime-sandbox/resourcelease-*.json" >&2
  FAIL=1
else
  echo "[khr_resourcelease_freeze_check] OK: sandbox lease examples"
fi

# Evidence anchors.
require_file "docs/evidence/khr-resourcelease-dryrun/summary.json" "dry-run evidence summary"
require_file "docs/evidence/khr-resourcelease-guarded-apply/summary.json" "guarded apply evidence summary"
require_file "docs/evidence/khr-provenance/summary.json" "provenance evidence summary"

if compgen -G "${ROOT}/docs/evidence/khr-resourcelease-guarded-apply/*/rollback.json" >/dev/null; then
  echo "[khr_resourcelease_freeze_check] OK: rollback evidence"
else
  echo "[khr_resourcelease_freeze_check] FAIL: no rollback.json under guarded-apply evidence" >&2
  FAIL=1
fi

require_json_field "docs/evidence/khr-resourcelease-dryrun/summary.json" ".noResourceLeaseApply" "true" "dry-run no cluster lease apply"
require_json_field "docs/evidence/khr-resourcelease-dryrun/summary.json" ".noProductionMutation" "true" "dry-run no production mutation"
require_json_field "docs/evidence/khr-resourcelease-guarded-apply/summary.json" ".noProductionMutation" "true" "guarded apply no production mutation"
require_json_field "docs/evidence/khr-provenance/summary.json" ".noAutonomousOrchestration" "true" "provenance no autonomous orchestration"
require_json_field "docs/evidence/khr-provenance/summary.json" ".readOnly" "true" "provenance readOnly"

# Dashboard projection version >= readonly-m (current readonly-y).
PROJ_FILE="${ROOT}/../Karl-Dashboard/Karl-Dashboard-dashboard/kubernetes-console/pkg/server/hyperdensity_parent_fabric_khr_projection_v1.go"
if [[ -f "${PROJ_FILE}" ]]; then
  if rg -q 'khr-projection-v1alpha1-readonly-y' "${PROJ_FILE}"; then
    echo "[khr_resourcelease_freeze_check] OK: projection contract readonly-y (>= readonly-m)"
  else
    echo "[khr_resourcelease_freeze_check] FAIL: expected khr-projection-v1alpha1-readonly-y in Dashboard projection" >&2
    FAIL=1
  fi
  if rg -q 'HyperdensityResourceLeaseSummaryV1' "${PROJ_FILE}"; then
    echo "[khr_resourcelease_freeze_check] OK: ResourceLeaseSummary projection type"
  else
    echo "[khr_resourcelease_freeze_check] FAIL: missing ResourceLeaseSummary projection" >&2
    FAIL=1
  fi
else
  echo "[khr_resourcelease_freeze_check] WARN: Dashboard repo not sibling; skip projection source check"
fi

# Docs wording (freeze candidate doc only).
TP_GLOB="${ROOT}/docs/khr/RESOURCELEASE_TP_FREEZE_CANDIDATE.md"
while IFS= read -r match; do
  [[ -z "${match}" ]] && continue
  line="${match#*:}"
  line="${line#*:}"
  line="${line#*:}"
  shopt -s nocasematch
  if [[ "${line}" == *"NOT production ready"* ]] || [[ "${line}" == *"not production"* ]] \
    || [[ "${line}" == *"No autonomous"* ]] || [[ "${line}" == *"forbidden"* ]] \
    || [[ "${line}" == *"Forbidden"* ]] || [[ "${line}" == *"Experimental"* ]] \
    || [[ "${line}" == *"sandbox only"* ]] || [[ "${line}" == *"Imply production"* ]]; then
    shopt -u nocasematch
    continue
  fi
  shopt -u nocasematch
  if rg -qi 'production[- ]ready|generally available|autonomous orchestration.{0,20}(enabled|production)' <<<"${line}"; then
    echo "[khr_resourcelease_freeze_check] FAIL (wording): ${match}" >&2
    FAIL=1
  fi
done < <(rg -n -i 'production[- ]ready|generally available|autonomous orchestration' "${TP_GLOB}" 2>/dev/null || true)

if [[ -x "${ROOT}/scripts/guard_khr_docs_scope.sh" ]]; then
  "${ROOT}/scripts/guard_khr_docs_scope.sh"
fi

if [[ "${FAIL}" -ne 0 ]]; then
  echo "[khr_resourcelease_freeze_check] See docs/khr/RESOURCELEASE_TP_FREEZE_CANDIDATE.md" >&2
  exit 1
fi

echo "[khr_resourcelease_freeze_check] PASS — ResourceLease TP freeze candidate"
