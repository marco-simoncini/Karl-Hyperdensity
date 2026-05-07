# windows-fluid-future-apply-unlock-readiness-v1

Formal non-executable readiness specification for a future Windows FluidVirt apply-unlock milestone.

## Safety baseline

- Executor remains hard-disabled in this milestone.
- No runtime mutation path is introduced.
- No CPU or RAM apply logic is implemented.
- No QMP mutating command flow is enabled.

## WindowsFluidFutureApplyUnlockReadiness

This specification defines the conceptual model `WindowsFluidFutureApplyUnlockReadiness`.

Required fields:

- `readinessId`
- `targetScope` in:
  - `cpu-lease-only-lab`
  - `memory-lease-lab`
  - `return-to-floor-lab`
  - `rollback-lab`
- `unlockStatus` in:
  - `NOT_READY`
  - `SPEC_READY`
  - `BLOCKED`
  - `QUARANTINED`
- `executorMustRemainDisabled=true`
- `requiredMilestones`
- `requiredProofs`
- `requiredAttestations`
- `requiredOperatorControls`
- `requiredKillSwitchControls`
- `requiredRollbackProofs`
- `requiredReturnToFloorProofs`
- `requiredGuestProofs`
- `requiredQmpProofs`
- `requiredClusterProofs`
- `requiredNegativeTests`
- `explicitNonGoals`
- `blockers`
- `createdAt`

### Meaning of `SPEC_READY`

- Means the formal unlock checklist exists and is auditable.
- Does not mean runtime unlock has happened.
- Does not permit executor activation.
- Does not permit CPU or RAM runtime apply.

## Signed-attestation trust model

Conceptual model: `WindowsFluidSignedAttestationTrustModel` (future-signable only in this phase).

### Attestation subjects

- dry-run result
- admission decision
- governance contract
- invariant set
- pre-apply revalidation
- executor guard
- kill switch state

### Trust boundaries

- Hyperdensity controller
- `karl-fluid-sidecar`
- KARL Agent `fluidShell`
- operator policy plane
- audit store

### Future signer roles

- controller signer
- sidecar signer
- guest agent signer
- policy signer

### Verification requirements

- canonical subject serialization
- signer role-to-subject binding
- freshness window checks
- replay detection across attestation IDs and proof hashes
- deny-on-missing-link behavior in trust chain

### Replay protection and freshness

- attestations must include deterministic content hash references
- each attestation must bind evidence window and timestamp
- stale or replayed attestations are hard blockers

### Future-only security operations

- key rotation plan is mandatory before unlock
- revocation propagation is mandatory before unlock
- `unsigned-dev` compatibility is kept only for pre-unlock replay

No KMS, certificates, private keys, or tokens are introduced in this milestone.

## Keyless verification strategy

The readiness phase uses deterministic local verification, not cryptographic signing.

### Strategy elements

- content-addressed evidence objects
- canonical JSON normalization before hashing
- deterministic hash generation for replay
- chained evidence references (`evidenceRef`, `proofHash`, `auditHash`)
- strict timestamp windows for freshness validation
- append-only audit assumptions for replay history
- local replay verification via CLI and fixture packs

### Trust-on-first-proof limitations

- deterministic hashes provide integrity hints, not signer identity
- hash continuity does not replace key-backed attestation
- local replay verification is not sufficient for production execution unlock

## Unlock criteria matrix

### 1) CPU +lease lab unlock

- Status: `SPEC_READY`
- Required evidence: fresh identity, QMP read-only truth, guest ACK, rollback/return proofs
- Required negative tests: stale evidence, identity mismatch, kill switch missing
- Kill switch: required and operational proof mandatory
- Operator approval: explicit per-run
- Blast radius: single VM only
- Blockers: any P0/P1 blocker, missing trust-chain link
- Why not ready today: execution milestone intentionally separated

### 2) CPU return-to-floor lab unlock

- Status: `NOT_READY`
- Required evidence: post-lease convergence and safe return evidence
- Required negative tests: return-to-floor unsafe, rollback unavailable
- Kill switch: required
- Operator approval: required
- Blast radius: single VM only
- Blockers: return proof missing
- Why not ready today: return procedure verification pack incomplete

### 3) RAM +lease lab unlock

- Status: `NOT_READY`
- Required evidence: memory safety and driver support proofs
- Required negative tests: memory driver unsupported, stale memory truth
- Kill switch: required
- Operator approval: required
- Blast radius: single VM only
- Blockers: `memory_driver_unverified`, `memory_return_not_safe`
- Why not ready today: memory safety proof set incomplete

### 4) RAM return-to-floor lab unlock

- Status: `NOT_READY`
- Required evidence: validated safe memory return sequence
- Required negative tests: return failure and memory safety regressions
- Kill switch: required
- Operator approval: required
- Blast radius: single VM only
- Blockers: any memory return risk
- Why not ready today: cannot proceed without proven memory safety

### 5) rollback lab unlock

- Status: `NOT_READY`
- Required evidence: deterministic rollback path and verification trail
- Required negative tests: rollback unavailable, rollback partial success
- Kill switch: required
- Operator approval: required
- Blast radius: single VM only
- Blockers: rollback proof missing
- Why not ready today: dedicated rollback drills still pending

### 6) multi-VM batch unlock

- Status: `NOT_READY`
- Required evidence: per-VM trust chain and aggregate blast-radius controls
- Required negative tests: wrong-target and cross-VM contamination tests
- Kill switch: required per cohort
- Operator approval: required multi-party
- Blast radius: tightly bounded, non-default
- Blockers: any unresolved single-VM blocker
- Why not ready today: architecture is currently single-node/single-VM safety-first

### 7) production autonomous apply unlock

- Status: `NOT_READY`
- Required evidence: full trust-chain signing, hardened controls, complete negative matrix
- Required negative tests: comprehensive adversarial and chaos sets
- Kill switch: mandatory and continuously verifiable
- Operator approval: governance board-level gate
- Blast radius: production policy-defined
- Blockers: any missing readiness dimension
- Why not ready today: explicitly out of scope for this program phase

No automatic unlock is allowed.

## Execution unlock risk model

For every risk, unlock remains blocked unless mitigation and negative proof exist.

- wrong target VM
- stale evidence
- QEMU PID mismatch
- Windows reboot undetected
- guest ACK spoof or missing
- QMP socket spoof or wrong socket
- return-to-floor unsafe
- rollback unavailable
- memory driver unsupported
- LiveMigration or VMIM race
- virt-launcher pod recreate
- node drain race
- operator bypass
- kill switch unavailable
- evidence store unavailable
- clock skew
- replay attack

Each risk must define: severity, likelihood, detectability, mandatory blocker, mitigation, negative test, and unlock implication.

## Formal safety test plan

### Categories

- unit tests
- fixture replay tests
- sidecar QMP mock tests
- guest evidence mock tests
- integration lab tests
- chaos and negative tests
- stale evidence tests
- identity mismatch tests
- kill switch tests
- rollback tests
- return-to-floor tests
- audit immutability tests
- no-mutation dry-run tests
- executor-remains-disabled tests

### Status split

- Already implemented: hard-disabled executor behavior, deterministic replay, no-mutation assertions.
- Missing for first CPU lab apply proposal: signed trust-chain verification simulations, expanded negative replay set, operator approval trace tests.
- Missing before any RAM path: memory safety and return-to-floor proof suites, failure-injection for memory path, rollback compatibility tests.

## Negative test matrix baseline

Mandatory cases include:

- missing QMP
- wrong QMP socket
- QEMU PID changed
- virt-launcher pod changed
- node changed
- last boot changed
- machine identity changed
- pending reboot
- guest ACK missing
- agent module missing
- rollback not ready
- return-to-floor not ready
- memory driver unverified
- critical Windows event
- LiveMigration object observed
- VMIM observed
- pool replica target
- generic Windows VM
- stale evidence
- kill switch missing
- attestation missing
- attestation malformed
- replayed old evidence

For each case, the matrix must define expected decision, blocker, phase, and whether quarantine is required.

## Rollout gates for future apply milestone

- Gate 0: hard-disabled executor proof
- Gate 1: read-only evidence completeness on `master-win11`
- Gate 2: future-signable attestation replay verification
- Gate 3: sidecar live read-only QMP socket proof
- Gate 4: guest `fluidShell` live ACK proof
- Gate 5: fresh pre-apply revalidation proof
- Gate 6: kill switch operational proof
- Gate 7: rollback and return-to-floor lab proof
- Gate 8: first CPU lease apply proposal review
- Gate 9: first CPU lease apply in separate milestone only
- Gate 10: CPU return-to-floor lab
- Gate 11: RAM research only
- Gate 12: autonomous apply remains disabled

Each gate requires entry criteria, exit criteria, blockers, evidence artifacts, and explicit anti-goals.

## What is required before first CPU path

- full Gate 0 through Gate 8 completion
- formal trust-chain readiness with future-signable compatibility
- expanded negative matrix pass set
- verified kill switch and operator control proofs

This document remains non-executable and does not modify runtime behavior.
