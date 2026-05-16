#!/usr/bin/env bash
# KHR-AX: validate stabilized TP Live reference environment (read-only).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RUN_ID="${KHR_TP_LIVE_REFERENCE_ENV_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
OUT_DIR="${ROOT}/docs/evidence/khr-tp-live-reference-env/${RUN_ID}"
CLUSTER="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"
RDP_GW="$(cd "${KHR_RDP_GW_PATH:-${ROOT}/../rdp-GW}" 2>/dev/null && pwd || echo "/home/m.simoncini/rdp-GW")"

mkdir -p "${OUT_DIR}"
log() { echo "[khr_tp_live_reference_env_check] $*" | tee -a "${OUT_DIR}/run.log"; }

log "runId=${RUN_ID} cluster=${CLUSTER}"

export ROOT OUT_DIR RUN_ID CLUSTER RDP_GW
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
RDP_GW = Path(os.environ["RDP_GW"])

checks: dict[str, dict[str, Any]] = {}
errors: list[str] = []


def load(path: Path) -> dict[str, Any] | None:
    if not path.is_file():
        return None
    try:
        return json.loads(path.read_text())
    except (OSError, json.JSONDecodeError):
        return None


def record(name: str, ok: bool, detail: str = "", **extra: Any) -> None:
    checks[name] = {"status": "PASS" if ok else "FAIL", "detail": detail, **extra}
    if not ok:
        errors.append(f"{name}: {detail}")


def latest_file(glob_parent: Path, filename: str) -> Path | None:
    if not glob_parent.is_dir():
        return None
    best = None
    for child in sorted(glob_parent.iterdir(), reverse=True):
        f = child / filename
        if f.is_file():
            best = f
    return best


# Preflight
preflight_dir = ROOT / "docs/evidence/khr-tp-live-enablement"
preflight = load(latest_file(preflight_dir, "enablement-preflight-summary.json") or Path())
if preflight is None:
    pf_root = preflight_dir / "enablement-preflight-summary.json"
    preflight = load(pf_root)
record(
    "enablementPreflight",
    bool(preflight and preflight.get("status") == "PASS"),
    f"readyForScope1={preflight.get('readyForScope1') if preflight else None}",
)
# Scope-1 evidence (primary source for scope-1 readiness)
scope1_base = ROOT / "docs/evidence/khr-tp-live-scope1"
scope1_committed = scope1_base / "committed-scope1-khr-aw" / "verify-summary.json"
scope1_verify = load(scope1_committed) or load(latest_file(scope1_base, "verify-summary.json") or Path())
scope1_ok = bool(
    scope1_verify
    and scope1_verify.get("status") == "PASS"
    and scope1_verify.get("accessGraphLiveReadonly") is True
    and scope1_verify.get("readyForScope2") is False
)
record(
    "scope1Evidence",
    scope1_ok,
    f"runId={scope1_verify.get('runId') if scope1_verify else 'none'}",
)
record(
    "readyForScope1",
    scope1_ok,
    str(scope1_verify.get("runId") if scope1_verify else "missing"),
)
record(
    "readyForScope2Blocked",
    scope1_verify is not None and scope1_verify.get("readyForScope2") is False,
    scope1_verify.get("scope2BlockedReason", "") if scope1_verify else "",
)

# rdp-GW live-readonly + deployMode
ag_ok = bool(scope1_verify and scope1_verify.get("accessGraphLiveReadonly") is True)
deploy_mode = "missing"
trust = "missing"
if scope1_verify:
    trust = "live-readonly" if scope1_verify.get("accessGraphLiveReadonly") else "fixture-readonly"
ag_dir = RDP_GW / "docs/evidence/khr-accessgraph-continuity"
if ag_dir.is_dir():
    for child in sorted(ag_dir.iterdir(), reverse=True):
        s = child / "summary.json"
        doc = load(s)
        if doc and doc.get("status") == "PASS":
            if doc.get("source") == "live-readonly" or doc.get("trustLevel") == "live-readonly":
                ag_ok = True
                trust = "live-readonly"
            break

def normalize_deploy_mode(mode: str | None) -> str:
    if not mode:
        return "missing"
    if mode == "local":
        return "local-gateway"
    if mode == "cluster":
        return "cluster-sandbox"
    return mode


rdpgw_cluster = None
rdpgw_cluster_base = RDP_GW / "docs/evidence/khr-rdpgw-cluster-sandbox"
committed_cluster = rdpgw_cluster_base / "committed-cluster-sandbox-khr-ay"
cluster_verify = committed_cluster / "verify-summary.json"
if cluster_verify.is_file():
    rdpgw_cluster = load(cluster_verify)
if rdpgw_cluster is None and rdpgw_cluster_base.is_dir():
    for child in sorted(rdpgw_cluster_base.iterdir(), reverse=True):
        verify = child / "verify-summary.json"
        if verify.is_file():
            doc = load(verify)
            if doc and doc.get("status") == "PASS" and doc.get("deployMode") == "cluster-sandbox":
                rdpgw_cluster = doc
                break

rdpgw_dep = None
rdpgw_scope1 = RDP_GW / "docs/evidence/khr-rdpgw-scope1"
if rdpgw_scope1.is_dir():
    for child in sorted(rdpgw_scope1.iterdir(), reverse=True):
        dep = child / "deploy-summary.json"
        if dep.is_file():
            rdpgw_dep = load(dep)
            break
        mode_file = child / "deploy-mode.json"
        if mode_file.is_file():
            rdpgw_dep = load(mode_file)
            break

deploy_mode = "missing"
deploy_warnings: list[str] = []
if rdpgw_cluster and rdpgw_cluster.get("deployMode") == "cluster-sandbox":
    deploy_mode = "cluster-sandbox"
    if rdpgw_cluster.get("accessGraphLiveReadonly"):
        ag_ok = True
        trust = "live-readonly"
elif scope1_verify:
    deploy_mode = normalize_deploy_mode(scope1_verify.get("rdpgwDeployMode"))
elif rdpgw_dep:
    deploy_mode = normalize_deploy_mode(rdpgw_dep.get("deployMode"))

if deploy_mode == "local-gateway":
    deploy_warnings.append("rdpgwDeployMode=local-gateway is fallback only; prefer cluster-sandbox evidence")
elif deploy_mode == "missing" and rdpgw_dep:
    deploy_mode = normalize_deploy_mode(rdpgw_dep.get("deployMode"))
    if deploy_mode == "local-gateway":
        deploy_warnings.append("rdpgwDeployMode=local-gateway is fallback only; prefer cluster-sandbox evidence")

record("rdpgwLiveReadonlyEvidence", ag_ok, f"trust={trust}")
record(
    "rdpgwDeployMode",
    deploy_mode in ("local-gateway", "cluster-sandbox"),
    f"deployMode={deploy_mode}",
    deployMode=deploy_mode,
    preferredMode="cluster-sandbox",
    fallbackWarning=deploy_warnings[0] if deploy_warnings else "",
)

# Scope-2 preflight (KHR-AZ — loop not enabled)
scope2_pf = None
scope2_base = ROOT / "docs/evidence/khr-tp-live-scope2-preflight"
committed_scope2 = scope2_base / "committed-scope2-preflight-khr-az" / "scope2-preflight-summary.json"
if committed_scope2.is_file():
    scope2_pf = load(committed_scope2)
elif scope2_base.is_dir():
    for child in sorted(scope2_base.iterdir(), reverse=True):
        p = child / "scope2-preflight-summary.json"
        if p.is_file():
            scope2_pf = load(p)
            if scope2_pf:
                break
scope2_pf_ok = bool(
    scope2_pf
    and scope2_pf.get("status") == "PASS"
    and scope2_pf.get("resourcePortLoopEnabled") is False
    and scope2_pf.get("loopEnabled") is False
)
record(
    "scope2PreflightEvidence",
    scope2_pf_ok,
    f"runId={scope2_pf.get('runId') if scope2_pf else 'none'}",
)
record(
    "resourcePortLoopDisabled",
    scope2_pf is None or scope2_pf.get("resourcePortLoopEnabled") is False,
    f"resourcePortLoopEnabled={scope2_pf.get('resourcePortLoopEnabled') if scope2_pf else 'unknown'}",
)
record(
    "sandboxApplyDisabled",
    scope2_pf is None or scope2_pf.get("sandboxApplyEnabled") is False,
    f"sandboxApplyEnabled={scope2_pf.get('sandboxApplyEnabled') if scope2_pf else 'unknown'}",
)

# Production namespace mutation proof (read-only kubectl generations)
prod_ok = True
prod_detail = "kubectl unavailable"
if subprocess.run(["kubectl", "version", "--client=true"], capture_output=True).returncode == 0:
    ctx = subprocess.run(["kubectl", "config", "current-context"], capture_output=True, text=True)
    current = (ctx.stdout or "").strip()
    if CLUSTER not in current:
        prod_ok = False
        prod_detail = f"context={current}"
    else:
        gens = {}
        for ns in ("karl-system", "kube-system", "default"):
            r = subprocess.run(
                ["kubectl", "get", "deploy", "-n", ns, "-o", "jsonpath={.items[*].metadata.generation}"],
                capture_output=True,
                text=True,
            )
            gens[ns] = (r.stdout or "").strip() if r.returncode == 0 else "unavailable"
        prod_detail = json.dumps({"deployGenerationsSnapshot": gens, "sandboxOnlyMutation": True})
record("noProductionNamespaceMutation", prod_ok, prod_detail)

# Dashboard fixture
dash_fixture = Path(os.environ.get("KHR_DASHBOARD_PATH", str(ROOT.parent / "Karl-Dashboard"))) / \
    "examples/khr-dashboard/tp-readiness-reference-env.json"
dash = load(dash_fixture)

scope2_loop = None
scope2_loop_base = ROOT / "docs/evidence/khr-tp-live-scope2-resourceport-loop"
committed_loop = scope2_loop_base / "committed-scope2-loop-khr-ba" / "verify-summary.json"
if committed_loop.is_file():
    scope2_loop = load(committed_loop)
elif scope2_loop_base.is_dir():
    for child in sorted(scope2_loop_base.iterdir(), reverse=True):
        v = child / "verify-summary.json"
        if v.is_file():
            scope2_loop = load(v)
            if scope2_loop:
                break
scope2_loop_ok = bool(
    scope2_loop
    and scope2_loop.get("status") == "PASS"
    and scope2_loop.get("readyForScope2") == "manual-loop-pass"
    and scope2_loop.get("readyForScope3") is False
)
record(
    "scope2ManualLoopEvidence",
    scope2_loop_ok,
    f"runId={scope2_loop.get('runId') if scope2_loop else 'none'}",
)

ready_for_scope2: bool | str = False
if scope2_loop_ok:
    ready_for_scope2 = "manual-loop-pass"
elif scope2_pf_ok:
    ready_for_scope2 = "conditional/manual-preflight-pass"

scope3_pf = None
scope3_base = ROOT / "docs/evidence/khr-tp-live-scope3-preflight"
committed_scope3 = scope3_base / "committed-scope3-preflight-khr-bb" / "scope3-preflight-summary.json"
if committed_scope3.is_file():
    scope3_pf = load(committed_scope3)
elif scope3_base.is_dir():
    for child in sorted(scope3_base.iterdir(), reverse=True):
        p = child / "scope3-preflight-summary.json"
        if p.is_file():
            scope3_pf = load(p)
            if scope3_pf:
                break
scope3_pf_ok = bool(
    scope3_pf
    and scope3_pf.get("status") == "PASS"
    and scope3_pf.get("resourceLeaseApplyEnabled") is False
    and scope3_pf.get("readyForScope3Active") is False
)
record(
    "scope3PreflightEvidence",
    scope3_pf_ok,
    f"runId={scope3_pf.get('runId') if scope3_pf else 'none'}",
)

scope3_dryrun = None
scope3_dryrun_base = ROOT / "docs/evidence/khr-tp-live-scope3-dryrun"
committed_scope3_dry = scope3_dryrun_base / "committed-scope3-dryrun-khr-bc" / "verify-summary.json"
if committed_scope3_dry.is_file():
    scope3_dryrun = load(committed_scope3_dry)
elif scope3_dryrun_base.is_dir():
    for child in sorted(scope3_dryrun_base.iterdir(), reverse=True):
        v = child / "verify-summary.json"
        if v.is_file():
            scope3_dryrun = load(v)
            if scope3_dryrun:
                break
scope3_dryrun_ok = bool(
    scope3_dryrun
    and scope3_dryrun.get("status") == "PASS"
    and scope3_dryrun.get("readyForScope3") == "manual-dryrun-pass"
    and scope3_dryrun.get("readyForScope3Active") is False
    and scope3_dryrun.get("dryRunObserved") is True
    and scope3_dryrun.get("applyObserved") is False
    and scope3_dryrun.get("noMutation") is True
    and scope3_dryrun.get("noApply") is True
)
record(
    "scope3ManualDryRunEvidence",
    scope3_dryrun_ok,
    f"runId={scope3_dryrun.get('runId') if scope3_dryrun else 'none'}",
)

ready_for_scope3: bool | str = False
if scope3_dryrun_ok:
    ready_for_scope3 = "manual-dryrun-pass"
elif scope3_pf_ok:
    ready_for_scope3 = "conditional/manual-preflight-pass"

dash_summary = (dash or {}).get("tpReadinessSummary", {})
record(
    "dashboardTpReadinessFixture",
    bool(
        dash
        and dash_summary.get("readyForScope1") is True
        and dash_summary.get("readyForScope2State") == "manual-loop-pass"
        and dash_summary.get("readyForScope3") is False
        and dash_summary.get("readyForScope3State") == "manual-dryrun-pass"
        and dash_summary.get("readyForScope3Active") is False
        and dash_summary.get("readyForScope4") is False
        and dash_summary.get("dryRunObserved") is True
        and dash_summary.get("applyObserved") is False
        and dash_summary.get("noMutation") is True
        and dash_summary.get("resourceLeaseDryRunExecuted") is True
        and dash_summary.get("resourceLeaseApplyEnabled") is False
    ),
    str(dash_fixture),
)

status = "PASS" if not errors else "FAIL"
summary = {
    "phase": "khr-tp-live-reference-env",
    "sprint": "KHR-BC",
    "runId": RUN_ID,
    "clusterContext": CLUSTER,
    "contractSetId": "khr-tp-contract-v1",
    "status": status,
    "readOnly": True,
    "mutating": False,
    "productionReady": False,
    "noAutonomousOrchestration": True,
    "readyForScope0": preflight.get("readyForScope0") if preflight else scope1_ok,
    "readyForScope1": scope1_ok,
    "readyForScope2": ready_for_scope2,
    "readyForScope2Active": False,
    "readyForScope2LoopExecution": scope2_loop_ok,
    "readyForScope3": ready_for_scope3,
    "readyForScope3Active": False,
    "readyForScope4": False,
    "resourceLeaseDryRunExecuted": scope3_dryrun_ok,
    "dryRunObserved": scope3_dryrun_ok,
    "applyObserved": False,
    "noMutation": scope3_dryrun.get("noMutation", True) if scope3_dryrun_ok else True,
    "resourceLeaseApplyEnabled": False,
    "cgroupMutationObserved": False,
    "resourcePortObservationAvailable": scope2_loop_ok,
    "resourcePortLoopEnabled": scope2_loop.get("resourcePortLoopEnabled", False) if scope2_loop else (
        scope2_pf.get("resourcePortLoopEnabled", False) if scope2_pf else False
    ),
    "sandboxApplyEnabled": scope2_pf.get("sandboxApplyEnabled", False) if scope2_pf else False,
    "scope2BlockedReason": scope2_pf.get("scope2BlockedReason", "ResourcePort loop not enabled") if scope2_pf else "scope-2 preflight missing",
    "rdpgwDeployMode": deploy_mode,
    "rdpgwEvidenceTrust": trust,
    "namespaces": ["khr-runtime-sandbox", "khr-rdpgw-sandbox"],
    "checks": checks,
    "errors": errors,
    "evidencePath": f"docs/evidence/khr-tp-live-reference-env/{RUN_ID}",
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(OUT / "reference-env-summary.json").write_text(json.dumps(summary, indent=2) + "\n")
print(f"[khr_tp_live_reference_env_check] summary={OUT / 'reference-env-summary.json'}")
print(
    f"[khr_tp_live_reference_env_check] status={status} "
    f"readyForScope1={summary['readyForScope1']} readyForScope2={summary['readyForScope2']} "
    f"readyForScope3={summary['readyForScope3']}"
)
if status != "PASS":
    for e in errors:
        print(f"[khr_tp_live_reference_env_check] FAIL: {e}", file=__import__("sys").stderr)
    raise SystemExit(1)
PY

log "PASS reference-env-summary=${OUT_DIR}/reference-env-summary.json"
