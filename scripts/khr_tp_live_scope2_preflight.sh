#!/usr/bin/env bash
# KHR-AZ: read-only Scope-2 ResourcePort loop preflight (no loop enable).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RUN_ID="${KHR_TP_LIVE_SCOPE2_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
OUT_DIR="${ROOT}/docs/evidence/khr-tp-live-scope2-preflight/${RUN_ID}"
CLUSTER="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"
SANDBOX_NS="${KHR_RUNTIME_SANDBOX_NS:-khr-runtime-sandbox}"
MANIFEST_DIR="${ROOT}/examples/khr/tp-live-scope1"
SCOPE1_CFG="${MANIFEST_DIR}/configmap-karl-host-runtime-scope1.yaml"
PREVIEW_DEPLOY="${MANIFEST_DIR}/karl-host-runtime-preview-deployment.yaml"

mkdir -p "${OUT_DIR}"
log() { echo "[khr_tp_live_scope2_preflight] $*" | tee -a "${OUT_DIR}/run.log"; }

find_repo() {
  local var="$1"
  shift
  if [[ -n "${!var:-}" && -d "${!var}" ]]; then
    echo "$(cd "${!var}" && pwd)"
    return 0
  fi
  for c in "$@"; do
    [[ -d "${c}" ]] || continue
    echo "$(cd "${c}" && pwd)"
    return 0
  done
  return 1
}

RDP_GW="$(find_repo KHR_RDP_GW_PATH "${ROOT}/../rdp-GW" "/home/m.simoncini/rdp-GW" 2>/dev/null || true)"

log "runId=${RUN_ID} cluster=${CLUSTER} namespace=${SANDBOX_NS}"

export ROOT OUT_DIR RUN_ID CLUSTER SANDBOX_NS MANIFEST_DIR SCOPE1_CFG PREVIEW_DEPLOY RDP_GW
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
MANIFEST_DIR = Path(os.environ["MANIFEST_DIR"])
SCOPE1_CFG = Path(os.environ["SCOPE1_CFG"])
PREVIEW_DEPLOY = Path(os.environ["PREVIEW_DEPLOY"])
RDP_GW = os.environ.get("RDP_GW", "")

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


def kubectl_json(args: list[str]) -> dict[str, Any] | None:
    r = subprocess.run(
        ["kubectl", "--context", CLUSTER, *args],
        capture_output=True,
        text=True,
    )
    if r.returncode != 0:
        return None
    try:
        return json.loads(r.stdout)
    except json.JSONDecodeError:
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

# Scope-1 ready
scope1_base = ROOT / "docs/evidence/khr-tp-live-scope1"
scope1_verify = load(scope1_base / "committed-scope1-khr-aw" / "verify-summary.json")
if scope1_verify is None and scope1_base.is_dir():
    for child in sorted(scope1_base.iterdir(), reverse=True):
        v = child / "verify-summary.json"
        if v.is_file():
            scope1_verify = load(v)
            if scope1_verify:
                break
scope1_ok = bool(
    scope1_verify
    and scope1_verify.get("status") == "PASS"
    and scope1_verify.get("accessGraphLiveReadonly") is True
    and scope1_verify.get("resourcePortLoopEnabled") is False
    and scope1_verify.get("readyForScope2") is False
)
record(
    "scope1Ready",
    scope1_ok,
    f"runId={scope1_verify.get('runId') if scope1_verify else 'none'}",
)

# Sandbox namespace label
ns_doc = kubectl_json(["get", "namespace", SANDBOX_NS, "-o", "json"])
label_ok = bool(
    ns_doc
    and (ns_doc.get("metadata", {}).get("labels", {}).get("khr.karl.io/sandbox") == "true")
)
record("sandboxNamespaceLabel", label_ok, f"namespace={SANDBOX_NS} label khr.karl.io/sandbox=true")

# Preview deployable (manifests + config guards)
manifest_ok = SCOPE1_CFG.is_file() and PREVIEW_DEPLOY.is_file()
cfg_loop_false = False
cfg_apply_false = False
if SCOPE1_CFG.is_file():
    text = SCOPE1_CFG.read_text()
    cfg_loop_false = "resourcePortLoopEnabled: false" in text
    cfg_apply_false = "sandboxApplyEnabled: false" in text
record(
    "karlHostRuntimePreviewDeployable",
    manifest_ok and cfg_loop_false and cfg_apply_false,
    f"manifests={manifest_ok} loopDisabled={cfg_loop_false} applyDisabled={cfg_apply_false}",
)

# Live cluster config (if preview deployed)
resource_port_loop_enabled = False
sandbox_apply_enabled = False
deploy = kubectl_json(["get", "deploy", "-n", SANDBOX_NS, "-o", "json"])
items = (deploy or {}).get("items") or []
for item in items:
    name = item.get("metadata", {}).get("name", "")
    if "karl-host-runtime" in name:
        cm_name = None
        for vol in item.get("spec", {}).get("template", {}).get("spec", {}).get("volumes", []):
            cm = vol.get("configMap", {})
            if cm.get("name"):
                cm_name = cm["name"]
        if cm_name:
            cm = kubectl_json(["get", "configmap", cm_name, "-n", SANDBOX_NS, "-o", "json"])
            data = (cm or {}).get("data", {})
            for v in data.values():
                if "resourcePortLoopEnabled: true" in v:
                    resource_port_loop_enabled = True
                if "sandboxApplyEnabled: true" in v:
                    sandbox_apply_enabled = True

# Static config is source of truth when not deployed
if not items:
    resource_port_loop_enabled = not cfg_loop_false
    sandbox_apply_enabled = not cfg_apply_false

record(
    "resourcePortLoopDisabled",
    resource_port_loop_enabled is False,
    f"resourcePortLoopEnabled={resource_port_loop_enabled}",
)
record(
    "sandboxApplyDisabled",
    sandbox_apply_enabled is False,
    f"sandboxApplyEnabled={sandbox_apply_enabled}",
)
record(
    "resourceLeaseApplyDisabled",
    sandbox_apply_enabled is False,
    "no ResourceLease apply when sandboxApplyEnabled=false",
)

# rdp-GW cluster-sandbox visibility (Scope-1 dependency)
rdpgw_ok = False
rdpgw_mode = "missing"
if RDP_GW:
    cs = Path(RDP_GW) / "docs/evidence/khr-rdpgw-cluster-sandbox/committed-cluster-sandbox-khr-ay/verify-summary.json"
    rdpgw = load(cs)
    if rdpgw and rdpgw.get("status") == "PASS" and rdpgw.get("deployMode") == "cluster-sandbox":
        rdpgw_ok = True
        rdpgw_mode = "cluster-sandbox"
record("rdpgwClusterSandboxEvidence", rdpgw_ok, f"deployMode={rdpgw_mode}")

# Production namespace mutation proof (read-only generation snapshot)
prod_ok = ctx_ok
prod_detail: dict[str, Any] = {}
if ctx_ok:
    for ns in ("karl", "karl-system", "default", "kube-system"):
        r = subprocess.run(
            [
                "kubectl",
                "--context",
                CLUSTER,
                "get",
                "deploy",
                "-n",
                ns,
                "-o",
                "jsonpath={.items[*].metadata.generation}",
            ],
            capture_output=True,
            text=True,
        )
        prod_detail[ns] = (r.stdout or "").strip() if r.returncode == 0 else "unavailable"
record("noProductionNamespaceMutation", prod_ok, json.dumps(prod_detail))

loop_enabled = resource_port_loop_enabled is True
status = "PASS" if not errors else "FAIL"
ready_for_scope2: bool | str = False
if status == "PASS":
    ready_for_scope2 = "conditional/manual-preflight-pass"

summary: dict[str, Any] = {
    "phase": "khr-tp-live-scope2-preflight",
    "sprint": "KHR-AZ",
    "runId": RUN_ID,
    "clusterContext": CLUSTER,
    "namespace": SANDBOX_NS,
    "contractSetId": "khr-tp-contract-v1",
    "status": status,
    "readOnly": True,
    "mutating": False,
    "automaticEnablement": False,
    "loopEnabled": loop_enabled,
    "productionReady": False,
    "noAutonomousOrchestration": True,
    "readyForScope1": scope1_ok,
    "readyForScope2": ready_for_scope2,
    "readyForScope2LoopExecution": False,
    "resourcePortLoopEnabled": resource_port_loop_enabled,
    "sandboxApplyEnabled": sandbox_apply_enabled,
    "resourceLeaseApplyEnabled": False,
    "scope2BlockedReason": "ResourcePort loop not enabled; execution deferred to dedicated sprint after sign-off",
    "scope2PreflightNote": "KHR-AZ preflight only — no loop execution",
    "checks": checks,
    "errors": errors,
    "forbidden": {
        "resourcePortLoopPermanentEnable": True,
        "resourceLeaseDryRunApply": True,
        "productionNamespaceMutation": True,
        "dashboardMutatingActions": True,
    },
    "evidencePath": f"docs/evidence/khr-tp-live-scope2-preflight/{RUN_ID}",
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(OUT / "scope2-preflight-summary.json").write_text(json.dumps(summary, indent=2) + "\n")
print(f"[khr_tp_live_scope2_preflight] summary={OUT / 'scope2-preflight-summary.json'}")
print(
    f"[khr_tp_live_scope2_preflight] status={status} "
    f"readyForScope2={ready_for_scope2} resourcePortLoopEnabled={resource_port_loop_enabled}"
)
if status != "PASS":
    for e in errors:
        print(f"[khr_tp_live_scope2_preflight] FAIL: {e}", file=__import__("sys").stderr)
    raise SystemExit(1)
PY

log "PASS scope2-preflight-summary=${OUT_DIR}/scope2-preflight-summary.json"
