# Windows FluidVirt Productization Roadmap v1

## Current confirmed state

- definitive live proof confirmed on `master-win11` (`WINDOWS_FLUIDVIRT_PRODUCT_PATH_CONFIRMED`)
- Node Fluid Actuator MVP hardened and validated (`windows_node_fluid_actuator_mvp_ready`)

## Meaning of product-path confirmed

It confirms that, for the validated shell:

- CPU is entitlement liquidity via controlled node actuator (`cpu.max`)
- RAM is balloon liquidity via QMP
- continuity invariants stayed true
- rollback and return-to-floor are evidenced

It does not imply production readiness.

## What is still missing for MVP controller

- orchestration controller for scheduled lease transitions
- durable policy distribution and allowlist lifecycle
- operational ownership model and admission integration
- production observability and incident workflows

## Integrating Node Fluid Actuator

- keep strict allowlist/identity pinning
- preserve no-arbitrary-write and no-parent-write invariants
- require dry-run acceptance before apply plans
- persist deterministic actuator audit hashes

## Integrating QMP RAM path

- require balloon capability and floor/ceiling envelope
- include guest-side memory safety witness before down transitions
- persist QMP before/after evidence

## Integrating Inventory/fluidShell

- maintain guest ACK and continuity witness as hard preconditions
- include guest evidence references inside lease/audit models

## Road to production readiness (without declaring it now)

1. implement controller execution engine with strict safety gates
2. add staged rollout policies and canary shell cohorts
3. integrate append-only audit store
4. expand chaos/failure-mode verification
5. formalize release sign-off and security review
