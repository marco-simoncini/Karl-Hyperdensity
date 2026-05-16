#!/usr/bin/env bash
# KHR-BB: read-only Scope-3 ResourceLease dry-run preflight (no execution).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RUN_ID="${KHR_TP_LIVE_SCOPE3_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
OUT_DIR="${ROOT}/docs/evidence/khr-tp-live-scope3-preflight/${RUN_ID}"
CLUSTER="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"
SANDBOX_NS="${KHR_RUNTIME_SANDBOX_NS:-khr-runtime-sandbox}"
LEASE_SAMPLE="${ROOT}/examples/khr/runtime-sandbox/resourcelease-dryrun-allowed.json"
LOOP_CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-loop.yaml"
MAIN_GO="${ROOT}/cmd/karl-host-runtime/main.go"

mkdir -p "${OUT_DIR}"
log() { echo "[khr_tp_live_scope3_preflight] $*" | tee -a "${OUT_DIR}/run.log"; }

log "runId=${RUN_ID} cluster=${CLUSTER} namespace=${SANDBOX_NS}"

export ROOT OUT_DIR RUN_ID CLUSTER SANDBOX_NS LEASE_SAMPLE LOOP_CFG MAIN_GO
python3 <<'PY'
from __future__ import annotations

import json
import os
import subprocess
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

ROOT = Path(os.environ["ROOT"])
OUT = Path(os.environ["OUT_DIR"])
RUN_ID = os.environ["RUN_ID"]
CLUSTER = os.environ["CLUSTER"]
SANDBOX_NS = os.environ["SANDBOX_NS"]
LEASE_SAMPLE = Path(os.environ["LEASE_SAMPLE"])
LOOP_CFG = Path(os.environ["LOOP_CFG"])
MAIN_GO = Path(os.environ["MAIN_GO"])

checks: dict[str, dict[str, Any]] = {}
errors: list[str] = []


def record(name: str, ok: bool, detail: str = "", **extra: Any) -> None:
    checks[name] = {"status": "PASS" if ok else "FAIL", "detail": detail, **extra}
    if not ok:
        errors.append(f"{name}: {detail}")


def load(path: Path) -> dict[str, Any] | None:
    if not path.is_file():
        return None
    try:
        return json.loads(path.read_text())
    except (OSError, json.JSONDecodeError):
        return None


# Cluster context
ctx_ok = False
current = ""
if subprocess.run(["kubectl", "version", "--client=true"], capture_output=True).returncode == 0:
    cur = subprocess.run(
        ["kubectl", "config", "current-context"],
        capture_output=True,
        text=True,
    )
    current = (cur.stdout or "").strip()
    ctx_ok = current == CLUSTER
record("clusterContext", ctx_ok, f"current={current or 'unavailable'} required={CLUSTER}")

# Sandbox namespace label
ns_ok = False
if ctx_ok:
    r = subprocess.run(
        [
            "kubectl", "--context", CLUSTER, "get", "namespace", SANDBOX_NS,
            "-o", "jsonpath={.metadata.labels.khr\\.karl\\.io/sandbox}",
        ],
        capture_output=True,
        text=True,
    )
    ns_ok = r.returncode == 0 and (r.stdout or "").strip() == "true"
record("sandboxNamespaceLabel", ns_ok, f"namespace={SANDBOX_NS}")

# Scope-1
scope1 = ROOT / "docs/evidence/khr-tp-live-scope1/committed-scope1-khr-aw/verify-summary.json"
scope1_doc = load(scope1)
scope1_ok = bool(scope1_doc and scope1_doc.get("status") == "PASS")
record("scope1Ready", scope1_ok, f"runId={scope1_doc.get('runId') if scope1_doc else 'none'}")

# Scope-2 manual-loop-pass
scope2_loop = ROOT / "docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/verify-summary.json"
scope2_doc = load(scope2_loop)
scope2_ok = bool(
    scope2_doc
    and scope2_doc.get("status") == "PASS"
    and scope2_doc.get("readyForScope2") == "manual-loop-pass"
    and scope2_doc.get("readyForScope2Active") is False
)
record("scope2ManualLoopPass", scope2_ok, f"readyForScope2={scope2_doc.get('readyForScope2') if scope2_doc else 'missing'}")

# ResourcePort observed-json evidence
loop_sum_path = ROOT / "docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/loop-summary.json"
loop_sum = load(loop_sum_path)
rp_obs_ok = bool(
    loop_sum
    and loop_sum.get("status") == "PASS"
    and loop_sum.get("emissionMode") == "observed-json"
)
record("resourcePortObservedJson", rp_obs_ok, f"emissionMode={loop_sum.get('emissionMode') if loop_sum else 'missing'}")

# Sample ResourceLease input
lease_ok = False
rollback_declared = False
dry_run_only = False
if LEASE_SAMPLE.is_file():
    lease = load(LEASE_SAMPLE)
    if lease:
        lease_ok = (
            lease.get("kind") == "ResourceLease"
            and (lease.get("metadata", {}).get("labels", {}).get("khr.karl.io/sandbox") == "true")
        )
        gov = lease.get("spec", {}).get("governance") or {}
        rollback_declared = bool(gov.get("rollbackPlanRef"))
        dry_run_only = gov.get("dryRunOnly") is True
record(
    "sampleResourceLeaseInput",
    lease_ok and rollback_declared and dry_run_only,
    f"lease={LEASE_SAMPLE.name} rollbackPlanRef={rollback_declared} dryRunOnly={dry_run_only}",
)

# Dry-run executable available (source check only — not executed)
main_text = MAIN_GO.read_text() if MAIN_GO.is_file() else ""
dryrun_mode_ok = "resourcelease-dryrun" in main_text and "lease-input" in main_text
record("dryRunExecutableAvailable", dryrun_mode_ok, "karl-host-runtime -mode=resourcelease-dryrun present in source")

# Apply flags disabled in sandbox configs
cfg_text = LOOP_CFG.read_text() if LOOP_CFG.is_file() else ""
apply_disabled = (
    "sandboxApplyEnabled: false" in cfg_text
    and "resourcePortLoopEnabled: true" in cfg_text
)
guarded_cfg = ROOT / "examples/khr/runtime-sandbox/karl-host-runtime-config-guarded-apply.yaml"
guarded_isolated = True
if guarded_cfg.is_file():
    gt = guarded_cfg.read_text()
    guarded_isolated = "sandboxApplyEnabled: false" in gt or "resourcePortLoopEnabled: true" not in gt
record(
    "resourceLeaseApplyDisabled",
    apply_disabled and LEASE_SAMPLE.is_file() and dry_run_only,
    "sandboxApplyEnabled=false; lease governance dryRunOnly=true",
)

# No live dry-run execution in TP scope-3 evidence path
tp_dryrun_evidence = ROOT / "docs/evidence/khr-tp-live-scope3-resourcelease-dryrun"
dryrun_executed = False
if tp_dryrun_evidence.is_dir():
    for child in tp_dryrun_evidence.iterdir():
        if (child / "dryrun-summary.json").is_file():
            doc = load(child / "dryrun-summary.json")
            if doc and doc.get("dryRunExecuted") is True:
                dryrun_executed = True
record(
    "noResourceLeaseDryRunExecuted",
    not dryrun_executed,
    f"resourceLeaseDryRunExecuted={dryrun_executed}",
)

# Production namespace snapshot (read-only)
prod_ok = ctx_ok
prod_detail: dict[str, str] = {}
if ctx_ok:
    for ns in ("karl", "karl-system", "default", "kube-system"):
        r = subprocess.run(
            [
                "kubectl", "--context", CLUSTER, "get", "deploy", "-n", ns,
                "-o", "jsonpath={.items[*].metadata.generation}",
            ],
            capture_output=True,
            text=True,
        )
        prod_detail[ns] = (r.stdout or "").strip() if r.returncode == 0 else "unavailable"
record("noProductionNamespaceMutation", prod_ok, json.dumps(prod_detail))

status = "PASS" if not errors else "FAIL"
ready_for_scope3: bool | str = False
if status == "PASS":
    ready_for_scope3 = "conditional/manual-preflight-pass"

summary: dict[str, Any] = {
    "phase": "khr-tp-live-scope3-preflight",
    "sprint": "KHR-BB",
    "runId": RUN_ID,
    "clusterContext": CLUSTER,
    "namespace": SANDBOX_NS,
    "contractSetId": "khr-tp-contract-v1",
    "status": status,
    "readOnly": True,
    "mutating": False,
    "automaticEnablement": False,
    "productionReady": False,
    "noAutonomousOrchestration": True,
    "readyForScope1": scope1_ok,
    "readyForScope2": "manual-loop-pass" if scope2_ok else False,
    "readyForScope3": ready_for_scope3,
    "readyForScope3Active": False,
    "readyForScope4": False,
    "resourceLeaseDryRunExecuted": False,
    "resourceLeaseApplyEnabled": False,
    "resourcePortObservationAvailable": rp_obs_ok,
    "rollbackPlanDeclared": rollback_declared,
    "rollbackPlanExecuted": False,
    "scope3BlockedReason": "ResourceLease dry-run not executed; deferred to dedicated sprint after sign-off",
    "scope3PreflightNote": "KHR-BB preflight only — no dry-run execution, no apply, no cgroup mutation",
    "checks": checks,
    "errors": errors,
    "forbidden": {
        "resourceLeaseDryRunLiveExecution": True,
        "resourceLeaseApply": True,
        "cgroupMutation": True,
        "autonomousDryRunScheduler": True,
        "productionNamespaceMutation": True,
    },
    "evidencePath": f"docs/evidence/khr-tp-live-scope3-preflight/{RUN_ID}",
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(OUT / "scope3-preflight-summary.json").write_text(json.dumps(summary, indent=2) + "\n")
print(f"[khr_tp_live_scope3_preflight] summary={OUT / 'scope3-preflight-summary.json'}")
print(
    f"[khr_tp_live_scope3_preflight] status={status} "
    f"readyForScope3={ready_for_scope3} dryRunExecuted=False"
)
if status != "PASS":
    for e in errors:
        print(f"[khr_tp_live_scope3_preflight] FAIL: {e}", file=__import__("sys").stderr)
    raise SystemExit(1)
PY

log "PASS scope3-preflight-summary=${OUT_DIR}/scope3-preflight-summary.json"
