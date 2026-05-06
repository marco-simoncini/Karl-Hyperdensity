# auto-rollback-controller-v1

Contract ID: `hyperdensity_auto_rollback_controller_v1`

Defines the automatic rollback control-plane readiness surface for KARL Hyperdensity:
`executionEngine.hyperdensityAutoRollbackController`.

## Purpose

- Model when guarded sandbox actions are eligible for automatic rollback.
- Keep rollback scope strictly inside evidence namespace.
- Preserve projection/readiness-first and safety-first semantics.

## Core safety boundary

- `autoRollbackControllerId=hyperdensity_auto_rollback_controller_v1`
- `autoRollbackControllerVersion=v1`
- `releaseTrack=technical_preview`
- `controllerMode=guarded_sandbox_rollback_readiness`
- `allowedNamespace=karl-hyperdensity-evidence`
- `productionRollbackAllowed=false`
- `productionMutationAllowed=false`
- `enforcementMode=disabled`
- `autonomousProductionApplyAllowed=false`
- `operatorKillSwitchRequired=true`
- `rollbackRequired=true`
- `verificationRequired=true`
- `auditRequired=true`

Automatic rollback may be considered ready only when all required gates pass.

## Mandatory rollback triggers

- `verification_failed`
- `runtime_not_converged`
- `cgroup_not_converged`
- `qga_libvirt_qmp_not_converged`
- `warning_event_detected`
- `restart_count_changed`
- `receiver_not_improved`
- `donor_became_pressured`
- `slo_degraded`
- `operator_kill_switch_triggered`
- `action_expired_before_verify`

## Mandatory gates

- `evidence_namespace_only`
- `production_rollback_disabled`
- `production_mutation_disabled`
- `enforcement_disabled`
- `autonomous_production_apply_disabled`
- `operator_kill_switch_available`
- `rollback_source_required`
- `verification_required`
- `audit_required`
- `low_risk_required`
- `support_boundary_required`
- `policy_consistency_required`
- `blast_radius_budget_required`
- `raw_runtime_controls_not_exposed`
- `raw_resource_creation_not_allowed`
- `windows_out_of_scope`

## Rollback plan constraints

Rollback plans must be projection/readiness entries derived from guarded sandbox candidates and must never bypass:

- evidence namespace scope
- rollback source availability/proof
- verification requirement
- audit requirement
- low-risk requirement
- support boundary validity
- policy consistency validity
- blast-radius budget validity

Production rollback remains disabled in this milestone.

## Explicit non-claims

- Not production rollback.
- Not production autonomous apply.
- Not enforcement.
- Not production mutation.
- Not raw runtime control exposure.
- Not raw resource creation.
- Not GA.
- Not HA from single-node proof.
- Not Windows.
- Not generic VM RAM template mutation.

## Required safety copy

- "Automatic rollback is evidence-namespace only."
- "Production rollback is disabled."
- "Production autonomous apply is disabled."
- "Enforcement is disabled."
- "Rollback source is required."
- "Verification is required."
- "Operator kill switch is required."
- "No raw runtime controls are exposed."
- "No raw resource creation."
- "Technical Preview boundary active."
