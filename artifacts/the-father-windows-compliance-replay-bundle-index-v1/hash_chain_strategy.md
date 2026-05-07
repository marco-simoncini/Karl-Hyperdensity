# Hash Chain Strategy

Chain mode: `local-deterministic-hash-chain`

For each run:

- compute `evidenceHash` from replay input
- compute `replayHash` from replay output contract
- compute optional `attestationHash`
- compute `runHash` from run payload including `previousRunHash`

Bundle chain:

- first run has empty `previousRunHash`
- each subsequent run must reference previous `runHash`
- validator recomputes run hashes and link consistency

This is an audit hash chain, not a cryptographic signature system.
