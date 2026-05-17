#!/usr/bin/env bash
# KHR-BE: mandatory rollback after Scope-4 guarded apply; verify baseline restored.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope4_lib.sh"

ROOT="$(khr_scope4_root)"
cd "${ROOT}"
RUN_ID="${KHR_TP_LIVE_SCOPE4_APPLY_RUN_ID:-$(ls -1t "$(khr_scope4_apply_evidence_base)" 2>/dev/null | head -1)}"
RUN_DIR="$(khr_scope4_apply_evidence_base)/${RUN_ID}"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-guarded-apply.yaml"
TIMEOUT_SEC="${KHR_SCOPE4_APPLY_TIMEOUT_SEC}"
BASELINE_ID="scope4-${RUN_ID}"

[[ -f "${RUN_DIR}/apply-output.json" ]] || { echo "missing apply-output in ${RUN_DIR}" >&2; exit 1; }
[[ -f "${RUN_DIR}/verify-summary.json" ]] || { echo "run verify before rollback" >&2; exit 1; }

khr_scope4_assert_cluster_context
BIN="$(khr_scope4_build_binaries)"

log() { khr_scope4_log "$*" | tee -a "${RUN_DIR}/run.log"; }

if [[ -f "${RUN_DIR}/rollback-summary.json" ]] && [[ "${KHR_TP_LIVE_SCOPE4_ROLLBACK_FORCE_RERUN:-}" != "true" ]]; then
  python3 -c "import json,sys; d=json.load(open('${RUN_DIR}/rollback-summary.json')); sys.exit(0 if d.get('status')=='PASS' else 1)" \
    && { khr_scope4_log "SKIP existing PASS rollback in ${RUN_DIR}"; exit 0; }
fi

log "rollback baseline=${BASELINE_ID}"

set +e
timeout "${TIMEOUT_SEC}" "${BIN}" \
  -mode=resourcelease-rollback \
  -config="${CFG}" \
  -sandbox-dir="${RUN_DIR}/sandbox-work" \
  -baseline-id="${BASELINE_ID}" \
  >"${RUN_DIR}/rollback-output.json" 2>"${RUN_DIR}/rollback-stderr.log"
RC=$?
set -e
[[ "${RC}" -ne 124 ]] || { echo "FAIL: rollback timed out" >&2; exit 1; }
[[ "${RC}" -eq 0 ]] || { echo "FAIL: rollback exit ${RC}" >&2; cat "${RUN_DIR}/rollback-stderr.log" >&2; exit 1; }

export RUN_DIR BASELINE_ID
python3 <<'PY'
import json, os
from datetime import datetime, timezone
from pathlib import Path

run = Path(os.environ["RUN_DIR"])
out = json.loads((run / "rollback-output.json").read_text())
ver = out.get("verification") or {}
ok = (
    out.get("mode") == "resourcelease-rollback"
    and out.get("rolledBack") is True
    and out.get("rollbackState") == "restored"
    and ver.get("state") == "pass"
)
doc = {
    "phase": "khr-tp-live-scope4-guarded-apply-rollback",
    "sprint": "KHR-BE",
    "runId": run.name,
    "status": "PASS" if ok else "FAIL",
    "rollbackObserved": bool(ok),
    "rollbackVerified": bool(ok),
    "baselineId": os.environ["BASELINE_ID"],
    "observedCpuMax": ver.get("observedCpuMax"),
    "expectedCpuMax": ver.get("expectedCpuMax"),
    "cgroupMutationObserved": False,
    "guardedApplyExecuted": True,
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(run / "rollback-summary.json").write_text(json.dumps(doc, indent=2) + "\n")
# Final readiness after mandatory rollback
vpath = run / "verify-summary.json"
if vpath.is_file() and ok:
    v = json.loads(vpath.read_text())
    v["rollbackObserved"] = True
    v["rollbackVerified"] = True
    v["readyForScope4"] = "manual-guarded-apply-pass"
    vpath.write_text(json.dumps(v, indent=2) + "\n")
if not ok:
    raise SystemExit(1)
PY

log "PASS rollback-summary=${RUN_DIR}/rollback-summary.json"
echo "[khr_tp_live_scope4_guarded_apply_rollback] PASS"
