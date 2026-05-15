# windows-fluid-apply-governance-v1

Formal governance contract for future apply phases, without runtime execution.

## Why admission is not apply

- Admission is a policy gate for eligibility review only.
- Admission output never executes runtime mutations.
- `ADMITTED_FOR_FUTURE_APPLY` does not authorize immediate apply.

## Why governance contract is not apply

- Governance contract is a formal prerequisite bundle.
- Contract output does not execute QMP actions, CPU/RAM resize, or controller calls.
- `CONTRACT_PREPARED` means "contract complete", not "apply executable".

## Governance phases

- `CONTRACT_PREPARED`
- `CONTRACT_BLOCKED`
- `CONTRACT_QUARANTINED`
- `NEEDS_REVALIDATION`

No `ACTIVE`, no `APPLYING`, no apply success states in this phase.

## Transition proofs

`WindowsFluidTransitionProof` records:

- transition source and target phases;
- required/observed/missing inputs;
- blocker list and invariant check snapshot;
- deterministic proof timestamp/hash for audit.

The model supports theoretical transition framing such as `FUTURE_APPLY_ELIGIBLE`, but still non-executable in this phase.

## Runtime invariant set

`WindowsFluidRuntimeInvariantSet` captures mandatory continuity and readiness invariants:

- no migration/VMIM/recreate/rollout/reboot;
- same node/pod/qemu/boot/machine identity;
- QMP ACK + guest ACK;
- rollback + return-to-floor readiness;
- kill-switch readiness;
- evidence freshness;
- QMP read-only until future apply phase.

Identity invariant failures lead to quarantine; readiness failures lead to blocked.

## Pre-apply revalidation contract

`WindowsFluidPreApplyRevalidationContract` requires fresh evidence and unchanged identity/runtime comparisons before any future apply review.

Allowed outputs:

- `REVALIDATION_READY`
- `REVALIDATION_BLOCKED`
- `REVALIDATION_QUARANTINED`
- `REVALIDATION_STALE`

## Rollback and return-to-floor requirements

- Rollback requirement is mandatory.
- Return-to-floor requirement is mandatory.
- Missing either requirement blocks governance contract preparation.

## Kill-switch requirements

- Kill-switch readiness is required at governance level.
- Mutation remains disabled in this phase (`mutationAllowed=false`, `applyAllowed=false`).

## Policy attestation

`WindowsFluidPolicyAttestation` captures attestable snapshots for:

- governance contract
- transition proof
- invariant set
- pre-apply revalidation context

Signature mode is `unsigned-dev` or `future-signable`; no real signing/key management in this phase.

## QMP command rule in this phase

- QMP mutating commands remain forbidden.
- QMP must stay read-only until a future explicit apply-phase implementation exists.

## Prerequisites before first future +CPU trial

- admitted decision for CPU path
- contract prepared and not blocked/quarantined/stale
- fresh identity/QMP/guest evidence
- rollback and return-to-floor ready
- zero P0/P1 blockers
- explicit separate apply executor phase (out of scope here)
