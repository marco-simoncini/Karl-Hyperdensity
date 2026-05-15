# KubeVirt Legacy — Shell / Cell / KubeVirt Mapping

**Sprint 3 — contract only**

## Layered model

| Layer | Responsibility | KubeVirt legacy role |
|-------|------------------|----------------------|
| **Shell** | Product / user / billing / policy identity | Declares *intent* (desktop, app-like VM experience, Linux VM workspace) and optional `spec.kubeVirtLegacy` **hints** for binding. |
| **Cell** | Runtime abstraction on a node / cluster slice | Holds **provider handle** pointing at KubeVirt `VirtualMachine` / `VirtualMachineInstance` (and optional derived ids). |
| **KubeVirt** | Implementation detail | Source of truth for VM process, disks, migration, and upstream labels today. |

## Direction of truth (initial / read-only phase)

1. **Existing VMs** remain owned by KubeVirt APIs and customer workflows.  
2. **Shell** may be introduced as a **read model** or **parallel declaration** without changing VM spec.  
3. **Cell** materialization in a future controller links Shell ↔ existing or new VM objects via **handle** + **labels** (opt-in).

## Mapping table (conceptual)

| Shell kind (example) | ShellClass `complianceFamily` (from CRD enum) | Typical Cell `spec.runtimeProviderRef` | KubeVirt object |
|----------------------|-----------------------------------------------|------------------------------------------|-----------------|
| Windows desktop on legacy VM | `windows.session` | `kubevirt-legacy-v1` | `VirtualMachine` + `VirtualMachineInstance` |
| Linux VM workspace | `linux.kubevirt.legacy.vm` | `kubevirt-legacy-v1` | `VirtualMachine` + `VirtualMachineInstance` |
| Windows app (future) | `windows.app` | *Not KubeVirt primary* | N/A for RemoteApp path |

## Handle summary

`Cell.status.runtimeHandle` (future population) carries normalized reference to KubeVirt API objects. See `api/providers/kubevirt/kubevirt-handle-contract.yaml` and `docs/providers/KUBEVIRT_LABEL_HANDLE_CONTRACT.md`.

## Compatibility with Sprint 2 CRDs

Examples under `examples/providers/kubevirt/` mirror `runtime.karl.io/v1alpha1` and `hyperdensity.karl.io/v1alpha1` shapes from `api/crds/` without requiring new CRD versions in Sprint 3.
