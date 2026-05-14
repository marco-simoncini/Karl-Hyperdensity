# KARL Hyperdensity — Production Kernel Boundary v1

**Contract ID:** `hyperdensity_production_kernel_boundary_v1`  
**Milestone:** `hyperdensity_production_kernel_boundary_v1`  
**Release track:** `technical_preview` (kernel boundary lock; not GA autonomous apply)

## Product formula (context)

> Hyperdensity non aspetta di scalare.  
> Hyperdensity mantiene un mercato di risorse già prevalidato.

This sprint freezes **who decides**, **who actuates**, **who displays evidence**, and **what may be claimed** — without enabling production autonomous apply, guaranteed savings, or universal performance claims.

## Allowed claim (this sprint only)

> KARL Hyperdensity is a governed runtime resource market powered by FluidVirt runtime actuation, with evidence-backed claim boundaries.

## Repository responsibilities

### FluidVirt (`marco-simoncini/FluidVirt`)

- Runtime actuator (node actuator, libvirt, cgroup).
- CPU entitlement actuation.
- RAM envelope / guest usable memory actuation.
- Guest evidence collection (QGA/agent/guest probe where available).
- Produces **Runtime Mutation Result** and **rollback result** artifacts.
- Does **not** own donor/receiver market indexing, leases, or claim policy.

### Hyperdensity (`marco-simoncini/Karl-Hyperdensity`)

- Donor index, receiver index, risk index, priority index.
- Blocked/remediable index.
- Resource futures, action slate, resource lease contracts.
- SLO guard decision and value/accounting decision (contracted; evidence-gated for GA claims).
- Claim policy registry and release gate definitions.
- Does **not** execute raw libvirt/cgroup mutations directly in product paths.

### Dashboard (`marco-simoncini/Karl-Dashboard`)

- Read-only cockpit and operator UI projection.
- Surfaces evidence, claim policy, component responsibility map, blockers.
- **Must not** invent runtime state or expose raw runtime controls.
- **Must not** be source of truth for runtime mutation.

### Karl-Hyperdensity (schemas)

- Contracts, JSON schemas, reference examples, validators.
- Claim policy v2 conservative defaults.

### Inventory / Warden (`marco-simoncini/Karl-Inventory`)

- Identity, access, endpoint/guest signals, Windows readiness signals.
- **Must not** become the Hyperdensity runtime engine or declare production runtime apply.

## Safety invariants (Sprint 1)

| Invariant | Sprint 1 value |
|-----------|----------------|
| `productionAutonomousApplyAllowed` | `false` |
| `guaranteedSavingsAllowed` | `false` |
| `universalPerformanceImprovementAllowed` | `false` |
| `logicalVcpuHotplugClaimAllowed` | `false` (unless separately proven) |
| `windowsTotalRamHotplugClaimAllowed` | `false` (unless separately proven) |
| `syntheticFleetProductionClaimAllowed` | `false` |
| Dashboard source of truth for mutation | forbidden |
| Inventory Hyperdensity engine | forbidden |
| Raw runtime controls in Dashboard | forbidden |

## Required surfaces

- `hyperdensity_production_kernel_boundary_v1` (ConfigMap `karl-system/hyperdensity-production-kernel-boundary-v1`, key `surface.json`)
- `hyperdensity_runtime_actuator_boundary_v1` (FluidVirt declaration)
- Shell Passport, Runtime Mutation Result, Resource Lease, Claim Policy v2 schemas

## Related schemas

- `schemas/production-kernel-boundary-v1.schema.json`
- `schemas/hyperdensity-claim-policy-v2.schema.json`
- `schemas/shell-passport-v1.schema.json`
- `schemas/runtime-mutation-result-v1.schema.json`
- `schemas/resource-lease-v1.schema.json`
