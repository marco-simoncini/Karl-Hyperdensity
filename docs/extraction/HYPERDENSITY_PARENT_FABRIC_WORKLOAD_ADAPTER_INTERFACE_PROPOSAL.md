# Hyperdensity Parent Fabric — workload adapter interface proposal (Sprint 50)

## Status

**Document-only.** No Go `interface` types in `parentfabric` packages, no Dashboard implementation, no runtime wiring, no copy from `workload_helpers.go` in Sprint 50.

## WorkloadPathAdapter (conceptual)

Responsible for **all** K8s/KubeVirt REST path construction currently in `workload_helpers.go`.

| Method (conceptual) | Dashboard function today |
|---------------------|--------------------------|
| `AppsWorkloadPath(kind, namespace, name string) (string, error)` | `hyperdensityAppsWorkloadAPIPath` |
| `PodPath(namespace, name string) string` | `hyperdensityPodAPIPath` |
| `PodResizePath(namespace, name string) string` | `hyperdensityPodResizeAPIPath` |
| `VirtualMachinePath(namespace, name string) string` | `hyperdensityVirtualMachineAPIPath` |
| `VirtualMachineInstancePath(namespace, name string) string` | `hyperdensityVirtualMachineInstanceAPIPath` |
| `GuestOSInfoPath(namespace, vmiName string) string` | `hyperdensityVirtualMachineInstanceGuestOSInfoAPIPath` |

**Owner:** Karl-Dashboard adapter package (future sprint), behind explicit allowlist.

## WorkloadObservationAdapter (conceptual)

Responsible for **observed-state** extraction from untyped API object maps.

| Method (conceptual) | Dashboard functions today |
|---------------------|---------------------------|
| `ExtractWorkloadObservedState(kind string, obj map[string]interface{}, containerName string, now time.Time) (HyperdensityPilotObservedState, error)` | `hyperdensityPilotObservedStateFromWorkload` |
| `ExtractStatefulSetObservedState(obj map[string]interface{}, containerName string, now time.Time) (HyperdensityPilotObservedState, error)` | `hyperdensityPilotObservedStateFromStatefulSet` |
| `ExtractPodCPURequest(pod map[string]interface{}, containerName string) (HyperdensityCPUQuantity, error)` | `hyperdensityObservedPodCPURequest` (+ siblings) |

Return types remain **Dashboard DTOs** until a dedicated pure DTO migration sprint.

## What is NOT in these adapters (Sprint 50)

- Execution mode / ready-reason / mechanism selection (remain **do-not-move** in Dashboard).
- The three **pure-candidate** kind helpers (`hyperdensityAppsWorkloadResource`, …) — candidates for Hyperdensity `primitives` or `workload` **after** adapter boundary is implemented.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REAUDIT_CRITERIA.md`
