# windows-fluid-future-apply-executor-v1

Formal skeleton for the future Windows FluidVirt apply executor, intentionally hard-disabled.

## Why this executor is hard-disabled

- This milestone is evidence-only and safety-first.
- Executor output is always non-executable.
- `applyAttempted`, `mutationPerformed`, `qmpCommandSent`, and `clusterMutationSent` are always `false`.

## Governance contract vs executor

- Governance contract proves preconditions and invariants.
- Executor consumes governance artifacts and emits denial evidence.
- `CONTRACT_PREPARED` still does not authorize runtime apply.

## Pre-apply guard

`WindowsFluidPreApplyGuard` binds:

- governance contract reference
- revalidation reference
- attestation reference
- kill-switch readiness proof

Guard phases:

- `GUARD_READY_BUT_EXECUTOR_DISABLED`
- `GUARD_BLOCKED`
- `GUARD_QUARANTINED`
- `GUARD_NEEDS_REVALIDATION`

The guard can be ready, but executor stays disabled.

## Kill switch model

`WindowsFluidKillSwitch` is mandatory for any future apply candidate.

In this phase defaults are:

- `enabled=true`
- `mode=hard-disabled`
- reason: `future apply executor disabled by policy`

If kill switch proof is missing/unobservable, execution result is blocked.

## Command envelope is preview-only

`WindowsFluidExecutorCommandEnvelope` is a formal envelope only:

- `commandPreviewOnly=true`
- `containsExecutableCommand=false`
- `qmpCommands=[]`
- `clusterMutations=[]`
- `guestMutations=[]`

No `device_add`, `qom-set`, `cpu-add`, `object-add`, `migrate`, or other mutating commands are present or executed.

## Attestation integration

Executor references `WindowsFluidPolicyAttestation` but does not sign.

Accepted signature modes:

- `unsigned-dev`
- `future-signable`

In this phase:

- `signature.value` must remain empty
- no private keys
- no certificates
- no KMS/token integration

## Why no QMP and no CPU/RAM apply

- Executor implementation denies execution by design.
- No QMP mutations are emitted.
- No cluster mutation path exists.
- No CPU/RAM runtime apply is implemented.

## Future unlock requirements

A separate milestone must explicitly introduce:

- signed attestation flow
- explicit mutation authorization model
- audited runtime mutation executor implementation
- dedicated safety and rollback verification pack

Until then, executor remains hard-disabled.
