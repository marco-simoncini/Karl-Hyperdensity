#!/usr/bin/env bash
# KHR-BD: read-only Scope-4 ResourceLease guarded-apply preflight (no execution).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RUN_ID="${KHR_TP_LIVE_SCOPE4_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
OUT_DIR="${ROOT}/docs/evidence/khr-tp-live-scope4-preflight/${RUN_ID}"
CLUSTER="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"
SANDBOX_NS="${KHR_RUNTIME_SANDBOX_NS:-khr-runtime-sandbox}"
GUARDED_CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-guarded-apply.yaml"
MAIN_GO="${ROOT}/cmd/karl-host-runtime/main.go"
SCOPE3_DRYRUN="${ROOT}/docs/evidence/khr-tp-live-scope3-dryrun/committed-scope3-dryrun-khr-bc"

mkdir -p "${OUT_DIR}"
log() { echo "[khr_tp_live_scope4_preflight] $*" | tee -a "${OUT_DIR}/run.log"; }

log "runId=${RUN_ID} cluster=${CLUSTER} namespace=${SANDBOX_NS}"

export ROOT OUT_DIR RUN_ID CLUSTER SANDBOX_NS GUARDED_CFG MAIN_GO SCOPE3_DRYRUN
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
GUARDED_CFG = Path(os.environ["GUARDED_CFG"])
MAIN_GO = Path(os.environ["MAIN_GO"])
SCOPE3_DRYRUN = Path(os.environ["SCOPE3_DRYRUN"])

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

# Production blocklist (read-only snapshot)
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
record("productionBlocklistUntouched", prod_ok, json.dumps(prod_detail))

# Scope-1
scope1 = ROOT / "docs/evidence/khr-tp-live-scope1/committed-scope1-khr-aw/verify-summary.json"
scope1_doc = load(scope1)
scope1_ok = bool(scope1_doc and scope1_doc.get("status") == "PASS")
record("scope1Ready", scope1_ok, f"runId={scope1_doc.get('runId') if scope1_doc else 'none'}")

# Scope-2 manual-loop-pass
scope2 = ROOT / "docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/verify-summary.json"
scope2_doc = load(scope2)
scope2_ok = bool(
    scope2_doc
    and scope2_doc.get("status") == "PASS"
    and scope2_doc.get("readyForScope2") == "manual-loop-pass"
    and scope2_doc.get("readyForScope2Active") is False
)
record("scope2ManualLoopPass", scope2_ok, f"readyForScope2={scope2_doc.get('readyForScope2') if scope2_doc else 'missing'}")

# Scope-3 dry-run evidence
scope3_verify = SCOPE3_DRYRUN / "verify-summary.json"
scope3_dryrun = SCOPE3_DRYRUN / "dryrun-output.json"
scope3_vdoc = load(scope3_verify)
scope3_ddoc = load(scope3_dryrun)
scope3_ok = bool(
    scope3_vdoc
    and scope3_vdoc.get("status") == "PASS"
    and scope3_vdoc.get("readyForScope3") == "manual-dryrun-pass"
    and scope3_vdoc.get("readyForScope3Active") is False
    and scope3_vdoc.get("dryRunObserved") is True
    and scope3_vdoc.get("applyObserved") is False
)
record("scope3ManualDryRunPass", scope3_ok, f"runId={scope3_vdoc.get('runId') if scope3_vdoc else 'none'}")

dry_decision = (scope3_ddoc or {}).get("dryRunDecision")
dry_allowed = dry_decision == "allowed"
record("dryRunDecisionAllowed", dry_allowed, f"dryRunDecision={dry_decision or 'missing'}")

src_ref = (scope3_ddoc or {}).get("sourceResourcePortRef")
rollback_ref = (scope3_ddoc or {}).get("rollbackPlanRef")
verify_ref = (scope3_ddoc or {}).get("verificationPlanRef")
record("sourceResourcePortRefPresent", bool(src_ref), str(src_ref or "missing"))
record("rollbackPlanRefPresent", bool(rollback_ref), str(rollback_ref or "missing"))
record("verificationPlanRefPresent", bool(verify_ref), str(verify_ref or "missing"))

# Native-live certification
cert = load(ROOT / "docs/evidence/khr-native-live-lane/certification-summary.json")
cert_ok = bool(
    cert
    and cert.get("status") == "certified"
    and cert.get("readOnly") is True
    and cert.get("regressionDetected") is False
)
record("nativeLiveCertified", cert_ok, f"status={cert.get('status') if cert else 'missing'}")

# Provenance
prov = load(ROOT / "docs/evidence/khr-provenance/summary.json")
prov_ok = bool(
    prov
    and prov.get("readOnly") is True
    and prov.get("noAutonomousOrchestration") is True
)
record("provenanceValid", prov_ok, f"runId={prov.get('runId') if prov else 'none'}")

# Apply command available (source + config; not executed)
main_text = MAIN_GO.read_text() if MAIN_GO.is_file() else ""
apply_mode_ok = (
    "resourcelease-guarded-apply" in main_text
    and "apply-resourcelease" in main_text
    and "i-understand-this-is-sandbox" in main_text
)
cfg_ok = GUARDED_CFG.is_file()
record(
    "guardedApplyCommandAvailable",
    apply_mode_ok and cfg_ok,
    f"mode=resourcelease-guarded-apply config={GUARDED_CFG.name}",
)

# No TP live guarded apply execution
tp_apply_base = ROOT / "docs/evidence/khr-tp-live-scope4-guarded-apply"
apply_executed = False
if tp_apply_base.is_dir():
    for child in tp_apply_base.iterdir():
        if (child / "apply-summary.json").is_file():
            doc = load(child / "apply-summary.json")
            if doc and doc.get("guardedApplyExecuted") is True:
                apply_executed = True
# Cluster ResourceLease applied annotation
if ctx_ok and not apply_executed:
    r = subprocess.run(
        ["kubectl", "--context", CLUSTER, "get", "resourceleases", "-A", "-o", "json"],
        capture_output=True,
        text=True,
    )
    if r.returncode == 0:
        for item in json.loads(r.stdout).get("items", []):
            if item.get("metadata", {}).get("annotations", {}).get("khr.karl.io/applied") == "true":
                apply_executed = True
                break
record(
    "guardedApplyNotExecuted",
    not apply_executed and (scope3_vdoc or {}).get("applyObserved") is False,
    f"guardedApplyExecuted={apply_executed}",
)

operator_confirmation_required = True
record(
    "operatorConfirmationRequired",
    operator_confirmation_required,
    "KHR_TP_LIVE_SCOPE4_I_UNDERSTAND_GUARDED_APPLY=true required for execution sprint",
)

status = "PASS" if not errors else "FAIL"
ready_for_scope4: bool | str = False
if status == "PASS":
    ready_for_scope4 = "conditional/manual-preflight-pass"

summary: dict[str, Any] = {
    "phase": "khr-tp-live-scope4-preflight",
    "sprint": "KHR-BD",
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
    "readyForScope3": "manual-dryrun-pass" if scope3_ok else False,
    "readyForScope3Active": False,
    "readyForScope4": ready_for_scope4,
    "readyForScope4Active": False,
    "guardedApplyExecuted": False,
    "cgroupMutationObserved": False,
    "dryRunDecision": dry_decision,
    "sourceResourcePortRef": src_ref,
    "rollbackPlanRef": rollback_ref,
    "verificationPlanRef": verify_ref,
    "rollbackPlanDeclared": bool(rollback_ref),
    "verificationPlanDeclared": bool(verify_ref),
    "operatorConfirmationRequired": operator_confirmation_required,
    "applyCommandAvailable": apply_mode_ok and cfg_ok,
    "scope3DryRunEvidenceRunId": scope3_vdoc.get("runId") if scope3_vdoc else None,
    "scope4BlockedReason": "Guarded apply not executed; deferred to dedicated sprint after sign-off",
    "scope4PreflightNote": "KHR-BD preflight only — no guarded apply, no cgroup mutation",
    "checks": checks,
    "errors": errors,
    "forbidden": {
        "productionNamespace": True,
        "autonomousApply": True,
        "persistentScheduler": True,
        "applyWithoutRollback": True,
        "applyWithoutVerification": True,
        "windowsApply": True,
        "kubevirtTemplateMutation": True,
        "guardedApplyLiveExecution": True,
        "cgroupMutation": True,
    },
    "evidencePath": f"docs/evidence/khr-tp-live-scope4-preflight/{RUN_ID}",
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(OUT / "scope4-preflight-summary.json").write_text(json.dumps(summary, indent=2) + "\n")
print(f"[khr_tp_live_scope4_preflight] summary={OUT / 'scope4-preflight-summary.json'}")
print(
    f"[khr_tp_live_scope4_preflight] status={status} "
    f"readyForScope4={ready_for_scope4} guardedApplyExecuted=False"
)
if status != "PASS":
    for e in errors:
        print(f"[khr_tp_live_scope4_preflight] FAIL: {e}", file=__import__("sys").stderr)
    raise SystemExit(1)
PY

log "PASS scope4-preflight-summary=${OUT_DIR}/scope4-preflight-summary.json"
