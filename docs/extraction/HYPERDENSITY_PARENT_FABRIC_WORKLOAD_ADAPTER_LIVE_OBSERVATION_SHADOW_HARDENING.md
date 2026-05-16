# Hyperdensity Parent Fabric — live observation shadow hardening (Sprint 60)

## Summary

**Sprint 60** adds **shadow tests** proving Sprint 59 live observation wrappers are **equivalent** to legacy helpers. **Sprint 61** flips **`LiveObservationWiredV1 = true`**; shadow tests remain PASS. **Sprint 62** extends with semantic candidate shadow (10 cases). Full **`workload_helpers.go`** verdict remains **`copy-deferred`**.

---

## 1. Scope

| Item | Sprint 60 |
|------|-----------|
| Live wrapper shadow tests (7 cases) | **Yes** |
| Enable `LiveObservationWiredV1` | **No** |
| Sprint 59 staged wrappers in `live.go` | **Unchanged** |
| `ObservationWiredV1` / `ProductionWiredV1` | **`false`** |
| Hyperdensity Go adapter code | **No** |

---

## 2. Non-goals

- Flip `hyperdensityWorkloadAdapterLiveObservationWiredV1` to `true`
- Set `hyperdensityWorkloadAdapterObservationWiredV1 = true`
- Wire apply, resource_exchange, rollback, VM runtime observation
- Change API responses or JSON ordering

---

## 3. Live wrapper shadow matrix

| # | Case | Wrapper | Legacy |
|---|------|---------|--------|
| 1 | Pod UID | `hyperdensityWorkloadLiveObservedPodUIDV1` | `hyperdensityObservedPodUID` |
| 2 | CPU request | `...CPURequestFromObservationV1` | `hyperdensityObservedPodCPURequestFromObservation` |
| 3 | CPU limit | `...CPULimitFromObservationV1` | `hyperdensityObservedPodCPULimitFromObservation` |
| 4 | Memory request | `...MemoryRequestFromObservationV1` | `hyperdensityObservedPodMemoryRequestFromObservation` |
| 5 | Memory limit | `...MemoryLimitFromObservationV1` | `hyperdensityObservedPodMemoryLimitFromObservation` |
| 6 | Missing container | all quantity wrappers | zero-value parity |
| 7 | Nil/empty pod map | UID wrapper | legacy UID |

---

## 4. Future flip criteria

See **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_FLIP_CRITERIA.md`**. Sprint 61 may flip **`LiveObservationWiredV1 = true`** only if Sprint 60 PASS.

---

## 5. Rollback

Remove shadow test + golden; revert docs. Sprint 59 staged wiring unchanged.

---

## 6. Test coverage

| Artifact | Owner |
|----------|-------|
| `hyperdensity_parent_fabric_workload_adapter_live_observation_shadow_test.go` | Dashboard |
| `testdata/..._live_observation_shadow.golden.json` | Dashboard |

---

## 7. Risks

| Risk | Mitigation |
|------|------------|
| Shadow PASS with identical legacy delegate in both branches | Flip criteria require dedicated sprint + updated golden |
| Premature flip | `liveObservationFlipAllowed: false` in golden; wiring guard |

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_STAGED_WIRING.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_FLIP_CRITERIA.md`
