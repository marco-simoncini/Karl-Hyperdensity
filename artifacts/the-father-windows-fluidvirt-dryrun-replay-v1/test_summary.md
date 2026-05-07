# test_summary

Executed checks:

- `go test ./...` (pass)
- `python3 scripts/validate_json.py` (pass)
- `git diff --check` (pass)
- `go run ./cmd/karl-fluid-dryrun` replay for:
  - `master-win11-certification-ready.ready.json`
  - `win11-pool-context-only.blocked.json`

Stored logs:

- `build_or_test_logs/go-test.log`
- `build_or_test_logs/validate-json.log`
- `build_or_test_logs/git-diff-check.log`
- `build_or_test_logs/master-win11-dryrun-output.json`
- `build_or_test_logs/win11-pool-dryrun-output.json`
