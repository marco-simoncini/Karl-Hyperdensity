# Windows Fluid Compliance Replay Attestation v1

`karl-fluid-compliance-replay` is a read-only CLI for replaying `EvaluateWindowsHyperdensityReadyCompliance` against JSON evidence fixtures.

## Read-only scope

- no runtime mutation
- no CPU apply
- no RAM apply
- no actuator apply
- no cluster calls
- no QMP calls

This contract is for verification and audit replay only.

## Input contract

The CLI accepts:

- fixture payloads with `input` field (existing compliance replay fixtures), or
- direct `EvaluateWindowsHyperdensityReadyComplianceInput` JSON.

## Output contract

Replay output includes:

- `replayId`
- `inputRef`
- `evaluationTime`
- `compliancePhase`
- `vmRef`
- `namespace`
- `shellRef`
- `evidenceSummary`
- `blockers`
- `remediationActions`
- `automatableActions`
- `manualActions`
- `risk`
- `poolContext`
- `poolScalingMechanismBlocked`
- `hyperdensityReady`
- `evidenceHash`
- `replayHash`
- `auditRefs`
- `mutationFlags` (all false in this CLI)

## Hash behavior

- `evidenceHash`: deterministic local SHA-256 over canonical replay input JSON.
- `replayHash`: deterministic local SHA-256 over replay decision payload.

These hashes are deterministic audit references, not signatures.

## Future-signable attestation

Optional envelope `WindowsComplianceReplayAttestation`:

- `subjectType=windows-hyperdensity-ready-compliance-replay`
- `signature.mode`: `unsigned-dev` or `future-signable`
- `signature.value`: always empty string

No real cryptographic signing is performed:

- no key generation
- no certificate material
- no KMS integration
- no secret material

## Security and release note

This replay attestation path does not enable apply behavior and must not be interpreted as production-ready signing or runtime control.
