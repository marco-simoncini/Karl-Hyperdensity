# mutate-preview-apply-dry-run-v1

Contract ID: `hyperdensity_mutate_preview_apply_dry_run_v1`

Server-side dry-run evidence contract for Admission Guard Mutate Preview remediation patches.

## Purpose

This contract closes the remediation proof chain without applying mutations:

`raw object detected -> admission guard classifies -> mutate preview proposes patch -> enforce simulation projects reject -> mutate preview apply dry-run validates Kubernetes API acceptance`

The contract is strictly dry-run only:

- `dryRunMode=server_side_apply_dry_run_only`
- `mutationAllowed=false`
- `productionMutationAllowed=false`
- `enforcementMode=disabled`
- `admissionGuardMode=audit_only`
- `mutatePreviewMode=audit_preview_only`
- `autonomousApplyAllowed=false`

## Required top-level fields

- `dryRunId` (must be `hyperdensity_mutate_preview_apply_dry_run_v1`)
- `dryRunVersion` (must be `v1`)
- `dryRunMode` (must be `server_side_apply_dry_run_only`)
- `mutationAllowed` (must be `false`)
- `productionMutationAllowed` (must be `false`)
- `enforcementMode` (must be `disabled`)
- `admissionGuardMode` (must be `audit_only`)
- `mutatePreviewMode` (must be `audit_preview_only`)
- `autonomousApplyAllowed` (must be `false`)
- `evidenceScope` (must be `evidence_namespace_only`)
- `policyPackId` (must be `hyperdensity_policy_pack_v1`)
- `policyConsistencyRequired` (must be `true`)
- `sourceSurface` (must be `admission_guard_mutate_preview`)
- `dryRunTargets`
- `dryRunResults`
- `safetyGates`
- `summary`
- `nextDryRunAction`
- `dryRunBlocker`

## Required target/result fields

Each element in `dryRunTargets[]` and `dryRunResults[]` must contain:

- `targetId`
- `objectRef`
- `namespace`
- `kind`
- `name`
- `sourcePreviewId`
- `factoryManaged`
- `shellKind`
- `shellProfile`
- `currentReadinessState`
- `previewPatchAvailable`
- `patchType`
- `patchSummary`
- `serverSideDryRunAttempted`
- `serverSideDryRunAccepted`
- `serverSideDryRunRejected`
- `rejectionReason`
- `validationWarnings`
- `remediationHints`
- `fieldManager`
- `dryRunRequestMode`
- `mutationObserved`
- `annotationsBeforeHash`
- `annotationsAfterHash`
- `objectBeforeHash`
- `objectAfterHash`
- `noMutationVerified`
- `cleanupRequired`
- `cleanupVerified`
- `sourceRuleIds`
- `safetyNotes`

## Contract semantics

- Dry-run means Kubernetes API dry-run only.
- The remediation patch must not be persisted.
- Object annotations/spec must remain unchanged after dry-run.
- `mutationObserved` must remain `false`.
- Dry-run success is not readiness certification.
- Dry-run rejection must be surfaced as remediation evidence, not hidden.
- Production namespace objects are out-of-scope.
- Windows remains out-of-scope.
- Enforcement remains disabled.
- Mutate Preview remains `audit_preview_only`.
- Admission Guard remains `audit_only`.
- Autonomous apply remains `false`.
