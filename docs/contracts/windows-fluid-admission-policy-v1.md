# windows-fluid-admission-policy-v1

Admission policy contract for Windows FluidVirt future apply governance.

## Admission vs dry-run

- Dry-run (`windows-fluid-dryrun-evaluation-v1`) evaluates runtime evidence state.
- Admission evaluates whether a future apply request may be considered in a separate phase.
- Admission output is governance-only and never mutates runtime.

## Admission vs apply

- `ADMITTED_FOR_FUTURE_APPLY` is not apply.
- `ADMITTED_FOR_FUTURE_APPLY` only means prerequisites are strong enough for a future gated apply review.
- This contract never emits `ACTIVE`, `APPLYING`, or apply success.

## Policy pack

`WindowsFluidPolicyPack` defines conservative controls:

- `allowedRuntimeMode=in-place-qmp`
- no migration/reboot/recreate
- same node/pod/qemu/boot/machine identity requirements
- mandatory QMP ACK and guest ACK
- mandatory rollback and return-to-floor readiness
- `allowPoolReplicaModel=false`
- `allowGenericWindowsVm=false`
- `allowMutationInThisPhase=false`

## Evidence scoring

Scoring produces:

- numeric score (0-100)
- evidence level:
  - `insufficient`
  - `partial`
  - `dryrun-ready`
  - `future-apply-admissible`
- missing evidence list
- hard blockers and soft unknowns

Hard blockers always override high scores.

## Blocker priority

- `P0_QUARANTINE` => `QUARANTINED`
- `P1_HARD_BLOCK` => `BLOCKED`
- `P2_CAPABILITY_BLOCK` => `BLOCKED` or `NEEDS_MORE_EVIDENCE`
- `P3_ENVIRONMENT_BLOCK` => release/deployment gate concerns, not runtime mutation path

## Blast-radius, rollback, return-to-floor

- Blast radius is constrained to single VM governance.
- Rollback readiness is mandatory for future apply admission.
- Return-to-floor readiness is mandatory, with memory safety gates for RAM intents.

## Denied models

- Pool replica model (`win11-pool-*`) is context-only and denied as FluidVirt mechanism.
- Generic Windows VM without certified annotations/contracts is denied.

## Prerequisites before future +CPU phase

- dry-run must be `READY` or `LEASE_PREPARED`
- complete continuity proofs (node/pod/qemu/boot/machine)
- no migration/recreate/reboot
- guest ACK + QMP read-only/ack
- rollback and return-to-floor ready
- evidence score above policy threshold
