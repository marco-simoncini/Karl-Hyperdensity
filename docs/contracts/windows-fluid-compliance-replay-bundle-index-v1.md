# Windows Fluid Compliance Replay Bundle Index v1

`WindowsComplianceReplayBundleIndex` links replay input, replay output, and optional attestation envelopes into a deterministic local hash chain.

## Purpose

Create an audit trail across compliance replay runs:

- input fixture/evidence linkage
- replay decision linkage
- attestation linkage
- deterministic hash-chain continuity

## Model overview

Top-level fields:

- `bundleId`
- `bundleVersion`
- `createdAt`
- `subjectRef`
- `subjectType`
- `runCount`
- `runs`
- `chain`
- `aggregateStatus`
- `latestCompliancePhase`
- `latestHyperdensityReady`
- `auditRefs`

Run fields:

- `runId`
- `inputRef`
- `outputRef`
- `attestationRef`
- `evidenceHash`
- `replayHash`
- `attestationHash`
- `previousRunHash`
- `runHash`
- `evaluationTime`
- `compliancePhase`
- `hyperdensityReady`
- `blockers`
- `remediationActions`

## Hash linkage

Each run forms:

`evidenceHash -> replayHash -> attestationHash -> runHash`

`runHash` includes `previousRunHash` in canonical payload, so each run binds to prior chain state.

Chain metadata:

- `chainMode=local-deterministic-hash-chain`
- `firstRunHash`
- `latestRunHash`
- `chainValid`
- `brokenAtRunId`
- `notes`

## Validation rules

Validator checks:

- bundle version present
- run count consistency
- required run hashes present
- deterministic run hash recomputation
- previous hash linkage integrity
- first/latest hash alignment
- attestation mode restricted to `unsigned-dev` or `future-signable`
- attestation signature value must remain empty

## Security posture

This is a local deterministic hash chain for audit correlation.

It is **not**:

- a cryptographic signature
- a PKI or key-backed attestation
- a KMS integration
- an apply/execute runtime control path

## Runtime boundary

Bundle indexing does not enable:

- CPU/RAM apply
- actuator apply
- cluster mutation calls

## Future work

- signed attestation integration (separate trust plane)
- append-only external store for chain anchoring
- multi-run bundle append CLI workflows
