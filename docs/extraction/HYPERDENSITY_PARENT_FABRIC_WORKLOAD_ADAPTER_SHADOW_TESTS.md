# Hyperdensity Parent Fabric — workload adapter shadow tests (Sprint 54)

## Summary

**Sprint 54** adds **Dashboard-only shadow tests** that compare **legacy helper output** vs **adapter v1** on representative path and observed-state samples. **Sprint 58** adds end-to-end pilot observation hardening tests (`hyperdensityPilotObservedStateForPlan`) — no new wiring beyond Sprint 57. **Karl-Hyperdensity** receives **no** Go adapter code. Adapter v1 remains **not wired** (`hyperdensityWorkloadAdapterProductionWiredV1 = false`). Full **`workload_helpers.go`** verdict remains **`copy-deferred`**.

## What Sprint 54 delivers

| Item | Location |
|------|----------|
| Shadow parity test | `hyperdensity_parent_fabric_workload_adapter_shadow_test.go` |
| Accidental wiring guard | `TestHyperdensityParentFabricWorkloadAdapterNotWired` |
| Golden manifest | `hyperdensity_parent_fabric_workload_adapter_shadow.golden.json` |

## Shadow scope

| Category | Cases | Comparison |
|----------|------:|------------|
| **Path** | 7 | Legacy `hyperdensity*APIPath` vs `hyperdensityWorkloadPathAdapterV1` |
| **Observed-state** | 3 | Legacy builders vs `extract*ObservedStateAt` / `ExtractPodObservedSnapshot` |

Fixed time for observed-state: `2026-05-16T12:00:00Z`.

## What is explicitly **not** done

| Item | Status |
|------|--------|
| Handler / call-site wiring | **No** |
| Dashboard import of `pkg/hyperdensity/parentfabric` | **Forbidden** |
| API response / payload / runtime behavior change | **No** |
| Hyperdensity adapter Go code | **No** |
| Full `workload_helpers.go` copy | **Still deferred** |

## Sprint timeline

| Sprint | Scope |
|--------|--------|
| **53** | Production-internal adapter v1 — not wired |
| **54** | Shadow tests (legacy vs adapter) — still not wired |
| **55** | Wiring **proposal** + call-site inventory |
| **56** | **Path-only** wiring (approved non-apply files) |
| **57** | **Pilot-only** observed-state wiring — general observation **not** wired |
| **Future** | Broader observation wiring (non-pilot files) |

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_HARDENING.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_SHADOW_TESTS_M40.md`
