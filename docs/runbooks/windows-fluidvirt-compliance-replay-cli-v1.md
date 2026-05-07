# Windows FluidVirt Compliance Replay CLI Runbook v1

## Purpose

Run deterministic read-only compliance replay for Windows FluidVirt evidence.

## Basic usage

Standalone ready replay:

`go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-real-evidence.ready.json -evaluation-time 2026-05-07T21:00:00Z -pretty`

Pool-child ready replay:

`go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-pool-child-real-evidence.ready.json -evaluation-time 2026-05-07T21:00:00Z -pretty`

Pool-scaling blocked replay:

`go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-pool-scaling-mechanism.blocked.json -evaluation-time 2026-05-07T21:00:00Z -pretty`

Emit future-signable attestation:

`go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-real-evidence.ready.json -evaluation-time 2026-05-07T21:00:00Z -emit-attestation -attestation-mode future-signable -pretty`

## Result interpretation

- `HYPERDENSITY_READY_WINDOWS_SHELL`: ready shell target.
- `BLOCKED_WITH_REMEDIATION`: not ready yet; use blockers/remediation lists.
- pool-child may be ready when pool is provisioning context only.
- pool scaling as mechanism must remain blocked.

## Attestation interpretation

- `signature.mode` may be `unsigned-dev` or `future-signable`.
- `signature.value` must remain empty.
- hashes are deterministic local audit references and are not real signatures.

## Forbidden actions

- no runtime mutation
- no CPU/RAM apply
- no actuator apply
- no cluster mutation/deploy
- no vCPU hotplug claim
- no production-ready claim
