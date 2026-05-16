# Hyperdensity Parent Fabric — live observation semantic prototype (Sprint 62)

## Summary

**Sprint 62** introduces a **semantic candidate** for live pod observation. **Sprint 63** connects the candidate to wrapper **true branch** in `observation_wiring_v1.go`. **`live.go`** still calls wrappers only. Broad observation remains **false**.

---

## 1. Scope

| Item | Sprint 62 |
|------|-----------|
| Candidate helpers in `*_live_observation_candidate_v1.go` | **Yes** |
| Candidate shadow tests (10 cases) | **Yes** |
| Runtime wrapper branch swap | **No** |
| `ObservationWiredV1` / `ProductionWiredV1` | **`false`** |
| Hyperdensity Go adapter code | **No** |

---

## 2. Non-goals

- Wire candidate into `observation_wiring_v1.go` true branch
- Set `hyperdensityWorkloadAdapterLiveObservationCandidateRuntimeUsedV1 = true`
- Broad observation phase (`ObservationWiredV1 = true`)
- apply / resource_exchange / rollback / VM runtime wiring

---

## 3. Prototype shape

| Function | Role |
|----------|------|
| `hyperdensityWorkloadLiveCandidateObservedPodUIDV1` | Pod UID from `map[string]interface{}` |
| `hyperdensityWorkloadLiveCandidateObservedPod*FromObservationV1` | Quantities via struct container lookup |

Constants: `LiveObservationCandidateV1 = true`, `LiveObservationCandidateRuntimeUsedV1 = false`.

---

## 4. Shadow matrix (10 cases)

Pod UID (normal, nil/empty), CPU/memory request/limit, missing container, multi-container selection, empty container name, zero-value observation. Each case: **candidate ≡ legacy ≡ wrapper**.

---

## 5. Criteria before runtime branch substitution

See **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_BRANCH_SWAP_CRITERIA.md`**. Proposed **Sprint 63** for branch swap only after Sprint 62 PASS.

---

## 6. Rollback

Remove candidate file and tests; no runtime change required (candidate not wired).

---

## 7. Test coverage

| Artifact | Owner |
|----------|-------|
| `hyperdensity_parent_fabric_workload_adapter_live_observation_candidate_v1.go` | Dashboard |
| `*_candidate_shadow_test.go` | Dashboard |
| `*_branch_swap_guard_test.go` | Dashboard |

---

## 8. Risks

| Risk | Mitigation |
|------|------------|
| Candidate duplicates legacy logic | Shadow triple-assert (candidate, wrapper, legacy) |
| Accidental runtime wiring | Branch-swap guard + audit script |

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_FLIP.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_BRANCH_SWAP_CRITERIA.md`
