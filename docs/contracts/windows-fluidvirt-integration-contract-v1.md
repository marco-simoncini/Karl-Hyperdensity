# windows-fluidvirt-integration-contract-v1

Defines how `Karl-Inventory/modules.fluidShell` evidence is consumed by `Karl-Hyperdensity`.

## Mapping

- `modules.fluidShell.agentModule` -> `WindowsFluidEvidence.guestEvidence.agentModule`
- `modules.fluidShell.agentModuleVersion` -> `WindowsFluidEvidence.guestEvidence.agentModuleVersion`
- `modules.fluidShell.guestAck` -> `FluidResourceLease.status.guestAck`
- `modules.fluidShell.processorCount` -> `WindowsFluidEvidence.guestEvidence.processorCount`
- `modules.fluidShell.visibleMemoryBytes` -> `WindowsFluidEvidence.guestEvidence.visibleMemoryBytes`
- `modules.fluidShell.lastBootTime` -> `WindowsFluidEvidence.lastBootAfter`
- `modules.fluidShell.machineGuidHash` -> `WindowsFluidEvidence.machineGuidAfter`
- `modules.fluidShell.pendingReboot` -> blocker projection `pending_reboot_detected`
- `modules.fluidShell.driverTruth.memoryAdapterVerified=false` -> blocker `memory_driver_unverified`
- `modules.fluidShell.criticalEvents.detected=true` -> blocker `critical_windows_event_detected`
- `modules.fluidShell.returnToFloor.ready=false` -> blockers `return_to_floor_not_ready`, `memory_return_not_safe`

## Guest Blocker Projection

- `pending_reboot_detected` -> BLOCKED
- `guest_ack_missing` -> BLOCKED
- `cpu_topology_not_confirmed` -> BLOCKED
- `guest_memory_not_confirmed` -> BLOCKED
- `memory_driver_unverified` -> BLOCKED
- `critical_windows_event_detected` -> BLOCKED
- `last_boot_changed` -> QUARANTINED when continuity comparison fails
- `machine_guid_changed` -> QUARANTINED when continuity comparison fails

## Mandatory Fields for READY

- `agentModule=fluidShell`
- `agentModuleVersion`
- `guestAck=true`
- `processorCount>0`
- `visibleMemoryBytes>0`
- `lastBootTime` present
- `machineGuidHash` present
- `pendingReboot=false`
- no critical blockers

## Mandatory Fields for Future ACTIVE

- all READY fields, plus:
- `qmpAck=true`
- `lastBootBefore == lastBootAfter`
- `qemuPidBefore == qemuPidAfter`
- `nodeBefore == nodeAfter`
- `virtLauncherPodBefore == virtLauncherPodAfter`
- `returnToFloor.ready=true`
- `rollbackResult.ready=true`

## Optional Fields

- `commitChargeBytes`
- `workingSetSummary`
- `pagePressure`
- `driverTruth.virtioDriversObserved`
- `driverTruth.balloonOrMemoryDriverObserved`

## Redaction / Privacy

- raw Machine GUID must stay guest-side in logs.
- exported contract uses `machineGuidHash` (sha256 lowercase hex) for continuity proof.
- no secrets/tokens are emitted in evidence payload.
