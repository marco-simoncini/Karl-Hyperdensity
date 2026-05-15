# windows-fluid-runtime-gate-v1

Evidence-only runtime gate for Windows Fluid Shell certification readiness.

## Scope

- Read-only evaluation only.
- No CPU/RAM runtime apply.
- No QMP mutating commands.
- No live migration/recreate/reboot rollout execution.

## Required VM annotations

- `hyperdensity.karl.io/fluid-runtime: "true"`
- `hyperdensity.karl.io/no-live-migration: "true"`
- `hyperdensity.karl.io/no-reboot: "required"`
- `hyperdensity.karl.io/no-recreate: "required"`
- `hyperdensity.karl.io/runtime-mode: "in-place-qmp"`
- `hyperdensity.karl.io/single-node-compatible: "true"`

Missing annotation gates produce `BLOCKED`.

## Runtime conditions produced

- `fluidRuntimeReady`
- `qmpReady`
- `guestAckReady`
- `noMigrationRequired`
- `noRebootProof`
- `sameQemuProcess`
- `sameNode`
- `sameVirtLauncherPod`
- `returnToFloorReady`
- `rollbackReady`

`READY` is reachable only when all required conditions are true.

## Certification classifications

- `SUPPORTED_CANDIDATE`
- `READY_FOR_FLUID_SHELL_CERTIFICATION`
- `BLOCKED_GENERIC_WINDOWS_VM`
- `BLOCKED_POOL_REPLICA_MODEL`
- `BLOCKED_MISSING_QMP`
- `BLOCKED_MISSING_GUEST_ACK`
- `BLOCKED_LIVE_MIGRATION_REQUIRED`
- `QUARANTINED_IDENTITY_CHANGED`
