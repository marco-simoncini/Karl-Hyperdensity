# Test Summary

Added test coverage for:

- single-run bundle deterministic output
- hash presence (`evidenceHash`, `replayHash`, `attestationHash`, `runHash`)
- `runHash` dependence on `previousRunHash`
- first run previous hash empty
- two-run chain valid
- broken previous hash invalid
- latest hash and latest phase semantics
- runCount mismatch invalid
- latest ready/blocked status correctness
- attestation mode restrictions and empty signature value
- CLI deterministic output with `-emit-bundle-index`
