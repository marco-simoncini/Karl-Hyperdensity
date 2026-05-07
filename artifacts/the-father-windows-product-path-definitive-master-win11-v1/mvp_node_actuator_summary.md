# MVP Node Actuator Summary

MVP component implemented:

- `cmd/karl-node-fluid-actuator`
- contract models in `pkg/windowsfluidvirt/node_actuator_mvp.go`

MVP capabilities used in this run:

- dry-run policy validation
- apply (`cpu.max`)
- return-to-floor (`cpu.max`)
- rollback (`cpu.max`)
- evidence JSON output
- request/allowlist identity checks
- stale request check
- kill-switch support

MVP was executed inside a temporary privileged node-local pod pinned to `karl-lab-metal-01`, then removed after test.
