#!/usr/bin/env bash
# KHR-G: validate sandbox evidence or run live when KHR_RUNTIME_SANDBOX_LIVE=1.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
EVIDENCE="${ROOT}/docs/evidence/khr-runtime-sandbox/summary.json"

if [[ "${KHR_RUNTIME_SANDBOX_LIVE:-}" == "1" ]]; then
  bash "${ROOT}/scripts/khr_runtime_sandbox_execute.sh"
fi

if [[ ! -f "${EVIDENCE}" ]]; then
  echo "[validate_khr_runtime_sandbox] FAIL: missing ${EVIDENCE}" >&2
  echo "Run: KHR_RUNTIME_SANDBOX_LIVE=1 ./scripts/khr_runtime_sandbox_execute.sh" >&2
  exit 1
fi

if ! jq -e '.status == "PASS" and .noProductionMutation == true' "${EVIDENCE}" >/dev/null 2>&1; then
  echo "[validate_khr_runtime_sandbox] FAIL: summary not PASS" >&2
  cat "${EVIDENCE}" >&2
  exit 1
fi

echo "[validate_khr_runtime_sandbox] PASS"
