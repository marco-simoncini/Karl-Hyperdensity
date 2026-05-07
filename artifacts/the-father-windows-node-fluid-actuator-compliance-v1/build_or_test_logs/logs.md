# Validation Logs

`go test ./...`

- PASS for `pkg/contracts`
- PASS for `pkg/windowsfluidvirt`
- PASS for `pkg/windowsfluidvirt/sidecar`

`python3 scripts/validate_json.py`

- `[validate_json] OK: parsed 24 schema files and 20 example files`

`git diff --check`

- no whitespace errors

Safety scan:

- scoped scan completed
- only benign substring matches (`RequireActuatorAck`, safety documentation strings)
