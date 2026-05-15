# Windows FluidVirt Node Actuator Lab-to-Product Runbook v1

## Goal

Move from temporary privileged lab pod mutations to controlled node-local `KARLNodeFluidActuator`.

## Required permissions

- Node-local runtime visibility for VM -> pod -> cgroup -> QEMU PID mapping.
- Write access constrained to allowlisted cgroup knobs: `cpu.max` and optional `cpu.weight`.
- Audit capability for before/after/rollback/return evidence.

## Allowed paths

- Target VM cgroup path only, validated against resolved mapping.
- No parent cgroup writes unless explicitly policy-authorized.
- No symlink traversal.

## Migration steps

1. Preserve previous full rerun evidence as baseline.
2. Deploy node-local actuator with kill switch on by default.
3. Register shell allowlist and target path allowlist.
4. Enable compliance engine checks for Windows readiness phases.
5. Run fixture matrix and dry-run safety checks.
6. Validate rollback and return-to-floor behavior.

## Tests to rerun

- `go test ./...`
- `python3 scripts/validate_json.py`
- `git diff --check`
- Compliance and lease fixture matrix under `examples/windows-fluid-compliance-fixtures`.

## Explicitly forbidden

- vCPU hotplug/unplug.
- VM spec patch as entitlement mechanism.
- LiveMigration/VMIM as scaling mechanism.
- Pool replica scaling as mechanism.
