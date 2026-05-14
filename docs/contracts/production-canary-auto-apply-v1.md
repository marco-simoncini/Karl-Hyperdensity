# KARL Hyperdensity — Production Canary Auto Apply v1

**Contract ID:** `hyperdensity_production_canary_auto_apply_v1`  
**Milestone:** `hyperdensity_production_canary_auto_apply_v1`  
**Release track:** `technical_preview`

## Product definition

> KARL Hyperdensity can execute guarded production canary auto-apply for explicitly allowlisted workloads under strict blast-radius, SLO, rollback, kill-switch, circuit-breaker, rate-limit and immutable-audit controls, while keeping general production auto disabled.

## Sprint 9 invariants

- `productionCanaryAutoApplyAllowed=true`
- `productionCanaryScope=true`
- `generalProductionAutoAllowed=false`
- `productionAutoWithPolicy=false`
- `productionAutonomousApplyAllowed=true` (canary scope only)
- `productionScope=true` (canary scope only)
- `rawRuntimeControlsExposed=false`

## Strict rules

- Only `production_canary_ready` actions may execute
- `productionScope=true` only when `productionCanaryScope=true` and allowlisted
- `production_auto_with_policy` forbidden
- General production auto forbidden
- Windows evidence-gated and synthetic/shadow blocked
- FluidVirt is actuator; Dashboard is projection-only

## Source-of-truth map

| Domain | Owner |
|--------|-------|
| Production canary execution / policy | Karl-Hyperdensity |
| Runtime actuator invocation | FluidVirt |
| Canary projection | Karl-Dashboard (read-only) |
| Identity / signals | Karl-Inventory |

## Forbidden claims

- general production auto
- production_auto_with_policy
- broad autonomous production mutation
- guaranteed savings active
- universal performance improvement
- Windows total RAM hotplug / logical vCPU hotplug
- Dashboard as executor / raw runtime controls
- FluidVirt as policy authority / market controller
