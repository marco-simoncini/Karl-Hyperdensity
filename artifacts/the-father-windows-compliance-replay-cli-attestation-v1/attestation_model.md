# Attestation Model

Model: `WindowsComplianceReplayAttestation`

Required fields:

- `attestationId`
- `replayId`
- `subjectType=windows-hyperdensity-ready-compliance-replay`
- `subjectRef`
- `policyVersion`
- `evaluatorVersion`
- `evidenceHash`
- `replayHash`
- `decisionSnapshot`
- `blockerSnapshot`
- `remediationSnapshot`
- `createdAt`
- `attestor`
- `signature`

Signature rules:

- mode only `unsigned-dev` or `future-signable`
- `signature.value` always empty string
- no cryptographic signing, no keys, no certs, no KMS
