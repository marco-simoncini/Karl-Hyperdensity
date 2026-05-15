# Test Summary

Added coverage for:

- CLI replay ready standalone
- CLI replay ready pool-child
- CLI replay blocked pool-scaling
- deterministic output with fixed evaluation time
- stable `replayHash` and `evidenceHash`
- attestation emission and mode checks
- invalid attestation mode rejection
- mutation flags all false

Validation commands also executed:

- `go test ./...`
- `python3 scripts/validate_json.py`
- `git diff --check`
