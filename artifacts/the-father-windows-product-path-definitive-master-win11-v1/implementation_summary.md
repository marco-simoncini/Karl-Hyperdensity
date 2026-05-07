# Implementation Summary

Executed a controlled live product-path validation on `master-win11` using:

- KARL Windows Prearmed Fluid Envelope v2
- MVP `karl-node-fluid-actuator` for CPU entitlement liquidity via `cgroup v2 cpu.max`
- QMP balloon liquidity for RAM

Outcome in this run:

- compliance replay before apply: `HYPERDENSITY_READY_WINDOWS_SHELL`
- actuator dry-run: accepted
- CPU entitlement up/down via product actuator: confirmed
- RAM entitlement up/down via QMP balloon: confirmed
- rollback and return-to-floor: confirmed
- runtime continuity invariants: preserved (same QEMU, boot, pod, node, no VMIM/migration/rollout/recreate)
