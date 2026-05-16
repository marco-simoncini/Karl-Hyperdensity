#!/usr/bin/env bash
# KHR-Z: docs language guard — infrastructure scope, no datacenter-only / GA / production-ready claims.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

FAIL=0
SELF="scripts/guard_khr_docs_scope.sh"

SCAN=(
  docs/khr
  docs/adr
  docs/architecture
)

rg_excludes=(
  --glob "!${SELF}"
  --glob "!scripts/guard_khr_iso_boundaries.sh"
  --glob "!**/KHR_RELEASE_READINESS_MAP.md"
  --glob "!**/TECHNICAL_PREVIEW_READINESS.md"
  --glob "!**/TECHNICAL_PREVIEW_READINESS_SUMMARY.md"
  --glob "!**/TECHNICAL_PREVIEW_READINESS_OBSERVATION.md"
  --glob "!**/KHR_TECHNICAL_PREVIEW_PROFILE.md"
  --glob "!**/TECHNICAL_PREVIEW_PACKAGE.md"
  --glob "!**/TECHNICAL_PREVIEW_DASHBOARD_GUIDE.md"
  --glob "!**/TECHNICAL_PREVIEW_INVENTORY_GUIDE.md"
  --glob "!**/TECHNICAL_PREVIEW_ISO_GUIDE.md"
)

line_allowed() {
  local line="$1"
  shopt -s nocasematch
  [[ "${line}" == *"not GA"* ]] && return 0
  [[ "${line}" == *"not production-ready"* ]] && return 0
  [[ "${line}" == *"not production ready"* ]] && return 0
  [[ "${line}" == *"Forbidden"* ]] && return 0
  [[ "${line}" == *"forbidden"* ]] && return 0
  [[ "${line}" == *"not only a datacenter"* ]] && return 0
  [[ "${line}" == *"not only datacenter"* ]] && return 0
  [[ "${line}" == *"infrastructure operating layer"* ]] && return 0
  [[ "${line}" == *"infrastructure OS"* ]] && return 0
  [[ "${line}" == *"infrastructure control plane"* ]] && return 0
  [[ "${line}" == *"deployment environment"* ]] && return 0
  [[ "${line}" == *"Avoid as sole framing"* ]] && return 0
  [[ "${line}" == *"Use instead"* ]] && return 0
  [[ "${line}" == *"Unqualified GA"* ]] && return 0
  [[ "${line}" == *"Sole datacenter"* ]] && return 0
  [[ "${line}" == *"datacenter / bare metal"* ]] && return 0
  [[ "${line}" == *"on-prem"* ]] && return 0
  [[ "${line}" == *"NOT production ready"* ]] && return 0
  [[ "${line}" == *"not production ready"* ]] && return 0
  [[ "${line}" == *"productionReady"* && "${line}" == *"false"* ]] && return 0
  [[ "${line}" == *"Hidden production enable"* ]] && return 0
  [[ "${line}" == *"No hidden production"* ]] && return 0
  shopt -u nocasematch
  return 1
}

check_pattern() {
  local pattern="$1"
  local label="$2"
  local match file line_num line
  while IFS= read -r match; do
    [[ -z "${match}" ]] && continue
    file="${match%%:*}"
    rest="${match#*:}"
    line_num="${rest%%:*}"
    line="${rest#*:}"
    if line_allowed "${line}"; then
      continue
    fi
    echo "[guard_khr_docs_scope] FAIL (${label}): ${file}:${line_num}: ${line}" >&2
    FAIL=1
  done < <(rg -n -i "${pattern}" "${rg_excludes[@]}" "${SCAN[@]}" 2>/dev/null || true)
}

# Datacenter-only framing (not when listed as one environment among others).
check_pattern 'datacenter[- ]only' 'datacenter-only'
check_pattern 'datacenter OS only' 'datacenter-os-only'
check_pattern 'only a datacenter OS' 'only-datacenter-os'
check_pattern 'datacenter-oriented product' 'datacenter-oriented-only'

# Unqualified GA / production-ready claims in KHR docs (preview/sandbox docs must negate explicitly).
check_pattern 'generally available' 'ga-claim'
check_pattern '\bGA\b.{0,20}(KHR|native-live|certification)' 'ga-khr-claim'
check_pattern 'production[- ]ready.{0,30}(KHR|native-live|certified)' 'production-ready-claim'

# Autonomous orchestration as product claim.
check_pattern 'autonomous orchestration.{0,20}(enabled|production)' 'autonomous-orchestration-claim'

if [[ "${FAIL}" -ne 0 ]]; then
  echo "[guard_khr_docs_scope] See docs/khr/KARL_INFRASTRUCTURE_SCOPE.md" >&2
  exit 1
fi

echo "[guard_khr_docs_scope] PASS"
