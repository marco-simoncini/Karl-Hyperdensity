# Inventory Guest Evidence Model

Source: `KarlInventoryAgent.FluidShell/FluidShellContracts.cs`

Fields implemented:

- `agentModule`, `agentModuleVersion`
- `guestAck`
- `hostname`, `osName`, `osVersion`, `architecture`
- `processorCount`
- `visibleMemoryBytes`, `freePhysicalMemoryBytes`
- `commitChargeBytes` (best effort)
- `workingSetSummary` (optional)
- `pagePressure` (optional)
- `lastBootTime`
- `machineGuidHash` (redacted/hashed continuity key)
- `pendingReboot`, `pendingRebootReasons`
- `driverTruth` (`qgaPresent`, `virtioDriversObserved`, `balloonOrMemoryDriverObserved`, `memoryAdapterVerified`)
- `criticalEvents` (`detected`, `summary`)
- `returnToFloor` (`ready`, `reason`, `conservativeBlockers`)
- `timestamps`
- `blockers`

Conservative rules:

- pending reboot => BLOCKED
- missing machine guid/last boot/cpu/memory truth => BLOCKED
- `memoryAdapterVerified=false` emits `memory_driver_unverified`
- `returnToFloor.ready` forced false in this phase
