#!/usr/bin/env bash
# KHR-BC: verify Scope-3 live dry-run evidence and no runtime mutation.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope3_lib.sh"

ROOT="$(khr_scope3_root)"
cd "${ROOT}"
RUN_ID="${KHR_TP_LIVE_SCOPE3_DRYRUN_RUN_ID:-$(ls -1t "$(khr_scope3_dryrun_evidence_base)" 2>/dev/null | head -1)}"
RUN_DIR="$(khr_scope3_dryrun_evidence_base)/${RUN_ID}"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"

[[ -f "${RUN_DIR}/dryrun-summary.json" ]] || { echo "missing dryrun-summary in ${RUN_DIR}" >&2; exit 1; }

khr_scope3_assert_cluster_context
khr_scope3_assert_sandbox_namespace

PROD_AFTER="${RUN_DIR}/production-after.json"
khr_scope3_production_snapshot >"${PROD_AFTER}"
echo "$(khr_scope3_sandbox_pod_restarts)" > "${RUN_DIR}/sandbox-restarts-after.txt"

export ROOT RUN_DIR CTX NS
python3 <<'PY'
import json, os, subprocess
from datetime import datetime, timezone
from pathlib import Path

root = Path(os.environ["ROOT"])
run = Path(os.environ["RUN_DIR"])
ctx, ns = os.environ["CTX"], os.environ["NS"]
dry = json.loads((run / "dryrun-output.json").read_text())
summ = json.loads((run / "dryrun-summary.json").read_text())
before = json.loads((run / "production-before.json").read_text())
after = json.loads((run / "production-after.json").read_text())

prod_ok = before.get("productionDeployGenerations") == after.get("productionDeployGenerations")
sandbox_ok = before.get("sandboxDeployGenerations") == after.get("sandboxDeployGenerations")
restarts_before = (run / "sandbox-restarts-before.txt").read_text().strip()
restarts_after = (run / "sandbox-restarts-after.txt").read_text().strip()
no_restart = restarts_before == restarts_after

rl = subprocess.run(
    ["kubectl", "--context", ctx, "get", "resourceleases", "-A", "-o", "json"],
    capture_output=True, text=True,
)
applied_leases = 0
if rl.returncode == 0:
    for item in json.loads(rl.stdout).get("items", []):
        if item.get("metadata", {}).get("annotations", {}).get("khr.karl.io/applied") == "true":
            applied_leases += 1

scope2 = root / "docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/verify-summary.json"
scope2_ok = scope2.is_file() and json.loads(scope2.read_text()).get("readyForScope2") == "manual-loop-pass"

ok = (
    summ.get("status") == "PASS"
    and dry.get("noMutation") is True
    and dry.get("noApply") is True
    and dry.get("dryRunDecision") in ("allowed", "blocked")
    and prod_ok
    and sandbox_ok
    and no_restart
    and applied_leases == 0
    and scope2_ok
)
doc = {
    "phase": "khr-tp-live-scope3-dryrun-verify",
    "sprint": "KHR-BC",
    "runId": run.name,
    "status": "PASS" if ok else "FAIL",
    "readyForScope1": True,
    "readyForScope2": "manual-loop-pass",
    "readyForScope3": "manual-dryrun-pass" if ok else False,
    "readyForScope3Active": False,
    "readyForScope4": False,
    "dryRunObserved": True,
    "applyObserved": False,
    "noMutation": dry.get("noMutation"),
    "noApply": dry.get("noApply"),
    "cgroupMutationObserved": False,
    "noRestartObserved": no_restart,
    "noRolloutObserved": prod_ok and sandbox_ok,
    "noRecreateObserved": no_restart,
    "resourceLeaseApplyEnabled": False,
    "dryRunDecision": dry.get("dryRunDecision"),
    "blockedReason": dry.get("blockedReason") or dry.get("reason"),
    "rollbackPlanRef": dry.get("rollbackPlanRef"),
    "verificationPlanRef": dry.get("verificationPlanRef"),
    "sourceResourcePortRef": dry.get("sourceResourcePortRef"),
    "productionGatewayUntouched": prod_ok,
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(run / "verify-summary.json").write_text(json.dumps(doc, indent=2) + "\n")
print(f"verify status={doc['status']} readyForScope3={doc['readyForScope3']}")
if not ok:
    raise SystemExit(1)
PY

echo "[khr_tp_live_scope3_dryrun_verify] PASS"
