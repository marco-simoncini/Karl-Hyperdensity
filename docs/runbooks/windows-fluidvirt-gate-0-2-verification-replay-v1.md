# windows-fluidvirt-gate-0-2-verification-replay-v1

Runbook for deterministic replay of Gate 0, Gate 1, Gate 2 verification in non-executable mode.

## Safety posture

- non-executable verification only
- executor remains hard-disabled
- no runtime mutation
- no QMP commands are sent
- no CPU/RAM apply is performed

## Fixture usage

Fixtures are under:

- `examples/windows-fluid-unlock-gate-fixtures/`

Core fixtures:

- Gate 0 pass/fail
- Gate 1 pass/blocked/quarantined/pool-context blocked
- Gate 2 pass/blocked malformed/blocked stale-replayed
- GateSet pass/quarantined

## CLI replay (`karl-fluid-gates`)

Single gate:

`go run ./cmd/karl-fluid-gates -fixture examples/windows-fluid-unlock-gate-fixtures/gate0-executor-hard-disabled.passed.json -mode gate -evaluation-time 2026-05-07T18:20:00Z`

GateSet:

`go run ./cmd/karl-fluid-gates -fixture examples/windows-fluid-unlock-gate-fixtures/gateset-0-1-2.passed.json -mode gateset -evaluation-time 2026-05-07T18:20:00Z`

## Interpreting statuses

- `PASSED`: gate conditions verified for replay input only
- `BLOCKED`: missing/invalid required readiness evidence
- `QUARANTINED`: identity continuity risk detected
- `FAILED`: safety flags contradicted hard-disable expectations

Even with `PASSED`, executor remains disabled and no unlock is granted.

## Verify hard-disabled behavior

Confirm output always keeps:

- `executorMustRemainDisabled=true`
- `mutationAllowed=false`
- `applyAllowed=false`

Also confirm no command envelope execution payload appears.

## Preparing first real `master-win11` evidence report

- keep report read-only and sanitized
- bind evidence to deterministic timestamp window
- include identity continuity and freshness proofs
- include Gate 1 completeness and Gate 2 attestation coherence checks
- record unresolved risks as blockers/unknowns

## What remains forbidden

- enabling executor
- introducing runtime mutation paths
- CPU/RAM apply execution
- mutating QMP commands
- deploy/cluster mutation operations
- frontend/dashboard scope changes
