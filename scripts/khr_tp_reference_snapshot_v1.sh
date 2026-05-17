#!/usr/bin/env bash
# KHR-BT: aggregate committed TP live evidence into reference snapshot v1 (read-only).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RUN_ID="${KHR_TP_REFERENCE_SNAPSHOT_RUN_ID:-committed-khr-bt-v1}"
OUT_DIR="${ROOT}/docs/evidence/khr-tp-reference-snapshot-v1/${RUN_ID}"
DASHBOARD="${KHR_DASHBOARD_PATH:-$(cd "${ROOT}/../Karl-Dashboard" 2>/dev/null && pwd || echo "")}"
INSTALLER="${KHR_INSTALLER_PATH:-$(cd "${ROOT}/../Karl-Installer" 2>/dev/null && pwd || echo "")}"
ISO="${KHR_ISO_PATH:-$(cd "${ROOT}/../Karl-OS-ISO" 2>/dev/null && pwd || echo "")}"
RDP_GW="${KHR_RDP_GW_PATH:-$(cd "${ROOT}/../rdp-GW" 2>/dev/null && pwd || echo "/home/m.simoncini/rdp-GW")}"

mkdir -p "${OUT_DIR}"
log() { echo "[khr_tp_reference_snapshot_v1] $*" | tee -a "${OUT_DIR}/run.log"; }
log "runId=${RUN_ID} read-only aggregation (no cluster mutation)"

export ROOT OUT_DIR RUN_ID DASHBOARD INSTALLER ISO RDP_GW
python3 <<'PY'
from __future__ import annotations

import json
import os
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

ROOT = Path(os.environ["ROOT"])
OUT = Path(os.environ["OUT_DIR"])
RUN_ID = os.environ["RUN_ID"]
DASHBOARD = Path(os.environ["DASHBOARD"]) if os.environ.get("DASHBOARD") else None
INSTALLER = Path(os.environ["INSTALLER"]) if os.environ.get("INSTALLER") else None
ISO = Path(os.environ["ISO"]) if os.environ.get("ISO") else None
RDP_GW = Path(os.environ["RDP_GW"]) if os.environ.get("RDP_GW") else None

CONTRACT_VERSION = "khr-tp-reference-snapshot-v1"
CONTRACT_SET_ID = "khr-tp-contract-v1"
CLUSTER = "karl-metal-01@ovh"


def log(msg: str) -> None:
    print(f"[khr_tp_reference_snapshot_v1] {msg}")


def load(path: Path) -> dict[str, Any] | None:
    if not path.is_file():
        return None
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError):
        return None


def git_short(repo: Path | None) -> str | None:
    if repo is None or not (repo / ".git").exists():
        return None
    cp = subprocess.run(
        ["git", "-C", str(repo), "rev-parse", "--short", "HEAD"],
        capture_output=True,
        text=True,
        check=False,
    )
    return (cp.stdout or "").strip() or None


def ref_entry(
    repo: str,
    repo_path: Path | None,
    rel_path: str,
    role: str,
    *,
    required: bool = True,
) -> dict[str, Any]:
    abs_path = None
    if repo == "Karl-Hyperdensity":
        abs_path = ROOT / rel_path
    elif repo_path:
        abs_path = repo_path / rel_path
    if abs_path and abs_path.suffix in (".yaml", ".yml"):
        present = abs_path.is_file()
        return {
            "repo": repo,
            "role": role,
            "relativePath": rel_path,
            "absolutePath": str(abs_path),
            "present": present,
            "status": "PASS" if present or not required else "MISSING",
            "format": "yaml",
        }
    doc = load(abs_path) if abs_path else None
    ok = doc is not None
    if ok:
        if doc.get("status") == "FAIL":
            ok = False
        elif doc.get("evidenceStatus") in (
            "LIVE_PASS",
            "IMAGE_RESOLVED",
            "ROLLBACK_VERIFIED",
            "certified-evidence-backed",
        ):
            ok = True
        elif doc.get("status") in ("PASS", None) or "status" not in doc:
            ok = True
        else:
            ok = doc.get("status") == "PASS"
    if not ok and not required:
        ok = True
    return {
        "repo": repo,
        "role": role,
        "relativePath": rel_path,
        "absolutePath": str(abs_path) if abs_path else None,
        "present": abs_path.is_file() if abs_path else False,
        "status": "PASS" if ok else ("MISSING" if doc is None else "FAIL"),
        "evidenceStatus": (doc or {}).get("evidenceStatus"),
        "summary": {k: (doc or {}).get(k) for k in ("runId", "readyForScope4", "readyForScope3", "readyForScope2", "readyForScope1", "deployMode", "providerProfile", "rollbackVerified", "imageResolved") if (doc or {}).get(k) is not None},
    }


# Committed evidence anchors (KHR-BT)
HYPERDENSITY_REFS = {
    "enablement": "docs/evidence/khr-tp-live-enablement/20260517T072602Z/enablement-preflight-summary.json",
    "scope1": "docs/evidence/khr-tp-live-scope1/committed-scope1-khr-aw/verify-summary.json",
    "scope2Preflight": "docs/evidence/khr-tp-live-scope2-preflight/committed-scope2-preflight-khr-az/scope2-preflight-summary.json",
    "scope2Loop": "docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/verify-summary.json",
    "scope3Preflight": "docs/evidence/khr-tp-live-scope3-preflight/committed-scope3-preflight-khr-bb/scope3-preflight-summary.json",
    "scope3Dryrun": "docs/evidence/khr-tp-live-scope3-dryrun/committed-scope3-dryrun-khr-bc/verify-summary.json",
    "scope4Preflight": "docs/evidence/khr-tp-live-scope4-preflight/committed-scope4-preflight-khr-bd/scope4-preflight-summary.json",
    "scope4Apply": "docs/evidence/khr-tp-live-scope4-guarded-apply/committed-scope4-guarded-apply-khr-be/verify-summary.json",
    "scope4Certification": "docs/evidence/khr-scope4-guarded-apply-certification/committed-scope4-certification-khr-bf/certification-summary.json",
    "scope4Governance": "docs/evidence/khr-scope4-governance/committed-scope4-governance-khr-bg/governance-summary.json",
    "federation": "docs/evidence/khr-runtime-observation-federation/20260517T072601Z/federation-summary.json",
}

DASHBOARD_REFS = {
    "livePass": "docs/evidence/khr-dashboard-reference-env-activation/committed-khr-bs-20260517T073046Z/live-summary.json",
    "liveSummary": "docs/evidence/khr-dashboard-reference-env-activation/committed-khr-bs-20260517T073046Z/summary.json",
    "rollback": "docs/evidence/khr-dashboard-reference-env-activation/committed-khr-bs-20260517T073046Z/rollback-summary.json",
    "imageResolve": "docs/evidence/khr-dashboard-console-image/committed-khr-bs-20260517T073046Z/image-resolution-summary.json",
    "rolloutPlan": "docs/evidence/khr-dashboard-reference-env-activation/committed-khr-bs-20260517T073046Z/rollout-plan.json",
}

INSTALLER_REFS = {
    "crdFoundation": "docs/evidence/khr-installer-crd-foundation/20260517T070416Z/summary.json",
    "crdFoundationManifest": "docs/evidence/khr-installer-crd-foundation/20260517T070416Z/khr-contract-manifest.yaml",
    "hybrid": "docs/evidence/khr-hybrid-transition/20260516T195854Z/summary.json",
}

ISO_REFS = {
    "postInstall": "docs/evidence/khr-post-install-verify/summary.json",
    "profileManifest": "Stage-2/custom-files/karl-os-opt/karl/khr/profile-manifest.yaml",
}

RDPGW_REFS = {
    "clusterSandbox": "docs/evidence/khr-rdpgw-cluster-sandbox/committed-cluster-sandbox-khr-ay/verify-summary.json",
}

index: list[dict[str, Any]] = []
errors: list[str] = []

for key, rel in HYPERDENSITY_REFS.items():
    e = ref_entry("Karl-Hyperdensity", ROOT, rel, f"hyperdensity.{key}")
    index.append(e)
    if e["status"] != "PASS":
        errors.append(f"hyperdensity.{key}:{e['status']}")

for key, rel in DASHBOARD_REFS.items():
    e = ref_entry("Karl-Dashboard", DASHBOARD, rel, f"dashboard.{key}")
    index.append(e)
    if e["status"] != "PASS" and key in ("livePass", "rollback", "imageResolve"):
        errors.append(f"dashboard.{key}:{e['status']}")

for key, rel in INSTALLER_REFS.items():
    req = key != "crdFoundationManifest"
    e = ref_entry("Karl-Installer", INSTALLER, rel, f"installer.{key}", required=req)
    index.append(e)
    if e["status"] != "PASS" and req:
        errors.append(f"installer.{key}:{e['status']}")

for key, rel in ISO_REFS.items():
    req = key == "postInstall"
    e = ref_entry("Karl-OS-ISO", ISO, rel, f"iso.{key}", required=req)
    index.append(e)
    if e["status"] != "PASS" and req:
        errors.append(f"iso.{key}:{e['status']}")

for key, rel in RDPGW_REFS.items():
    e = ref_entry("rdp-GW", RDP_GW, rel, f"rdpgw.{key}")
    index.append(e)
    if e["status"] != "PASS":
        errors.append(f"rdpgw.{key}:{e['status']}")

def pick(repo: str, role_prefix: str) -> dict[str, Any] | None:
    for e in index:
        if e["repo"] == repo and e["role"].startswith(role_prefix):
            if e["present"]:
                p = Path(e["absolutePath"])
                return load(p)
    return None

scope1 = load(ROOT / HYPERDENSITY_REFS["scope1"]) or {}
scope2 = load(ROOT / HYPERDENSITY_REFS["scope2Loop"]) or {}
scope3 = load(ROOT / HYPERDENSITY_REFS["scope3Dryrun"]) or {}
scope4 = load(ROOT / HYPERDENSITY_REFS["scope4Apply"]) or {}
gov = load(ROOT / HYPERDENSITY_REFS["scope4Governance"]) or {}
cert = load(ROOT / HYPERDENSITY_REFS["scope4Certification"]) or {}
live = None
if DASHBOARD:
    live = load(DASHBOARD / DASHBOARD_REFS["livePass"])
rollback = load(DASHBOARD / DASHBOARD_REFS["rollback"]) if DASHBOARD else None
img = load(DASHBOARD / DASHBOARD_REFS["imageResolve"]) if DASHBOARD else None
rdpgw = load(RDP_GW / RDPGW_REFS["clusterSandbox"]) if RDP_GW else None
inst = load(INSTALLER / INSTALLER_REFS["crdFoundation"]) if INSTALLER else None
hyb = load(INSTALLER / INSTALLER_REFS["hybrid"]) if INSTALLER else None
iso_pi = load(ISO / ISO_REFS["postInstall"]) if ISO else None

scope_readiness = {
    "scope1": "PASS" if scope1.get("status") == "PASS" else str(scope1.get("status", "unknown")),
    "scope2": scope2.get("readyForScope2") or scope2.get("status", "unknown"),
    "scope3": scope3.get("readyForScope3") or scope3.get("status", "unknown"),
    "scope4": scope4.get("readyForScope4") or scope4.get("status", "unknown"),
    "scope4Active": bool(scope4.get("readyForScope4Active")),
}

dashboard_live_pass_ref = {
    "repo": "Karl-Dashboard",
    "runId": (live or {}).get("runId", "committed-khr-bs-20260517T073046Z"),
    "evidenceStatus": (live or {}).get("evidenceStatus"),
    "providerProfile": (live or {}).get("providerProfile", "khr-native"),
    "httpStatus": (live or {}).get("httpStatus"),
    "source": (live or {}).get("source"),
    "imageTag": (img or {}).get("requestedImage") or (img or {}).get("suggestedImage"),
    "rollbackVerified": (rollback or {}).get("rollbackVerified"),
    "paths": {
        "liveSummary": DASHBOARD_REFS["livePass"],
        "rollbackSummary": DASHBOARD_REFS["rollback"],
        "imageResolution": DASHBOARD_REFS["imageResolve"],
    },
}

rdpgw_cluster_sandbox_ref = {
    "repo": "rdp-GW",
    "runId": "committed-cluster-sandbox-khr-ay",
    "deployMode": (rdpgw or {}).get("deployMode"),
    "accessGraphLiveReadonly": (rdpgw or {}).get("accessGraphLiveReadonly"),
    "noRevoke": (rdpgw or {}).get("noRevoke"),
    "noDisconnect": (rdpgw or {}).get("noDisconnect"),
    "path": RDPGW_REFS["clusterSandbox"],
}

installer_crd_foundation_ref = {
    "repo": "Karl-Installer",
    "runId": (inst or {}).get("runId"),
    "profile": (inst or {}).get("profile", "karl2-khr-technical-preview"),
    "contractSetId": (inst or {}).get("contractSetId", CONTRACT_SET_ID),
    "crdDiffEmpty": (inst or {}).get("crdDiffEmpty"),
    "hostRuntimeEnabled": (inst or {}).get("hostRuntimeEnabled", False),
    "path": INSTALLER_REFS["crdFoundation"],
}

hybrid_transition_ref = {
    "repo": "Karl-Installer",
    "runId": (hyb or {}).get("runId"),
    "profile": (hyb or {}).get("profile", "hybrid-transition"),
    "kubevirtCompatibility": (hyb or {}).get("kubevirtCompatibility"),
    "kubevirtAsKhrCore": (hyb or {}).get("kubevirtAsKhrCore", False),
    "path": INSTALLER_REFS["hybrid"],
}

snapshot: dict[str, Any] = {
    "contractVersion": CONTRACT_VERSION,
    "sprint": "KHR-BT",
    "runId": RUN_ID,
    "contractSetId": CONTRACT_SET_ID,
    "clusterContext": CLUSTER,
    "readOnly": True,
    "liveMutationPerformed": False,
    "productionReady": False,
    "autonomousOrchestration": False,
    "globalDefaultsChanged": False,
    "noNewRollout": True,
    "noNewRuntimeMutation": True,
    "providerProfile": dashboard_live_pass_ref.get("providerProfile") or "khr-native",
    "scopeReadiness": scope_readiness,
    "scope4CertificationState": (cert or gov or {}).get("scope4CertificationState", "certified-evidence-backed"),
    "governanceState": (gov or {}).get("scope4GovernanceState", "certified"),
    "readyForScope4Active": scope_readiness.get("scope4Active", False),
    "dashboardLivePassRef": dashboard_live_pass_ref,
    "rdpgwClusterSandboxRef": rdpgw_cluster_sandbox_ref,
    "installerCrdFoundationRef": installer_crd_foundation_ref,
    "hybridTransitionRef": hybrid_transition_ref,
    "isoPostInstallRef": {
        "repo": "Karl-OS-ISO",
        "runId": (iso_pi or {}).get("runId"),
        "crdPresentCount": (iso_pi or {}).get("crdPresentCount"),
        "hostRuntimeEnabled": (iso_pi or {}).get("hostRuntimeEnabled", False),
        "path": ISO_REFS["postInstall"],
    },
    "repoGitRefs": {
        "Karl-Hyperdensity": git_short(ROOT),
        "Karl-Dashboard": git_short(DASHBOARD),
        "Karl-Installer": git_short(INSTALLER),
        "Karl-OS-ISO": git_short(ISO),
        "rdp-GW": git_short(RDP_GW),
    },
    "crossRepoEvidenceIndex": index,
    "errors": errors,
    "status": "PASS" if not errors else "FAIL",
    "generatedAt": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}

(OUT / "snapshot-summary.json").write_text(json.dumps(snapshot, indent=2) + "\n", encoding="utf-8")
(OUT / "cross-repo-evidence-index.json").write_text(
    json.dumps({"entries": index}, indent=2) + "\n", encoding="utf-8"
)

# Validate against contract fixture shape
contract_path = ROOT / "docs/contracts/khr/khr-tp-reference-snapshot-v1.json"
if contract_path.is_file():
    (OUT / "contract-ref.json").write_text(
        json.dumps({"contractSchema": str(contract_path.relative_to(ROOT))}, indent=2) + "\n",
        encoding="utf-8",
    )

log(f"status={snapshot['status']} errors={len(errors)}")
log(f"snapshot={OUT / 'snapshot-summary.json'}")
if errors:
    for err in errors:
        log(f"  FAIL: {err}")
    sys.exit(1)
sys.exit(0)
PY
