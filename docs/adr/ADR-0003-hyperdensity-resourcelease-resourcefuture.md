# ADR-0003 — Hyperdensity uses ResourceLease and ResourceFuture as resource market primitives

| Field | Value |
|-------|-------|
| **Status** | Accepted (architecture foundation — Sprint 1) |
| **Date** | 2026-05-15 |
| **Applies to** | `marco-simoncini/Karl-Hyperdensity`, future `KHR`, `Karl-Dashboard` runtime extraction |

---

## Context

`marco-simoncini/Karl-Hyperdensity` already encodes market-like concepts in code and JSON artifacts:

- **Donor / receiver** dyads appear in canary cohort and movement reference surfaces (`donorShellId`, `receiverShellId`).
- **Execution leases** and **rollback windows** appear in production canary cohort schemas and validators.
- Packages include `leaseslate`, `marketcontroller`, `applygate`, `autoapply`, `guardedauto`, and dry-run admission patterns—evidence of movement toward **contractual transfers** rather than ad-hoc patches.

However, **ResourceLease** and **ResourceFuture** are not yet universally exposed as stable, KHR-consumable APIs across all execution paths.

Hyperdensity’s README positions the system as **live, evidence-driven CPU and RAM control** with strict non-disruptive guarantees for Linux shells—this ADR generalizes the **verb model** for all providers that opt in honestly.

---

## Decision

1. Standardize on **ResourceLease** as the atomic **committed transfer** or **envelope mutation** instruction between parties (donor/receiver or intra-Shell plane shifts) with rollback and verification hooks.
2. Standardize on **ResourceFuture** as the **conditional / scheduled** commitment that becomes a ResourceLease when predicates and gates pass—powering **predictive prevalidation**.
3. Standardize on **ResourcePort** as the **capability truth** gate: no lease or future may be scheduled if ports forbid the mode (see `KARL_HYPERDENSITY_KHR_FUSION.md`).
4. Position **Grande Padre** as the **decisioning** authority that produces and graduates leases/futures; **KHR** as the **execution** authority that applies them on nodes (once implemented).

---

## Grande Padre relationship

Grande Padre owns:

- Donor/receiver indices and risk/priority scoring.
- Action slate construction and gate orchestration (SLO gates, rollback gates, blast radius, allowlists—patterns already reflected in cohort selection validators).
- **Dry-run** promotion paths (admission simulation packages exist toward this discipline).

Grande Padre does **not** silently bypass KHR refusal when ResourcePort or local node health blocks apply.

---

## Dry-run, rollback, verification

- **Dry-run:** must produce a **shadow slate** with predicted post-state and evidence requirements **before** mutating Cells (aligns with mutate-preview-apply patterns already named in repo contracts list).
- **Rollback:** every lease must declare rollback readiness where safety class demands it; rollback windows remain **first-class evidence** (see existing canary cohort reference JSON patterns).
- **Verification:** post-apply probes and telemetry convergence close the lease or trigger compensating rollback.

---

## KHR apply path (forward)

Execution order (target):

1. Lease promoted to **APPROVED** with attached evidence bundle addresses.
2. KHR receives **ApplyLease** RPC/message with Cell handles and ResourcePort snapshot.
3. KHR executes provider-specific apply, streaming **Observed** state back to Hyperdensity evidence plane.
4. Verification hooks clear or rollback triggers.

Until KHR ships, **legacy apply executors** remain authoritative but should converge vocabulary toward ResourceLease to reduce rewrite cost.

---

## Cloud adaptive constraints

On public cloud, ResourceLease verbs may map to:

- Provider resize APIs with cooldowns.
- **Replacement** Cells with restored Shell attachment (different blast radius—must be explicit in risk index).
- **No-op** futures when stock is unavailable—Hyperdensity must surface **blocked/remediable** states honestly.

ResourceFuture liquidity may emphasize **reserved capacity purchases** rather than host-level donor slicing.

---

## Consequences

### Positive

- Unified language for engineering, sales engineering, and compliance audits.
- Clear interface between **market** (Hyperdensity) and **hands** (KHR).

### Negative

- Requires refactoring duplicated “movement” concepts across Dashboard runtime and Hyperdensity repo during extraction phases documented in `docs/migration/dashboard-to-hyperdensity-extraction-plan.md`.

---

## Alternatives considered

1. **Keep ad-hoc per-provider tuning without lease abstraction** — rejected (does not scale to multi-provider KHR).
2. **Implement full market only inside Dashboard forever** — rejected (contradicts `Karl-Hyperdensity` README ownership direction).

---

## References

- `docs/architecture/KARL_HYPERDENSITY_KHR_FUSION.md`
- `pkg/leaseslate`, `pkg/marketcontroller` (implementation pointers)
- `examples/canary-cohort-execution-lease-reference.json` (example pointer)
