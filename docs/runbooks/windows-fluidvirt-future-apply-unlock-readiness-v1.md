# windows-fluidvirt-future-apply-unlock-readiness-v1

Runbook for using the future apply unlock readiness specification in strictly non-executable mode.

## Purpose

- Use the readiness specification to evaluate unlock prerequisites.
- Keep executor hard-disabled while preparing evidence and trust artifacts.
- Plan next lab prompts without introducing runtime mutation.

## How to use the readiness spec

1. Select target scope (`cpu-lease-only-lab`, `memory-lease-lab`, `return-to-floor-lab`, `rollback-lab`).
2. Map current evidence and blockers to readiness fields.
3. Confirm `executorMustRemainDisabled=true`.
4. Evaluate trust model readiness and keyless verification limits.
5. Evaluate unlock criteria matrix row for the selected scope.
6. Record missing proofs and unresolved blockers as explicit unknowns.

## Verify executor remains hard-disabled

- Run `go test ./...`.
- Replay executor fixtures with deterministic timestamp.
- Confirm outputs keep:
  - `applyAttempted=false`
  - `mutationPerformed=false`
  - `qmpCommandSent=false`
  - `clusterMutationSent=false`
- Confirm no path introduces executable QMP/cluster/guest mutation lists.

## How to use the matrices

- Use unlock criteria matrix to decide readiness state by scope.
- Use risk model to map hazards to mandatory blockers and mitigations.
- Use negative matrix to verify expected phase/decision behavior.
- Use rollout gates to sequence readiness into a separate future milestone.

## Preparing the next lab evidence prompt

- State selected gate and target scope.
- Require fresh evidence refs and deterministic evaluation timestamp.
- Require explicit kill switch proof and operator approval evidence.
- Require negative test subset relevant to that gate.
- Keep anti-goals explicit: no unlock, no apply, no runtime mutation.

## What remains forbidden

- enabling executor
- introducing runtime mutation paths
- CPU or RAM apply logic
- mutating QMP command flow
- deploy or cluster mutation actions
- frontend/dashboard scope changes

This runbook is readiness-only and does not authorize execution.
