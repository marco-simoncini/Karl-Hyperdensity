# CLI Usage

New flags:

- `-bundle-index`
- `-bundle-subject`
- `-previous-run-hash`
- `-emit-bundle-index`

Example:

`go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-real-evidence.ready.json -evaluation-time 2026-05-07T21:30:00Z -emit-attestation -attestation-mode future-signable -emit-bundle-index -bundle-subject windows-shell/karl/master-win11 -pretty`

Output can include:

- replay result
- optional attestation envelope
- optional single-run bundle index
