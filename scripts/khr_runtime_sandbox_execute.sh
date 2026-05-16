#!/usr/bin/env bash
# KHR-G: full sandbox execution on karl-metal-01@ovh (operator-driven, non-autonomous).
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export KHR_RUNTIME_RUN_ID="${KHR_RUNTIME_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"

for s in \
  khr_runtime_sandbox_preflight.sh \
  khr_runtime_sandbox_dry_run.sh \
  khr_runtime_sandbox_guarded_apply.sh \
  khr_runtime_sandbox_rollback.sh \
  khr_runtime_sandbox_collect_evidence.sh; do
  bash "${ROOT}/scripts/${s}"
done

echo "[khr_runtime_sandbox_execute] PASS run=${KHR_RUNTIME_RUN_ID}"
