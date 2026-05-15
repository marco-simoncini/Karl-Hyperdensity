# test_summary

Executed:

- `go test ./...`
- `python3 scripts/validate_json.py`
- `git diff --check`
- `go run ./cmd/karl-fluid-admission` on:
  - `admission-master-win11-cpu.future-apply-admissible.json`
  - `admission-win11-pool.denied.json`

Logs:

- `build_or_test_logs/go-test.log`
- `build_or_test_logs/validate-json.log`
- `build_or_test_logs/git-diff-check.log`
- `build_or_test_logs/admission-master-win11-cpu-output.json`
- `build_or_test_logs/admission-win11-pool-output.json`
