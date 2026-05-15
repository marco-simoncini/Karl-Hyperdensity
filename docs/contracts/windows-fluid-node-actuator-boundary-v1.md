# Windows Fluid Node Actuator Boundary v1

This contract defines the **boundary-only** shape of the Windows Node Fluid Actuator in KARL Hyperdensity.

## Release Position

- release track: `technical_preview`
- lane status: `gated_preview`
- contract mode: `boundary_only`
- technical preview candidate only

## What This Milestone Does

- defines the contract boundary model (`WindowsFluidVirtNodeActuatorContract`)
- records required gates and blockers
- records safety defaults and forbidden actions
- records future MVP porting plan inputs

## What This Milestone Does Not Do

- does not port runtime actuator executable
- does not enable cgroup writes
- does not execute QMP/QGA commands
- does not enable autonomous apply
- does not enable production apply
- does not declare Windows GA
- does not declare Windows production-ready
- does not declare Windows execution-ready by default

## CPU Boundary

- mechanism: `cgroup_v2_cpu_max_entitlement_liquidity`
- mechanism state: `contract_defined_not_runtime_enabled`
- host scope: `node_local`
- write path allowed: `false`
- cgroup write enabled: `false`
- cpu.max mutation enabled: `false`
- requires allowlist, manual approval, guest witness, same-boot/same-QEMU proof
- requires rollback, return-to-floor, audit hash chain, kill switch, and lease TTL

## RAM Boundary

- mechanism: `qmp_balloon_liquidity_model`
- mechanism state: `contract_reference_only`
- QMP command execution allowed: `false`
- raw QMP control exposed: `false`
- memory apply enabled: `false`
- generic KubeVirt VM RAM template mutation allowed: `false`

## Separation Of Milestones

- Node Fluid Actuator MVP runtime porting: separate milestone
- OS-ISO packaging/systemd/DaemonSet: separate milestone
- Inventory fluidShell integration: separate milestone

## Safety Statement

All runtime mutation paths remain disabled in this boundary contract.
