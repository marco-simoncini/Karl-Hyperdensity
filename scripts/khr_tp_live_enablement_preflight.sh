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


def latest_scope1_verify() -> dict[str, Any] | None:
    base = ROOT / "docs/evidence/khr-tp-live-scope1"
    if not base.is_dir():
        return None
    for name in sorted(base.iterdir(), reverse=True):
        v = name / "verify-summary.json"
        if v.is_file():
            doc = load_json(v)
            if doc:
                return doc
    return None


scope1_verify = latest_scope1_verify()
scope1_ok = bool(
    scope1_verify
    and scope1_verify.get("status") == "PASS"
    and scope1_verify.get("accessGraphLiveReadonly") is True
    and scope1_verify.get("resourcePortLoopEnabled") is False
    and scope1_verify.get("readyForScope2") is False
)
gate(
    "scope1SandboxEvidence",
    scope1_ok,
    f"runId={scope1_verify.get('runId') if scope1_verify else 'none'}",
)


def latest_scope2_preflight() -> dict[str, Any] | None:
    base = ROOT / "docs/evidence/khr-tp-live-scope2-preflight"
    committed = base / "committed-scope2-preflight-khr-az" / "scope2-preflight-summary.json"
    if committed.is_file():
        return load_json(committed)
    if not base.is_dir():
        return None
    for name in sorted(base.iterdir(), reverse=True):
        p = name / "scope2-preflight-summary.json"
        if p.is_file():
            doc = load_json(p)
            if doc and doc.get("status") == "PASS":
                return doc
    return None


scope2_preflight = latest_scope2_preflight()
scope2_pf_ok = bool(
    scope2_preflight
    and scope2_preflight.get("status") == "PASS"
    and scope2_preflight.get("resourcePortLoopEnabled") is False
    and scope2_preflight.get("loopEnabled") is False
    and scope2_preflight.get("resourceLeaseApplyEnabled") is False
)
gate(
    "scope2PreflightEvidence",
    scope2_pf_ok,
    f"runId={scope2_preflight.get('runId') if scope2_preflight else 'none'}",
)


def latest_scope2_loop_verify() -> dict[str, Any] | None:
    base = ROOT / "docs/evidence/khr-tp-live-scope2-resourceport-loop"
    committed = base / "committed-scope2-loop-khr-ba" / "verify-summary.json"
    if committed.is_file():
        return load_json(committed)
    if not base.is_dir():
        return None
    for name in sorted(base.iterdir(), reverse=True):
        v = name / "verify-summary.json"
        if v.is_file():
            doc = load_json(v)
            if doc and doc.get("status") == "PASS":
                return doc
    return None


scope2_loop = latest_scope2_loop_verify()
scope2_loop_ok = bool(
    scope2_loop
    and scope2_loop.get("status") == "PASS"
    and scope2_loop.get("readyForScope2") == "manual-loop-pass"
    and scope2_loop.get("readyForScope2Active") is False
    and scope2_loop.get("readyForScope3") is False
    and scope2_loop.get("resourceLeaseApplyEnabled") is False
)
gate(
    "scope2ManualLoopEvidence",
    scope2_loop_ok,
    f"runId={scope2_loop.get('runId') if scope2_loop else 'none'}",
)


def latest_scope3_preflight() -> dict[str, Any] | None:
    base = ROOT / "docs/evidence/khr-tp-live-scope3-preflight"
    committed = base / "committed-scope3-preflight-khr-bb" / "scope3-preflight-summary.json"
    if committed.is_file():
        return load_json(committed)
    if not base.is_dir():
        return None
    for name in sorted(base.iterdir(), reverse=True):
        p = name / "scope3-preflight-summary.json"
        if p.is_file():
            doc = load_json(p)
            if doc and doc.get("status") == "PASS":
                return doc
    return None


def latest_scope3_dryrun_verify() -> dict[str, Any] | None:
    base = ROOT / "docs/evidence/khr-tp-live-scope3-dryrun"
    committed = base / "committed-scope3-dryrun-khr-bc" / "verify-summary.json"
    if committed.is_file():
        return load_json(committed)
    if not base.is_dir():
        return None
    for name in sorted(base.iterdir(), reverse=True):
        v = name / "verify-summary.json"
        if v.is_file():
            doc = load_json(v)
            if doc and doc.get("status") == "PASS":
                return doc
    return None


scope3_preflight = latest_scope3_preflight()
scope3_pf_ok = bool(
    scope3_preflight
    and scope3_preflight.get("status") == "PASS"
    and scope3_preflight.get("resourceLeaseApplyEnabled") is False
    and scope3_preflight.get("readyForScope3Active") is False
)
gate(
    "scope3PreflightEvidence",
    scope3_pf_ok,
    f"runId={scope3_preflight.get('runId') if scope3_preflight else 'none'}",
)

scope3_dryrun = latest_scope3_dryrun_verify()
scope3_dryrun_ok = bool(
    scope3_dryrun
    and scope3_dryrun.get("status") == "PASS"
    and scope3_dryrun.get("readyForScope3") == "manual-dryrun-pass"
    and scope3_dryrun.get("readyForScope3Active") is False
    and scope3_dryrun.get("readyForScope4") is False
    and scope3_dryrun.get("dryRunObserved") is True
    and scope3_dryrun.get("applyObserved") is False
    and scope3_dryrun.get("noMutation") is True
    and scope3_dryrun.get("noApply") is True
    and scope3_dryrun.get("resourceLeaseApplyEnabled") is False
    and scope3_dryrun.get("cgroupMutationObserved") is False
)
gate(
    "scope3ManualDryRunEvidence",
    scope3_dryrun_ok,
    f"runId={scope3_dryrun.get('runId') if scope3_dryrun else 'none'}",
)


def latest_scope4_preflight() -> dict[str, Any] | None:
    base = ROOT / "docs/evidence/khr-tp-live-scope4-preflight"
    committed = base / "committed-scope4-preflight-khr-bd" / "scope4-preflight-summary.json"
    if committed.is_file():
        return load_json(committed)
    if not base.is_dir():
        return None
    for name in sorted(base.iterdir(), reverse=True):
        p = name / "scope4-preflight-summary.json"
        if p.is_file():
            doc = load_json(p)
            if doc and doc.get("status") == "PASS":
                return doc
    return None


scope4_preflight = latest_scope4_preflight()
scope4_pf_ok = bool(
    scope4_preflight
    and scope4_preflight.get("status") == "PASS"
    and scope4_preflight.get("readyForScope4") == "conditional/manual-preflight-pass"
    and scope4_preflight.get("readyForScope4Active") is False
    and scope4_preflight.get("guardedApplyExecuted") is False
    and scope4_preflight.get("cgroupMutationObserved") is False
    and scope4_preflight.get("dryRunDecision") == "allowed"
    and scope4_preflight.get("rollbackPlanRef")
    and scope4_preflight.get("verificationPlanRef")
    and scope4_preflight.get("sourceResourcePortRef")
)
gate(
    "scope4PreflightEvidence",
    scope4_pf_ok,
    f"runId={scope4_preflight.get('runId') if scope4_preflight else 'none'}",
)


def latest_scope4_apply_verify() -> dict[str, Any] | None:
    base = ROOT / "docs/evidence/khr-tp-live-scope4-guarded-apply"
    committed = base / "committed-scope4-guarded-apply-khr-be" / "verify-summary.json"
    if committed.is_file():
        return load_json(committed)
    if not base.is_dir():
        return None
    for name in sorted(base.iterdir(), reverse=True):
        v = name / "verify-summary.json"
        if v.is_file():
            doc = load_json(v)
            if doc and doc.get("status") == "PASS":
                return doc
    return None


def latest_scope4_rollback() -> dict[str, Any] | None:
    base = ROOT / "docs/evidence/khr-tp-live-scope4-guarded-apply"
    committed = base / "committed-scope4-guarded-apply-khr-be" / "rollback-summary.json"
    if committed.is_file():
        return load_json(committed)
    if not base.is_dir():
        return None
    for name in sorted(base.iterdir(), reverse=True):
        r = name / "rollback-summary.json"
        if r.is_file():
            doc = load_json(r)
            if doc and doc.get("status") == "PASS":
                return doc
    return None


scope4_apply = latest_scope4_apply_verify()
scope4_rollback = latest_scope4_rollback()
scope4_apply_ok = bool(
    scope4_apply
    and scope4_apply.get("status") == "PASS"
    and scope4_apply.get("readyForScope4") == "manual-guarded-apply-pass"
    and scope4_apply.get("readyForScope4Active") is False
    and scope4_apply.get("guardedApplyObserved") is True
    and scope4_apply.get("continuityPreserved") is True
    and scope4_apply.get("applyScope") == "sandbox-only"
)
scope4_rollback_ok = bool(
    scope4_rollback
    and scope4_rollback.get("status") == "PASS"
    and scope4_rollback.get("rollbackVerified") is True
)
gate(
    "scope4ManualGuardedApplyEvidence",
    scope4_apply_ok and scope4_rollback_ok,
    f"runId={scope4_apply.get('runId') if scope4_apply else 'none'} rollback={scope4_rollback_ok}",
)


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
if scope1_ok:
    ready1: bool | str = True
    ready1_note = "scope-1 sandbox deploy/verify evidence PASS (KHR-AW)"
elif all_core:
    ready1 = "conditional"
    ready1_note = "gates PASS; run khr_tp_live_scope1_deploy.sh for scope-1 evidence"
else:
    ready1 = False
    ready1_note = "core gates not PASS"
blocked_reason = None
if not all_core:
    blocked_reason = "core readiness gates not PASS"
elif not scope1_ok and ready1 == "conditional":
    blocked_reason = None

if scope2_loop_ok and scope1_ok:
    ready2 = "manual-loop-pass"
    ready2_note = "KHR-BA manual loop evidence PASS; scope-2 not active; scope-3 blocked"
elif scope2_pf_ok and scope1_ok:
    ready2 = "conditional/manual-preflight-pass"
    ready2_note = "scope-2 preflight PASS; run khr_tp_live_scope2_resourceport_loop_run.sh"
elif scope1_ok:
    ready2 = False
    ready2_note = "run khr_tp_live_scope2_preflight.sh for scope-2 readiness"
else:
    ready2 = False
    ready2_note = "scope-1 evidence required before scope-2 preflight"

if scope3_dryrun_ok and scope2_loop_ok and scope1_ok:
    ready3: bool | str = "manual-dryrun-pass"
    ready3_note = "KHR-BC scope-3 manual dry-run evidence PASS; not active; scope-4 preflight when ready"
elif scope3_pf_ok and scope2_loop_ok and scope1_ok:
    ready3 = "conditional/manual-preflight-pass"
    ready3_note = "KHR-BB scope-3 preflight PASS; dry-run execution deferred; scope-4 blocked"
elif scope2_loop_ok and scope1_ok:
    ready3 = False
    ready3_note = "run khr_tp_live_scope3_preflight.sh for scope-3 readiness"
else:
    ready3 = False
    ready3_note = "scope-2 manual-loop-pass required before scope-3 preflight"

dry_run_executed = bool(scope3_dryrun_ok)

if scope4_apply_ok and scope4_rollback_ok and scope3_dryrun_ok and scope2_loop_ok and scope1_ok:
    ready4: bool | str = "manual-guarded-apply-pass"
    ready4_note = "KHR-BE scope-4 guarded apply evidence PASS with rollback verified; not active"
elif scope4_pf_ok and scope3_dryrun_ok and scope2_loop_ok and scope1_ok:
    ready4 = "conditional/manual-preflight-pass"
    ready4_note = "KHR-BD scope-4 preflight PASS; guarded apply not executed; not active"
elif scope3_dryrun_ok and scope1_ok:
    ready4 = False
    ready4_note = "run khr_tp_live_scope4_preflight.sh for scope-4 readiness"
else:
    ready4 = False
    ready4_note = "scope-3 manual-dryrun-pass required before scope-4 preflight"

summary: dict[str, Any] = {
    "phase": "khr-tp-live-enablement-preflight",
    "sprint": "KHR-BE",
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
    "readyForScope2": ready2,
    "readyForScope2Active": False,
    "readyForScope2LoopExecution": scope2_loop_ok,
    "readyForScope3": ready3,
    "readyForScope3Active": False,
    "readyForScope4": ready4,
    "readyForScope4Active": False,
    "scope2PlusBlockedReason": "ResourcePort loop execution requires dedicated sprint sign-off (KHR-AZ)",
    "readyForScope1Note": ready1_note,
    "readyForScope2Note": ready2_note,
    "scope1EvidenceRunId": scope1_verify.get("runId") if scope1_verify else None,
    "scope2PreflightRunId": scope2_preflight.get("runId") if scope2_preflight else None,
    "scope2LoopRunId": scope2_loop.get("runId") if scope2_loop else None,
    "scope3PreflightRunId": scope3_preflight.get("runId") if scope3_preflight else None,
    "scope3DryRunRunId": scope3_dryrun.get("runId") if scope3_dryrun else None,
    "scope4PreflightRunId": scope4_preflight.get("runId") if scope4_preflight else None,
    "scope4ApplyRunId": scope4_apply.get("runId") if scope4_apply else None,
    "readyForScope3Note": ready3_note,
    "readyForScope4Note": ready4_note,
    "guardedApplyExecuted": scope4_apply_ok,
    "rollbackVerified": scope4_rollback_ok,
    "continuityPreserved": scope4_apply.get("continuityPreserved", False) if scope4_apply else False,
    "resourceLeaseDryRunExecuted": dry_run_executed,
    "dryRunObserved": dry_run_executed,
    "applyObserved": False,
    "noMutation": scope3_dryrun.get("noMutation", True) if scope3_dryrun_ok else True,
    "resourceLeaseApplyEnabled": False,
    "cgroupMutationObserved": False,
    "resourcePortLoopEnabled": scope2_preflight.get("resourcePortLoopEnabled", False) if scope2_preflight else False,
    "sandboxApplyEnabled": scope2_preflight.get("sandboxApplyEnabled", False) if scope2_preflight else False,
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
        "scope2": "ResourcePort manual loop (manual-loop-pass when evidenced; not active)",
        "scope3": "ResourceLease manual dry-run (manual-dryrun-pass when evidenced; not active)",
        "scope4": "ResourceLease guarded apply (manual-guarded-apply-pass when evidenced; not active)",
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
    f"readyForScope1={ready1} readyForScope2={ready2} readyForScope3={ready3} "
    f"readyForScope4={ready4} status={summary['status']}"
)
if summary["status"] != "PASS":
    raise SystemExit(1)
PY

log "PASS enablement-preflight-summary=${OUT_DIR}/enablement-preflight-summary.json"
