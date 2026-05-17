#!/usr/bin/env bash
# KHR-BF: read-only Scope-4 guarded-apply certification from KHR-BE evidence (no live mutation).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

RUN_ID="${KHR_SCOPE4_CERTIFICATION_RUN_ID:-committed-scope4-certification-khr-bf}"
OUT_DIR="${ROOT}/docs/evidence/khr-scope4-guarded-apply-certification/${RUN_ID}"
EVIDENCE_RUN_ID="${KHR_SCOPE4_EVIDENCE_RUN_ID:-committed-scope4-guarded-apply-khr-be}"
EVIDENCE_DIR="${ROOT}/docs/evidence/khr-tp-live-scope4-guarded-apply/${EVIDENCE_RUN_ID}"

mkdir -p "${OUT_DIR}"
log() { echo "[khr_scope4_certification_check] $*" | tee -a "${OUT_DIR}/run.log"; }

log "runId=${RUN_ID} evidenceRunId=${EVIDENCE_RUN_ID}"

export ROOT OUT_DIR RUN_ID EVIDENCE_DIR EVIDENCE_RUN_ID
python3 <<'PY'
from __future__ import annotations

import json
import os
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

ROOT = Path(os.environ["ROOT"])
OUT = Path(os.environ["OUT_DIR"])
RUN_ID = os.environ["RUN_ID"]
EVIDENCE_DIR = Path(os.environ["EVIDENCE_DIR"])
EVIDENCE_RUN_ID = os.environ["EVIDENCE_RUN_ID"]

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


def require_file(path: Path, label: str) -> dict[str, Any] | None:
    doc = load(path)
    record(f"artifact_{label}", path.is_file(), str(path))
    if not doc:
        errors.append(f"missing {path}")
    return doc


apply_sum = require_file(EVIDENCE_DIR / "apply-summary.json", "applySummary")
verify_sum = require_file(EVIDENCE_DIR / "verify-summary.json", "verifySummary")
rollback_sum = require_file(EVIDENCE_DIR / "rollback-summary.json", "rollbackSummary")
apply_out = require_file(EVIDENCE_DIR / "apply-output.json", "applyOutput")
prod_before = require_file(EVIDENCE_DIR / "production-before.json", "productionBefore")
prod_after = require_file(EVIDENCE_DIR / "production-after-apply.json", "productionAfterApply")

apply_ok = bool(
    apply_sum
    and apply_sum.get("status") == "PASS"
    and apply_sum.get("guardedApplyExecuted") is True
    and apply_sum.get("mutationScope") == "cpu.max"
    and apply_sum.get("lane") == "native-live"
    and apply_sum.get("applyScope") == "sandbox-only"
    and apply_sum.get("noAutonomousOrchestration") is True
    and apply_sum.get("noPersistentRuntimeLoop") is True
)
record("applySummaryValid", apply_ok, f"mutationScope={apply_sum.get('mutationScope') if apply_sum else 'missing'}")

verify_ok = bool(
    verify_sum
    and verify_sum.get("status") == "PASS"
    and verify_sum.get("readyForScope4") == "manual-guarded-apply-pass"
    and verify_sum.get("readyForScope4Active") is False
    and verify_sum.get("guardedApplyObserved") is True
    and verify_sum.get("continuityPreserved") is True
    and verify_sum.get("rollbackVerified") is True
    and verify_sum.get("noRestartObserved") is True
    and verify_sum.get("noRolloutObserved") is True
    and verify_sum.get("noRecreateObserved") is True
    and verify_sum.get("productionGatewayUntouched") is True
)
record("verifySummaryValid", verify_ok)

rollback_ok = bool(
    rollback_sum
    and rollback_sum.get("status") == "PASS"
    and rollback_sum.get("rollbackVerified") is True
    and rollback_sum.get("rollbackObserved") is True
)
record("rollbackSummaryValid", rollback_ok)

prod_ok = bool(
    prod_before
    and prod_after
    and prod_before.get("productionDeployGenerations")
    == prod_after.get("productionDeployGenerations")
)
record("noProductionMutation", prod_ok)

rdpgw_ok = True
rdpgw = load(EVIDENCE_DIR / "rdpgw-continuity-summary.json")
if rdpgw:
    rdpgw_ok = (
        rdpgw.get("noRevoke") is True
        and rdpgw.get("noDisconnect") is True
        and rdpgw.get("continuityObserved") is True
    )
record("rdpgwContinuityEvidence", rdpgw_ok, "optional rdpgw-continuity-summary.json")

dry_run_ok = bool(
    apply_out
    and apply_out.get("applied") is True
    and apply_out.get("dryRun", {}).get("dryRunDecision") == "allowed"
    and apply_out.get("dryRun", {}).get("rollbackPlanRef")
    and apply_out.get("dryRun", {}).get("verificationPlanRef")
)
record("applyOutputValid", dry_run_ok)

# Failure semantics fixtures exist (simulate only)
fixture_base = ROOT / "examples/khr/scope4-failure-semantics"
fixture_names = [
    "missing-rollback-plan.json",
    "stale-provenance.json",
    "failed-verification.json",
    "rollback-failure.json",
    "continuity-regression.json",
]
fixtures_ok = all((fixture_base / n).is_file() for n in fixture_names)
record("failureSemanticsFixtures", fixtures_ok, f"count={len(fixture_names)}")

all_ok = not errors and apply_ok and verify_ok and rollback_ok and prod_ok and dry_run_ok and fixtures_ok
cert_state = "certified-evidence-backed" if all_ok else "not-certified"

summary: dict[str, Any] = {
    "phase": "khr-scope4-guarded-apply-certification",
    "sprint": "KHR-BF",
    "runId": RUN_ID,
    "status": "PASS" if all_ok else "FAIL",
    "scope4CertificationState": cert_state,
    "evidenceRef": f"docs/evidence/khr-tp-live-scope4-guarded-apply/{EVIDENCE_RUN_ID}",
    "evidenceRunId": EVIDENCE_RUN_ID,
    "mutationType": "cpu.max",
    "targetLane": "native-live",
    "rollbackVerified": rollback_ok,
    "continuityPreserved": verify_sum.get("continuityPreserved", False) if verify_sum else False,
    "noRestart": verify_sum.get("noRestartObserved", False) if verify_sum else False,
    "noRollout": verify_sum.get("noRolloutObserved", False) if verify_sum else False,
    "noRecreate": verify_sum.get("noRecreateObserved", False) if verify_sum else False,
    "noDisconnect": rdpgw_ok,
    "noRevoke": rdpgw_ok,
    "notPersistent": apply_sum.get("noPersistentRuntimeLoop", True) if apply_sum else False,
    "notAutonomous": apply_sum.get("noAutonomousOrchestration", True) if apply_sum else False,
    "readyForScope4": "manual-guarded-apply-pass" if all_ok else False,
    "readyForScope4Active": False,
    "readyForScope4Enabled": False,
    "guardedApplyEnabled": False,
    "guardedApplyAutonomous": False,
    "readOnly": True,
    "liveMutationPerformed": False,
    "checks": checks,
    "errors": errors,
    "evidencePath": f"docs/evidence/khr-scope4-guarded-apply-certification/{RUN_ID}",
    "at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
(OUT / "certification-summary.json").write_text(json.dumps(summary, indent=2) + "\n")
print(
    f"[khr_scope4_certification_check] status={summary['status']} "
    f"scope4CertificationState={summary['scope4CertificationState']}"
)
if not all_ok:
    for e in errors:
        print(f"[khr_scope4_certification_check] FAIL: {e}", file=__import__("sys").stderr)
    raise SystemExit(1)
PY

log "PASS certification-summary=${OUT_DIR}/certification-summary.json"
