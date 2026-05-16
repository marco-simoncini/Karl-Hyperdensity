#!/usr/bin/env bash
# KHR-BA: cleanup after Scope-2 manual loop (no persistent loop enable).
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope2_lib.sh"

ROOT="$(khr_scope2_root)"
cd "${ROOT}"
RUN_ID="${KHR_TP_LIVE_SCOPE2_LOOP_RUN_ID:-$(ls -1t "$(khr_scope2_evidence_base)" 2>/dev/null | head -1)}"
RUN_DIR="$(khr_scope2_evidence_base)/${RUN_ID}"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG="${RUN_DIR}/config-loop-manual-run.yaml"

[[ -d "${RUN_DIR}" ]] || { echo "missing run dir ${RUN_DIR}" >&2; exit 1; }

khr_scope2_assert_cluster_context
khr_scope2_assert_sandbox_namespace

BIN="$(khr_scope2_build_binary)"
if [[ -f "${CFG}" ]]; then
  "${BIN}" -mode=resourceport-cleanup -config="${CFG}" -namespace="${NS}" \
    -cluster-context="${CTX}" >"${RUN_DIR}/cleanup-output.json" 2>>"${RUN_DIR}/run.log" || true
fi

CLUSTER_LOOP="$(khr_scope2_cluster_loop_disabled)"
PROD_AFTER="${RUN_DIR}/production-after-cleanup.json"
khr_scope2_production_snapshot >"${PROD_AFTER}"

export RUN_DIR CLUSTER_LOOP
python3 <<'PY'
import json, os
from datetime import datetime, timezone
from pathlib import Path

run = Path(os.environ["RUN_DIR"])
cluster_loop = os.environ["CLUSTER_LOOP"] == "true"
ok = not cluster_loop
doc = {
    "phase": "khr-tp-live-scope2-resourceport-loop-cleanup",
    "sprint": "KHR-BA",
    "runId": run.name,
    "status": "PASS" if ok else "FAIL",
    "resourcePortLoopEnabledPersistent": cluster_loop,
    "resourcePortLoopEnabledAfterCleanup": cluster_loop,
    "sandboxApplyEnabled": False,
    "resourceLeaseApplyEnabled": False,
    "productionReady": False,
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(run / "cleanup-summary.json").write_text(json.dumps(doc, indent=2) + "\n")
if not ok:
    print("WARN: cluster config still shows resourcePortLoopEnabled=true", file=__import__("sys").stderr)
    raise SystemExit(1)
PY

khr_scope2_log "cleanup complete resourcePortLoopEnabled persistent=false"
echo "[khr_tp_live_scope2_resourceport_loop_cleanup] PASS"
