# QMP Command Allowlist

Source:

- `pkg/windowsfluidvirt/qmp_contract.go`
- `docs/contracts/windows-fluid-qmp-command-policy-v1.md`

Read-only allowlist:

- `qmp_capabilities`
- `query-status`
- `query-cpus-fast`
- `query-hotpluggable-cpus`
- `query-memory-devices`
- `query-memory-size-summary`
- `query-machines`
- `query-version`

Explicitly forbidden:

- `device_add`, `device_del`, `qom-set`, `set_link`
- `system_powerdown`, `stop`, `cont`, `quit`
- `migrate`, `migrate_cancel`
- `object-add`, `object-del`
- `cpu-add`, `balloon`, `memsave`, `pmemsave`
