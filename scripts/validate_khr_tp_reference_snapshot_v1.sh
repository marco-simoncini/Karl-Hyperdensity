#!/usr/bin/env bash
# KHR-BT: validate reference snapshot v1 doc, contract, and aggregator output.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

for f in \
  docs/khr/KHR_TP_REFERENCE_SNAPSHOT_V1.md \
  docs/contracts/khr/khr-tp-reference-snapshot-v1.json \
  scripts/khr_tp_reference_snapshot_v1.sh; do
  [[ -f "${ROOT}/${f}" ]] || {
    echo "[validate_khr_tp_reference_snapshot_v1] FAIL missing ${f}" >&2
    exit 1
  }
done

[[ -x "${ROOT}/scripts/khr_tp_reference_snapshot_v1.sh" ]] || {
  echo "[validate_khr_tp_reference_snapshot_v1] FAIL script not executable" >&2
  exit 1
}

RUN_ID="${KHR_TP_REFERENCE_SNAPSHOT_RUN_ID:-fixture-khr-bt-validate}"
export KHR_TP_REFERENCE_SNAPSHOT_RUN_ID="${RUN_ID}"
"${ROOT}/scripts/khr_tp_reference_snapshot_v1.sh"

SUMMARY="${ROOT}/docs/evidence/khr-tp-reference-snapshot-v1/${RUN_ID}/snapshot-summary.json"
export ROOT SUMMARY
python3 <<'PY'
import json
import os
import sys
from pathlib import Path

ROOT = Path(os.environ["ROOT"])
summary = json.loads(Path(os.environ["SUMMARY"]).read_text(encoding="utf-8"))
required = [
    "contractVersion",
    "contractSetId",
    "scopeReadiness",
    "scope4CertificationState",
    "dashboardLivePassRef",
    "rdpgwClusterSandboxRef",
    "installerCrdFoundationRef",
    "hybridTransitionRef",
    "crossRepoEvidenceIndex",
]
for k in required:
    if k not in summary:
        print(f"[validate_khr_tp_reference_snapshot_v1] FAIL missing {k}", file=sys.stderr)
        sys.exit(1)
if summary.get("contractVersion") != "khr-tp-reference-snapshot-v1":
    sys.exit(1)
if summary.get("globalDefaultsChanged") is not False:
    sys.exit(1)
if summary.get("liveMutationPerformed") is not False:
    sys.exit(1)
live = summary.get("dashboardLivePassRef") or {}
if live.get("evidenceStatus") != "LIVE_PASS":
    print(f"[validate_khr_tp_reference_snapshot_v1] FAIL dashboard LIVE_PASS", file=sys.stderr)
    sys.exit(1)
print("[validate_khr_tp_reference_snapshot_v1] PASS")
PY

echo "[validate_khr_tp_reference_snapshot_v1] OK"
