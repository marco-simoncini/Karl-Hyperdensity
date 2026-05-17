#!/usr/bin/env bash
# KHR-BV: Beta Candidate 0 release marker check (offline; no cluster).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

MANIFEST="${ROOT}/docs/contracts/khr/khr-beta-candidate-0-manifest.json"

for doc in \
  docs/khr/KHR_BETA_CANDIDATE_0_SCOPE.md \
  docs/khr/KHR_BETA_CANDIDATE_0_RELEASE_MARKER.md \
  docs/khr/KHR_BETA_READINESS_PLAN.md \
  docs/khr/KHR_VALIDATION_MODES.md \
  docs/khr/KHR_SNAPSHOT_V1_FREEZE_POLICY.md; do
  [[ -f "${ROOT}/${doc}" ]] || {
    echo "[khr_beta_candidate_0_check] FAIL: missing ${doc}" >&2
    exit 1
  }
done

[[ -f "${MANIFEST}" ]] || {
  echo "[khr_beta_candidate_0_check] FAIL: missing manifest" >&2
  exit 1
}

KHR_TP_REFERENCE_SNAPSHOT_RUN_ID="${KHR_TP_REFERENCE_SNAPSHOT_RUN_ID:-committed-khr-bt-v1}" \
  "${ROOT}/scripts/khr_validate_reference_snapshot.sh"

export ROOT MANIFEST
python3 <<'PY'
import json
import os
import subprocess
import sys
from pathlib import Path

ROOT = Path(os.environ["ROOT"])
MANIFEST = Path(os.environ["MANIFEST"])
errors: list[str] = []

manifest = json.loads(MANIFEST.read_text(encoding="utf-8"))

if manifest.get("contractVersion") != "khr-beta-candidate-0":
    errors.append("contractVersion")
if manifest.get("releaseMarker") != "khr-beta-candidate-0":
    errors.append("releaseMarker")
if manifest.get("productionReady") is not False:
    errors.append("productionReady must be false")
if manifest.get("autonomousOrchestration") is not False:
    errors.append("autonomousOrchestration must be false")
if manifest.get("globalDefaultsChanged") is not False:
    errors.append("globalDefaultsChanged must be false")

snap_ref = manifest.get("snapshotRef") or {}
snap_path = ROOT / snap_ref.get("summaryPath", "")
if not snap_path.is_file():
    errors.append(f"missing snapshot: {snap_path}")
else:
    snap = json.loads(snap_path.read_text(encoding="utf-8"))
    if snap.get("status") != "PASS":
        errors.append(f"snapshot status={snap.get('status')}")
    if snap.get("runId") != snap_ref.get("runId"):
        errors.append("snapshot runId mismatch")
    live = snap.get("dashboardLivePassRef") or {}
    if live.get("evidenceStatus") != "LIVE_PASS":
        errors.append("dashboard LIVE_PASS")
    if live.get("rollbackVerified") is not True:
        errors.append("dashboard rollbackVerified")
    if snap.get("scope4CertificationState") != "certified-evidence-backed":
        errors.append("scope4CertificationState")
    if snap.get("governanceState") != "certified":
        errors.append("governanceState")

def _repo_root(name: str, env_key: str, *extra: str) -> Path:
    if os.environ.get(env_key):
        return Path(os.environ[env_key])
    for c in (ROOT.parent / name, *(Path(p) for p in extra)):
        if c.is_dir():
            return c
    return ROOT.parent / name


_repo_roots = {
    "Karl-Hyperdensity": ROOT,
    "Karl-Dashboard": _repo_root("Karl-Dashboard", "KHR_DASHBOARD_PATH"),
    "Karl-Installer": _repo_root("Karl-Installer", "KHR_INSTALLER_PATH"),
    "Karl-OS-ISO": _repo_root("Karl-OS-ISO", "KHR_ISO_PATH"),
    "rdp-GW": _repo_root("rdp-GW", "KHR_RDP_GW_PATH", "/home/m.simoncini/rdp-GW"),
    "Karl-Inventory": _repo_root("Karl-Inventory", "KHR_INVENTORY_PATH"),
}

for ref_name, ref in (manifest.get("evidenceRefs") or {}).items():
    repo = ref.get("repo")
    rel = ref.get("path")
    if not rel:
        continue
    base = _repo_roots.get(repo)
    p = (base / rel) if base and base.is_dir() else None
    if p is None or not p.is_file():
        errors.append(f"evidence missing: {repo}:{rel}")
        continue
    if ref_name == "dashboardLivePass":
        doc = json.loads(p.read_text(encoding="utf-8"))
        if doc.get("evidenceStatus") != "LIVE_PASS":
            errors.append("dashboard live evidenceStatus")

repo_commits = manifest.get("repoCommits") or {}
for repo, pinned in repo_commits.items():
    rpath = _repo_roots.get(repo)
    if not rpath or not rpath.is_dir():
        continue
    head = subprocess.run(
        ["git", "-C", str(rpath), "rev-parse", "HEAD"],
        capture_output=True,
        text=True,
    )
    if head.returncode != 0:
        errors.append(f"git HEAD unavailable: {repo}")
        continue
    head_sha = head.stdout.strip()
    # Pinned commit must be ancestor of current HEAD (release marker minimum)
    anc = subprocess.run(
        ["git", "-C", str(rpath), "merge-base", "--is-ancestor", pinned, head_sha],
    )
    if anc.returncode != 0 and head_sha != pinned:
        errors.append(f"{repo}: HEAD {head_sha[:12]} not descendant of pin {pinned[:12]}")

if errors:
    for e in errors:
        print(f"[khr_beta_candidate_0_check] FAIL: {e}", file=sys.stderr)
    sys.exit(1)

print("[khr_beta_candidate_0_check] PASS releaseMarker=khr-beta-candidate-0")
PY

echo "[khr_beta_candidate_0_check] OK"
