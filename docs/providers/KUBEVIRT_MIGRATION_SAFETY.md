# KubeVirt Legacy — Migration and Safety

**Sprint 3 — governance document.** No operational changes in this milestone.

## Invariants

1. **No KubeVirt removal** — all ISO / installer / Dashboard paths that install or assume KubeVirt stay unchanged in this sprint.
2. **No VM behavior change** — CPU, memory, disk, migration, and guest OS behavior are unchanged; Shell/Cell documents are **overlay** semantics.
3. **Read-only mapping first** — initial integration is **observation and correlation**: Shell/Cell objects may reference existing VMs via `kubeVirtLegacy` hints and `runtimeHandle` without mutating VM specs.
4. **Labels opt-in later** — `karl.io/*` and `hyperdensity.karl.io/legacy-provider` labels are **not** applied by any shipping component in Sprint 3; they are **reserved** for a gated rollout (feature flag / canary pool).
5. **Rollback semantics** — removing KARL labels or unlinking a Shell must **never** delete or stop a VM; rollback is **detach overlay**, not destructive teardown.
6. **Dual-model risk** — operators may see both native KubeVirt objects and KARL Shell/Cell CRs; documentation and UI must avoid conflicting “source of truth” until ownership mode is explicitly promoted in a future ADR.

## Migration Factory

`marco-simoncini/Karl-Migration-Factory` remains the strategic bridge for attested moves from external estates. This contract **does not** replace Migration Factory workflows; it provides **stable KARL-side identifiers** (`Shell`, `Cell`, handle) that Migration Factory can emit or consume in later integration phases.

## Hyperdensity interaction

See `KUBEVIRT_LEGACY_PROVIDER_CONTRACT.md` and `examples/providers/kubevirt/resourcelease-kubevirt-guarded-example.yaml`. Under `evidenceLevel: provider-constrained`, Hyperdensity must prefer **recommendation**, **dry-run**, and **evidence requests** over blind apply until KHR or provider-specific executors certify paths.

## Example bundle overlap

`examples/crds/` and `examples/providers/kubevirt/` both ship a `RuntimeProvider` named `kubevirt-legacy-v1`. Applying the entire repository example set without deduplication will fail on second create. Prefer **provider-scoped** examples when working on KubeVirt legacy contracts.

## Escalation path

If a lease would imply guest-visible mutation unsupported by KubeVirt for a given VM class, status must move to **blocked/remediable** (future controller) rather than partial apply.
