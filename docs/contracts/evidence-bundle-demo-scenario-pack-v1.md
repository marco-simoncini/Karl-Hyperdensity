# evidence-bundle-demo-scenario-pack-v1

Contract ID: `hyperdensity_evidence_bundle_demo_scenario_pack_v1`

Defines the Technical Preview evidence/demo packaging surface that projects a credible 10-15 minute Grande Padre flow using existing evidence-backed milestones.

## Product principle

- No raw resource creation.
- Only Hyperdensity-ready shell creation.
- Evidence bundle packaging only (`evidence_bundle_only`).
- Guided operator demo mode (`guided_operator_demo`).

## Scope and safety invariants

- `releaseTrack=technical_preview`
- `enforcementMode=disabled`
- `autonomousApplyAllowed=false`
- `productionMutationAllowed=false`
- `evidenceScope=evidence_namespace_only`
- `supportMatrixId=hyperdensity_release_support_matrix_v1`
- `policyPackId=hyperdensity_policy_pack_v1`
- `profilePackId=hyperdensity_shell_claim_templates_profile_pack_v1`

## Canonical objects

- `HyperdensityEvidenceBundleDemoScenarioPack`
  - top-level bundle with supported claims, demo scenarios, evidence catalog, artifact index, proof checks, safety gates, and runbook.
- `HyperdensityEvidenceBundleSupportedClaim`
  - approved/rejected claim wording with proof references and safety boundaries.
- `HyperdensityEvidenceBundleDemoScenario`
  - objective-driven scenario definition with structured steps and success criteria.
- `HyperdensityEvidenceBundleEvidenceCatalogEntry`
  - mapping from milestone claim to evidence source and artifact path.
- `HyperdensityEvidenceBundleArtifactIndexEntry`
  - artifact freshness/validation projection for demo readiness.
- `HyperdensityEvidenceBundleProofCheck`
  - consistency checks across required source surfaces and claim boundaries.
- `HyperdensityEvidenceBundleSafetyGate`
  - release-demo safety gate status with blockers.

## Required claim semantics

- Linux container CPU/RAM live up/down is supported only within proven Technical Preview bounds.
- Linux VM CPU/RAM live adjustment is evidence-backed/object-specific within proven bounds.
- VM RAM live path wording must be runtime overlay `virtio-mem/QMP/QOM requested-size` where proven.
- Generic KubeVirt memory template mutation must not be approved wording.
- Windows lane remains frozen out-of-scope.
- No production mutation, no autonomous apply, no enforcement.
- `warming_up`, `partial`, and `blocked` are not ready states.

## Bundle semantics

- Packaging/projection only; no runtime behavior widening.
- Existing surfaces remain source-of-truth:
  - Release Support Matrix
  - Policy Pack + Policy Consistency
  - Shell Claim Template/Profile Pack
  - Admission Guard + Mutate Preview + Enforce Simulation + Mutate Preview Apply Dry-Run
  - Shell Factory + Shell Claim + Resource Exchange
- Missing required source surfaces must degrade/block bundle state.
