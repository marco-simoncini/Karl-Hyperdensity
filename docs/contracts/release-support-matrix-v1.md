# release-support-matrix-v1

Contract ID: `hyperdensity_release_support_matrix_v1`

Defines the official Technical Preview release boundary for KARL Hyperdensity using evidence-backed support claims.

## Product principle

- No raw resource creation.
- Only Hyperdensity-ready shell creation.
- Release boundary visibility only (`release_boundary_visibility_only`) with evidence-backed claims (`evidence_backed_only`).

## Scope and safety invariants

- `releaseTrack=technical_preview`
- `enforcementMode=disabled`
- `autonomousApplyAllowed=false`
- `productionMutationAllowed=false`
- `evidenceScope=evidence_namespace_only`
- `policyPackId=hyperdensity_policy_pack_v1`
- `policyConsistencyRequired=true`
- `profilePackId=hyperdensity_shell_claim_templates_profile_pack_v1`

## Required support-level semantics

- `supported_for_technical_preview`
- `evidence_only`
- `preview_only`
- `simulation_only`
- `dry_run_only`
- `recommendation_only`
- `operator_controlled_only`
- `blocked`
- `out_of_scope`
- `future`

## Required release-boundary statements

- Linux container lane is supported only within proven bounds.
- Linux VM lane is supported only within proven bounds and object-specific evidence.
- Windows lane is out-of-scope/frozen.
- Production workload mutation is not supported.
- Enforcement remains disabled.
- Autonomous apply remains disabled.
- Evidence namespace is the live proof namespace.
- No raw resource creation is a supported product path.
- Shell Claim/Profile Pack is the official creation catalog.
- `warming_up`, `partial`, and `blocked` are not ready states.

## Canonical objects

- `HyperdensityReleaseSupportMatrix`
  - top-level release-boundary catalog with shell kinds/profiles/capabilities/operations/surfaces/safety/proof/limitations.
- `HyperdensityReleaseSupportShellKind`
  - support classification per shell kind and supported operations.
- `HyperdensityReleaseSupportProfile`
  - profile-level support and readiness semantics aligned to profile pack.
- `HyperdensityReleaseSupportCapability`
  - capability-level support, proof references, safety gates, and follow-up actions.
- `HyperdensityReleaseSupportOperation`
  - operation-level support mode, mutation scope, and proof status.
- `HyperdensityReleaseSupportSurface`
  - API surface-level mode/state and safety posture.
- `HyperdensityReleaseSupportProofCatalogEntry`
  - explicit evidence mapping to milestones, artifacts, and commits.
- `HyperdensityReleaseSupportLimitation`
  - known constraints with remediation and target release.
- `HyperdensityReleaseSupportOutOfScope`
  - frozen/out-of-scope lanes and future entry conditions.
- `HyperdensityReleaseSupportMatrixSafetyGate`
  - consistency and safety gates with blocker state.

## Technical Preview support boundary

- Supported shell kinds:
  - `linux_container`
  - `linux_vm`
- Required supported profiles:
  - `linux_container_donor`
  - `linux_container_receiver`
  - `linux_container_service`
  - `linux_container_batch`
  - `linux_vm_desktop`
  - `linux_vm_service`
  - `linux_vm_batch`
- Windows shell kinds/profiles are out-of-scope metadata only and must not be counted as supported.

## Consistency expectations

- Matrix support claims must align with:
  - Policy Pack (`hyperdensity_policy_pack_v1`)
  - Policy Pack Consistency Checker
  - Shell Claim Template/Profile Pack (`hyperdensity_shell_claim_templates_profile_pack_v1`)
  - Shell Claim Generator surfaces
  - Admission Guard + Mutate Preview + Enforce Simulation + Mutate Preview Apply Dry-Run
  - Shell Factory
  - Resource Exchange
- Missing required source surfaces must degrade or block matrix state.
