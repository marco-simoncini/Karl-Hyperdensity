# Validation Summary

Executed:

- `go test ./pkg/windowsfluidvirt -run "TestRealEvidenceMasterWin11Ready|TestRealEvidenceMasterWin11PoolChildReady|TestRealEvidenceMasterWin11PoolScalingMechanismBlocked|TestRealEvidenceMissingCPUActuatorBlocked|TestRealEvidenceMissingRAMBalloonBlocked" -v`
- `go test ./...`
- `python3 scripts/validate_json.py`
- `git diff --check`

All commands passed for this replay update.
