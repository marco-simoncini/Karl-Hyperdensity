# transition_proof_model

Model: `WindowsFluidTransitionProof`

Captured fields:

- from/to phase
- allowed flag and reason
- required/observed/missing inputs
- blocker list
- invariant check snapshot
- proof timestamp and deterministic proof hash

Behavior:

- no transition produces `ACTIVE` or `APPLYING`;
- stale evidence transitions to `NEEDS_REVALIDATION`;
- identity/P0 failures transition to `CONTRACT_QUARANTINED`;
- P1 failures transition to `CONTRACT_BLOCKED`;
- `FUTURE_APPLY_ELIGIBLE` is represented as theoretical transition context only.
