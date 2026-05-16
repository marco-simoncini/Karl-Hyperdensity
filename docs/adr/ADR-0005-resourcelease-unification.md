# ADR-0005 — Unified ResourceLease (runtime + transfer)

| Field | Value |
|-------|-------|
| **Status** | Accepted |
| **Date** | 2026-05-16 |
| **Sprint** | KHR-B |
| **Supersedes** | Divergence documented in `RESOURCELEASE_SCHEMA_DELTA.md` |

---

## Context

Sprint KHR-A found two contract generations:

1. **JSON Schema (Sprint 91)** — unified workload lease (`shell`, `cell`, `provider`, storage, network).
2. **CRD OpenAPI** — donor/receiver transfer only.

Splitting into **TransferLease** + **RuntimeLease** would duplicate Shell/Cell binding, provider binding, status lifecycle, and Dashboard projection logic.

---

## Decision

Adopt a **single unified `ResourceLease` CRD** with:

| `spec.leaseKind` | Purpose |
|------------------|---------|
| **`runtime`** | Bind Shell + Cell + provider + resources/storage/network/policy |
| **`transfer`** | Move resource envelope between parties via `spec.transfer` |

### Always required (both kinds)

- `spec.leaseKind`
- `spec.shell`
- `spec.cell`
- `spec.provider`

### Kind-specific

| Kind | Additional required |
|------|---------------------|
| `runtime` | `resources`, `storage`, `network`, `policy` |
| `transfer` | `spec.transfer` (`donor`, `receiver`, `resource`, `mode`) |

### Governance (optional both kinds)

`spec.governance`: `dryRunOnly`, `rollbackPlanRef`, `rollbackRequired`, `verificationHooks`, `noRestart`, `guestVisible`, `telemetryConvergedRequired`, `durationSeconds`, `ttlSeconds`.

### Status

Structured `status.phase` (not opaque): `Pending`, `DryRunValidated`, `Bound`, `Active`, `Completing`, `Completed`, `Failed`, `RolledBack`.

### Rejected alternative

**Split TransferLease + RuntimeLease** — rejected: doubles CRD/controller/projection surfaces without reducing complexity.

---

## Consequences

- `resourcelease.schema.json` and `api/crds/.../resourcelease.yaml` **must match** ADR-0005 shape.
- `pkg/khr/crdv1alpha1` keeps **legacy flat transfer fields** as deprecated aliases; parsers accept `spec.transfer` or inline `donor`/`receiver` for v1alpha1 fixtures.
- **No controller** in KHR-B — contract + validation only.
- Dashboard projection remains read-only; observation stub lists Shell/Cell when flag enabled.

---

## Related

- `docs/khr/RESOURCELEASE_LIFECYCLE.md`
- `docs/extraction/RESOURCELEASE_SCHEMA_DELTA.md` (historical)
