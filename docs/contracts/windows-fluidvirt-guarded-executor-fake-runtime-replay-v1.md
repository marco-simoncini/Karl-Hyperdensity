# Windows FluidVirt Guarded Executor Fake-Runtime Replay v1

## Scope

This contract defines a deterministic fake-runtime replay surface for the Windows FluidVirt guarded executor.
It is governance-only and replay-only.

## Included

- deterministic replay of executor boundary checks
- replayed input/output boundary validations
- replayed approval/lease/witness/rollback/return-to-floor/kill-switch gate evaluations
- replayed blocking reasons and replayed audit events
- fake-runtime/no-mutation validation result

## Explicitly Excluded

- no executor runtime process
- no executor runtime loop
- no controlled apply execution path
- no controlled apply enablement
- no real cgroup writes
- no QMP command execution
- no QGA command execution
- no host runtime touch
- no Dashboard/Inventory/OS-ISO changes

## Safety Claims

- `executorEnabled=false`
- `executorRuntimeAvailable=false`
- `controlledApplyEnabled=false`
- `controlledApplyReady=false`
- `runtimeMutationEnabled=false`
- `cgroupWriteEnabled=false`
- `qmpCommandExecutionAllowed=false`
- `qgaCommandExecutionAllowed=false`
- `autonomousApplyAllowed=false`
- `productionMutationAllowed=false`
- `windowsGaClaimAllowed=false`
- `windowsProductionReadyClaimAllowed=false`
- `windowsExecutionReadyByDefault=false`

`executorFakeRuntimeReplayExecuted=true` is valid only for deterministic fake-runtime replay semantics and does not imply runtime readiness.

## Boundary Semantics

The fake runtime boundary must enforce:

- temporary-file-only fixture loading
- rejection of `/sys/fs/cgroup` paths
- rejection of raw QMP/QGA material
- rejection of secret/token material
- no runtime privileges required

## Forward Path

Runtime MVP remains a separate milestone.
Inventory fluidShell witness integration is expected before any runtime-readiness promotion.
