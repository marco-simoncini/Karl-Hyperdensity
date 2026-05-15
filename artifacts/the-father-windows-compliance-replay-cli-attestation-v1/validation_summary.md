# Validation Summary

Executed:

- `go test ./pkg/windowsfluidvirt -run "TestComplianceReplayCLI|TestComplianceReplayHashAndEvidenceHashStable|TestComplianceReplayAttestationFutureSignable|TestComplianceReplayInvalidAttestationModeRejected|TestComplianceReplayNoMutationFlagsTrue" -v`
- `go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-real-evidence.ready.json -evaluation-time 2026-05-07T21:00:00Z -pretty`
- `go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-pool-child-real-evidence.ready.json -evaluation-time 2026-05-07T21:00:00Z -pretty`
- `go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-pool-scaling-mechanism.blocked.json -evaluation-time 2026-05-07T21:00:00Z -pretty`
- `go run ./cmd/karl-fluid-compliance-replay -input examples/windows-fluid-compliance-fixtures/master-win11-real-evidence.ready.json -evaluation-time 2026-05-07T21:00:00Z -emit-attestation -attestation-mode future-signable -pretty`
- `go test ./...`
- `python3 scripts/validate_json.py`
- `git diff --check`

All passed for this scope.
