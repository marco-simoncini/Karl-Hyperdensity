# Scope-4 Guarded Apply Certification (KHR-BF)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-BF |
| **Capability** | TP live ResourceLease guarded apply (sandbox, evidence-backed) |
| **Evidence baseline** | KHR-BE `committed-scope4-guarded-apply-khr-be` |
| **Live mutation in BF** | **None** — certification reads committed evidence only |

---

## Purpose

Certify Scope-4 guarded apply as an **evidence-backed TP capability** using KHR-BE live sandbox evidence. Scope-4 remains **manual**, **not active**, **not enabled**, and **not autonomous**.

---

## Certification fields

| Field | KHR-BF certified value |
|-------|------------------------|
| `scope4CertificationState` | `certified-evidence-backed` |
| `evidenceRef` | `docs/evidence/khr-tp-live-scope4-guarded-apply/committed-scope4-guarded-apply-khr-be` |
| `mutationType` | `cpu.max` |
| `targetLane` | `native-live` |
| `rollbackVerified` | `true` |
| `continuityPreserved` | `true` |
| `noRestart` | `true` |
| `noRollout` | `true` |
| `noRecreate` | `true` |
| `noDisconnect` | `true` |
| `noRevoke` | `true` |
| `notPersistent` | `true` |
| `notAutonomous` | `true` |
| `readyForScope4` | `manual-guarded-apply-pass` |
| `readyForScope4Active` | `false` |
| `guardedApplyEnabled` | `false` |
| `guardedApplyAutonomous` | `false` |

---

## Evidence requirements (KHR-BE)

| Artifact | Required |
|----------|----------|
| `apply-summary.json` | `status=PASS`, `mutationScope=cpu.max`, `lane=native-live` |
| `verify-summary.json` | `readyForScope4=manual-guarded-apply-pass`, `readyForScope4Active=false`, continuity + no-restart flags |
| `rollback-summary.json` | `rollbackVerified=true` |
| `production-before.json` / `production-after-apply.json` | Production deploy generations unchanged |
| `apply-output.json` | `applied=true`, dry-run allowed, rollback/verification plan refs |

---

## Certification check (read-only)

```bash
./scripts/khr_scope4_certification_check.sh
```

Output: `docs/evidence/khr-scope4-guarded-apply-certification/committed-scope4-certification-khr-bf/certification-summary.json`

---

## Non-goals (KHR-BF)

| Forbidden | Reason |
|-----------|--------|
| Re-run guarded apply | Uses KHR-BE evidence only |
| Live failure injection | See `SCOPE4_FAILURE_SEMANTICS.md` (simulate/document only) |
| Production namespace mutation | Out of scope |
| Autonomous / persistent apply | Explicitly false in certification |
| Scope-5 | Does not exist |

---

## Related

- `KHR_TP_LIVE_SCOPE4_GUARDED_APPLY_PLAN.md`
- `SCOPE4_FAILURE_SEMANTICS.md`
- `scripts/khr_scope4_certification_check.sh`
