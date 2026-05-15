# test_summary

Executed:

- `go test ./...`
- `python3 scripts/validate_json.py`
- `git diff --check`
- governance CLI replay:
  - `governance-master-win11-cpu.contract-prepared.json`
  - `governance-master-win11-stale.needs-revalidation.json`

Logs:

- `build_or_test_logs/go-test.log`
- `build_or_test_logs/validate-json.log`
- `build_or_test_logs/git-diff-check.log`
- `build_or_test_logs/governance-master-win11-cpu-output.json`
- `build_or_test_logs/governance-master-win11-stale-output.json`
