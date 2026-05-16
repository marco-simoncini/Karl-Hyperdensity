# Hyperdensity Parent Fabric — live observation flip (Sprint 61)

## Summary

**Sprint 61** sets **`hyperdensityWorkloadAdapterLiveObservationWiredV1 = true`** for staged wrappers in **`hyperdensity_parent_fabric_live.go`**. The **true branch remains legacy-equivalent** (same delegate helpers). **Sprint 62** semantic candidate shadow. **Sprint 63** branch swap to candidate — see **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_BRANCH_SWAP.md`**. **`ObservationWiredV1`** and **`ProductionWiredV1`** stay **`false`**.

---

## 1. Scope

| Item | Sprint 61 |
|------|-----------|
| Flip `LiveObservationWiredV1` | **`true`** |
| Live wrappers in `live.go` only | **Yes** (7 call sites, unchanged) |
| Semantic adapter change | **No** (legacy delegate in true branch) |
| Broad `ObservationWiredV1` | **`false`** |

---

## 2. Non-goals

- Set `hyperdensityWorkloadAdapterObservationWiredV1 = true`
- Set `hyperdensityWorkloadAdapterProductionWiredV1 = true`
- Wire apply, resource_exchange, rollback, VM runtime observation
- Dashboard import of `pkg/hyperdensity/parentfabric`
- Change API responses or JSON ordering

---

## 3. Preconditions satisfied from Sprint 60

| Criterion | Status |
|-----------|--------|
| Live shadow hardening PASS (7 cases) | **Yes** |
| Call-site audit PASS | **Yes** |
| `liveObservationFlipAllowed: true` in golden | **Yes** (Sprint 61) |
| Dedicated flip sprint | **Yes** |

---

## 4. Gate / constants

| Constant | Sprint 61 |
|----------|-----------|
| `hyperdensityWorkloadAdapterLiveObservationWiredV1` | **`true`** |
| `hyperdensityWorkloadAdapterObservationWiredV1` | **`false`** |
| `hyperdensityWorkloadAdapterProductionWiredV1` | **`false`** |
| `hyperdensityWorkloadAdapterPilotObservationWiredV1` | **`true`** |
| `hyperdensityWorkloadAdapterPathWiredV1` | **`true`** |

---

## 5. Behavior equivalence statement

With flag **`true`**, live wrappers execute the **same legacy helpers** as the **`false`** branch (Sprint 59–60). Runtime output for staged paths is **unchanged**; Sprint 61 validates **governance and gate wiring**, not a new observation implementation.

---

## 6. Rollback

Set `hyperdensityWorkloadAdapterLiveObservationWiredV1 = false` in `adapter_v1.go`; revert shadow/flip goldens; re-run parity.

---

## 7. Test coverage

| Artifact | Owner |
|----------|-------|
| `hyperdensity_parent_fabric_workload_adapter_live_observation_flip_test.go` | Dashboard |
| Updated shadow test + goldens | Dashboard |
| `testdata/..._live_observation_flip.golden.json` | Dashboard |

---

## 8. Risks

| Risk | Mitigation |
|------|------------|
| Flip perceived as behavior change | Document legacy-equivalent true branch |
| Premature broad observation | `ObservationWiredV1` remains **false** |

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_FLIP_CRITERIA.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_SHADOW_HARDENING.md`
