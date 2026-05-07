# Hyperdensity Contract Model Summary

Implemented package: `pkg/windowsfluidvirt`

## Contracts

- `WindowsFluidShell`
  - enforces `runtimeMode=in-place-qmp`
  - enforces no migration/reboot/recreate
  - tracks floor/envelope/runtimeTarget/runtimeActual
  - guest contract requires `agentModule=fluidShell` and `requireAck=true`
- `FluidResourceLease`
  - in-place lease spec (`mode=in-place`, ttl, rollback target)
  - hard guarantees (same node/pod/qemu/machine/boot + no migration/reboot/recreate)
  - readiness gates in status (`qmpAck`, `guestAck`, continuity, rollback, returnToFloor)
- `WindowsFluidEvidence`
  - before/after continuity points for pid/pod/node/boot/machine/vmi
  - qmp + guest evidence buckets
  - noReboot/noRecreate/noMigration proof booleans

## IDs

Added to `pkg/contracts/contract_ids.go`:

- `hyperdensity_windows_fluid_shell_v1`
- `hyperdensity_fluid_resource_lease_v1`
- `hyperdensity_windows_fluid_evidence_v1`
- `hyperdensity_windows_fluid_blocker_v1`
