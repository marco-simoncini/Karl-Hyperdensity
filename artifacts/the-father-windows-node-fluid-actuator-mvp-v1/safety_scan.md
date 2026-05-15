# Safety Scan

Scan executed on modified actuator, replay, docs, fixtures, and artifact files using required patterns.

Result:

- no forbidden matches found for secrets/credentials patterns
- no forbidden frontend/dashboard terms detected
- no forbidden port mentions (`443`, `8888`) detected
- no forbidden success claims detected

Scan scope included:

- `cmd/karl-node-fluid-actuator/main.go`
- `pkg/windowsfluidvirt/node_actuator_mvp*.go`
- `pkg/windowsfluidvirt/compliance_replay_cli*.go`
- `cmd/karl-fluid-compliance-replay/main.go`
- `docs/contracts/windows-fluid-node-actuator-v1.md`
- `docs/runbooks/windows-fluidvirt-node-actuator-mvp-v1.md`
- `examples/windows-fluid-actuator-mvp-fixtures/*.json`
- `artifacts/the-father-windows-node-fluid-actuator-mvp-v1/**`
