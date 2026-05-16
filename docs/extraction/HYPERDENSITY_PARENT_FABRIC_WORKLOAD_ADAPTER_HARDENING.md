# Hyperdensity Parent Fabric — workload adapter hardening (Sprint 53)

## Summary

**Sprint 53** introduces a **production-internal** workload adapter in **Karl-Dashboard** (`hyperdensity_parent_fabric_workload_adapter_v1.go`). **Karl-Hyperdensity** receives **no** Go adapter code. **`parentfabric/workload`** retains only the **three pure-candidate** copy-contract (Sprint 52). Full **`workload_helpers.go`** verdict remains **`copy-deferred`**.

## What Sprint 53 delivers (Dashboard)

| Item | Status |
|------|--------|
| `hyperdensityWorkloadPathAdapterV1` | Production `.go` — delegates to legacy `hyperdensity*APIPath` |
| `hyperdensityWorkloadObservationAdapterV1` | Production `.go` — delegates to legacy observed-state builders |
| Adapter constants | `hyperdensityWorkloadAdapterVersionV1`, `productionWired=false`, `parentFabricImportAllowed=false` |
| Dedicated tests | `hyperdensity_parent_fabric_workload_adapter_v1_test.go` |
| Golden manifest | `hyperdensity_parent_fabric_workload_adapter_v1.golden.json` |

## What is explicitly **not** done

| Item | Status |
|------|--------|
| Handler / call-site wiring | **No** — `hyperdensityWorkloadAdapterProductionWiredV1 = false` |
| Dashboard import of `pkg/hyperdensity/parentfabric` | **Forbidden** |
| API response / payload change | **No** |
| Hyperdensity adapter Go code | **No** |
| Full `workload_helpers.go` copy | **Still deferred** |

## Sprint timeline

| Sprint | Scope |
|--------|--------|
| **50** | Adapter boundary classification |
| **51** | Test-only adapter stubs (`*_test.go`) |
| **52** | Three pure-candidates in Hyperdensity `parentfabric/workload` |
| **53** | Dashboard production-internal adapter — **not wired** |
| **54** | Dashboard shadow tests (legacy vs adapter v1) — **still not wired** |
| **55** | Wiring proposal + call-site inventory |
| **56** | Path-only wiring via wrappers (Dashboard) |
| **57** | Pilot-only observation wiring (Dashboard) |
| **Future** | Broader observation wiring |

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_STUB_READINESS.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_PURE_CANDIDATES_CONTRACT.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_HARDENING_M39.md`
