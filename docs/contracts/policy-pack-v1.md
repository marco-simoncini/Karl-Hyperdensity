# policy-pack-v1

Contract ID: `hyperdensity_policy_pack_v1`

Versioned visibility-only policy catalog for KARL Hyperdensity.

## Purpose

`policy-pack-v1` unifies product policy gates that were previously scattered across:

- Shell Factory
- Shell Claim generate/dry-run/create surfaces
- Shell Claim evidence create history
- Admission Guard audit + mutate preview
- Resource Exchange dry-run/apply/rollback/history
- telemetry freshness/confidence and warmup semantics

This contract is packaging/visibility only. It is not an enforcement switch.

## Non-negotiable contract assertions

- `policyPackMode=visibility_only`
- `enforcementMode=disabled`
- `autonomousApplyAllowed=false`
- no mutation path is introduced by this contract
- no production workload mutation
- evidence namespace only for live proof/create paths
- operator-controlled only for apply/create paths
- Windows lane remains out-of-scope
- `factory_warming_up` is not `factory_ready`

## Required top-level fields

- `policyPackId`
- `policyPackVersion`
- `policyPackMode`
- `enforcementMode`
- `autonomousApplyAllowed`
- `supportedShellKinds`
- `supportedProfiles`
- `factoryRequirements`
- `claimValidationRules`
- `admissionGuardRules`
- `mutatePreviewDefaults`
- `exchangeEligibilityRules`
- `stageApplyRules`
- `shellClaimEvidenceCreateRules`
- `safetyGates`
- `warmupPolicy`

## Safety gate catalog (semantic minimum)

The contract must cover:

- `evidence_namespace_only`
- `operator_controlled_only`
- `autonomous_apply_disabled`
- `dry_run_required_before_create_or_apply`
- `rollback_required_where_applicable`
- `cleanup_required_where_applicable`
- `warning_events_clean_required`
- `no_production_mutation`
- `windows_lane_out_of_scope`
- `no_rollout_fallback_for_live_staged_apply`
- `pod_resize_capability_required_for_container_cpu_stage_apply`
- `sustained_idle_evidence_required`
- `safe_band_stage_splitting_required`
- `factory_managed_required`
- `shell_profile_required`
- `resource_envelope_required`
- `telemetry_required_for_ready`
- `compliance_required_for_ready`
- `exchange_eligibility_required_for_ready`
- `warming_up_is_not_ready`

## Warmup policy

`factory_warming_up`, `factory_partial`, and `factory_blocked` are visibility states and are never projected as ready.
