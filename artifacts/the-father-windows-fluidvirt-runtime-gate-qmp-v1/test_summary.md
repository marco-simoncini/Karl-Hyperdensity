# Test Summary

## Hyperdensity

- `go test ./...` -> PASS
- `python3 scripts/validate_json.py` -> PASS
- `git diff --check` + `git diff --cached --check` -> PASS

Logs:

- `build_or_test_logs/hyperdensity-go-test.log`
- `build_or_test_logs/hyperdensity-validate-json.log`
- `build_or_test_logs/hyperdensity-diff-check.log`

## Sidecar/QMP fixture coverage

- handshake and capabilities success
- missing socket -> blocker path
- QMP error -> blocker path
- mutating command rejected
- evidence remains `qmpReadOnly=true`
