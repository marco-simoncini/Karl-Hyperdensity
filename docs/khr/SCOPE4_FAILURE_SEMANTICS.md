# Scope-4 Failure Semantics (KHR-BF)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-BF |
| **Mode** | **Read-only / simulate / document** — **no live failure injection** |

---

## Purpose

Define how Scope-4 guarded apply **would block or fail** without executing mutation on cluster. Operators and consumers use these semantics for Dashboard projection, Inventory observation stubs, and certification regression review.

---

## Failure classes (documented)

| Class | Symptom | Simulated fixture | Live injection |
|-------|---------|-------------------|----------------|
| `missingRollbackPlan` | Dry-run/apply blocked: no `rollbackPlanRef` | `examples/khr/scope4-failure-semantics/missing-rollback-plan.json` | **Forbidden** |
| `staleProvenance` | Provenance gate fails before apply | `examples/khr/scope4-failure-semantics/stale-provenance.json` | **Forbidden** |
| `failedVerification` | Post-apply `verification.state != pass` | `examples/khr/scope4-failure-semantics/failed-verification.json` | **Forbidden** |
| `rollbackFailure` | Rollback cannot restore baseline | `examples/khr/scope4-failure-semantics/rollback-failure.json` | **Forbidden** |
| `continuityRegression` | Restart/rollout/recreate or session break | `examples/khr/scope4-failure-semantics/continuity-regression.json` | **Forbidden** |

---

## Projection mapping (Dashboard / Inventory)

| Semantic field | When failure class active |
|----------------|---------------------------|
| `failedVerification` | `verification.state=fail` or cgroup mismatch |
| `rollbackFailure` | `rollbackVerified=false` |
| `staleEvidence` | Provenance/registry stale vs evidence runId |
| `continuityRegressionObserved` | `noRestart`, `noRollout`, `noRecreate`, or gateway continuity false |

All failure projections remain **read-only** — no apply buttons, no autonomous retry.

---

## PASS baseline (KHR-BE)

KHR-BE evidence demonstrates **none** of the failure classes above at certification time:

- Evidence: `docs/evidence/khr-tp-live-scope4-guarded-apply/committed-scope4-guarded-apply-khr-be/`
- Certification: `docs/evidence/khr-scope4-guarded-apply-certification/committed-scope4-certification-khr-bf/`

---

## Related

- `SCOPE4_GUARDED_APPLY_CERTIFICATION.md`
- Karl-Dashboard `DASHBOARD_TP_READINESS_REFERENCE_ENV.md` (failure projection section)
