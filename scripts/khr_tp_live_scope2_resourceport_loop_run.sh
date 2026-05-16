#!/usr/bin/env bash
# KHR-BA: bounded manual ResourcePort loop (observed-json only; no ResourceLease).
set -euo pipefail
# shellcheck source=scripts/khr_tp_live_scope2_lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope2_lib.sh"

ROOT="$(khr_scope2_root)"
cd "${ROOT}"
RUN_DIR="$(khr_scope2_run_dir)"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
ITERATIONS="${KHR_SCOPE2_LOOP_ITERATIONS}"
TIMEOUT_SEC="${KHR_SCOPE2_LOOP_TIMEOUT_SEC}"
INTERVAL_MS="${KHR_SCOPE2_LOOP_INTERVAL_MS}"
CFG_SRC="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-loop.yaml"
PREFLIGHT="${ROOT}/docs/evidence/khr-tp-live-scope2-preflight/committed-scope2-preflight-khr-az/scope2-preflight-summary.json"

mkdir -p "${RUN_DIR}"
log() { khr_scope2_log "$*" | tee -a "${RUN_DIR}/run.log"; }

if [[ ! "${ITERATIONS}" =~ ^[1-3]$ ]]; then
  echo "BLOCKED: KHR_SCOPE2_LOOP_ITERATIONS must be 1-3 (got ${ITERATIONS})" >&2
  exit 2
fi

khr_scope2_require_manual_loop_confirmation
khr_scope2_assert_cluster_context
khr_scope2_assert_sandbox_namespace

if [[ ! -f "${PREFLIGHT}" ]]; then
  echo "BLOCKED: missing scope-2 preflight ${PREFLIGHT}" >&2
  exit 2
fi
python3 -c "import json,sys; d=json.load(open('${PREFLIGHT}')); sys.exit(0 if d.get('status')=='PASS' else 1)" \
  || { echo "BLOCKED: scope-2 preflight not PASS" >&2; exit 2; }

PROD_BEFORE="${RUN_DIR}/production-before.json"
khr_scope2_production_snapshot >"${PROD_BEFORE}"

RUN_CFG="${RUN_DIR}/config-loop-manual-run.yaml"
cp "${CFG_SRC}" "${RUN_CFG}"
khr_scope2_assert_config_safe "${RUN_CFG}"
khr_scope2_assert_no_resourcelease_path "${RUN_CFG}"

BIN="$(khr_scope2_build_binary)"
log "runId=$(basename "${RUN_DIR}") iterations=${ITERATIONS} timeoutSec=${TIMEOUT_SEC} emission=observed-json"

set +e
timeout "${TIMEOUT_SEC}" "${BIN}" \
  -mode=resourceport-loop \
  -config="${RUN_CFG}" \
  -namespace="${NS}" \
  -cluster-context="${CTX}" \
  -loop-iterations="${ITERATIONS}" \
  -loop-interval-ms="${INTERVAL_MS}" \
  >"${RUN_DIR}/loop-output.json" 2>"${RUN_DIR}/loop-stderr.log"
LOOP_RC=$?
set -e

if [[ "${LOOP_RC}" -eq 124 ]]; then
  echo "FAIL: loop timed out after ${TIMEOUT_SEC}s" >&2
  exit 1
fi
if [[ "${LOOP_RC}" -ne 0 ]]; then
  echo "FAIL: loop exit code ${LOOP_RC}" >&2
  cat "${RUN_DIR}/loop-stderr.log" >&2 || true
  exit 1
fi

export RUN_DIR CTX NS ITERATIONS TIMEOUT_SEC LOOP_RC
python3 <<'PY'
import json, os
from datetime import datetime, timezone
from pathlib import Path

run = Path(os.environ["RUN_DIR"])
out = json.loads((run / "loop-output.json").read_text())
ok = (
    out.get("blocked") is False
    and out.get("emissionMode") == "observed-json"
    and out.get("applyCR") is False
    and out.get("applyCRBlocked") is not True
)
doc = {
    "phase": "khr-tp-live-scope2-resourceport-loop-run",
    "sprint": "KHR-BA",
    "runId": run.name,
    "clusterContext": os.environ["CTX"],
    "namespace": os.environ["NS"],
    "status": "PASS" if ok else "FAIL",
    "deployMode": "manual-loop",
    "emissionMode": out.get("emissionMode", "unknown"),
    "emitCR": bool(out.get("emitCR")),
    "applyCR": bool(out.get("applyCR")),
    "loopIterations": int(os.environ["ITERATIONS"]),
    "loopTimeoutSec": int(os.environ["TIMEOUT_SEC"]),
    "loopExitCode": int(os.environ["LOOP_RC"]),
    "resourcePortLoopEnabledDuringRun": True,
    "resourcePortLoopEnabledPersistent": False,
    "sandboxApplyEnabled": False,
    "resourceLeaseApplyEnabled": False,
    "readOnly": True,
    "mutating": False,
    "noResourceLeaseApply": True,
    "noAutonomousOrchestration": True,
    "productionReady": False,
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(run / "loop-summary.json").write_text(json.dumps(doc, indent=2) + "\n")
if not ok:
    raise SystemExit(1)
PY

echo "${RUN_DIR}" > "${RUN_DIR}/run-dir.txt"
log "PASS loop-summary=${RUN_DIR}/loop-summary.json"
echo "[khr_tp_live_scope2_resourceport_loop_run] PASS"
