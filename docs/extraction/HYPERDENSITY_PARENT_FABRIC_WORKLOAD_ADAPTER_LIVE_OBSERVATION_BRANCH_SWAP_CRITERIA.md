# Hyperdensity Parent Fabric — live observation branch swap criteria

## Summary

**Sprint 63** executed the dedicated branch swap. **Sprint 64** documents remaining surfaces and broad observation policy (`ObservationWiredV1` stays **false**).

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

## Sprint 62–63 status

| Item | Sprint 62 | Sprint 63 |
|------|-----------|-----------|
| Candidate present | **Yes** | **Yes** |
| Candidate runtime used | **No** | **Yes** |
| Branch swap allowed | **No** | **Yes** |
| Runtime true branch | Legacy delegate | **Candidate** (≡ legacy) |

---

## Forbidden with branch swap (until criteria met)

- `hyperdensityWorkloadAdapterObservationWiredV1 = true`
- Candidate calls from files other than candidate module + tests
- Dashboard `parentfabric` import

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_SEMANTIC_PROTOTYPE.md`
