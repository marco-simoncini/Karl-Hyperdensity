# Hyperdensity Parent Fabric — live observation branch swap (Sprint 63)

## Summary

**Sprint 63** swaps the **true branch** of live observation wrappers in `observation_wiring_v1.go` to use **semantic candidate helpers** introduced in Sprint 62. **`CandidateRuntimeUsedV1 = true`**, **`branchSwapAllowed = true`**. Broad **`ObservationWiredV1`** and **`ProductionWiredV1`** remain **`false`**. **Sprint 64** re-audit confirms this does **not** complete general observation wiring.

---

## 1. Scope

| Item | Sprint 63 |
|------|-----------|
| True branch → candidate helpers | **Yes** |
| `live.go` call sites | **Unchanged** (still use wrappers only) |
| Broad observation | **No** |
| Hyperdensity Go code | **No** |

---

## 2. Non-goals

- `hyperdensityWorkloadAdapterObservationWiredV1 = true`
- Wire candidate into apply / resource_exchange / rollback / VM runtime
- Dashboard `parentfabric` import

---

## 3. Preconditions from Sprint 62

- 10-case candidate shadow PASS
- `branchSwapAllowed` criteria documented
- Dedicated branch-swap sprint (this sprint)

---

## 4. Gate / constants

| Constant | Value |
|----------|-------|
| `LiveObservationWiredV1` | **`true`** |
| `LiveObservationCandidateRuntimeUsedV1` | **`true`** |
| `ObservationWiredV1` | **`false`** |
| `ProductionWiredV1` | **`false`** |

---

## 5. Branch swap details

When **`LiveObservationWiredV1 && CandidateRuntimeUsedV1`**, wrappers delegate to `hyperdensityWorkloadLiveCandidate*V1`. **`live.go`** continues calling wrappers only — not candidate directly.

---

## 6. Behavior equivalence statement

Candidate helpers were proven ≡ legacy in Sprint 62 shadow. Post-swap runtime output for live observation paths is **unchanged**.

---

## 7. Rollback

Set `CandidateRuntimeUsedV1 = false`; revert `observation_wiring_v1.go` true branch to legacy helpers; set `branchSwapAllowed = false` in goldens.

---

## 8. Test coverage

| Artifact | Owner |
|----------|-------|
| `*_branch_swap_test.go` | Dashboard |
| `*_branch_swap.golden.json` | Dashboard |
| Updated candidate shadow golden | Dashboard |

---

## 9. Risks

| Risk | Mitigation |
|------|------------|
| Future candidate drift from legacy | Keep 10-case shadow in parity |
| Confusion with broad observation | `ObservationWiredV1` stays false |

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_BRANCH_SWAP_CRITERIA.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_SEMANTIC_PROTOTYPE.md`
