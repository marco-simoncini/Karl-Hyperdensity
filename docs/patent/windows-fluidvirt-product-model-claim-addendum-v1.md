# Windows FluidVirt Product Model Claim Addendum v1

Draft tecnico per revisione brevettuale.

## Technical focus

- compliance engine classifies standalone and pool-child Windows shells as Hyperdensity Ready
- prearmed envelope runtime model for in-place lease transitions
- CPU entitlement via node-local actuator (`cgroup v2 cpu.max`)
- RAM entitlement via QMP balloon
- lease lifecycle modeled as non-disruptive in-place transitions

## Continuity constraints

- same QEMU process identity
- same Windows boot identity
- no reboot, no recreate, no rollout
- no LiveMigration/VMIM mechanism

## Safety and reversibility

- mandatory rollback target
- mandatory return-to-floor target
- guest ACK witness requirement
- strict actuator allowlist and request identity pinning

## Audit and evidentiary chain

- deterministic replay hash model
- deterministic bundle append hash chain
- chain validation and duplicate/broken-chain rejection
