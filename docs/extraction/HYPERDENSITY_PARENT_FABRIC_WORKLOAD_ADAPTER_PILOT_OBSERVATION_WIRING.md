# Hyperdensity Parent Fabric — pilot observed-state adapter wiring (Sprint 57)

## Summary

**Sprint 57** enables **pilot-only** observed-state wiring: `hyperdensity_parent_fabric_pilot.go` calls **`hyperdensityWorkloadPilotObservedStateV1`** when **`hyperdensityWorkloadAdapterPilotObservationWiredV1 = true`**. **General** observation wiring remains **`false`**. **Sprint 56** path wiring unchanged. Full **`workload_helpers.go`** verdict remains **`copy-deferred`**.

---

## 1. Scope

| Item | Sprint 57 |
|------|-----------|
| Pilot `FromWorkload` → adapter wrapper | **Yes** |
| Observed-state in apply/live/VM/rollback | **No** (legacy) |
| Hyperdensity Go adapter code | **No** |
| Dashboard → `parentfabric` import | **Forbidden** |

---

## 2. Allowed file

- **`hyperdensity_parent_fabric_pilot.go`** only (single `hyperdensityPilotObservedStateFromWorkload` call site in `hyperdensityPilotObservedStateForPlan`)

---

## 3. Explicitly excluded files

- `hyperdensity_parent_fabric_apply.go`
- `hyperdensity_parent_fabric_resource_exchange_*`
- `hyperdensity_parent_fabric_admission_guard_*`
- `hyperdensity_parent_fabric_live.go`
- `hyperdensity_parent_fabric_rollback.go`
- All `hyperdensity_parent_fabric_vm_linux_*` runtime files
- Direct observed pod helpers in pilot (post-wrapper enrichment unchanged)

---

## 4. Gate / constants

| Constant | Sprint 57 |
|----------|-----------|
| `hyperdensityWorkloadAdapterPilotObservationWiredV1` | **`true`** |
| `hyperdensityWorkloadAdapterObservationWiredV1` | **`false`** |
| `hyperdensityWorkloadAdapterPathWiredV1` | **`true`** (Sprint 56) |
| `hyperdensityWorkloadAdapterProductionWiredV1` | **`false`** |

---

## 5. Fallback behavior

Wrapper uses `extractWorkloadObservedStateAt` when pilot flag true; on `ok=false` falls back to **`hyperdensityPilotObservedStateFromWorkload`**. Adapter delegates to same legacy helpers as Sprint 54 shadow.

---

## 6. Rollback

1. Set `hyperdensityWorkloadAdapterPilotObservationWiredV1 = false`.
2. Revert `pilot.go` to direct `hyperdensityPilotObservedStateFromWorkload` call.
3. Run parity + audit script.

---

## 7. Test coverage

| Test / audit | Role |
|--------------|------|
| `TestHyperdensityParentFabricWorkloadAdapterPilotObservationWiring` | Wrapper ≡ legacy (Deployment/StatefulSet) |
| `TestHyperdensityParentFabricWorkloadAdapterShadow` | Unchanged baseline |
| `TestHyperdensityParentFabricWorkloadAdapterWiringGuard` | Adapter instantiation guards |
| `audit_workload_adapter_call_sites.sh` | Pilot-only wrapper enforcement |

---

## 8. Risks

| Risk | Mitigation |
|------|------------|
| Pod enrichment after wrapper differs | Only `FromWorkload` replaced; pod helpers unchanged |
| `time.Now()` in non-pilot adapter methods | Pilot wrapper uses explicit `now` parameter |
| Premature general observation wire | `ObservationWiredV1` stays false + audit |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_PATH_WIRING.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_PILOT_OBSERVATION_WIRING_M44.md`
