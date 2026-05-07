# QMP Evidence Contract

Implemented in:

- `pkg/windowsfluidvirt/qmp_contract.go`
- `schemas/windows-fluid-qmp-evidence-v1.schema.json`
- `docs/contracts/windows-fluid-qmp-evidence-v1.md`

Mandatory contract controls:

- `qmpReadOnly=true`
- read-only command execution only
- mutating command attempts are rejected
- QMP connection/capabilities/errors map to canonical blockers:
  - `qmp_socket_unavailable`
  - `qmp_ack_missing`
  - `hotplug_error_detected` (denylist violation)
