# Build/Test Logs

Targeted replay command:

`go test ./pkg/windowsfluidvirt -run "TestRealEvidenceMasterWin11Ready|TestRealEvidenceMasterWin11PoolChildReady|TestRealEvidenceMasterWin11PoolScalingMechanismBlocked|TestRealEvidenceMissingCPUActuatorBlocked|TestRealEvidenceMissingRAMBalloonBlocked" -v`

Observed replay outputs:

- `master-win11-real-evidence` -> `HYPERDENSITY_READY_WINDOWS_SHELL`
- `master-win11-pool-child-real-evidence` -> `HYPERDENSITY_READY_WINDOWS_SHELL`
- `master-win11-pool-scaling-mechanism` -> `BLOCKED_WITH_REMEDIATION`
- `missing-cpu-actuator` -> `BLOCKED_WITH_REMEDIATION`
- `missing-ram-balloon` -> `BLOCKED_WITH_REMEDIATION`

Validation commands:

- `go test ./...` -> PASS
- `python3 scripts/validate_json.py` -> PASS
- `git diff --check` -> PASS
