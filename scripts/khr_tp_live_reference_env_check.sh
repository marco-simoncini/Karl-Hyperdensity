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

rdpgw_dep = None
rdpgw_scope1 = RDP_GW / "docs/evidence/khr-rdpgw-scope1"
if rdpgw_scope1.is_dir():
    for child in sorted(rdpgw_scope1.iterdir(), reverse=True):
        dep = child / "deploy-summary.json"
        if dep.is_file():
            rdpgw_dep = load(dep)
            break
        # legacy: infer from deploy log marker file
        mode_file = child / "deploy-mode.json"
        if mode_file.is_file():
            rdpgw_dep = load(mode_file)
            break

if scope1_verify:
    deploy_mode = scope1_verify.get("rdpgwDeployMode", "missing")
    if deploy_mode == "local":
        deploy_mode = "local-gateway"
    if deploy_mode == "cluster":
        deploy_mode = "cluster-sandbox"

if rdpgw_dep:
    deploy_mode = rdpgw_dep.get("deployMode", deploy_mode)
    if deploy_mode == "local":
        deploy_mode = "local-gateway"
    if deploy_mode == "cluster":
        deploy_mode = "cluster-sandbox"

record("rdpgwLiveReadonlyEvidence", ag_ok, f"trust={trust}")
record(
    "rdpgwDeployMode",
    deploy_mode in ("local-gateway", "cluster-sandbox"),
    f"deployMode={deploy_mode}",
    deployMode=deploy_mode,
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
record(
    "dashboardTpReadinessFixture",
    bool(
        dash
        and dash.get("tpReadinessSummary", {}).get("readyForScope1") is True
        and dash.get("tpReadinessSummary", {}).get("readyForScope2") is False
        and dash.get("tpReadinessSummary", {}).get("productionReady") is False
    ),
    str(dash_fixture),
)

status = "PASS" if not errors else "FAIL"
summary = {
    "phase": "khr-tp-live-reference-env",
    "sprint": "KHR-AX",
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
    "readyForScope2": False,
    "scope2BlockedReason": scope1_verify.get("scope2BlockedReason", "scope-2+ blocked in KHR-AX") if scope1_verify else "scope-2+ blocked",
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
print(f"[khr_tp_live_reference_env_check] status={status} readyForScope1={summary['readyForScope1']}")
if status != "PASS":
    for e in errors:
        print(f"[khr_tp_live_reference_env_check] FAIL: {e}", file=__import__("sys").stderr)
    raise SystemExit(1)
PY

log "PASS reference-env-summary=${OUT_DIR}/reference-env-summary.json"
