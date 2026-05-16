#!/usr/bin/env bash
# KHR-BA: verify Scope-2 manual ResourcePort loop evidence.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope2_lib.sh"

ROOT="$(khr_scope2_root)"
cd "${ROOT}"
RUN_ID="${KHR_TP_LIVE_SCOPE2_LOOP_RUN_ID:-$(ls -1t "$(khr_scope2_evidence_base)" 2>/dev/null | head -1)}"
RUN_DIR="$(khr_scope2_evidence_base)/${RUN_ID}"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"

[[ -f "${RUN_DIR}/loop-summary.json" ]] || { echo "missing loop-summary in ${RUN_DIR}" >&2; exit 1; }

khr_scope2_assert_cluster_context
khr_scope2_assert_sandbox_namespace

PROD_AFTER="${RUN_DIR}/production-after.json"
khr_scope2_production_snapshot >"${PROD_AFTER}"
PROD_BEFORE="${RUN_DIR}/production-before.json"
khr_scope2_assert_production_untouched "${PROD_BEFORE}" "${PROD_AFTER}"

CLUSTER_LOOP="$(khr_scope2_cluster_loop_disabled)"

export ROOT RUN_DIR CTX NS CLUSTER_LOOP
python3 <<'PY'
import json, os, subprocess
from datetime import datetime, timezone
from pathlib import Path

root = Path(os.environ["ROOT"])
run = Path(os.environ["RUN_DIR"])
loop_sum = json.loads((run / "loop-summary.json").read_text())
loop_out = json.loads((run / "loop-output.json").read_text())
ctx, ns = os.environ["CTX"], os.environ["NS"]
cluster_loop = os.environ["CLUSTER_LOOP"] == "true"

scope1 = root / "docs/evidence/khr-tp-live-scope1/committed-scope1-khr-aw/verify-summary.json"
scope1_ok = scope1.is_file() and json.loads(scope1.read_text()).get("status") == "PASS"

rl = subprocess.run(
    ["kubectl", "--context", ctx, "get", "resourceleases", "-A", "-o", "json"],
    capture_output=True, text=True,
)
rl_count = 0
if rl.returncode == 0:
    items = json.loads(rl.stdout).get("items", [])
    rl_count = len([i for i in items if ns in (i.get("metadata", {}).get("namespace") or "")])

ok = (
    loop_sum.get("status") == "PASS"
    and loop_sum.get("emissionMode") == "observed-json"
    and loop_sum.get("resourceLeaseApplyEnabled") is False
    and loop_sum.get("sandboxApplyEnabled") is False
    and loop_out.get("blocked") is False
    and not cluster_loop
    and scope1_ok
)
doc = {
    "phase": "khr-tp-live-scope2-resourceport-loop-verify",
    "sprint": "KHR-BA",
    "runId": run.name,
    "status": "PASS" if ok else "FAIL",
    "readyForScope1": scope1_ok,
    "readyForScope2": "manual-loop-pass" if ok else False,
    "readyForScope2Active": False,
    "readyForScope3": False,
    "scope3BlockedReason": "ResourceLease dry-run/apply blocked until dedicated scope-3 sprint",
    "resourcePortObservationAvailable": loop_out.get("emissionMode") == "observed-json",
    "resourcePortLoopEnabled": False,
    "clusterConfigLoopEnabled": cluster_loop,
    "resourceLeaseApplyEnabled": False,
    "sandboxApplyEnabled": False,
    "emissionMode": loop_sum.get("emissionMode"),
    "productionGatewayUntouched": True,
    "noResourceLeaseDryRunApply": rl_count == 0,
    "readOnly": True,
    "mutating": False,
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(run / "verify-summary.json").write_text(json.dumps(doc, indent=2) + "\n")
print(f"verify status={doc['status']} readyForScope2={doc['readyForScope2']}")
if not ok:
    raise SystemExit(1)
PY

echo "[khr_tp_live_scope2_resourceport_loop_verify] PASS"
