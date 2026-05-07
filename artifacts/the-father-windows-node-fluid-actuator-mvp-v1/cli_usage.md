# CLI Usage

Binary:

- `cmd/karl-node-fluid-actuator`

Supported flags:

- `-mode dry-run|apply|rollback|return-to-floor`
- `-request <path>`
- `-allowlist <path>`
- `-kill-switch <path>`
- `-evidence-out <path>`
- `-evaluation-time <rfc3339>`
- optional `-dry-run`

Constraints:

- local file operations only
- no cluster calls
- no QMP calls
- write scope constrained to allowlisted `cpu.max`
