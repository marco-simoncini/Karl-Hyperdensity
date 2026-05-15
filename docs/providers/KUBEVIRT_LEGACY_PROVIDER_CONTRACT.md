# KubeVirt Legacy Runtime Provider — Contract (Sprint 3)

**Status:** Contract and documentation only. No controllers, no ISO/installer changes, no VM mutation.  
**Repository:** `marco-simoncini/Karl-Hyperdensity`  
**Branch:** `KHR`

## Positioning

KubeVirt is a **compatibility / legacy RuntimeProvider** under the KARL Shell/Cell model (ADR-0002). It is **not** the primary hotplug or economic story; **Hyperdensity** and future **KHR-native** paths own differentiated value. KubeVirt remains installed and authoritative for VM lifecycle exactly as today until explicit migration milestones.

## Declarative contract bundle

Machine-readable stubs live under `api/providers/kubevirt/`:

| File | Purpose |
|------|---------|
| `kubevirt-legacy-provider.yaml` | `RuntimeProvider` identity and policy fields for `kubevirt-legacy` |
| `kubevirt-label-contract.yaml` | Label and annotation keys KARL may emit in future **opt-in** reconcilers |
| `kubevirt-handle-contract.yaml` | Shape of `Cell.status.runtimeHandle` for KubeVirt bindings |

## RuntimeProvider — required semantics

The legacy provider **declares** (see `kubevirt-legacy-provider.yaml`):

| Field | Value | Meaning |
|-------|-------|---------|
| `providerType` | `kubevirt-legacy` | Stable type discriminator for controllers and RBAC. |
| `maturity` | `legacy` | No hero positioning; compatibility track. |
| `ownershipMode` | `hybrid-legacy` | KARL Shell/Cell overlay coexists with native KubeVirt objects; no forced ownership flip in Sprint 3+. |
| `supportedCellTypes` | `kubevirt-legacy` | Cell kind discriminator for materialization contracts. |
| `evidenceLevel` | `provider-constrained` | Hyperdensity evidence must respect KubeVirt API and guest limits; no universal live claims. |
| `resourcePortTemplateRef` | `kubevirt-legacy` | Default `ResourcePort` profile name for this provider (see `resourceport-kubevirt-legacy` example). |

## Relationship to existing product code (read-only alignment)

- `marco-simoncini/Karl-Dashboard` Hyperdensity control-plane docs reference **`vm.kubevirt.io/name`** and **`vmi.kubevirt.io/id`** for discovery — this contract **preserves** those labels as first-class **upstream** identifiers while introducing **parallel KARL keys** for future mappers.
- `marco-simoncini/Karl-OS-ISO` / `Karl-Installer` continue to own install matrices; this sprint **does not** change manifests or versions.

## Non-goals (Sprint 3)

- No production controller reconciliation.
- No apply path that mutates `VirtualMachine`, `VirtualMachineInstance`, or launcher pods.
- No removal or downgrade of KubeVirt components.
- No change to guest-visible behavior of existing VMs.

## References

- `docs/providers/KUBEVIRT_SHELL_CELL_MAPPING.md`
- `docs/providers/KUBEVIRT_LABEL_HANDLE_CONTRACT.md`
- `docs/providers/KUBEVIRT_MIGRATION_SAFETY.md`
- `docs/adr/ADR-0002-kubevirt-as-legacy-provider.md`

## Hyperdensity and `ResourceLease` on kubevirt-legacy

Under `evidenceLevel: provider-constrained`, Hyperdensity / Grande Padre interactions with KubeVirt-backed Cells are limited to:

1. **Guarded recommendation** — surface suggested CPU/RAM/disk moves with explicit dependency on KubeVirt and guest capability.
2. **Dry-run** — `ResourceLease.spec.dryRunOnly: true` (or slate-level dry-run) produces projected state without API mutation.
3. **Evidence requests** — ask for guest-assisted probes, `guestosinfo`, or other **provider-scoped** artifacts before promotion.
4. **Apply** — real mutation remains **legacy/provider-constrained** until a certified KHR or KubeVirt-specific executor is shipped; Sprint 3 does **not** implement apply.

No universal hotplug promise: even when `ResourcePort` lists a mode, KubeVirt and guest truth may **narrow** effective modes per instance.
