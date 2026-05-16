#!/usr/bin/env bash
# KHR-BC: cleanup temporary Scope-3 dry-run artifacts only (no CRD/cluster cleanup).
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope3_lib.sh"

ROOT="$(khr_scope3_root)"
RUN_ID="${KHR_TP_LIVE_SCOPE3_DRYRUN_RUN_ID:-$(ls -1t "$(khr_scope3_dryrun_evidence_base)" 2>/dev/null | head -1)}"
RUN_DIR="$(khr_scope3_dryrun_evidence_base)/${RUN_ID}"

[[ -d "${RUN_DIR}" ]] || { echo "missing ${RUN_DIR}" >&2; exit 1; }

rm -rf "${RUN_DIR}/sandbox-work" 2>/dev/null || true
rm -f "${RUN_DIR}/dryrun-stderr.log" 2>/dev/null || true

# Evidence JSON retained; no resourceport-cleanup / no CRD deletion
export RUN_DIR
python3 <<'PY'
import json, os
from datetime import datetime, timezone
from pathlib import Path

run = Path(os.environ["RUN_DIR"])
doc = {
    "phase": "khr-tp-live-scope3-dryrun-cleanup",
    "sprint": "KHR-BC",
    "runId": run.name,
    "status": "PASS",
    "crdCleanupPerformed": False,
    "temporaryArtifactsRemoved": True,
    "resourceLeaseApplyEnabled": False,
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(run / "cleanup-summary.json").write_text(json.dumps(doc, indent=2) + "\n")
PY

khr_scope3_log "cleanup complete (temporary only; no CRD cleanup)"
echo "[khr_tp_live_scope3_dryrun_cleanup] PASS"
