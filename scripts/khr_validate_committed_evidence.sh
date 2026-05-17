#!/usr/bin/env bash
# KHR-BU: offline validation of committed TP evidence (no cluster / kubectl).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RUN_ID="${KHR_TP_REFERENCE_SNAPSHOT_RUN_ID:-committed-khr-bt-v1}"
SNAPSHOT="${ROOT}/docs/evidence/khr-tp-reference-snapshot-v1/${RUN_ID}/snapshot-summary.json"

[[ -f "${SNAPSHOT}" ]] || {
  echo "[khr_validate_committed_evidence] FAIL: missing snapshot ${SNAPSHOT}" >&2
  echo "Run: KHR_TP_REFERENCE_SNAPSHOT_RUN_ID=${RUN_ID} ./scripts/khr_tp_reference_snapshot_v1.sh" >&2
  exit 1
}

export ROOT SNAPSHOT RUN_ID
python3 <<'PY'
import json
import sys
from pathlib import Path

ROOT = Path(__import__("os").environ["ROOT"])
SNAPSHOT = Path(__import__("os").environ["SNAPSHOT"])
RUN_ID = __import__("os").environ["RUN_ID"]

snap = json.loads(SNAPSHOT.read_text(encoding="utf-8"))
errors: list[str] = []


def check_json(rel: str, *, status_key: str = "status", want_status: str = "PASS", alt_ok=None) -> None:
    path = ROOT / rel
    if not path.is_file():
        errors.append(f"missing {rel}")
        return
    doc = json.loads(path.read_text(encoding="utf-8"))
    if alt_ok and alt_ok(doc):
        return
    if doc.get(status_key) != want_status:
        errors.append(f"{rel}: {status_key}={doc.get(status_key)!r} want {want_status!r}")


def scope3_dryrun_ok(doc: dict) -> bool:
    return doc.get("status") == "PASS" and doc.get("noMutation") is True and doc.get("noApply") is True


# Snapshot invariants
if snap.get("status") != "PASS":
    errors.append(f"snapshot status={snap.get('status')}")
if snap.get("liveMutationPerformed") is not False:
    errors.append("snapshot liveMutationPerformed must be false")
if snap.get("globalDefaultsChanged") is not False:
    errors.append("snapshot globalDefaultsChanged must be false")

readiness = snap.get("scopeReadiness") or {}
for scope, want in (
    ("scope1", "PASS"),
    ("scope2", "manual-loop-pass"),
    ("scope3", "manual-dryrun-pass"),
    ("scope4", "manual-guarded-apply-pass"),
):
    if readiness.get(scope) != want:
        errors.append(f"scopeReadiness.{scope}={readiness.get(scope)!r} want {want!r}")

# Hyperdensity committed scope evidence (authoritative for offline CI)
check_json(
    "docs/evidence/khr-tp-live-scope2-preflight/committed-scope2-preflight-khr-az/scope2-preflight-summary.json"
)
check_json(
    "docs/evidence/khr-tp-live-scope2-resourceport-loop/committed-scope2-loop-khr-ba/loop-summary.json"
)
check_json(
    "docs/evidence/khr-tp-live-scope3-preflight/committed-scope3-preflight-khr-bb/scope3-preflight-summary.json"
)
# Scope-3 execute summary (PASS); verify-summary may be FAIL after stale live re-validation — not used offline
check_json(
    "docs/evidence/khr-tp-live-scope3-dryrun/committed-scope3-dryrun-khr-bc/dryrun-summary.json",
    alt_ok=scope3_dryrun_ok,
)
check_json(
    "docs/evidence/khr-scope4-guarded-apply-certification/committed-scope4-certification-khr-bf/certification-summary.json"
)
check_json(
    "docs/evidence/khr-scope4-governance/committed-scope4-governance-khr-bg/governance-summary.json"
)

# Hyperdensity index entries from snapshot
index = snap.get("crossRepoEvidenceIndex") or {}
entries = index.get("entries") if isinstance(index, dict) else index
for entry in entries or []:
    if entry.get("repo") != "Karl-Hyperdensity":
        continue
    rel = entry.get("relativePath")
    if not rel:
        continue
    if not (ROOT / rel).is_file():
        errors.append(f"index missing file: {rel}")
    elif entry.get("present") is False:
        errors.append(f"index present=false: {rel}")
    elif entry.get("status") not in (None, "PASS"):
        errors.append(f"index status={entry.get('status')}: {rel}")

if errors:
    for e in errors:
        print(f"[khr_validate_committed_evidence] FAIL: {e}", file=sys.stderr)
    sys.exit(1)

print(f"[khr_validate_committed_evidence] PASS runId={RUN_ID} (offline committed evidence)")
PY
