# KARL Hyperdensity — Operator-Controlled Apply Gate v1

**Contract ID:** `hyperdensity_operator_controlled_apply_gate_v1`  
**Milestone:** `hyperdensity_operator_controlled_apply_gate_v1`  
**Release track:** `technical_preview`

## Product definition

> KARL Hyperdensity can execute operator-approved resource movements through a guarded apply gate, producing runtime mutation, post-verify and rollback evidence without enabling autonomous production apply.

## Allowed Sprint 4 claim

Operator-controlled apply is permitted with `operatorControlledApplyAllowed=true` while `autoApplyAllowed=false`, `productionAutonomousApplyAllowed=false`, and `productionScope=false`.

## State machine

```
prepared → dry_run_valid → operator_approval_required → operator_approved
  → apply_requested → fluidvirt_invocation_recorded → mutation_observed
  → post_verify_passed → rollback_window_open → accounted_preview → closed
```

Failure states: `approval_denied`, `apply_blocked`, `fluidvirt_invocation_failed`, `mutation_not_observed`, `post_verify_failed`, `rollback_required`, `rollback_executed`, `rollback_failed`, `expired`.

## Concepts

| Concept | Owner | Description |
|---------|-------|-------------|
| Operator apply gate | Karl-Hyperdensity | Controlled transition from `operator_controlled_ready` action to apply evidence |
| Operator approval record | Karl-Hyperdensity | Explicit operator identity, timestamp, reason, scope, risk acceptance |
| Apply request | Karl-Hyperdensity | References action, lease, approval, dry-run, rollback, SLO, risk |
| FluidVirt invocation record | FluidVirt | Guarded actuator invocation evidence |
| Runtime mutation observation | FluidVirt | Donor/receiver runtime and guest observations |
| Post-verify result | Karl-Hyperdensity + FluidVirt | Runtime delta, SLO guard, health, rollback readiness |
| Rollback window | Karl-Hyperdensity | Open window with rollback plan and trigger policy |
| Apply audit event | Karl-Hyperdensity | Immutable audit trail |
| Apply gate projection | Karl-Dashboard | Read-only ConfigMap surface |

## Strict rules

- Only Sprint 3 actions with `readyForApply=operator_controlled_ready` may enter the gate.
- Operator approval must be explicit with `approvalMode=operator_required`.
- FluidVirt is the sole actuator; Dashboard and Inventory must not apply mutation.
- `mutationScope=technical_preview_operator_controlled` for all Sprint 4 mutations.
- No reboot/recreate/rollout/migration for in-place movement certification.

## Forbidden claims

- guaranteed savings active
- universal performance improvement
- production autonomous apply
- Windows total RAM hotplug / logical vCPU hotplug
- Dashboard raw runtime controls
- Inventory/Warden runtime apply
