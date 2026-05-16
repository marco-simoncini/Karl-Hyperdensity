# Hyperdensity Parent Fabric — workload adapter path-only wiring (Sprint 56)

## Summary

**Sprint 56** was the **first minimal runtime wiring** (path-only). **Sprint 57** adds pilot-only observed-state — see **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_PILOT_OBSERVATION_WIRING.md`**.: Dashboard routes **approved non-apply** path helper call sites through **`hyperdensityWorkloadPath*V1`** wrappers when **`hyperdensityWorkloadAdapterPathWiredV1 = true`**. **Observed-state** remains legacy. **Karl-Hyperdensity** receives **no** Go adapter code. Full **`workload_helpers.go`** verdict remains **`copy-deferred`**.

**Sprint 54 shadow PASS does not imply full production wiring.** `hyperdensityWorkloadAdapterProductionWiredV1` stays **`false`** until all approved phases complete.

---

## 1. Scope

| Item | Sprint 56 |
|------|-----------|
| Path helpers via adapter v1 | **Yes** (approved files only) |
| Observed-state via adapter | **No** |
| Apply / resource-exchange / admission_guard | **Excluded** |
| Dashboard → Hyperdensity `parentfabric` import | **Forbidden** |

---

## 2. Files allowed

- `hyperdensity_parent_fabric_pilot.go`
- `hyperdensity_parent_fabric_live.go`
- `hyperdensity_parent_fabric_vm_linux_cpu_burst_runtime.go`
- `hyperdensity_parent_fabric_vm_linux_cpu_guest_assisted_executor.go`
- `hyperdensity_parent_fabric_vm_linux_memory_guest_assisted.go`
- `hyperdensity_parent_fabric_vm_linux_memory_runtime.go`

Wired families: `appsWorkloadPath`, `virtualMachinePath`, `virtualMachineInstancePath`, `guestOSInfoPath`.

---

## 3. Files explicitly excluded

- `hyperdensity_parent_fabric_apply.go`
- `hyperdensity_parent_fabric_resource_exchange_*`
- `hyperdensity_parent_fabric_admission_guard_*`
- All observed-state / execution / apply mode helpers

---

## 4. Gate / constants

| Constant | Sprint 56 |
|----------|-----------|
| `hyperdensityWorkloadAdapterPathWiredV1` | **`true`** |
| `hyperdensityWorkloadAdapterObservationWiredV1` | **`false`** |
| `hyperdensityWorkloadAdapterProductionWiredV1` | **`false`** |
| `hyperdensityWorkloadAdapterParentFabricImportAllowedV1` | **`false`** |

Wrappers: `hyperdensityWorkloadPathAppsV1`, `hyperdensityWorkloadPathVMV1`, `hyperdensityWorkloadPathVMIV1`, `hyperdensityWorkloadPathGuestOSV1`, etc. — fallback to legacy when adapter returns `ok=false`.

---

## 5. Rollback

1. Set `hyperdensityWorkloadAdapterPathWiredV1 = false`.
2. Revert allowed-file call-site commits (wrappers → direct legacy).
3. Keep adapter + shadow tests.
4. Run parity + `audit_workload_adapter_call_sites.sh`.

---

## 6. Test coverage

| Test / audit | Role |
|--------------|------|
| `TestHyperdensityParentFabricWorkloadAdapterPathWiring` | Wrapper ≡ legacy on samples |
| `TestHyperdensityParentFabricWorkloadAdapterShadow` | Adapter ≡ legacy (unchanged) |
| `TestHyperdensityParentFabricWorkloadAdapterWiringGuard` | No stray adapter instantiation |
| `audit_workload_adapter_call_sites.sh` | Allowed/forbidden file guards |

---

## 7. Risks

| Risk | Mitigation |
|------|------------|
| Subtle path string drift | Shadow + path wiring tests |
| Premature apply wiring | Forbidden file audit |
| Observation wired early | `ObservationWiredV1 = false` + guard |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_WIRING_PROPOSAL.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_PATH_WIRING_M43.md`
