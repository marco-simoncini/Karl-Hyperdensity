# Hyperdensity Parent Fabric — live observation branch swap criteria

## Summary

Criteria before changing runtime wrapper **true branch** from legacy delegate to **semantic candidate** helpers. **`branchSwapAllowed`** remains **`false`** until a **dedicated branch-swap sprint** (proposed Sprint 63).

---

## Minimum criteria before branch swap

| # | Criterion |
|---|-----------|
| 1 | **Semantic prototype shadow PASS** — 10 cases, candidate ≡ legacy |
| 2 | **Missing/empty pod** cases PASS | 
| 3 | **Multiple container** selection PASS |
| 4 | **Call-site audit PASS** |
| 5 | **`ObservationWiredV1`** remains **`false`** |
| 6 | No apply / resource_exchange / rollback / VM runtime wiring |
| 7 | **Dedicated sprint** for branch swap only |
| 8 | **Golden update** — `branchSwapAllowed: true` in candidate shadow golden |
| 9 | **Rollback** — set `LiveObservationCandidateRuntimeUsedV1 = false` + revert wrapper true branch to legacy |
| 10 | **Parity runner** green |

---

## Sprint 62 status

| Item | Status |
|------|--------|
| Candidate present | **Yes** |
| Candidate runtime used | **No** |
| Branch swap allowed | **No** |
| Runtime true branch | **Legacy-equivalent** |

---

## Forbidden with branch swap (until criteria met)

- `hyperdensityWorkloadAdapterObservationWiredV1 = true`
- Candidate calls from files other than candidate module + tests
- Dashboard `parentfabric` import

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_SEMANTIC_PROTOTYPE.md`
