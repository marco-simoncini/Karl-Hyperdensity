# windows-fluid-qmp-command-policy-v1

QMP command policy for `karl-fluid-sidecar` in evidence-only mode.

## Allowlist (read-only)

- `qmp_capabilities`
- `query-status`
- `query-cpus-fast`
- `query-hotpluggable-cpus`
- `query-memory-devices`
- `query-memory-size-summary`
- `query-machines`
- `query-version`

## Forbidden (must be rejected)

- `device_add`
- `device_del`
- `qom-set`
- `set_link`
- `system_powerdown`
- `stop`
- `cont`
- `quit`
- `migrate`
- `migrate_cancel`
- `object-add`
- `object-del`
- `cpu-add`
- `balloon`
- `memsave`
- `pmemsave`

## Enforcement

- Requests outside allowlist are hard-rejected.
- Rejected command attempts are captured in sidecar error evidence.
- No mutating command execution is allowed in this phase.
