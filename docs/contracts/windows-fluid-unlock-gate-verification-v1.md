# windows-fluid-unlock-gate-verification-v1

Deterministic non-executable verification contract for readiness Gate 0, Gate 1, Gate 2, and Hyperdensity parity completion.

## Purpose

- Verify readiness gates without enabling apply execution.
- Keep executor hard-disabled by contract.
- Produce deterministic replayable gate decisions.

## Gate definitions

### Gate 0 — executor hard-disabled verified

Passes only when executor proof confirms:

- `applyAttempted=false`
- `mutationPerformed=false`
- `qmpCommandSent=false`
- `clusterMutationSent=false`
- guard keeps `executorEnabled=false`
- guard keeps mutation window and mutation flags closed
- command envelope remains non-executable and empty
- blocker `future_apply_executor_disabled` is present

### Gate 1 — lab read-only evidence complete (`master-win11`)

Passes only when:

- target is `master-win11` candidate context
- pool replicas are not used as target
- identity evidence is complete
- QMP evidence is present and read-only
- guest evidence is present and healthy
- no pending reboot
- no migration/VMIM/recreate/rollout evidence
- continuity proofs hold (`sameNode`, `samePod`, `sameQemu`)
- rollback and return-to-floor readiness are present
- evidence freshness is valid

### Gate 2 — future-signable attestation replay verification

Passes only when:

- signature mode is `future-signable` or `unsigned-dev`
- `signature.value` is empty in this phase
- subject refs and evidence refs are coherent
- deterministic replay hash can be reproduced
- attestation is neither stale nor replayed nor malformed

### Hyperdensity parity gate — `GATE_HYPERDENSITY_PARITY_COMPLETE`

This is the final non-numbered parity gate used to validate Hyperdensity completeness.

Passes only when all four parity proofs are present and complete:

- `cpu_scale_up`
- `cpu_scale_down`
- `ram_scale_up`
- `ram_scale_down`

Each proof must confirm all invariants:

- QMP-confirmed runtime state
- Windows guest-confirmed actual state
- same VM, namespace, node, virt-launcher pod, QEMU process, Windows boot, and machine identity
- no reboot, no rollout, no recreate, no migration, no destructive migration
- rollback verified
- return-to-floor verified
- evidence-backed audit

If any one proof is missing or fails, parity is blocked and emits:

- `hyperdensity_parity_partial_success_not_total_feasibility`

`partial success` is not `total feasibility`.

### Deprecated legacy alias

- `Gate3HyperdensityParity` is accepted only as a deprecated compatibility alias.
- The alias is normalized to `GATE_HYPERDENSITY_PARITY_COMPLETE` during evaluation.
- The alias does not redefine architectural Gate 3 semantics.
- Architectural Gate 3 remains reserved for QMP sidecar read-only socket proof.

## Gate status semantics

- `PASSED`
- `FAILED`
- `BLOCKED`
- `QUARANTINED`
- `NOT_APPLICABLE`

`PASSED` means the gate verification condition is met for readiness evidence.

`PASSED` does **not** mean:

- executor unlock is granted
- runtime mutation is allowed
- CPU or RAM apply is authorized

## Aggregate gate-set semantics

`GATE_SET_PASSED` means Gate 0, Gate 1, Gate 2, and `GATE_HYPERDENSITY_PARITY_COMPLETE` are verified for current replay inputs.

`GATE_SET_PASSED` does **not** mean unlock execution.

Executor remains hard-disabled and non-executable.

## Negative matrix mapping

Gate evaluation maps negative cases to deterministic blockers and statuses:

- missing QMP or guest evidence -> `BLOCKED`
- identity drift (`node/pod/qemu/boot/machine`) -> `QUARANTINED`
- malformed or replayed attestation -> `BLOCKED`
- any accidental execution-enabling flag -> `FAILED`

## Safety invariants

- `executorMustRemainDisabled=true`
- `mutationAllowed=false`
- `applyAllowed=false`
- no runtime mutation path is introduced

## Prerequisites before next prompt

- deterministic fixtures for Gate 0/1/2 and GateSet are green
- negative mapping tests are exhaustive and green
- replay logs are archived in artifacts
- unresolved blockers are recorded as unknowns, not bypassed
