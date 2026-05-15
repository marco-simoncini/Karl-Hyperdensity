# windows-fluidvirt-governed-inplace-resource-lease-claim-draft-v1

Draft tecnico per revisione brevettuale.

## Technical framing

Proposed technical scope:

- certified Windows VM represented as a fluid shell;
- in-place runtime resource lease concept for CPU/RAM;
- evidence-governed lifecycle across dry-run, admission, governance contract, and future apply review;
- strict continuity constraints:
  - same node
  - same virt-launcher pod
  - same QEMU process
  - same Windows boot
  - same machine identity
- explicit no-migration/no-reboot/no-recreate constraints;
- mandatory guest ACK and QMP evidence;
- mandatory rollback and return-to-floor readiness.

## Candidate technical claim themes (non-legal)

1. **Governed In-Place Lease Envelope**
   - A method where resource lease intents are evaluated against continuity evidence before any execution path is permitted.

2. **Transition Proof-Gated Eligibility**
   - A deterministic transition proof model linking dry-run/admission/governance states and blocking execution when invariants are violated or stale.

3. **Invariant-Backed Runtime Safety Lattice**
   - A mandatory invariant set with quarantine vs blocked semantics tied to identity and readiness failures.

4. **Pre-Apply Revalidation Contract**
   - A freshness and identity comparison contract that must pass immediately before any future apply executor is even considered.

5. **Attestation-Ready Governance Snapshot**
   - A policy attestation model designed for future signing, capturing decision snapshots and blocker/invariant evidence.

## Novelty-oriented technical differentiators to evaluate

- single-node, same-process continuity proofs as first-class runtime constraints;
- lease governance that remains non-executable until all proof layers align;
- explicit separation between:
  - admission eligibility
  - governance preparedness
  - future execution phase (not implemented here).

## Out-of-scope in this draft

- no legal assertions of patent grantability;
- no final legal claim language;
- no production readiness claim;
- no runtime execution implementation in this phase.
