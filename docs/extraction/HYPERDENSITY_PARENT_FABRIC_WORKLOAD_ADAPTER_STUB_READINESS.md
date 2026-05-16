# Hyperdensity Parent Fabric — workload adapter stub readiness (Sprint 51)

## Summary

**Sprint 51** introduces **test-only adapter stubs** in **Karl-Dashboard** only. **Karl-Hyperdensity** receives **no** Go adapter code. **`parentfabric/workload`** remains a **placeholder** (`doc.go`). **`hyperdensity_parent_fabric_workload_helpers.go`** remains **`copy-deferred`**.

## What Sprint 51 delivers

| Item | Location | Production wiring |
|------|----------|-------------------|
| `testWorkloadPathAdapter` / `testWorkloadObservationAdapter` | Dashboard `*_test.go` | **No** |
| `testDashboardWorkloadPathAdapter` | delegates to existing path builders in tests | **No** |
| `testDashboardWorkloadObservationAdapter` | stub delegates + minimal fixtures in tests | **No** |
| Golden manifest | `hyperdensity_parent_fabric_workload_adapter_stub.golden.json` | **No** |

## What Hyperdensity does **not** receive

- No `WorkloadPathAdapter` Go interface in `pkg/hyperdensity/parentfabric/workload`
- No K8s/KubeVirt path literals copied from Dashboard
- No observed-state builders copied from Dashboard
- No production import from Dashboard runtime

## Sprint 52 note

Three **pure-candidate** functions are now copy-contracted in Hyperdensity `parentfabric/workload`. Adapter stubs and full `workload_helpers.go` remain **copy-deferred**.

## Re-audit gate (unchanged from Sprint 50 criteria)

`workload_helpers.go` may be **re-audited** only when:

1. Dashboard adapter stub tests **PASS** (Sprint 51 baseline)
2. Adapter boundary classification remains complete (Sprint 50)
3. Pure allowlist narrowed to **3** functions
4. `parentfabric/primitives` contract stable
5. Golden tests available for any copied pure-core slice
6. **Explicit** sprint approves copy — not Sprint 51

## Behavior guarantees

- **No** API response change
- **No** Parent Fabric runtime behavior change
- **No** JSON ordering change
- **No** execution/apply changes
- **No** `pkg/hyperdensity/parentfabric` import in Dashboard production `.go`

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_INTERFACE_PROPOSAL.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REAUDIT_CRITERIA.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_STUB_M37.md`
