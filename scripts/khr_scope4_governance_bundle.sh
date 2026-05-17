#!/usr/bin/env bash
# KHR-BG: read-only Scope-4 operational governance bundle (no live mutation).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RUN_ID="${KHR_SCOPE4_GOVERNANCE_RUN_ID:-committed-scope4-governance-khr-bg}"
OUT_DIR="${ROOT}/docs/evidence/khr-scope4-governance/${RUN_ID}"
CERT_RUN_ID="${KHR_SCOPE4_CERTIFICATION_RUN_ID:-committed-scope4-certification-khr-bf}"
APPLY_RUN_ID="${KHR_SCOPE4_EVIDENCE_RUN_ID:-committed-scope4-guarded-apply-khr-be}"

mkdir -p "${OUT_DIR}"
log() { echo "[khr_scope4_governance_bundle] $*" | tee -a "${OUT_DIR}/run.log"; }

log "runId=${RUN_ID} certificationRunId=${CERT_RUN_ID} applyRunId=${APPLY_RUN_ID}"

export ROOT OUT_DIR RUN_ID CERT_RUN_ID APPLY_RUN_ID
python3 <<'PY'
from __future__ import annotations

import json
import os
from datetime import datetime, timedelta, timezone
from pathlib import Path
from typing import Any

ROOT = Path(os.environ["ROOT"])
OUT = Path(os.environ["OUT_DIR"])
RUN_ID = os.environ["RUN_ID"]
CERT_RUN_ID = os.environ["CERT_RUN_ID"]
APPLY_RUN_ID = os.environ["APPLY_RUN_ID"]

CERT_DAYS_VALID = 180
CERT_DAYS_STALE = 90


def load(path: Path) -> dict[str, Any] | None:
    if not path.is_file():
        return None
    try:
        return json.loads(path.read_text())
    except (OSError, json.JSONDecodeError):
        return None


def latest_federation() -> tuple[Path | None, dict[str, Any] | None]:
    base = ROOT / "docs/evidence/khr-runtime-observation-federation"
    if not base.is_dir():
        return None, None
    best: tuple[Path, dict[str, Any]] | None = None
    for child in sorted(base.iterdir(), reverse=True):
        p = child / "federation-summary.json"
        if p.is_file():
            doc = load(p)
            if doc and doc.get("status") == "PASS":
                return p, doc
            if doc and best is None:
                best = (p, doc)
    return best


checks: dict[str, dict[str, Any]] = {}
errors: list[str] = []


def record(name: str, ok: bool, detail: str = "", ref: str = "") -> None:
    checks[name] = {"status": "PASS" if ok else "FAIL", "detail": detail, "ref": ref}
    if not ok:
        errors.append(f"{name}: {detail}")


# Certification (KHR-BF)
cert_path = ROOT / "docs/evidence/khr-scope4-guarded-apply-certification" / CERT_RUN_ID / "certification-summary.json"
cert = load(cert_path)
record("certificationSummary", bool(cert and cert.get("status") == "PASS"), str(cert_path), str(cert_path.relative_to(ROOT)))

apply_dir = ROOT / "docs/evidence/khr-tp-live-scope4-guarded-apply" / APPLY_RUN_ID
apply_sum = load(apply_dir / "apply-summary.json")
verify_sum = load(apply_dir / "verify-summary.json")
rollback_sum = load(apply_dir / "rollback-summary.json")
continuity = load(apply_dir / "continuity-proof-apply.json")
rdpgw_cont = load(apply_dir / "rdpgw-continuity-summary.json")

record("applySummary", bool(apply_sum and apply_sum.get("status") == "PASS"), "KHR-BE apply", str((apply_dir / "apply-summary.json").relative_to(ROOT)))
record("verifySummary", bool(verify_sum and verify_sum.get("status") == "PASS"), "KHR-BE verify", str((apply_dir / "verify-summary.json").relative_to(ROOT)))
record(
    "rollbackSummary",
    bool(rollback_sum and rollback_sum.get("rollbackVerified") is True),
    "rollbackVerified",
    str((apply_dir / "rollback-summary.json").relative_to(ROOT)),
)
record(
    "continuityEvidence",
    bool(verify_sum and verify_sum.get("continuityPreserved") is True),
    "continuityPreserved",
    str((apply_dir / "verify-summary.json").relative_to(ROOT)),
)

prov = load(ROOT / "docs/evidence/khr-provenance/summary.json")
record(
    "provenanceEvidence",
    bool(prov and prov.get("readOnly") is True and prov.get("noAutonomousOrchestration") is True),
    f"runId={prov.get('runId') if prov else 'none'}",
    "docs/evidence/khr-provenance/summary.json",
)

fed_path, fed = latest_federation()
record(
    "federationEvidence",
    bool(fed and fed.get("status") == "PASS" and fed.get("readOnly") is True),
    str(fed_path.relative_to(ROOT)) if fed_path else "missing",
    str(fed_path.relative_to(ROOT)) if fed_path else "",
)

cg = load(ROOT / "docs/evidence/khr-control-graph/summary.json")
record(
    "policyGates",
    bool(cg and cg.get("readOnly") is True and cg.get("noApply") is True),
    "control-graph summary",
    "docs/evidence/khr-control-graph/summary.json",
)

approval = load(ROOT / "docs/evidence/khr-action-approval/summary.json")
record(
    "approvalEvidence",
    bool(approval and approval.get("readOnly") is True and approval.get("noApply") is True),
    "action-approval summary",
    "docs/evidence/khr-action-approval/summary.json",
)

# Lifecycle timestamps
cert_at: datetime | None = None
if cert and cert.get("at"):
    cert_at = datetime.strptime(cert["at"], "%Y-%m-%dT%H:%M:%SZ").replace(tzinfo=timezone.utc)
now = datetime.now(timezone.utc)
certification_expiry = None
stale_certification = False
if cert_at:
    certification_expiry = (cert_at + timedelta(days=CERT_DAYS_VALID)).strftime("%Y-%m-%dT%H:%M:%SZ")
    stale_certification = (now - cert_at).days > CERT_DAYS_STALE

expired = bool(cert_at and now > cert_at + timedelta(days=CERT_DAYS_VALID))
revoked = bool(prov and prov.get("registryIntegrity") is False)
regression = bool(
    verify_sum
    and (
        verify_sum.get("continuityPreserved") is False
        or verify_sum.get("noRestartObserved") is False
        or verify_sum.get("rollbackVerified") is False
    )
)

core_ok = all(c["status"] == "PASS" for c in checks.values())

if revoked:
    governance_state = "revoked"
elif regression:
    governance_state = "regression-detected"
elif expired:
    governance_state = "expired"
elif stale_certification or (prov and prov.get("registryIntegrity") is not True):
    governance_state = "stale"
elif core_ok and cert and cert.get("scope4CertificationState") == "certified-evidence-backed":
    governance_state = "certified"
else:
    governance_state = "stale"

operator_revalidation = governance_state in (
    "stale",
    "expired",
    "revoked",
    "regression-detected",
)

summary: dict[str, Any] = {
    "phase": "khr-scope4-operational-governance",
    "sprint": "KHR-BG",
    "runId": RUN_ID,
    "status": "PASS" if governance_state == "certified" and core_ok else "FAIL",
    "readOnly": True,
    "liveMutationPerformed": False,
    "scope4GovernanceState": governance_state,
    "scope4CertificationState": cert.get("scope4CertificationState") if cert else None,
    "certificationExpiry": certification_expiry,
    "staleCertification": stale_certification or governance_state == "stale",
    "regressionDetected": regression or governance_state == "regression-detected",
    "operatorRevalidationRequired": operator_revalidation,
    "readyForScope4": cert.get("readyForScope4") if cert else False,
    "readyForScope4Active": False,
    "readyForScope4Enabled": False,
    "guardedApplyEnabled": False,
    "guardedApplyAutonomous": False,
    "notAutonomous": True,
    "notPersistent": True,
    "mutationType": "cpu.max",
    "targetLane": "native-live",
    "applyScope": "sandbox-only",
    "evidenceRetention": {
        "applyAnchor": f"docs/evidence/khr-tp-live-scope4-guarded-apply/{APPLY_RUN_ID}",
        "certificationAnchor": f"docs/evidence/khr-scope4-guarded-apply-certification/{CERT_RUN_ID}",
        "governanceBundle": f"docs/evidence/khr-scope4-governance/{RUN_ID}",
    },
    "aggregatedEvidence": {
        "certificationSummary": str(cert_path.relative_to(ROOT)) if cert_path.is_file() else None,
        "rollbackSummary": str((apply_dir / "rollback-summary.json").relative_to(ROOT)),
        "continuityEvidence": str((apply_dir / "continuity-proof-apply.json").relative_to(ROOT)) if continuity else str((apply_dir / "verify-summary.json").relative_to(ROOT)),
        "provenanceEvidence": "docs/evidence/khr-provenance/summary.json",
        "federationEvidence": str(fed_path.relative_to(ROOT)) if fed_path else None,
        "policyGates": "docs/evidence/khr-control-graph/summary.json",
        "approvalEvidence": "docs/evidence/khr-action-approval/summary.json",
        "rdpgwContinuity": str((apply_dir / "rdpgw-continuity-summary.json").relative_to(ROOT)) if rdpgw_cont else None,
    },
    "checks": checks,
    "errors": errors,
    "evidencePath": f"docs/evidence/khr-scope4-governance/{RUN_ID}",
    "at": now.strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(OUT / "governance-summary.json").write_text(json.dumps(summary, indent=2) + "\n")
print(
    f"[khr_scope4_governance_bundle] status={summary['status']} "
    f"scope4GovernanceState={governance_state} operatorRevalidationRequired={operator_revalidation}"
)
if summary["status"] != "PASS":
    for e in errors:
        print(f"[khr_scope4_governance_bundle] FAIL: {e}", file=__import__("sys").stderr)
    if governance_state != "certified":
        print(f"[khr_scope4_governance_bundle] governance state={governance_state}", file=__import__("sys").stderr)
    raise SystemExit(1)
PY

log "PASS governance-summary=${OUT_DIR}/governance-summary.json"
