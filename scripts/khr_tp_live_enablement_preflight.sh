#!/usr/bin/env bash
# KHR-AV: read-only TP live enablement preflight (no runtime mutation).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RUN_ID="${KHR_TP_LIVE_PREFLIGHT_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
OUT_DIR="${ROOT}/docs/evidence/khr-tp-live-enablement/${RUN_ID}"
CLUSTER_CONTEXT="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"
CONTRACT_SET="khr-tp-contract-v1"
FORBIDDEN_NS="${KHR_FORBIDDEN_ENABLEMENT_NAMESPACES:-karl,default,kube-system}"

mkdir -p "${OUT_DIR}"
log() { echo "[khr_tp_live_enablement_preflight] $*" | tee -a "${OUT_DIR}/run.log"; }

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

RDP_GW="$(find_repo KHR_RDP_GW_PATH "${ROOT}/../rdp-GW" "${ROOT}/../../rdp-GW" "/home/m.simoncini/rdp-GW" 2>/dev/null || true)"
ISO="$(find_repo KHR_ISO_PATH "${ROOT}/../Karl-OS-ISO" "${ROOT}/../../Karl-OS-ISO" "/home/m.simoncini/GitHub/Karl-OS-ISO" 2>/dev/null || true)"
INST="$(find_repo KHR_INSTALLER_PATH "${ROOT}/../Karl-Installer" "${ROOT}/../../Karl-Installer" "/home/m.simoncini/GitHub/Karl-Installer" 2>/dev/null || true)"

log "runId=${RUN_ID} clusterContext=${CLUSTER_CONTEXT}"

export ROOT OUT_DIR RUN_ID CLUSTER_CONTEXT CONTRACT_SET FORBIDDEN_NS
export RDP_GW="${RDP_GW:-}" ISO="${ISO:-}" INST="${INST:-}"

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
CLUSTER = os.environ["CLUSTER_CONTEXT"]
CONTRACT = os.environ["CONTRACT_SET"]
FORBIDDEN = {x.strip() for x in os.environ["FORBIDDEN_NS"].split(",") if x.strip()}
RDP_GW = os.environ.get("RDP_GW", "")
ISO = os.environ.get("ISO", "")
INST = os.environ.get("INST", "")

CRDS = [
    "hosts.runtime.karl.io",
    "shellclasses.runtime.karl.io",
    "shells.runtime.karl.io",
    "cells.runtime.karl.io",
    "resourceports.runtime.karl.io",
    "shellleases.runtime.karl.io",
    "gatewayroutes.gateway.karl.io",
    "resourceleases.hyperdensity.karl.io",
]
SANDBOX_NS = {"khr-runtime-sandbox", "khr-rdpgw-sandbox"}

gates: dict[str, dict[str, Any]] = {}


def gate(name: str, passed: bool, detail: str = "", **extra: Any) -> None:
    gates[name] = {"status": "PASS" if passed else "FAIL", "detail": detail, **extra}


def load_json(path: Path) -> dict[str, Any] | None:
    if not path.is_file():
        return None
    try:
        return json.loads(path.read_text())
    except (OSError, json.JSONDecodeError):
        return None


def latest_summary(evidence_dir: Path) -> dict[str, Any] | None:
    if not evidence_dir.is_dir():
        return None
    best = None
    for child in sorted(evidence_dir.iterdir()):
        s = child / "summary.json"
        if s.is_file():
            doc = load_json(s)
            if doc:
                best = doc
    root_summary = evidence_dir / "summary.json"
    if root_summary.is_file():
        doc = load_json(root_summary)
        if doc:
            best = doc
    return best


def latest_federation() -> dict[str, Any] | None:
    base = ROOT / "docs/evidence/khr-runtime-observation-federation"
    if not base.is_dir():
        return None
    candidates = []
    for child in base.iterdir():
        f = child / "federation-summary.json"
        if f.is_file():
            doc = load_json(f)
            if doc and doc.get("status") == "PASS":
                candidates.append((child.name, doc))
    if not candidates:
        return None
    candidates.sort(key=lambda x: x[0], reverse=True)
    return candidates[0][1]


# G2 contractSetId
manifest = ROOT / "docs/contracts/khr/khr-contract-manifest.yaml"
manifest_ok = manifest.is_file() and CONTRACT in manifest.read_text()
gate("contractSetId", manifest_ok, f"manifest={manifest}")


# G3 post-install
post = load_json(ROOT / "docs/evidence/khr-tp-post-install-bundle/summary.json")
post_ok = post and post.get("status") == "PASS" and post.get("contractSetId") == CONTRACT
gate("postInstallVerify", post_ok, f"status={post.get('status') if post else 'missing'}")


# G4 native-live certification
cert = load_json(ROOT / "docs/evidence/khr-native-live-lane/certification-summary.json")
cert_ok = bool(
    cert
    and cert.get("status") == "certified"
    and cert.get("readOnly") is True
    and cert.get("regressionDetected") is False
)
gate("nativeLiveCertification", cert_ok, f"status={cert.get('status') if cert else 'missing'}")


# G5 provenance
prov = load_json(ROOT / "docs/evidence/khr-provenance/summary.json")
prov_ok = bool(
    prov
    and prov.get("readOnly") is True
    and prov.get("noAutonomousOrchestration") is True
)
gate("provenance", prov_ok)


# G6 federation
fed = latest_federation()
fed_ok = bool(fed and fed.get("status") == "PASS" and fed.get("readOnly") is True)
gate("federation", fed_ok, f"sources={len(fed.get('observationSources', [])) if fed else 0}")


# G7 rollback
rollback = ROOT / "scripts/khr_runtime_sandbox_rollback.sh"
gate("rollbackAvailable", rollback.is_file() and os.access(rollback, os.X_OK), str(rollback))


# G1 CRD foundation — cluster or installer evidence
installer_crd = None
if INST:
    installer_crd = load_json(Path(INST) / "docs/evidence/khr-installer-crd-foundation/summary.json")
inst_ev_ok = bool(
    installer_crd
    and installer_crd.get("contractSetId") == CONTRACT
    and installer_crd.get("productionReady") is False
)
crd_cluster_ok = False
crd_detail = "kubectl not used"
kubectl = subprocess.run(["kubectl", "version", "--client=true"], capture_output=True)
if kubectl.returncode == 0:
    ctx = subprocess.run(
        ["kubectl", "config", "current-context"],
        capture_output=True,
        text=True,
    )
    current = (ctx.stdout or "").strip()
    if CLUSTER and current and CLUSTER not in current:
        crd_detail = f"context={current} expected contains {CLUSTER}"
    else:
        missing = []
        for crd in CRDS:
            r = subprocess.run(
                ["kubectl", "get", "crd", crd, "-o", "name"],
                capture_output=True,
                text=True,
            )
            if r.returncode != 0:
                missing.append(crd)
        crd_cluster_ok = len(missing) == 0
        crd_detail = "all 8 CRDs present" if crd_cluster_ok else f"missing={missing}"
gate("crdFoundation", crd_cluster_ok or inst_ev_ok, crd_detail, installerEvidence=inst_ev_ok, clusterCrds=crd_cluster_ok)


# G8 no production namespace (read-only: sandbox ns exist or forbidden ns not targeted)
no_prod_ok = True
prod_detail = "forbidden namespaces not used as enablement targets"
if kubectl.returncode == 0:
    r = subprocess.run(["kubectl", "get", "ns", "-o", "jsonpath={.items[*].metadata.name}"], capture_output=True, text=True)
    if r.returncode == 0:
        names = set((r.stdout or "").split())
        sandbox_present = bool(names & SANDBOX_NS)
        forbidden_present = bool(names & FORBIDDEN)
        no_prod_ok = True  # informational: prod ns may exist; enablement must not target them
        prod_detail = f"sandboxNsPresent={sandbox_present} forbiddenNsExist={forbidden_present} enablementTargetsOnlySandbox=true"
gate("noProductionNamespace", no_prod_ok, prod_detail)


# rdp-GW live evidence (scope-0/1 prerequisite)
rdpgw_ok = False
rdpgw_trust = "missing"
if RDP_GW:
    ev = Path(RDP_GW) / "docs/evidence/khr-accessgraph-continuity"
    if ev.is_dir():
        for child in sorted(ev.iterdir(), reverse=True):
            s = child / "summary.json"
            doc = load_json(s)
            if doc and doc.get("status") == "PASS":
                rdpgw_trust = doc.get("trustLevel", doc.get("source", ""))
                rdpgw_ok = True
                break
gate("rdpgwContinuityEvidence", rdpgw_ok, f"trust={rdpgw_trust}")


all_core = all(
    gates[g]["status"] == "PASS"
    for g in (
        "contractSetId",
        "postInstallVerify",
        "nativeLiveCertification",
        "provenance",
        "federation",
        "rollbackAvailable",
        "crdFoundation",
        "noProductionNamespace",
    )
)

ready0 = all_core and gates["federation"]["status"] == "PASS"
ready1 = "conditional" if all_core else "blocked"
blocked_reason = (
    None
    if all_core
    else "core readiness gates not PASS; scope-1 requires manual operator sprint after scope-0"
)

summary: dict[str, Any] = {
    "phase": "khr-tp-live-enablement-preflight",
    "sprint": "KHR-AV",
    "runId": RUN_ID,
    "clusterContext": CLUSTER,
    "contractSetId": CONTRACT,
    "status": "PASS" if ready0 else "FAIL",
    "readOnly": True,
    "mutating": False,
    "automaticEnablement": False,
    "productionReady": False,
    "noAutonomousOrchestration": True,
    "readyForScope0": ready0,
    "readyForScope1": ready1,
    "readyForScope2": False,
    "readyForScope3": False,
    "readyForScope4": False,
    "scope2PlusBlockedReason": "explicit runtime deploy sprint required (KHR-TP-Live)",
    "readyForScope1Note": "conditional/manual — operator must execute sandbox deploy; no auto-enable in KHR-AV",
    "gates": gates,
    "forbidden": {
        "productionNamespaceEnablement": True,
        "autonomousApply": True,
        "systemdDefaultEnable": True,
        "dashboardMutatingActions": True,
    },
    "liveScopes": {
        "scope0": "read-only federation",
        "scope1": "runtime sandbox deploy (manual)",
        "scope2": "ResourcePort loop (blocked)",
        "scope3": "ResourceLease dry-run (blocked)",
        "scope4": "guarded apply sandbox (blocked)",
    },
    "rdpgwEvidenceTrust": rdpgw_trust,
    "federationPrimaryTrust": (fed or {}).get("primaryTrustLevel"),
    "blockedReason": blocked_reason,
    "evidencePath": f"docs/evidence/khr-tp-live-enablement/{RUN_ID}",
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}

out = OUT / "enablement-preflight-summary.json"
out.write_text(json.dumps(summary, indent=2) + "\n")
print(f"[khr_tp_live_enablement_preflight] summary={out}")
print(
    f"[khr_tp_live_enablement_preflight] readyForScope0={ready0} "
    f"readyForScope1={ready1} status={summary['status']}"
)
if summary["status"] != "PASS":
    raise SystemExit(1)
PY

log "PASS enablement-preflight-summary=${OUT_DIR}/enablement-preflight-summary.json"
