# QMP Sidecar Summary

Read-only skeleton implemented:

- `cmd/karl-fluid-sidecar/main.go`
- `pkg/windowsfluidvirt/sidecar/readonly_executor.go`
- `pkg/windowsfluidvirt/sidecar/socket_transport.go`

Capabilities:

- QMP socket path input
- greeting/capabilities flow
- read-only query execution
- evidence emission in JSON
- command deny on non-allowlisted operations

No mutating QMP action is executable in this milestone.
