# Bundle Append Summary

Added multi-run append workflow for compliance audit chain.

Core function:

- `AppendWindowsComplianceReplayBundleRun(existing, run, evaluationTime)`

Behavior:

- validates existing chain before append
- rejects append if chain broken
- rejects duplicate run (`runId` or `runHash`)
- sets `previousRunHash` to prior `latestRunHash`
- recomputes deterministic appended `runHash`
- rebuilds bundle with incremented `runCount` and updated `latestRunHash`

CLI support:

- `karl-fluid-compliance-replay -append-bundle -append-bundle-in <bundle.json>`
