# ResourceLease and ResourceFuture Contract (v1alpha1)

**Group:** `hyperdensity.karl.io`  
**Kinds:** `ResourceLease`, `ResourceFuture`  
**Scope:** Namespaced

## ResourceLease

### Intent

A `ResourceLease` is an **atomic market commitment**: transfer or reshape resources between a **donor** and a **receiver**, where both sides are expressed as **`Shell` or `Cell` references** (API-evolvable in later revisions).

### Required fields

- `spec.donor` / `spec.receiver`: `{ kind: Shell|Cell, name, namespace? , apiGroup? }`
- `spec.resource`: `cpu` | `memory` | `disk` | `network` | `gpu`
- `spec.mode`: string (intentionally open in v1alpha1—constrained further once ResourcePort indexing is wired)

### Safety and evidence hooks

| Field | Meaning |
|-------|---------|
| `noRestart` | **Intent** only; KHR must refuse if `ResourcePort` forbids the implied transition. |
| `guestVisible` | Distinguishes host-only envelope vs guest-visible changes. |
| `telemetryConvergedRequired` | Lease cannot close until telemetry gates succeed (future controller). |
| `dryRunOnly` | Simulation / slate membership without execution promotion. |
| `rollbackRequired` / `rollbackPlanRef` | Rollback discipline per ADR-0003. |
| `verificationHooks` | Opaque hook payloads for probes / evidence bundle addresses (preserved unknown fields in v1alpha1). |

### Amount

`spec.amount` is an **opaque object** in v1alpha1 to allow parallel evolution with existing JSON evidence artifacts in this repository. Tighten to `Quantity` or structured `{milliCpu, memoryMi}` in a later API revision.

## ResourceFuture

### Intent

A `ResourceFuture` holds a **prevalidated or tentatively validated** `leaseTemplate` plus **activation** rules (schedule and/or predicates). It is the CR representation of **ResourceFuture** from the architecture fusion document: liquidity **before** pressure becomes outage—**subject to honest capability and cloud constraints**.

### Required fields

- `spec.leaseTemplate`: mirrors `ResourceLease.spec` shape (preserved-unknown in v1alpha1 to avoid premature lock-in).

### Activation

`spec.activation` supports:

- `schedule.notBefore` / `notAfter` (RFC3339 strings)
- `predicates[]` with `type` enum (`resourcePressure`, `sloRisk`, `costBand`, `maintenanceWindow`, `manualApproval`) and opaque `params`.

### Priority and prevalidation

- `priority`: integer hint for slate ordering.
- `prevalidationRequired`: forces dry-run / gate discipline before activation (aligns with Hyperdensity admission patterns).

## Public cloud adaptive mode

On public cloud, many `mode` values map to **provider APIs** or are **unsupported**. Controllers must downgrade or mark `ResourceFuture` / `ResourceLease` as **blocked/remediable** with explicit status reasons (status schema to be tightened post–Sprint 2).

## References

- `docs/architecture/KARL_HYPERDENSITY_KHR_FUSION.md`
- `docs/adr/ADR-0003-hyperdensity-resourcelease-resourcefuture.md`
