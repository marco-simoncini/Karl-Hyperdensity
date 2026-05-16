# Hyperdensity Parent Fabric — pilot observation hardening (Sprint 58)

## Summary

**Sprint 58** hardens **Sprint 57** pilot-only observation wiring with stronger end-to-end tests and documentation. **Sprint 59** stages live observation wrappers (flag false). Full **`workload_helpers.go`** verdict remains **`copy-deferred`**.

---

## 1. Scope

| Item | Sprint 58 |
|------|-----------|
| New production wiring | **No** |
| Pilot observation hardening (tests + docs) | **Yes** |
| Live observation inventory / proposal | **Yes** (docs only) |
| Hyperdensity Go adapter code | **No** |
| Dashboard → `parentfabric` import | **Forbidden** |

---

## 2. Non-goals

- Wire live observation in `hyperdensity_parent_fabric_live.go`
- Set `hyperdensityWorkloadAdapterObservationWiredV1 = true`
- Touch `apply.go`, `resource_exchange_*`, `admission_guard_*`, rollback, VM runtime
- Change API responses or JSON ordering
- Move runtime code from Dashboard to Hyperdensity
- New ContractKit tag (remains **v0.1.9-khr-m1-m19**)

---

## 3. Pilot observed-state hardening matrix

| Case | Path | Expected |
|------|------|----------|
| Deployment workload-template | `hyperdensityPilotObservedStateForPlan` | ≡ legacy reference via wrapper |
| StatefulSet workload-template | same | ≡ reference |
| VirtualMachine / VM live-update | VM legacy builder | **Not** workload-only wrapper |
| Pod resize enrichment | pod_resize_subresource | enrichment unchanged |
| Unknown kind | legacy fallback | preserved |

---

## 4. Live observation readiness, proposal-only

See **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_PROPOSAL.md`**. Sprint 58 inventories `live.go` call sites only.

---

## 5. Rollback

Same as Sprint 57: set `hyperdensityWorkloadAdapterPilotObservationWiredV1 = false` and revert `pilot.go` call site to `hyperdensityPilotObservedStateFromWorkload`.

---

## 6. Test coverage

| Repo | Artifact |
|------|----------|
| Dashboard | `hyperdensity_parent_fabric_workload_adapter_pilot_observation_hardening_test.go` |
| Dashboard | `testdata/hyperdensity_parent_fabric_workload_adapter_pilot_observation_hardening.golden.json` |
| Dashboard | Updated `audit_workload_adapter_call_sites.sh` |

---

## 7. Risks

| Risk | Mitigation |
|------|------------|
| False confidence from wrapper-only tests | End-to-end `ForPlan` tests in Sprint 58 |
| Premature live wiring | Inventory + explicit `LiveObservationWiredV1` flag reserved for future sprint |
| Broad `ObservationWiredV1` flip | Remains **false** until all observation phases complete |

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_PILOT_OBSERVATION_WIRING.md` (Sprint 57)
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_PROPOSAL.md` (Sprint 58)
