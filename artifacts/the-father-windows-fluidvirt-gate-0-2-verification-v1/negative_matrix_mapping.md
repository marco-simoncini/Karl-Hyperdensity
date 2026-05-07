# Negative Matrix Mapping

Exhaustive mapping tests cover required negative cases including:

- QMP/guest missing and malformed conditions
- identity drift and quarantine paths
- rollback/return and memory safety blockers
- migration/VMIM/pool/generic-target blockers
- attestation missing/malformed/replayed/stale
- accidental execution-enable flags

Each case validates gate target, status, blocker, and hard-disabled invariants.
