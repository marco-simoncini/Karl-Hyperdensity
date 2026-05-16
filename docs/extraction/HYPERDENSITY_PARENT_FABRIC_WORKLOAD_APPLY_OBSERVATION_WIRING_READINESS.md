# Hyperdensity Parent Fabric — apply observation wiring readiness (Sprint 69)

## Summary

**Sprint 69** formalizes **readiness** for future `apply.go` call-site wiring. **Readiness-only** — no `apply.go` changes, no flag flips. **Sprint 70** is the first allowed call-site wiring sprint.

---

## 1. Scope

| Item | Sprint 69 |
|------|-----------|
| Wiring readiness doc + Dashboard golden/test | **Yes** |
| `readyForApplyGoCallSiteWiring` certification | **Yes** |
| `apply.go` call-site repoint | **No** (Sprint 70+) |
| `ApplyObservationWiredV1 = true` | **No** |

---

## 2. Non-goals

- Modifying `apply.go` call sites.
- `ApplyObservationCandidateRuntimeUsedV1 = true`.
- `ObservationWiredV1` or `ProductionWiredV1` flip.
- resource_exchange, rollback, VM, admission_guard wiring.
- Hyperdensity Go adapter code or Dashboard `parentfabric` import.

---

## 3. Preconditions satisfied

| Precondition | Sprint |
|--------------|--------|
| Apply observation proposal + audit | 65 |
| Shadow matrix (candidate) | 66 |
| Staged wrappers | 67 |
| Wrapper hardening 8×4 | 68 |
| Wrappers ≡ legacy ≡ candidate (flags false) | 67–68 |

---

## 4. Remaining no-touch invariants

- `apply.go`: four legacy helper call sites only.
- No wrapper/candidate names in `apply.go`.
- No apply helpers in resource_exchange, rollback, VM runtime.
- All apply observation flags **false**.

---

## 5. Allowed future Sprint 70 change

**Sprint 70** may repoint the four `apply.go` call sites from legacy helpers to:

- `hyperdensityWorkloadApplyObservedPodMemoryRequestV1`
- `hyperdensityWorkloadApplyObservedPodMemoryLimitV1`
- `hyperdensityWorkloadApplyObservedPodCPURequestV1`
- `hyperdensityWorkloadApplyObservedPodCPULimitV1`

With `ApplyObservationWiredV1` and `ApplyObservationCandidateRuntimeUsedV1` still **false**, wrappers remain legacy-equivalent.

**Flag flips** (`ApplyObservationWiredV1`, candidate runtime) require **separate** dedicated sprints after call-site wiring is stable.

---

## 6. Rollback plan

Sprint 69 is docs/test only — no runtime rollback.

Sprint 70 rollback (if wired): restore four legacy call sites in `apply.go`; delete wrapper references; flags unchanged if never flipped.

---

## 7. Readiness checklist

- [x] Proposal + criteria (Sprint 65)
- [x] Shadow matrix (Sprint 66)
- [x] Staged wrappers (Sprint 67)
- [x] Hardening 8×4 (Sprint 68)
- [x] Readiness golden `readyForApplyGoCallSiteWiring: true` (Sprint 69)
- [ ] Sprint 70 call-site wiring (not started)
- [ ] `ApplyObservationWiredV1` flip (not started)

---

## 8. Risks

Premature Sprint 70 without parity green could change apply plan target fields. Broad observation flip must remain blocked.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_WRAPPER_HARDENING.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CRITERIA.md`
