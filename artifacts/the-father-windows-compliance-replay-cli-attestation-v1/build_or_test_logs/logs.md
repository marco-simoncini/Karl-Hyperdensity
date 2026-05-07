# Logs

Replay command outputs:

- `replay-master-win11-ready.json`
- `replay-pool-child-ready.json`
- `replay-pool-scaling-blocked.json`
- `replay-master-win11-attestation.json`
- `replay-deterministic-a.json`
- `replay-deterministic-b.json`

Validation logs:

- `go-test-all.log`
- `validate-json.log`
- `git-diff-check.log`

Determinism check:

- `cmp replay-deterministic-a.json replay-deterministic-b.json` passed.
