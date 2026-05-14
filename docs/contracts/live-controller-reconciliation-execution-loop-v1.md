# Live Controller Reconciliation Execution Loop v1

**Milestone:** `hyperdensity_live_controller_reconciliation_execution_loop_v1`

## Purpose

Sprint 13 transforms the Sprint 12 deterministic continuous resource market controller into a **live reconciliation loop** with persistent controller state, scheduled ticks, lease/action/future lifecycles, permitted-scope execution selection, post-execution accounting hooks, and realized idle-compression tracking.

## Core product rule

Sprint 12 generated prioritized actions and futures. Sprint 13 **keeps them live**:

1. Load previous controller state
2. Collect observed market state
3. Compute desired market state
4. Compute reconciliation diff
5. Invalidate stale leases/actions/futures
6. Refresh indices, action slate, resource futures
7. Select executable actions within permitted scopes only
8. Record lifecycle transitions and audit trail
9. Track realized compression only from movements with mutation/post-verify evidence

## Permitted execution scopes

- `operator_controlled`
- `sandbox_auto`
- `nonprod_auto`
- `production_canary_auto`

## Forbidden execution scopes

- `general_production_auto`
- `production_auto_with_policy`

## Source of truth

| Concern | Authority |
|---|---|
| Live controller / reconciliation | Karl-Hyperdensity |
| Runtime actuator / mutation evidence | FluidVirt |
| Operator projection | Karl-Dashboard (read-only) |
| Identity / signals | Karl-Inventory |

## Claim boundaries

- Projected compression and projected moved idle value are **not** realized
- Realized compression requires `mutationObserved=true` and `postVerifyPassed=true`
- `projectedCompressionCountedAsRealized=false` always in Sprint 13
- `generalProductionAutoAllowed=false`, `productionAutoWithPolicy=false`
- Synthetic/reference fleet is not production proof

## Sprint 13 allowed flags

- `liveReconciliationEnabled=true`
- `stateStoreEnabled=true`
- `scheduledTicksEnabled=true`
- `leaseLifecycleEnabled=true`
- `actionLifecycleEnabled=true`
- `futuresRefreshEnabled=true`
- `executionSelectionEnabled=true`
- `realizedCompressionTrackingEnabled=true`

## Sprint 13 forbidden flags

- `generalProductionAutoAllowed=false`
- `productionAutoWithPolicy=false`
- `universalGuaranteedSavingsAllowed=false`
- `estimatedIdleCountedAsMoved=false`
- `projectedCompressionCountedAsRealized=false`
- `dashboardExecutor=false`
- `fluidvirtPolicyAuthority=false`
- `inventoryRuntimeExecutor=false`
