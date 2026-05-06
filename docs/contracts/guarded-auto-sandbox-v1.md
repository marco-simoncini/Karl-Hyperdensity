# guarded-auto-sandbox-v1

Contract ID: `hyperdensity_guarded_auto_sandbox_v1`

Defines the first bounded autonomous-readiness surface for KARL Hyperdensity:
`executionEngine.hyperdensityGuardedAutoSandbox`.

## Purpose

- Model guarded autonomous execution readiness in Technical Preview.
- Keep scope strictly inside evidence namespace sandbox.
- Preserve projection-first and safety-first semantics.

## Core safety boundary

- `allowedNamespace=karl-hyperdensity-evidence`
- `sandboxMode=guarded_auto_evidence_namespace_only`
- `releaseTrack=technical_preview`
- `productionMutationAllowed=false`
- `enforcementMode=disabled`
- `autonomousProductionApplyAllowed=false`
- `operatorKillSwitchRequired=true`
- `dryRunRequired=true`
- `rollbackRequired=true`
- `auditRequired=true`
- `verificationRequired=true`

Guarded auto may be considered ready only when all required gates pass.

## Mandatory gates

- `evidence_namespace_only`
- `production_mutation_disabled`
- `enforcement_disabled`
- `autonomous_production_apply_disabled`
- `operator_kill_switch_available`
- `action_slate_ready_required`
- `dry_run_ready_required`
- `rollback_ready_required`
- `low_risk_required`
- `support_boundary_required`
- `policy_consistency_required`
- `blast_radius_budget_required`
- `audit_required`
- `verification_required`
- `raw_runtime_controls_not_exposed`
- `raw_resource_creation_not_allowed`
- `windows_out_of_scope`

## Candidate action constraints

Candidate actions must come from Action Slate references and must never bypass:

- evidence namespace scope
- ready action state
- ready dry-run
- ready rollback proof
- low risk
- support boundary validity
- policy consistency validity
- blast radius budget validity

Blocked/expired/high-risk/non-sandbox candidates are non-executable.

## Explicit non-claims

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

- "Guarded auto is evidence-namespace only."
- "Production autonomous apply is disabled."
- "Enforcement is disabled."
- "No production mutation."
- "Dry-run is required."
- "Rollback proof is required."
- "Low risk is required."
- "Operator kill switch is required."
- "No raw runtime controls are exposed."
- "No raw resource creation."
- "Technical Preview boundary active."
