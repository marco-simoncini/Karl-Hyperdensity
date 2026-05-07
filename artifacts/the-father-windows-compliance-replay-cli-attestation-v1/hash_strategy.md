# Hash Strategy

Deterministic local hashing is used for replay audit references:

- `evidenceHash`: SHA-256 over canonical replay input JSON
- `replayHash`: SHA-256 over canonical replay output snapshot

Canonicalization relies on deterministic JSON marshaling for stable payloads with fixed evaluation time.

These hashes are:

- replay/audit references
- not signatures
- not trust anchors
- not cryptographic attestation proof
