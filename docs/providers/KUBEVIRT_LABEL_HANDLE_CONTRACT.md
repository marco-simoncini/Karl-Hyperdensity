# KubeVirt Legacy — Label and Handle Contract

**Sprint 3 — specification only.** No reconciler applies these labels in this sprint.

## A) KARL keys (future opt-in reconciler)

These keys are **reserved** for mappers from Shell/Cell to KubeVirt-owned objects:

| Key | Value example | Semantics |
|-----|----------------|-----------|
| `karl.io/shell-id` | `finance-desktop-legacy-01` | Stable Shell identity (namespaced uniqueness is Shell metadata). |
| `karl.io/cell-id` | `fin-vm-legacy-01-cell` | Stable Cell identity. |
| `karl.io/runtime-provider` | `kubevirt-legacy` | Must match `providerType` in legacy provider contract. |
| `karl.io/provider-handle-kind` | `VirtualMachine` | API kind of primary handle (or `VirtualMachineInstance` when VMI-only). |
| `karl.io/provider-handle-namespace` | `karl-sandbox` | Namespace of primary object. |
| `karl.io/provider-handle-name` | `fin-vm-legacy-01` | Name of primary object. |
| `hyperdensity.karl.io/evidence-scope` | `vm-linux` / `vm-windows` / … | Hyperdensity evidence partition (aligns with Dashboard lanes). |
| `hyperdensity.karl.io/legacy-provider` | `true` | Marks objects participating in legacy provider evidence paths. |

**Placement:** Typically on `VirtualMachine`, optionally mirrored on `VirtualMachineInstance` when policy requires — exact placement is a **Sprint 4+** controller decision.

## B) Upstream KubeVirt labels (preserve)

Existing discovery and Hyperdensity pilots rely on:

- `vm.kubevirt.io/name` — canonical VM name reference on VMI and related objects where present.
- `vmi.kubevirt.io/id` — VMI identifier surface used in parent-fabric style discovery.

**Rule:** KARL controllers **must not** remove or rewrite these labels. KARL may **add** keys under §A without colliding with KubeVirt reserved prefixes.

## C) `Cell.status.runtimeHandle` (contract shape)

Populated in future reconciliation; schema defined in `api/providers/kubevirt/kubevirt-handle-contract.yaml`.

Minimum logical fields:

- `provider`: `kubevirt-legacy`
- `apiGroup`: `kubevirt.io`
- `kind`: `VirtualMachine` (primary) or `VirtualMachineInstance`
- `namespace`, `name`, `uid`
- `vmiName` when distinct from VM controller naming
- `launcherPodName` when observable and stable enough for evidence (optional; provider-constrained)

**Note:** UID is required for idempotent bind/unbind and audit correlation with Hyperdensity evidence bundles.
