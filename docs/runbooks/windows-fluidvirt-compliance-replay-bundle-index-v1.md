# Windows FluidVirt Compliance Replay Bundle Index Runbook v1

## Goal

Generate deterministic replay bundle index outputs for continuous audit of Windows compliance replay.

## Single-run bundle index generation

Master-win11 ready:

`go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-real-evidence.ready.json -evaluation-time 2026-05-07T21:30:00Z -emit-attestation -attestation-mode future-signable -emit-bundle-index -bundle-subject windows-shell/karl/master-win11 -pretty`

Pool-child ready:

`go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-pool-child-real-evidence.ready.json -evaluation-time 2026-05-07T21:30:00Z -emit-attestation -attestation-mode future-signable -emit-bundle-index -bundle-subject windows-shell/karl/master-win11 -pretty`

Pool-scaling blocked:

`go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-pool-scaling-mechanism.blocked.json -evaluation-time 2026-05-07T21:30:00Z -emit-attestation -attestation-mode future-signable -emit-bundle-index -bundle-subject windows-shell/karl/master-win11 -pretty`

## Chain validation

Use `ValidateWindowsComplianceReplayBundleIndex` from backend tests/helpers to validate:

- run linkage via `previousRunHash`
- deterministic run hash recomputation
- first/latest hash consistency
- attestation mode/value safety constraints

## Status interpretation

- `latestCompliancePhase` + `latestHyperdensityReady` reflect latest run outcome
- `chainValid=true` indicates deterministic chain integrity across provided runs
- `chainValid=false` with `brokenAtRunId` identifies the first inconsistent link

## Continuous audit usage

Use bundles as replay audit snapshots in chronological order with fixed evaluation times for deterministic comparison.

## Forbidden operations

- runtime mutation
- CPU/RAM apply
- actuator apply
- cluster mutation/deploy
- any production-ready signature claim
