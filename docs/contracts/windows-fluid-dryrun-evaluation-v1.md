# windows-fluid-dryrun-evaluation-v1

Evidence-only dry-run contract for Windows FluidVirt runtime integration.

## Purpose

- Compose runtime evidence into one deterministic backend decision.
- Never mutate CPU/RAM, never issue mutating QMP commands, never declare runtime apply success.
- Produce only: `READY`, `BLOCKED`, `QUARANTINED`, `LEASE_PREPARED`.

## Evidence Bundle

`WindowsFluidRuntimeEvidenceBundle` aggregates:

- `WindowsFluidShell` declared/runtime target/actual signals.
- KubeVirt runtime identity continuity (`before` and `after`).
- Read-only QMP evidence.
- Guest ACK/runtime proof from `modules.fluidShell`.
- Optional lease intent (`prepare-cpu-lease`, `prepare-memory-lease`).
- Policy gates (annotations, runtime mode, pool-context flag).
- Observed blockers, timestamps, source metadata, sanitization status.

## Evaluation Input

- Bundle (mandatory).
- Optional evaluation options (timestamp for deterministic replay output).
- Optional lease intent embedded in the bundle.

## Evaluation Output

- `phase`: `READY | BLOCKED | QUARANTINED | LEASE_PREPARED`.
- `classification`: includes:
  - `READY_FOR_FLUID_SHELL_CERTIFICATION`
  - `BLOCKED_GENERIC_WINDOWS_VM`
  - `BLOCKED_POOL_REPLICA_MODEL`
  - `BLOCKED_MISSING_QMP`
  - `BLOCKED_MISSING_GUEST_ACK`
  - `QUARANTINED_IDENTITY_CHANGED`
- `conditions`, `blockers`, and evidence summary.
- Non-mutating Action Slate entry.
- Recommended next safe step.

## Blocker Semantics

- Missing guest evidence => `guest_ack_missing` (or `karl_agent_fluid_module_missing` when module declaration/proof is absent).
- Missing QMP evidence => `qmp_socket_unavailable` (or `qmp_ack_missing`).
- Identity continuity break (node/pod/qemu/boot/machine identity) => `QUARANTINED`.
- Incomplete evidence cannot produce `READY`.
- Pool replica context (`win11-pool-*`) is always `BLOCKED_POOL_REPLICA_MODEL`.

## Action Slate (Non-Mutating)

- `mutationAllowed=false`.
- `applyAllowed=false`.
- Runtime mode fixed to `in-place-qmp`.
- Contains only dry-run metadata/proofs/blockers; no runtime commands.

## READY vs LEASE_PREPARED

- `READY`: certification-ready evidence is complete and coherent.
- `LEASE_PREPARED`: same as `READY` plus a valid non-mutating lease intent with:
  - `rollbackReady=true`
  - `returnToFloorReady=true`
  - QMP/guest readiness true

`LEASE_PREPARED` is still dry-run only and does not imply `APPLYING`/`ACTIVE`.

## Why ACTIVE Is Impossible Here

- Dry-run pipeline never executes mutation paths.
- Action Slate explicitly disables apply.
- No QMP mutating commands are accepted or emitted.

## Generic Windows VM and Pool Rules

- Generic Windows VM without required annotations/certification signals => `BLOCKED_GENERIC_WINDOWS_VM`.
- `win11-pool-*` is context-only evidence, not FluidVirt mechanism => `BLOCKED_POOL_REPLICA_MODEL`.

## Proofs Required Before Future Apply Phase

- Stable same-node / same-pod / same-qemu continuity.
- No migration / no recreate / no reboot proofs.
- Guest ACK and read-only QMP evidence.
- Rollback readiness and return-to-floor readiness.
