#!/usr/bin/env bash
# KHR-AB: Technical Preview package readiness check (docs + evidence anchors + wording).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

FAIL=0
SELF="scripts/khr_tp_package_check.sh"

require_file() {
  local path="$1"
  local label="$2"
  if [[ ! -f "${path}" ]]; then
    echo "[khr_tp_package_check] FAIL: missing ${label}: ${path}" >&2
    FAIL=1
  else
    echo "[khr_tp_package_check] OK: ${label}"
  fi
}

echo "[khr_tp_package_check] Checking TP package anchors..."

require_file "docs/khr/TECHNICAL_PREVIEW_READINESS.md" "TP readiness doc"
require_file "docs/khr/TECHNICAL_PREVIEW_PACKAGE.md" "TP package doc"
require_file "docs/khr/TECHNICAL_PREVIEW_OPERATOR_RUNBOOK.md" "TP operator runbook"
require_file "docs/evidence/khr-native-live-lane/certification-summary.json" "native-live certification summary"
require_file "docs/evidence/khr-certification-registry/summary.json" "certification registry summary"
require_file "docs/evidence/khr-provenance/summary.json" "provenance evidence summary"

# Evidence summaries must assert read-only / no autonomous apply.
for summary in \
  docs/evidence/khr-certification-registry/summary.json \
  docs/evidence/khr-provenance/summary.json; do
  if ! rg -q '"readOnly"[[:space:]]*:[[:space:]]*true' "${summary}" 2>/dev/null; then
    echo "[khr_tp_package_check] FAIL: ${summary} must include readOnly: true" >&2
    FAIL=1
  fi
  if ! rg -q 'noAutonomousOrchestration' "${summary}" 2>/dev/null; then
    echo "[khr_tp_package_check] FAIL: ${summary} must document noAutonomousOrchestration" >&2
    FAIL=1
  fi
done

if [[ -x "${ROOT}/scripts/guard_khr_docs_scope.sh" ]]; then
  "${ROOT}/scripts/guard_khr_docs_scope.sh"
else
  echo "[khr_tp_package_check] FAIL: guard_khr_docs_scope.sh not executable" >&2
  FAIL=1
fi

# TP doc wording guard — TECHNICAL_PREVIEW*.md in docs/khr only.
TP_GLOB='docs/khr/TECHNICAL_PREVIEW*.md'
if ! compgen -G "${TP_GLOB}" >/dev/null; then
  echo "[khr_tp_package_check] FAIL: no ${TP_GLOB} files" >&2
  FAIL=1
fi

tp_line_allowed() {
  local line="$1"
  shopt -s nocasematch
  [[ "${line}" == *"NOT production ready"* ]] && return 0
  [[ "${line}" == *"not production ready"* ]] && return 0
  [[ "${line}" == *"not production-ready"* ]] && return 0
  [[ "${line}" == *"no production enable"* ]] && return 0
  [[ "${line}" == *"No production enable"* ]] && return 0
  [[ "${line}" == *"no autonomous"* ]] && return 0
  [[ "${line}" == *"No autonomous"* ]] && return 0
  [[ "${line}" == *"Autonomous orchestration"*"Absent"* ]] && return 0
  [[ "${line}" == *"Autonomous orchestration"*"Disabled"* ]] && return 0
  [[ "${line}" == *"Autonomous orchestration"*"Forbidden"* ]] && return 0
  [[ "${line}" == *"not GA"* ]] && return 0
  [[ "${line}" == *"Forbidden"* ]] && return 0
  [[ "${line}" == *"disabled by default"* ]] && return 0
  [[ "${line}" == *"Disabled by default"* ]] && return 0
  [[ "${line}" == *"sandbox/manual only"* ]] && return 0
  [[ "${line}" == *"sandbox only"* ]] && return 0
  [[ "${line}" == *"no production enable"* ]] && return 0
  [[ "${line}" == *"production enable"* && "${line}" == *"not"* ]] && return 0
  [[ "${line}" == *"production enable"* && "${line}" == *"No "* ]] && return 0
  [[ "${line}" == *"without production enable"* ]] && return 0
  [[ "${line}" == *"Hidden production enable"* ]] && return 0
  [[ "${line}" == *"not production-enabled"* ]] && return 0
  [[ "${line}" == *"Any production enable"* ]] && return 0
  [[ "${line}" == *"no production enable"* ]] && return 0
  shopt -u nocasematch
  return 1
}

tp_check_pattern() {
  local pattern="$1"
  local label="$2"
  local match file line_num line
  while IFS= read -r match; do
    [[ -z "${match}" ]] && continue
    file="${match%%:*}"
    rest="${match#*:}"
    line_num="${rest%%:*}"
    line="${rest#*:}"
    if tp_line_allowed "${line}"; then
      continue
    fi
    echo "[khr_tp_package_check] FAIL (${label}): ${file}:${line_num}: ${line}" >&2
    FAIL=1
  done < <(rg -n -i "${pattern}" ${TP_GLOB} 2>/dev/null || true)
}

echo "[khr_tp_package_check] Scanning ${TP_GLOB} for forbidden wording..."
tp_check_pattern 'production[- ]ready' 'production-ready-claim'
tp_check_pattern 'production[- ]enable' 'production-enable-claim'
tp_check_pattern 'autonomous orchestration.{0,25}(enabled|production)' 'autonomous-orchestration-enabled'
tp_check_pattern 'generally available' 'ga-claim'
tp_check_pattern 'systemctl enable.{0,30}karl-host-runtime' 'systemd-enable-default'

if [[ "${FAIL}" -ne 0 ]]; then
  echo "[khr_tp_package_check] See docs/khr/TECHNICAL_PREVIEW_PACKAGE.md" >&2
  exit 1
fi

echo "[khr_tp_package_check] PASS"
