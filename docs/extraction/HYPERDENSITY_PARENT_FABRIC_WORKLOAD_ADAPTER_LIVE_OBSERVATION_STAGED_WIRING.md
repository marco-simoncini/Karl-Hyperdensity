# Hyperdensity Parent Fabric — live observation staged wiring (Sprint 59)

## Summary

**Sprint 59** introduces **staged** live read-only observation wiring in Karl-Dashboard: `hyperdensity_parent_fabric_live.go` calls dedicated wrappers while **`hyperdensityWorkloadAdapterLiveObservationWiredV1 = false`** preserves **legacy fallback** behavior. **Sprint 60** adds shadow hardening. **Sprint 61** flips **`LiveObservationWiredV1 = true`** — see **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_FLIP.md`**.

---

## 1. Scope

| Item | Sprint 59 |
|------|-----------|
| Live observation wrappers in `observation_wiring_v1.go` | **Yes** |
| `live.go` call-site replacement (7 sites) | **Yes** |
| `LiveObservationWiredV1` enabled | **No** (`false`) |
| `ObservationWiredV1` broad flag | **No** (`false`) |
| Hyperdensity Go adapter code | **No** |
| Dashboard → `parentfabric` import | **Forbidden** |

---

## 2. Non-goals

- Activate live adapter behavior (`LiveObservationWiredV1 = true`)
- Set `hyperdensityWorkloadAdapterObservationWiredV1 = true`
- Wire apply, resource_exchange, rollback, VM runtime observation
- Change API responses or JSON ordering
- New ContractKit tag

---

## 3. Allowed file: `hyperdensity_parent_fabric_live.go`

| Wrapper | Legacy delegate (flag false) |
|---------|---------------------------|
| `hyperdensityWorkloadLiveObservedPodUIDV1` | `hyperdensityObservedPodUID` |
| `hyperdensityWorkloadLiveObservedPodCPURequestFromObservationV1` | `hyperdensityObservedPodCPURequestFromObservation` |
| `hyperdensityWorkloadLiveObservedPodCPULimitFromObservationV1` | `hyperdensityObservedPodCPULimitFromObservation` |
| `hyperdensityWorkloadLiveObservedPodMemoryRequestFromObservationV1` | `hyperdensityObservedPodMemoryRequestFromObservation` |
| `hyperdensityWorkloadLiveObservedPodMemoryLimitFromObservationV1` | `hyperdensityObservedPodMemoryLimitFromObservation` |

**Staged call sites:** 7

---

## 4. Explicitly excluded files

- `hyperdensity_parent_fabric_apply.go`
- `hyperdensity_parent_fabric_resource_exchange_*`
- `hyperdensity_parent_fabric_admission_guard_*`
- `hyperdensity_parent_fabric_rollback.go`
- All `hyperdensity_parent_fabric_vm_linux_*` runtime files
- `hyperdensity_parent_fabric_pilot.go` (Sprint 57 pilot wiring unchanged)

---

## 5. Gate / constants

| Constant | Sprint 59 |
|----------|-----------|
| `hyperdensityWorkloadAdapterPathWiredV1` | **`true`** |
| `hyperdensityWorkloadAdapterPilotObservationWiredV1` | **`true`** |
| `hyperdensityWorkloadAdapterLiveObservationWiredV1` | **`false`** |
| `hyperdensityWorkloadAdapterObservationWiredV1` | **`false`** |
| `hyperdensityWorkloadAdapterProductionWiredV1` | **`false`** |

---

## 6. Fallback behavior

When **`LiveObservationWiredV1 == false`**, every live wrapper calls the **same legacy helper** as before Sprint 59. Runtime output is **bit-for-bit equivalent** for the staged paths.

When a future sprint sets **`LiveObservationWiredV1 = true`**, the `true` branch is structured for adapter substitution after shadow tests pass.

---

## 7. Rollback

1. Revert `live.go` call sites to direct legacy helpers.
2. Remove live wrappers from `observation_wiring_v1.go`.
3. Remove `hyperdensityWorkloadAdapterLiveObservationWiredV1` constant.

---

## 8. Test coverage

| Artifact | Owner |
|----------|-------|
| `hyperdensity_parent_fabric_workload_adapter_live_observation_staged_test.go` | Dashboard |
| `testdata/..._live_observation_staged.golden.json` | Dashboard |
| `audit_workload_adapter_call_sites.sh` (Sprint 59 guards) | Dashboard |

---

## 9. Risks

| Risk | Mitigation |
|------|------------|
| Accidental `LiveObservationWiredV1 = true` | Audit + wiring guard + golden manifest |
| Live wrappers spread beyond `live.go` | Scope audit + wiring guard |
| False sense of completion | `ObservationWiredV1` remains **false** until all phases done |

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_PROPOSAL.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_PILOT_OBSERVATION_HARDENING.md`
