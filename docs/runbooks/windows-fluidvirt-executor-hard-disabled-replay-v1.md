# windows-fluidvirt-executor-hard-disabled-replay-v1

Local replay runbook for the hard-disabled future apply executor.

## Safety posture

- hard-disabled executor only
- no runtime mutation
- no QMP commands are sent
- no CPU/RAM apply is performed
- local replay only (no cluster access)

## Fixture replay

Replay a fixture:

`go run ./cmd/karl-fluid-executor -fixture examples/windows-fluid-executor-fixtures/executor-master-win11-cpu.hard-disabled.json -evaluation-time 2026-05-07T16:00:00Z`

Expected outcomes:

- `EXECUTION_HARD_DISABLED`
- `EXECUTION_BLOCKED`
- `EXECUTION_QUARANTINED`
- `EXECUTION_DENIED`

No `ACTIVE`, `APPLYING`, or `EXECUTED` states are allowed.

## Direct JSON replay

`go run ./cmd/karl-fluid-executor -governance <contract.json> -revalidation <revalidation.json> -attestation <attestation.json> [-killswitch <killswitch.json>]`

If `-killswitch` is missing, execution must remain blocked.

## Denial proof checklist

Verify output contains:

- `applyAttempted=false`
- `mutationPerformed=false`
- `qmpCommandSent=false`
- `clusterMutationSent=false`
- `commandEnvelope.commandPreviewOnly=true`
- `commandEnvelope.containsExecutableCommand=false`
- empty `qmpCommands`, `clusterMutations`, `guestMutations`

## What remains forbidden

- any runtime CPU/RAM apply
- any mutating QMP command
- any cluster mutation call
- any deploy action
- any dashboard/frontend scope change
