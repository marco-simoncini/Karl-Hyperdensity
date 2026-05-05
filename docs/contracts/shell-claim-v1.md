# shell-claim-v1

Contract ID: `hyperdensity_shell_claim_v1`

Defines the canonical declarative claim used to generate Hyperdensity-ready shell manifests.

## Product principle

- No raw resource creation.
- Only Hyperdensity-ready shell creation from claim + profile.
- Generator-first milestone: `generate_only` and `evidence_reference_only` are active; `enforce` is reserved.

## Canonical objects

- `HyperdensityShellClaim`
  - structured intent with profile, namespace/name, exchange policy, envelope, telemetry, rollback, runtime-overlay, and creation mode.
- `HyperdensityGeneratedShellManifest`
  - generated Kubernetes manifest preview from one claim (dry-run artifact).
- `ShellClaimValidationResult`
  - validation state, score, blocker, missing requirements, remediation lane.

## Validation baseline

- Namespace is required and must be evidence/reference for active modes.
- Profile must exist and match shell kind.
- Envelope must be complete and ordered (`floor <= baseline <= ceiling`).
- Burst steps must be positive.
- Donation/receive policy must align with profile role.
- Telemetry + rollback + runtime-overlay + factory-managed are required.
- Unsupported shell kinds and unsupported profiles must return remediation lanes.
