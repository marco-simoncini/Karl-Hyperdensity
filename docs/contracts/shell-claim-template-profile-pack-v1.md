# shell-claim-template-profile-pack-v1

Contract ID: `hyperdensity_shell_claim_templates_profile_pack_v1`

Defines the versioned golden template/profile catalog used by Shell Claim generation to create Hyperdensity-ready Linux shells without raw resource creation.

## Product principle

- No raw resource creation.
- Only Hyperdensity-ready shell creation.
- Template catalog only (`template_catalog_only`) and generate-only default (`generate_only`).

## Scope and safety

- `liveCreateAllowedByDefault=false`
- `dryRunRequiredBeforeCreate=true`
- `evidenceScope=evidence_namespace_only`
- `productionMutationAllowed=false`
- `autonomousApplyAllowed=false`
- `enforcementMode=disabled`
- Windows lane remains out-of-scope and is not counted as supported.

## Canonical objects

- `HyperdensityShellClaimTemplateProfilePack`
  - top-level catalog surface with policy/consistency projection, supported shell kinds/profiles, templates, validation rules, and safety gates.
- `HyperdensityShellClaimTemplateProfile`
  - profile metadata for donor/receiver/service/batch container and desktop/service/batch VM lanes.
- `HyperdensityShellClaimTemplate`
  - versioned template descriptors including required inputs, generated annotations/labels, resource envelopes, telemetry/compliance defaults, and dry-run/create safety defaults.
- `HyperdensityShellClaimTemplateValidationRule`
  - catalog validation rules for dry-run-before-create, evidence scope, windows out-of-scope posture, and factory compliance.
- `HyperdensityShellClaimTemplateSafetyGate`
  - runtime gate projection for consistency/alignment/scope defaults.

## Required supported shell kinds

- `linux_container`
- `linux_vm`

## Required supported profiles

- `linux_container_donor`
- `linux_container_receiver`
- `linux_container_service`
- `linux_container_batch`
- `linux_vm_desktop`
- `linux_vm_service`
- `linux_vm_batch`

## Alignment semantics

- Shell Claim generator must align to this catalog via explicit profile mapping.
- Every supported profile must have at least one template.
- Every template must include Hyperdensity annotations/labels plus CPU and memory resource envelope.
- Warming-up states are not ready states.
