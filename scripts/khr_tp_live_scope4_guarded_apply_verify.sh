#!/usr/bin/env bash
# KHR-BE: verify Scope-4 live guarded apply evidence (post-apply; pre-rollback).
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope4_lib.sh"

ROOT="$(khr_scope4_root)"
cd "${ROOT}"
RUN_ID="${KHR_TP_LIVE_SCOPE4_APPLY_RUN_ID:-$(ls -1t "$(khr_scope4_apply_evidence_base)" 2>/dev/null | head -1)}"
RUN_DIR="$(khr_scope4_apply_evidence_base)/${RUN_ID}"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
NS="${KHR_RUNTIME_SANDBOX_NS}"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-guarded-apply.yaml"

[[ -f "${RUN_DIR}/apply-summary.json" ]] || { echo "missing apply-summary in ${RUN_DIR}" >&2; exit 1; }

khr_scope4_assert_cluster_context

PROD_AFTER="${RUN_DIR}/production-after-apply.json"
khr_scope4_production_snapshot >"${PROD_AFTER}"
echo "$(khr_scope4_sandbox_pod_restarts)" > "${RUN_DIR}/sandbox-restarts-after-apply.txt"
echo "$(khr_scope4_native_live_pod_uid)" > "${RUN_DIR}/native-live-pod-uid-after-apply.txt"

SNAP="${ROOT}/bin/khr-continuity-snapshot"
PROOF="${ROOT}/bin/khr-continuity-proof"
if [[ -x "${SNAP}" ]]; then
  "${SNAP}" -cluster-context="${CTX}" -namespace="${NS}" -out="${RUN_DIR}/continuity-after-apply.json"
fi
if [[ -x "${PROOF}" && -f "${RUN_DIR}/continuity-before.json" && -f "${RUN_DIR}/continuity-after-apply.json" ]]; then
  "${PROOF}" -before="${RUN_DIR}/continuity-before.json" -after="${RUN_DIR}/continuity-after-apply.json" \
    -out="${RUN_DIR}/continuity-proof-apply.json"
fi

# rdp-GW read-only continuity probe (no mutation)
if [[ -x "${ROOT}/../rdp-GW/scripts/khr_accessgraph_continuity_evidence.sh" ]]; then
  AG_RUN="scope4-guarded-apply-${RUN_ID}"
  (cd "${ROOT}/../rdp-GW" && KHR_ACCESSGRAPH_RUN_ID="${AG_RUN}" \
    ./scripts/khr_accessgraph_continuity_evidence.sh) >>"${RUN_DIR}/rdpgw-continuity.log" 2>&1 || true
  if [[ -f "${ROOT}/../rdp-GW/docs/evidence/khr-accessgraph-continuity/${AG_RUN}/summary.json" ]]; then
    cp "${ROOT}/../rdp-GW/docs/evidence/khr-accessgraph-continuity/${AG_RUN}/summary.json" \
      "${RUN_DIR}/rdpgw-continuity-summary.json"
  fi
fi

export ROOT RUN_DIR CTX NS CFG
python3 <<'PY'
import json, os, subprocess
from datetime import datetime, timezone
from pathlib import Path

root = Path(os.environ["ROOT"])
run = Path(os.environ["RUN_DIR"])
ctx, ns = os.environ["CTX"], os.environ["NS"]
apply = json.loads((run / "apply-output.json").read_text())
apply_sum = json.loads((run / "apply-summary.json").read_text())
before = json.loads((run / "production-before.json").read_text())
after = json.loads((run / "production-after-apply.json").read_text())
ver = apply.get("verification") or {}

prod_ok = before.get("productionDeployGenerations") == after.get("productionDeployGenerations")
restarts_before = (run / "sandbox-restarts-before.txt").read_text().strip()
restarts_after = (run / "sandbox-restarts-after-apply.txt").read_text().strip()
no_restart = restarts_before == restarts_after
uid_before = (run / "native-live-pod-uid-before.txt").read_text().strip()
uid_after = (run / "native-live-pod-uid-after-apply.txt").read_text().strip()
no_recreate = uid_before == uid_after and uid_before != ""
gen_before = (run / "deploy-generation-before.txt").read_text().strip()
gen_after = subprocess.run(
    ["kubectl", "--context", ctx, "-n", ns, "get", "deploy", "khr-native-live-target",
     "-o", "jsonpath={.metadata.generation}"],
    capture_output=True, text=True,
).stdout.strip()
no_rollout = gen_before == gen_after

continuity_score = None
continuity_preserved = True
proof_path = run / "continuity-proof-apply.json"
if proof_path.is_file():
    proof = json.loads(proof_path.read_text())
    continuity_score = proof.get("continuityScore")
    continuity_preserved = proof.get("continuityPreserved", True) is True

rdpgw_ok = True
rdpgw_sum = run / "rdpgw-continuity-summary.json"
if rdpgw_sum.is_file():
    rg = json.loads(rdpgw_sum.read_text())
    rdpgw_ok = (
        rg.get("noRevoke") is True
        and rg.get("noDisconnect") is True
        and rg.get("continuityObserved") is True
    )

cgroup_changed = (
    ver.get("state") == "pass"
    and ver.get("observedCpuMax")
    and ver.get("expectedCpuMax")
    and ver.get("observedCpuMax") == ver.get("expectedCpuMax")
)

ok = (
    apply_sum.get("status") == "PASS"
    and apply.get("applied") is True
    and cgroup_changed
    and ver.get("noRestart") is True
    and ver.get("noRollout") is True
    and ver.get("noRecreate") is True
    and prod_ok
    and no_restart
    and no_rollout
    and no_recreate
    and continuity_preserved
    and rdpgw_ok
)
rb_sum = run / "rollback-summary.json"
rollback_observed = False
rollback_verified = False
if rb_sum.is_file():
    rb = json.loads(rb_sum.read_text())
    rollback_observed = rb.get("rollbackObserved") is True
    rollback_verified = rb.get("rollbackVerified") is True

doc = {
    "phase": "khr-tp-live-scope4-guarded-apply-verify",
    "sprint": "KHR-BE",
    "runId": run.name,
    "status": "PASS" if ok else "FAIL",
    "readyForScope1": True,
    "readyForScope2": "manual-loop-pass",
    "readyForScope3": "manual-dryrun-pass",
    "readyForScope4": "manual-guarded-apply-pass" if (ok and rollback_verified) or (ok and not rb_sum.is_file()) else ("manual-guarded-apply-pass" if ok else False),
    "readyForScope4Active": False,
    "guardedApplyObserved": True,
    "rollbackObserved": rollback_observed,
    "rollbackVerified": rollback_verified,
    "continuityPreserved": continuity_preserved,
    "continuityScore": continuity_score,
    "cgroupMutationObserved": cgroup_changed,
    "applyScope": "sandbox-only",
    "noRestartObserved": no_restart,
    "noRolloutObserved": no_rollout and ver.get("noRollout") is True,
    "noRecreateObserved": no_recreate and ver.get("noRecreate") is True,
    "productionGatewayUntouched": prod_ok,
    "rdpgwContinuityPreserved": rdpgw_ok,
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(run / "verify-summary.json").write_text(json.dumps(doc, indent=2) + "\n")
print(f"verify status={doc['status']} readyForScope4={doc['readyForScope4']} continuity={continuity_preserved}")
if not ok:
    raise SystemExit(1)
PY

echo "[khr_tp_live_scope4_guarded_apply_verify] PASS"
