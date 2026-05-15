# Validation Summary

Executed:

- `go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-real-evidence.ready.json -evaluation-time 2026-05-07T22:00:00Z -emit-attestation -attestation-mode future-signable -emit-bundle-index -bundle-subject windows-shell/karl/master-win11 -pretty`
- `go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-pool-child-real-evidence.ready.json -evaluation-time 2026-05-07T22:00:00Z -emit-attestation -attestation-mode future-signable -emit-bundle-index -bundle-subject windows-shell/karl/master-win11 -pretty`
- `go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-pool-scaling-mechanism.blocked.json -evaluation-time 2026-05-07T22:00:00Z -emit-attestation -attestation-mode future-signable -emit-bundle-index -bundle-subject windows-shell/karl/master-win11 -pretty`
- `go test ./...`
- `python3 scripts/validate_json.py`
- `git diff --check`
- replay CLI runs with `-emit-attestation -emit-bundle-index`
- deterministic replay output comparison with fixed `-evaluation-time`

All checks passed.
