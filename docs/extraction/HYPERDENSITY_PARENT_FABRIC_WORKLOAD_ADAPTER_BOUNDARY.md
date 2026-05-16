# Hyperdensity Parent Fabric — workload adapter boundary (Sprint 50)

## Purpose

Define what **must remain in Karl-Dashboard** (adapter-bound / runtime-bound) vs what **may** eventually live in **Karl-Hyperdensity** pure-core when `hyperdensity_parent_fabric_workload_helpers.go` is re-audited.

| Sprint | Scope |
|--------|--------|
| **50** | Boundary docs + Dashboard classification fixture — no Go adapter |
| **51** | Dashboard **test-only** adapter stubs (`*_test.go`) — **no production wiring** |
| **52** | Three pure-candidates in Hyperdensity `parentfabric/workload` |
| **53** | Dashboard **production-internal** adapter v1 — **not wired** |
| **54** | Shadow tests (legacy vs adapter) — **still not wired** |
| **55** | Wiring proposal + call-site inventory |
| **56** | Path-only wiring (6 approved files) — `PathWiredV1=true` |
| **Future** | Observation wiring sprint |

## Adapter-bound (stays in Dashboard)

| Class | Examples in `workload_helpers.go` | Why |
|-------|-----------------------------------|-----|
| **API path builders** | `hyperdensityAppsWorkloadAPIPath`, `hyperdensityPodAPIPath`, `hyperdensityVirtualMachine*APIPath` | K8s/KubeVirt REST path strings |
| **KubeVirt / K8s path literals** | `/apis/kubevirt.io/…`, `/apis/apps/v1/…` | Product retains KubeVirt; paths are not pure-core |
| **Observed-state from live objects** | `hyperdensityObservedPod*`, `hyperdensityPilotObservedStateFrom*` | Untyped `map[string]interface{}` from API responses |
| **Pilot observed state** | `HyperdensityPilotObservedState` assembly | Runtime coupling to console server types |
| **Runtime / apply mode selection** | `hyperdensityExecutionModeForKind`, `hyperdensityExecutionReadyReasonForMode`, `hyperdensityExecutionMechanismForMode` | Apply orchestration surface — **do not move** |

## May become pure-core in Hyperdensity (future)

| Class | Current home | Notes |
|-------|--------------|-------|
| **Nested map access** | `parentfabric/primitives` | `StringAt`, `MapAt`, … — stdlib contract |
| **Quantity normalization** | `parentfabric/primitives` | Narrow contract; not full `resource.ParseQuantity` |
| **Pure DTOs** | `parentfabric/executiontypes` (partial) | Copy-contract slices only |
| **Pure enums / kind tables** | Future, after adapter strips runtime refs | Only when no Dashboard `server` symbols |

## Rule for `parentfabric/workload`

`pkg/hyperdensity/parentfabric/workload` may receive **only** functions that:

- Do **not** build API paths (K8s/KubeVirt/HTTP).
- Do **not** read live object maps tied to Dashboard server responses.
- Do **not** reference execution/apply action constants from runtime packages.

Until re-audit criteria pass, **`workload` stays placeholder-only** (`doc.go`).

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_INTERFACE_PROPOSAL.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REAUDIT_CRITERIA.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_BOUNDARY_M35.md`
