# KARL Hyperdensity — Guarded Auto Apply Sandbox / NonProd v1

**Contract ID:** `hyperdensity_guarded_auto_apply_sandbox_nonprod_v1`  
**Milestone:** `hyperdensity_guarded_auto_apply_sandbox_nonprod_v1`  
**Release track:** `technical_preview`

## Product definition

> KARL Hyperdensity can execute guarded auto-apply in sandbox/non-production scope for policy-approved actions, with kill switch, circuit breaker, blast radius, rate limit, cooldown, SLO guard, rollback and audit, while keeping production autonomous apply disabled.

## State machine

`policy_candidate` → `sandbox_ready` → `auto_preflight_recheck` → `kill_switch_checked` → `circuit_breaker_checked` → `blast_radius_reserved` → `rate_limit_reserved` → `cooldown_reserved` → `auto_apply_requested` → `fluidvirt_invoked` → `mutation_observed` → `post_verify_evaluated` → `rollback_window_open` → `accounted_preview` → `auto_closed`

Failure states: `auto_blocked`, `kill_switch_blocked`, `circuit_breaker_blocked`, `blast_radius_blocked`, `rate_limit_blocked`, `cooldown_blocked`, `preflight_failed`, `fluidvirt_invocation_failed`, `mutation_not_observed`, `post_verify_failed`, `auto_rollback_required`, `auto_rollback_executed`, `auto_rollback_failed`, `expired`

## Sprint 8 invariants

- `sandboxAutoApplyAllowed=true`
- `nonProdAutoApplyAllowed=true`
- `autoApplyExecutionEnabled=true` (sandbox/nonprod only)
- `productionAutonomousApplyAllowed=false`
- `productionScope=false`
- `productionMutationAllowed=false`
- `rawRuntimeControlsExposed=false`

## Execution rules

- Only `sandbox_ready` or `nonprod_ready` candidates may auto-execute
- `candidate_only` is not executable
- FluidVirt is the actuator; Dashboard is not executor
- Windows evidence-gated and synthetic/shadow actions blocked
- Regression after apply triggers automatic rollback

## Source-of-truth map

| Domain | Owner |
|--------|-------|
| Auto apply execution / policy | Karl-Hyperdensity |
| Runtime actuator invocation | FluidVirt |
| Execution projection | Karl-Dashboard (read-only) |
| Identity / signals | Karl-Inventory |

## Forbidden claims

- production autonomous apply / autonomous production mutation
- productionScope=true auto apply
- guaranteed savings active
- universal performance improvement
- Windows total RAM hotplug / logical vCPU hotplug
- Dashboard as executor / raw runtime controls
- FluidVirt as policy authority / market controller
