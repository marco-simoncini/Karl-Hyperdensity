# blast-radius-policy-v1

Contract ID: `hyperdensity_blast_radius_policy_v1`

Defines the blast-radius and safety-budget readiness surface:
`executionEngine.hyperdensityBlastRadiusPolicy`.

## Purpose

- Define action scope, speed, and budget constraints for guarded auto readiness.
- Bind autonomous readiness to explicit freeze, stop, and escalation policy.
- Keep policy/readiness-first semantics in Technical Preview.

## Core safety boundary

- `blastRadiusPolicyId=hyperdensity_blast_radius_policy_v1`
- `blastRadiusPolicyVersion=v1`
- `releaseTrack=technical_preview`
- `policyMode=guarded_auto_safety_budget`
- `productionAutonomousApplyAllowed=false`
- `productionMutationAllowed=false`
- `enforcementMode=disabled`
- `evidenceNamespace=karl-hyperdensity-evidence`
- `operatorKillSwitchRequired=true`

## Mandatory model components

- global budget model
- namespace budgets
- resource budgets
- concurrency limits
- rate limits
- freeze conditions
- stop conditions
- escalation rules
- safety gates

## Mandatory freeze conditions

- `warning_event_detected`
- `rollback_failure`
- `verification_failure`
- `donor_pressure_increased`
- `receiver_not_improved`
- `policy_inconsistency`
- `support_boundary_missing`
- `unknown_risk_detected`
- `kill_switch_triggered`
- `audit_gap_detected`

## Mandatory stop conditions

- `max_budget_exceeded`
- `high_risk_action_detected`
- `blocked_action_detected`
- `expired_action_detected`
- `unsupported_shell_detected`
- `windows_lane_detected`
- `raw_runtime_control_requested`
- `production_scope_requested`

## Mandatory escalation rules

- `medium_risk_requires_operator_review`
- `high_risk_blocks_auto`
- `unknown_risk_blocks_auto`
- `production_scope_blocks_auto`
- `rollback_missing_blocks_auto`
- `dry_run_missing_blocks_auto`

## Explicit requirements

- no production autonomous apply
- no enforcement
- no production mutation
- evidence namespace only for guarded-auto readiness
- budgets required
- rate limits required
- freeze conditions required
- escalation rules required
- stop conditions required
- kill switch required
- rollback required
- verification required
- audit required

## Explicit non-claims

- Not GA.
- Not HA from single-node proof.
- Not Windows.
- Not generic VM RAM template mutation.

## Required safety copy

- "Blast radius limits are enforced before any guarded auto action."
- "Production autonomous apply is disabled."
- "Enforcement is disabled."
- "No production mutation."
- "Dry-run is required."
- "Rollback proof is required."
- "Verification is required."
- "Audit is required."
- "Operator kill switch is required."
- "High and unknown risk block auto."
- "Technical Preview boundary active."
